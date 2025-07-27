package tplimpl_test

import (
	"context"
	"io"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/hugolib"
	"github.com/gohugoio/hugo/resources/kinds"
	"github.com/gohugoio/hugo/resources/page"
	"github.com/gohugoio/hugo/tpl/tplimpl"
)

// Old as in before Hugo v0.146.0.
func TestLayoutsOldSetup(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
defaultContentLanguage = "en"
defaultContentLanguageInSubdir = true
[languages]
[languages.en]
title = "Title in English"
weight = 1
[languages.nn]
title = "Tittel på nynorsk"
weight = 2
-- layouts/index.html --
Home.
{{ template "_internal/twitter_cards.html" . }}
-- layouts/_default/single.html --
Single.
-- layouts/_default/single.nn.html --
Single NN.
-- layouts/_default/list.html --
List HTML.
-- layouts/docs/list-baseof.html --
Docs Baseof List HTML.
{{ block "main" . }}Docs Baseof List HTML main block.{{ end }}
-- layouts/docs/list.section.html --
{{ define "main" }}
Docs List HTML.
{{ end }}
-- layouts/_default/list.json --
List JSON.
-- layouts/_default/list.rss.xml --
List RSS.
-- layouts/_default/list.nn.rss.xml --
List NN RSS.
-- layouts/_default/baseof.html --
Base.
-- layouts/partials/mypartial.html --
Partial.
-- layouts/shortcodes/myshortcode.html --
Shortcode.
-- content/docs/p1.md --
---
title: "P1"
---

	`

	b := hugolib.Test(t, files)

	//	b.DebugPrint("", tplimpl.CategoryBaseof)

	b.AssertFileContent("public/en/docs/index.html", "Docs Baseof List HTML.\n\nDocs List HTML.")
}

func TestLayoutsOldSetupBaseofPrefix(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
-- layouts/_default/layout1-baseof.html --
Baseof layout1. {{ block "main" . }}{{ end }}
-- layouts/_default/layout2-baseof.html --
Baseof layout2. {{ block "main" . }}{{ end }}
-- layouts/_default/layout1.html --
{{ define "main" }}Layout1. {{ .Title }}{{ end }}
-- layouts/_default/layout2.html --
{{ define "main" }}Layout2. {{ .Title }}{{ end }}
-- content/p1.md --
---
title: "P1"
layout: "layout1"
---
-- content/p2.md --
---
title: "P2"
layout: "layout2"
---
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/p1/index.html", "Baseof layout1. Layout1. P1")
	b.AssertFileContent("public/p2/index.html", "Baseof layout2. Layout2. P2")
}

func TestLayoutsOldSetupTaxonomyAndTerm(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
[taxonomies]
cat = 'cats'
dog = 'dogs'
# Templates for term taxonomy, old setup.
-- layouts/dogs/terms.html --
Dogs Terms. Most specific taxonomy template.
-- layouts/taxonomy/terms.html --
Taxonomy Terms. Down the list.
# Templates for term term, old setup.
-- layouts/dogs/term.html --
Dogs Term. Most specific term template.
-- layouts/term/term.html --
Term Term. Down the list.
-- layouts/dogs/max/list.html --
max: {{ .Title }}
-- layouts/_default/list.html --
Default list.
-- layouts/_default/single.html --
Default single.
-- content/p1.md --
---
title: "P1"
dogs: ["luna", "daisy", "max"]
---

`
	b := hugolib.Test(t, files, hugolib.TestOptWarn())

	b.AssertLogContains("! WARN")

	b.AssertFileContent("public/dogs/index.html", "Dogs Terms. Most specific taxonomy template.")
	b.AssertFileContent("public/dogs/luna/index.html", "Dogs Term. Most specific term template.")
	b.AssertFileContent("public/dogs/max/index.html", "max: Max") // layouts/dogs/max/list.html wins over layouts/term/term.html because of distance.
}

func TestLayoutsOldSetupCustomRSS(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ["taxonomy", "term", "page"]
[outputs]
home = ["rss"]
-- layouts/_default/list.rss.xml --
List RSS.
`
	b := hugolib.Test(t, files)
	b.AssertFileContent("public/index.xml", "List RSS.")
}

var newSetupTestSites = `
-- hugo.toml --
defaultContentLanguage = "en"
defaultContentLanguageInSubdir = true
[languages]
[languages.en]
title = "Title in English"
weight = 1
[languages.nn]
title = "Tittel på nynorsk"
weight = 2
[languages.fr]
title = "Titre en français"
weight = 3

[outputs]
home     = ["html", "rss", "redir"]

[outputFormats]
[outputFormats.redir]
mediatype   = "text/plain"
baseName    = "_redirects"
isPlainText = true
-- layouts/404.html --
{{ define "main" }}
404.
{{ end }}
-- layouts/home.html --
{{ define "main" }}
Home: {{ .Title }}|{{ .Content }}|
Inline Partial: {{ partial "my-inline-partial.html" . }}
{{ end }}
{{ define "hero" }}
Home hero.
{{ end }}
{{ define "partials/my-inline-partial.html" }}
{{ $value := 32 }}
{{ return $value }}
{{ end }}
-- layouts/index.redir --
Redir.
-- layouts/single.html --
{{ define "main" }}
Single needs base.
{{ end }}
-- layouts/foo/bar/single.html --
{{ define "main" }}
Single sub path.
{{ end }}
-- layouts/_markup/render-codeblock.html --
Render codeblock.
-- layouts/_markup/render-blockquote.html --
Render blockquote.
-- layouts/_markup/render-codeblock-go.html --
 Render codeblock go.
-- layouts/_markup/render-link.html --
Link: {{ .Destination | safeURL }}
-- layouts/foo/baseof.html --
Base sub path.{{ block "main" . }}{{ end }}
-- layouts/foo/bar/baseof.page.html --
Base sub path.{{ block "main" . }}{{ end }}
-- layouts/list.html --
{{ define "main" }}
List needs base.
{{ end }}
-- layouts/section.html --
Section.
-- layouts/mysectionlayout.section.fr.amp.html --
Section with layout.
-- layouts/baseof.html --
Base.{{ block "main" . }}{{ end }}
Hero:{{ block "hero" . }}{{ end }}:
{{ with (templates.Defer (dict "key" "global")) }}
Defer Block.
{{ end }}
-- layouts/baseof.fr.html --
Base fr.{{ block "main" . }}{{ end }}
-- layouts/baseof.term.html --
Base term.
-- layouts/baseof.section.fr.amp.html --
Base with identifiers.{{ block "main" . }}{{ end }}
-- layouts/partials/mypartial.html --
Partial. {{ partial "_inline/my-inline-partial-in-partial-with-no-ext" . }}
{{ define "partials/_inline/my-inline-partial-in-partial-with-no-ext" }}
Partial in partial.
{{ end }}
-- layouts/partials/returnfoo.html --
{{ $v := "foo" }}
{{ return $v }}
-- layouts/shortcodes/myshortcode.html --
Shortcode. {{ partial "mypartial.html" . }}|return:{{ partial "returnfoo.html" . }}|
-- content/_index.md --
---
title: Home sweet home!
---

{{< myshortcode >}}

> My blockquote.


Markdown link: [Foo](/foo)
-- content/p1.md --
---
title: "P1"
---
-- content/foo/bar/index.md --
---
title: "Foo Bar"
---

{{< myshortcode >}}

-- content/single-list.md --
---
title: "Single List"
layout: "list"
---

`

func TestLayoutsType(t *testing.T) {
	files := `
-- hugo.toml --
disableKinds = ["taxonomy", "term"]
-- layouts/list.html --
List.
-- layouts/mysection/single.html --
mysection/single|{{ .Title }}
-- layouts/mytype/single.html --
mytype/single|{{ .Title }}
-- content/mysection/_index.md --
-- content/mysection/mysubsection/_index.md --
-- content/mysection/mysubsection/p1.md --
---
title: "P1"
---
-- content/mysection/mysubsection/p2.md --
---
title: "P2"
type: "mytype"
---

`

	b := hugolib.Test(t, files, hugolib.TestOptWarn())

	b.AssertLogContains("! WARN")

	b.AssertFileContent("public/mysection/mysubsection/p1/index.html", "mysection/single|P1")
	b.AssertFileContent("public/mysection/mysubsection/p2/index.html", "mytype/single|P2")
}

// New, as in from Hugo v0.146.0.
func TestLayoutsNewSetup(t *testing.T) {
	const numIterations = 1
	for range numIterations {

		b := hugolib.Test(t, newSetupTestSites, hugolib.TestOptWarn())

		b.AssertLogContains("! WARN")

		b.AssertFileContent("public/en/index.html",
			"Base.\nHome: Home sweet home!|",
			"|Shortcode.\n|",
			"<p>Markdown link: Link: /foo</p>",
			"|return:foo|",
			"Defer Block.",
			"Home hero.",
			"Render blockquote.",
		)

		b.AssertFileContent("public/en/p1/index.html", "Base.\nSingle needs base.\n\nHero::\n\nDefer Block.")
		b.AssertFileContent("public/en/404.html", "404.")
		b.AssertFileContent("public/nn/404.html", "404.")
		b.AssertFileContent("public/fr/404.html", "404.")

	}
}

func TestHomeRSSAndHTMLWithHTMLOnlyShortcode(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ["taxonomy", "term"]
[outputs]
home = ["html", "rss"]
-- layouts/home.html --
Home: {{ .Title }}|{{ .Content }}|
-- layouts/single.html --
Single: {{ .Title }}|{{ .Content }}|
-- layouts/shortcodes/myshortcode.html --
Myshortcode: Count: {{ math.Counter }}|
-- content/p1.md --
---
title: "P1"
---

{{< myshortcode >}}
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/p1/index.html", "Single: P1|Myshortcode: Count: 1|")
	b.AssertFileContent("public/index.xml", "Myshortcode: Count: 1")
}

func TestHomeRSSAndHTMLWithHTMLOnlyRenderHook(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ["taxonomy", "term"]
[outputs]
home = ["html", "rss"]
-- layouts/home.html --
Home: {{ .Title }}|{{ .Content }}|
-- layouts/single.html --
Single: {{ .Title }}|{{ .Content }}|
-- layouts/_markup/render-link.html --
Render Link: {{ math.Counter }}|
-- content/p1.md --
---
title: "P1"
---

Link: [Foo](/foo)
`

	for range 2 {
		b := hugolib.Test(t, files)
		b.AssertFileContent("public/index.xml", "Link: Render Link: 1|")
		b.AssertFileContent("public/p1/index.html", "Single: P1|<p>Link: Render Link: 1|<")
	}
}

func TestCodeblockIssue13864(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['page','rss','section','sitemap','taxonomy','term']
-- content/_index.md --
---
title: home
---

~~~
LANG: none
~~~

~~~go
LANG: go
~~~

~~~go-html-template
LANG: go-html-template
~~~

~~~xy
LANG: xy
~~~

~~~x-y
LANG: x-y
~~~
-- layouts/home.html --
{{ .Content }}
-- layouts/_markup/render-codeblock.html --
{{ .Inner }} LAYOUT: render-codeblock.html|
-- layouts/_markup/render-codeblock-go.html --
{{ .Inner }} LAYOUT: render-codeblock-go.html|
-- layouts/_markup/render-codeblock-go-html-template.html --
{{ .Inner }} LAYOUT: render-codeblock-go-html-template.html|
-- layouts/_markup/render-codeblock-xy.html --
{{ .Inner }} LAYOUT: render-codeblock-xy.html|
-- layouts/_markup/render-codeblock-x-y.html.html --
{{ .Inner }} LAYOUT: render-codeblock-x-y.html|
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html",
		"LANG: none LAYOUT: render-codeblock.html|",                              // pass
		"LANG: go LAYOUT: render-codeblock-go.html|",                             // fail: uses render-codeblock-go-html-template.html
		"LANG: go-html-template LAYOUT: render-codeblock-go-html-template.html|", // fail: uses render-codeblock.html
		"LANG: xy LAYOUT: render-codeblock-xy.html|",                             // pass
		"LANG: x-y LAYOUT: render-codeblock-x-y.html|",                           // fail: uses render-codeblock.html
	)
}

func TestRenderCodeblockSpecificity(t *testing.T) {
	files := `
-- hugo.toml --
-- layouts/_markup/render-codeblock.html --
Render codeblock.|{{ .Inner }}|
-- layouts/_markup/render-codeblock-go.html --
Render codeblock go.|{{ .Inner }}|
-- layouts/single.html --
{{ .Title }}|{{ .Content }}|
-- content/p1.md --
---
title: "P1"
---

§§§
Basic
§§§

§§§ go
Go
§§§

`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/p1/index.html", "P1|Render codeblock.|Basic|Render codeblock go.|Go|")
}

func TestPrintUnusedTemplates(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
baseURL = 'http://example.com/'
printUnusedTemplates=true
-- content/p1.md --
---
title: "P1"
---
{{< usedshortcode >}}
-- layouts/baseof.html --
{{ block "main" . }}{{ end }}
-- layouts/baseof.json --
{{ block "main" . }}{{ end }}
-- layouts/index.html --
{{ define "main" }}FOO{{ end }}
-- layouts/_default/single.json --
-- layouts/_default/single.html --
{{ define "main" }}MAIN /_default/single.html{{ end }}
-- layouts/post/single.html --
{{ define "main" }}MAIN{{ end }}
-- layouts/_partials/usedpartial.html --
-- layouts/_partials/unusedpartial.html --
-- layouts/_shortcodes/usedshortcode.html --
{{ partial "usedpartial.html" }}
-- layouts/shortcodes/unusedshortcode.html --

	`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			NeedsOsFS:   true,
		},
	)
	b.Build()

	b.AssertFileContent("public/p1/index.html", "MAIN /_default/single.html")

	unused := b.H.GetTemplateStore().UnusedTemplates()
	var names []string
	for _, tmpl := range unused {
		if fi := tmpl.Fi; fi != nil {
			names = append(names, fi.Meta().PathInfo.PathNoLeadingSlash())
		}
	}

	b.Assert(names, qt.DeepEquals, []string{"_partials/unusedpartial.html", "shortcodes/unusedshortcode.html", "baseof.json", "post/single.html", "_default/single.json"})
	b.Assert(len(unused), qt.Equals, 5, qt.Commentf("%#v", names))
}

func TestCreateManyTemplateStores(t *testing.T) {
	t.Parallel()
	b := hugolib.Test(t, newSetupTestSites)
	store := b.H.TemplateStore

	for range 70 {
		newStore, err := store.NewFromOpts()
		b.Assert(err, qt.IsNil)
		b.Assert(newStore, qt.Not(qt.IsNil))
	}
}

func BenchmarkLookupPagesLayout(b *testing.B) {
	files := `
-- hugo.toml --
-- layouts/single.html --
{{ define "main" }}
 Main.
{{ end }}
-- layouts/baseof.html --
baseof: {{ block "main" . }}{{ end }}
-- layouts/foo/bar/single.html --
{{ define "main" }}
 Main.
{{ end }}

`
	bb := hugolib.Test(b, files)
	store := bb.H.TemplateStore

	b.ResetTimer()
	b.Run("Single root", func(b *testing.B) {
		q := tplimpl.TemplateQuery{
			Path:     "/baz",
			Category: tplimpl.CategoryLayout,
			Desc:     tplimpl.TemplateDescriptor{Kind: kinds.KindPage, LayoutFromTemplate: "single", OutputFormat: "html"},
		}
		for i := 0; i < b.N; i++ {
			store.LookupPagesLayout(q)
		}
	})

	b.Run("Single sub folder", func(b *testing.B) {
		q := tplimpl.TemplateQuery{
			Path:     "/foo/bar",
			Category: tplimpl.CategoryLayout,
			Desc:     tplimpl.TemplateDescriptor{Kind: kinds.KindPage, LayoutFromTemplate: "single", OutputFormat: "html"},
		}
		for i := 0; i < b.N; i++ {
			store.LookupPagesLayout(q)
		}
	})
}

func BenchmarkNewTemplateStore(b *testing.B) {
	bb := hugolib.Test(b, newSetupTestSites)
	store := bb.H.TemplateStore

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		newStore, err := store.NewFromOpts()
		if err != nil {
			b.Fatal(err)
		}
		if newStore == nil {
			b.Fatal("newStore is nil")
		}
	}
}

func TestLayoutsLookupVariants(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
defaultContentLanguage = "en"
defaultContentLanguageInSubdir = true
[outputs]
home = ["html", "rss"]
page = ["html", "rss",  "amp"]
section = ["html", "rss"]

[languages]
[languages.en]
title = "Title in English"
weight = 1
[languages.nn]
title = "Tittel på nynorsk"
weight = 2
-- layouts/list.xml --
layouts/list.xml
-- layouts/_shortcodes/myshortcode.html --
layouts/shortcodes/myshortcode.html
-- layouts/foo/bar/_shortcodes/myshortcode.html --
layouts/foo/bar/_shortcodes/myshortcode.html
-- layouts/_markup/render-codeblock.html --
layouts/_markup/render-codeblock.html|{{ .Type }}|
-- layouts/_markup/render-codeblock-go.html --
layouts/_markup/render-codeblock-go.html|{{ .Type }}|
-- layouts/single.xml --
layouts/single.xml
-- layouts/single.rss.xml --
layouts/single.rss.xml
-- layouts/single.nn.rss.xml --
layouts/single.nn.rss.xml
-- layouts/list.html --
layouts/list.html
-- layouts/single.html --
layouts/single.html
{{ .Content }}
-- layouts/mylayout.html --
layouts/mylayout.html
-- layouts/mylayout.nn.html --
layouts/mylayout.nn.html
-- layouts/foo/single.rss.xml --
layouts/foo/single.rss.xml
-- layouts/foo/single.amp.html --
layouts/foo/single.amp.html
-- layouts/foo/bar/page.html --
layouts/foo/bar/page.html
-- layouts/foo/bar/baz/single.html --
layouts/foo/bar/baz/single.html
{{ .Content }}
-- layouts/qux/mylayout.html --
layouts/qux/mylayout.html
-- layouts/qux/single.xml --
layouts/qux/single.xml
-- layouts/qux/mylayout.section.html --
layouts/qux/mylayout.section.html
-- content/p.md --
---
---
§§§
code
§§§

§§§ go
code
§§§

{{< myshortcode >}}
-- content/foo/p.md --
-- content/foo/p.nn.md --
-- content/foo/bar/p.md --
-- content/foo/bar/withmylayout.md --
---
layout: mylayout
---
-- content/foo/bar/_index.md --
-- content/foo/bar/baz/p.md --
---
---
{{< myshortcode >}}
-- content/qux/p.md --
-- content/qux/_index.md --
---
layout: mylayout
---
-- content/qux/quux/p.md --
-- content/qux/quux/withmylayout.md --
---
layout: mylayout
---
-- content/qux/quux/withmylayout.nn.md --
---
layout: mylayout
---


`

	b := hugolib.Test(t, files, hugolib.TestOptWarn())

	b.AssertLogContains("! WARN")

	// Single pages.
	// output format: html.
	b.AssertFileContent("public/en/p/index.html", "layouts/single.html",
		"layouts/_markup/render-codeblock.html|",
		"layouts/_markup/render-codeblock-go.html|go|",
		"layouts/shortcodes/myshortcode.html",
	)
	b.AssertFileContent("public/en/foo/p/index.html", "layouts/single.html")
	b.AssertFileContent("public/en/foo/bar/p/index.html", "layouts/foo/bar/page.html")
	b.AssertFileContent("public/en/foo/bar/withmylayout/index.html", "layouts/mylayout.html")
	b.AssertFileContent("public/en/foo/bar/baz/p/index.html", "layouts/foo/bar/baz/single.html", "layouts/foo/bar/_shortcodes/myshortcode.html")
	b.AssertFileContent("public/en/qux/quux/withmylayout/index.html", "layouts/qux/mylayout.html")
	// output format: amp.
	b.AssertFileContent("public/en/amp/p/index.html", "layouts/single.html")
	b.AssertFileContent("public/en/amp/foo/p/index.html", "layouts/foo/single.amp.html")
	// output format: rss.
	b.AssertFileContent("public/en/p/index.xml", "layouts/single.rss.xml")
	b.AssertFileContent("public/en/foo/p/index.xml", "layouts/foo/single.rss.xml")
	b.AssertFileContent("public/nn/foo/p/index.xml", "layouts/single.nn.rss.xml")

	// Note: There is qux/single.xml that's closer, but the one in the root is used becaulse of the output format match.
	b.AssertFileContent("public/en/qux/p/index.xml", "layouts/single.rss.xml")

	// Note.
	b.AssertFileContent("public/nn/qux/quux/withmylayout/index.html", "layouts/mylayout.nn.html")

	// Section pages.
	// output format: html.
	b.AssertFileContent("public/en/foo/index.html", "layouts/list.html")
	b.AssertFileContent("public/en/qux/index.html", "layouts/qux/mylayout.section.html")
	// output format: rss.
	b.AssertFileContent("public/en/foo/index.xml", "layouts/list.xml")
}

func TestLookupOrderIssue13636(t *testing.T) {
	t.Parallel()

	filesTemplate := `
-- hugo.toml --
defaultContentLanguage = "en"
defaultContentLanguageInSubdir = true
[languages]
[languages.en]
weight = 1
[languages.nn]
weight = 2
-- content/s1/p1.en.md --
---
outputs: ["html", "amp", "json"]
---
-- content/s1/p1.nn.md --
---
outputs: ["html", "amp", "json"]
---
-- layouts/L1 --
L1
-- layouts/L2 --
L2
-- layouts/L3 --
L3

`

	tests := []struct {
		Lang       string
		L1         string
		L2         string
		L3         string
		ExpectHTML string
		ExpectAmp  string
		ExpectJSON string
	}{
		{"en", "all.en.html", "all.html", "single.html", "single.html", "single.html", ""},
		{"en", "all.amp.html", "all.html", "page.html", "page.html", "all.amp.html", ""},
		{"en", "all.amp.html", "all.html", "list.html", "all.html", "all.amp.html", ""},
		{"en", "all.en.html", "all.json", "single.html", "single.html", "single.html", "all.json"},
		{"en", "all.en.html", "single.json", "single.html", "single.html", "single.html", "single.json"},
		{"en", "all.en.html", "all.html", "list.html", "all.en.html", "all.en.html", ""},
		{"en", "list.en.html", "list.html", "list.en.html", "", "", ""},
		{"nn", "all.en.html", "all.html", "single.html", "single.html", "single.html", ""},
		{"nn", "all.en.html", "all.nn.html", "single.html", "single.html", "single.html", ""},
		{"nn", "all.en.html", "all.nn.html", "single.nn.html", "single.nn.html", "single.nn.html", ""},
		{"nn", "single.json", "single.nn.json", "all.json", "", "", "single.nn.json"},
		{"nn", "single.json", "single.en.json", "all.nn.json", "", "", "single.json"},
	}

	for i, test := range tests {
		if i != 8 {
			// continue
		}
		files := strings.ReplaceAll(filesTemplate, "L1", test.L1)
		files = strings.ReplaceAll(files, "L2", test.L2)
		files = strings.ReplaceAll(files, "L3", test.L3)
		t.Logf("Test %d: %s %s %s %s", i, test.Lang, test.L1, test.L2, test.L3)

		for range 3 {
			b := hugolib.Test(t, files)
			b.Assert(len(b.H.Sites), qt.Equals, 2)

			var (
				pubhHTML = "public/LANG/s1/p1/index.html"
				pubhAmp  = "public/LANG/amp/s1/p1/index.html"
				pubhJSON = "public/LANG/s1/p1/index.json"
			)

			pubhHTML = strings.ReplaceAll(pubhHTML, "LANG", test.Lang)
			pubhAmp = strings.ReplaceAll(pubhAmp, "LANG", test.Lang)
			pubhJSON = strings.ReplaceAll(pubhJSON, "LANG", test.Lang)

			if test.ExpectHTML != "" {
				b.AssertFileContent(pubhHTML, test.ExpectHTML)
			} else {
				b.AssertFileExists(pubhHTML, false)
			}

			if test.ExpectAmp != "" {
				b.AssertFileContent(pubhAmp, test.ExpectAmp)
			} else {
				b.AssertFileExists(pubhAmp, false)
			}

			if test.ExpectJSON != "" {
				b.AssertFileContent(pubhJSON, test.ExpectJSON)
			} else {
				b.AssertFileExists(pubhJSON, false)
			}
		}
	}
}

func TestLookupShortcodeDepth(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
-- layouts/_shortcodes/myshortcode.html --
layouts/_shortcodes/myshortcode.html
-- layouts/foo/_shortcodes/myshortcode.html --
layouts/foo/_shortcodes/myshortcode.html
-- layouts/single.html --
{{ .Content }}|
-- content/p.md --
---
---
{{< myshortcode >}}
-- content/foo/p.md --
---
---
{{< myshortcode >}}
-- content/foo/bar/p.md --
---
---
{{< myshortcode >}}
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/p/index.html", "layouts/_shortcodes/myshortcode.html")
	b.AssertFileContent("public/foo/p/index.html", "layouts/foo/_shortcodes/myshortcode.html")
	b.AssertFileContent("public/foo/bar/p/index.html", "layouts/foo/_shortcodes/myshortcode.html")
}

func TestLookupShortcodeLayout(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
-- layouts/_shortcodes/myshortcode.single.html --
layouts/_shortcodes/myshortcode.single.html
-- layouts/_shortcodes/myshortcode.list.html --
layouts/_shortcodes/myshortcode.list.html
-- layouts/single.html --
{{ .Content }}|
-- layouts/list.html --
{{ .Content }}|
-- content/_index.md --
---
---
{{< myshortcode >}}
-- content/p.md --
---
---
{{< myshortcode >}}
-- content/foo/p.md --
---
---
{{< myshortcode >}}
-- content/foo/bar/p.md --
---
---
{{< myshortcode >}}
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/p/index.html", "layouts/_shortcodes/myshortcode.single.html")
	b.AssertFileContent("public/index.html", "layouts/_shortcodes/myshortcode.list.html")
}

func TestLayoutAll(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
-- layouts/single.html --
Single.
-- layouts/all.html --
All.
-- content/p1.md --

`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/p1/index.html", "Single.")
	b.AssertFileContent("public/index.html", "All.")
}

func TestLayoutAllNested(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['rss','sitemap','taxonomy','term']
-- content/s1/p1.md --
---
title: p1
---
-- content/s2/p2.md --
---
title: p2
---
-- layouts/single.html --
layouts/single.html
-- layouts/list.html --
layouts/list.html
-- layouts/s1/all.html --
layouts/s1/all.html
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html", "layouts/list.html")
	b.AssertFileContent("public/s1/index.html", "layouts/s1/all.html")
	b.AssertFileContent("public/s1/p1/index.html", "layouts/s1/all.html")
	b.AssertFileContent("public/s2/index.html", "layouts/list.html")
	b.AssertFileContent("public/s2/p2/index.html", "layouts/single.html")
}

func TestPartialHTML(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
-- layouts/all.html --
<html>
<head>
{{ partial "css.html" .}}
</head>
</html>
-- layouts/partials/css.html --
<link rel="stylesheet" href="/css/style.css">
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html", "<link rel=\"stylesheet\" href=\"/css/style.css\">")
}

func TestPartialPlainTextInHTML(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
-- layouts/all.html --
<html>
<head>
{{ partial "mypartial.txt" . }}
</head>
</html>
-- layouts/partials/mypartial.txt --
My <div>partial</div>.
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html", "My &lt;div&gt;partial&lt;/div&gt;.")
}

// Issue #13593.
func TestGoatAndNoGoat(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
-- content/_index.md --
---
title: "Home"
---


§§§
printf "Hello, world!"
§§§


§§§ goat
.---.     .-.       .-.       .-.     .---.
| A +--->| 1 |<--->| 2 |<--->| 3 |<---+ B |
'---'     '-'       '+'       '+'     '---'
§§§



-- layouts/all.html --
{{ .Content }}

`

	b := hugolib.Test(t, files)

	// Basic code block.
	b.AssertFileContent("public/index.html", "<code>printf &#34;Hello, world!&#34;\n</code>")

	// Goat code block.
	b.AssertFileContent("public/index.html", "Menlo,Lucida")
}

// Issue #13595.
func TestGoatAndNoGoatCustomTemplate(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
-- content/_index.md --
---
title: "Home"
---

§§§
printf "Hello, world!"
§§§

§§§ goat
.---.     .-.       .-.       .-.     .---.
| A +--->| 1 |<--->| 2 |<--->| 3 |<---+ B |
'---'     '-'       '+'       '+'     '---'
§§§



-- layouts/_markup/render-codeblock.html --
_markup/render-codeblock.html
-- layouts/all.html --
{{ .Content }}

`

	b := hugolib.Test(t, files)

	// Basic code block.
	b.AssertFileContent("public/index.html", "_markup/render-codeblock.html")

	// Goat code block.
	b.AssertFileContent("public/index.html", "Menlo,Lucida")
}

func TestGoatcustom(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
-- content/_index.md --
---
title: "Home"
---

§§§
printf "Hello, world!"
§§§

§§§ goat
.---.     .-.       .-.       .-.     .---.
| A +--->| 1 |<--->| 2 |<--->| 3 |<---+ B |
'---'     '-'       '+'       '+'     '---'
§§§



-- layouts/_markup/render-codeblock.html --
_markup/render-codeblock.html
-- layouts/_markup/render-codeblock-goat.html --
_markup/render-codeblock-goat.html
-- layouts/all.html --
{{ .Content }}

`

	b := hugolib.Test(t, files)

	// Basic code block.
	b.AssertFileContent("public/index.html", "_markup/render-codeblock.html")

	// Custom Goat code block.
	b.AssertFileContent("public/index.html", "_markup/render-codeblock.html_markup/render-codeblock-goat.html")
}

func TestLookupCodeblockIssue13651(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
-- layouts/all.html --
{{ .Content }}|
-- layouts/_markup/render-codeblock-foo.html --
render-codeblock-foo.html
-- content/_index.md --
---
---

§§§
printf "Hello, world!"
§§§

§§§ foo
printf "Hello, world again!"
§§§
`

	b := hugolib.Test(t, files)

	content := b.FileContent("public/index.html")
	fooCount := strings.Count(content, "render-codeblock-foo.html")
	b.Assert(fooCount, qt.Equals, 1)
}

// Issue #13515
func TestPrintPathWarningOnDotRemoval(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com"
printPathWarnings = true
-- content/v0.124.0.md --
-- content/v0.123.0.md --
-- layouts/all.html --
All.
-- layouts/_default/single.html --
{{ .Title }}|
`

	b := hugolib.Test(t, files, hugolib.TestOptWarn())

	b.AssertLogContains("! Duplicate content path")
}

// Issue #13577.
func TestPrintPathWarningOnInvalidMarkupFilename(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
-- layouts/all.html --
All.
-- layouts/_markup/sitemap.xml --
`
	b := hugolib.Test(t, files, hugolib.TestOptWarn())

	b.AssertLogContains("unrecognized render hook")
}

func TestLayoutNotFound(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
-- layouts/single.html --
Single.
`
	b := hugolib.Test(t, files, hugolib.TestOptWarn())
	b.AssertLogContains("WARN  found no layout file for \"html\" for kind \"home\"")
}

func TestLayoutOverrideThemeWhenThemeOnOldFormatIssue13715(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
theme = "mytheme"
-- layouts/list.html --
 layouts/list.html
-- themes/mytheme/layouts/_default/list.html --
mytheme/layouts/_default/list.html

`

	b := hugolib.Test(t, files)
	b.AssertFileContent("public/index.html", "layouts/list.html")
}

func BenchmarkExecuteWithContext(b *testing.B) {
	files := `
-- hugo.toml --
disableKinds = ["taxonomy", "term", "home"]
-- layouts/all.html --
{{ .Title }}|
{{ partial "p1.html" . }}
-- layouts/_partials/p1.html --
 p1.
{{ partial "p2.html" . }}
{{ partial "p2.html" . }}
{{ partial "p3.html" . }}
{{ partial "p2.html" . }}
{{ partial "p2.html" . }}
{{ partial "p2.html" . }}
{{ partial "p3.html" . }}
-- layouts/_partials/p2.html --
{{ partial "p3.html" . }}
-- layouts/_partials/p3.html --
p3
-- content/p1.md --
`

	bb := hugolib.Test(b, files)

	store := bb.H.TemplateStore

	ti := store.LookupByPath("/all.html")
	bb.Assert(ti, qt.Not(qt.IsNil))
	p := bb.H.Sites[0].RegularPages()[0]
	bb.Assert(p, qt.Not(qt.IsNil))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := store.ExecuteWithContext(context.Background(), ti, io.Discard, p)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkLookupPartial(b *testing.B) {
	files := `
-- hugo.toml --
disableKinds = ["taxonomy", "term", "home"]
-- layouts/all.html --
{{ .Title }}|
-- layouts/_partials/p1.html --
-- layouts/_partials/p2.html --
-- layouts/_partials/p2.json --
-- layouts/_partials/p3.html --
`
	bb := hugolib.Test(b, files)

	store := bb.H.TemplateStore

	for i := 0; i < b.N; i++ {
		fi := store.LookupPartial("p3.html")
		if fi == nil {
			b.Fatal("not found")
		}
	}
}

// Implemented by pageOutput.
type getDescriptorProvider interface {
	GetInternalTemplateBasePathAndDescriptor() (string, tplimpl.TemplateDescriptor)
}

func BenchmarkLookupShortcode(b *testing.B) {
	files := `
-- hugo.toml --
disableKinds = ["taxonomy", "term", "home"]
-- content/toplevelpage.md --
-- content/a/b/c/nested.md --
-- layouts/all.html --
{{ .Title }}|
-- layouts/_shortcodes/s.html --
s1.
-- layouts/_shortcodes/a/b/s.html --
s2.

`
	bb := hugolib.Test(b, files)
	store := bb.H.TemplateStore

	runOne := func(p page.Page) {
		pth, desc := p.(getDescriptorProvider).GetInternalTemplateBasePathAndDescriptor()
		q := tplimpl.TemplateQuery{
			Path:     pth,
			Name:     "s",
			Category: tplimpl.CategoryShortcode,
			Desc:     desc,
		}
		v, err := store.LookupShortcode(q)
		if v == nil || err != nil {
			b.Fatal("not found")
		}
	}

	b.Run("toplevelpage", func(b *testing.B) {
		toplevelpage, _ := bb.H.Sites[0].GetPage("/toplevelpage")
		for i := 0; i < b.N; i++ {
			runOne(toplevelpage)
		}
	})

	b.Run("nestedpage", func(b *testing.B) {
		toplevelpage, _ := bb.H.Sites[0].GetPage("/a/b/c/nested")
		for i := 0; i < b.N; i++ {
			runOne(toplevelpage)
		}
	})
}

func TestStandardLayoutInFrontMatter13588(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['home','page','rss','sitemap','taxonomy','term']
-- content/s1/_index.md --
---
title: s1
---
-- content/s2/_index.md --
---
title: s2
layout: list
---
-- content/s3/_index.md --
---
title: s3
layout: single
---
-- layouts/list.html --
list.html
-- layouts/section.html --
section.html
-- layouts/single.html --
single.html
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/s1/index.html", "section.html")
	b.AssertFileContent("public/s2/index.html", "list.html")   // fail
	b.AssertFileContent("public/s3/index.html", "single.html") // fail
}

func TestIssue13605(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['home','rss','section','sitemap','taxonomy','term']
-- content/s1/p1.md --
---
title: p1
---
{{< sc >}}
-- layouts/s1/_shortcodes/sc.html --
layouts/s1/_shortcodes/sc.html
-- layouts/single.html --
{{ .Content }}
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/s1/p1/index.html", "layouts/s1/_shortcodes/sc.html")
}

func TestSkipDotFiles(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
-- layouts/all.html --
All.
-- layouts/.DS_Store --
{{ foo }}
`

	// Just make sure it doesn't fail.
	hugolib.Test(t, files)
}

func TestPartialsLangIssue13612(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['page','section','sitemap','taxonomy','term']

defaultContentLanguage = 'ru'
defaultContentLanguageInSubdir = true

[languages.ru]
weight = 1

[languages.en]
weight = 2

[outputs]
home = ['html','rss']

-- layouts/_partials/comment.en.html --
layouts/_partials/comment.en.html
-- layouts/_partials/comment.en.xml --
layouts/_partials/comment.en.xml
-- layouts/_partials/comment.ru.html --
layouts/_partials/comment.ru.html
-- layouts/_partials/comment.ru.xml --
layouts/_partials/comment.ru.xml
-- layouts/home.html --
{{ partial (print "comment." (default "ru" .Lang) ".html") . }}
-- layouts/home.rss.xml --
{{ partial (print "comment." (default "ru" .Lang) ".xml") . }}
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/en/index.html", "layouts/_partials/comment.en.html")
	b.AssertFileContent("public/en/index.xml", "layouts/_partials/comment.en.xml")   // fail
	b.AssertFileContent("public/ru/index.html", "layouts/_partials/comment.ru.html") // fail
	b.AssertFileContent("public/ru/index.xml", "layouts/_partials/comment.ru.xml")   // fail
}

func TestLayoutIssue13628(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['home','rss','sitemap','taxonomy','term']
-- content/p1.md --
---
title: p1
layout: foo
---
-- layouts/single.html --
layouts/single.html
-- layouts/list.html --
layouts/list.html
`

	for range 5 {

		b := hugolib.Test(t, files)

		b.AssertFileContent("public/p1/index.html", "layouts/single.html")
	}
}

func TestTemplateLoop(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
-- layouts/_partials/p.html --
p: {{ partial "p.html" . }}
-- layouts/all.html --
{{ partial "p.html" . }}

`
	b, err := hugolib.TestE(t, files)
	b.Assert(err, qt.IsNotNil)
	b.Assert(err.Error(), qt.Contains, "error calling partial: maximum template call stack size exceeded")
}

func TestIssue13630(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['rss','sitemap']
-- content/p1.md --
---
title: p1
layout: foo
---
-- layouts/list.html --
layouts/list.html
-- layouts/taxononmy.html.html --
layouts/taxononmy.html.html
`

	var b *hugolib.IntegrationTestBuilder

	for range 3 {
		b = hugolib.Test(t, files)
		b.AssertFileExists("public/p1/index.html", false)
	}
}

func TestTemplateLoopBlogVsBlogrollIssue13672(t *testing.T) {
	t.Parallel()
	files := `
-- hugo.toml --
-- layouts/blog/_shortcodes/myshortcode.html --
layouts/blog/_shortcodes/myshortcode.html
-- layouts/blog/baseof.html --
blog/baseof.html {{ block "main" . }}{{ end }}
-- layouts/blog/all.html --
{{ define "main" }}blog/all.html|{{ .Content}}{{ end }}
-- layouts/blogroll/_shortcodes/myshortcode.html --
layouts/blogroll/myshortcode.html
-- layouts/blogroll/baseof.html --
{{ block "main" . }}blogroll/baseof.html{{ end }}
-- layouts/blogroll/all.html --
{{ define "main" }}blogroll/all.html|{{ .Content}}{{ end }}
-- content/blog/p1.md --
---
title: p1
---
{{< myshortcode >}}
-- content/blogroll/p1.md --
---
title: p1
---
{{< myshortcode >}}
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/blog/p1/index.html", "blog/baseof.html blog/all.html|layouts/blog/_shortcodes/myshortcode.html")
	b.AssertFileContent("public/blogroll/p1/index.html", "blogroll/all.html|layouts/blogroll/myshortcode.html")
}

// See issue #13668.
func TestPartialPlainTextVsHTML(t *testing.T) {
	t.Parallel()

	/*
		Note that in the below, there's no output format named txt,
		so the isPlainText is fetched from the only output format with that extension.
	*/
	files := `
-- hugo.toml --
-- layouts/_partials/myhtml.html --
<div>myhtml</div>
-- layouts/_partials/mytext.txt --
<div>mytext</div>
-- layouts/all.html --
myhtml: {{ partial "myhtml.html" . }}
mytext: {{ partial "mytext.txt" . }}
mytexts|safeHTML: {{ partial "mytext.txt" . | safeHTML }}
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html",
		"myhtml: <div>myhtml</div>",
		"mytext: &lt;div&gt;mytext&lt;/div&gt;",
		"mytexts|safeHTML: <div>mytext</div>",
	)
}

func TestIssue13351(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['page','rss','section','sitemap','taxonomy','term']
[outputs]
home = ['html','json']
[outputFormats.html]
weight = 1
[outputFormats.json]
weight = 2
-- content/_index.md --
---
title: home
---
a|b
:--|:--
1|2
-- layouts/index.html --
{{ .Content }}
-- layouts/index.json --
{{ .Content }}
`

	b := hugolib.Test(t, files)
	b.AssertFileContent("public/index.html", "<table>")
	b.AssertFileContent("public/index.json", "<table>")

	f := strings.ReplaceAll(files, "weight = 1", "weight = 0")
	b = hugolib.Test(t, f)
	b.AssertFileContent("public/index.html", "<table>")
	b.AssertFileContent("public/index.json", "<table>")

	f = strings.ReplaceAll(files, "weight = 1", "")
	b = hugolib.Test(t, f)
	b.AssertFileContent("public/index.html", "<table>")
	b.AssertFileContent("public/index.json", "<table>")
}

func TestPageKindIssue13868(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['home','rss','section','sitemap','taxonomy','term']
-- content/p1.md --
---
title: p1
---
-- layouts/page.html --
layouts/page.html
-- layouts/p1/page.html --
layouts/p1/page.html
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/p1/index.html", "layouts/p1/page.html")
}
