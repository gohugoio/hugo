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
		{"a: Easy!\nb:\n  c: 2\n  d: [3, 4]", YAML, map[string]any{"a": "Easy!", "b": map[string]any{"c": uint64(2), "d": []any{uint64(3), uint64(4)}}}},
		{"a:\n  true: 1\n  false: 2", YAML, map[string]any{"a": map[string]any{"true": uint64(1), "false": uint64(2)}}},
		{`{ "a": "b" }`, JSON, expect},
		{`<root><a>b</a></root>`, XML, expect},
		{`#+a: b`, ORG, expect},
		// errors
		{`a = b`, TOML, false},
		{`a,b,c`, CSV, false}, // Use Unmarshal for CSV
		{`<root>just a string</root>`, XML, false},
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
		{[]byte(`#+FILETAGS: :work:`), ORG, map[string]any{"filetags": []string{"work"}}},
		{[]byte(`#+FILETAGS: :work:fun:`), ORG, map[string]any{"filetags": []string{"work", "fun"}}},
		{[]byte(`#+PUBLISHDATE: <2020-06-26 Fri>`), ORG, map[string]any{"publishdate": "2020-06-26"}},
		{[]byte(`#+EXPIRYDATE: <2020-06-26 Fri>`), ORG, map[string]any{"expirydate": "2020-06-26"}},
		{[]byte(`a = "b"`), TOML, expect},
		{[]byte(`a: "b"`), YAML, expect},
		{[]byte(`<root><a>b</a></root>`), XML, expect},
		{[]byte(`a,b,c`), CSV, [][]string{{"a", "b", "c"}}},
		{[]byte("a: Easy!\nb:\n  c: 2\n  d: [3, 4]"), YAML, map[string]any{"a": "Easy!", "b": map[string]any{"c": uint64(2), "d": []any{uint64(3), uint64(4)}}}},
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
		{"[3,7,9]", []any{}, []any{uint64(3), uint64(7), uint64(9)}},
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

func TestCalculateAliasLimit(t *testing.T) {
	c := qt.New(t)

	const kb = 1024

	c.Assert(calculateCollectionAliasLimit(0), qt.Equals, 100)
	c.Assert(calculateCollectionAliasLimit(500), qt.Equals, 100)
	c.Assert(calculateCollectionAliasLimit(1*kb), qt.Equals, 100)
	c.Assert(calculateCollectionAliasLimit(2*kb), qt.Equals, 5000)
	c.Assert(calculateCollectionAliasLimit(8*kb), qt.Equals, 5000)
	c.Assert(calculateCollectionAliasLimit(12*kb), qt.Equals, 10000)
	c.Assert(calculateCollectionAliasLimit(10000*kb), qt.Equals, 10000)
}

func BenchmarkDecodeYAMLToMap(b *testing.B) {
	d := Default

	data := []byte(`
a:
  v1: 32
  v2: 43
  v3: "foo"
b:
  - a
  - b
c: "d"

`)

	for i := 0; i < b.N; i++ {
		_, err := d.UnmarshalToMap(data, YAML)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUnmarshalBillionLaughs(b *testing.B) {
	yamlBillionLaughs := []byte(`
a: &a [_, _, _, _, _, _, _, _, _, _, _, _, _, _, _]
b: &b [*a, *a, *a, *a, *a, *a, *a, *a, *a, *a]
c: &c [*b, *b, *b, *b, *b, *b, *b, *b, *b, *b]
d: &d [*c, *c, *c, *c, *c, *c, *c, *c, *c, *c]
e: &e [*d, *d, *d, *d, *d, *d, *d, *d, *d, *d]
f: &f [*e, *e, *e, *e, *e, *e, *e, *e, *e, *e]
g: &g [*f, *f, *f, *f, *f, *f, *f, *f, *f, *f]
h: &h [*g, *g, *g, *g, *g, *g, *g, *g, *g, *g]
i: &i [*h, *h, *h, *h, *h, *h, *h, *h, *h, *h]
`)

	yamlFrontMatter := []byte(`
title: mysect
tags: [tag1, tag2]
params:
  color: blue
`)

	yamlTests := []struct {
		Title                      string
		Content                    []byte
		IsExpectedToFailValidation bool
	}{
		{"Billion Laughs", yamlBillionLaughs, true},
		{"YAML Front Matter", yamlFrontMatter, false},
	}

	for _, tt := range yamlTests {

		b.Run(tt.Title+" no validation", func(b *testing.B) {
			for range b.N {
				var v any
				if err := unmarshalYamlNoValidation(tt.Content, &v); err != nil {
					b.Fatal(err)
				}
			}
		})

		b.Run(tt.Title+" with validation", func(b *testing.B) {
			for range b.N {
				var v any
				err := UnmarshalYaml(tt.Content, &v)
				if tt.IsExpectedToFailValidation {
					if err == nil {
						b.Fatal("expected to fail validation but did not")
					}
				} else {
					if err != nil {
						b.Fatal(err)
					}
				}
			}
		})
	}
}
