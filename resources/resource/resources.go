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
	"strings"

	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/hugofs/glob"
	"github.com/spf13/cast"
)

var _ ResourceFinder = (*Resources)(nil)

// Resources represents a slice of resources, which can be a mix of different types.
// I.e. both pages and images etc.
type Resources []Resource

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
	namestr = strings.ToLower(namestr)

	// First check the Name.
	// Note that this can be modified by the user in the front matter,
	// also, it does not contain any language code.
	for _, resource := range r {
		if strings.EqualFold(namestr, resource.Name()) {
			return resource
		}
	}

	// Finally, check the original name.
	for _, resource := range r {
		if nop, ok := resource.(NameOriginalProvider); ok {
			if strings.EqualFold(namestr, nop.NameOriginal()) {
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

	patternstr = paths.NormalizePathStringBasic(patternstr)

	g, err := glob.GetGlob(patternstr)
	if err != nil {
		panic(err)
	}

	for _, resource := range r {
		if g.Match(paths.NormalizePathStringBasic(resource.Name())) {
			return resource
		}
	}

	// Finally, check the original name.
	for _, resource := range r {
		if nop, ok := resource.(NameOriginalProvider); ok {
			if g.Match(paths.NormalizePathStringBasic(nop.NameOriginal())) {
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

	g, err := glob.GetGlob(patternstr)
	if err != nil {
		panic(err)
	}

	var matches Resources
	for _, resource := range r {
		if g.Match(strings.ToLower(resource.Name())) {
			matches = append(matches, resource)
		}
	}
	if len(matches) == 0 {
		// 	Fall back to the original name.
		for _, resource := range r {
			if nop, ok := resource.(NameOriginalProvider); ok {
				if g.Match(strings.ToLower(nop.NameOriginal())) {
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

// ResourceFinder provides methods to find Resources.
// Note that GetRemote (as found in resources.GetRemote) is
// not covered by this interface, as this is only available as a global template function.
type ResourceFinder interface {
	// Get locates the Resource with the given name in the current context (e.g. in .Page.Resources).
	//
	// It returns nil if no Resource could found, panics if name is invalid.
	Get(name any) Resource

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
