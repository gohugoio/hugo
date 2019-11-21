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

	"github.com/spf13/cast"

	"github.com/gohugoio/hugo/common/maps"
)

// Index returns the result of indexing its first argument by the following
// arguments. Thus "index x 1 2 3" is, in Go syntax, x[1][2][3]. Each
// indexed item must be a map, slice, or array.
//
// Copied from Go stdlib src/text/template/funcs.go.
//
// We deviate from the stdlib due to https://github.com/golang/go/issues/14751.
//
// TODO(moorereason): merge upstream changes.
func (ns *Namespace) Index(item interface{}, args ...interface{}) (interface{}, error) {
	v := reflect.ValueOf(item)
	if !v.IsValid() {
		return nil, errors.New("index of untyped nil")
	}

	lowerm, ok := item.(maps.Params)
	if ok {
		return lowerm.Get(cast.ToStringSlice(args)...), nil
	}

	var indices []interface{}

	if len(args) == 1 {
		v := reflect.ValueOf(args[0])
		if v.Kind() == reflect.Slice {
			for i := 0; i < v.Len(); i++ {
				indices = append(indices, v.Index(i).Interface())
			}
		}
	}

	if indices == nil {
		indices = args
	}

	for _, i := range indices {
		index := reflect.ValueOf(i)
		var isNil bool
		if v, isNil = indirect(v); isNil {
			return nil, errors.New("index of nil pointer")
		}
		switch v.Kind() {
		case reflect.Array, reflect.Slice, reflect.String:
			var x int64
			switch index.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				x = index.Int()
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				x = int64(index.Uint())
			case reflect.Invalid:
				return nil, errors.New("cannot index slice/array with nil")
			default:
				return nil, fmt.Errorf("cannot index slice/array with type %s", index.Type())
			}
			if x < 0 || x >= int64(v.Len()) {
				// We deviate from stdlib here.  Don't return an error if the
				// index is out of range.
				return nil, nil
			}
			v = v.Index(int(x))
		case reflect.Map:
			index, err := prepareArg(index, v.Type().Key())
			if err != nil {
				return nil, err
			}

			if x := v.MapIndex(index); x.IsValid() {
				v = x
			} else {
				v = reflect.Zero(v.Type().Elem())
			}
		case reflect.Invalid:
			// the loop holds invariant: v.IsValid()
			panic("unreachable")
		default:
			return nil, fmt.Errorf("can't index item of type %s", v.Type())
		}
	}
	return v.Interface(), nil
}

// prepareArg checks if value can be used as an argument of type argType, and
// converts an invalid value to appropriate zero if possible.
//
// Copied from Go stdlib src/text/template/funcs.go.
func prepareArg(value reflect.Value, argType reflect.Type) (reflect.Value, error) {
	if !value.IsValid() {
		if !canBeNil(argType) {
			return reflect.Value{}, fmt.Errorf("value is nil; should be of type %s", argType)
		}
		value = reflect.Zero(argType)
	}
	if !value.Type().AssignableTo(argType) {
		return reflect.Value{}, fmt.Errorf("value has type %s; should be %s", value.Type(), argType)
	}
	return value, nil
}

// canBeNil reports whether an untyped nil can be assigned to the type. See reflect.Zero.
//
// Copied from Go stdlib src/text/template/exec.go.
func canBeNil(typ reflect.Type) bool {
	switch typ.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return true
	}
	return false
}
