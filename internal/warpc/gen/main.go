//go:generate go run main.go
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
)

var scripts = []string{
	"greet.js",
	"renderkatex.js",
}

func main() {
	for _, script := range scripts {
		filename := filepath.Join("../js", script)
		err := buildJSBundle(filename)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func buildJSBundle(filename string) error {
	minify := true
	result := api.Build(
		api.BuildOptions{
			EntryPoints:       []string{filename},
			Bundle:            true,
			MinifyWhitespace:  minify,
			MinifyIdentifiers: minify,
			MinifySyntax:      minify,
			Target:            api.ES2020,
			Outfile:           strings.Replace(filename, ".js", ".bundle.js", 1),
			SourceRoot:        "../js",
		})

	if len(result.Errors) > 0 {
		return fmt.Errorf("build failed: %v", result.Errors)
	}
	if len(result.OutputFiles) != 1 {
		return fmt.Errorf("expected 1 output file, got %d", len(result.OutputFiles))
	}

	of := result.OutputFiles[0]
	if err := os.WriteFile(filepath.FromSlash(of.Path), of.Contents, 0o644); err != nil {
		return fmt.Errorf("write file failed: %v", err)
	}
	return nil
}
