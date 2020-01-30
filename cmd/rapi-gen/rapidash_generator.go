package main

import (
	"log"
	"path/filepath"

	generator "github.com/kanataxa/rapidash-generator"
)

func main() {
	gen, err := generator.Parse(filepath.Join("testdata", "entity.go"))
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	if err := gen.Generate(); err != nil {
		log.Fatalf("%v\n", err)
	}
}
