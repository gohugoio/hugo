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
	"sync"
)

// Cache is a simple thread safe cache backed by a map.
type Cache[K comparable, T any] struct {
	m                  map[K]T
	hasBeenInitialized bool
	sync.RWMutex
}

// NewCache creates a new Cache.
func NewCache[K comparable, T any]() *Cache[K, T] {
	return &Cache[K, T]{m: make(map[K]T)}
}

// Delete deletes the given key from the cache.
// If c is nil, this method is a no-op.
func (c *Cache[K, T]) Get(key K) (T, bool) {
	if c == nil {
		var zero T
		return zero, false
	}
	c.RLock()
	v, found := c.get(key)
	c.RUnlock()
	return v, found
}

func (c *Cache[K, T]) get(key K) (T, bool) {
	v, found := c.m[key]
	return v, found
}

// GetOrCreate gets the value for the given key if it exists, or creates it if not.
func (c *Cache[K, T]) GetOrCreate(key K, create func() (T, error)) (T, error) {
	c.RLock()
	v, found := c.m[key]
	c.RUnlock()
	if found {
		return v, nil
	}
	c.Lock()
	defer c.Unlock()
	v, found = c.m[key]
	if found {
		return v, nil
	}
	v, err := create()
	if err != nil {
		return v, err
	}
	c.m[key] = v
	return v, nil
}

// Contains returns whether the given key exists in the cache.
func (c *Cache[K, T]) Contains(key K) bool {
	c.RLock()
	_, found := c.m[key]
	c.RUnlock()
	return found
}

// InitAndGet initializes the cache if not already done and returns the value for the given key.
// The init state will be reset on Reset or Drain.
func (c *Cache[K, T]) InitAndGet(key K, init func(get func(key K) (T, bool), set func(key K, value T)) error) (T, error) {
	var v T
	c.RLock()
	if !c.hasBeenInitialized {
		c.RUnlock()
		if err := func() error {
			c.Lock()
			defer c.Unlock()
			// Double check in case another goroutine has initialized it in the meantime.
			if !c.hasBeenInitialized {
				err := init(c.get, c.set)
				if err != nil {
					return err
				}
				c.hasBeenInitialized = true
			}
			return nil
		}(); err != nil {
			return v, err
		}
		// Reacquire the read lock.
		c.RLock()
	}

	v = c.m[key]
	c.RUnlock()

	return v, nil
}

// Set sets the given key to the given value.
func (c *Cache[K, T]) Set(key K, value T) {
	c.Lock()
	c.set(key, value)
	c.Unlock()
}

// SetIfAbsent sets the given key to the given value if the key does not already exist in the cache.
func (c *Cache[K, T]) SetIfAbsent(key K, value T) {
	c.RLock()
	if _, found := c.get(key); !found {
		c.RUnlock()
		c.Set(key, value)
	} else {
		c.RUnlock()
	}
}

func (c *Cache[K, T]) set(key K, value T) {
	c.m[key] = value
}

// ForEeach calls the given function for each key/value pair in the cache.
// If the function returns false, the iteration stops.
func (c *Cache[K, T]) ForEeach(f func(K, T) bool) {
	c.RLock()
	defer c.RUnlock()
	for k, v := range c.m {
		if !f(k, v) {
			return
		}
	}
}

func (c *Cache[K, T]) Drain() map[K]T {
	c.Lock()
	m := c.m
	c.m = make(map[K]T)
	c.hasBeenInitialized = false
	c.Unlock()
	return m
}

func (c *Cache[K, T]) Len() int {
	c.RLock()
	defer c.RUnlock()
	return len(c.m)
}

func (c *Cache[K, T]) Reset() {
	c.Lock()
	clear(c.m)
	c.hasBeenInitialized = false
	c.Unlock()
}

// SliceCache is a simple thread safe cache backed by a map.
type SliceCache[T any] struct {
	m map[string][]T
	sync.RWMutex
}

func NewSliceCache[T any]() *SliceCache[T] {
	return &SliceCache[T]{m: make(map[string][]T)}
}

func (c *SliceCache[T]) Get(key string) ([]T, bool) {
	c.RLock()
	v, found := c.m[key]
	c.RUnlock()
	return v, found
}

func (c *SliceCache[T]) Append(key string, values ...T) {
	c.Lock()
	c.m[key] = append(c.m[key], values...)
	c.Unlock()
}

func (c *SliceCache[T]) Reset() {
	c.Lock()
	c.m = make(map[string][]T)
	c.Unlock()
}
