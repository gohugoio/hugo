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
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/gohugoio/hugo/hugofs/glob"

	"github.com/spf13/afero"

	qt "github.com/frankban/quicktest"
)

func TestFilenameFilterFs(t *testing.T) {
	c := qt.New(t)

	base := filepath.FromSlash("/mybase")

	fs := NewBaseFileDecorator(afero.NewMemMapFs())

	for _, letter := range []string{"a", "b", "c"} {
		for i := 1; i <= 3; i++ {
			c.Assert(afero.WriteFile(fs, filepath.Join(base, letter, fmt.Sprintf("my%d.txt", i)), []byte("some text file for"+letter), 0o755), qt.IsNil)
			c.Assert(afero.WriteFile(fs, filepath.Join(base, letter, fmt.Sprintf("my%d.json", i)), []byte("some json file for"+letter), 0o755), qt.IsNil)
		}
	}

	fs = NewBasePathFs(fs, base)

	filter, err := glob.NewFilenameFilter(nil, []string{"/b/**.txt"})
	c.Assert(err, qt.IsNil)

	fs = newFilenameFilterFs(fs, base, filter)

	assertExists := func(filename string, shouldExist bool) {
		filename = filepath.Clean(filename)
		_, err1 := fs.Stat(filename)
		f, err2 := fs.Open(filename)
		if shouldExist {
			c.Assert(err1, qt.IsNil)
			c.Assert(err2, qt.IsNil)
			defer f.Close()

		} else {
			for _, err := range []error{err1, err2} {
				c.Assert(err, qt.Not(qt.IsNil))
				c.Assert(errors.Is(err, os.ErrNotExist), qt.IsTrue)
			}
		}
	}

	assertExists("/a/my1.txt", true)
	assertExists("/b/my1.txt", false)

	dirB, err := fs.Open("/b")
	c.Assert(err, qt.IsNil)
	defer dirB.Close()
	dirBEntries, err := dirB.Readdirnames(-1)
	c.Assert(err, qt.IsNil)
	c.Assert(dirBEntries, qt.DeepEquals, []string{"my1.json", "my2.json", "my3.json"})

	dirC, err := fs.Open("/c")
	c.Assert(err, qt.IsNil)
	defer dirC.Close()
	dirCEntries, err := dirC.Readdirnames(-1)
	c.Assert(err, qt.IsNil)
	c.Assert(dirCEntries, qt.DeepEquals, []string{"my1.json", "my1.txt", "my2.json", "my2.txt", "my3.json", "my3.txt"})
}
