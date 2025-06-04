// Copyright 2025 The Hugo Authors. All rights reserved.
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

package dimensions

import (
	"cmp"
	"fmt"

	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/common/predicate"
	"github.com/gohugoio/hugo/config"
)

// IntSets holds the ordered sets of integers for the dimensions,
// which is used for fast membership testing of files, resources and pages.
type IntSets struct {
	Languages *maps.OrderedIntSet `mapstructure:"-" json:"-"`
	Versions  *maps.OrderedIntSet `mapstructure:"-" json:"-"`
	Roles     *maps.OrderedIntSet `mapstructure:"-" json:"-"`
}

func (s *IntSets) String() string {
	return fmt.Sprintf("Languages: %v, Versions: %v, Roles: %v", s.Languages, s.Versions, s.Roles)
}

func (s *IntSets) SetFrom(other *IntSets) {
	if other == nil {
		return
	}
	if other.Languages != nil {
		if s.Languages == nil {
			s.Languages = maps.NewOrderedIntSet()
		}
		s.Languages.SetFrom(other.Languages)
	}
	if other.Versions != nil {
		if s.Versions == nil {
			s.Versions = maps.NewOrderedIntSet()
		}
		s.Versions.SetFrom(other.Versions)
	}
	if other.Roles != nil {
		if s.Roles == nil {
			s.Roles = maps.NewOrderedIntSet()
		}
		s.Roles.SetFrom(other.Roles)
	}
}

// NewIntSets creates a new DimensionsIntSets with nil sets for languages, roles, and versions.
func NewIntSets() *IntSets {
	return &IntSets{}
}

// TODO1 name etc.
func NewIntSets2(cfg config.ConfiguredDimensions, applyDefault bool, languages, versions, roles []string) (*IntSets, error) {
	// TODO1
	/*if p.Lang != "" {
		// Merge into the languages slice.
		p.Languages = append(p.Languages, p.Lang)
		p.Languages = hstrings.UniqueStringsReuse(p.Languages)
	}*/

	applyFilter := func(what string, values []string, matcher config.ConfiguredDimension) (*maps.OrderedIntSet, error) {
		if len(values) == 0 {
			if applyDefault {
				result := maps.NewOrderedIntSet()
				result.Set(matcher.IndexDefault())
				return result, nil
			}
			return nil, nil
		}
		var result *maps.OrderedIntSet
		filter, err := predicate.NewFilterFromGlobs(values)
		if err != nil {
			return nil, fmt.Errorf("failed to create filter for %s: %w", what, err)
		}
		for _, pattern := range values {
			iter, err := matcher.IndexMatch(filter)
			if err != nil {
				return nil, fmt.Errorf("failed to match %s %q: %w", what, pattern, err)
			}
			for i := range iter {
				if result == nil {
					result = maps.NewOrderedIntSet()
				}
				result.Set(i)
			}
		}

		return result, nil
	}

	sets := NewIntSets()
	l, err1 := applyFilter("languages", languages, cfg.ConfiguredLanguages)
	v, err2 := applyFilter("versions", versions, cfg.ConfiguredVersions)
	r, err3 := applyFilter("roles", roles, cfg.ConfiguredRoles)

	if err := cmp.Or(err1, err2, err3); err != nil {
		return nil, fmt.Errorf("failed to apply filters: %w", err)
	}
	sets.Languages = l
	sets.Versions = v
	sets.Roles = r

	return sets, nil
}

type DimensionsStringSlices struct {
	Languages []string `mapstructure:"languages" json:"languages"`
	Versions  []string `mapstructure:"versions" json:"versions"`
	Roles     []string `mapstructure:"roles" json:"roles"`
}

func (d DimensionsStringSlices) IsZero() bool {
	return len(d.Languages) == 0 && len(d.Versions) == 0 && len(d.Roles) == 0
}
