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

type StructWithSlice struct {
	A string
	B []string
}

type StructWithSlicePointers []*StructWithSlice

func TestComplement(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	ns := New(&deps.Deps{})

	s1 := []TstX{{A: "a"}, {A: "b"}, {A: "d"}, {A: "e"}}
	s2 := []TstX{{A: "b"}, {A: "e"}}

	xa, xb, xd, xe := &StructWithSlice{A: "a"}, &StructWithSlice{A: "b"}, &StructWithSlice{A: "d"}, &StructWithSlice{A: "e"}

	sp1 := []*StructWithSlice{xa, xb, xd, xe}
	sp2 := []*StructWithSlice{xb, xe}

	sp1_2 := StructWithSlicePointers{xa, xb, xd, xe}
	sp2_2 := StructWithSlicePointers{xb, xe}

	for i, test := range []struct {
		s        interface{}
		t        []interface{}
		expected interface{}
	}{
		{[]string{"a", "b", "c"}, []interface{}{[]string{"c", "d"}}, []string{"a", "b"}},
		{[]string{"a", "b", "c"}, []interface{}{[]string{"c", "d"}, []string{"a", "b"}}, []string{}},
		{[]interface{}{"a", "b", nil}, []interface{}{[]string{"a", "d"}}, []interface{}{"b", nil}},
		{[]int{1, 2, 3, 4, 5}, []interface{}{[]int{1, 3}, []string{"a", "b"}, []int{1, 2}}, []int{4, 5}},
		{[]int{1, 2, 3, 4, 5}, []interface{}{[]int64{1, 3}}, []int{2, 4, 5}},
		{s1, []interface{}{s2}, []TstX{{A: "a"}, {A: "d"}}},
		{sp1, []interface{}{sp2}, []*StructWithSlice{xa, xd}},
		{sp1_2, []interface{}{sp2_2}, StructWithSlicePointers{xa, xd}},

		// Errors
		{[]string{"a", "b", "c"}, []interface{}{"error"}, false},
		{"error", []interface{}{[]string{"c", "d"}, []string{"a", "b"}}, false},
		{[]string{"a", "b", "c"}, []interface{}{[][]string{{"c", "d"}}}, false},
		{
			[]interface{}{[][]string{{"c", "d"}}}, []interface{}{[]string{"c", "d"}, []string{"a", "b"}},
			[]interface{}{[][]string{{"c", "d"}}},
		},
	} {

		errMsg := qt.Commentf("[%d]", i)

		args := append(test.t, test.s)

		result, err := ns.Complement(args...)

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
