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
	"fmt"
	"iter"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gohugoio/hugo/common/hashing"
	"github.com/gohugoio/hugo/common/hstrings"
	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/hugofs/hglob"
	"github.com/gohugoio/hugo/hugolib/sitesmatrix"
	"github.com/gohugoio/hugo/resources/kinds"
	"github.com/mitchellh/mapstructure"
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
	// Deprecated: use Sites.Matrix instead.
	Lang string

	// The sites to apply this to.
	// Note that we currently only use the Matrix field for cascade matching.
	Sites sitesmatrix.Sites

	// A Glob pattern matching the Page's Environment, e.g. "{production,development}".
	Environment string

	// Compiled values.
	// The site vectors to apply this to.
	SitesMatrixCompiled sitesmatrix.VectorProvider `mapstructure:"-"`
}

func (m PageMatcher) Matches(p Page) bool {
	return m.Match(p.Kind(), p.Path(), p.Site().Hugo().Environment, nil)
}

func (m PageMatcher) Match(kind, path, environment string, sitesMatrix sitesmatrix.VectorProvider) bool {
	if sitesMatrix != nil {
		if m.SitesMatrixCompiled != nil && !m.SitesMatrixCompiled.HasAnyVector(sitesMatrix) {
			return false
		}
	}
	if m.Kind != "" {
		g, err := hglob.GetGlob(m.Kind)
		if err == nil && !g.Match(kind) {
			return false
		}
	}

	if m.Path != "" {
		g, err := hglob.GetGlob(m.Path)
		// TODO(bep) Path() vs filepath vs leading slash.
		p := strings.ToLower(filepath.ToSlash(path))
		if !(strings.HasPrefix(p, "/")) {
			p = "/" + p
		}
		if err == nil && !g.Match(p) {
			return false
		}
	}

	if m.Environment != "" {
		g, err := hglob.GetGlob(m.Environment)
		if err == nil && !g.Match(environment) {
			return false
		}
	}

	return true
}

var disallowedCascadeKeys = map[string]bool{
	// These define the structure of the page tree and cannot
	// currently be set in the cascade.
	"kind":    true,
	"path":    true,
	"lang":    true,
	"cascade": true,
}

// See issue 11977.
func isGlobWithExtension(s string) bool {
	pathParts := strings.Split(s, "/")
	last := pathParts[len(pathParts)-1]
	return strings.Count(last, ".") > 0
}

func checkCascadePattern(logger loggers.Logger, m PageMatcher) {
	if m.Lang != "" {
		hugo.Deprecate("cascade.target.language", "cascade.target.sites.matrix instead, see https://gohugo.io/content-management/front-matter/#target", "v0.150.0")
	}
}

func AddLangToCascadeTargetMap(lang string, m maps.Params) {
	maps.SetNestedParamIfNotSet("target.sites.matrix.languages", ".", lang, m)
}

func DecodeCascadeConfig(in any) (*PageMatcherParamsConfigs, error) {
	buildConfig := func(in any) (CascadeConfig, any, error) {
		dec := cascadeConfigDecoder{}

		if in == nil {
			return CascadeConfig{}, []map[string]any{}, nil
		}

		ms, err := maps.ToSliceStringMap(in)
		if err != nil {
			return CascadeConfig{}, nil, err
		}

		var cfgs []PageMatcherParamsConfig

		for _, m := range ms {
			m = maps.CleanConfigStringMap(m)
			var (
				c   PageMatcherParamsConfig
				err error
			)
			c, err = dec.mapToPageMatcherParamsConfig(m)
			if err != nil {
				return CascadeConfig{}, nil, err
			}
			for k := range m {
				if disallowedCascadeKeys[k] {
					return CascadeConfig{}, nil, fmt.Errorf("key %q not allowed in cascade config", k)
				}
			}
			cfgs = append(cfgs, c)
		}

		if len(cfgs) == 0 {
			return CascadeConfig{}, nil, nil
		}

		var n int
		for _, cfg := range cfgs {
			if len(cfg.Params) > 0 || len(cfg.Fields) > 0 {
				cfgs[n] = cfg
				n++
			}
		}

		if n == 0 {
			return CascadeConfig{}, nil, nil
		}

		cfgs = cfgs[:n]

		return CascadeConfig{Cascades: cfgs}, cfgs, nil
	}

	c, err := config.DecodeNamespace[[]PageMatcherParamsConfig](in, buildConfig)
	if err != nil || len(c.Config.Cascades) == 0 {
		return nil, err
	}

	return &PageMatcherParamsConfigs{c: []*config.ConfigNamespace[[]PageMatcherParamsConfig, CascadeConfig]{c}}, nil
}

type cascadeConfigDecoder struct{}

func (d cascadeConfigDecoder) mapToPageMatcherParamsConfig(m map[string]any) (PageMatcherParamsConfig, error) {
	var pcfg PageMatcherParamsConfig
	if pcfg.Fields == nil {
		pcfg.Fields = make(maps.Params)
	}
	if pcfg.Params == nil {
		pcfg.Params = make(maps.Params)
	}

	for k, v := range m {
		switch strings.ToLower(k) {
		case "_target", "target":
			var target PageMatcher
			if err := d.decodePageMatcher(v, &target); err != nil {
				return pcfg, err
			}
			pcfg.Target = target
		case "params":
			params := maps.ToStringMap(v)
			for k, v := range params {
				if _, found := pcfg.Params[k]; !found {
					pcfg.Params[k] = v
				}
			}
		default:
			pcfg.Fields[k] = v
		}
	}
	return pcfg, pcfg.init()
}

// decodePageMatcher decodes m into v.
func (d cascadeConfigDecoder) decodePageMatcher(m any, v *PageMatcher) error {
	if err := mapstructure.WeakDecode(m, v); err != nil {
		return err
	}

	v.Kind = strings.ToLower(v.Kind)
	if v.Kind != "" {
		g, _ := hglob.GetGlob(v.Kind)
		found := slices.ContainsFunc(kinds.AllKindsInPages, g.Match)
		if !found {
			return fmt.Errorf("%q did not match a valid Page Kind", v.Kind)
		}
	}

	v.Path = filepath.ToSlash(strings.ToLower(v.Path))

	if v.Lang != "" {
		v.Sites.Matrix.Languages = append(v.Sites.Matrix.Languages, v.Lang)
		v.Sites.Matrix.Languages = hstrings.UniqueStringsReuse(v.Sites.Matrix.Languages)
	}

	return nil
}

// DecodeCascadeConfigOptions
func (v *PageMatcher) compileSitesMatrix(configuredDimensions *sitesmatrix.ConfiguredDimensions) error {
	if v.Sites.Matrix.IsZero() {
		// Nothing to do.
		v.SitesMatrixCompiled = nil
		return nil
	}
	intSetsCfg := sitesmatrix.IntSetsConfig{
		Globs: v.Sites.Matrix,
	}
	b := sitesmatrix.NewIntSetsBuilder(configuredDimensions).WithConfig(intSetsCfg).WithAllIfNotSet()

	v.SitesMatrixCompiled = b.Build()
	return nil
}

type CascadeConfig struct {
	Cascades []PageMatcherParamsConfig
}

type PageMatcherParamsConfig struct {
	// Apply Params to all Pages matching Target.
	Params maps.Params
	// Fields holds all fields but Params.
	Fields maps.Params
	// Target is the PageMatcher that this config applies to.
	Target PageMatcher
}

func (p *PageMatcherParamsConfig) init() error {
	maps.PrepareParams(p.Params)
	maps.PrepareParams(p.Fields)

	return nil
}

type PageMatcherParamsConfigs struct {
	c []*config.ConfigNamespace[[]PageMatcherParamsConfig, CascadeConfig]
}

func (c *PageMatcherParamsConfigs) Append(other *PageMatcherParamsConfigs) *PageMatcherParamsConfigs {
	if c == nil || len(c.c) == 0 {
		return other
	}
	if other == nil || len(other.c) == 0 {
		return c
	}
	return &PageMatcherParamsConfigs{c: slices.Concat(c.c, other.c)}
}

func (c *PageMatcherParamsConfigs) Prepend(other *PageMatcherParamsConfigs) *PageMatcherParamsConfigs {
	if c == nil || len(c.c) == 0 {
		return other
	}
	if other == nil || len(other.c) == 0 {
		return c
	}
	return &PageMatcherParamsConfigs{c: slices.Concat(other.c, c.c)}
}

func (c *PageMatcherParamsConfigs) All() iter.Seq[PageMatcherParamsConfig] {
	if c == nil {
		return func(func(PageMatcherParamsConfig) bool) {}
	}
	return func(yield func(PageMatcherParamsConfig) bool) {
		if c == nil {
			return
		}
		for _, v := range c.c {
			for _, vv := range v.Config.Cascades {
				if !yield(vv) {
					return
				}
			}
		}
	}
}

func (c *PageMatcherParamsConfigs) Len() int {
	if c == nil {
		return 0
	}
	var n int
	for _, v := range c.c {
		n += len(v.Config.Cascades)
	}
	return n
}

func (c *PageMatcherParamsConfigs) SourceHash() uint64 {
	if c == nil {
		return 0
	}
	h := hashing.XxHasher()
	defer h.Close()

	for _, v := range c.c {
		h.WriteString(v.SourceHash)
	}
	return h.Sum64()
}

func (c *PageMatcherParamsConfigs) InitConfig(logger loggers.Logger, _ sitesmatrix.VectorStore, configuredDimensions *sitesmatrix.ConfiguredDimensions) error {
	if c == nil {
		return nil
	}
	for _, cc := range c.c {
		for i := range cc.Config.Cascades {
			checkCascadePattern(logger, cc.Config.Cascades[i].Target)
			if err := cc.Config.Cascades[i].Target.compileSitesMatrix(configuredDimensions); err != nil {
				return fmt.Errorf("failed to compile cascade target %d: %w", i, err)
			}
		}
	}
	return nil
}
