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
	"io"
	"strings"
)

// ReadSeeker wraps io.Reader and io.Seeker.
type ReadSeeker interface {
	io.Reader
	io.Seeker
}

// ReadSeekCloser is implemented by afero.File. We use this as the common type for
// content in Resource objects, even for strings.
type ReadSeekCloser interface {
	ReadSeeker
	io.Closer
}

// ReadSeekCloserProvider provides a ReadSeekCloser.
type ReadSeekCloserProvider interface {
	ReadSeekCloser() (ReadSeekCloser, error)
}

// readSeekerNopCloser implements ReadSeekCloser by doing nothing in Close.
type readSeekerNopCloser struct {
	ReadSeeker
}

// Close does nothing.
func (r readSeekerNopCloser) Close() error {
	return nil
}

// NewReadSeekerNoOpCloser creates a new ReadSeekerNoOpCloser with the given ReadSeeker.
func NewReadSeekerNoOpCloser(r ReadSeeker) ReadSeekCloser {
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

// NewReadSeekerNoOpCloserFromString uses strings.NewReader to create a new ReadSeekerNoOpCloser
// from the given bytes slice.
func NewReadSeekerNoOpCloserFromBytes(content []byte) readSeekerNopCloser {
	return readSeekerNopCloser{bytes.NewReader(content)}
}

// NewReadSeekCloser creates a new ReadSeekCloser from the given ReadSeeker.
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
	b, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
