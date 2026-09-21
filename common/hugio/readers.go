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

package hugio

import (
	"bytes"
	"errors"
	"io"
	"strings"

	"github.com/spf13/afero"
)

const BOM = "\xef\xbb\xbf"

// ReadSeekCloser is implemented by afero.File. We use this as the common type for
// content in Resource objects, even for strings.
type ReadSeekCloser interface {
	io.ReadSeeker
	io.Closer
}

// Sizer provides the size of, typically, a io.Reader.
// As implemented by e.g. os.File and io.SectionReader.
type Sizer interface {
	Size() int64
}

type SizeReader interface {
	io.Reader
	Sizer
}

// ToSizeReader converts the given io.Reader to a SizeReader.
// Note that if r is not a SizeReader, the entire content will be read into memory
func ToSizeReader(r io.Reader) (SizeReader, error) {
	if sr, ok := r.(SizeReader); ok {
		return sr, nil
	}
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}

// CloserFunc is an adapter to allow the use of ordinary functions as io.Closers.
type CloserFunc func() error

func (f CloserFunc) Close() error {
	return f()
}

// ReadSeekCloserProvider provides a ReadSeekCloser.
type ReadSeekCloserProvider interface {
	ReadSeekCloser() (ReadSeekCloser, error)
}

// readSeekerNopCloser implements ReadSeekCloser by doing nothing in Close.
type readSeekerNopCloser struct {
	io.ReadSeeker
}

// Close does nothing.
func (r readSeekerNopCloser) Close() error {
	return nil
}

// NewReadSeekerNoOpCloser creates a new ReadSeekerNoOpCloser with the given ReadSeeker.
func NewReadSeekerNoOpCloser(r io.ReadSeeker) ReadSeekCloser {
	return readSeekerNopCloser{r}
}

// NewReadSeekerNoOpCloserFromString uses strings.NewReader to create a new ReadSeekerNoOpCloser
// from the given string.
func NewReadSeekerNoOpCloserFromString(content string) ReadSeekCloser {
	return stringReadSeeker{s: content, readSeekerNopCloser: readSeekerNopCloser{strings.NewReader(content)}}
}

var _ StringReader = (*stringReadSeeker)(nil)

type stringReadSeeker struct {
	s string
	readSeekerNopCloser
}

func (s *stringReadSeeker) ReadString() string {
	return s.s
}

// StringReader provides a way to read a string.
type StringReader interface {
	ReadString() string
}

// NewReadSeekerNoOpCloserFromBytes uses bytes.NewReader to create a new ReadSeekerNoOpCloser
// from the given bytes slice.
func NewReadSeekerNoOpCloserFromBytes(content []byte) readSeekerNopCloser {
	return readSeekerNopCloser{bytes.NewReader(content)}
}

// NewReadSeekerNoOpCloserFromReader creates a new ReadSeekerNoOpCloser from the given io.Reader.
// If the given io.Reader is not an io.ReadSeeker, the entire content will be read into memory.
func NewReadSeekerNoOpCloserFromReader(r io.Reader) (readSeekerNopCloser, error) {
	var rs io.ReadSeeker
	if s, ok := r.(io.ReadSeeker); ok {
		rs = s
	} else {
		b, err := io.ReadAll(r)
		if err != nil {
			return readSeekerNopCloser{rs}, err
		}
		rs = bytes.NewReader(b)
	}
	return readSeekerNopCloser{rs}, nil
}

// NewOpenReadSeekCloser creates a new ReadSeekCloser from the given ReadSeeker.
// The ReadSeeker will be seeked to the beginning before returned.
func NewOpenReadSeekCloser(r ReadSeekCloser) OpenReadSeekCloser {
	return func() (ReadSeekCloser, error) {
		r.Seek(0, io.SeekStart)
		return r, nil
	}
}

// OpenReadSeekCloser allows setting some other way (than reading from a filesystem)
// to open or create a ReadSeekCloser.
type OpenReadSeekCloser func() (ReadSeekCloser, error)

// ReadString reads from the given reader and returns the content as a string.
func ReadString(r io.Reader) (string, error) {
	if sr, ok := r.(StringReader); ok {
		return sr.ReadString(), nil
	}
	b, err := ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ReadAll reads all content from the given reader and returns it as a byte slice.
// It also removes the UTF-8 BOM if present.
func ReadAll(r io.Reader) ([]byte, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return trimBOM(b), nil
}

// ReadFile reads the content of the given file from the provided afero.Fs and returns it as a byte slice.
// It also removes the UTF-8 BOM if present.
func ReadFile(fs afero.Fs, filename string) ([]byte, error) {
	b, err := afero.ReadFile(fs, filename)
	if err != nil {
		return nil, err
	}
	return trimBOM(b), nil
}

func trimBOM(b []byte) []byte {
	if len(b) >= 3 && b[0] == 0xef && b[1] == 0xbb && b[2] == 0xbf {
		return b[3:]
	}
	return b
}

var bomArr = [3]byte{0xef, 0xbb, 0xbf}

// NewSkipBOMReader returns r positioned after the UTF-8 BOM if present, else r as is.
func NewSkipBOMReader(r ReadSeekCloser) (ReadSeekCloser, error) {
	var buf [3]byte
	n, err := io.ReadFull(r, buf[:])
	if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
		return nil, err
	}
	if n == 3 && buf == bomArr {
		return &skipBOMReader{ReadSeekCloser: r}, nil
	}
	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	return r, nil
}

// skipBOMReader hides the 3 byte BOM at the start of the underlying reader.
// It is only created by NewSkipBOMReader after the BOM has been read.
type skipBOMReader struct {
	ReadSeekCloser
}

func (r *skipBOMReader) Seek(offset int64, whence int) (int64, error) {
	var abs int64
	switch whence {
	case io.SeekStart:
		abs = offset
	case io.SeekCurrent, io.SeekEnd:
		pos, err := r.ReadSeekCloser.Seek(0, whence)
		if err != nil {
			return 0, err
		}
		abs = pos - 3 + offset
	default:
		return 0, errors.New("skipBOMReader.Seek: invalid whence")
	}
	if abs < 0 {
		return 0, errors.New("skipBOMReader.Seek: negative position")
	}
	if _, err := r.ReadSeekCloser.Seek(abs+3, io.SeekStart); err != nil {
		return 0, err
	}
	return abs, nil
}
