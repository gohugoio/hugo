// Copyright 2021 The Hugo Authors. All rights reserved.
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

package hglob

import (
	"path/filepath"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestResolveRootDir(t *testing.T) {
	c := qt.New(t)

	for _, test := range []struct {
		input    string
		expected string
	}{
		{"data/foo.json", "data"},
		{"a/b/**/foo.json", "a/b"},
		{"dat?a/foo.json", ""},
		{"a/b[a-c]/foo.json", "a"},
	} {
		c.Assert(ResolveRootDir(test.input), qt.Equals, test.expected)
	}
}

func TestFilterGlobParts(t *testing.T) {
	c := qt.New(t)

	for _, test := range []struct {
		input    []string
		expected []string
	}{
		{[]string{"a", "*", "c"}, []string{"a", "c"}},
	} {
		c.Assert(FilterGlobParts(test.input), qt.DeepEquals, test.expected)
	}
}

func TestNormalizePath(t *testing.T) {
	c := qt.New(t)

	for _, test := range []struct {
		input    string
		expected string
	}{
		{filepath.FromSlash("data/FOO.json"), "data/foo.json"},
		{filepath.FromSlash("/data/FOO.json"), "data/foo.json"},
		{filepath.FromSlash("./FOO.json"), "foo.json"},
		{"//", ""},
	} {
		c.Assert(NormalizePath(test.input), qt.Equals, test.expected)
	}
}

func TestGetGlob(t *testing.T) {
	for _, cache := range []*pathGlobCache{defaultGlobCache} {
		c := qt.New(t)
		g, err := cache.GetGlob("**.JSON")
		c.Assert(err, qt.IsNil)
		c.Assert(g.Match("data/my.jSon"), qt.Equals, true)
	}
}

func BenchmarkGetGlob(b *testing.B) {
	runBench := func(name string, cache *pathGlobCache, search string) {
		b.Run(name, func(b *testing.B) {
			g, err := GetGlob("**/foo")
			if err != nil {
				b.Fatal(err)
			}
			for b.Loop() {
				_ = g.Match(search)
			}
		})
	}

	runBench("Default cache", defaultGlobCache, "abcde")
	runBench("Filenames cache, lowercase searches", defaultGlobCache, "abcde")
	runBench("Filenames cache, mixed case searches", defaultGlobCache, "abCDe")

	b.Run("GetGlob", func(b *testing.B) {
		for b.Loop() {
			_, err := GetGlob("**/foo")
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
