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

package sitesmatrix

import (
	"cmp"
	"fmt"
	"iter"
	xmaps "maps"
	"slices"
	"sort"
	"sync"

	"github.com/gohugoio/hashstructure"
	"github.com/gohugoio/hugo/common/hashing"
	"github.com/gohugoio/hugo/common/hdebug"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/common/predicate"
	"github.com/gohugoio/hugo/common/types"
	"github.com/gohugoio/hugo/hugofs/glob"
)

var (
	_ VectorStore            = &IntSets{}
	_ VectorStore            = &vectorStoreMap{}
	_ hashstructure.Hashable = &IntSets{}
	_ hashstructure.Hashable = &vectorStoreMap{}
)

func newVectorStoreMap(cap int) *vectorStoreMap {
	return &vectorStoreMap{
		sets: make(map[Vector]struct{}, cap),
		h:    &hashOnce{},
	}
}

// A vector store backed by a map.
type vectorStoreMap struct {
	sets map[Vector]struct{}
	h    *hashOnce
}

func (m *vectorStoreMap) initHash() {
	m.h.once.Do(func() {
		var err error
		m.h.hash, err = hashing.Hash(m.sets)
		if err != nil {
			panic(fmt.Errorf("failed to calculate hash for MapVectorStore: %w", err))
		}
	})
}

func (m *vectorStoreMap) setVector(vec Vector) {
	m.sets[vec] = struct{}{}
}

func (m *vectorStoreMap) Ordinal() int {
	return 0
}

func (m *vectorStoreMap) KeysSorted() ([]int, []int, []int) {
	var k0, k1, k2 []int
	for v := range m.sets {
		k0 = append(k0, v.Language())
		k1 = append(k1, v.Version())
		k2 = append(k2, v.Role())
	}
	sort.Ints(k0)
	sort.Ints(k1)
	sort.Ints(k2)
	k0 = slices.Compact(k0)
	k1 = slices.Compact(k1)
	k2 = slices.Compact(k2)

	return k0, k1, k2
}

func (m *vectorStoreMap) Hash() (uint64, error) {
	m.initHash()
	return m.h.hash, nil
}

func (m *vectorStoreMap) MustHash() uint64 {
	i, _ := m.Hash()
	return i
}

func (m *vectorStoreMap) HasVector(v Vector) bool {
	if _, ok := m.sets[v]; ok {
		return true
	}
	return false
}

func (m *vectorStoreMap) HasAnyVector(v VectorProvider) bool {
	if v == nil || m.LenVectors() == 0 || v.LenVectors() == 0 {
		return false
	}

	return !v.ForEeachVector(func(vec Vector) bool {
		if m.HasVector(vec) {
			return false // stop iteration
		}
		return true // continue iteration
	})
}

func (m *vectorStoreMap) LenVectors() int {
	return len(m.sets)
}

func (m *vectorStoreMap) Complement(is ...VectorProvider) VectorStore {
	panic("TODO1: Implement me.")
}

func (m *vectorStoreMap) EqualsVector(other VectorProvider) bool {
	if other == nil {
		return false
	}
	if m == other {
		return true
	}
	if m.LenVectors() != other.LenVectors() {
		return false
	}
	return other.ForEeachVector(func(v Vector) bool {
		_, ok := m.sets[v]
		return ok
	})
}

func (m *vectorStoreMap) FirstVector() Vector {
	if len(m.sets) == 0 {
		panic("no vectors available")
	}
	for v := range m.sets {
		return v
	}
	panic("unreachable")
}

func (m *vectorStoreMap) ForEeachVector(yield func(v Vector) bool) bool {
	if len(m.sets) == 0 {
		return true
	}
	for v := range m.sets {
		if !yield(v) {
			return false
		}
	}
	return true
}

func (m *vectorStoreMap) Vectors() []Vector {
	if m == nil || len(m.sets) == 0 {
		return nil
	}
	var vectors []Vector
	for v := range m.sets {
		vectors = append(vectors, v)
	}
	sort.Slice(vectors, func(i, j int) bool {
		v1, v2 := vectors[i], vectors[j]
		return v1.Compare(v2) < 0
	})
	return vectors
}

func (m *vectorStoreMap) WithLanguageIndex(i int) VectorStore {
	c := m.clone()

	for v := range c.sets {
		v[Language.Index()] = i
		c.sets[v] = struct{}{}
	}

	return c
}

func (m *vectorStoreMap) HasLanguage(i int) bool {
	for v := range m.sets {
		if v.Language() == i {
			return true
		}
	}
	return false
}

func (m *vectorStoreMap) HasVersion(i int) bool {
	for v := range m.sets {
		if v.Version() == i {
			return true
		}
	}
	return false
}

func (m *vectorStoreMap) HasRole(i int) bool {
	for v := range m.sets {
		if v.Role() == i {
			return true
		}
	}
	return false
}

func (m *vectorStoreMap) clone() *vectorStoreMap {
	c := *m
	c.h = &hashOnce{}
	c.sets = xmaps.Clone(m.sets)
	return &c
}

// NewIntSets creates a new NewIntSets with nil sets for languages, roles, and versions.
// TODO1 remove me.
func NewIntSets() *IntSets {
	return &IntSets{h: &hashOnce{}}
}

func NewIntSetsBuilder() *IntSetsBuilder {
	return &IntSetsBuilder{s: &IntSets{h: &hashOnce{}}}
}

type ConfiguredDimension interface {
	predicate.IndexMatcher
	IndexDefault() int
	ResolveIndex(string) int
	ResolveName(int) string
	ForEachIndex() iter.Seq[int]
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
	languages *maps.OrderedIntSet `mapstructure:"-" json:"-"` // TODO1 does this need to be ordered?
	versions  *maps.OrderedIntSet `mapstructure:"-" json:"-"`
	roles     *maps.OrderedIntSet `mapstructure:"-" json:"-"`

	h *hashOnce
}

type hashOnce struct {
	once sync.Once
	hash uint64
}

func (s *IntSets) ToVectorStoreMap() *vectorStoreMap {
	if s == nil {
		return nil
	}

	result := newVectorStoreMap(36)

	s.ForEeachVector(func(vec Vector) bool {
		result.setVector(vec)
		return true
	})

	return result
}

func (s *IntSets) IsSuperSet(other *IntSets) bool {
	return s.languages.IsSuperSet(other.languages) &&
		s.versions.IsSuperSet(other.versions) &&
		s.roles.IsSuperSet(other.roles)
}

func (s *IntSets) DifferenceCardinality(other *IntSets) Vector {
	return Vector{
		int(s.languages.Values().DifferenceCardinality(other.languages.Values())),
		int(s.versions.Values().DifferenceCardinality(other.versions.Values())),
		int(s.roles.Values().DifferenceCardinality(other.roles.Values())),
	}
}

// Complement returns a new VectorStore that contains all vectors in s that are not in any of ss.
func (s *IntSets) Complement(ss ...VectorProvider) VectorStore {
	if len(ss) == 0 || (len(ss) == 1 && ss[0] == s) {
		return nil
	}

	for _, v := range ss {
		vv, ok := v.(*IntSets)
		if !ok {
			continue
		}
		if vv.IsSuperSet(s) {
			var s *IntSets
			return s
		}
	}

	result := newVectorStoreMap(36)

	s.ForEeachVector(func(vec Vector) bool {
		var found bool
		for _, v := range ss {
			if v.HasVector(vec) {
				found = true
				break
			}
		}

		if !found {
			result.setVector(vec)
		}

		return true
	})

	return result
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
		s.languages.Next(0),
		s.versions.Next(0),
		s.roles.Next(0),
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
				return yield(Vector{lang, ver, role})
			})
		})
	})

	return b
}

func (s *IntSets) Vectors() []Vector {
	if s.LenVectors() == 0 {
		return nil
	}

	var vectors []Vector
	s.ForEeachVector(func(v Vector) bool {
		vectors = append(vectors, v)
		return true
	})

	sort.Slice(vectors, func(i, j int) bool {
		v1, v2 := vectors[i], vectors[j]
		return v1.Compare(v2) < 0
	})

	return vectors
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

func (s *IntSets) HasAnyVector(v VectorProvider) bool {
	if s == nil || v == nil {
		return false
	}
	if s.LenVectors() == 0 || v.LenVectors() == 0 {
		return false
	}

	return !v.ForEeachVector(func(vec Vector) bool {
		if s.HasVector(vec) {
			return false // stop iteration
		}
		return true // continue iteration
	})
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
func (s *IntSets) WithLanguageIndex(i int) VectorStore {
	c := s.shallowClone()
	c.languages = maps.NewOrderedIntSet(i)
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

func (s *IntSets) setDefaultsAndAllLAnguagesIfNotSet(cfg ConfiguredDimensions) {
	if s.languages == nil {
		s.languages = maps.NewOrderedIntSet()
		for i := range cfg.ConfiguredLanguages.ForEachIndex() {
			s.languages.Set(i)
		}
	}
	s.setDefaultsIfNotSet(cfg)
}

func (s *IntSets) setAllIfNotSet(cfg ConfiguredDimensions) {
	if s.languages == nil {
		s.languages = maps.NewOrderedIntSet()
		for i := range cfg.ConfiguredLanguages.ForEachIndex() {
			s.languages.Set(i)
		}
	}
	if s.versions == nil {
		s.versions = maps.NewOrderedIntSet()
		for i := range cfg.ConfiguredVersions.ForEachIndex() {
			s.versions.Set(i)
		}
	}
	if s.roles == nil {
		s.roles = maps.NewOrderedIntSet()
		for i := range cfg.ConfiguredRoles.ForEachIndex() {
			s.roles.Set(i)
		}
	}
}

func (s *IntSets) initHash() {
	s.h.once.Do(func() {
		var err error
		s.h.hash, err = hashing.Hash(s.languages.Words(), s.versions.Words(), s.roles.Words())
		if err != nil {
			panic(fmt.Errorf("failed to calculate hash for IntSets: %w", err))
		}
	})
}

func (s *IntSets) init() *IntSets {
	return s
}

func (s *IntSets) setDimensionsFromOtherIfNotSet(other VectorStore) {
	if other == nil {
		return
	}
	setLang := s.languages == nil
	setVer := s.versions == nil
	setRole := s.roles == nil

	if !(setLang || setVer || setRole) {
		return
	}

	other.ForEeachVector(func(v Vector) bool {
		if !s.HasVector(v) {
			s.setValuesInNilSets(v, setLang, setVer, setRole)
		}
		return true
	})
}

func (s *IntSets) setValuesInNilSets(vec Vector, setLang, setVer, setRole bool) {
	if setLang {
		if s.languages == nil {
			s.languages = maps.NewOrderedIntSet()
		}
		s.languages.Set(vec.Language())
	}
	if setVer {
		if s.versions == nil {
			s.versions = maps.NewOrderedIntSet()
		}
		s.versions.Set(vec.Version())
	}

	if setRole {
		if s.roles == nil {
			s.roles = maps.NewOrderedIntSet()
		}
		s.roles.Set(vec.Role())
	}
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
		var result *maps.OrderedIntSet
		if len(values) == 0 {

			if cfg.ApplyDefaults > 0 {
				result = maps.NewOrderedIntSet()
			}
			switch cfg.ApplyDefaults {
			case IntSetsConfigApplyDefaultsIfNotSet:
				result.Set(matcher.IndexDefault())
			case IntSetsConfigApplyDefaultsAndAllLanguagesIfNotSet:
				if what == "languages" {
					for i := range matcher.ForEachIndex() {
						result.Set(i)
					}
				} else {
					result.Set(matcher.IndexDefault())
				}
			}

			return result, nil
		}

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

	// TODO1 defaults.if l == nil &&

	b.s.languages = l
	b.s.versions = v
	b.s.roles = r

	return b
}

func (s *IntSetsBuilder) WithLanguageIndex(i int) *IntSetsBuilder {
	if s.s.languages == nil {
		s.s.languages = maps.NewOrderedIntSet()
	}
	s.s.languages.Set(i)
	return s
}

func (s *IntSetsBuilder) WithDefaultsAndAllLanguagesIfNotSet(cfg ConfiguredDimensions) *IntSetsBuilder {
	s.s.setDefaultsAndAllLAnguagesIfNotSet(cfg)
	return s
}

func (s *IntSetsBuilder) WithAllIfNotSet(cfg ConfiguredDimensions) *IntSetsBuilder {
	s.s.setAllIfNotSet(cfg)
	return s
}

func (s *IntSetsBuilder) WithDefaultsIfNotSet(cfg ConfiguredDimensions) *IntSetsBuilder {
	s.s.setDefaultsIfNotSet(cfg)
	return s
}

func (s *IntSetsBuilder) WithDimensionsFromOtherIfNotSet(other VectorStore) *IntSetsBuilder {
	s.s.setDimensionsFromOtherIfNotSet(other)
	return s
}

// TODO1 remove me.
func (s *IntSetsBuilder) Debug() {
	hdebug.Printf("IntSetsBuilder: languages: %v, versions: %v, roles: %v", s.s.languages, s.s.versions, s.s.roles)
}

func (b *IntSetsBuilder) WithSets(languages, versions, roles *maps.OrderedIntSet) *IntSetsBuilder {
	b.s.languages = languages
	b.s.versions = versions
	b.s.roles = roles
	return b
}

type IntSetsConfigApplyDefaults int

const (
	IntSetsConfigApplyDefaultsNone IntSetsConfigApplyDefaults = iota
	IntSetsConfigApplyDefaultsIfNotSet
	IntSetsConfigApplyDefaultsAndAllLanguagesIfNotSet
)

type IntSetsConfig struct {
	Cfg           ConfiguredDimensions
	ApplyDefaults IntSetsConfigApplyDefaults
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
