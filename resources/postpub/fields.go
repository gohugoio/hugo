// Copyright 2020 The Hugo Authors. All rights reserved.
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

package postpub

import (
	"reflect"
)

const (
	FieldNotSupported = "__field_not_supported"
)

func structToMapWithPlaceholders(root string, in interface{}, createPlaceholder func(s string) string) map[string]interface{} {
	m := structToMap(in)
	insertFieldPlaceholders(root, m, createPlaceholder)
	return m
}

func structToMap(s interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	t := reflect.TypeOf(s)

	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		if method.PkgPath != "" {
			continue
		}
		if method.Type.NumIn() == 1 {
			m[method.Name] = ""
		}
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.PkgPath != "" {
			continue
		}
		m[field.Name] = ""
	}
	return m
}

// insert placeholder for the templates. Do it very shallow for now.
func insertFieldPlaceholders(root string, m map[string]interface{}, createPlaceholder func(s string) string) {
	for k, _ := range m {
		m[k] = createPlaceholder(root + "." + k)
	}
}
