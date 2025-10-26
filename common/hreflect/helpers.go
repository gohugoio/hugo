// Copyright 2025 The Hugo Authors. All rights reserved.
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

// Package hreflect contains reflect helpers.
package hreflect

import (
	"context"
	"reflect"
	"sync"
	"time"

	"github.com/gohugoio/hugo/common/htime"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/common/types"
)

// IsInterfaceOrPointer returns whether the given kind is an interface or a pointer.
func IsInterfaceOrPointer(kind reflect.Kind) bool {
	return kind == reflect.Interface || kind == reflect.Ptr
}

// TODO(bep) replace the private versions in /tpl with these.
// IsNumber returns whether the given kind is a number.
func IsNumber(kind reflect.Kind) bool {
	return IsInt(kind) || IsUint(kind) || IsFloat(kind)
}

// IsInt returns whether the given kind is an int.
func IsInt(kind reflect.Kind) bool {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	default:
		return false
	}
}

// IsUint returns whether the given kind is an uint.
func IsUint(kind reflect.Kind) bool {
	switch kind {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	default:
		return false
	}
}

// IsFloat returns whether the given kind is a float.
func IsFloat(kind reflect.Kind) bool {
	switch kind {
	case reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

// IsString returns whether the given kind is a string.
func IsString(kind reflect.Kind) bool {
	return kind == reflect.String
}

// IsTruthful returns whether in represents a truthful value.
// See IsTruthfulValue
func IsTruthful(in any) bool {
	switch v := in.(type) {
	case reflect.Value:
		return IsTruthfulValue(v)
	default:
		return IsTruthfulValue(reflect.ValueOf(in))
	}
}

// IsMap reports whether v is a map.
func IsMap(v any) bool {
	return reflect.ValueOf(v).Kind() == reflect.Map
}

// IsSlice reports whether v is a slice.
func IsSlice(v any) bool {
	return reflect.ValueOf(v).Kind() == reflect.Slice
}

var zeroType = reflect.TypeOf((*types.Zeroer)(nil)).Elem()

var isZeroCache sync.Map

func implementsIsZero(tp reflect.Type) bool {
	v, ok := isZeroCache.Load(tp)
	if ok {
		return v.(bool)
	}
	implements := tp.Implements(zeroType)
	isZeroCache.Store(tp, implements)
	return implements
}

// IsTruthfulValue returns whether the given value has a meaningful truth value.
// This is based on template.IsTrue in Go's stdlib, but also considers
// IsZero and any interface value will be unwrapped before it's considered
// for truthfulness.
//
// Based on:
// https://github.com/golang/go/blob/178a2c42254166cffed1b25fb1d3c7a5727cada6/src/text/template/exec.go#L306
func IsTruthfulValue(val reflect.Value) (truth bool) {
	val, isNil := Indirect(val)

	if !val.IsValid() {
		// Something like: var x any, never set. It's a form of nil.
		return
	}

	if val.Kind() == reflect.Pointer && isNil {
		return
	}

	if implementsIsZero(val.Type()) {
		return !val.Interface().(types.Zeroer).IsZero()
	}

	switch val.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		truth = val.Len() > 0
	case reflect.Bool:
		truth = val.Bool()
	case reflect.Complex64, reflect.Complex128:
		truth = val.Complex() != 0
	case reflect.Chan, reflect.Func, reflect.Ptr, reflect.Interface:
		truth = !val.IsNil()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		truth = val.Int() != 0
	case reflect.Float32, reflect.Float64:
		truth = val.Float() != 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		truth = val.Uint() != 0
	case reflect.Struct:
		truth = true // Struct values are always true.
	default:
		return
	}

	return
}

type methodKey struct {
	typ  reflect.Type
	name string
}

var (
	methodIndexCache sync.Map
	methodCache      sync.Map
)

// GetMethodByNameForType returns the method with the given name for the given type,
// or a zero Method if no such method exists.
// It panics if tp is an interface type.
// It caches the lookup.
func GetMethodByNameForType(tp reflect.Type, name string) reflect.Method {
	if tp.Kind() == reflect.Interface {
		// Func field is nil for interface types.
		panic("not supported for interface types")
	}
	k := methodKey{tp, name}
	v, found := methodCache.Load(k)
	if found {
		return v.(reflect.Method)
	}
	m, _ := tp.MethodByName(name)
	methodCache.Store(k, m)
	return m
}

// GetMethodByName is the same as reflect.Value.MethodByName, but it caches the lookup.
func GetMethodByName(v reflect.Value, name string) reflect.Value {
	index := GetMethodIndexByName(v.Type(), name)

	if index == -1 {
		return reflect.Value{}
	}

	return v.Method(index)
}

// GetMethodIndexByName returns the index of the method with the given name, or
// -1 if no such method exists.
func GetMethodIndexByName(tp reflect.Type, name string) int {
	k := methodKey{tp, name}
	v, found := methodIndexCache.Load(k)
	if found {
		return v.(int)
	}
	m, ok := tp.MethodByName(name)
	index := m.Index
	if !ok {
		index = -1
	}
	methodIndexCache.Store(k, index)

	if !ok {
		return -1
	}

	return m.Index
}

var (
	timeType           = reflect.TypeOf((*time.Time)(nil)).Elem()
	asTimeProviderType = reflect.TypeOf((*htime.AsTimeProvider)(nil)).Elem()
)

// IsTime returns whether tp is a time.Time type or if it can be converted into one
// in ToTime.
func IsTime(tp reflect.Type) bool {
	if tp == timeType {
		return true
	}

	if tp.Implements(asTimeProviderType) {
		return true
	}
	return false
}

// IsValid returns whether v is not nil and a valid value.
func IsValid(v reflect.Value) bool {
	if !v.IsValid() {
		return false
	}

	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return !v.IsNil()
	}

	return true
}

// AsTime returns v as a time.Time if possible.
// The given location is only used if the value implements AsTimeProvider (e.g. go-toml local).
// A zero Time and false is returned if this isn't possible.
// Note that this function does not accept string dates.
func AsTime(v reflect.Value, loc *time.Location) (time.Time, bool) {
	if v.Kind() == reflect.Interface {
		return AsTime(v.Elem(), loc)
	}

	if v.Type() == timeType {
		return v.Interface().(time.Time), true
	}

	if v.Type().Implements(asTimeProviderType) {
		return v.Interface().(htime.AsTimeProvider).AsTime(loc), true
	}

	return time.Time{}, false
}

// ToSliceAny converts the given value to a slice of any if possible.
func ToSliceAny(v any) ([]any, bool) {
	if v == nil {
		return nil, false
	}
	switch vv := v.(type) {
	case []any:
		return vv, true
	default:
		vvv := reflect.ValueOf(v)
		if vvv.Kind() == reflect.Slice {
			out := make([]any, vvv.Len())
			for i := range vvv.Len() {
				out[i] = vvv.Index(i).Interface()
			}
			return out, true
		}
	}
	return nil, false
}

// CallMethodByName calls the method with the given name on v.
func CallMethodByName(cxt context.Context, name string, v reflect.Value) []reflect.Value {
	fn := v.MethodByName(name)
	var args []reflect.Value
	tp := fn.Type()
	if tp.NumIn() > 0 {
		if tp.NumIn() > 1 {
			panic("not supported")
		}
		first := tp.In(0)
		if IsContextType(first) {
			args = append(args, reflect.ValueOf(cxt))
		}
	}

	return fn.Call(args)
}

// Indirect unwraps interfaces and pointers until it finds a non-interface/pointer value.
// If a nil is encountered, the second return value is true.
// If a pointer to a struct is encountered, it is not unwrapped.
func Indirect(v reflect.Value) (vv reflect.Value, isNil bool) {
	for ; IsInterfaceOrPointer(v.Kind()); v = v.Elem() {
		if IsNil(v) {
			return v, true
		}
		if v.Kind() != reflect.Interface {
			// A pointer.
			if v.NumMethod() > 0 {
				break
			}
			if v.Elem().Kind() == reflect.Struct {
				// Avoid unwrapping pointers to structs.
				break
			}
		}
	}
	return v, false
}

// IndirectElem is like Indirect, but if the final value is a pointer, it unwraps it.
func IndirectElem(v reflect.Value) (vv reflect.Value, isNil bool) {
	vv, isNil = Indirect(v)
	if isNil {
		return vv, isNil
	}
	if vv.Kind() == reflect.Pointer {
		vv = vv.Elem()
	}
	return vv, isNil
}

// IsNil reports whether v is nil.
// Based on reflect.Value.IsNil, but also considers invalid values as nil.
func IsNil(v reflect.Value) bool {
	if !v.IsValid() {
		return true
	}
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return v.IsNil()
	}
	return false
}

var contextInterface = reflect.TypeOf((*context.Context)(nil)).Elem()

var isContextCache = maps.NewCache[reflect.Type, bool]()

type k string

var contextTypeValue = reflect.TypeOf(context.WithValue(context.Background(), k("key"), 32))

// IsContextType returns whether tp is a context.Context type.
func IsContextType(tp reflect.Type) bool {
	if tp == contextTypeValue {
		return true
	}
	if tp == contextInterface {
		return true
	}

	isContext, _ := isContextCache.GetOrCreate(tp, func() (bool, error) {
		return tp.Implements(contextInterface), nil
	})
	return isContext
}
