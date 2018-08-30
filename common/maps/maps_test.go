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

package maps

import (
	"reflect"
	"testing"
)

func TestToLower(t *testing.T) {

	tests := []struct {
		input    map[string]interface{}
		expected map[string]interface{}
	}{
		{
			map[string]interface{}{
				"abC": 32,
			},
			map[string]interface{}{
				"abc": 32,
			},
		},
		{
			map[string]interface{}{
				"abC": 32,
				"deF": map[interface{}]interface{}{
					23: "A value",
					24: map[string]interface{}{
						"AbCDe": "A value",
						"eFgHi": "Another value",
					},
				},
				"gHi": map[string]interface{}{
					"J": 25,
				},
			},
			map[string]interface{}{
				"abc": 32,
				"def": map[string]interface{}{
					"23": "A value",
					"24": map[string]interface{}{
						"abcde": "A value",
						"efghi": "Another value",
					},
				},
				"ghi": map[string]interface{}{
					"j": 25,
				},
			},
		},
	}

	for i, test := range tests {
		// ToLower modifies input.
		ToLower(test.input)
		if !reflect.DeepEqual(test.expected, test.input) {
			t.Errorf("[%d] Expected\n%#v, got\n%#v\n", i, test.expected, test.input)
		}
	}
}
