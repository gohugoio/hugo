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

package sitematrix

import (
	"cmp"
	"fmt"

	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/common/predicate"
	"github.com/gohugoio/hugo/common/types"
	"github.com/gohugoio/hugo/hugofs/glob"
)

var (
	_ VectorProvider       = &IntSets{}
	_ types.WeightProvider = &IntSets{}
)

// IntSets holds the ordered sets of integers for the dimensions,
// which is used for fast membership testing of files, resources and pages.
type IntSets struct {
	weight    int
	ordinal   int                 // Any non-zero value will be considered when sorting, lesser weights comes first.
	Languages *maps.OrderedIntSet `mapstructure:"-" json:"-"`
	Versions  *maps.OrderedIntSet `mapstructure:"-" json:"-"`
	Roles     *maps.OrderedIntSet `mapstructure:"-" json:"-"`
}

func (s *IntSets) String() string {
	return fmt.Sprintf("Languages: %v, Versions: %v, Roles: %v", s.Languages, s.Versions, s.Roles)
}

func (s *IntSets) Weight() int {
	if s == nil {
		return 0
	}
	return s.weight
}

func (s *IntSets) Ordinal() int {
	if s == nil {
		return 0
	}
	return s.ordinal
}

// HasVector checks if the given vector is contained in the sets.
func (s *IntSets) HasVector(v Vector) bool {
	if s == nil {
		return false
	}
	if !s.Languages.Has(v.Language()) {
		return false
	}
	if !s.Versions.Has(v.Version()) {
		return false
	}
	if !s.Roles.Has(v.Role()) {
		return false
	}
	return true
}

func (s *IntSets) FirstVector() Vector {
	if s.LenVectors() == 0 {
		panic("no vectors available")
	}

	return Vector{
		s.Languages.Get(0),
		s.Versions.Get(0),
		s.Roles.Get(0),
	}
}

func (s *IntSets) LenVectors() int {
	if s == nil {
		return 0
	}
	return s.Languages.Len() * s.Versions.Len() * s.Roles.Len()
}

// The reason we don't use iter.Seq is https://github.com/golang/go/issues/69015
// This is 60% faster and allocation free.
// The yield function should return false to stop iteration.
func (s *IntSets) ForEeachVector(yield func(v Vector) bool) bool {
	if s.LenVectors() == 0 {
		return true
	}

	b := s.Languages.ForEachKey(func(lang int) bool {
		return s.Versions.ForEachKey(func(ver int) bool {
			return s.Roles.ForEachKey(func(role int) bool {
				if !yield(Vector{lang, ver, role}) {
					return false
				}
				return true
			})
		})
	})

	return b
}

func (s *IntSets) EqualsVector(other VectorProvider) bool {
	if s == nil && other == nil {
		return true
	}
	if s == nil || other == nil {
		return false
	}
	if s == other {
		return true
	}
	if s.LenVectors() != other.LenVectors() {
		return false
	}

	return other.ForEeachVector(func(v Vector) bool {
		return s.HasVector(v)
	})
}

// ApplyDefaultsIfNotSet applies default values to the IntSets if they are not already set.
func (s *IntSets) SetDefaultsIfNotSet(cfg ConfiguredDimensions) {
	if s.Languages == nil {
		s.Languages = maps.NewOrderedIntSet()
		s.Languages.Set(cfg.ConfiguredLanguages.IndexDefault())
	}
	if s.Versions == nil {
		s.Versions = maps.NewOrderedIntSet()
		s.Versions.Set(cfg.ConfiguredVersions.IndexDefault())
	}
	if s.Roles == nil {
		s.Roles = maps.NewOrderedIntSet()
		s.Roles.Set(cfg.ConfiguredRoles.IndexDefault())
	}
}

func (s *IntSets) SetFromOtherIfNotSet(other *IntSets) {
	if other == nil {
		return
	}
	if s.Languages == nil && other.Languages != nil {
		s.Languages = maps.NewOrderedIntSet()
		s.Languages.SetFrom(other.Languages)
	}
	if s.Versions == nil && other.Versions != nil {
		s.Versions = maps.NewOrderedIntSet()
		s.Versions.SetFrom(other.Versions)
	}
	if s.Roles == nil && other.Roles != nil {
		s.Roles = maps.NewOrderedIntSet()
		s.Roles.SetFrom(other.Roles)
	}
}

func (s IntSets) WithDefaultsIfNotSet(cfg ConfiguredDimensions) *IntSets {
	s.SetDefaultsIfNotSet(cfg)
	return &s
}

// WithLanguageIndex replaces the current language set with a single language index.
func (s IntSets) WithLanguageIndex(i int) *IntSets {
	s.Languages = maps.NewOrderedIntSet(i)
	return &s
}

type IntSetsConfig struct {
	Cfg          ConfiguredDimensions
	Weight       int
	Ordinal      int
	ApplyDefault bool
	Languages    []string
	Versions     []string
	Roles        []string
}

// NewIntSets creates a new DimensionsIntSets with nil sets for languages, roles, and versions.
func NewIntSets(weight, ordinal int) *IntSets {
	return &IntSets{weight: weight, ordinal: ordinal}
}

// NewIntSetsFromConfig creates a new IntSets from the given IntSetsConfig.
// It applies the filters based on the provided languages, versions, and roles.
func NewIntSetsFromConfig(cfg IntSetsConfig) (*IntSets, error) {
	applyFilter := func(what string, values []string, matcher ConfiguredDimension) (*maps.OrderedIntSet, error) {
		if len(values) == 0 {
			if cfg.ApplyDefault {
				result := maps.NewOrderedIntSet()
				result.Set(matcher.IndexDefault())
				return result, nil
			}
			return nil, nil
		}
		var result *maps.OrderedIntSet
		// Dot separated globs.
		filter, err := predicate.NewFilterFromGlobs(values, glob.GetGlobDot)
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

	sets := NewIntSets(cfg.Weight, cfg.Ordinal)
	l, err1 := applyFilter("languages", cfg.Languages, cfg.Cfg.ConfiguredLanguages)
	v, err2 := applyFilter("versions", cfg.Versions, cfg.Cfg.ConfiguredVersions)
	r, err3 := applyFilter("roles", cfg.Roles, cfg.Cfg.ConfiguredRoles)

	if err := cmp.Or(err1, err2, err3); err != nil {
		return nil, fmt.Errorf("failed to apply filters: %w", err)
	}
	sets.Languages = l
	sets.Versions = v
	sets.Roles = r

	return sets, nil
}

type SitesConfig struct {
	// Matrix defines the main build matrix.
	Matrix StringSlices `mapstructure:"matrix" json:"matrix"`
	// Fallbacks defines the fallback matrix.
	Fallbacks StringSlices `mapstructure:"fallbacks" json:"fallbacks"`
}

// IsZero returns true if all slices are empty.
func (s SitesConfig) IsZero() bool {
	return s.Matrix.IsZero() && s.Fallbacks.IsZero()
}

type StringSlices struct {
	Languages []string `mapstructure:"languages" json:"languages"`
	Versions  []string `mapstructure:"versions" json:"versions"`
	Roles     []string `mapstructure:"roles" json:"roles"`
}

func (d StringSlices) IsZero() bool {
	return len(d.Languages) == 0 && len(d.Versions) == 0 && len(d.Roles) == 0
}

type ConfiguredDimension interface {
	predicate.IndexMatcher
	IndexDefault() int
	ResolveName(int) string
}

type ConfiguredDimensions struct {
	ConfiguredLanguages ConfiguredDimension
	ConfiguredVersions  ConfiguredDimension
	ConfiguredRoles     ConfiguredDimension
}

func (c ConfiguredDimensions) ResolveNames(v Vector) types.Strings3 {
	return types.Strings3{
		c.ConfiguredLanguages.ResolveName(v.Language()),
		c.ConfiguredVersions.ResolveName(v.Version()),
		c.ConfiguredRoles.ResolveName(v.Role()),
	}
}
