package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"text/template"

	_ "github.com/kanataxa/rapidash-generator/statik"
	"github.com/rakyll/statik/fs"
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

func (g *GoSourceGenerator) Generate() ([]byte, error) {
	statikFS, err := fs.New()
	if err != nil {
		return nil, xerrors.Errorf("failed to init statik: %w", err)
	}
	r, err := statikFS.Open("/rapidash_function.tmpl")
	if err != nil {
		return nil, xerrors.Errorf("failed to open statik file: %w", err)
	}
	defer r.Close()
	contents, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, xerrors.Errorf("failed to read file: %w", err)
	}
	tmpl, err := template.New("rapidash_generator").Parse(string(contents))
	if err != nil {
		return nil, xerrors.Errorf("failed to parse template: %w", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, g); err != nil {
		return nil, xerrors.Errorf("failed to exec template: %w", err)
	}
	source, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, xerrors.Errorf("failed to run format: %w", err)
	}
	return source, nil
}

type FunctionGenerator interface {
	Generate() ([]byte, error)
}

func Generate(path string, config *Config) error {
	generator, err := Parse(path, config)
	if err != nil {
		return xerrors.Errorf("failed to parse go source: %w", err)
	}
	source, err := generator.Generate()
	if err != nil {
		return xerrors.Errorf("failed to generate go source: %w", err)
	}

	writer := os.Stdout
	if config.Output != "" {
		exists, err := existsFile(config.Output)
		if err != nil {
			return xerrors.Errorf("failed to check exists: %w", err)
		}
		if exists && !config.ShouldOverwrite {
			return xerrors.New(fmt.Sprintf("file: [%s] is already exists. please give -w option", config.Output))
		}

		writer.Close()
		writer, err = os.Create(config.Output)
		if err != nil {
			return xerrors.Errorf("failed to create file: %w", err)
		}
		defer writer.Close()
	}
	if _, err := writer.Write(source); err != nil {
		return xerrors.Errorf("failed to write source: %w", err)
	}
	return nil
}

func existsFile(fpath string) (bool, error) {
	f, err := os.Open(fpath)
	if err != nil {
		if !os.IsNotExist(err) {
			return false, xerrors.Errorf("failed to open file: %w", err)
		}
		return false, nil
	}
	defer f.Close()
	return true, nil
}
