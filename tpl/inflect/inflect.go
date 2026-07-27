// Copyright 2017 The Hugo Authors. All rights reserved.
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

// Package inflect provides template functions for the inflection of words.
package inflect

import (
	"regexp"
	"strconv"
	"strings"

	_inflect "github.com/gobuffalo/flect"
	"github.com/spf13/cast"
)

// New returns a new instance of the inflect-namespaced template functions.
func New() *Namespace {
	return &Namespace{}
}

// Namespace provides template functions for the "inflect" namespace.
type Namespace struct{}

// quoteSpaceRe matches a spurious space that flect.Humanize inserts after an
// opening quotation mark when the quoted text begins with an uppercase letter.
// flect incorrectly treats quote+uppercase as a camelCase word boundary.
var quoteSpaceRe = regexp.MustCompile("([\u201e\u201a\u201c\u2018\u00ab\u2039\"']) ")

// Humanize returns the humanized form of v.
//
// If v is either an integer or a string containing an integer
// value, the behavior is to add the appropriate ordinal.
func (ns *Namespace) Humanize(v any) (string, error) {
	word, err := cast.ToStringE(v)
	if err != nil {
		return "", err
	}

	if word == "" {
		return "", nil
	}

	_, ok := v.(int)            // original param was literal int value
	_, err = strconv.Atoi(word) // original param was string containing an int value
	if ok || err == nil {
		return _inflect.Ordinalize(word), nil
	}

	str := _inflect.Humanize(word)
	result := _inflect.Humanize(strings.ToLower(str))
	return quoteSpaceRe.ReplaceAllString(result, "$1"), nil
}

// Pluralize returns the plural form of the single word in v.
func (ns *Namespace) Pluralize(v any) (string, error) {
	word, err := cast.ToStringE(v)
	if err != nil {
		return "", err
	}

	return _inflect.Pluralize(word), nil
}

// Singularize returns the singular form of a single word in v.
func (ns *Namespace) Singularize(v any) (string, error) {
	word, err := cast.ToStringE(v)
	if err != nil {
		return "", err
	}

	return _inflect.Singularize(word), nil
}
