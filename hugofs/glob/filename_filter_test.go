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

package glob

import (
	"path/filepath"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestFilenameFilter(t *testing.T) {
	c := qt.New(t)

	excludeAlmostAllJSON, err := NewFilenameFilter([]string{"/a/b/c/foo.json"}, []string{"**.json"})
	c.Assert(err, qt.IsNil)
	c.Assert(excludeAlmostAllJSON.Match(filepath.FromSlash("/data/my.json"), false), qt.IsFalse)
	c.Assert(excludeAlmostAllJSON.Match(filepath.FromSlash("/a/b/c/foo.json"), false), qt.IsTrue)
	c.Assert(excludeAlmostAllJSON.Match(filepath.FromSlash("/a/b/c/foo.bar"), false), qt.IsFalse)
	c.Assert(excludeAlmostAllJSON.Match(filepath.FromSlash("/a/b/c"), true), qt.IsTrue)
	c.Assert(excludeAlmostAllJSON.Match(filepath.FromSlash("/a/b"), true), qt.IsTrue)
	c.Assert(excludeAlmostAllJSON.Match(filepath.FromSlash("/a"), true), qt.IsTrue)
	c.Assert(excludeAlmostAllJSON.Match(filepath.FromSlash("/"), true), qt.IsTrue)
	c.Assert(excludeAlmostAllJSON.Match("", true), qt.IsTrue)

	excludeAllButFooJSON, err := NewFilenameFilter([]string{"/a/**/foo.json"}, []string{"**.json"})
	c.Assert(excludeAllButFooJSON.Match(filepath.FromSlash("/data/my.json"), false), qt.IsFalse)
	c.Assert(excludeAllButFooJSON.Match(filepath.FromSlash("/a/b/c/d/e/foo.json"), false), qt.IsTrue)
	c.Assert(excludeAllButFooJSON.Match(filepath.FromSlash("/a/b/c"), true), qt.IsTrue)
	c.Assert(excludeAllButFooJSON.Match(filepath.FromSlash("/a/b/"), true), qt.IsTrue)
	c.Assert(excludeAllButFooJSON.Match(filepath.FromSlash("/"), true), qt.IsTrue)
	c.Assert(excludeAllButFooJSON.Match(filepath.FromSlash("/b"), true), qt.IsFalse)

	excludeAllButFooJSONMixedCasePattern, err := NewFilenameFilter([]string{"/**/Foo.json"}, nil)
	c.Assert(excludeAllButFooJSONMixedCasePattern.Match(filepath.FromSlash("/a/b/c/d/e/foo.json"), false), qt.IsTrue)
	c.Assert(excludeAllButFooJSONMixedCasePattern.Match(filepath.FromSlash("/a/b/c/d/e/FOO.json"), false), qt.IsTrue)

	c.Assert(err, qt.IsNil)

	nopFilter, err := NewFilenameFilter(nil, nil)
	c.Assert(err, qt.IsNil)
	c.Assert(nopFilter.Match("ab.txt", false), qt.IsTrue)

	includeOnlyFilter, err := NewFilenameFilter([]string{"**.json", "**.jpg"}, nil)
	c.Assert(err, qt.IsNil)
	c.Assert(includeOnlyFilter.Match("ab.json", false), qt.IsTrue)
	c.Assert(includeOnlyFilter.Match("ab.jpg", false), qt.IsTrue)
	c.Assert(includeOnlyFilter.Match("ab.gif", false), qt.IsFalse)

	excludeOnlyFilter, err := NewFilenameFilter(nil, []string{"**.json", "**.jpg"})
	c.Assert(err, qt.IsNil)
	c.Assert(excludeOnlyFilter.Match("ab.json", false), qt.IsFalse)
	c.Assert(excludeOnlyFilter.Match("ab.jpg", false), qt.IsFalse)
	c.Assert(excludeOnlyFilter.Match("ab.gif", false), qt.IsTrue)

	var nilFilter *FilenameFilter
	c.Assert(nilFilter.Match("ab.gif", false), qt.IsTrue)

	funcFilter := NewFilenameFilterForInclusionFunc(func(s string) bool { return strings.HasSuffix(s, ".json") })
	c.Assert(funcFilter.Match("ab.json", false), qt.IsTrue)
	c.Assert(funcFilter.Match("ab.bson", false), qt.IsFalse)

}
