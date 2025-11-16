// Copyright 2024 The Hugo Authors. All rights reserved.
// Some functions in this file (see comments) is based on the Go source code,
// copyright The Go Authors and  governed by a BSD-style license.
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

package loggers

import (
	"bytes"
	"testing"

	"github.com/bep/logg"
	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/common/terminal"
)

func TestNoAnsiEscapeHandler(t *testing.T) {
	c := qt.New(t)

	test := func(s string) {
		c.Assert(stripANSI(terminal.Notice(s)), qt.Equals, s)
	}
	test(`error in "file.md:1:2"`)

	var buf bytes.Buffer
	h := newNoAnsiEscapeHandler(&buf, &buf, false, nil)
	h.HandleLog(&logg.Entry{Message: terminal.Notice(`error in "file.md:1:2"`), Level: logg.LevelInfo})

	c.Assert(buf.String(), qt.Equals, "INFO  error in \"file.md:1:2\"\n")
}
