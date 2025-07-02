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
	"path/filepath"
	"slices"
	"strings"

	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/hugofs/glob"
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
	Lang string

	// A Glob pattern matching the Page's Environment, e.g. "{production,development}".
	Environment string
}

// Matches returns whether p matches this matcher.
func (m PageMatcher) Matches(p Page) bool {
	if m.Kind != "" {
		g, err := glob.GetGlob(m.Kind)
		if err == nil && !g.Match(p.Kind()) {
			return false
		}
	}

	if m.Lang != "" {
		g, err := glob.GetGlob(m.Lang)
		if err == nil && !g.Match(p.Lang()) {
			return false
		}
	}

	if m.Path != "" {
		g, err := glob.GetGlob(m.Path)
		// TODO(bep) Path() vs filepath vs leading slash.
		p := strings.ToLower(filepath.ToSlash(p.Path()))
		if !(strings.HasPrefix(p, "/")) {
			p = "/" + p
		}
		if err == nil && !g.Match(p) {
			return false
		}
	}

	if m.Environment != "" {
		g, err := glob.GetGlob(m.Environment)
		if err == nil && !g.Match(p.Site().Hugo().Environment) {
			return false
		}
	}

	return true
}

var disallowedCascadeKeys = map[string]bool{
	// These define the structure of the page tree and cannot
	// currently be set in the cascade.
	"kind": true,
	"path": true,
	"lang": true,
}

// See issue 11977.
func isGlobWithExtension(s string) bool {
	pathParts := strings.Split(s, "/")
	last := pathParts[len(pathParts)-1]
	return strings.Count(last, ".") > 0
}

func CheckCascadePattern(logger loggers.Logger, m PageMatcher) {
	if logger != nil && isGlobWithExtension(m.Path) {
		logger.Erroridf("cascade-pattern-with-extension", "cascade target path %q looks like a path with an extension; since Hugo v0.123.0 this will not match anything, see  https://gohugo.io/methods/page/path/", m.Path)
	}
}

func DecodeCascadeConfig(logger loggers.Logger, handleLegacyFormat bool, in any) (*config.ConfigNamespace[[]PageMatcherParamsConfig, *maps.Ordered[PageMatcher, PageMatcherParamsConfig]], error) {
	buildConfig := func(in any) (*maps.Ordered[PageMatcher, PageMatcherParamsConfig], any, error) {
		cascade := maps.NewOrdered[PageMatcher, PageMatcherParamsConfig]()
		if in == nil {
			return cascade, []map[string]any{}, nil
		}
		ms, err := maps.ToSliceStringMap(in)
		if err != nil {
			return nil, nil, err
		}

		var cfgs []PageMatcherParamsConfig

		for _, m := range ms {
			m = maps.CleanConfigStringMap(m)
			var (
				c   PageMatcherParamsConfig
				err error
			)
			c, err = mapToPageMatcherParamsConfig(m)
			if err != nil {
				return nil, nil, err
			}
			for k := range m {
				if disallowedCascadeKeys[k] {
					return nil, nil, fmt.Errorf("key %q not allowed in cascade config", k)
				}
			}
			cfgs = append(cfgs, c)
		}

		for _, cfg := range cfgs {
			m := cfg.Target
			CheckCascadePattern(logger, m)
			c, found := cascade.Get(m)
			if found {
				// Merge
				for k, v := range cfg.Params {
					if _, found := c.Params[k]; !found {
						c.Params[k] = v
					}
				}
				for k, v := range cfg.Fields {
					if _, found := c.Fields[k]; !found {
						c.Fields[k] = v
					}
				}
			} else {
				cascade.Set(m, cfg)
			}
		}

		return cascade, cfgs, nil
	}

	return config.DecodeNamespace[[]PageMatcherParamsConfig, *maps.Ordered[PageMatcher, PageMatcherParamsConfig]](in, buildConfig)
}

// DecodeCascade decodes in which could be either a map or a slice of maps.
func DecodeCascade(logger loggers.Logger, handleLegacyFormat bool, in any) (*maps.Ordered[PageMatcher, PageMatcherParamsConfig], error) {
	conf, err := DecodeCascadeConfig(logger, handleLegacyFormat, in)
	if err != nil {
		return nil, err
	}
	return conf.Config, nil
}

func mapToPageMatcherParamsConfig(m map[string]any) (PageMatcherParamsConfig, error) {
	var pcfg PageMatcherParamsConfig
	if pcfg.Fields == nil {
		pcfg.Fields = make(maps.Params)
	}
	for k, v := range m {
		switch strings.ToLower(k) {
		case "_target", "target":
			var target PageMatcher
			if err := decodePageMatcher(v, &target); err != nil {
				return pcfg, err
			}
			pcfg.Target = target
		case "params":
			if pcfg.Params == nil {
				pcfg.Params = make(maps.Params)
			}
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
func decodePageMatcher(m any, v *PageMatcher) error {
	if err := mapstructure.WeakDecode(m, v); err != nil {
		return err
	}

	v.Kind = strings.ToLower(v.Kind)
	if v.Kind != "" {
		g, _ := glob.GetGlob(v.Kind)
		found := slices.ContainsFunc(kinds.AllKindsInPages, g.Match)
		if !found {
			return fmt.Errorf("%q did not match a valid Page Kind", v.Kind)
		}
	}

	v.Path = filepath.ToSlash(strings.ToLower(v.Path))

	return nil
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
