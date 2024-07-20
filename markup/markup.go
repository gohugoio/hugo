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

// Package markup contains the markup handling (e.g. Markdown).
package markup

import (
	"fmt"
	"strings"

	"github.com/gohugoio/hugo/markup/highlight"
	"github.com/gohugoio/hugo/media"

	"github.com/gohugoio/hugo/markup/markup_config"

	"github.com/gohugoio/hugo/markup/goldmark"

	"github.com/gohugoio/hugo/markup/org"

	"github.com/gohugoio/hugo/markup/asciidocext"
	"github.com/gohugoio/hugo/markup/converter"
	"github.com/gohugoio/hugo/markup/pandoc"
	"github.com/gohugoio/hugo/markup/rst"
)

func NewConverterProvider(cfg converter.ProviderConfig) (ConverterProvider, error) {
	converters := make(map[string]converter.Provider)

	mcfg := cfg.MarkupConfig()

	if cfg.Highlighter == nil {
		cfg.Highlighter = highlight.New(mcfg.Highlight)
	}

	defaultHandler := mcfg.DefaultMarkdownHandler
	var defaultFound bool

	add := func(p converter.ProviderProvider, subType string, aliases ...string) error {
		c, err := p.New(cfg)
		if err != nil {
			return err
		}

		name := c.Name()

		aliases = append(aliases, name)
		aliases = append(aliases, subType)

		if strings.EqualFold(name, defaultHandler) {
			aliases = append(aliases, "markdown")
			defaultFound = true
		}

		addConverter(converters, c, aliases...)
		return nil
	}

	contentTypes := cfg.Conf.ContentTypes().(media.ContentTypes)

	if err := add(goldmark.Provider, contentTypes.Markdown.SubType, contentTypes.Markdown.Suffixes()...); err != nil {
		return nil, err
	}
	if err := add(asciidocext.Provider, contentTypes.AsciiDoc.SubType, contentTypes.AsciiDoc.Suffixes()...); err != nil {
		return nil, err
	}
	if err := add(rst.Provider, contentTypes.ReStructuredText.SubType, contentTypes.ReStructuredText.Suffixes()...); err != nil {
		return nil, err
	}
	if err := add(pandoc.Provider, contentTypes.Pandoc.SubType, contentTypes.Pandoc.Suffixes()...); err != nil {
		return nil, err
	}
	if err := add(org.Provider, contentTypes.EmacsOrgMode.SubType, contentTypes.EmacsOrgMode.Suffixes()...); err != nil {
		return nil, err
	}

	if !defaultFound {
		msg := "markup: Configured defaultMarkdownHandler %q not found."
		if defaultHandler == "blackfriday" {
			msg += " Did you mean to use goldmark? Blackfriday was removed in Hugo v0.100.0."
		}
		return nil, fmt.Errorf(msg, defaultHandler)
	}

	return &converterRegistry{
		config:     cfg,
		converters: converters,
	}, nil
}

type ConverterProvider interface {
	Get(name string) converter.Provider
	IsGoldmark(name string) bool
	// Default() converter.Provider
	GetMarkupConfig() markup_config.Config
	GetHighlighter() highlight.Highlighter
}

type converterRegistry struct {
	// Maps name (md, markdown, goldmark etc.) to a converter provider.
	// Note that this is also used for aliasing, so the same converter
	// may be registered multiple times.
	// All names are lower case.
	converters map[string]converter.Provider

	config converter.ProviderConfig
}

func (r *converterRegistry) IsGoldmark(name string) bool {
	cp := r.Get(name)
	return cp != nil && cp.Name() == "goldmark"
}

func (r *converterRegistry) Get(name string) converter.Provider {
	return r.converters[strings.ToLower(name)]
}

func (r *converterRegistry) GetHighlighter() highlight.Highlighter {
	return r.config.Highlighter
}

func (r *converterRegistry) GetMarkupConfig() markup_config.Config {
	return r.config.MarkupConfig()
}

func addConverter(m map[string]converter.Provider, c converter.Provider, aliases ...string) {
	for _, alias := range aliases {
		m[alias] = c
	}
}

// ResolveMarkup returns the markup type.
func ResolveMarkup(s string) string {
	s = strings.ToLower(s)
	switch s {
	case "goldmark":
		return media.DefaultContentTypes.Markdown.SubType
	case "asciidocext":
		return media.DefaultContentTypes.AsciiDoc.SubType
	default:
		return s
	}
}
