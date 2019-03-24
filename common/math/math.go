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

package math

import (
	"errors"
	"reflect"
)

// DoArithmetic performs arithmetic operations (+,-,*,/) using reflection to
// determine the type of the two terms.
func DoArithmetic(a, b interface{}, op rune) (interface{}, error) {
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)
	var ai, bi int64
	var af, bf float64
	var au, bu uint64
	switch av.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		ai = av.Int()
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			bi = bv.Int()
		case reflect.Float32, reflect.Float64:
			af = float64(ai) // may overflow
			ai = 0
			bf = bv.Float()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			bu = bv.Uint()
			if ai >= 0 {
				au = uint64(ai)
				ai = 0
			} else {
				bi = int64(bu) // may overflow
				bu = 0
			}
		default:
			return nil, errors.New("can't apply the operator to the values")
		}
	case reflect.Float32, reflect.Float64:
		af = av.Float()
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			bf = float64(bv.Int()) // may overflow
		case reflect.Float32, reflect.Float64:
			bf = bv.Float()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			bf = float64(bv.Uint()) // may overflow
		default:
			return nil, errors.New("can't apply the operator to the values")
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		au = av.Uint()
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			bi = bv.Int()
			if bi >= 0 {
				bu = uint64(bi)
				bi = 0
			} else {
				ai = int64(au) // may overflow
				au = 0
			}
		case reflect.Float32, reflect.Float64:
			af = float64(au) // may overflow
			au = 0
			bf = bv.Float()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			bu = bv.Uint()
		default:
			return nil, errors.New("can't apply the operator to the values")
		}
	case reflect.String:
		as := av.String()
		if bv.Kind() == reflect.String && op == '+' {
			bs := bv.String()
			return as + bs, nil
		}
		return nil, errors.New("can't apply the operator to the values")
	default:
		return nil, errors.New("can't apply the operator to the values")
	}

	switch op {
	case '+':
		if ai != 0 || bi != 0 {
			return ai + bi, nil
		} else if af != 0 || bf != 0 {
			return af + bf, nil
		} else if au != 0 || bu != 0 {
			return au + bu, nil
		}
		return 0, nil
	case '-':
		if ai != 0 || bi != 0 {
			return ai - bi, nil
		} else if af != 0 || bf != 0 {
			return af - bf, nil
		} else if au != 0 || bu != 0 {
			return au - bu, nil
		}
		return 0, nil
	case '*':
		if ai != 0 || bi != 0 {
			return ai * bi, nil
		} else if af != 0 || bf != 0 {
			return af * bf, nil
		} else if au != 0 || bu != 0 {
			return au * bu, nil
		}
		return 0, nil
	case '/':
		if bi != 0 {
			return ai / bi, nil
		} else if bf != 0 {
			return af / bf, nil
		} else if bu != 0 {
			return au / bu, nil
		}
		return nil, errors.New("can't divide the value by 0")
	default:
		return nil, errors.New("there is no such an operation")
	}
}
