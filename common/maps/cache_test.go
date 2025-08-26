// Copyright 2024 The Hugo Authors. All rights reserved.
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
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestCacheSize(t *testing.T) {
	c := qt.New(t)

	cache := NewCacheWithOptions[string, string](CacheOptions{Size: 10})

	for i := 0; i < 30; i++ {
		cache.Set(string(rune('a'+i)), "value")
	}

	c.Assert(len(cache.m), qt.Equals, 10)

	for i := 20; i < 50; i++ {
		cache.GetOrCreate(string(rune('a'+i)), func() (string, error) {
			return "value", nil
		})
	}

	c.Assert(len(cache.m), qt.Equals, 10)

	for i := 100; i < 200; i++ {
		cache.SetIfAbsent(string(rune('a'+i)), "value")
	}

	c.Assert(len(cache.m), qt.Equals, 10)

	cache.InitAndGet("foo", func(
		get func(key string) (string, bool), set func(key string, value string),
	) error {
		for i := 50; i < 100; i++ {
			set(string(rune('a'+i)), "value")
		}
		return nil
	})

	c.Assert(len(cache.m), qt.Equals, 10)
}
