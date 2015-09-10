// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Licensed under the Simple Public License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://opensource.org/licenses/Simple-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package source

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// LazyFileReader is an io.Reader implementation to postpone reading the file
// contents until it is really needed. It keeps filename and file contents once
// it is read.
type LazyFileReader struct {
	filename string
	contents *bytes.Reader
	pos      int64
}

// NewLazyFileReader creates and initializes a new LazyFileReader of filename.
// It checks whether the file can be opened. If it fails, it returns nil and an
// error.
func NewLazyFileReader(filename string) (*LazyFileReader, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return &LazyFileReader{filename: filename, contents: nil, pos: 0}, nil
}

// Filename returns a file name which LazyFileReader keeps
func (l *LazyFileReader) Filename() string {
	return l.filename
}

// Read reads up to len(p) bytes from the LazyFileReader's file and copies them
// into p. It returns the number of bytes read and any error encountered. If
// the file is once read, it returns its contents from cache, doesn't re-read
// the file.
func (l *LazyFileReader) Read(p []byte) (n int, err error) {
	if l.contents == nil {
		b, err := ioutil.ReadFile(l.filename)
		if err != nil {
			return 0, fmt.Errorf("failed to read content from %s: %s", l.filename, err.Error())
		}
		l.contents = bytes.NewReader(b)
	}
	l.contents.Seek(l.pos, 0)
	if err != nil {
		return 0, errors.New("failed to set read position: " + err.Error())
	}
	n, err = l.contents.Read(p)
	l.pos += int64(n)
	return n, err
}

// Seek implements the io.Seeker interface. Once reader contents is consumed by
// Read, WriteTo etc, to read it again, it must be rewinded by this function
func (l *LazyFileReader) Seek(offset int64, whence int) (pos int64, err error) {
	if l.contents == nil {
		switch whence {
		case 0:
			pos = offset
		case 1:
			pos = l.pos + offset
		case 2:
			fi, err := os.Stat(l.filename)
			if err != nil {
				return 0, fmt.Errorf("failed to get %q info: %s", l.filename, err.Error())
			}
			pos = fi.Size() + offset
		default:
			return 0, errors.New("invalid whence")
		}
		if pos < 0 {
			return 0, errors.New("negative position")
		}
	} else {
		pos, err = l.contents.Seek(offset, whence)
		if err != nil {
			return 0, err
		}
	}
	l.pos = pos
	return pos, nil
}

// WriteTo writes data to w until all the LazyFileReader's file contents is
// drained or an error occurs. If the file is once read, it just writes its
// read cache to w, doesn't re-read the file but this method itself doesn't try
// to keep the contents in cache.
func (l *LazyFileReader) WriteTo(w io.Writer) (n int64, err error) {
	if l.contents != nil {
		l.contents.Seek(l.pos, 0)
		if err != nil {
			return 0, errors.New("failed to set read position: " + err.Error())
		}
		n, err = l.contents.WriteTo(w)
		l.pos += n
		return n, err
	}
	f, err := os.Open(l.filename)
	if err != nil {
		return 0, fmt.Errorf("failed to open %s to read content: %s", l.filename, err.Error())
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return 0, fmt.Errorf("failed to get %q info: %s", l.filename, err.Error())
	}

	if l.pos >= fi.Size() {
		return 0, nil
	}

	// following code is taken from io.Copy in 'io/io.go'
	buf := make([]byte, 32*1024)
	for {
		nr, er := f.Read(buf)
		if nr > 0 {
			nw, ew := w.Write(buf[0:nr])
			if nw > 0 {
				l.pos += int64(nw)
				n += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er == io.EOF {
			break
		}
		if er != nil {
			err = er
			break
		}
	}
	return n, err
}
