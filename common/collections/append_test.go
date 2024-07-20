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
	"html/template"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestAppend(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	for i, test := range []struct {
		start    any
		addend   []any
		expected any
	}{
		{[]string{"a", "b"}, []any{"c"}, []string{"a", "b", "c"}},
		{[]string{"a", "b"}, []any{"c", "d", "e"}, []string{"a", "b", "c", "d", "e"}},
		{[]string{"a", "b"}, []any{[]string{"c", "d", "e"}}, []string{"a", "b", "c", "d", "e"}},
		{[]string{"a"}, []any{"b", template.HTML("c")}, []any{"a", "b", template.HTML("c")}},
		{nil, []any{"a", "b"}, []string{"a", "b"}},
		{nil, []any{nil}, []any{nil}},
		{[]any{}, []any{[]string{"c", "d", "e"}}, []string{"c", "d", "e"}},
		{
			tstSlicers{&tstSlicer{"a"}, &tstSlicer{"b"}},
			[]any{&tstSlicer{"c"}},
			tstSlicers{&tstSlicer{"a"}, &tstSlicer{"b"}, &tstSlicer{"c"}},
		},
		{
			&tstSlicers{&tstSlicer{"a"}, &tstSlicer{"b"}},
			[]any{&tstSlicer{"c"}},
			tstSlicers{
				&tstSlicer{"a"},
				&tstSlicer{"b"},
				&tstSlicer{"c"},
			},
		},
		{
			testSlicerInterfaces{&tstSlicerIn1{"a"}, &tstSlicerIn1{"b"}},
			[]any{&tstSlicerIn1{"c"}},
			testSlicerInterfaces{&tstSlicerIn1{"a"}, &tstSlicerIn1{"b"}, &tstSlicerIn1{"c"}},
		},
		// https://github.com/gohugoio/hugo/issues/5361
		{
			[]string{"a", "b"},
			[]any{tstSlicers{&tstSlicer{"a"}, &tstSlicer{"b"}}},
			[]any{"a", "b", &tstSlicer{"a"}, &tstSlicer{"b"}},
		},
		{
			[]string{"a", "b"},
			[]any{&tstSlicer{"a"}},
			[]any{"a", "b", &tstSlicer{"a"}},
		},
		// Errors
		{"", []any{[]string{"a", "b"}}, false},
		// No string concatenation.
		{
			"ab",
			[]any{"c"},
			false,
		},
		{[]string{"a", "b"}, []any{nil}, []any{"a", "b", nil}},
		{[]string{"a", "b"}, []any{nil, "d", nil}, []any{"a", "b", nil, "d", nil}},
		{[]any{"a", nil, "c"}, []any{"d", nil, "f"}, []any{"a", nil, "c", "d", nil, "f"}},
	} {

		result, err := Append(test.start, test.addend...)

		if b, ok := test.expected.(bool); ok && !b {

			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.DeepEquals, test.expected, qt.Commentf("test: [%d] %v", i, test))
	}
}

// #11093
func TestAppendToMultiDimensionalSlice(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	for _, test := range []struct {
		to       any
		from     []any
		expected any
	}{
		{
			[][]string{{"a", "b"}},
			[]any{[]string{"c", "d"}},
			[][]string{
				{"a", "b"},
				{"c", "d"},
			},
		},
		{
			[][]string{{"a", "b"}},
			[]any{[]string{"c", "d"}, []string{"e", "f"}},
			[][]string{
				{"a", "b"},
				{"c", "d"},
				{"e", "f"},
			},
		},
		{
			[][]string{{"a", "b"}},
			[]any{[]int{1, 2}},
			false,
		},
	} {
		result, err := Append(test.to, test.from...)
		if b, ok := test.expected.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil))
		} else {
			c.Assert(err, qt.IsNil)
			c.Assert(result, qt.DeepEquals, test.expected)
		}
	}
}

func TestAppendShouldMakeACopyOfTheInputSlice(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	slice := make([]string, 0, 100)
	slice = append(slice, "a", "b")
	result, err := Append(slice, "c")
	c.Assert(err, qt.IsNil)
	slice[0] = "d"
	c.Assert(result, qt.DeepEquals, []string{"a", "b", "c"})
	c.Assert(slice, qt.DeepEquals, []string{"d", "b"})
}
