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

// Package pandoc converts content to HTML using Pandoc as an external helper.
package pandoc

import (
	"github.com/cli/safeexec"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/markup/internal"

	"github.com/gohugoio/hugo/markup/converter"
)

// Provider is the package entry point.
var Provider converter.ProviderProvider = provider{}

type provider struct {
}

func (p provider) New(cfg converter.ProviderConfig) (converter.Provider, error) {
	return converter.NewProvider("pandoc", func(ctx converter.DocumentContext) (converter.Converter, error) {
		return &pandocConverter{
			ctx: ctx,
			cfg: cfg,
		}, nil
	}), nil

}

type pandocConverter struct {
	ctx converter.DocumentContext
	cfg converter.ProviderConfig
}

func (c *pandocConverter) Convert(ctx converter.RenderContext) (converter.Result, error) {
	return converter.Bytes(c.getPandocContent(ctx.Src, c.ctx)), nil
}

func (c *pandocConverter) Supports(feature identity.Identity) bool {
	return false
}

// getPandocContent calls pandoc as an external helper to convert pandoc markdown to HTML.
func (c *pandocConverter) getPandocContent(src []byte, ctx converter.DocumentContext) []byte {
	logger := c.cfg.Logger
	path := getPandocExecPath()
	if path == "" {
		logger.Println("pandoc not found in $PATH: Please install.\n",
			"                 Leaving pandoc content unrendered.")
		return src
	}
	args := []string{"--mathjax"}
	return internal.ExternallyRenderContent(c.cfg, ctx, src, path, args)
}

func getPandocExecPath() string {
	path, err := safeexec.LookPath("pandoc")
	if err != nil {
		return ""
	}

	return path
}

// Supports returns whether Pandoc is installed on this computer.
func Supports() bool {
	return getPandocExecPath() != ""
}
