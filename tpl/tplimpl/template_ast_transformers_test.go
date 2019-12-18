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
	"github.com/gohugoio/hugo/hugofs/files"

	"testing"

	template "github.com/gohugoio/hugo/tpl/internal/go_templates/htmltemplate"
	"github.com/gohugoio/hugo/tpl/internal/go_templates/texttemplate/parse"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/tpl"
)

// Issue #2927
func TestTransformRecursiveTemplate(t *testing.T) {
	c := qt.New(t)

	recursive := `
{{ define "menu-nodes" }}
{{ template "menu-node" }}
{{ end }}
{{ define "menu-nßode" }}
{{ template "menu-node" }}
{{ end }}
{{ template "menu-nodes" }}
`

	templ, err := template.New("foo").Parse(recursive)
	c.Assert(err, qt.IsNil)
	parseInfo := tpl.DefaultParseInfo

	ctx := newTemplateContext(
		newTemplateInfo("test").(identity.Manager),
		&parseInfo,
		createGetTemplateInfoTree(templ.Tree),
	)
	ctx.applyTransformations(templ.Tree.Root)

}

func createGetTemplateInfoTree(tree *parse.Tree) func(name string) *templateInfoTree {
	return func(name string) *templateInfoTree {
		return &templateInfoTree{
			tree: tree,
		}
	}
}

type I interface {
	Method0()
}

type T struct {
	NonEmptyInterfaceTypedNil I
}

func (T) Method0() {
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
			parseInfo := tpl.DefaultParseInfo

			ctx := newTemplateContext(
				newTemplateInfo("test").(identity.Manager), &parseInfo, createGetTemplateInfoTree(templ.Tree))
			ctx.typ = templateShortcode
			ctx.applyTransformations(templ.Tree.Root)
			c.Assert(ctx.parseInfo, qt.DeepEquals, &test.expected)
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

			_, err = applyTemplateTransformers(
				templatePartial,
				&templateInfoTree{tree: templ.Tree, info: tpl.DefaultParseInfo},
				createGetTemplateInfoTree(templ.Tree))

			// Just check that it doesn't fail in this test. We have functional tests
			// in hugoblib.
			c.Assert(err, qt.IsNil)

		})
	}

}

func newTemplateInfo(name string) tpl.Info {
	return tpl.NewInfo(
		identity.NewManager(identity.NewPathIdentity(files.ComponentFolderLayouts, name)),
		tpl.DefaultParseInfo,
	)
}
