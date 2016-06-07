// Copyright 2016 The Hugo Authors. All rights reserved.
//
// Portions Copyright The Go Authors.

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

package tpl

import (
	"bytes"
	_md5 "crypto/md5"
	_sha1 "crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"html/template"
	"math/rand"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/spf13/afero"
	"github.com/spf13/hugo/hugofs"

	"github.com/bep/inflect"

	"github.com/spf13/cast"
	"github.com/spf13/hugo/helpers"
	jww "github.com/spf13/jwalterweatherman"
)

var funcMap template.FuncMap

// eq returns the boolean truth of arg1 == arg2.
func eq(x, y interface{}) bool {
	normalize := func(v interface{}) interface{} {
		vv := reflect.ValueOf(v)
		switch vv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return vv.Int()
		case reflect.Float32, reflect.Float64:
			return vv.Float()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return vv.Uint()
		default:
			return v
		}
	}
	x = normalize(x)
	y = normalize(y)
	return reflect.DeepEqual(x, y)
}

// ne returns the boolean truth of arg1 != arg2.
func ne(x, y interface{}) bool {
	return !eq(x, y)
}

// ge returns the boolean truth of arg1 >= arg2.
func ge(a, b interface{}) bool {
	left, right := compareGetFloat(a, b)
	return left >= right
}

// gt returns the boolean truth of arg1 > arg2.
func gt(a, b interface{}) bool {
	left, right := compareGetFloat(a, b)
	return left > right
}

// le returns the boolean truth of arg1 <= arg2.
func le(a, b interface{}) bool {
	left, right := compareGetFloat(a, b)
	return left <= right
}

// lt returns the boolean truth of arg1 < arg2.
func lt(a, b interface{}) bool {
	left, right := compareGetFloat(a, b)
	return left < right
}

// dictionary creates a map[string]interface{} from the given parameters by
// walking the parameters and treating them as key-value pairs.  The number
// of parameters must be even.
func dictionary(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dict call")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

// slice returns a slice of all passed arguments
func slice(args ...interface{}) []interface{} {
	return args
}

func compareGetFloat(a interface{}, b interface{}) (float64, float64) {
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

// slicestr slices a string by specifying a half-open range with
// two indices, start and end. 1 and 4 creates a slice including elements 1 through 3.
// The end index can be omitted, it defaults to the string's length.
func slicestr(a interface{}, startEnd ...interface{}) (string, error) {
	aStr, err := cast.ToStringE(a)
	if err != nil {
		return "", err
	}

	var argStart, argEnd int

	argNum := len(startEnd)

	if argNum > 0 {
		if argStart, err = cast.ToIntE(startEnd[0]); err != nil {
			return "", errors.New("start argument must be integer")
		}
	}
	if argNum > 1 {
		if argEnd, err = cast.ToIntE(startEnd[1]); err != nil {
			return "", errors.New("end argument must be integer")
		}
	}

	if argNum > 2 {
		return "", errors.New("too many arguments")
	}

	asRunes := []rune(aStr)

	if argNum > 0 && (argStart < 0 || argStart >= len(asRunes)) {
		return "", errors.New("slice bounds out of range")
	}

	if argNum == 2 {
		if argEnd < 0 || argEnd > len(asRunes) {
			return "", errors.New("slice bounds out of range")
		}
		return string(asRunes[argStart:argEnd]), nil
	} else if argNum == 1 {
		return string(asRunes[argStart:]), nil
	} else {
		return string(asRunes[:]), nil
	}

}

// substr extracts parts of a string, beginning at the character at the specified
// position, and returns the specified number of characters.
//
// It normally takes two parameters: start and length.
// It can also take one parameter: start, i.e. length is omitted, in which case
// the substring starting from start until the end of the string will be returned.
//
// To extract characters from the end of the string, use a negative start number.
//
// In addition, borrowing from the extended behavior described at http://php.net/substr,
// if length is given and is negative, then that many characters will be omitted from
// the end of string.
func substr(a interface{}, nums ...interface{}) (string, error) {
	aStr, err := cast.ToStringE(a)
	if err != nil {
		return "", err
	}

	var start, length int

	asRunes := []rune(aStr)

	switch len(nums) {
	case 0:
		return "", errors.New("too less arguments")
	case 1:
		if start, err = cast.ToIntE(nums[0]); err != nil {
			return "", errors.New("start argument must be integer")
		}
		length = len(asRunes)
	case 2:
		if start, err = cast.ToIntE(nums[0]); err != nil {
			return "", errors.New("start argument must be integer")
		}
		if length, err = cast.ToIntE(nums[1]); err != nil {
			return "", errors.New("length argument must be integer")
		}
	default:
		return "", errors.New("too many arguments")
	}

	if start < -len(asRunes) {
		start = 0
	}
	if start > len(asRunes) {
		return "", fmt.Errorf("start position out of bounds for %d-byte string", len(aStr))
	}

	var s, e int
	if start >= 0 && length >= 0 {
		s = start
		e = start + length
	} else if start < 0 && length >= 0 {
		s = len(asRunes) + start - length + 1
		e = len(asRunes) + start + 1
	} else if start >= 0 && length < 0 {
		s = start
		e = len(asRunes) + length
	} else {
		s = len(asRunes) + start
		e = len(asRunes) + length
	}

	if s > e {
		return "", fmt.Errorf("calculated start position greater than end position: %d > %d", s, e)
	}
	if e > len(asRunes) {
		e = len(asRunes)
	}

	return string(asRunes[s:e]), nil
}

// split slices an input string into all substrings separated by delimiter.
func split(a interface{}, delimiter string) ([]string, error) {
	aStr, err := cast.ToStringE(a)
	if err != nil {
		return []string{}, err
	}
	return strings.Split(aStr, delimiter), nil
}

// intersect returns the common elements in the given sets, l1 and l2.  l1 and
// l2 must be of the same type and may be either arrays or slices.
func intersect(l1, l2 interface{}) (interface{}, error) {
	if l1 == nil || l2 == nil {
		return make([]interface{}, 0), nil
	}

	l1v := reflect.ValueOf(l1)
	l2v := reflect.ValueOf(l2)

	switch l1v.Kind() {
	case reflect.Array, reflect.Slice:
		switch l2v.Kind() {
		case reflect.Array, reflect.Slice:
			r := reflect.MakeSlice(l1v.Type(), 0, 0)
			for i := 0; i < l1v.Len(); i++ {
				l1vv := l1v.Index(i)
				for j := 0; j < l2v.Len(); j++ {
					l2vv := l2v.Index(j)
					switch l1vv.Kind() {
					case reflect.String:
						if l1vv.Type() == l2vv.Type() && l1vv.String() == l2vv.String() && !in(r.Interface(), l2vv.Interface()) {
							r = reflect.Append(r, l2vv)
						}
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						switch l2vv.Kind() {
						case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
							if l1vv.Int() == l2vv.Int() && !in(r.Interface(), l2vv.Interface()) {
								r = reflect.Append(r, l2vv)
							}
						}
					case reflect.Float32, reflect.Float64:
						switch l2vv.Kind() {
						case reflect.Float32, reflect.Float64:
							if l1vv.Float() == l2vv.Float() && !in(r.Interface(), l2vv.Interface()) {
								r = reflect.Append(r, l2vv)
							}
						}
					}
				}
			}
			return r.Interface(), nil
		default:
			return nil, errors.New("can't iterate over " + reflect.ValueOf(l2).Type().String())
		}
	default:
		return nil, errors.New("can't iterate over " + reflect.ValueOf(l1).Type().String())
	}
}

// in returns whether v is in the set l.  l may be an array or slice.
func in(l interface{}, v interface{}) bool {
	lv := reflect.ValueOf(l)
	vv := reflect.ValueOf(v)

	switch lv.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < lv.Len(); i++ {
			lvv := lv.Index(i)
			lvv, isNil := indirect(lvv)
			if isNil {
				continue
			}
			switch lvv.Kind() {
			case reflect.String:
				if vv.Type() == lvv.Type() && vv.String() == lvv.String() {
					return true
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				switch vv.Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					if vv.Int() == lvv.Int() {
						return true
					}
				}
			case reflect.Float32, reflect.Float64:
				switch vv.Kind() {
				case reflect.Float32, reflect.Float64:
					if vv.Float() == lvv.Float() {
						return true
					}
				}
			}
		}
	case reflect.String:
		if vv.Type() == lv.Type() && strings.Contains(lv.String(), vv.String()) {
			return true
		}
	}
	return false
}

// first returns the first N items in a rangeable list.
func first(limit interface{}, seq interface{}) (interface{}, error) {
	if limit == nil || seq == nil {
		return nil, errors.New("both limit and seq must be provided")
	}

	limitv, err := cast.ToIntE(limit)

	if err != nil {
		return nil, err
	}

	if limitv < 1 {
		return nil, errors.New("can't return negative/empty count of items from sequence")
	}

	seqv := reflect.ValueOf(seq)
	seqv, isNil := indirect(seqv)
	if isNil {
		return nil, errors.New("can't iterate over a nil value")
	}

	switch seqv.Kind() {
	case reflect.Array, reflect.Slice, reflect.String:
		// okay
	default:
		return nil, errors.New("can't iterate over " + reflect.ValueOf(seq).Type().String())
	}
	if limitv > seqv.Len() {
		limitv = seqv.Len()
	}
	return seqv.Slice(0, limitv).Interface(), nil
}

// findRE returns a list of strings that match the regular expression. By default all matches
// will be included. The number of matches can be limitted with an optional third parameter.
func findRE(expr string, content interface{}, limit ...int) ([]string, error) {
	re, err := reCache.Get(expr)
	if err != nil {
		return nil, err
	}

	conv := cast.ToString(content)
	if len(limit) > 0 {
		return re.FindAllString(conv, limit[0]), nil
	}

	return re.FindAllString(conv, -1), nil
}

// last returns the last N items in a rangeable list.
func last(limit interface{}, seq interface{}) (interface{}, error) {
	if limit == nil || seq == nil {
		return nil, errors.New("both limit and seq must be provided")
	}

	limitv, err := cast.ToIntE(limit)

	if err != nil {
		return nil, err
	}

	if limitv < 1 {
		return nil, errors.New("can't return negative/empty count of items from sequence")
	}

	seqv := reflect.ValueOf(seq)
	seqv, isNil := indirect(seqv)
	if isNil {
		return nil, errors.New("can't iterate over a nil value")
	}

	switch seqv.Kind() {
	case reflect.Array, reflect.Slice, reflect.String:
		// okay
	default:
		return nil, errors.New("can't iterate over " + reflect.ValueOf(seq).Type().String())
	}
	if limitv > seqv.Len() {
		limitv = seqv.Len()
	}
	return seqv.Slice(seqv.Len()-limitv, seqv.Len()).Interface(), nil
}

// after returns all the items after the first N in a rangeable list.
func after(index interface{}, seq interface{}) (interface{}, error) {
	if index == nil || seq == nil {
		return nil, errors.New("both limit and seq must be provided")
	}

	indexv, err := cast.ToIntE(index)

	if err != nil {
		return nil, err
	}

	if indexv < 1 {
		return nil, errors.New("can't return negative/empty count of items from sequence")
	}

	seqv := reflect.ValueOf(seq)
	seqv, isNil := indirect(seqv)
	if isNil {
		return nil, errors.New("can't iterate over a nil value")
	}

	switch seqv.Kind() {
	case reflect.Array, reflect.Slice, reflect.String:
		// okay
	default:
		return nil, errors.New("can't iterate over " + reflect.ValueOf(seq).Type().String())
	}
	if indexv >= seqv.Len() {
		return nil, errors.New("no items left")
	}
	return seqv.Slice(indexv, seqv.Len()).Interface(), nil
}

// shuffle returns the given rangeable list in a randomised order.
func shuffle(seq interface{}) (interface{}, error) {
	if seq == nil {
		return nil, errors.New("both count and seq must be provided")
	}

	seqv := reflect.ValueOf(seq)
	seqv, isNil := indirect(seqv)
	if isNil {
		return nil, errors.New("can't iterate over a nil value")
	}

	switch seqv.Kind() {
	case reflect.Array, reflect.Slice, reflect.String:
		// okay
	default:
		return nil, errors.New("can't iterate over " + reflect.ValueOf(seq).Type().String())
	}

	shuffled := reflect.MakeSlice(reflect.TypeOf(seq), seqv.Len(), seqv.Len())

	rand.Seed(time.Now().UTC().UnixNano())
	randomIndices := rand.Perm(seqv.Len())

	for index, value := range randomIndices {
		shuffled.Index(value).Set(seqv.Index(index))
	}

	return shuffled.Interface(), nil
}

func evaluateSubElem(obj reflect.Value, elemName string) (reflect.Value, error) {
	if !obj.IsValid() {
		return zero, errors.New("can't evaluate an invalid value")
	}
	typ := obj.Type()
	obj, isNil := indirect(obj)

	// first, check whether obj has a method. In this case, obj is
	// an interface, a struct or its pointer. If obj is a struct,
	// to check all T and *T method, use obj pointer type Value
	objPtr := obj
	if objPtr.Kind() != reflect.Interface && objPtr.CanAddr() {
		objPtr = objPtr.Addr()
	}
	mt, ok := objPtr.Type().MethodByName(elemName)
	if ok {
		if mt.PkgPath != "" {
			return zero, fmt.Errorf("%s is an unexported method of type %s", elemName, typ)
		}
		// struct pointer has one receiver argument and interface doesn't have an argument
		if mt.Type.NumIn() > 1 || mt.Type.NumOut() == 0 || mt.Type.NumOut() > 2 {
			return zero, fmt.Errorf("%s is a method of type %s but doesn't satisfy requirements", elemName, typ)
		}
		if mt.Type.NumOut() == 1 && mt.Type.Out(0).Implements(errorType) {
			return zero, fmt.Errorf("%s is a method of type %s but doesn't satisfy requirements", elemName, typ)
		}
		if mt.Type.NumOut() == 2 && !mt.Type.Out(1).Implements(errorType) {
			return zero, fmt.Errorf("%s is a method of type %s but doesn't satisfy requirements", elemName, typ)
		}
		res := objPtr.Method(mt.Index).Call([]reflect.Value{})
		if len(res) == 2 && !res[1].IsNil() {
			return zero, fmt.Errorf("error at calling a method %s of type %s: %s", elemName, typ, res[1].Interface().(error))
		}
		return res[0], nil
	}

	// elemName isn't a method so next start to check whether it is
	// a struct field or a map value. In both cases, it mustn't be
	// a nil value
	if isNil {
		return zero, fmt.Errorf("can't evaluate a nil pointer of type %s by a struct field or map key name %s", typ, elemName)
	}
	switch obj.Kind() {
	case reflect.Struct:
		ft, ok := obj.Type().FieldByName(elemName)
		if ok {
			if ft.PkgPath != "" && !ft.Anonymous {
				return zero, fmt.Errorf("%s is an unexported field of struct type %s", elemName, typ)
			}
			return obj.FieldByIndex(ft.Index), nil
		}
		return zero, fmt.Errorf("%s isn't a field of struct type %s", elemName, typ)
	case reflect.Map:
		kv := reflect.ValueOf(elemName)
		if kv.Type().AssignableTo(obj.Type().Key()) {
			return obj.MapIndex(kv), nil
		}
		return zero, fmt.Errorf("%s isn't a key of map type %s", elemName, typ)
	}
	return zero, fmt.Errorf("%s is neither a struct field, a method nor a map element of type %s", elemName, typ)
}

func checkCondition(v, mv reflect.Value, op string) (bool, error) {
	v, vIsNil := indirect(v)
	if !v.IsValid() {
		vIsNil = true
	}
	mv, mvIsNil := indirect(mv)
	if !mv.IsValid() {
		mvIsNil = true
	}
	if vIsNil || mvIsNil {
		switch op {
		case "", "=", "==", "eq":
			return vIsNil == mvIsNil, nil
		case "!=", "<>", "ne":
			return vIsNil != mvIsNil, nil
		}
		return false, nil
	}

	if v.Kind() == reflect.Bool && mv.Kind() == reflect.Bool {
		switch op {
		case "", "=", "==", "eq":
			return v.Bool() == mv.Bool(), nil
		case "!=", "<>", "ne":
			return v.Bool() != mv.Bool(), nil
		}
		return false, nil
	}

	var ivp, imvp *int64
	var svp, smvp *string
	var slv, slmv interface{}
	var ima []int64
	var sma []string
	if mv.Type() == v.Type() {
		switch v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			iv := v.Int()
			ivp = &iv
			imv := mv.Int()
			imvp = &imv
		case reflect.String:
			sv := v.String()
			svp = &sv
			smv := mv.String()
			smvp = &smv
		case reflect.Struct:
			switch v.Type() {
			case timeType:
				iv := toTimeUnix(v)
				ivp = &iv
				imv := toTimeUnix(mv)
				imvp = &imv
			}
		case reflect.Array, reflect.Slice:
			slv = v.Interface()
			slmv = mv.Interface()
		}
	} else {
		if mv.Kind() != reflect.Array && mv.Kind() != reflect.Slice {
			return false, nil
		}

		if mv.Len() == 0 {
			return false, nil
		}

		if v.Kind() != reflect.Interface && mv.Type().Elem().Kind() != reflect.Interface && mv.Type().Elem() != v.Type() {
			return false, nil
		}
		switch v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			iv := v.Int()
			ivp = &iv
			for i := 0; i < mv.Len(); i++ {
				if anInt := toInt(mv.Index(i)); anInt != -1 {
					ima = append(ima, anInt)
				}

			}
		case reflect.String:
			sv := v.String()
			svp = &sv
			for i := 0; i < mv.Len(); i++ {
				if aString := toString(mv.Index(i)); aString != "" {
					sma = append(sma, aString)
				}
			}
		case reflect.Struct:
			switch v.Type() {
			case timeType:
				iv := toTimeUnix(v)
				ivp = &iv
				for i := 0; i < mv.Len(); i++ {
					ima = append(ima, toTimeUnix(mv.Index(i)))
				}
			}
		}
	}

	switch op {
	case "", "=", "==", "eq":
		if ivp != nil && imvp != nil {
			return *ivp == *imvp, nil
		} else if svp != nil && smvp != nil {
			return *svp == *smvp, nil
		}
	case "!=", "<>", "ne":
		if ivp != nil && imvp != nil {
			return *ivp != *imvp, nil
		} else if svp != nil && smvp != nil {
			return *svp != *smvp, nil
		}
	case ">=", "ge":
		if ivp != nil && imvp != nil {
			return *ivp >= *imvp, nil
		} else if svp != nil && smvp != nil {
			return *svp >= *smvp, nil
		}
	case ">", "gt":
		if ivp != nil && imvp != nil {
			return *ivp > *imvp, nil
		} else if svp != nil && smvp != nil {
			return *svp > *smvp, nil
		}
	case "<=", "le":
		if ivp != nil && imvp != nil {
			return *ivp <= *imvp, nil
		} else if svp != nil && smvp != nil {
			return *svp <= *smvp, nil
		}
	case "<", "lt":
		if ivp != nil && imvp != nil {
			return *ivp < *imvp, nil
		} else if svp != nil && smvp != nil {
			return *svp < *smvp, nil
		}
	case "in", "not in":
		var r bool
		if ivp != nil && len(ima) > 0 {
			r = in(ima, *ivp)
		} else if svp != nil {
			if len(sma) > 0 {
				r = in(sma, *svp)
			} else if smvp != nil {
				r = in(*smvp, *svp)
			}
		} else {
			return false, nil
		}
		if op == "not in" {
			return !r, nil
		}
		return r, nil
	case "intersect":
		r, err := intersect(slv, slmv)
		if err != nil {
			return false, err
		}

		if reflect.TypeOf(r).Kind() == reflect.Slice {
			s := reflect.ValueOf(r)

			if s.Len() > 0 {
				return true, nil
			}
			return false, nil
		} else {
			return false, errors.New("invalid intersect values")
		}
	default:
		return false, errors.New("no such operator")
	}
	return false, nil
}

// where returns a filtered subset of a given data type.
func where(seq, key interface{}, args ...interface{}) (r interface{}, err error) {
	seqv := reflect.ValueOf(seq)
	kv := reflect.ValueOf(key)

	var mv reflect.Value
	var op string
	switch len(args) {
	case 1:
		mv = reflect.ValueOf(args[0])
	case 2:
		var ok bool
		if op, ok = args[0].(string); !ok {
			return nil, errors.New("operator argument must be string type")
		}
		op = strings.TrimSpace(strings.ToLower(op))
		mv = reflect.ValueOf(args[1])
	default:
		return nil, errors.New("can't evaluate the array by no match argument or more than or equal to two arguments")
	}

	seqv, isNil := indirect(seqv)
	if isNil {
		return nil, errors.New("can't iterate over a nil value of type " + reflect.ValueOf(seq).Type().String())
	}

	var path []string
	if kv.Kind() == reflect.String {
		path = strings.Split(strings.Trim(kv.String(), "."), ".")
	}

	switch seqv.Kind() {
	case reflect.Array, reflect.Slice:
		rv := reflect.MakeSlice(seqv.Type(), 0, 0)
		for i := 0; i < seqv.Len(); i++ {
			var vvv reflect.Value
			rvv := seqv.Index(i)
			if kv.Kind() == reflect.String {
				vvv = rvv
				for _, elemName := range path {
					vvv, err = evaluateSubElem(vvv, elemName)
					if err != nil {
						return nil, err
					}
				}
			} else {
				vv, _ := indirect(rvv)
				if vv.Kind() == reflect.Map && kv.Type().AssignableTo(vv.Type().Key()) {
					vvv = vv.MapIndex(kv)
				}
			}
			if ok, err := checkCondition(vvv, mv, op); ok {
				rv = reflect.Append(rv, rvv)
			} else if err != nil {
				return nil, err
			}
		}
		return rv.Interface(), nil
	default:
		return nil, fmt.Errorf("can't iterate over %v", seq)
	}
}

// apply takes a map, array, or slice and returns a new slice with the function fname applied over it.
func apply(seq interface{}, fname string, args ...interface{}) (interface{}, error) {
	if seq == nil {
		return make([]interface{}, 0), nil
	}

	if fname == "apply" {
		return nil, errors.New("can't apply myself (no turtles allowed)")
	}

	seqv := reflect.ValueOf(seq)
	seqv, isNil := indirect(seqv)
	if isNil {
		return nil, errors.New("can't iterate over a nil value")
	}

	fn, found := funcMap[fname]
	if !found {
		return nil, errors.New("can't find function " + fname)
	}

	fnv := reflect.ValueOf(fn)

	switch seqv.Kind() {
	case reflect.Array, reflect.Slice:
		r := make([]interface{}, seqv.Len())
		for i := 0; i < seqv.Len(); i++ {
			vv := seqv.Index(i)

			vvv, err := applyFnToThis(fnv, vv, args...)

			if err != nil {
				return nil, err
			}

			r[i] = vvv.Interface()
		}

		return r, nil
	default:
		return nil, fmt.Errorf("can't apply over %v", seq)
	}
}

func applyFnToThis(fn, this reflect.Value, args ...interface{}) (reflect.Value, error) {
	n := make([]reflect.Value, len(args))
	for i, arg := range args {
		if arg == "." {
			n[i] = this
		} else {
			n[i] = reflect.ValueOf(arg)
		}
	}

	num := fn.Type().NumIn()

	if fn.Type().IsVariadic() {
		num--
	}

	// TODO(bep) see #1098 - also see template_tests.go
	/*if len(args) < num {
		return reflect.ValueOf(nil), errors.New("Too few arguments")
	} else if len(args) > num {
		return reflect.ValueOf(nil), errors.New("Too many arguments")
	}*/

	for i := 0; i < num; i++ {
		if xt, targ := n[i].Type(), fn.Type().In(i); !xt.AssignableTo(targ) {
			return reflect.ValueOf(nil), errors.New("called apply using " + xt.String() + " as type " + targ.String())
		}
	}

	res := fn.Call(n)

	if len(res) == 1 || res[1].IsNil() {
		return res[0], nil
	}
	return reflect.ValueOf(nil), res[1].Interface().(error)
}

// delimit takes a given sequence and returns a delimited HTML string.
// If last is passed to the function, it will be used as the final delimiter.
func delimit(seq, delimiter interface{}, last ...interface{}) (template.HTML, error) {
	d, err := cast.ToStringE(delimiter)
	if err != nil {
		return "", err
	}

	var dLast *string
	for _, l := range last {
		dStr, err := cast.ToStringE(l)
		if err != nil {
			dLast = nil
		}
		dLast = &dStr
		break
	}

	seqv := reflect.ValueOf(seq)
	seqv, isNil := indirect(seqv)
	if isNil {
		return "", errors.New("can't iterate over a nil value")
	}

	var str string
	switch seqv.Kind() {
	case reflect.Map:
		sortSeq, err := sortSeq(seq)
		if err != nil {
			return "", err
		}
		seqv = reflect.ValueOf(sortSeq)
		fallthrough
	case reflect.Array, reflect.Slice, reflect.String:
		for i := 0; i < seqv.Len(); i++ {
			val := seqv.Index(i).Interface()
			valStr, err := cast.ToStringE(val)
			if err != nil {
				continue
			}
			switch {
			case i == seqv.Len()-2 && dLast != nil:
				str += valStr + *dLast
			case i == seqv.Len()-1:
				str += valStr
			default:
				str += valStr + d
			}
		}

	default:
		return "", fmt.Errorf("can't iterate over %v", seq)
	}

	return template.HTML(str), nil
}

// sortSeq returns a sorted sequence.
func sortSeq(seq interface{}, args ...interface{}) (interface{}, error) {
	if seq == nil {
		return nil, errors.New("sequence must be provided")
	}

	seqv := reflect.ValueOf(seq)
	seqv, isNil := indirect(seqv)
	if isNil {
		return nil, errors.New("can't iterate over a nil value")
	}

	switch seqv.Kind() {
	case reflect.Array, reflect.Slice, reflect.Map:
		// ok
	default:
		return nil, errors.New("can't sort " + reflect.ValueOf(seq).Type().String())
	}

	// Create a list of pairs that will be used to do the sort
	p := pairList{SortAsc: true, SliceType: reflect.SliceOf(seqv.Type().Elem())}
	p.Pairs = make([]pair, seqv.Len())

	var sortByField string
	for i, l := range args {
		dStr, err := cast.ToStringE(l)
		switch {
		case i == 0 && err != nil:
			sortByField = ""
		case i == 0 && err == nil:
			sortByField = dStr
		case i == 1 && err == nil && dStr == "desc":
			p.SortAsc = false
		case i == 1:
			p.SortAsc = true
		}
	}
	path := strings.Split(strings.Trim(sortByField, "."), ".")

	switch seqv.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < seqv.Len(); i++ {
			p.Pairs[i].Key = reflect.ValueOf(i)
			p.Pairs[i].Value = seqv.Index(i)
			if sortByField == "" || sortByField == "value" {
				p.Pairs[i].SortByValue = p.Pairs[i].Value
			} else {
				v := p.Pairs[i].Value
				var err error
				for _, elemName := range path {
					v, err = evaluateSubElem(v, elemName)
					if err != nil {
						return nil, err
					}
				}
				p.Pairs[i].SortByValue = v
			}
		}

	case reflect.Map:
		keys := seqv.MapKeys()
		for i := 0; i < seqv.Len(); i++ {
			p.Pairs[i].Key = keys[i]
			p.Pairs[i].Value = seqv.MapIndex(keys[i])
			if sortByField == "" {
				p.Pairs[i].SortByValue = p.Pairs[i].Key
			} else if sortByField == "value" {
				p.Pairs[i].SortByValue = p.Pairs[i].Value
			} else {
				v := p.Pairs[i].Value
				var err error
				for _, elemName := range path {
					v, err = evaluateSubElem(v, elemName)
					if err != nil {
						return nil, err
					}
				}
				p.Pairs[i].SortByValue = v
			}
		}
	}
	return p.sort(), nil
}

// Credit for pair sorting method goes to Andrew Gerrand
// https://groups.google.com/forum/#!topic/golang-nuts/FT7cjmcL7gw
// A data structure to hold a key/value pair.
type pair struct {
	Key         reflect.Value
	Value       reflect.Value
	SortByValue reflect.Value
}

// A slice of pairs that implements sort.Interface to sort by Value.
type pairList struct {
	Pairs     []pair
	SortAsc   bool
	SliceType reflect.Type
}

func (p pairList) Swap(i, j int) { p.Pairs[i], p.Pairs[j] = p.Pairs[j], p.Pairs[i] }
func (p pairList) Len() int      { return len(p.Pairs) }
func (p pairList) Less(i, j int) bool {
	iv := p.Pairs[i].SortByValue
	jv := p.Pairs[j].SortByValue

	if iv.IsValid() {
		if jv.IsValid() {
			// can only call Interface() on valid reflect Values
			return lt(iv.Interface(), jv.Interface())
		}
		// if j is invalid, test i against i's zero value
		return lt(iv.Interface(), reflect.Zero(iv.Type()))
	}

	if jv.IsValid() {
		// if i is invalid, test j against j's zero value
		return lt(reflect.Zero(jv.Type()), jv.Interface())
	}

	return false
}

// sorts a pairList and returns a slice of sorted values
func (p pairList) sort() interface{} {
	if p.SortAsc {
		sort.Sort(p)
	} else {
		sort.Sort(sort.Reverse(p))
	}
	sorted := reflect.MakeSlice(p.SliceType, len(p.Pairs), len(p.Pairs))
	for i, v := range p.Pairs {
		sorted.Index(i).Set(v.Value)
	}

	return sorted.Interface()
}

// isSet returns whether a given array, channel, slice, or map has a key
// defined.
func isSet(a interface{}, key interface{}) bool {
	av := reflect.ValueOf(a)
	kv := reflect.ValueOf(key)

	switch av.Kind() {
	case reflect.Array, reflect.Chan, reflect.Slice:
		if int64(av.Len()) > kv.Int() {
			return true
		}
	case reflect.Map:
		if kv.Type() == av.Type().Key() {
			return av.MapIndex(kv).IsValid()
		}
	}

	return false
}

// returnWhenSet returns a given value if it set.  Otherwise, it returns an
// empty string.
func returnWhenSet(a, k interface{}) interface{} {
	av, isNil := indirect(reflect.ValueOf(a))
	if isNil {
		return ""
	}

	var avv reflect.Value
	switch av.Kind() {
	case reflect.Array, reflect.Slice:
		index, ok := k.(int)
		if ok && av.Len() > index {
			avv = av.Index(index)
		}
	case reflect.Map:
		kv := reflect.ValueOf(k)
		if kv.Type().AssignableTo(av.Type().Key()) {
			avv = av.MapIndex(kv)
		}
	}

	if avv.IsValid() {
		switch avv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return avv.Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return avv.Uint()
		case reflect.Float32, reflect.Float64:
			return avv.Float()
		case reflect.String:
			return avv.String()
		}
	}

	return ""
}

// highlight returns an HTML string with syntax highlighting applied.
func highlight(in interface{}, lang, opts string) (template.HTML, error) {
	str, err := cast.ToStringE(in)

	if err != nil {
		return "", err
	}

	return template.HTML(helpers.Highlight(html.UnescapeString(str), lang, opts)), nil
}

var markdownTrimPrefix = []byte("<p>")
var markdownTrimSuffix = []byte("</p>\n")

// markdownify renders a given string from Markdown to HTML.
func markdownify(in interface{}) template.HTML {
	text := cast.ToString(in)
	m := helpers.RenderBytes(&helpers.RenderingContext{Content: []byte(text), PageFmt: "markdown"})
	m = bytes.TrimPrefix(m, markdownTrimPrefix)
	m = bytes.TrimSuffix(m, markdownTrimSuffix)
	return template.HTML(m)
}

// jsonify encodes a given object to JSON.
func jsonify(v interface{}) (template.HTML, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return template.HTML(b), nil

}

// emojify "emojifies" the given string.
//
// See http://www.emoji-cheat-sheet.com/
func emojify(in interface{}) (template.HTML, error) {
	str, err := cast.ToStringE(in)

	if err != nil {
		return "", err
	}

	return template.HTML(helpers.Emojify([]byte(str))), nil
}

// plainify strips any HTML and returns the plain text version.
func plainify(in interface{}) (string, error) {
	s, err := cast.ToStringE(in)

	if err != nil {
		return "", err
	}

	return helpers.StripHTML(s), nil
}

func refPage(page interface{}, ref, methodName string) template.HTML {
	value := reflect.ValueOf(page)

	method := value.MethodByName(methodName)

	if method.IsValid() && method.Type().NumIn() == 1 && method.Type().NumOut() == 2 {
		result := method.Call([]reflect.Value{reflect.ValueOf(ref)})

		url, err := result[0], result[1]

		if !err.IsNil() {
			jww.ERROR.Printf("%s", err.Interface())
			return template.HTML(fmt.Sprintf("%s", err.Interface()))
		}

		if url.String() == "" {
			jww.ERROR.Printf("ref %s could not be found\n", ref)
			return template.HTML(ref)
		}

		return template.HTML(url.String())
	}

	jww.ERROR.Printf("Can only create references from Page and Node objects.")
	return template.HTML(ref)
}

// ref returns the absolute URL path to a given content item.
func ref(page interface{}, ref string) template.HTML {
	return refPage(page, ref, "Ref")
}

// relRef returns the relative URL path to a given content item.
func relRef(page interface{}, ref string) template.HTML {
	return refPage(page, ref, "RelRef")
}

// chomp removes trailing newline characters from a string.
func chomp(text interface{}) (template.HTML, error) {
	s, err := cast.ToStringE(text)
	if err != nil {
		return "", err
	}

	return template.HTML(strings.TrimRight(s, "\r\n")), nil
}

// trim leading/trailing characters defined by b from a
func trim(a interface{}, b string) (string, error) {
	aStr, err := cast.ToStringE(a)
	if err != nil {
		return "", err
	}
	return strings.Trim(aStr, b), nil
}

// replace all occurrences of b with c in a
func replace(a, b, c interface{}) (string, error) {
	aStr, err := cast.ToStringE(a)
	if err != nil {
		return "", err
	}
	bStr, err := cast.ToStringE(b)
	if err != nil {
		return "", err
	}
	cStr, err := cast.ToStringE(c)
	if err != nil {
		return "", err
	}
	return strings.Replace(aStr, bStr, cStr, -1), nil
}

// regexpCache represents a cache of regexp objects protected by a mutex.
type regexpCache struct {
	mu sync.RWMutex
	re map[string]*regexp.Regexp
}

// Get retrieves a regexp object from the cache based upon the pattern.
// If the pattern is not found in the cache, create one
func (rc *regexpCache) Get(pattern string) (re *regexp.Regexp, err error) {
	var ok bool

	if re, ok = rc.get(pattern); !ok {
		re, err = regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}
		rc.set(pattern, re)
	}

	return re, nil
}

func (rc *regexpCache) get(key string) (re *regexp.Regexp, ok bool) {
	rc.mu.RLock()
	re, ok = rc.re[key]
	rc.mu.RUnlock()
	return
}

func (rc *regexpCache) set(key string, re *regexp.Regexp) {
	rc.mu.Lock()
	rc.re[key] = re
	rc.mu.Unlock()
}

var reCache = regexpCache{re: make(map[string]*regexp.Regexp)}

// replaceRE exposes a regular expression replacement function to the templates.
func replaceRE(pattern, repl, src interface{}) (_ string, err error) {
	patternStr, err := cast.ToStringE(pattern)
	if err != nil {
		return
	}

	replStr, err := cast.ToStringE(repl)
	if err != nil {
		return
	}

	srcStr, err := cast.ToStringE(src)
	if err != nil {
		return
	}

	re, err := reCache.Get(patternStr)
	if err != nil {
		return "", err
	}
	return re.ReplaceAllString(srcStr, replStr), nil
}

// dateFormat converts the textual representation of the datetime string into
// the other form or returns it of the time.Time value. These are formatted
// with the layout string
func dateFormat(layout string, v interface{}) (string, error) {
	t, err := cast.ToTimeE(v)
	if err != nil {
		return "", err
	}
	return t.Format(layout), nil
}

// dfault checks whether a given value is set and returns a default value if it
// is not.  "Set" in this context means non-zero for numeric types and times;
// non-zero length for strings, arrays, slices, and maps;
// any boolean or struct value; or non-nil for any other types.
func dfault(dflt interface{}, given ...interface{}) (interface{}, error) {
	// given is variadic because the following construct will not pass a piped
	// argument when the key is missing:  {{ index . "key" | default "foo" }}
	// The Go template will complain that we got 1 argument when we expectd 2.

	if given == nil || len(given) == 0 {
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

// canBeNil reports whether an untyped nil can be assigned to the type. See reflect.Zero.
//
// Copied from Go stdlib src/text/template/exec.go.
func canBeNil(typ reflect.Type) bool {
	switch typ.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return true
	}
	return false
}

// prepareArg checks if value can be used as an argument of type argType, and
// converts an invalid value to appropriate zero if possible.
//
// Copied from Go stdlib src/text/template/funcs.go.
func prepareArg(value reflect.Value, argType reflect.Type) (reflect.Value, error) {
	if !value.IsValid() {
		if !canBeNil(argType) {
			return reflect.Value{}, fmt.Errorf("value is nil; should be of type %s", argType)
		}
		value = reflect.Zero(argType)
	}
	if !value.Type().AssignableTo(argType) {
		return reflect.Value{}, fmt.Errorf("value has type %s; should be %s", value.Type(), argType)
	}
	return value, nil
}

// index returns the result of indexing its first argument by the following
// arguments. Thus "index x 1 2 3" is, in Go syntax, x[1][2][3]. Each
// indexed item must be a map, slice, or array.
//
// Copied from Go stdlib src/text/template/funcs.go.
// Can hopefully be removed in Go 1.7, see https://github.com/golang/go/issues/14751
func index(item interface{}, indices ...interface{}) (interface{}, error) {
	v := reflect.ValueOf(item)
	if !v.IsValid() {
		return nil, fmt.Errorf("index of untyped nil")
	}
	for _, i := range indices {
		index := reflect.ValueOf(i)
		var isNil bool
		if v, isNil = indirect(v); isNil {
			return nil, fmt.Errorf("index of nil pointer")
		}
		switch v.Kind() {
		case reflect.Array, reflect.Slice, reflect.String:
			var x int64
			switch index.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				x = index.Int()
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				x = int64(index.Uint())
			case reflect.Invalid:
				return nil, fmt.Errorf("cannot index slice/array with nil")
			default:
				return nil, fmt.Errorf("cannot index slice/array with type %s", index.Type())
			}
			if x < 0 || x >= int64(v.Len()) {
				// We deviate from stdlib here.  Don't return an error if the
				// index is out of range.
				return nil, nil
			}
			v = v.Index(int(x))
		case reflect.Map:
			index, err := prepareArg(index, v.Type().Key())
			if err != nil {
				return nil, err
			}
			if x := v.MapIndex(index); x.IsValid() {
				v = x
			} else {
				v = reflect.Zero(v.Type().Elem())
			}
		case reflect.Invalid:
			// the loop holds invariant: v.IsValid()
			panic("unreachable")
		default:
			return nil, fmt.Errorf("can't index item of type %s", v.Type())
		}
	}
	return v.Interface(), nil
}

// readFile reads the file named by filename relative to the given basepath
// and returns the contents as a string.
// There is a upper size limit set at 1 megabytes.
func readFile(fs *afero.BasePathFs, filename string) (string, error) {
	if filename == "" {
		return "", errors.New("readFile needs a filename")
	}

	if info, err := fs.Stat(filename); err == nil {
		if info.Size() > 1000000 {
			return "", fmt.Errorf("File %q is too big", filename)
		}
	} else {
		return "", err
	}
	b, err := afero.ReadFile(fs, filename)

	if err != nil {
		return "", err
	}

	return string(b), nil
}

// readFileFromWorkingDir reads the file named by filename relative to the
// configured WorkingDir.
// It returns the contents as a string.
// There is a upper size limit set at 1 megabytes.
func readFileFromWorkingDir(i interface{}) (string, error) {
	return readFile(hugofs.WorkingDir(), cast.ToString(i))
}

// readDirFromWorkingDir listst the directory content relative to the
// configured WorkingDir.
func readDirFromWorkingDir(i interface{}) ([]os.FileInfo, error) {

	path := cast.ToString(i)

	list, err := afero.ReadDir(hugofs.WorkingDir(), path)

	if err != nil {
		return nil, fmt.Errorf("Failed to read Directory %s with error message %s", path, err)
	}

	return list, nil
}

// safeHTMLAttr returns a given string as html/template HTMLAttr content.
//
// safeHTMLAttr is currently disabled, pending further discussion
// on its use case.  2015-01-19
func safeHTMLAttr(a interface{}) template.HTMLAttr {
	return template.HTMLAttr(cast.ToString(a))
}

// safeCSS returns a given string as html/template CSS content.
func safeCSS(a interface{}) template.CSS {
	return template.CSS(cast.ToString(a))
}

// safeURL returns a given string as html/template URL content.
func safeURL(a interface{}) template.URL {
	return template.URL(cast.ToString(a))
}

// safeHTML returns a given string as html/template HTML content.
func safeHTML(a interface{}) template.HTML { return template.HTML(cast.ToString(a)) }

// safeJS returns the given string as a html/template JS content.
func safeJS(a interface{}) template.JS { return template.JS(cast.ToString(a)) }

// mod returns a % b.
func mod(a, b interface{}) (int64, error) {
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)
	var ai, bi int64

	switch av.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		ai = av.Int()
	default:
		return 0, errors.New("Modulo operator can't be used with non integer value")
	}

	switch bv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		bi = bv.Int()
	default:
		return 0, errors.New("Modulo operator can't be used with non integer value")
	}

	if bi == 0 {
		return 0, errors.New("The number can't be divided by zero at modulo operation")
	}

	return ai % bi, nil
}

// modBool returns the boolean of a % b.  If a % b == 0, return true.
func modBool(a, b interface{}) (bool, error) {
	res, err := mod(a, b)
	if err != nil {
		return false, err
	}
	return res == int64(0), nil
}

// base64Decode returns the base64 decoding of the given content.
func base64Decode(content interface{}) (string, error) {
	conv, err := cast.ToStringE(content)

	if err != nil {
		return "", err
	}

	dec, err := base64.StdEncoding.DecodeString(conv)

	return string(dec), err
}

// base64Encode returns the base64 encoding of the given content.
func base64Encode(content interface{}) (string, error) {
	conv, err := cast.ToStringE(content)

	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString([]byte(conv)), nil
}

// countWords returns the approximate word count of the given content.
func countWords(content interface{}) (int, error) {
	conv, err := cast.ToStringE(content)

	if err != nil {
		return 0, fmt.Errorf("Failed to convert content to string: %s", err.Error())
	}

	counter := 0
	for _, word := range strings.Fields(helpers.StripHTML(conv)) {
		runeCount := utf8.RuneCountInString(word)
		if len(word) == runeCount {
			counter++
		} else {
			counter += runeCount
		}
	}

	return counter, nil
}

// countRunes returns the approximate rune count of the given content.
func countRunes(content interface{}) (int, error) {
	conv, err := cast.ToStringE(content)

	if err != nil {
		return 0, fmt.Errorf("Failed to convert content to string: %s", err.Error())
	}

	counter := 0
	for _, r := range helpers.StripHTML(conv) {
		if !helpers.IsWhitespace(r) {
			counter++
		}
	}

	return counter, nil
}

// humanize returns the humanized form of a single word.
// Example:  "my-first-post" -> "My first post"
func humanize(in interface{}) (string, error) {
	word, err := cast.ToStringE(in)
	if err != nil {
		return "", err
	}
	return inflect.Humanize(word), nil
}

// pluralize returns the plural form of a single word.
func pluralize(in interface{}) (string, error) {
	word, err := cast.ToStringE(in)
	if err != nil {
		return "", err
	}
	return inflect.Pluralize(word), nil
}

// singularize returns the singular form of a single word.
func singularize(in interface{}) (string, error) {
	word, err := cast.ToStringE(in)
	if err != nil {
		return "", err
	}
	return inflect.Singularize(word), nil
}

// md5 hashes the given input and returns its MD5 checksum
func md5(in interface{}) (string, error) {
	conv, err := cast.ToStringE(in)
	if err != nil {
		return "", err
	}

	hash := _md5.Sum([]byte(conv))
	return hex.EncodeToString(hash[:]), nil
}

// sha1 hashes the given input and returns its SHA1 checksum
func sha1(in interface{}) (string, error) {
	conv, err := cast.ToStringE(in)
	if err != nil {
		return "", err
	}

	hash := _sha1.Sum([]byte(conv))
	return hex.EncodeToString(hash[:]), nil
}

func init() {
	funcMap = template.FuncMap{
		"absURL":       func(a string) template.HTML { return template.HTML(helpers.AbsURL(a)) },
		"add":          func(a, b interface{}) (interface{}, error) { return helpers.DoArithmetic(a, b, '+') },
		"after":        after,
		"apply":        apply,
		"base64Decode": base64Decode,
		"base64Encode": base64Encode,
		"chomp":        chomp,
		"countrunes":   countRunes,
		"countwords":   countWords,
		"default":      dfault,
		"dateFormat":   dateFormat,
		"delimit":      delimit,
		"dict":         dictionary,
		"div":          func(a, b interface{}) (interface{}, error) { return helpers.DoArithmetic(a, b, '/') },
		"echoParam":    returnWhenSet,
		"emojify":      emojify,
		"eq":           eq,
		"findRE":       findRE,
		"first":        first,
		"ge":           ge,
		"getCSV":       getCSV,
		"getJSON":      getJSON,
		"getenv":       func(varName string) string { return os.Getenv(varName) },
		"gt":           gt,
		"hasPrefix":    func(a, b string) bool { return strings.HasPrefix(a, b) },
		"highlight":    highlight,
		"humanize":     humanize,
		"in":           in,
		"index":        index,
		"int":          func(v interface{}) int { return cast.ToInt(v) },
		"intersect":    intersect,
		"isSet":        isSet,
		"isset":        isSet,
		"jsonify":      jsonify,
		"last":         last,
		"le":           le,
		"lower":        func(a string) string { return strings.ToLower(a) },
		"lt":           lt,
		"markdownify":  markdownify,
		"md5":          md5,
		"mod":          mod,
		"modBool":      modBool,
		"mul":          func(a, b interface{}) (interface{}, error) { return helpers.DoArithmetic(a, b, '*') },
		"ne":           ne,
		"partial":      partial,
		"plainify":     plainify,
		"pluralize":    pluralize,
		"readDir":      readDirFromWorkingDir,
		"readFile":     readFileFromWorkingDir,
		"ref":          ref,
		"relURL":       func(a string) template.HTML { return template.HTML(helpers.RelURL(a)) },
		"relref":       relRef,
		"replace":      replace,
		"replaceRE":    replaceRE,
		"safeCSS":      safeCSS,
		"safeHTML":     safeHTML,
		"safeJS":       safeJS,
		"safeURL":      safeURL,
		"sanitizeURL":  helpers.SanitizeURL,
		"sanitizeurl":  helpers.SanitizeURL,
		"seq":          helpers.Seq,
		"sha1":         sha1,
		"shuffle":      shuffle,
		"singularize":  singularize,
		"slice":        slice,
		"slicestr":     slicestr,
		"sort":         sortSeq,
		"split":        split,
		"string":       func(v interface{}) string { return cast.ToString(v) },
		"sub":          func(a, b interface{}) (interface{}, error) { return helpers.DoArithmetic(a, b, '-') },
		"substr":       substr,
		"title":        func(a string) string { return strings.Title(a) },
		"trim":         trim,
		"upper":        func(a string) string { return strings.ToUpper(a) },
		"urlize":       helpers.URLize,
		"where":        where,
	}
}
