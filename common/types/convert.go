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
	"encoding/json"
	"fmt"
	"html/template"
	"reflect"
	"time"

	"github.com/spf13/cast"
)

// ToDuration converts v to time.Duration.
// See ToDurationE if you need to handle errors.
func ToDuration(v any) time.Duration {
	d, _ := ToDurationE(v)
	return d
}

// ToDurationE converts v to time.Duration.
func ToDurationE(v any) (time.Duration, error) {
	if n := cast.ToInt(v); n > 0 {
		return time.Duration(n) * time.Millisecond, nil
	}
	d, err := time.ParseDuration(cast.ToString(v))
	if err != nil {
		return 0, fmt.Errorf("cannot convert %v to time.Duration", v)
	}
	return d, nil
}

// ToStringSlicePreserveString is the same as ToStringSlicePreserveStringE,
// but it never fails.
func ToStringSlicePreserveString(v any) []string {
	vv, _ := ToStringSlicePreserveStringE(v)
	return vv
}

// ToStringSlicePreserveStringE converts v to a string slice.
// If v is a string, it will be wrapped in a string slice.
func ToStringSlicePreserveStringE(v any) ([]string, error) {
	if v == nil {
		return nil, nil
	}
	if sds, ok := v.(string); ok {
		return []string{sds}, nil
	}
	result, err := cast.ToStringSliceE(v)
	if err == nil {
		return result, nil
	}

	// Probably []int or similar. Fall back to reflect.
	vv := reflect.ValueOf(v)

	switch vv.Kind() {
	case reflect.Slice, reflect.Array:
		result = make([]string, vv.Len())
		for i := 0; i < vv.Len(); i++ {
			s, err := cast.ToStringE(vv.Index(i).Interface())
			if err != nil {
				return nil, err
			}
			result[i] = s
		}
		return result, nil
	default:
		return nil, fmt.Errorf("failed to convert %T to a string slice", v)
	}
}

// TypeToString converts v to a string if it's a valid string type.
// Note that this will not try to convert numeric values etc.,
// use ToString for that.
func TypeToString(v any) (string, bool) {
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
func ToString(v any) string {
	s, _ := ToStringE(v)
	return s
}

// ToStringE converts v to a string.
func ToStringE(v any) (string, error) {
	if s, ok := TypeToString(v); ok {
		return s, nil
	}

	switch s := v.(type) {
	case json.RawMessage:
		return string(s), nil
	default:
		return cast.ToStringE(v)
	}
}
