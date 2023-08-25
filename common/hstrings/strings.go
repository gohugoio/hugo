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
	"regexp"
	"strings"
	"sync"

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

// regexpCache represents a cache of regexp objects protected by a mutex.
type regexpCache struct {
	mu sync.RWMutex
	re map[string]*regexp.Regexp
}

func (rc *regexpCache) getOrCompileRegexp(pattern string) (re *regexp.Regexp, err error) {
	var ok bool

	if re, ok = rc.get(pattern); !ok {
		re, err = regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}
		rc.set(pattern, re)
	}

	return re, nil
}

func (rc *regexpCache) get(key string) (re *regexp.Regexp, ok bool) {
	rc.mu.RLock()
	re, ok = rc.re[key]
	rc.mu.RUnlock()
	return
}

func (rc *regexpCache) set(key string, re *regexp.Regexp) {
	rc.mu.Lock()
	rc.re[key] = re
	rc.mu.Unlock()
}

var reCache = regexpCache{re: make(map[string]*regexp.Regexp)}

// GetOrCompileRegexp retrieves a regexp object from the cache based upon the pattern.
// If the pattern is not found in the cache, the pattern is compiled and added to
// the cache.
func GetOrCompileRegexp(pattern string) (re *regexp.Regexp, err error) {
	return reCache.getOrCompileRegexp(pattern)
}
