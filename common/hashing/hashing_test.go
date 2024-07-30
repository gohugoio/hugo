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
	"strings"
	"testing"

	"github.com/cespare/xxhash/v2"
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

func xxHashFromString(f string) uint64 {
	h := xxhash.New()
	h.WriteString(f)
	return h.Sum64()
}
