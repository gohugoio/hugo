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
	"io"
	"os"
	"testing"

	"github.com/spf13/afero"
)

func TestNewLazyFileReader(t *testing.T) {
	fs := afero.NewOsFs()
	filename := "itdoesnotexistfile"
	_, err := NewLazyFileReader(fs, filename)
	if err == nil {
		t.Errorf("NewLazyFileReader %s: error expected but no error is returned", filename)
	}

	filename = "lazy_file_reader_test.go"
	_, err = NewLazyFileReader(fs, filename)
	if err != nil {
		t.Errorf("NewLazyFileReader %s: %v", filename, err)
	}
}

func TestFilename(t *testing.T) {
	fs := afero.NewOsFs()
	filename := "lazy_file_reader_test.go"
	rd, err := NewLazyFileReader(fs, filename)
	if err != nil {
		t.Fatalf("NewLazyFileReader %s: %v", filename, err)
	}
	if rd.Filename() != filename {
		t.Errorf("Filename: expected filename %q, got %q", filename, rd.Filename())
	}
}

func TestRead(t *testing.T) {
	fs := afero.NewOsFs()
	filename := "lazy_file_reader_test.go"
	fi, err := fs.Stat(filename)
	if err != nil {
		t.Fatalf("os.Stat: %v", err)
	}

	b, err := afero.ReadFile(fs, filename)
	if err != nil {
		t.Fatalf("afero.ReadFile: %v", err)
	}

	rd, err := NewLazyFileReader(fs, filename)
	if err != nil {
		t.Fatalf("NewLazyFileReader %s: %v", filename, err)
	}

	tst := func(testcase string) {
		p := make([]byte, fi.Size())
		n, err := rd.Read(p)
		if err != nil {
			t.Fatalf("Read %s case: %v", testcase, err)
		}
		if int64(n) != fi.Size() {
			t.Errorf("Read %s case: read bytes length expected %d, got %d", testcase, fi.Size(), n)
		}
		if !bytes.Equal(b, p) {
			t.Errorf("Read %s case: read bytes are different from expected", testcase)
		}
	}
	tst("No cache")
	_, err = rd.Seek(0, 0)
	if err != nil {
		t.Fatalf("Seek: %v", err)
	}
	tst("Cache")
}

func TestSeek(t *testing.T) {
	type testcase struct {
		seek     int
		offset   int64
		length   int
		moveto   int64
		expected []byte
	}
	fs := afero.NewOsFs()
	filename := "lazy_file_reader_test.go"
	b, err := afero.ReadFile(fs, filename)
	if err != nil {
		t.Fatalf("afero.ReadFile: %v", err)
	}

	// no cache case
	for i, this := range []testcase{
		{seek: os.SEEK_SET, offset: 0, length: 10, moveto: 0, expected: b[:10]},
		{seek: os.SEEK_SET, offset: 5, length: 10, moveto: 5, expected: b[5:15]},
		{seek: os.SEEK_CUR, offset: 5, length: 10, moveto: 5, expected: b[5:15]}, // current pos = 0
		{seek: os.SEEK_END, offset: -1, length: 1, moveto: int64(len(b) - 1), expected: b[len(b)-1:]},
		{seek: 3, expected: nil},
		{seek: os.SEEK_SET, offset: -1, expected: nil},
	} {
		rd, err := NewLazyFileReader(fs, filename)
		if err != nil {
			t.Errorf("[%d] NewLazyFileReader %s: %v", i, filename, err)
			continue
		}

		pos, err := rd.Seek(this.offset, this.seek)
		if this.expected == nil {
			if err == nil {
				t.Errorf("[%d] Seek didn't return an expected error", i)
			}
		} else {
			if err != nil {
				t.Errorf("[%d] Seek failed unexpectedly: %v", i, err)
				continue
			}
			if pos != this.moveto {
				t.Errorf("[%d] Seek failed to move the pointer: got %d, expected: %d", i, pos, this.moveto)
			}

			buf := make([]byte, this.length)
			n, err := rd.Read(buf)
			if err != nil {
				t.Errorf("[%d] Read failed unexpectedly: %v", i, err)
			}
			if !bytes.Equal(this.expected, buf[:n]) {
				t.Errorf("[%d] Seek and Read got %q but expected %q", i, buf[:n], this.expected)
			}
		}
	}

	// cache case
	rd, err := NewLazyFileReader(fs, filename)
	if err != nil {
		t.Fatalf("NewLazyFileReader %s: %v", filename, err)
	}
	dummy := make([]byte, len(b))
	_, err = rd.Read(dummy)
	if err != nil {
		t.Fatalf("Read failed unexpectedly: %v", err)
	}

	for i, this := range []testcase{
		{seek: os.SEEK_SET, offset: 0, length: 10, moveto: 0, expected: b[:10]},
		{seek: os.SEEK_SET, offset: 5, length: 10, moveto: 5, expected: b[5:15]},
		{seek: os.SEEK_CUR, offset: 1, length: 10, moveto: 16, expected: b[16:26]}, // current pos = 15
		{seek: os.SEEK_END, offset: -1, length: 1, moveto: int64(len(b) - 1), expected: b[len(b)-1:]},
		{seek: 3, expected: nil},
		{seek: os.SEEK_SET, offset: -1, expected: nil},
	} {
		pos, err := rd.Seek(this.offset, this.seek)
		if this.expected == nil {
			if err == nil {
				t.Errorf("[%d] Seek didn't return an expected error", i)
			}
		} else {
			if err != nil {
				t.Errorf("[%d] Seek failed unexpectedly: %v", i, err)
				continue
			}
			if pos != this.moveto {
				t.Errorf("[%d] Seek failed to move the pointer: got %d, expected: %d", i, pos, this.moveto)
			}

			buf := make([]byte, this.length)
			n, err := rd.Read(buf)
			if err != nil {
				t.Errorf("[%d] Read failed unexpectedly: %v", i, err)
			}
			if !bytes.Equal(this.expected, buf[:n]) {
				t.Errorf("[%d] Seek and Read got %q but expected %q", i, buf[:n], this.expected)
			}
		}
	}
}

func TestWriteTo(t *testing.T) {
	fs := afero.NewOsFs()
	filename := "lazy_file_reader_test.go"
	fi, err := fs.Stat(filename)
	if err != nil {
		t.Fatalf("os.Stat: %v", err)
	}

	b, err := afero.ReadFile(fs, filename)
	if err != nil {
		t.Fatalf("afero.ReadFile: %v", err)
	}

	rd, err := NewLazyFileReader(fs, filename)
	if err != nil {
		t.Fatalf("NewLazyFileReader %s: %v", filename, err)
	}

	tst := func(testcase string, expectedSize int64, checkEqual bool) {
		buf := bytes.NewBuffer(make([]byte, 0, bytes.MinRead))
		n, err := rd.WriteTo(buf)
		if err != nil {
			t.Fatalf("WriteTo %s case: %v", testcase, err)
		}
		if n != expectedSize {
			t.Errorf("WriteTo %s case: written bytes length expected %d, got %d", testcase, expectedSize, n)
		}
		if checkEqual && !bytes.Equal(b, buf.Bytes()) {
			t.Errorf("WriteTo %s case: written bytes are different from expected", testcase)
		}
	}
	tst("No cache", fi.Size(), true)
	tst("No cache 2nd", 0, false)

	p := make([]byte, fi.Size())
	_, err = rd.Read(p)
	if err != nil && err != io.EOF {
		t.Fatalf("Read: %v", err)
	}
	_, err = rd.Seek(0, 0)
	if err != nil {
		t.Fatalf("Seek: %v", err)
	}

	tst("Cache", fi.Size(), true)
}
