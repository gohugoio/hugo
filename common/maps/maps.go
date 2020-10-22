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

package maps

import (
	"fmt"
	"strings"

	"github.com/gobwas/glob"

	"github.com/spf13/cast"
)

// ToLower makes all the keys in the given map lower cased and will do so
// recursively.
// Notes:
// * This will modify the map given.
// * Any nested map[interface{}]interface{} will be converted to Params.
func ToLower(m Params) {
	for k, v := range m {
		var retyped bool
		switch v.(type) {
		case map[interface{}]interface{}:
			var p Params = cast.ToStringMap(v)
			v = p
			ToLower(p)
			retyped = true
		case map[string]interface{}:
			var p Params = v.(map[string]interface{})
			v = p
			ToLower(p)
			retyped = true
		}

		lKey := strings.ToLower(k)
		if retyped || k != lKey {
			delete(m, k)
			m[lKey] = v
		}
	}
}

func ToStringMapE(in interface{}) (map[string]interface{}, error) {
	switch in.(type) {
	case Params:
		return in.(Params), nil
	default:
		return cast.ToStringMapE(in)
	}
}

func ToStringMap(in interface{}) map[string]interface{} {
	m, _ := ToStringMapE(in)
	return m
}

func ToSliceStringMap(in interface{}) ([]map[string]interface{}, error) {
	switch v := in.(type) {
	case []map[string]interface{}:
		return v, nil
	case []interface{}:
		var s []map[string]interface{}
		for _, entry := range v {
			if vv, ok := entry.(map[string]interface{}); ok {
				s = append(s, vv)
			}
		}
		return s, nil
	default:
		return nil, fmt.Errorf("unable to cast %#v of type %T to []map[string]interface{}", in, in)
	}
}

type keyRename struct {
	pattern glob.Glob
	newKey  string
}

// KeyRenamer supports renaming of keys in a map.
type KeyRenamer struct {
	renames []keyRename
}

// NewKeyRenamer creates a new KeyRenamer given a list of pattern and new key
// value pairs.
func NewKeyRenamer(patternKeys ...string) (KeyRenamer, error) {
	var renames []keyRename
	for i := 0; i < len(patternKeys); i += 2 {
		g, err := glob.Compile(strings.ToLower(patternKeys[i]), '/')
		if err != nil {
			return KeyRenamer{}, err
		}
		renames = append(renames, keyRename{pattern: g, newKey: patternKeys[i+1]})
	}

	return KeyRenamer{renames: renames}, nil
}

func (r KeyRenamer) getNewKey(keyPath string) string {
	for _, matcher := range r.renames {
		if matcher.pattern.Match(keyPath) {
			return matcher.newKey
		}
	}

	return ""
}

// Rename renames the keys in the given map according
// to the patterns in the current KeyRenamer.
func (r KeyRenamer) Rename(m map[string]interface{}) {
	r.renamePath("", m)
}

func (KeyRenamer) keyPath(k1, k2 string) string {
	k1, k2 = strings.ToLower(k1), strings.ToLower(k2)
	if k1 == "" {
		return k2
	} else {
		return k1 + "/" + k2
	}
}

func (r KeyRenamer) renamePath(parentKeyPath string, m map[string]interface{}) {
	for key, val := range m {
		keyPath := r.keyPath(parentKeyPath, key)
		switch val.(type) {
		case map[interface{}]interface{}:
			val = cast.ToStringMap(val)
			r.renamePath(keyPath, val.(map[string]interface{}))
		case map[string]interface{}:
			r.renamePath(keyPath, val.(map[string]interface{}))
		}

		newKey := r.getNewKey(keyPath)

		if newKey != "" {
			delete(m, key)
			m[newKey] = val
		}
	}
}
