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

// package loggers contains some basic logging setup.
package loggers

import (
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/bep/logg"

	"github.com/fatih/color"
)

// levelColor mapping.
var levelColor = [...]*color.Color{
	logg.LevelTrace: color.New(color.FgWhite),
	logg.LevelDebug: color.New(color.FgWhite),
	logg.LevelInfo:  color.New(color.FgBlue),
	logg.LevelWarn:  color.New(color.FgYellow),
	logg.LevelError: color.New(color.FgRed),
}

// levelString mapping.
var levelString = [...]string{
	logg.LevelTrace: "TRACE",
	logg.LevelDebug: "DEBUG",
	logg.LevelInfo:  "INFO ",
	logg.LevelWarn:  "WARN ",
	logg.LevelError: "ERROR",
}

// newDefaultHandler handler.
func newDefaultHandler(outWriter, errWriter io.Writer) logg.Handler {
	return &defaultHandler{
		outWriter: outWriter,
		errWriter: errWriter,
		Padding:   0,
	}
}

// Default Handler implementation.
// Based on https://github.com/apex/log/blob/master/handlers/cli/cli.go
type defaultHandler struct {
	mu        sync.Mutex
	outWriter io.Writer // Defaults to os.Stdout.
	errWriter io.Writer // Defaults to os.Stderr.

	Padding int
}

// HandleLog implements logg.Handler.
func (h *defaultHandler) HandleLog(e *logg.Entry) error {
	color := levelColor[e.Level]
	level := levelString[e.Level]

	h.mu.Lock()
	defer h.mu.Unlock()

	var w io.Writer
	if e.Level > logg.LevelInfo {
		w = h.errWriter
	} else {
		w = h.outWriter
	}

	var prefix string
	for _, field := range e.Fields {
		if field.Name == FieldNameCmd {
			prefix = fmt.Sprint(field.Value)
			break
		}
	}

	if prefix != "" {
		prefix = prefix + ": "
	}

	color.Fprintf(w, "%s %s%s", fmt.Sprintf("%*s", h.Padding+1, level), color.Sprint(prefix), e.Message)

	for _, field := range e.Fields {
		if strings.HasPrefix(field.Name, reservedFieldNamePrefix) {
			continue
		}
		fmt.Fprintf(w, " %s %v", color.Sprint(field.Name), field.Value)
	}

	fmt.Fprintln(w)

	return nil
}
