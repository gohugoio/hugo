// Copyright 2018 The Hugo Authors. All rights reserved.
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

// The metadecoders package contains functions to decode metadata (e.g. page front matter)
// from different formats: TOML, YAML, JSON.
package metadecoders

import (
	"fmt"

	"github.com/spf13/cast"
	yaml "gopkg.in/yaml.v2"
)

// HandleYAMLData unmarshals YAML-encoded datum and returns a Go interface
// representing the encoded data structure.
func HandleYAMLData(datum []byte) (interface{}, error) {
	var m interface{}
	err := yaml.Unmarshal(datum, &m)
	if err != nil {
		return nil, err
	}

	// To support boolean keys, the `yaml` package unmarshals maps to
	// map[interface{}]interface{}. Here we recurse through the result
	// and change all maps to map[string]interface{} like we would've
	// gotten from `json`.
	if mm, changed := stringifyMapKeys(m); changed {
		return mm, nil
	}

	return m, nil
}

// stringifyMapKeys recurses into in and changes all instances of
// map[interface{}]interface{} to map[string]interface{}. This is useful to
// work around the impedence mismatch between JSON and YAML unmarshaling that's
// described here: https://github.com/go-yaml/yaml/issues/139
//
// Inspired by https://github.com/stripe/stripe-mock, MIT licensed
func stringifyMapKeys(in interface{}) (interface{}, bool) {
	switch in := in.(type) {
	case []interface{}:
		for i, v := range in {
			if vv, replaced := stringifyMapKeys(v); replaced {
				in[i] = vv
			}
		}
	case map[interface{}]interface{}:
		res := make(map[string]interface{})
		var (
			ok  bool
			err error
		)
		for k, v := range in {
			var ks string

			if ks, ok = k.(string); !ok {
				ks, err = cast.ToStringE(k)
				if err != nil {
					ks = fmt.Sprintf("%v", k)
				}
			}
			if vv, replaced := stringifyMapKeys(v); replaced {
				res[ks] = vv
			} else {
				res[ks] = v
			}
		}
		return res, true
	}

	return nil, false
}
