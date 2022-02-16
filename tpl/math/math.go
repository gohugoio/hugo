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

// Add adds two numbers.
func (ns *Namespace) Add(a, b interface{}) (interface{}, error) {
	return _math.DoArithmetic(a, b, '+')
}

// Ceil returns the least integer value greater than or equal to x.
func (ns *Namespace) Ceil(x interface{}) (float64, error) {
	xf, err := cast.ToFloat64E(x)
	if err != nil {
		return 0, errors.New("Ceil operator can't be used with non-float value")
	}

	return math.Ceil(xf), nil
}

// Div divides two numbers.
func (ns *Namespace) Div(a, b interface{}) (interface{}, error) {
	return _math.DoArithmetic(a, b, '/')
}

// Floor returns the greatest integer value less than or equal to x.
func (ns *Namespace) Floor(x interface{}) (float64, error) {
	xf, err := cast.ToFloat64E(x)
	if err != nil {
		return 0, errors.New("Floor operator can't be used with non-float value")
	}

	return math.Floor(xf), nil
}

// Log returns the natural logarithm of a number.
func (ns *Namespace) Log(a interface{}) (float64, error) {
	af, err := cast.ToFloat64E(a)
	if err != nil {
		return 0, errors.New("Log operator can't be used with non integer or float value")
	}

	return math.Log(af), nil
}

// Max returns the greater of two numbers.
func (ns *Namespace) Max(a, b interface{}) (float64, error) {
	af, erra := cast.ToFloat64E(a)
	bf, errb := cast.ToFloat64E(b)

	if erra != nil || errb != nil {
		return 0, errors.New("Max operator can't be used with non-float value")
	}

	return math.Max(af, bf), nil
}

// Min returns the smaller of two numbers.
func (ns *Namespace) Min(a, b interface{}) (float64, error) {
	af, erra := cast.ToFloat64E(a)
	bf, errb := cast.ToFloat64E(b)

	if erra != nil || errb != nil {
		return 0, errors.New("Min operator can't be used with non-float value")
	}

	return math.Min(af, bf), nil
}

// Mod returns a % b.
func (ns *Namespace) Mod(a, b interface{}) (int64, error) {
	ai, erra := cast.ToInt64E(a)
	bi, errb := cast.ToInt64E(b)

	if erra != nil || errb != nil {
		return 0, errors.New("modulo operator can't be used with non integer value")
	}

	if bi == 0 {
		return 0, errors.New("the number can't be divided by zero at modulo operation")
	}

	return ai % bi, nil
}

// ModBool returns the boolean of a % b.  If a % b == 0, return true.
func (ns *Namespace) ModBool(a, b interface{}) (bool, error) {
	res, err := ns.Mod(a, b)
	if err != nil {
		return false, err
	}

	return res == int64(0), nil
}

// Mul multiplies two numbers.
func (ns *Namespace) Mul(a, b interface{}) (interface{}, error) {
	return _math.DoArithmetic(a, b, '*')
}

// Pow returns a raised to the power of b.
func (ns *Namespace) Pow(a, b interface{}) (float64, error) {
	af, erra := cast.ToFloat64E(a)
	bf, errb := cast.ToFloat64E(b)

	if erra != nil || errb != nil {
		return 0, errors.New("Pow operator can't be used with non-float value")
	}

	return math.Pow(af, bf), nil
}

// Round returns the nearest integer, rounding half away from zero.
func (ns *Namespace) Round(x interface{}) (float64, error) {
	xf, err := cast.ToFloat64E(x)
	if err != nil {
		return 0, errors.New("Round operator can't be used with non-float value")
	}

	return _round(xf), nil
}

// Sqrt returns the square root of a number.
func (ns *Namespace) Sqrt(a interface{}) (float64, error) {
	af, err := cast.ToFloat64E(a)
	if err != nil {
		return 0, errors.New("Sqrt operator can't be used with non integer or float value")
	}

	return math.Sqrt(af), nil
}

// Sub subtracts two numbers.
func (ns *Namespace) Sub(a, b interface{}) (interface{}, error) {
	return _math.DoArithmetic(a, b, '-')
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
