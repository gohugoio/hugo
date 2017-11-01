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

package types

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEvictingStringQueue(t *testing.T) {
	assert := require.New(t)

	queue := NewEvictingStringQueue(3)

	assert.Equal("", queue.Peek())
	queue.Add("a")
	queue.Add("b")
	queue.Add("a")
	assert.Equal("b", queue.Peek())
	queue.Add("b")
	assert.Equal("b", queue.Peek())

	queue.Add("a")
	queue.Add("b")

	assert.Equal([]string{"b", "a"}, queue.PeekAll())
	assert.Equal("b", queue.Peek())
	queue.Add("c")
	queue.Add("d")
	// Overflowed, a should now be removed.
	assert.Equal([]string{"d", "c", "b"}, queue.PeekAll())
	assert.Len(queue.PeekAllSet(), 3)
	assert.True(queue.PeekAllSet()["c"])
}

func TestEvictingStringQueueConcurrent(t *testing.T) {
	var wg sync.WaitGroup
	val := "someval"

	queue := NewEvictingStringQueue(3)

	for j := 0; j < 100; j++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			queue.Add(val)
			v := queue.Peek()
			if v != val {
				t.Error("wrong val")
			}
			vals := queue.PeekAll()
			if len(vals) != 1 || vals[0] != val {
				t.Error("wrong val")
			}
		}()
	}
	wg.Wait()
}
