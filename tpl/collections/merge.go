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

package collections

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/gohugoio/hugo/common/hreflect"
	"github.com/gohugoio/hugo/common/maps"
)

// Merge creates a copy of the final parameter in params and merges the preceding
// parameters into it in reverse order.
//
// Currently only maps are supported. Key handling is case insensitive.
func (ns *Namespace) Merge(params ...any) (any, error) {
	if len(params) < 2 {
		return nil, errors.New("merge requires at least two parameters")
	}

	var err error
	result := params[len(params)-1]

	for i := len(params) - 2; i >= 0; i-- {
		result, err = ns.merge(params[i], result)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

// merge creates a copy of dst and merges src into it.
func (ns *Namespace) merge(src, dst any) (any, error) {
	vdst, vsrc := reflect.ValueOf(dst), reflect.ValueOf(src)

	if vdst.Kind() != reflect.Map {
		return nil, fmt.Errorf("destination must be a map, got %T", dst)
	}

	if !hreflect.IsTruthfulValue(vsrc) {
		return dst, nil
	}

	if vsrc.Kind() != reflect.Map {
		return nil, fmt.Errorf("source must be a map, got %T", src)
	}

	if vsrc.Type().Key() != vdst.Type().Key() {
		return nil, fmt.Errorf("incompatible map types, got %T to %T", src, dst)
	}

	return mergeMap(vdst, vsrc).Interface(), nil
}

func caseInsensitiveLookup(m, k reflect.Value) (reflect.Value, bool) {
	if m.Type().Key().Kind() != reflect.String || k.Kind() != reflect.String {
		// Fall back to direct lookup.
		v := m.MapIndex(k)
		return v, hreflect.IsTruthfulValue(v)
	}

	k2 := reflect.New(m.Type().Key()).Elem()

	iter := m.MapRange()
	for iter.Next() {
		k2.SetIterKey(iter)
		if strings.EqualFold(k.String(), k2.String()) {
			return iter.Value(), true
		}
	}

	return reflect.Value{}, false
}

func mergeMap(dst, src reflect.Value) reflect.Value {
	out := reflect.MakeMap(dst.Type())

	// If the destination is Params, we must lower case all keys.
	_, lowerCase := dst.Interface().(maps.Params)

	k := reflect.New(dst.Type().Key()).Elem()
	v := reflect.New(dst.Type().Elem()).Elem()

	// Copy the destination map.
	iter := dst.MapRange()
	for iter.Next() {
		k.SetIterKey(iter)
		v.SetIterValue(iter)
		out.SetMapIndex(k, v)
	}

	// Add all keys in src not already in destination.
	// Maps of the same type will be merged.
	k = reflect.New(src.Type().Key()).Elem()
	sv := reflect.New(src.Type().Elem()).Elem()

	iter = src.MapRange()
	for iter.Next() {
		sv.SetIterValue(iter)
		k.SetIterKey(iter)

		dv, found := caseInsensitiveLookup(dst, k)

		if found {
			// If both are the same map key type, merge.
			dve := dv.Elem()
			if dve.Kind() == reflect.Map {
				sve := sv.Elem()
				if sve.Kind() != reflect.Map {
					continue
				}

				if dve.Type().Key() == sve.Type().Key() {
					out.SetMapIndex(k, mergeMap(dve, sve))
				}
			}
		} else {
			kk := k
			if lowerCase && k.Kind() == reflect.String {
				kk = reflect.ValueOf(strings.ToLower(k.String()))
			}
			out.SetMapIndex(kk, sv)
		}
	}

	return out
}
