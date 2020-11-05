package mount_field

import (
	"bytes"
	"errors"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"strconv"
	"strings"

	"golang.org/x/tools/imports"
)

type structMounter struct {
	fileSet    *token.FileSet
	astFile    *ast.File
	data       []byte
	structPath string
	structName string
	object     *ast.Object
}

type StructMounter interface {
	Write() (err error)
	MountTypeField(fieldType string, name, pkg string) error
}

func NewStructMounter(structPath, structName string) (set StructMounter, err error) {
	b, err := ioutil.ReadFile(structPath)
	if err != nil {
		return
	}
	m := &structMounter{
		fileSet:    token.NewFileSet(),
		data:       b,
		structPath: structPath,
		structName: structName,
	}
	m.astFile, err = parser.ParseFile(m.fileSet, "", m.data, parser.ParseComments)
	if err != nil {
		return
	}
	m.object = m.astFile.Scope.Lookup(m.structName)
	if m.object == nil {
		err = errors.New("struct not found")
		return
	}
	set = m
	return
}

func (sSet *structMounter) Write() (err error) {
	ret, err := imports.Process("", sSet.data, nil)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(sSet.structPath, ret, 0664)
	return
}

func (sSet *structMounter) MountTypeField(fieldType, name, pkgPath string) (err error) {
	var (
		typ    = sSet.object.Decl.(*ast.TypeSpec)
		fields = typ.Type.(*ast.StructType)
		list   = fields.Fields.List
	)

	pkgPath = strings.Trim(pkgPath, `"`)

	var needImport *string

	// check imports
	if len(pkgPath) > 0 {
		sp := strings.Split(fieldType, ".")
		if len(sp) != 2 {
			err = errors.New("invalid type with import package path")
			return
		}
		if p, ok := sSet.getPackageName(pkgPath); ok {
			sp[0] = p
			fieldType = strings.Join(sp, ".")
		} else {
			needImport = &sp[0]
		}
	}

	// check duplicate field
	for _, f := range list {
		se, ok := f.Type.(*ast.SelectorExpr)
		if !ok {
			continue
		}
		fpkg, ok := se.X.(*ast.Ident)
		if !ok {
			continue
		}
		if fieldType == fpkg.Name+"."+se.Sel.Name {
			return
		}
	}

	sSet.insertField(fieldType, name, pkgPath, fields.Fields)

	// append new import
	if needImport != nil {
		importDecl, ok := sSet.astFile.Decls[0].(*ast.GenDecl)
		if ok {
			importDecl.Specs = append(importDecl.Specs, &ast.ImportSpec{
				Doc: nil,
				Name: &ast.Ident{
					Name: *needImport,
				},
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: strconv.Quote(pkgPath),
				},
			})
		}
	}
	err = sSet.freshSet()
	return
}

func (sSet *structMounter) freshSet() (err error) {
	bf := new(bytes.Buffer)
	if err = format.Node(bf, sSet.fileSet, sSet.astFile); err != nil {
		return
	}
	sSet.data = bf.Bytes()
	sSet.fileSet = token.NewFileSet()
	if sSet.astFile, err = parser.ParseFile(sSet.fileSet, "", sSet.data, parser.ParseComments); err != nil {
		return
	}
	sSet.object = sSet.astFile.Scope.Lookup(sSet.structName)
	return
}

func (sSet *structMounter) getPackageName(pkgPath string) (name string, ok bool) {
	if len(pkgPath) == 0 {
		return
	}
	for _, imp := range sSet.astFile.Imports {
		v := strings.Trim(imp.Path.Value, `"`)
		if v == pkgPath {
			if imp.Name != nil {
				return imp.Name.Name, true
			} else {
				sp := strings.Split(pkgPath, "/")
				return sp[len(sp)-1], true
			}
		}
	}
	return
}

func (sSet *structMounter) insertField(fieldType, name, pkg string, list *ast.FieldList) {
	sp := strings.Split(fieldType, ".")
	if len(name) == 0 {
		name = sp[1]
	}
	addStructField(sSet.fileSet, list, fieldType, name, pkg)
}

// todo:fix comment
func addStructField(fset *token.FileSet, fields *ast.FieldList, typ string, name, pkg string) {
	f := &ast.Field{
		Names: []*ast.Ident{{Name: name}},
		Type:  &ast.Ident{Name: typ},
	}
	fields.List = append(fields.List, f)
}
