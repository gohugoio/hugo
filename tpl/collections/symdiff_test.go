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

package collections

import (
	"reflect"
	"testing"

	"github.com/gohugoio/hugo/deps"

	qt "github.com/frankban/quicktest"
)

func TestSymDiff(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	ns := New(&deps.Deps{})

	s1 := []TstX{{A: "a"}, {A: "b"}}
	s2 := []TstX{{A: "a"}, {A: "e"}}

	xa, xb, xd, xe := &StructWithSlice{A: "a"}, &StructWithSlice{A: "b"}, &StructWithSlice{A: "d"}, &StructWithSlice{A: "e"}

	sp1 := []*StructWithSlice{xa, xb, xd, xe}
	sp2 := []*StructWithSlice{xb, xe}

	for i, test := range []struct {
		s1       interface{}
		s2       interface{}
		expected interface{}
	}{
		{[]string{"a", "x", "b", "c"}, []string{"a", "b", "y", "c"}, []string{"x", "y"}},
		{[]string{"a", "b", "c"}, []string{"a", "b", "c"}, []string{}},
		{[]interface{}{"a", "b", nil}, []interface{}{"a"}, []interface{}{"b", nil}},
		{[]int{1, 2, 3}, []int{3, 4}, []int{1, 2, 4}},
		{[]int{1, 2, 3}, []int64{3, 4}, []int{1, 2, 4}},
		{s1, s2, []TstX{{A: "b"}, {A: "e"}}},
		{sp1, sp2, []*StructWithSlice{xa, xd}},

		// Errors
		{"error", "error", false},
		{[]int{1, 2, 3}, []string{"3", "4"}, false},
	} {

		errMsg := qt.Commentf("[%d]", i)

		result, err := ns.SymDiff(test.s2, test.s1)

		if b, ok := test.expected.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil), errMsg)
			continue
		}

		c.Assert(err, qt.IsNil, errMsg)

		if !reflect.DeepEqual(test.expected, result) {
			t.Fatalf("%s got\n%T: %v\nexpected\n%T: %v", errMsg, result, result, test.expected, test.expected)
		}
	}

	_, err := ns.Complement()
	c.Assert(err, qt.Not(qt.IsNil))
	_, err = ns.Complement([]string{"a", "b"})
	c.Assert(err, qt.Not(qt.IsNil))

}
