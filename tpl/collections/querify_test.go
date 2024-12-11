// Copyright 2024 The Hugo Authors. All rights reserved.
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

package collections

import (
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/common/maps"
)

func TestQuerify(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ns := newNs()

	for _, test := range []struct {
		name   string
		params []any
		expect any
	}{
		// map
		{"01", []any{maps.Params{"a": "foo", "b": "bar"}}, `a=foo&b=bar`},
		{"02", []any{maps.Params{"a": 6, "b": 7}}, `a=6&b=7`},
		{"03", []any{maps.Params{"a": "foo", "b": 7}}, `a=foo&b=7`},
		{"04", []any{map[string]any{"a": "foo", "b": "bar"}}, `a=foo&b=bar`},
		{"05", []any{map[string]any{"a": 6, "b": 7}}, `a=6&b=7`},
		{"06", []any{map[string]any{"a": "foo", "b": 7}}, `a=foo&b=7`},
		// slice
		{"07", []any{[]string{"a", "foo", "b", "bar"}}, `a=foo&b=bar`},
		{"08", []any{[]any{"a", 6, "b", 7}}, `a=6&b=7`},
		{"09", []any{[]any{"a", "foo", "b", 7}}, `a=foo&b=7`},
		// sequence of scalar values
		{"10", []any{"a", "foo", "b", "bar"}, `a=foo&b=bar`},
		{"11", []any{"a", 6, "b", 7}, `a=6&b=7`},
		{"12", []any{"a", "foo", "b", 7}, `a=foo&b=7`},
		// empty map
		{"13", []any{map[string]any{}}, ``},
		// empty slice
		{"14", []any{[]string{}}, ``},
		{"15", []any{[]any{}}, ``},
		// no arguments
		{"16", []any{}, ``},
		// errors: zero key length
		{"17", []any{maps.Params{"": "foo"}}, false},
		{"18", []any{map[string]any{"": "foo"}}, false},
		{"19", []any{[]string{"", "foo"}}, false},
		{"20", []any{[]any{"", 6}}, false},
		{"21", []any{"", "foo"}, false},
		// errors: odd number of values
		{"22", []any{[]string{"a", "foo", "b"}}, false},
		{"23", []any{[]any{"a", 6, "b"}}, false},
		{"24", []any{"a", "foo", "b"}, false},
		// errors: value cannot be cast to string
		{"25", []any{map[string]any{"a": "foo", "b": tstNoStringer{}}}, false},
		{"26", []any{[]any{"a", "foo", "b", tstNoStringer{}}}, false},
		{"27", []any{"a", "foo", "b", tstNoStringer{}}, false},
	} {
		errMsg := qt.Commentf("[%s] %v", test.name, test.params)

		result, err := ns.Querify(test.params...)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil), errMsg)
			continue
		}

		c.Assert(err, qt.IsNil, errMsg)
		c.Assert(result, qt.Equals, test.expect, errMsg)
	}
}

func BenchmarkQuerify(b *testing.B) {
	ns := newNs()
	params := []any{"a", "b", "c", "d", "f", " &"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ns.Querify(params...)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkQuerifySlice(b *testing.B) {
	ns := newNs()
	params := []string{"a", "b", "c", "d", "f", " &"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ns.Querify(params)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkQuerifyMap(b *testing.B) {
	ns := newNs()
	params := map[string]any{"a": "b", "c": "d", "f": " &"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ns.Querify(params)
		if err != nil {
			b.Fatal(err)
		}
	}
}
