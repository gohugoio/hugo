// Copyright 2017 The Hugo Authors. All rights reserved.
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

// Package fmt provides template functions for formatting strings.
package fmt

import (
	_fmt "fmt"

	"github.com/gohugoio/hugo/common/loggers"

	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/helpers"
)

// New returns a new instance of the fmt-namespaced template functions.
func New(d *deps.Deps) *Namespace {
	ignorableLogger, ok := d.Log.(loggers.IgnorableLogger)
	if !ok {
		ignorableLogger = loggers.NewIgnorableLogger(d.Log, nil)
	}

	distinctLogger := helpers.NewDistinctLogger(d.Log)
	ns := &Namespace{
		distinctLogger: ignorableLogger.Apply(distinctLogger),
	}

	d.BuildStartListeners.Add(func() {
		ns.distinctLogger.Reset()
	})

	return ns
}

// Namespace provides template functions for the "fmt" namespace.
type Namespace struct {
	distinctLogger loggers.IgnorableLogger
}

// Print returns a string representation of args.
func (ns *Namespace) Print(args ...any) string {
	return _fmt.Sprint(args...)
}

// Printf returns string representation of args formatted with the layouut in format.
func (ns *Namespace) Printf(format string, args ...any) string {
	return _fmt.Sprintf(format, args...)
}

// Println returns string representation of args  ending with a newline.
func (ns *Namespace) Println(args ...any) string {
	return _fmt.Sprintln(args...)
}

// Errorf formats args according to a format specifier and logs an ERROR.
// It returns an empty string.
func (ns *Namespace) Errorf(format string, args ...any) string {
	ns.distinctLogger.Errorf(format, args...)
	return ""
}

// Erroridf formats args according to a format specifier and logs an ERROR and
// an information text that the error with the given id can be suppressed in config.
// It returns an empty string.
func (ns *Namespace) Erroridf(id, format string, args ...any) string {
	ns.distinctLogger.Errorsf(id, format, args...)
	return ""
}

// Warnf formats args according to a format specifier and logs a WARNING.
// It returns an empty string.
func (ns *Namespace) Warnf(format string, args ...any) string {
	ns.distinctLogger.Warnf(format, args...)
	return ""
}
