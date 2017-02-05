// Copyright 2015 The Hugo Authors. All rights reserved.
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

package source

import (
	"bytes"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestEmptySourceFilesystem(t *testing.T) {
	ss := newTestSourceSpec()
	src := ss.NewFilesystem("Empty")
	if len(src.Files()) != 0 {
		t.Errorf("new filesystem should contain 0 files.")
	}
}

type TestPath struct {
	filename string
	logical  string
	content  string
	section  string
	dir      string
}

func TestAddFile(t *testing.T) {
	ss := newTestSourceSpec()
	tests := platformPaths
	for _, test := range tests {
		base := platformBase
		srcDefault := ss.NewFilesystem("")
		srcWithBase := ss.NewFilesystem(base)

		for _, src := range []*Filesystem{srcDefault, srcWithBase} {

			p := test.filename
			if !filepath.IsAbs(test.filename) {
				p = filepath.Join(src.Base, test.filename)
			}

			if err := src.add(p, bytes.NewReader([]byte(test.content))); err != nil {
				if err.Error() == "source: missing base directory" {
					continue
				}
				t.Fatalf("%s add returned an error: %s", p, err)
			}

			if len(src.Files()) != 1 {
				t.Fatalf("%s Files() should return 1 file", p)
			}

			f := src.Files()[0]
			if f.LogicalName() != test.logical {
				t.Errorf("Filename (Base: %q) expected: %q, got: %q", src.Base, test.logical, f.LogicalName())
			}

			b := new(bytes.Buffer)
			b.ReadFrom(f.Contents)
			if b.String() != test.content {
				t.Errorf("File (Base: %q) contents should be %q, got: %q", src.Base, test.content, b.String())
			}

			if f.Section() != test.section {
				t.Errorf("File section (Base: %q) expected: %q, got: %q", src.Base, test.section, f.Section())
			}

			if f.Dir() != test.dir {
				t.Errorf("Dir path (Base: %q) expected: %q, got: %q", src.Base, test.dir, f.Dir())
			}
		}
	}
}

func TestUnicodeNorm(t *testing.T) {
	if runtime.GOOS != "darwin" {
		// Normalization code is only for Mac OS, since it is not necessary for other OSes.
		return
	}

	paths := []struct {
		NFC string
		NFD string
	}{
		{NFC: "å", NFD: "\x61\xcc\x8a"},
		{NFC: "é", NFD: "\x65\xcc\x81"},
	}

	ss := newTestSourceSpec()

	for _, path := range paths {
		src := ss.NewFilesystem("")
		_ = src.add(path.NFD, strings.NewReader(""))
		f := src.Files()[0]
		if f.BaseFileName() != path.NFC {
			t.Fatalf("file name in NFD form should be normalized (%s)", path.NFC)
		}
	}

}
