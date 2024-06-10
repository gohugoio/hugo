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

func TestUnmarshalXML(t *testing.T) {
	c := qt.New(t)

	xmlDoc := `<?xml version="1.0" encoding="utf-8" standalone="yes"?>
	<rss version="2.0"
		xmlns:atom="http://www.w3.org/2005/Atom">
		<channel>
			<title>Example feed</title>
			<link>https://example.com/</link>
			<description>Example feed</description>
			<generator>Hugo -- gohugo.io</generator>
			<language>en-us</language>
			<copyright>Example</copyright>
			<lastBuildDate>Fri, 08 Jan 2021 14:44:10 +0000</lastBuildDate>
			<atom:link href="https://example.com/feed.xml" rel="self" type="application/rss+xml"/>
			<item>
				<title>Example title</title>
				<link>https://example.com/2021/11/30/example-title/</link>
				<pubDate>Tue, 30 Nov 2021 15:00:00 +0000</pubDate>
				<guid>https://example.com/2021/11/30/example-title/</guid>
				<description>Example description</description>
			</item>
		</channel>
	</rss>`

	expect := map[string]any{
		"-atom": "http://www.w3.org/2005/Atom", "-version": "2.0",
		"channel": map[string]any{
			"copyright":   "Example",
			"description": "Example feed",
			"generator":   "Hugo -- gohugo.io",
			"item": map[string]any{
				"description": "Example description",
				"guid":        "https://example.com/2021/11/30/example-title/",
				"link":        "https://example.com/2021/11/30/example-title/",
				"pubDate":     "Tue, 30 Nov 2021 15:00:00 +0000",
				"title":       "Example title",
			},
			"language":      "en-us",
			"lastBuildDate": "Fri, 08 Jan 2021 14:44:10 +0000",
			"link": []any{"https://example.com/", map[string]any{
				"-href": "https://example.com/feed.xml",
				"-rel":  "self",
				"-type": "application/rss+xml",
			}},
			"title": "Example feed",
		},
	}

	d := Default

	m, err := d.Unmarshal([]byte(xmlDoc), XML)
	c.Assert(err, qt.IsNil)
	c.Assert(m, qt.DeepEquals, expect)
}

func TestUnmarshalToMap(t *testing.T) {
	c := qt.New(t)

	expect := map[string]any{"a": "b"}

	d := Default

	for i, test := range []struct {
		data   string
		format Format
		expect any
	}{
		{`a = "b"`, TOML, expect},
		{`a: "b"`, YAML, expect},
		// Make sure we get all string keys, even for YAML
		{"a: Easy!\nb:\n  c: 2\n  d: [3, 4]", YAML, map[string]any{"a": "Easy!", "b": map[string]any{"c": 2, "d": []any{3, 4}}}},
		{"a:\n  true: 1\n  false: 2", YAML, map[string]any{"a": map[string]any{"true": 1, "false": 2}}},
		{`{ "a": "b" }`, JSON, expect},
		{`<root><a>b</a></root>`, XML, expect},
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

	expect := map[string]any{"a": "b"}

	d := Default

	for i, test := range []struct {
		data   []byte
		format Format
		expect any
	}{
		{[]byte(`[ "Brecker", "Blake", "Redman" ]`), JSON, []any{"Brecker", "Blake", "Redman"}},
		{[]byte(`{ "a": "b" }`), JSON, expect},
		{[]byte(``), JSON, map[string]any{}},
		{[]byte(nil), JSON, map[string]any{}},
		{[]byte(`#+a: b`), ORG, expect},
		{[]byte("#+a: foo bar\n#+a: baz"), ORG, map[string]any{"a": []string{string("foo bar"), string("baz")}}},
		{[]byte(`#+DATE: <2020-06-26 Fri>`), ORG, map[string]any{"date": "2020-06-26"}},
		{[]byte(`#+LASTMOD: <2020-06-26 Fri>`), ORG, map[string]any{"lastmod": "2020-06-26"}},
		{[]byte(`#+PUBLISHDATE: <2020-06-26 Fri>`), ORG, map[string]any{"publishdate": "2020-06-26"}},
		{[]byte(`#+EXPIRYDATE: <2020-06-26 Fri>`), ORG, map[string]any{"expirydate": "2020-06-26"}},
		{[]byte(`a = "b"`), TOML, expect},
		{[]byte(`a: "b"`), YAML, expect},
		{[]byte(`<root><a>b</a></root>`), XML, expect},
		{[]byte(`a,b,c`), CSV, [][]string{{"a", "b", "c"}}},
		{[]byte("a: Easy!\nb:\n  c: 2\n  d: [3, 4]"), YAML, map[string]any{"a": "Easy!", "b": map[string]any{"c": 2, "d": []any{3, 4}}}},
		// errors
		{[]byte(`a = "`), TOML, false},
	} {
		msg := qt.Commentf("%d: %s", i, test.format)
		m, err := d.Unmarshal(test.data, test.format)
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

	expectMap := map[string]any{"a": "b"}

	for i, test := range []struct {
		data   string
		to     any
		expect any
	}{
		{"a string", "string", "a string"},
		{`{ "a": "b" }`, make(map[string]any), expectMap},
		{"32", int64(1234), int64(32)},
		{"32", int(1234), int(32)},
		{"3.14159", float64(1), float64(3.14159)},
		{"[3,7,9]", []any{}, []any{3, 7, 9}},
		{"[3.1,7.2,9.3]", []any{}, []any{3.1, 7.2, 9.3}},
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
		input    any
		want     any
		replaced bool
	}{
		{
			map[any]any{"a": 1, "b": 2},
			map[string]any{"a": 1, "b": 2},
			true,
		},
		{
			map[any]any{"a": []any{1, map[any]any{"b": 2}}},
			map[string]any{"a": []any{1, map[string]any{"b": 2}}},
			true,
		},
		{
			map[any]any{true: 1, "b": false},
			map[string]any{"true": 1, "b": false},
			true,
		},
		{
			map[any]any{1: "a", 2: "b"},
			map[string]any{"1": "a", "2": "b"},
			true,
		},
		{
			map[any]any{"a": map[any]any{"b": 1}},
			map[string]any{"a": map[string]any{"b": 1}},
			true,
		},
		{
			map[string]any{"a": map[string]any{"b": 1}},
			map[string]any{"a": map[string]any{"b": 1}},
			false,
		},
		{
			[]any{map[any]any{1: "a", 2: "b"}},
			[]any{map[string]any{"1": "a", "2": "b"}},
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
	maps := make([]map[any]any, b.N)
	for i := 0; i < b.N; i++ {
		maps[i] = map[any]any{
			"a": map[any]any{
				"b": 32,
				"c": 43,
				"d": map[any]any{
					"b": 32,
					"c": 43,
				},
			},
			"b": []any{"a", "b"},
			"c": "d",
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stringifyMapKeys(maps[i])
	}
}

func BenchmarkStringifyMapKeysStringsOnlyStringMaps(b *testing.B) {
	m := map[string]any{
		"a": map[string]any{
			"b": 32,
			"c": 43,
			"d": map[string]any{
				"b": 32,
				"c": 43,
			},
		},
		"b": []any{"a", "b"},
		"c": "d",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stringifyMapKeys(m)
	}
}

func BenchmarkStringifyMapKeysIntegers(b *testing.B) {
	maps := make([]map[any]any, b.N)
	for i := 0; i < b.N; i++ {
		maps[i] = map[any]any{
			1: map[any]any{
				4: 32,
				5: 43,
				6: map[any]any{
					7: 32,
					8: 43,
				},
			},
			2: []any{"a", "b"},
			3: "d",
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stringifyMapKeys(maps[i])
	}
}
