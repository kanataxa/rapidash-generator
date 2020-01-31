package main

//go:generate statik -src=../../template -dest=../../
import (
	"log"
	"os"

	"github.com/jessevdk/go-flags"
	generator "github.com/kanataxa/rapidash-generator"
)

type Options struct {
	DependenceFiles []string `short:"d" long:"deps" description:"dependency files"`
	ShouldOverwrite bool     `short:"w" description:"force write if file is already exists"`
	Output          string   `short:"o" long:"output" description:"output file name. default: os.Stdout"`
	Tag             string   `short:"t" long:"tag" default:"db" description:"use tag name"`
}

var opt Options

func run(path string) error {
	if err := generator.Generate(path, &generator.Config{
		DependenceFiles: opt.DependenceFiles,
		ShouldOverwrite: opt.ShouldOverwrite,
		Output:          opt.Output,
		Tag:             opt.Tag,
	}); err != nil {
		return err
	}
	return nil
}

func main() {
	args, err := flags.Parse(&opt)
	if err != nil {
		e, ok := err.(*flags.Error)
		if ok && e.Type == flags.ErrHelp {
			os.Exit(0)
		}
		log.Fatal(err)
	}
	if len(args) == 0 {
		log.Fatal("cannot get input file path, please pass path arguments")
	}
	if err := run(args[0]); err != nil {
		log.Fatal(err)
	}
}
