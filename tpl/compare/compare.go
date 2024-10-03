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

// Package compare provides template functions for comparing values.
package compare

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"time"

	"github.com/gohugoio/hugo/compare"
	"github.com/gohugoio/hugo/langs"

	"github.com/gohugoio/hugo/common/hreflect"
	"github.com/gohugoio/hugo/common/htime"
	"github.com/gohugoio/hugo/common/types"
)

// New returns a new instance of the compare-namespaced template functions.
func New(loc *time.Location, caseInsensitive bool) *Namespace {
	return &Namespace{loc: loc, caseInsensitive: caseInsensitive}
}

// Namespace provides template functions for the "compare" namespace.
type Namespace struct {
	loc *time.Location
	// Enable to do case insensitive string compares.
	caseInsensitive bool
}

// Default checks whether a givenv is set and returns the default value defaultv if it
// is not.  "Set" in this context means non-zero for numeric types and times;
// non-zero length for strings, arrays, slices, and maps;
// any boolean or struct value; or non-nil for any other types.
func (*Namespace) Default(defaultv any, givenv ...any) (any, error) {
	// given is variadic because the following construct will not pass a piped
	// argument when the key is missing:  {{ index . "key" | default "foo" }}
	// The Go template will complain that we got 1 argument when we expected 2.

	if len(givenv) == 0 {
		return defaultv, nil
	}
	if len(givenv) != 1 {
		return nil, fmt.Errorf("wrong number of args for default: want 2 got %d", len(givenv)+1)
	}

	g := reflect.ValueOf(givenv[0])
	if !g.IsValid() {
		return defaultv, nil
	}

	set := false

	switch g.Kind() {
	case reflect.Bool:
		set = true
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
		set = g.Len() != 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		set = g.Int() != 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		set = g.Uint() != 0
	case reflect.Float32, reflect.Float64:
		set = g.Float() != 0
	case reflect.Complex64, reflect.Complex128:
		set = g.Complex() != 0
	case reflect.Struct:
		switch actual := givenv[0].(type) {
		case time.Time:
			set = !actual.IsZero()
		default:
			set = true
		}
	default:
		set = !g.IsNil()
	}

	if set {
		return givenv[0], nil
	}

	return defaultv, nil
}

// Eq returns the boolean truth of arg1 == arg2 || arg1 == arg3 || arg1 == arg4.
func (n *Namespace) Eq(first any, others ...any) bool {
	if n.caseInsensitive {
		panic("caseInsensitive not implemented for Eq")
	}
	n.checkComparisonArgCount(1, others...)
	normalize := func(v any) any {
		if types.IsNil(v) {
			return nil
		}

		if at, ok := v.(htime.AsTimeProvider); ok {
			return at.AsTime(n.loc)
		}

		vv := reflect.ValueOf(v)
		switch vv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return vv.Int()
		case reflect.Float32, reflect.Float64:
			return vv.Float()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			i := vv.Uint()
			// If it can fit in an int, convert it.
			if i <= math.MaxInt64 {
				return int64(i)
			}
			return i
		case reflect.String:
			return vv.String()
		default:
			return v
		}
	}

	normFirst := normalize(first)
	for _, other := range others {
		if e, ok := first.(compare.Eqer); ok {
			if e.Eq(other) {
				return true
			}
			continue
		}

		if e, ok := other.(compare.Eqer); ok {
			if e.Eq(first) {
				return true
			}
			continue
		}

		other = normalize(other)
		if reflect.DeepEqual(normFirst, other) {
			return true
		}
	}

	return false
}

// Ne returns the boolean truth of arg1 != arg2 && arg1 != arg3 && arg1 != arg4.
func (n *Namespace) Ne(first any, others ...any) bool {
	n.checkComparisonArgCount(1, others...)
	for _, other := range others {
		if n.Eq(first, other) {
			return false
		}
	}
	return true
}

// Ge returns the boolean truth of arg1 >= arg2 && arg1 >= arg3 && arg1 >= arg4.
func (n *Namespace) Ge(first any, others ...any) bool {
	n.checkComparisonArgCount(1, others...)
	for _, other := range others {
		left, right := n.compareGet(first, other)
		if !(left >= right) {
			return false
		}
	}
	return true
}

// Gt returns the boolean truth of arg1 > arg2 && arg1 > arg3 && arg1 > arg4.
func (n *Namespace) Gt(first any, others ...any) bool {
	n.checkComparisonArgCount(1, others...)
	for _, other := range others {
		left, right := n.compareGet(first, other)
		if !(left > right) {
			return false
		}
	}
	return true
}

// Le returns the boolean truth of arg1 <= arg2 && arg1 <= arg3 && arg1 <= arg4.
func (n *Namespace) Le(first any, others ...any) bool {
	n.checkComparisonArgCount(1, others...)
	for _, other := range others {
		left, right := n.compareGet(first, other)
		if !(left <= right) {
			return false
		}
	}
	return true
}

// LtCollate returns the boolean truth of arg1 < arg2 && arg1 < arg3 && arg1 < arg4.
// The provided collator will be used for string comparisons.
// This is for internal use.
func (n *Namespace) LtCollate(collator *langs.Collator, first any, others ...any) bool {
	n.checkComparisonArgCount(1, others...)
	for _, other := range others {
		left, right := n.compareGetWithCollator(collator, first, other)
		if !(left < right) {
			return false
		}
	}
	return true
}

// Lt returns the boolean truth of arg1 < arg2 && arg1 < arg3 && arg1 < arg4.
func (n *Namespace) Lt(first any, others ...any) bool {
	return n.LtCollate(nil, first, others...)
}

func (n *Namespace) checkComparisonArgCount(min int, others ...any) bool {
	if len(others) < min {
		panic("missing arguments for comparison")
	}
	return true
}

// Conditional can be used as a ternary operator.
//
// It returns v1 if cond is true, else v2.
func (n *Namespace) Conditional(cond any, v1, v2 any) any {
	if hreflect.IsTruthful(cond) {
		return v1
	}
	return v2
}

func (ns *Namespace) compareGet(a any, b any) (float64, float64) {
	return ns.compareGetWithCollator(nil, a, b)
}

func (ns *Namespace) compareTwoUints(a uint64, b uint64) (float64, float64) {
	if a < b {
		return 1, 0
	} else if a == b {
		return 0, 0
	} else {
		return 0, 1
	}
}

func (ns *Namespace) compareGetWithCollator(collator *langs.Collator, a any, b any) (float64, float64) {
	if ac, ok := a.(compare.Comparer); ok {
		c := ac.Compare(b)
		if c < 0 {
			return 1, 0
		} else if c == 0 {
			return 0, 0
		} else {
			return 0, 1
		}
	}

	if bc, ok := b.(compare.Comparer); ok {
		c := bc.Compare(a)
		if c < 0 {
			return 0, 1
		} else if c == 0 {
			return 0, 0
		} else {
			return 1, 0
		}
	}

	var left, right float64
	var leftStr, rightStr *string
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)

	switch av.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		left = float64(av.Len())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if hreflect.IsUint(bv.Kind()) {
			return ns.compareTwoUints(uint64(av.Int()), bv.Uint())
		}
		left = float64(av.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		left = float64(av.Uint())
	case reflect.Uint64:
		if hreflect.IsUint(bv.Kind()) {
			return ns.compareTwoUints(av.Uint(), bv.Uint())
		}
	case reflect.Float32, reflect.Float64:
		left = av.Float()
	case reflect.String:
		var err error
		left, err = strconv.ParseFloat(av.String(), 64)
		// Check if float is a special floating value and cast value as string.
		if math.IsInf(left, 0) || math.IsNaN(left) || err != nil {
			str := av.String()
			leftStr = &str
		}
	case reflect.Struct:
		if hreflect.IsTime(av.Type()) {
			left = float64(ns.toTimeUnix(av))
		}
	case reflect.Bool:
		left = 0
		if av.Bool() {
			left = 1
		}
	}

	switch bv.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		right = float64(bv.Len())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if hreflect.IsUint(av.Kind()) {
			return ns.compareTwoUints(av.Uint(), uint64(bv.Int()))
		}
		right = float64(bv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		right = float64(bv.Uint())
	case reflect.Uint64:
		if hreflect.IsUint(av.Kind()) {
			return ns.compareTwoUints(av.Uint(), bv.Uint())
		}
	case reflect.Float32, reflect.Float64:
		right = bv.Float()
	case reflect.String:
		var err error
		right, err = strconv.ParseFloat(bv.String(), 64)
		// Check if float is a special floating value and cast value as string.
		if math.IsInf(right, 0) || math.IsNaN(right) || err != nil {
			str := bv.String()
			rightStr = &str
		}
	case reflect.Struct:
		if hreflect.IsTime(bv.Type()) {
			right = float64(ns.toTimeUnix(bv))
		}
	case reflect.Bool:
		right = 0
		if bv.Bool() {
			right = 1
		}
	}

	if (ns.caseInsensitive || collator != nil) && leftStr != nil && rightStr != nil {
		var c int
		if collator != nil {
			c = collator.CompareStrings(*leftStr, *rightStr)
		} else {
			c = compare.Strings(*leftStr, *rightStr)
		}
		if c < 0 {
			return 0, 1
		} else if c > 0 {
			return 1, 0
		} else {
			return 0, 0
		}
	}

	switch {
	case leftStr == nil || rightStr == nil:
	case *leftStr < *rightStr:
		return 0, 1
	case *leftStr > *rightStr:
		return 1, 0
	default:
		return 0, 0
	}

	return left, right
}

func (ns *Namespace) toTimeUnix(v reflect.Value) int64 {
	t, ok := hreflect.AsTime(v, ns.loc)
	if !ok {
		panic("coding error: argument must be time.Time type reflect Value")
	}
	return t.Unix()
}
