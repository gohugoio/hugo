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
	"bytes"
	"strconv"
	"strings"
	"sync"

	"github.com/gohugoio/hugo/common/hexec"
	"github.com/gohugoio/hugo/htesting"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/markup/converter"
	"github.com/gohugoio/hugo/markup/internal"
)

// Provider is the package entry point.
var Provider converter.ProviderProvider = provider{}

type provider struct{}

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

func (c *pandocConverter) Convert(ctx converter.RenderContext) (converter.ResultRender, error) {
	b, err := c.getPandocContent(ctx.Src, c.ctx)
	if err != nil {
		return nil, err
	}
	return converter.Bytes(b), nil
}

func (c *pandocConverter) Supports(feature identity.Identity) bool {
	return false
}

// getPandocContent calls pandoc as an external helper to convert pandoc markdown to HTML.
func (c *pandocConverter) getPandocContent(src []byte, ctx converter.DocumentContext) ([]byte, error) {
	logger := c.cfg.Logger
	binaryName := getPandocBinaryName()
	if binaryName == "" {
		logger.Println("pandoc not found in $PATH: Please install.\n",
			"                 Leaving pandoc content unrendered.")
		return src, nil
	}
	args := []string{"--mathjax"}
	if supportsCitations(c.cfg) {
		args = append(args[:], "--citeproc")
	}
	return internal.ExternallyRenderContent(c.cfg, ctx, src, binaryName, args)
}

const pandocBinary = "pandoc"

func getPandocBinaryName() string {
	if hexec.InPath(pandocBinary) {
		return pandocBinary
	}
	return ""
}

type pandocVersion struct {
	major, minor int64
}

func (left pandocVersion) greaterThanOrEqual(right pandocVersion) bool {
	return left.major > right.major || (left.major == right.major && left.minor >= right.minor)
}

var versionOnce sync.Once
var foundPandocVersion pandocVersion

// getPandocVersion parses the pandoc version output
func getPandocVersion(cfg converter.ProviderConfig) (pandocVersion, error) {
	var err error

	versionOnce.Do(func() {
		argsv := []any{"--version"}

		var out bytes.Buffer
		argsv = append(argsv, hexec.WithStdout(&out))

		cmd, err := cfg.Exec.New(pandocBinary, argsv...)
		if err != nil {
			cfg.Logger.Errorf("Could not call pandoc: %v", err)
			foundPandocVersion = pandocVersion{0, 0}
			return
		}

		err = cmd.Run()
		if err != nil {
			cfg.Logger.Errorf("%s --version: %v", pandocBinary, err)
			foundPandocVersion = pandocVersion{0, 0}
			return
		}

		outbytes := bytes.Replace(out.Bytes(), []byte("\r"), []byte(""), -1)
		output := strings.Split(string(outbytes), "\n")[0]
		// Split, e.g., "pandoc 2.5" into 2 and 5 and convert them to integers
		versionStrings := strings.Split(strings.Split(output, " ")[1], ".")
		majorVersion, err := strconv.ParseInt(versionStrings[0], 10, 64)
		if err != nil {
			println(err)
		}
		minorVersion, err := strconv.ParseInt(versionStrings[1], 10, 64)
		if err != nil {
			println(err)
		}
		foundPandocVersion = pandocVersion{majorVersion, minorVersion}
	})

	return foundPandocVersion, err
}

// SupportsCitations returns true for pandoc versions >= 2.11, which include citeproc
func supportsCitations(cfg converter.ProviderConfig) bool {
	if Supports() {
		foundPandocVersion, err := getPandocVersion(cfg)
		supportsCitations := foundPandocVersion.greaterThanOrEqual(pandocVersion{2, 11}) && err == nil
		return supportsCitations
	}
	return false
}

// Supports returns whether Pandoc is installed on this computer.
func Supports() bool {
	hasBin := getPandocBinaryName() != ""
	if htesting.SupportsAll() {
		if !hasBin {
			panic("pandoc not installed")
		}
		return true
	}
	return hasBin
}
