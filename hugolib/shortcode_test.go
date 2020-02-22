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
	"path/filepath"
	"reflect"

	"github.com/gohugoio/hugo/markup/asciidoc"
	"github.com/gohugoio/hugo/markup/rst"

	"github.com/spf13/viper"

	"github.com/gohugoio/hugo/parser/pageparser"
	"github.com/gohugoio/hugo/resources/page"

	"strings"
	"testing"

	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/tpl"
	"github.com/spf13/cast"

	qt "github.com/frankban/quicktest"
)

func CheckShortCodeMatch(t *testing.T, input, expected string, withTemplate func(templ tpl.TemplateManager) error) {
	t.Helper()
	CheckShortCodeMatchAndError(t, input, expected, withTemplate, false)
}

func CheckShortCodeMatchAndError(t *testing.T, input, expected string, withTemplate func(templ tpl.TemplateManager) error, expectError bool) {
	t.Helper()
	cfg, fs := newTestCfg()

	cfg.Set("markup", map[string]interface{}{
		"defaultMarkdownHandler": "blackfriday", // TODO(bep)
	})

	c := qt.New(t)

	// Need some front matter, see https://github.com/gohugoio/hugo/issues/2337
	contentFile := `---
title: "Title"
---
` + input

	writeSource(t, fs, "content/simple.md", contentFile)

	b := newTestSitesBuilderFromDepsCfg(t, deps.DepsCfg{Fs: fs, Cfg: cfg, WithTemplate: withTemplate}).WithNothingAdded()
	err := b.BuildE(BuildCfg{})

	if err != nil && !expectError {
		t.Fatalf("Shortcode rendered error %s.", err)
	}

	if err == nil && expectError {
		t.Fatalf("No error from shortcode")
	}

	h := b.H
	c.Assert(len(h.Sites), qt.Equals, 1)

	c.Assert(len(h.Sites[0].RegularPages()), qt.Equals, 1)

	output := strings.TrimSpace(content(h.Sites[0].RegularPages()[0]))
	output = strings.TrimPrefix(output, "<p>")
	output = strings.TrimSuffix(output, "</p>")

	expected = strings.TrimSpace(expected)

	if output != expected {
		t.Fatalf("Shortcode render didn't match. got \n%q but expected \n%q", output, expected)
	}
}

func TestNonSC(t *testing.T) {
	t.Parallel()
	// notice the syntax diff from 0.12, now comment delims must be added
	CheckShortCodeMatch(t, "{{%/* movie 47238zzb */%}}", "{{% movie 47238zzb %}}", nil)
}

// Issue #929
func TestHyphenatedSC(t *testing.T) {
	t.Parallel()
	wt := func(tem tpl.TemplateManager) error {

		tem.AddTemplate("_internal/shortcodes/hyphenated-video.html", `Playing Video {{ .Get 0 }}`)
		return nil
	}

	CheckShortCodeMatch(t, "{{< hyphenated-video 47238zzb >}}", "Playing Video 47238zzb", wt)
}

// Issue #1753
func TestNoTrailingNewline(t *testing.T) {
	t.Parallel()
	wt := func(tem tpl.TemplateManager) error {
		tem.AddTemplate("_internal/shortcodes/a.html", `{{ .Get 0 }}`)
		return nil
	}

	CheckShortCodeMatch(t, "ab{{< a c >}}d", "abcd", wt)
}

func TestPositionalParamSC(t *testing.T) {
	t.Parallel()
	wt := func(tem tpl.TemplateManager) error {
		tem.AddTemplate("_internal/shortcodes/video.html", `Playing Video {{ .Get 0 }}`)
		return nil
	}

	CheckShortCodeMatch(t, "{{< video 47238zzb >}}", "Playing Video 47238zzb", wt)
	CheckShortCodeMatch(t, "{{< video 47238zzb 132 >}}", "Playing Video 47238zzb", wt)
	CheckShortCodeMatch(t, "{{<video 47238zzb>}}", "Playing Video 47238zzb", wt)
	CheckShortCodeMatch(t, "{{<video 47238zzb    >}}", "Playing Video 47238zzb", wt)
	CheckShortCodeMatch(t, "{{<   video   47238zzb    >}}", "Playing Video 47238zzb", wt)
}

func TestPositionalParamIndexOutOfBounds(t *testing.T) {
	t.Parallel()
	wt := func(tem tpl.TemplateManager) error {
		tem.AddTemplate("_internal/shortcodes/video.html", `Playing Video {{ with .Get 1 }}{{ . }}{{ else }}Missing{{ end }}`)
		return nil
	}
	CheckShortCodeMatch(t, "{{< video 47238zzb >}}", "Playing Video Missing", wt)
}

// #5071
func TestShortcodeRelated(t *testing.T) {
	t.Parallel()
	wt := func(tem tpl.TemplateManager) error {
		tem.AddTemplate("_internal/shortcodes/a.html", `{{ len (.Site.RegularPages.Related .Page) }}`)
		return nil
	}

	CheckShortCodeMatch(t, "{{< a >}}", "0", wt)
}

func TestShortcodeInnerMarkup(t *testing.T) {
	t.Parallel()
	wt := func(tem tpl.TemplateManager) error {
		tem.AddTemplate("shortcodes/a.html", `<div>{{ .Inner }}</div>`)
		tem.AddTemplate("shortcodes/b.html", `**Bold**: <div>{{ .Inner }}</div>`)
		return nil
	}

	CheckShortCodeMatch(t,
		"{{< a >}}B: <div>{{% b %}}**Bold**{{% /b %}}</div>{{< /a >}}",
		// This assertion looks odd, but is correct: for inner shortcodes with
		// the {{% we treats the .Inner content as markup, but not the shortcode
		// itself.
		"<div>B: <div>**Bold**: <div><strong>Bold</strong></div></div></div>",
		wt)

	CheckShortCodeMatch(t,
		"{{% b %}}This is **B**: {{< b >}}This is B{{< /b>}}{{% /b %}}",
		"<strong>Bold</strong>: <div>This is <strong>B</strong>: <strong>Bold</strong>: <div>This is B</div></div>",
		wt)
}

// some repro issues for panics in Go Fuzz testing

func TestNamedParamSC(t *testing.T) {
	t.Parallel()
	wt := func(tem tpl.TemplateManager) error {
		tem.AddTemplate("_internal/shortcodes/img.html", `<img{{ with .Get "src" }} src="{{.}}"{{end}}{{with .Get "class"}} class="{{.}}"{{end}}>`)
		return nil
	}
	CheckShortCodeMatch(t, `{{< img src="one" >}}`, `<img src="one">`, wt)
	CheckShortCodeMatch(t, `{{< img class="aspen" >}}`, `<img class="aspen">`, wt)
	CheckShortCodeMatch(t, `{{< img src= "one" >}}`, `<img src="one">`, wt)
	CheckShortCodeMatch(t, `{{< img src ="one" >}}`, `<img src="one">`, wt)
	CheckShortCodeMatch(t, `{{< img src = "one" >}}`, `<img src="one">`, wt)
	CheckShortCodeMatch(t, `{{< img src = "one" class = "aspen grove" >}}`, `<img src="one" class="aspen grove">`, wt)
}

// Issue #2294
func TestNestedNamedMissingParam(t *testing.T) {
	t.Parallel()
	wt := func(tem tpl.TemplateManager) error {
		tem.AddTemplate("_internal/shortcodes/acc.html", `<div class="acc">{{ .Inner }}</div>`)
		tem.AddTemplate("_internal/shortcodes/div.html", `<div {{with .Get "class"}} class="{{ . }}"{{ end }}>{{ .Inner }}</div>`)
		tem.AddTemplate("_internal/shortcodes/div2.html", `<div {{with .Get 0}} class="{{ . }}"{{ end }}>{{ .Inner }}</div>`)
		return nil
	}
	CheckShortCodeMatch(t,
		`{{% acc %}}{{% div %}}d1{{% /div %}}{{% div2 %}}d2{{% /div2 %}}{{% /acc %}}`,
		"<div class=\"acc\"><div >d1</div><div >d2</div></div>", wt)
}

func TestIsNamedParamsSC(t *testing.T) {
	t.Parallel()
	wt := func(tem tpl.TemplateManager) error {
		tem.AddTemplate("_internal/shortcodes/bynameorposition.html", `{{ with .Get "id" }}Named: {{ . }}{{ else }}Pos: {{ .Get 0 }}{{ end }}`)
		tem.AddTemplate("_internal/shortcodes/ifnamedparams.html", `<div id="{{ if .IsNamedParams }}{{ .Get "id" }}{{ else }}{{ .Get 0 }}{{end}}">`)
		return nil
	}
	CheckShortCodeMatch(t, `{{< ifnamedparams id="name" >}}`, `<div id="name">`, wt)
	CheckShortCodeMatch(t, `{{< ifnamedparams position >}}`, `<div id="position">`, wt)
	CheckShortCodeMatch(t, `{{< bynameorposition id="name" >}}`, `Named: name`, wt)
	CheckShortCodeMatch(t, `{{< bynameorposition position >}}`, `Pos: position`, wt)
}

func TestInnerSC(t *testing.T) {
	t.Parallel()
	wt := func(tem tpl.TemplateManager) error {
		tem.AddTemplate("_internal/shortcodes/inside.html", `<div{{with .Get "class"}} class="{{.}}"{{end}}>{{ .Inner }}</div>`)
		return nil
	}
	CheckShortCodeMatch(t, `{{< inside class="aspen" >}}`, `<div class="aspen"></div>`, wt)
	CheckShortCodeMatch(t, `{{< inside class="aspen" >}}More Here{{< /inside >}}`, "<div class=\"aspen\">More Here</div>", wt)
	CheckShortCodeMatch(t, `{{< inside >}}More Here{{< /inside >}}`, "<div>More Here</div>", wt)
}

func TestInnerSCWithMarkdown(t *testing.T) {
	t.Parallel()
	wt := func(tem tpl.TemplateManager) error {
		// Note: In Hugo 0.55 we made it so any outer {{%'s inner content was rendered as part of the surrounding
		// markup. This solved lots of problems, but it also meant that this test had to be adjusted.
		tem.AddTemplate("_internal/shortcodes/wrapper.html", `<div{{with .Get "class"}} class="{{.}}"{{end}}>{{ .Inner }}</div>`)
		tem.AddTemplate("_internal/shortcodes/inside.html", `{{ .Inner }}`)
		return nil
	}
	CheckShortCodeMatch(t, `{{< wrapper >}}{{% inside %}}
# More Here

[link](http://spf13.com) and text

{{% /inside %}}{{< /wrapper >}}`, "<div><h1 id=\"more-here\">More Here</h1>\n\n<p><a href=\"http://spf13.com\">link</a> and text</p>\n</div>", wt)
}

func TestEmbeddedSC(t *testing.T) {
	t.Parallel()
	CheckShortCodeMatch(t, `{{% figure src="/found/here" class="bananas orange" %}}`, "<figure class=\"bananas orange\">\n    <img src=\"/found/here\"/> \n</figure>", nil)
	CheckShortCodeMatch(t, `{{% figure src="/found/here" class="bananas orange" caption="This is a caption" %}}`, "<figure class=\"bananas orange\">\n    <img src=\"/found/here\"\n         alt=\"This is a caption\"/> <figcaption>\n            <p>This is a caption</p>\n        </figcaption>\n</figure>", nil)
}

func TestNestedSC(t *testing.T) {
	t.Parallel()
	wt := func(tem tpl.TemplateManager) error {
		tem.AddTemplate("_internal/shortcodes/scn1.html", `<div>Outer, inner is {{ .Inner }}</div>`)
		tem.AddTemplate("_internal/shortcodes/scn2.html", `<div>SC2</div>`)
		return nil
	}
	CheckShortCodeMatch(t, `{{% scn1 %}}{{% scn2 %}}{{% /scn1 %}}`, "<div>Outer, inner is <div>SC2</div></div>", wt)

	CheckShortCodeMatch(t, `{{< scn1 >}}{{% scn2 %}}{{< /scn1 >}}`, "<div>Outer, inner is <div>SC2</div></div>", wt)
}

func TestNestedComplexSC(t *testing.T) {
	t.Parallel()
	wt := func(tem tpl.TemplateManager) error {
		tem.AddTemplate("_internal/shortcodes/row.html", `-row-{{ .Inner}}-rowStop-`)
		tem.AddTemplate("_internal/shortcodes/column.html", `-col-{{.Inner    }}-colStop-`)
		tem.AddTemplate("_internal/shortcodes/aside.html", `-aside-{{    .Inner  }}-asideStop-`)
		return nil
	}
	CheckShortCodeMatch(t, `{{< row >}}1-s{{% column %}}2-**s**{{< aside >}}3-**s**{{< /aside >}}4-s{{% /column %}}5-s{{< /row >}}6-s`,
		"-row-1-s-col-2-<strong>s</strong>-aside-3-<strong>s</strong>-asideStop-4-s-colStop-5-s-rowStop-6-s", wt)

	// turn around the markup flag
	CheckShortCodeMatch(t, `{{% row %}}1-s{{< column >}}2-**s**{{% aside %}}3-**s**{{% /aside %}}4-s{{< /column >}}5-s{{% /row %}}6-s`,
		"-row-1-s-col-2-<strong>s</strong>-aside-3-<strong>s</strong>-asideStop-4-s-colStop-5-s-rowStop-6-s", wt)
}

func TestParentShortcode(t *testing.T) {
	t.Parallel()
	wt := func(tem tpl.TemplateManager) error {
		tem.AddTemplate("_internal/shortcodes/r1.html", `1: {{ .Get "pr1" }} {{ .Inner }}`)
		tem.AddTemplate("_internal/shortcodes/r2.html", `2: {{ .Parent.Get "pr1" }}{{ .Get "pr2" }} {{ .Inner }}`)
		tem.AddTemplate("_internal/shortcodes/r3.html", `3: {{ .Parent.Parent.Get "pr1" }}{{ .Parent.Get "pr2" }}{{ .Get "pr3" }} {{ .Inner }}`)
		return nil
	}
	CheckShortCodeMatch(t, `{{< r1 pr1="p1" >}}1: {{< r2 pr2="p2" >}}2: {{< r3 pr3="p3" >}}{{< /r3 >}}{{< /r2 >}}{{< /r1 >}}`,
		"1: p1 1: 2: p1p2 2: 3: p1p2p3 ", wt)

}

func TestFigureOnlySrc(t *testing.T) {
	t.Parallel()
	CheckShortCodeMatch(t, `{{< figure src="/found/here" >}}`, "<figure>\n    <img src=\"/found/here\"/> \n</figure>", nil)
}

func TestFigureCaptionAttrWithMarkdown(t *testing.T) {
	t.Parallel()
	CheckShortCodeMatch(t, `{{< figure src="/found/here" caption="Something **bold** _italic_" >}}`, "<figure>\n    <img src=\"/found/here\"\n         alt=\"Something bold italic\"/> <figcaption>\n            <p>Something <strong>bold</strong> <em>italic</em></p>\n        </figcaption>\n</figure>", nil)
	CheckShortCodeMatch(t, `{{< figure src="/found/here" attr="Something **bold** _italic_" >}}`, "<figure>\n    <img src=\"/found/here\"/> <figcaption>\n            <p>Something <strong>bold</strong> <em>italic</em></p>\n        </figcaption>\n</figure>", nil)
}

func TestFigureImgWidth(t *testing.T) {
	t.Parallel()
	CheckShortCodeMatch(t, `{{% figure src="/found/here" class="bananas orange" alt="apple" width="100px" %}}`, "<figure class=\"bananas orange\">\n    <img src=\"/found/here\"\n         alt=\"apple\" width=\"100px\"/> \n</figure>", nil)
}

func TestFigureImgHeight(t *testing.T) {
	t.Parallel()
	CheckShortCodeMatch(t, `{{% figure src="/found/here" class="bananas orange" alt="apple" height="100px" %}}`, "<figure class=\"bananas orange\">\n    <img src=\"/found/here\"\n         alt=\"apple\" height=\"100px\"/> \n</figure>", nil)
}

func TestFigureImgWidthAndHeight(t *testing.T) {
	t.Parallel()
	CheckShortCodeMatch(t, `{{% figure src="/found/here" class="bananas orange" alt="apple" width="50" height="100" %}}`, "<figure class=\"bananas orange\">\n    <img src=\"/found/here\"\n         alt=\"apple\" width=\"50\" height=\"100\"/> \n</figure>", nil)
}

func TestFigureLinkNoTarget(t *testing.T) {
	t.Parallel()
	CheckShortCodeMatch(t, `{{< figure src="/found/here" link="/jump/here/on/clicking" >}}`, "<figure><a href=\"/jump/here/on/clicking\">\n    <img src=\"/found/here\"/> </a>\n</figure>", nil)
}

func TestFigureLinkWithTarget(t *testing.T) {
	t.Parallel()
	CheckShortCodeMatch(t, `{{< figure src="/found/here" link="/jump/here/on/clicking" target="_self" >}}`, "<figure><a href=\"/jump/here/on/clicking\" target=\"_self\">\n    <img src=\"/found/here\"/> </a>\n</figure>", nil)
}

func TestFigureLinkWithTargetAndRel(t *testing.T) {
	t.Parallel()
	CheckShortCodeMatch(t, `{{< figure src="/found/here" link="/jump/here/on/clicking" target="_blank" rel="noopener" >}}`, "<figure><a href=\"/jump/here/on/clicking\" target=\"_blank\" rel=\"noopener\">\n    <img src=\"/found/here\"/> </a>\n</figure>", nil)
}

// #1642
func TestShortcodeWrappedInPIssue(t *testing.T) {
	t.Parallel()
	wt := func(tem tpl.TemplateManager) error {
		tem.AddTemplate("_internal/shortcodes/bug.html", `xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx`)
		return nil
	}
	CheckShortCodeMatch(t, `
{{< bug >}}

{{< bug >}}
`, "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx\n\nxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", wt)
}

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

	/*errCheck := func(s string) func(name string, assert *require.Assertions, shortcode *shortcode, err error) {
		return func(name string, assert *require.Assertions, shortcode *shortcode, err error) {
			c.Assert(err, name, qt.Not(qt.IsNil))
			c.Assert(err.Error(), name, qt.Equals, s)
		}
	}*/

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
		{"nested inner", `{{< inner >}}Inner Content->{{% inner2 param1 %}}inner2txt{{% /inner2 %}}Inner close->{{< / inner >}}`,
			regexpCheck("inner;.*inner:{Inner Content->.*Inner close->}")},
		{"nested, nested inner", `{{< inner >}}inner2->{{% inner2 param1 %}}inner2txt->inner3{{< inner3>}}inner3txt{{</ inner3 >}}{{% /inner2 %}}final close->{{< / inner >}}`,
			regexpCheck("inner:{inner2-> inner2.*{{inner2txt->inner3.*final close->}")},
		{"closed without content", `{{< inner param1 >}}{{< / inner >}}`, regexpCheck("inner.*inner:{}")},
		{"inline", `{{< my.inline >}}Hi{{< /my.inline >}}`, regexpCheck("my.inline;inline:true;closing:true;inner:{Hi};")},
	} {

		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			c := qt.New(t)

			counter := 0
			placeholderFunc := func() string {
				counter++
				return fmt.Sprintf("HAHA%s-%dHBHB", shortcodePlaceholderPrefix, counter)
			}

			p, err := pageparser.ParseMain(strings.NewReader(test.input), pageparser.Config{})
			c.Assert(err, qt.IsNil)
			handler := newShortcodeHandler(nil, s, placeholderFunc)
			iter := p.Iterator()

			short, err := handler.extractShortcode(0, 0, iter)

			test.check(c, short, err)

		})
	}

}

func TestShortcodesInSite(t *testing.T) {
	baseURL := "http://foo/bar"

	tests := []struct {
		contentPath string
		content     string
		outFile     string
		expected    interface{}
	}{
		{"sect/doc1.md", `a{{< b >}}c`,
			filepath.FromSlash("public/sect/doc1/index.html"), "<p>abc</p>\n"},
		// Issue #1642: Multiple shortcodes wrapped in P
		// Deliberately forced to pass even if they maybe shouldn't.
		{"sect/doc2.md", `a

{{< b >}}		
{{< c >}}
{{< d >}}

e`,
			filepath.FromSlash("public/sect/doc2/index.html"),
			"<p>a</p>\n\n<p>b<br />\nc\nd</p>\n\n<p>e</p>\n"},
		{"sect/doc3.md", `a

{{< b >}}		
{{< c >}}

{{< d >}}

e`,
			filepath.FromSlash("public/sect/doc3/index.html"),
			"<p>a</p>\n\n<p>b<br />\nc</p>\n\nd\n\n<p>e</p>\n"},
		{"sect/doc4.md", `a
{{< b >}}
{{< b >}}
{{< b >}}
{{< b >}}
{{< b >}}










`,
			filepath.FromSlash("public/sect/doc4/index.html"),
			"<p>a\nb\nb\nb\nb\nb</p>\n"},
		// #2192 #2209: Shortcodes in markdown headers
		{"sect/doc5.md", `# {{< b >}}	
## {{% c %}}`,
			filepath.FromSlash("public/sect/doc5/index.html"), `-hbhb">b</h1>`},
		// #2223 pygments
		{"sect/doc6.md", "\n```bash\nb = {{< b >}} c = {{% c %}}\n```\n",
			filepath.FromSlash("public/sect/doc6/index.html"),
			`<span class="nv">b</span>`},
		// #2249
		{"sect/doc7.ad", `_Shortcodes:_ *b: {{< b >}} c: {{% c %}}*`,
			filepath.FromSlash("public/sect/doc7/index.html"),
			"<div class=\"paragraph\">\n<p><em>Shortcodes:</em> <strong>b: b c: c</strong></p>\n</div>\n"},
		{"sect/doc8.rst", `**Shortcodes:** *b: {{< b >}} c: {{% c %}}*`,
			filepath.FromSlash("public/sect/doc8/index.html"),
			"<div class=\"document\">\n\n\n<p><strong>Shortcodes:</strong> <em>b: b c: c</em></p>\n</div>"},
		{"sect/doc9.mmark", `
---
menu:
  main:
    parent: 'parent'
---
**Shortcodes:** *b: {{< b >}} c: {{% c %}}*`,
			filepath.FromSlash("public/sect/doc9/index.html"),
			"<p><strong>Shortcodes:</strong> <em>b: b c: c</em></p>\n"},
		// Issue #1229: Menus not available in shortcode.
		{"sect/doc10.md", `---
menu:
  main:
    identifier: 'parent'
tags:
- Menu
---
**Menus:** {{< menu >}}`,
			filepath.FromSlash("public/sect/doc10/index.html"),
			"<p><strong>Menus:</strong> 1</p>\n"},
		// Issue #2323: Taxonomies not available in shortcode.
		{"sect/doc11.md", `---
tags:
- Bugs
---
**Tags:** {{< tags >}}`,
			filepath.FromSlash("public/sect/doc11/index.html"),
			"<p><strong>Tags:</strong> 2</p>\n"},
		{"sect/doc12.md", `---
title: "Foo"
---

{{% html-indented-v1 %}}`,
			"public/sect/doc12/index.html",
			"<h1>Hugo!</h1>"},
	}

	temp := tests[:0]
	for _, test := range tests {
		if strings.HasSuffix(test.contentPath, ".ad") && !asciidoc.Supports() {
			t.Log("Skip Asciidoc test case as no Asciidoc present.")
			continue
		} else if strings.HasSuffix(test.contentPath, ".rst") && !rst.Supports() {
			t.Log("Skip Rst test case as no rst2html present.")
			continue
		}
		temp = append(temp, test)
	}
	tests = temp

	sources := make([][2]string, len(tests))

	for i, test := range tests {
		sources[i] = [2]string{filepath.FromSlash(test.contentPath), test.content}
	}

	addTemplates := func(templ tpl.TemplateManager) error {
		templ.AddTemplate("_default/single.html", "{{.Content}} Word Count: {{ .WordCount }}")

		templ.AddTemplate("_internal/shortcodes/b.html", `b`)
		templ.AddTemplate("_internal/shortcodes/c.html", `c`)
		templ.AddTemplate("_internal/shortcodes/d.html", `d`)
		templ.AddTemplate("_internal/shortcodes/html-indented-v1.html", "{{ $_hugo_config := `{ \"version\": 1 }` }}"+`
    <h1>Hugo!</h1>
`)
		templ.AddTemplate("_internal/shortcodes/menu.html", `{{ len (index .Page.Menus "main").Children }}`)
		templ.AddTemplate("_internal/shortcodes/tags.html", `{{ len .Page.Site.Taxonomies.tags }}`)

		return nil

	}

	cfg, fs := newTestCfg()

	cfg.Set("defaultContentLanguage", "en")
	cfg.Set("baseURL", baseURL)
	cfg.Set("uglyURLs", false)
	cfg.Set("verbose", true)

	cfg.Set("pygmentsUseClasses", true)
	cfg.Set("pygmentsCodefences", true)
	cfg.Set("markup", map[string]interface{}{
		"defaultMarkdownHandler": "blackfriday", // TODO(bep)
	})

	writeSourcesToSource(t, "content", fs, sources...)

	s := buildSingleSite(t, deps.DepsCfg{WithTemplate: addTemplates, Fs: fs, Cfg: cfg}, BuildCfg{})

	for i, test := range tests {
		test := test
		t.Run(fmt.Sprintf("test=%d;contentPath=%s", i, test.contentPath), func(t *testing.T) {
			t.Parallel()

			th := newTestHelper(s.Cfg, s.Fs, t)

			expected := cast.ToStringSlice(test.expected)

			th.assertFileContent(filepath.FromSlash(test.outFile), expected...)
		})

	}

}

func TestShortcodeMultipleOutputFormats(t *testing.T) {
	t.Parallel()

	siteConfig := `
baseURL = "http://example.com/blog"

paginate = 1

disableKinds = ["section", "taxonomy", "taxonomyTerm", "RSS", "sitemap", "robotsTXT", "404"]

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
	home := s.getPage(page.KindHome)
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
		replacements map[string]string
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

	var in = make([]input, b.N*len(data))
	var cnt = 0
	for i := 0; i < b.N; i++ {
		for _, this := range data {
			in[cnt] = input{[]byte(this.input), this.replacements, this.expect}
			cnt++
		}
	}

	b.ResetTimer()
	cnt = 0
	for i := 0; i < b.N; i++ {
		for j := range data {
			currIn := in[cnt]
			cnt++
			results, err := replaceShortcodeTokens(currIn.in, currIn.replacements)

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

func TestReplaceShortcodeTokens(t *testing.T) {
	t.Parallel()
	for i, this := range []struct {
		input        string
		prefix       string
		replacements map[string]string
		expect       interface{}
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
		{"Hello HAHAHUGOSHORTCODE-1HBHB. HAHAHUGOSHORTCODE-1HBHB-HAHAHUGOSHORTCODE-1HBHB HAHAHUGOSHORTCODE-1HBHB HAHAHUGOSHORTCODE-1HBHB HAHAHUGOSHORTCODE-1HBHB END", "P", map[string]string{"HAHAHUGOSHORTCODE-1HBHB": strings.Repeat("BC", 100)},
			fmt.Sprintf("Hello %s. %s-%s %s %s %s END",
				strings.Repeat("BC", 100), strings.Repeat("BC", 100), strings.Repeat("BC", 100), strings.Repeat("BC", 100), strings.Repeat("BC", 100), strings.Repeat("BC", 100))},
	} {

		results, err := replaceShortcodeTokens([]byte(this.input), this.replacements)

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

// https://github.com/gohugoio/hugo/issues/6504
func TestShortcodeEmoji(t *testing.T) {
	t.Parallel()

	v := viper.New()
	v.Set("enableEmoji", true)

	builder := newTestSitesBuilder(t).WithViper(v)

	builder.WithContent("page.md", `---
title: "Hugo Rocks!"
---

# doc

{{< event >}}10:30-11:00 My :smile: Event {{< /event >}}


`).WithTemplatesAdded(
		"layouts/shortcodes/event.html", `<div>{{ "\u29BE" }} {{ .Inner }} </div>`)

	builder.Build(BuildCfg{})
	builder.AssertFileContent("public/page/index.html",
		"⦾ 10:30-11:00 My 😄 Event",
	)
}

func TestShortcodeTypedParams(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	builder := newTestSitesBuilder(t).WithSimpleConfigFile()

	builder.WithContent("page.md", `---
title: "Hugo Rocks!"
---

# doc

types positional: {{< hello true false 33 3.14 >}}
types named: {{< hello b1=true b2=false i1=33 f1=3.14 >}}
types string: {{< hello "true" trues "33" "3.14" >}}


`).WithTemplatesAdded(
		"layouts/shortcodes/hello.html",
		`{{ range $i, $v := .Params }}
-  {{ printf "%v: %v (%T)" $i $v $v }}
{{ end }}
{{ $b1 := .Get "b1" }}
Get: {{ printf "%v (%T)" $b1 $b1 | safeHTML }}
`).Build(BuildCfg{})

	s := builder.H.Sites[0]
	c.Assert(len(s.RegularPages()), qt.Equals, 1)

	builder.AssertFileContent("public/page/index.html",
		"types positional: - 0: true (bool) - 1: false (bool) - 2: 33 (int) - 3: 3.14 (float64)",
		"types named: - b1: true (bool) - b2: false (bool) - f1: 3.14 (float64) - i1: 33 (int) Get: true (bool) ",
		"types string: - 0: true (string) - 1: trues (string) - 2: 33 (string) - 3: 3.14 (string) ",
	)
}

func TestShortcodeRef(t *testing.T) {
	for _, plainIDAnchors := range []bool{false, true} {
		plainIDAnchors := plainIDAnchors
		t.Run(fmt.Sprintf("plainIDAnchors=%t", plainIDAnchors), func(t *testing.T) {
			t.Parallel()

			v := viper.New()
			v.Set("baseURL", "https://example.org")
			v.Set("blackfriday", map[string]interface{}{
				"plainIDAnchors": plainIDAnchors,
			})
			v.Set("markup", map[string]interface{}{
				"defaultMarkdownHandler": "blackfriday", // TODO(bep)
			})

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

			if plainIDAnchors {
				builder.AssertFileContent("public/page2/index.html",
					`
<a href="/page1/#doc">Page 1 with anchor</a>
<a href="https://example.org/page2/">Page 2</a>
<a href="/page2/#doc">Page 2 with anchor</a></p>

<h2 id="doc">Doc</h2>
`,
				)
			} else {
				builder.AssertFileContent("public/page2/index.html",
					`
<p><a href="https://example.org/page1/">Page 1</a>
<a href="/page1/#doc:45ca767ba77bc1445a0acab74f80812f">Page 1 with anchor</a>
<a href="https://example.org/page2/">Page 2</a>
<a href="/page2/#doc:8e3cdf52fa21e33270c99433820e46bd">Page 2 with anchor</a></p>
<h2 id="doc:8e3cdf52fa21e33270c99433820e46bd">Doc</h2>
`,
				)
			}

		})
	}

}
