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

	m := map[string]interface{}{
		"string":          "value",
		"first":           1,
		"with_underscore": 2,
		"nested": map[string]interface{}{
			"color": "blue",
			"nestednested": map[string]interface{}{
				"color": "green",
			},
		},
	}

	c := qt.New(t)

	must := func(keyStr, separator string, candidates ...Params) interface{} {
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

}
