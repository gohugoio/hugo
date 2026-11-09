// Copyright 2019 The Hugo Authors. All rights reserved.
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
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"

	qt "github.com/frankban/quicktest"
)

func TestGlob(t *testing.T) {
	c := qt.New(t)

	var visited []string
	fs := NewBaseFileDecorator(afero.NewMemMapFs(), func(fi FileMetaInfo) {
		visited = append(visited, filepath.ToSlash(fi.Meta().Filename))
	})

	create := func(filename string) {
		filename = filepath.FromSlash(filename)
		dir := filepath.Dir(filename)
		if dir != "." {
			err := fs.MkdirAll(dir, 0o777)
			c.Assert(err, qt.IsNil)
		}
		err := afero.WriteFile(fs, filename, []byte("content "+filename), 0o777)
		c.Assert(err, qt.IsNil)
	}

	collect := func(pattern string) []string {
		var paths []string
		h := func(fi FileMetaInfo) (bool, error) {
			p := filepath.ToSlash(fi.Meta().Filename)
			paths = append(paths, p)
			return false, nil
		}
		err := Glob(fs, pattern, h)
		c.Assert(err, qt.IsNil)
		return paths
	}

	create("/root.json")
	create("/jsonfiles/d1.json")
	create("/jsonfiles/d2.json")
	create("/jsonfiles/sub/d3.json")
	create("/jsonfiles/d1.xml")
	create("/a/b/c/e/f.json")
	create("/UPPER/sub/style.css")
	create("/root/UPPER/sub/style.css")
	create("/İmages/test.txt")

	afero.Walk(fs, "/", func(path string, info os.FileInfo, err error) error {
		c.Assert(err, qt.IsNil)
		return nil
	})

	c.Assert(collect(filepath.FromSlash("/jsonfiles/*.json")), qt.HasLen, 2)
	c.Assert(collect("/*.json"), qt.HasLen, 1)
	c.Assert(collect("**.json"), qt.HasLen, 5)
	c.Assert(collect("**"), qt.HasLen, 9)
	c.Assert(collect(""), qt.HasLen, 0)
	c.Assert(collect("jsonfiles/*.json"), qt.HasLen, 2)
	c.Assert(collect("*.json"), qt.HasLen, 1)
	c.Assert(collect("**.xml"), qt.HasLen, 1)

	c.Assert(collect("root/UPPER/sub/style.css"), qt.HasLen, 1)
	c.Assert(collect("UPPER/sub/style.css"), qt.HasLen, 1)
	c.Assert(collect("upper/sub/STYLE.css"), qt.DeepEquals, []string{"/UPPER/sub/style.css"})
	visited = nil
	c.Assert(collect("ROOT/upper/sub/*.CSS"), qt.DeepEquals, []string{"/root/UPPER/sub/style.css"})
	c.Assert(visited, qt.Not(qt.Contains), "/UPPER/sub/style.css")
	// strings.ToLower("İ") is "i", but strings.EqualFold("İ", "i") is false.
	c.Assert(collect("images/*.txt"), qt.DeepEquals, []string{"/İmages/test.txt"})
}
