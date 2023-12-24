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

package reflect

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

var ns = New()

func TestIsMap(t *testing.T) {
	c := qt.New(t)
	for _, test := range []struct {
		v      any
		expect any
	}{
		{map[int]int{1: 1}, true},
		{"foo", false},
		{nil, false},
	} {
		result := ns.IsMap(test.v)
		c.Assert(result, qt.Equals, test.expect)
	}
}

func TestIsSlice(t *testing.T) {
	c := qt.New(t)
	for _, test := range []struct {
		v      any
		expect any
	}{
		{[]int{1, 2}, true},
		{"foo", false},
		{nil, false},
	} {
		result := ns.IsSlice(test.v)
		c.Assert(result, qt.Equals, test.expect)
	}
}
