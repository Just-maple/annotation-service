package swaggen

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"

	annotation_service "github.com/Just-maple/annotation-service"
	"golang.org/x/tools/imports"
)

const swagTemplate = `
// @Summary  {{ .Summary }}
// @Produce  {{ .ContentType }} 
{{ .Params }}
// @Success 200 {{ .RetType }}	
// @Failure {{ .DefErrorCode }} {{ .DefErrorRet }} 
// @Router {{ .Route }} [{{ .Method }}]
func _() { _ = {{ .Jumper }} }
`

type swagApi struct {
	Summary      string
	ContentType  string
	Params       string
	RetType      retType
	RetExample   string
	DefErrorCode int
	DefErrorRet  retType
	ErrExample   string
	Route        string
	Method       string
	FuncTypes    string
	Jumper       string
}

type retType struct {
	Type       string
	ObjectType string
}

func (r retType) String() string {
	return fmt.Sprintf("%s %s", r.Type, r.ObjectType)
}

var swagT = template.Must(template.New("").Parse(swagTemplate))

func newDefSwagApi() swagApi {
	return swagApi{
		Summary:      "",
		ContentType:  "json",
		Params:       "//",
		RetExample:   "",
		DefErrorCode: http.StatusBadRequest,
		ErrExample:   "",
		Route:        "/",
		Method:       "get",
		FuncTypes:    "",
	}
}

type Opt func(*opt)
type Opts []Opt
type opt struct {
	RetPack    string
	RetPackKey string
}

func (opts Opts) Apply() *opt {
	o := new(opt)
	for _, f := range opts {
		f(o)
	}
	return o
}

func WithRetPack(pack interface{}, key string) Opt {
	t := reflect.TypeOf(pack)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return func(i *opt) {
		i.RetPack = t.String()
		i.RetPackKey = key
	}
}

func GenDoc(apiDir, servicePath string, opts ...Opt) {
	o := Opts(opts).Apply()
	res, err := annotation_service.GetAllService(servicePath)
	if err != nil {
		panic(err)
	}
	// 所有接口定义
	for _, r := range res {
		baseRoute, groupDir, groupPkg := r.GetPath(apiDir)
		head := "package " + groupPkg + "\n"
		fileName := filepath.Join(groupDir, r.ServiceName+"_doc.go")

		httpApis, ok := r.ApiAnnotates["http"]
		if !ok {
			continue
		}
		var apis []string
		// 所有http注解
		for _, api := range httpApis.Apis {
			bf := new(bytes.Buffer)
			t := newDefSwagApi()
			t.Method = getMethod(api)
			t.RetType = o.getRet(api)
			t.DefErrorRet = o.getErrRet(api)
			t.Summary = api.Title
			t.Route = getRoute(api, baseRoute)
			t.FuncTypes = getFuncTypes(api)
			t.Jumper = r.Pkg + "." + r.InterfaceName + "." + api.Handler
			t.Params = getParams(api, t.Route)
			err := swagT.Execute(bf, &t)
			if err != nil {
				panic(err)
			}
			apis = append(apis, bf.String())
		}
		if len(apis) == 0 {
			continue
		}
		data := fmt.Sprintf("%s\n\n%s", head, strings.Join(apis, "\n"))
		ret, err := imports.Process("", []byte(data), nil)
		if err != nil {
			panic(err)
		}
		_ = os.MkdirAll(groupDir, 0775)
		_ = ioutil.WriteFile(fileName, ret, 0664)
	}
}

var typeMap = map[string]string{
	"int": "integer",
}

func getParams(api annotation_service.ApiAnnotateItem, route string) (params string) {
	var paramList []string
	var isJson bool
	defer func() {
		if len(paramList) == 0 {
			params = "//"
			return
		}
		for i := range paramList {
			paramList[i] = fmt.Sprintf("// @Param %s", paramList[i])
		}
		if isJson {
			paramList = append(paramList, "// @Accept json")
		}
		params = strings.Join(paramList, "\n")
	}()
	if len(api.Params) > 0 && api.Params[0] == "context.Context" {
		api.Params = api.Params[1:]
	}
	if len(api.Params) > 0 {
		if strings.EqualFold(api.Method, "get") {
			paramList = append(paramList, fmt.Sprintf(`query query %s true "%s"`, api.Params[0], api.Params[0]))
		} else {
			isJson = true
			paramList = append(paramList, fmt.Sprintf(`json body %s true "%s"`, api.Params[0], api.Params[0]))
		}
	}
	spRoutes := strings.Split(route, "/")
	for _, r := range spRoutes {
		if strings.HasPrefix(r, "{") && strings.HasSuffix(r, "}") {
			p := r[1 : len(r)-1]
			paramList = append(paramList, fmt.Sprintf(`%s path string true "%s"`, p, p))
		}
	}
	return
}
func (o opt) getRet(api annotation_service.ApiAnnotateItem) (ret retType) {
	ret.Type = "null"
	ret.ObjectType = "_"
	var retType string
	if len(api.Returns) > 1 {
		retType = api.Returns[0]
		if len(retType) > 0 {
			if retType[0] <= 'z' && retType[0] >= 'a' {
				newRetType, ok := typeMap[retType]
				if ok {
					retType = newRetType
				}
				ret.Type = retType
				ret.ObjectType = "_"
			} else {
				ret.Type = "{object}"
				ret.ObjectType = retType
			}
		}
	}
	if len(o.RetPack) > 0 {
		ret.Type = "{object}"
		if len(retType) == 0 {
			ret.ObjectType = o.RetPack
		} else {
			ret.ObjectType = fmt.Sprintf("%s{%s=%s}", o.RetPack, o.RetPackKey, retType)
		}
	}
	return
}

func (o opt) getErrRet(api annotation_service.ApiAnnotateItem) (ret retType) {
	if len(o.RetPack) > 0 {
		ret.Type = "{object}"
		ret.ObjectType = o.RetPack
	} else {
		ret.Type = "string"
		ret.ObjectType = "_"
	}
	return
}

func getMethod(api annotation_service.ApiAnnotateItem) (method string) {
	method = api.Method
	if len(method) == 0 {
		method = api.Options["method"]
	}
	if strings.ContainsRune(method, '.') {
		method = strings.Split(method, ".")[0]
	}
	return
}

func getFuncTypes(api annotation_service.ApiAnnotateItem) (types string) {
	return
}

func getRoute(api annotation_service.ApiAnnotateItem, baseRoute string) (route string) {
	baseRoute = strings.Trim(baseRoute, `"`)
	route = api.Options["route"]
	if len(route) == 0 && len(api.Args) > 0 {
		route = api.Args[0]
	}
	route = strings.Trim(route, `"`)
	route = path.Join(baseRoute, route)
	sp := strings.Split(route, "/")
	for i, r := range sp {
		if strings.HasPrefix(r, ":") {
			sp[i] = fmt.Sprintf("{%s}", r[1:])
		}
	}
	route = path.Join("/", path.Join(sp...))
	return
}
