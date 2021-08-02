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

package loggers

import (
	"bytes"
	"fmt"
	"log"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestLogger(t *testing.T) {
	c := qt.New(t)
	l := NewWarningLogger()

	l.Errorln("One error")
	l.Errorln("Two error")
	l.Warnln("A warning")

	c.Assert(l.LogCounters().ErrorCounter.Count(), qt.Equals, uint64(2))
}

func TestLoggerToWriterWithPrefix(t *testing.T) {
	c := qt.New(t)

	var b bytes.Buffer

	logger := log.New(&b, "", 0)

	w := LoggerToWriterWithPrefix(logger, "myprefix")

	fmt.Fprint(w, "Hello Hugo!")

	c.Assert(b.String(), qt.Equals, "myprefix: Hello Hugo!\n")
}
