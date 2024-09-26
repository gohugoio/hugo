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

	r := Resources{
		testResource{name: "/a/b/c.txt"},
		// testResource{name: "/a/b/d.txt"},
	}

	var m ResourceGetter

	check := func(in, expect string) {
		c.Helper()
		r := m.Get(in)
		c.Assert(r, qt.Not(qt.IsNil))
		c.Assert(r.Name(), qt.Equals, expect)
	}

	m = r.Mount("", "/b")

	check("/b/c.txt", "/a/b/c.txt")
	c.Assert(m.Get("b/c.txt"), qt.IsNil)

	if true {
		return
	}

	m = r.Mount("", "/a")

	check("b/c.txt", "/a/b/c.txt")
	check("/a/b/c.txt", "/a/b/c.txt")

	m = r.Mount("/", "/a/b")
	c.Assert(m.Get("c.txt").Name(), qt.Equals, "/a/b/c.txt")
	c.Assert(m.Get("./c.txt").Name(), qt.Equals, "/a/b/c.txt")
	c.Assert(m.Get("../b/c.txt").Name(), qt.Equals, "/a/b/c.txt")
	c.Assert(m.Get("../b/d.txt").Name(), qt.Equals, "/a/b/d.txt")
	c.Assert(m.Get("../b/e.txt"), qt.IsNil)

	m = r.Mount("/a", "/e")
	c.Assert(m.Get("/e/b/c.txt").Name(), qt.Equals, "/a/b/c.txt")

	m = r.Mount("/", "/a/b/")
	c.Assert(m.Get("c.txt").Name(), qt.Equals, "/a/b/c.txt")
	c.Assert(m.Get("./c.txt").Name(), qt.Equals, "/a/b/c.txt")
	c.Assert(m.Get("../b/c.txt").Name(), qt.Equals, "/a/b/c.txt")
	c.Assert(m.Get("../b/d.txt").Name(), qt.Equals, "/a/b/d.txt")
	c.Assert(m.Get("../b/e.txt"), qt.IsNil)
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
