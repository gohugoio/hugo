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
	"sync"

	"github.com/gohugoio/hashstructure"
	"github.com/gohugoio/hugo/common/hashing"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/common/predicate"
	"github.com/gohugoio/hugo/common/types"
	"github.com/gohugoio/hugo/hugofs/glob"
)

var (
	_ VectorProvider         = &IntSets{}
	_ hashstructure.Hashable = &IntSets{}
)

// NewIntSets creates a new DimensionsIntSets with nil sets for languages, roles, and versions.
// TODO1 remove me.
func NewIntSets(ordinal int) *IntSets {
	return &IntSets{ordinal: ordinal, h: &hashOnce{}}
}

func NewIntSetsBuilder(ordinal int) *IntSetsBuilder {
	return &IntSetsBuilder{s: &IntSets{ordinal: ordinal, h: &hashOnce{}}}
}

type ConfiguredDimension interface {
	predicate.IndexMatcher
	IndexDefault() int
	ResolveIndex(string) int
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

// IntSets holds the ordered sets of integers for the dimensions,
// which is used for fast membership testing of files, resources and pages.
type IntSets struct {
	ordinal   int
	languages *maps.OrderedIntSet `mapstructure:"-" json:"-"` // TODO1 does this need to be ordered?
	versions  *maps.OrderedIntSet `mapstructure:"-" json:"-"`
	roles     *maps.OrderedIntSet `mapstructure:"-" json:"-"`

	h *hashOnce
}

type hashOnce struct {
	once sync.Once
	hash uint64
}

// Complement returns a new IntSets that is the complement of the IntSets passed in is.
// This will return nil if the resulting set is empty.
func (s *IntSets) Complement(is ...*IntSets) *IntSets {
	if len(is) == 0 || (len(is) == 1 && is[0] == s) {
		return nil
	}

	// TODO1 see bitsets.IsSuperSet etc.

	result := NewIntSets(s.ordinal)

	s.ForEeachVector(func(vec Vector) bool {
		var found bool
		for _, v := range is {
			if v.HasVector(vec) {
				found = true
				break
			}
		}

		if !found {
			result.initSets()
			result.setVector(vec)
		}

		return true
	})

	return result.init()
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

func (s *IntSets) LenVectors() int {
	if s == nil {
		return 0
	}
	return s.languages.Len() * s.versions.Len() * s.roles.Len()
}

func (s *IntSets) Ordinal() int {
	if s == nil {
		return 0
	}
	return s.ordinal
}

func (s *IntSets) HasRole(role int) bool {
	if s == nil {
		return false
	}
	return s.roles.Has(role)
}

func (s *IntSets) String() string {
	return fmt.Sprintf("Languages: %v, Versions: %v, Roles: %v", s.languages, s.versions, s.roles)
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

func (s *IntSets) HasVersion(ver int) bool {
	if s == nil {
		return false
	}
	return s.versions.Has(ver)
}

func (s IntSets) shallowClone() *IntSets {
	s.h = &hashOnce{}
	return &s
}

// WithLanguageIndex replaces the current language set with a single language index.
func (s *IntSets) WithLanguageIndex(i int) *IntSets {
	c := s.shallowClone()
	c.languages = maps.NewOrderedIntSet(i)
	return c.init()
}

func (s *IntSets) WithOrdinal(i int) *IntSets {
	c := s.shallowClone()
	c.ordinal = i
	return c.init()
}

func (s *IntSets) Hash() (uint64, error) {
	s.initHash()
	return s.h.hash, nil
}

func (s *IntSets) MustHash() uint64 {
	if s == nil {
		return 0
	}
	hash, err := s.Hash()
	if err != nil {
		panic(fmt.Errorf("failed to calculate hash for IntSets: %w", err))
	}
	return hash
}

// setDefaultsIfNotSet applies default values to the IntSets if they are not already set.
func (s *IntSets) setDefaultsIfNotSet(cfg ConfiguredDimensions) {
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

func (s *IntSets) initHash() {
	s.h.once.Do(func() {
		var err error
		s.h.hash, err = hashing.Hash(s.ordinal, s.languages.Words(), s.versions.Words(), s.roles.Words())
		if err != nil {
			panic(fmt.Errorf("failed to calculate hash for IntSets: %w", err))
		}
	})
}

func (s *IntSets) init() *IntSets {
	return s
}

func (s *IntSets) setFromOtherIfNotSet(other *IntSets) {
	if other == nil {
		return
	}

	// TODO1 clean up these nil checks.
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

func (s *IntSets) initSets() {
	if s.languages == nil {
		s.languages = maps.NewOrderedIntSet()
	}
	if s.versions == nil {
		s.versions = maps.NewOrderedIntSet()
	}
	if s.roles == nil {
		s.roles = maps.NewOrderedIntSet()
	}
}

func (s *IntSets) setVector(vec Vector) {
	s.languages.Set(vec.Language())
	s.versions.Set(vec.Version())
	s.roles.Set(vec.Role())
}

type IntSetsBuilder struct {
	s *IntSets
}

func (b *IntSetsBuilder) Build() *IntSets {
	b.s.init()
	return b.s
}

func (b *IntSetsBuilder) WithConfig(cfg IntSetsConfig) *IntSetsBuilder {
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

	l, err1 := applyFilter("languages", cfg.Globs.Languages, cfg.Cfg.ConfiguredLanguages)
	v, err2 := applyFilter("versions", cfg.Globs.Versions, cfg.Cfg.ConfiguredVersions)
	r, err3 := applyFilter("roles", cfg.Globs.Roles, cfg.Cfg.ConfiguredRoles)

	if err := cmp.Or(err1, err2, err3); err != nil {
		panic(fmt.Errorf("failed to apply filters: %w", err))
	}
	b.s.languages = l
	b.s.versions = v
	b.s.roles = r

	return b
}

func (s *IntSetsBuilder) WithDefaultsIfNotSet(cfg ConfiguredDimensions) *IntSetsBuilder {
	s.s.setDefaultsIfNotSet(cfg)
	return s
}

func (s *IntSetsBuilder) WithFromOtherIfNotSet(other *IntSets) *IntSetsBuilder {
	s.s.setFromOtherIfNotSet(other)
	return s
}

func (b *IntSetsBuilder) WithSets(languages, versions, roles *maps.OrderedIntSet) *IntSetsBuilder {
	b.s.languages = languages
	b.s.versions = versions
	b.s.roles = roles
	return b
}

type IntSetsConfig struct {
	Cfg           ConfiguredDimensions
	ApplyDefaults bool
	Globs         StringSlices
}

// Sites holds configuration about which sites a file/content/page/resource belongs to.
type Sites struct {
	// Matrix defines the main build matrix.
	Matrix StringSlices `mapstructure:"matrix" json:"matrix"`
	// Fallbacks defines the fallback matrix.
	Fallbacks StringSlices `mapstructure:"fallbacks" json:"fallbacks"`
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

// IsZero returns true if all slices are empty.
func (s Sites) IsZero() bool {
	return s.Matrix.IsZero() && s.Fallbacks.IsZero()
}

// StringSlices holds slices of Glob patterns for languages, versions, and roles.
type StringSlices struct {
	Languages []string `mapstructure:"languages" json:"languages"`
	Versions  []string `mapstructure:"versions" json:"versions"`
	Roles     []string `mapstructure:"roles" json:"roles"`
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

func (d StringSlices) IsZero() bool {
	return len(d.Languages) == 0 && len(d.Versions) == 0 && len(d.Roles) == 0
}
