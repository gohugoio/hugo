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

	"github.com/gohugoio/hugo/common/hreflect"
)

// Regexp definitions
var (
	keyMatchRegex       = regexp.MustCompile(`\"(\w+)\":`)
	nullEnableBoolRegex = regexp.MustCompile(`\"(enable\w+)\":null`)
)

type NullBoolJSONMarshaller struct {
	Wrapped json.Marshaler
}

func (c NullBoolJSONMarshaller) MarshalJSON() ([]byte, error) {
	b, err := c.Wrapped.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return nullEnableBoolRegex.ReplaceAll(b, []byte(`"$1": false`)), nil
}

// Code adapted from https://gist.github.com/piersy/b9934790a8892db1a603820c0c23e4a7
type LowerCaseCamelJSONMarshaller struct {
	Value any
}

var preserveUpperCaseKeyRe = regexp.MustCompile(`^"HTTP`)

func preserveUpperCaseKey(match []byte) bool {
	return preserveUpperCaseKeyRe.Match(match)
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
			if len(match) > 2 && !preserveUpperCaseKey(match) {
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

type ReplacingJSONMarshaller struct {
	Value any

	KeysToLower bool
	OmitEmpty   bool
}

func (c ReplacingJSONMarshaller) MarshalJSON() ([]byte, error) {
	converted, err := json.Marshal(c.Value)

	if c.KeysToLower {
		converted = keyMatchRegex.ReplaceAllFunc(
			converted,
			func(match []byte) []byte {
				return bytes.ToLower(match)
			},
		)
	}

	if c.OmitEmpty {
		// It's tricky to do this with a regexp, so convert it to a map, remove zero values and convert back.
		var m map[string]interface{}
		err = json.Unmarshal(converted, &m)
		if err != nil {
			return nil, err
		}
		var removeZeroVAlues func(m map[string]any)
		removeZeroVAlues = func(m map[string]any) {
			for k, v := range m {
				if !hreflect.IsTruthful(v) {
					delete(m, k)
				} else {
					switch vv := v.(type) {
					case map[string]interface{}:
						removeZeroVAlues(vv)
					case []interface{}:
						for _, vvv := range vv {
							if m, ok := vvv.(map[string]any); ok {
								removeZeroVAlues(m)
							}
						}
					}
				}
			}
		}
		removeZeroVAlues(m)
		converted, err = json.Marshal(m)

	}

	return converted, err
}
