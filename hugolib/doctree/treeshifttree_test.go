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

package doctree_test

import (
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/hugolib/doctree"
	"github.com/gohugoio/hugo/hugolib/sitesmatrix"
)

func TestTreeShiftTree(t *testing.T) {
	c := qt.New(t)

	tree := doctree.NewTreeShiftTree[string](sitesmatrix.Vector{10, 1, 1})
	c.Assert(tree, qt.IsNotNil)
}

func BenchmarkTreeShiftTreeSlice(b *testing.B) {
	v := sitesmatrix.Vector{10, 10, 10}
	t := doctree.NewTreeShiftTree[string](v)
	b.Run("New", func(b *testing.B) {
		for b.Loop() {
			for l1 := 0; l1 < v[0]; l1++ {
				for l2 := 0; l2 < v[1]; l2++ {
					for l3 := 0; l3 < v[2]; l3++ {
						_ = doctree.NewTreeShiftTree[string](sitesmatrix.Vector{l1, l2, l3})
					}
				}
			}
		}
	})
	b.Run("Shape", func(b *testing.B) {
		for b.Loop() {
			for l1 := 0; l1 < v[0]; l1++ {
				for l2 := 0; l2 < v[1]; l2++ {
					for l3 := 0; l3 < v[2]; l3++ {
						t.Shape(sitesmatrix.Vector{l1, l2, l3})
					}
				}
			}
		}
	})
}
