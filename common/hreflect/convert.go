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

package hreflect

import (
	"fmt"
	"math"
	"reflect"
)

var (
	typeInt64   = reflect.TypeFor[int64]()
	typeFloat64 = reflect.TypeFor[float64]()
	typeString  = reflect.TypeFor[string]()
)

// ToInt64 converts v to int64 if possible, returning an error if not.
func ToInt64E(v reflect.Value) (int64, error) {
	if v, ok := ConvertIfPossible(v, typeInt64); ok {
		return v.Int(), nil
	}
	return 0, errConvert(v, "int64")
}

// ToInt64 converts v to int64 if possible. It panics if the conversion is not possible.
func ToInt64(v reflect.Value) int64 {
	vv, err := ToInt64E(v)
	if err != nil {
		panic(err)
	}
	return vv
}

// ToFloat64E converts v to float64 if possible, returning an error if not.
func ToFloat64E(v reflect.Value) (float64, error) {
	if v, ok := ConvertIfPossible(v, typeFloat64); ok {
		return v.Float(), nil
	}
	return 0, errConvert(v, "float64")
}

// ToFloat64 converts v to float64 if possible, panicking if not.
func ToFloat64(v reflect.Value) float64 {
	vv, err := ToFloat64E(v)
	if err != nil {
		panic(err)
	}
	return vv
}

// ToStringE converts v to string if possible, returning an error if not.
func ToStringE(v reflect.Value) (string, error) {
	vv, err := ToStringValueE(v)
	if err != nil {
		return "", err
	}
	return vv.String(), nil
}

func ToStringValueE(v reflect.Value) (reflect.Value, error) {
	if v, ok := ConvertIfPossible(v, typeString); ok {
		return v, nil
	}
	return reflect.Value{}, errConvert(v, "string")
}

// ToString converts v to string if possible, panicking if not.
func ToString(v reflect.Value) string {
	vv, err := ToStringE(v)
	if err != nil {
		panic(err)
	}
	return vv
}

func errConvert(v reflect.Value, s string) error {
	return fmt.Errorf("unable to convert value of type %q to %q", v.Type().String(), s)
}

// ConvertIfPossible tries to convert val to typ if possible.
// This is currently only implemented for int kinds,
// added to handle the move to a new YAML library which produces uint64 for unsigned integers.
// We can expand on this later if needed.
// This conversion is lossless.
// See Issue 14079.
func ConvertIfPossible(val reflect.Value, typ reflect.Type) (reflect.Value, bool) {
	switch val.Kind() {
	case reflect.Ptr, reflect.Interface:
		if val.IsNil() {
			// Return typ's zero value.
			return reflect.Zero(typ), true
		}
		val = val.Elem()
	}

	if val.Type().AssignableTo(typ) {
		// No conversion needed.
		return val, true
	}

	if IsInt(typ.Kind()) {
		return convertToIntIfPossible(val, typ)
	}
	if IsFloat(typ.Kind()) {
		return convertToFloatIfPossible(val, typ)
	}
	if IsUint(typ.Kind()) {
		return convertToUintIfPossible(val, typ)
	}
	if IsString(typ.Kind()) && IsString(val.Kind()) {
		return val.Convert(typ), true
	}

	return reflect.Value{}, false
}

func convertToUintIfPossible(val reflect.Value, typ reflect.Type) (reflect.Value, bool) {
	if IsInt(val.Kind()) {
		i := val.Int()
		if i < 0 {
			return reflect.Value{}, false
		}
		u := uint64(i)
		if typ.OverflowUint(u) {
			return reflect.Value{}, false
		}
		return reflect.ValueOf(u).Convert(typ), true
	}
	if IsUint(val.Kind()) {
		if typ.OverflowUint(val.Uint()) {
			return reflect.Value{}, false
		}
		return val.Convert(typ), true
	}
	if IsFloat(val.Kind()) {
		f := val.Float()
		if f < 0 || f > float64(math.MaxUint64) {
			return reflect.Value{}, false
		}
		if f != math.Trunc(f) {
			return reflect.Value{}, false
		}
		u := uint64(f)
		if typ.OverflowUint(u) {
			return reflect.Value{}, false
		}
		return reflect.ValueOf(u).Convert(typ), true
	}
	return reflect.Value{}, false
}

func convertToFloatIfPossible(val reflect.Value, typ reflect.Type) (reflect.Value, bool) {
	if IsInt(val.Kind()) {
		i := val.Int()
		f := float64(i)
		if typ.OverflowFloat(f) {
			return reflect.Value{}, false
		}
		return reflect.ValueOf(f).Convert(typ), true
	}
	if IsUint(val.Kind()) {
		u := val.Uint()
		f := float64(u)
		if typ.OverflowFloat(f) {
			return reflect.Value{}, false
		}
		return reflect.ValueOf(f).Convert(typ), true
	}
	if IsFloat(val.Kind()) {
		if typ.OverflowFloat(val.Float()) {
			return reflect.Value{}, false
		}
		return val.Convert(typ), true
	}

	return reflect.Value{}, false
}

func convertToIntIfPossible(val reflect.Value, typ reflect.Type) (reflect.Value, bool) {
	if IsInt(val.Kind()) {
		if typ.OverflowInt(val.Int()) {
			return reflect.Value{}, false
		}
		return val.Convert(typ), true
	}
	if IsUint(val.Kind()) {
		if val.Uint() > uint64(math.MaxInt64) {
			return reflect.Value{}, false
		}
		if typ.OverflowInt(int64(val.Uint())) {
			return reflect.Value{}, false
		}
		return val.Convert(typ), true
	}
	if IsFloat(val.Kind()) {
		f := val.Float()
		if f < float64(math.MinInt64) || f > float64(math.MaxInt64) {
			return reflect.Value{}, false
		}
		if f != math.Trunc(f) {
			return reflect.Value{}, false
		}
		if typ.OverflowInt(int64(f)) {
			return reflect.Value{}, false
		}
		return reflect.ValueOf(int64(f)).Convert(typ), true

	}

	return reflect.Value{}, false
}
