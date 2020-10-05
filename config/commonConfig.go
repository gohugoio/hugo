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
	"github.com/pkg/errors"

	"sort"
	"strings"
	"sync"

	"github.com/gohugoio/hugo/common/types"

	"github.com/gobwas/glob"
	"github.com/gohugoio/hugo/common/herrors"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
	jww "github.com/spf13/jwalterweatherman"
)

var DefaultBuild = Build{
	UseResourceCacheWhen: "fallback",
	WriteStats:           false,
}

// Build holds some build related condfiguration.
type Build struct {
	UseResourceCacheWhen string // never, fallback, always. Default is fallback

	// When enabled, will collect and write a hugo_stats.json with some build
	// related aggregated data (e.g. CSS class names).
	WriteStats bool

	// Can be used to toggle off writing of the intellinsense /assets/jsconfig.js
	// file.
	NoJSConfigInAssets bool
}

func (b Build) UseResourceCache(err error) bool {
	if b.UseResourceCacheWhen == "never" {
		return false
	}

	if b.UseResourceCacheWhen == "fallback" {
		return err == herrors.ErrFeatureNotAvailable
	}

	return true
}

func DecodeBuild(cfg Provider) Build {
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

// Sitemap configures the sitemap to be generated.
type Sitemap struct {
	ChangeFreq string
	Priority   float64
	Filename   string
}

func DecodeSitemap(prototype Sitemap, input map[string]interface{}) Sitemap {

	for key, value := range input {
		switch key {
		case "changefreq":
			prototype.ChangeFreq = cast.ToString(value)
		case "priority":
			prototype.Priority = cast.ToFloat64(value)
		case "filename":
			prototype.Filename = cast.ToString(value)
		default:
			jww.WARN.Printf("Unknown Sitemap field: %s\n", key)
		}
	}

	return prototype
}

// Config for the dev server.
type Server struct {
	Headers   []Headers
	Redirects []Redirect

	compiledInit      sync.Once
	compiledHeaders   []glob.Glob
	compiledRedirects []glob.Glob
}

func (s *Server) init() {

	s.compiledInit.Do(func() {
		for _, h := range s.Headers {
			s.compiledHeaders = append(s.compiledHeaders, glob.MustCompile(h.For))
		}
		for _, r := range s.Redirects {
			s.compiledRedirects = append(s.compiledRedirects, glob.MustCompile(r.From))
		}
	})
}

func (s *Server) MatchHeaders(pattern string) []types.KeyValueStr {
	s.init()

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
	s.init()

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
	Values map[string]interface{}
}

type Redirect struct {
	From   string
	To     string
	Status int
	Force  bool
}

func (r Redirect) IsZero() bool {
	return r.From == ""
}

func DecodeServer(cfg Provider) (*Server, error) {
	m := cfg.GetStringMap("server")
	s := &Server{}
	if m == nil {
		return s, nil
	}

	_ = mapstructure.WeakDecode(m, s)

	for i, redir := range s.Redirects {
		// Get it in line with the Hugo server.
		redir.To = strings.TrimSuffix(redir.To, "index.html")
		if !strings.HasPrefix(redir.To, "https") && !strings.HasSuffix(redir.To, "/") {
			// There are some tricky infinite loop situations when dealing
			// when the target does not have a trailing slash.
			// This can certainly be handled better, but not time for that now.
			return nil, errors.Errorf("unsupported redirect to value %q in server config; currently this must be either a remote destination or a local folder, e.g. \"/blog/\" or \"/blog/index.html\"", redir.To)
		}
		s.Redirects[i] = redir
	}

	return s, nil
}
