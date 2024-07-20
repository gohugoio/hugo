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

package segments

import (
	"fmt"

	"github.com/gobwas/glob"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/common/predicate"
	"github.com/gohugoio/hugo/config"
	hglob "github.com/gohugoio/hugo/hugofs/glob"
	"github.com/mitchellh/mapstructure"
)

// Segments is a collection of named segments.
type Segments struct {
	s map[string]excludeInclude
}

type excludeInclude struct {
	exclude predicate.P[SegmentMatcherFields]
	include predicate.P[SegmentMatcherFields]
}

// ShouldExcludeCoarse returns whether the given fields should be excluded.
// This is used for the coarser grained checks, e.g. language and output format.
// Note that ShouldExcludeCoarse(fields) == ShouldExcludeFine(fields) may
// not always be true, but ShouldExcludeCoarse(fields) == true == ShouldExcludeFine(fields)
// will always be truthful.
func (e excludeInclude) ShouldExcludeCoarse(fields SegmentMatcherFields) bool {
	return e.exclude != nil && e.exclude(fields)
}

// ShouldExcludeFine returns whether the given fields should be excluded.
// This is used for the finer grained checks, e.g. on invididual pages.
func (e excludeInclude) ShouldExcludeFine(fields SegmentMatcherFields) bool {
	if e.exclude != nil && e.exclude(fields) {
		return true
	}
	return e.include != nil && !e.include(fields)
}

type SegmentFilter interface {
	// ShouldExcludeCoarse returns whether the given fields should be excluded on a coarse level.
	ShouldExcludeCoarse(SegmentMatcherFields) bool

	// ShouldExcludeFine returns whether the given fields should be excluded on a fine level.
	ShouldExcludeFine(SegmentMatcherFields) bool
}

type segmentFilter struct {
	coarse predicate.P[SegmentMatcherFields]
	fine   predicate.P[SegmentMatcherFields]
}

func (f segmentFilter) ShouldExcludeCoarse(field SegmentMatcherFields) bool {
	return f.coarse(field)
}

func (f segmentFilter) ShouldExcludeFine(fields SegmentMatcherFields) bool {
	return f.fine(fields)
}

var (
	matchAll     = func(SegmentMatcherFields) bool { return true }
	matchNothing = func(SegmentMatcherFields) bool { return false }
)

// Get returns a SegmentFilter for the given segments.
func (sms Segments) Get(onNotFound func(s string), ss ...string) SegmentFilter {
	if ss == nil {
		return segmentFilter{coarse: matchNothing, fine: matchNothing}
	}
	var sf segmentFilter
	for _, s := range ss {
		if seg, ok := sms.s[s]; ok {
			if sf.coarse == nil {
				sf.coarse = seg.ShouldExcludeCoarse
			} else {
				sf.coarse = sf.coarse.Or(seg.ShouldExcludeCoarse)
			}
			if sf.fine == nil {
				sf.fine = seg.ShouldExcludeFine
			} else {
				sf.fine = sf.fine.Or(seg.ShouldExcludeFine)
			}
		} else if onNotFound != nil {
			onNotFound(s)
		}
	}

	if sf.coarse == nil {
		sf.coarse = matchAll
	}
	if sf.fine == nil {
		sf.fine = matchAll
	}

	return sf
}

type SegmentConfig struct {
	Excludes []SegmentMatcherFields
	Includes []SegmentMatcherFields
}

// SegmentMatcherFields is a matcher for a segment include or exclude.
// All of these are Glob patterns.
type SegmentMatcherFields struct {
	Kind   string
	Path   string
	Lang   string
	Output string
}

func getGlob(s string) (glob.Glob, error) {
	if s == "" {
		return nil, nil
	}
	g, err := hglob.GetGlob(s)
	if err != nil {
		return nil, fmt.Errorf("failed to compile Glob %q: %w", s, err)
	}
	return g, nil
}

func compileSegments(f []SegmentMatcherFields) (predicate.P[SegmentMatcherFields], error) {
	if f == nil {
		return func(SegmentMatcherFields) bool { return false }, nil
	}
	var (
		result  predicate.P[SegmentMatcherFields]
		section predicate.P[SegmentMatcherFields]
	)

	addToSection := func(matcherFields SegmentMatcherFields, f func(fields SegmentMatcherFields) string) error {
		s1 := f(matcherFields)
		g, err := getGlob(s1)
		if err != nil {
			return err
		}
		matcher := func(fields SegmentMatcherFields) bool {
			s2 := f(fields)
			if s2 == "" {
				return false
			}
			return g.Match(s2)
		}
		if section == nil {
			section = matcher
		} else {
			section = section.And(matcher)
		}
		return nil
	}

	for _, fields := range f {
		if fields.Kind != "" {
			if err := addToSection(fields, func(fields SegmentMatcherFields) string { return fields.Kind }); err != nil {
				return result, err
			}
		}
		if fields.Path != "" {
			if err := addToSection(fields, func(fields SegmentMatcherFields) string { return fields.Path }); err != nil {
				return result, err
			}
		}
		if fields.Lang != "" {
			if err := addToSection(fields, func(fields SegmentMatcherFields) string { return fields.Lang }); err != nil {
				return result, err
			}
		}
		if fields.Output != "" {
			if err := addToSection(fields, func(fields SegmentMatcherFields) string { return fields.Output }); err != nil {
				return result, err
			}
		}

		if result == nil {
			result = section
		} else {
			result = result.Or(section)
		}
		section = nil

	}

	return result, nil
}

func DecodeSegments(in map[string]any) (*config.ConfigNamespace[map[string]SegmentConfig, Segments], error) {
	buildConfig := func(in any) (Segments, any, error) {
		sms := Segments{
			s: map[string]excludeInclude{},
		}
		m, err := maps.ToStringMapE(in)
		if err != nil {
			return sms, nil, err
		}
		if m == nil {
			m = map[string]any{}
		}
		m = maps.CleanConfigStringMap(m)

		var scfgm map[string]SegmentConfig
		if err := mapstructure.Decode(m, &scfgm); err != nil {
			return sms, nil, err
		}

		for k, v := range scfgm {
			var (
				include predicate.P[SegmentMatcherFields]
				exclude predicate.P[SegmentMatcherFields]
				err     error
			)
			if v.Excludes != nil {
				exclude, err = compileSegments(v.Excludes)
				if err != nil {
					return sms, nil, err
				}
			}
			if v.Includes != nil {
				include, err = compileSegments(v.Includes)
				if err != nil {
					return sms, nil, err
				}
			}

			ei := excludeInclude{
				exclude: exclude,
				include: include,
			}
			sms.s[k] = ei

		}

		return sms, nil, nil
	}

	ns, err := config.DecodeNamespace[map[string]SegmentConfig](in, buildConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to decode segments: %w", err)
	}
	return ns, nil
}
