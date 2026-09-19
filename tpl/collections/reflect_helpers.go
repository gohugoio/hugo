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
	"math"
	"reflect"
	"strconv"

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
	if v.Kind() == reflect.Interface {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}
	k := v.Kind()
	switch {
	case !v.Type().Comparable():
		return hashing.HashUint64(v.Interface())
	case hreflect.IsNumber(k):
		return numericKey(v)
	}

	vv := types.Unwrapv(v.Interface())
	if ip, ok := vv.(resource.TransientIdentifier); ok {
		return ip.TransientKey()
	}

	return vv
}

// numericKey returns a key that unifies equal numbers across int, uint and
// float types exactly. Integer formatting keeps values beyond the float64
// exact range distinct, so e.g. 2^53 and 2^53+1 no longer collide.
// See issue 15322.
func numericKey(v reflect.Value) string {
	switch {
	case hreflect.IsInt(v.Kind()):
		return "i" + strconv.FormatInt(v.Int(), 10)
	case hreflect.IsUint(v.Kind()):
		u := v.Uint()
		if u <= math.MaxInt64 {
			return "i" + strconv.FormatInt(int64(u), 10)
		}
		return "u" + strconv.FormatUint(u, 10)
	default:
		f := v.Float()
		if !math.IsInf(f, 0) && f == math.Trunc(f) && f >= math.MinInt64 && f < math.MaxUint64 {
			if f < math.MaxInt64 {
				return "i" + strconv.FormatInt(int64(f), 10)
			}
			return "u" + strconv.FormatUint(uint64(f), 10)
		}
		return "f" + strconv.FormatUint(math.Float64bits(f), 16)
	}
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
	if cv, ok := hreflect.ConvertIfPossible(v, typ); ok {
		return cv, nil
	}
	if v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	if hreflect.IsFloat(v.Kind()) && hreflect.IsFloat(typ.Kind()) {
		// Narrowing float64 to float32 loses precision, but that is expected.
		return v.Convert(typ), nil
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
		if tp.Kind() == reflect.Pointer {
			tp = tp.Elem()
		}

		return reflect.New(tp).Interface()
	}
	return nil
}
