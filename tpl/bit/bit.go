// Copyright 2024 The Hugo Authors. All rights reserved.
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

// Package bit provides template functions for bitwise operations.
package bit

import (
	"errors"
	"math/bits"

	"github.com/spf13/cast"
)

var (
	errMustTwoNumbersError = errors.New("must provide at least two numbers")
	errAndNotIntegerError = errors.New("And operator can't be used with non integer value")
	errClearNotIntegerError = errors.New("Clear operator can't be used with non integer value")
	errOrNotIntegerError = errors.New("Or operator can't be used with non integer value")
	errXorNotIntegerError = errors.New("Xor operator can't be used with non integer value")
	errXnorNotIntegerError = errors.New("Xnor operator can't be used with non integer value")
)

// New returns a new instance of the bit-namespaced template functions.
func New() *Namespace {
	return &Namespace{}
}

// Namespace provides template functions for the "bit" namespace.
type Namespace struct{}

// And returns the bitwise AND of n1 and n2 or more values.
func (ns *Namespace) And(inputs ...any) (int64, error) {
	if len(inputs) < 2 {
		return 0, errMustTwoNumbersError
	}
	value, err := cast.ToInt64E(inputs[0])
	if err != nil {
		return 0, errAndNotIntegerError
	}
	for i := 1; i < len(inputs); i++ {
		andval, anderr := cast.ToInt64E(inputs[i])
		if anderr != nil {
			return 0, errAndNotIntegerError
		}
		value &= andval
	}
	return value, nil
}

// Clear returns the bitwise AND of n1 with the NOT of n2 or following values.
func (ns *Namespace) Clear(inputs ...any) (int64, error) {
	if len(inputs) < 2 {
		return 0, errMustTwoNumbersError
	}
	value, err := cast.ToInt64E(inputs[0])
	if err != nil {
		return 0, errAndNotIntegerError
	}
	for i := 1; i < len(inputs); i++ {
		clearval, clearerr := cast.ToInt64E(inputs[i])
		if clearerr != nil {
			return 0, errAndNotIntegerError
		}
		value &= ^clearval
	}
	return value, nil
}

// Extract returns the last l bits of (n >> s).
func (ns *Namespace) Extract(n any, l any, s any) (int64, error) {
	ni, nerr := cast.ToInt64E(n)
	li, lerr := cast.ToUint8E(l)
	si, serr := cast.ToUint8E(s)
	if nerr != nil || serr != nil || lerr != nil {
		return 0, errors.New("Extract operator can't be used with non integer value")
	}
	return (ni >> si) & ((int64(1) << li) - 1), nil
}

// LeadingZeros returns the number of leading 0 bits in n (when cast to uint64).
func (ns *Namespace) LeadingZeros(n any) (int64, error) {
	value, err := cast.ToInt64E(n)
	if err != nil {
		return 0, errors.New("LeadingZeros operator can't be used with non integer value")
	}
	return int64(bits.LeadingZeros64(uint64(value))), nil
}

// Not returns the bitwise NOT (over 64 bits) of n.
func (ns *Namespace) Not(n any) (int64, error) {
	value, err := cast.ToInt64E(n)
	if err != nil {
		return 0, errors.New("Not operator can't be used with non integer value")
	}
	return ^value, nil
}

// OnesCount returns the number of bits set to 1 (known as population count) in n.
func (ns *Namespace) OnesCount(n any) (int64, error) {
	value, err := cast.ToInt64E(n)
	if err != nil {
		return 0, errors.New("OnesCount operator can't be used with non integer value")
	}
	return int64(bits.OnesCount64(uint64(value))), nil
}

// Or returns the bitwise OR of n1 and n2 or more values.
func (ns *Namespace) Or(inputs ...any) (int64, error) {
	if len(inputs) < 2 {
		return 0, errMustTwoNumbersError
	}
	value, err := cast.ToInt64E(inputs[0])
	if err != nil {
		return 0, errOrNotIntegerError
	}
	for i := 1; i < len(inputs); i++ {
		orval, orerr := cast.ToInt64E(inputs[i])
		if orerr != nil {
			return 0, errOrNotIntegerError
		}
		value |= orval
	}
	return value, nil
}

// ShiftLeft returns n shifted s bits to the left.
func (ns *Namespace) ShiftLeft(n any, s any) (int64, error) {
	ni, nerr := cast.ToInt64E(n)
	si, serr := cast.ToUint8E(s)
	if nerr != nil || serr != nil {
		return 0, errors.New("ShiftLeft operator can't be used with non integer value")
	}
	return ni << si, nil
}

// ShiftRight returns n (arethmetically) shifted s bits to the right.
func (ns *Namespace) ShiftRight(n any, s any) (int64, error) {
	ni, nerr := cast.ToInt64E(n)
	si, serr := cast.ToUint8E(s)
	if nerr != nil || serr != nil {
		return 0, errors.New("ShiftRight operator can't be used with non integer value")
	}
	return ni >> si, nil
}

// TrailingZeros returns the number of trailing 0 bits in n (when cast to uint64).
func (ns *Namespace) TrailingZeros(n any) (int64, error) {
	value, err := cast.ToInt64E(n)
	if err != nil {
		return 0, errors.New("TrailingZeros operator can't be used with non integer value")
	}
	return int64(bits.TrailingZeros64(uint64(value))), nil
}

// Xnor returns the bitwise XNOR (Exclusive Not-OR) of n1 and n2 or more values. Equivalent to Not (Xor inputs...).
func (ns *Namespace) Xnor(inputs ...any) (int64, error) {
	if len(inputs) < 2 {
		return 0, errMustTwoNumbersError
	}
	value, err := cast.ToInt64E(inputs[0])
	if err != nil {
		return 0, errXnorNotIntegerError
	}
	for i := 1; i < len(inputs); i++ {
		xorval, xorerr := cast.ToInt64E(inputs[i])
		if xorerr != nil {
			return 0, errXnorNotIntegerError
		}
		value ^= xorval
	}
	return ^value, nil
}

// Xor returns the bitwise XOR (Exclusive OR) of n1 and n2 or more values.
func (ns *Namespace) Xor(inputs ...any) (int64, error) {
	if len(inputs) < 2 {
		return 0, errMustTwoNumbersError
	}
	value, err := cast.ToInt64E(inputs[0])
	if err != nil {
		return 0, errXorNotIntegerError
	}
	for i := 1; i < len(inputs); i++ {
		xorval, xorerr := cast.ToInt64E(inputs[i])
		if xorerr != nil {
			return 0, errXorNotIntegerError
		}
		value ^= xorval
	}
	return value, nil
}
