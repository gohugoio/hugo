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

package maps

import (
	"reflect"
	"sync"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestScratchAdd(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	scratch := NewScratch()
	scratch.Add("int1", 10)
	scratch.Add("int1", 20)
	scratch.Add("int2", 20)

	c.Assert(scratch.Get("int1"), qt.Equals, int64(30))
	c.Assert(scratch.Get("int2"), qt.Equals, 20)

	scratch.Add("float1", float64(10.5))
	scratch.Add("float1", float64(20.1))

	c.Assert(scratch.Get("float1"), qt.Equals, float64(30.6))

	scratch.Add("string1", "Hello ")
	scratch.Add("string1", "big ")
	scratch.Add("string1", "World!")

	c.Assert(scratch.Get("string1"), qt.Equals, "Hello big World!")

	scratch.Add("scratch", scratch)
	_, err := scratch.Add("scratch", scratch)

	if err == nil {
		t.Errorf("Expected error from invalid arithmetic")
	}

}

func TestScratchAddSlice(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	scratch := NewScratch()

	_, err := scratch.Add("intSlice", []int{1, 2})
	c.Assert(err, qt.IsNil)
	_, err = scratch.Add("intSlice", 3)
	c.Assert(err, qt.IsNil)

	sl := scratch.Get("intSlice")
	expected := []int{1, 2, 3}

	if !reflect.DeepEqual(expected, sl) {
		t.Errorf("Slice difference, go %q expected %q", sl, expected)
	}
	_, err = scratch.Add("intSlice", []int{4, 5})

	c.Assert(err, qt.IsNil)

	sl = scratch.Get("intSlice")
	expected = []int{1, 2, 3, 4, 5}

	if !reflect.DeepEqual(expected, sl) {
		t.Errorf("Slice difference, go %q expected %q", sl, expected)
	}
}

// https://github.com/gohugoio/hugo/issues/5275
func TestScratchAddTypedSliceToInterfaceSlice(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	scratch := NewScratch()
	scratch.Set("slice", []interface{}{})

	_, err := scratch.Add("slice", []int{1, 2})
	c.Assert(err, qt.IsNil)
	c.Assert(scratch.Get("slice"), qt.DeepEquals, []int{1, 2})

}

// https://github.com/gohugoio/hugo/issues/5361
func TestScratchAddDifferentTypedSliceToInterfaceSlice(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	scratch := NewScratch()
	scratch.Set("slice", []string{"foo"})

	_, err := scratch.Add("slice", []int{1, 2})
	c.Assert(err, qt.IsNil)
	c.Assert(scratch.Get("slice"), qt.DeepEquals, []interface{}{"foo", 1, 2})

}

func TestScratchSet(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	scratch := NewScratch()
	scratch.Set("key", "val")
	c.Assert(scratch.Get("key"), qt.Equals, "val")
}

func TestScratchDelete(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	scratch := NewScratch()
	scratch.Set("key", "val")
	scratch.Delete("key")
	scratch.Add("key", "Lucy Parsons")
	c.Assert(scratch.Get("key"), qt.Equals, "Lucy Parsons")
}

// Issue #2005
func TestScratchInParallel(t *testing.T) {
	var wg sync.WaitGroup
	scratch := NewScratch()

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
	scratch := NewScratch()
	nothing := scratch.Get("nothing")
	if nothing != nil {
		t.Errorf("Should not return anything, but got %v", nothing)
	}
}

func TestScratchSetInMap(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	scratch := NewScratch()
	scratch.SetInMap("key", "lux", "Lux")
	scratch.SetInMap("key", "abc", "Abc")
	scratch.SetInMap("key", "zyx", "Zyx")
	scratch.SetInMap("key", "abc", "Abc (updated)")
	scratch.SetInMap("key", "def", "Def")
	c.Assert(scratch.GetSortedMapValues("key"), qt.DeepEquals, []interface{}{0: "Abc (updated)", 1: "Def", 2: "Lux", 3: "Zyx"})
}

func TestScratchGetSortedMapValues(t *testing.T) {
	t.Parallel()
	scratch := NewScratch()
	nothing := scratch.GetSortedMapValues("nothing")
	if nothing != nil {
		t.Errorf("Should not return anything, but got %v", nothing)
	}
}

func BenchmarkScratchGet(b *testing.B) {
	scratch := NewScratch()
	scratch.Add("A", 1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scratch.Get("A")
	}
}
