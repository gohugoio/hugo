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

package asciidoc

import (
	"strings"
	"testing"

	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/markup/converter"

	"github.com/spf13/viper"

	qt "github.com/frankban/quicktest"

	"path/filepath"
)

func TestAsciidoctorDefaultArgs(t *testing.T) {
	c := qt.New(t)
	cfg, _ := config.FromConfigString("", "toml")
	p, err := Provider.New(converter.ProviderConfig{Logger: loggers.NewErrorLogger(), Cfg: cfg})
	c.Assert(err, qt.IsNil)

	conv, err := p.New(converter.DocumentContext{})
	c.Assert(err, qt.IsNil)

	ac := conv.(*asciidocConverter)
	c.Assert(ac, qt.Not(qt.IsNil))

	args := ac.getAsciidoctorArgs(converter.DocumentContext{})
	c.Assert(args, qt.Not(qt.IsNil))
	c.Assert(strings.Join(args, " "), qt.Equals, "--no-header-footer --safe --trace")
}

func TestAsciidoctorDiagramArgs(t *testing.T) {
	c := qt.New(t)
	cfg := viper.New()
	cfg.Set("asciidoctor.Args", []string{"--no-header-footer", "-r", "asciidoctor-html5s", "-b", "html5s", "-r", "asciidoctor-diagram"})
	p, err := Provider.New(converter.ProviderConfig{Logger: loggers.NewErrorLogger(), Cfg: cfg})
	c.Assert(err, qt.IsNil)

	conv, err := p.New(converter.DocumentContext{})
	c.Assert(err, qt.IsNil)

	ac := conv.(*asciidocConverter)
	c.Assert(ac, qt.Not(qt.IsNil))

	args := ac.getAsciidoctorArgs(converter.DocumentContext{})
	c.Assert(len(args), qt.Equals, 7)
	c.Assert(strings.Join(args, " "), qt.Equals, "--no-header-footer -r asciidoctor-html5s -b html5s -r asciidoctor-diagram")
}

func TestAsciidoctorCurrentContent(t *testing.T) {
	c := qt.New(t)
	cfg := viper.New()
	cfg.Set("asciidoctor.CurrentContent", true)
	p, err := Provider.New(converter.ProviderConfig{Logger: loggers.NewErrorLogger(), Cfg: cfg})
	c.Assert(err, qt.IsNil)

	ctx := converter.DocumentContext{FileName: "/tmp/hugo_asciidoc_ddd/docs/chapter2/index.adoc", DocumentName: "chapter2/index.adoc"}
	conv, err := p.New(ctx)
	c.Assert(err, qt.IsNil)

	ac := conv.(*asciidocConverter)
	c.Assert(ac, qt.Not(qt.IsNil))

	args := ac.getAsciidoctorArgs(ctx)
	c.Assert(len(args), qt.Equals, 7)
	c.Assert(args[0], qt.Equals, "--no-header-footer")
	c.Assert(args[1], qt.Equals, "--safe")
	c.Assert(args[2], qt.Equals, "--trace")
	c.Assert(args[3], qt.Equals, "--base-dir")
	c.Assert(filepath.ToSlash(args[4]), qt.Matches, "/tmp/hugo_asciidoc_ddd/docs/chapter2")
	c.Assert(args[5], qt.Equals, "-a")
	c.Assert(args[6], qt.Matches, `outdir=.*[/\\]{1,2}asciidoc[/\\]{1,2}chapter2`)
}

func TestAsciidoctorCurrentContentAndArgs(t *testing.T) {
	c := qt.New(t)
	cfg := viper.New()
	cfg.Set("asciidoctor.Args", []string{"--no-header-footer", "-r", "asciidoctor-html5s", "-b", "html5s", "-r", "asciidoctor-diagram"})
	cfg.Set("asciidoctor.CurrentContent", true)
	p, err := Provider.New(converter.ProviderConfig{Logger: loggers.NewErrorLogger(), Cfg: cfg})
	c.Assert(err, qt.IsNil)

	conv, err := p.New(converter.DocumentContext{})
	c.Assert(err, qt.IsNil)

	ac := conv.(*asciidocConverter)
	c.Assert(ac, qt.Not(qt.IsNil))

	args := ac.getAsciidoctorArgs(converter.DocumentContext{})
	c.Assert(len(args), qt.Equals, 11)
	c.Assert(args[0], qt.Equals, "--no-header-footer")
	c.Assert(args[1], qt.Equals, "-r")
	c.Assert(args[2], qt.Equals, "asciidoctor-html5s")
	c.Assert(args[3], qt.Equals, "-b")
	c.Assert(args[4], qt.Equals, "html5s")
	c.Assert(args[5], qt.Equals, "-r")
	c.Assert(args[6], qt.Equals, "asciidoctor-diagram")
	c.Assert(args[7], qt.Equals, "--base-dir")
	c.Assert(args[8], qt.Equals, ".")
	c.Assert(args[9], qt.Equals, "-a")
	c.Assert(args[10], qt.Contains, "outdir=")
}

func TestConvert(t *testing.T) {
	if !Supports() {
		t.Skip("asciidoc/asciidoctor not installed")
	}
	c := qt.New(t)
	p, err := Provider.New(converter.ProviderConfig{Logger: loggers.NewErrorLogger(), Cfg: viper.New()})
	c.Assert(err, qt.IsNil)

	conv, err := p.New(converter.DocumentContext{})
	c.Assert(err, qt.IsNil)

	b, err := conv.Convert(converter.RenderContext{Src: []byte("testContent")})
	c.Assert(err, qt.IsNil)
	c.Assert(string(b.Bytes()), qt.Equals, "<div class=\"paragraph\">\n<p>testContent</p>\n</div>\n")
}
