package generator

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/xerrors"
)

func parseDir(path string, fset *token.FileSet, resolveFiles []string) (map[string]*ast.Package, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, xerrors.Errorf("failed to open file: %w", err)
	}
	defer f.Close()
	in, err := f.Stat()
	if err != nil {
		return nil, xerrors.Errorf("failed to get stat: %w", err)
	}
	dir := path
	if !in.IsDir() {
		dir = filepath.Dir(path)
	}
	pkgs, err := parser.ParseDir(fset, dir, func(info os.FileInfo) bool {
		if strings.Contains(info.Name(), "_test.go") {
			return false
		}
		if in.IsDir() {
			return true
		}
		for _, f := range resolveFiles {
			if strings.Contains(f, info.Name()) {
				return true
			}
		}
		return strings.Contains(path, info.Name())
	}, 0)
	if err != nil {
		return nil, xerrors.Errorf("failed to parse dir: %w", err)
	}
	return pkgs, nil
}

func parsePkg(fset *token.FileSet, astPkg *ast.Package, tagField string) (FunctionGenerator, error) {
	var files []*ast.File
	for _, f := range astPkg.Files {
		files = append(files, f)
	}
	conf := types.Config{
		Importer: importer.Default(),
		Error: func(err error) {
			fmt.Printf("Types.Config run error: %+v\n", err)
		},
	}
	info := &types.Info{}
	pkg, err := conf.Check(astPkg.Name, fset, files, info)
	if err != nil {
		return nil, xerrors.Errorf("failed to check: %w", err)
	}
	scope := pkg.Scope()
	var structs []*Struct
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
			structs = append(structs, &Struct{
				obj:    obj,
				Fields: fields,
			})
		}
	}
	return &GoSourceGenerator{Structs: structs}, nil
}

func Parse(path string, config *Config) (FunctionGenerator, error) {
	fset := token.NewFileSet()
	pkgs, err := parseDir(path, fset, config.DependenceFiles)
	if err != nil {
		return nil, xerrors.Errorf("failed to parse dir: %w", err)
	}
	if len(pkgs) > 1 {
		return nil, xerrors.New("contain multiple package")
	}
	for _, astPkg := range pkgs {
		generator, err := parsePkg(fset, astPkg, config.Tag)
		if err != nil {
			return nil, xerrors.Errorf("failed to parse pkg: %w", err)
		}
		return generator, nil
	}
	return nil, nil
}
