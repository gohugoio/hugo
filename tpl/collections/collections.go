// Copyright 2019 The Hugo Authors. All rights reserved.
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

// Package collections provides template functions for manipulating collections
// such as arrays, maps, and slices.
package collections

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/gohugoio/hugo/common/collections"
	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/common/types"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/langs"
	"github.com/gohugoio/hugo/tpl/compare"
	"github.com/spf13/cast"
)

// New returns a new instance of the collections-namespaced template functions.
func New(deps *deps.Deps) *Namespace {
	language := deps.Conf.Language()
	if language == nil {
		panic("language must be set")
	}
	loc := langs.GetLocation(language)

	return &Namespace{
		loc:      loc,
		sortComp: compare.New(loc, true),
		deps:     deps,
	}
}

// Namespace provides template functions for the "collections" namespace.
type Namespace struct {
	loc      *time.Location
	sortComp *compare.Namespace
	deps     *deps.Deps
}

// After returns all the items after the first n items in list l.
func (ns *Namespace) After(n any, l any) (any, error) {
	if n == nil || l == nil {
		return nil, errors.New("both limit and seq must be provided")
	}

	nv, err := cast.ToIntE(n)
	if err != nil {
		return nil, err
	}

	if nv < 0 {
		return nil, errors.New("sequence bounds out of range [" + cast.ToString(nv) + ":]")
	}

	lv := reflect.ValueOf(l)
	lv, isNil := indirect(lv)
	if isNil {
		return nil, errors.New("can't iterate over a nil value")
	}

	switch lv.Kind() {
	case reflect.Array, reflect.Slice, reflect.String:
		// okay
	default:
		return nil, errors.New("can't iterate over " + reflect.ValueOf(l).Type().String())
	}

	if nv >= lv.Len() {
		return lv.Slice(0, 0).Interface(), nil
	}

	return lv.Slice(nv, lv.Len()).Interface(), nil
}

// Delimit takes a given list l and returns a string delimited by sep.
// If last is passed to the function, it will be used as the final delimiter.
func (ns *Namespace) Delimit(ctx context.Context, l, sep any, last ...any) (string, error) {
	d, err := cast.ToStringE(sep)
	if err != nil {
		return "", err
	}

	var dLast *string
	if len(last) > 0 {
		l := last[0]
		dStr, err := cast.ToStringE(l)
		if err != nil {
			dLast = nil
		} else {
			dLast = &dStr
		}
	}

	lv := reflect.ValueOf(l)
	lv, isNil := indirect(lv)
	if isNil {
		return "", errors.New("can't iterate over a nil value")
	}

	var str string
	switch lv.Kind() {
	case reflect.Map:
		sortSeq, err := ns.Sort(ctx, l)
		if err != nil {
			return "", err
		}
		lv = reflect.ValueOf(sortSeq)
		fallthrough
	case reflect.Array, reflect.Slice, reflect.String:
		for i := 0; i < lv.Len(); i++ {
			val := lv.Index(i).Interface()
			valStr, err := cast.ToStringE(val)
			if err != nil {
				continue
			}
			switch {
			case i == lv.Len()-2 && dLast != nil:
				str += valStr + *dLast
			case i == lv.Len()-1:
				str += valStr
			default:
				str += valStr + d
			}
		}

	default:
		return "", fmt.Errorf("can't iterate over %T", l)
	}

	return str, nil
}

// Dictionary creates a new map from the given parameters by
// treating values as key-value pairs.  The number of values must be even.
// The keys can be string slices, which will create the needed nested structure.
func (ns *Namespace) Dictionary(values ...any) (map[string]any, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dictionary call")
	}

	root := make(map[string]any)

	for i := 0; i < len(values); i += 2 {
		dict := root
		var key string
		switch v := values[i].(type) {
		case string:
			key = v
		case []string:
			for i := 0; i < len(v)-1; i++ {
				key = v[i]
				var m map[string]any
				v, found := dict[key]
				if found {
					m = v.(map[string]any)
				} else {
					m = make(map[string]any)
					dict[key] = m
				}
				dict = m
			}
			key = v[len(v)-1]
		default:
			return nil, errors.New("invalid dictionary key")
		}
		dict[key] = values[i+1]
	}

	return root, nil
}

// EchoParam returns the value in the collection c with key k if is set; otherwise, it returns an
// empty string.
// Deprecated: Use the index function instead.
func (ns *Namespace) EchoParam(c, k any) any {
	hugo.Deprecate("collections.EchoParam", "Use the index function instead.", "v0.120.0")
	av, isNil := indirect(reflect.ValueOf(c))
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
		case reflect.Bool:
			return avv.Bool()
		}
	}

	return ""
}

// First returns the first limit items in list l.
func (ns *Namespace) First(limit any, l any) (any, error) {
	if limit == nil || l == nil {
		return nil, errors.New("both limit and seq must be provided")
	}

	limitv, err := cast.ToIntE(limit)
	if err != nil {
		return nil, err
	}

	if limitv < 0 {
		return nil, errors.New("sequence length must be non-negative")
	}

	lv := reflect.ValueOf(l)
	lv, isNil := indirect(lv)
	if isNil {
		return nil, errors.New("can't iterate over a nil value")
	}

	switch lv.Kind() {
	case reflect.Array, reflect.Slice, reflect.String:
		// okay
	default:
		return nil, errors.New("can't iterate over " + reflect.ValueOf(l).Type().String())
	}

	if limitv > lv.Len() {
		limitv = lv.Len()
	}

	return lv.Slice(0, limitv).Interface(), nil
}

// In returns whether v is in the list l.  l may be an array or slice.
func (ns *Namespace) In(l any, v any) (bool, error) {
	if l == nil || v == nil {
		return false, nil
	}

	lv := reflect.ValueOf(l)
	vv := reflect.ValueOf(v)

	vvk := normalize(vv)

	switch lv.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < lv.Len(); i++ {
			lvv, isNil := indirectInterface(lv.Index(i))
			if isNil {
				continue
			}

			lvvk := normalize(lvv)

			if lvvk == vvk {
				return true, nil
			}
		}
	}
	ss, err := cast.ToStringE(l)
	if err != nil {
		return false, nil
	}

	su, err := cast.ToStringE(v)
	if err != nil {
		return false, nil
	}
	return strings.Contains(ss, su), nil
}

// Intersect returns the common elements in the given sets, l1 and l2.  l1 and
// l2 must be of the same type and may be either arrays or slices.
func (ns *Namespace) Intersect(l1, l2 any) (any, error) {
	if l1 == nil || l2 == nil {
		return make([]any, 0), nil
	}

	var ins *intersector

	l1v := reflect.ValueOf(l1)
	l2v := reflect.ValueOf(l2)

	switch l1v.Kind() {
	case reflect.Array, reflect.Slice:
		ins = &intersector{r: reflect.MakeSlice(l1v.Type(), 0, 0), seen: make(map[any]bool)}
		switch l2v.Kind() {
		case reflect.Array, reflect.Slice:
			for i := 0; i < l1v.Len(); i++ {
				l1vv := l1v.Index(i)
				if !l1vv.Type().Comparable() {
					return make([]any, 0), errors.New("intersect does not support slices or arrays of uncomparable types")
				}

				for j := 0; j < l2v.Len(); j++ {
					l2vv := l2v.Index(j)
					if !l2vv.Type().Comparable() {
						return make([]any, 0), errors.New("intersect does not support slices or arrays of uncomparable types")
					}

					ins.handleValuePair(l1vv, l2vv)
				}
			}
			return ins.r.Interface(), nil
		default:
			return nil, errors.New("can't iterate over " + reflect.ValueOf(l2).Type().String())
		}
	default:
		return nil, errors.New("can't iterate over " + reflect.ValueOf(l1).Type().String())
	}
}

// Group groups a set of items by the given key.
// This is currently only supported for Pages.
func (ns *Namespace) Group(key any, items any) (any, error) {
	if key == nil {
		return nil, errors.New("nil is not a valid key to group by")
	}

	if g, ok := items.(collections.Grouper); ok {
		return g.Group(key, items)
	}

	in := newSliceElement(items)

	if g, ok := in.(collections.Grouper); ok {
		return g.Group(key, items)
	}

	return nil, fmt.Errorf("grouping not supported for type %T %T", items, in)
}

// IsSet returns whether a given array, channel, slice, or map in c has the given key
// defined.
func (ns *Namespace) IsSet(c any, key any) (bool, error) {
	av := reflect.ValueOf(c)
	kv := reflect.ValueOf(key)

	switch av.Kind() {
	case reflect.Array, reflect.Chan, reflect.Slice:
		k, err := cast.ToIntE(key)
		if err != nil {
			return false, fmt.Errorf("isset unable to use key of type %T as index", key)
		}
		if av.Len() > k {
			return true, nil
		}
	case reflect.Map:
		if kv.Type() == av.Type().Key() {
			return av.MapIndex(kv).IsValid(), nil
		}
	default:
		ns.deps.Log.Warnf("calling IsSet with unsupported type %q (%T) will always return false.\n", av.Kind(), c)
	}

	return false, nil
}

// Last returns the last limit items in the list l.
func (ns *Namespace) Last(limit any, l any) (any, error) {
	if limit == nil || l == nil {
		return nil, errors.New("both limit and seq must be provided")
	}

	limitv, err := cast.ToIntE(limit)
	if err != nil {
		return nil, err
	}

	if limitv < 0 {
		return nil, errors.New("sequence length must be non-negative")
	}

	seqv := reflect.ValueOf(l)
	seqv, isNil := indirect(seqv)
	if isNil {
		return nil, errors.New("can't iterate over a nil value")
	}

	switch seqv.Kind() {
	case reflect.Array, reflect.Slice, reflect.String:
		// okay
	default:
		return nil, errors.New("can't iterate over " + reflect.ValueOf(l).Type().String())
	}

	if limitv > seqv.Len() {
		limitv = seqv.Len()
	}

	return seqv.Slice(seqv.Len()-limitv, seqv.Len()).Interface(), nil
}

// Querify encodes the given params in URL-encoded form ("bar=baz&foo=quux") sorted by key.
func (ns *Namespace) Querify(params ...any) (string, error) {
	qs := url.Values{}

	if len(params) == 1 {
		switch v := params[0].(type) {
		case []string:
			if len(v)%2 != 0 {
				return "", errors.New("invalid query")
			}

			for i := 0; i < len(v); i += 2 {
				qs.Add(v[i], v[i+1])
			}

			return qs.Encode(), nil

		case []any:
			params = v

		default:
			return "", errors.New("query keys must be strings")
		}
	}

	if len(params)%2 != 0 {
		return "", errors.New("invalid query")
	}

	for i := 0; i < len(params); i += 2 {
		switch v := params[i].(type) {
		case string:
			qs.Add(v, fmt.Sprintf("%v", params[i+1]))
		default:
			return "", errors.New("query keys must be strings")
		}
	}

	return qs.Encode(), nil
}

// Reverse creates a copy of the list l and reverses it.
func (ns *Namespace) Reverse(l any) (any, error) {
	if l == nil {
		return nil, nil
	}
	v := reflect.ValueOf(l)

	switch v.Kind() {
	case reflect.Slice:
	default:
		return nil, errors.New("argument must be a slice")
	}

	sliceCopy := reflect.MakeSlice(v.Type(), v.Len(), v.Len())

	for i := v.Len() - 1; i >= 0; i-- {
		element := sliceCopy.Index(i)
		element.Set(v.Index(v.Len() - 1 - i))
	}

	return sliceCopy.Interface(), nil
}

// Seq creates a sequence of integers from args. It's named and used as GNU's seq.
//
// Examples:
//
//	3 => 1, 2, 3
//	1 2 4 => 1, 3
//	-3 => -1, -2, -3
//	1 4 => 1, 2, 3, 4
//	1 -2 => 1, 0, -1, -2
func (ns *Namespace) Seq(args ...any) ([]int, error) {
	if len(args) < 1 || len(args) > 3 {
		return nil, errors.New("invalid number of arguments to Seq")
	}

	intArgs := cast.ToIntSlice(args)
	if len(intArgs) < 1 || len(intArgs) > 3 {
		return nil, errors.New("invalid arguments to Seq")
	}

	inc := 1
	var last int
	first := intArgs[0]

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

// Shuffle returns list l in a randomized order.
func (ns *Namespace) Shuffle(l any) (any, error) {
	if l == nil {
		return nil, errors.New("both count and seq must be provided")
	}

	lv := reflect.ValueOf(l)
	lv, isNil := indirect(lv)
	if isNil {
		return nil, errors.New("can't iterate over a nil value")
	}

	switch lv.Kind() {
	case reflect.Array, reflect.Slice, reflect.String:
		// okay
	default:
		return nil, errors.New("can't iterate over " + reflect.ValueOf(l).Type().String())
	}

	shuffled := reflect.MakeSlice(reflect.TypeOf(l), lv.Len(), lv.Len())

	randomIndices := rand.Perm(lv.Len())

	for index, value := range randomIndices {
		shuffled.Index(value).Set(lv.Index(index))
	}

	return shuffled.Interface(), nil
}

// Slice returns a slice of all passed arguments.
func (ns *Namespace) Slice(args ...any) any {
	if len(args) == 0 {
		return args
	}

	return collections.Slice(args...)
}

type intersector struct {
	r    reflect.Value
	seen map[any]bool
}

func (i *intersector) appendIfNotSeen(v reflect.Value) {
	k := normalize(v)
	if !i.seen[k] {
		i.r = reflect.Append(i.r, v)
		i.seen[k] = true
	}
}

func (i *intersector) handleValuePair(l1vv, l2vv reflect.Value) {
	switch kind := l1vv.Kind(); {
	case kind == reflect.String:
		l2t, err := toString(l2vv)
		if err == nil && l1vv.String() == l2t {
			i.appendIfNotSeen(l1vv)
		}
	case isNumber(kind):
		f1, err1 := numberToFloat(l1vv)
		f2, err2 := numberToFloat(l2vv)
		if err1 == nil && err2 == nil && f1 == f2 {
			i.appendIfNotSeen(l1vv)
		}
	case kind == reflect.Ptr, kind == reflect.Struct:
		if types.Unwrapv(l1vv.Interface()) == types.Unwrapv(l2vv.Interface()) {
			i.appendIfNotSeen(l1vv)
		}
	case kind == reflect.Interface:
		i.handleValuePair(reflect.ValueOf(l1vv.Interface()), l2vv)
	}
}

// Union returns the union of the given sets, l1 and l2. l1 and
// l2 must be of the same type and may be either arrays or slices.
// If l1 and l2 aren't of the same type then l1 will be returned.
// If either l1 or l2 is nil then the non-nil list will be returned.
func (ns *Namespace) Union(l1, l2 any) (any, error) {
	if l1 == nil && l2 == nil {
		return []any{}, nil
	} else if l1 == nil && l2 != nil {
		return l2, nil
	} else if l1 != nil && l2 == nil {
		return l1, nil
	}

	l1v := reflect.ValueOf(l1)
	l2v := reflect.ValueOf(l2)

	var ins *intersector

	switch l1v.Kind() {
	case reflect.Array, reflect.Slice:
		switch l2v.Kind() {
		case reflect.Array, reflect.Slice:
			ins = &intersector{r: reflect.MakeSlice(l1v.Type(), 0, 0), seen: make(map[any]bool)}

			if l1v.Type() != l2v.Type() &&
				l1v.Type().Elem().Kind() != reflect.Interface &&
				l2v.Type().Elem().Kind() != reflect.Interface {
				return ins.r.Interface(), nil
			}

			var (
				l1vv  reflect.Value
				isNil bool
			)

			for i := 0; i < l1v.Len(); i++ {
				l1vv, isNil = indirectInterface(l1v.Index(i))

				if !l1vv.Type().Comparable() {
					return []any{}, errors.New("union does not support slices or arrays of uncomparable types")
				}

				if !isNil {
					ins.appendIfNotSeen(l1vv)
				}
			}

			if !l1vv.IsValid() {
				// The first slice may be empty. Pick the first value of the second
				// to use as a prototype.
				if l2v.Len() > 0 {
					l1vv = l2v.Index(0)
				}
			}

			for j := 0; j < l2v.Len(); j++ {
				l2vv := l2v.Index(j)

				switch kind := l1vv.Kind(); {
				case kind == reflect.String:
					l2t, err := toString(l2vv)
					if err == nil {
						ins.appendIfNotSeen(reflect.ValueOf(l2t))
					}
				case isNumber(kind):
					var err error
					l2vv, err = convertNumber(l2vv, kind)
					if err == nil {
						ins.appendIfNotSeen(l2vv)
					}
				case kind == reflect.Interface, kind == reflect.Struct, kind == reflect.Ptr:
					ins.appendIfNotSeen(l2vv)

				}
			}

			return ins.r.Interface(), nil
		default:
			return nil, errors.New("can't iterate over " + reflect.ValueOf(l2).Type().String())
		}
	default:
		return nil, errors.New("can't iterate over " + reflect.ValueOf(l1).Type().String())
	}
}

// Uniq returns a new list with duplicate elements in the list l removed.
func (ns *Namespace) Uniq(l any) (any, error) {
	if l == nil {
		return make([]any, 0), nil
	}

	v := reflect.ValueOf(l)
	var slice reflect.Value

	switch v.Kind() {
	case reflect.Slice:
		slice = reflect.MakeSlice(v.Type(), 0, 0)

	case reflect.Array:
		slice = reflect.MakeSlice(reflect.SliceOf(v.Type().Elem()), 0, 0)
	default:
		return nil, fmt.Errorf("type %T not supported", l)
	}

	seen := make(map[any]bool)

	for i := 0; i < v.Len(); i++ {
		ev, _ := indirectInterface(v.Index(i))

		key := normalize(ev)

		if _, found := seen[key]; !found {
			slice = reflect.Append(slice, ev)
			seen[key] = true
		}
	}

	return slice.Interface(), nil
}

// KeyVals creates a key and values wrapper.
func (ns *Namespace) KeyVals(key any, values ...any) (types.KeyValues, error) {
	return types.KeyValues{Key: key, Values: values}, nil
}

// NewScratch creates a new Scratch which can be used to store values in a
// thread safe way.
func (ns *Namespace) NewScratch() *maps.Scratch {
	return maps.NewScratch()
}
