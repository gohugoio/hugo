// Copyright 2024 The Hugo Authors. All rights reserved.
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

// Package resource contains Resource related types.
package resource

import (
	"fmt"
	"path"
	"strings"

	"github.com/gohugoio/hugo/common/hreflect"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/hugofs/glob"
	"github.com/spf13/cast"
)

var _ ResourceFinder = (*Resources)(nil)

// Resources represents a slice of resources, which can be a mix of different types.
// I.e. both pages and images etc.
type Resources []Resource

// Mount mounts the given resources from base to the given target path.
// Note that leading slashes in target marks an absolute path.
// This method is currently only useful in js.Batch.
func (r Resources) Mount(base, target string) ResourceGetter {
	return resourceGetterFunc(func(namev any) Resource {
		name1, err := cast.ToStringE(namev)
		if err != nil {
			panic(err)
		}

		isTargetAbs := strings.HasPrefix(target, "/")

		if target != "" {
			name1 = strings.TrimPrefix(name1, target)
			if !isTargetAbs {
				name1 = paths.TrimLeading(name1)
			}
		}

		if base != "" && isTargetAbs {
			name1 = path.Join(base, name1)
		}

		for _, res := range r {
			name2 := res.Name()

			if base != "" && !isTargetAbs {
				name2 = paths.TrimLeading(strings.TrimPrefix(name2, base))
			}

			if strings.EqualFold(name1, name2) {
				return res
			}

		}

		return nil
	})
}

type ResourcesProvider interface {
	// Resources returns a list of all resources.
	Resources() Resources
}

// var _ resource.ResourceFinder = (*Namespace)(nil)
// ResourcesConverter converts a given slice of Resource objects to Resources.
type ResourcesConverter interface {
	// For internal use.
	ToResources() Resources
}

// ByType returns resources of a given resource type (e.g. "image").
func (r Resources) ByType(typ any) Resources {
	tpstr, err := cast.ToStringE(typ)
	if err != nil {
		panic(err)
	}
	var filtered Resources

	for _, resource := range r {
		if resource.ResourceType() == tpstr {
			filtered = append(filtered, resource)
		}
	}
	return filtered
}

// Get locates the name given in Resources.
// The search is case insensitive.
func (r Resources) Get(name any) Resource {
	if r == nil {
		return nil
	}
	namestr, err := cast.ToStringE(name)
	if err != nil {
		panic(err)
	}

	isDotCurrent := strings.HasPrefix(namestr, "./")
	if isDotCurrent {
		namestr = strings.TrimPrefix(namestr, "./")
	} else {
		namestr = paths.AddLeadingSlash(namestr)
	}

	check := func(name string) bool {
		if !isDotCurrent {
			name = paths.AddLeadingSlash(name)
		}
		return strings.EqualFold(namestr, name)
	}

	// First check the Name.
	// Note that this can be modified by the user in the front matter,
	// also, it does not contain any language code.
	for _, resource := range r {
		if check(resource.Name()) {
			return resource
		}
	}

	// Finally, check the normalized name.
	for _, resource := range r {
		if nop, ok := resource.(NameNormalizedProvider); ok {
			if check(nop.NameNormalized()) {
				return resource
			}
		}
	}

	return nil
}

// GetMatch finds the first Resource matching the given pattern, or nil if none found.
// See Match for a more complete explanation about the rules used.
func (r Resources) GetMatch(pattern any) Resource {
	patternstr, err := cast.ToStringE(pattern)
	if err != nil {
		panic(err)
	}

	g, err := glob.GetGlob(paths.AddLeadingSlash(patternstr))
	if err != nil {
		panic(err)
	}

	for _, resource := range r {
		if g.Match(paths.AddLeadingSlash(resource.Name())) {
			return resource
		}
	}

	// Finally, check the normalized name.
	for _, resource := range r {
		if nop, ok := resource.(NameNormalizedProvider); ok {
			if g.Match(paths.AddLeadingSlash(nop.NameNormalized())) {
				return resource
			}
		}
	}

	return nil
}

// Match gets all resources matching the given base filename prefix, e.g
// "*.png" will match all png files. The "*" does not match path delimiters (/),
// so if you organize your resources in sub-folders, you need to be explicit about it, e.g.:
// "images/*.png". To match any PNG image anywhere in the bundle you can do "**.png", and
// to match all PNG images below the images folder, use "images/**.jpg".
// The matching is case insensitive.
// Match matches by using the value of Resource.Name, which, by default, is a filename with
// path relative to the bundle root with Unix style slashes (/) and no leading slash, e.g. "images/logo.png".
// See https://github.com/gobwas/glob for the full rules set.
func (r Resources) Match(pattern any) Resources {
	patternstr, err := cast.ToStringE(pattern)
	if err != nil {
		panic(err)
	}

	g, err := glob.GetGlob(paths.AddLeadingSlash(patternstr))
	if err != nil {
		panic(err)
	}

	var matches Resources
	for _, resource := range r {
		if g.Match(paths.AddLeadingSlash(resource.Name())) {
			matches = append(matches, resource)
		}
	}
	if len(matches) == 0 {
		// 	Fall back to the normalized name.
		for _, resource := range r {
			if nop, ok := resource.(NameNormalizedProvider); ok {
				if g.Match(paths.AddLeadingSlash(nop.NameNormalized())) {
					matches = append(matches, resource)
				}
			}
		}
	}
	return matches
}

type translatedResource interface {
	TranslationKey() string
}

// MergeByLanguage adds missing translations in r1 from r2.
func (r Resources) MergeByLanguage(r2 Resources) Resources {
	result := append(Resources(nil), r...)
	m := make(map[string]bool)
	for _, rr := range r {
		if translated, ok := rr.(translatedResource); ok {
			m[translated.TranslationKey()] = true
		}
	}

	for _, rr := range r2 {
		if translated, ok := rr.(translatedResource); ok {
			if _, found := m[translated.TranslationKey()]; !found {
				result = append(result, rr)
			}
		}
	}
	return result
}

// MergeByLanguageInterface is the generic version of MergeByLanguage. It
// is here just so it can be called from the tpl package.
// This is for internal use.
func (r Resources) MergeByLanguageInterface(in any) (any, error) {
	r2, ok := in.(Resources)
	if !ok {
		return nil, fmt.Errorf("%T cannot be merged by language", in)
	}
	return r.MergeByLanguage(r2), nil
}

// Source is an internal template and not meant for use in the templates. It
// may change without notice.
type Source interface {
	Publish() error
}

type ResourceGetter interface {
	// Get locates the Resource with the given name in the current context (e.g. in .Page.Resources).
	//
	// It returns nil if no Resource could found, panics if name is invalid.
	Get(name any) Resource
}

type IsProbablySameResourceGetter interface {
	IsProbablySameResourceGetter(other ResourceGetter) bool
}

// StaleInfoResourceGetter is a ResourceGetter that also provides information about
// whether the underlying resources are stale.
type StaleInfoResourceGetter interface {
	StaleInfo
	ResourceGetter
}

type resourceGetterFunc func(name any) Resource

func (f resourceGetterFunc) Get(name any) Resource {
	return f(name)
}

// ResourceFinder provides methods to find Resources.
// Note that GetRemote (as found in resources.GetRemote) is
// not covered by this interface, as this is only available as a global template function.
type ResourceFinder interface {
	ResourceGetter

	// GetMatch finds the first Resource matching the given pattern, or nil if none found.
	//
	// See Match for a more complete explanation about the rules used.
	//
	// It returns nil if no Resource could found, panics if pattern is invalid.
	GetMatch(pattern any) Resource

	// Match gets all resources matching the given base path prefix, e.g
	// "*.png" will match all png files. The "*" does not match path delimiters (/),
	// so if you organize your resources in sub-folders, you need to be explicit about it, e.g.:
	// "images/*.png". To match any PNG image anywhere in the bundle you can do "**.png", and
	// to match all PNG images below the images folder, use "images/**.jpg".
	//
	// The matching is case insensitive.
	//
	// Match matches by using a relative pathwith Unix style slashes (/) and no
	// leading slash, e.g. "images/logo.png".
	//
	// See https://github.com/gobwas/glob for the full rules set.
	//
	// See Match for a more complete explanation about the rules used.
	//
	// It returns nil if no Resources could found, panics if pattern is invalid.
	Match(pattern any) Resources

	// ByType returns resources of a given resource type (e.g. "image").
	// It returns nil if no Resources could found, panics if typ is invalid.
	ByType(typ any) Resources
}

// NewCachedResourceGetter creates a new ResourceGetter from the given objects.
// If multiple objects are provided, they are merged into one where
// the first match wins.
func NewCachedResourceGetter(os ...any) *cachedResourceGetter {
	var getters multiResourceGetter
	for _, o := range os {
		if g, ok := unwrapResourceGetter(o); ok {
			getters = append(getters, g)
		}
	}

	return &cachedResourceGetter{
		cache:    maps.NewCache[string, Resource](),
		delegate: getters,
	}
}

type multiResourceGetter []ResourceGetter

func (m multiResourceGetter) Get(name any) Resource {
	for _, g := range m {
		if res := g.Get(name); res != nil {
			return res
		}
	}
	return nil
}

var (
	_ ResourceGetter               = (*cachedResourceGetter)(nil)
	_ IsProbablySameResourceGetter = (*cachedResourceGetter)(nil)
)

type cachedResourceGetter struct {
	cache    *maps.Cache[string, Resource]
	delegate ResourceGetter
}

func (c *cachedResourceGetter) Get(name any) Resource {
	namestr, err := cast.ToStringE(name)
	if err != nil {
		panic(err)
	}
	v, _ := c.cache.GetOrCreate(namestr, func() (Resource, error) {
		v := c.delegate.Get(name)
		return v, nil
	})
	return v
}

func (c *cachedResourceGetter) IsProbablySameResourceGetter(other ResourceGetter) bool {
	isProbablyEq := true
	c.cache.ForEeach(func(k string, v Resource) bool {
		if v != other.Get(k) {
			isProbablyEq = false
			return false
		}
		return true
	})

	return isProbablyEq
}

func unwrapResourceGetter(v any) (ResourceGetter, bool) {
	if v == nil {
		return nil, false
	}
	switch vv := v.(type) {
	case ResourceGetter:
		return vv, true
	case ResourcesProvider:
		return vv.Resources(), true
	case func(name any) Resource:
		return resourceGetterFunc(vv), true
	default:
		vvv, ok := hreflect.ToSliceAny(v)
		if !ok {
			return nil, false
		}
		var getters multiResourceGetter
		for _, vv := range vvv {
			if g, ok := unwrapResourceGetter(vv); ok {
				getters = append(getters, g)
			}
		}
		return getters, len(getters) > 0
	}
}
