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

package parser

import (
	"bytes"
	"encoding/json"
	"regexp"
	"unicode"
	"unicode/utf8"
)

// Regexp definitions
var keyMatchRegex = regexp.MustCompile(`\"(\w+)\":`)
var wordBarrierRegex = regexp.MustCompile(`(\w)([A-Z])`)

// Code adapted from https://gist.github.com/piersy/b9934790a8892db1a603820c0c23e4a7
type LowerCaseCamelJSONMarshaller struct {
	Value interface{}
}

func (c LowerCaseCamelJSONMarshaller) MarshalJSON() ([]byte, error) {
	marshalled, err := json.Marshal(c.Value)

	converted := keyMatchRegex.ReplaceAllFunc(
		marshalled,
		func(match []byte) []byte {

			// Attributes on the form XML, JSON etc.
			if bytes.Equal(match, bytes.ToUpper(match)) {
				return bytes.ToLower(match)
			}

			// Empty keys are valid JSON, only lowercase if we do not have an
			// empty key.
			if len(match) > 2 {
				// Decode first rune after the double quotes
				r, width := utf8.DecodeRune(match[1:])
				r = unicode.ToLower(r)
				utf8.EncodeRune(match[1:width+1], r)
			}
			return match
		},
	)

	return converted, err
}
