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

package types

import (
	"html/template"

	"github.com/spf13/cast"
)

// ToStringSlicePreserveString converts v to a string slice.
// If v is a string, it will be wrapped in a string slice.
func ToStringSlicePreserveString(v interface{}) []string {
	if v == nil {
		return nil
	}
	if sds, ok := v.(string); ok {
		return []string{sds}
	}
	return cast.ToStringSlice(v)
}

// TypeToString converts v to a string if it's a valid string type.
// Note that this will not try to convert numeric values etc.,
// use ToString for that.
func TypeToString(v interface{}) (string, bool) {
	switch s := v.(type) {
	case string:
		return s, true
	case template.HTML:
		return string(s), true
	case template.CSS:
		return string(s), true
	case template.HTMLAttr:
		return string(s), true
	case template.JS:
		return string(s), true
	case template.JSStr:
		return string(s), true
	case template.URL:
		return string(s), true
	case template.Srcset:
		return string(s), true
	}

	return "", false
}

// ToString converts v to a string.
func ToString(v interface{}) string {
	if s, ok := TypeToString(v); ok {
		return s
	}

	return cast.ToString(v)

}
