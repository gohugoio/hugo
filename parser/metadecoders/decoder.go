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

package metadecoders

import (
	"encoding/json"
	"fmt"

	"github.com/gohugoio/hugo/common/herrors"

	"github.com/BurntSushi/toml"
	"github.com/chaseadamsio/goorgeous"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	yaml "gopkg.in/yaml.v2"
)

// UnmarshalToMap will unmarshall data in format f into a new map. This is
// what's needed for Hugo's front matter decoding.
func UnmarshalToMap(data []byte, f Format) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	if data == nil {
		return m, nil
	}

	err := unmarshal(data, f, &m)

	return m, err

}

// Unmarshal will unmarshall data in format f into an interface{}.
// This is what's needed for Hugo's /data handling.
func Unmarshal(data []byte, f Format) (interface{}, error) {
	if data == nil {
		return make(map[string]interface{}), nil
	}
	var v interface{}
	err := unmarshal(data, f, &v)

	return v, err
}

// unmarshal unmarshals data in format f into v.
func unmarshal(data []byte, f Format, v interface{}) error {

	var err error

	switch f {
	case ORG:
		vv, err := goorgeous.OrgHeaders(data)
		if err != nil {
			return toFileError(f, errors.Wrap(err, "failed to unmarshal ORG headers"))
		}
		switch v.(type) {
		case *map[string]interface{}:
			*v.(*map[string]interface{}) = vv
		default:
			*v.(*interface{}) = vv
		}
	case JSON:
		err = json.Unmarshal(data, v)
	case TOML:
		err = toml.Unmarshal(data, v)
	case YAML:
		err = yaml.Unmarshal(data, v)
		if err != nil {
			return toFileError(f, errors.Wrap(err, "failed to unmarshal YAML"))
		}

		// To support boolean keys, the YAML package unmarshals maps to
		// map[interface{}]interface{}. Here we recurse through the result
		// and change all maps to map[string]interface{} like we would've
		// gotten from `json`.
		var ptr interface{}
		switch v.(type) {
		case *map[string]interface{}:
			ptr = *v.(*map[string]interface{})
		case *interface{}:
			ptr = *v.(*interface{})
		default:
			return errors.Errorf("unknown type %T in YAML unmarshal", v)
		}

		if mm, changed := stringifyMapKeys(ptr); changed {
			switch v.(type) {
			case *map[string]interface{}:
				*v.(*map[string]interface{}) = mm.(map[string]interface{})
			case *interface{}:
				*v.(*interface{}) = mm
			}
		}
	default:
		return errors.Errorf("unmarshal of format %q is not supported", f)
	}

	if err == nil {
		return nil
	}

	return toFileError(f, errors.Wrap(err, "unmarshal failed"))

}

func toFileError(f Format, err error) error {
	return herrors.ToFileError(string(f), err)
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
	case map[string]interface{}:
		for k, v := range in {
			if vv, changed := stringifyMapKeys(v); changed {
				in[k] = vv
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
