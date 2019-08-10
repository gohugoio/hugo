// Copyright 2018 The Hugo Authors. All rights reserved.
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

package metadecoders

import (
	"reflect"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestUnmarshalToMap(t *testing.T) {
	c := qt.New(t)

	expect := map[string]interface{}{"a": "b"}

	d := Default

	for i, test := range []struct {
		data   string
		format Format
		expect interface{}
	}{
		{`a = "b"`, TOML, expect},
		{`a: "b"`, YAML, expect},
		// Make sure we get all string keys, even for YAML
		{"a: Easy!\nb:\n  c: 2\n  d: [3, 4]", YAML, map[string]interface{}{"a": "Easy!", "b": map[string]interface{}{"c": 2, "d": []interface{}{3, 4}}}},
		{"a:\n  true: 1\n  false: 2", YAML, map[string]interface{}{"a": map[string]interface{}{"true": 1, "false": 2}}},
		{`{ "a": "b" }`, JSON, expect},
		{`#+a: b`, ORG, expect},
		// errors
		{`a = b`, TOML, false},
		{`a,b,c`, CSV, false}, // Use Unmarshal for CSV
	} {
		msg := qt.Commentf("%d: %s", i, test.format)
		m, err := d.UnmarshalToMap([]byte(test.data), test.format)
		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil), msg)
		} else {
			c.Assert(err, qt.IsNil, msg)
			c.Assert(m, qt.DeepEquals, test.expect, msg)
		}
	}
}

func TestUnmarshalToInterface(t *testing.T) {
	c := qt.New(t)

	expect := map[string]interface{}{"a": "b"}

	d := Default

	for i, test := range []struct {
		data   string
		format Format
		expect interface{}
	}{
		{`[ "Brecker", "Blake", "Redman" ]`, JSON, []interface{}{"Brecker", "Blake", "Redman"}},
		{`{ "a": "b" }`, JSON, expect},
		{`#+a: b`, ORG, expect},
		{`a = "b"`, TOML, expect},
		{`a: "b"`, YAML, expect},
		{`a,b,c`, CSV, [][]string{{"a", "b", "c"}}},
		{"a: Easy!\nb:\n  c: 2\n  d: [3, 4]", YAML, map[string]interface{}{"a": "Easy!", "b": map[string]interface{}{"c": 2, "d": []interface{}{3, 4}}}},
		// errors
		{`a = "`, TOML, false},
	} {
		msg := qt.Commentf("%d: %s", i, test.format)
		m, err := d.Unmarshal([]byte(test.data), test.format)
		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil), msg)
		} else {
			c.Assert(err, qt.IsNil, msg)
			c.Assert(m, qt.DeepEquals, test.expect, msg)
		}

	}

}

func TestUnmarshalStringTo(t *testing.T) {
	c := qt.New(t)

	d := Default

	expectMap := map[string]interface{}{"a": "b"}

	for i, test := range []struct {
		data   string
		to     interface{}
		expect interface{}
	}{
		{"a string", "string", "a string"},
		{`{ "a": "b" }`, make(map[string]interface{}), expectMap},
		{"32", int64(1234), int64(32)},
		{"32", int(1234), int(32)},
		{"3.14159", float64(1), float64(3.14159)},
		{"[3,7,9]", []interface{}{}, []interface{}{3, 7, 9}},
		{"[3.1,7.2,9.3]", []interface{}{}, []interface{}{3.1, 7.2, 9.3}},
	} {
		msg := qt.Commentf("%d: %T", i, test.to)
		m, err := d.UnmarshalStringTo(test.data, test.to)
		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil), msg)
		} else {
			c.Assert(err, qt.IsNil, msg)
			c.Assert(m, qt.DeepEquals, test.expect, msg)
		}

	}
}

func TestStringifyYAMLMapKeys(t *testing.T) {
	cases := []struct {
		input    interface{}
		want     interface{}
		replaced bool
	}{
		{
			map[interface{}]interface{}{"a": 1, "b": 2},
			map[string]interface{}{"a": 1, "b": 2},
			true,
		},
		{
			map[interface{}]interface{}{"a": []interface{}{1, map[interface{}]interface{}{"b": 2}}},
			map[string]interface{}{"a": []interface{}{1, map[string]interface{}{"b": 2}}},
			true,
		},
		{
			map[interface{}]interface{}{true: 1, "b": false},
			map[string]interface{}{"true": 1, "b": false},
			true,
		},
		{
			map[interface{}]interface{}{1: "a", 2: "b"},
			map[string]interface{}{"1": "a", "2": "b"},
			true,
		},
		{
			map[interface{}]interface{}{"a": map[interface{}]interface{}{"b": 1}},
			map[string]interface{}{"a": map[string]interface{}{"b": 1}},
			true,
		},
		{
			map[string]interface{}{"a": map[string]interface{}{"b": 1}},
			map[string]interface{}{"a": map[string]interface{}{"b": 1}},
			false,
		},
		{
			[]interface{}{map[interface{}]interface{}{1: "a", 2: "b"}},
			[]interface{}{map[string]interface{}{"1": "a", "2": "b"}},
			false,
		},
	}

	for i, c := range cases {
		res, replaced := stringifyMapKeys(c.input)

		if c.replaced != replaced {
			t.Fatalf("[%d] Replaced mismatch: %t", i, replaced)
		}
		if !c.replaced {
			res = c.input
		}
		if !reflect.DeepEqual(res, c.want) {
			t.Errorf("[%d] given %q\nwant: %q\n got: %q", i, c.input, c.want, res)
		}
	}
}

func BenchmarkStringifyMapKeysStringsOnlyInterfaceMaps(b *testing.B) {
	maps := make([]map[interface{}]interface{}, b.N)
	for i := 0; i < b.N; i++ {
		maps[i] = map[interface{}]interface{}{
			"a": map[interface{}]interface{}{
				"b": 32,
				"c": 43,
				"d": map[interface{}]interface{}{
					"b": 32,
					"c": 43,
				},
			},
			"b": []interface{}{"a", "b"},
			"c": "d",
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stringifyMapKeys(maps[i])
	}
}

func BenchmarkStringifyMapKeysStringsOnlyStringMaps(b *testing.B) {
	m := map[string]interface{}{
		"a": map[string]interface{}{
			"b": 32,
			"c": 43,
			"d": map[string]interface{}{
				"b": 32,
				"c": 43,
			},
		},
		"b": []interface{}{"a", "b"},
		"c": "d",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stringifyMapKeys(m)
	}
}

func BenchmarkStringifyMapKeysIntegers(b *testing.B) {
	maps := make([]map[interface{}]interface{}, b.N)
	for i := 0; i < b.N; i++ {
		maps[i] = map[interface{}]interface{}{
			1: map[interface{}]interface{}{
				4: 32,
				5: 43,
				6: map[interface{}]interface{}{
					7: 32,
					8: 43,
				},
			},
			2: []interface{}{"a", "b"},
			3: "d",
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stringifyMapKeys(maps[i])
	}
}
