// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// According to https://github.com/golang/go/issues/20100, the Go stdlib will
// include math.Round beginning with Go 1.10.
//
// The following implementation was taken from https://golang.org/cl/43652.

package math

import "math"

const (
	mask  = 0x7FF
	shift = 64 - 11 - 1
	bias  = 1023
)

// Round returns the nearest integer, rounding half away from zero.
//
// Special cases are:
//
//	Round(±0) = ±0
//	Round(±Inf) = ±Inf
//	Round(NaN) = NaN
func _round(x float64) float64 {
	// Round is a faster implementation of:
	//
	//	func Round(x float64) float64 {
	//		t := Trunc(x)
	//		if Abs(x-t) >= 0.5 {
	//			return t + Copysign(1, x)
	//		}
	//		return t
	//	}
	const (
		signMask = 1 << 63
		fracMask = 1<<shift - 1
		half     = 1 << (shift - 1)
		one      = bias << shift
	)

	bits := math.Float64bits(x)
	e := uint(bits>>shift) & mask
	if e < bias {
		// Round abs(x) < 1 including denormals.
		bits &= signMask // +-0
		if e == bias-1 {
			bits |= one // +-1
		}
	} else if e < bias+shift {
		// Round any abs(x) >= 1 containing a fractional component [0,1).
		//
		// Numbers with larger exponents are returned unchanged since they
		// must be either an integer, infinity, or NaN.
		e -= bias
		bits += half >> e
		bits &^= fracMask >> e
	}
	return math.Float64frombits(bits)
}
