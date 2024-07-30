// Copyright 2024 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pagesfromdata_test

import (
	"fmt"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/hugolib"
	"github.com/gohugoio/hugo/markup/asciidocext"
	"github.com/gohugoio/hugo/markup/pandoc"
	"github.com/gohugoio/hugo/markup/rst"
)

const filesPagesFromDataTempleBasic = `
-- hugo.toml --
disableKinds = ["taxonomy", "term", "rss", "sitemap"]
baseURL = "https://example.com"
disableLiveReload = true
-- assets/a/pixel.png --
iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==
-- assets/mydata.yaml --
p1: "p1"
draft: false
-- layouts/partials/get-value.html --
{{ $val := "p1" }}
{{ return $val }}
-- layouts/_default/baseof.html --
Baseof:
{{ block "main" . }}{{ end }}
-- layouts/_default/single.html --
{{ define "main" }}
Single: {{ .Title }}|{{ .Content }}|Params: {{ .Params.param1 }}|Path: {{ .Path }}|
Dates: Date: {{ .Date.Format "2006-01-02" }}|Lastmod: {{ .Lastmod.Format "2006-01-02" }}|PublishDate: {{ .PublishDate.Format "2006-01-02" }}|ExpiryDate: {{ .ExpiryDate.Format "2006-01-02" }}|
Len Resources: {{ .Resources | len }}
Resources: {{ range .Resources }}RelPermalink: {{ .RelPermalink }}|Name: {{ .Name }}|Title: {{ .Title }}|Params: {{ .Params }}|{{ end }}$
{{ with .Resources.Get "featured.png" }}
Featured Image: {{ .RelPermalink }}|{{ .Name }}|
{{ with .Resize "10x10" }}
Resized Featured Image: {{ .RelPermalink }}|{{ .Width }}|
{{ end}}
{{ end }}
{{ end }}
-- layouts/_default/list.html --
List: {{ .Title }}|{{ .Content }}|
RegularPagesRecursive: {{ range .RegularPagesRecursive }}{{ .Title }}:{{ .Path }}|{{ end }}$
Sections: {{ range .Sections }}{{ .Title }}:{{ .Path }}|{{ end }}$
-- content/docs/pfile.md --
---
title: "pfile"
date: 2023-03-01
---
Pfile Content
-- content/docs/_content.gotmpl --
{{ $pixel := resources.Get "a/pixel.png" }}
{{ $dataResource := resources.Get "mydata.yaml" }}
{{ $data := $dataResource | transform.Unmarshal }}
{{ $pd := $data.p1 }}
{{ $pp := partial "get-value.html" }}
{{ $title := printf "%s:%s" $pd $pp }}
{{ $date := "2023-03-01" | time.AsTime }}
{{ $dates := dict "date" $date }}
{{ $contentMarkdown := dict "value" "**Hello World**"  "mediaType" "text/markdown" }}
{{ $contentMarkdownDefault := dict "value" "**Hello World Default**" }}
{{ $contentHTML := dict "value" "<b>Hello World!</b> No **markdown** here." "mediaType" "text/html" }}
{{ $.AddPage  (dict "kind" "page" "path" "P1" "title" $title "dates" $dates "content" $contentMarkdown "params" (dict "param1" "param1v" ) ) }}
{{ $.AddPage  (dict "kind" "page" "path" "p2" "title" "p2title" "dates" $dates "content" $contentHTML ) }}
{{ $.AddPage  (dict "kind" "page" "path" "p3" "title" "p3title" "dates" $dates "content" $contentMarkdownDefault "draft" false ) }}
{{ $.AddPage  (dict "kind" "page" "path" "p4" "title" "p4title" "dates" $dates "content" $contentMarkdownDefault "draft" $data.draft ) }}
ADD_MORE_PLACEHOLDER


{{ $resourceContent := dict "value" $dataResource }}
{{ $.AddResource (dict "path" "p1/data1.yaml" "content" $resourceContent) }}
{{ $.AddResource (dict "path" "p1/mytext.txt" "content" (dict "value" "some text") "name" "textresource" "title" "My Text Resource" "params" (dict "param1" "param1v") )}}
{{ $.AddResource (dict "path" "p1/sub/mytex2.txt" "content" (dict "value" "some text") "title" "My Text Sub Resource" ) }}
{{ $.AddResource (dict "path" "P1/Sub/MyMixCaseText2.txt" "content" (dict "value" "some text") "title" "My Text Sub Mixed Case Path Resource" ) }}
{{ $.AddResource (dict "path" "p1/sub/data1.yaml" "content" $resourceContent "title" "Sub data") }}
{{ $resourceParams := dict "data2ParaM1" "data2Param1v" }}
{{ $.AddResource (dict "path" "p1/data2.yaml" "name" "data2.yaml" "title" "My data 2" "params" $resourceParams "content" $resourceContent) }}
{{ $.AddResource (dict "path" "p1/featuredimage.png" "name" "featured.png" "title" "My Featured Image" "params" $resourceParams "content" (dict "value" $pixel ))}}
`

func TestPagesFromGoTmplMisc(t *testing.T) {
	t.Parallel()
	b := hugolib.Test(t, filesPagesFromDataTempleBasic)
	b.AssertPublishDir(`
docs/p1/mytext.txt
docs/p1/sub/mytex2.tx
docs/p1/sub/mymixcasetext2.txt
	`)

	// Page from markdown file.
	b.AssertFileContent("public/docs/pfile/index.html", "Dates: Date: 2023-03-01|Lastmod: 2023-03-01|PublishDate: 2023-03-01|ExpiryDate: 0001-01-01|")
	// Pages from gotmpl.
	b.AssertFileContent("public/docs/p1/index.html",
		"Single: p1:p1|",
		"Path: /docs/p1|",
		"<strong>Hello World</strong>",
		"Params: param1v|",
		"Len Resources: 7",
		"RelPermalink: /mydata.yaml|Name: data1.yaml|Title: data1.yaml|Params: map[]|",
		"RelPermalink: /mydata.yaml|Name: data2.yaml|Title: My data 2|Params: map[data2param1:data2Param1v]|",
		"RelPermalink: /a/pixel.png|Name: featured.png|Title: My Featured Image|Params: map[data2param1:data2Param1v]|",
		"RelPermalink: /docs/p1/sub/mytex2.txt|Name: sub/mytex2.txt|",
		"RelPermalink: /docs/p1/sub/mymixcasetext2.txt|Name: sub/mymixcasetext2.txt|",
		"RelPermalink: /mydata.yaml|Name: sub/data1.yaml|Title: Sub data|Params: map[]|",
		"Featured Image: /a/pixel.png|featured.png|",
		"Resized Featured Image: /a/pixel_hu16809842526914527184.png|10|",
		// Resource from string
		"RelPermalink: /docs/p1/mytext.txt|Name: textresource|Title: My Text Resource|Params: map[param1:param1v]|",
		// Dates
		"Dates: Date: 2023-03-01|Lastmod: 2023-03-01|PublishDate: 2023-03-01|ExpiryDate: 0001-01-01|",
	)
	b.AssertFileContent("public/docs/p2/index.html", "Single: p2title|", "<b>Hello World!</b> No **markdown** here.")
	b.AssertFileContent("public/docs/p3/index.html", "<strong>Hello World Default</strong>")
}

func TestPagesFromGoTmplAsciidocAndSimilar(t *testing.T) {
	files := `
-- hugo.toml --
disableKinds = ["taxonomy", "term", "rss", "sitemap"]
baseURL = "https://example.com"
[security]
[security.exec]
allow = ['asciidoctor', 'pandoc','rst2html', 'python']
-- layouts/_default/single.html --
|Content: {{ .Content }}|Title: {{ .Title }}|Path: {{ .Path }}|
-- content/docs/_content.gotmpl --
{{ $.AddPage (dict "path" "asciidoc" "content" (dict "value" "Mark my words, #automation is essential#." "mediaType" "text/asciidoc" )) }}
{{ $.AddPage (dict "path" "pandoc" "content" (dict "value" "This ~~is deleted text.~~" "mediaType" "text/pandoc" )) }}
{{ $.AddPage (dict "path" "rst" "content" (dict "value" "This is *bold*." "mediaType" "text/rst" )) }}
{{ $.AddPage (dict "path" "org" "content" (dict "value" "the ability to use +strikethrough+ is a plus" "mediaType" "text/org" )) }}
{{ $.AddPage (dict "path" "nocontent" "title" "No Content" ) }}

	`

	b := hugolib.Test(t, files)

	if asciidocext.Supports() {
		b.AssertFileContent("public/docs/asciidoc/index.html",
			"Mark my words, <mark>automation is essential</mark>",
			"Path: /docs/asciidoc|",
		)
	}
	if pandoc.Supports() {
		b.AssertFileContent("public/docs/pandoc/index.html",
			"This <del>is deleted text.</del>",
			"Path: /docs/pandoc|",
		)
	}

	if rst.Supports() {
		b.AssertFileContent("public/docs/rst/index.html",
			"This is <em>bold</em>",
			"Path: /docs/rst|",
		)
	}

	b.AssertFileContent("public/docs/org/index.html",
		"the ability to use <del>strikethrough</del> is a plus",
		"Path: /docs/org|",
	)

	b.AssertFileContent("public/docs/nocontent/index.html", "|Content: |Title: No Content|Path: /docs/nocontent|")
}

func TestPagesFromGoTmplAddPageErrors(t *testing.T) {
	filesTemplate := `
-- hugo.toml --
disableKinds = ["taxonomy", "term", "rss", "sitemap"]
baseURL = "https://example.com"
-- content/docs/_content.gotmpl --
{{ $.AddPage  DICT }}
`

	t.Run("AddPage, missing Path", func(t *testing.T) {
		files := strings.ReplaceAll(filesTemplate, "DICT", `(dict "kind" "page" "title" "p1")`)
		b, err := hugolib.TestE(t, files)
		b.Assert(err, qt.IsNotNil)
		b.Assert(err.Error(), qt.Contains, "_content.gotmpl:1:4")
		b.Assert(err.Error(), qt.Contains, "error calling AddPage: path not set")
	})

	t.Run("AddPage, path starting with slash", func(t *testing.T) {
		files := strings.ReplaceAll(filesTemplate, "DICT", `(dict "kind" "page" "title" "p1" "path" "/foo")`)
		b, err := hugolib.TestE(t, files)
		b.Assert(err, qt.IsNotNil)
		b.Assert(err.Error(), qt.Contains, `path "/foo" must not start with a /`)
	})

	t.Run("AddPage, lang set", func(t *testing.T) {
		files := strings.ReplaceAll(filesTemplate, "DICT", `(dict "kind" "page" "path" "p1" "lang" "en")`)
		b, err := hugolib.TestE(t, files)
		b.Assert(err, qt.IsNotNil)
		b.Assert(err.Error(), qt.Contains, "_content.gotmpl:1:4")
		b.Assert(err.Error(), qt.Contains, "error calling AddPage: lang must not be set")
	})

	t.Run("Site methods not ready", func(t *testing.T) {
		filesTemplate := `
-- hugo.toml --
disableKinds = ["taxonomy", "term", "rss", "sitemap"]
baseURL = "https://example.com"
-- content/docs/_content.gotmpl --
{{ .Site.METHOD }}
`

		for _, method := range []string{"RegularPages", "Pages", "AllPages", "AllRegularPages", "Home", "Sections", "GetPage", "Menus", "MainSections", "Taxonomies"} {
			t.Run(method, func(t *testing.T) {
				files := strings.ReplaceAll(filesTemplate, "METHOD", method)
				b, err := hugolib.TestE(t, files)
				b.Assert(err, qt.IsNotNil)
				b.Assert(err.Error(), qt.Contains, fmt.Sprintf("error calling %s: this method cannot be called before the site is fully initialized", method))
			})
		}
	})
}

func TestPagesFromGoTmplAddResourceErrors(t *testing.T) {
	filesTemplate := `
-- hugo.toml --
disableKinds = ["taxonomy", "term", "rss", "sitemap"]
baseURL = "https://example.com"
-- content/docs/_content.gotmpl --
{{ $.AddResource  DICT }}
`

	t.Run("missing Path", func(t *testing.T) {
		files := strings.ReplaceAll(filesTemplate, "DICT", `(dict "name" "r1")`)
		b, err := hugolib.TestE(t, files)
		b.Assert(err, qt.IsNotNil)
		b.Assert(err.Error(), qt.Contains, "error calling AddResource: path not set")
	})
}

func TestPagesFromGoTmplEditGoTmpl(t *testing.T) {
	t.Parallel()
	b := hugolib.TestRunning(t, filesPagesFromDataTempleBasic)
	b.EditFileReplaceAll("content/docs/_content.gotmpl", `"title" "p2title"`, `"title" "p2titleedited"`).Build()
	b.AssertFileContent("public/docs/p2/index.html", "Single: p2titleedited|")
	b.AssertFileContent("public/docs/index.html", "p2titleedited")
}

func TestPagesFromGoTmplEditDataResource(t *testing.T) {
	t.Parallel()
	b := hugolib.TestRunning(t, filesPagesFromDataTempleBasic)
	b.AssertRenderCountPage(7)
	b.EditFileReplaceAll("assets/mydata.yaml", "p1: \"p1\"", "p1: \"p1edited\"").Build()
	b.AssertFileContent("public/docs/p1/index.html", "Single: p1edited:p1|")
	b.AssertFileContent("public/docs/index.html", "p1edited")
	b.AssertRenderCountPage(3)
}

func TestPagesFromGoTmplEditPartial(t *testing.T) {
	t.Parallel()
	b := hugolib.TestRunning(t, filesPagesFromDataTempleBasic)
	b.EditFileReplaceAll("layouts/partials/get-value.html", "p1", "p1edited").Build()
	b.AssertFileContent("public/docs/p1/index.html", "Single: p1:p1edited|")
	b.AssertFileContent("public/docs/index.html", "p1edited")
}

func TestPagesFromGoTmplRemovePage(t *testing.T) {
	t.Parallel()
	b := hugolib.TestRunning(t, filesPagesFromDataTempleBasic)
	b.EditFileReplaceAll("content/docs/_content.gotmpl", `{{ $.AddPage  (dict "kind" "page" "path" "p2" "title" "p2title" "dates" $dates "content" $contentHTML ) }}`, "").Build()
	b.AssertFileContent("public/index.html", "RegularPagesRecursive: p1:p1:/docs/p1|p3title:/docs/p3|p4title:/docs/p4|pfile:/docs/pfile|$")
}

func TestPagesFromGoTmplAddPage(t *testing.T) {
	t.Parallel()
	b := hugolib.TestRunning(t, filesPagesFromDataTempleBasic)
	b.EditFileReplaceAll("content/docs/_content.gotmpl", "ADD_MORE_PLACEHOLDER", `{{ $.AddPage  (dict "kind" "page" "path" "page_added" "title" "page_added_title" "dates" $dates "content" $contentHTML ) }}`).Build()
	b.AssertFileExists("public/docs/page_added/index.html", true)
	b.AssertFileContent("public/index.html", "RegularPagesRecursive: p1:p1:/docs/p1|p2title:/docs/p2|p3title:/docs/p3|p4title:/docs/p4|page_added_title:/docs/page_added|pfile:/docs/pfile|$")
}

func TestPagesFromGoTmplDraftPage(t *testing.T) {
	t.Parallel()
	b := hugolib.TestRunning(t, filesPagesFromDataTempleBasic)
	b.EditFileReplaceAll("content/docs/_content.gotmpl", `"draft" false`, `"draft" true`).Build()
	b.AssertFileContent("public/index.html", "RegularPagesRecursive: p1:p1:/docs/p1|p2title:/docs/p2|p4title:/docs/p4|pfile:/docs/pfile|$")
}

func TestPagesFromGoTmplDraftFlagFromResource(t *testing.T) {
	t.Parallel()
	b := hugolib.TestRunning(t, filesPagesFromDataTempleBasic)
	b.EditFileReplaceAll("assets/mydata.yaml", `draft: false`, `draft: true`).Build()
	b.AssertFileContent("public/index.html", "RegularPagesRecursive: p1:p1:/docs/p1|p2title:/docs/p2|p3title:/docs/p3|pfile:/docs/pfile|$")
	b.EditFileReplaceAll("assets/mydata.yaml", `draft: true`, `draft: false`).Build()
	b.AssertFileContent("public/index.html", "RegularPagesRecursive: p1:p1:/docs/p1|p2title:/docs/p2|p3title:/docs/p3|p4title:/docs/p4|pfile:/docs/pfile|$")
}

func TestPagesFromGoTmplMovePage(t *testing.T) {
	t.Parallel()
	b := hugolib.TestRunning(t, filesPagesFromDataTempleBasic)
	b.AssertFileContent("public/index.html", "RegularPagesRecursive: p1:p1:/docs/p1|p2title:/docs/p2|p3title:/docs/p3|p4title:/docs/p4|pfile:/docs/pfile|$")
	b.EditFileReplaceAll("content/docs/_content.gotmpl", `"path" "p2"`, `"path" "p2moved"`).Build()
	b.AssertFileContent("public/index.html", "RegularPagesRecursive: p1:p1:/docs/p1|p2title:/docs/p2moved|p3title:/docs/p3|p4title:/docs/p4|pfile:/docs/pfile|$")
}

func TestPagesFromGoTmplRemoveGoTmpl(t *testing.T) {
	t.Parallel()
	b := hugolib.TestRunning(t, filesPagesFromDataTempleBasic)
	b.AssertFileContent("public/index.html",
		"RegularPagesRecursive: p1:p1:/docs/p1|p2title:/docs/p2|p3title:/docs/p3|p4title:/docs/p4|pfile:/docs/pfile|$",
		"Sections: Docs:/docs|",
	)
	b.AssertFileContent("public/docs/index.html", "RegularPagesRecursive: p1:p1:/docs/p1|p2title:/docs/p2|p3title:/docs/p3|p4title:/docs/p4|pfile:/docs/pfile|$")
	b.RemoveFiles("content/docs/_content.gotmpl").Build()
	// One regular page left.
	b.AssertFileContent("public/index.html",
		"RegularPagesRecursive: pfile:/docs/pfile|$",
		"Sections: Docs:/docs|",
	)
	b.AssertFileContent("public/docs/index.html", "RegularPagesRecursive: pfile:/docs/pfile|$")
}

func TestPagesFromGoTmplLanguagePerFile(t *testing.T) {
	filesTemplate := `
-- hugo.toml --
defaultContentLanguage = "en"
defaultContentLanguageInSubdir = true
[languages]
[languages.en]
weight = 1
title = "Title"
[languages.fr]
weight = 2
title = "Titre"
disabled = DISABLE
-- layouts/_default/single.html --
Single: {{ .Title }}|{{ .Content }}|
-- content/docs/_content.gotmpl --
{{ $.AddPage  (dict "kind" "page" "path" "p1" "title" "Title" ) }}
-- content/docs/_content.fr.gotmpl --
{{ $.AddPage  (dict "kind" "page" "path" "p1" "title" "Titre" ) }}
`

	for _, disable := range []bool{false, true} {
		t.Run(fmt.Sprintf("disable=%t", disable), func(t *testing.T) {
			b := hugolib.Test(t, strings.ReplaceAll(filesTemplate, "DISABLE", fmt.Sprintf("%t", disable)))
			b.AssertFileContent("public/en/docs/p1/index.html", "Single: Title||")
			b.AssertFileExists("public/fr/docs/p1/index.html", !disable)
			if !disable {
				b.AssertFileContent("public/fr/docs/p1/index.html", "Single: Titre||")
			}
		})
	}
}

func TestPagesFromGoTmplEnableAllLanguages(t *testing.T) {
	t.Parallel()

	filesTemplate := `
-- hugo.toml --
defaultContentLanguage = "en"
defaultContentLanguageInSubdir = true
[languages]
[languages.en]
weight = 1
title = "Title"
[languages.fr]
title = "Titre"
weight = 2
disabled = DISABLE
-- i18n/en.yaml --
title: Title
-- i18n/fr.yaml --
title: Titre
-- content/docs/_content.gotmpl --
{{ .EnableAllLanguages }}
{{ $titleFromStore := .Store.Get "title" }}
{{ if not $titleFromStore }}
	{{ $titleFromStore = "notfound"}}
	{{ .Store.Set "title" site.Title }}
{{ end }}
{{ $title := printf "%s:%s:%s" site.Title (i18n "title") $titleFromStore }}
{{ $.AddPage  (dict "kind" "page" "path" "p1" "title" $title ) }}
--  layouts/_default/single.html --
Single: {{ .Title }}|{{ .Content }}|

`

	for _, disable := range []bool{false, true} {
		t.Run(fmt.Sprintf("disable=%t", disable), func(t *testing.T) {
			b := hugolib.Test(t, strings.ReplaceAll(filesTemplate, "DISABLE", fmt.Sprintf("%t", disable)))
			b.AssertFileExists("public/fr/docs/p1/index.html", !disable)
			if !disable {
				b.AssertFileContent("public/en/docs/p1/index.html", "Single: Title:Title:notfound||")
				b.AssertFileContent("public/fr/docs/p1/index.html", "Single: Titre:Titre:Title||")
			}
		})
	}
}

func TestPagesFromGoTmplMarkdownify(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ["taxonomy", "term", "rss", "sitemap"]
baseURL = "https://example.com"
-- layouts/_default/single.html --
|Content: {{ .Content }}|Title: {{ .Title }}|Path: {{ .Path }}|
-- content/docs/_content.gotmpl --
{{ $content := "**Hello World**" | markdownify }}
{{ $.AddPage (dict "path" "p1" "content" (dict "value" $content "mediaType" "text/html" )) }}
`

	b, err := hugolib.TestE(t, files)

	// This currently fails. We should fix this, but that is not a trivial task, so do it later.
	b.Assert(err, qt.IsNotNil)
	b.Assert(err.Error(), qt.Contains, "error calling markdownify: this method cannot be called before the site is fully initialized")
}

func TestPagesFromGoTmplResourceWithoutExtensionWithMediaTypeProvided(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ["taxonomy", "term", "rss", "sitemap"]
baseURL = "https://example.com"
-- layouts/_default/single.html --
|Content: {{ .Content }}|Title: {{ .Title }}|Path: {{ .Path }}|
{{ range .Resources }}
|RelPermalink: {{ .RelPermalink }}|Name: {{ .Name }}|Title: {{ .Title }}|Params: {{ .Params }}|MediaType: {{ .MediaType }}|
{{ end }}
-- content/docs/_content.gotmpl --
{{ $.AddPage (dict "path" "p1" "content" (dict "value" "**Hello World**" "mediaType" "text/markdown" )) }}
{{ $.AddResource (dict "path" "p1/myresource" "content" (dict "value" "abcde" "mediaType" "text/plain" )) }}
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/docs/p1/index.html", "RelPermalink: /docs/p1/myresource|Name: myresource|Title: myresource|Params: map[]|MediaType: text/plain|")
}

func TestPagesFromGoTmplCascade(t *testing.T) {
	t.Parallel()

	files := ` 
-- hugo.toml --
disableKinds = ["taxonomy", "term", "rss", "sitemap"]
baseURL = "https://example.com"
-- layouts/_default/single.html --
|Content: {{ .Content }}|Title: {{ .Title }}|Path: {{ .Path }}|Params: {{ .Params }}|
-- content/_content.gotmpl --
{{ $cascade := dict "params" (dict "cascadeparam1" "cascadeparam1value" ) }}
{{ $.AddPage (dict "path" "docs" "kind" "section" "cascade" $cascade ) }}
{{ $.AddPage (dict "path" "docs/p1" "content" (dict "value" "**Hello World**" "mediaType" "text/markdown" )) }}

`
	b := hugolib.Test(t, files)

	b.AssertFileContent("public/docs/p1/index.html", "|Path: /docs/p1|Params: map[cascadeparam1:cascadeparam1value")
}

func TestPagesFromGoBuildOptions(t *testing.T) {
	t.Parallel()

	files := ` 
-- hugo.toml --
disableKinds = ["taxonomy", "term", "rss", "sitemap"]
baseURL = "https://example.com"
-- layouts/_default/single.html --
|Content: {{ .Content }}|Title: {{ .Title }}|Path: {{ .Path }}|Params: {{ .Params }}|
-- content/_content.gotmpl --
{{ $.AddPage (dict "path" "docs/p1" "content" (dict "value" "**Hello World**" "mediaType" "text/markdown" )) }}
{{ $never := dict "list" "never"  "publishResources" false "render" "never"  }}
{{ $.AddPage (dict "path" "docs/p2" "content" (dict "value" "**Hello World**" "mediaType" "text/markdown" ) "build" $never ) }}


`
	b := hugolib.Test(t, files)

	b.AssertFileExists("public/docs/p1/index.html", true)
	b.AssertFileExists("public/docs/p2/index.html", false)
}

func TestPagesFromGoPathsWithDotsIssue12493(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['home','section','rss','sitemap','taxonomy','term']
-- content/_content.gotmpl --
{{ .AddPage (dict "path" "s-1.2.3/p-4.5.6" "title" "p-4.5.6") }}
-- layouts/_default/single.html --
{{ .Title }}
`

	b := hugolib.Test(t, files)

	b.AssertFileExists("public/s-1.2.3/p-4.5.6/index.html", true)
}

func TestPagesFromGoParamsIssue12497(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['home','section','rss','sitemap','taxonomy','term']
-- content/_content.gotmpl --
{{ .AddPage (dict "path" "p1" "title" "p1" "params" (dict "paraM1" "param1v" )) }}
{{ .AddResource (dict "path" "p1/data1.yaml" "content" (dict "value" "data1" ) "params" (dict "paraM1" "param1v" )) }}
-- layouts/_default/single.html --
{{ .Title }}|{{ .Params.paraM1 }}
{{ range .Resources }}
{{ .Name }}|{{ .Params.paraM1 }}
{{ end }}
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/p1/index.html",
		"p1|param1v",
		"data1.yaml|param1v",
	)
}

func TestPagesFromGoTmplPathWarningsPathPage(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com"
disableKinds = ['home','section','rss','sitemap','taxonomy','term']
printPathWarnings = true
-- content/_content.gotmpl --
{{ .AddPage (dict "path" "p1" "title" "p1" ) }}
{{ .AddPage (dict "path" "p2" "title" "p2" ) }}
-- content/p1.md --
---
title: "p1"
---
-- layouts/_default/single.html --
{{ .Title }}|
`

	b := hugolib.Test(t, files, hugolib.TestOptWarn())

	b.AssertFileContent("public/p1/index.html", "p1|")

	b.AssertLogContains("Duplicate content path")

	files = strings.ReplaceAll(files, `"path" "p1"`, `"path" "p1new"`)

	b = hugolib.Test(t, files, hugolib.TestOptWarn())

	b.AssertLogContains("! WARN")
}

func TestPagesFromGoTmplPathWarningsPathResource(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com"
disableKinds = ['home','section','rss','sitemap','taxonomy','term']
printPathWarnings = true
-- content/_content.gotmpl --
{{ .AddResource (dict "path" "p1/data1.yaml" "content" (dict "value" "data1" ) ) }}
{{ .AddResource (dict "path" "p1/data2.yaml" "content" (dict "value" "data2" ) ) }}

-- content/p1/index.md --
---
title: "p1"
---
-- content/p1/data1.yaml --
value: data1
-- layouts/_default/single.html --
{{ .Title }}|
`

	b := hugolib.Test(t, files, hugolib.TestOptWarn())

	b.AssertFileContent("public/p1/index.html", "p1|")

	b.AssertLogContains("Duplicate resource path")

	files = strings.ReplaceAll(files, `"path" "p1/data1.yaml"`, `"path" "p1/data1new.yaml"`)

	b = hugolib.Test(t, files, hugolib.TestOptWarn())

	b.AssertLogContains("! WARN")
}

func TestPagesFromGoTmplShortcodeNoPreceddingCharacterIssue12544(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['home','rss','section','sitemap','taxonomy','term']
-- content/_content.gotmpl --
{{ $content := dict "mediaType" "text/html" "value" "x{{< sc >}}" }}
{{ .AddPage (dict "content" $content "path" "a") }}

{{ $content := dict "mediaType" "text/html" "value" "{{< sc >}}" }}
{{ .AddPage (dict "content" $content "path" "b") }}
-- layouts/_default/single.html --
|{{ .Content }}|
-- layouts/shortcodes/sc.html --
foo
{{- /**/ -}}
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/a/index.html", "|xfoo|")
	b.AssertFileContent("public/b/index.html", "|foo|") // fails
}

func TestPagesFromGoTmplMenus(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['rss','section','sitemap','taxonomy','term']

[menus]
[[menus.main]]
name = "Main"
[[menus.footer]]
name = "Footer"
-- content/_content.gotmpl --
{{ .AddPage (dict "path" "p1" "title" "p1" "menus" "main" ) }}
{{ .AddPage (dict "path" "p2" "title" "p2" "menus" (slice "main" "footer")) }}
-- layouts/index.html --
Main: {{ range index site.Menus.main }}{{ .Name }}|{{ end }}|
Footer: {{ range index site.Menus.footer }}{{ .Name }}|{{ end }}|

`
	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html",
		"Main: Main|p1|p2||",
		"Footer: Footer|p2||",
	)
}

func TestPagesFromGoTmplMore(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['home','rss','section','sitemap','taxonomy','term']
[markup.goldmark.renderer]
unsafe = true
-- content/s1/_content.gotmpl --
{{ $page := dict
	"content" (dict "mediaType" "text/markdown" "value" "aaa <!--more--> bbb")
	"title" "p1"
	"path" "p1"
  }}
  {{ .AddPage $page }}
-- layouts/_default/single.html --
summary: {{ .Summary }}|content: {{ .Content}}
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/s1/p1/index.html",
		"<p>aaa</p>|content: <p>aaa</p>\n<p>bbb</p>",
	)
}
