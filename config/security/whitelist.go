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

	"github.com/gohugoio/hugo/hugofs/hglob"
)

const acceptNoneKeyword = "none"

// Whitelist holds a whitelist.
//
// Patterns are regular expressions. A pattern prefixed with "! "
// (see hglob.NegationPrefix) is a deny rule: a name that matches any
// deny rule is rejected even if it matches an allow rule.
// A whitelist made up exclusively of deny rules implicitly allows
// names that do not match any of them.
type Whitelist struct {
	acceptNone bool
	allow      []*regexp.Regexp
	deny       []*regexp.Regexp

	// Store this for debugging/error reporting.
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
// a whitelist that will Accept none. Patterns prefixed with "! " act as
// deny rules; see Whitelist.
func NewWhitelist(patterns ...string) (Whitelist, error) {
	if len(patterns) == 0 {
		return Whitelist{acceptNone: true}, nil
	}

	var (
		acceptSome      bool
		patternsStrings []string
	)

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
		return Whitelist{acceptNone: true}, nil
	}

	var allow, deny []*regexp.Regexp
	for _, p := range patternsStrings {
		raw := p
		negate := strings.HasPrefix(p, hglob.NegationPrefix)
		if negate {
			raw = p[len(hglob.NegationPrefix):]
		}
		re, err := regexp.Compile(raw)
		if err != nil {
			return Whitelist{}, fmt.Errorf("failed to compile whitelist pattern %q: %w", p, err)
		}
		if negate {
			deny = append(deny, re)
		} else {
			allow = append(allow, re)
		}
	}

	return Whitelist{allow: allow, deny: deny, patternsStrings: patternsStrings}, nil
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

	for _, p := range w.deny {
		if p.MatchString(name) {
			return false
		}
	}

	if len(w.allow) == 0 {
		// A whitelist with only deny rules implicitly allows everything
		// that is not denied. An empty (zero-value) whitelist rejects.
		return len(w.deny) > 0
	}

	for _, p := range w.allow {
		if p.MatchString(name) {
			return true
		}
	}
	return false
}

func (w Whitelist) String() string {
	return fmt.Sprint(w.patternsStrings)
}
