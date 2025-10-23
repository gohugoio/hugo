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

package collections

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/gohugoio/hugo/common/hashing"
	"github.com/gohugoio/hugo/common/hreflect"
	"github.com/gohugoio/hugo/common/types"
	"github.com/gohugoio/hugo/resources/resource"
)

var (
	zero      reflect.Value
	errorType = reflect.TypeFor[error]()
)

// normalizes different numeric types if isNumber
// or get the hash values if not Comparable (such as map or struct)
// to make them comparable
func normalize(v reflect.Value) any {
	k := v.Kind()
	switch {
	case !v.Type().Comparable():
		return hashing.HashUint64(v.Interface())
	case hreflect.IsNumber(k):
		f, err := hreflect.ToFloat64E(v)
		if err == nil {
			return f
		}
	}

	vv := types.Unwrapv(v.Interface())
	if ip, ok := vv.(resource.TransientIdentifier); ok {
		return ip.TransientKey()
	}

	return vv
}

// collects identities from the slices in seqs into a set. Numeric values are normalized,
// pointers unwrapped.
func collectIdentities(seqs ...any) (map[any]bool, error) {
	seen := make(map[any]bool)
	for _, seq := range seqs {
		v := reflect.ValueOf(seq)
		switch v.Kind() {
		case reflect.Array, reflect.Slice:
			for i := range v.Len() {
				ev, _ := hreflect.Indirect(v.Index(i))

				if !ev.Type().Comparable() {
					return nil, errors.New("elements must be comparable")
				}

				seen[normalize(ev)] = true
			}
		default:
			return nil, fmt.Errorf("arguments must be slices or arrays")
		}
	}

	return seen, nil
}

// We have some different numeric and string types that we try to behave like
// they were the same.
func convertValue(v reflect.Value, to reflect.Type) (reflect.Value, error) {
	if v.Type().AssignableTo(to) {
		return v, nil
	}
	switch kind := to.Kind(); {
	case kind == reflect.String:
		return hreflect.ToStringValueE(v)
	case hreflect.IsNumber(kind):
		return convertNumber(v, to)
	default:
		return reflect.Value{}, fmt.Errorf("%s is not assignable to %s", v.Type(), to)
	}
}

func convertNumber(v reflect.Value, typ reflect.Type) (reflect.Value, error) {
	if v, ok := hreflect.ConvertIfPossible(v, typ); ok {
		return v, nil
	}
	return reflect.Value{}, fmt.Errorf("unable to convert value of type %q to %q", v.Type().String(), typ.String())
}

func newSliceElement(items any) any {
	tp := reflect.TypeOf(items)
	if tp == nil {
		return nil
	}
	switch tp.Kind() {
	case reflect.Array, reflect.Slice:
		tp = tp.Elem()
		if tp.Kind() == reflect.Ptr {
			tp = tp.Elem()
		}

		return reflect.New(tp).Interface()
	}
	return nil
}
