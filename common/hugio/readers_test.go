// Copyright 2026 The Hugo Authors. All rights reserved.
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
	"io"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestNewSkipBOMReader(t *testing.T) {
	c := qt.New(t)

	newReader := func(s string) ReadSeekCloser {
		r, err := NewSkipBOMReader(NewReadSeekerNoOpCloserFromString(s))
		c.Assert(err, qt.IsNil)
		return r
	}

	readAll := func(r io.Reader) string {
		b, err := io.ReadAll(r)
		c.Assert(err, qt.IsNil)
		return string(b)
	}

	r := newReader(BOM + "abc")
	c.Assert(readAll(r), qt.Equals, "abc")

	pos, err := r.Seek(0, io.SeekStart)
	c.Assert(err, qt.IsNil)
	c.Assert(pos, qt.Equals, int64(0))
	c.Assert(readAll(r), qt.Equals, "abc")

	pos, err = r.Seek(1, io.SeekStart)
	c.Assert(err, qt.IsNil)
	c.Assert(pos, qt.Equals, int64(1))
	c.Assert(readAll(r), qt.Equals, "bc")

	pos, err = r.Seek(-1, io.SeekEnd)
	c.Assert(err, qt.IsNil)
	c.Assert(pos, qt.Equals, int64(2))
	c.Assert(readAll(r), qt.Equals, "c")

	pos, err = r.Seek(0, io.SeekCurrent)
	c.Assert(err, qt.IsNil)
	c.Assert(pos, qt.Equals, int64(3))

	_, err = r.Seek(-1, io.SeekStart)
	c.Assert(err, qt.ErrorMatches, ".*negative position")
	_, err = r.Seek(-4, io.SeekEnd)
	c.Assert(err, qt.ErrorMatches, ".*negative position")
	_, err = r.Seek(-4, io.SeekCurrent)
	c.Assert(err, qt.ErrorMatches, ".*negative position")

	pos, err = r.Seek(-3, io.SeekEnd)
	c.Assert(err, qt.IsNil)
	c.Assert(pos, qt.Equals, int64(0))
	c.Assert(readAll(r), qt.Equals, "abc")

	c.Assert(r.Close(), qt.IsNil)

	for _, s := range []string{"abc", "ab", "a", "", BOM[:2], BOM} {
		got := readAll(newReader(s))
		if s == BOM {
			c.Assert(got, qt.Equals, "")
		} else {
			c.Assert(got, qt.Equals, s)
		}
	}
}
