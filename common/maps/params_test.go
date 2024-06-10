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

package maps

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestGetNestedParam(t *testing.T) {
	m := map[string]any{
		"string":          "value",
		"first":           1,
		"with_underscore": 2,
		"nested": map[string]any{
			"color": "blue",
			"nestednested": map[string]any{
				"color": "green",
			},
		},
	}

	c := qt.New(t)

	must := func(keyStr, separator string, candidates ...Params) any {
		v, err := GetNestedParam(keyStr, separator, candidates...)
		c.Assert(err, qt.IsNil)
		return v
	}

	c.Assert(must("first", "_", m), qt.Equals, 1)
	c.Assert(must("First", "_", m), qt.Equals, 1)
	c.Assert(must("with_underscore", "_", m), qt.Equals, 2)
	c.Assert(must("nested_color", "_", m), qt.Equals, "blue")
	c.Assert(must("nested.nestednested.color", ".", m), qt.Equals, "green")
	c.Assert(must("string.name", ".", m), qt.IsNil)
	c.Assert(must("nested.foo", ".", m), qt.IsNil)
}

// https://github.com/gohugoio/hugo/issues/7903
func TestGetNestedParamFnNestedNewKey(t *testing.T) {
	c := qt.New(t)

	nested := map[string]any{
		"color": "blue",
	}
	m := map[string]any{
		"nested": nested,
	}

	existing, nestedKey, owner, err := GetNestedParamFn("nested.new", ".", func(key string) any {
		return m[key]
	})

	c.Assert(err, qt.IsNil)
	c.Assert(existing, qt.IsNil)
	c.Assert(nestedKey, qt.Equals, "new")
	c.Assert(owner, qt.DeepEquals, nested)
}

func TestParamsSetAndMerge(t *testing.T) {
	c := qt.New(t)

	createParamsPair := func() (Params, Params) {
		p1 := Params{"a": "av", "c": "cv", "nested": Params{"al2": "al2v", "cl2": "cl2v"}}
		p2 := Params{"b": "bv", "a": "abv", "nested": Params{"bl2": "bl2v", "al2": "al2bv"}, MergeStrategyKey: ParamsMergeStrategyDeep}
		return p1, p2
	}

	p1, p2 := createParamsPair()

	SetParams(p1, p2)

	c.Assert(p1, qt.DeepEquals, Params{
		"a": "abv",
		"c": "cv",
		"nested": Params{
			"al2": "al2bv",
			"cl2": "cl2v",
			"bl2": "bl2v",
		},
		"b":              "bv",
		MergeStrategyKey: ParamsMergeStrategyDeep,
	})

	p1, p2 = createParamsPair()

	MergeParamsWithStrategy("", p1, p2)

	// Default is to do a shallow merge.
	c.Assert(p1, qt.DeepEquals, Params{
		"c": "cv",
		"nested": Params{
			"al2": "al2v",
			"cl2": "cl2v",
		},
		"b": "bv",
		"a": "av",
	})

	p1, p2 = createParamsPair()
	p1.SetMergeStrategy(ParamsMergeStrategyNone)
	MergeParamsWithStrategy("", p1, p2)
	p1.DeleteMergeStrategy()

	c.Assert(p1, qt.DeepEquals, Params{
		"a": "av",
		"c": "cv",
		"nested": Params{
			"al2": "al2v",
			"cl2": "cl2v",
		},
	})

	p1, p2 = createParamsPair()
	p1.SetMergeStrategy(ParamsMergeStrategyShallow)
	MergeParamsWithStrategy("", p1, p2)
	p1.DeleteMergeStrategy()

	c.Assert(p1, qt.DeepEquals, Params{
		"a": "av",
		"c": "cv",
		"nested": Params{
			"al2": "al2v",
			"cl2": "cl2v",
		},
		"b": "bv",
	})

	p1, p2 = createParamsPair()
	p1.SetMergeStrategy(ParamsMergeStrategyDeep)
	MergeParamsWithStrategy("", p1, p2)
	p1.DeleteMergeStrategy()

	c.Assert(p1, qt.DeepEquals, Params{
		"nested": Params{
			"al2": "al2v",
			"cl2": "cl2v",
			"bl2": "bl2v",
		},
		"b": "bv",
		"a": "av",
		"c": "cv",
	})
}

func TestParamsIsZero(t *testing.T) {
	c := qt.New(t)

	var nilParams Params

	c.Assert(Params{}.IsZero(), qt.IsTrue)
	c.Assert(nilParams.IsZero(), qt.IsTrue)
	c.Assert(Params{"foo": "bar"}.IsZero(), qt.IsFalse)
	c.Assert(Params{"_merge": "foo", "foo": "bar"}.IsZero(), qt.IsFalse)
	c.Assert(Params{"_merge": "foo"}.IsZero(), qt.IsTrue)
}
