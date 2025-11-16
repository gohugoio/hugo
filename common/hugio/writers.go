// Copyright 2018 The Hugo Authors. All rights reserved.
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
	"io"
)

// As implemented by strings.Builder.
type FlexiWriter interface {
	io.Writer
	io.ByteWriter
	WriteString(s string) (int, error)
	WriteRune(r rune) (int, error)
}

type multiWriteCloser struct {
	io.Writer
	closers []io.WriteCloser
}

func (m multiWriteCloser) Close() error {
	var err error
	for _, c := range m.closers {
		if closeErr := c.Close(); closeErr != nil {
			err = closeErr
		}
	}
	return err
}

// NewMultiWriteCloser creates a new io.WriteCloser that duplicates its writes to all the
// provided writers.
func NewMultiWriteCloser(writeClosers ...io.WriteCloser) io.WriteCloser {
	writers := make([]io.Writer, len(writeClosers))
	for i, w := range writeClosers {
		writers[i] = w
	}
	return multiWriteCloser{Writer: io.MultiWriter(writers...), closers: writeClosers}
}

// ToWriteCloser creates an io.WriteCloser from the given io.Writer.
// If it's not already, one will be created with a Close method that does nothing.
func ToWriteCloser(w io.Writer) io.WriteCloser {
	if rw, ok := w.(io.WriteCloser); ok {
		return rw
	}

	return struct {
		io.Writer
		io.Closer
	}{
		w,
		io.NopCloser(nil),
	}
}

// ToReadCloser creates an io.ReadCloser from the given io.Reader.
// If it's not already, one will be created with a Close method that does nothing.
func ToReadCloser(r io.Reader) io.ReadCloser {
	if rc, ok := r.(io.ReadCloser); ok {
		return rc
	}

	return struct {
		io.Reader
		io.Closer
	}{
		r,
		io.NopCloser(nil),
	}
}

type ReadWriteCloser interface {
	io.Reader
	io.Writer
	io.Closer
}

// PipeReadWriteCloser is a convenience type to create a pipe with a ReadCloser and a WriteCloser.
type PipeReadWriteCloser struct {
	*io.PipeReader
	*io.PipeWriter
}

// NewPipeReadWriteCloser creates a new PipeReadWriteCloser.
func NewPipeReadWriteCloser() PipeReadWriteCloser {
	pr, pw := io.Pipe()
	return PipeReadWriteCloser{pr, pw}
}

func (c PipeReadWriteCloser) Close() (err error) {
	if err = c.PipeReader.Close(); err != nil {
		return
	}
	err = c.PipeWriter.Close()
	return
}

func (c PipeReadWriteCloser) WriteString(s string) (int, error) {
	return c.PipeWriter.Write([]byte(s))
}
