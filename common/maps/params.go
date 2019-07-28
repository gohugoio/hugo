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

// GetNestedParam gets the first match of the keyStr in the candidates given.
// It will first try the exact match and then try to find it as a nested map value,
// using the given separator, e.g. "mymap.name".
// It assumes that all the maps given have lower cased keys.
func GetNestedParam(keyStr, separator string, candidates ...map[string]interface{}) (interface{}, error) {
	keyStr = strings.ToLower(keyStr)

	lookupFn := func(key string) interface{} {
		for _, m := range candidates {
			if v, ok := m[key]; ok {
				return v
			}
		}

		return nil
	}

	v, _, _, err := GetNestedParamFn(keyStr, separator, lookupFn)
	return v, err
}

func GetNestedParamFn(keyStr, separator string, lookupFn func(key string) interface{}) (interface{}, string, map[string]interface{}, error) {
	result, _ := traverseDirectParams(keyStr, lookupFn)
	if result != nil {
		return result, keyStr, nil, nil
	}

	keySegments := strings.Split(keyStr, separator)
	if len(keySegments) == 1 {
		return nil, keyStr, nil, nil
	}

	return traverseNestedParams(keySegments, lookupFn)
}

func traverseDirectParams(keyStr string, lookupFn func(key string) interface{}) (interface{}, error) {
	return lookupFn(keyStr), nil
}

func traverseNestedParams(keySegments []string, lookupFn func(key string) interface{}) (interface{}, string, map[string]interface{}, error) {
	firstKey, rest := keySegments[0], keySegments[1:]
	result := lookupFn(firstKey)
	if result == nil || len(rest) == 0 {
		return result, firstKey, nil, nil
	}

	switch m := result.(type) {
	case map[string]interface{}:
		v, key, owner := traverseParams(rest, m)
		return v, key, owner, nil
	default:
		return nil, "", nil, nil
	}
}

func traverseParams(keys []string, m map[string]interface{}) (interface{}, string, map[string]interface{}) {
	// Shift first element off.
	firstKey, rest := keys[0], keys[1:]
	result := m[firstKey]

	// No point in continuing here.
	if result == nil {
		return result, "", nil
	}

	if len(rest) == 0 {
		// That was the last key.
		return result, firstKey, m
	}

	// That was not the last key.
	return traverseParams(rest, cast.ToStringMap(result))
}
