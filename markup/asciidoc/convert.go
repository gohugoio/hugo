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

// Package asciidoc converts Asciidoc to HTML using Asciidoc or Asciidoctor
// external binaries.
package asciidoc

import (
	"os/exec"
	"path/filepath"

	"github.com/gohugoio/hugo/markup/converter"
	"github.com/gohugoio/hugo/markup/internal"
)

// Provider is the package entry point.
var Provider converter.ProviderProvider = provider{}

type provider struct {
}

func (p provider) New(cfg converter.ProviderConfig) (converter.Provider, error) {
	return converter.NewProvider("asciidoc", func(ctx converter.DocumentContext) (converter.Converter, error) {
		config, _ := DecodeConfig(cfg.Cfg)
		return &asciidocConverter{
			ctx:    ctx,
			cfg:    cfg,
			config: config,
		}, nil
	}), nil
}

type asciidocConverter struct {
	ctx    converter.DocumentContext
	cfg    converter.ProviderConfig
	config Config
}

func (a *asciidocConverter) Convert(ctx converter.RenderContext) (converter.Result, error) {
	return converter.Bytes(a.getAsciidocContent(ctx.Src, a.ctx)), nil
}

// getAsciidocContent calls asciidoctor or asciidoc as an external helper
// to convert AsciiDoc content to HTML.
func (a *asciidocConverter) getAsciidocContent(src []byte, ctx converter.DocumentContext) []byte {
	var isAsciidoctor bool
	path := getAsciidoctorExecPath()
	if path == "" {
		path = getAsciidocExecPath()
		if path == "" {
			a.cfg.Logger.ERROR.Println("asciidoctor / asciidoc not found in $PATH: Please install.\n",
				"                 Leaving AsciiDoc content unrendered.")
			return src
		}
	} else {
		isAsciidoctor = true
	}

	var args []string
	if isAsciidoctor {
		args = a.getAsciidoctorArgs(ctx)
	} else {
		args = []string{"--no-header-footer", "--safe"}
	}
	if args[len(args)-1] != "-" {
		args = append(args, "-")
	}

	a.cfg.Logger.INFO.Println("Rendering", ctx.DocumentName, "with", path, "using asciidoc args", args, "...")

	return internal.ExternallyRenderContent(a.cfg, ctx, src, path, args)
}

func (a *asciidocConverter) getAsciidoctorArgs(ctx converter.DocumentContext) []string {
	args := a.config.Args
	currentContent := a.config.CurrentContent

	if currentContent {
		contentDir := filepath.Dir(ctx.FileName)
		destinationDir := a.cfg.Cfg.GetString("destination")
		outDir, err := filepath.Abs(filepath.Dir(filepath.Join(destinationDir, ctx.DocumentName)))
		if err != nil {
			a.cfg.Logger.ERROR.Println("asciidoctor outDir", err)
		}
		args = append(args, "--base-dir", contentDir, "-a", "outdir="+outDir)
	}

	return args
}

func getAsciidocExecPath() string {
	path, err := exec.LookPath("asciidoc")
	if err != nil {
		return ""
	}
	return path
}

func getAsciidoctorExecPath() string {
	path, err := exec.LookPath("asciidoctor")
	if err != nil {
		return ""
	}
	return path
}

// Supports returns whether Asciidoc or Asciidoctor is installed on this computer.
func Supports() bool {
	return (getAsciidoctorExecPath() != "" ||
		getAsciidocExecPath() != "")
}
