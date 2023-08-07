// Copyright 2019 The Hugo Authors. All rights reserved.
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

package markup_config

import (
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/markup/asciidocext/asciidocext_config"
	"github.com/gohugoio/hugo/markup/goldmark/goldmark_config"
	"github.com/gohugoio/hugo/markup/highlight"
	"github.com/gohugoio/hugo/markup/tableofcontents"
	"github.com/mitchellh/mapstructure"
)

type Config struct {
	// Default markdown handler for md/markdown extensions.
	// Default is "goldmark".
	DefaultMarkdownHandler string

	// The configuration used by code highlighters.
	Highlight highlight.Config

	// Table of contents configuration
	TableOfContents tableofcontents.Config

	// Configuration for the Goldmark markdown engine.
	Goldmark goldmark_config.Config

	// Configuration for the Asciidoc external markdown engine.
	AsciidocExt asciidocext_config.Config
}

func Decode(cfg config.Provider) (conf Config, err error) {
	conf = Default

	m := cfg.GetStringMap("markup")
	if m == nil {
		return
	}
	m = maps.CleanConfigStringMap(m)

	normalizeConfig(m)

	err = mapstructure.WeakDecode(m, &conf)
	if err != nil {
		return
	}

	if err = highlight.ApplyLegacyConfig(cfg, &conf.Highlight); err != nil {
		return
	}

	return
}

func normalizeConfig(m map[string]any) {
	v, err := maps.GetNestedParam("goldmark.parser", ".", m)
	if err == nil {
		vm := maps.ToStringMap(v)
		// Changed from a bool in 0.81.0
		if vv, found := vm["attribute"]; found {
			if vvb, ok := vv.(bool); ok {
				vm["attribute"] = goldmark_config.ParserAttribute{
					Title: vvb,
				}
			}
		}
	}

	// Changed from a bool in 0.112.0.
	v, err = maps.GetNestedParam("goldmark.extensions", ".", m)
	if err == nil {
		vm := maps.ToStringMap(v)
		const typographerKey = "typographer"
		if vv, found := vm[typographerKey]; found {
			if vvb, ok := vv.(bool); ok {
				if !vvb {
					vm[typographerKey] = goldmark_config.Typographer{
						Disable: true,
					}
				} else {
					delete(vm, typographerKey)
				}
			}
		}
	}
}

var Default = Config{
	DefaultMarkdownHandler: "goldmark",

	TableOfContents: tableofcontents.DefaultConfig,
	Highlight:       highlight.DefaultConfig,

	Goldmark:    goldmark_config.Default,
	AsciidocExt: asciidocext_config.Default,
}
