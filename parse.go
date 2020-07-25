package annotation_service

import (
	"errors"
	"go/ast"
	"strings"
)

func isPackageIdent(t ast.Expr) bool {
	id, isIdent := t.(*ast.Ident)
	return isIdent && id.Name[0] <= 'Z' && id.Name[0] >= 'A'
}

func (p AnnotateParser) parseMethod(method *ast.Field) (err error) {
	// check doc
	if method.Doc == nil {
		return
	}
	var params, results []string
	if ft, ok := method.Type.(*ast.FuncType); ok {
		// get param results
		if ft.Params != nil {
			collectList(p.fileData, &params, ft.Params.List, p.file)
		}
		if ft.Results != nil {
			collectList(p.fileData, &results, ft.Results.List, p.file)
		}
	} else {
		return
	}
	var title string
	for i, cm := range method.Doc.List {
		// match annotation like
		// @namespace:type.method(opt1=xxx,opt2=xxx,opt3=xxx)
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
		// parse options
		err := parseKV(match[2], newApi.Options)
		if err != nil {
			return err
		}
		apiType := match[1]
		// check namespace
		if strings.ContainsRune(apiType, ':') {
			sp := strings.Split(apiType, ":")
			if sp[0] != p.namespace {
				continue
			}
			apiType = sp[1]
		}
		if strings.ContainsRune(apiType, '.') {
			sp := strings.SplitN(apiType, ".", 2)
			apiType = sp[0]
			newApi.Method = sp[1]
		}
		// split method
		if p.m[apiType] == nil {
			p.m[apiType] = &ApiAnnotate{
				Interface: p.file.Name.String() + "." + p.serviceName,
			}
		}
		p.m[apiType].Apis = append(p.m[apiType].Apis, newApi)
	}
	return
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
