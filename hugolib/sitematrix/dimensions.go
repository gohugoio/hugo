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

package sitematrix

import (
	"fmt"
)

const (
	// Dimensions in the Hugo build matrix.
	Language Dimension = iota + 1
	Version
	Role
)

// Vector represents a dimension vector in the Hugo build matrix from the three dimensions:
// Language, Version and Role.
type Vector [3]int

// Compare returns -1 if v1 is less than v2, 0 if they are equal, and 1 if v1 is greater than v2.
// This adds a impicit weighting to the dimensions, where the first dimension is the most important,
// but this is just used for sorting to get stable output.
func (v1 Vector) Compare(v2 Vector) int {
	// note that a and b will never be equal.
	minusOneOrOne := func(a, b int) int {
		if a < b {
			return -1
		}
		return 1
	}
	if v1[0] != v2[0] {
		return minusOneOrOne(v1[0], v2[0])
	}
	if v1[1] != v2[1] {
		return minusOneOrOne(v1[1], v2[1])
	}
	if v1[2] != v2[2] {
		return minusOneOrOne(v1[2], v2[2])
	}
	// They are equal.
	return 0
}

// DotProduct returns the dot product of two vectors.
func (v1 Vector) DotProduct(v2 Vector) int {
	dot := 0
	for i := range v1 {
		dot += v1[i] * v2[i]
	}
	return dot
}

// Distance returns the distance between v1 and v2
// ussing the first dimension that is different.
func (v1 Vector) Distance(v2 Vector) int {
	if v1[0] != v2[0] {
		return v1[0] - v2[0]
	}
	if v1[1] != v2[1] {
		return v1[1] - v2[1]
	}
	if v1[2] != v2[2] {
		return v1[2] - v2[2]
	}
	return 0
}

// EuclideanDistanceSquared returns the Euclidean distance between two vectors as the sum of the squared differences.
func (v1 Vector) EuclideanDistanceSquared(v2 Vector) int {
	sum := 0
	for i := range v1 {
		diff := v1[i] - v2[i]
		sum += diff * diff
	}
	return sum
}

func (v1 Vector) HasVector(v2 Vector) bool {
	return v1 == v2
}

func (v1 Vector) LenVectors() int {
	return 1
}

func (v1 Vector) FirstVector() Vector {
	return v1
}

func (v1 Vector) EqualsVector(other VectorProvider) bool {
	if other.LenVectors() != 1 {
		return false
	}
	return other.FirstVector() == v1
}

func (v1 Vector) ForEeachVector(yield func(v Vector) bool) bool {
	return yield(v1)
}

// Language returns the language dimension.
func (v1 Vector) Language() int {
	return v1[Language.Index()]
}

// Version returns the version dimension.
func (v1 Vector) Version() int {
	return v1[Version.Index()]
}

// Role returns the role dimension.
func (v1 Vector) Role() int {
	return v1[Role.Index()]
}

func (v1 Vector) Weight() int {
	return 0
}

func (v1 Vector) Ordinal() int {
	return 0
}

type VectorProvider interface {
	// HasVector returns true if the given vector is contained in the provider.
	// Used for membership testing of files, resources and pages.
	HasVector(v Vector) bool

	// LenVectors returns the number of vectors in the provider.
	LenVectors() int

	// Get returns the first vector in the provider.
	// This will panic if the provider is empty.
	FirstVector() Vector

	// Equals returns true if this provider is equal to the other provider.
	EqualsVector(other VectorProvider) bool

	// ForEeachVector iterates over all vectors in the provider.
	// It returns false if the iteration was stopped early.
	ForEeachVector(func(v Vector) bool) bool

	// Weight returns the weight of this vector provider, 0 if not set.
	Weight() int

	// Zero based ordinal. Lower is earlier.
	Ordinal() int
}

// Dimension is a dimension in the Hugo build matrix.
type Dimension int8

// Has returns whether the given dimension is set.
func (d Dimension) Has(o Dimension) bool {
	return d&o == o
}

// Set sets the given dimension.
func (d Dimension) Set(o Dimension) Dimension {
	return d | o
}

// Index returns this dimension's index in the Vector.
func (d Dimension) Index() int {
	if d == 0 {
		panic("dimension flag not set")
	}
	return int(d - 1)
}

func (d Dimension) Name() string {
	switch d {
	case Language:
		return "language"
	case Version:
		return "version"
	case Role:
		return "role"
	default:
		panic(fmt.Sprintf("unknown dimension flag: %d", d))
	}
}

func ParseDimension(s string) (Dimension, error) {
	switch s {
	case "language":
		return Language, nil
	case "version":
		return Version, nil
	case "role":
		return Role, nil
	default:
		return 0, fmt.Errorf("unknown dimension flag: %s", s)
	}
}

// Connon information provided by all of labguage, version and role.
type DimensionInfo interface {
	// The name. This corresponds to the key in the config, e.g. "en", "v1.2.3", "guest".
	Name() string

	// Whether this is the default value for this dimension.
	IsDefault() bool
}
