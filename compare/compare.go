// Copyright 2019 The Hugo Authors. All rights reserved.
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

// Eq returns whether v1 is equal to v2.
// It will use the Eqer interface if implemented, which
// defines equals when two value are interchangeable
// in the Hugo templates.
func Eq(v1, v2 any) bool {
	if v1 == nil || v2 == nil {
		return v1 == v2
	}

	if eqer, ok := v1.(Eqer); ok {
		return eqer.Eq(v2)
	}

	return v1 == v2
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
