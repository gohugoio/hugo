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
	"reflect"
	"sort"
	"strings"

	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/tpl/compare"
	"github.com/spf13/cast"
)

var sortComp = compare.New(true)

// Sort returns a sorted sequence.
func (ns *Namespace) Sort(seq interface{}, args ...interface{}) (interface{}, error) {
	if seq == nil {
		return nil, errors.New("sequence must be provided")
	}

	seqv, isNil := indirect(reflect.ValueOf(seq))
	if isNil {
		return nil, errors.New("can't iterate over a nil value")
	}

	var sliceType reflect.Type
	switch seqv.Kind() {
	case reflect.Array, reflect.Slice:
		sliceType = seqv.Type()
	case reflect.Map:
		sliceType = reflect.SliceOf(seqv.Type().Elem())
	default:
		return nil, errors.New("can't sort " + reflect.ValueOf(seq).Type().String())
	}

	// Create a list of pairs that will be used to do the sort
	p := pairList{SortAsc: true, SliceType: sliceType}
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
			p.Pairs[i].Value = seqv.Index(i)
			if sortByField == "" || sortByField == "value" {
				p.Pairs[i].Key = p.Pairs[i].Value
			} else {
				v := p.Pairs[i].Value
				var err error
				for i, elemName := range path {
					v, err = evaluateSubElem(v, elemName)
					if err != nil {
						return nil, err
					}
					if !v.IsValid() {
						continue
					}
					// Special handling of lower cased maps.
					if params, ok := v.Interface().(maps.Params); ok {
						v = reflect.ValueOf(params.Get(path[i+1:]...))
						break
					}
				}
				p.Pairs[i].Key = v
			}
		}

	case reflect.Map:
		keys := seqv.MapKeys()
		for i := 0; i < seqv.Len(); i++ {
			p.Pairs[i].Value = seqv.MapIndex(keys[i])

			if sortByField == "" {
				p.Pairs[i].Key = keys[i]
			} else if sortByField == "value" {
				p.Pairs[i].Key = p.Pairs[i].Value
			} else {
				v := p.Pairs[i].Value
				var err error
				for i, elemName := range path {
					v, err = evaluateSubElem(v, elemName)
					if err != nil {
						return nil, err
					}
					if !v.IsValid() {
						continue
					}
					// Special handling of lower cased maps.
					if params, ok := v.Interface().(maps.Params); ok {
						v = reflect.ValueOf(params.Get(path[i+1:]...))
						break
					}
				}
				p.Pairs[i].Key = v
			}
		}
	}
	return p.sort(), nil
}

// Credit for pair sorting method goes to Andrew Gerrand
// https://groups.google.com/forum/#!topic/golang-nuts/FT7cjmcL7gw
// A data structure to hold a key/value pair.
type pair struct {
	Key   reflect.Value
	Value reflect.Value
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
	iv := p.Pairs[i].Key
	jv := p.Pairs[j].Key

	if iv.IsValid() {
		if jv.IsValid() {
			// can only call Interface() on valid reflect Values
			return sortComp.Lt(iv.Interface(), jv.Interface())
		}

		// if j is invalid, test i against i's zero value
		return sortComp.Lt(iv.Interface(), reflect.Zero(iv.Type()))
	}

	if jv.IsValid() {
		// if i is invalid, test j against j's zero value
		return sortComp.Lt(reflect.Zero(jv.Type()), jv.Interface())
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
