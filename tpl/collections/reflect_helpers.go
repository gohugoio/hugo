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
	"time"
)

var (
	zero      reflect.Value
	errorType = reflect.TypeOf((*error)(nil)).Elem()
	timeType  = reflect.TypeOf((*time.Time)(nil)).Elem()
)

func numberToFloat(v reflect.Value) (float64, error) {
	switch kind := v.Kind(); {
	case isFloat(kind):
		return v.Float(), nil
	case isInt(kind):
		return float64(v.Int()), nil
	case isUint(kind):
		return float64(v.Uint()), nil
	case kind == reflect.Interface:
		return numberToFloat(v.Elem())
	default:
		return 0, fmt.Errorf("Invalid kind %s in numberToFloat", kind)
	}
}

// There are potential overflows in this function, but the downconversion of
// int64 etc. into int8 etc. is coming from the synthetic unit tests for Union etc.
// TODO(bep) We should consider normalizing the slices to int64 etc.
func convertNumber(v reflect.Value, to reflect.Kind) (reflect.Value, error) {
	var n reflect.Value
	if isFloat(to) {
		f, err := toFloat(v)
		if err != nil {
			return n, err
		}
		switch to {
		case reflect.Float32:
			n = reflect.ValueOf(float32(f))
		default:
			n = reflect.ValueOf(float64(f))
		}
	} else if isInt(to) {
		i, err := toInt(v)
		if err != nil {
			return n, err
		}
		switch to {
		case reflect.Int:
			n = reflect.ValueOf(int(i))
		case reflect.Int8:
			n = reflect.ValueOf(int8(i))
		case reflect.Int16:
			n = reflect.ValueOf(int16(i))
		case reflect.Int32:
			n = reflect.ValueOf(int32(i))
		case reflect.Int64:
			n = reflect.ValueOf(int64(i))
		}
	} else if isUint(to) {
		i, err := toUint(v)
		if err != nil {
			return n, err
		}
		switch to {
		case reflect.Uint:
			n = reflect.ValueOf(uint(i))
		case reflect.Uint8:
			n = reflect.ValueOf(uint8(i))
		case reflect.Uint16:
			n = reflect.ValueOf(uint16(i))
		case reflect.Uint32:
			n = reflect.ValueOf(uint32(i))
		case reflect.Uint64:
			n = reflect.ValueOf(uint64(i))
		}

	}

	if !n.IsValid() {
		return n, errors.New("invalid values")
	}

	return n, nil

}

func isNumber(kind reflect.Kind) bool {
	return isInt(kind) || isUint(kind) || isFloat(kind)
}

func isInt(kind reflect.Kind) bool {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	default:
		return false
	}
}

func isUint(kind reflect.Kind) bool {
	switch kind {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	default:
		return false
	}
}

func isFloat(kind reflect.Kind) bool {
	switch kind {
	case reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}
