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

package hashing

import (
	"fmt"
	"math"
	"strings"
	"sync"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestXxHashFromReader(t *testing.T) {
	c := qt.New(t)
	s := "Hello World"
	r := strings.NewReader(s)
	got, size, err := XXHashFromReader(r)
	c.Assert(err, qt.IsNil)
	c.Assert(size, qt.Equals, int64(len(s)))
	c.Assert(got, qt.Equals, uint64(7148569436472236994))
}

func TestXxHashFromReaderPara(t *testing.T) {
	c := qt.New(t)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				s := strings.Repeat("Hello ", i+j+1*42)
				r := strings.NewReader(s)
				got, size, err := XXHashFromReader(r)
				c.Assert(size, qt.Equals, int64(len(s)))
				c.Assert(err, qt.IsNil)
				expect, _ := XXHashFromString(s)
				c.Assert(got, qt.Equals, expect)
			}
		}()
	}

	wg.Wait()
}

func TestXxHashFromString(t *testing.T) {
	c := qt.New(t)
	s := "Hello World"
	got, err := XXHashFromString(s)
	c.Assert(err, qt.IsNil)
	c.Assert(got, qt.Equals, uint64(7148569436472236994))
}

func TestXxHashFromStringHexEncoded(t *testing.T) {
	c := qt.New(t)
	s := "The quick brown fox jumps over the lazy dog"
	got := XxHashFromStringHexEncoded(s)
	// Facit: https://asecuritysite.com/encryption/xxhash?val=The%20quick%20brown%20fox%20jumps%20over%20the%20lazy%20dog
	c.Assert(got, qt.Equals, "0b242d361fda71bc")
}

func BenchmarkXXHashFromReader(b *testing.B) {
	r := strings.NewReader("Hello World")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		XXHashFromReader(r)
		r.Seek(0, 0)
	}
}

func BenchmarkXXHashFromString(b *testing.B) {
	s := "Hello World"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		XXHashFromString(s)
	}
}

func BenchmarkXXHashFromStringHexEncoded(b *testing.B) {
	s := "The quick brown fox jumps over the lazy dog"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		XxHashFromStringHexEncoded(s)
	}
}

func TestHashString(t *testing.T) {
	c := qt.New(t)

	c.Assert(HashString("a", "b"), qt.Equals, "3176555414984061461")
	c.Assert(HashString("ab"), qt.Equals, "7347350983217793633")

	var vals []any = []any{"a", "b", tstKeyer{"c"}}

	c.Assert(HashString(vals...), qt.Equals, "4438730547989914315")
	c.Assert(vals[2], qt.Equals, tstKeyer{"c"})
}

type tstKeyer struct {
	key string
}

func (t tstKeyer) Key() string {
	return t.key
}

func (t tstKeyer) String() string {
	return "key: " + t.key
}

func BenchmarkHashString(b *testing.B) {
	word := " hello "

	var tests []string

	for i := 1; i <= 5; i++ {
		sentence := strings.Repeat(word, int(math.Pow(4, float64(i))))
		tests = append(tests, sentence)
	}

	b.ResetTimer()

	for _, test := range tests {
		b.Run(fmt.Sprintf("n%d", len(test)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				HashString(test)
			}
		})
	}
}
