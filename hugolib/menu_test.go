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
	"testing"

	"fmt"

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

func TestSectionPagesMenu(t *testing.T) {
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

func TestMenuFrontMatter(t *testing.T) {

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
func TestMenuPageMultipleOutputFormats(t *testing.T) {

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
func TestMenuPageSortByDate(t *testing.T) {

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
