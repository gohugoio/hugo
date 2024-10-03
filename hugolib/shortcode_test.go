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
	"context"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/resources/kinds"

	"github.com/gohugoio/hugo/parser/pageparser"

	qt "github.com/frankban/quicktest"
)

func TestExtractShortcodes(t *testing.T) {
	b := newTestSitesBuilder(t).WithSimpleConfigFile()

	b.WithTemplates(
		"default/single.html", `EMPTY`,
		"_internal/shortcodes/tag.html", `tag`,
		"_internal/shortcodes/legacytag.html", `{{ $_hugo_config := "{ \"version\": 1 }" }}tag`,
		"_internal/shortcodes/sc1.html", `sc1`,
		"_internal/shortcodes/sc2.html", `sc2`,
		"_internal/shortcodes/inner.html", `{{with .Inner }}{{ . }}{{ end }}`,
		"_internal/shortcodes/inner2.html", `{{.Inner}}`,
		"_internal/shortcodes/inner3.html", `{{.Inner}}`,
	).WithContent("page.md", `---
title: "Shortcodes Galore!"
---
`)

	b.CreateSites().Build(BuildCfg{})

	s := b.H.Sites[0]

	// Make it more regexp friendly
	strReplacer := strings.NewReplacer("[", "{", "]", "}")

	str := func(s *shortcode) string {
		if s == nil {
			return "<nil>"
		}

		var version int
		if s.info != nil {
			version = s.info.ParseInfo().Config.Version
		}
		return strReplacer.Replace(fmt.Sprintf("%s;inline:%t;closing:%t;inner:%v;params:%v;ordinal:%d;markup:%t;version:%d;pos:%d",
			s.name, s.isInline, s.isClosing, s.inner, s.params, s.ordinal, s.doMarkup, version, s.pos))
	}

	regexpCheck := func(re string) func(c *qt.C, shortcode *shortcode, err error) {
		return func(c *qt.C, shortcode *shortcode, err error) {
			c.Assert(err, qt.IsNil)
			c.Assert(str(shortcode), qt.Matches, ".*"+re+".*")
		}
	}

	for _, test := range []struct {
		name  string
		input string
		check func(c *qt.C, shortcode *shortcode, err error)
	}{
		{"one shortcode, no markup", "{{< tag >}}", regexpCheck("tag.*closing:false.*markup:false")},
		{"one shortcode, markup", "{{% tag %}}", regexpCheck("tag.*closing:false.*markup:true;version:2")},
		{"one shortcode, markup, legacy", "{{% legacytag %}}", regexpCheck("tag.*closing:false.*markup:true;version:1")},
		{"outer shortcode markup", "{{% inner %}}{{< tag >}}{{% /inner %}}", regexpCheck("inner.*closing:true.*markup:true")},
		{"inner shortcode markup", "{{< inner >}}{{% tag %}}{{< /inner >}}", regexpCheck("inner.*closing:true.*;markup:false;version:2")},
		{"one pos param", "{{% tag param1 %}}", regexpCheck("tag.*params:{param1}")},
		{"two pos params", "{{< tag param1 param2>}}", regexpCheck("tag.*params:{param1 param2}")},
		{"one named param", `{{% tag param1="value" %}}`, regexpCheck("tag.*params:map{param1:value}")},
		{"two named params", `{{< tag param1="value1" param2="value2" >}}`, regexpCheck("tag.*params:map{param\\d:value\\d param\\d:value\\d}")},
		{"inner", `{{< inner >}}Inner Content{{< / inner >}}`, regexpCheck("inner;inline:false;closing:true;inner:{Inner Content};")},
		// issue #934
		{"inner self-closing", `{{< inner />}}`, regexpCheck("inner;.*inner:{}")},
		{
			"nested inner", `{{< inner >}}Inner Content->{{% inner2 param1 %}}inner2txt{{% /inner2 %}}Inner close->{{< / inner >}}`,
			regexpCheck("inner;.*inner:{Inner Content->.*Inner close->}"),
		},
		{
			"nested, nested inner", `{{< inner >}}inner2->{{% inner2 param1 %}}inner2txt->inner3{{< inner3>}}inner3txt{{</ inner3 >}}{{% /inner2 %}}final close->{{< / inner >}}`,
			regexpCheck("inner:{inner2-> inner2.*{{inner2txt->inner3.*final close->}"),
		},
		{"closed without content", `{{< inner param1 >}}{{< / inner >}}`, regexpCheck("inner.*inner:{}")},
		{"inline", `{{< my.inline >}}Hi{{< /my.inline >}}`, regexpCheck("my.inline;inline:true;closing:true;inner:{Hi};")},
	} {

		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			c := qt.New(t)

			p, err := pageparser.ParseMain(strings.NewReader(test.input), pageparser.Config{})
			c.Assert(err, qt.IsNil)
			handler := newShortcodeHandler("", s)
			iter := p.Iterator()

			short, err := handler.extractShortcode(0, 0, p.Input(), iter)

			test.check(c, short, err)
		})
	}
}

func TestShortcodeMultipleOutputFormats(t *testing.T) {
	t.Parallel()

	siteConfig := `
baseURL = "http://example.com/blog"

disableKinds = ["section", "term", "taxonomy", "RSS", "sitemap", "robotsTXT", "404"]

[pagination]
pagerSize = 1

[outputs]
home = [ "HTML", "AMP", "Calendar" ]
page =  [ "HTML", "AMP", "JSON" ]

`

	pageTemplate := `---
title: "%s"
---
# Doc

{{< myShort >}}
{{< noExt >}}
{{%% onlyHTML %%}}

{{< myInner >}}{{< myShort >}}{{< /myInner >}}

`

	pageTemplateCSVOnly := `---
title: "%s"
outputs: ["CSV"]
---
# Doc

CSV: {{< myShort >}}
`

	b := newTestSitesBuilder(t).WithConfigFile("toml", siteConfig)
	b.WithTemplates(
		"layouts/_default/single.html", `Single HTML: {{ .Title }}|{{ .Content }}`,
		"layouts/_default/single.json", `Single JSON: {{ .Title }}|{{ .Content }}`,
		"layouts/_default/single.csv", `Single CSV: {{ .Title }}|{{ .Content }}`,
		"layouts/index.html", `Home HTML: {{ .Title }}|{{ .Content }}`,
		"layouts/index.amp.html", `Home AMP: {{ .Title }}|{{ .Content }}`,
		"layouts/index.ics", `Home Calendar: {{ .Title }}|{{ .Content }}`,
		"layouts/shortcodes/myShort.html", `ShortHTML`,
		"layouts/shortcodes/myShort.amp.html", `ShortAMP`,
		"layouts/shortcodes/myShort.csv", `ShortCSV`,
		"layouts/shortcodes/myShort.ics", `ShortCalendar`,
		"layouts/shortcodes/myShort.json", `ShortJSON`,
		"layouts/shortcodes/noExt", `ShortNoExt`,
		"layouts/shortcodes/onlyHTML.html", `ShortOnlyHTML`,
		"layouts/shortcodes/myInner.html", `myInner:--{{- .Inner -}}--`,
	)

	b.WithContent("_index.md", fmt.Sprintf(pageTemplate, "Home"),
		"sect/mypage.md", fmt.Sprintf(pageTemplate, "Single"),
		"sect/mycsvpage.md", fmt.Sprintf(pageTemplateCSVOnly, "Single CSV"),
	)

	b.Build(BuildCfg{})
	h := b.H
	b.Assert(len(h.Sites), qt.Equals, 1)

	s := h.Sites[0]
	home := s.getPageOldVersion(kinds.KindHome)
	b.Assert(home, qt.Not(qt.IsNil))
	b.Assert(len(home.OutputFormats()), qt.Equals, 3)

	b.AssertFileContent("public/index.html",
		"Home HTML",
		"ShortHTML",
		"ShortNoExt",
		"ShortOnlyHTML",
		"myInner:--ShortHTML--",
	)

	b.AssertFileContent("public/amp/index.html",
		"Home AMP",
		"ShortAMP",
		"ShortNoExt",
		"ShortOnlyHTML",
		"myInner:--ShortAMP--",
	)

	b.AssertFileContent("public/index.ics",
		"Home Calendar",
		"ShortCalendar",
		"ShortNoExt",
		"ShortOnlyHTML",
		"myInner:--ShortCalendar--",
	)

	b.AssertFileContent("public/sect/mypage/index.html",
		"Single HTML",
		"ShortHTML",
		"ShortNoExt",
		"ShortOnlyHTML",
		"myInner:--ShortHTML--",
	)

	b.AssertFileContent("public/sect/mypage/index.json",
		"Single JSON",
		"ShortJSON",
		"ShortNoExt",
		"ShortOnlyHTML",
		"myInner:--ShortJSON--",
	)

	b.AssertFileContent("public/amp/sect/mypage/index.html",
		// No special AMP template
		"Single HTML",
		"ShortAMP",
		"ShortNoExt",
		"ShortOnlyHTML",
		"myInner:--ShortAMP--",
	)

	b.AssertFileContent("public/sect/mycsvpage/index.csv",
		"Single CSV",
		"ShortCSV",
	)
}

func BenchmarkReplaceShortcodeTokens(b *testing.B) {
	type input struct {
		in           []byte
		tokenHandler func(ctx context.Context, token string) ([]byte, error)
		expect       []byte
	}

	data := []struct {
		input        string
		replacements map[string]string
		expect       []byte
	}{
		{"Hello HAHAHUGOSHORTCODE-1HBHB.", map[string]string{"HAHAHUGOSHORTCODE-1HBHB": "World"}, []byte("Hello World.")},
		{strings.Repeat("A", 100) + " HAHAHUGOSHORTCODE-1HBHB.", map[string]string{"HAHAHUGOSHORTCODE-1HBHB": "Hello World"}, []byte(strings.Repeat("A", 100) + " Hello World.")},
		{strings.Repeat("A", 500) + " HAHAHUGOSHORTCODE-1HBHB.", map[string]string{"HAHAHUGOSHORTCODE-1HBHB": "Hello World"}, []byte(strings.Repeat("A", 500) + " Hello World.")},
		{strings.Repeat("ABCD ", 500) + " HAHAHUGOSHORTCODE-1HBHB.", map[string]string{"HAHAHUGOSHORTCODE-1HBHB": "Hello World"}, []byte(strings.Repeat("ABCD ", 500) + " Hello World.")},
		{strings.Repeat("A ", 3000) + " HAHAHUGOSHORTCODE-1HBHB." + strings.Repeat("BC ", 1000) + " HAHAHUGOSHORTCODE-1HBHB.", map[string]string{"HAHAHUGOSHORTCODE-1HBHB": "Hello World"}, []byte(strings.Repeat("A ", 3000) + " Hello World." + strings.Repeat("BC ", 1000) + " Hello World.")},
	}

	cnt := 0
	in := make([]input, b.N*len(data))
	for i := 0; i < b.N; i++ {
		for _, this := range data {
			replacements := make(map[string]shortcodeRenderer)
			for k, v := range this.replacements {
				replacements[k] = prerenderedShortcode{s: v}
			}
			tokenHandler := func(ctx context.Context, token string) ([]byte, error) {
				return []byte(this.replacements[token]), nil
			}
			in[cnt] = input{[]byte(this.input), tokenHandler, this.expect}
			cnt++
		}
	}

	b.ResetTimer()
	cnt = 0
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		for j := range data {
			currIn := in[cnt]
			cnt++
			results, err := expandShortcodeTokens(ctx, currIn.in, currIn.tokenHandler)
			if err != nil {
				b.Fatalf("[%d] failed: %s", i, err)
				continue
			}
			if len(results) != len(currIn.expect) {
				b.Fatalf("[%d] replaceShortcodeTokens, got \n%q but expected \n%q", j, results, currIn.expect)
			}

		}
	}
}

func BenchmarkShortcodesInSite(b *testing.B) {
	files := `
-- config.toml --
-- layouts/shortcodes/mark1.md --
{{ .Inner }}
-- layouts/shortcodes/mark2.md --
1. Item Mark2 1
1. Item Mark2 2
   1. Item Mark2 2-1
1. Item Mark2 3
-- layouts/_default/single.html --
{{ .Content }}
`

	content := `
---
title: "Markdown Shortcode"
---

## List

1. List 1
	{{§ mark1 §}}
	1. Item Mark1 1
	1. Item Mark1 2
	{{§ mark2 §}}
	{{§ /mark1 §}}

`

	for i := 1; i < 100; i++ {
		files += fmt.Sprintf("\n-- content/posts/p%d.md --\n"+content, i+1)
	}
	files = strings.ReplaceAll(files, "§", "%")

	cfg := IntegrationTestConfig{
		T:           b,
		TxtarString: files,
	}
	builders := make([]*IntegrationTestBuilder, b.N)

	for i := range builders {
		builders[i] = NewIntegrationTestBuilder(cfg)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		builders[i].Build()
	}
}

func TestReplaceShortcodeTokens(t *testing.T) {
	t.Parallel()
	for i, this := range []struct {
		input        string
		prefix       string
		replacements map[string]string
		expect       any
	}{
		{"Hello HAHAHUGOSHORTCODE-1HBHB.", "PREFIX", map[string]string{"HAHAHUGOSHORTCODE-1HBHB": "World"}, "Hello World."},
		{"Hello HAHAHUGOSHORTCODE-1@}@.", "PREFIX", map[string]string{"HAHAHUGOSHORTCODE-1HBHB": "World"}, false},
		{"HAHAHUGOSHORTCODE2-1HBHB", "PREFIX2", map[string]string{"HAHAHUGOSHORTCODE2-1HBHB": "World"}, "World"},
		{"Hello World!", "PREFIX2", map[string]string{}, "Hello World!"},
		{"!HAHAHUGOSHORTCODE-1HBHB", "PREFIX", map[string]string{"HAHAHUGOSHORTCODE-1HBHB": "World"}, "!World"},
		{"HAHAHUGOSHORTCODE-1HBHB!", "PREFIX", map[string]string{"HAHAHUGOSHORTCODE-1HBHB": "World"}, "World!"},
		{"!HAHAHUGOSHORTCODE-1HBHB!", "PREFIX", map[string]string{"HAHAHUGOSHORTCODE-1HBHB": "World"}, "!World!"},
		{"_{_PREFIX-1HBHB", "PREFIX", map[string]string{"HAHAHUGOSHORTCODE-1HBHB": "World"}, "_{_PREFIX-1HBHB"},
		{"Hello HAHAHUGOSHORTCODE-1HBHB.", "PREFIX", map[string]string{"HAHAHUGOSHORTCODE-1HBHB": "To You My Old Friend Who Told Me This Fantastic Story"}, "Hello To You My Old Friend Who Told Me This Fantastic Story."},
		{"A HAHAHUGOSHORTCODE-1HBHB asdf HAHAHUGOSHORTCODE-2HBHB.", "A", map[string]string{"HAHAHUGOSHORTCODE-1HBHB": "v1", "HAHAHUGOSHORTCODE-2HBHB": "v2"}, "A v1 asdf v2."},
		{"Hello HAHAHUGOSHORTCODE2-1HBHB. Go HAHAHUGOSHORTCODE2-2HBHB, Go, Go HAHAHUGOSHORTCODE2-3HBHB Go Go!.", "PREFIX2", map[string]string{"HAHAHUGOSHORTCODE2-1HBHB": "Europe", "HAHAHUGOSHORTCODE2-2HBHB": "Jonny", "HAHAHUGOSHORTCODE2-3HBHB": "Johnny"}, "Hello Europe. Go Jonny, Go, Go Johnny Go Go!."},
		{"A HAHAHUGOSHORTCODE-2HBHB HAHAHUGOSHORTCODE-1HBHB.", "PREFIX", map[string]string{"HAHAHUGOSHORTCODE-1HBHB": "A", "HAHAHUGOSHORTCODE-2HBHB": "B"}, "A B A."},
		{"A HAHAHUGOSHORTCODE-1HBHB HAHAHUGOSHORTCODE-2", "PREFIX", map[string]string{"HAHAHUGOSHORTCODE-1HBHB": "A"}, false},
		{"A HAHAHUGOSHORTCODE-1HBHB but not the second.", "PREFIX", map[string]string{"HAHAHUGOSHORTCODE-1HBHB": "A", "HAHAHUGOSHORTCODE-2HBHB": "B"}, "A A but not the second."},
		{"An HAHAHUGOSHORTCODE-1HBHB.", "PREFIX", map[string]string{"HAHAHUGOSHORTCODE-1HBHB": "A", "HAHAHUGOSHORTCODE-2HBHB": "B"}, "An A."},
		{"An HAHAHUGOSHORTCODE-1HBHB HAHAHUGOSHORTCODE-2HBHB.", "PREFIX", map[string]string{"HAHAHUGOSHORTCODE-1HBHB": "A", "HAHAHUGOSHORTCODE-2HBHB": "B"}, "An A B."},
		{"A HAHAHUGOSHORTCODE-1HBHB HAHAHUGOSHORTCODE-2HBHB HAHAHUGOSHORTCODE-3HBHB HAHAHUGOSHORTCODE-1HBHB HAHAHUGOSHORTCODE-3HBHB.", "PREFIX", map[string]string{"HAHAHUGOSHORTCODE-1HBHB": "A", "HAHAHUGOSHORTCODE-2HBHB": "B", "HAHAHUGOSHORTCODE-3HBHB": "C"}, "A A B C A C."},
		{"A HAHAHUGOSHORTCODE-1HBHB HAHAHUGOSHORTCODE-2HBHB HAHAHUGOSHORTCODE-3HBHB HAHAHUGOSHORTCODE-1HBHB HAHAHUGOSHORTCODE-3HBHB.", "PREFIX", map[string]string{"HAHAHUGOSHORTCODE-1HBHB": "A", "HAHAHUGOSHORTCODE-2HBHB": "B", "HAHAHUGOSHORTCODE-3HBHB": "C"}, "A A B C A C."},
		// Issue #1148 remove p-tags 10 =>
		{"Hello <p>HAHAHUGOSHORTCODE-1HBHB</p>. END.", "PREFIX", map[string]string{"HAHAHUGOSHORTCODE-1HBHB": "World"}, "Hello World. END."},
		{"Hello <p>HAHAHUGOSHORTCODE-1HBHB</p>. <p>HAHAHUGOSHORTCODE-2HBHB</p> END.", "PREFIX", map[string]string{"HAHAHUGOSHORTCODE-1HBHB": "World", "HAHAHUGOSHORTCODE-2HBHB": "THE"}, "Hello World. THE END."},
		{"Hello <p>HAHAHUGOSHORTCODE-1HBHB. END</p>.", "PREFIX", map[string]string{"HAHAHUGOSHORTCODE-1HBHB": "World"}, "Hello <p>World. END</p>."},
		{"<p>Hello HAHAHUGOSHORTCODE-1HBHB</p>. END.", "PREFIX", map[string]string{"HAHAHUGOSHORTCODE-1HBHB": "World"}, "<p>Hello World</p>. END."},
		{"Hello <p>HAHAHUGOSHORTCODE-1HBHB12", "PREFIX", map[string]string{"HAHAHUGOSHORTCODE-1HBHB": "World"}, "Hello <p>World12"},
		{
			"Hello HAHAHUGOSHORTCODE-1HBHB. HAHAHUGOSHORTCODE-1HBHB-HAHAHUGOSHORTCODE-1HBHB HAHAHUGOSHORTCODE-1HBHB HAHAHUGOSHORTCODE-1HBHB HAHAHUGOSHORTCODE-1HBHB END", "P",
			map[string]string{"HAHAHUGOSHORTCODE-1HBHB": strings.Repeat("BC", 100)},
			fmt.Sprintf("Hello %s. %s-%s %s %s %s END",
				strings.Repeat("BC", 100), strings.Repeat("BC", 100), strings.Repeat("BC", 100), strings.Repeat("BC", 100), strings.Repeat("BC", 100), strings.Repeat("BC", 100)),
		},
	} {

		replacements := make(map[string]shortcodeRenderer)
		for k, v := range this.replacements {
			replacements[k] = prerenderedShortcode{s: v}
		}
		tokenHandler := func(ctx context.Context, token string) ([]byte, error) {
			return []byte(this.replacements[token]), nil
		}

		ctx := context.Background()
		results, err := expandShortcodeTokens(ctx, []byte(this.input), tokenHandler)

		if b, ok := this.expect.(bool); ok && !b {
			if err == nil {
				t.Errorf("[%d] replaceShortcodeTokens didn't return an expected error", i)
			}
		} else {
			if err != nil {
				t.Errorf("[%d] failed: %s", i, err)
				continue
			}
			if !reflect.DeepEqual(results, []byte(this.expect.(string))) {
				t.Errorf("[%d] replaceShortcodeTokens, got \n%q but expected \n%q", i, results, this.expect)
			}
		}

	}
}

func TestShortcodeGetContent(t *testing.T) {
	t.Parallel()

	contentShortcode := `
{{- $t := .Get 0 -}}
{{- $p := .Get 1 -}}
{{- $k := .Get 2 -}}
{{- $page := $.Page.Site.GetPage "page" $p -}}
{{ if $page }}
{{- if eq $t "bundle" -}}
{{- .Scratch.Set "p" ($page.Resources.GetMatch (printf "%s*" $k)) -}}
{{- else -}}
{{- $.Scratch.Set "p" $page -}}
{{- end -}}P1:{{ .Page.Content }}|P2:{{ $p := ($.Scratch.Get "p") }}{{ $p.Title }}/{{ $p.Content }}|
{{- else -}}
{{- errorf "Page %s is nil" $p -}}
{{- end -}}
`

	var templates []string
	var content []string

	contentWithShortcodeTemplate := `---
title: doc%s
weight: %d
---
Logo:{{< c "bundle" "b1" "logo.png" >}}:P1: {{< c "page" "section1/p1" "" >}}:BP1:{{< c "bundle" "b1" "bp1" >}}`

	simpleContentTemplate := `---
title: doc%s
weight: %d
---
C-%s`

	templates = append(templates, []string{"shortcodes/c.html", contentShortcode}...)
	templates = append(templates, []string{"_default/single.html", "Single Content: {{ .Content }}"}...)
	templates = append(templates, []string{"_default/list.html", "List Content: {{ .Content }}"}...)

	content = append(content, []string{"b1/index.md", fmt.Sprintf(contentWithShortcodeTemplate, "b1", 1)}...)
	content = append(content, []string{"b1/logo.png", "PNG logo"}...)
	content = append(content, []string{"b1/bp1.md", fmt.Sprintf(simpleContentTemplate, "bp1", 1, "bp1")}...)

	content = append(content, []string{"section1/_index.md", fmt.Sprintf(contentWithShortcodeTemplate, "s1", 2)}...)
	content = append(content, []string{"section1/p1.md", fmt.Sprintf(simpleContentTemplate, "s1p1", 2, "s1p1")}...)

	content = append(content, []string{"section2/_index.md", fmt.Sprintf(simpleContentTemplate, "b1", 1, "b1")}...)
	content = append(content, []string{"section2/s2p1.md", fmt.Sprintf(contentWithShortcodeTemplate, "bp1", 1)}...)

	builder := newTestSitesBuilder(t).WithDefaultMultiSiteConfig()

	builder.WithContent(content...).WithTemplates(templates...).CreateSites().Build(BuildCfg{})
	s := builder.H.Sites[0]
	builder.Assert(len(s.RegularPages()), qt.Equals, 3)

	builder.AssertFileContent("public/en/section1/index.html",
		"List Content: <p>Logo:P1:|P2:logo.png/PNG logo|:P1: P1:|P2:docs1p1/<p>C-s1p1</p>\n|",
		"BP1:P1:|P2:docbp1/<p>C-bp1</p>",
	)

	builder.AssertFileContent("public/en/b1/index.html",
		"Single Content: <p>Logo:P1:|P2:logo.png/PNG logo|:P1: P1:|P2:docs1p1/<p>C-s1p1</p>\n|",
		"P2:docbp1/<p>C-bp1</p>",
	)

	builder.AssertFileContent("public/en/section2/s2p1/index.html",
		"Single Content: <p>Logo:P1:|P2:logo.png/PNG logo|:P1: P1:|P2:docs1p1/<p>C-s1p1</p>\n|",
		"P2:docbp1/<p>C-bp1</p>",
	)
}

// https://github.com/gohugoio/hugo/issues/5833
func TestShortcodeParentResourcesOnRebuild(t *testing.T) {
	t.Parallel()

	b := newTestSitesBuilder(t).Running().WithSimpleConfigFile()
	b.WithTemplatesAdded(
		"index.html", `
{{ $b := .Site.GetPage "b1" }}
b1 Content: {{ $b.Content }}
{{$p := $b.Resources.GetMatch "p1*" }}
Content: {{ $p.Content }}
{{ $article := .Site.GetPage "blog/article" }}
Article Content: {{ $article.Content }}
`,
		"shortcodes/c.html", `
{{ range .Page.Parent.Resources }}
* Parent resource: {{ .Name }}: {{ .RelPermalink }}
{{ end }}
`)

	pageContent := `
---
title: MyPage
---

SHORTCODE: {{< c >}}

`

	b.WithContent("b1/index.md", pageContent,
		"b1/logo.png", "PNG logo",
		"b1/p1.md", pageContent,
		"blog/_index.md", pageContent,
		"blog/logo-article.png", "PNG logo",
		"blog/article.md", pageContent,
	)

	b.Build(BuildCfg{})

	assert := func(matchers ...string) {
		allMatchers := append(matchers, "Parent resource: logo.png: /b1/logo.png",
			"Article Content: <p>SHORTCODE: \n\n* Parent resource: logo-article.png: /blog/logo-article.png",
		)

		b.AssertFileContent("public/index.html",
			allMatchers...,
		)
	}

	assert()

	b.EditFiles("content/b1/index.md", pageContent+" Edit.")

	b.Build(BuildCfg{})

	assert("Edit.")
}

func TestShortcodePreserveOrder(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	contentTemplate := `---
title: doc%d
weight: %d
---
# doc

{{< s1 >}}{{< s2 >}}{{< s3 >}}{{< s4 >}}{{< s5 >}}

{{< nested >}}
{{< ordinal >}} {{< scratch >}}
{{< ordinal >}} {{< scratch >}}
{{< ordinal >}} {{< scratch >}}
{{< /nested >}}

`

	ordinalShortcodeTemplate := `ordinal: {{ .Ordinal }}{{ .Page.Scratch.Set "ordinal" .Ordinal }}`

	nestedShortcode := `outer ordinal: {{ .Ordinal }} inner: {{ .Inner }}`
	scratchGetShortcode := `scratch ordinal: {{ .Ordinal }} scratch get ordinal: {{ .Page.Scratch.Get "ordinal" }}`
	shortcodeTemplate := `v%d: {{ .Ordinal }} sgo: {{ .Page.Scratch.Get "o2" }}{{ .Page.Scratch.Set "o2" .Ordinal }}|`

	var shortcodes []string
	var content []string

	shortcodes = append(shortcodes, []string{"shortcodes/nested.html", nestedShortcode}...)
	shortcodes = append(shortcodes, []string{"shortcodes/ordinal.html", ordinalShortcodeTemplate}...)
	shortcodes = append(shortcodes, []string{"shortcodes/scratch.html", scratchGetShortcode}...)

	for i := 1; i <= 5; i++ {
		sc := fmt.Sprintf(shortcodeTemplate, i)
		sc = strings.Replace(sc, "%%", "%", -1)
		shortcodes = append(shortcodes, []string{fmt.Sprintf("shortcodes/s%d.html", i), sc}...)
	}

	for i := 1; i <= 3; i++ {
		content = append(content, []string{fmt.Sprintf("p%d.md", i), fmt.Sprintf(contentTemplate, i, i)}...)
	}

	builder := newTestSitesBuilder(t).WithDefaultMultiSiteConfig()

	builder.WithContent(content...).WithTemplatesAdded(shortcodes...).CreateSites().Build(BuildCfg{})

	s := builder.H.Sites[0]
	c.Assert(len(s.RegularPages()), qt.Equals, 3)

	builder.AssertFileContent("public/en/p1/index.html", `v1: 0 sgo: |v2: 1 sgo: 0|v3: 2 sgo: 1|v4: 3 sgo: 2|v5: 4 sgo: 3`)
	builder.AssertFileContent("public/en/p1/index.html", `outer ordinal: 5 inner:
ordinal: 0 scratch ordinal: 1 scratch get ordinal: 0
ordinal: 2 scratch ordinal: 3 scratch get ordinal: 2
ordinal: 4 scratch ordinal: 5 scratch get ordinal: 4`)
}

func TestShortcodeVariables(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	builder := newTestSitesBuilder(t).WithSimpleConfigFile()

	builder.WithContent("page.md", `---
title: "Hugo Rocks!"
---

# doc

   {{< s1 >}}

`).WithTemplatesAdded("layouts/shortcodes/s1.html", `
Name: {{ .Name }}
{{ with .Position }}
File: {{ .Filename }}
Offset: {{ .Offset }}
Line: {{ .LineNumber }}
Column: {{ .ColumnNumber }}
String: {{ . | safeHTML }}
{{ end }}

`).CreateSites().Build(BuildCfg{})

	s := builder.H.Sites[0]
	c.Assert(len(s.RegularPages()), qt.Equals, 1)

	builder.AssertFileContent("public/page/index.html",
		filepath.FromSlash("File: content/page.md"),
		"Line: 7", "Column: 4", "Offset: 40",
		filepath.FromSlash("String: \"content/page.md:7:4\""),
		"Name: s1",
	)
}

func TestInlineShortcodes(t *testing.T) {
	for _, enableInlineShortcodes := range []bool{true, false} {
		enableInlineShortcodes := enableInlineShortcodes
		t.Run(fmt.Sprintf("enableInlineShortcodes=%t", enableInlineShortcodes),
			func(t *testing.T) {
				t.Parallel()
				conf := fmt.Sprintf(`
baseURL = "https://example.com"
enableInlineShortcodes = %t
`, enableInlineShortcodes)

				b := newTestSitesBuilder(t)
				b.WithConfigFile("toml", conf)

				shortcodeContent := `FIRST:{{< myshort.inline "first" >}}
Page: {{ .Page.Title }}
Seq: {{ seq 3 }}
Param: {{ .Get 0 }}
{{< /myshort.inline >}}:END:

SECOND:{{< myshort.inline "second" />}}:END
NEW INLINE:  {{< n1.inline "5" >}}W1: {{ seq (.Get 0) }}{{< /n1.inline >}}:END:
INLINE IN INNER: {{< outer >}}{{< n2.inline >}}W2: {{ seq 4 }}{{< /n2.inline >}}{{< /outer >}}:END:
REUSED INLINE IN INNER: {{< outer >}}{{< n1.inline "3" />}}{{< /outer >}}:END:
## MARKDOWN DELIMITER: {{% mymarkdown.inline %}}**Hugo Rocks!**{{% /mymarkdown.inline %}}
`

				b.WithContent("page-md-shortcode.md", `---
title: "Hugo"
---
`+shortcodeContent)

				b.WithContent("_index.md", `---
title: "Hugo Home"
---

`+shortcodeContent)

				b.WithTemplatesAdded("layouts/_default/single.html", `
CONTENT:{{ .Content }}
TOC: {{ .TableOfContents }}
`)

				b.WithTemplatesAdded("layouts/index.html", `
CONTENT:{{ .Content }}
TOC: {{ .TableOfContents }}
`)

				b.WithTemplatesAdded("layouts/shortcodes/outer.html", `Inner: {{ .Inner }}`)

				b.CreateSites().Build(BuildCfg{})

				shouldContain := []string{
					"Seq: [1 2 3]",
					"Param: first",
					"Param: second",
					"NEW INLINE:  W1: [1 2 3 4 5]",
					"INLINE IN INNER: Inner: W2: [1 2 3 4]",
					"REUSED INLINE IN INNER: Inner: W1: [1 2 3]",
					`<li><a href="#markdown-delimiter-hugo-rocks">MARKDOWN DELIMITER: <strong>Hugo Rocks!</strong></a></li>`,
				}

				if enableInlineShortcodes {
					b.AssertFileContent("public/page-md-shortcode/index.html",
						shouldContain...,
					)
					b.AssertFileContent("public/index.html",
						shouldContain...,
					)
				} else {
					b.AssertFileContent("public/page-md-shortcode/index.html",
						"FIRST::END",
						"SECOND::END",
						"NEW INLINE:  :END",
						"INLINE IN INNER: Inner: :END:",
						"REUSED INLINE IN INNER: Inner: :END:",
					)
				}
			})

	}
}

// https://github.com/gohugoio/hugo/issues/5863
func TestShortcodeNamespaced(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	builder := newTestSitesBuilder(t).WithSimpleConfigFile()

	builder.WithContent("page.md", `---
title: "Hugo Rocks!"
---

# doc

   hello: {{< hello >}}
   test/hello: {{< test/hello >}}

`).WithTemplatesAdded(
		"layouts/shortcodes/hello.html", `hello`,
		"layouts/shortcodes/test/hello.html", `test/hello`).CreateSites().Build(BuildCfg{})

	s := builder.H.Sites[0]
	c.Assert(len(s.RegularPages()), qt.Equals, 1)

	builder.AssertFileContent("public/page/index.html",
		"hello: hello",
		"test/hello: test/hello",
	)
}

func TestShortcodeParams(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.org"
-- layouts/shortcodes/hello.html --
{{ range $i, $v := .Params }}{{ printf "- %v: %v (%T) " $i $v $v -}}{{ end }}
-- content/page.md --
title: "Hugo Rocks!"
summary: "Foo"
---

# doc

types positional: {{< hello true false 33 3.14 >}}
types named: {{< hello b1=true b2=false i1=33 f1=3.14 >}}
types string: {{< hello "true" trues "33" "3.14" >}}
escaped quoute: {{< hello "hello \"world\"." >}}
-- layouts/_default/single.html --
Content: {{ .Content }}|
`

	b := Test(t, files)

	b.AssertFileContent("public/page/index.html",
		"types positional: - 0: true (bool) - 1: false (bool) - 2: 33 (int) - 3: 3.14 (float64)",
		"types named: - b1: true (bool) - b2: false (bool) - f1: 3.14 (float64) - i1: 33 (int)",
		"types string: - 0: true (string) - 1: trues (string) - 2: 33 (string) - 3: 3.14 (string) ",
		"hello &#34;world&#34;. (string)",
	)
}

func TestShortcodeRef(t *testing.T) {
	t.Parallel()

	v := config.New()
	v.Set("baseURL", "https://example.org")

	builder := newTestSitesBuilder(t).WithViper(v)

	for i := 1; i <= 2; i++ {
		builder.WithContent(fmt.Sprintf("page%d.md", i), `---
title: "Hugo Rocks!"
---



[Page 1]({{< ref "page1.md" >}})
[Page 1 with anchor]({{< relref "page1.md#doc" >}})
[Page 2]({{< ref "page2.md" >}})
[Page 2 with anchor]({{< relref "page2.md#doc" >}})


## Doc


`)
	}

	builder.Build(BuildCfg{})

	builder.AssertFileContent("public/page2/index.html", `
<a href="/page1/#doc">Page 1 with anchor</a>
<a href="https://example.org/page2/">Page 2</a>
<a href="/page2/#doc">Page 2 with anchor</a></p>

<h2 id="doc">Doc</h2>
`,
	)
}

// https://github.com/gohugoio/hugo/issues/6857
func TestShortcodeNoInner(t *testing.T) {
	t.Parallel()

	b := newTestSitesBuilder(t)

	b.WithContent("mypage.md", `---
title: "No Inner!"
---
{{< noinner >}}{{< /noinner >}}


`).WithTemplatesAdded(
		"layouts/shortcodes/noinner.html", `No inner here.`)

	err := b.BuildE(BuildCfg{})
	b.Assert(err.Error(), qt.Contains, filepath.FromSlash(`"content/mypage.md:4:16": failed to extract shortcode: shortcode "noinner" does not evaluate .Inner or .InnerDeindent, yet a closing tag was provided`))
}

func TestShortcodeStableOutputFormatTemplates(t *testing.T) {
	t.Parallel()

	for i := 0; i < 5; i++ {

		b := newTestSitesBuilder(t)

		const numPages = 10

		for i := 0; i < numPages; i++ {
			b.WithContent(fmt.Sprintf("page%d.md", i), `---
title: "Page"
outputs: ["html", "css", "csv", "json"]
---
{{< myshort >}}

`)
		}

		b.WithTemplates(
			"_default/single.html", "{{ .Content }}",
			"_default/single.css", "{{ .Content }}",
			"_default/single.csv", "{{ .Content }}",
			"_default/single.json", "{{ .Content }}",
			"shortcodes/myshort.html", `Short-HTML`,
			"shortcodes/myshort.csv", `Short-CSV`,
		)

		b.Build(BuildCfg{})

		// helpers.PrintFs(b.Fs.Destination, "public", os.Stdout)

		for i := 0; i < numPages; i++ {
			b.AssertFileContent(fmt.Sprintf("public/page%d/index.html", i), "Short-HTML")
			b.AssertFileContent(fmt.Sprintf("public/page%d/index.csv", i), "Short-CSV")
			b.AssertFileContent(fmt.Sprintf("public/page%d/index.json", i), "Short-HTML")

		}

		for i := 0; i < numPages; i++ {
			b.AssertFileContent(fmt.Sprintf("public/page%d/styles.css", i), "Short-HTML")
		}

	}
}

// #9821
func TestShortcodeMarkdownOutputFormat(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
-- content/p1.md --
---
title: "p1"
---
{{< foo >}}
# The below would have failed using the HTML template parser.
-- layouts/shortcodes/foo.md --
§§§
<x
§§§
-- layouts/_default/single.html --
{{ .Content }}
`

	b := Test(t, files)

	b.AssertFileContent("public/p1/index.html", `
<x
	`)
}

func TestShortcodePreserveIndentation(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
-- content/p1.md --
---
title: "p1"
---

## List With Indented Shortcodes

1. List 1
    {{% mark1 %}}
	1. Item Mark1 1
	1. Item Mark1 2
	{{% mark2 %}}
	{{% /mark1 %}}
-- layouts/shortcodes/mark1.md --
{{ .Inner }}
-- layouts/shortcodes/mark2.md --
1. Item Mark2 1
1. Item Mark2 2
   1. Item Mark2 2-1
1. Item Mark2 3
-- layouts/_default/single.html --
{{ .Content }}
`

	b := Test(t, files)

	b.AssertFileContent("public/p1/index.html", "<ol>\n<li>\n<p>List 1</p>\n<ol>\n<li>Item Mark1 1</li>\n<li>Item Mark1 2</li>\n<li>Item Mark2 1</li>\n<li>Item Mark2 2\n<ol>\n<li>Item Mark2 2-1</li>\n</ol>\n</li>\n<li>Item Mark2 3</li>\n</ol>\n</li>\n</ol>")
}

func TestShortcodeCodeblockIndent(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
-- content/p1.md --
---
title: "p1"
---

## Code block

    {{% code %}}

-- layouts/shortcodes/code.md --
echo "foo";
-- layouts/_default/single.html --
{{ .Content }}
`

	b := Test(t, files)

	b.AssertFileContent("public/p1/index.html", "<pre><code>echo &quot;foo&quot;;\n</code></pre>")
}

func TestShortcodeHighlightDeindent(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
[markup]
[markup.highlight]
codeFences = true
noClasses = false
-- content/p1.md --
---
title: "p1"
---

## Indent 5 Spaces

     {{< highlight bash >}}
     line 1;
     line 2;
     line 3;
     {{< /highlight >}}

-- layouts/_default/single.html --
{{ .Content }}
`

	b := Test(t, files)

	b.AssertFileContent("public/p1/index.html", `
<pre><code> <div class="highlight"><pre tabindex="0" class="chroma"><code class="language-bash" data-lang="bash"><span class="line"><span class="cl">line 1<span class="p">;</span>
</span></span><span class="line"><span class="cl">line 2<span class="p">;</span>
</span></span><span class="line"><span class="cl">line 3<span class="p">;</span></span></span></code></pre></div>
</code></pre>

	`)
}

// Issue 10236.
func TestShortcodeParamEscapedQuote(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
-- content/p1.md --
---
title: "p1"
---

{{< figure src="/media/spf13.jpg" title="Steve \"Francia\"." >}}

-- layouts/shortcodes/figure.html --
Title: {{ .Get "title" | safeHTML }}
-- layouts/_default/single.html --
{{ .Content }}
`

	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,

			Verbose: true,
		},
	).Build()

	b.AssertFileContent("public/p1/index.html", `Title: Steve "Francia".`)
}

// Issue 10391.
func TestNestedShortcodeCustomOutputFormat(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --

[outputFormats.Foobar]
baseName = "foobar"
isPlainText = true
mediaType = "application/json"
notAlternative = true

[languages.en]
languageName = "English"

[languages.en.outputs]
home = [ "HTML", "RSS", "Foobar" ]

[languages.fr]
languageName = "Français"

[[module.mounts]]
source = "content/en"
target = "content"
lang = "en"

[[module.mounts]]
source = "content/fr"
target = "content"
lang = "fr"

-- layouts/_default/list.foobar.json --
{{- $.Scratch.Add "data" slice -}}
{{- range (where .Site.AllPages "Kind" "!=" "home") -}}
	{{- $.Scratch.Add "data" (dict "content" (.Plain | truncate 10000) "type" .Type "full_url" .Permalink) -}}
{{- end -}}
{{- $.Scratch.Get "data" | jsonify -}}
-- content/en/p1.md --
---
title: "p1"
---

### More information

{{< tabs >}}
{{% tab "Test" %}}

It's a test

{{% /tab %}}
{{< /tabs >}}

-- content/fr/p2.md --
---
title: Test
---

### Plus d'informations

{{< tabs >}}
{{% tab "Test" %}}

C'est un test

{{% /tab %}}
{{< /tabs >}}

-- layouts/shortcodes/tabs.html --
<div>
  <div class="tab-content">{{ .Inner }}</div>
</div>

-- layouts/shortcodes/tab.html --
<div>{{ .Inner }}</div>

-- layouts/_default/single.html --
{{ .Content }}
`

	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,

			Verbose: true,
		},
	).Build()

	b.AssertFileContent("public/fr/p2/index.html", `plus-dinformations`)
}

// Issue 10671.
func TestShortcodeInnerShouldBeEmptyWhenNotClosed(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
disableKinds = ["home", "taxonomy", "term"]
-- content/p1.md --
---
title: "p1"
---

{{< sc "self-closing" />}}

Text.

{{< sc "closing-no-newline" >}}{{< /sc >}}

-- layouts/shortcodes/sc.html --
Inner: {{ .Get 0 }}: {{ len .Inner }}
InnerDeindent: {{ .Get 0 }}: {{ len .InnerDeindent }}
-- layouts/_default/single.html --
{{ .Content }}
`

	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,

			Verbose: true,
		},
	).Build()

	b.AssertFileContent("public/p1/index.html", `
Inner: self-closing: 0
InnerDeindent: self-closing: 0
Inner: closing-no-newline: 0
InnerDeindent: closing-no-newline: 0

`)
}

// Issue 10675.
func TestShortcodeErrorWhenItShouldBeClosed(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
disableKinds = ["home", "taxonomy", "term"]
-- content/p1.md --
---
title: "p1"
---

{{< sc >}}

Text.

-- layouts/shortcodes/sc.html --
Inner: {{ .Get 0 }}: {{ len .Inner }}
-- layouts/_default/single.html --
{{ .Content }}
`

	b, err := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,

			Verbose: true,
		},
	).BuildE()

	b.Assert(err, qt.Not(qt.IsNil))
	b.Assert(err.Error(), qt.Contains, `p1.md:5:1": failed to extract shortcode: shortcode "sc" must be closed or self-closed`)
}

// Issue 10819.
func TestShortcodeInCodeFenceHyphen(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
disableKinds = ["home", "taxonomy", "term"]
-- content/p1.md --
---
title: "p1"
---

§§§go
{{< sc >}}
§§§

Text.

-- layouts/shortcodes/sc.html --
Hello.
-- layouts/_default/single.html --
{{ .Content }}
`

	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,

			Verbose: true,
		},
	).Build()

	b.AssertFileContent("public/p1/index.html", "<span style=\"color:#a6e22e\">Hello.</span>")
}
