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

package resource

import (
	"fmt"
	"strings"

	"github.com/gohugoio/hugo/hugofs/glob"
)

// Resources represents a slice of resources, which can be a mix of different types.
// I.e. both pages and images etc.
type Resources []Resource

// ResourcesConverter converts a given slice of Resource objects to Resources.
type ResourcesConverter interface {
	ToResources() Resources
}

// ByType returns resources of a given resource type (ie. "image").
func (r Resources) ByType(tp string) Resources {
	var filtered Resources

	for _, resource := range r {
		if resource.ResourceType() == tp {
			filtered = append(filtered, resource)
		}
	}
	return filtered
}

// GetMatch finds the first Resource matching the given pattern, or nil if none found.
// See Match for a more complete explanation about the rules used.
func (r Resources) GetMatch(pattern string) Resource {
	g, err := glob.GetGlob(pattern)
	if err != nil {
		return nil
	}

	for _, resource := range r {
		if g.Match(strings.ToLower(resource.Name())) {
			return resource
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
func (r Resources) Match(pattern string) Resources {
	g, err := glob.GetGlob(pattern)
	if err != nil {
		return nil
	}

	var matches Resources
	for _, resource := range r {
		if g.Match(strings.ToLower(resource.Name())) {
			matches = append(matches, resource)
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
func (r Resources) MergeByLanguageInterface(in interface{}) (interface{}, error) {
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
