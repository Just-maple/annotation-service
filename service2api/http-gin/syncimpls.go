package gin

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"strings"

	"golang.org/x/tools/imports"
)

func SyncItf2Svc(recvTypeName, itfPkg string, fs *token.FileSet, file *ast.File, interfaceType *ast.InterfaceType) (res []byte) {
	f := SyncTypeMap(recvTypeName, itfPkg, GetInterfaceTypeMap(interfaceType), GetServiceTypeMap(recvTypeName, file.Decls))
	var bf bytes.Buffer
	_ = format.Node(&bf, fs, file)
	ret := bf.String()
	ret += "\n" + strings.Join(f, "\n\n")
	res, err := imports.Process("", []byte(ret), nil)
	if err != nil {
		panic(err)
	}
	return
}

func GetServiceTypeMap(serviceName string, decls []ast.Decl) (res map[string]*ast.FuncDecl) {
	res = map[string]*ast.FuncDecl{}
	for _, decl := range decls {
		f, ok := decl.(*ast.FuncDecl)
		if !ok || f.Recv == nil || len(f.Recv.List) != 1 || f.Recv.List[0].Type == nil {
			continue
		}
		rident, ok := f.Recv.List[0].Type.(*ast.Ident)
		if !ok || rident.Name != serviceName {
			continue
		}
		res[f.Name.Name] = f
	}
	return
}

const implFuncTemplate = `func (s %s) %s { 
panic("implement me") 
} `

func SyncTypeMap(recvName, itfPkg string, itfMap map[string]*ast.FuncType, funcMap map[string]*ast.FuncDecl) (appendStr []string) {
	type ap struct {
		name string
		f    *ast.FuncType
	}
	var toAppend []ap
	for k, m := range itfMap {
		f, ok := funcMap[k]
		if !ok {
			toAppend = append(toAppend, ap{
				name: k,
				f:    m,
			})
			continue
		}
		f.Type.Params = copyFieldList(itfPkg, m.Params)
		f.Type.Results = copyFieldList(itfPkg, m.Results)
	}

	for _, a := range toAppend {
		a.f.Params = copyFieldList(itfPkg, a.f.Params)
		a.f.Results = copyFieldList(itfPkg, a.f.Results)
		var bf bytes.Buffer
		_ = format.Node(&bf, token.NewFileSet(), a.f)
		if bf.Len() > 4 {
			appendStr = append(appendStr, fmt.Sprintf(implFuncTemplate, recvName, a.name+bf.String()[4:]))
		}
	}
	return
}

func copyFieldList(itfPkg string, src *ast.FieldList) (dst *ast.FieldList) {
	if src == nil {
		return nil
	}
	src.Opening = 0
	src.Closing = 0
	for i := range src.List {
		for j := range src.List[i].Names {
			n := src.List[i].Names[j]
			o := n.Obj
			n.NamePos = 0
			n.Obj = ast.NewObj(o.Kind, o.Name)
		}
		switch s := src.List[i].Type.(type) {
		case *ast.SelectorExpr:
			if s.Sel.Obj != nil {
				s.Sel.Obj = ast.NewObj(s.Sel.Obj.Kind, s.Sel.Obj.Name)
			}
			s.Sel.NamePos = 0
		case *ast.Ident:
			if s.Obj == nil {
				continue
			}
			if s.Obj.Kind == ast.Typ {
				s.NamePos = 0
				src.List[i].Type = &ast.SelectorExpr{
					X:   ast.NewIdent(itfPkg),
					Sel: s,
				}
			}
		}
	}
	return src
}

func GetInterfaceTypeMap(itf *ast.InterfaceType) (res map[string]*ast.FuncType) {
	res = map[string]*ast.FuncType{}
	for _, m := range itf.Methods.List {
		ft, ok := m.Type.(*ast.FuncType)
		if !ok {
			continue
		}
		res[m.Names[0].Name] = ft
	}
	return
}
