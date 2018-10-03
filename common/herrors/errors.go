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

// Package herrors contains common Hugo errors and error related utilities.
package herrors

import (
	"errors"
	"fmt"
	"io"
	"os"

	_errors "github.com/pkg/errors"
)

// As defined in https://godoc.org/github.com/pkg/errors
type causer interface {
	Cause() error
}

type stackTracer interface {
	StackTrace() _errors.StackTrace
}

// PrintStackTrace prints the error's stack trace to stdoud.
func PrintStackTrace(err error) {
	FprintStackTrace(os.Stdout, err)
}

// FprintStackTrace prints the error's stack trace to w.
func FprintStackTrace(w io.Writer, err error) {
	if err, ok := err.(stackTracer); ok {
		for _, f := range err.StackTrace() {
			fmt.Fprintf(w, "%+s:%d\n", f, f)
		}
	}
}

// ErrFeatureNotAvailable denotes that a feature is unavailable.
//
// We will, at least to begin with, make some Hugo features (SCSS with libsass) optional,
// and this error is used to signal those situations.
var ErrFeatureNotAvailable = errors.New("this feature is not available in your current Hugo version")
