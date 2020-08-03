package mount_field

import (
	"bytes"
	"errors"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
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
	MountTypeField(fieldType string, name string) error
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

func (sSet *structMounter) MountTypeField(fieldType, name string) (err error) {
	var (
		typ    = sSet.object.Decl.(*ast.TypeSpec)
		fields = typ.Type.(*ast.StructType)
		list   = fields.Fields.List
	)
	var has bool
	for _, f := range list {
		se, ok := f.Type.(*ast.SelectorExpr)
		if !ok {
			continue
		}
		pkg, ok := se.X.(*ast.Ident)
		if !ok {
			continue
		}
		if fieldType == pkg.Name+"."+se.Sel.Name {
			has = true
		}
	}
	if has {
		return
	}
	sSet.insertField(fieldType, name, fields.Fields)
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

func (sSet *structMounter) insertField(fieldType, name string, list *ast.FieldList) {
	sp := strings.Split(fieldType, ".")
	if len(name) == 0 {
		name = sp[1]
	}
	addStructField(sSet.fileSet, list, fieldType, name)
}

// todo:fix comment
func addStructField(fset *token.FileSet, fields *ast.FieldList, typ string, name string) {
	f := &ast.Field{
		Names: []*ast.Ident{{Name: name}},
		Type:  &ast.Ident{Name: typ},
	}
	fields.List = append(fields.List, f)
}
