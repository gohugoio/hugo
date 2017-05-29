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

package collections

import (
	"errors"
	"fmt"
	"html/template"
	"math/rand"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/spf13/cast"
	"github.com/spf13/hugo/deps"
	"github.com/spf13/hugo/helpers"
)

// New returns a new instance of the collections-namespaced template functions.
func New(deps *deps.Deps) *Namespace {
	return &Namespace{
		deps: deps,
	}
}

// Namespace provides template functions for the "collections" namespace.
type Namespace struct {
	deps *deps.Deps
}

// After returns all the items after the first N in a rangeable list.
func (ns *Namespace) After(index interface{}, seq interface{}) (interface{}, error) {
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

// Delimit takes a given sequence and returns a delimited HTML string.
// If last is passed to the function, it will be used as the final delimiter.
func (ns *Namespace) Delimit(seq, delimiter interface{}, last ...interface{}) (template.HTML, error) {
	d, err := cast.ToStringE(delimiter)
	if err != nil {
		return "", err
	}

	var dLast *string
	if len(last) > 0 {
		l := last[0]
		dStr, err := cast.ToStringE(l)
		if err != nil {
			dLast = nil
		}
		dLast = &dStr
	}

	seqv := reflect.ValueOf(seq)
	seqv, isNil := indirect(seqv)
	if isNil {
		return "", errors.New("can't iterate over a nil value")
	}

	var str string
	switch seqv.Kind() {
	case reflect.Map:
		sortSeq, err := ns.Sort(seq)
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

// Dictionary creates a map[string]interface{} from the given parameters by
// walking the parameters and treating them as key-value pairs.  The number
// of parameters must be even.
func (ns *Namespace) Dictionary(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dictionary call")
	}

	dict := make(map[string]interface{}, len(values)/2)

	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dictionary keys must be strings")
		}
		dict[key] = values[i+1]
	}

	return dict, nil
}

// EchoParam returns a given value if it is set; otherwise, it returns an
// empty string.
func (ns *Namespace) EchoParam(a, key interface{}) interface{} {
	av, isNil := indirect(reflect.ValueOf(a))
	if isNil {
		return ""
	}

	var avv reflect.Value
	switch av.Kind() {
	case reflect.Array, reflect.Slice:
		index, ok := key.(int)
		if ok && av.Len() > index {
			avv = av.Index(index)
		}
	case reflect.Map:
		kv := reflect.ValueOf(key)
		if kv.Type().AssignableTo(av.Type().Key()) {
			avv = av.MapIndex(kv)
		}
	}

	avv, isNil = indirect(avv)

	if isNil {
		return ""
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

// First returns the first N items in a rangeable list.
func (ns *Namespace) First(limit interface{}, seq interface{}) (interface{}, error) {
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

// In returns whether v is in the set l.  l may be an array or slice.
func (ns *Namespace) In(l interface{}, v interface{}) bool {
	if l == nil || v == nil {
		return false
	}

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

// Intersect returns the common elements in the given sets, l1 and l2.  l1 and
// l2 must be of the same type and may be either arrays or slices.
func (ns *Namespace) Intersect(l1, l2 interface{}) (interface{}, error) {
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
						l2t, err := toString(l2vv)
						if err == nil && l1vv.String() == l2t && !ns.In(r.Interface(), l1vv.Interface()) {
							r = reflect.Append(r, l1vv)
						}
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						l2t, err := toInt(l2vv)
						if err == nil && l1vv.Int() == l2t && !ns.In(r.Interface(), l1vv.Interface()) {
							r = reflect.Append(r, l1vv)
						}
					case reflect.Float32, reflect.Float64:
						l2t, err := toFloat(l2vv)
						if err == nil && l1vv.Float() == l2t && !ns.In(r.Interface(), l1vv.Interface()) {
							r = reflect.Append(r, l1vv)
						}
					case reflect.Interface:
						switch l1vvActual := l1vv.Interface().(type) {
						case string:
							switch l2vvActual := l2vv.Interface().(type) {
							case string:
								if l1vvActual == l2vvActual && !ns.In(r.Interface(), l1vvActual) {
									r = reflect.Append(r, l1vv)
								}
							}
						case int, int8, int16, int32, int64:
							switch l2vvActual := l2vv.Interface().(type) {
							case int, int8, int16, int32, int64:
								if l1vvActual == l2vvActual && !ns.In(r.Interface(), l1vvActual) {
									r = reflect.Append(r, l1vv)
								}
							}
						case uint, uint8, uint16, uint32, uint64:
							switch l2vvActual := l2vv.Interface().(type) {
							case uint, uint8, uint16, uint32, uint64:
								if l1vvActual == l2vvActual && !ns.In(r.Interface(), l1vvActual) {
									r = reflect.Append(r, l1vv)
								}
							}
						case float32, float64:
							switch l2vvActual := l2vv.Interface().(type) {
							case float32, float64:
								if l1vvActual == l2vvActual && !ns.In(r.Interface(), l1vvActual) {
									r = reflect.Append(r, l1vv)
								}
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

// IsSet returns whether a given array, channel, slice, or map has a key
// defined.
func (ns *Namespace) IsSet(a interface{}, key interface{}) (bool, error) {
	av := reflect.ValueOf(a)
	kv := reflect.ValueOf(key)

	switch av.Kind() {
	case reflect.Array, reflect.Chan, reflect.Slice:
		if int64(av.Len()) > kv.Int() {
			return true, nil
		}
	case reflect.Map:
		if kv.Type() == av.Type().Key() {
			return av.MapIndex(kv).IsValid(), nil
		}
	default:
		helpers.DistinctFeedbackLog.Printf("WARNING: calling IsSet with unsupported type %q (%T) will always return false.\n", av.Kind(), a)
	}

	return false, nil
}

// Last returns the last N items in a rangeable list.
func (ns *Namespace) Last(limit interface{}, seq interface{}) (interface{}, error) {
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

// Querify encodes the given parameters in URL-encoded form ("bar=baz&foo=quux") sorted by key.
func (ns *Namespace) Querify(params ...interface{}) (string, error) {
	qs := url.Values{}
	vals, err := ns.Dictionary(params...)
	if err != nil {
		return "", errors.New("querify keys must be strings")
	}

	for name, value := range vals {
		qs.Add(name, fmt.Sprintf("%v", value))
	}

	return qs.Encode(), nil
}

// Seq creates a sequence of integers.  It's named and used as GNU's seq.
//
// Examples:
//     3 => 1, 2, 3
//     1 2 4 => 1, 3
//     -3 => -1, -2, -3
//     1 4 => 1, 2, 3, 4
//     1 -2 => 1, 0, -1, -2
func (ns *Namespace) Seq(args ...interface{}) ([]int, error) {
	if len(args) < 1 || len(args) > 3 {
		return nil, errors.New("invalid number of arguments to Seq")
	}

	intArgs := cast.ToIntSlice(args)
	if len(intArgs) < 1 || len(intArgs) > 3 {
		return nil, errors.New("invalid arguments to Seq")
	}

	var inc = 1
	var last int
	var first = intArgs[0]

	if len(intArgs) == 1 {
		last = first
		if last == 0 {
			return []int{}, nil
		} else if last > 0 {
			first = 1
		} else {
			first = -1
			inc = -1
		}
	} else if len(intArgs) == 2 {
		last = intArgs[1]
		if last < first {
			inc = -1
		}
	} else {
		inc = intArgs[1]
		last = intArgs[2]
		if inc == 0 {
			return nil, errors.New("'increment' must not be 0")
		}
		if first < last && inc < 0 {
			return nil, errors.New("'increment' must be > 0")
		}
		if first > last && inc > 0 {
			return nil, errors.New("'increment' must be < 0")
		}
	}

	// sanity check
	if last < -100000 {
		return nil, errors.New("size of result exceeds limit")
	}
	size := ((last - first) / inc) + 1

	// sanity check
	if size <= 0 || size > 2000 {
		return nil, errors.New("size of result exceeds limit")
	}

	seq := make([]int, size)
	val := first
	for i := 0; ; i++ {
		seq[i] = val
		val += inc
		if (inc < 0 && val < last) || (inc > 0 && val > last) {
			break
		}
	}

	return seq, nil
}

// Shuffle returns the given rangeable list in a randomised order.
func (ns *Namespace) Shuffle(seq interface{}) (interface{}, error) {
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

// Slice returns a slice of all passed arguments.
func (ns *Namespace) Slice(args ...interface{}) []interface{} {
	return args
}

// Union returns the union of the given sets, l1 and l2. l1 and
// l2 must be of the same type and may be either arrays or slices.
// If l1 and l2 aren't of the same type then l1 will be returned.
// If either l1 or l2 is nil then the non-nil list will be returned.
func (ns *Namespace) Union(l1, l2 interface{}) (interface{}, error) {
	if l1 == nil && l2 == nil {
		return nil, errors.New("both arrays/slices have to be of the same type")
	} else if l1 == nil && l2 != nil {
		return l2, nil
	} else if l1 != nil && l2 == nil {
		return l1, nil
	}

	l1v := reflect.ValueOf(l1)
	l2v := reflect.ValueOf(l2)

	switch l1v.Kind() {
	case reflect.Array, reflect.Slice:
		switch l2v.Kind() {
		case reflect.Array, reflect.Slice:
			r := reflect.MakeSlice(l1v.Type(), 0, 0)

			if l1v.Type() != l2v.Type() {
				return r.Interface(), nil
			}

			for i := 0; i < l1v.Len(); i++ {
				elem := l1v.Index(i)
				if !ns.In(r.Interface(), elem.Interface()) {
					r = reflect.Append(r, elem)
				}
			}

			for j := 0; j < l2v.Len(); j++ {
				elem := l2v.Index(j)
				if !ns.In(r.Interface(), elem.Interface()) {
					r = reflect.Append(r, elem)
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

// Uniq takes in a slice or array and returns a slice with subsequent
// duplicate elements removed.
func (ns *Namespace) Uniq(l interface{}) (interface{}, error) {
	if l == nil {
		return make([]interface{}, 0), nil
	}

	lv := reflect.ValueOf(l)
	lv, isNil := indirect(lv)
	if isNil {
		return nil, errors.New("invalid nil argument to Uniq")
	}

	var ret reflect.Value

	switch lv.Kind() {
	case reflect.Slice:
		ret = reflect.MakeSlice(lv.Type(), 0, 0)
	case reflect.Array:
		ret = reflect.MakeSlice(reflect.SliceOf(lv.Type().Elem()), 0, 0)
	default:
		return nil, errors.New("Can't use Uniq on " + reflect.ValueOf(lv).Type().String())
	}

	for i := 0; i != lv.Len(); i++ {
		lvv := lv.Index(i)
		lvv, isNil := indirect(lvv)
		if isNil {
			continue
		}

		if !ns.In(ret.Interface(), lvv.Interface()) {
			ret = reflect.Append(ret, lvv)
		}
	}
	return ret.Interface(), nil
}
