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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ns = New()

type tstNoStringer struct{}

func TestIsMap(t *testing.T) {
	for i, test := range []struct {
		v      interface{}
		expect interface{}
	}{
		{map[int]int{1: 1}, true},
		{"foo", false},
		{nil, false},
	} {
		errMsg := fmt.Sprintf("[%d] %v", i, test)
		result := ns.IsMap(test.v)
		assert.Equal(t, test.expect, result, errMsg)
	}
}

func TestIsSlice(t *testing.T) {
	for i, test := range []struct {
		v      interface{}
		expect interface{}
	}{
		{[]int{1, 2}, true},
		{"foo", false},
		{nil, false},
	} {
		errMsg := fmt.Sprintf("[%d] %v", i, test)
		result := ns.IsSlice(test.v)
		assert.Equal(t, test.expect, result, errMsg)
	}
}
