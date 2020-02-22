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
	"reflect"
	"strconv"
	"time"

	"github.com/gohugoio/hugo/compare"

	"github.com/gohugoio/hugo/common/types"
)

// New returns a new instance of the compare-namespaced template functions.
func New(caseInsensitive bool) *Namespace {
	return &Namespace{caseInsensitive: caseInsensitive}
}

// Namespace provides template functions for the "compare" namespace.
type Namespace struct {
	// Enable to do case insensitive string compares.
	caseInsensitive bool
}

// Default checks whether a given value is set and returns a default value if it
// is not.  "Set" in this context means non-zero for numeric types and times;
// non-zero length for strings, arrays, slices, and maps;
// any boolean or struct value; or non-nil for any other types.
func (*Namespace) Default(dflt interface{}, given ...interface{}) (interface{}, error) {
	// given is variadic because the following construct will not pass a piped
	// argument when the key is missing:  {{ index . "key" | default "foo" }}
	// The Go template will complain that we got 1 argument when we expectd 2.

	if len(given) == 0 {
		return dflt, nil
	}
	if len(given) != 1 {
		return nil, fmt.Errorf("wrong number of args for default: want 2 got %d", len(given)+1)
	}

	g := reflect.ValueOf(given[0])
	if !g.IsValid() {
		return dflt, nil
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
		switch actual := given[0].(type) {
		case time.Time:
			set = !actual.IsZero()
		default:
			set = true
		}
	default:
		set = !g.IsNil()
	}

	if set {
		return given[0], nil
	}

	return dflt, nil
}

// Eq returns the boolean truth of arg1 == arg2 || arg1 == arg3 || arg1 == arg4.
func (n *Namespace) Eq(first interface{}, others ...interface{}) bool {
	if n.caseInsensitive {
		panic("caseInsensitive not implemented for Eq")
	}
	if len(others) == 0 {
		panic("missing arguments for comparison")
	}

	normalize := func(v interface{}) interface{} {
		if types.IsNil(v) {
			return nil
		}
		vv := reflect.ValueOf(v)
		switch vv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return vv.Int()
		case reflect.Float32, reflect.Float64:
			return vv.Float()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return vv.Uint()
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
func (n *Namespace) Ne(first interface{}, others ...interface{}) bool {
	for _, other := range others {
		if n.Eq(first, other) {
			return false
		}
	}
	return true
}

// Ge returns the boolean truth of arg1 >= arg2 && arg1 >= arg3 && arg1 >= arg4.
func (n *Namespace) Ge(first interface{}, others ...interface{}) bool {
	for _, other := range others {
		left, right := n.compareGet(first, other)
		if !(left >= right) {
			return false
		}
	}
	return true
}

// Gt returns the boolean truth of arg1 > arg2 && arg1 > arg3 && arg1 > arg4.
func (n *Namespace) Gt(first interface{}, others ...interface{}) bool {
	for _, other := range others {
		left, right := n.compareGet(first, other)
		if !(left > right) {
			return false
		}
	}
	return true
}

// Le returns the boolean truth of arg1 <= arg2 && arg1 <= arg3 && arg1 <= arg4.
func (n *Namespace) Le(first interface{}, others ...interface{}) bool {
	for _, other := range others {
		left, right := n.compareGet(first, other)
		if !(left <= right) {
			return false
		}
	}
	return true
}

// Lt returns the boolean truth of arg1 < arg2 && arg1 < arg3 && arg1 < arg4.
func (n *Namespace) Lt(first interface{}, others ...interface{}) bool {
	for _, other := range others {
		left, right := n.compareGet(first, other)
		if !(left < right) {
			return false
		}
	}
	return true
}

// Conditional can be used as a ternary operator.
// It returns a if condition, else b.
func (n *Namespace) Conditional(condition bool, a, b interface{}) interface{} {
	if condition {
		return a
	}
	return b
}

func (ns *Namespace) compareGet(a interface{}, b interface{}) (float64, float64) {
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

	switch av.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		left = float64(av.Len())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		left = float64(av.Int())
	case reflect.Float32, reflect.Float64:
		left = av.Float()
	case reflect.String:
		var err error
		left, err = strconv.ParseFloat(av.String(), 64)
		if err != nil {
			str := av.String()
			leftStr = &str
		}
	case reflect.Struct:
		switch av.Type() {
		case timeType:
			left = float64(toTimeUnix(av))
		}
	}

	bv := reflect.ValueOf(b)

	switch bv.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		right = float64(bv.Len())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		right = float64(bv.Int())
	case reflect.Float32, reflect.Float64:
		right = bv.Float()
	case reflect.String:
		var err error
		right, err = strconv.ParseFloat(bv.String(), 64)
		if err != nil {
			str := bv.String()
			rightStr = &str
		}
	case reflect.Struct:
		switch bv.Type() {
		case timeType:
			right = float64(toTimeUnix(bv))
		}
	}

	if ns.caseInsensitive && leftStr != nil && rightStr != nil {
		c := compare.Strings(*leftStr, *rightStr)
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

var timeType = reflect.TypeOf((*time.Time)(nil)).Elem()

func toTimeUnix(v reflect.Value) int64 {
	if v.Kind() == reflect.Interface {
		return toTimeUnix(v.Elem())
	}
	if v.Type() != timeType {
		panic("coding error: argument must be time.Time type reflect Value")
	}
	return v.MethodByName("Unix").Call([]reflect.Value{})[0].Int()
}
