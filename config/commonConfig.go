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
}

// Build holds some build related condfiguration.
type Build struct {
	UseResourceCacheWhen string // never, fallback, always. Default is fallback
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
	Headers []Headers

	compiledInit sync.Once
	compiled     []glob.Glob
}

func (s *Server) Match(pattern string) []types.KeyValueStr {
	s.compiledInit.Do(func() {
		for _, h := range s.Headers {
			s.compiled = append(s.compiled, glob.MustCompile(h.For))
		}
	})

	if s.compiled == nil {
		return nil
	}

	var matches []types.KeyValueStr

	for i, g := range s.compiled {
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

type Headers struct {
	For    string
	Values map[string]interface{}
}

func DecodeServer(cfg Provider) *Server {
	m := cfg.GetStringMap("server")
	s := &Server{}
	if m == nil {
		return s
	}

	_ = mapstructure.WeakDecode(m, s)
	return s
}
