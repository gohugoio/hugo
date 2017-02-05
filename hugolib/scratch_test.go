// Copyright 2015 The Hugo Authors. All rights reserved.
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
	"reflect"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScratchAdd(t *testing.T) {
	t.Parallel()
	scratch := newScratch()
	scratch.Add("int1", 10)
	scratch.Add("int1", 20)
	scratch.Add("int2", 20)

	assert.Equal(t, int64(30), scratch.Get("int1"))
	assert.Equal(t, 20, scratch.Get("int2"))

	scratch.Add("float1", float64(10.5))
	scratch.Add("float1", float64(20.1))

	assert.Equal(t, float64(30.6), scratch.Get("float1"))

	scratch.Add("string1", "Hello ")
	scratch.Add("string1", "big ")
	scratch.Add("string1", "World!")

	assert.Equal(t, "Hello big World!", scratch.Get("string1"))

	scratch.Add("scratch", scratch)
	_, err := scratch.Add("scratch", scratch)

	if err == nil {
		t.Errorf("Expected error from invalid arithmetic")
	}

}

func TestScratchAddSlice(t *testing.T) {
	t.Parallel()
	scratch := newScratch()

	_, err := scratch.Add("intSlice", []int{1, 2})
	assert.Nil(t, err)
	_, err = scratch.Add("intSlice", 3)
	assert.Nil(t, err)

	sl := scratch.Get("intSlice")
	expected := []int{1, 2, 3}

	if !reflect.DeepEqual(expected, sl) {
		t.Errorf("Slice difference, go %q expected %q", sl, expected)
	}

	_, err = scratch.Add("intSlice", []int{4, 5})

	assert.Nil(t, err)

	sl = scratch.Get("intSlice")
	expected = []int{1, 2, 3, 4, 5}

	if !reflect.DeepEqual(expected, sl) {
		t.Errorf("Slice difference, go %q expected %q", sl, expected)
	}

}

func TestScratchSet(t *testing.T) {
	t.Parallel()
	scratch := newScratch()
	scratch.Set("key", "val")
	assert.Equal(t, "val", scratch.Get("key"))
}

// Issue #2005
func TestScratchInParallel(t *testing.T) {
	var wg sync.WaitGroup
	scratch := newScratch()
	key := "counter"
	scratch.Set(key, int64(1))
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(j int) {
			for k := 0; k < 10; k++ {
				newVal := int64(k + j)

				_, err := scratch.Add(key, newVal)
				if err != nil {
					t.Errorf("Got err %s", err)
				}

				scratch.Set(key, newVal)

				val := scratch.Get(key)

				if counter, ok := val.(int64); ok {
					if counter < 1 {
						t.Errorf("Got %d", counter)
					}
				} else {
					t.Errorf("Got %T", val)
				}
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func TestScratchGet(t *testing.T) {
	t.Parallel()
	scratch := newScratch()
	nothing := scratch.Get("nothing")
	if nothing != nil {
		t.Errorf("Should not return anything, but got %v", nothing)
	}
}

func TestScratchSetInMap(t *testing.T) {
	t.Parallel()
	scratch := newScratch()
	scratch.SetInMap("key", "lux", "Lux")
	scratch.SetInMap("key", "abc", "Abc")
	scratch.SetInMap("key", "zyx", "Zyx")
	scratch.SetInMap("key", "abc", "Abc (updated)")
	scratch.SetInMap("key", "def", "Def")
	assert.Equal(t, []interface{}{0: "Abc (updated)", 1: "Def", 2: "Lux", 3: "Zyx"}, scratch.GetSortedMapValues("key"))
}

func TestScratchGetSortedMapValues(t *testing.T) {
	t.Parallel()
	scratch := newScratch()
	nothing := scratch.GetSortedMapValues("nothing")
	if nothing != nil {
		t.Errorf("Should not return anything, but got %v", nothing)
	}
}

func BenchmarkScratchGet(b *testing.B) {
	scratch := newScratch()
	scratch.Add("A", 1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scratch.Get("A")
	}
}
