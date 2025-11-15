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
	"iter"
	"sync"
)

func NewMap[K comparable, T any]() *Map[K, T] {
	return &Map[K, T]{
		m: make(map[K]T),
	}
}

// Map is a thread safe map backed by a Go map.
type Map[K comparable, T any] struct {
	m  map[K]T
	mu sync.RWMutex
}

// Get gets the value for the given key.
// It returns the zero value of T if the key is not found.
func (m *Map[K, T]) Get(key K) T {
	v, _ := m.Lookup(key)
	return v
}

// Lookup looks up the given key in the map.
// It returns the value and a boolean indicating whether the key was found.
func (m *Map[K, T]) Lookup(key K) (T, bool) {
	m.mu.RLock()
	v, found := m.m[key]
	m.mu.RUnlock()
	return v, found
}

// GetOrCreate gets the value for the given key if it exists, or creates it if not.
func (m *Map[K, T]) GetOrCreate(key K, create func() (T, error)) (T, error) {
	v, found := m.Lookup(key)
	if found {
		return v, nil
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	v, found = m.m[key]
	if found {
		return v, nil
	}
	v, err := create()
	if err != nil {
		return v, err
	}
	m.m[key] = v
	return v, nil
}

// Set sets the given key to the given value.
func (m *Map[K, T]) Set(key K, value T) {
	m.mu.Lock()
	m.m[key] = value
	m.mu.Unlock()
}

// WithWriteLock executes the given function with a write lock on the map.
func (m *Map[K, T]) WithWriteLock(f func(m map[K]T)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	f(m.m)
}

// SetIfAbsent sets the given key to the given value if the key does not already exist in the map.
// It returns true if the value was set, false otherwise.
func (m *Map[K, T]) SetIfAbsent(key K, value T) bool {
	m.mu.RLock()
	if _, found := m.m[key]; !found {
		m.mu.RUnlock()
		return m.doSetIfAbsent(key, value)
	}
	m.mu.RUnlock()
	return false
}

func (m *Map[K, T]) doSetIfAbsent(key K, value T) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, found := m.m[key]; !found {
		m.m[key] = value
		return true
	}
	return false
}

// All returns an iterator over all key/value pairs in the map.
// A read lock is held during the iteration.
func (m *Map[K, T]) All() iter.Seq2[K, T] {
	return func(yield func(K, T) bool) {
		m.mu.RLock()
		defer m.mu.RUnlock()
		for k, v := range m.m {
			if !yield(k, v) {
				return
			}
		}
	}
}
