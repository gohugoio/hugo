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
package tplimpl

import (
	"bytes"
	"fmt"
	"github.com/gohugoio/hugo/tpl/internal/go_templates/htmltemplate"

	"testing"
	"time"

	"github.com/gohugoio/hugo/tpl"

	"github.com/spf13/cast"

	qt "github.com/frankban/quicktest"
)

type paramsHolder struct {
	params map[string]interface{}
	page   *paramsHolder
}

func (p paramsHolder) Params() map[string]interface{} {
	return p.params
}

func (p paramsHolder) GetPage(arg string) *paramsHolder {
	return p.page
}

var (
	testFuncs = map[string]interface{}{
		"getif":  func(v interface{}) interface{} { return v },
		"ToTime": func(v interface{}) interface{} { return cast.ToTime(v) },
		"First":  func(v ...interface{}) interface{} { return v[0] },
		"Echo":   func(v interface{}) interface{} { return v },
		"where": func(seq, key interface{}, args ...interface{}) (interface{}, error) {
			return map[string]interface{}{
				"ByWeight": fmt.Sprintf("%v:%v:%v", seq, key, args),
			}, nil
		},
		"site": func() paramsHolder {
			return paramsHolder{
				params: map[string]interface{}{
					"lower": "global-site",
				},
				page: &paramsHolder{
					params: map[string]interface{}{
						"lower": "page",
					},
				},
			}
		},
	}

	paramsData = map[string]interface{}{

		"NotParam": "Hi There",
		"Slice":    []int{1, 3},
		"Params": map[string]interface{}{
			"lower":  "P1L",
			"slice":  []int{1, 3},
			"mydate": "1972-01-28",
		},
		"Pages": map[string]interface{}{
			"ByWeight": []int{1, 3},
		},
		"CurrentSection": map[string]interface{}{
			"Params": map[string]interface{}{
				"lower": "pcurrentsection",
			},
		},
		"Site": map[string]interface{}{
			"Params": map[string]interface{}{
				"lower": "P2L",
				"slice": []int{1, 3},
			},
			"Language": map[string]interface{}{
				"Params": map[string]interface{}{
					"lower": "P22L",
					"nested": map[string]interface{}{
						"lower": "P22L_nested",
					},
				},
			},
			"Data": map[string]interface{}{
				"Params": map[string]interface{}{
					"NOLOW": "P3H",
				},
			},
		},
		"Site2": paramsHolder{
			params: map[string]interface{}{
				"lower": "global-site",
			},
			page: &paramsHolder{
				params: map[string]interface{}{
					"lower": "page",
				},
			},
		},
	}

	paramsTempl = `
{{ $page := . }}
{{ $pages := .Pages }}
{{ $pageParams := .Params }}
{{ $site := .Site }}
{{ $siteParams := .Site.Params }}
{{ $data := .Site.Data }}
{{ $notparam := .NotParam }}

PCurrentSection: {{ .CurrentSection.Params.LOWER }}
P1: {{ .Params.LOWER }}
P1_2: {{ $.Params.LOWER }}
P1_3: {{ $page.Params.LOWER }}
P1_4: {{ $pageParams.LOWER }}
P2: {{ .Site.Params.LOWER }}
P2_2: {{ $.Site.Params.LOWER }}
P2_3: {{ $site.Params.LOWER }}
P2_4: {{ $siteParams.LOWER }}
P22: {{ .Site.Language.Params.LOWER }}
P22_nested: {{ .Site.Language.Params.NESTED.LOWER }}
P3: {{ .Site.Data.Params.NOLOW }}
P3_2: {{ $.Site.Data.Params.NOLOW }}
P3_3: {{ $site.Data.Params.NOLOW }}
P3_4: {{ $data.Params.NOLOW }}
P4: {{ range $i, $e := .Site.Params.SLICE }}{{ $e }}{{ end }}
P5: {{ Echo .Params.LOWER }}
P5_2: {{ Echo $site.Params.LOWER }}
{{ if .Params.LOWER }}
IF: {{ .Params.LOWER }}
{{ end }}
{{ if .Params.NOT_EXIST }}
{{ else }}
ELSE: {{ .Params.LOWER }}
{{ end }}


{{ with .Params.LOWER }}
WITH: {{ . }}
{{ end }}


{{ range .Slice }}
RANGE: {{ . }}: {{ $.Params.LOWER }}
{{ end }}
{{ index .Slice 1 }}
{{ .NotParam }}
{{ .NotParam }}
{{ .NotParam }}
{{ .NotParam }}
{{ .NotParam }}
{{ .NotParam }}
{{ .NotParam }}
{{ .NotParam }}
{{ .NotParam }}
{{ .NotParam }}
{{ $notparam }}


{{ $lower := .Site.Params.LOWER }}
F1: {{ printf "themes/%s-theme" .Site.Params.LOWER }}
F2: {{ Echo (printf "themes/%s-theme" $lower) }}
F3: {{ Echo (printf "themes/%s-theme" .Site.Params.LOWER) }}

PSLICE: {{ range .Params.SLICE }}PSLICE{{.}}|{{ end }}

{{ $pages := "foo" }}
{{ $pages := where $pages ".Params.toc_hide" "!=" true }}
PARAMS STRING: {{ $pages.ByWeight }}
PARAMS STRING2: {{ with $pages }}{{ .ByWeight }}{{ end }}
{{ $pages3 := where ".Params.TOC_HIDE" "!=" .Params.LOWER }}
PARAMS STRING3: {{ $pages3.ByWeight }}
{{ $first := First .Pages .Site.Params.LOWER }}
PARAMS COMPOSITE: {{ $first.ByWeight }}


{{ $time := $.Params.MyDate | ToTime }}
{{ $time = $time.AddDate 0 1 0 }}
PARAMS TIME: {{ $time.Format "2006-01-02" }}

{{ $_x :=  $.Params.MyDate | ToTime }}
PARAMS TIME2: {{ $_x.AddDate 0 1 0 }}

PARAMS SITE GLOBAL1: {{ site.Params.LOwER }}
{{ $lower := site.Params.LOwER }}
{{ $site := site }}
PARAMS SITE GLOBAL2: {{ $lower }}
PARAMS SITE GLOBAL3: {{ $site.Params.LOWER }}

{{ $p := $site.GetPage "foo" }}
PARAMS GETPAGE: {{ $p.Params.LOWER }}
{{ $p := .Site2.GetPage "foo" }}
PARAMS GETPAGE2: {{ $p.Params.LOWER }}
`
)

func TestParamsKeysToLower(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	_, err := applyTemplateTransformers(templateUndefined, nil, nil)
	c.Assert(err, qt.Not(qt.IsNil))

	templ, err := template.New("foo").Funcs(testFuncs).Parse(paramsTempl)

	c.Assert(err, qt.IsNil)

	ctx := newTemplateContext(createParseTreeLookup(templ))

	c.Assert(ctx.decl.indexOfReplacementStart([]string{}), qt.Equals, -1)

	ctx.applyTransformations(templ.Tree.Root)

	var b bytes.Buffer

	c.Assert(templ.Execute(&b, paramsData), qt.IsNil)

	result := b.String()

	c.Assert(result, qt.Contains, "P1: P1L")
	c.Assert(result, qt.Contains, "P1_2: P1L")
	c.Assert(result, qt.Contains, "P1_3: P1L")
	c.Assert(result, qt.Contains, "P1_4: P1L")
	c.Assert(result, qt.Contains, "P2: P2L")
	c.Assert(result, qt.Contains, "P2_2: P2L")
	c.Assert(result, qt.Contains, "P2_3: P2L")
	c.Assert(result, qt.Contains, "P2_4: P2L")
	c.Assert(result, qt.Contains, "P22: P22L")
	c.Assert(result, qt.Contains, "P22_nested: P22L_nested")
	c.Assert(result, qt.Contains, "P3: P3H")
	c.Assert(result, qt.Contains, "P3_2: P3H")
	c.Assert(result, qt.Contains, "P3_3: P3H")
	c.Assert(result, qt.Contains, "P3_4: P3H")
	c.Assert(result, qt.Contains, "P4: 13")
	c.Assert(result, qt.Contains, "P5: P1L")
	c.Assert(result, qt.Contains, "P5_2: P2L")

	c.Assert(result, qt.Contains, "IF: P1L")
	c.Assert(result, qt.Contains, "ELSE: P1L")

	c.Assert(result, qt.Contains, "WITH: P1L")

	c.Assert(result, qt.Contains, "RANGE: 3: P1L")

	c.Assert(result, qt.Contains, "Hi There")

	// Issue #2740
	c.Assert(result, qt.Contains, "F1: themes/P2L-theme")
	c.Assert(result, qt.Contains, "F2: themes/P2L-theme")
	c.Assert(result, qt.Contains, "F3: themes/P2L-theme")

	c.Assert(result, qt.Contains, "PSLICE: PSLICE1|PSLICE3|")
	c.Assert(result, qt.Contains, "PARAMS STRING: foo:.Params.toc_hide:[!= true]")
	c.Assert(result, qt.Contains, "PARAMS STRING2: foo:.Params.toc_hide:[!= true]")
	c.Assert(result, qt.Contains, "PARAMS STRING3: .Params.TOC_HIDE:!=:[P1L]")

	// Issue #5094
	c.Assert(result, qt.Contains, "PARAMS COMPOSITE: [1 3]")

	// Issue #5068
	c.Assert(result, qt.Contains, "PCurrentSection: pcurrentsection")

	// Issue #5541
	c.Assert(result, qt.Contains, "PARAMS TIME: 1972-02-28")
	c.Assert(result, qt.Contains, "PARAMS TIME2: 1972-02-28")

	// Issue ##5615
	c.Assert(result, qt.Contains, "PARAMS SITE GLOBAL1: global-site")
	c.Assert(result, qt.Contains, "PARAMS SITE GLOBAL2: global-site")
	c.Assert(result, qt.Contains, "PARAMS SITE GLOBAL3: global-site")

	//
	c.Assert(result, qt.Contains, "PARAMS GETPAGE: page")
	c.Assert(result, qt.Contains, "PARAMS GETPAGE2: page")

}

func BenchmarkTemplateParamsKeysToLower(b *testing.B) {
	templ, err := template.New("foo").Funcs(testFuncs).Parse(paramsTempl)

	if err != nil {
		b.Fatal(err)
	}

	templates := make([]*template.Template, b.N)

	for i := 0; i < b.N; i++ {
		templates[i], err = templ.Clone()
		if err != nil {
			b.Fatal(err)
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c := newTemplateContext(createParseTreeLookup(templates[i]))
		c.applyTransformations(templ.Tree.Root)
	}
}

func TestParamsKeysToLowerVars(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	var (
		data = map[string]interface{}{
			"Params": map[string]interface{}{
				"colors": map[string]interface{}{
					"blue": "Amber",
					"pretty": map[string]interface{}{
						"first": "Indigo",
					},
				},
			},
		}

		// This is how Amber behaves:
		paramsTempl = `
{{$__amber_1 := .Params.Colors}}
{{$__amber_2 := $__amber_1.Blue}}
{{$__amber_3 := $__amber_1.Pretty}}
{{$__amber_4 := .Params}}

Color: {{$__amber_2}}
Blue: {{ $__amber_1.Blue}}
Pretty First1: {{ $__amber_3.First}}
Pretty First2: {{ $__amber_1.Pretty.First}}
Pretty First3: {{ $__amber_4.COLORS.PRETTY.FIRST}}
`
	)

	templ, err := template.New("foo").Parse(paramsTempl)

	c.Assert(err, qt.IsNil)

	ctx := newTemplateContext(createParseTreeLookup(templ))

	ctx.applyTransformations(templ.Tree.Root)

	var b bytes.Buffer

	c.Assert(templ.Execute(&b, data), qt.IsNil)

	result := b.String()

	c.Assert(result, qt.Contains, "Color: Amber")
	c.Assert(result, qt.Contains, "Blue: Amber")
	c.Assert(result, qt.Contains, "Pretty First1: Indigo")
	c.Assert(result, qt.Contains, "Pretty First2: Indigo")
	c.Assert(result, qt.Contains, "Pretty First3: Indigo")

}

func TestParamsKeysToLowerInBlockTemplate(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	var (
		data = map[string]interface{}{
			"Params": map[string]interface{}{
				"lower": "P1L",
			},
		}

		master = `
P1: {{ .Params.LOWER }}
{{ block "main" . }}DEFAULT{{ end }}`
		overlay = `
{{ define "main" }}
P2: {{ .Params.LOWER }}
{{ end }}`
	)

	masterTpl, err := template.New("foo").Parse(master)
	c.Assert(err, qt.IsNil)

	overlayTpl, err := template.Must(masterTpl.Clone()).Parse(overlay)
	c.Assert(err, qt.IsNil)
	overlayTpl = overlayTpl.Lookup(overlayTpl.Name())

	ctx := newTemplateContext(createParseTreeLookup(overlayTpl))

	ctx.applyTransformations(overlayTpl.Tree.Root)

	var b bytes.Buffer

	c.Assert(overlayTpl.Execute(&b, data), qt.IsNil)

	result := b.String()

	c.Assert(result, qt.Contains, "P1: P1L")
	c.Assert(result, qt.Contains, "P2: P1L")
}

// Issue #2927
func TestTransformRecursiveTemplate(t *testing.T) {
	c := qt.New(t)

	recursive := `
{{ define "menu-nodes" }}
{{ template "menu-node" }}
{{ end }}
{{ define "menu-node" }}
{{ template "menu-node" }}
{{ end }}
{{ template "menu-nodes" }}
`

	templ, err := template.New("foo").Parse(recursive)
	c.Assert(err, qt.IsNil)

	ctx := newTemplateContext(createParseTreeLookup(templ))
	ctx.applyTransformations(templ.Tree.Root)

}

type I interface {
	Method0()
}

type T struct {
	NonEmptyInterfaceTypedNil I
}

func (T) Method0() {
}

func TestInsertIsZeroFunc(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	var (
		ctx = map[string]interface{}{
			"True":     true,
			"Now":      time.Now(),
			"TimeZero": time.Time{},
			"T":        &T{NonEmptyInterfaceTypedNil: (*T)(nil)},
		}

		templ1 = `
{{ if .True }}.True: TRUE{{ else }}.True: FALSE{{ end }}
{{ if .TimeZero }}.TimeZero1: TRUE{{ else }}.TimeZero1: FALSE{{ end }}
{{ if (.TimeZero) }}.TimeZero2: TRUE{{ else }}.TimeZero2: FALSE{{ end }}
{{ if not .TimeZero }}.TimeZero3: TRUE{{ else }}.TimeZero3: FALSE{{ end }}
{{ if .Now }}.Now: TRUE{{ else }}.Now: FALSE{{ end }}
{{ with .TimeZero }}.TimeZero1 with: {{ . }}{{ else }}.TimeZero1 with: FALSE{{ end }}
{{ template "mytemplate" . }}
{{ if .T.NonEmptyInterfaceTypedNil }}.NonEmptyInterfaceTypedNil: TRUE{{ else }}.NonEmptyInterfaceTypedNil: FALSE{{ end }}

{{ template "other-file-template" . }}

{{ define "mytemplate" }}
{{ if .TimeZero }}.TimeZero1: mytemplate: TRUE{{ else }}.TimeZero1: mytemplate: FALSE{{ end }}
{{ end }}

`

		// https://github.com/gohugoio/hugo/issues/5865
		templ2 = `{{ define "other-file-template" }}
{{ if .TimeZero }}.TimeZero1: other-file-template: TRUE{{ else }}.TimeZero1: other-file-template: FALSE{{ end }}
{{ end }}		
`
	)

	d := newD(c)
	h := d.Tmpl.(tpl.TemplateHandler)

	// HTML templates
	c.Assert(h.AddTemplate("mytemplate.html", templ1), qt.IsNil)
	c.Assert(h.AddTemplate("othertemplate.html", templ2), qt.IsNil)

	// Text templates
	c.Assert(h.AddTemplate("_text/mytexttemplate.txt", templ1), qt.IsNil)
	c.Assert(h.AddTemplate("_text/myothertexttemplate.txt", templ2), qt.IsNil)

	c.Assert(h.MarkReady(), qt.IsNil)

	for _, name := range []string{"mytemplate.html", "mytexttemplate.txt"} {
		tt, _ := d.Tmpl.Lookup(name)
		result, err := tt.(tpl.TemplateExecutor).ExecuteToString(ctx)
		c.Assert(err, qt.IsNil)

		c.Assert(result, qt.Contains, ".True: TRUE")
		c.Assert(result, qt.Contains, ".TimeZero1: FALSE")
		c.Assert(result, qt.Contains, ".TimeZero2: FALSE")
		c.Assert(result, qt.Contains, ".TimeZero3: TRUE")
		c.Assert(result, qt.Contains, ".Now: TRUE")
		c.Assert(result, qt.Contains, "TimeZero1 with: FALSE")
		c.Assert(result, qt.Contains, ".TimeZero1: mytemplate: FALSE")
		c.Assert(result, qt.Contains, ".TimeZero1: other-file-template: FALSE")
		c.Assert(result, qt.Contains, ".NonEmptyInterfaceTypedNil: FALSE")
	}

}

func TestCollectInfo(t *testing.T) {

	configStr := `{ "version": 42 }`

	tests := []struct {
		name      string
		tplString string
		expected  tpl.Info
	}{
		{"Basic Inner", `{{ .Inner }}`, tpl.Info{IsInner: true, Config: tpl.DefaultConfig}},
		{"Basic config map", "{{ $_hugo_config := `" + configStr + "`  }}", tpl.Info{
			Config: tpl.Config{
				Version: 42,
			},
		}},
	}

	echo := func(in interface{}) interface{} {
		return in
	}

	funcs := template.FuncMap{
		"highlight": echo,
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := qt.New(t)

			templ, err := template.New("foo").Funcs(funcs).Parse(test.tplString)
			c.Assert(err, qt.IsNil)

			ctx := newTemplateContext(createParseTreeLookup(templ))
			ctx.typ = templateShortcode
			ctx.applyTransformations(templ.Tree.Root)

			c.Assert(ctx.Info, qt.Equals, test.expected)
		})
	}

}

func TestPartialReturn(t *testing.T) {

	tests := []struct {
		name      string
		tplString string
		expected  bool
	}{
		{"Basic", `
{{ $a := "Hugo Rocks!" }}
{{ return $a }}
`, true},
		{"Expression", `
{{ return add 32 }}
`, true},
	}

	echo := func(in interface{}) interface{} {
		return in
	}

	funcs := template.FuncMap{
		"return": echo,
		"add":    echo,
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := qt.New(t)

			templ, err := template.New("foo").Funcs(funcs).Parse(test.tplString)
			c.Assert(err, qt.IsNil)

			_, err = applyTemplateTransformers(templatePartial, templ.Tree, createParseTreeLookup(templ))

			// Just check that it doesn't fail in this test. We have functional tests
			// in hugoblib.
			c.Assert(err, qt.IsNil)

		})
	}

}
