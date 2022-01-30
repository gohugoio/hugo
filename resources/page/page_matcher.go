// Copyright 2020 The Hugo Authors. All rights reserved.
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

package page

import (
	"path/filepath"
	"strings"

	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/hugofs/glob"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

// A PageMatcher can be used to match a Page with Glob patterns.
// Note that the pattern matching is case insensitive.
type PageMatcher struct {
	// A Glob pattern matching the content path below /content.
	// Expects Unix-styled slashes.
	// Note that this is the virtual path, so it starts at the mount root
	// with a leading "/".
	Path string

	// A Glob pattern matching the Page's Kind(s), e.g. "{home,section}"
	Kind string

	// A Glob pattern matching the Page's language, e.g. "{en,sv}".
	Lang string
}

// Matches returns whether p matches this matcher.
func (m PageMatcher) Matches(p Page) bool {
	if m.Kind != "" {
		g, err := glob.GetGlob(m.Kind)
		if err == nil && !g.Match(p.Kind()) {
			return false
		}
	}

	if m.Lang != "" {
		g, err := glob.GetGlob(m.Lang)
		if err == nil && !g.Match(p.Lang()) {
			return false
		}
	}

	if m.Path != "" {
		g, err := glob.GetGlob(m.Path)
		// TODO(bep) Path() vs filepath vs leading slash.
		p := strings.ToLower(filepath.ToSlash(p.Pathc()))
		if !(strings.HasPrefix(p, "/")) {
			p = "/" + p
		}
		if err == nil && !g.Match(p) {
			return false
		}
	}

	return true
}

// DecodeCascade decodes in which could be eiter a map or a slice of maps.
func DecodeCascade(in interface{}) (map[PageMatcher]maps.Params, error) {
	m, err := maps.ToSliceStringMap(in)
	if err != nil {
		return map[PageMatcher]maps.Params{
			{}: maps.ToStringMap(in),
		}, nil
	}

	cascade := make(map[PageMatcher]maps.Params)

	for _, vv := range m {
		var m PageMatcher
		if mv, found := vv["_target"]; found {
			err := DecodePageMatcher(mv, &m)
			if err != nil {
				return nil, err
			}
		}
		c, found := cascade[m]
		if found {
			// Merge
			for k, v := range vv {
				if _, found := c[k]; !found {
					c[k] = v
				}
			}
		} else {
			cascade[m] = vv
		}
	}

	return cascade, nil
}

// DecodePageMatcher decodes m into v.
func DecodePageMatcher(m interface{}, v *PageMatcher) error {
	if err := mapstructure.WeakDecode(m, v); err != nil {
		return err
	}

	v.Kind = strings.ToLower(v.Kind)
	if v.Kind != "" {
		g, _ := glob.GetGlob(v.Kind)
		found := false
		for _, k := range kindMap {
			if g.Match(k) {
				found = true
				break
			}
		}
		if !found {
			return errors.Errorf("%q did not match a valid Page Kind", v.Kind)
		}
	}

	v.Path = filepath.ToSlash(strings.ToLower(v.Path))

	return nil
}
