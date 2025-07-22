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

var _ VectorProvider = &IntSets{}

// IntSets holds the ordered sets of integers for the dimensions,
// which is used for fast membership testing of files, resources and pages.
type IntSets struct {
	ordinal   int                 // Any non-zero value will be considered when sorting, lesser weights comes first.
	languages *maps.OrderedIntSet `mapstructure:"-" json:"-"`
	versions  *maps.OrderedIntSet `mapstructure:"-" json:"-"`
	roles     *maps.OrderedIntSet `mapstructure:"-" json:"-"`
}

func (s *IntSets) String() string {
	return fmt.Sprintf("Languages: %v, Versions: %v, Roles: %v", s.languages, s.versions, s.roles)
}

func (s *IntSets) Ordinal() int {
	if s == nil {
		return 0
	}
	return s.ordinal
}

func (s *IntSets) KeysSorted() ([]int, []int, []int) {
	if s == nil {
		return nil, nil, nil
	}
	languages := s.languages.KeysSorted()
	versions := s.versions.KeysSorted()
	roles := s.roles.KeysSorted()
	return languages, versions, roles
}

func (s *IntSets) HasLanguage(lang int) bool {
	if s == nil {
		return false
	}
	return s.languages.Has(lang)
}

func (s *IntSets) HasVersion(ver int) bool {
	if s == nil {
		return false
	}
	return s.versions.Has(ver)
}

func (s *IntSets) HasRole(role int) bool {
	if s == nil {
		return false
	}
	return s.roles.Has(role)
}

// HasVector checks if the given vector is contained in the sets.
func (s *IntSets) HasVector(v Vector) bool {
	if s == nil {
		return false
	}
	if !s.languages.Has(v.Language()) {
		return false
	}
	if !s.versions.Has(v.Version()) {
		return false
	}
	if !s.roles.Has(v.Role()) {
		return false
	}
	return true
}

func (s *IntSets) FirstVector() Vector {
	if s.LenVectors() == 0 {
		panic("no vectors available")
	}

	return Vector{
		s.languages.Get(0),
		s.versions.Get(0),
		s.roles.Get(0),
	}
}

func (s *IntSets) LenVectors() int {
	if s == nil {
		return 0
	}
	return s.languages.Len() * s.versions.Len() * s.roles.Len()
}

// The reason we don't use iter.Seq is https://github.com/golang/go/issues/69015
// This is 60% faster and allocation free.
// The yield function should return false to stop iteration.
func (s *IntSets) ForEeachVector(yield func(v Vector) bool) bool {
	if s.LenVectors() == 0 {
		return true
	}

	b := s.languages.ForEachKey(func(lang int) bool {
		return s.versions.ForEachKey(func(ver int) bool {
			return s.roles.ForEachKey(func(role int) bool {
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
	if s.languages == nil {
		s.languages = maps.NewOrderedIntSet()
		s.languages.Set(cfg.ConfiguredLanguages.IndexDefault())
	}
	if s.versions == nil {
		s.versions = maps.NewOrderedIntSet()
		s.versions.Set(cfg.ConfiguredVersions.IndexDefault())
	}
	if s.roles == nil {
		s.roles = maps.NewOrderedIntSet()
		s.roles.Set(cfg.ConfiguredRoles.IndexDefault())
	}
}

func (s *IntSets) SetFromOtherIfNotSet(other *IntSets) {
	if other == nil {
		return
	}
	if s.languages == nil && other.languages != nil {
		s.languages = maps.NewOrderedIntSet()
		s.languages.SetFrom(other.languages)
	}
	if s.versions == nil && other.versions != nil {
		s.versions = maps.NewOrderedIntSet()
		s.versions.SetFrom(other.versions)
	}
	if s.roles == nil && other.roles != nil {
		s.roles = maps.NewOrderedIntSet()
		s.roles.SetFrom(other.roles)
	}
}

func (s IntSets) WithOrdinal(i int) *IntSets {
	s.ordinal = i
	return &s
}

func (s IntSets) Clone() *IntSets {
	if s.languages == nil && s.versions == nil && s.roles == nil {
		return nil
	}
	s.languages = s.languages.Clone()
	s.versions = s.versions.Clone()
	s.roles = s.roles.Clone()
	return &s
}

// Complement returns a new IntSets that is the complement of the IntSets passed in is.
// This will return nil if the resulting set is empty.
func (s *IntSets) Complement(is ...*IntSets) *IntSets {
	if len(is) == 0 || (len(is) == 1 && is[0] == s) {
		return nil
	}

	// If all keys in s are present in is, we return nil.
	var notAllPresent bool
	for _, v := range is {
		s.ForEeachVector(func(vec Vector) bool {
			if !v.HasVector(vec) {
				notAllPresent = true
				return false
			}
			return true
		})
	}

	if !notAllPresent {
		return nil
	}

	result := s.Clone()
	for _, i := range is {
		if i == nil {
			continue
		}
		if result.languages != nil {
			result.languages.Complement(i.languages)
		}
		if result.versions != nil {
			result.versions.Complement(i.versions)
		}
		if result.roles != nil {
			result.roles.Complement(i.roles)
		}
	}
	return result
}

func (s IntSets) WithDefaultsIfNotSet(cfg ConfiguredDimensions) *IntSets {
	s.SetDefaultsIfNotSet(cfg)
	return &s
}

// WithLanguageIndex replaces the current language set with a single language index.
func (s IntSets) WithLanguageIndex(i int) *IntSets {
	s.languages = maps.NewOrderedIntSet(i)
	return &s
}

type IntSetsConfig struct {
	Cfg           ConfiguredDimensions
	Ordinal       int
	ApplyDefaults bool
	Globs         StringSlices
}

// NewIntSets creates a new DimensionsIntSets with nil sets for languages, roles, and versions.
func NewIntSets(ordinal int) *IntSets {
	return &IntSets{ordinal: ordinal}
}

// NewIntSetsFromConfig creates a new IntSets from the given IntSetsConfig.
// It applies the filters based on the provided languages, versions, and roles.
func NewIntSetsFromConfig(cfg IntSetsConfig) (*IntSets, error) {
	applyFilter := func(what string, values []string, matcher ConfiguredDimension) (*maps.OrderedIntSet, error) {
		if len(values) == 0 {
			if cfg.ApplyDefaults {
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

	sets := NewIntSets(cfg.Ordinal)
	l, err1 := applyFilter("languages", cfg.Globs.Languages, cfg.Cfg.ConfiguredLanguages)
	v, err2 := applyFilter("versions", cfg.Globs.Versions, cfg.Cfg.ConfiguredVersions)
	r, err3 := applyFilter("roles", cfg.Globs.Roles, cfg.Cfg.ConfiguredRoles)

	if err := cmp.Or(err1, err2, err3); err != nil {
		return nil, fmt.Errorf("failed to apply filters: %w", err)
	}
	sets.languages = l
	sets.versions = v
	sets.roles = r

	return sets, nil
}

// Sites holds configuration about which sites a file/content/page/resource belongs to.
type Sites struct {
	// Matrix defines the main build matrix.
	Matrix StringSlices `mapstructure:"matrix" json:"matrix"`
	// Fallbacks defines the fallback matrix.
	Fallbacks StringSlices `mapstructure:"fallbacks" json:"fallbacks"`
}

// IsZero returns true if all slices are empty.
func (s Sites) IsZero() bool {
	return s.Matrix.IsZero() && s.Fallbacks.IsZero()
}

func (s *Sites) SetFromParamsIfNotSet(params maps.Params) {
	const (
		matrixKey    = "matrix"
		fallbacksKey = "fallbacks"
	)

	if m, ok := params[matrixKey]; ok {
		s.Matrix.SetFromParamsIfNotSet(m.(maps.Params))
	}
	if f, ok := params[fallbacksKey]; ok {
		s.Fallbacks.SetFromParamsIfNotSet(f.(maps.Params))
	}
}

// StringSlices holds slices of Glob patterns for languages, versions, and roles.
type StringSlices struct {
	Languages []string `mapstructure:"languages" json:"languages"`
	Versions  []string `mapstructure:"versions" json:"versions"`
	Roles     []string `mapstructure:"roles" json:"roles"`
}

func (d StringSlices) IsZero() bool {
	return len(d.Languages) == 0 && len(d.Versions) == 0 && len(d.Roles) == 0
}

func (d *StringSlices) SetFromParamsIfNotSet(params maps.Params) {
	const (
		languagesKey = "languages"
		versionsKey  = "versions"
		rolesKey     = "roles"
	)

	if len(d.Languages) == 0 {
		if v, ok := params[languagesKey]; ok {
			d.Languages = types.ToStringSlicePreserveString(v)
		}
	}

	if len(d.Versions) == 0 {
		if v, ok := params[versionsKey]; ok {
			d.Versions = types.ToStringSlicePreserveString(v)
		}
	}

	if len(d.Roles) == 0 {
		if v, ok := params[rolesKey]; ok {
			d.Roles = types.ToStringSlicePreserveString(v)
		}
	}
}

type ConfiguredDimension interface {
	predicate.IndexMatcher
	IndexDefault() int
	ResolveName(int) string
	ResolveIndex(string) int
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

func (c ConfiguredDimensions) ResolveVector(names types.Strings3) Vector {
	var vec Vector
	if s := names[0]; s != "" {
		vec[0] = c.ConfiguredLanguages.ResolveIndex(s)
	} else {
		vec[0] = c.ConfiguredLanguages.IndexDefault()
	}
	if s := names[1]; s != "" {
		vec[1] = c.ConfiguredVersions.ResolveIndex(s)
	} else {
		vec[1] = c.ConfiguredVersions.IndexDefault()
	}
	if s := names[2]; s != "" {
		vec[2] = c.ConfiguredRoles.ResolveIndex(s)
	} else {
		vec[2] = c.ConfiguredRoles.IndexDefault()
	}
	return vec
}
