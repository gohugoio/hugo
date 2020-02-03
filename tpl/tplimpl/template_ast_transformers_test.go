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
	"testing"

	template "github.com/gohugoio/hugo/tpl/internal/go_templates/htmltemplate"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/tpl"
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
	ts := newTestTemplate(templ)

	ctx := newTemplateContext(
		ts,
		newTestTemplateLookup(ts),
	)
	ctx.applyTransformations(templ.Tree.Root)

}

func newTestTemplate(templ tpl.Template) *templateState {
	return newTemplateState(
		templ,
		templateInfo{
			name: templ.Name(),
		},
	)
}

func newTestTemplateLookup(in *templateState) func(name string) *templateState {
	m := make(map[string]*templateState)
	return func(name string) *templateState {
		if in.Name() == name {
			return in
		}

		if ts, found := m[name]; found {
			return ts
		}

		if templ, found := findTemplateIn(name, in); found {
			ts := newTestTemplate(templ)
			m[name] = ts
			return ts
		}

		return nil
	}
}

func TestCollectInfo(t *testing.T) {

	configStr := `{ "version": 42 }`

	tests := []struct {
		name      string
		tplString string
		expected  tpl.ParseInfo
	}{
		{"Basic Inner", `{{ .Inner }}`, tpl.ParseInfo{IsInner: true, Config: tpl.DefaultParseConfig}},
		{"Basic config map", "{{ $_hugo_config := `" + configStr + "`  }}", tpl.ParseInfo{Config: tpl.ParseConfig{Version: 42}}},
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
			ts := newTestTemplate(templ)
			ts.typ = templateShortcode
			ctx := newTemplateContext(
				ts,
				newTestTemplateLookup(ts),
			)
			ctx.applyTransformations(templ.Tree.Root)
			c.Assert(ctx.t.parseInfo, qt.DeepEquals, test.expected)
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
			ts := newTestTemplate(templ)
			ctx := newTemplateContext(
				ts,
				newTestTemplateLookup(ts),
			)

			_, err = ctx.applyTransformations(templ.Tree.Root)

			// Just check that it doesn't fail in this test. We have functional tests
			// in hugoblib.
			c.Assert(err, qt.IsNil)

		})
	}

}
