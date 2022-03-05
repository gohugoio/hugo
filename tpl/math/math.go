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

// Package math provides template functions for mathematical operations.
package math

import (
	"errors"
	"math"
	"sync/atomic"

	_math "github.com/gohugoio/hugo/common/math"

	"github.com/spf13/cast"
)

// New returns a new instance of the math-namespaced template functions.
func New() *Namespace {
	return &Namespace{}
}

// Namespace provides template functions for the "math" namespace.
type Namespace struct{}

// Add adds the two numbers num1 and num2.
func (ns *Namespace) Add(num1, num2 interface{}) (interface{}, error) {
	return _math.DoArithmetic(num1, num2, '+')
}

// Ceil returns the least integer value greater than or equal to num.
func (ns *Namespace) Ceil(num interface{}) (float64, error) {
	xf, err := cast.ToFloat64E(num)
	if err != nil {
		return 0, errors.New("Ceil operator can't be used with non-float value")
	}

	return math.Ceil(xf), nil
}

// Div divides num1 by num2.
func (ns *Namespace) Div(num1, num2 interface{}) (interface{}, error) {
	return _math.DoArithmetic(num1, num2, '/')
}

// Floor returns the greatest integer value less than or equal to num.
func (ns *Namespace) Floor(num interface{}) (float64, error) {
	xf, err := cast.ToFloat64E(num)
	if err != nil {
		return 0, errors.New("Floor operator can't be used with non-float value")
	}

	return math.Floor(xf), nil
}

// Log returns the natural logarithm of num.
func (ns *Namespace) Log(num interface{}) (float64, error) {
	af, err := cast.ToFloat64E(num)
	if err != nil {
		return 0, errors.New("Log operator can't be used with non integer or float value")
	}

	return math.Log(af), nil
}

// Max returns the greater of num1 or num2.
func (ns *Namespace) Max(num1, num2 interface{}) (float64, error) {
	af, erra := cast.ToFloat64E(num1)
	bf, errb := cast.ToFloat64E(num2)

	if erra != nil || errb != nil {
		return 0, errors.New("Max operator can't be used with non-float value")
	}

	return math.Max(af, bf), nil
}

// Min returns the smaller of two num1 or num2.
func (ns *Namespace) Min(num1, num2 interface{}) (float64, error) {
	af, erra := cast.ToFloat64E(num1)
	bf, errb := cast.ToFloat64E(num2)

	if erra != nil || errb != nil {
		return 0, errors.New("Min operator can't be used with non-float value")
	}

	return math.Min(af, bf), nil
}

// Mod returns num1 % num2.
func (ns *Namespace) Mod(num1, num2 interface{}) (int64, error) {
	ai, erra := cast.ToInt64E(num1)
	bi, errb := cast.ToInt64E(num2)

	if erra != nil || errb != nil {
		return 0, errors.New("modulo operator can't be used with non integer value")
	}

	if bi == 0 {
		return 0, errors.New("the number can't be divided by zero at modulo operation")
	}

	return ai % bi, nil
}

// ModBool returns the boolean of num1 % num2.  If  num1 % num2 == 0, return true.
func (ns *Namespace) ModBool(num1, num2 interface{}) (bool, error) {
	res, err := ns.Mod(num1, num2)
	if err != nil {
		return false, err
	}

	return res == int64(0), nil
}

// Mul multiplies num1 with num2.
func (ns *Namespace) Mul(num1, num2 interface{}) (interface{}, error) {
	return _math.DoArithmetic(num1, num2, '*')
}

// Pow returns num1 raised to the power of num2.
func (ns *Namespace) Pow(num1, num2 interface{}) (float64, error) {
	af, erra := cast.ToFloat64E(num1)
	bf, errb := cast.ToFloat64E(num2)

	if erra != nil || errb != nil {
		return 0, errors.New("Pow operator can't be used with non-float value")
	}

	return math.Pow(af, bf), nil
}

// Round returns the nearest integer of num, rounding half away from zero.
func (ns *Namespace) Round(num interface{}) (float64, error) {
	xf, err := cast.ToFloat64E(num)
	if err != nil {
		return 0, errors.New("Round operator can't be used with non-float value")
	}

	return _round(xf), nil
}

// Sqrt returns the square root of a num.
func (ns *Namespace) Sqrt(num interface{}) (float64, error) {
	af, err := cast.ToFloat64E(num)
	if err != nil {
		return 0, errors.New("Sqrt operator can't be used with non integer or float value")
	}

	return math.Sqrt(af), nil
}

// Sub subtracts num2 fron num1.
func (ns *Namespace) Sub(num1, num2 interface{}) (interface{}, error) {
	return _math.DoArithmetic(num1, num2, '-')
}

var counter uint64

// Counter increments and returns a global counter.
// This was originally added to be used in tests where now.UnixNano did not
// have the needed precision (especially on Windows).
// Note that given the parallel nature of Hugo, you cannot use this to get sequences of numbers,
// and the counter will reset on new builds.
func (ns *Namespace) Counter() uint64 {
	return atomic.AddUint64(&counter, uint64(1))
}
