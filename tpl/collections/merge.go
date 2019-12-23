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
	"reflect"
	"strings"

	"github.com/gohugoio/hugo/common/maps"

	"github.com/gohugoio/hugo/common/hreflect"

	"github.com/pkg/errors"
)

// Merge creates a copy of dst and merges src into it.
// Currently only maps supported. Key handling is case insensitive.
func (ns *Namespace) Merge(src, dst interface{}) (interface{}, error) {

	vdst, vsrc := reflect.ValueOf(dst), reflect.ValueOf(src)

	if vdst.Kind() != reflect.Map {
		return nil, errors.Errorf("destination must be a map, got %T", dst)
	}

	if !hreflect.IsTruthfulValue(vsrc) {
		return dst, nil
	}

	if vsrc.Kind() != reflect.Map {
		return nil, errors.Errorf("source must be a map, got %T", src)
	}

	if vsrc.Type().Key() != vdst.Type().Key() {
		return nil, errors.Errorf("incompatible map types, got %T to %T", src, dst)
	}

	return mergeMap(vdst, vsrc).Interface(), nil
}

func caseInsensitiveLookup(m, k reflect.Value) (reflect.Value, bool) {
	if m.Type().Key().Kind() != reflect.String || k.Kind() != reflect.String {
		// Fall back to direct lookup.
		v := m.MapIndex(k)
		return v, hreflect.IsTruthfulValue(v)
	}

	for _, key := range m.MapKeys() {
		if strings.EqualFold(k.String(), key.String()) {
			return m.MapIndex(key), true
		}

	}

	return reflect.Value{}, false
}

func mergeMap(dst, src reflect.Value) reflect.Value {

	out := reflect.MakeMap(dst.Type())

	// If the destination is Params, we must lower case all keys.
	_, lowerCase := dst.Interface().(maps.Params)

	// Copy the destination map.
	for _, key := range dst.MapKeys() {
		v := dst.MapIndex(key)
		out.SetMapIndex(key, v)
	}

	// Add all keys in src not already in destination.
	// Maps of the same type will be merged.
	for _, key := range src.MapKeys() {
		sv := src.MapIndex(key)
		dv, found := caseInsensitiveLookup(dst, key)

		if found {
			// If both are the same map key type, merge.
			dve := dv.Elem()
			if dve.Kind() == reflect.Map {
				sve := sv.Elem()
				if dve.Type().Key() == sve.Type().Key() {
					out.SetMapIndex(key, mergeMap(dve, sve))
				}
			}
		} else {
			if lowerCase && key.Kind() == reflect.String {
				key = reflect.ValueOf(strings.ToLower(key.String()))
			}
			out.SetMapIndex(key, sv)
		}
	}

	return out
}
