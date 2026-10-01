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
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestFileMetaInfoUnwrapFileInfo(t *testing.T) {
	c := qt.New(t)
	filename := filepath.Join(t.TempDir(), "a.txt")
	c.Assert(os.WriteFile(filename, []byte("a"), 0o644), qt.IsNil)
	fi, err := os.Stat(filename)
	c.Assert(err, qt.IsNil)

	unwrap := func(fi fs.FileInfo) fs.FileInfo {
		for {
			u, ok := fi.(interface{ UnwrapFileInfo() fs.FileInfo })
			if !ok {
				return fi
			}
			fi = u.UnwrapFileInfo()
		}
	}

	fim := NewFileMetaInfo(fi, &FileMeta{Filename: filename})
	c.Assert(os.SameFile(fi, fim), qt.IsFalse)
	c.Assert(os.SameFile(fi, unwrap(fim)), qt.IsTrue)
	nested := NewFileMetaInfo(fim, &FileMeta{})
	c.Assert(os.SameFile(fi, unwrap(nested)), qt.IsTrue)
}

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
