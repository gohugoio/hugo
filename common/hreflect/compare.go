// Copyright 2026 The Hugo Authors. All rights reserved.
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
	"cmp"
	"math"
	"reflect"
)

// CompareNumbers compares a and b exactly, without first converting integers to float64.
// It returns -1, 0 or 1. ok is false if a or b is not a number or is NaN.
func CompareNumbers(a, b reflect.Value) (c int, ok bool) {
	ak, bk := a.Kind(), b.Kind()
	switch {
	case IsInt(ak):
		switch {
		case IsInt(bk):
			return cmp.Compare(a.Int(), b.Int()), true
		case IsUint(bk):
			return compareIntUint(a.Int(), b.Uint()), true
		case IsFloat(bk):
			return compareIntFloat(a.Int(), b.Float())
		}
	case IsUint(ak):
		switch {
		case IsInt(bk):
			return -compareIntUint(b.Int(), a.Uint()), true
		case IsUint(bk):
			return cmp.Compare(a.Uint(), b.Uint()), true
		case IsFloat(bk):
			return compareUintFloat(a.Uint(), b.Float())
		}
	case IsFloat(ak):
		switch {
		case IsInt(bk):
			c, ok := compareIntFloat(b.Int(), a.Float())
			return -c, ok
		case IsUint(bk):
			c, ok := compareUintFloat(b.Uint(), a.Float())
			return -c, ok
		case IsFloat(bk):
			af, bf := a.Float(), b.Float()
			if math.IsNaN(af) || math.IsNaN(bf) {
				return 0, false
			}
			return cmp.Compare(af, bf), true
		}
	}
	return 0, false
}

func compareIntUint(i int64, u uint64) int {
	if i < 0 {
		return -1
	}
	return cmp.Compare(uint64(i), u)
}

func compareIntFloat(i int64, f float64) (int, bool) {
	switch {
	case math.IsNaN(f):
		return 0, false
	case f >= 1<<63:
		return -1, true
	case f < -(1 << 63):
		return 1, true
	}
	if c := cmp.Compare(i, int64(f)); c != 0 {
		return c, true
	}
	return cmp.Compare(0, f-math.Trunc(f)), true
}

func compareUintFloat(u uint64, f float64) (int, bool) {
	switch {
	case math.IsNaN(f):
		return 0, false
	case f >= 1<<64:
		return -1, true
	case f < 0:
		return 1, true
	}
	if c := cmp.Compare(u, uint64(f)); c != 0 {
		return c, true
	}
	return cmp.Compare(0, f-math.Trunc(f)), true
}
