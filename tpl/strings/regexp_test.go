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
		content interface{}
		limit   interface{}
		expect  interface{}
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
		c.Assert(result, qt.DeepEquals, test.expect)
	}
}

func TestReplaceRE(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	for _, test := range []struct {
		pattern interface{}
		repl    interface{}
		s       interface{}
		expect  interface{}
	}{
		{"^https?://([^/]+).*", "$1", "http://gohugo.io/docs", "gohugo.io"},
		{"^https?://([^/]+).*", "$2", "http://gohugo.io/docs", ""},
		{"(ab)", "AB", "aabbaab", "aABbaAB"},
		// errors
		{"(ab", "AB", "aabb", false}, // invalid re
		{tstNoStringer{}, "$2", "http://gohugo.io/docs", false},
		{"^https?://([^/]+).*", tstNoStringer{}, "http://gohugo.io/docs", false},
		{"^https?://([^/]+).*", "$2", tstNoStringer{}, false},
	} {

		result, err := ns.ReplaceRE(test.pattern, test.repl, test.s)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
	}
}
