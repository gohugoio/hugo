// Copyright 2020 The Hugo Authors. All rights reserved.
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
	"errors"
	"html/template"

	bp "github.com/gohugoio/hugo/bufferpool"

	"github.com/gohugoio/hugo/common/maps"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
)

// New returns a new instance of the encoding-namespaced template functions.
func New() *Namespace {
	return &Namespace{}
}

// Namespace provides template functions for the "encoding" namespace.
type Namespace struct{}

// Base64Decode returns the base64 decoding of the given content.
func (ns *Namespace) Base64Decode(content any) (string, error) {
	conv, err := cast.ToStringE(content)
	if err != nil {
		return "", err
	}

	dec, err := base64.StdEncoding.DecodeString(conv)
	return string(dec), err
}

// Base64Encode returns the base64 encoding of the given content.
func (ns *Namespace) Base64Encode(content any) (string, error) {
	conv, err := cast.ToStringE(content)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString([]byte(conv)), nil
}

// Jsonify encodes a given object to JSON.  To pretty print the JSON, pass a map
// or dictionary of options as the first value in args.  Supported options are
// "prefix" and "indent".  Each JSON element in the output will begin on a new
// line beginning with prefix followed by one or more copies of indent according
// to the indentation nesting.
func (ns *Namespace) Jsonify(args ...any) (template.HTML, error) {
	var (
		b    []byte
		err  error
		obj  any
		opts jsonifyOpts
	)

	switch len(args) {
	case 0:
		return "", nil
	case 1:
		obj = args[0]
	case 2:
		var m map[string]any
		m, err = maps.ToStringMapE(args[0])
		if err != nil {
			break
		}
		if err = mapstructure.WeakDecode(m, &opts); err != nil {
			break
		}
		obj = args[1]
	default:
		err = errors.New("too many arguments to jsonify")
	}

	if err != nil {
		return "", err
	}

	buff := bp.GetBuffer()
	defer bp.PutBuffer(buff)
	e := json.NewEncoder(buff)
	e.SetEscapeHTML(!opts.NoHTMLEscape)
	e.SetIndent(opts.Prefix, opts.Indent)
	if err = e.Encode(obj); err != nil {
		return "", err
	}
	b = buff.Bytes()
	// See https://github.com/golang/go/issues/37083
	// Hugo changed from MarshalIndent/Marshal. To make the output
	// the same, we need to trim the trailing newline.
	b = b[:len(b)-1]

	return template.HTML(b), nil
}

type jsonifyOpts struct {
	Prefix       string
	Indent       string
	NoHTMLEscape bool
}
