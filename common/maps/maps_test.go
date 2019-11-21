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
	"fmt"
	"reflect"
	"testing"

	qt "github.com/frankban/quicktest"
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
			Params{
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
			Params{
				"abc": 32,
				"def": Params{
					"23": "A value",
					"24": Params{
						"abcde": "A value",
						"efghi": "Another value",
					},
				},
				"ghi": Params{
					"j": 25,
				},
			},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			// ToLower modifies input.
			ToLower(test.input)
			if !reflect.DeepEqual(test.expected, test.input) {
				t.Errorf("[%d] Expected\n%#v, got\n%#v\n", i, test.expected, test.input)
			}
		})
	}
}

func TestRenameKeys(t *testing.T) {
	c := qt.New(t)

	m := map[string]interface{}{
		"a":    32,
		"ren1": "m1",
		"ren2": "m1_2",
		"sub": map[string]interface{}{
			"subsub": map[string]interface{}{
				"REN1": "m2",
				"ren2": "m2_2",
			},
		},
		"no": map[string]interface{}{
			"ren1": "m2",
			"ren2": "m2_2",
		},
	}

	expected := map[string]interface{}{
		"a":    32,
		"new1": "m1",
		"new2": "m1_2",
		"sub": map[string]interface{}{
			"subsub": map[string]interface{}{
				"new1": "m2",
				"ren2": "m2_2",
			},
		},
		"no": map[string]interface{}{
			"ren1": "m2",
			"ren2": "m2_2",
		},
	}

	renamer, err := NewKeyRenamer(
		"{ren1,sub/*/ren1}", "new1",
		"{Ren2,sub/ren2}", "new2",
	)
	c.Assert(err, qt.IsNil)

	renamer.Rename(m)

	if !reflect.DeepEqual(expected, m) {
		t.Errorf("Expected\n%#v, got\n%#v\n", expected, m)
	}

}
