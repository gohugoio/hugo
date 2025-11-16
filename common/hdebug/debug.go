// Copyright 2025 The Hugo Authors. All rights reserved.
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

package hdebug

import (
	"fmt"
	"strings"

	"github.com/gohugoio/hugo/common/types"
	"github.com/gohugoio/hugo/htesting"
)

// Printf is a debug print function that should be removed before committing code to the repository.
func Printf(format string, args ...any) {
	panicIfRealCI()
	if len(args) == 1 && !strings.Contains(format, "%") {
		format = format + ": %v"
	}
	if !strings.HasSuffix(format, "\n") {
		format = format + "\n"
	}
	fmt.Printf(format, args...)
}

func AssertNotNil(a ...any) {
	panicIfRealCI()
	for _, v := range a {
		if types.IsNil(v) {
			panic("hdebug.AssertNotNil: value is nil")
		}
	}
}

func Panicf(format string, args ...any) {
	panicIfRealCI()
	// fmt.Println(stack())
	if len(args) == 1 && !strings.Contains(format, "%") {
		format = format + ": %v"
	}
	if !strings.HasSuffix(format, "\n") {
		format = format + "\n"
	}
	panic(fmt.Sprintf(format, args...))
}

func panicIfRealCI() {
	if htesting.IsRealCI() {
		panic("This debug statement should be removed before committing code!")
	}
}
