// Copyright 2020 The Hugo Authors. All rights reserved.
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

// Package asciidocext converts Asciidoc to HTML using Asciidoc or Asciidoctor
// external binaries. The `asciidoc` module is reserved for a future golang
// implementation.

package asciidocext

import (
	"github.com/gohugoio/hugo/markup/markup_config"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/markup/converter"
	"github.com/spf13/viper"

	qt "github.com/frankban/quicktest"
)

func TestAsciidoctorDefaultArgs(t *testing.T) {
	c := qt.New(t)
	cfg := viper.New()
	mconf := markup_config.Default

	p, err := Provider.New(
		converter.ProviderConfig{
			Cfg:          cfg,
			MarkupConfig: mconf,
			Logger:       loggers.NewErrorLogger(),
		},
	)
	c.Assert(err, qt.IsNil)

	conv, err := p.New(converter.DocumentContext{})
	c.Assert(err, qt.IsNil)

	ac := conv.(*asciidocConverter)
	c.Assert(ac, qt.Not(qt.IsNil))

	args := ac.parseArgs(converter.DocumentContext{})
	c.Assert(args, qt.Not(qt.IsNil))
	c.Assert(strings.Join(args, " "), qt.Equals, "--no-header-footer")
}

func TestAsciidoctorDiagramArgs(t *testing.T) {
	c := qt.New(t)
	cfg := viper.New()
	mconf := markup_config.Default
	mconf.AsciidocExt.NoHeaderOrFooter = true
	mconf.AsciidocExt.Extensions = []string{"asciidoctor-html5s", "asciidoctor-diagram"}
	mconf.AsciidocExt.Backend = "html5s"

	p, err := Provider.New(
		converter.ProviderConfig{
			Cfg:          cfg,
			MarkupConfig: mconf,
			Logger:       loggers.NewErrorLogger(),
		},
	)
	c.Assert(err, qt.IsNil)

	conv, err := p.New(converter.DocumentContext{})
	c.Assert(err, qt.IsNil)

	ac := conv.(*asciidocConverter)
	c.Assert(ac, qt.Not(qt.IsNil))

	args := ac.parseArgs(converter.DocumentContext{})
	c.Assert(len(args), qt.Equals, 7)
	c.Assert(strings.Join(args, " "), qt.Equals, "-b html5s -r asciidoctor-html5s -r asciidoctor-diagram --no-header-footer")
}

func TestAsciidoctorWorkingFolderCurrent(t *testing.T) {
	c := qt.New(t)
	cfg := viper.New()
	mconf := markup_config.Default
	mconf.AsciidocExt.WorkingFolderCurrent = true
	p, err := Provider.New(
		converter.ProviderConfig{
			Cfg:          cfg,
			MarkupConfig: mconf,
			Logger:       loggers.NewErrorLogger(),
		},
	)
	c.Assert(err, qt.IsNil)

	ctx := converter.DocumentContext{Filename: "/tmp/hugo_asciidoc_ddd/docs/chapter2/index.adoc", DocumentName: "chapter2/index.adoc"}
	conv, err := p.New(ctx)
	c.Assert(err, qt.IsNil)

	ac := conv.(*asciidocConverter)
	c.Assert(ac, qt.Not(qt.IsNil))

	args := ac.parseArgs(ctx)
	c.Assert(len(args), qt.Equals, 5)
	c.Assert(args[0], qt.Equals, "--base-dir")
	c.Assert(filepath.ToSlash(args[1]), qt.Matches, "/tmp/hugo_asciidoc_ddd/docs/chapter2")
	c.Assert(args[2], qt.Equals, "-a")
	c.Assert(args[3], qt.Matches, `outdir=.*[/\\]{1,2}asciidocext[/\\]{1,2}chapter2`)
	c.Assert(args[4], qt.Equals, "--no-header-footer")
}

func TestAsciidoctorWorkingFolderCurrentAndExtensions(t *testing.T) {
	c := qt.New(t)
	cfg := viper.New()
	mconf := markup_config.Default
	mconf.AsciidocExt.NoHeaderOrFooter = true
	mconf.AsciidocExt.Extensions = []string{"asciidoctor-html5s", "asciidoctor-diagram"}
	mconf.AsciidocExt.Backend = "html5s"
	mconf.AsciidocExt.WorkingFolderCurrent = true
	p, err := Provider.New(
		converter.ProviderConfig{
			Cfg:          cfg,
			MarkupConfig: mconf,
			Logger:       loggers.NewErrorLogger(),
		},
	)
	c.Assert(err, qt.IsNil)

	conv, err := p.New(converter.DocumentContext{})
	c.Assert(err, qt.IsNil)

	ac := conv.(*asciidocConverter)
	c.Assert(ac, qt.Not(qt.IsNil))

	args := ac.parseArgs(converter.DocumentContext{})
	c.Assert(len(args), qt.Equals, 11)
	c.Assert(args[0], qt.Equals, "-b")
	c.Assert(args[1], qt.Equals, "html5s")
	c.Assert(args[2], qt.Equals, "-r")
	c.Assert(args[3], qt.Equals, "asciidoctor-html5s")
	c.Assert(args[4], qt.Equals, "-r")
	c.Assert(args[5], qt.Equals, "asciidoctor-diagram")
	c.Assert(args[6], qt.Equals, "--base-dir")
	c.Assert(args[7], qt.Equals, ".")
	c.Assert(args[8], qt.Equals, "-a")
	c.Assert(args[9], qt.Contains, "outdir=")
	c.Assert(args[10], qt.Equals, "--no-header-footer")
}

func TestAsciidoctorAttributes(t *testing.T) {
	c := qt.New(t)
	cfg := viper.New()
	mconf := markup_config.Default
	mconf.AsciidocExt.Attributes = map[string]string{"my-base-url": "https://gohugo.io/", "my-attribute-name": "my value"}
	p, err := Provider.New(
		converter.ProviderConfig{
			Cfg:          cfg,
			MarkupConfig: mconf,
			Logger:       loggers.NewErrorLogger(),
		},
	)
	c.Assert(err, qt.IsNil)

	conv, err := p.New(converter.DocumentContext{})
	c.Assert(err, qt.IsNil)

	ac := conv.(*asciidocConverter)
	c.Assert(ac, qt.Not(qt.IsNil))

	expectedValues := map[string]bool{
		"my-base-url=https://gohugo.io/": true,
		"my-attribute-name=my value":     true,
	}

	args := ac.parseArgs(converter.DocumentContext{})
	c.Assert(len(args), qt.Equals, 5)
	c.Assert(args[0], qt.Equals, "-a")
	c.Assert(expectedValues[args[1]], qt.Equals, true)
	c.Assert(args[2], qt.Equals, "-a")
	c.Assert(expectedValues[args[3]], qt.Equals, true)
	c.Assert(args[4], qt.Equals, "--no-header-footer")

}

func TestConvert(t *testing.T) {
	if !Supports() {
		t.Skip("asciidoc/asciidoctor not installed")
	}
	c := qt.New(t)

	mconf := markup_config.Default
	p, err := Provider.New(
		converter.ProviderConfig{
			MarkupConfig: mconf,
			Logger:       loggers.NewErrorLogger(),
		},
	)
	c.Assert(err, qt.IsNil)

	conv, err := p.New(converter.DocumentContext{})
	c.Assert(err, qt.IsNil)

	b, err := conv.Convert(converter.RenderContext{Src: []byte("testContent")})
	c.Assert(err, qt.IsNil)
	c.Assert(string(b.Bytes()), qt.Equals, "<div class=\"paragraph\">\n<p>testContent</p>\n</div>\n")
}
