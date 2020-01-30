package generator

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"go/types"
	"log"
	"strings"
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
