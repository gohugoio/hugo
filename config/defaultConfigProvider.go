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
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/spf13/cast"

	"github.com/gohugoio/hugo/common/maps"
)

var (

	// ConfigRootKeysSet contains all of the config map root keys.
	ConfigRootKeysSet = map[string]bool{
		"build":         true,
		"caches":        true,
		"cascade":       true,
		"frontmatter":   true,
		"languages":     true,
		"imaging":       true,
		"markup":        true,
		"mediatypes":    true,
		"menus":         true,
		"minify":        true,
		"module":        true,
		"outputformats": true,
		"params":        true,
		"permalinks":    true,
		"related":       true,
		"sitemap":       true,
		"privacy":       true,
		"security":      true,
		"taxonomies":    true,
	}

	// ConfigRootKeys is a sorted version of ConfigRootKeysSet.
	ConfigRootKeys []string
)

func init() {
	for k := range ConfigRootKeysSet {
		ConfigRootKeys = append(ConfigRootKeys, k)
	}
	sort.Strings(ConfigRootKeys)
}

// New creates a Provider backed by an empty maps.Params.
func New() Provider {
	return &defaultConfigProvider{
		root: make(maps.Params),
	}
}

// NewFrom creates a Provider backed by params.
func NewFrom(params maps.Params) Provider {
	maps.PrepareParams(params)
	return &defaultConfigProvider{
		root: params,
	}
}

// defaultConfigProvider is a Provider backed by a map where all keys are lower case.
// All methods are thread safe.
type defaultConfigProvider struct {
	mu   sync.RWMutex
	root maps.Params

	keyCache sync.Map
}

func (c *defaultConfigProvider) Get(k string) interface{} {
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

func (c *defaultConfigProvider) GetParams(k string) maps.Params {
	v := c.Get(k)
	if v == nil {
		return nil
	}
	return v.(maps.Params)
}

func (c *defaultConfigProvider) GetStringMap(k string) map[string]interface{} {
	v := c.Get(k)
	return maps.ToStringMap(v)
}

func (c *defaultConfigProvider) GetStringMapString(k string) map[string]string {
	v := c.Get(k)
	return maps.ToStringMapString(v)
}

func (c *defaultConfigProvider) GetStringSlice(k string) []string {
	v := c.Get(k)
	return cast.ToStringSlice(v)
}

func (c *defaultConfigProvider) Set(k string, v interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	k = strings.ToLower(k)

	if k == "" {
		if p, ok := maps.ToParamsAndPrepare(v); ok {
			// Set the values directly in root.
			c.root.Set(p)
		} else {
			c.root[k] = v
		}

		return
	}

	switch vv := v.(type) {
	case map[string]interface{}, map[interface{}]interface{}, map[string]string:
		p := maps.MustToParamsAndPrepare(vv)
		v = p
	}

	key, m := c.getNestedKeyAndMap(k, true)
	if m == nil {
		return
	}

	if existing, found := m[key]; found {
		if p1, ok := existing.(maps.Params); ok {
			if p2, ok := v.(maps.Params); ok {
				p1.Set(p2)
				return
			}
		}
	}

	m[key] = v
}

// SetDefaults will set values from params if not already set.
func (c *defaultConfigProvider) SetDefaults(params maps.Params) {
	maps.PrepareParams(params)
	for k, v := range params {
		if _, found := c.root[k]; !found {
			c.root[k] = v
		}
	}
}

func (c *defaultConfigProvider) Merge(k string, v interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	k = strings.ToLower(k)

	const (
		languagesKey = "languages"
		paramsKey    = "params"
		menusKey     = "menus"
	)

	if k == "" {
		rs, f := c.root.GetMergeStrategy()
		if f && rs == maps.ParamsMergeStrategyNone {
			// The user has set a "no merge" strategy on this,
			// nothing more to do.
			return
		}

		if p, ok := maps.ToParamsAndPrepare(v); ok {
			// As there may be keys in p not in root, we need to handle
			// those as a special case.
			var keysToDelete []string
			for kk, vv := range p {
				if pp, ok := vv.(maps.Params); ok {
					if pppi, ok := c.root[kk]; ok {
						ppp := pppi.(maps.Params)
						if kk == languagesKey {
							// Languages is currently a special case.
							// We may have languages with menus or params in the
							// right map that is not present in the left map.
							// With the default merge strategy those items will not
							// be passed over.
							var hasParams, hasMenus bool
							for _, rv := range pp {
								if lkp, ok := rv.(maps.Params); ok {
									_, hasMenus = lkp[menusKey]
									_, hasParams = lkp[paramsKey]
								}
							}

							if hasMenus || hasParams {
								for _, lv := range ppp {
									if lkp, ok := lv.(maps.Params); ok {
										if hasMenus {
											if _, ok := lkp[menusKey]; !ok {
												p := maps.Params{}
												p.SetDefaultMergeStrategy(maps.ParamsMergeStrategyShallow)
												lkp[menusKey] = p
											}
										}
										if hasParams {
											if _, ok := lkp[paramsKey]; !ok {
												p := maps.Params{}
												p.SetDefaultMergeStrategy(maps.ParamsMergeStrategyShallow)
												lkp[paramsKey] = p
											}
										}
									}
								}
							}
						}
						ppp.Merge(pp)
					} else {
						// We need to use the default merge strategy for
						// this key.
						np := make(maps.Params)
						strategy := c.determineMergeStrategy(KeyParams{Key: "", Params: c.root}, KeyParams{Key: kk, Params: np})
						np.SetDefaultMergeStrategy(strategy)
						np.Merge(pp)
						c.root[kk] = np
						if np.IsZero() {
							// Just keep it until merge is done.
							keysToDelete = append(keysToDelete, kk)
						}
					}
				}
			}
			// Merge the rest.
			c.root.MergeRoot(p)
			for _, k := range keysToDelete {
				delete(c.root, k)
			}
		} else {
			panic(fmt.Sprintf("unsupported type %T received in Merge", v))
		}

		return
	}

	switch vv := v.(type) {
	case map[string]interface{}, map[interface{}]interface{}, map[string]string:
		p := maps.MustToParamsAndPrepare(vv)
		v = p
	}

	key, m := c.getNestedKeyAndMap(k, true)
	if m == nil {
		return
	}

	if existing, found := m[key]; found {
		if p1, ok := existing.(maps.Params); ok {
			if p2, ok := v.(maps.Params); ok {
				p1.Merge(p2)
			}
		}
	} else {
		m[key] = v
	}
}

func (c *defaultConfigProvider) WalkParams(walkFn func(params ...KeyParams) bool) {
	var walk func(params ...KeyParams)
	walk = func(params ...KeyParams) {
		if walkFn(params...) {
			return
		}
		p1 := params[len(params)-1]
		i := len(params)
		for k, v := range p1.Params {
			if p2, ok := v.(maps.Params); ok {
				paramsplus1 := make([]KeyParams, i+1)
				copy(paramsplus1, params)
				paramsplus1[i] = KeyParams{Key: k, Params: p2}
				walk(paramsplus1...)
			}
		}
	}
	walk(KeyParams{Key: "", Params: c.root})
}

func (c *defaultConfigProvider) determineMergeStrategy(params ...KeyParams) maps.ParamsMergeStrategy {
	if len(params) == 0 {
		return maps.ParamsMergeStrategyNone
	}

	var (
		strategy   maps.ParamsMergeStrategy
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
		strategy = maps.ParamsMergeStrategyDeep
	case "outputformats", "mediatypes":
		if prevIsRoot {
			strategy = maps.ParamsMergeStrategyShallow
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
			strategy = maps.ParamsMergeStrategyShallow
		}
	default:
		if strategy == "" {
			strategy = maps.ParamsMergeStrategyNone
		}
	}

	return strategy
}

type KeyParams struct {
	Key    string
	Params maps.Params
}

func (c *defaultConfigProvider) SetDefaultMergeStrategy() {
	c.WalkParams(func(params ...KeyParams) bool {
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
			p.SetDefaultMergeStrategy(strategy)
		}
		return false
	})

}

func (c *defaultConfigProvider) getNestedKeyAndMap(key string, create bool) (string, maps.Params) {
	var parts []string
	v, ok := c.keyCache.Load(key)
	if ok {
		parts = v.([]string)
	} else {
		parts = strings.Split(key, ".")
		c.keyCache.Store(key, parts)
	}
	current := c.root
	for i := 0; i < len(parts)-1; i++ {
		next, found := current[parts[i]]
		if !found {
			if create {
				next = make(maps.Params)
				current[parts[i]] = next
			} else {
				return "", nil
			}
		}
		var ok bool
		current, ok = next.(maps.Params)
		if !ok {
			// E.g. a string, not a map that we can store values in.
			return "", nil
		}
	}
	return parts[len(parts)-1], current
}
