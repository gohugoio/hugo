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

	for _, test := range []struct {
		start    interface{}
		addend   []interface{}
		expected interface{}
	}{
		{[]string{"a", "b"}, []interface{}{"c"}, []string{"a", "b", "c"}},
		{[]string{"a", "b"}, []interface{}{"c", "d", "e"}, []string{"a", "b", "c", "d", "e"}},
		{[]string{"a", "b"}, []interface{}{[]string{"c", "d", "e"}}, []string{"a", "b", "c", "d", "e"}},
		{[]string{"a"}, []interface{}{"b", template.HTML("c")}, []interface{}{"a", "b", template.HTML("c")}},
		{nil, []interface{}{"a", "b"}, []string{"a", "b"}},
		{nil, []interface{}{nil}, []interface{}{nil}},
		{[]interface{}{}, []interface{}{[]string{"c", "d", "e"}}, []string{"c", "d", "e"}},
		{tstSlicers{&tstSlicer{"a"}, &tstSlicer{"b"}},
			[]interface{}{&tstSlicer{"c"}},
			tstSlicers{&tstSlicer{"a"}, &tstSlicer{"b"}, &tstSlicer{"c"}}},
		{&tstSlicers{&tstSlicer{"a"}, &tstSlicer{"b"}},
			[]interface{}{&tstSlicer{"c"}},
			tstSlicers{&tstSlicer{"a"},
				&tstSlicer{"b"},
				&tstSlicer{"c"}}},
		{testSlicerInterfaces{&tstSlicerIn1{"a"}, &tstSlicerIn1{"b"}},
			[]interface{}{&tstSlicerIn1{"c"}},
			testSlicerInterfaces{&tstSlicerIn1{"a"}, &tstSlicerIn1{"b"}, &tstSlicerIn1{"c"}}},
		//https://github.com/gohugoio/hugo/issues/5361
		{[]string{"a", "b"}, []interface{}{tstSlicers{&tstSlicer{"a"}, &tstSlicer{"b"}}},
			[]interface{}{"a", "b", &tstSlicer{"a"}, &tstSlicer{"b"}}},
		{[]string{"a", "b"}, []interface{}{&tstSlicer{"a"}},
			[]interface{}{"a", "b", &tstSlicer{"a"}}},
		// Errors
		{"", []interface{}{[]string{"a", "b"}}, false},
		// No string concatenation.
		{"ab",
			[]interface{}{"c"},
			false},
	} {

		result, err := Append(test.start, test.addend...)

		if b, ok := test.expected.(bool); ok && !b {

			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.DeepEquals, test.expected)
	}

}
