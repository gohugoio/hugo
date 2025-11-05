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

package predicate

import (
	"iter"
	"strings"

	"github.com/gobwas/glob"
	"github.com/gohugoio/hugo/hugofs/hglob"
)

// Match represents the result of a predicate evaluation.
type Match interface {
	OK() bool
}

var (
	// Predefined Match values for common cases.
	True  = BoolMatch(true)
	False = BoolMatch(false)
)

// BoolMatch is a simple Match implementation based on a boolean value.
type BoolMatch bool

func (b BoolMatch) OK() bool {
	return bool(b)
}

// breakMatch is a Match implementation that always returns false for OK() and signals to break evaluation.
type breakMatch struct{}

func (b breakMatch) OK() bool {
	return false
}

var matchBreak = breakMatch{}

// P is a predicate function that tests whether a value of type T satisfies some condition.
type P[T any] func(T) bool

// Or returns a predicate that is a short-circuiting logical OR of this and the given predicates.
// Note that P[T] only supports Or. For chained AND/OR logic, use PR[T].
func (p P[T]) Or(ps ...P[T]) P[T] {
	return func(v T) bool {
		if p != nil && p(v) {
			return true
		}
		for _, pp := range ps {
			if pp(v) {
				return true
			}
		}
		return false
	}
}

// PR is a predicate function that tests whether a value of type T satisfies some condition and returns a Match result.
type PR[T any] func(T) Match

// BoolFunc returns a P[T] version of this predicate.
func (p PR[T]) BoolFunc() P[T] {
	return func(v T) bool {
		if p == nil {
			return false
		}
		return p(v).OK()
	}
}

// And returns a predicate that is a short-circuiting logical AND of this and the given predicates.
func (p PR[T]) And(ps ...PR[T]) PR[T] {
	return func(v T) Match {
		if p != nil {
			m := p(v)
			if !m.OK() || shouldBreak(m) {
				return matchBreak
			}
		}
		for _, pp := range ps {
			m := pp(v)
			if !m.OK() || shouldBreak(m) {
				return matchBreak
			}
		}
		return BoolMatch(true)
	}
}

// Or returns a predicate that is a short-circuiting logical OR of this and the given predicates.
func (p PR[T]) Or(ps ...PR[T]) PR[T] {
	return func(v T) Match {
		if p != nil {
			m := p(v)
			if m.OK() {
				return m
			}
			if shouldBreak(m) {
				return matchBreak
			}
		}
		for _, pp := range ps {
			m := pp(v)
			if m.OK() {
				return m
			}
			if shouldBreak(m) {
				return matchBreak
			}
		}
		return BoolMatch(false)
	}
}

func shouldBreak(m Match) bool {
	_, ok := m.(breakMatch)
	return ok
}

// Filter returns a new slice holding only the elements of s that satisfy p.
// Filter modifies the contents of the slice s and returns the modified slice, which may have a smaller length.
func (p PR[T]) Filter(s []T) []T {
	var n int
	for _, v := range s {
		if p(v).OK() {
			s[n] = v
			n++
		}
	}
	return s[:n]
}

// FilterCopy returns a new slice holding only the elements of s that satisfy p.
func (p PR[T]) FilterCopy(s []T) []T {
	var result []T
	for _, v := range s {
		if p(v).OK() {
			result = append(result, v)
		}
	}
	return result
}

// NewStringPredicateFromGlobs creates a string predicate from the given glob patterns.
// A glob pattern starting with "!" is a negation pattern which will be ANDed with the rest.
func NewStringPredicateFromGlobs(patterns []string, getGlob func(pattern string) (glob.Glob, error)) (P[string], error) {
	var p PR[string]
	for _, pattern := range patterns {
		pattern = strings.TrimSpace(pattern)
		if pattern == "" {
			continue
		}
		negate := strings.HasPrefix(pattern, hglob.NegationPrefix)
		if negate {
			pattern = pattern[2:]
			g, err := getGlob(pattern)
			if err != nil {
				return nil, err
			}
			p = p.And(func(s string) Match {
				return BoolMatch(!g.Match(s))
			})
		} else {
			g, err := getGlob(pattern)
			if err != nil {
				return nil, err
			}
			p = p.Or(func(s string) Match {
				return BoolMatch(g.Match(s))
			})

		}
	}

	return p.BoolFunc(), nil
}

type IndexMatcher interface {
	IndexMatch(match P[string]) (iter.Seq[int], error)
}
