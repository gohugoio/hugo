// Copyright 2017 The Hugo Authors. All rights reserved.
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

package tplimpl

import (
	"reflect"
	"time"
)

// toInt returns the int value if possible, -1 if not.
func toInt(v reflect.Value) int64 {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int()
	case reflect.Interface:
		return toInt(v.Elem())
	}
	return -1
}

// toString returns the string value if possible, "" if not.
func toString(v reflect.Value) string {
	switch v.Kind() {
	case reflect.String:
		return v.String()
	case reflect.Interface:
		return toString(v.Elem())
	}
	return ""
}

var (
	zero      reflect.Value
	errorType = reflect.TypeOf((*error)(nil)).Elem()
	timeType  = reflect.TypeOf((*time.Time)(nil)).Elem()
)

func toTimeUnix(v reflect.Value) int64 {
	if v.Kind() == reflect.Interface {
		return toTimeUnix(v.Elem())
	}
	if v.Type() != timeType {
		panic("coding error: argument must be time.Time type reflect Value")
	}
	return v.MethodByName("Unix").Call([]reflect.Value{})[0].Int()
}

// indirect is taken from 'text/template/exec.go'
func indirect(v reflect.Value) (rv reflect.Value, isNil bool) {
	for ; v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface; v = v.Elem() {
		if v.IsNil() {
			return v, true
		}
		if v.Kind() == reflect.Interface && v.NumMethod() > 0 {
			break
		}
	}
	return v, false
}
