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

func TestPrepareParams(t *testing.T) {
	tests := []struct {
		input    Params
		expected Params
	}{
		{
			map[string]any{
				"abC": 32,
			},
			Params{
				"abc": 32,
			},
		},
		{
			map[string]any{
				"abC": 32,
				"deF": map[any]any{
					23: "A value",
					24: map[string]any{
						"AbCDe": "A value",
						"eFgHi": "Another value",
					},
				},
				"gHi": map[string]any{
					"J": 25,
				},
				"jKl": map[string]string{
					"M": "26",
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
				"jkl": Params{
					"m": "26",
				},
			},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			// PrepareParams modifies input.
			PrepareParams(test.input)
			if !reflect.DeepEqual(test.expected, test.input) {
				t.Errorf("[%d] Expected\n%#v, got\n%#v\n", i, test.expected, test.input)
			}
		})
	}
}

func TestToSliceStringMap(t *testing.T) {
	c := qt.New(t)

	tests := []struct {
		input    any
		expected []map[string]any
	}{
		{
			input: []map[string]any{
				{"abc": 123},
			},
			expected: []map[string]any{
				{"abc": 123},
			},
		}, {
			input: []any{
				map[string]any{
					"def": 456,
				},
			},
			expected: []map[string]any{
				{"def": 456},
			},
		},
	}

	for _, test := range tests {
		v, err := ToSliceStringMap(test.input)
		c.Assert(err, qt.IsNil)
		c.Assert(v, qt.DeepEquals, test.expected)
	}
}

func TestToParamsAndPrepare(t *testing.T) {
	c := qt.New(t)
	_, err := ToParamsAndPrepare(map[string]any{"A": "av"})
	c.Assert(err, qt.IsNil)

	params, err := ToParamsAndPrepare(nil)
	c.Assert(err, qt.IsNil)
	c.Assert(params, qt.DeepEquals, Params{})
}

func TestRenameKeys(t *testing.T) {
	c := qt.New(t)

	m := map[string]any{
		"a":    32,
		"ren1": "m1",
		"ren2": "m1_2",
		"sub": map[string]any{
			"subsub": map[string]any{
				"REN1": "m2",
				"ren2": "m2_2",
			},
		},
		"no": map[string]any{
			"ren1": "m2",
			"ren2": "m2_2",
		},
	}

	expected := map[string]any{
		"a":    32,
		"new1": "m1",
		"new2": "m1_2",
		"sub": map[string]any{
			"subsub": map[string]any{
				"new1": "m2",
				"ren2": "m2_2",
			},
		},
		"no": map[string]any{
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

func TestLookupEqualFold(t *testing.T) {
	c := qt.New(t)

	m1 := map[string]any{
		"a": "av",
		"B": "bv",
	}

	v, k, found := LookupEqualFold(m1, "b")
	c.Assert(found, qt.IsTrue)
	c.Assert(v, qt.Equals, "bv")
	c.Assert(k, qt.Equals, "B")

	m2 := map[string]string{
		"a": "av",
		"B": "bv",
	}

	v, k, found = LookupEqualFold(m2, "b")
	c.Assert(found, qt.IsTrue)
	c.Assert(k, qt.Equals, "B")
	c.Assert(v, qt.Equals, "bv")
}
