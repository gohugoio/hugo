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

package collections

import (
	"fmt"
	"reflect"
)

// Append appends from to a slice to and returns the resulting slice.
// If lenght of from is one and the only element is a slice of same type as to,
// it will be appended.
func Append(to interface{}, from ...interface{}) (interface{}, error) {
	tov, toIsNil := indirect(reflect.ValueOf(to))

	toIsNil = toIsNil || to == nil
	var tot reflect.Type

	if !toIsNil {
		if tov.Kind() != reflect.Slice {
			return nil, fmt.Errorf("expected a slice, got %T", to)
		}

		tot = tov.Type().Elem()
		toIsNil = tov.Len() == 0

		if len(from) == 1 {
			fromv := reflect.ValueOf(from[0])
			if fromv.Kind() == reflect.Slice {
				if toIsNil {
					// If we get nil []string, we just return the []string
					return from[0], nil
				}

				fromt := reflect.TypeOf(from[0]).Elem()

				// If we get []string []string, we append the from slice to to
				if tot == fromt {
					return reflect.AppendSlice(tov, fromv).Interface(), nil
				}
			}
		}
	}

	if toIsNil {
		return Slice(from...), nil
	}

	for _, f := range from {
		fv := reflect.ValueOf(f)
		if tot != fv.Type() {
			return nil, fmt.Errorf("append element type mismatch: expected %v, got %v", tot, fv.Type())
		}
		tov = reflect.Append(tov, fv)
	}

	return tov.Interface(), nil
}

// indirect is borrowed from the Go stdlib: 'text/template/exec.go'
// TODO(bep) consolidate
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
