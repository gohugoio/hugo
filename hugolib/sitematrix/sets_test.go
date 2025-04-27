// Copyright 2025 The Hugo Authors. All rights reserved.
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

package sitematrix

import (
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/common/maps"
)

func TestIntSetsVectorProvider(t *testing.T) {
	c := qt.New(t)

	sets := &IntSets{
		Languages: maps.NewOrderedIntSet(1, 2),
		Versions:  maps.NewOrderedIntSet(1, 2, 3),
		Roles:     maps.NewOrderedIntSet(1, 2, 3),
	}

	c.Assert(sets.HasVector(Vector{1, 2, 3}), qt.Equals, true)
	c.Assert(sets.HasVector(Vector{3, 2, 3}), qt.Equals, false)
	c.Assert(sets.FirstVector(), qt.Equals, Vector{1, 1, 1})

	alllCount := 0
	seen := make(map[Vector]bool)
	for v := range sets.AllVectors() {
		c.Assert(seen[v], qt.IsFalse)
		seen[v] = true
		alllCount++
	}
	// 2 languages * 3 versions * 3 roles = 18 combinations.
	c.Assert(alllCount, qt.Equals, 18)
}
