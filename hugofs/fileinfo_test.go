// Copyright 2021 The Hugo Authors. All rights reserved.
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

package hugofs

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestFileMeta(t *testing.T) {
	c := qt.New(t)

	c.Run("Merge", func(c *qt.C) {
		src := &FileMeta{
			Filename: "fs1",
		}
		dst := &FileMeta{
			Filename: "fd1",
		}

		dst.Merge(src)

		c.Assert(dst.Filename, qt.Equals, "fd1")
	})

	c.Run("Copy", func(c *qt.C) {
		src := &FileMeta{
			Filename: "fs1",
		}
		dst := src.Copy()

		c.Assert(dst, qt.Not(qt.Equals), src)
		c.Assert(dst, qt.DeepEquals, src)
	})
}
