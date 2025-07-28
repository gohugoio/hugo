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
	"slices"

	"github.com/bits-and-blooms/bitset"
)

type OrderedIntSet struct {
	keys   []int
	values *bitset.BitSet
}

// NewOrderedIntSet creates a new OrderedIntSet.
// Note that this is backed by https://github.com/bits-and-blooms/bitset
func NewOrderedIntSet(vals ...int) *OrderedIntSet {
	m := &OrderedIntSet{
		keys:   make([]int, 0, len(vals)),
		values: bitset.New(uint(len(vals))),
	}
	for _, v := range vals {
		m.Set(v)
	}
	return m
}

// Set sets the value for the given key.
// Note that insertion order is not affected if a key is re-inserted into the set.
func (m *OrderedIntSet) Set(key int) {
	if m == nil {
		panic("nil OrderedIntSet")
	}
	keyu := uint(key)
	if m.values.Test(keyu) {
		return
	}
	m.values.Set(keyu)
	m.keys = append(m.keys, key)
}

// SetFrom sets the values from another OrderedIntSet.
func (m *OrderedIntSet) SetFrom(other *OrderedIntSet) {
	if m == nil || other == nil {
		return
	}
	for _, key := range other.keys {
		m.Set(key)
	}
}

func (m *OrderedIntSet) Clone() *OrderedIntSet {
	if m == nil {
		return nil
	}
	newSet := &OrderedIntSet{
		keys:   slices.Clone(m.keys),
		values: m.values.Clone(),
	}
	return newSet
}

// Next returns the next key in the set possibly including the given key.
// It returns -1 if the key is not found or if there are no keys greater than the given key.
func (m *OrderedIntSet) Next(i int) int {
	n, ok := m.values.NextSet(uint(i))
	if !ok {
		return -1
	}
	return int(n)
}

// The reason we don't use iter.Seq is https://github.com/golang/go/issues/69015
// This is 70% faster than using iter.Seq2[int, int] for the keys.
// It returns false if the iteration was stopped early.
func (m *OrderedIntSet) ForEachKey(yield func(int) bool) bool {
	if m == nil {
		return true
	}
	for _, key := range m.keys {
		if !yield(key) {
			return false
		}
	}
	return true
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

// KeysSorted returns the keys in sorted order.
func (m *OrderedIntSet) KeysSorted() []int {
	if m == nil {
		return nil
	}
	keys := slices.Clone(m.keys)
	slices.Sort(keys)
	return m.keys
}

func (m *OrderedIntSet) String() string {
	if m == nil {
		return "[]"
	}
	return fmt.Sprintf("%v", m.keys)
}

func (m *OrderedIntSet) Values() *bitset.BitSet {
	if m == nil {
		return nil
	}
	return m.values
}

func (m *OrderedIntSet) IsSuperSet(other *OrderedIntSet) bool {
	if m == nil || other == nil {
		return false
	}
	return m.values.IsSuperSet(other.values)
}

// Words returns the bitset as array of 64-bit words, giving direct access to the internal representation.
// It is not a copy, so changes to the returned slice will affect the bitset.
// It is meant for advanced users.
func (m *OrderedIntSet) Words() []uint64 {
	if m == nil {
		return nil
	}
	return m.values.Words()
}
