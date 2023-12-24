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

package loggers_test

import (
	"io"
	"strings"
	"testing"

	"github.com/bep/logg"
	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/common/loggers"
)

func TestLogDistinct(t *testing.T) {
	c := qt.New(t)

	opts := loggers.Options{
		DistinctLevel: logg.LevelWarn,
		StoreErrors:   true,
		Stdout:        io.Discard,
		Stderr:        io.Discard,
	}

	l := loggers.New(opts)

	for i := 0; i < 10; i++ {
		l.Errorln("error 1")
		l.Errorln("error 2")
		l.Warnln("warn 1")
	}
	c.Assert(strings.Count(l.Errors(), "error 1"), qt.Equals, 1)
	c.Assert(l.LoggCount(logg.LevelError), qt.Equals, 2)
	c.Assert(l.LoggCount(logg.LevelWarn), qt.Equals, 1)
}

func TestHookLast(t *testing.T) {
	c := qt.New(t)

	opts := loggers.Options{
		HandlerPost: func(e *logg.Entry) error {
			panic(e.Message)
		},
		Stdout: io.Discard,
		Stderr: io.Discard,
	}

	l := loggers.New(opts)

	c.Assert(func() { l.Warnln("warn 1") }, qt.PanicMatches, "warn 1")
}

func TestOptionStoreErrors(t *testing.T) {
	c := qt.New(t)

	var sb strings.Builder

	opts := loggers.Options{
		StoreErrors: true,
		Stderr:      &sb,
		Stdout:      &sb,
	}

	l := loggers.New(opts)
	l.Errorln("error 1")
	l.Errorln("error 2")

	errorsStr := l.Errors()

	c.Assert(errorsStr, qt.Contains, "error 1")
	c.Assert(errorsStr, qt.Not(qt.Contains), "ERROR")

	c.Assert(sb.String(), qt.Contains, "error 1")
	c.Assert(sb.String(), qt.Contains, "ERROR")
}

func TestLogCount(t *testing.T) {
	c := qt.New(t)

	opts := loggers.Options{
		StoreErrors: true,
	}

	l := loggers.New(opts)
	l.Errorln("error 1")
	l.Errorln("error 2")
	l.Warnln("warn 1")

	c.Assert(l.LoggCount(logg.LevelError), qt.Equals, 2)
	c.Assert(l.LoggCount(logg.LevelWarn), qt.Equals, 1)
	c.Assert(l.LoggCount(logg.LevelInfo), qt.Equals, 0)
}

func TestSuppressStatements(t *testing.T) {
	c := qt.New(t)

	opts := loggers.Options{
		StoreErrors: true,
		SuppressStatements: map[string]bool{
			"error-1": true,
		},
	}

	l := loggers.New(opts)
	l.Error().WithField(loggers.FieldNameStatementID, "error-1").Logf("error 1")
	l.Errorln("error 2")

	errorsStr := l.Errors()

	c.Assert(errorsStr, qt.Not(qt.Contains), "error 1")
	c.Assert(errorsStr, qt.Contains, "error 2")
	c.Assert(l.LoggCount(logg.LevelError), qt.Equals, 1)
}

func TestReset(t *testing.T) {
	c := qt.New(t)

	opts := loggers.Options{
		StoreErrors:   true,
		DistinctLevel: logg.LevelWarn,
		Stdout:        io.Discard,
		Stderr:        io.Discard,
	}

	l := loggers.New(opts)

	for i := 0; i < 3; i++ {
		l.Errorln("error 1")
		l.Errorln("error 2")
		l.Errorln("error 1")
		c.Assert(l.LoggCount(logg.LevelError), qt.Equals, 2)

		l.Reset()

		errorsStr := l.Errors()

		c.Assert(errorsStr, qt.Equals, "")
		c.Assert(l.LoggCount(logg.LevelError), qt.Equals, 0)

	}
}
