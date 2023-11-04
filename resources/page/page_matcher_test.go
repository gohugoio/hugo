// Copyright 2020 The Hugo Authors. All rights reserved.
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

package page

import (
	"path/filepath"
	"testing"

	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/common/maps"

	qt "github.com/frankban/quicktest"
)

func TestPageMatcher(t *testing.T) {
	c := qt.New(t)
	developmentTestSite := testSite{h: hugo.NewInfo(testConfig{environment: "development"}, nil)}
	productionTestSite := testSite{h: hugo.NewInfo(testConfig{environment: "production"}, nil)}

	p1, p2, p3 :=
		&testPage{path: "/p1", kind: "section", lang: "en", site: developmentTestSite},
		&testPage{path: "p2", kind: "page", lang: "no", site: productionTestSite},
		&testPage{path: "p3", kind: "page", lang: "en"}

	c.Run("Matches", func(c *qt.C) {
		m := PageMatcher{Kind: "section"}

		c.Assert(m.Matches(p1), qt.Equals, true)
		c.Assert(m.Matches(p2), qt.Equals, false)

		m = PageMatcher{Kind: "page"}
		c.Assert(m.Matches(p1), qt.Equals, false)
		c.Assert(m.Matches(p2), qt.Equals, true)
		c.Assert(m.Matches(p3), qt.Equals, true)

		m = PageMatcher{Kind: "page", Path: "/p2"}
		c.Assert(m.Matches(p1), qt.Equals, false)
		c.Assert(m.Matches(p2), qt.Equals, true)
		c.Assert(m.Matches(p3), qt.Equals, false)

		m = PageMatcher{Path: "/p*"}
		c.Assert(m.Matches(p1), qt.Equals, true)
		c.Assert(m.Matches(p2), qt.Equals, true)
		c.Assert(m.Matches(p3), qt.Equals, true)

		m = PageMatcher{Lang: "en"}
		c.Assert(m.Matches(p1), qt.Equals, true)
		c.Assert(m.Matches(p2), qt.Equals, false)
		c.Assert(m.Matches(p3), qt.Equals, true)

		m = PageMatcher{Environment: "development"}
		c.Assert(m.Matches(p1), qt.Equals, true)
		c.Assert(m.Matches(p2), qt.Equals, false)
		c.Assert(m.Matches(p3), qt.Equals, false)

		m = PageMatcher{Environment: "production"}
		c.Assert(m.Matches(p1), qt.Equals, false)
		c.Assert(m.Matches(p2), qt.Equals, true)
		c.Assert(m.Matches(p3), qt.Equals, false)
	})

	c.Run("Decode", func(c *qt.C) {
		var v PageMatcher
		c.Assert(decodePageMatcher(map[string]any{"kind": "foo"}, &v), qt.Not(qt.IsNil))
		c.Assert(decodePageMatcher(map[string]any{"kind": "{foo,bar}"}, &v), qt.Not(qt.IsNil))
		c.Assert(decodePageMatcher(map[string]any{"kind": "taxonomy"}, &v), qt.IsNil)
		c.Assert(decodePageMatcher(map[string]any{"kind": "{taxonomy,foo}"}, &v), qt.IsNil)
		c.Assert(decodePageMatcher(map[string]any{"kind": "{taxonomy,term}"}, &v), qt.IsNil)
		c.Assert(decodePageMatcher(map[string]any{"kind": "*"}, &v), qt.IsNil)
		c.Assert(decodePageMatcher(map[string]any{"kind": "home", "path": filepath.FromSlash("/a/b/**")}, &v), qt.IsNil)
		c.Assert(v, qt.Equals, PageMatcher{Kind: "home", Path: "/a/b/**"})
	})

	c.Run("mapToPageMatcherParamsConfig", func(c *qt.C) {
		fn := func(m map[string]any) PageMatcherParamsConfig {
			v, err := mapToPageMatcherParamsConfig(m)
			c.Assert(err, qt.IsNil)
			return v
		}
		// Legacy.
		c.Assert(fn(map[string]any{"_target": map[string]any{"kind": "page"}, "foo": "bar"}), qt.DeepEquals, PageMatcherParamsConfig{
			Params: maps.Params{
				"foo": "bar",
			},
			Target: PageMatcher{Path: "", Kind: "page", Lang: "", Environment: ""},
		})

		// Current format.
		c.Assert(fn(map[string]any{"target": map[string]any{"kind": "page"}, "params": map[string]any{"foo": "bar"}}), qt.DeepEquals, PageMatcherParamsConfig{
			Params: maps.Params{
				"foo": "bar",
			},
			Target: PageMatcher{Path: "", Kind: "page", Lang: "", Environment: ""},
		})
	})
}

func TestDecodeCascadeConfig(t *testing.T) {
	c := qt.New(t)

	in := []map[string]any{
		{
			"params": map[string]any{
				"a": "av",
			},
			"target": map[string]any{
				"kind":        "page",
				"Environment": "production",
			},
		},
		{
			"params": map[string]any{
				"b": "bv",
			},
			"target": map[string]any{
				"kind": "page",
			},
		},
	}

	got, err := DecodeCascadeConfig(in)

	c.Assert(err, qt.IsNil)
	c.Assert(got, qt.IsNotNil)
	c.Assert(got.Config, qt.DeepEquals,
		map[PageMatcher]maps.Params{
			{Path: "", Kind: "page", Lang: "", Environment: ""}: {
				"b": "bv",
			},
			{Path: "", Kind: "page", Lang: "", Environment: "production"}: {
				"a": "av",
			},
		},
	)
	c.Assert(got.SourceStructure, qt.DeepEquals, []PageMatcherParamsConfig{
		{
			Params: maps.Params{"a": string("av")},
			Target: PageMatcher{Kind: "page", Environment: "production"},
		},
		{Params: maps.Params{"b": string("bv")}, Target: PageMatcher{Kind: "page"}},
	})

	got, err = DecodeCascadeConfig(nil)
	c.Assert(err, qt.IsNil)
	c.Assert(got, qt.IsNotNil)

}

type testConfig struct {
	environment string
	running     bool
	workingDir  string
}

func (c testConfig) Environment() string {
	return c.environment
}

func (c testConfig) Running() bool {
	return c.running
}

func (c testConfig) WorkingDir() string {
	return c.workingDir
}
