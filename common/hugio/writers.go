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

type multiWriteCloser struct {
	io.Writer
	closers []io.WriteCloser
}

func (m multiWriteCloser) Close() error {
	var err error
	for _, c := range m.closers {
		if closeErr := c.Close(); err != nil {
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
