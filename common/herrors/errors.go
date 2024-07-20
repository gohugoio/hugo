// Copyright 2024 The Hugo Authors. All rights reserved.
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
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
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

// IsTimeoutError returns true if the given error is or contains a TimeoutError.
func IsTimeoutError(err error) bool {
	return errors.Is(err, &TimeoutError{})
}

type TimeoutError struct {
	Duration time.Duration
}

func (e *TimeoutError) Error() string {
	return fmt.Sprintf("timeout after %s", e.Duration)
}

func (e *TimeoutError) Is(target error) bool {
	_, ok := target.(*TimeoutError)
	return ok
}

// errMessage wraps an error with a message.
type errMessage struct {
	msg string
	err error
}

func (e *errMessage) Error() string {
	return e.msg
}

func (e *errMessage) Unwrap() error {
	return e.err
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

var nilPointerErrRe = regexp.MustCompile(`at <(.*)>: error calling (.*?): runtime error: invalid memory address or nil pointer dereference`)

const deferredPrefix = "__hdeferred/"

var deferredStringToRemove = regexp.MustCompile(`executing "__hdeferred/.*" `)

// ImproveRenderErr improves the error message for rendering errors.
func ImproveRenderErr(inErr error) (outErr error) {
	outErr = inErr
	msg := improveIfNilPointerMsg(inErr)
	if msg != "" {
		outErr = &errMessage{msg: msg, err: outErr}
	}

	if strings.Contains(inErr.Error(), deferredPrefix) {
		msg := deferredStringToRemove.ReplaceAllString(inErr.Error(), "executing ")
		outErr = &errMessage{msg: msg, err: outErr}
	}
	return
}

func improveIfNilPointerMsg(inErr error) string {
	m := nilPointerErrRe.FindStringSubmatch(inErr.Error())
	if len(m) == 0 {
		return ""
	}
	call := m[1]
	field := m[2]
	parts := strings.Split(call, ".")
	if len(parts) < 2 {
		return ""
	}
	receiverName := parts[len(parts)-2]
	receiver := strings.Join(parts[:len(parts)-1], ".")
	s := fmt.Sprintf("– %s is nil; wrap it in if or with: {{ with %s }}{{ .%s }}{{ end }}", receiverName, receiver, field)
	return nilPointerErrRe.ReplaceAllString(inErr.Error(), s)
}
