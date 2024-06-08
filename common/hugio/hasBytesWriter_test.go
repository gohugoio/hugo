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

	for i := 0; i < 22; i++ {
		h, w := neww()
		fmt.Fprintf(w, rndStr()+"abc __foobar"+rndStr())
		c.Assert(h.Patterns[0].Match, qt.Equals, true)

		h, w = neww()
		fmt.Fprintf(w, rndStr()+"abc __f")
		fmt.Fprintf(w, "oo bar"+rndStr())
		c.Assert(h.Patterns[0].Match, qt.Equals, true)

		h, w = neww()
		fmt.Fprintf(w, rndStr()+"abc __moo bar")
		c.Assert(h.Patterns[0].Match, qt.Equals, false)
	}

	h, w := neww()
	fmt.Fprintf(w, "__foo")
	c.Assert(h.Patterns[0].Match, qt.Equals, true)
}
