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
type Params map[string]any

// KeyParams is an utility struct for the WalkParams method.
type KeyParams struct {
	Key    string
	Params Params
}

// GetNested does a lower case and nested search in this map.
// It will return nil if none found.
// Make all of these methods internal somehow.
func (p Params) GetNested(indices ...string) any {
	v, _, _ := getNested(p, indices)
	return v
}

// SetParams overwrites values in dst with values in src for common or new keys.
// This is done recursively.
func SetParams(dst, src Params) {
	for k, v := range src {
		vv, found := dst[k]
		if !found {
			dst[k] = v
		} else {
			switch vvv := vv.(type) {
			case Params:
				if pv, ok := v.(Params); ok {
					SetParams(vvv, pv)
				} else {
					dst[k] = v
				}
			default:
				dst[k] = v
			}
		}
	}
}

// IsZero returns true if p is considered empty.
func (p Params) IsZero() bool {
	if len(p) == 0 {
		return true
	}

	if len(p) > 1 {
		return false
	}

	for k := range p {
		return k == MergeStrategyKey
	}

	return false
}

// MergeParamsWithStrategy transfers values from src to dst for new keys using the merge strategy given.
// This is done recursively.
func MergeParamsWithStrategy(strategy string, dst, src Params) {
	dst.merge(ParamsMergeStrategy(strategy), src)
}

// MergeParams transfers values from src to dst for new keys using the merge encoded in dst.
// This is done recursively.
func MergeParams(dst, src Params) {
	ms, _ := dst.GetMergeStrategy()
	dst.merge(ms, src)
}

func (p Params) merge(ps ParamsMergeStrategy, pp Params) {
	ns, found := p.GetMergeStrategy()

	ms := ns
	if !found && ps != "" {
		ms = ps
	}

	noUpdate := ms == ParamsMergeStrategyNone
	noUpdate = noUpdate || (ps != "" && ps == ParamsMergeStrategyShallow)

	for k, v := range pp {

		if k == MergeStrategyKey {
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

// For internal use.
func (p Params) GetMergeStrategy() (ParamsMergeStrategy, bool) {
	if v, found := p[MergeStrategyKey]; found {
		if s, ok := v.(ParamsMergeStrategy); ok {
			return s, true
		}
	}
	return ParamsMergeStrategyShallow, false
}

// For internal use.
func (p Params) DeleteMergeStrategy() bool {
	if _, found := p[MergeStrategyKey]; found {
		delete(p, MergeStrategyKey)
		return true
	}
	return false
}

// For internal use.
func (p Params) SetMergeStrategy(s ParamsMergeStrategy) {
	switch s {
	case ParamsMergeStrategyDeep, ParamsMergeStrategyNone, ParamsMergeStrategyShallow:
	default:
		panic(fmt.Sprintf("invalid merge strategy %q", s))
	}
	p[MergeStrategyKey] = s
}

func getNested(m map[string]any, indices []string) (any, string, map[string]any) {
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
	case map[string]any:
		return getNested(m2, indices[1:])
	default:
		return nil, "", nil
	}
}

// GetNestedParam gets the first match of the keyStr in the candidates given.
// It will first try the exact match and then try to find it as a nested map value,
// using the given separator, e.g. "mymap.name".
// It assumes that all the maps given have lower cased keys.
func GetNestedParam(keyStr, separator string, candidates ...Params) (any, error) {
	keyStr = strings.ToLower(keyStr)

	// Try exact match first
	for _, m := range candidates {
		if v, ok := m[keyStr]; ok {
			return v, nil
		}
	}

	keySegments := strings.Split(keyStr, separator)
	for _, m := range candidates {
		if v := m.GetNested(keySegments...); v != nil {
			return v, nil
		}
	}

	return nil, nil
}

func GetNestedParamFn(keyStr, separator string, lookupFn func(key string) any) (any, string, map[string]any, error) {
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
	case map[string]any:
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

	MergeStrategyKey = "_merge"
)

// CleanConfigStringMapString removes any processing instructions from m,
// m will never be modified.
func CleanConfigStringMapString(m map[string]string) map[string]string {
	if len(m) == 0 {
		return m
	}
	if _, found := m[MergeStrategyKey]; !found {
		return m
	}
	// Create a new map and copy all the keys except the merge strategy key.
	m2 := make(map[string]string, len(m)-1)
	for k, v := range m {
		if k != MergeStrategyKey {
			m2[k] = v
		}
	}
	return m2
}

// CleanConfigStringMap is the same as CleanConfigStringMapString but for
// map[string]any.
func CleanConfigStringMap(m map[string]any) map[string]any {
	if len(m) == 0 {
		return m
	}
	if _, found := m[MergeStrategyKey]; !found {
		return m
	}
	// Create a new map and copy all the keys except the merge strategy key.
	m2 := make(map[string]any, len(m)-1)
	for k, v := range m {
		if k != MergeStrategyKey {
			m2[k] = v
		}
		switch v2 := v.(type) {
		case map[string]any:
			m2[k] = CleanConfigStringMap(v2)
		case Params:
			var p Params = CleanConfigStringMap(v2)
			m2[k] = p
		case map[string]string:
			m2[k] = CleanConfigStringMapString(v2)
		}

	}
	return m2
}

func toMergeStrategy(v any) ParamsMergeStrategy {
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
		if lKey == MergeStrategyKey {
			v = toMergeStrategy(v)
			retyped = true
		} else {
			switch vv := v.(type) {
			case map[any]any:
				var p Params = cast.ToStringMap(v)
				v = p
				PrepareParams(p)
				retyped = true
			case map[string]any:
				var p Params = v.(map[string]any)
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
