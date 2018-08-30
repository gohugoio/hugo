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

package hugolib

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOrderedMap(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	m := newOrderedMap()
	m.Add("b", "vb")
	m.Add("c", "vc")
	m.Add("a", "va")
	b, f1 := m.Get("b")

	assert.True(f1)
	assert.Equal(b, "vb")
	assert.True(m.Contains("b"))
	assert.False(m.Contains("e"))

	assert.Equal([]interface{}{"b", "c", "a"}, m.Keys())

}

func TestOrderedMapConcurrent(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	var wg sync.WaitGroup

	m := newOrderedMap()

	for i := 1; i < 20; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", id)
			val := key + "val"
			m.Add(key, val)
			v, found := m.Get(key)
			assert.True(found)
			assert.Equal(v, val)
			assert.True(m.Contains(key))
			assert.True(m.Len() > 0)
			assert.True(len(m.Keys()) > 0)
		}(i)

	}

	wg.Wait()
}
