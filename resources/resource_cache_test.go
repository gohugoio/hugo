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

package resources

import (
	"path/filepath"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestResourceKeyPartitions(t *testing.T) {
	c := qt.New(t)

	for _, test := range []struct {
		input    string
		expected []string
	}{
		{"a.js", []string{"js"}},
		{"a.scss", []string{"sass", "scss"}},
		{"a.sass", []string{"sass", "scss"}},
		{"d/a.js", []string{"d", "js"}},
		{"js/a.js", []string{"js"}},
		{"D/a.JS", []string{"d", "js"}},
		{"d/a", []string{"d"}},
		{filepath.FromSlash("/d/a.js"), []string{"d", "js"}},
		{filepath.FromSlash("/d/e/a.js"), []string{"d", "js"}},
	} {
		c.Assert(ResourceKeyPartitions(test.input), qt.DeepEquals, test.expected, qt.Commentf(test.input))
	}
}

func TestResourceKeyContainsAny(t *testing.T) {
	c := qt.New(t)

	for _, test := range []struct {
		key      string
		filename string
		expected bool
	}{
		{"styles/css", "asdf.css", true},
		{"styles/css", "styles/asdf.scss", true},
		{"js/foo.bar", "asdf.css", false},
	} {
		c.Assert(ResourceKeyContainsAny(test.key, ResourceKeyPartitions(test.filename)), qt.Equals, test.expected)
	}
}
