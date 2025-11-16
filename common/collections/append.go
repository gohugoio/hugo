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

package collections

import (
	"fmt"
	"reflect"

	"github.com/gohugoio/hugo/common/hreflect"
)

// Append appends from to a slice to and returns the resulting slice.
// If length of from is one and the only element is a slice of same type as to,
// it will be appended.
func Append(to any, from ...any) (any, error) {
	if len(from) == 0 {
		return to, nil
	}
	tov, toIsNil := hreflect.Indirect(reflect.ValueOf(to))

	toIsNil = toIsNil || to == nil
	var tot reflect.Type

	if !toIsNil {
		if tov.Kind() == reflect.Slice {
			// Create a copy of tov, so we don't modify the original.
			c := reflect.MakeSlice(tov.Type(), tov.Len(), tov.Len()+len(from))
			reflect.Copy(c, tov)
			tov = c
		}

		if tov.Kind() != reflect.Slice {
			return nil, fmt.Errorf("expected a slice, got %T", to)
		}

		tot = tov.Type().Elem()
		if tot.Kind() == reflect.Slice {
			totvt := tot.Elem()
			fromvs := make([]reflect.Value, len(from))
			for i, f := range from {
				fromv := reflect.ValueOf(f)
				fromt := fromv.Type()
				if fromt.Kind() == reflect.Slice {
					fromt = fromt.Elem()
				}
				if totvt != fromt {
					return nil, fmt.Errorf("cannot append slice of %s to slice of %s", fromt, totvt)
				} else {
					fromvs[i] = fromv
				}
			}
			return reflect.Append(tov, fromvs...).Interface(), nil

		}

		toIsNil = tov.Len() == 0

		if len(from) == 1 {
			fromv := reflect.ValueOf(from[0])
			if !fromv.IsValid() {
				// from[0] is nil
				return appendToInterfaceSliceFromValues(tov, fromv)
			}
			fromt := fromv.Type()
			if fromt.Kind() == reflect.Slice {
				fromt = fromt.Elem()
			}
			if fromv.Kind() == reflect.Slice {
				if toIsNil {
					// If we get nil []string, we just return the []string
					return from[0], nil
				}

				// If we get []string []string, we append the from slice to to
				if tot == fromt {
					return reflect.AppendSlice(tov, fromv).Interface(), nil
				} else if !fromt.AssignableTo(tot) {
					// Fall back to a []interface{} slice.
					return appendToInterfaceSliceFromValues(tov, fromv)
				}

			}
		}
	}

	if toIsNil {
		return Slice(from...), nil
	}

	for _, f := range from {
		fv := reflect.ValueOf(f)
		if !fv.IsValid() || !fv.Type().AssignableTo(tot) {
			// Fall back to a []interface{} slice.
			tov, _ := hreflect.Indirect(reflect.ValueOf(to))
			return appendToInterfaceSlice(tov, from...)
		}
		tov = reflect.Append(tov, fv)
	}

	return tov.Interface(), nil
}

func appendToInterfaceSliceFromValues(slice1, slice2 reflect.Value) ([]any, error) {
	var tos []any

	for _, slice := range []reflect.Value{slice1, slice2} {
		if !slice.IsValid() {
			tos = append(tos, nil)
			continue
		}
		for i := range slice.Len() {
			tos = append(tos, slice.Index(i).Interface())
		}
	}

	return tos, nil
}

func appendToInterfaceSlice(tov reflect.Value, from ...any) ([]any, error) {
	var tos []any

	for i := range tov.Len() {
		tos = append(tos, tov.Index(i).Interface())
	}

	tos = append(tos, from...)

	return tos, nil
}
