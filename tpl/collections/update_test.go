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

package collections

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestUpdate(t *testing.T) {
	ns := newNs()

	for i, test := range []struct {
		name   string
		params [3]any
		expect any
		isErr  bool
	}{
		{
			"map",
			[3]any{"c", 3, map[string]any{"a": 1, "b": 2}},
			map[string]any{"a": 1, "b": 2, "c": 3},
			false,
		},
		{
			"mapdelete",
			[3]any{"b", nil, map[string]any{"a": 1, "b": 2}},
			map[string]any{"a": 1},
			false,
		},
		{
			"slice",
			[3]any{1, 100, []any{1, 2, 3}},
			[]any{1, 100, 3},
			false,
		},
		{
			"sliceerror",
			[3]any{[]any{100}, 100, []any{1, 2, 3}},
			nil,
			true,
		},
	} {

		test := test
		i := i

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			errMsg := qt.Commentf("[%d] %v", i, test)

			c := qt.New(t)

			_, err := ns.Update(test.params[0], test.params[1], test.params[2])

			if test.isErr {
				c.Assert(err, qt.Not(qt.IsNil), errMsg)
				return
			}

			c.Assert(err, qt.IsNil)
			c.Assert(test.params[2], qt.DeepEquals, test.expect, errMsg)
		})
	}
}
