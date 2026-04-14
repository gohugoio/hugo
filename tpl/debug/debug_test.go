// Copyright 2026 The Hugo Authors. All rights reserved.
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

package debug

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

type testMethodStruct struct {
	Name string
}

func (testMethodStruct) Foo() string { return "" }

func TestList(t *testing.T) {
	c := qt.New(t)
	ns := &Namespace{}

	type Embedded struct {
		Inner string
	}

	type testStruct struct {
		Embedded
		Name string
	}

	// Struct: exported non-anonymous fields (including promoted).
	c.Assert(ns.List(testStruct{Name: "x"}), qt.DeepEquals, []string{"Inner", "Name"})

	// Pointer to struct.
	c.Assert(ns.List(&testStruct{Name: "x"}), qt.DeepEquals, []string{"Inner", "Name"})

	// Struct with methods.
	c.Assert(ns.List(testMethodStruct{Name: "x"}), qt.DeepEquals, []string{"Foo()", "Name"})

	// Pointer to struct with methods.
	c.Assert(ns.List(&testMethodStruct{Name: "x"}), qt.DeepEquals, []string{"Foo()", "Name"})

	// Nil pointer.
	c.Assert(ns.List((*testStruct)(nil)), qt.IsNil)

	// Nil.
	c.Assert(ns.List(nil), qt.IsNil)

	// Unsupported type.
	c.Assert(ns.List(42), qt.IsNil)

	// Map.
	c.Assert(ns.List(map[string]int{"b": 1, "a": 2}), qt.DeepEquals, []string{"a", "b"})

	// Slice: order preserved.
	c.Assert(ns.List([]string{"c", "a", "b"}), qt.DeepEquals, []string{"c", "a", "b"})
}
