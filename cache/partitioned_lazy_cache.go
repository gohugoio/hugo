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

type lazyPartition struct {
	initSync sync.Once
	cache    map[string]interface{}
	load     func() (map[string]interface{}, error)
}

func (l *lazyPartition) init() error {
	var err error
	l.initSync.Do(func() {
		var c map[string]interface{}
		c, err = l.load()
		l.cache = c
	})

	return err
}

// PartitionedLazyCache is a lazily loaded cache paritioned by a supplied string key.
type PartitionedLazyCache struct {
	partitions map[string]*lazyPartition
}

// NewPartitionedLazyCache creates a new NewPartitionedLazyCache with the supplied
// partitions.
func NewPartitionedLazyCache(partitions ...Partition) *PartitionedLazyCache {
	lazyPartitions := make(map[string]*lazyPartition, len(partitions))
	for _, partition := range partitions {
		lazyPartitions[partition.Key] = &lazyPartition{load: partition.Load}
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

	if err := p.init(); err != nil {
		return nil, err
	}

	if v, found := p.cache[key]; found {
		return v, nil
	}

	return nil, nil

}
