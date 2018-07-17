// Copyright 2017-present The Hugo Authors. All rights reserved.
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

package cache

import (
	"sync"
)

// Partition represents a cache partition where Load is the callback
// for when the partition is needed.
type Partition struct {
	Key  string
	Load func() (map[string]interface{}, error)
}

// Lazy represents a lazily loaded cache.
type Lazy struct {
	initSync sync.Once
	initErr  error
	cache    map[string]interface{}
	load     func() (map[string]interface{}, error)
}

// NewLazy creates a lazy cache with the given load func.
func NewLazy(load func() (map[string]interface{}, error)) *Lazy {
	return &Lazy{load: load}
}

func (l *Lazy) init() error {
	l.initSync.Do(func() {
		c, err := l.load()
		l.cache = c
		l.initErr = err

	})

	return l.initErr
}

// Get initializes the cache if not already initialized, then looks up the
// given key.
func (l *Lazy) Get(key string) (interface{}, bool, error) {
	l.init()
	if l.initErr != nil {
		return nil, false, l.initErr
	}
	v, found := l.cache[key]
	return v, found, nil
}

// PartitionedLazyCache is a lazily loaded cache paritioned by a supplied string key.
type PartitionedLazyCache struct {
	partitions map[string]*Lazy
}

// NewPartitionedLazyCache creates a new NewPartitionedLazyCache with the supplied
// partitions.
func NewPartitionedLazyCache(partitions ...Partition) *PartitionedLazyCache {
	lazyPartitions := make(map[string]*Lazy, len(partitions))
	for _, partition := range partitions {
		lazyPartitions[partition.Key] = NewLazy(partition.Load)
	}
	cache := &PartitionedLazyCache{partitions: lazyPartitions}

	return cache
}

// Get initializes the partition if not already done so, then looks up the given
// key in the given partition, returns nil if no value found.
func (c *PartitionedLazyCache) Get(partition, key string) (interface{}, error) {
	p, found := c.partitions[partition]

	if !found {
		return nil, nil
	}

	v, found, err := p.Get(key)
	if err != nil {
		return nil, err
	}

	if found {
		return v, nil
	}

	return nil, nil

}
