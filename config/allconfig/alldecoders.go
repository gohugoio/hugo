// Copyright 2023 The Hugo Authors. All rights reserved.
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

package allconfig

import (
	"fmt"
	"strings"

	"github.com/gohugoio/hugo/cache/filecache"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/common/types"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/config/privacy"
	"github.com/gohugoio/hugo/config/security"
	"github.com/gohugoio/hugo/config/services"
	"github.com/gohugoio/hugo/deploy"
	"github.com/gohugoio/hugo/langs"
	"github.com/gohugoio/hugo/markup/markup_config"
	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/minifiers"
	"github.com/gohugoio/hugo/modules"
	"github.com/gohugoio/hugo/navigation"
	"github.com/gohugoio/hugo/output"
	"github.com/gohugoio/hugo/related"
	"github.com/gohugoio/hugo/resources/images"
	"github.com/gohugoio/hugo/resources/page"
	"github.com/gohugoio/hugo/resources/page/pagemeta"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/afero"
	"github.com/spf13/cast"
)

type decodeConfig struct {
	p    config.Provider
	c    *Config
	fs   afero.Fs
	bcfg config.BaseConfig
}

type decodeWeight struct {
	key         string
	decode      func(decodeWeight, decodeConfig) error
	getCompiler func(c *Config) configCompiler
	weight      int
}

var allDecoderSetups = map[string]decodeWeight{
	"": {
		key:    "",
		weight: -100, // Always first.
		decode: func(d decodeWeight, p decodeConfig) error {
			return mapstructure.WeakDecode(p.p.Get(""), &p.c.RootConfig)
		},
	},
	"imaging": {
		key: "imaging",
		decode: func(d decodeWeight, p decodeConfig) error {
			var err error
			p.c.Imaging, err = images.DecodeConfig(p.p.GetStringMap(d.key))
			return err
		},
	},
	"caches": {
		key: "caches",
		decode: func(d decodeWeight, p decodeConfig) error {
			var err error
			p.c.Caches, err = filecache.DecodeConfig(p.fs, p.bcfg, p.p.GetStringMap(d.key))
			if p.c.IgnoreCache {
				// Set MaxAge in all caches to 0.
				for k, cache := range p.c.Caches {
					cache.MaxAge = 0
					p.c.Caches[k] = cache
				}
			}
			return err
		},
	},
	"build": {
		key: "build",
		decode: func(d decodeWeight, p decodeConfig) error {
			p.c.Build = config.DecodeBuildConfig(p.p)
			return nil
		},
		getCompiler: func(c *Config) configCompiler {
			return &c.Build
		},
	},
	"frontmatter": {
		key: "frontmatter",
		decode: func(d decodeWeight, p decodeConfig) error {
			var err error
			p.c.Frontmatter, err = pagemeta.DecodeFrontMatterConfig(p.p)
			return err
		},
	},
	"markup": {
		key: "markup",
		decode: func(d decodeWeight, p decodeConfig) error {
			var err error
			p.c.Markup, err = markup_config.Decode(p.p)
			return err
		},
	},
	"server": {
		key: "server",
		decode: func(d decodeWeight, p decodeConfig) error {
			var err error
			p.c.Server, err = config.DecodeServer(p.p)
			return err
		},
		getCompiler: func(c *Config) configCompiler {
			return &c.Server
		},
	},
	"minify": {
		key: "minify",
		decode: func(d decodeWeight, p decodeConfig) error {
			var err error
			p.c.Minify, err = minifiers.DecodeConfig(p.p.Get(d.key))
			return err
		},
	},
	"mediaTypes": {
		key: "mediaTypes",
		decode: func(d decodeWeight, p decodeConfig) error {
			var err error
			p.c.MediaTypes, err = media.DecodeTypes(p.p.GetStringMap(d.key))
			return err
		},
	},
	"outputs": {
		key: "outputs",
		decode: func(d decodeWeight, p decodeConfig) error {
			defaults := createDefaultOutputFormats(p.c.OutputFormats.Config)
			m := p.p.GetStringMap("outputs")
			p.c.Outputs = make(map[string][]string)
			for k, v := range m {
				s := types.ToStringSlicePreserveString(v)
				for i, v := range s {
					s[i] = strings.ToLower(v)
				}
				p.c.Outputs[k] = s
			}
			// Apply defaults.
			for k, v := range defaults {
				if _, found := p.c.Outputs[k]; !found {
					p.c.Outputs[k] = v
				}
			}
			return nil
		},
	},
	"outputFormats": {
		key: "outputFormats",
		decode: func(d decodeWeight, p decodeConfig) error {
			var err error
			p.c.OutputFormats, err = output.DecodeConfig(p.c.MediaTypes.Config, p.p.Get(d.key))
			return err
		},
	},
	"params": {
		key: "params",
		decode: func(d decodeWeight, p decodeConfig) error {
			p.c.Params = maps.CleanConfigStringMap(p.p.GetStringMap("params"))
			if p.c.Params == nil {
				p.c.Params = make(map[string]any)
			}

			// Before Hugo 0.112.0 this was configured via site Params.
			if mainSections, found := p.c.Params["mainsections"]; found {
				p.c.MainSections = types.ToStringSlicePreserveString(mainSections)
				if p.c.MainSections == nil {
					p.c.MainSections = []string{}
				}
			}

			return nil
		},
	},
	"module": {
		key: "module",
		decode: func(d decodeWeight, p decodeConfig) error {
			var err error
			p.c.Module, err = modules.DecodeConfig(p.p)
			return err
		},
	},
	"permalinks": {
		key: "permalinks",
		decode: func(d decodeWeight, p decodeConfig) error {
			p.c.Permalinks = maps.CleanConfigStringMapString(p.p.GetStringMapString(d.key))
			return nil
		},
	},
	"sitemap": {
		key: "sitemap",
		decode: func(d decodeWeight, p decodeConfig) error {
			var err error
			p.c.Sitemap, err = config.DecodeSitemap(config.SitemapConfig{Priority: -1, Filename: "sitemap.xml"}, p.p.GetStringMap(d.key))
			return err
		},
	},
	"taxonomies": {
		key: "taxonomies",
		decode: func(d decodeWeight, p decodeConfig) error {
			p.c.Taxonomies = maps.CleanConfigStringMapString(p.p.GetStringMapString(d.key))
			return nil
		},
	},
	"related": {
		key:    "related",
		weight: 100, // This needs to be decoded after taxonomies.
		decode: func(d decodeWeight, p decodeConfig) error {
			if p.p.IsSet(d.key) {
				var err error
				p.c.Related, err = related.DecodeConfig(p.p.GetParams(d.key))
				if err != nil {
					return fmt.Errorf("failed to decode related config: %w", err)
				}
			} else {
				p.c.Related = related.DefaultConfig
				if _, found := p.c.Taxonomies["tag"]; found {
					p.c.Related.Add(related.IndexConfig{Name: "tags", Weight: 80})
				}
			}
			return nil
		},
	},
	"languages": {
		key: "languages",
		decode: func(d decodeWeight, p decodeConfig) error {
			var err error
			p.c.Languages, err = langs.DecodeConfig(p.p.GetStringMap(d.key))
			return err
		},
	},
	"cascade": {
		key: "cascade",
		decode: func(d decodeWeight, p decodeConfig) error {
			var err error
			p.c.Cascade, err = page.DecodeCascadeConfig(p.p.Get(d.key))
			return err
		},
	},
	"menus": {
		key: "menus",
		decode: func(d decodeWeight, p decodeConfig) error {
			var err error
			p.c.Menus, err = navigation.DecodeConfig(p.p.Get(d.key))
			return err
		},
	},
	"privacy": {
		key: "privacy",
		decode: func(d decodeWeight, p decodeConfig) error {
			var err error
			p.c.Privacy, err = privacy.DecodeConfig(p.p)
			return err
		},
	},
	"security": {
		key: "security",
		decode: func(d decodeWeight, p decodeConfig) error {
			var err error
			p.c.Security, err = security.DecodeConfig(p.p)
			return err
		},
	},
	"services": {
		key: "services",
		decode: func(d decodeWeight, p decodeConfig) error {
			var err error
			p.c.Services, err = services.DecodeConfig(p.p)
			return err
		},
	},
	"deployment": {
		key: "deployment",
		decode: func(d decodeWeight, p decodeConfig) error {
			var err error
			p.c.Deployment, err = deploy.DecodeConfig(p.p)
			return err
		},
	},
	"author": {
		key: "author",
		decode: func(d decodeWeight, p decodeConfig) error {
			p.c.Author = p.p.GetStringMap(d.key)
			return nil
		},
	},
	"social": {
		key: "social",
		decode: func(d decodeWeight, p decodeConfig) error {
			p.c.Social = p.p.GetStringMapString(d.key)
			return nil
		},
	},
	"uglyurls": {
		key: "uglyurls",
		decode: func(d decodeWeight, p decodeConfig) error {
			v := p.p.Get(d.key)
			switch vv := v.(type) {
			case bool:
				p.c.UglyURLs = vv
			case string:
				p.c.UglyURLs = vv == "true"
			default:
				p.c.UglyURLs = cast.ToStringMapBool(v)
			}
			return nil
		},
	},
	"internal": {
		key: "internal",
		decode: func(d decodeWeight, p decodeConfig) error {
			return mapstructure.WeakDecode(p.p.GetStringMap(d.key), &p.c.Internal)
		},
	},
}
