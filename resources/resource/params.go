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

package resource

import (
	"strings"

	"github.com/spf13/cast"
)

func Param(r ResourceParamsProvider, fallback map[string]interface{}, key interface{}) (interface{}, error) {
	keyStr, err := cast.ToStringE(key)
	if err != nil {
		return nil, err
	}

	keyStr = strings.ToLower(keyStr)
	result, _ := traverseDirectParams(r, fallback, keyStr)
	if result != nil {
		return result, nil
	}

	keySegments := strings.Split(keyStr, ".")
	if len(keySegments) == 1 {
		return nil, nil
	}

	return traverseNestedParams(r, fallback, keySegments)
}

func traverseDirectParams(r ResourceParamsProvider, fallback map[string]interface{}, key string) (interface{}, error) {
	keyStr := strings.ToLower(key)
	if val, ok := r.Params()[keyStr]; ok {
		return val, nil
	}

	if fallback == nil {
		return nil, nil
	}

	return fallback[keyStr], nil
}

func traverseNestedParams(r ResourceParamsProvider, fallback map[string]interface{}, keySegments []string) (interface{}, error) {
	result := traverseParams(keySegments, r.Params())
	if result != nil {
		return result, nil
	}

	if fallback != nil {
		result = traverseParams(keySegments, fallback)
		if result != nil {
			return result, nil
		}
	}

	// Didn't find anything, but also no problems.
	return nil, nil
}

func traverseParams(keys []string, m map[string]interface{}) interface{} {
	// Shift first element off.
	firstKey, rest := keys[0], keys[1:]
	result := m[firstKey]

	// No point in continuing here.
	if result == nil {
		return result
	}

	if len(rest) == 0 {
		// That was the last key.
		return result
	}

	// That was not the last key.
	return traverseParams(rest, cast.ToStringMap(result))
}
