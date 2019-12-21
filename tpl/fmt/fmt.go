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

	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/helpers"
)

// New returns a new instance of the fmt-namespaced template functions.
func New(d *deps.Deps) *Namespace {
	return &Namespace{
		errorLogger: helpers.NewDistinctLogger(d.Log.ERROR),
		warnLogger:  helpers.NewDistinctLogger(d.Log.WARN),
	}
}

// Namespace provides template functions for the "fmt" namespace.
type Namespace struct {
	errorLogger *helpers.DistinctLogger
	warnLogger  *helpers.DistinctLogger
}

// Print returns string representation of the passed arguments.
func (ns *Namespace) Print(a ...interface{}) string {
	return _fmt.Sprint(a...)
}

// Printf returns a formatted string representation of the passed arguments.
func (ns *Namespace) Printf(format string, a ...interface{}) string {
	return _fmt.Sprintf(format, a...)

}

// Println returns string representation of the passed arguments ending with a newline.
func (ns *Namespace) Println(a ...interface{}) string {
	return _fmt.Sprintln(a...)
}

// Errorf formats according to a format specifier and logs an ERROR.
// It returns an empty string.
func (ns *Namespace) Errorf(format string, a ...interface{}) string {
	ns.errorLogger.Printf(format, a...)
	return ""
}

// Warnf formats according to a format specifier and logs a WARNING.
// It returns an empty string.
func (ns *Namespace) Warnf(format string, a ...interface{}) string {
	ns.warnLogger.Printf(format, a...)
	return ""
}
