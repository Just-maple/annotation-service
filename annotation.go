package annotation_service

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	apiAnnotateRegex             = regexp.MustCompile(`@(!?[A-Za-z0-9_.:]+?)\((.+?)\)`)
	serviceAnnotateRegexTemplate = `@%s\((.+?)\)`

	logger = log.New(os.Stdout, "[SVC] ", 0)
)

type (
	ApiAnnotateItem struct {
		Handler string
		Params  []string
		Returns []string
		Title   string
		Method  string
		Args    []string
		Options map[string]string
		Doc     []string
	}

	Service struct {
		InterfaceName string
		ServiceName   string
		Pkg           string
		OtherOptions  map[string]string
		ApiAnnotates  map[string]*ApiAnnotate
	}

	ApiAnnotate struct {
		Interface string
		Apis      []ApiAnnotateItem
	}

	AnnotateParser struct {
		m           map[string]*ApiAnnotate
		namespace   string
		serviceName string
		fileData    []byte
		file        *ast.File
	}

	Option func(*opt)

	opt struct {
		Ident string
	}

	Options []Option
)

func WithIdent(ident string) Option {
	return func(o *opt) {
		o.Ident = ident
	}
}

func (os Options) Apply() *opt {
	op := &opt{
		Ident: "service",
	}
	for _, o := range os {
		o(op)
	}
	return op
}

func (r Service) GetPath(apiDir string) (groupRoute, dir, pkg string) {
	group := r.OtherOptions["group"]
	if len(group) == 0 {
		group = r.ServiceName
	}
	group = strings.Trim(group, `"`)
	// 组路由前缀
	groupRoute = r.OtherOptions["route"]
	if len(groupRoute) == 0 {
		groupRoute = fmt.Sprintf(`"%s"`, group)
	}
	// 创建组目录
	dir = filepath.Join(apiDir, group)
	// package名
	pkg = filepath.Base(dir)
	return
}

func GetAllService(file string, opts ...Option) (res []Service, err error) {
	fileData, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	f, err := parser.ParseFile(token.NewFileSet(), "", fileData, parser.ParseComments)
	if err != nil {
		return
	}

	o := Options(opts).Apply()
	serviceAnnotateRegex := regexp.MustCompile(fmt.Sprintf(serviceAnnotateRegexTemplate, o.Ident))

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
					ApiAnnotates:  apis,
					Pkg:           f.Name.Name,
				}
				if len(annotate) > 1 {
					_, svc.OtherOptions, err = parseKV(strings.Join(annotate[1:], ","))
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

func AnalysisServiceWithFileToken(fileData []byte, serviceName, namespace string) (apiAnnotate map[string]*ApiAnnotate, err error) {
	f, err := parser.ParseFile(token.NewFileSet(), "", fileData, parser.ParseComments)
	if err != nil {
		return
	}
	aParser := AnnotateParser{
		m:           make(map[string]*ApiAnnotate),
		namespace:   namespace,
		serviceName: serviceName,
		fileData:    fileData,
		file:        f,
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
			if sp.Name.Name != serviceName {
				continue
			}
			itf, ok := sp.Type.(*ast.InterfaceType)
			if !ok {
				continue
			}
			for _, method := range itf.Methods.List {
				err = aParser.parseMethod(method)
				if err != nil {
					return
				}
			}
		}
	}
	apiAnnotate = aParser.m
	logger.Printf("analysis finished: %s", serviceName)
	return
}

func AnalysisFileService(file string, serviceName, namespace string) (apiAnnotate map[string]*ApiAnnotate, err error) {
	fileData, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	return AnalysisServiceWithFileToken(fileData, serviceName, namespace)
}
