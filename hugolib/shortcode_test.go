// Copyright 2016 The Hugo Authors. All rights reserved.
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
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/spf13/viper"

	"github.com/spf13/afero"

	"github.com/gohugoio/hugo/output"

	"github.com/gohugoio/hugo/media"

	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/tpl"

	"github.com/stretchr/testify/require"
)

// TODO(bep) remove
func pageFromString(in, filename string, shortcodePlaceholderFn func() string, withTemplate ...func(templ tpl.TemplateHandler) error) (*Page, error) {
	var err error
	cfg, fs := newTestCfg()

	d := deps.DepsCfg{Cfg: cfg, Fs: fs, WithTemplate: withTemplate[0]}

	s, err := NewSiteForCfg(d)
	if err != nil {
		return nil, err
	}

	s.shortcodePlaceholderFunc = shortcodePlaceholderFn

	return s.newPageFrom(strings.NewReader(in), filename)
}

func CheckShortCodeMatch(t *testing.T, input, expected string, withTemplate func(templ tpl.TemplateHandler) error) {
	CheckShortCodeMatchAndError(t, input, expected, withTemplate, false)
}

func CheckShortCodeMatchAndError(t *testing.T, input, expected string, withTemplate func(templ tpl.TemplateHandler) error, expectError bool) {

	cfg, fs := newTestCfg()

	// Need some front matter, see https://github.com/gohugoio/hugo/issues/2337
	contentFile := `---
title: "Title"
---
` + input

	writeSource(t, fs, "content/simple.md", contentFile)

	h, err := NewHugoSites(deps.DepsCfg{Fs: fs, Cfg: cfg, WithTemplate: withTemplate})

	require.NoError(t, err)
	require.Len(t, h.Sites, 1)

	err = h.Build(BuildCfg{})

	if err != nil && !expectError {
		t.Fatalf("Shortcode rendered error %s.", err)
	}

	if err == nil && expectError {
		t.Fatalf("No error from shortcode")
	}

	require.Len(t, h.Sites[0].RegularPages, 1)

	output := strings.TrimSpace(string(h.Sites[0].RegularPages[0].content()))
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
	wt := func(tem tpl.TemplateHandler) error {

		tem.AddTemplate("_internal/shortcodes/hyphenated-video.html", `Playing Video {{ .Get 0 }}`)
		return nil
	}

	CheckShortCodeMatch(t, "{{< hyphenated-video 47238zzb >}}", "Playing Video 47238zzb", wt)
}

// Issue #1753
func TestNoTrailingNewline(t *testing.T) {
	t.Parallel()
	wt := func(tem tpl.TemplateHandler) error {
		tem.AddTemplate("_internal/shortcodes/a.html", `{{ .Get 0 }}`)
		return nil
	}

	CheckShortCodeMatch(t, "ab{{< a c >}}d", "abcd", wt)
}

func TestPositionalParamSC(t *testing.T) {
	t.Parallel()
	wt := func(tem tpl.TemplateHandler) error {
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
	wt := func(tem tpl.TemplateHandler) error {
		tem.AddTemplate("_internal/shortcodes/video.html", `Playing Video {{ with .Get 1 }}{{ . }}{{ else }}Missing{{ end }}`)
		return nil
	}
	CheckShortCodeMatch(t, "{{< video 47238zzb >}}", "Playing Video Missing", wt)
}

// #5071
func TestShortcodeRelated(t *testing.T) {
	t.Parallel()
	wt := func(tem tpl.TemplateHandler) error {
		tem.AddTemplate("_internal/shortcodes/a.html", `{{ len (.Site.RegularPages.Related .Page) }}`)
		return nil
	}

	CheckShortCodeMatch(t, "{{< a >}}", "0", wt)
}

// some repro issues for panics in Go Fuzz testing

func TestNamedParamSC(t *testing.T) {
	t.Parallel()
	wt := func(tem tpl.TemplateHandler) error {
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
	wt := func(tem tpl.TemplateHandler) error {
		tem.AddTemplate("_internal/shortcodes/acc.html", `<div class="acc">{{ .Inner }}</div>`)
		tem.AddTemplate("_internal/shortcodes/div.html", `<div {{with .Get "class"}} class="{{ . }}"{{ end }}>{{ .Inner }}</div>`)
		tem.AddTemplate("_internal/shortcodes/div2.html", `<div {{with .Get 0}} class="{{ . }}"{{ end }}>{{ .Inner }}</div>`)
		return nil
	}
	CheckShortCodeMatch(t,
		`{{% acc %}}{{% div %}}d1{{% /div %}}{{% div2 %}}d2{{% /div2 %}}{{% /acc %}}`,
		"<div class=\"acc\"><div >d1</div><div >d2</div>\n</div>", wt)
}

func TestIsNamedParamsSC(t *testing.T) {
	t.Parallel()
	wt := func(tem tpl.TemplateHandler) error {
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
	wt := func(tem tpl.TemplateHandler) error {
		tem.AddTemplate("_internal/shortcodes/inside.html", `<div{{with .Get "class"}} class="{{.}}"{{end}}>{{ .Inner }}</div>`)
		return nil
	}
	CheckShortCodeMatch(t, `{{< inside class="aspen" >}}`, `<div class="aspen"></div>`, wt)
	CheckShortCodeMatch(t, `{{< inside class="aspen" >}}More Here{{< /inside >}}`, "<div class=\"aspen\">More Here</div>", wt)
	CheckShortCodeMatch(t, `{{< inside >}}More Here{{< /inside >}}`, "<div>More Here</div>", wt)
}

func TestInnerSCWithMarkdown(t *testing.T) {
	t.Parallel()
	wt := func(tem tpl.TemplateHandler) error {
		tem.AddTemplate("_internal/shortcodes/inside.html", `<div{{with .Get "class"}} class="{{.}}"{{end}}>{{ .Inner }}</div>`)
		return nil
	}
	CheckShortCodeMatch(t, `{{% inside %}}
# More Here

[link](http://spf13.com) and text

{{% /inside %}}`, "<div><h1 id=\"more-here\">More Here</h1>\n\n<p><a href=\"http://spf13.com\">link</a> and text</p>\n</div>", wt)
}

func TestInnerSCWithAndWithoutMarkdown(t *testing.T) {
	t.Parallel()
	wt := func(tem tpl.TemplateHandler) error {
		tem.AddTemplate("_internal/shortcodes/inside.html", `<div{{with .Get "class"}} class="{{.}}"{{end}}>{{ .Inner }}</div>`)
		return nil
	}
	CheckShortCodeMatch(t, `{{% inside %}}
# More Here

[link](http://spf13.com) and text

{{% /inside %}}

And then:

{{< inside >}}
# More Here

This is **plain** text.

{{< /inside >}}
`, "<div><h1 id=\"more-here\">More Here</h1>\n\n<p><a href=\"http://spf13.com\">link</a> and text</p>\n</div>\n\n<p>And then:</p>\n\n<div>\n# More Here\n\nThis is **plain** text.\n\n</div>", wt)
}

func TestEmbeddedSC(t *testing.T) {
	t.Parallel()
	CheckShortCodeMatch(t, `{{% figure src="/found/here" class="bananas orange" %}}`, "<figure class=\"bananas orange\">\n    <img src=\"/found/here\"/> \n</figure>", nil)
	CheckShortCodeMatch(t, `{{% figure src="/found/here" class="bananas orange" caption="This is a caption" %}}`, "<figure class=\"bananas orange\">\n    <img src=\"/found/here\"\n         alt=\"This is a caption\"/> <figcaption>\n            <p>This is a caption</p>\n        </figcaption>\n</figure>", nil)
}

func TestNestedSC(t *testing.T) {
	t.Parallel()
	wt := func(tem tpl.TemplateHandler) error {
		tem.AddTemplate("_internal/shortcodes/scn1.html", `<div>Outer, inner is {{ .Inner }}</div>`)
		tem.AddTemplate("_internal/shortcodes/scn2.html", `<div>SC2</div>`)
		return nil
	}
	CheckShortCodeMatch(t, `{{% scn1 %}}{{% scn2 %}}{{% /scn1 %}}`, "<div>Outer, inner is <div>SC2</div>\n</div>", wt)

	CheckShortCodeMatch(t, `{{< scn1 >}}{{% scn2 %}}{{< /scn1 >}}`, "<div>Outer, inner is <div>SC2</div></div>", wt)
}

func TestNestedComplexSC(t *testing.T) {
	t.Parallel()
	wt := func(tem tpl.TemplateHandler) error {
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
	wt := func(tem tpl.TemplateHandler) error {
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
	wt := func(tem tpl.TemplateHandler) error {
		tem.AddTemplate("_internal/shortcodes/bug.html", `xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx`)
		return nil
	}
	CheckShortCodeMatch(t, `
{{< bug >}}

{{< bug >}}
`, "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx\n\nxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", wt)
}

const testScPlaceholderRegexp = "HAHAHUGOSHORTCODE-\\d+HBHB"

func TestExtractShortcodes(t *testing.T) {
	t.Parallel()

	for i, this := range []struct {
		name             string
		input            string
		expectShortCodes string
		expect           interface{}
		expectErrorMsg   string
	}{
		{"text", "Some text.", "map[]", "Some text.", ""},
		{"invalid right delim", "{{< tag }}", "", false, "unrecognized character"},
		{"invalid close", "\n{{< /tag >}}", "", false, "got closing shortcode, but none is open"},
		{"invalid close2", "\n\n{{< tag >}}{{< /anotherTag >}}", "", false, "closing tag for shortcode 'anotherTag' does not match start tag"},
		{"unterminated quote 1", `{{< figure src="im caption="S" >}}`, "", false, "got pos"},
		{"unterminated quote 1", `{{< figure src="im" caption="S >}}`, "", false, "unterm"},
		{"one shortcode, no markup", "{{< tag >}}", "", testScPlaceholderRegexp, ""},
		{"one shortcode, markup", "{{% tag %}}", "", testScPlaceholderRegexp, ""},
		{"one pos param", "{{% tag param1 %}}", `tag([\"param1\"], true){[]}"]`, testScPlaceholderRegexp, ""},
		{"two pos params", "{{< tag param1 param2>}}", `tag([\"param1\" \"param2\"], false){[]}"]`, testScPlaceholderRegexp, ""},
		{"one named param", `{{% tag param1="value" %}}`, `tag([\"param1:value\"], true){[]}`, testScPlaceholderRegexp, ""},
		{"two named params", `{{< tag param1="value1" param2="value2" >}}`, `tag([\"param1:value1\" \"param2:value2\"], false){[]}"]`,
			testScPlaceholderRegexp, ""},
		{"inner", `Some text. {{< inner >}}Inner Content{{< / inner >}}. Some more text.`, `inner([], false){[Inner Content]}`,
			fmt.Sprintf("Some text. %s. Some more text.", testScPlaceholderRegexp), ""},
		// issue #934
		{"inner self-closing", `Some text. {{< inner />}}. Some more text.`, `inner([], false){[]}`,
			fmt.Sprintf("Some text. %s. Some more text.", testScPlaceholderRegexp), ""},
		{"close, but not inner", "{{< tag >}}foo{{< /tag >}}", "", false, `shortcode "tag" has no .Inner, yet a closing tag was provided`},
		{"nested inner", `Inner->{{< inner >}}Inner Content->{{% inner2 param1 %}}inner2txt{{% /inner2 %}}Inner close->{{< / inner >}}<-done`,
			`inner([], false){[Inner Content-> inner2([\"param1\"], true){[inner2txt]} Inner close->]}`,
			fmt.Sprintf("Inner->%s<-done", testScPlaceholderRegexp), ""},
		{"nested, nested inner", `Inner->{{< inner >}}inner2->{{% inner2 param1 %}}inner2txt->inner3{{< inner3>}}inner3txt{{</ inner3 >}}{{% /inner2 %}}final close->{{< / inner >}}<-done`,
			`inner([], false){[inner2-> inner2([\"param1\"], true){[inner2txt->inner3 inner3(%!q(<nil>), false){[inner3txt]}]} final close->`,
			fmt.Sprintf("Inner->%s<-done", testScPlaceholderRegexp), ""},
		{"two inner", `Some text. {{% inner %}}First **Inner** Content{{% / inner %}} {{< inner >}}Inner **Content**{{< / inner >}}. Some more text.`,
			`map["HAHAHUGOSHORTCODE-1HBHB:inner([], true){[First **Inner** Content]}" "HAHAHUGOSHORTCODE-2HBHB:inner([], false){[Inner **Content**]}"]`,
			fmt.Sprintf("Some text. %s %s. Some more text.", testScPlaceholderRegexp, testScPlaceholderRegexp), ""},
		{"closed without content", `Some text. {{< inner param1 >}}{{< / inner >}}. Some more text.`, `inner([\"param1\"], false){[]}`,
			fmt.Sprintf("Some text. %s. Some more text.", testScPlaceholderRegexp), ""},
		{"two shortcodes", "{{< sc1 >}}{{< sc2 >}}",
			`map["HAHAHUGOSHORTCODE-1HBHB:sc1([], false){[]}" "HAHAHUGOSHORTCODE-2HBHB:sc2([], false){[]}"]`,
			testScPlaceholderRegexp + testScPlaceholderRegexp, ""},
		{"mix of shortcodes", `Hello {{< sc1 >}}world{{% sc2 p2="2"%}}. And that's it.`,
			`map["HAHAHUGOSHORTCODE-1HBHB:sc1([], false){[]}" "HAHAHUGOSHORTCODE-2HBHB:sc2([\"p2:2\"]`,
			fmt.Sprintf("Hello %sworld%s. And that's it.", testScPlaceholderRegexp, testScPlaceholderRegexp), ""},
		{"mix with inner", `Hello {{< sc1 >}}world{{% inner p2="2"%}}Inner{{%/ inner %}}. And that's it.`,
			`map["HAHAHUGOSHORTCODE-1HBHB:sc1([], false){[]}" "HAHAHUGOSHORTCODE-2HBHB:inner([\"p2:2\"], true){[Inner]}"]`,
			fmt.Sprintf("Hello %sworld%s. And that's it.", testScPlaceholderRegexp, testScPlaceholderRegexp), ""},
	} {

		pageInput := simplePage + this.input

		counter := 0
		placeholderFunc := func() string {
			counter++
			return fmt.Sprintf("HAHA%s-%dHBHB", shortcodePlaceholderPrefix, counter)
		}

		p, err := pageFromString(pageInput, "simple.md", placeholderFunc, func(templ tpl.TemplateHandler) error {
			templ.AddTemplate("_internal/shortcodes/tag.html", `tag`)
			templ.AddTemplate("_internal/shortcodes/sc1.html", `sc1`)
			templ.AddTemplate("_internal/shortcodes/sc2.html", `sc2`)
			templ.AddTemplate("_internal/shortcodes/inner.html", `{{with .Inner }}{{ . }}{{ end }}`)
			templ.AddTemplate("_internal/shortcodes/inner2.html", `{{.Inner}}`)
			templ.AddTemplate("_internal/shortcodes/inner3.html", `{{.Inner}}`)
			return nil
		})

		if b, ok := this.expect.(bool); ok && !b {
			if err == nil {
				t.Fatalf("[%d] %s: ExtractShortcodes didn't return an expected error", i, this.name)
			} else {
				r := regexp.MustCompile(this.expectErrorMsg)
				if !r.MatchString(err.Error()) {
					t.Fatalf("[%d] %s: ExtractShortcodes didn't return an expected error message, got\n%s but expected\n%s",
						i, this.name, err.Error(), this.expectErrorMsg)
				}
			}
			continue
		} else {
			if err != nil {
				t.Fatalf("[%d] %s: failed: %q", i, this.name, err)
			}
		}

		shortCodes := p.shortcodeState.shortcodes
		contentReplaced := string(p.workContent)

		var expected string
		av := reflect.ValueOf(this.expect)
		switch av.Kind() {
		case reflect.String:
			expected = av.String()
		}

		r, err := regexp.Compile(expected)

		if err != nil {
			t.Fatalf("[%d] %s: Failed to compile regexp %q: %q", i, this.name, expected, err)
		}

		if strings.Count(contentReplaced, shortcodePlaceholderPrefix) != shortCodes.Len() {
			t.Fatalf("[%d] %s: Not enough placeholders, found %d", i, this.name, shortCodes.Len())
		}

		if !r.MatchString(contentReplaced) {
			t.Fatalf("[%d] %s: Shortcode extract didn't match. got %q but expected %q", i, this.name, contentReplaced, expected)
		}

		for _, placeHolder := range shortCodes.Keys() {
			sc := shortCodes.getShortcode(placeHolder)
			if !strings.Contains(contentReplaced, placeHolder.(string)) {
				t.Fatalf("[%d] %s: Output does not contain placeholder %q", i, this.name, placeHolder)
			}

			if sc.params == nil {
				t.Fatalf("[%d] %s: Params is nil for shortcode '%s'", i, this.name, sc.name)
			}
		}

		if this.expectShortCodes != "" {
			shortCodesAsStr := fmt.Sprintf("map%q", collectAndSortShortcodes(shortCodes))
			if !strings.Contains(shortCodesAsStr, this.expectShortCodes) {
				t.Fatalf("[%d] %s: Shortcodes not as expected, got\n%s but expected\n%s", i, this.name, shortCodesAsStr, this.expectShortCodes)
			}
		}
	}
}

func TestShortcodesInSite(t *testing.T) {
	t.Parallel()
	baseURL := "http://foo/bar"

	tests := []struct {
		contentPath string
		content     string
		outFile     string
		expected    string
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
			filepath.FromSlash("public/sect/doc5/index.html"), "\n\n<h1 id=\"hahahugoshortcode-1hbhb\">b</h1>\n\n<h2 id=\"hahahugoshortcode-2hbhb\">c</h2>\n"},
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
	}

	sources := make([][2]string, len(tests))

	for i, test := range tests {
		sources[i] = [2]string{filepath.FromSlash(test.contentPath), test.content}
	}

	addTemplates := func(templ tpl.TemplateHandler) error {
		templ.AddTemplate("_default/single.html", "{{.Content}}")

		templ.AddTemplate("_internal/shortcodes/b.html", `b`)
		templ.AddTemplate("_internal/shortcodes/c.html", `c`)
		templ.AddTemplate("_internal/shortcodes/d.html", `d`)
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

	writeSourcesToSource(t, "content", fs, sources...)

	s := buildSingleSite(t, deps.DepsCfg{WithTemplate: addTemplates, Fs: fs, Cfg: cfg}, BuildCfg{})
	th := testHelper{s.Cfg, s.Fs, t}

	for _, test := range tests {
		if strings.HasSuffix(test.contentPath, ".ad") && !helpers.HasAsciidoc() {
			fmt.Println("Skip Asciidoc test case as no Asciidoc present.")
			continue
		} else if strings.HasSuffix(test.contentPath, ".rst") && !helpers.HasRst() {
			fmt.Println("Skip Rst test case as no rst2html present.")
			continue
		} else if strings.Contains(test.expected, "code") {
			fmt.Println("Skip Pygments test case as no pygments present.")
			continue
		}

		th.assertFileContent(test.outFile, test.expected)
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

	mf := afero.NewMemMapFs()

	th, h := newTestSitesFromConfig(t, mf, siteConfig,
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

	fs := th.Fs

	writeSource(t, fs, "content/_index.md", fmt.Sprintf(pageTemplate, "Home"))
	writeSource(t, fs, "content/sect/mypage.md", fmt.Sprintf(pageTemplate, "Single"))
	writeSource(t, fs, "content/sect/mycsvpage.md", fmt.Sprintf(pageTemplateCSVOnly, "Single CSV"))

	err := h.Build(BuildCfg{})
	require.NoError(t, err)
	require.Len(t, h.Sites, 1)

	s := h.Sites[0]
	home := s.getPage(KindHome)
	require.NotNil(t, home)
	require.Len(t, home.outputFormats, 3)

	th.assertFileContent("public/index.html",
		"Home HTML",
		"ShortHTML",
		"ShortNoExt",
		"ShortOnlyHTML",
		"myInner:--ShortHTML--",
	)

	th.assertFileContent("public/amp/index.html",
		"Home AMP",
		"ShortAMP",
		"ShortNoExt",
		"ShortOnlyHTML",
		"myInner:--ShortAMP--",
	)

	th.assertFileContent("public/index.ics",
		"Home Calendar",
		"ShortCalendar",
		"ShortNoExt",
		"ShortOnlyHTML",
		"myInner:--ShortCalendar--",
	)

	th.assertFileContent("public/sect/mypage/index.html",
		"Single HTML",
		"ShortHTML",
		"ShortNoExt",
		"ShortOnlyHTML",
		"myInner:--ShortHTML--",
	)

	th.assertFileContent("public/sect/mypage/index.json",
		"Single JSON",
		"ShortJSON",
		"ShortNoExt",
		"ShortOnlyHTML",
		"myInner:--ShortJSON--",
	)

	th.assertFileContent("public/amp/sect/mypage/index.html",
		// No special AMP template
		"Single HTML",
		"ShortAMP",
		"ShortNoExt",
		"ShortOnlyHTML",
		"myInner:--ShortAMP--",
	)

	th.assertFileContent("public/sect/mycsvpage/index.csv",
		"Single CSV",
		"ShortCSV",
	)

}

func collectAndSortShortcodes(shortcodes *orderedMap) []string {
	var asArray []string

	for _, key := range shortcodes.Keys() {
		sc := shortcodes.getShortcode(key)
		asArray = append(asArray, fmt.Sprintf("%s:%s", key, sc))
	}

	sort.Strings(asArray)
	return asArray

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
			results, err := replaceShortcodeTokens(currIn.in, "HUGOSHORTCODE", currIn.replacements)

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
		{"Hello HAHAPREFIX-1HBHB.", "PREFIX", map[string]string{"HAHAPREFIX-1HBHB": "World"}, "Hello World."},
		{"Hello HAHAPREFIX-1@}@.", "PREFIX", map[string]string{"HAHAPREFIX-1HBHB": "World"}, false},
		{"HAHAPREFIX2-1HBHB", "PREFIX2", map[string]string{"HAHAPREFIX2-1HBHB": "World"}, "World"},
		{"Hello World!", "PREFIX2", map[string]string{}, "Hello World!"},
		{"!HAHAPREFIX-1HBHB", "PREFIX", map[string]string{"HAHAPREFIX-1HBHB": "World"}, "!World"},
		{"HAHAPREFIX-1HBHB!", "PREFIX", map[string]string{"HAHAPREFIX-1HBHB": "World"}, "World!"},
		{"!HAHAPREFIX-1HBHB!", "PREFIX", map[string]string{"HAHAPREFIX-1HBHB": "World"}, "!World!"},
		{"_{_PREFIX-1HBHB", "PREFIX", map[string]string{"HAHAPREFIX-1HBHB": "World"}, "_{_PREFIX-1HBHB"},
		{"Hello HAHAPREFIX-1HBHB.", "PREFIX", map[string]string{"HAHAPREFIX-1HBHB": "To You My Old Friend Who Told Me This Fantastic Story"}, "Hello To You My Old Friend Who Told Me This Fantastic Story."},
		{"A HAHAA-1HBHB asdf HAHAA-2HBHB.", "A", map[string]string{"HAHAA-1HBHB": "v1", "HAHAA-2HBHB": "v2"}, "A v1 asdf v2."},
		{"Hello HAHAPREFIX2-1HBHB. Go HAHAPREFIX2-2HBHB, Go, Go HAHAPREFIX2-3HBHB Go Go!.", "PREFIX2", map[string]string{"HAHAPREFIX2-1HBHB": "Europe", "HAHAPREFIX2-2HBHB": "Jonny", "HAHAPREFIX2-3HBHB": "Johnny"}, "Hello Europe. Go Jonny, Go, Go Johnny Go Go!."},
		{"A HAHAPREFIX-2HBHB HAHAPREFIX-1HBHB.", "PREFIX", map[string]string{"HAHAPREFIX-1HBHB": "A", "HAHAPREFIX-2HBHB": "B"}, "A B A."},
		{"A HAHAPREFIX-1HBHB HAHAPREFIX-2", "PREFIX", map[string]string{"HAHAPREFIX-1HBHB": "A"}, false},
		{"A HAHAPREFIX-1HBHB but not the second.", "PREFIX", map[string]string{"HAHAPREFIX-1HBHB": "A", "HAHAPREFIX-2HBHB": "B"}, "A A but not the second."},
		{"An HAHAPREFIX-1HBHB.", "PREFIX", map[string]string{"HAHAPREFIX-1HBHB": "A", "HAHAPREFIX-2HBHB": "B"}, "An A."},
		{"An HAHAPREFIX-1HBHB HAHAPREFIX-2HBHB.", "PREFIX", map[string]string{"HAHAPREFIX-1HBHB": "A", "HAHAPREFIX-2HBHB": "B"}, "An A B."},
		{"A HAHAPREFIX-1HBHB HAHAPREFIX-2HBHB HAHAPREFIX-3HBHB HAHAPREFIX-1HBHB HAHAPREFIX-3HBHB.", "PREFIX", map[string]string{"HAHAPREFIX-1HBHB": "A", "HAHAPREFIX-2HBHB": "B", "HAHAPREFIX-3HBHB": "C"}, "A A B C A C."},
		{"A HAHAPREFIX-1HBHB HAHAPREFIX-2HBHB HAHAPREFIX-3HBHB HAHAPREFIX-1HBHB HAHAPREFIX-3HBHB.", "PREFIX", map[string]string{"HAHAPREFIX-1HBHB": "A", "HAHAPREFIX-2HBHB": "B", "HAHAPREFIX-3HBHB": "C"}, "A A B C A C."},
		// Issue #1148 remove p-tags 10 =>
		{"Hello <p>HAHAPREFIX-1HBHB</p>. END.", "PREFIX", map[string]string{"HAHAPREFIX-1HBHB": "World"}, "Hello World. END."},
		{"Hello <p>HAHAPREFIX-1HBHB</p>. <p>HAHAPREFIX-2HBHB</p> END.", "PREFIX", map[string]string{"HAHAPREFIX-1HBHB": "World", "HAHAPREFIX-2HBHB": "THE"}, "Hello World. THE END."},
		{"Hello <p>HAHAPREFIX-1HBHB. END</p>.", "PREFIX", map[string]string{"HAHAPREFIX-1HBHB": "World"}, "Hello <p>World. END</p>."},
		{"<p>Hello HAHAPREFIX-1HBHB</p>. END.", "PREFIX", map[string]string{"HAHAPREFIX-1HBHB": "World"}, "<p>Hello World</p>. END."},
		{"Hello <p>HAHAPREFIX-1HBHB12", "PREFIX", map[string]string{"HAHAPREFIX-1HBHB": "World"}, "Hello <p>World12"},
		{"Hello HAHAP-1HBHB. HAHAP-1HBHB-HAHAP-1HBHB HAHAP-1HBHB HAHAP-1HBHB HAHAP-1HBHB END", "P", map[string]string{"HAHAP-1HBHB": strings.Repeat("BC", 100)},
			fmt.Sprintf("Hello %s. %s-%s %s %s %s END",
				strings.Repeat("BC", 100), strings.Repeat("BC", 100), strings.Repeat("BC", 100), strings.Repeat("BC", 100), strings.Repeat("BC", 100), strings.Repeat("BC", 100))},
	} {

		results, err := replaceShortcodeTokens([]byte(this.input), this.prefix, this.replacements)

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

func TestScKey(t *testing.T) {
	require.Equal(t, scKey{Suffix: "xml", ShortcodePlaceholder: "ABCD"},
		newScKey(media.XMLType, "ABCD"))
	require.Equal(t, scKey{Lang: "en", Suffix: "html", OutputFormat: "AMP", ShortcodePlaceholder: "EFGH"},
		newScKeyFromLangAndOutputFormat("en", output.AMPFormat, "EFGH"))
	require.Equal(t, scKey{Suffix: "html", ShortcodePlaceholder: "IJKL"},
		newDefaultScKey("IJKL"))

}

func TestShortcodeGetContent(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

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

	v := viper.New()

	v.Set("timeout", 500)

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

	builder.WithViper(v).WithContent(content...).WithTemplates(templates...).CreateSites().Build(BuildCfg{})
	s := builder.H.Sites[0]
	assert.Equal(3, len(s.RegularPages))

	builder.AssertFileContent("public/section1/index.html",
		"List Content: <p>Logo:P1:|P2:logo.png/PNG logo|:P1: P1:|P2:docs1p1/<p>C-s1p1</p>\n|",
		"BP1:P1:|P2:docbp1/<p>C-bp1</p>",
	)

	builder.AssertFileContent("public/b1/index.html",
		"Single Content: <p>Logo:P1:|P2:logo.png/PNG logo|:P1: P1:|P2:docs1p1/<p>C-s1p1</p>\n|",
		"P2:docbp1/<p>C-bp1</p>",
	)

	builder.AssertFileContent("public/section2/s2p1/index.html",
		"Single Content: <p>Logo:P1:|P2:logo.png/PNG logo|:P1: P1:|P2:docs1p1/<p>C-s1p1</p>\n|",
		"P2:docbp1/<p>C-bp1</p>",
	)

}

func TestShortcodePreserveOrder(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

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
	assert.Equal(3, len(s.RegularPages))

	builder.AssertFileContent("public/en/p1/index.html", `v1: 0 sgo: |v2: 1 sgo: 0|v3: 2 sgo: 1|v4: 3 sgo: 2|v5: 4 sgo: 3`)
	builder.AssertFileContent("public/en/p1/index.html", `outer ordinal: 5 inner: 
ordinal: 0 scratch ordinal: 1 scratch get ordinal: 0
ordinal: 2 scratch ordinal: 3 scratch get ordinal: 2
ordinal: 4 scratch ordinal: 5 scratch get ordinal: 4`)

}

func TestShortcodeVariables(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

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
	assert.Equal(1, len(s.RegularPages))

	builder.AssertFileContent("public/page/index.html",
		filepath.FromSlash("File: content/page.md"),
		"Line: 7", "Column: 4", "Offset: 40",
		filepath.FromSlash("String: \"content/page.md:7:4\""),
		"Name: s1",
	)

}

func TestInlineShortcodes(t *testing.T) {
	for _, enableInlineShortcodes := range []bool{true, false} {
		t.Run(fmt.Sprintf("enableInlineShortcodes=%t", enableInlineShortcodes),
			func(t *testing.T) {
				conf := fmt.Sprintf(`
baseURL = "https://example.com"
enableInlineShortcodes = %t
`, enableInlineShortcodes)

				b := newTestSitesBuilder(t)
				b.WithConfigFile("toml", conf)
				b.WithContent("page-md-shortcode.md", `---
title: "Hugo"
---

FIRST:{{< myshort.inline "first" >}}
Page: {{ .Page.Title }}
Seq: {{ seq 3 }}
Param: {{ .Get 0 }}
{{< /myshort.inline >}}:END:

SECOND:{{< myshort.inline "second" />}}:END

`)

				b.WithTemplatesAdded("layouts/_default/single.html", `
CONTENT:{{ .Content }}
`)

				b.CreateSites().Build(BuildCfg{})

				if enableInlineShortcodes {
					b.AssertFileContent("public/page-md-shortcode/index.html",
						"Page: Hugo",
						"Seq: [1 2 3]",
						"Param: first",
						"Param: second",
					)
				} else {
					b.AssertFileContent("public/page-md-shortcode/index.html",
						"FIRST::END",
						"SECOND::END",
					)
				}

			})

	}
}
