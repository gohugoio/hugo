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
	"path/filepath"
	"testing"

	"github.com/spf13/afero"

	qt "github.com/frankban/quicktest"
)

func TestGlob(t *testing.T) {
	c := qt.New(t)

	fs := NewBaseFileDecorator(afero.NewMemMapFs())

	create := func(filename string) {
		err := afero.WriteFile(fs, filepath.FromSlash(filename), []byte("content "+filename), 0777)
		c.Assert(err, qt.IsNil)
	}

	collect := func(pattern string) []string {
		var paths []string
		h := func(fi FileMetaInfo) (bool, error) {
			paths = append(paths, fi.Meta().Path())
			return false, nil
		}
		err := Glob(fs, pattern, h)
		c.Assert(err, qt.IsNil)
		return paths
	}

	create("root.json")
	create("jsonfiles/d1.json")
	create("jsonfiles/d2.json")
	create("jsonfiles/sub/d3.json")
	create("jsonfiles/d1.xml")
	create("a/b/c/e/f.json")

	c.Assert(collect("**.json"), qt.HasLen, 5)
	c.Assert(collect("**"), qt.HasLen, 6)
	c.Assert(collect(""), qt.HasLen, 0)
	c.Assert(collect("jsonfiles/*.json"), qt.HasLen, 2)
	c.Assert(collect("*.json"), qt.HasLen, 1)
	c.Assert(collect("**.xml"), qt.HasLen, 1)
	c.Assert(collect(filepath.FromSlash("/jsonfiles/*.json")), qt.HasLen, 2)

}
