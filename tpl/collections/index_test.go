// Copyright 2017 The Hugo Authors. All rights reserved.
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
	"fmt"
	"testing"

	"github.com/gohugoio/hugo/common/maps"

	qt "github.com/frankban/quicktest"
)

func TestIndex(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ns := newNs()

	var (
		emptyInterface any
		nilPointer     *int
	)

	for i, test := range []struct {
		item    any
		indices []any
		expect  any
		isErr   bool
	}{
		{[]int{0, 1}, []any{0}, 0, false},
		{[]int{0, 1}, []any{9}, nil, false}, // index out of range
		{[]uint{0, 1}, nil, []uint{0, 1}, false},
		{[][]int{{1, 2}, {3, 4}}, []any{0, 0}, 1, false},
		{map[int]int{1: 10, 2: 20}, []any{1}, 10, false},
		{map[int]int{1: 10, 2: 20}, []any{0}, 0, false},
		{map[string]map[string]string{"a": {"b": "c"}}, []any{"a", "b"}, "c", false},
		{[]map[string]map[string]string{{"a": {"b": "c"}}}, []any{0, "a", "b"}, "c", false},
		{map[string]map[string]any{"a": {"b": []string{"c", "d"}}}, []any{"a", "b", 1}, "d", false},
		{maps.Params{"a": "av"}, []any{"A"}, "av", false},
		{maps.Params{"a": map[string]any{"b": "bv"}}, []any{"A", "B"}, "bv", false},

		// These used to be errors.
		// See issue 10489.
		{nil, nil, nil, false},
		{nil, []any{0}, nil, false},
		{emptyInterface, []any{0}, nil, false},
		{nilPointer, []any{0}, nil, false},

		// errors
		{[]int{0, 1}, []any{"1"}, nil, true},
		{[]int{0, 1}, []any{nil}, nil, true},
		{tstNoStringer{}, []any{0}, nil, true},
	} {

		c.Run(fmt.Sprintf("vararg %d", i), func(c *qt.C) {
			errMsg := qt.Commentf("[%d] %v", i, test)

			result, err := ns.Index(test.item, test.indices...)

			if test.isErr {
				c.Assert(err, qt.Not(qt.IsNil), errMsg)
				return
			}
			c.Assert(err, qt.IsNil, errMsg)
			c.Assert(result, qt.DeepEquals, test.expect, errMsg)
		})

		c.Run(fmt.Sprintf("slice %d", i), func(c *qt.C) {
			errMsg := qt.Commentf("[%d] %v", i, test)

			result, err := ns.Index(test.item, test.indices)

			if test.isErr {
				c.Assert(err, qt.Not(qt.IsNil), errMsg)
				return
			}
			c.Assert(err, qt.IsNil, errMsg)
			c.Assert(result, qt.DeepEquals, test.expect, errMsg)
		})
	}
}
