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

package hugio

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
)

func TestHasBytesWriter(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	c := qt.New((t))

	neww := func() (*HasBytesWriter, io.Writer) {
		var b bytes.Buffer

		h := &HasBytesWriter{
			Patterns: []*HasBytesPattern{
				{Pattern: []byte("__foo")},
			},
		}

		return h, io.MultiWriter(&b, h)
	}

	rndStr := func() string {
		return strings.Repeat("ab cfo", r.Intn(33))
	}

	for range 22 {
		h, w := neww()
		fmt.Fprint(w, rndStr()+"abc __foobar"+rndStr())
		c.Assert(h.Patterns[0].Match, qt.Equals, true)

		h, w = neww()
		fmt.Fprint(w, rndStr()+"abc __f")
		fmt.Fprint(w, "oo bar"+rndStr())
		c.Assert(h.Patterns[0].Match, qt.Equals, true)

		h, w = neww()
		fmt.Fprint(w, rndStr()+"abc __moo bar")
		c.Assert(h.Patterns[0].Match, qt.Equals, false)
	}

	h, w := neww()
	fmt.Fprintf(w, "__foo")
	c.Assert(h.Patterns[0].Match, qt.Equals, true)
}

func TestHasBytesWriterMultiplePatterns(t *testing.T) {
	c := qt.New(t)

	neww := func() (*HasBytesWriter, io.Writer) {
		var b bytes.Buffer
		h := &HasBytesWriter{
			Patterns: []*HasBytesPattern{
				{Pattern: []byte("__hdeferred/")},
				{Pattern: []byte("__h_pp_l1")},
			},
		}
		return h, io.MultiWriter(&b, h)
	}

	// Neither pattern present.
	h, w := neww()
	fmt.Fprint(w, "the quick brown fox jumps over the lazy dog")
	c.Assert(h.Patterns[0].Match, qt.Equals, false)
	c.Assert(h.Patterns[1].Match, qt.Equals, false)
	c.Assert(h.done, qt.Equals, false)

	// Only the second pattern present; the writer must not report a match
	// for the first, and must not prematurely mark itself done.
	h, w = neww()
	fmt.Fprint(w, "prefix __h_pp_l1 suffix")
	c.Assert(h.Patterns[0].Match, qt.Equals, false)
	c.Assert(h.Patterns[1].Match, qt.Equals, true)
	c.Assert(h.done, qt.Equals, false)

	// Both patterns present across multiple writes; done once all match.
	h, w = neww()
	fmt.Fprint(w, "aaa __hdef")
	fmt.Fprint(w, "erred/xyz bbb __h_p")
	fmt.Fprint(w, "p_l1 ccc")
	c.Assert(h.Patterns[0].Match, qt.Equals, true)
	c.Assert(h.Patterns[1].Match, qt.Equals, true)
	c.Assert(h.done, qt.Equals, true)
}

func BenchmarkHasBytesWriter(b *testing.B) {
	// A large chunk of output containing neither pattern is the common case
	// (a normal rendered page): the writer must scan all of it.
	content := []byte(strings.Repeat("<div class=\"nav\"><a href=\"/foo/bar\">baz</a></div>\n", 4000))

	b.ResetTimer()
	for range b.N {
		h := &HasBytesWriter{
			Patterns: []*HasBytesPattern{
				{Pattern: []byte("__hdeferred/")},
				{Pattern: []byte("__h_pp_l1")},
			},
		}
		if _, err := h.Write(content); err != nil {
			b.Fatal(err)
		}
	}
}
