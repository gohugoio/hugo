// Copyright 2021 The Hugo Authors. All rights reserved.
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

package config

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/gohugoio/hugo/common/hmaps"
	"github.com/spf13/cast"
)

// New creates a Provider backed by an empty hmaps.Params.
func New() Provider {
	return &defaultConfigProvider{
		root: make(hmaps.Params),
	}
}

// NewFrom creates a Provider backed by params.
func NewFrom(params hmaps.Params) Provider {
	hmaps.PrepareParams(params)
	return &defaultConfigProvider{
		root: params,
	}
}

// defaultConfigProvider is a Provider backed by a map where all keys are lower case.
// All methods are thread safe.
type defaultConfigProvider struct {
	mu   sync.RWMutex
	root hmaps.Params

	keyCache sync.Map
}

func (c *defaultConfigProvider) Get(k string) any {
	if k == "" {
		return c.root
	}
	c.mu.RLock()
	key, m := c.getNestedKeyAndMap(strings.ToLower(k), false)
	if m == nil {
		c.mu.RUnlock()
		return nil
	}
	v := m[key]
	c.mu.RUnlock()
	return v
}

func (c *defaultConfigProvider) GetBool(k string) bool {
	v := c.Get(k)
	return cast.ToBool(v)
}

func (c *defaultConfigProvider) GetInt(k string) int {
	v := c.Get(k)
	return cast.ToInt(v)
}

func (c *defaultConfigProvider) IsSet(k string) bool {
	var found bool
	c.mu.RLock()
	key, m := c.getNestedKeyAndMap(strings.ToLower(k), false)
	if m != nil {
		_, found = m[key]
	}
	c.mu.RUnlock()
	return found
}

func (c *defaultConfigProvider) GetString(k string) string {
	v := c.Get(k)
	return cast.ToString(v)
}

func (c *defaultConfigProvider) GetParams(k string) hmaps.Params {
	v := c.Get(k)
	if v == nil {
		return nil
	}
	return v.(hmaps.Params)
}

func (c *defaultConfigProvider) GetStringMap(k string) map[string]any {
	v := c.Get(k)
	return hmaps.ToStringMap(v)
}

func (c *defaultConfigProvider) GetStringMapString(k string) map[string]string {
	v := c.Get(k)
	return hmaps.ToStringMapString(v)
}

func (c *defaultConfigProvider) GetStringSlice(k string) []string {
	v := c.Get(k)
	return cast.ToStringSlice(v)
}

func (c *defaultConfigProvider) Set(k string, v any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	k = strings.ToLower(k)

	if k == "" {
		if p, err := hmaps.ToParamsAndPrepare(v); err == nil {
			// Set the values directly in root.
			hmaps.SetParams(c.root, p)
		} else {
			c.root[k] = v
		}

		return
	}

	switch vv := v.(type) {
	case map[string]any, map[any]any, map[string]string:
		p := hmaps.MustToParamsAndPrepare(vv)
		v = p
	}

	key, m := c.getNestedKeyAndMap(k, true)
	if m == nil {
		return
	}

	if existing, found := m[key]; found {
		if p1, ok := existing.(hmaps.Params); ok {
			if p2, ok := v.(hmaps.Params); ok {
				hmaps.SetParams(p1, p2)
				return
			}
		}
	}

	m[key] = v
}

// SetDefaults will set values from params if not already set.
func (c *defaultConfigProvider) SetDefaults(params hmaps.Params) {
	hmaps.PrepareParams(params)
	for k, v := range params {
		if _, found := c.root[k]; !found {
			c.root[k] = v
		}
	}
}

func (c *defaultConfigProvider) Merge(k string, v any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	k = strings.ToLower(k)

	if k == "" {
		rs, f := c.root.GetMergeStrategy()
		if f && rs == hmaps.ParamsMergeStrategyNone {
			// The user has set a "no merge" strategy on this,
			// nothing more to do.
			return
		}

		if p, err := hmaps.ToParamsAndPrepare(v); err == nil {
			// As there may be keys in p not in root, we need to handle
			// those as a special case.
			var keysToDelete []string
			for kk, vv := range p {
				if pp, ok := vv.(hmaps.Params); ok {
					if pppi, ok := c.root[kk]; ok {
						ppp := pppi.(hmaps.Params)
						hmaps.MergeParamsWithStrategy("", ppp, pp)
					} else {
						// We need to use the default merge strategy for
						// this key.
						np := make(hmaps.Params)
						strategy := c.determineMergeStrategy(hmaps.KeyParams{Key: "", Params: c.root}, hmaps.KeyParams{Key: kk, Params: np})
						np.SetMergeStrategy(strategy)
						hmaps.MergeParamsWithStrategy("", np, pp)
						c.root[kk] = np
						if np.IsZero() {
							// Just keep it until merge is done.
							keysToDelete = append(keysToDelete, kk)
						}
					}
				}
			}
			// Merge the rest.
			hmaps.MergeParams(c.root, p)
			for _, k := range keysToDelete {
				delete(c.root, k)
			}
		} else {
			panic(fmt.Sprintf("unsupported type %T received in Merge", v))
		}

		return
	}

	switch vv := v.(type) {
	case map[string]any, map[any]any, map[string]string:
		p := hmaps.MustToParamsAndPrepare(vv)
		v = p
	}

	key, m := c.getNestedKeyAndMap(k, true)
	if m == nil {
		return
	}

	if existing, found := m[key]; found {
		if p1, ok := existing.(hmaps.Params); ok {
			if p2, ok := v.(hmaps.Params); ok {
				hmaps.MergeParamsWithStrategy("", p1, p2)
			}
		}
	} else {
		m[key] = v
	}
}

func (c *defaultConfigProvider) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	var keys []string
	for k := range c.root {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (c *defaultConfigProvider) WalkParams(walkFn func(params ...hmaps.KeyParams) bool) {
	maxDepth := 1000
	var walk func(depth int, params ...hmaps.KeyParams)
	walk = func(depth int, params ...hmaps.KeyParams) {
		if depth > maxDepth {
			panic(errors.New("max depth exceeded"))
		}
		if walkFn(params...) {
			return
		}
		p1 := params[len(params)-1]
		i := len(params)
		for k, v := range p1.Params {
			if p2, ok := v.(hmaps.Params); ok {
				paramsplus1 := make([]hmaps.KeyParams, i+1)
				copy(paramsplus1, params)
				paramsplus1[i] = hmaps.KeyParams{Key: k, Params: p2}
				walk(depth+1, paramsplus1...)
			}
		}
	}
	walk(0, hmaps.KeyParams{Key: "", Params: c.root})
}

func (c *defaultConfigProvider) determineMergeStrategy(params ...hmaps.KeyParams) hmaps.ParamsMergeStrategy {
	if len(params) == 0 {
		return hmaps.ParamsMergeStrategyNone
	}

	var (
		strategy   hmaps.ParamsMergeStrategy
		prevIsRoot bool
		curr       = params[len(params)-1]
	)

	if len(params) > 1 {
		prev := params[len(params)-2]
		prevIsRoot = prev.Key == ""

		// Inherit from parent (but not from the root unless it's set by user).
		s, found := prev.Params.GetMergeStrategy()
		if !prevIsRoot && !found {
			panic("invalid state, merge strategy not set on parent")
		}
		if found || !prevIsRoot {
			strategy = s
		}
	}

	switch curr.Key {
	case "":
	// Don't set a merge strategy on the root unless set by user.
	// This will be handled as a special case.
	case "params":
		strategy = hmaps.ParamsMergeStrategyDeep
	case "outputformats", "mediatypes":
		if prevIsRoot {
			strategy = hmaps.ParamsMergeStrategyShallow
		}
	case "menus":
		isMenuKey := prevIsRoot
		if !isMenuKey {
			// Can also be set below languages.
			// root > languages > en > menus
			if len(params) == 4 && params[1].Key == "languages" {
				isMenuKey = true
			}
		}
		if isMenuKey {
			strategy = hmaps.ParamsMergeStrategyShallow
		}
	default:
		if strategy == "" {
			strategy = hmaps.ParamsMergeStrategyNone
		}
	}

	return strategy
}

func (c *defaultConfigProvider) SetDefaultMergeStrategy() {
	c.WalkParams(func(params ...hmaps.KeyParams) bool {
		if len(params) == 0 {
			return false
		}
		p := params[len(params)-1].Params
		var found bool
		if _, found = p.GetMergeStrategy(); found {
			// Set by user.
			return false
		}
		strategy := c.determineMergeStrategy(params...)
		if strategy != "" {
			p.SetMergeStrategy(strategy)
		}
		return false
	})
}

func (c *defaultConfigProvider) getNestedKeyAndMap(key string, create bool) (string, hmaps.Params) {
	var parts []string
	v, ok := c.keyCache.Load(key)
	if ok {
		parts = v.([]string)
	} else {
		parts = strings.Split(key, ".")
		c.keyCache.Store(key, parts)
	}
	current := c.root
	for i := range len(parts) - 1 {
		next, found := current[parts[i]]
		if !found {
			if create {
				next = make(hmaps.Params)
				current[parts[i]] = next
			} else {
				return "", nil
			}
		}
		var ok bool
		current, ok = next.(hmaps.Params)
		if !ok {
			// E.g. a string, not a map that we can store values in.
			return "", nil
		}
	}
	return parts[len(parts)-1], current
}
