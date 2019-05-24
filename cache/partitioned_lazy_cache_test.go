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

	"github.com/stretchr/testify/require"
)

func TestNewPartitionedLazyCache(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

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
	assert.NoError(err)
	assert.Equal("p1v1", v)

	v, err = cache.Get("p1", "p2_1")
	assert.NoError(err)
	assert.Nil(v)

	v, err = cache.Get("p1", "p1_nil")
	assert.NoError(err)
	assert.Nil(v)

	v, err = cache.Get("p2", "p2_3")
	assert.NoError(err)
	assert.Equal("p2v3", v)

	v, err = cache.Get("doesnotexist", "p1_1")
	assert.NoError(err)
	assert.Nil(v)

	v, err = cache.Get("p1", "doesnotexist")
	assert.NoError(err)
	assert.Nil(v)

	errorP := Partition{
		Key: "p3",
		Load: func() (map[string]interface{}, error) {
			return nil, errors.New("Failed")
		},
	}

	cache = NewPartitionedLazyCache(errorP)

	v, err = cache.Get("p1", "doesnotexist")
	assert.NoError(err)
	assert.Nil(v)

	_, err = cache.Get("p3", "doesnotexist")
	assert.Error(err)

}

func TestConcurrentPartitionedLazyCache(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

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
				assert.NoError(err)
				assert.Equal("p1v1", v)
			}
		}()
	}
	wg.Wait()
}
