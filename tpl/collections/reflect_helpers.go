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
	"reflect"
)

func numberToFloat(v reflect.Value) float64 {
	switch kind := v.Kind(); {
	case isFloat(kind):
		return v.Float()
	case isInt(kind):
		return float64(v.Int())
	case isUInt(kind):
		return float64(v.Uint())
	case kind == reflect.Interface:
		return numberToFloat(v.Elem())
	default:
		panic("Invalid type in numberToFloat")
	}
}

func isNumber(kind reflect.Kind) bool {
	return isInt(kind) || isUInt(kind) || isFloat(kind)
}

func isInt(kind reflect.Kind) bool {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	default:
		return false
	}
}

func isUInt(kind reflect.Kind) bool {
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
