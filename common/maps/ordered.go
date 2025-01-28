// Copyright 2024 The Hugo Authors. All rights reserved.
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
	"github.com/gohugoio/hugo/common/hashing"
)

// Ordered is a map that can be iterated in the order of insertion.
// Note that insertion order is not affected if a key is re-inserted into the map.
// In a nil map, all operations are no-ops.
// This is not thread safe.
type Ordered[K comparable, T any] struct {
	// The keys in the order they were added.
	keys []K
	// The values.
	values map[K]T
}

// NewOrdered creates a new Ordered map.
func NewOrdered[K comparable, T any]() *Ordered[K, T] {
	return &Ordered[K, T]{values: make(map[K]T)}
}

// Contains returns whether the map contains the given key.
func (m *Ordered[K, T]) Contains(key K) bool {
	if m == nil {
		return false
	}
	_, found := m.values[key]
	return found
}

// Set sets the value for the given key.
// Note that insertion order is not affected if a key is re-inserted into the map.
func (m *Ordered[K, T]) Set(key K, value T) {
	if m == nil {
		return
	}
	// Check if key already exists.
	if _, found := m.values[key]; !found {
		m.keys = append(m.keys, key)
	}
	m.values[key] = value
}

// Get gets the value for the given key.
func (m *Ordered[K, T]) Get(key K) (T, bool) {
	if m == nil {
		var v T
		return v, false
	}
	value, found := m.values[key]
	return value, found
}

// Delete deletes the value for the given key.
func (m *Ordered[K, T]) Delete(key K) {
	if m == nil {
		return
	}
	delete(m.values, key)
	for i, k := range m.keys {
		if k == key {
			m.keys = append(m.keys[:i], m.keys[i+1:]...)
			break
		}
	}
}

// Clone creates a shallow copy of the map.
func (m *Ordered[K, T]) Clone() *Ordered[K, T] {
	if m == nil {
		return nil
	}
	clone := NewOrdered[K, T]()
	for _, k := range m.keys {
		clone.Set(k, m.values[k])
	}
	return clone
}

// Keys returns the keys in the order they were added.
func (m *Ordered[K, T]) Keys() []K {
	if m == nil {
		return nil
	}
	return m.keys
}

// Values returns the values in the order they were added.
func (m *Ordered[K, T]) Values() []T {
	if m == nil {
		return nil
	}
	var values []T
	for _, k := range m.keys {
		values = append(values, m.values[k])
	}
	return values
}

// Len returns the number of items in the map.
func (m *Ordered[K, T]) Len() int {
	if m == nil {
		return 0
	}
	return len(m.keys)
}

// Range calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
// TODO(bep) replace with iter.Seq2 when we bump go Go 1.24.
func (m *Ordered[K, T]) Range(f func(key K, value T) bool) {
	if m == nil {
		return
	}
	for _, k := range m.keys {
		if !f(k, m.values[k]) {
			return
		}
	}
}

// Hash calculates a hash from the values.
func (m *Ordered[K, T]) Hash() (uint64, error) {
	if m == nil {
		return 0, nil
	}
	return hashing.Hash(m.values)
}
