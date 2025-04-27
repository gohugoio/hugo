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

package maps

import (
	"fmt"

	"github.com/bits-and-blooms/bitset"
)

type OrderedIntSet struct {
	keys   []int
	values bitset.BitSet
}

// NewOrderedIntSet creates a new OrderedIntSet.
// Note that this is backed by https://github.com/bits-and-blooms/bitset
func NewOrderedIntSet(vals ...int) *OrderedIntSet {
	m := &OrderedIntSet{}
	for _, v := range vals {
		m.Set(v)
	}
	return m
}

// Set sets the value for the given key.
// Note that insertion order is not affected if a key is re-inserted into the set.
func (m *OrderedIntSet) Set(key int) {
	if m == nil {
		return
	}
	keyu := uint(key)
	if m.values.Test(keyu) {
		return
	}
	m.values.Set(keyu)
	m.keys = append(m.keys, key)
}

func (m *OrderedIntSet) Has(key int) bool {
	if m == nil {
		return false
	}
	return m.values.Test(uint(key))
}

func (m *OrderedIntSet) Len() int {
	if m == nil {
		return 0
	}
	return len(m.keys)
}

func (m *OrderedIntSet) String() string {
	if m == nil {
		return "[]"
	}
	return fmt.Sprintf("%v", m.keys)
}
