package generator

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io/ioutil"
	"strings"

	"golang.org/x/xerrors"
)

func Parse(filepath string) (FunctionGenerator, error) {
	src, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, xerrors.Errorf("failed to read file: %w", err)
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, parser.Mode(0))
	if err != nil {
		return nil, xerrors.Errorf("failed to parse file: %w", err)
	}
	conf := types.Config{
		Importer: importer.Default(),
		Error: func(err error) {
			fmt.Printf("Types.Config run error: %+v\n", err)
		},
	}
	pkg, err := conf.Check("pkg", fset, []*ast.File{f}, nil)
	if err != nil {
		return nil, xerrors.Errorf("failed to check: %w", err)
	}
	scope := pkg.Scope()
	var strcuts []*Struct
	for _, name := range scope.Names() {
		obj := scope.Lookup(name)
		structType, ok := obj.Type().Underlying().(*types.Struct)
		if !ok {
			continue
		}
		var fields []*Field
		for i := 0; i < structType.NumFields(); i++ {
			tag := structType.Tag(i)
			if !strings.Contains(tag, "db:") {
				continue
			}
			fields = append(fields, &Field{
				v:      structType.Field(i),
				DBName: strings.Split(tag, "db:")[1],
			})
		}
		if len(fields) > 0 {
			strcuts = append(strcuts, &Struct{
				obj:    obj,
				Fields: fields,
			})
		}
	}
	return &GoSourceGenerator{
		Structs: strcuts,
	}, nil
}
