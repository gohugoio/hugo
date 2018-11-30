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
	"strconv"

	_inflect "github.com/markbates/inflect"
	"github.com/spf13/cast"
)

// New returns a new instance of the inflect-namespaced template functions.
func New() *Namespace {
	return &Namespace{}
}

// Namespace provides template functions for the "inflect" namespace.
type Namespace struct{}

// Humanize returns the humanized form of a single parameter.
//
// If the parameter is either an integer or a string containing an integer
// value, the behavior is to add the appropriate ordinal.
//
//     Example:  "my-first-post" -> "My first post"
//     Example:  "103" -> "103rd"
//     Example:  52 -> "52nd"
func (ns *Namespace) Humanize(in interface{}) (string, error) {
	word, err := cast.ToStringE(in)
	if err != nil {
		return "", err
	}

	if word == "" {
		return "", nil
	}

	_, ok := in.(int)           // original param was literal int value
	_, err = strconv.Atoi(word) // original param was string containing an int value
	if ok || err == nil {
		return _inflect.Ordinalize(word), nil
	}

	return _inflect.Humanize(word), nil
}

// Pluralize returns the plural form of a single word.
func (ns *Namespace) Pluralize(in interface{}) (string, error) {
	word, err := cast.ToStringE(in)
	if err != nil {
		return "", err
	}

	return _inflect.Pluralize(word), nil
}

// Singularize returns the singular form of a single word.
func (ns *Namespace) Singularize(in interface{}) (string, error) {
	word, err := cast.ToStringE(in)
	if err != nil {
		return "", err
	}

	return _inflect.Singularize(word), nil
}
