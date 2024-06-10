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

package langs

import (
	"sync"
	"testing"

	qt "github.com/frankban/quicktest"
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

func TestCollator(t *testing.T) {
	c := qt.New(t)

	var wg sync.WaitGroup

	coll := &Collator{c: collate.New(language.English, collate.Loose)}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			coll.Lock()
			defer coll.Unlock()
			defer wg.Done()
			for j := 0; j < 10; j++ {
				k := coll.CompareStrings("abc", "def")
				c.Assert(k, qt.Equals, -1)
			}
		}()
	}
	wg.Wait()
}

func BenchmarkCollator(b *testing.B) {
	s := []string{"foo", "bar", "éntre", "baz", "qux", "quux", "corge", "grault", "garply", "waldo", "fred", "plugh", "xyzzy", "thud"}

	doWork := func(coll *Collator) {
		for i := 0; i < len(s); i++ {
			for j := i + 1; j < len(s); j++ {
				_ = coll.CompareStrings(s[i], s[j])
			}
		}
	}

	b.Run("Single", func(b *testing.B) {
		coll := &Collator{c: collate.New(language.English, collate.Loose)}
		for i := 0; i < b.N; i++ {
			doWork(coll)
		}
	})

	b.Run("Para", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			coll := &Collator{c: collate.New(language.English, collate.Loose)}

			for pb.Next() {
				coll.Lock()
				doWork(coll)
				coll.Unlock()
			}
		})
	})
}
