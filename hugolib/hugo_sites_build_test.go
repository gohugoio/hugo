package hugolib

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/htesting"
	"github.com/gohugoio/hugo/resources/kinds"

	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/spf13/afero"
)

func TestMultiSitesMainLangInRoot(t *testing.T) {
	files := `
-- hugo.toml --
defaultContentLanguage = "fr"
defaultContentLanguageInSubdir = false
disableKinds = ["taxonomy", "term"]
[languages]
[languages.en]
weight = 1
[languages.fr]
weight = 2
-- content/sect/doc1.en.md --
---
title: doc1 en
---
-- content/sect/doc1.fr.md --
---
title: doc1 fr
slug: doc1-fr
---
-- layouts/_default/single.html --
Single: {{ .Title }}|{{ .Lang }}|{{ .RelPermalink }}|

`
	b := Test(t, files)
	b.AssertFileContent("public/sect/doc1-fr/index.html", "Single: doc1 fr|fr|/sect/doc1-fr/|")
	b.AssertFileContent("public/en/sect/doc1/index.html", "Single: doc1 en|en|/en/sect/doc1/|")
}

func TestMultiSitesWithTwoLanguages(t *testing.T) {
	t.Parallel()

	c := qt.New(t)
	b := newTestSitesBuilder(t).WithConfigFile("toml", `

defaultContentLanguage = "nn"

[languages]
[languages.nn]
languageName = "Nynorsk"
weight = 1
title = "Tittel på Nynorsk"
[languages.nn.params]
p1 = "p1nn"

[languages.en]
title = "Title in English"
languageName = "English"
weight = 2
[languages.en.params]
p1 = "p1en"
`)

	b.CreateSites()
	b.Build(BuildCfg{SkipRender: true})
	sites := b.H.Sites

	c.Assert(len(sites), qt.Equals, 2)

	nnSite := sites[0]
	nnHome := nnSite.getPageOldVersion(kinds.KindHome)
	c.Assert(len(nnHome.AllTranslations()), qt.Equals, 2)
	c.Assert(len(nnHome.Translations()), qt.Equals, 1)
	c.Assert(nnHome.IsTranslated(), qt.Equals, true)

	enHome := sites[1].getPageOldVersion(kinds.KindHome)

	p1, err := enHome.Param("p1")
	c.Assert(err, qt.IsNil)
	c.Assert(p1, qt.Equals, "p1en")

	p1, err = nnHome.Param("p1")
	c.Assert(err, qt.IsNil)
	c.Assert(p1, qt.Equals, "p1nn")
}

// https://github.com/gohugoio/hugo/issues/4706
func TestContentStressTest(t *testing.T) {
	b := newTestSitesBuilder(t)

	numPages := 500

	contentTempl := `
---
%s
title: %q
weight: %d
multioutput: %t
---

# Header

CONTENT

The End.
`

	contentTempl = strings.Replace(contentTempl, "CONTENT", strings.Repeat(`
	
## Another header

Some text. Some more text.

`, 100), -1)

	var content []string
	defaultOutputs := `outputs: ["html", "json", "rss" ]`

	for i := 1; i <= numPages; i++ {
		outputs := defaultOutputs
		multioutput := true
		if i%3 == 0 {
			outputs = `outputs: ["json"]`
			multioutput = false
		}
		section := "s1"
		if i%10 == 0 {
			section = "s2"
		}
		content = append(content, []string{fmt.Sprintf("%s/page%d.md", section, i), fmt.Sprintf(contentTempl, outputs, fmt.Sprintf("Title %d", i), i, multioutput)}...)
	}

	content = append(content, []string{"_index.md", fmt.Sprintf(contentTempl, defaultOutputs, fmt.Sprintf("Home %d", 0), 0, true)}...)
	content = append(content, []string{"s1/_index.md", fmt.Sprintf(contentTempl, defaultOutputs, fmt.Sprintf("S %d", 1), 1, true)}...)
	content = append(content, []string{"s2/_index.md", fmt.Sprintf(contentTempl, defaultOutputs, fmt.Sprintf("S %d", 2), 2, true)}...)

	b.WithSimpleConfigFile()
	b.WithTemplates("layouts/_default/single.html", `Single: {{ .Content }}|RelPermalink: {{ .RelPermalink }}|Permalink: {{ .Permalink }}`)
	b.WithTemplates("layouts/_default/myview.html", `View: {{ len .Content }}`)
	b.WithTemplates("layouts/_default/single.json", `Single JSON: {{ .Content }}|RelPermalink: {{ .RelPermalink }}|Permalink: {{ .Permalink }}`)
	b.WithTemplates("layouts/_default/list.html", `
Page: {{ .Paginator.PageNumber }}
P: {{ with .File }}{{ path.Join .Path }}{{ end }}
List: {{ len .Paginator.Pages }}|List Content: {{ len .Content }}
{{ $shuffled :=  where .Site.RegularPages "Params.multioutput" true | shuffle }}
{{ $first5 := $shuffled | first 5 }}
L1: {{ len .Site.RegularPages }} L2: {{ len $first5 }}
{{ range $i, $e := $first5 }}
Render {{ $i }}: {{ .Render "myview" }}
{{ end }}
END
`)

	b.WithContent(content...)

	b.CreateSites().Build(BuildCfg{})

	contentMatchers := []string{"<h2 id=\"another-header\">Another header</h2>", "<h2 id=\"another-header-99\">Another header</h2>", "<p>The End.</p>"}

	for i := 1; i <= numPages; i++ {
		if i%3 != 0 {
			section := "s1"
			if i%10 == 0 {
				section = "s2"
			}
			checkContent(b, fmt.Sprintf("public/%s/page%d/index.html", section, i), contentMatchers...)
		}
	}

	for i := 1; i <= numPages; i++ {
		section := "s1"
		if i%10 == 0 {
			section = "s2"
		}
		checkContent(b, fmt.Sprintf("public/%s/page%d/index.json", section, i), contentMatchers...)
	}

	checkContent(b, "public/s1/index.html", "P: s1/_index.md\nList: 10|List Content: 8132\n\n\nL1: 500 L2: 5\n\nRender 0: View: 8132\n\nRender 1: View: 8132\n\nRender 2: View: 8132\n\nRender 3: View: 8132\n\nRender 4: View: 8132\n\nEND\n")
	checkContent(b, "public/s2/index.html", "P: s2/_index.md\nList: 10|List Content: 8132", "Render 4: View: 8132\n\nEND")
	checkContent(b, "public/index.html", "P: _index.md\nList: 10|List Content: 8132", "4: View: 8132\n\nEND")

	// Check paginated pages
	for i := 2; i <= 9; i++ {
		checkContent(b, fmt.Sprintf("public/page/%d/index.html", i), fmt.Sprintf("Page: %d", i), "Content: 8132\n\n\nL1: 500 L2: 5\n\nRender 0: View: 8132", "Render 4: View: 8132\n\nEND")
	}
}

func checkContent(s *sitesBuilder, filename string, matches ...string) {
	s.T.Helper()
	content := readWorkingDir(s.T, s.Fs, filename)
	for _, match := range matches {
		if !strings.Contains(content, match) {
			s.Fatalf("No match for\n%q\nin content for %s\n%q\nDiff:\n%s", match, filename, content, htesting.DiffStrings(content, match))
		}
	}
}

func TestTranslationsFromContentToNonContent(t *testing.T) {
	b := newTestSitesBuilder(t)
	b.WithConfigFile("toml", `

baseURL = "http://example.com/"

defaultContentLanguage = "en"

[languages]
[languages.en]
weight = 10
contentDir = "content/en"
[languages.nn]
weight = 20
contentDir = "content/nn"


`)

	b.WithContent("en/mysection/_index.md", `
---
Title: My Section
---

`)

	b.WithContent("en/_index.md", `
---
Title: My Home
---

`)

	b.WithContent("en/categories/mycat/_index.md", `
---
Title: My MyCat
---

`)

	b.WithContent("en/categories/_index.md", `
---
Title: My categories
---

`)

	for _, lang := range []string{"en", "nn"} {
		b.WithContent(lang+"/mysection/page.md", `
---
Title: My Page
categories: ["mycat"]
---

`)
	}

	b.Build(BuildCfg{})

	for _, path := range []string{
		"/",
		"/mysection",
		"/categories",
		"/categories/mycat",
	} {
		t.Run(path, func(t *testing.T) {
			c := qt.New(t)

			s1, _ := b.H.Sites[0].getPage(nil, path)
			s2, _ := b.H.Sites[1].getPage(nil, path)

			c.Assert(s1, qt.Not(qt.IsNil))
			c.Assert(s2, qt.Not(qt.IsNil))

			c.Assert(len(s1.Translations()), qt.Equals, 1)
			c.Assert(len(s2.Translations()), qt.Equals, 1)
			c.Assert(s1.Translations()[0], qt.Equals, s2)
			c.Assert(s2.Translations()[0], qt.Equals, s1)

			m1 := s1.Translations().MergeByLanguage(s2.Translations())
			m2 := s2.Translations().MergeByLanguage(s1.Translations())

			c.Assert(len(m1), qt.Equals, 1)
			c.Assert(len(m2), qt.Equals, 1)
		})
	}
}

func writeSource(t testing.TB, fs *hugofs.Fs, filename, content string) {
	t.Helper()
	writeToFs(t, fs.Source, filename, content)
}

func writeToFs(t testing.TB, fs afero.Fs, filename, content string) {
	t.Helper()
	if err := afero.WriteFile(fs, filepath.FromSlash(filename), []byte(content), 0o755); err != nil {
		t.Fatalf("Failed to write file: %s", err)
	}
}

func readWorkingDir(t testing.TB, fs *hugofs.Fs, filename string) string {
	t.Helper()
	return readFileFromFs(t, fs.WorkingDirReadOnly, filename)
}

func workingDirExists(fs *hugofs.Fs, filename string) bool {
	b, err := helpers.Exists(filename, fs.WorkingDirReadOnly)
	if err != nil {
		panic(err)
	}
	return b
}

func readFileFromFs(t testing.TB, fs afero.Fs, filename string) string {
	t.Helper()
	filename = filepath.Clean(filename)
	b, err := afero.ReadFile(fs, filename)
	if err != nil {
		// Print some debug info
		hadSlash := strings.HasPrefix(filename, helpers.FilePathSeparator)
		start := 0
		if hadSlash {
			start = 1
		}
		end := start + 1

		parts := strings.Split(filename, helpers.FilePathSeparator)
		if parts[start] == "work" {
			end++
		}

		/*
			root := filepath.Join(parts[start:end]...)
			if hadSlash {
				root = helpers.FilePathSeparator + root
			}

			helpers.PrintFs(fs, root, os.Stdout)
		*/

		t.Fatalf("Failed to read file: %s", err)
	}
	return string(b)
}

const testPageTemplate = `---
title: "%s"
publishdate: "%s"
weight: %d
---
# Doc %s
`

func newTestPage(title, date string, weight int) string {
	return fmt.Sprintf(testPageTemplate, title, date, weight, title)
}

func TestRebuildOnAssetChange(t *testing.T) {
	b := newTestSitesBuilder(t).Running().WithLogger(loggers.NewDefault())
	b.WithTemplatesAdded("index.html", `
{{ (resources.Get "data.json").Content }}
`)
	b.WithSourceFile("assets/data.json", "orig data")

	b.Build(BuildCfg{})
	b.AssertFileContent("public/index.html", `orig data`)

	b.EditFiles("assets/data.json", "changed data")

	b.Build(BuildCfg{})
	b.AssertFileContent("public/index.html", `changed data`)
}
