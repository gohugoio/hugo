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
	"github.com/spf13/cast"
)

// GetString tries to get a value with key from map m and convert it to a string.
// It will return an empty string if not found or if it cannot be convertd to a string.
func GetString(m map[string]interface{}, key string) string {
	if m == nil {
		return ""
	}
	v, found := m[key]
	if !found {
		return ""
	}
	return cast.ToString(v)
}
