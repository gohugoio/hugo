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

// Package encoding provides template functions for encoding content.
package encoding

import (
	"encoding/base64"
	"encoding/json"
	"html/template"

	"github.com/spf13/cast"
)

// New returns a new instance of the encoding-namespaced template functions.
func New() *Namespace {
	return &Namespace{}
}

// Namespace provides template functions for the "encoding" namespace.
type Namespace struct{}

// Base64Decode returns the base64 decoding of the given content.
func (ns *Namespace) Base64Decode(content interface{}) (string, error) {
	conv, err := cast.ToStringE(content)
	if err != nil {
		return "", err
	}

	dec, err := base64.StdEncoding.DecodeString(conv)
	return string(dec), err
}

// Base64Encode returns the base64 encoding of the given content.
func (ns *Namespace) Base64Encode(content interface{}) (string, error) {
	conv, err := cast.ToStringE(content)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString([]byte(conv)), nil
}

// Jsonify encodes a given object to JSON.
func (ns *Namespace) Jsonify(v interface{}) (template.HTML, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	return template.HTML(b), nil
}
