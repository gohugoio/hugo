package hugolib

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/fortytw2/leaktest"
	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/common/types"
	"github.com/gohugoio/hugo/htesting"
	"github.com/gohugoio/hugo/markup/asciidocext"
	"github.com/gohugoio/hugo/resources/resource_transformers/tocss/dartsass"
	"github.com/gohugoio/hugo/resources/resource_transformers/tocss/scss"
)

const rebuildFilesSimple = `
-- hugo.toml --
baseURL = "https://example.com"
disableKinds = ["term", "taxonomy", "sitemap", "robotstxt", "404"]
disableLiveReload = true
[outputs]
home = ["html"]
section = ["html"]
page = ["html"]
-- content/mysection/_index.md --
---
title: "My Section"
---
-- content/mysection/mysectionbundle/index.md --
---
title: "My Section Bundle"
---
My Section Bundle Content.
-- content/mysection/mysectionbundle/mysectionbundletext.txt --
My Section Bundle Text 2 Content.
-- content/mysection/mysectionbundle/mysectionbundlecontent.md --
---
title: "My Section Bundle Content"
---
My Section Bundle Content Content.
-- content/mysection/_index.md --
---
title: "My Section"
---
-- content/mysection/mysectiontext.txt --
-- content/_index.md --
---
title: "Home"
---
Home Content.
-- content/hometext.txt --
Home Text Content.
-- content/myothersection/myothersectionpage.md --
---
title: "myothersectionpage"
---
myothersectionpage Content.
-- layouts/_default/single.html --
Single: {{ .Title }}|{{ .Content }}$
Resources: {{ range $i, $e := .Resources }}{{ $i }}:{{ .RelPermalink }}|{{ .Content }}|{{ end }}$
Len Resources: {{ len .Resources }}|
-- layouts/_default/list.html --
List: {{ .Title }}|{{ .Content }}$
Len Resources: {{ len .Resources }}|
Resources: {{ range $i, $e := .Resources }}{{ $i }}:{{ .RelPermalink }}|{{ .Content }}|{{ end }}$
-- layouts/shortcodes/foo.html --
Foo.

`

func TestRebuildEditTextFileInLeafBundle(t *testing.T) {
	b := TestRunning(t, rebuildFilesSimple)
	b.AssertFileContent("public/mysection/mysectionbundle/index.html",
		"Resources: 0:/mysection/mysectionbundle/mysectionbundletext.txt|My Section Bundle Text 2 Content.|1:|<p>My Section Bundle Content Content.</p>\n|$")

	b.EditFileReplaceAll("content/mysection/mysectionbundle/mysectionbundletext.txt", "Content.", "Content Edited.").Build()
	b.AssertFileContent("public/mysection/mysectionbundle/index.html",
		"Text 2 Content Edited")
	b.AssertRenderCountPage(1)
	b.AssertRenderCountContent(1)
}

func TestRebuiEditUnmarshaledYamlFileInLeafBundle(t *testing.T) {
	files := `
-- hugo.toml --
baseURL = "https://example.com"
disableLiveReload = true
disableKinds = ["taxonomy", "term", "sitemap", "robotsTXT", "404", "rss"]
-- content/mybundle/index.md --
-- content/mybundle/mydata.yml --
foo: bar
-- layouts/_default/single.html --
MyData: {{ .Resources.Get "mydata.yml" | transform.Unmarshal }}|
`
	b := TestRunning(t, files)

	b.AssertFileContent("public/mybundle/index.html", "MyData: map[foo:bar]")

	b.EditFileReplaceAll("content/mybundle/mydata.yml", "bar", "bar edited").Build()

	b.AssertFileContent("public/mybundle/index.html", "MyData: map[foo:bar edited]")
}

func TestRebuildEditTextFileInHomeBundle(t *testing.T) {
	b := TestRunning(t, rebuildFilesSimple)
	b.AssertFileContent("public/index.html", "Home Content.")
	b.AssertFileContent("public/index.html", "Home Text Content.")

	b.EditFileReplaceAll("content/hometext.txt", "Content.", "Content Edited.").Build()
	b.AssertFileContent("public/index.html", "Home Content.")
	b.AssertFileContent("public/index.html", "Home Text Content Edited.")
	b.AssertRenderCountPage(1)
	b.AssertRenderCountContent(1)
}

func TestRebuildEditTextFileInBranchBundle(t *testing.T) {
	b := TestRunning(t, rebuildFilesSimple)
	b.AssertFileContent("public/mysection/index.html", "My Section")

	b.EditFileReplaceAll("content/mysection/mysectiontext.txt", "Content.", "Content Edited.").Build()
	b.AssertFileContent("public/mysection/index.html", "My Section")
	b.AssertRenderCountPage(1)
	b.AssertRenderCountContent(1)
}

func testRebuildBothWatchingAndRunning(t *testing.T, files string, withB func(b *IntegrationTestBuilder)) {
	t.Helper()
	for _, opt := range []TestOpt{TestOptWatching(), TestOptRunning()} {
		b := Test(t, files, opt)
		withB(b)
	}
}

func TestRebuildRenameTextFileInLeafBundle(t *testing.T) {
	testRebuildBothWatchingAndRunning(t, rebuildFilesSimple, func(b *IntegrationTestBuilder) {
		b.AssertFileContent("public/mysection/mysectionbundle/index.html", "My Section Bundle Text 2 Content.", "Len Resources: 2|")

		b.RenameFile("content/mysection/mysectionbundle/mysectionbundletext.txt", "content/mysection/mysectionbundle/mysectionbundletext2.txt").Build()
		b.AssertFileContent("public/mysection/mysectionbundle/index.html", "mysectionbundletext2", "My Section Bundle Text 2 Content.", "Len Resources: 2|")
		b.AssertRenderCountPage(5)
		b.AssertRenderCountContent(6)
	})
}

func TestRebuilEditContentFileInLeafBundle(t *testing.T) {
	b := TestRunning(t, rebuildFilesSimple)
	b.AssertFileContent("public/mysection/mysectionbundle/index.html", "My Section Bundle Content Content.")
	b.EditFileReplaceAll("content/mysection/mysectionbundle/mysectionbundlecontent.md", "Content Content.", "Content Content Edited.").Build()
	b.AssertFileContent("public/mysection/mysectionbundle/index.html", "My Section Bundle Content Content Edited.")
}

func TestRebuilEditContentFileThenAnother(t *testing.T) {
	b := TestRunning(t, rebuildFilesSimple)
	b.EditFileReplaceAll("content/mysection/mysectionbundle/mysectionbundlecontent.md", "Content Content.", "Content Content Edited.").Build()
	b.AssertFileContent("public/mysection/mysectionbundle/index.html", "My Section Bundle Content Content Edited.")
	b.AssertRenderCountPage(1)
	b.AssertRenderCountContent(2)

	b.EditFileReplaceAll("content/myothersection/myothersectionpage.md", "myothersectionpage Content.", "myothersectionpage Content Edited.").Build()
	b.AssertFileContent("public/myothersection/myothersectionpage/index.html", "myothersectionpage Content Edited")
	b.AssertRenderCountPage(1)
	b.AssertRenderCountContent(1)
}

func TestRebuildRenameTextFileInBranchBundle(t *testing.T) {
	b := TestRunning(t, rebuildFilesSimple)
	b.AssertFileContent("public/mysection/index.html", "My Section")

	b.RenameFile("content/mysection/mysectiontext.txt", "content/mysection/mysectiontext2.txt").Build()
	b.AssertFileContent("public/mysection/index.html", "mysectiontext2", "My Section")
	b.AssertRenderCountPage(2)
	b.AssertRenderCountContent(2)
}

func TestRebuildRenameTextFileInHomeBundle(t *testing.T) {
	b := TestRunning(t, rebuildFilesSimple)
	b.AssertFileContent("public/index.html", "Home Text Content.")

	b.RenameFile("content/hometext.txt", "content/hometext2.txt").Build()
	b.AssertFileContent("public/index.html", "hometext2", "Home Text Content.")
	b.AssertRenderCountPage(3)
}

func TestRebuildRenameDirectoryWithLeafBundle(t *testing.T) {
	b := TestRunning(t, rebuildFilesSimple)
	b.RenameDir("content/mysection/mysectionbundle", "content/mysection/mysectionbundlerenamed").Build()
	b.AssertFileContent("public/mysection/mysectionbundlerenamed/index.html", "My Section Bundle")
	b.AssertRenderCountPage(1)
}

func TestRebuildRenameDirectoryWithBranchBundle(t *testing.T) {
	b := TestRunning(t, rebuildFilesSimple)
	b.RenameDir("content/mysection", "content/mysectionrenamed").Build()
	b.AssertFileContent("public/mysectionrenamed/index.html", "My Section")
	b.AssertFileContent("public/mysectionrenamed/mysectionbundle/index.html", "My Section Bundle")
	b.AssertFileContent("public/mysectionrenamed/mysectionbundle/mysectionbundletext.txt", "My Section Bundle Text 2 Content.")
	b.AssertRenderCountPage(3)
}

func TestRebuildRenameDirectoryWithRegularPageUsedInHome(t *testing.T) {
	files := `
-- hugo.toml --
baseURL = "https://example.com"
disableLiveReload = true
-- content/foo/p1.md --
---
title: "P1"
---
-- layouts/index.html --
Pages: {{ range .Site.RegularPages }}{{ .RelPermalink }}|{{ end }}$
`
	b := TestRunning(t, files)

	b.AssertFileContent("public/index.html", "Pages: /foo/p1/|$")

	b.RenameDir("content/foo", "content/bar").Build()

	b.AssertFileContent("public/index.html", "Pages: /bar/p1/|$")
}

func TestRebuildAddRegularFileRegularPageUsedInHomeMultilingual(t *testing.T) {
	files := `
-- hugo.toml --
baseURL = "https://example.com"
disableLiveReload = true
[languages]
[languages.en]
weight = 1
[languages.nn]
weight = 2
[languages.fr]
weight = 3
[languages.a]
weight = 4
[languages.b]
weight = 5
[languages.c]
weight = 6
[languages.d]
weight = 7
[languages.e]
weight = 8
[languages.f]
weight = 9
[languages.g]
weight = 10
[languages.h]
weight = 11
[languages.i]
weight = 12
[languages.j]
weight = 13
-- content/foo/_index.md --
-- content/foo/data.txt --
-- content/foo/p1.md --
-- content/foo/p1.nn.md --
-- content/foo/p1.fr.md --
-- content/foo/p1.a.md --
-- content/foo/p1.b.md --
-- content/foo/p1.c.md --
-- content/foo/p1.d.md --
-- content/foo/p1.e.md --
-- content/foo/p1.f.md --
-- content/foo/p1.g.md --
-- content/foo/p1.h.md --
-- content/foo/p1.i.md --
-- content/foo/p1.j.md --
-- layouts/index.html --
RegularPages: {{ range .Site.RegularPages }}{{ .RelPermalink }}|{{ end }}$
`
	b := TestRunning(t, files)

	b.AssertFileContent("public/index.html", "RegularPages: /foo/p1/|$")
	b.AssertFileContent("public/nn/index.html", "RegularPages: /nn/foo/p1/|$")
	b.AssertFileContent("public/i/index.html", "RegularPages: /i/foo/p1/|$")

	b.AddFiles("content/foo/p2.md", ``).Build()

	b.AssertFileContent("public/index.html", "RegularPages: /foo/p1/|/foo/p2/|$")
	b.AssertFileContent("public/fr/index.html", "RegularPages: /fr/foo/p1/|$")

	b.AddFiles("content/foo/p2.fr.md", ``).Build()
	b.AssertFileContent("public/fr/index.html", "RegularPages: /fr/foo/p1/|/fr/foo/p2/|$")

	b.AddFiles("content/foo/p2.i.md", ``).Build()
	b.AssertFileContent("public/i/index.html", "RegularPages: /i/foo/p1/|/i/foo/p2/|$")
}

func TestRebuildRenameDirectoryWithBranchBundleFastRender(t *testing.T) {
	recentlyVisited := types.NewEvictingStringQueue(10).Add("/a/b/c/")
	b := TestRunning(t, rebuildFilesSimple, func(cfg *IntegrationTestConfig) { cfg.BuildCfg = BuildCfg{RecentlyVisited: recentlyVisited} })
	b.RenameDir("content/mysection", "content/mysectionrenamed").Build()
	b.AssertFileContent("public/mysectionrenamed/index.html", "My Section")
	b.AssertFileContent("public/mysectionrenamed/mysectionbundle/index.html", "My Section Bundle")
	b.AssertFileContent("public/mysectionrenamed/mysectionbundle/mysectionbundletext.txt", "My Section Bundle Text 2 Content.")
	b.AssertRenderCountPage(3)
}

func TestRebuilErrorRecovery(t *testing.T) {
	b := TestRunning(t, rebuildFilesSimple)
	_, err := b.EditFileReplaceAll("content/mysection/mysectionbundle/index.md", "My Section Bundle Content.", "My Section Bundle Content\n\n\n\n{{< foo }}.").BuildE()

	b.Assert(err, qt.Not(qt.IsNil))
	b.Assert(err.Error(), qt.Contains, filepath.FromSlash(`"/content/mysection/mysectionbundle/index.md:8:9": unrecognized character`))

	// Fix the error
	b.EditFileReplaceAll("content/mysection/mysectionbundle/index.md", "{{< foo }}", "{{< foo >}}").Build()
}

func TestRebuildAddPageListPagesInHome(t *testing.T) {
	files := `
-- hugo.toml --
baseURL = "https://example.com"
disableLiveReload = true
-- content/asection/s1.md --
-- content/p1.md --
---
title: "P1"
weight: 1
---
-- layouts/_default/single.html --
Single: {{ .Title }}|{{ .Content }}|
-- layouts/index.html --
Pages: {{ range .RegularPages }}{{ .RelPermalink }}|{{ end }}$
`

	b := TestRunning(t, files)
	b.AssertFileContent("public/index.html", "Pages: /p1/|$")
	b.AddFiles("content/p2.md", ``).Build()
	b.AssertFileContent("public/index.html", "Pages: /p1/|/p2/|$")
}

func TestRebuildAddPageWithSpaceListPagesInHome(t *testing.T) {
	files := `
-- hugo.toml --
baseURL = "https://example.com"
disableLiveReload = true
-- content/asection/s1.md --
-- content/p1.md --
---
title: "P1"
weight: 1
---
-- layouts/_default/single.html --
Single: {{ .Title }}|{{ .Content }}|
-- layouts/index.html --
Pages: {{ range .RegularPages }}{{ .RelPermalink }}|{{ end }}$
`

	b := TestRunning(t, files)
	b.AssertFileContent("public/index.html", "Pages: /p1/|$")
	b.AddFiles("content/test test/index.md", ``).Build()
	b.AssertFileContent("public/index.html", "Pages: /p1/|/test-test/|$")
}

func TestRebuildScopedToOutputFormat(t *testing.T) {
	files := `
-- hugo.toml --
baseURL = "https://example.com"
disableKinds = ["term", "taxonomy", "sitemap", "robotstxt", "404"]
disableLiveReload = true
-- content/p1.md --
---
title: "P1"
outputs: ["html", "json"]
---
P1 Content.

{{< myshort >}}
-- layouts/_default/single.html --
Single HTML: {{ .Title }}|{{ .Content }}|
-- layouts/_default/single.json --
Single JSON: {{ .Title }}|{{ .Content }}|
-- layouts/shortcodes/myshort.html --
My short.
`
	b := Test(t, files, TestOptRunning())
	b.AssertRenderCountPage(3)
	b.AssertRenderCountContent(1)
	b.AssertFileContent("public/p1/index.html", "Single HTML: P1|<p>P1 Content.</p>\n")
	b.AssertFileContent("public/p1/index.json", "Single JSON: P1|<p>P1 Content.</p>\n")
	b.EditFileReplaceAll("layouts/_default/single.html", "Single HTML", "Single HTML Edited").Build()
	b.AssertFileContent("public/p1/index.html", "Single HTML Edited: P1|<p>P1 Content.</p>\n")
	b.AssertRenderCountPage(1)

	// Edit shortcode. Note that this is reused across all output formats.
	b.EditFileReplaceAll("layouts/shortcodes/myshort.html", "My short", "My short edited").Build()
	b.AssertFileContent("public/p1/index.html", "My short edited")
	b.AssertFileContent("public/p1/index.json", "My short edited")
	b.AssertRenderCountPage(3) // rss (uses .Content) + 2 single pages.
}

func TestRebuildBaseof(t *testing.T) {
	files := `
-- hugo.toml --
title = "Hugo Site"
baseURL = "https://example.com"
disableKinds = ["term", "taxonomy"]
disableLiveReload = true
-- layouts/_default/baseof.html --
Baseof: {{ .Title }}|
{{ block "main" . }}default{{ end }}
-- layouts/index.html --
{{ define "main" }}
Home: {{ .Title }}|{{ .Content }}|
{{ end }}
`
	testRebuildBothWatchingAndRunning(t, files, func(b *IntegrationTestBuilder) {
		b.AssertFileContent("public/index.html", "Baseof: Hugo Site|", "Home: Hugo Site||")
		b.EditFileReplaceFunc("layouts/_default/baseof.html", func(s string) string {
			return strings.Replace(s, "Baseof", "Baseof Edited", 1)
		}).Build()
		b.AssertFileContent("public/index.html", "Baseof Edited: Hugo Site|", "Home: Hugo Site||")
	})
}

func TestRebuildSingleWithBaseof(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
title = "Hugo Site"
baseURL = "https://example.com"
disableKinds = ["term", "taxonomy"]
disableLiveReload = true
-- content/p1.md --
---
title: "P1"
---
P1 Content.
-- layouts/_default/baseof.html --
Baseof: {{ .Title }}|
{{ block "main" . }}default{{ end }}
-- layouts/index.html --
Home.
-- layouts/_default/single.html --
{{ define "main" }}
Single: {{ .Title }}|{{ .Content }}|
{{ end }}
`
	b := Test(t, files, TestOptRunning())
	b.AssertFileContent("public/p1/index.html", "Baseof: P1|\n\nSingle: P1|<p>P1 Content.</p>\n|")
	b.EditFileReplaceFunc("layouts/_default/single.html", func(s string) string {
		return strings.Replace(s, "Single", "Single Edited", 1)
	}).Build()
	b.AssertFileContent("public/p1/index.html", "Baseof: P1|\n\nSingle Edited: P1|<p>P1 Content.</p>\n|")
}

func TestRebuildFromString(t *testing.T) {
	files := `
-- hugo.toml --
baseURL = "https://example.com"
disableKinds = ["term", "taxonomy", "sitemap", "robotstxt", "404"]
disableLiveReload = true
-- content/p1.md --
---
title: "P1"
layout: "l1"
---
P1 Content.
-- content/p2.md --
---
title: "P2"
layout: "l2"
---
P2 Content.
-- assets/mytext.txt --
My Text
-- layouts/_default/l1.html --
{{ $r := partial "get-resource.html" . }}
L1: {{ .Title }}|{{ .Content }}|R: {{ $r.Content }}|
-- layouts/_default/l2.html --
L2.
-- layouts/partials/get-resource.html --
{{ $mytext := resources.Get "mytext.txt" }}
{{ $txt := printf "Text: %s" $mytext.Content }}
{{ $r := resources.FromString "r.txt"  $txt }}
{{ return $r }}

`
	b := TestRunning(t, files)

	b.AssertFileContent("public/p1/index.html", "L1: P1|<p>P1 Content.</p>\n|R: Text: My Text|")

	b.EditFileReplaceAll("assets/mytext.txt", "My Text", "My Text Edited").Build()

	b.AssertFileContent("public/p1/index.html", "L1: P1|<p>P1 Content.</p>\n|R: Text: My Text Edited|")

	b.AssertRenderCountPage(1)
}

func TestRebuildDeeplyNestedLink(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com/"
disableKinds = ["term", "taxonomy", "sitemap", "robotstxt", "404"]
disableLiveReload = true
-- content/s/p1.md --
---
title: "P1"
---
-- content/s/p2.md --
---
title: "P2"
---
-- content/s/p3.md --
---
title: "P3"
---
-- content/s/p4.md --
---
title: "P4"
---
-- content/s/p5.md --
---
title: "P5"
---
-- content/s/p6.md --
---
title: "P6"
---
-- content/s/p7.md --
---
title: "P7"
---
-- layouts/_default/list.html --
List.
-- layouts/_default/single.html --
Single.
-- layouts/_default/single.html --
Next: {{ with  .PrevInSection }}{{ .Title }}{{ end }}|
Prev: {{ with  .NextInSection }}{{ .Title }}{{ end }}|


`

	b := TestRunning(t, files)

	b.AssertFileContent("public/s/p1/index.html", "Next: P2|")
	b.EditFileReplaceAll("content/s/p7.md", "P7", "P7 Edited").Build()
	b.AssertFileContent("public/s/p6/index.html", "Next: P7 Edited|")
}

func TestRebuildVariations(t *testing.T) {
	// t.Parallel() not supported, see https://github.com/fortytw2/leaktest/issues/4
	// This leaktest seems to be a little bit shaky on Travis.
	if !htesting.IsCI() {
		defer leaktest.CheckTimeout(t, 10*time.Second)()
	}

	files := `
-- hugo.toml --
baseURL = "https://example.com"
disableKinds = ["term", "taxonomy"]
disableLiveReload = true
defaultContentLanguage = "nn"
[pagination]
pagerSize = 20
[security]
enableInlineShortcodes = true
[languages]
[languages.en]
weight = 1
[languages.nn]
weight = 2
-- content/mysect/p1/index.md --
---
title: "P1"
---
P1 Content.
{{< include "mysect/p2" >}}
§§§go { page="mysect/p3" }
hello
§§§

{{< foo.inline >}}Foo{{< /foo.inline >}}
-- content/mysect/p2/index.md --
---
title: "P2"
---
P2 Content.
-- content/mysect/p3/index.md --
---
title: "P3"
---
P3 Content.
-- content/mysect/sub/_index.md --
-- content/mysect/sub/p4/index.md --
---
title: "P4"
---
P4 Content.
-- content/mysect/sub/p5/index.md --
---
title: "P5"
lastMod: 2019-03-02
---
P5 Content.
-- content/myothersect/_index.md --
---
cascade:
- _target:
  cascadeparam: "cascadevalue"
---
-- content/myothersect/sub/_index.md --
-- content/myothersect/sub/p6/index.md --
---
title: "P6"
---
P6 Content.
-- content/translations/p7.en.md --
---
title: "P7 EN"
---
P7 EN Content.
-- content/translations/p7.nn.md --
---
title: "P7 NN"
---
P7 NN Content.
-- layouts/index.html --
Home: {{ .Title }}|{{ .Content }}|
RegularPages: {{ range .RegularPages }}{{ .RelPermalink }}|{{ end }}$
Len RegularPagesRecursive: {{ len .RegularPagesRecursive }}
Site.Lastmod: {{ .Site.Lastmod.Format "2006-01-02" }}|
Paginate: {{ range (.Paginate .Site.RegularPages).Pages }}{{ .RelPermalink }}|{{ .Title }}|{{ end }}$
-- layouts/_default/single.html --
Single: {{ .Title }}|{{ .Content }}|
Single Partial Cached: {{ partialCached "pcached" . }}|
Page.Lastmod: {{ .Lastmod.Format "2006-01-02" }}|
Cascade param: {{ .Params.cascadeparam }}|
-- layouts/_default/list.html --
List: {{ .Title }}|{{ .Content }}|
RegularPages: {{ range .RegularPages }}{{ .Title }}|{{ end }}$
Len RegularPagesRecursive: {{ len .RegularPagesRecursive }}
RegularPagesRecursive: {{ range .RegularPagesRecursive }}{{ .RelPermalink }}|{{ end }}$
List Partial P1: {{ partial "p1" . }}|
Page.Lastmod: {{ .Lastmod.Format "2006-01-02" }}|
Cascade param: {{ .Params.cascadeparam }}|
-- layouts/partials/p1.html --
Partial P1.
-- layouts/partials/pcached.html --
Partial Pcached.
-- layouts/shortcodes/include.html --
{{ $p := site.GetPage (.Get 0)}}
{{ with $p }}
Shortcode Include: {{ .Title }}|
{{ end }}
Shortcode .Page.Title: {{ .Page.Title }}|
Shortcode Partial P1: {{ partial "p1" . }}|
-- layouts/_default/_markup/render-codeblock.html --
{{ $p := site.GetPage (.Attributes.page)}}
{{ with $p }}
Codeblock Include: {{ .Title }}|
{{ end }}



`

	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			Running:     true,
			BuildCfg: BuildCfg{
				testCounters: &buildCounters{},
			},
			// Verbose:     true,
			// LogLevel: logg.LevelTrace,
		},
	).Build()

	// When running the server, this is done on shutdown.
	// Do this here to satisfy the leak detector above.
	defer func() {
		b.Assert(b.H.Close(), qt.IsNil)
	}()

	contentRenderCount := b.counters.contentRenderCounter.Load()
	pageRenderCount := b.counters.pageRenderCounter.Load()

	b.Assert(contentRenderCount > 0, qt.IsTrue)
	b.Assert(pageRenderCount > 0, qt.IsTrue)

	// Test cases:
	// - Edit content file direct
	// - Edit content file transitive shortcode
	// - Edit content file transitive render hook
	// - Rename one language version of a content file
	// - Delete content file, check site.RegularPages and section.RegularPagesRecursive (length)
	// - Add content file (see above).
	// - Edit shortcode
	// - Edit inline shortcode
	// - Edit render hook
	// - Edit partial used in template
	// - Edit partial used in shortcode
	// - Edit partial cached.
	// - Edit lastMod date in content file, check site.Lastmod.
	editFile := func(filename string, replacementFunc func(s string) string) {
		b.EditFileReplaceFunc(filename, replacementFunc).Build()
		b.Assert(b.counters.contentRenderCounter.Load() < contentRenderCount, qt.IsTrue, qt.Commentf("count %d < %d", b.counters.contentRenderCounter.Load(), contentRenderCount))
		b.Assert(b.counters.pageRenderCounter.Load() < pageRenderCount, qt.IsTrue, qt.Commentf("count %d < %d", b.counters.pageRenderCounter.Load(), pageRenderCount))
	}

	b.AssertFileContent("public/index.html", "RegularPages: $", "Len RegularPagesRecursive: 7", "Site.Lastmod: 2019-03-02")

	b.AssertFileContent("public/mysect/p1/index.html",
		"Single: P1|<p>P1 Content.",
		"Shortcode Include: P2|",
		"Codeblock Include: P3|")

	editFile("content/mysect/p1/index.md", func(s string) string {
		return strings.ReplaceAll(s, "P1", "P1 Edited")
	})

	b.AssertFileContent("public/mysect/p1/index.html", "Single: P1 Edited|<p>P1 Edited Content.")
	b.AssertFileContent("public/index.html", "RegularPages: $", "Len RegularPagesRecursive: 7", "Paginate: /mysect/sub/p5/|P5|/mysect/p1/|P1 Edited")
	b.AssertFileContent("public/mysect/index.html", "RegularPages: P1 Edited|P2|P3|$", "Len RegularPagesRecursive: 5")

	// p2 is included in p1 via shortcode.
	editFile("content/mysect/p2/index.md", func(s string) string {
		return strings.ReplaceAll(s, "P2", "P2 Edited")
	})

	b.AssertFileContent("public/mysect/p1/index.html", "Shortcode Include: P2 Edited|")

	// p3 is included in p1 via codeblock hook.
	editFile("content/mysect/p3/index.md", func(s string) string {
		return strings.ReplaceAll(s, "P3", "P3 Edited")
	})

	b.AssertFileContent("public/mysect/p1/index.html", "Codeblock Include: P3 Edited|")

	// Remove a content file in a nested section.
	b.RemoveFiles("content/mysect/sub/p4/index.md").Build()
	b.AssertFileContent("public/mysect/index.html", "RegularPages: P1 Edited|P2 Edited|P3 Edited", "Len RegularPagesRecursive: 4")
	b.AssertFileContent("public/mysect/sub/index.html", "RegularPages: P5|$", "RegularPagesRecursive: 1")

	// Rename one of the translations.
	b.AssertFileContent("public/translations/index.html", "RegularPagesRecursive: /translations/p7/")
	b.AssertFileContent("public/en/translations/index.html", "RegularPagesRecursive: /en/translations/p7/")
	b.RenameFile("content/translations/p7.nn.md", "content/translations/p7rename.nn.md").Build()
	b.AssertFileContent("public/translations/index.html", "RegularPagesRecursive: /translations/p7rename/")
	b.AssertFileContent("public/en/translations/index.html", "RegularPagesRecursive: /en/translations/p7/")

	// Edit shortcode
	editFile("layouts/shortcodes/include.html", func(s string) string {
		return s + "\nShortcode Include Edited."
	})
	b.AssertFileContent("public/mysect/p1/index.html", "Shortcode Include Edited.")

	// Edit render hook
	editFile("layouts/_default/_markup/render-codeblock.html", func(s string) string {
		return s + "\nCodeblock Include Edited."
	})
	b.AssertFileContent("public/mysect/p1/index.html", "Codeblock Include Edited.")

	// Edit partial p1
	editFile("layouts/partials/p1.html", func(s string) string {
		return strings.Replace(s, "Partial P1", "Partial P1 Edited", 1)
	})
	b.AssertFileContent("public/mysect/index.html", "List Partial P1: Partial P1 Edited.")
	b.AssertFileContent("public/mysect/p1/index.html", "Shortcode Partial P1: Partial P1 Edited.")

	// Edit partial cached.
	editFile("layouts/partials/pcached.html", func(s string) string {
		return strings.Replace(s, "Partial Pcached", "Partial Pcached Edited", 1)
	})
	b.AssertFileContent("public/mysect/p1/index.html", "Pcached Edited.")

	// Edit lastMod date in content file, check site.Lastmod.
	editFile("content/mysect/sub/p5/index.md", func(s string) string {
		return strings.Replace(s, "2019-03-02", "2020-03-10", 1)
	})
	b.AssertFileContent("public/index.html", "Site.Lastmod: 2020-03-10|")
	b.AssertFileContent("public/mysect/index.html", "Page.Lastmod: 2020-03-10|")

	// Adjust the date back a few days.
	editFile("content/mysect/sub/p5/index.md", func(s string) string {
		return strings.Replace(s, "2020-03-10", "2019-03-08", 1)
	})
	b.AssertFileContent("public/mysect/index.html", "Page.Lastmod: 2019-03-08|")
	b.AssertFileContent("public/index.html", "Site.Lastmod: 2019-03-08|")

	// Check cascade mods.
	b.AssertFileContent("public/myothersect/index.html", "Cascade param: cascadevalue|")
	b.AssertFileContent("public/myothersect/sub/index.html", "Cascade param: cascadevalue|")
	b.AssertFileContent("public/myothersect/sub/p6/index.html", "Cascade param: cascadevalue|")

	editFile("content/myothersect/_index.md", func(s string) string {
		return strings.Replace(s, "cascadevalue", "cascadevalue edited", 1)
	})
	b.AssertFileContent("public/myothersect/index.html", "Cascade param: cascadevalue edited|")
	b.AssertFileContent("public/myothersect/sub/p6/index.html", "Cascade param: cascadevalue edited|")

	// Repurpose the cascadeparam to set the title.
	editFile("content/myothersect/_index.md", func(s string) string {
		return strings.Replace(s, "cascadeparam:", "title:", 1)
	})
	b.AssertFileContent("public/myothersect/sub/index.html", "Cascade param: |", "List: cascadevalue edited|")

	// Revert it.
	editFile("content/myothersect/_index.md", func(s string) string {
		return strings.Replace(s, "title:", "cascadeparam:", 1)
	})
	b.AssertFileContent("public/myothersect/sub/index.html", "Cascade param: cascadevalue edited|", "List: |")
}

func TestRebuildVariationsJSNoneFingerprinted(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com/"
disableKinds = ["term", "taxonomy", "sitemap", "robotsTXT", "404", "rss"]
disableLiveReload = true
-- content/p1/index.md --
---
title: "P1"
---
P1.
-- content/p2/index.md --
---
title: "P2"
---
P2.
-- content/p3/index.md --
---
title: "P3"
---
P3.
-- content/p4/index.md --
---
title: "P4"
---
P4.
-- assets/main.css --
body {
	background: red;
}
-- layouts/default/list.html --
List.
-- layouts/_default/single.html --
Single.
{{ $css := resources.Get "main.css" | minify }}
RelPermalink: {{ $css.RelPermalink }}|

`

	b := TestRunning(t, files)

	b.AssertFileContent("public/p1/index.html", "RelPermalink: /main.min.css|")
	b.AssertFileContent("public/main.min.css", "body{background:red}")

	b.EditFileReplaceAll("assets/main.css", "red", "blue")
	b.RemoveFiles("content/p2/index.md")
	b.RemoveFiles("content/p3/index.md")
	b.Build()

	b.AssertFileContent("public/main.min.css", "body{background:blue}")
}

func TestRebuildVariationsJSInNestedCachedPartialFingerprinted(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com/"
disableKinds = ["term", "taxonomy", "sitemap", "robotsTXT", "404", "rss"]
disableLiveReload = true
-- content/p1/index.md --
---
title: "P1"
---
P1.
-- content/p2/index.md --
---
title: "P2"
---
P2.
-- content/p3/index.md --
---
title: "P3"
---
P3.
-- content/p4/index.md --
---
title: "P4"
---
P4.
-- assets/js/main.js --
console.log("Hello");
-- layouts/_default/list.html --
List. {{ partial "head.html" . }}$
-- layouts/_default/single.html --
Single. {{ partial "head.html" . }}$
-- layouts/partials/head.html --
{{ partialCached "js.html" . }}$
-- layouts/partials/js.html --
{{ $js := resources.Get "js/main.js" | js.Build | fingerprint }}
RelPermalink: {{ $js.RelPermalink }}|
`

	b := TestRunning(t, files)

	b.AssertFileContent("public/p1/index.html", "/js/main.712a50b59d0f0dedb4e3606eaa3860b1f1a5305f6c42da30a2985e473ba314eb.js")
	b.AssertFileContent("public/index.html", "/js/main.712a50b59d0f0dedb4e3606eaa3860b1f1a5305f6c42da30a2985e473ba314eb.js")

	b.EditFileReplaceAll("assets/js/main.js", "Hello", "Hello is Edited").Build()

	for i := 1; i < 5; i++ {
		b.AssertFileContent(fmt.Sprintf("public/p%d/index.html", i), "/js/main.6535698cec9a21875f40ae03e96f30c4bee41a01e979224761e270b9034b2424.js")
	}

	b.AssertFileContent("public/index.html", "/js/main.6535698cec9a21875f40ae03e96f30c4bee41a01e979224761e270b9034b2424.js")
}

func TestRebuildVariationsJSInNestedPartialFingerprintedInBase(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com/"
disableKinds = ["term", "taxonomy", "sitemap", "robotsTXT", "404", "rss"]
disableLiveReload = true
-- assets/js/main.js --
console.log("Hello");
-- layouts/_default/baseof.html --
Base. {{ partial "common/head.html" . }}$
{{ block "main" . }}default{{ end }}
-- layouts/_default/list.html --
{{ define "main" }}main{{ end }}
-- layouts/partials/common/head.html --
{{ partial "myfiles/js.html" . }}$
-- layouts/partials/myfiles/js.html --
{{ $js := resources.Get "js/main.js" | js.Build | fingerprint }}
RelPermalink: {{ $js.RelPermalink }}|
`

	b := TestRunning(t, files)

	b.AssertFileContent("public/index.html", "/js/main.712a50b59d0f0dedb4e3606eaa3860b1f1a5305f6c42da30a2985e473ba314eb.js")

	b.EditFileReplaceAll("assets/js/main.js", "Hello", "Hello is Edited").Build()

	b.AssertFileContent("public/index.html", "/js/main.6535698cec9a21875f40ae03e96f30c4bee41a01e979224761e270b9034b2424.js")
}

func TestRebuildVariationsJSBundled(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com"
disableKinds = ["term", "taxonomy", "sitemap", "robotsTXT", "404", "rss"]
disableLiveReload = true
-- content/_index.md --
---
title: "Home"
---
-- content/p1.md --
---
title: "P1"
layout: "main"
---
-- content/p2.md --
---
title: "P2"
---
{{< jsfingerprinted >}}
-- content/p3.md --
---
title: "P3"
layout: "plain"
---
{{< jsfingerprinted >}}
-- content/main.js --
console.log("Hello");
-- content/foo.js --
console.log("Foo");
-- layouts/index.html --
Home.
{{ $js := site.Home.Resources.Get "main.js"  }}
{{ with $js }}
<script src="{{ .RelPermalink }}"></script>
{{ end }}
-- layouts/_default/single.html --
Single. Deliberately no .Content in here.
-- layouts/_default/plain.html --
Content: {{ .Content }}|
-- layouts/_default/main.html --
{{ $js := site.Home.Resources.Get "main.js"  }}
{{ with $js }}
<script>
{{ .Content }}
</script>
{{ end }}
-- layouts/shortcodes/jsfingerprinted.html --
{{ $js := site.Home.Resources.Get "foo.js" | fingerprint  }}
<script src="{{ $js.RelPermalink }}"></script>
`

	testCounters := &buildCounters{}

	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			Running:     true,
			// LogLevel:    logg.LevelTrace,
			// Verbose:     true,
			BuildCfg: BuildCfg{
				testCounters: testCounters,
			},
		},
	).Build()

	b.AssertFileContent("public/index.html", `<script src="/main.js"></script>`)
	b.AssertFileContent("public/p1/index.html", "<script>\n\"console.log(\\\"Hello\\\");\"\n</script>")
	b.AssertFileContent("public/p2/index.html", "Single. Deliberately no .Content in here.")
	b.AssertFileContent("public/p3/index.html", "foo.57b4465b908531b43d4e4680ab1063d856b475cb1ae81ad43e0064ecf607bec1.js")
	b.AssertRenderCountPage(4)

	// Edit JS file.
	b.EditFileReplaceFunc("content/main.js", func(s string) string {
		return strings.Replace(s, "Hello", "Hello is Edited", 1)
	}).Build()

	b.AssertFileContent("public/p1/index.html", "<script>\n\"console.log(\\\"Hello is Edited\\\");\"\n</script>")
	// The p1 (the one inlining the JS) should be rebuilt.
	b.AssertRenderCountPage(2)
	// But not the content file.
	b.AssertRenderCountContent(0)

	// This is included with RelPermalink in a shortcode used in p3, but it's fingerprinted
	// so we need to rebuild on change.
	b.EditFileReplaceFunc("content/foo.js", func(s string) string {
		return strings.Replace(s, "Foo", "Foo Edited", 1)
	}).Build()

	// Verify that the hash has changed.
	b.AssertFileContent("public/p3/index.html", "foo.3a332a088521231e5fc9bd22f15e0ccf507faa7b373fbff22959005b9a80481c.js")

	b.AssertRenderCountPage(1)
	b.AssertRenderCountContent(1)
}

func TestRebuildEditData(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableLiveReload = true
[security]
enableInlineShortcodes=true
-- data/mydata.yaml --
foo: bar
-- content/_index.md --
---
title: "Home"
---
{{< data "mydata.foo" >}}}
-- content/p1.md --
---
title: "P1"
---

Foo inline: {{< foo.inline >}}{{ site.Data.mydata.foo }}|{{< /foo.inline >}}
-- layouts/shortcodes/data.html --
{{ $path := split (.Get 0) "." }}
{{ $data := index site.Data $path }}
Foo: {{ $data }}|
-- layouts/index.html --
Content: {{ .Content }}|
-- layouts/_default/single.html --
Single: {{ .Content }}|
`
	b := TestRunning(t, files)

	b.AssertFileContent("public/index.html", "Foo: bar|")
	b.AssertFileContent("public/p1/index.html", "Foo inline: bar|")
	b.EditFileReplaceFunc("data/mydata.yaml", func(s string) string {
		return strings.Replace(s, "bar", "bar edited", 1)
	}).Build()
	b.AssertFileContent("public/index.html", "Foo: bar edited|")
	b.AssertFileContent("public/p1/index.html", "Foo inline: bar edited|")
}

func TestRebuildEditHomeContent(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com"
disableLiveReload = true
-- content/_index.md --
---
title: "Home"
---
Home.
-- layouts/index.html --
Content: {{ .Content }}
`
	b := TestRunning(t, files)

	b.AssertFileContent("public/index.html", "Content: <p>Home.</p>")
	b.EditFileReplaceAll("content/_index.md", "Home.", "Home").Build()
	b.AssertFileContent("public/index.html", "Content: <p>Home</p>")
}

func TestRebuildVariationsAssetsJSImport(t *testing.T) {
	t.Parallel()
	files := `
-- hugo.toml --
baseURL = "https://example.com"
disableKinds = ["term", "taxonomy"]
disableLiveReload = true
-- layouts/index.html --
Home. {{ now }}
{{ with (resources.Get "js/main.js" | js.Build | fingerprint) }}
<script>{{ .Content | safeJS }}</script>
{{ end }}
-- assets/js/lib/foo.js --
export function foo() {
	console.log("Foo");
}
-- assets/js/main.js --
import { foo } from "./lib/foo.js";
console.log("Hello");
foo();
`

	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			Running:     true,
			// LogLevel:    logg.LevelTrace,
			NeedsOsFS: true,
		},
	).Build()

	b.AssertFileContent("public/index.html", "Home.", "Hello", "Foo")
	// Edit the imported file.
	b.EditFileReplaceAll("assets/js/lib/foo.js", "Foo", "Foo Edited").Build()
	b.AssertFileContent("public/index.html", "Home.", "Hello", "Foo Edited")
}

func TestRebuildVariationsAssetsPostCSSImport(t *testing.T) {
	if !htesting.IsCI() {
		t.Skip("skip CI only")
	}

	files := `
-- hugo.toml --
baseURL = "https://example.com"
disableKinds = ["term", "taxonomy", "sitemap", "rss"]
disableLiveReload = true
-- assets/css/lib/foo.css --
body {
	background: red;
}
-- assets/css/main.css --
@import "lib/foo.css";
-- package.json --
{
	"devDependencies": {
		"postcss-cli": "^9.0.1"
	}
}
-- content/p1.md --
---
title: "P1"
---
-- content/p2.md --
---
title: "P2"
layout: "foo"
---
{{< fingerprinted >}}
-- content/p3.md --
---
title: "P3"
layout: "foo"
---
{{< notfingerprinted >}}
-- layouts/shortcodes/fingerprinted.html --
Fingerprinted.
{{ $opts := dict "inlineImports" true "noMap" true }}
{{ with (resources.Get "css/main.css" | postCSS $opts | fingerprint) }}
<style src="{{ .RelPermalink }}"></style>
{{ end }}
-- layouts/shortcodes/notfingerprinted.html --
Fingerprinted.
{{ $opts := dict "inlineImports" true "noMap" true }}
{{ with (resources.Get "css/main.css" | postCSS $opts) }}
<style src="{{ .RelPermalink }}"></style>
{{ end }}
-- layouts/index.html --
Home.
{{ $opts := dict "inlineImports" true "noMap" true }}
{{ with (resources.Get "css/main.css" | postCSS $opts) }}
<style>{{ .Content | safeCSS }}</style>
{{ end }}
-- layouts/_default/foo.html --
Foo.
{{ .Title }}|{{ .Content }}|
-- layouts/_default/single.html --
Single.
{{ $opts := dict "inlineImports" true "noMap" true }}
{{ with (resources.Get "css/main.css" | postCSS $opts) }}
<style src="{{ .RelPermalink }}"></style>
{{ end }}
`

	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:               t,
			TxtarString:     files,
			Running:         true,
			NeedsOsFS:       true,
			NeedsNpmInstall: true,
			// LogLevel:        logg.LevelDebug,
		},
	).Build()

	b.AssertFileContent("public/index.html", "Home.", "<style>body {\n\tbackground: red;\n}</style>")
	b.AssertFileContent("public/p1/index.html", "Single.", "/css/main.css")
	b.AssertRenderCountPage(4)

	// Edit the imported file.
	b.EditFileReplaceFunc("assets/css/lib/foo.css", func(s string) string {
		return strings.Replace(s, "red", "blue", 1)
	}).Build()

	b.AssertRenderCountPage(3)

	b.AssertFileContent("public/index.html", "Home.", "<style>body {\n\tbackground: blue;\n}</style>")
}

func TestRebuildI18n(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com"
disableLiveReload = true
-- i18n/en.toml --
hello = "Hello"
-- layouts/index.html --
Hello: {{ i18n "hello" }}
`

	b := TestRunning(t, files)

	b.AssertFileContent("public/index.html", "Hello: Hello")

	b.EditFileReplaceAll("i18n/en.toml", "Hello", "Hugo").Build()

	b.AssertFileContent("public/index.html", "Hello: Hugo")
}

func TestRebuildEditContentNonDefaultLanguage(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com"
disableLiveReload = true
defaultContentLanguage = "en"
defaultContentLanguageInSubdir = true
[languages]
[languages.en]
weight = 1
[languages.nn]
weight = 2
-- content/p1/index.en.md --
---
title: "P1 en"
---
P1 en.
-- content/p1/b.en.md --
---
title: "B en"
---
B en.
-- content/p1/f1.en.txt --
F1 en
-- content/p1/index.nn.md --
---
title: "P1 nn"
---
P1 nn.
-- content/p1/b.nn.md --
---
title: "B nn"
---
B nn.
-- content/p1/f1.nn.txt --
F1 nn
-- layouts/_default/single.html --
Single: {{ .Title }}|{{ .Content }}|Bundled File: {{ with .Resources.GetMatch "f1.*" }}{{ .Content }}{{ end }}|Bundled Page: {{ with .Resources.GetMatch "b.*" }}{{ .Content }}{{ end }}|
`

	b := TestRunning(t, files)

	b.AssertFileContent("public/nn/p1/index.html", "Single: P1 nn|<p>P1 nn.</p>", "F1 nn|")
	b.EditFileReplaceAll("content/p1/index.nn.md", "P1 nn.", "P1 nn edit.").Build()
	b.AssertFileContent("public/nn/p1/index.html", "Single: P1 nn|<p>P1 nn edit.</p>\n|")
	b.EditFileReplaceAll("content/p1/f1.nn.txt", "F1 nn", "F1 nn edit.").Build()
	b.AssertFileContent("public/nn/p1/index.html", "Bundled File: F1 nn edit.")
	b.EditFileReplaceAll("content/p1/b.nn.md", "B nn.", "B nn edit.").Build()
	b.AssertFileContent("public/nn/p1/index.html", "B nn edit.")
}

func TestRebuildEditContentNonDefaultLanguageDifferentBundles(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com"
disableLiveReload = true
defaultContentLanguage = "en"
defaultContentLanguageInSubdir = true
[languages]
[languages.en]
weight = 1
contentDir = "content/en"
[languages.nn]
weight = 2
contentDir = "content/nn"
-- content/en/p1en/index.md --
---
title: "P1 en"
---
-- content/nn/p1nn/index.md --
---
title: "P1 nn"
---
P1 nn.
-- layouts/_default/single.html --
Single: {{ .Title }}|{{ .Content }}|
`

	b := TestRunning(t, files)

	b.AssertFileContent("public/nn/p1nn/index.html", "Single: P1 nn|<p>P1 nn.</p>")
	b.EditFileReplaceAll("content/nn/p1nn/index.md", "P1 nn.", "P1 nn edit.").Build()
	b.AssertFileContent("public/nn/p1nn/index.html", "Single: P1 nn|<p>P1 nn edit.</p>\n|")
	b.AssertFileContent("public/nn/p1nn/index.html", "P1 nn edit.")
}

func TestRebuildVariationsAssetsSassImport(t *testing.T) {
	if !htesting.IsCI() {
		t.Skip("skip CI only")
	}

	filesTemplate := `
-- hugo.toml --
baseURL = "https://example.com"
disableKinds = ["term", "taxonomy"]
disableLiveReload = true
-- assets/css/lib/foo.scss --
body {
	background: red;
}
-- assets/css/main.scss --
@import "lib/foo";
-- layouts/index.html --
Home.
{{ $opts := dict "transpiler" "TRANSPILER" }}
{{ with (resources.Get "css/main.scss" | toCSS $opts) }}
<style>{{ .Content | safeCSS }}</style>
{{ end }}
`

	runTest := func(transpiler string) {
		t.Run(transpiler, func(t *testing.T) {
			files := strings.Replace(filesTemplate, "TRANSPILER", transpiler, 1)
			b := NewIntegrationTestBuilder(
				IntegrationTestConfig{
					T:           t,
					TxtarString: files,
					Running:     true,
					NeedsOsFS:   true,
				},
			).Build()

			b.AssertFileContent("public/index.html", "Home.", "background: red")

			// Edit the imported file.
			b.EditFileReplaceFunc("assets/css/lib/foo.scss", func(s string) string {
				return strings.Replace(s, "red", "blue", 1)
			}).Build()

			b.AssertFileContent("public/index.html", "Home.", "background: blue")
		})
	}

	if scss.Supports() {
		runTest("libsass")
	}

	if dartsass.Supports() {
		runTest("dartsass")
	}
}

func benchmarkFilesEdit(count int) string {
	files := `
-- hugo.toml --
baseURL = "https://example.com"
disableKinds = ["term", "taxonomy"]
disableLiveReload = true
-- layouts/_default/single.html --
Single: {{ .Title }}|{{ .Content }}|
-- layouts/_default/list.html --
List: {{ .Title }}|{{ .Content }}|
-- content/mysect/_index.md --
---
title: "My Sect"
---
	`

	contentTemplate := `
---
title: "P%d"
---
P%d Content.
`

	for i := 0; i < count; i++ {
		files += fmt.Sprintf("-- content/mysect/p%d/index.md --\n%s", i, fmt.Sprintf(contentTemplate, i, i))
	}

	return files
}

func BenchmarkRebuildContentFileChange(b *testing.B) {
	files := benchmarkFilesEdit(500)

	cfg := IntegrationTestConfig{
		T:           b,
		TxtarString: files,
		Running:     true,
		// Verbose:     true,
		// LogLevel: logg.LevelInfo,
	}
	builders := make([]*IntegrationTestBuilder, b.N)

	for i := range builders {
		builders[i] = NewIntegrationTestBuilder(cfg)
		builders[i].Build()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bb := builders[i]
		bb.EditFileReplaceFunc("content/mysect/p123/index.md", func(s string) string {
			return s + "... Edited"
		}).Build()
		// fmt.Println(bb.LogString())
	}
}

func TestRebuildConcat(t *testing.T) {
	files := `
-- hugo.toml --
baseURL = "https://example.com"
disableLiveReload = true
disableKinds = ["taxonomy", "term", "sitemap", "robotsTXT", "404", "rss"]
-- assets/a.css --
a
-- assets/b.css --
b
-- assets/c.css --
c
-- assets/common/c1.css --
c1
-- assets/common/c2.css --
c2
-- layouts/index.html --
{{ $a := resources.Get "a.css" }}
{{ $b := resources.Get "b.css" }}
{{ $common := resources.Match "common/*.css" | resources.Concat "common.css" | minify }}
{{ $ab := slice $a $b $common | resources.Concat "ab.css" }}
all: {{ $ab.RelPermalink }}
`
	b := TestRunning(t, files)

	b.AssertFileContent("public/ab.css", "abc1c2")
	b.EditFileReplaceAll("assets/common/c2.css", "c2", "c2 edited").Build()
	b.AssertFileContent("public/ab.css", "abc1c2 edited")
	b.AddFiles("assets/common/c3.css", "c3").Build()
	b.AssertFileContent("public/ab.css", "abc1c2 editedc3")
}

func TestRebuildEditArchetypeFile(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com"
disableLiveReload = true
-- archetypes/default.md --
---
title: "Default"
---
`

	b := TestRunning(t, files)
	// Just make sure that it doesn't panic.
	b.EditFileReplaceAll("archetypes/default.md", "Default", "Default Edited").Build()
}

func TestRebuildEditMixedCaseTemplateFileIssue12165(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com"
disableLiveReload = true
-- layouts/partials/MyTemplate.html --
MyTemplate
-- layouts/index.html --
MyTemplate: {{ partial "MyTemplate.html" . }}|


`

	b := TestRunning(t, files)

	b.AssertFileContent("public/index.html", "MyTemplate: MyTemplate")

	b.EditFileReplaceAll("layouts/partials/MyTemplate.html", "MyTemplate", "MyTemplate Edited").Build()

	b.AssertFileContent("public/index.html", "MyTemplate: MyTemplate Edited")
}

func TestRebuildEditAsciidocContentFile(t *testing.T) {
	if !asciidocext.Supports() {
		t.Skip("skip asciidoc")
	}
	files := `
-- hugo.toml --
baseURL = "https://example.com"
disableLiveReload = true
disableKinds = ["taxonomy", "term", "sitemap", "robotsTXT", "404", "rss", "home", "section"]
[security]
[security.exec]
allow = ['^python$', '^rst2html.*', '^asciidoctor$']
-- content/posts/p1.adoc --
---
title: "P1"
---
P1 Content.
-- content/posts/p2.adoc --
---
title: "P2"
---
P2 Content.
-- layouts/_default/single.html --
Single: {{ .Title }}|{{ .Content }}|
`
	b := TestRunning(t, files)
	b.AssertFileContent("public/posts/p1/index.html",
		"Single: P1|<div class=\"paragraph\">\n<p>P1 Content.</p>\n</div>\n|")
	b.AssertRenderCountPage(2)
	b.AssertRenderCountContent(2)

	b.EditFileReplaceAll("content/posts/p1.adoc", "P1 Content.", "P1 Content Edited.").Build()

	b.AssertFileContent("public/posts/p1/index.html", "Single: P1|<div class=\"paragraph\">\n<p>P1 Content Edited.</p>\n</div>\n|")
	b.AssertRenderCountPage(1)
	b.AssertRenderCountContent(1)
}

func TestRebuildEditSingleListChangeUbuntuIssue12362(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['rss','section','sitemap','taxonomy','term']
disableLiveReload = true
-- layouts/_default/list.html --
{{ range .Pages }}{{ .Title }}|{{ end }}
-- layouts/_default/single.html --
{{ .Title }}
-- content/p1.md --
---
title: p1
---
`

	b := TestRunning(t, files)
	b.AssertFileContent("public/index.html", "p1|")

	b.AddFiles("content/p2.md", "---\ntitle: p2\n---").Build()
	b.AssertFileContent("public/index.html", "p1|p2|") // this test passes, which doesn't match reality
}

func TestRebuildHomeThenPageIssue12436(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com"
disableKinds = ['sitemap','taxonomy','term']
disableLiveReload = true
-- layouts/_default/list.html --
{{ .Content }}
-- layouts/_default/single.html --
{{ .Content }}
-- content/_index.md --
---
title: home
---
home-content|
-- content/p1/index.md --
---
title: p1
---
p1-content|
`

	b := TestRunning(t, files)

	b.AssertFileContent("public/index.html", "home-content|")
	b.AssertFileContent("public/p1/index.html", "p1-content|")
	b.AssertRenderCountPage(3)

	b.EditFileReplaceAll("content/_index.md", "home-content", "home-content-foo").Build()
	b.AssertFileContent("public/index.html", "home-content-foo")
	b.AssertRenderCountPage(2) // Home page rss + html

	b.EditFileReplaceAll("content/p1/index.md", "p1-content", "p1-content-foo").Build()
	b.AssertFileContent("public/p1/index.html", "p1-content-foo")
}
