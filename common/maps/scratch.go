// Copyright 2019 The Hugo Authors. All rights reserved.
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
	"reflect"
	"sort"
	"sync"

	"github.com/gohugoio/hugo/common/collections"
	"github.com/gohugoio/hugo/common/math"
)

// Scratch is a writable context used for stateful operations in Page/Node rendering.
type Scratch struct {
	values map[string]any
	mu     sync.RWMutex
}

// Scratcher provides a scratching service.
type Scratcher interface {
	// Scratch returns a "scratch pad" that can be used to store state.
	Scratch() *Scratch
}

type scratcher struct {
	s *Scratch
}

func (s scratcher) Scratch() *Scratch {
	return s.s
}

// NewScratcher creates a new Scratcher.
func NewScratcher() Scratcher {
	return scratcher{s: NewScratch()}
}

// Add will, for single values, add (using the + operator) the addend to the existing addend (if found).
// Supports numeric values and strings.
//
// If the first add for a key is an array or slice, then the next value(s) will be appended.
func (c *Scratch) Add(key string, newAddend any) (string, error) {
	var newVal any
	c.mu.RLock()
	existingAddend, found := c.values[key]
	c.mu.RUnlock()
	if found {
		var err error

		addendV := reflect.TypeOf(existingAddend)

		if addendV.Kind() == reflect.Slice || addendV.Kind() == reflect.Array {
			newVal, err = collections.Append(existingAddend, newAddend)
			if err != nil {
				return "", err
			}
		} else {
			newVal, err = math.DoArithmetic(existingAddend, newAddend, '+')
			if err != nil {
				return "", err
			}
		}
	} else {
		newVal = newAddend
	}
	c.mu.Lock()
	c.values[key] = newVal
	c.mu.Unlock()
	return "", nil // have to return something to make it work with the Go templates
}

// Set stores a value with the given key in the Node context.
// This value can later be retrieved with Get.
func (c *Scratch) Set(key string, value any) string {
	c.mu.Lock()
	c.values[key] = value
	c.mu.Unlock()
	return ""
}

// Delete deletes the given key.
func (c *Scratch) Delete(key string) string {
	c.mu.Lock()
	delete(c.values, key)
	c.mu.Unlock()
	return ""
}

// Get returns a value previously set by Add or Set.
func (c *Scratch) Get(key string) any {
	c.mu.RLock()
	val := c.values[key]
	c.mu.RUnlock()

	return val
}

// Values returns the raw backing map. Note that you should just use
// this method on the locally scoped Scratch instances you obtain via newScratch, not
// .Page.Scratch etc., as that will lead to concurrency issues.
func (c *Scratch) Values() map[string]any {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.values
}

// SetInMap stores a value to a map with the given key in the Node context.
// This map can later be retrieved with GetSortedMapValues.
func (c *Scratch) SetInMap(key string, mapKey string, value any) string {
	c.mu.Lock()
	_, found := c.values[key]
	if !found {
		c.values[key] = make(map[string]any)
	}

	c.values[key].(map[string]any)[mapKey] = value
	c.mu.Unlock()
	return ""
}

// DeleteInMap deletes a value to a map with the given key in the Node context.
func (c *Scratch) DeleteInMap(key string, mapKey string) string {
	c.mu.Lock()
	_, found := c.values[key]
	if found {
		delete(c.values[key].(map[string]any), mapKey)
	}
	c.mu.Unlock()
	return ""
}

// GetSortedMapValues returns a sorted map previously filled with SetInMap.
func (c *Scratch) GetSortedMapValues(key string) any {
	c.mu.RLock()

	if c.values[key] == nil {
		c.mu.RUnlock()
		return nil
	}

	unsortedMap := c.values[key].(map[string]any)
	c.mu.RUnlock()
	var keys []string
	for mapKey := range unsortedMap {
		keys = append(keys, mapKey)
	}

	sort.Strings(keys)

	sortedArray := make([]any, len(unsortedMap))
	for i, mapKey := range keys {
		sortedArray[i] = unsortedMap[mapKey]
	}

	return sortedArray
}

// NewScratch returns a new instance of Scratch.
func NewScratch() *Scratch {
	return &Scratch{values: make(map[string]any)}
}
