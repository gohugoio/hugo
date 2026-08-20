package main

import (
	"fmt"
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
	/*
		Previously: with 2dc996f71b0ebafb77e64433e58333e049488a3c go1.26.3
		Current:  8af21751f0 [release-branch.go1.27] go1.27.0

		Note that the upgrade here is mostly automatic, but:

		* testenv.go is a stubbed Hugo version and is never overwritten; if the template tests start using new helpers from it, stub them in by hand.
		* Some of the replacements below match exact upstream source; if upstream drifts, the build breaks and they need updating.
		* Some test code depends on Go internals we don't fork; remove or stub these out to make it build.
		* We're only patching the execution part of the template packages, so it's also good to check the execution package's Git history to check for valuable changes in our patched files.
	*/
	fmt.Println("Forking ...")
	defer fmt.Println("Done ...")

	cleanFork()

	htmlRoot := filepath.Join(forkRoot, "htmltemplate")

	for _, pkg := range goPackages {
		copyGoPackage(pkg.dstPkg, pkg.srcPkg, pkg.skip)
	}

	for _, pkg := range goPackages {
		doWithGoFiles(pkg.dstPkg, pkg.rewriter, pkg.replacer)
	}

	goimports(htmlRoot)
	gofmt(forkRoot)
}

const (
	goSource = "/Users/bep/dev/go/misc/go/src"
	forkRoot = "../../tpl/internal/go_templates"
)

type goPackage struct {
	srcPkg   string
	dstPkg   string
	replacer func(name, content string) string
	rewriter func(name string)
	skip     func(name string) bool
}

var (
	textTemplateReplacers = strings.NewReplacer(
		`"text/template/`, `"github.com/gohugoio/hugo/tpl/internal/go_templates/texttemplate/`,
		`"internal/fmtsort"`, `"github.com/gohugoio/hugo/tpl/internal/go_templates/fmtsort"`,
		`"internal/testenv"`, `"github.com/gohugoio/hugo/tpl/internal/go_templates/testenv"`,
		"TestLinkerGC", "_TestLinkerGC",
		"{new(int), true},", "//{new(int), true}, // Commented out for Hugo. We have a slightly different view on ... the truth.",
		// Rename types and function that we want to overload.
		"type state struct", "type stateOld struct",
		"func (s *state) evalFunction", "func (s *state) evalFunctionOld",
		"func (s *state) evalField(", "func (s *state) evalFieldOld(",
		"func (s *state) evalCall(", "func (s *state) evalCallOld(",
		"func (s *state) walkTemplate(", "func (s *state) walkTemplateOld(",
		"func (s *state) validateType(", "func (s *state) _validateType(",
		"func isTrue(val reflect.Value) (truth, ok bool) {", "func isTrueOld(val reflect.Value) (truth, ok bool) {",
	)

	testEnvReplacers = strings.NewReplacer(
		`"internal/cfg"`, `"github.com/gohugoio/hugo/tpl/internal/go_templates/cfg"`,
	)

	htmlTemplateReplacers = strings.NewReplacer(
		`. "html/template"`, `. "github.com/gohugoio/hugo/tpl/internal/go_templates/htmltemplate"`,
		`"html/template"`, `template "github.com/gohugoio/hugo/tpl/internal/go_templates/htmltemplate"`,
		"\"text/template\"\n", "template \"github.com/gohugoio/hugo/tpl/internal/go_templates/texttemplate\"\n",
		`"html/template"`, `htmltemplate "html/template"`,
		`"fmt"`, `htmltemplate "html/template"`,
		`"internal/testenv"`, `"github.com/gohugoio/hugo/tpl/internal/go_templates/testenv"`,
		// Renamed so hugo_template.go can wrap it.
		"func indirect(", "func doIndirect(",
		`t.Skip("this test currently fails with -race; see issue #39807")`, `// t.Skip("this test currently fails with -race; see issue #39807")`,
	)

	// We don't want the internal/godebug dependency; meta content URL escaping is always on.
	escapeGodebugReplacers = strings.NewReplacer(
		`var debugAllowActionJSTmpl = godebug.New("jstmpllitinterp")`, ``,
		`var htmlmetacontenturlescape = godebug.New("htmlmetacontenturlescape")`, `var htmlmetacontenturlescape = true`,
		`if htmlmetacontenturlescape.Value() != "0" {`, `if htmlmetacontenturlescape {`,
	)
)

func commonReplace(name, content string) string {
	if strings.HasSuffix(name, "_test.go") {
		content = strings.Replace(content, "package template\n", `//go:build !windows

package template
`, 1)
	}

	return content
}

var goPackages = []goPackage{
	{
		srcPkg: "text/template", dstPkg: "texttemplate",
		replacer: func(name, content string) string { return textTemplateReplacers.Replace(commonReplace(name, content)) },
	},
	{
		srcPkg: "html/template", dstPkg: "htmltemplate", replacer: func(name, content string) string {
			if strings.HasSuffix(name, "content.go") {
				// Remove template.HTML types. We need to use the Go types.
				content = removeAll(`(?s)// Strings of content.*?\)\n`, content)
			}

			if strings.HasSuffix(name, "escape.go") {
				content = escapeGodebugReplacers.Replace(content)
				// Drop the else branch calling IncNonDefault; goimports removes the then-unused godebug import.
				content = removeAll(` else \{\n(?:\t*//.*\n)*\t*htmlmetacontenturlescape\.IncNonDefault\(\)\n\t*\}`, content)
			}

			if strings.HasSuffix(name, "escape_test.go") {
				// Tests the GODEBUG=htmlmetacontenturlescape=0 path, which we hard code to on.
				content = removeAll(`(?s)func TestMetaContentEscapeGODEBUG.*?\n\}\n`, content)
			}

			content = commonReplace(name, content)

			return htmlTemplateReplacers.Replace(content)
		},
		rewriter: func(name string) {
			for _, s := range []string{"CSS", "HTML", "HTMLAttr", "JS", "JSStr", "URL", "Srcset"} {
				rewrite(name, fmt.Sprintf("%s -> htmltemplate.%s", s, s))
			}
			rewrite(name, `"text/template/parse" -> "github.com/gohugoio/hugo/tpl/internal/go_templates/texttemplate/parse"`)
		},
	},
	{srcPkg: "internal/fmtsort", dstPkg: "fmtsort", rewriter: func(name string) {
		rewrite(name, `"internal/fmtsort" -> "github.com/gohugoio/hugo/tpl/internal/go_templates/fmtsort"`)
	}},
	{
		srcPkg: "internal/testenv", dstPkg: "testenv",
		replacer: func(name, content string) string { return testEnvReplacers.Replace(content) }, rewriter: func(name string) {
			rewrite(name, `"internal/testenv" -> "github.com/gohugoio/hugo/tpl/internal/go_templates/testenv"`)
		},
		// testenv.go is a heavily stubbed Hugo version; keep it. The tests test the stubbed away parts.
		skip: func(name string) bool {
			base := filepath.Base(name)
			return base == "testenv.go" || base == "testenv_test.go"
		},
	},
	{srcPkg: "internal/cfg", dstPkg: "cfg", rewriter: func(name string) {
		rewrite(name, `"internal/cfg" -> "github.com/gohugoio/hugo/tpl/internal/go_templates/cfg"`)
	}},
}

var fs = afero.NewOsFs()

// Removes all non-Hugo files in the go_templates folder.
func cleanFork() {
	keepRe := regexp.MustCompile(`(?i)hugo|staticcheck\.conf|^testenv\.go$`)
	must(filepath.Walk(filepath.Join(forkRoot), func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || len(path) <= 10 {
			return nil
		}

		if keepRe.MatchString(info.Name()) {
			return nil
		}

		must(fs.Remove(path))

		return nil
	}))
}

func must(err error, what ...string) {
	if err != nil {
		log.Fatal(what, " ERROR: ", err)
	}
}

func copyGoPackage(dst, src string, skip func(name string) bool) {
	from := filepath.Join(goSource, src)
	to := filepath.Join(forkRoot, dst)
	fmt.Println("Copy", from, "to", to)
	must(hugio.CopyDir(fs, from, to, func(s string) bool { return skip == nil || !skip(s) }))
}

func doWithGoFiles(dir string,
	rewrite func(name string),
	transform func(name, in string) string,
) {
	if rewrite == nil && transform == nil {
		return
	}
	must(filepath.Walk(filepath.Join(forkRoot, dir), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
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

		data, err := os.ReadFile(path)
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
	// Needs go install golang.org/x/tools/cmd/goimports@latest
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
