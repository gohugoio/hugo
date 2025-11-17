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
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/common/predicate"
	"github.com/gohugoio/hugo/common/types"
	"github.com/gohugoio/hugo/hugofs/hglob"
)

var (
	_ VectorStore            = &IntSets{}
	_ VectorStore            = &vectorStoreMap{}
	_ hashstructure.Hashable = &IntSets{}
	_ hashstructure.Hashable = &vectorStoreMap{}
)

func newVectorStoreMap(cap int) *vectorStoreMap {
	return &vectorStoreMap{
		sets: make(Vectors, cap),
		h:    &hashOnce{},
	}
}

func newVectorStoreMapFromVectors(v Vectors) *vectorStoreMap {
	return &vectorStoreMap{
		sets: v,
		h:    &hashOnce{},
	}
}

// A vector store backed by a map.
type vectorStoreMap struct {
	sets Vectors
	h    *hashOnce
}

func (s *vectorStoreMap) initHash() {
	s.h.once.Do(func() {
		var err error
		s.h.hash, err = hashing.Hash(s.sets)
		if err != nil {
			panic(fmt.Errorf("failed to calculate hash for MapVectorStore: %w", err))
		}
	})
}

func (s *vectorStoreMap) setVector(vec Vector) {
	s.sets[vec] = struct{}{}
}

func (s *vectorStoreMap) KeysSorted() ([]int, []int, []int) {
	var k0, k1, k2 []int
	for v := range s.sets {
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

func (s *vectorStoreMap) Hash() (uint64, error) {
	s.initHash()
	return s.h.hash, nil
}

func (s *vectorStoreMap) MustHash() uint64 {
	i, _ := s.Hash()
	return i
}

func (s *vectorStoreMap) HasVector(v Vector) bool {
	if _, ok := s.sets[v]; ok {
		return true
	}
	return false
}

func (s *vectorStoreMap) HasAnyVector(v VectorProvider) bool {
	if v == nil || s.LenVectors() == 0 || v.LenVectors() == 0 {
		return false
	}

	return !v.ForEachVector(func(vec Vector) bool {
		if s.HasVector(vec) {
			return false // stop iteration
		}
		return true // continue iteration
	})
}

func (s *vectorStoreMap) LenVectors() int {
	return len(s.sets)
}

// Complement returns a new VectorStore that contains all vectors in s that are not in any of ss.
func (s *vectorStoreMap) Complement(ss ...VectorProvider) VectorStore {
	if len(ss) == 0 || (len(ss) == 1 && ss[0] == s) {
		return nil
	}

	for _, v := range ss {
		if v == s {
			var s *vectorStoreMap
			return s
		}
	}

	result := newVectorStoreMap(36)

	s.ForEachVector(func(vec Vector) bool {
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

func (s *vectorStoreMap) EqualsVector(other VectorProvider) bool {
	if other == nil {
		return false
	}
	if s == other {
		return true
	}
	if s.LenVectors() != other.LenVectors() {
		return false
	}
	return other.ForEachVector(func(v Vector) bool {
		_, ok := s.sets[v]
		return ok
	})
}

func (s *vectorStoreMap) VectorSample() Vector {
	if len(s.sets) == 0 {
		panic("no vectors available")
	}
	for v := range s.sets {
		return v
	}
	panic("unreachable")
}

func (s *vectorStoreMap) ForEachVector(yield func(v Vector) bool) bool {
	if len(s.sets) == 0 {
		return true
	}
	for v := range s.sets {
		if !yield(v) {
			return false
		}
	}
	return true
}

func (s *vectorStoreMap) Vectors() []Vector {
	if s == nil || len(s.sets) == 0 {
		return nil
	}
	var vectors []Vector
	for v := range s.sets {
		vectors = append(vectors, v)
	}
	sort.Slice(vectors, func(i, j int) bool {
		v1, v2 := vectors[i], vectors[j]
		return v1.Compare(v2) < 0
	})
	return vectors
}

func (s *vectorStoreMap) WithLanguageIndices(i int) VectorStore {
	c := s.clone()

	for v := range c.sets {
		v[Language] = i
		c.sets[v] = struct{}{}
	}

	return c
}

func (s *vectorStoreMap) HasLanguage(i int) bool {
	for v := range s.sets {
		if v.Language() == i {
			return true
		}
	}
	return false
}

func (s *vectorStoreMap) HasVersion(i int) bool {
	for v := range s.sets {
		if v.Version() == i {
			return true
		}
	}
	return false
}

func (s *vectorStoreMap) HasRole(i int) bool {
	for v := range s.sets {
		if v.Role() == i {
			return true
		}
	}
	return false
}

func (s *vectorStoreMap) clone() *vectorStoreMap {
	c := *s
	c.h = &hashOnce{}
	c.sets = xmaps.Clone(s.sets)
	return &c
}

func NewIntSetsBuilder(cfg *ConfiguredDimensions) *IntSetsBuilder {
	if cfg == nil {
		panic("cfg is required")
	}
	return &IntSetsBuilder{cfg: cfg, s: &IntSets{h: &hashOnce{}}}
}

type ConfiguredDimension interface {
	predicate.IndexMatcher
	IndexDefault() int
	ResolveIndex(string) int
	ResolveName(int) string
	ForEachIndex() iter.Seq[int]
	Len() int
}

// ConfiguredDimensions holds the configured dimensions for the site matrix.
type ConfiguredDimensions struct {
	ConfiguredLanguages ConfiguredDimension
	ConfiguredVersions  ConfiguredDimension
	ConfiguredRoles     ConfiguredDimension
	CommonSitesMatrix   CommonSitestMatrix

	singleVectorStoreCache *maps.Cache[Vector, *IntSets]
}

func (c *ConfiguredDimensions) IsSingleVector() bool {
	return c.ConfiguredLanguages.Len() == 1 && c.ConfiguredRoles.Len() == 1 && c.ConfiguredVersions.Len() == 1
}

// GetOrCreateSingleVectorStore returns a VectorStore for the given vector.
func (c *ConfiguredDimensions) GetOrCreateSingleVectorStore(vec Vector) *IntSets {
	store, _ := c.singleVectorStoreCache.GetOrCreate(vec, func() (*IntSets, error) {
		is := &IntSets{}
		is.setValuesInNilSets(vec, true, true, true)
		return is, nil
	})
	return store
}

func (c *ConfiguredDimensions) Init() error {
	c.singleVectorStoreCache = maps.NewCache[Vector, *IntSets]()
	b := NewIntSetsBuilder(c).WithDefaultsIfNotSet().Build()
	defaultVec := b.VectorSample()
	c.singleVectorStoreCache.Set(defaultVec, b)
	c.CommonSitesMatrix.DefaultSite = b

	return nil
}

type CommonSitestMatrix struct {
	DefaultSite VectorStore
}

func (c *ConfiguredDimensions) ResolveNames(v Vector) types.Strings3 {
	return types.Strings3{
		c.ConfiguredLanguages.ResolveName(v.Language()),
		c.ConfiguredVersions.ResolveName(v.Version()),
		c.ConfiguredRoles.ResolveName(v.Role()),
	}
}

func (c *ConfiguredDimensions) ResolveVector(names types.Strings3) Vector {
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
	languages *maps.OrderedIntSet `mapstructure:"-" json:"-"`
	versions  *maps.OrderedIntSet `mapstructure:"-" json:"-"`
	roles     *maps.OrderedIntSet `mapstructure:"-" json:"-"`

	h *hashOnce
}

var NilStore *IntSets = nil

type hashOnce struct {
	once sync.Once
	hash uint64
}

func (s *IntSets) ToVectorStoreMap() *vectorStoreMap {
	if s == nil {
		return nil
	}

	result := newVectorStoreMap(36)

	s.ForEachVector(func(vec Vector) bool {
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

func (s *IntSets) Intersects(other *IntSets) bool {
	if s == nil || other == nil {
		return false
	}
	return s.languages.Values().IntersectionCardinality(other.languages.Values()) > 0 &&
		s.versions.Values().IntersectionCardinality(other.versions.Values()) > 0 &&
		s.roles.Values().IntersectionCardinality(other.roles.Values()) > 0
}

// Complement returns a new VectorStore that contains all vectors in s that are not in any of ss.
func (s *IntSets) Complement(ss ...VectorProvider) VectorStore {
	if len(ss) == 0 || (len(ss) == 1 && ss[0] == s) {
		return NilStore
	}

	for _, v := range ss {
		vv, ok := v.(*IntSets)
		if !ok {
			continue
		}
		if vv.IsSuperSet(s) {
			return NilStore
		}
	}

	result := newVectorStoreMap(36)

	s.ForEachVector(func(vec Vector) bool {
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

	return other.ForEachVector(func(v Vector) bool {
		return s.HasVector(v)
	})
}

func (s *IntSets) VectorSample() Vector {
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
func (s *IntSets) ForEachVector(yield func(v Vector) bool) bool {
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
	s.ForEachVector(func(v Vector) bool {
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

// LenVectors returns the total number of vectors represented by the IntSets.
// This is the Cartesian product of the lengths of the individual sets.
// This will be 0 if s is nil or any of the sets is empty.
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
	if v.LenVectors() == 1 {
		// Fast path.
		return s.HasVector(v.VectorSample())
	}

	if vs, ok := v.(*IntSets); ok {
		// Fast path.
		return s.Intersects(vs)
	}

	return !v.ForEachVector(func(vec Vector) bool {
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

// WithLanguageIndices replaces the current language set with a single language index.
func (s *IntSets) WithLanguageIndices(i int) VectorStore {
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
func (s *IntSets) setDefaultsIfNotSet(cfg *ConfiguredDimensions) {
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

func (s *IntSets) setDefaultsAndAllLAnguagesIfNotSet(cfg *ConfiguredDimensions) {
	if s.languages == nil {
		s.languages = maps.NewOrderedIntSet()
		for i := range cfg.ConfiguredLanguages.ForEachIndex() {
			s.languages.Set(i)
		}
	}
	s.setDefaultsIfNotSet(cfg)
}

func (s *IntSets) setAllIfNotSet(cfg *ConfiguredDimensions) {
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

func (s *IntSets) setDimensionsFromOtherIfNotSet(other VectorIterator) {
	if other == nil {
		return
	}
	setLang := s.languages == nil
	setVer := s.versions == nil
	setRole := s.roles == nil

	if !(setLang || setVer || setRole) {
		return
	}

	other.ForEachVector(func(v Vector) bool {
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
	cfg *ConfiguredDimensions
	s   *IntSets

	// Set when a Glob (e.g. "en") filter is provided but no matches are found.
	GlobFilterMisses Bools
}

func (b *IntSetsBuilder) Build() *IntSets {
	b.s.init()

	if b.s.LenVectors() == 1 {
		// Cache it or use the existing cached version, which will allow b.s to be GCed.
		bb, _ := b.cfg.singleVectorStoreCache.GetOrCreate(b.s.VectorSample(), func() (*IntSets, error) {
			return b.s, nil
		})
		return bb
	}
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
		filter, err := predicate.NewStringPredicateFromGlobs(values, hglob.GetGlobDot)
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

	l, err1 := applyFilter("languages", cfg.Globs.Languages, b.cfg.ConfiguredLanguages)
	v, err2 := applyFilter("versions", cfg.Globs.Versions, b.cfg.ConfiguredVersions)
	r, err3 := applyFilter("roles", cfg.Globs.Roles, b.cfg.ConfiguredRoles)

	if err := cmp.Or(err1, err2, err3); err != nil {
		panic(fmt.Errorf("failed to apply filters: %w", err))
	}

	b.GlobFilterMisses = Bools{
		len(cfg.Globs.Languages) > 0 && l == nil,
		len(cfg.Globs.Versions) > 0 && v == nil,
		len(cfg.Globs.Roles) > 0 && r == nil,
	}

	b.s.languages = l
	b.s.versions = v
	b.s.roles = r

	return b
}

func (s *IntSetsBuilder) WithLanguageIndices(idxs ...int) *IntSetsBuilder {
	if len(idxs) == 0 {
		return s
	}
	if s.s.languages == nil {
		s.s.languages = maps.NewOrderedIntSet()
	}
	for _, i := range idxs {
		s.s.languages.Set(i)
	}
	return s
}

func (s *IntSetsBuilder) WithDefaultsAndAllLanguagesIfNotSet() *IntSetsBuilder {
	s.s.setDefaultsAndAllLAnguagesIfNotSet(s.cfg)
	return s
}

func (s *IntSetsBuilder) WithAllIfNotSet() *IntSetsBuilder {
	s.s.setAllIfNotSet(s.cfg)
	return s
}

func (s *IntSetsBuilder) WithDefaultsIfNotSet() *IntSetsBuilder {
	s.s.setDefaultsIfNotSet(s.cfg)
	return s
}

func (s *IntSetsBuilder) WithDimensionsFromOtherIfNotSet(other VectorIterator) *IntSetsBuilder {
	s.s.setDimensionsFromOtherIfNotSet(other)
	return s
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
	ApplyDefaults IntSetsConfigApplyDefaults
	Globs         StringSlices
}

// Sites holds configuration about which sites a file/content/page/resource belongs to.
type Sites struct {
	// Matrix defines what sites to build this content for.
	Matrix StringSlices `mapstructure:"matrix" json:"matrix"`
	// Complements defines what sites to complement with this content.
	Complements StringSlices `mapstructure:"complements" json:"complements"`
}

func (s *Sites) Equal(other Sites) bool {
	return s.Matrix.Equal(other.Matrix) && s.Complements.Equal(other.Complements)
}

func (s *Sites) SetFromParamsIfNotSet(params maps.Params) {
	const (
		matrixKey      = "matrix"
		complementsKey = "complements"
	)

	if m, ok := params[matrixKey]; ok {
		s.Matrix.SetFromParamsIfNotSet(m.(maps.Params))
	}
	if f, ok := params[complementsKey]; ok {
		s.Complements.SetFromParamsIfNotSet(f.(maps.Params))
	}
}

// IsZero returns true if all slices are empty.
func (s Sites) IsZero() bool {
	return s.Matrix.IsZero() && s.Complements.IsZero()
}

// StringSlices holds slices of Glob patterns for languages, versions, and roles.
type StringSlices struct {
	Languages []string `mapstructure:"languages" json:"languages"`
	Versions  []string `mapstructure:"versions" json:"versions"`
	Roles     []string `mapstructure:"roles" json:"roles"`
}

func (d StringSlices) Equal(other StringSlices) bool {
	return slices.Equal(d.Languages, other.Languages) &&
		slices.Equal(d.Versions, other.Versions) &&
		slices.Equal(d.Roles, other.Roles)
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

// Used in tests.
type testDimension struct {
	names []string
}

func (m testDimension) Len() int {
	return len(m.names)
}

func (m testDimension) IndexDefault() int {
	return 0
}

func (m *testDimension) ResolveIndex(name string) int {
	for i, n := range m.names {
		if n == name {
			return i
		}
	}
	return -1
}

func (m *testDimension) ResolveName(i int) string {
	if i < 0 || i >= len(m.names) {
		return ""
	}
	return m.names[i]
}

func (m *testDimension) ForEachIndex() iter.Seq[int] {
	return func(yield func(i int) bool) {
		for i := range m.names {
			if !yield(i) {
				return
			}
		}
	}
}

func (m *testDimension) IndexMatch(match predicate.P[string]) (iter.Seq[int], error) {
	return func(yield func(i int) bool) {
		for i, n := range m.names {
			if match(n) {
				if !yield(i) {
					return
				}
			}
		}
	}, nil
}

// NewTestingDimensions creates a new ConfiguredDimensions for testing.
func NewTestingDimensions(languages, versions, roles []string) *ConfiguredDimensions {
	c := &ConfiguredDimensions{
		ConfiguredLanguages: &testDimension{names: languages},
		ConfiguredVersions:  &testDimension{names: versions},
		ConfiguredRoles:     &testDimension{names: roles},
	}
	if err := c.Init(); err != nil {
		panic(err)
	}
	return c
}
