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

var (
	errMustTwoValueError = errors.New("A minimum of two values is REQUIRED")
)

// New returns a new instance of the math-namespaced template functions.
func New() *Namespace {
	return &Namespace{}
}

// Namespace provides template functions for the "math" namespace.
type Namespace struct{}

// Add adds the multivalued addends n1 and n2 or more values.
func (ns *Namespace) Add(inputs ...any) (any, error) {
	if len(inputs) < 2 {
		return nil, errMustTwoValueError
	}
	return ns.doArithmetic(inputs, '+')
}

// Ceil returns the least integer value greater than or equal to n.
func (ns *Namespace) Ceil(n any) (float64, error) {
	xf, err := cast.ToFloat64E(n)
	if err != nil {
		return 0, errors.New("Ceil operator can't be used with non-float value")
	}

	return math.Ceil(xf), nil
}

// Div divides n1 by n2.
func (ns *Namespace) Div(inputs ...any) (any, error) {
	if len(inputs) < 2 {
		return nil, errMustTwoValueError
	}
	return ns.doArithmetic(inputs, '/')
}

// Floor returns the greatest integer value less than or equal to n.
func (ns *Namespace) Floor(n any) (float64, error) {
	xf, err := cast.ToFloat64E(n)
	if err != nil {
		return 0, errors.New("Floor operator can't be used with non-float value")
	}

	return math.Floor(xf), nil
}

// Log returns the natural logarithm of the number n.
func (ns *Namespace) Log(n any) (float64, error) {
	af, err := cast.ToFloat64E(n)
	if err != nil {
		return 0, errors.New("Log operator can't be used with non integer or float value")
	}

	return math.Log(af), nil
}

// Max returns the greater of the multivalued numbers n1 and n2 or more values.
func (ns *Namespace) Max(inputs ...any) (maximum float64, err error) {
	if len(inputs) < 2 {
		err = errMustTwoValueError
		return
	}
	var value float64
	for index, input := range inputs {
		value, err = cast.ToFloat64E(input)
		if err != nil {
			err = errors.New("Max operator can't be used with non-float value")
			return
		}
		if index == 0 {
			maximum = value
			continue
		}
		maximum = math.Max(value, maximum)
	}
	return
}

// Min returns the smaller of multivalued numbers n1 and n2 or more values.
func (ns *Namespace) Min(inputs ...any) (minimum float64, err error) {
	if len(inputs) < 2 {
		err = errMustTwoValueError
		return
	}
	var value float64
	for index, input := range inputs {
		value, err = cast.ToFloat64E(input)
		if err != nil {
			err = errors.New("Max operator can't be used with non-float value")
			return
		}
		if index == 0 {
			minimum = value
			continue
		}
		minimum = math.Min(value, minimum)
	}
	return
}

// Mod returns n1 % n2.
func (ns *Namespace) Mod(n1, n2 any) (int64, error) {
	ai, erra := cast.ToInt64E(n1)
	bi, errb := cast.ToInt64E(n2)

	if erra != nil || errb != nil {
		return 0, errors.New("modulo operator can't be used with non integer value")
	}

	if bi == 0 {
		return 0, errors.New("the number can't be divided by zero at modulo operation")
	}

	return ai % bi, nil
}

// ModBool returns the boolean of n1 % n2.  If n1 % n2 == 0, return true.
func (ns *Namespace) ModBool(n1, n2 any) (bool, error) {
	res, err := ns.Mod(n1, n2)
	if err != nil {
		return false, err
	}

	return res == int64(0), nil
}

// Mul multiplies the multivalued numbers n1 and n2 or more values.
func (ns *Namespace) Mul(inputs ...any) (any, error) {
	if len(inputs) < 2 {
		return nil, errMustTwoValueError
	}
	return ns.doArithmetic(inputs, '*')
}

// Pow returns n1 raised to the power of n2.
func (ns *Namespace) Pow(n1, n2 any) (float64, error) {
	af, erra := cast.ToFloat64E(n1)
	bf, errb := cast.ToFloat64E(n2)

	if erra != nil || errb != nil {
		return 0, errors.New("Pow operator can't be used with non-float value")
	}

	return math.Pow(af, bf), nil
}

// Round returns the integer nearest to n, rounding half away from zero.
func (ns *Namespace) Round(n any) (float64, error) {
	xf, err := cast.ToFloat64E(n)
	if err != nil {
		return 0, errors.New("Round operator can't be used with non-float value")
	}

	return _round(xf), nil
}

// Sqrt returns the square root of the number n.
func (ns *Namespace) Sqrt(n any) (float64, error) {
	af, err := cast.ToFloat64E(n)
	if err != nil {
		return 0, errors.New("Sqrt operator can't be used with non integer or float value")
	}

	return math.Sqrt(af), nil
}

// Sub subtracts multivalued.
func (ns *Namespace) Sub(inputs ...any) (any, error) {
	if len(inputs) < 2 {
		return nil, errMustTwoValueError
	}
	return ns.doArithmetic(inputs, '-')
}

func (ns *Namespace) doArithmetic(inputs []any, operation rune) (value any, err error) {
	value = inputs[0]
	for i := 1; i < len(inputs); i++ {
		value, err = _math.DoArithmetic(value, inputs[i], operation)
		if err != nil {
			return
		}
	}
	return
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
