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
	"errors"
	"fmt"
	"reflect"
)

// Complement gives the elements in the last element of ls that are not in
// any of the others.
//
// All elements of ls must be slices or arrays of comparable types.
//
// The reasoning behind this rather clumsy API is so we can do this in the templates:
//
//	{{ $c := .Pages | complement $last4 }}
func (ns *Namespace) Complement(ls ...any) (any, error) {
	if len(ls) < 2 {
		return nil, errors.New("complement needs at least two arguments")
	}

	universe := ls[len(ls)-1]
	as := ls[:len(ls)-1]

	aset, err := collectIdentities(as...)
	if err != nil {
		return nil, err
	}

	v := reflect.ValueOf(universe)
	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		sl := reflect.MakeSlice(v.Type(), 0, 0)
		for i := range v.Len() {
			ev, _ := indirectInterface(v.Index(i))
			if _, found := aset[normalize(ev)]; !found {
				sl = reflect.Append(sl, ev)
			}
		}
		return sl.Interface(), nil
	default:
		return nil, fmt.Errorf("arguments to complement must be slices or arrays")
	}
}
