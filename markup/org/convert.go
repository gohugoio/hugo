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

// Package org converts Emacs Org-Mode to HTML.
package org

import (
	"bytes"

	"github.com/gohugoio/hugo/identity"

	"github.com/gohugoio/hugo/markup/converter"
	"github.com/niklasfasching/go-org/org"
	"github.com/spf13/afero"
)

// Provider is the package entry point.
var Provider converter.ProviderProvider = provide{}

type provide struct {
}

func (p provide) New(cfg converter.ProviderConfig) (converter.Provider, error) {
	return converter.NewProvider("org", func(ctx converter.DocumentContext) (converter.Converter, error) {
		return &orgConverter{
			ctx: ctx,
			cfg: cfg,
		}, nil
	}), nil
}

type orgConverter struct {
	ctx converter.DocumentContext
	cfg converter.ProviderConfig
}

func (c *orgConverter) Convert(ctx converter.RenderContext) (converter.Result, error) {
	logger := c.cfg.Logger
	config := org.New()
	config.Log = logger.Warn()
	config.ReadFile = func(filename string) ([]byte, error) {
		return afero.ReadFile(c.cfg.ContentFs, filename)
	}
	writer := org.NewHTMLWriter()
	writer.HighlightCodeBlock = func(source, lang string, inline bool) string {
		highlightedSource, err := c.cfg.Highlight(source, lang, "")
		if err != nil {
			logger.Errorf("Could not highlight source as lang %s. Using raw source.", lang)
			return source
		}
		return highlightedSource
	}

	html, err := config.Parse(bytes.NewReader(ctx.Src), c.ctx.DocumentName).Write(writer)
	if err != nil {
		logger.Errorf("Could not render org: %s. Using unrendered content.", err)
		return converter.Bytes(ctx.Src), nil
	}
	return converter.Bytes([]byte(html)), nil
}

func (c *orgConverter) Supports(feature identity.Identity) bool {
	return false
}
