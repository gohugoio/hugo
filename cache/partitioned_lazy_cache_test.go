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
	"errors"
	"sync"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestNewPartitionedLazyCache(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	p1 := Partition{
		Key: "p1",
		Load: func() (map[string]interface{}, error) {
			return map[string]interface{}{
				"p1_1":   "p1v1",
				"p1_2":   "p1v2",
				"p1_nil": nil,
			}, nil
		},
	}

	p2 := Partition{
		Key: "p2",
		Load: func() (map[string]interface{}, error) {
			return map[string]interface{}{
				"p2_1": "p2v1",
				"p2_2": "p2v2",
				"p2_3": "p2v3",
			}, nil
		},
	}

	cache := NewPartitionedLazyCache(p1, p2)

	v, err := cache.Get("p1", "p1_1")
	c.Assert(err, qt.IsNil)
	c.Assert(v, qt.Equals, "p1v1")

	v, err = cache.Get("p1", "p2_1")
	c.Assert(err, qt.IsNil)
	c.Assert(v, qt.IsNil)

	v, err = cache.Get("p1", "p1_nil")
	c.Assert(err, qt.IsNil)
	c.Assert(v, qt.IsNil)

	v, err = cache.Get("p2", "p2_3")
	c.Assert(err, qt.IsNil)
	c.Assert(v, qt.Equals, "p2v3")

	v, err = cache.Get("doesnotexist", "p1_1")
	c.Assert(err, qt.IsNil)
	c.Assert(v, qt.IsNil)

	v, err = cache.Get("p1", "doesnotexist")
	c.Assert(err, qt.IsNil)
	c.Assert(v, qt.IsNil)

	errorP := Partition{
		Key: "p3",
		Load: func() (map[string]interface{}, error) {
			return nil, errors.New("Failed")
		},
	}

	cache = NewPartitionedLazyCache(errorP)

	v, err = cache.Get("p1", "doesnotexist")
	c.Assert(err, qt.IsNil)
	c.Assert(v, qt.IsNil)

	_, err = cache.Get("p3", "doesnotexist")
	c.Assert(err, qt.Not(qt.IsNil))

}

func TestConcurrentPartitionedLazyCache(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	var wg sync.WaitGroup

	p1 := Partition{
		Key: "p1",
		Load: func() (map[string]interface{}, error) {
			return map[string]interface{}{
				"p1_1":   "p1v1",
				"p1_2":   "p1v2",
				"p1_nil": nil,
			}, nil
		},
	}

	p2 := Partition{
		Key: "p2",
		Load: func() (map[string]interface{}, error) {
			return map[string]interface{}{
				"p2_1": "p2v1",
				"p2_2": "p2v2",
				"p2_3": "p2v3",
			}, nil
		},
	}

	cache := NewPartitionedLazyCache(p1, p2)

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				v, err := cache.Get("p1", "p1_1")
				c.Assert(err, qt.IsNil)
				c.Assert(v, qt.Equals, "p1v1")
			}
		}()
	}
	wg.Wait()
}
