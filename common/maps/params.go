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

package maps

import (
	"fmt"
	"strings"

	"github.com/spf13/cast"
)

// Params is a map where all keys are lower case.
type Params map[string]interface{}

// Get does a lower case and nested search in this map.
// It will return nil if none found.
func (p Params) Get(indices ...string) interface{} {
	v, _, _ := getNested(p, indices)
	return v
}

// Set overwrites values in p with values in pp for common or new keys.
// This is done recursively.
func (p Params) Set(pp Params) {
	for k, v := range pp {
		vv, found := p[k]
		if !found {
			p[k] = v
		} else {
			switch vvv := vv.(type) {
			case Params:
				if pv, ok := v.(Params); ok {
					vvv.Set(pv)
				} else {
					p[k] = v
				}
			default:
				p[k] = v
			}
		}
	}
}

// IsZero returns true if p is considered empty.
func (p Params) IsZero() bool {
	if p == nil || len(p) == 0 {
		return true
	}

	if len(p) > 1 {
		return false
	}

	for k, _ := range p {
		return k == mergeStrategyKey
	}

	return false

}

// Merge transfers values from pp to p for new keys.
// This is done recursively.
func (p Params) Merge(pp Params) {
	p.merge("", pp)
}

// MergeRoot transfers values from pp to p for new keys where p is the
// root of the tree.
// This is done recursively.
func (p Params) MergeRoot(pp Params) {
	ms, _ := p.GetMergeStrategy()
	p.merge(ms, pp)
}

func (p Params) merge(ps ParamsMergeStrategy, pp Params) {
	ns, found := p.GetMergeStrategy()

	var ms = ns
	if !found && ps != "" {
		ms = ps
	}

	noUpdate := ms == ParamsMergeStrategyNone
	noUpdate = noUpdate || (ps != "" && ps == ParamsMergeStrategyShallow)

	for k, v := range pp {

		if k == mergeStrategyKey {
			continue
		}
		vv, found := p[k]

		if found {
			// Key matches, if both sides are Params, we try to merge.
			if vvv, ok := vv.(Params); ok {
				if pv, ok := v.(Params); ok {
					vvv.merge(ms, pv)
				}
			}
		} else if !noUpdate {
			p[k] = v
		}

	}
}

func (p Params) GetMergeStrategy() (ParamsMergeStrategy, bool) {
	if v, found := p[mergeStrategyKey]; found {
		if s, ok := v.(ParamsMergeStrategy); ok {
			return s, true
		}
	}
	return ParamsMergeStrategyShallow, false
}

func (p Params) DeleteMergeStrategy() bool {
	if _, found := p[mergeStrategyKey]; found {
		delete(p, mergeStrategyKey)
		return true
	}
	return false
}

func (p Params) SetDefaultMergeStrategy(s ParamsMergeStrategy) {
	switch s {
	case ParamsMergeStrategyDeep, ParamsMergeStrategyNone, ParamsMergeStrategyShallow:
	default:
		panic(fmt.Sprintf("invalid merge strategy %q", s))
	}
	p[mergeStrategyKey] = s
}

func getNested(m map[string]interface{}, indices []string) (interface{}, string, map[string]interface{}) {
	if len(indices) == 0 {
		return nil, "", nil
	}

	first := indices[0]
	v, found := m[strings.ToLower(cast.ToString(first))]
	if !found {
		if len(indices) == 1 {
			return nil, first, m
		}
		return nil, "", nil

	}

	if len(indices) == 1 {
		return v, first, m
	}

	switch m2 := v.(type) {
	case Params:
		return getNested(m2, indices[1:])
	case map[string]interface{}:
		return getNested(m2, indices[1:])
	default:
		return nil, "", nil
	}
}

// GetNestedParam gets the first match of the keyStr in the candidates given.
// It will first try the exact match and then try to find it as a nested map value,
// using the given separator, e.g. "mymap.name".
// It assumes that all the maps given have lower cased keys.
func GetNestedParam(keyStr, separator string, candidates ...Params) (interface{}, error) {
	keyStr = strings.ToLower(keyStr)

	// Try exact match first
	for _, m := range candidates {
		if v, ok := m[keyStr]; ok {
			return v, nil
		}
	}

	keySegments := strings.Split(keyStr, separator)
	for _, m := range candidates {
		if v := m.Get(keySegments...); v != nil {
			return v, nil
		}
	}

	return nil, nil
}

func GetNestedParamFn(keyStr, separator string, lookupFn func(key string) interface{}) (interface{}, string, map[string]interface{}, error) {
	keySegments := strings.Split(keyStr, separator)
	if len(keySegments) == 0 {
		return nil, "", nil, nil
	}

	first := lookupFn(keySegments[0])
	if first == nil {
		return nil, "", nil, nil
	}

	if len(keySegments) == 1 {
		return first, keySegments[0], nil, nil
	}

	switch m := first.(type) {
	case map[string]interface{}:
		v, key, owner := getNested(m, keySegments[1:])
		return v, key, owner, nil
	case Params:
		v, key, owner := getNested(m, keySegments[1:])
		return v, key, owner, nil
	}

	return nil, "", nil, nil
}

// ParamsMergeStrategy tells what strategy to use in Params.Merge.
type ParamsMergeStrategy string

const (
	// Do not merge.
	ParamsMergeStrategyNone ParamsMergeStrategy = "none"
	// Only add new keys.
	ParamsMergeStrategyShallow ParamsMergeStrategy = "shallow"
	// Add new keys, merge existing.
	ParamsMergeStrategyDeep ParamsMergeStrategy = "deep"

	mergeStrategyKey = "_merge"
)

func toMergeStrategy(v interface{}) ParamsMergeStrategy {
	s := ParamsMergeStrategy(cast.ToString(v))
	switch s {
	case ParamsMergeStrategyDeep, ParamsMergeStrategyNone, ParamsMergeStrategyShallow:
		return s
	default:
		return ParamsMergeStrategyDeep
	}
}

// PrepareParams
// * makes all the keys in the given map lower cased and will do so
// * This will modify the map given.
// * Any nested map[interface{}]interface{}, map[string]interface{},map[string]string  will be converted to Params.
// * Any _merge value will be converted to proper type and value.
func PrepareParams(m Params) {
	for k, v := range m {
		var retyped bool
		lKey := strings.ToLower(k)
		if lKey == mergeStrategyKey {
			v = toMergeStrategy(v)
			retyped = true
		} else {
			switch vv := v.(type) {
			case map[interface{}]interface{}:
				var p Params = cast.ToStringMap(v)
				v = p
				PrepareParams(p)
				retyped = true
			case map[string]interface{}:
				var p Params = v.(map[string]interface{})
				v = p
				PrepareParams(p)
				retyped = true
			case map[string]string:
				p := make(Params)
				for k, v := range vv {
					p[k] = v
				}
				v = p
				PrepareParams(p)
				retyped = true
			}
		}

		if retyped || k != lKey {
			delete(m, k)
			m[lKey] = v
		}
	}
}
