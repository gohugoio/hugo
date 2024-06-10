// Copyright 2017 The Hugo Authors. All rights reserved.
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

package strings

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestFindRE(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	for _, test := range []struct {
		expr    string
		content any
		limit   any
		expect  any
	}{
		{"[G|g]o", "Hugo is a static site generator written in Go.", 2, []string{"go", "Go"}},
		{"[G|g]o", "Hugo is a static site generator written in Go.", -1, []string{"go", "Go"}},
		{"[G|g]o", "Hugo is a static site generator written in Go.", 1, []string{"go"}},
		{"[G|g]o", "Hugo is a static site generator written in Go.", "1", []string{"go"}},
		{"[G|g]o", "Hugo is a static site generator written in Go.", nil, []string(nil)},
		// errors
		{"[G|go", "Hugo is a static site generator written in Go.", nil, false},
		{"[G|g]o", t, nil, false},
	} {
		result, err := ns.FindRE(test.expr, test.content, test.limit)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Check(result, qt.DeepEquals, test.expect)
	}
}

func TestFindRESubmatch(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	for _, test := range []struct {
		expr    string
		content any
		limit   any
		expect  any
	}{
		{`<a\s*href="(.+?)">(.+?)</a>`, `<li><a href="#foo">Foo</a></li><li><a href="#bar">Bar</a></li>`, -1, [][]string{
			{"<a href=\"#foo\">Foo</a>", "#foo", "Foo"},
			{"<a href=\"#bar\">Bar</a>", "#bar", "Bar"},
		}},
		// Some simple cases.
		{"([G|g]o)", "Hugo is a static site generator written in Go.", -1, [][]string{{"go", "go"}, {"Go", "Go"}}},
		{"([G|g]o)", "Hugo is a static site generator written in Go.", 1, [][]string{{"go", "go"}}},

		// errors
		{"([G|go", "Hugo is a static site generator written in Go.", nil, false},
		{"([G|g]o)", t, nil, false},
	} {
		result, err := ns.FindRESubmatch(test.expr, test.content, test.limit)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Check(result, qt.DeepEquals, test.expect)
	}
}

func TestReplaceRE(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	for _, test := range []struct {
		pattern any
		repl    any
		s       any
		n       []any
		expect  any
	}{
		{"^https?://([^/]+).*", "$1", "http://gohugo.io/docs", nil, "gohugo.io"},
		{"^https?://([^/]+).*", "$2", "http://gohugo.io/docs", nil, ""},
		{"(ab)", "AB", "aabbaab", nil, "aABbaAB"},
		{"(ab)", "AB", "aabbaab", []any{1}, "aABbaab"},
		// errors
		{"(ab", "AB", "aabb", nil, false}, // invalid re
		{tstNoStringer{}, "$2", "http://gohugo.io/docs", nil, false},
		{"^https?://([^/]+).*", tstNoStringer{}, "http://gohugo.io/docs", nil, false},
		{"^https?://([^/]+).*", "$2", tstNoStringer{}, nil, false},
	} {

		var (
			result string
			err    error
		)
		if len(test.n) > 0 {
			result, err = ns.ReplaceRE(test.pattern, test.repl, test.s, test.n...)
		} else {
			result, err = ns.ReplaceRE(test.pattern, test.repl, test.s)
		}

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Check(result, qt.Equals, test.expect)
	}
}
