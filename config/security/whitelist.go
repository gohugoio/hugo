// Copyright 2021 The Hugo Authors. All rights reserved.
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

package security

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

const (
	acceptNoneKeyword = "none"
)

// Whitelist holds a whitelist.
type Whitelist struct {
	acceptNone bool
	patterns   []*regexp.Regexp

	// Store this for debugging/error reporting
	patternsStrings []string
}

// MarshalJSON is for internal use only.
func (w Whitelist) MarshalJSON() ([]byte, error) {
	if w.acceptNone {
		return json.Marshal(acceptNoneKeyword)
	}

	return json.Marshal(w.patternsStrings)
}

// NewWhitelist creates a new Whitelist from zero or more patterns.
// An empty patterns list or a pattern with the value 'none' will create
// a whitelist that will Accept none.
func NewWhitelist(patterns ...string) (Whitelist, error) {
	if len(patterns) == 0 {
		return Whitelist{acceptNone: true}, nil
	}

	var acceptSome bool
	var patternsStrings []string

	for _, p := range patterns {
		if p == acceptNoneKeyword {
			acceptSome = false
			break
		}

		if ps := strings.TrimSpace(p); ps != "" {
			acceptSome = true
			patternsStrings = append(patternsStrings, ps)
		}
	}

	if !acceptSome {
		return Whitelist{
			acceptNone: true,
		}, nil
	}

	var patternsr []*regexp.Regexp

	for i := 0; i < len(patterns); i++ {
		p := strings.TrimSpace(patterns[i])
		if p == "" {
			continue
		}
		re, err := regexp.Compile(p)
		if err != nil {
			return Whitelist{}, fmt.Errorf("failed to compile whitelist pattern %q: %w", p, err)
		}
		patternsr = append(patternsr, re)
	}

	return Whitelist{patterns: patternsr, patternsStrings: patternsStrings}, nil
}

// MustNewWhitelist creates a new Whitelist from zero or more patterns and panics on error.
func MustNewWhitelist(patterns ...string) Whitelist {
	w, err := NewWhitelist(patterns...)
	if err != nil {
		panic(err)
	}
	return w
}

// Accept reports whether name is whitelisted.
func (w Whitelist) Accept(name string) bool {
	if w.acceptNone {
		return false
	}

	for _, p := range w.patterns {
		if p.MatchString(name) {
			return true
		}
	}
	return false
}

func (w Whitelist) String() string {
	return fmt.Sprint(w.patternsStrings)
}
