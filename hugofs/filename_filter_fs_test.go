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
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/spf13/afero"
)

// ref. https://github.com/spf13/afero/blob/master/ro_regexp_test.go

func TestFilenameFilter(t *testing.T) {
	c := qt.New(t)

	pred := func(name string) bool {
		return name == "/only.html"
	}

	fs := NewFilenameFilterFs(&afero.MemMapFs{}, pred)
	_, err := fs.Create("/another.html")

	c.Assert(err, qt.Not(qt.IsNil))
	// go 1.14.x:  "The system cannot find the file specified."
	// go 1.15.x:  "no such file or directory"
}

func TestFilterROFilenameFilterChain(t *testing.T) {
	c := qt.New(t)
	pred := func(name string) bool {
		return name == "/only.html"
	}

	rofs := afero.NewReadOnlyFs(&afero.MemMapFs{})
	fs := &FilenameFilterFs{pred: pred, source: rofs}
	_, err := fs.Create("/file.txt")
	c.Assert(err, qt.Not(qt.IsNil))
	// go 1.14.x:  "The system cannot find the file specified."
	// go 1.15.x:  "no such file or directory"
}

func TestFilenameFilterReadDir(t *testing.T) {
	c := qt.New(t)

	txtExts := func(name string) bool {
		return strings.HasSuffix(name, ".txt")
	}

	hasPrefixA := func(name string) bool {
		return strings.HasPrefix(name, "a")
	}

	mfs := &afero.MemMapFs{}
	fs1 := &FilenameFilterFs{pred: txtExts, source: mfs}
	fs := &FilenameFilterFs{pred: hasPrefixA, source: fs1}

	mfs.MkdirAll("/dir/sub", 0777)
	for _, name := range []string{"afile.txt", "afile.html", "bfile.txt"} {
		for _, dir := range []string{"/dir/", "/dir/sub/"} {
			fh, _ := mfs.Create(dir + name)
			fh.Close()
		}
	}

	files, _ := afero.ReadDir(fs, "/dir")

	// afile.txt, sub
	c.Assert(len(files), qt.Equals, 2)

	f, _ := fs.Open("/dir/sub")
	names, _ := f.Readdirnames(-1)
	c.Assert(len(names), qt.Equals, 1)
}
