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
	"strings"

	"github.com/spf13/cast"
)

// ToLower makes all the keys in the given map lower cased and will do so
// recursively.
// Notes:
// * This will modify the map given.
// * Any nested map[interface{}]interface{} will be converted to map[string]interface{}.
func ToLower(m map[string]interface{}) {
	for k, v := range m {
		switch v.(type) {
		case map[interface{}]interface{}:
			v = cast.ToStringMap(v)
			ToLower(v.(map[string]interface{}))
		case map[string]interface{}:
			ToLower(v.(map[string]interface{}))
		}

		lKey := strings.ToLower(k)
		if k != lKey {
			delete(m, k)
			m[lKey] = v
		}

	}
}
