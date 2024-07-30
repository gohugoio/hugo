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
	"fmt"
	"strings"
	"sync"

	"github.com/bep/logg"
	"github.com/gohugoio/hugo/common/hashing"
)

// PanicOnWarningHook panics on warnings.
var PanicOnWarningHook = func(e *logg.Entry) error {
	if e.Level != logg.LevelWarn {
		return nil
	}
	panic(e.Message)
}

func newLogLevelCounter() *logLevelCounter {
	return &logLevelCounter{
		counters: make(map[logg.Level]int),
	}
}

func newLogOnceHandler(threshold logg.Level) *logOnceHandler {
	return &logOnceHandler{
		threshold: threshold,
		seen:      make(map[uint64]bool),
	}
}

func newStopHandler(h ...logg.Handler) *stopHandler {
	return &stopHandler{
		handlers: h,
	}
}

func newSuppressStatementsHandler(statements map[string]bool) *suppressStatementsHandler {
	return &suppressStatementsHandler{
		statements: statements,
	}
}

type logLevelCounter struct {
	mu       sync.RWMutex
	counters map[logg.Level]int
}

func (h *logLevelCounter) HandleLog(e *logg.Entry) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.counters[e.Level]++
	return nil
}

var errStop = fmt.Errorf("stop")

type logOnceHandler struct {
	threshold logg.Level
	mu        sync.Mutex
	seen      map[uint64]bool
}

func (h *logOnceHandler) HandleLog(e *logg.Entry) error {
	if e.Level < h.threshold {
		// We typically only want to enable this for warnings and above.
		// The common use case is that many go routines may log the same error.
		return nil
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	hash := hashing.HashUint64(e.Level, e.Message, e.Fields)
	if h.seen[hash] {
		return errStop
	}
	h.seen[hash] = true
	return nil
}

func (h *logOnceHandler) reset() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.seen = make(map[uint64]bool)
}

type stopHandler struct {
	handlers []logg.Handler
}

// HandleLog implements logg.Handler.
func (h *stopHandler) HandleLog(e *logg.Entry) error {
	for _, handler := range h.handlers {
		if err := handler.HandleLog(e); err != nil {
			if err == errStop {
				return nil
			}
			return err
		}
	}
	return nil
}

type suppressStatementsHandler struct {
	statements map[string]bool
}

func (h *suppressStatementsHandler) HandleLog(e *logg.Entry) error {
	for _, field := range e.Fields {
		if field.Name == FieldNameStatementID {
			if h.statements[field.Value.(string)] {
				return errStop
			}
		}
	}
	return nil
}

// whiteSpaceTrimmer creates a new log handler that trims whitespace from log messages and string fields.
func whiteSpaceTrimmer() logg.Handler {
	return logg.HandlerFunc(func(e *logg.Entry) error {
		e.Message = strings.TrimSpace(e.Message)
		for i, field := range e.Fields {
			if s, ok := field.Value.(string); ok {
				e.Fields[i].Value = strings.TrimSpace(s)
			}
		}
		return nil
	})
}
