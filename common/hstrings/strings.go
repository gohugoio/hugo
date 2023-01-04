// Copyright 2023 The Hugo Authors. All rights reserved.
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

package hstrings

import (
	"fmt"
	"strings"

	"github.com/gohugoio/hugo/compare"
)

var _ compare.Eqer = StringEqualFold("")

// StringEqualFold is a string that implements the compare.Eqer interface and considers
// two strings equal if they are equal when folded to lower case.
// The compare.Eqer interface is used in Hugo to compare values in templates (e.g. using the eq template function).
type StringEqualFold string

func (s StringEqualFold) EqualFold(s2 string) bool {
	return strings.EqualFold(string(s), s2)
}

func (s StringEqualFold) String() string {
	return string(s)
}

func (s StringEqualFold) Eq(s2 any) bool {
	switch ss := s2.(type) {
	case string:
		return s.EqualFold(ss)
	case fmt.Stringer:
		return s.EqualFold(ss.String())
	}

	return false
}

// EqualAny returns whether a string is equal to any of the given strings.
func EqualAny(a string, b ...string) bool {
	for _, s := range b {
		if a == s {
			return true
		}
	}
	return false
}
