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
	"time"

	"github.com/gohugoio/hugo/helpers"
	"github.com/pelletier/go-toml/v2"

	"github.com/spf13/cast"
)

// GetParam will return the param with the given key from the Resource,
// nil if not found.
func GetParam(r Resource, key string) any {
	return getParam(r, key, false)
}

// GetParamToLower is the same as GetParam but it will lower case any string
// result, including string slices.
func GetParamToLower(r Resource, key string) any {
	return getParam(r, key, true)
}

func getParam(r Resource, key string, stringToLower bool) any {
	v := r.Params()[strings.ToLower(key)]

	if v == nil {
		return nil
	}

	switch val := v.(type) {
	case bool:
		return val
	case string:
		if stringToLower {
			return strings.ToLower(val)
		}
		return val
	case int64, int32, int16, int8, int:
		return cast.ToInt(v)
	case float64, float32:
		return cast.ToFloat64(v)
	case time.Time:
		return val
	case toml.LocalDate:
		return val.AsTime(time.UTC)
	case toml.LocalDateTime:
		return val.AsTime(time.UTC)
	case []string:
		if stringToLower {
			return helpers.SliceToLower(val)
		}
		return v
	case map[string]any:
		return v
	case map[any]any:
		return v
	}
	return nil
}
