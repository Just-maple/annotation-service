package annotation_service

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"regexp"
	"strings"
)

var (
	apiAnnotateRegex     = regexp.MustCompile(`@([A-Za-z0-9_.]+?)\((.+?)\)`)
	serviceAnnotateRegex = regexp.MustCompile(`@service\((.+?)\)`)
)

type (
	ApiAnnotateItem struct {
		Handler string
		Params  []string
		Returns []string
		Title   string
		Options map[string]string
	}

	Service struct {
		InterfaceName string
		ServiceName   string
		OtherOptions  map[string]string
		ApiAnnotates  map[string]*ApiAnnotate
	}

	ApiAnnotate struct {
		Interface string
		Apis      []ApiAnnotateItem
	}
)

func GetAllService(file string) (res []Service, err error) {
	fileData, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	f, err := parser.ParseFile(token.NewFileSet(), "", fileData, parser.ParseComments)
	if err != nil {
		return
	}
	for i := range f.Decls {
		decl, ok := f.Decls[i].(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spt := range decl.Specs {
			sp, ok := spt.(*ast.TypeSpec)
			if !ok {
				continue
			}
			_, ok = sp.Type.(*ast.InterfaceType)
			if !ok {
				continue
			}
			var doc *ast.CommentGroup
			if len(decl.Specs) == 1 {
				doc = decl.Doc
			} else {
				doc = sp.Doc
			}
			if doc == nil {
				continue
			}
			svcNameMap := map[string]struct{}{}
			for _, cm := range doc.List {
				match := serviceAnnotateRegex.FindStringSubmatch(cm.Text)
				if len(match) != 2 {
					continue
				}
				annotate := strings.Split(match[1], ",")
				serviceName := strings.TrimSpace(annotate[0])
				if _, dup := svcNameMap[serviceName]; dup {
					err = errors.New("duplicate service name")
					return nil, err
				}
				svcNameMap[serviceName] = struct{}{}
				apis, err := AnalysisServiceWithFileToken(fileData, sp.Name.String(), serviceName)
				if err != nil {
					return nil, err
				}
				svc := Service{
					InterfaceName: sp.Name.String(),
					ServiceName:   serviceName,
					OtherOptions:  map[string]string{},
					ApiAnnotates:  apis,
				}
				if len(annotate) > 1 {
					err = parseKV(strings.Join(annotate[1:], ","), svc.OtherOptions)
					if err != nil {
						return nil, err
					}
				}
				res = append(res, svc)
			}
		}
	}
	return
}

func isPackageIdent(t ast.Expr) bool {
	id, isIdent := t.(*ast.Ident)
	return isIdent && id.Name[0] <= 'Z' && id.Name[0] >= 'A'
}

func AnalysisServiceWithFileToken(fileData []byte, serviceName, namespace string) (apiAnnotate map[string]*ApiAnnotate, err error) {
	f, err := parser.ParseFile(token.NewFileSet(), "", fileData, parser.ParseComments)
	if err != nil {
		return
	}
	apiAnnotate = make(map[string]*ApiAnnotate)
	for i := range f.Decls {
		decl, ok := f.Decls[i].(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spt := range decl.Specs {
			sp, ok := spt.(*ast.TypeSpec)
			if !ok {
				continue
			}
			if sp.Name.Name != serviceName {
				continue
			}
			itf, ok := sp.Type.(*ast.InterfaceType)
			if !ok {
				continue
			}
			for _, method := range itf.Methods.List {
				if method.Doc == nil {
					continue
				}
				var params, results []string
				if ft, ok := method.Type.(*ast.FuncType); ok {
					collectList(fileData, &params, ft.Params.List, f)
					collectList(fileData, &results, ft.Results.List, f)
				}
				var title string
				for i, cm := range method.Doc.List {
					match := apiAnnotateRegex.FindStringSubmatch(cm.Text)
					if len(match) == 0 {
						if i == 0 {
							text := strings.TrimPrefix(cm.Text, "//")
							title = strings.TrimSpace(text)
						}
						continue
					}
					if len(match) != 3 {
						continue
					}
					newApi := ApiAnnotateItem{
						Options: make(map[string]string),
						Handler: method.Names[0].Name,
						Title:   title,
						Params:  params,
						Returns: results,
					}
					err = parseKV(match[2], newApi.Options)
					if err != nil {
						return
					}
					apiType := match[1]
					// check namespace
					if strings.ContainsRune(apiType, '.') {
						sp := strings.Split(apiType, ".")
						if sp[0] != namespace {
							continue
						}
						apiType = sp[1]
					}
					if apiAnnotate[apiType] == nil {
						apiAnnotate[apiType] = &ApiAnnotate{
							Interface: f.Name.String() + "." + serviceName,
						}
					}
					apiAnnotate[apiType].Apis = append(apiAnnotate[apiType].Apis, newApi)
				}
			}
		}
	}
	fmt.Printf("[SVC] %s Analysis Finished\n", serviceName)
	return
}

func AnalysisFileService(file string, serviceName, namespace string) (apiAnnotate map[string]*ApiAnnotate, err error) {
	fileData, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	return AnalysisServiceWithFileToken(fileData, serviceName, namespace)
}

func parseKV(raw string, resMap map[string]string) (err error) {
	options := strings.Split(raw, ",")
	if resMap == nil {
		resMap = make(map[string]string)
	}
	if len(options) == 0 {
		return
	}
	for _, o := range options {
		res := strings.Split(o, "=")
		if len(res) == 1 {
			res = append(res, "")
		}
		if len(res) != 2 {
			err = errors.New("invalid options key")
			return

		}
		if len(res[0]) == 0 {
			err = errors.New("invalid options key")
			return
		}
		if _, dup := resMap[res[0]]; dup {
			err = errors.New("duplicate options key")
			return
		}
		resMap[res[0]] = res[1]
	}
	return
}

func collectList(fileData []byte, collectList *[]string, fl []*ast.Field, f *ast.File) {
	for _, l := range fl {
		p := string(fileData[l.Type.Pos()-1 : l.Type.End()-1])
		if se, isSt := l.Type.(*ast.StarExpr); isSt {
			p = p[1:]
			l.Type = se.X
		}
		if isPackageIdent(l.Type) {
			p = f.Name.String() + "." + p
		}
		*collectList = append(*collectList, p)
	}
}
