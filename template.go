package generator

import (
	"fmt"
	"go/types"
	"log"
	"strings"

	"github.com/iancoleman/strcase"
)

type Field struct {
	v      *types.Var
	DBName string
}

func (f *Field) Name() string {
	return f.v.Name()
}

func (f *Field) FieldFunc() string {
	return f.function("Field", f.v.Type())
}

func (f *Field) EncodeFunc() string {
	return f.function("", f.v.Type())
}

func (f *Field) DecodeFunc() string {
	return f.function("", f.v.Type())
}

func (f *Field) IsWrapType() bool {
	unwrap := f.v.Type().Underlying()
	if _, ok := unwrap.(*types.Basic); !ok {
		return false
	}
	return f.v.Type().String() != unwrap.String()
}

func (f *Field) WrapType() string {
	if !f.IsWrapType() {
		return ""
	}
	return strings.Trim(f.v.Type().String(), fmt.Sprintf("%s.", f.v.Pkg().Path()))
}

func (f *Field) UnwrapType() string {
	if !f.IsWrapType() {
		return ""
	}
	return f.v.Type().Underlying().String()
}

func (f *Field) function(prefix string, t types.Type) string {
	switch v := t.Underlying().(type) {
	case *types.Struct:
		s := f.v.Type().String()
		if !strings.Contains(s, "time.Time") {
			break
		}
		return fmt.Sprintf("%sTime", prefix)
	case *types.Basic:
		return fmt.Sprintf("%s%s", prefix, strcase.ToCamel(v.String()))
	case *types.Pointer:
		return fmt.Sprintf("%sPtr", f.function(prefix, v.Elem()))
	default:
		log.Printf("invalid type field: %v", v)
	}
	return fmt.Sprintf("%s%s", prefix, strcase.ToCamel(t.String()))
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
