//go:generate go run main.go
package main

import (
	_ "embed"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/evanw/esbuild/pkg/api"
)

//go:embed livereload-hugo-plugin.js
var livereloadHugoPluginJS string

func main() {
	// 4.0.2
	// To upgrade to a new version, change to the commit hash of the version you want to upgrade to
	// then run mage generate from the root.
	const liveReloadCommit = "d803a41804d2d71e0814c4e9e3233e78991024d9"
	liveReloadSourceURL := fmt.Sprintf("https://raw.githubusercontent.com/livereload/livereload-js/%s/dist/livereload.js", liveReloadCommit)

	func() {
		resp, err := http.Get(liveReloadSourceURL)
		must(err)
		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		must(err)

		// Write the unminified livereload.js file.
		err = os.WriteFile("../livereload.js", b, 0o644)
		must(err)

		// Bundle and minify with ESBuild.
		result := api.Build(api.BuildOptions{
			Stdin: &api.StdinOptions{
				Contents: string(b) + livereloadHugoPluginJS,
			},
			Outfile:           "../livereload.min.js",
			Bundle:            true,
			Target:            api.ES2015,
			Write:             true,
			MinifyWhitespace:  true,
			MinifyIdentifiers: true,
			MinifySyntax:      true,
		})

		if len(result.Errors) > 0 {
			log.Fatal(result.Errors)
		}
	}()
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
