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
package tplimpl

import (
	"bytes"
	"fmt"
	"testing"

	"html/template"

	"github.com/spf13/cast"

	"github.com/stretchr/testify/require"
)

var (
	testFuncs = map[string]interface{}{
		"ToTime": func(v interface{}) interface{} { return cast.ToTime(v) },
		"First":  func(v ...interface{}) interface{} { return v[0] },
		"Echo":   func(v interface{}) interface{} { return v },
		"where": func(seq, key interface{}, args ...interface{}) (interface{}, error) {
			return map[string]interface{}{
				"ByWeight": fmt.Sprintf("%v:%v:%v", seq, key, args),
			}, nil
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
`
)

func TestParamsKeysToLower(t *testing.T) {
	t.Parallel()

	require.Error(t, applyTemplateTransformers(nil, nil))

	templ, err := template.New("foo").Funcs(testFuncs).Parse(paramsTempl)

	require.NoError(t, err)

	c := newTemplateContext(createParseTreeLookup(templ))

	require.Equal(t, -1, c.decl.indexOfReplacementStart([]string{}))

	c.paramsKeysToLower(templ.Tree.Root)

	var b bytes.Buffer

	require.NoError(t, templ.Execute(&b, paramsData))

	result := b.String()

	require.Contains(t, result, "P1: P1L")
	require.Contains(t, result, "P1_2: P1L")
	require.Contains(t, result, "P1_3: P1L")
	require.Contains(t, result, "P1_4: P1L")
	require.Contains(t, result, "P2: P2L")
	require.Contains(t, result, "P2_2: P2L")
	require.Contains(t, result, "P2_3: P2L")
	require.Contains(t, result, "P2_4: P2L")
	require.Contains(t, result, "P22: P22L")
	require.Contains(t, result, "P22_nested: P22L_nested")
	require.Contains(t, result, "P3: P3H")
	require.Contains(t, result, "P3_2: P3H")
	require.Contains(t, result, "P3_3: P3H")
	require.Contains(t, result, "P3_4: P3H")
	require.Contains(t, result, "P4: 13")
	require.Contains(t, result, "P5: P1L")
	require.Contains(t, result, "P5_2: P2L")

	require.Contains(t, result, "IF: P1L")
	require.Contains(t, result, "ELSE: P1L")

	require.Contains(t, result, "WITH: P1L")

	require.Contains(t, result, "RANGE: 3: P1L")

	require.Contains(t, result, "Hi There")

	// Issue #2740
	require.Contains(t, result, "F1: themes/P2L-theme")
	require.Contains(t, result, "F2: themes/P2L-theme")
	require.Contains(t, result, "F3: themes/P2L-theme")

	require.Contains(t, result, "PSLICE: PSLICE1|PSLICE3|")
	require.Contains(t, result, "PARAMS STRING: foo:.Params.toc_hide:[!= true]")
	require.Contains(t, result, "PARAMS STRING2: foo:.Params.toc_hide:[!= true]")
	require.Contains(t, result, "PARAMS STRING3: .Params.TOC_HIDE:!=:[P1L]")

	// Issue #5094
	require.Contains(t, result, "PARAMS COMPOSITE: [1 3]")

	// Issue #5068
	require.Contains(t, result, "PCurrentSection: pcurrentsection")

	// Issue #5541
	require.Contains(t, result, "PARAMS TIME: 1972-02-28")
	require.Contains(t, result, "PARAMS TIME2: 1972-02-28")

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
		c.paramsKeysToLower(templ.Tree.Root)
	}
}

func TestParamsKeysToLowerVars(t *testing.T) {
	t.Parallel()
	var (
		ctx = map[string]interface{}{
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

	require.NoError(t, err)

	c := newTemplateContext(createParseTreeLookup(templ))

	c.paramsKeysToLower(templ.Tree.Root)

	var b bytes.Buffer

	require.NoError(t, templ.Execute(&b, ctx))

	result := b.String()

	require.Contains(t, result, "Color: Amber")
	require.Contains(t, result, "Blue: Amber")
	require.Contains(t, result, "Pretty First1: Indigo")
	require.Contains(t, result, "Pretty First2: Indigo")
	require.Contains(t, result, "Pretty First3: Indigo")

}

func TestParamsKeysToLowerInBlockTemplate(t *testing.T) {
	t.Parallel()

	var (
		ctx = map[string]interface{}{
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
	require.NoError(t, err)

	overlayTpl, err := template.Must(masterTpl.Clone()).Parse(overlay)
	require.NoError(t, err)
	overlayTpl = overlayTpl.Lookup(overlayTpl.Name())

	c := newTemplateContext(createParseTreeLookup(overlayTpl))

	c.paramsKeysToLower(overlayTpl.Tree.Root)

	var b bytes.Buffer

	require.NoError(t, overlayTpl.Execute(&b, ctx))

	result := b.String()

	require.Contains(t, result, "P1: P1L")
	require.Contains(t, result, "P2: P1L")
}

// Issue #2927
func TestTransformRecursiveTemplate(t *testing.T) {

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
	require.NoError(t, err)

	c := newTemplateContext(createParseTreeLookup(templ))
	c.paramsKeysToLower(templ.Tree.Root)

}
