package main

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func main() {
	src, err := ioutil.ReadFile(filepath.Join("..","testdata", "entity.go"))
	if err != nil {
		panic(err)
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, parser.Mode(0))
	if err != nil {
		panic(err)
	}
	conf := types.Config{
		Importer: importer.Default(),
		Error: func(err error) {
			fmt.Printf("!!! %#v\n", err)
		},
	}
	pkg, err := conf.Check("rapi-gen", fset, []*ast.File{f}, nil)
	if err != nil {
		panic(err)
	}
	scope := pkg.Scope()
	for _, name := range scope.Names() {
		obj := scope.Lookup(name)
		structType, ok := obj.Type().Underlying().(*types.Struct)
		if !ok {
			continue
		}
		for i := 0; i < structType.NumFields(); i++ {
			tag := structType.Tag(i)
			if !strings.Contains(tag, "db:") {
				continue
			}
			fmt.Println(strings.Split(tag, "db:")[1])
			_, ok := structType.Field(i).Type().Underlying().(*types.Basic)
			fmt.Println(ok)
		}
	}
	fmt.Println(pkg.Scope().Names())
}
