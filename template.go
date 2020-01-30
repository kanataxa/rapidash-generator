package generator

import (
	"fmt"
	"go/types"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"golang.org/x/xerrors"
)

type Field struct {
	v      *types.Var
	DBName string
}

func (f *Field) FieldFunc() string {
	return f.fieldFunc(f.v.Type())
}

func (f *Field) fieldFunc(t types.Type) string {
	switch v := t.Underlying().(type) {
	case *types.Struct:
		s := f.v.Type().String()
		if !strings.Contains(s, "time.Time") {
			break
		}
		return fmt.Sprint("FieldTime")
	case *types.Basic:
		return fmt.Sprintf("Field%s", strcase.ToCamel(v.String()))
	case *types.Pointer:
		return fmt.Sprintf("%sPtr", f.fieldFunc(v.Elem()))
	default:
		log.Printf("invalid type field: %v", v)
	}
	return fmt.Sprintf("Field%s", strcase.ToCamel(t.String()))
}

type Struct struct {
	obj    types.Object
	Fields []*Field
}

func (s *Struct) TableName() string {
	return strcase.ToSnake(fmt.Sprintf("%ss", s.obj.Name()))
}

func (s *Struct) Name() string {
	return s.obj.Name()
}

func (s *Struct) Package() string {
	return s.obj.Pkg().Name()
}

type GoSourceGenerator struct {
	Structs []*Struct
}

func (g *GoSourceGenerator) Package() string {
	if len(g.Structs) == 0 {
		return ""
	}
	return g.Structs[0].Package()
}

func (g *GoSourceGenerator) Generate() error {
	tmpl, err := template.ParseFiles(filepath.Join("template", "rapidash_function.tmpl"))
	if err != nil {
		return xerrors.Errorf("failed to parse template: %w", err)
	}
	if err := tmpl.Execute(os.Stdout, g); err != nil {
		return xerrors.Errorf("failed to exec template: %w", err)
	}
	return nil
}

type FunctionGenerator interface {
	Generate() error
}
