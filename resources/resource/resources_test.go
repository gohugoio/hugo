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

package resource

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestResourcesMount(t *testing.T) {
	c := qt.New(t)
	c.Assert(true, qt.IsTrue)

	var m ResourceGetter
	var r Resources

	check := func(in, expect string) {
		c.Helper()
		r := m.Get(in)
		c.Assert(r, qt.Not(qt.IsNil))
		c.Assert(r.Name(), qt.Equals, expect)
	}

	checkNil := func(in string) {
		c.Helper()
		r := m.Get(in)
		c.Assert(r, qt.IsNil)
	}

	// Misc tests.
	r = Resources{
		testResource{name: "/foo/theme.css"},
	}

	m = r.Mount("/foo", ".")
	check("./theme.css", "/foo/theme.css")

	// Relative target.
	r = Resources{
		testResource{name: "/a/b/c/d.txt"},
		testResource{name: "/a/b/c/e/f.txt"},
		testResource{name: "/a/b/d.txt"},
		testResource{name: "/a/b/e.txt"},
	}

	m = r.Mount("/a/b/c", "z")
	check("z/d.txt", "/a/b/c/d.txt")
	check("z/e/f.txt", "/a/b/c/e/f.txt")

	m = r.Mount("/a/b", "")
	check("d.txt", "/a/b/d.txt")
	m = r.Mount("/a/b", ".")
	check("d.txt", "/a/b/d.txt")
	m = r.Mount("/a/b", "./")
	check("d.txt", "/a/b/d.txt")
	check("./d.txt", "/a/b/d.txt")

	m = r.Mount("/a/b", ".")
	check("./d.txt", "/a/b/d.txt")

	// Absolute target.
	m = r.Mount("/a/b/c", "/z")
	check("/z/d.txt", "/a/b/c/d.txt")
	check("/z/e/f.txt", "/a/b/c/e/f.txt")
	checkNil("/z/f.txt")

	m = r.Mount("/a/b", "/z")
	check("/z/c/d.txt", "/a/b/c/d.txt")
	check("/z/c/e/f.txt", "/a/b/c/e/f.txt")
	check("/z/d.txt", "/a/b/d.txt")
	checkNil("/z/f.txt")

	m = r.Mount("", "")
	check("/a/b/c/d.txt", "/a/b/c/d.txt")
	check("/a/b/c/e/f.txt", "/a/b/c/e/f.txt")
	check("/a/b/d.txt", "/a/b/d.txt")
	checkNil("/a/b/f.txt")

	m = r.Mount("/a/b", "/a/b")
	check("/a/b/c/d.txt", "/a/b/c/d.txt")
	check("/a/b/c/e/f.txt", "/a/b/c/e/f.txt")
	check("/a/b/d.txt", "/a/b/d.txt")
	checkNil("/a/b/f.txt")

	// Resources with relative paths.
	r = Resources{
		testResource{name: "a/b/c/d.txt"},
		testResource{name: "a/b/c/e/f.txt"},
		testResource{name: "a/b/d.txt"},
		testResource{name: "a/b/e.txt"},
		testResource{name: "n.txt"},
	}

	m = r.Mount("a/b", "z")
	check("z/d.txt", "a/b/d.txt")
	checkNil("/z/d.txt")
}

type testResource struct {
	Resource
	name string
}

func (r testResource) Name() string {
	return r.name
}

func (r testResource) NameNormalized() string {
	return r.name
}
