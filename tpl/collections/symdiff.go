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

// SymDiff returns the symmetric difference of s1 and s2.
// Arguments must be either a slice or an array of comparable types.
func (ns *Namespace) SymDiff(s2, s1 any) (any, error) {
	ids1, err := collectIdentities(s1)
	if err != nil {
		return nil, err
	}
	ids2, err := collectIdentities(s2)
	if err != nil {
		return nil, err
	}

	var slice reflect.Value
	var sliceElemType reflect.Type

	for i, s := range []any{s1, s2} {
		v := reflect.ValueOf(s)

		switch v.Kind() {
		case reflect.Array, reflect.Slice:
			if i == 0 {
				sliceType := v.Type()
				sliceElemType = sliceType.Elem()
				slice = reflect.MakeSlice(sliceType, 0, 0)
			}

			for i := 0; i < v.Len(); i++ {
				ev, _ := indirectInterface(v.Index(i))
				key := normalize(ev)

				// Append if the key is not in their intersection.
				if ids1[key] != ids2[key] {
					v, err := convertValue(ev, sliceElemType)
					if err != nil {
						return nil, fmt.Errorf("symdiff: failed to convert value: %w", err)
					}
					slice = reflect.Append(slice, v)
				}
			}
		default:
			return nil, fmt.Errorf("arguments to symdiff must be slices or arrays")
		}
	}

	return slice.Interface(), nil
}
