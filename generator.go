package generator

import (
	"bytes"
	"go/format"
	"os"
	"path/filepath"
	"text/template"

	"golang.org/x/xerrors"
)

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
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, g); err != nil {
		return xerrors.Errorf("failed to exec template: %w", err)
	}
	source, err := format.Source(buf.Bytes())
	if err != nil {
		return xerrors.Errorf("failed to run format: %w", err)
	}
	os.Stdout.Write(source)
	return nil
}

type FunctionGenerator interface {
	Generate() error
}
