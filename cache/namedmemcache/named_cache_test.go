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

package namedmemcache

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNamedCache(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	cache := New()

	counter := 0
	create := func() (interface{}, error) {
		counter++
		return counter, nil
	}

	for i := 0; i < 5; i++ {
		v1, err := cache.GetOrCreate("a1", create)
		assert.NoError(err)
		assert.Equal(1, v1)
		v2, err := cache.GetOrCreate("a2", create)
		assert.NoError(err)
		assert.Equal(2, v2)
	}

	cache.Clear()

	v3, err := cache.GetOrCreate("a2", create)
	assert.NoError(err)
	assert.Equal(3, v3)
}

func TestNamedCacheConcurrent(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	var wg sync.WaitGroup

	cache := New()

	create := func(i int) func() (interface{}, error) {
		return func() (interface{}, error) {
			return i, nil
		}
	}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				id := fmt.Sprintf("id%d", j)
				v, err := cache.GetOrCreate(id, create(j))
				assert.NoError(err)
				assert.Equal(j, v)
			}
		}()
	}
	wg.Wait()
}
