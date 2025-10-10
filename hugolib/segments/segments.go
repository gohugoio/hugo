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
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/common/predicate"
	"github.com/gohugoio/hugo/config"
	hglob "github.com/gohugoio/hugo/hugofs/glob"
	"github.com/gohugoio/hugo/hugolib/sitesmatrix"
	"github.com/mitchellh/mapstructure"
)

// Segments is a collection of named segments.
type Segments struct {
	builder *segmentsBuilder

	// IncludeSegment is the compiled filter for all segments to render.
	IncludeSegment predicate.P[SegmentMatcherQuery]
}

type segmentsBuilder struct {
	isConfigInit         bool
	configuredDimensions *sitesmatrix.ConfiguredDimensions
	segmentCfg           map[string]SegmentConfig
	segmentsToRender     []string
	compiled             map[string]predicate.P[SegmentMatcherQuery]
}

func (b *Segments) compile() error {
	filter, err := b.builder.build()
	if err != nil {
		return err
	}
	b.IncludeSegment = filter
	b.builder = nil
	return nil
}

func (s *segmentsBuilder) build() (predicate.P[SegmentMatcherQuery], error) {
	s.compiled = make(map[string]predicate.P[SegmentMatcherQuery], len(s.segmentCfg))

	var hasLegacyExcludesOrIncludes bool
	for k, v := range s.segmentCfg {
		// In Hugo v0.152.0 we reworked excludes and includes to use rules instead.
		// To make that backwards compatible, we convert the old fields to rules here.
		// We need to start with the excludes, then the includes.
		hasLegacyExcludesOrIncludes = hasLegacyExcludesOrIncludes || len(v.Excludes) > 0 || len(v.Includes) > 0

		// Grouped per rule set.
		var rules []func(q SegmentMatcherQuery) (bool, bool)
		for _, r := range v.Rules {
			incl, err := compileShouldIncludeFilter(r, s.configuredDimensions)
			if err != nil {
				return nil, fmt.Errorf("failed to compile segment %q: %w", k, err)
			}
			rules = append(rules, incl)
		}

		include := func(q SegmentMatcherQuery) bool {
			for _, shouldInclude := range rules {
				ok, terminate := shouldInclude(q)
				if ok || terminate {
					return ok
				}
			}
			return false
		}

		s.compiled[k] = include
	}

	if hasLegacyExcludesOrIncludes {
		// I tried my best, but it's not poossible/practical to make this a warning.
		// The upside is that the new setup should be much easier to understand.
		return nil, fmt.Errorf("the use of segments.[id].{excludes,includes} was deprecated and removed in v0.152.0 to add support for multidimensional filtering; use segments.[id].rules instead, see https://gohugo.io/configuration/segments/#segment-definition")
	}
	return s.compileSegmentFilter()
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
	matchAll     = func(SegmentMatcherQuery) bool { return true }
	matchNothing = func(SegmentMatcherQuery) bool { return false }
)

func (b *segmentsBuilder) compileSegmentFilter() (predicate.P[SegmentMatcherQuery], error) {
	if b.segmentsToRender == nil {
		return matchAll, nil
	}

	var sf predicate.P[SegmentMatcherQuery]
	for _, s := range b.segmentsToRender {
		if seg, ok := b.compiled[s]; ok {
			if sf == nil {
				sf = seg
			} else {
				sf = sf.Or(seg)
			}
		}
	}

	if sf == nil {
		sf = matchNothing
	}

	return sf, nil
}

type SegmentConfig struct {
	Excludes []SegmentMatcherFields `json:"-"` // Deprecated: Use Rules.
	Includes []SegmentMatcherFields `json:"-"` // Deprecated: Use Rules.

	Rules []SegmentMatcherRules
}

// SegmentMatcherFields is a matcher for a segment include or exclude.
// All of these are Glob patterns.
// Deprecated: Use SegmentMatcherRules and SegmentMatcherQuery.
type SegmentMatcherFields struct {
	Kind   string
	Path   string
	Lang   string
	Output string
}

type SegmentMatcherQuery struct {
	Kind   string
	Path   string
	Output string
	Site   *sitesmatrix.Vector // May be nil.

	// TODO1 remove
	Dodebug bool
}

// SegmentMatcherRules holds string slices of ordered filters for segment matching.
// The Glob patterns can be negated by prefixing with "! ".
// The first match wins (either include or exclude).
type SegmentMatcherRules struct {
	Kind   []string
	Path   []string
	Output []string
	Sites  sitesmatrix.Sites // Note that we only use Sites.Matrix for now.
}

func (r SegmentMatcherRules) compile(configuredDimensions *sitesmatrix.ConfiguredDimensions) (segmentMatcherRulesCompiled, error) {
	compileStringPredicate := func(what string, ss []string) (func(string) (bool, bool), error) {
		if ss == nil {
			return nil, nil
		}
		type globExclude struct {
			glob    glob.Glob
			s       string
			exclude bool
		}
		var patterns []globExclude
		for _, s := range ss {
			var exclude bool
			if strings.HasPrefix(s, hglob.NegationPrefix) {
				exclude = true
				s = strings.TrimPrefix(s, hglob.NegationPrefix)
			}
			g, err := getGlob(s)
			if err != nil {
				return nil, err
			}

			patterns = append(patterns, globExclude{glob: g, exclude: exclude, s: s})

		}

		matcher := func(s string) (bool, bool) {
			if s == "" {
				return false, false
			}

			for _, pe := range patterns {
				g := pe.glob
				m := g.Match(s)
				if m {
					return !pe.exclude, m && pe.exclude
				}
			}
			return false, false
		}

		return matcher, nil
	}

	filter := segmentMatcherRulesCompiled{}

	var err error
	filter.Kind, err = compileStringPredicate("kind", r.Kind)
	if err != nil {
		return filter, err
	}
	filter.Path, err = compileStringPredicate("path", r.Path)
	if err != nil {
		return filter, err
	}
	filter.Output, err = compileStringPredicate("output", r.Output)
	if err != nil {
		return filter, err
	}
	if !r.Sites.Matrix.IsZero() {
		intSetsCfg := sitesmatrix.IntSetsConfig{
			Globs: r.Sites.Matrix,
		}
		matrix := sitesmatrix.NewIntSetsBuilder(configuredDimensions).WithConfig(intSetsCfg).WithAllIfNotSet().Build()

		filter.Sites = func(vec sitesmatrix.Vector) (bool, bool) {
			return matrix.HasVector(vec), true
		}
	} else {
		// Match all.
		filter.Sites = func(sitesmatrix.Vector) (bool, bool) { return true, false }
	}

	return filter, nil
}

type segmentMatcherRulesCompiled struct {
	Kind   func(string) (bool, bool)
	Path   func(string) (bool, bool)
	Output func(string) (bool, bool)
	Sites  func(sitesmatrix.Vector) (bool, bool)
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

func compileShouldIncludeFilter(rules SegmentMatcherRules, configuredDimensions *sitesmatrix.ConfiguredDimensions) (func(q SegmentMatcherQuery) (include bool, terminate bool), error) {
	c, err := rules.compile(configuredDimensions)
	if err != nil {
		return nil, err
	}

	return func(q SegmentMatcherQuery) (include bool, terminate bool) {
		if q.Kind != "" && c.Kind != nil {
			if b, terminate := c.Kind(q.Kind); !b || terminate {
				return b, terminate
			}
			include = true
		}

		if q.Path != "" && c.Path != nil {
			if b, terminate := c.Path(q.Path); !b || terminate {
				return b, terminate
			}
			include = true
		}

		if q.Output != "" && c.Output != nil {
			if b, terminate := c.Output(q.Output); !b || terminate {
				return b, terminate
			}
			include = true
		}

		if q.Site != nil && c.Sites != nil {
			if b, terminate := c.Sites(*q.Site); !b || terminate {
				return b, terminate
			}
			include = true
		}

		return
	}, nil
}

func DecodeSegments(in map[string]any, segmentsToRender []string) (*config.ConfigNamespace[map[string]SegmentConfig, *Segments], error) {
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
