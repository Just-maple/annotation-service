package http_gin

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	annotation_service "github.com/Just-maple/annotation-service"
	"github.com/iancoleman/strcase"
	"golang.org/x/tools/imports"
)

type svcSet struct {
	svcFs        *token.FileSet
	svcMountF    *ast.File
	svcData      []byte
	serviceMount string
	svcUpdate    bool
}

func GetGoModFilePath() (modPath string) {
	cmd := exec.Command("go", "env", "GOMOD")
	stdout := &bytes.Buffer{}
	cmd.Stdout = stdout
	_ = cmd.Run()
	mod := stdout.String()
	mod = strings.Trim(mod, "\n")
	return mod
}

const svcTemplate = `package %s
		// @autowire(set=Server)
		type Service struct {
		}`

func Service2GinApi(apiDir, serviceFilePath, implDir, serviceMountFilePath, testerFunc string) {
	modDir := filepath.Dir(GetGoModFilePath())
	apiDir = filepath.Join(modDir, apiDir)
	implDir = filepath.Join(modDir, implDir)
	serviceMountFilePath = filepath.Join(modDir, serviceMountFilePath)

	fileData, err := ioutil.ReadFile(serviceFilePath)
	if err != nil {
		panic(err)
	}
	f, err := parser.ParseFile(token.NewFileSet(), "", fileData, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	res, err := annotation_service.GetAllService(serviceFilePath)
	if err != nil {
		panic(err)
	}
	sSet := &svcSet{svcFs: token.NewFileSet(), serviceMount: serviceMountFilePath}
	if len(serviceMountFilePath) > 0 {
		emptySvc := fmt.Sprintf(svcTemplate, filepath.Base(filepath.Dir(serviceMountFilePath)))
		sSet.svcData, err = ioutil.ReadFile(serviceMountFilePath)
		if err != nil {
			sSet.svcData = []byte(emptySvc)
		}

		sSet.svcMountF, err = parser.ParseFile(sSet.svcFs, "", sSet.svcData, parser.ParseComments)
		if err != nil {
			panic(err)
		}
		defer func() {
			if sSet.svcUpdate && err == nil {
				sSet.write()
			}
		}()
	}

	// 所有接口定义
	for _, r := range res {
		interfaceName := f.Name.Name + "." + r.InterfaceName
		sSet.UpdateSvcSet(interfaceName)
		if len(implDir) > 0 {
			implDir := filepath.Join(implDir, "svc_"+strcase.ToSnake(r.InterfaceName))
			implFileName := filepath.Join(implDir, strcase.ToSnake(r.InterfaceName)+".go")
			_ = os.MkdirAll(implDir, 0775)
			// 生成impl文件
			if _, err := ioutil.ReadFile(implFileName); err != nil {
				var bf bytes.Buffer
				err = implTemplate.Execute(&bf, &Impl{
					Service:   strcase.ToSnake(r.InterfaceName),
					Interface: interfaceName,
				})
				if err != nil {
					panic(err)
				}
				res, err := imports.Process("", bf.Bytes(), nil)
				if err != nil {
					panic(err)
				}
				err = ioutil.WriteFile(implFileName, res, 0664)
				if err != nil {
					panic(err)
				}
			}
			if t := f.Scope.Lookup(r.InterfaceName); t != nil {
				itfT := t.Decl.(*ast.TypeSpec).Type.(*ast.InterfaceType)
				fs := token.NewFileSet()
				s, _ := ioutil.ReadFile(implFileName)
				f2, _ := parser.ParseFile(fs, "", s, parser.ParseComments)
				err = ioutil.WriteFile(implFileName, SyncItf2Svc("Service", f.Name.Name, fs, f2, itfT), 0664)
				if err != nil {
					panic(err)
				}
			}
		}

		// 接口组
		group := r.OtherOptions["group"]
		if len(group) == 0 {
			group = r.ServiceName
		}
		group = strings.Trim(group, `"`)

		// 组路由前缀
		groupRoute := r.OtherOptions["route"]
		if len(groupRoute) == 0 {
			groupRoute = fmt.Sprintf(`"%s"`, r.ServiceName)
		}
		// 创建组目录
		dir := filepath.Join(apiDir, group)

		// package名
		pkg := filepath.Base(group)
		// 路由组
		route := GinRoute{
			Package:     pkg,
			GroupName:   strcase.ToCamel(r.ServiceName),
			ServiceName: r.InterfaceName,
			GroupRoute:  groupRoute,
		}
		// 测试用例组
		testCases := ApiTestTemplate{
			Package:       pkg,
			TesterFactory: testerFunc,
		}
		//  http服务
		httpApis, ok := r.ApiAnnotates["http"]
		if !ok {
			continue
		}
		route.ServiceName = httpApis.Interface
		// 所有http注解
		for _, api := range httpApis.Apis {
			// 命名空间过滤
			namespace := api.Options["ns"]
			if len(namespace) > 0 && namespace != r.ServiceName {
				continue
			}
			// 路由名
			routePath := api.Options["route"]
			route.Apis = append(route.Apis, GinApi{
				HttpMethod: strings.ToUpper(api.Options["method"]),
				Route:      routePath,
				Handler:    api.Handler,
				Title:      api.Title,
			})
			testFuncName := strings.Title(api.Options["method"]) + strings.Title(strcase.ToCamel(strings.Join(strings.Split(routePath, "/"), "_")))
			testCase := ApiTestCase{
				Router:   fmt.Sprintf(`"%s"`, path.Join(strings.Trim(groupRoute, `"`), strings.Trim(routePath, `"`))),
				Method:   strings.Title(api.Options["method"]),
				TestName: testFuncName,
				TestFunc: "tester.Test",
				Title:    api.Title,
			}
			testCase.Jumper = httpApis.Interface + "." + api.Handler
			if len(api.Params) > 1 {
				testCase.Param = api.Params[1]
			}
			if len(api.Returns) > 1 {
				testCase.Return = api.Returns[0]
			}
			testCases.Tests = append(testCases.Tests, testCase)
		}
		if len(route.Apis) == 0 {
			continue
		}

		_ = os.MkdirAll(dir, 0775)
		if err != nil {
			panic(err)
		}

		filename := r.OtherOptions["filename"]
		if len(filename) == 0 {
			filename = r.ServiceName
		}
		if err = WriteApiFiles(filepath.Join(dir, filename+".go"), &route); err != nil {
			panic(err)
		}
		if err = WriteTestFile(filepath.Join(dir, filename+"_test.go"), testCases, false); err != nil {
			panic(err)
		}
	}
}

func WriteApiFiles(path string, route *GinRoute) (err error) {
	// 写入数据
	var bf bytes.Buffer
	err = ginSvcTemplate.Execute(&bf, &route)
	if err != nil {
		panic(err)
	}
	ret, err := imports.Process("", bf.Bytes(), nil)
	if err != nil {
		panic(err)
	}
	if err = ioutil.WriteFile(path, ret, 0664); err != nil {
		panic(err)
	}
	return
}

func (sSet *svcSet) write() {
	ret, err := imports.Process("", sSet.svcData, nil)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(sSet.serviceMount, ret, 0664)
	if err != nil {
		panic(err)
	}
}

func (sSet *svcSet) UpdateSvcSet(interfaceName string) {
	if sSet.svcMountF == nil {
		return
	}
	svc := sSet.svcMountF.Scope.Lookup("Service")
	if svc != nil {
		typ := svc.Decl.(*ast.TypeSpec)
		fields := typ.Type.(*ast.StructType)
		fields.Incomplete = true
		list := fields.Fields.List
		var has bool
		for _, f := range list {
			if interfaceName == string(sSet.svcData[f.Type.Pos()-1:f.Type.End()-1]) {
				has = true
			}
		}
		if !has {
			// 插入新的服务
			sp := strings.Split(interfaceName, ".")
			fields.Fields.List = append(fields.Fields.List, &ast.Field{
				Names: []*ast.Ident{{
					Name: sp[1],
				}},
				Type: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: sp[0],
					},
					Sel: &ast.Ident{
						Name: sp[1],
					},
				},
			})
			sSet.svcUpdate = true
		}
	}
	if sSet.svcUpdate {
		var bf bytes.Buffer
		var err error
		err = format.Node(&bf, sSet.svcFs, sSet.svcMountF)
		if err != nil {
			panic(err)
		}
		sSet.svcFs = token.NewFileSet()
		sSet.svcMountF, err = parser.ParseFile(sSet.svcFs, "", bf.Bytes(), parser.ParseComments)
		if err != nil {
			panic(err)
		}
		sSet.svcData = bf.Bytes()
	}
}
