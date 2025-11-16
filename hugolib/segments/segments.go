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
	"strings"

	"github.com/gobwas/glob"
	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/common/predicate"
	"github.com/gohugoio/hugo/config"
	hglob "github.com/gohugoio/hugo/hugofs/hglob"
	"github.com/gohugoio/hugo/hugolib/sitesmatrix"
	"github.com/mitchellh/mapstructure"
)

// Segments is a collection of named segments.
type Segments struct {
	builder *segmentsBuilder

	// SegmentFilter is the compiled filter for all segments to render.
	SegmentFilter SegmentFilter
}

type SegmentFilter interface {
	// ShouldExcludeCoarse returns whether the given fields should be excluded on a coarse level.
	ShouldExcludeCoarse(SegmentQuery) bool

	// ShouldExcludeFine returns whether the given fields should be excluded on a fine level.
	ShouldExcludeFine(SegmentQuery) bool
}

type segmentFilter struct {
	exclude predicate.PR[SegmentQuery]
	include predicate.PR[SegmentQuery]
}

func (f segmentFilter) ShouldExcludeCoarse(q SegmentQuery) bool {
	return f.exclude(q).OK()
}

func (f segmentFilter) ShouldExcludeFine(q SegmentQuery) bool {
	return f.exclude(q).OK() || !f.include(q).OK()
}

type segmentsBuilder struct {
	logger               loggers.Logger
	isConfigInit         bool
	configuredDimensions *sitesmatrix.ConfiguredDimensions
	segmentCfg           map[string]SegmentConfig
	segmentsToRender     []string
}

func (b *Segments) compile() error {
	filter, err := b.builder.build()
	if err != nil {
		return err
	}
	b.SegmentFilter = filter
	b.builder = nil
	return nil
}

func (s *segmentsBuilder) buildOne(f []SegmentMatcherFields) (predicate.PR[SegmentQuery], error) {
	if f == nil {
		return matchNothing, nil
	}
	var (
		result  predicate.PR[SegmentQuery]
		section predicate.PR[SegmentQuery]
	)

	addSectionMatcher := func(matcher predicate.PR[SegmentQuery]) {
		if section == nil {
			section = matcher
		} else {
			section = section.And(matcher)
		}
	}

	addToSection := func(matcherFields SegmentMatcherFields, f1 func(fields SegmentMatcherFields) []string, f2 func(q SegmentQuery) string) error {
		s1 := f1(matcherFields)
		if s1 == nil {
			// Nothing to match against.
			return nil
		}
		var sliceMatcher predicate.PR[SegmentQuery]

		for _, s := range s1 {
			negate := strings.HasPrefix(s, hglob.NegationPrefix)
			if negate {
				s = strings.TrimPrefix(s, hglob.NegationPrefix)
			}

			g, err := getGlob(s)
			if err != nil {
				return err
			}

			m := func(fields SegmentQuery) predicate.Match {
				s2 := f2(fields)
				if s2 == "" {
					return predicate.False
				}
				return predicate.BoolMatch(g.Match(s2) != negate)
			}

			if negate {
				sliceMatcher = sliceMatcher.And(m)
			} else {
				sliceMatcher = sliceMatcher.Or(m)
			}
		}

		if sliceMatcher != nil {
			addSectionMatcher(sliceMatcher)
		}

		return nil
	}

	for _, fields := range f {
		if len(fields.Kind) > 0 {
			if err := addToSection(fields,
				func(fields SegmentMatcherFields) []string { return fields.Kind },
				func(fields SegmentQuery) string { return fields.Kind },
			); err != nil {
				return result, err
			}
		}
		if len(fields.Path) > 0 {
			if err := addToSection(fields,
				func(fields SegmentMatcherFields) []string { return fields.Path },
				func(fields SegmentQuery) string { return fields.Path },
			); err != nil {
				return result, err
			}
		}
		if fields.Lang != "" {
			hugo.DeprecateWithLogger("config segments.[...]lang ", "Use sites.matrix instead, see https://gohugo.io/configuration/segments/#segment-definition", "v0.153.0", s.logger.Logger())
			fields.Sites.Matrix.Languages = []string{fields.Lang}
		}
		if !fields.Sites.Matrix.IsZero() {
			intSetsCfg := sitesmatrix.IntSetsConfig{
				Globs: fields.Sites.Matrix,
			}
			matrix := sitesmatrix.NewIntSetsBuilder(s.configuredDimensions).WithConfig(intSetsCfg).WithAllIfNotSet().Build()

			addSectionMatcher(
				func(fields SegmentQuery) predicate.Match {
					return predicate.BoolMatch(matrix.HasVector(fields.Site))
				},
			)
		}
		if len(fields.Output) > 0 {
			if err := addToSection(fields,
				func(fields SegmentMatcherFields) []string { return fields.Output },
				func(fields SegmentQuery) string { return fields.Output },
			); err != nil {
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

func (s *segmentsBuilder) build() (SegmentFilter, error) {
	var sf segmentFilter

	for _, segID := range s.segmentsToRender {
		segCfg, ok := s.segmentCfg[segID]
		if !ok {
			continue
		}

		include, err := s.buildOne(segCfg.Includes)
		if err != nil {
			return nil, err
		}

		exclude, err := s.buildOne(segCfg.Excludes)
		if err != nil {
			return nil, err
		}
		if sf.include == nil {
			sf.include = include
		} else {
			sf.include = sf.include.Or(include)
		}
		if sf.exclude == nil {
			sf.exclude = exclude
		} else {
			sf.exclude = sf.exclude.Or(exclude)
		}
	}

	if sf.exclude == nil {
		sf.exclude = matchNothing
	}
	if sf.include == nil {
		sf.include = matchAll
	}

	return sf, nil
}

func (b *Segments) InitConfig(logger loggers.Logger, _ sitesmatrix.VectorStore, configuredDimensions *sitesmatrix.ConfiguredDimensions) error {
	if b.builder == nil || b.builder.isConfigInit {
		return nil
	}
	b.builder.isConfigInit = true
	b.builder.configuredDimensions = configuredDimensions
	return b.compile()
}

var (
	matchAll     = func(SegmentQuery) predicate.Match { return predicate.True }
	matchNothing = func(SegmentQuery) predicate.Match { return predicate.False }
)

type SegmentConfig struct {
	Excludes []SegmentMatcherFields
	Includes []SegmentMatcherFields
}

type SegmentQuery struct {
	Kind   string
	Path   string
	Output string
	Site   sitesmatrix.Vector
}

// SegmentMatcherFields holds string slices of ordered filters for segment matching.
// The Glob patterns can be negated by prefixing with "! ".
// The first match wins (either include or exclude).
type SegmentMatcherFields struct {
	Kind   []string
	Path   []string
	Output []string
	Lang   string            // Deprecated: use Sites.Matrix instead.
	Sites  sitesmatrix.Sites // Note that we only use Sites.Matrix for now.
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

func DecodeSegments(in map[string]any, segmentsToRender []string, logger loggers.Logger) (*config.ConfigNamespace[map[string]SegmentConfig, *Segments], error) {
	buildConfig := func(in any) (*Segments, any, error) {
		m, err := maps.ToStringMapE(in)
		if err != nil {
			return nil, nil, err
		}
		if m == nil {
			m = map[string]any{}
		}
		m = maps.CleanConfigStringMap(m)

		var segmentCfg map[string]SegmentConfig
		if err := mapstructure.WeakDecode(m, &segmentCfg); err != nil {
			return nil, nil, err
		}

		sms := &Segments{
			builder: &segmentsBuilder{
				logger:           logger,
				segmentCfg:       segmentCfg,
				segmentsToRender: segmentsToRender,
			},
		}

		return sms, nil, nil
	}

	ns, err := config.DecodeNamespace[map[string]SegmentConfig](in, buildConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to decode segments: %w", err)
	}
	return ns, nil
}
