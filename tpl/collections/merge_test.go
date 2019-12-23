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
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/gohugoio/hugo/common/maps"

	"github.com/gohugoio/hugo/parser"

	"github.com/gohugoio/hugo/parser/metadecoders"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/deps"
)

func TestMerge(t *testing.T) {

	ns := New(&deps.Deps{})

	simpleMap := map[string]interface{}{"a": 1, "b": 2}

	for i, test := range []struct {
		name   string
		dst    interface{}
		src    interface{}
		expect interface{}
		isErr  bool
	}{
		{
			"basic",
			map[string]interface{}{"a": 1, "b": 2},
			map[string]interface{}{"a": 42, "c": 3},
			map[string]interface{}{"a": 1, "b": 2, "c": 3}, false},
		{
			"basic case insensitive",
			map[string]interface{}{"a": 1, "b": 2},
			map[string]interface{}{"A": 42, "c": 3},
			map[string]interface{}{"a": 1, "b": 2, "c": 3}, false},
		{
			"nested",
			map[string]interface{}{"a": 1, "b": map[string]interface{}{"d": 1, "e": 2}},
			map[string]interface{}{"a": 42, "c": 3, "b": map[string]interface{}{"d": 55, "e": 66, "f": 3}},
			map[string]interface{}{"a": 1, "b": map[string]interface{}{"d": 1, "e": 2, "f": 3}, "c": 3}, false},
		{
			// https://github.com/gohugoio/hugo/issues/6633
			"params dst",
			maps.Params{"a": 1, "b": 2},
			map[string]interface{}{"a": 42, "c": 3},
			maps.Params{"a": int(1), "b": int(2), "c": int(3)}, false},
		{
			"params dst, upper case src",
			maps.Params{"a": 1, "b": 2},
			map[string]interface{}{"a": 42, "C": 3},
			maps.Params{"a": int(1), "b": int(2), "c": int(3)}, false},
		{
			"params src",
			map[string]interface{}{"a": 1, "c": 2},
			maps.Params{"a": 42, "c": 3},
			map[string]interface{}{"a": int(1), "c": int(2)}, false},
		{
			"params src, upper case dst",
			map[string]interface{}{"a": 1, "C": 2},
			maps.Params{"a": 42, "c": 3},
			map[string]interface{}{"a": int(1), "C": int(2)}, false},
		{
			"nested, params dst",
			maps.Params{"a": 1, "b": maps.Params{"d": 1, "e": 2}},
			map[string]interface{}{"a": 42, "c": 3, "b": map[string]interface{}{"d": 55, "e": 66, "f": 3}},
			maps.Params{"a": 1, "b": maps.Params{"d": 1, "e": 2, "f": 3}, "c": 3}, false},
		{"src nil", simpleMap, nil, simpleMap, false},
		// Error cases.
		{"dst not a map", "not a map", nil, nil, true},
		{"src not a map", simpleMap, "not a map", nil, true},
		{"different map types", simpleMap, map[int]interface{}{32: "a"}, nil, true},
		{"all nil", nil, nil, nil, true},
	} {

		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			errMsg := qt.Commentf("[%d] %v", i, test)

			c := qt.New(t)

			srcStr, dstStr := fmt.Sprint(test.src), fmt.Sprint(test.dst)

			result, err := ns.Merge(test.src, test.dst)

			if test.isErr {
				c.Assert(err, qt.Not(qt.IsNil), errMsg)
				return
			}

			c.Assert(err, qt.IsNil)
			c.Assert(result, qt.DeepEquals, test.expect, errMsg)

			// map sort in fmt was fixed in go 1.12.
			if !strings.HasPrefix(runtime.Version(), "go1.11") {
				// Verify that the original maps are preserved.
				c.Assert(fmt.Sprint(test.src), qt.Equals, srcStr)
				c.Assert(fmt.Sprint(test.dst), qt.Equals, dstStr)
			}

		})
	}
}

func TestMergeDataFormats(t *testing.T) {
	c := qt.New(t)
	ns := New(&deps.Deps{})

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
			map[string]interface{}{
				"V1": "v1_1", "V2": "v2_2",
				"V2s": map[string]interface{}{"V21": "v21_1", "V22": "v22_2"}})
	}

}

func TestCaseInsensitiveMapLookup(t *testing.T) {
	c := qt.New(t)

	m1 := reflect.ValueOf(map[string]interface{}{
		"a": 1,
		"B": 2,
	})

	m2 := reflect.ValueOf(map[int]interface{}{
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
