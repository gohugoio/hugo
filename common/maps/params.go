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

package maps

import (
	"strings"

	"github.com/spf13/cast"
)

// Params is a map where all keys are lower case.
type Params map[string]interface{}

// Get does a lower case and nested search in this map.
// It will return nil if none found.
func (p Params) Get(indices ...string) interface{} {
	v, _, _ := getNested(p, indices)
	return v
}

func getNested(m map[string]interface{}, indices []string) (interface{}, string, map[string]interface{}) {
	if len(indices) == 0 {
		return nil, "", nil
	}

	first := indices[0]
	v, found := m[strings.ToLower(cast.ToString(first))]
	if !found {
		return nil, "", nil
	}

	if len(indices) == 1 {
		return v, first, m
	}

	switch m2 := v.(type) {
	case Params:
		return getNested(m2, indices[1:])
	case map[string]interface{}:
		return getNested(m2, indices[1:])
	default:
		return nil, "", nil
	}
}

// GetNestedParam gets the first match of the keyStr in the candidates given.
// It will first try the exact match and then try to find it as a nested map value,
// using the given separator, e.g. "mymap.name".
// It assumes that all the maps given have lower cased keys.
func GetNestedParam(keyStr, separator string, candidates ...Params) (interface{}, error) {
	keyStr = strings.ToLower(keyStr)

	// Try exact match first
	for _, m := range candidates {
		if v, ok := m[keyStr]; ok {
			return v, nil
		}
	}

	keySegments := strings.Split(keyStr, separator)
	for _, m := range candidates {
		if v := m.Get(keySegments...); v != nil {
			return v, nil
		}
	}

	return nil, nil

}

func GetNestedParamFn(keyStr, separator string, lookupFn func(key string) interface{}) (interface{}, string, map[string]interface{}, error) {
	keySegments := strings.Split(strings.ToLower(keyStr), separator)
	if len(keySegments) == 0 {
		return nil, "", nil, nil
	}

	first := lookupFn(keySegments[0])
	if first == nil {
		return nil, "", nil, nil
	}

	if len(keySegments) == 1 {
		return first, keySegments[0], nil, nil
	}

	switch m := first.(type) {
	case map[string]interface{}:
		v, key, owner := getNested(m, keySegments[1:])
		return v, key, owner, nil
	case Params:
		v, key, owner := getNested(m, keySegments[1:])
		return v, key, owner, nil
	}

	return nil, "", nil, nil
}
