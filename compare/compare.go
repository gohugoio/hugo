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

package compare

import (
	"reflect"
	"time"

	"github.com/gohugoio/hugo/common/hreflect"
	"github.com/gohugoio/hugo/common/htime"
	"github.com/gohugoio/hugo/common/types"
)

// Eqer can be used to determine if this value is equal to the other.
// The semantics of equals is that the two value are interchangeable
// in the Hugo templates.
type Eqer interface {
	// Eq returns whether this value is equal to the other.
	// This is for internal use.
	Eq(other any) bool
}

// ProbablyEqer is an equal check that may return false positives, but never
// a false negative.
type ProbablyEqer interface {
	// For internal use.
	ProbablyEq(other any) bool
}

// Comparer can be used to compare two values.
// This will be used when using the le, ge etc. operators in the templates.
// Compare returns -1 if the given version is less than, 0 if equal and 1 if greater than
// the running version.
type Comparer interface {
	Compare(other any) int
}

// Eq returns the boolean truth of arg1 == arg2 || arg1 == arg3 || arg1 == arg4.
// Numbers are compared by value regardless of type, so e.g. 1 and 1.0 are considered equal.
// It will use the Eqer interface if implemented, which
// defines equals when two value are interchangeable
// in the Hugo templates.
func Eq(first any, others ...any) bool {
	return EqInLocation(time.UTC, first, others...)
}

// EqInLocation is like Eq, but allows specifying a time.Location for time comparisons.
func EqInLocation(loc *time.Location, first any, others ...any) bool {
	normalize := func(v any) any {
		if types.IsNil(v) {
			return nil
		}
		if at, ok := v.(htime.AsTimeProvider); ok {
			return at.AsTime(loc)
		}
		return v
	}

	normFirst := normalize(first)
	fv := reflect.ValueOf(normFirst)
	for _, other := range others {
		if e, ok := first.(Eqer); ok {
			if e.Eq(other) {
				return true
			}
			continue
		}

		if e, ok := other.(Eqer); ok {
			if e.Eq(first) {
				return true
			}
			continue
		}

		other = normalize(other)
		if normFirst == nil || other == nil {
			if normFirst == other {
				return true
			}
			continue
		}

		ov := reflect.ValueOf(other)

		if fv.Kind() == reflect.String && ov.Kind() == reflect.String {
			if fv.String() == ov.String() {
				return true
			}
			continue
		}

		if c, ok := hreflect.CompareNumbers(fv, ov); ok {
			if c == 0 {
				return true
			}
			continue
		}

		if reflect.DeepEqual(normFirst, other) {
			return true
		}
	}

	return false
}

// ProbablyEq returns whether v1 is probably equal to v2.
func ProbablyEq(v1, v2 any) bool {
	if Eq(v1, v2) {
		return true
	}

	if peqer, ok := v1.(ProbablyEqer); ok {
		return peqer.ProbablyEq(v2)
	}

	return false
}
