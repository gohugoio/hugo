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

package doctree

const (
	// Dimensions in the Hugo build matrix.
	DimensionLanguage DimensionFlag = iota + 1
	DimensionVersion
	DimensionRole
)

// Dimensions is a row in the Hugo build matrix which currently has three values: language, version and role, in that order.
type Dimensions [3]int

// Compare returns -1 if this dimension is less than the given dimension, 0 if they are equal, and 1 if this dimension is greater than the given dimension.
// This adds a impicit weighting to the dimensions, where the first dimension is the most important,
// but this is just used for sorting to get stable output.
func (d Dimensions) Compare(e Dimensions) int {
	// note that a and b will never be equal.
	minusOneOrOne := func(a, b int) int {
		if a < b {
			return -1
		}
		return 1
	}
	if d[0] != e[0] {
		return minusOneOrOne(d[0], e[0])
	}
	if d[1] != e[1] {
		return minusOneOrOne(d[1], e[1])
	}
	if d[2] != e[2] {
		return minusOneOrOne(d[2], e[2])
	}
	// They are equal.
	return 0
}

// Distance returns the distance between this dimension and the given dimension
// ussing the first dimension that is different.
func (d Dimensions) Distance(e Dimensions) int {
	if d[0] != e[0] {
		return d[0] - e[0]
	}
	if d[1] != e[1] {
		return d[1] - e[1]
	}
	if d[2] != e[2] {
		return d[2] - e[2]
	}
	return 0
}

// Language returns the language dimension.
func (d Dimensions) Language() int {
	return d[DimensionLanguage.Index()]
}

// Version returns the version dimension.
func (d Dimensions) Version() int {
	return d[DimensionVersion.Index()]
}

// Role returns the role dimension.
func (d Dimensions) Role() int {
	return d[DimensionRole.Index()]
}

// DimensionFlag is a flag in the Hugo build matrix.
type DimensionFlag int8

// Has returns whether the given flag is set.
func (d DimensionFlag) Has(o DimensionFlag) bool {
	return d&o == o
}

// Set sets the given flag.
func (d DimensionFlag) Set(o DimensionFlag) DimensionFlag {
	return d | o
}

// Index returns this flag's index in the Dimensions array.
func (d DimensionFlag) Index() int {
	if d == 0 {
		panic("dimension flag not set")
	}
	return int(d - 1)
}
