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
	"strings"

	template "github.com/gohugoio/hugo/tpl/internal/go_templates/htmltemplate"

	"testing"
	"time"

	"github.com/gohugoio/hugo/tpl"

	qt "github.com/frankban/quicktest"
)

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
	h := d.Tmpl.(tpl.TemplateManager)

	// HTML templates
	c.Assert(h.AddTemplate("mytemplate.html", templ1), qt.IsNil)
	c.Assert(h.AddTemplate("othertemplate.html", templ2), qt.IsNil)

	// Text templates
	c.Assert(h.AddTemplate("_text/mytexttemplate.txt", templ1), qt.IsNil)
	c.Assert(h.AddTemplate("_text/myothertexttemplate.txt", templ2), qt.IsNil)

	c.Assert(h.MarkReady(), qt.IsNil)

	for _, name := range []string{"mytemplate.html", "mytexttemplate.txt"} {
		tt, _ := d.Tmpl.Lookup(name)
		sb := &strings.Builder{}

		err := d.Tmpl.Execute(tt, sb, ctx)
		c.Assert(err, qt.IsNil)

		result := sb.String()

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
