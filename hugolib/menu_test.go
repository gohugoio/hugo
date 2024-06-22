// Copyright 2019 The Hugo Authors. All rights reserved.
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

package hugolib

import (
	"fmt"
	"testing"

	qt "github.com/frankban/quicktest"
)

const (
	menuPageTemplate = `---
title: %q
weight: %d
menu:
  %s:
    title: %s
    weight: %d
---
# Doc Menu
`
)

func TestMenusSectionPagesMenu(t *testing.T) {
	t.Parallel()

	siteConfig := `
baseurl = "http://example.com/"
title = "Section Menu"
sectionPagesMenu = "sect"
`

	b := newTestSitesBuilder(t).WithConfigFile("toml", siteConfig)

	b.WithTemplates(
		"partials/menu.html",
		`{{- $p := .page -}}
{{- $m := .menu -}}
{{ range (index $p.Site.Menus $m) -}}
{{- .URL }}|{{ .Name }}|{{ .Title }}|{{ .Weight -}}|
{{- if $p.IsMenuCurrent $m . }}IsMenuCurrent{{ else }}-{{ end -}}|
{{- if $p.HasMenuCurrent $m . }}HasMenuCurrent{{ else }}-{{ end -}}|
{{- end -}}
`,
		"_default/single.html",
		`Single|{{ .Title }}
Menu Sect:  {{ partial "menu.html" (dict "page" . "menu" "sect") }}
Menu Main:  {{ partial "menu.html" (dict "page" . "menu" "main") }}`,
		"_default/list.html", "List|{{ .Title }}|{{ .Content }}",
	)

	b.WithContent(
		"sect1/p1.md", fmt.Sprintf(menuPageTemplate, "p1", 1, "main", "atitle1", 40),
		"sect1/p2.md", fmt.Sprintf(menuPageTemplate, "p2", 2, "main", "atitle2", 30),
		"sect2/p3.md", fmt.Sprintf(menuPageTemplate, "p3", 3, "main", "atitle3", 20),
		"sect2/p4.md", fmt.Sprintf(menuPageTemplate, "p4", 4, "main", "atitle4", 10),
		"sect3/p5.md", fmt.Sprintf(menuPageTemplate, "p5", 5, "main", "atitle5", 5),
		"sect1/_index.md", newTestPage("Section One", "2017-01-01", 100),
		"sect5/_index.md", newTestPage("Section Five", "2017-01-01", 10),
	)

	b.Build(BuildCfg{})
	h := b.H

	s := h.Sites[0]

	b.Assert(len(s.Menus()), qt.Equals, 2)

	p1 := s.RegularPages()[0].Menus()

	// There is only one menu in the page, but it is "member of" 2
	b.Assert(len(p1), qt.Equals, 1)

	b.AssertFileContent("public/sect1/p1/index.html", "Single",
		"Menu Sect:  "+
			"/sect5/|Section Five|Section Five|10|-|-|"+
			"/sect1/|Section One|Section One|100|-|HasMenuCurrent|"+
			"/sect2/|Sect2s|Sect2s|0|-|-|"+
			"/sect3/|Sect3s|Sect3s|0|-|-|",
		"Menu Main:  "+
			"/sect3/p5/|p5|atitle5|5|-|-|"+
			"/sect2/p4/|p4|atitle4|10|-|-|"+
			"/sect2/p3/|p3|atitle3|20|-|-|"+
			"/sect1/p2/|p2|atitle2|30|-|-|"+
			"/sect1/p1/|p1|atitle1|40|IsMenuCurrent|-|",
	)

	b.AssertFileContent("public/sect2/p3/index.html", "Single",
		"Menu Sect:  "+
			"/sect5/|Section Five|Section Five|10|-|-|"+
			"/sect1/|Section One|Section One|100|-|-|"+
			"/sect2/|Sect2s|Sect2s|0|-|HasMenuCurrent|"+
			"/sect3/|Sect3s|Sect3s|0|-|-|")
}

// related issue #7594
func TestMenusSort(t *testing.T) {
	b := newTestSitesBuilder(t).WithSimpleConfigFile()

	b.WithTemplatesAdded("index.html", `
{{ range $k, $v := .Site.Menus.main }}
Default1|{{ $k }}|{{ $v.Weight }}|{{ $v.Name }}|{{ .URL }}|{{ $v.Page }}{{ end }}
{{ range $k, $v := .Site.Menus.main.ByWeight }}
ByWeight|{{ $k }}|{{ $v.Weight }}|{{ $v.Name }}|{{ .URL }}|{{ $v.Page }}{{ end }}
{{ range $k, $v := (.Site.Menus.main.ByWeight).Reverse }}
Reverse|{{ $k }}|{{ $v.Weight }}|{{ $v.Name }}|{{ .URL }}|{{ $v.Page }}{{ end }}
{{ range $k, $v := .Site.Menus.main }}
Default2|{{ $k }}|{{ $v.Weight }}|{{ $v.Name }}|{{ .URL }}|{{ $v.Page }}{{ end }}
{{ range $k, $v := .Site.Menus.main.ByWeight }}
ByWeight|{{ $k }}|{{ $v.Weight }}|{{ $v.Name }}|{{ .URL }}|{{ $v.Page }}{{ end }}
{{ range $k, $v := .Site.Menus.main }}
Default3|{{ $k }}|{{ $v.Weight }}|{{ $v.Name }}|{{ .URL }}|{{ $v.Page }}{{ end }}
`)

	b.WithContent("_index.md", `
---
title: Home
menu:
  main:
    weight: 100
---`)

	b.WithContent("blog/A.md", `
---
title: "A"
menu:
  main:
    weight: 10
---
`)

	b.WithContent("blog/B.md", `
---
title: "B"
menu:
  main:
    weight: 20
---
`)
	b.WithContent("blog/C.md", `
---
title: "C"
menu:
  main:
    weight: 30
---
`)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/index.html",
		`Default1|0|10|A|/blog/a/|Page(/blog/a)
        Default1|1|20|B|/blog/b/|Page(/blog/b)
        Default1|2|30|C|/blog/c/|Page(/blog/c)
        Default1|3|100|Home|/|Page(/)

        ByWeight|0|10|A|/blog/a/|Page(/blog/a)
        ByWeight|1|20|B|/blog/b/|Page(/blog/b)
        ByWeight|2|30|C|/blog/c/|Page(/blog/c)
        ByWeight|3|100|Home|/|Page(/)

        Reverse|0|100|Home|/|Page(/)
        Reverse|1|30|C|/blog/c/|Page(/blog/c)
        Reverse|2|20|B|/blog/b/|Page(/blog/b)
        Reverse|3|10|A|/blog/a/|Page(/blog/a)

        Default2|0|10|A|/blog/a/|Page(/blog/a)
        Default2|1|20|B|/blog/b/|Page(/blog/b)
        Default2|2|30|C|/blog/c/|Page(/blog/c)
        Default2|3|100|Home|/|Page(/)

        ByWeight|0|10|A|/blog/a/|Page(/blog/a)
        ByWeight|1|20|B|/blog/b/|Page(/blog/b)
        ByWeight|2|30|C|/blog/c/|Page(/blog/c)
        ByWeight|3|100|Home|/|Page(/)

        Default3|0|10|A|/blog/a/|Page(/blog/a)
        Default3|1|20|B|/blog/b/|Page(/blog/b)
        Default3|2|30|C|/blog/c/|Page(/blog/c)
        Default3|3|100|Home|/|Page(/)`,
	)
}

func TestMenusFrontMatter(t *testing.T) {
	b := newTestSitesBuilder(t).WithSimpleConfigFile()

	b.WithTemplatesAdded("index.html", `
Main: {{ len .Site.Menus.main }}
Other: {{ len .Site.Menus.other }}
{{ range .Site.Menus.main }}
* Main|{{ .Name }}: {{ .URL }}
{{ end }}
{{ range .Site.Menus.other }}
* Other|{{ .Name }}: {{ .URL }}
{{ end }}
`)

	// Issue #5828
	b.WithContent("blog/page1.md", `
---
title: "P1"
menu: main
---

`)

	b.WithContent("blog/page2.md", `
---
title: "P2"
menu: [main,other]
---

`)

	b.WithContent("blog/page3.md", `
---
title: "P3"
menu:
  main:
    weight: 30
---
`)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/index.html",
		"Main: 3", "Other: 1",
		"Main|P1: /blog/page1/",
		"Other|P2: /blog/page2/",
	)
}

// https://github.com/gohugoio/hugo/issues/5849
func TestMenusPageMultipleOutputFormats(t *testing.T) {
	config := `
baseURL = "https://example.com"

# DAMP is similar to AMP, but not permalinkable.
[outputFormats]
[outputFormats.damp]
mediaType = "text/html"
path = "damp"

`

	b := newTestSitesBuilder(t).WithConfigFile("toml", config)
	b.WithContent("_index.md", `
---
Title: Home Sweet Home
outputs: [ "html", "amp" ]
menu: "main"
---

`)

	b.WithContent("blog/html-amp.md", `
---
Title: AMP and HTML
outputs: [ "html", "amp" ]
menu: "main"
---

`)

	b.WithContent("blog/html.md", `
---
Title: HTML only
outputs: [ "html" ]
menu: "main"
---

`)

	b.WithContent("blog/amp.md", `
---
Title: AMP only
outputs: [ "amp" ]
menu: "main"
---

`)

	b.WithTemplatesAdded("index.html", `{{ range .Site.Menus.main }}{{ .Title }}|{{ .URL }}|{{ end }}`)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/index.html", "AMP and HTML|/blog/html-amp/|AMP only|/amp/blog/amp/|Home Sweet Home|/|HTML only|/blog/html/|")
	b.AssertFileContent("public/amp/index.html", "AMP and HTML|/amp/blog/html-amp/|AMP only|/amp/blog/amp/|Home Sweet Home|/amp/|HTML only|/blog/html/|")
}

// https://github.com/gohugoio/hugo/issues/5989
func TestMenusPageSortByDate(t *testing.T) {
	b := newTestSitesBuilder(t).WithSimpleConfigFile()

	b.WithContent("blog/a.md", `
---
Title: A
date: 2019-01-01
menu:
  main:
    identifier: "a"
    weight: 1
---

`)

	b.WithContent("blog/b.md", `
---
Title: B
date: 2018-01-02
menu:
  main:
    parent: "a"
    weight: 100
---

`)

	b.WithContent("blog/c.md", `
---
Title: C
date: 2019-01-03
menu:
  main:
    parent: "a"
    weight: 10
---

`)

	b.WithTemplatesAdded("index.html", `{{ range .Site.Menus.main }}{{ .Title }}|Children: 
{{- $children := sort .Children ".Page.Date" "desc" }}{{ range $children }}{{ .Title }}|{{ end }}{{ end }}
	
`)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/index.html", "A|Children:C|B|")
}

// Issue #8825
func TestMenuParamsEmptyYaml(t *testing.T) {
	b := newTestSitesBuilder(t).WithConfigFile("yaml", `

`)

	b.WithTemplates("index.html", `{{ site.Menus }}`)

	b.WithContent("p1.md", `---
menus:
  main: 
    identity: journal
    weight: 2
    params:
---	
`)
	b.Build(BuildCfg{})
}

func TestMenuParams(t *testing.T) {
	b := newTestSitesBuilder(t).WithConfigFile("toml", `
[[menus.main]]
identifier = "contact"
title = "Contact Us"
url = "mailto:noreply@example.com"
weight = 300
[menus.main.params]
foo = "foo_config"	
key2 = "key2_config"	
camelCase = "camelCase_config"	
`)

	b.WithTemplatesAdded("index.html", `
Main: {{ len .Site.Menus.main }}
{{ range .Site.Menus.main }}
foo: {{ .Params.foo }}
key2: {{ .Params.KEy2 }}
camelCase: {{ .Params.camelcase }}
{{ end }}
`)

	b.WithContent("_index.md", `
---
title: "Home"
menu:
  main:
    weight: 10
    params:
      foo: "foo_content"
      key2: "key2_content"
      camelCase: "camelCase_content"
---
`)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/index.html", `
Main: 2

foo: foo_content
key2: key2_content
camelCase: camelCase_content

foo: foo_config
key2: key2_config
camelCase: camelCase_config
`)
}

func TestMenusShadowMembers(t *testing.T) {
	b := newTestSitesBuilder(t).WithConfigFile("toml", `
[[menus.main]]
identifier = "contact"
pageRef = "contact"
title = "Contact Us"
url = "mailto:noreply@example.com"
weight = 1
[[menus.main]]
pageRef = "/blog/post3"
title = "My Post 3"
url = "/blog/post3"
	
`)

	commonTempl := `
Main: {{ len .Site.Menus.main }}
{{ range .Site.Menus.main }}
{{ .Title }}|HasMenuCurrent: {{ $.HasMenuCurrent "main" . }}|Page: {{ .Page }}
{{ .Title }}|IsMenuCurrent: {{ $.IsMenuCurrent "main" . }}|Page: {{ .Page }}
{{ end }}
`

	b.WithTemplatesAdded("index.html", commonTempl)
	b.WithTemplatesAdded("_default/single.html", commonTempl)

	b.WithContent("_index.md", `
---
title: "Home"
menu:
  main:
    weight: 10
---
`)

	b.WithContent("blog/_index.md", `
---
title: "Blog"
menu:
  main:
    weight: 20
---
`)

	b.WithContent("blog/post1.md", `
---
title: "My Post 1: With  No Menu Defined"
---
`)

	b.WithContent("blog/post2.md", `
---
title: "My Post 2: With Menu Defined"
menu:
  main:
    weight: 30
---
`)

	b.WithContent("blog/post3.md", `
---
title: "My Post 2: With  No Menu Defined"
---
`)

	b.WithContent("contact.md", `
---
title: "Contact: With  No Menu Defined"
---
`)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/index.html", `
Main: 5
Home|HasMenuCurrent: false|Page: Page(/)
Blog|HasMenuCurrent: false|Page: Page(/blog)
My Post 2: With Menu Defined|HasMenuCurrent: false|Page: Page(/blog/post2)
My Post 3|HasMenuCurrent: false|Page: Page(/blog/post3)
Contact Us|HasMenuCurrent: false|Page: Page(/contact)
`)

	b.AssertFileContent("public/blog/post1/index.html", `
Home|HasMenuCurrent: false|Page: Page(/)
Blog|HasMenuCurrent: true|Page: Page(/blog)
`)

	b.AssertFileContent("public/blog/post2/index.html", `
Home|HasMenuCurrent: false|Page: Page(/)
Blog|HasMenuCurrent: true|Page: Page(/blog)
Blog|IsMenuCurrent: false|Page: Page(/blog)
`)

	b.AssertFileContent("public/blog/post3/index.html", `
Home|HasMenuCurrent: false|Page: Page(/)
Blog|HasMenuCurrent: true|Page: Page(/blog)
`)

	b.AssertFileContent("public/contact/index.html", `
Contact Us|HasMenuCurrent: false|Page: Page(/contact)
Contact Us|IsMenuCurrent: true|Page: Page(/contact)
Blog|HasMenuCurrent: false|Page: Page(/blog)
Blog|IsMenuCurrent: false|Page: Page(/blog)
`)
}

// Issue 9846
func TestMenuHasMenuCurrentSection(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
disableKinds = ['RSS','sitemap','taxonomy','term']
[[menu.main]]
name = 'Home'
pageRef = '/'
weight = 1

[[menu.main]]
name = 'Tests'
pageRef = '/tests'
weight = 2
[[menu.main]]
name = 'Test 1'
pageRef = '/tests/test-1'
parent = 'Tests'
weight = 1

-- content/tests/test-1.md --
---
title: "Test 1"
---
-- layouts/_default/list.html --
{{ range site.Menus.main }}
{{ .Name }}|{{ .URL }}|IsMenuCurrent = {{ $.IsMenuCurrent "main" . }}|HasMenuCurrent = {{ $.HasMenuCurrent "main" . }}|
{{ range .Children }}
{{ .Name }}|{{ .URL }}|IsMenuCurrent = {{ $.IsMenuCurrent "main" . }}|HasMenuCurrent = {{ $.HasMenuCurrent "main" . }}|
{{ end }}
{{ end }}

{{/* Some tests for issue 9925 */}}
{{ $page := .Site.GetPage "tests/test-1" }}
{{ $section := site.GetPage "tests" }}

Home IsAncestor Self: {{ site.Home.IsAncestor site.Home }}
Home IsDescendant Self: {{ site.Home.IsDescendant site.Home }}
Section IsAncestor Self: {{ $section.IsAncestor $section }}
Section IsDescendant Self: {{ $section.IsDescendant $section}}
Page IsAncestor Self: {{ $page.IsAncestor $page }}
Page IsDescendant Self: {{ $page.IsDescendant $page}}
`

	b := Test(t, files)

	b.AssertFileContent("public/tests/index.html", `
Tests|/tests/|IsMenuCurrent = true|HasMenuCurrent = false
Home IsAncestor Self: false
Home IsDescendant Self: false
Section IsAncestor Self: false
Section IsDescendant Self: false
Page IsAncestor Self: false
Page IsDescendant Self: false
`)
}

func TestMenusNewConfigSetup(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com"
title = "Hugo Menu Test"
[menus]
[[menus.main]]
name = "Home"
url = "/"
pre = "<span>"
post = "</span>"
weight = 1
-- layouts/index.html --
{{ range $i, $e := site.Menus.main }}
Menu Item: {{ $i }}: {{ .Pre }}{{ .Name }}{{ .Post }}|{{ .URL }}|
{{ end }}
`

	b := Test(t, files)

	b.AssertFileContent("public/index.html", `
Menu Item: 0: <span>Home</span>|/|	
`)
}

// Issue #11062
func TestMenusSubDirInBaseURL(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com/foo/"
title = "Hugo Menu Test"
[menus]
[[menus.main]]
name = "Posts"
url = "/posts"
weight = 1
-- layouts/index.html --
{{ range $i, $e := site.Menus.main }}
Menu Item: {{ $i }}|{{ .URL }}|
{{ end }}
`

	b := Test(t, files)

	b.AssertFileContent("public/index.html", `
Menu Item: 0|/foo/posts|
`)
}

func TestSectionPagesMenuMultilingualWarningIssue12306(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['section','rss','sitemap','taxonomy','term']
defaultContentLanguageInSubdir = true
sectionPagesMenu = "main"
[languages.en]
[languages.fr]
-- layouts/_default/home.html --
{{- range site.Menus.main -}}
  <a href="{{ .URL }}">{{ .Name }}</a>
{{- end -}}
-- layouts/_default/single.html --
{{ .Title }}
-- content/p1.en.md --
---
title: p1
menu: main
---
-- content/p1.fr.md --
---
title: p1
menu: main
---
-- content/p2.en.md --
---
title: p2
menu: main
---
`

	b := Test(t, files, TestOptWarn())

	b.AssertFileContent("public/en/index.html", `<a href="/en/p1/">p1</a><a href="/en/p2/">p2</a>`)
	b.AssertFileContent("public/fr/index.html", `<a href="/fr/p1/">p1</a>`)
	b.AssertLogContains("! WARN")
}

func TestSectionPagesIssue12399(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['rss','sitemap','taxonomy','term']
capitalizeListTitles = false
pluralizeListTitles = false
sectionPagesMenu = 'main'
-- content/p1.md --
---
title: p1
---
-- content/s1/p2.md --
---
title: p2
menus: main
---
-- content/s1/p3.md --
---
title: p3
---
-- layouts/_default/list.html --
{{ range site.Menus.main }}<a href="{{ .URL }}">{{ .Name }}</a>{{ end }}
-- layouts/_default/single.html --
{{ .Title }}
`

	b := Test(t, files)

	b.AssertFileExists("public/index.html", true)
	b.AssertFileContent("public/index.html", `<a href="/s1/p2/">p2</a><a href="/s1/">s1</a>`)
}
