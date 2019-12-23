package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gohugoio/hugo/common/hugio"

	"github.com/spf13/afero"
)

func main() {
	// TODO(bep) git checkout tag
	// The current is built with Go version 9341fe073e6f7742c9d61982084874560dac2014 / go1.13.5
	fmt.Println("Forking ...")
	defer fmt.Println("Done ...")

	cleanFork()

	htmlRoot := filepath.Join(forkRoot, "htmltemplate")

	for _, pkg := range goPackages {
		copyGoPackage(pkg.dstPkg, pkg.srcPkg)
	}

	for _, pkg := range goPackages {
		doWithGoFiles(pkg.dstPkg, pkg.rewriter, pkg.replacer)
	}

	goimports(htmlRoot)
	gofmt(forkRoot)

}

const (
	// TODO(bep)
	goSource = "/Users/bep/dev/go/dump/go/src"
	forkRoot = "../../tpl/internal/go_templates"
)

type goPackage struct {
	srcPkg   string
	dstPkg   string
	replacer func(name, content string) string
	rewriter func(name string)
}

var (
	textTemplateReplacers = strings.NewReplacer(
		`"text/template/`, `"github.com/gohugoio/hugo/tpl/internal/go_templates/texttemplate/`,
		`"internal/fmtsort"`, `"github.com/gohugoio/hugo/tpl/internal/go_templates/fmtsort"`,
		// Rename types and function that we want to overload.
		"type state struct", "type stateOld struct",
		"func (s *state) evalFunction", "func (s *state) evalFunctionOld",
		"func (s *state) evalField(", "func (s *state) evalFieldOld(",
		"func (s *state) evalCall(", "func (s *state) evalCallOld(",
		"func isTrue(val reflect.Value) (truth, ok bool) {", "func isTrueOld(val reflect.Value) (truth, ok bool) {",
	)

	htmlTemplateReplacers = strings.NewReplacer(
		`. "html/template"`, `. "github.com/gohugoio/hugo/tpl/internal/go_templates/htmltemplate"`,
		`"html/template"`, `template "github.com/gohugoio/hugo/tpl/internal/go_templates/htmltemplate"`,
		"\"text/template\"\n", "template \"github.com/gohugoio/hugo/tpl/internal/go_templates/texttemplate\"\n",
		`"html/template"`, `htmltemplate "html/template"`,
		`"fmt"`, `htmltemplate "html/template"`,
	)
)

func commonReplace(name, content string) string {
	if strings.HasSuffix(name, "_test.go") {
		content = strings.Replace(content, "package template\n", `// +build go1.13,!windows

package template
`, 1)
		content = strings.Replace(content, "package template_test\n", `// +build go1.13

package template_test
`, 1)

		content = strings.Replace(content, "package parse\n", `// +build go1.13

package parse
`, 1)

	}

	return content

}

var goPackages = []goPackage{
	goPackage{srcPkg: "text/template", dstPkg: "texttemplate",
		replacer: func(name, content string) string { return textTemplateReplacers.Replace(commonReplace(name, content)) }},
	goPackage{srcPkg: "html/template", dstPkg: "htmltemplate", replacer: func(name, content string) string {
		if strings.HasSuffix(name, "content.go") {
			// Remove template.HTML types. We need to use the Go types.
			content = removeAll(`(?s)// Strings of content.*?\)\n`, content)
		}

		content = commonReplace(name, content)

		return htmlTemplateReplacers.Replace(content)
	},
		rewriter: func(name string) {
			for _, s := range []string{"CSS", "HTML", "HTMLAttr", "JS", "JSStr", "URL", "Srcset"} {
				rewrite(name, fmt.Sprintf("%s -> htmltemplate.%s", s, s))
			}
			rewrite(name, `"text/template/parse" -> "github.com/gohugoio/hugo/tpl/internal/go_templates/texttemplate/parse"`)
		}},
	goPackage{srcPkg: "internal/fmtsort", dstPkg: "fmtsort", rewriter: func(name string) {
		rewrite(name, `"internal/fmtsort" -> "github.com/gohugoio/hugo/tpl/internal/go_templates/fmtsort"`)
	}},
}

var fs = afero.NewOsFs()

// Removes all non-Hugo files in the go_templates folder.
func cleanFork() {
	must(filepath.Walk(filepath.Join(forkRoot), func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && len(path) > 10 && !strings.Contains(path, "hugo") {
			must(fs.Remove(path))
		}
		return nil
	}))
}

func must(err error, what ...string) {
	if err != nil {
		log.Fatal(what, " ERROR: ", err)
	}
}

func copyGoPackage(dst, src string) {
	from := filepath.Join(goSource, src)
	to := filepath.Join(forkRoot, dst)
	fmt.Println("Copy", from, "to", to)
	must(hugio.CopyDir(fs, from, to, func(s string) bool { return true }))
}

func doWithGoFiles(dir string,
	rewrite func(name string),
	transform func(name, in string) string) {
	if rewrite == nil && transform == nil {
		return
	}
	must(filepath.Walk(filepath.Join(forkRoot, dir), func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ".go") || strings.Contains(path, "hugo_") {
			return nil
		}

		fmt.Println("Handle", path)

		if rewrite != nil {
			rewrite(path)
		}

		if transform == nil {
			return nil
		}

		data, err := ioutil.ReadFile(path)
		must(err)
		f, err := os.Create(path)
		must(err)
		defer f.Close()
		_, err = f.WriteString(transform(path, string(data)))
		must(err)

		return nil
	}))
}

func removeAll(expression, content string) string {
	re := regexp.MustCompile(expression)
	return re.ReplaceAllString(content, "")

}

func rewrite(filename, rule string) {
	cmf := exec.Command("gofmt", "-w", "-r", rule, filename)
	out, err := cmf.CombinedOutput()
	if err != nil {
		log.Fatal("gofmt failed:", string(out))
	}
}

func goimports(dir string) {
	cmf := exec.Command("goimports", "-w", dir)
	out, err := cmf.CombinedOutput()
	if err != nil {
		log.Fatal("goimports failed:", string(out))
	}
}

func gofmt(dir string) {
	cmf := exec.Command("gofmt", "-w", dir)
	out, err := cmf.CombinedOutput()
	if err != nil {
		log.Fatal("gofmt failed:", string(out))
	}
}
