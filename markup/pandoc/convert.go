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
	"errors"
	"strings"

	"github.com/gohugoio/hugo/common/hexec"
	"github.com/gohugoio/hugo/htesting"
	"github.com/mitchellh/mapstructure"

	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/markup/bibliography"
	"github.com/gohugoio/hugo/markup/converter"
	"github.com/gohugoio/hugo/markup/internal"
	"github.com/gohugoio/hugo/markup/pandoc/pandoc_config"

	"path"
)

type paramer interface {
	Param(interface{}) (interface{}, error)
}

type searchPaths struct {
	Paths []string
}

func (s *searchPaths) AsResourcePath() string {
	return strings.Join(s.Paths, ":")
}

// Provider is the package entry point.
var Provider converter.ProviderProvider = provider{}

type provider struct {
}

func (p provider) New(cfg converter.ProviderConfig) (converter.Provider, error) {
	return converter.NewProvider("pandoc", func(ctx converter.DocumentContext) (converter.Converter, error) {
		return &pandocConverter{
			docCtx: ctx,
			cfg:    cfg,
		}, nil
	}), nil
}

type pandocConverter struct {
	docCtx converter.DocumentContext
	cfg    converter.ProviderConfig
}

func (c *pandocConverter) Convert(ctx converter.RenderContext) (converter.ResultRender, error) {
	b, err := c.getPandocContent(ctx.Src)
	if err != nil {
		return nil, err
	}
	return converter.Bytes(b), nil
}

func (c *pandocConverter) Supports(feature identity.Identity) bool {
	return false
}

// getPandocContent calls pandoc as an external helper to convert pandoc markdown to HTML.
func (c *pandocConverter) getPandocContent(src []byte) ([]byte, error) {
	pandocPath, pandocFound := getPandocBinaryName()
	if !pandocFound {
		return nil, errors.New("pandoc not found in $PATH: Please install.")
	}

	var pandocConfig pandoc_config.Config = c.cfg.MarkupConfig().Pandoc
	var bibConfig bibliography.Config = c.cfg.MarkupConfig().Bibliography

	if pageParameters, ok := c.docCtx.Document.(paramer); ok {
		if bibParam, err := pageParameters.Param("bibliography"); err == nil {
			mapstructure.WeakDecode(bibParam, &bibConfig)
		}

		if pandocParam, err := pageParameters.Param("pandoc"); err == nil {
			mapstructure.WeakDecode(pandocParam, &pandocConfig)
		}
	}

	arguments := pandocConfig.AsPandocArguments()

	if bibConfig.Source != "" {
		arguments = append(arguments, "--citeproc", "--bibliography", bibConfig.Source)
		if bibConfig.CitationStyle != "" {
			arguments = append(arguments, "--csl", bibConfig.CitationStyle)
		}
	}

	resourcePath := strings.Join([]string{path.Dir(c.docCtx.Filename), "static", "."}, ":")
	arguments = append(arguments, "--resource-path", resourcePath)

	renderedContent, _ := internal.ExternallyRenderContent(c.cfg, c.docCtx, src, pandocPath, arguments)
	return renderedContent, nil
}

const pandocBinary = "pandoc"

func getPandocBinaryName() (string, bool) {
	return pandocBinary, hexec.InPath(pandocBinary)
}

// Supports returns whether Pandoc is installed on this computer.
func Supports() bool {
	_, hasBin := getPandocBinaryName()
	if htesting.SupportsAll() {
		if !hasBin {
			panic("pandoc not installed")
		}
		return true
	}
	return hasBin
}
