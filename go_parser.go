package generator

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/xerrors"
)

func Parse(path, tagField string) (FunctionGenerator, error) {
	var structs []*Struct
	root := path
	if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return xerrors.Errorf("failed to walk: %w", err)
		}
		if info.IsDir() && root != path {
			return filepath.SkipDir
		}
		if filepath.Ext(path) != ".go" {
			return nil
		}
		s, err := parse(path, tagField)
		if err != nil {
			return xerrors.Errorf("failed to parse %s: %w", path, err)
		}
		structs = append(structs, s...)

		return nil
	}); err != nil {
		return nil, xerrors.Errorf("failed to walk: %w", err)
	}
	return &GoSourceGenerator{Structs: structs}, nil
}

func parse(fpath, tagField string) ([]*Struct, error) {
	src, err := ioutil.ReadFile(fpath)
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
			if !strings.Contains(tag, fmt.Sprintf("%s:", tagField)) {
				continue
			}
			leftTrimTag := strings.Split(tag, fmt.Sprintf("%s:", tagField))[1]
			fields = append(fields, &Field{
				v:      structType.Field(i),
				DBName: fmt.Sprintf("\"%s\"", strings.Split(leftTrimTag, "\"")[1]),
			})
		}
		if len(fields) > 0 {
			strcuts = append(strcuts, &Struct{
				obj:    obj,
				Fields: fields,
			})
		}
	}
	return strcuts, nil
}
