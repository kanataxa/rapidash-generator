package main

import (
	"log"
	"os"

	"github.com/jessevdk/go-flags"
	generator "github.com/kanataxa/rapidash-generator"
)

type Options struct {
	ShouldOverwrite bool   `short:"w" description:"force write if file is already exists"`
	Input           string `short:"f" long:"file" description:"input go file name"`
	Output          string `short:"o" long:"output" description:"output file name. default: os.Stdout"`
	Tag             string `short:"t" long:"tag" default:"db" description:"use tag name"`
}

var opt Options

func run() error {
	if err := generator.Generate(&generator.Config{
		ShouldOverwrite: opt.ShouldOverwrite,
		FilePath:        opt.Input,
		Output:          opt.Output,
		Tag:             opt.Tag,
	}); err != nil {
		return err
	}
	return nil
}

func main() {
	if _, err := flags.Parse(&opt); err != nil {
		e, ok := err.(*flags.Error)
		if ok && e.Type == flags.ErrHelp {
			os.Exit(0)
		}
		log.Fatal(err)
	}
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
