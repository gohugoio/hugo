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

// Package namedmemcache provides a memory cache with a named lock. This is suitable
// for situations where creating the cached resource can be time consuming or otherwise
// resource hungry, or in situations where a "once only per key" is a requirement.
package namedmemcache

import (
	"sync"

	"github.com/BurntSushi/locker"
)

// Cache holds the cached values.
type Cache struct {
	nlocker *locker.Locker
	cache   map[string]cacheEntry
	mu      sync.RWMutex
}

type cacheEntry struct {
	value interface{}
	err   error
}

// New creates a new cache.
func New() *Cache {
	return &Cache{
		nlocker: locker.NewLocker(),
		cache:   make(map[string]cacheEntry),
	}
}

// Clear clears the cache state.
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache = make(map[string]cacheEntry)
	c.nlocker = locker.NewLocker()

}

// GetOrCreate tries to get the value with the given cache key, if not found
// create will be called and cached.
// This method is thread safe. It also guarantees that the create func for a given
// key is invoced only once for this cache.
func (c *Cache) GetOrCreate(key string, create func() (interface{}, error)) (interface{}, error) {
	c.mu.RLock()
	entry, found := c.cache[key]
	c.mu.RUnlock()

	if found {
		return entry.value, entry.err
	}

	c.nlocker.Lock(key)
	defer c.nlocker.Unlock(key)

	// Create it.
	value, err := create()

	c.mu.Lock()
	c.cache[key] = cacheEntry{value: value, err: err}
	c.mu.Unlock()

	return value, err
}
