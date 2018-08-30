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

package transform

import (
	"bytes"
	"io"

	bp "github.com/gohugoio/hugo/bufferpool"
)

// Transformer is the func that needs to be implemented by a transformation step.
type Transformer func(ft FromTo) error

// BytesReader wraps the Bytes method, usually implemented by bytes.Buffer, and an
// io.Reader.
type BytesReader interface {
	// The slice given by Bytes is valid for use only until the next buffer modification.
	// That is, if you want to use this value outside of the current transformer step,
	// you need to take a copy.
	Bytes() []byte

	io.Reader
}

// FromTo is sent to each transformation step in the chain.
type FromTo interface {
	From() BytesReader
	To() io.Writer
}

// Chain is an ordered processing chain. The next transform operation will
// receive the output from the previous.
type Chain []Transformer

// New creates a content transformer chain given the provided transform funcs.
func New(trs ...Transformer) Chain {
	return trs
}

// NewEmpty creates a new slice of transformers with a capacity of 20.
func NewEmpty() Chain {
	return make(Chain, 0, 20)
}

// Implements contentTransformer
// Content is read from the from-buffer and rewritten to to the to-buffer.
type fromToBuffer struct {
	from *bytes.Buffer
	to   *bytes.Buffer
}

func (ft fromToBuffer) From() BytesReader {
	return ft.from
}

func (ft fromToBuffer) To() io.Writer {
	return ft.to
}

// Apply passes the given from io.Reader through the transformation chain.
// The result is written to to.
func (c *Chain) Apply(to io.Writer, from io.Reader) error {
	if len(*c) == 0 {
		_, err := io.Copy(to, from)
		return err
	}

	b1 := bp.GetBuffer()
	defer bp.PutBuffer(b1)

	if _, err := b1.ReadFrom(from); err != nil {
		return err
	}

	b2 := bp.GetBuffer()
	defer bp.PutBuffer(b2)

	fb := &fromToBuffer{from: b1, to: b2}

	for i, tr := range *c {
		if i > 0 {
			if fb.from == b1 {
				fb.from = b2
				fb.to = b1
				fb.to.Reset()
			} else {
				fb.from = b1
				fb.to = b2
				fb.to.Reset()
			}
		}

		if err := tr(fb); err != nil {
			return err
		}
	}

	_, err := fb.to.WriteTo(to)
	return err
}
