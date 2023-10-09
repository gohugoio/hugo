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

package collections

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/parser"
	"github.com/gohugoio/hugo/parser/metadecoders"

	qt "github.com/frankban/quicktest"
)

func TestMerge(t *testing.T) {
	ns := newNs()

	simpleMap := map[string]any{"a": 1, "b": 2}

	for i, test := range []struct {
		name   string
		params []any
		expect any
		isErr  bool
	}{
		{
			"basic",
			[]any{
				map[string]any{"a": 42, "c": 3},
				map[string]any{"a": 1, "b": 2},
			},
			map[string]any{"a": 1, "b": 2, "c": 3},
			false,
		},
		{
			"multi",
			[]any{
				map[string]any{"a": 42, "c": 3, "e": 11},
				map[string]any{"a": 1, "b": 2},
				map[string]any{"a": 9, "c": 4, "d": 7},
			},
			map[string]any{"a": 9, "b": 2, "c": 4, "d": 7, "e": 11},
			false,
		},
		{
			"basic case insensitive",
			[]any{
				map[string]any{"A": 42, "c": 3},
				map[string]any{"a": 1, "b": 2},
			},
			map[string]any{"a": 1, "b": 2, "c": 3},
			false,
		},
		{
			"nested",
			[]any{
				map[string]any{"a": 42, "c": 3, "b": map[string]any{"d": 55, "e": 66, "f": 3}},
				map[string]any{"a": 1, "b": map[string]any{"d": 1, "e": 2}},
			},
			map[string]any{"a": 1, "b": map[string]any{"d": 1, "e": 2, "f": 3}, "c": 3},
			false,
		},
		{
			// https://github.com/gohugoio/hugo/issues/6633
			"params dst",
			[]any{
				map[string]any{"a": 42, "c": 3},
				maps.Params{"a": 1, "b": 2},
			},
			maps.Params{"a": int(1), "b": int(2), "c": int(3)},
			false,
		},
		{
			"params dst, upper case src",
			[]any{
				map[string]any{"a": 42, "C": 3},
				maps.Params{"a": 1, "b": 2},
			},
			maps.Params{"a": int(1), "b": int(2), "c": int(3)},
			false,
		},
		{
			"params src",
			[]any{
				maps.Params{"a": 42, "c": 3},
				map[string]any{"a": 1, "c": 2},
			},
			map[string]any{"a": int(1), "c": int(2)},
			false,
		},
		{
			"params src, upper case dst",
			[]any{
				maps.Params{"a": 42, "c": 3},
				map[string]any{"a": 1, "C": 2},
			},
			map[string]any{"a": int(1), "C": int(2)},
			false,
		},
		{
			"nested, params dst",
			[]any{
				map[string]any{"a": 42, "c": 3, "b": map[string]any{"d": 55, "e": 66, "f": 3}},
				maps.Params{"a": 1, "b": maps.Params{"d": 1, "e": 2}},
			},
			maps.Params{"a": 1, "b": maps.Params{"d": 1, "e": 2, "f": 3}, "c": 3},
			false,
		},
		{
			// https://github.com/gohugoio/hugo/issues/7899
			"matching keys with non-map src value",
			[]any{
				map[string]any{"k": "v"},
				map[string]any{"k": map[string]any{"k2": "v2"}},
			},
			map[string]any{"k": map[string]any{"k2": "v2"}},
			false,
		},
		{"src nil", []any{nil, simpleMap}, simpleMap, false},
		// Error cases.
		{"dst not a map", []any{nil, "not a map"}, nil, true},
		{"src not a map", []any{"not a map", simpleMap}, nil, true},
		{"different map types", []any{map[int]any{32: "a"}, simpleMap}, nil, true},
		{"all nil", []any{nil, nil}, nil, true},
	} {

		test := test
		i := i

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			errMsg := qt.Commentf("[%d] %v", i, test)

			c := qt.New(t)

			result, err := ns.Merge(test.params...)

			if test.isErr {
				c.Assert(err, qt.Not(qt.IsNil), errMsg)
				return
			}

			c.Assert(err, qt.IsNil)
			c.Assert(result, qt.DeepEquals, test.expect, errMsg)
		})
	}
}

func TestMergeDataFormats(t *testing.T) {
	c := qt.New(t)
	ns := newNs()

	toml1 := `
V1 = "v1_1"

[V2s]
V21 = "v21_1"

`

	toml2 := `
V1 = "v1_2"
V2 = "v2_2"

[V2s]
V21 = "v21_2"
V22 = "v22_2"

`

	meta1, err := metadecoders.Default.UnmarshalToMap([]byte(toml1), metadecoders.TOML)
	c.Assert(err, qt.IsNil)
	meta2, err := metadecoders.Default.UnmarshalToMap([]byte(toml2), metadecoders.TOML)
	c.Assert(err, qt.IsNil)

	for _, format := range []metadecoders.Format{metadecoders.JSON, metadecoders.YAML, metadecoders.TOML} {

		var dataStr1, dataStr2 bytes.Buffer
		err = parser.InterfaceToConfig(meta1, format, &dataStr1)
		c.Assert(err, qt.IsNil)
		err = parser.InterfaceToConfig(meta2, format, &dataStr2)
		c.Assert(err, qt.IsNil)

		dst, err := metadecoders.Default.UnmarshalToMap(dataStr1.Bytes(), format)
		c.Assert(err, qt.IsNil)
		src, err := metadecoders.Default.UnmarshalToMap(dataStr2.Bytes(), format)
		c.Assert(err, qt.IsNil)

		merged, err := ns.Merge(src, dst)
		c.Assert(err, qt.IsNil)

		c.Assert(
			merged,
			qt.DeepEquals,
			map[string]any{
				"V1": "v1_1", "V2": "v2_2",
				"V2s": map[string]any{"V21": "v21_1", "V22": "v22_2"},
			})
	}
}

func TestCaseInsensitiveMapLookup(t *testing.T) {
	c := qt.New(t)

	m1 := reflect.ValueOf(map[string]any{
		"a": 1,
		"B": 2,
	})

	m2 := reflect.ValueOf(map[int]any{
		1: 1,
		2: 2,
	})

	var found bool

	a, found := caseInsensitiveLookup(m1, reflect.ValueOf("A"))
	c.Assert(found, qt.Equals, true)
	c.Assert(a.Interface(), qt.Equals, 1)

	b, found := caseInsensitiveLookup(m1, reflect.ValueOf("b"))
	c.Assert(found, qt.Equals, true)
	c.Assert(b.Interface(), qt.Equals, 2)

	two, found := caseInsensitiveLookup(m2, reflect.ValueOf(2))
	c.Assert(found, qt.Equals, true)
	c.Assert(two.Interface(), qt.Equals, 2)
}
