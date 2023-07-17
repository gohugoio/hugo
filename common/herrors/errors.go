// Copyright 2022 The Hugo Authors. All rights reserved.
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

// Package herrors contains common Hugo errors and error related utilities.
package herrors

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
)

// PrintStackTrace prints the current stacktrace to w.
func PrintStackTrace(w io.Writer) {
	buf := make([]byte, 1<<16)
	runtime.Stack(buf, true)
	fmt.Fprintf(w, "%s", buf)
}

// ErrorSender is a, typically, non-blocking error handler.
type ErrorSender interface {
	SendError(err error)
}

// Recover is a helper function that can be used to capture panics.
// Put this at the top of a method/function that crashes in a template:
//
//	defer herrors.Recover()
func Recover(args ...any) {
	if r := recover(); r != nil {
		fmt.Println("ERR:", r)
		args = append(args, "stacktrace from panic: \n"+string(debug.Stack()), "\n")
		fmt.Println(args...)
	}
}

// GetGID the current goroutine id. Used only for debugging.
func GetGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

// IsFeatureNotAvailableError returns true if the given error is or contains a FeatureNotAvailableError.
func IsFeatureNotAvailableError(err error) bool {
	return errors.Is(err, &FeatureNotAvailableError{})
}

// ErrFeatureNotAvailable denotes that a feature is unavailable.
//
// We will, at least to begin with, make some Hugo features (SCSS with libsass) optional,
// and this error is used to signal those situations.
var ErrFeatureNotAvailable = &FeatureNotAvailableError{Cause: errors.New("this feature is not available in your current Hugo version, see https://goo.gl/YMrWcn for more information")}

// FeatureNotAvailableError is an error type used to signal that a feature is not available.
type FeatureNotAvailableError struct {
	Cause error
}

func (e *FeatureNotAvailableError) Unwrap() error {
	return e.Cause
}

func (e *FeatureNotAvailableError) Error() string {
	return e.Cause.Error()
}

func (e *FeatureNotAvailableError) Is(target error) bool {
	_, ok := target.(*FeatureNotAvailableError)
	return ok
}

// Must panics if err != nil.
func Must(err error) {
	if err != nil {
		panic(err)
	}
}

// IsNotExist returns true if the error is a file not found error.
// Unlike os.IsNotExist, this also considers wrapped errors.
func IsNotExist(err error) bool {
	if os.IsNotExist(err) {
		return true
	}

	// os.IsNotExist does not consider wrapped errors.
	if os.IsNotExist(errors.Unwrap(err)) {
		return true
	}

	return false
}
