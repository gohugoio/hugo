// Copyright 2021 The Hugo Authors. All rights reserved.
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
)

// DistinctLogger ignores duplicate log statements.
type DistinctLogger struct {
	Logger
	sync.RWMutex
	m map[string]bool
}

func (l *DistinctLogger) Reset() {
	l.Lock()
	defer l.Unlock()

	l.m = make(map[string]bool)
}

// Println will log the string returned from fmt.Sprintln given the arguments,
// but not if it has been logged before.
func (l *DistinctLogger) Println(v ...interface{}) {
	// fmt.Sprint doesn't add space between string arguments
	logStatement := strings.TrimSpace(fmt.Sprintln(v...))
	l.printIfNotPrinted("println", logStatement, func() {
		l.Logger.Println(logStatement)
	})
}

// Printf will log the string returned from fmt.Sprintf given the arguments,
// but not if it has been logged before.
func (l *DistinctLogger) Printf(format string, v ...interface{}) {
	logStatement := fmt.Sprintf(format, v...)
	l.printIfNotPrinted("printf", logStatement, func() {
		l.Logger.Printf(format, v...)
	})
}

func (l *DistinctLogger) Debugf(format string, v ...interface{}) {
	logStatement := fmt.Sprintf(format, v...)
	l.printIfNotPrinted("debugf", logStatement, func() {
		l.Logger.Debugf(format, v...)
	})
}

func (l *DistinctLogger) Debugln(v ...interface{}) {
	logStatement := fmt.Sprint(v...)
	l.printIfNotPrinted("debugln", logStatement, func() {
		l.Logger.Debugln(v...)
	})
}

func (l *DistinctLogger) Infof(format string, v ...interface{}) {
	logStatement := fmt.Sprintf(format, v...)
	l.printIfNotPrinted("info", logStatement, func() {
		l.Logger.Infof(format, v...)
	})
}

func (l *DistinctLogger) Infoln(v ...interface{}) {
	logStatement := fmt.Sprint(v...)
	l.printIfNotPrinted("infoln", logStatement, func() {
		l.Logger.Infoln(v...)
	})
}

func (l *DistinctLogger) Warnf(format string, v ...interface{}) {
	logStatement := fmt.Sprintf(format, v...)
	l.printIfNotPrinted("warnf", logStatement, func() {
		l.Logger.Warnf(format, v...)
	})
}
func (l *DistinctLogger) Warnln(v ...interface{}) {
	logStatement := fmt.Sprint(v...)
	l.printIfNotPrinted("warnln", logStatement, func() {
		l.Logger.Warnln(v...)
	})
}
func (l *DistinctLogger) Errorf(format string, v ...interface{}) {
	logStatement := fmt.Sprint(v...)
	l.printIfNotPrinted("errorf", logStatement, func() {
		l.Logger.Errorf(format, v...)
	})
}

func (l *DistinctLogger) Errorln(v ...interface{}) {
	logStatement := fmt.Sprint(v...)
	l.printIfNotPrinted("errorln", logStatement, func() {
		l.Logger.Errorln(v...)
	})
}

func (l *DistinctLogger) hasPrinted(key string) bool {
	l.RLock()
	defer l.RUnlock()
	_, found := l.m[key]
	return found
}

func (l *DistinctLogger) printIfNotPrinted(level, logStatement string, print func()) {
	key := level + logStatement
	if l.hasPrinted(key) {
		return
	}
	l.Lock()
	print()
	l.m[key] = true
	l.Unlock()
}

// NewDistinctErrorLogger creates a new DistinctLogger that logs ERRORs
func NewDistinctErrorLogger() Logger {
	return &DistinctLogger{m: make(map[string]bool), Logger: NewErrorLogger()}
}

// NewDistinctLogger creates a new DistinctLogger that logs to the provided logger.
func NewDistinctLogger(logger Logger) Logger {
	return &DistinctLogger{m: make(map[string]bool), Logger: logger}
}

// NewDistinctWarnLogger creates a new DistinctLogger that logs WARNs
func NewDistinctWarnLogger() Logger {
	return &DistinctLogger{m: make(map[string]bool), Logger: NewWarningLogger()}
}
