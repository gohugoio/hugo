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

package config

import (
	"fmt"
	"sort"
	"strings"

	"github.com/gohugoio/hugo/common/types"

	"github.com/gobwas/glob"
	"github.com/gohugoio/hugo/common/herrors"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
)

type BaseConfig struct {
	WorkingDir string
	CacheDir   string
	ThemesDir  string
	PublishDir string
}

type CommonDirs struct {
	// The directory where Hugo will look for themes.
	ThemesDir string

	// Where to put the generated files.
	PublishDir string

	// The directory to put the generated resources files. This directory should in most situations be considered temporary
	// and not be committed to version control. But there may be cached content in here that you want to keep,
	// e.g. resources/_gen/images for performance reasons or CSS built from SASS when your CI server doesn't have the full setup.
	ResourceDir string

	// The project root directory.
	WorkingDir string

	// The root directory for all cache files.
	CacheDir string

	// The content source directory.
	// Deprecated: Use module mounts.
	ContentDir string
	// Deprecated: Use module mounts.
	// The data source directory.
	DataDir string
	// Deprecated: Use module mounts.
	// The layout source directory.
	LayoutDir string
	// Deprecated: Use module mounts.
	// The i18n source directory.
	I18nDir string
	// Deprecated: Use module mounts.
	// The archetypes source directory.
	ArcheTypeDir string
	// Deprecated: Use module mounts.
	// The assets source directory.
	AssetDir string
}

type LoadConfigResult struct {
	Cfg         Provider
	ConfigFiles []string
	BaseConfig  BaseConfig
}

var DefaultBuild = BuildConfig{
	UseResourceCacheWhen: "fallback",
	WriteStats:           false,
}

// BuildConfig holds some build related configuration.
type BuildConfig struct {
	UseResourceCacheWhen string // never, fallback, always. Default is fallback

	// When enabled, will collect and write a hugo_stats.json with some build
	// related aggregated data (e.g. CSS class names).
	WriteStats bool

	// Can be used to toggle off writing of the intellinsense /assets/jsconfig.js
	// file.
	NoJSConfigInAssets bool
}

func (b BuildConfig) UseResourceCache(err error) bool {
	if b.UseResourceCacheWhen == "never" {
		return false
	}

	if b.UseResourceCacheWhen == "fallback" {
		return err == herrors.ErrFeatureNotAvailable
	}

	return true
}

func DecodeBuildConfig(cfg Provider) BuildConfig {
	m := cfg.GetStringMap("build")
	b := DefaultBuild
	if m == nil {
		return b
	}

	err := mapstructure.WeakDecode(m, &b)
	if err != nil {
		return DefaultBuild
	}

	b.UseResourceCacheWhen = strings.ToLower(b.UseResourceCacheWhen)
	when := b.UseResourceCacheWhen
	if when != "never" && when != "always" && when != "fallback" {
		b.UseResourceCacheWhen = "fallback"
	}

	return b
}

// SitemapConfig configures the sitemap to be generated.
type SitemapConfig struct {
	// The page change frequency.
	ChangeFreq string
	// The priority of the page.
	Priority float64
	// The sitemap filename.
	Filename string
}

func DecodeSitemap(prototype SitemapConfig, input map[string]any) (SitemapConfig, error) {
	err := mapstructure.WeakDecode(input, &prototype)
	return prototype, err
}

// Config for the dev server.
type Server struct {
	Headers   []Headers
	Redirects []Redirect

	compiledHeaders   []glob.Glob
	compiledRedirects []glob.Glob
}

func (s *Server) CompileConfig() error {
	if s.compiledHeaders != nil {
		return nil
	}
	for _, h := range s.Headers {
		s.compiledHeaders = append(s.compiledHeaders, glob.MustCompile(h.For))
	}
	for _, r := range s.Redirects {
		s.compiledRedirects = append(s.compiledRedirects, glob.MustCompile(r.From))
	}
	return nil
}

func (s *Server) MatchHeaders(pattern string) []types.KeyValueStr {
	if s.compiledHeaders == nil {
		return nil
	}

	var matches []types.KeyValueStr

	for i, g := range s.compiledHeaders {
		if g.Match(pattern) {
			h := s.Headers[i]
			for k, v := range h.Values {
				matches = append(matches, types.KeyValueStr{Key: k, Value: cast.ToString(v)})
			}
		}
	}

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Key < matches[j].Key
	})

	return matches
}

func (s *Server) MatchRedirect(pattern string) Redirect {
	if s.compiledRedirects == nil {
		return Redirect{}
	}

	pattern = strings.TrimSuffix(pattern, "index.html")

	for i, g := range s.compiledRedirects {
		redir := s.Redirects[i]

		// No redirect to self.
		if redir.To == pattern {
			return Redirect{}
		}

		if g.Match(pattern) {
			return redir
		}
	}

	return Redirect{}
}

type Headers struct {
	For    string
	Values map[string]any
}

type Redirect struct {
	From string
	To   string

	// HTTP status code to use for the redirect.
	// A status code of 200 will trigger a URL rewrite.
	Status int

	// Forcode redirect, even if original request path exists.
	Force bool
}

func (r Redirect) IsZero() bool {
	return r.From == ""
}

func DecodeServer(cfg Provider) (Server, error) {
	s := &Server{}

	_ = mapstructure.WeakDecode(cfg.GetStringMap("server"), s)

	for i, redir := range s.Redirects {
		// Get it in line with the Hugo server for OK responses.
		// We currently treat the 404 as a special case, they are always "ugly", so keep them as is.
		if redir.Status != 404 {
			redir.To = strings.TrimSuffix(redir.To, "index.html")
			if !strings.HasPrefix(redir.To, "https") && !strings.HasSuffix(redir.To, "/") {
				// There are some tricky infinite loop situations when dealing
				// when the target does not have a trailing slash.
				// This can certainly be handled better, but not time for that now.
				return Server{}, fmt.Errorf("unsupported redirect to value %q in server config; currently this must be either a remote destination or a local folder, e.g. \"/blog/\" or \"/blog/index.html\"", redir.To)
			}
		}
		s.Redirects[i] = redir
	}

	if len(s.Redirects) == 0 {
		// Set up a default redirect for 404s.
		s.Redirects = []Redirect{
			{
				From:   "**",
				To:     "/404.html",
				Status: 404,
			},
		}

	}

	return *s, nil
}
