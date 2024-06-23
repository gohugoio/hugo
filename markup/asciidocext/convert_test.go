// Copyright 2024 The Hugo Authors. All rights reserved.
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

// Package asciidocext converts AsciiDoc to HTML using Asciidoctor
// external binary. The `asciidoc` module is reserved for a future golang
// implementation.

package asciidocext_test

import (
	"path/filepath"
	"testing"

	"github.com/gohugoio/hugo/common/collections"
	"github.com/gohugoio/hugo/common/hexec"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/config/security"
	"github.com/gohugoio/hugo/config/testconfig"
	"github.com/gohugoio/hugo/markup/asciidocext"
	"github.com/gohugoio/hugo/markup/asciidocext/internal"
	"github.com/gohugoio/hugo/markup/converter"
	"github.com/gohugoio/hugo/markup/markup_config"
	"github.com/spf13/afero"

	qt "github.com/frankban/quicktest"
)

func TestAsciidoctorDefaultArgs(t *testing.T) {
	c := qt.New(t)
	cfg := config.New()
	conf := testconfig.GetTestConfig(afero.NewMemMapFs(), cfg)

	p, err := asciidocext.Provider.New(
		converter.ProviderConfig{
			Conf:   conf,
			Logger: loggers.NewDefault(),
		},
	)
	c.Assert(err, qt.IsNil)

	conv, err := p.New(converter.DocumentContext{})
	c.Assert(err, qt.IsNil)

	ac := conv.(*internal.AsciidocConverter)
	c.Assert(ac, qt.Not(qt.IsNil))

	args := ac.ParseArgs(converter.DocumentContext{})
	expected := []string{"--no-header-footer"}
	c.Assert(args, qt.DeepEquals, expected)
}

func TestAsciidoctorNonDefaultArgs(t *testing.T) {
	c := qt.New(t)
	mconf := markup_config.Default
	mconf.AsciidocExt.Backend = "manpage"
	mconf.AsciidocExt.NoHeaderOrFooter = false
	mconf.AsciidocExt.SafeMode = "safe"
	mconf.AsciidocExt.SectionNumbers = true
	mconf.AsciidocExt.Verbose = true
	mconf.AsciidocExt.Trace = false
	mconf.AsciidocExt.FailureLevel = "warn"

	conf := testconfig.GetTestConfigSectionFromStruct("markup", mconf)

	p, err := asciidocext.Provider.New(
		converter.ProviderConfig{
			Conf:   conf,
			Logger: loggers.NewDefault(),
		},
	)
	c.Assert(err, qt.IsNil)

	conv, err := p.New(converter.DocumentContext{})
	c.Assert(err, qt.IsNil)

	ac := conv.(*internal.AsciidocConverter)
	c.Assert(ac, qt.Not(qt.IsNil))

	args := ac.ParseArgs(converter.DocumentContext{})
	expected := []string{"-b", "manpage", "--section-numbers", "--verbose", "--failure-level", "warn", "--safe-mode", "safe"}
	c.Assert(args, qt.DeepEquals, expected)
}

func TestAsciidoctorDisallowedArgs(t *testing.T) {
	c := qt.New(t)
	mconf := markup_config.Default
	mconf.AsciidocExt.Backend = "disallowed-backend"
	mconf.AsciidocExt.Extensions = []string{"./disallowed-extension"}
	mconf.AsciidocExt.Attributes = map[string]string{"outdir": "disallowed-attribute"}
	mconf.AsciidocExt.SafeMode = "disallowed-safemode"
	mconf.AsciidocExt.FailureLevel = "disallowed-failurelevel"

	conf := testconfig.GetTestConfigSectionFromStruct("markup", mconf)

	p, err := asciidocext.Provider.New(
		converter.ProviderConfig{
			Conf:   conf,
			Logger: loggers.NewDefault(),
		},
	)
	c.Assert(err, qt.IsNil)

	conv, err := p.New(converter.DocumentContext{})
	c.Assert(err, qt.IsNil)

	ac := conv.(*internal.AsciidocConverter)
	c.Assert(ac, qt.Not(qt.IsNil))

	args := ac.ParseArgs(converter.DocumentContext{})
	expected := []string{"--no-header-footer"}
	c.Assert(args, qt.DeepEquals, expected)
}

func TestAsciidoctorArbitraryExtension(t *testing.T) {
	c := qt.New(t)
	mconf := markup_config.Default
	mconf.AsciidocExt.Extensions = []string{"arbitrary-extension"}
	conf := testconfig.GetTestConfigSectionFromStruct("markup", mconf)
	p, err := asciidocext.Provider.New(
		converter.ProviderConfig{
			Conf:   conf,
			Logger: loggers.NewDefault(),
		},
	)
	c.Assert(err, qt.IsNil)

	conv, err := p.New(converter.DocumentContext{})
	c.Assert(err, qt.IsNil)

	ac := conv.(*internal.AsciidocConverter)
	c.Assert(ac, qt.Not(qt.IsNil))

	args := ac.ParseArgs(converter.DocumentContext{})
	expected := []string{"-r", "arbitrary-extension", "--no-header-footer"}
	c.Assert(args, qt.DeepEquals, expected)
}

func TestAsciidoctorDisallowedExtension(t *testing.T) {
	c := qt.New(t)

	for _, disallowedExtension := range []string{
		`foo-bar//`,
		`foo-bar\\ `,
		`../../foo-bar`,
		`/foo-bar`,
		`C:\foo-bar`,
		`foo-bar.rb`,
		`foo.bar`,
	} {
		mconf := markup_config.Default
		mconf.AsciidocExt.Extensions = []string{disallowedExtension}
		conf := testconfig.GetTestConfigSectionFromStruct("markup", mconf)
		p, err := asciidocext.Provider.New(
			converter.ProviderConfig{
				Conf:   conf,
				Logger: loggers.NewDefault(),
			},
		)
		c.Assert(err, qt.IsNil)

		conv, err := p.New(converter.DocumentContext{})
		c.Assert(err, qt.IsNil)

		ac := conv.(*internal.AsciidocConverter)
		c.Assert(ac, qt.Not(qt.IsNil))

		args := ac.ParseArgs(converter.DocumentContext{})
		expected := []string{"--no-header-footer"}
		c.Assert(args, qt.DeepEquals, expected)
	}
}

func TestAsciidoctorWorkingFolderCurrent(t *testing.T) {
	c := qt.New(t)
	cfg := config.FromTOMLConfigString(`
[markup]
[markup.asciidocext]
workingFolderCurrent = true
trace = false
`)

	conf := testconfig.GetTestConfig(afero.NewMemMapFs(), cfg)

	p, err := asciidocext.Provider.New(
		converter.ProviderConfig{
			Conf:   conf,
			Logger: loggers.NewDefault(),
		},
	)
	c.Assert(err, qt.IsNil)

	ctx := converter.DocumentContext{Filename: "/tmp/hugo_asciidoc_ddd/docs/chapter2/index.adoc", DocumentName: "chapter2/index.adoc"}
	conv, err := p.New(ctx)
	c.Assert(err, qt.IsNil)

	ac := conv.(*internal.AsciidocConverter)
	c.Assert(ac, qt.Not(qt.IsNil))

	args := ac.ParseArgs(ctx)
	c.Assert(len(args), qt.Equals, 5)
	c.Assert(args[0], qt.Equals, "--base-dir")
	c.Assert(filepath.ToSlash(args[1]), qt.Matches, "/tmp/hugo_asciidoc_ddd/docs/chapter2")
	c.Assert(args[2], qt.Equals, "-a")
	c.Assert(args[3], qt.Matches, `outdir=.*chapter2`)
	c.Assert(args[4], qt.Equals, "--no-header-footer")
}

func TestAsciidoctorWorkingFolderCurrentAndExtensions(t *testing.T) {
	c := qt.New(t)
	cfg := config.FromTOMLConfigString(`
[markup]
[markup.asciidocext]
backend = "html5s"
workingFolderCurrent = true
trace = false
noHeaderOrFooter = true
extensions = ["asciidoctor-html5s", "asciidoctor-diagram"]
`)
	conf := testconfig.GetTestConfig(afero.NewMemMapFs(), cfg)

	p, err := asciidocext.Provider.New(
		converter.ProviderConfig{
			Conf:   conf,
			Logger: loggers.NewDefault(),
		},
	)
	c.Assert(err, qt.IsNil)

	conv, err := p.New(converter.DocumentContext{})
	c.Assert(err, qt.IsNil)

	ac := conv.(*internal.AsciidocConverter)
	c.Assert(ac, qt.Not(qt.IsNil))

	args := ac.ParseArgs(converter.DocumentContext{})
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
	cfg := config.FromTOMLConfigString(`
[markup]
[markup.asciidocext]
trace = false
[markup.asciidocext.attributes]
my-base-url = "https://gohugo.io/"
my-attribute-name = "my value"
`)
	conf := testconfig.GetTestConfig(nil, cfg)
	p, err := asciidocext.Provider.New(
		converter.ProviderConfig{
			Conf:   conf,
			Logger: loggers.NewDefault(),
		},
	)
	c.Assert(err, qt.IsNil)

	conv, err := p.New(converter.DocumentContext{})
	c.Assert(err, qt.IsNil)

	ac := conv.(*internal.AsciidocConverter)
	c.Assert(ac, qt.Not(qt.IsNil))

	expectedValues := map[string]bool{
		"my-base-url=https://gohugo.io/": true,
		"my-attribute-name=my value":     true,
	}

	args := ac.ParseArgs(converter.DocumentContext{})
	c.Assert(len(args), qt.Equals, 5)
	c.Assert(args[0], qt.Equals, "-a")
	c.Assert(expectedValues[args[1]], qt.Equals, true)
	c.Assert(args[2], qt.Equals, "-a")
	c.Assert(expectedValues[args[3]], qt.Equals, true)
	c.Assert(args[4], qt.Equals, "--no-header-footer")
}

func getProvider(c *qt.C, mConfStr string) converter.Provider {
	confStr := `
[security]
[security.exec]
allow = ['asciidoctor']
`
	confStr += mConfStr

	cfg := config.FromTOMLConfigString(confStr)
	conf := testconfig.GetTestConfig(nil, cfg)
	securityConfig := conf.GetConfigSection("security").(security.Config)

	p, err := asciidocext.Provider.New(
		converter.ProviderConfig{
			Logger: loggers.NewDefault(),
			Conf:   conf,
			Exec:   hexec.New(securityConfig, ""),
		},
	)
	c.Assert(err, qt.IsNil)
	return p
}

func TestConvert(t *testing.T) {
	if !asciidocext.Supports() {
		t.Skip("asciidoctor not installed")
	}
	c := qt.New(t)

	p := getProvider(c, "")

	conv, err := p.New(converter.DocumentContext{})
	c.Assert(err, qt.IsNil)

	b, err := conv.Convert(converter.RenderContext{Src: []byte("testContent")})
	c.Assert(err, qt.IsNil)
	c.Assert(string(b.Bytes()), qt.Equals, "<div class=\"paragraph\">\n<p>testContent</p>\n</div>\n")
}

func TestTableOfContents(t *testing.T) {
	if !asciidocext.Supports() {
		t.Skip("asciidoctor not installed")
	}
	c := qt.New(t)
	p := getProvider(c, "")

	conv, err := p.New(converter.DocumentContext{})
	c.Assert(err, qt.IsNil)
	r, err := conv.Convert(converter.RenderContext{Src: []byte(`:toc: macro
:toclevels: 4
toc::[]

=== Introduction

== Section 1

=== Section 1.1

==== Section 1.1.1

=== Section 1.2

testContent

== Section 2
`)})
	c.Assert(err, qt.IsNil)
	toc, ok := r.(converter.TableOfContentsProvider)
	c.Assert(ok, qt.Equals, true)

	c.Assert(toc.TableOfContents().Identifiers, qt.DeepEquals, collections.SortedStringSlice{"_introduction", "_section_1", "_section_1_1", "_section_1_1_1", "_section_1_2", "_section_2"})
	// Although "Introduction" has a level 3 markup heading, AsciiDoc treats the first heading as level 2.
	c.Assert(toc.TableOfContents().HeadingsMap["_introduction"].Level, qt.Equals, 2)
	c.Assert(toc.TableOfContents().HeadingsMap["_section_1"].Level, qt.Equals, 2)
	c.Assert(toc.TableOfContents().HeadingsMap["_section_1_1"].Level, qt.Equals, 3)
	c.Assert(toc.TableOfContents().HeadingsMap["_section_1_1_1"].Level, qt.Equals, 4)
	c.Assert(toc.TableOfContents().HeadingsMap["_section_1_2"].Level, qt.Equals, 3)
	c.Assert(toc.TableOfContents().HeadingsMap["_section_2"].Level, qt.Equals, 2)
	c.Assert(string(r.Bytes()), qt.Not(qt.Contains), "<div id=\"toc\" class=\"toc\">")
}

func TestTableOfContentsWithCode(t *testing.T) {
	if !asciidocext.Supports() {
		t.Skip("asciidoctor not installed")
	}
	c := qt.New(t)
	p := getProvider(c, "")
	conv, err := p.New(converter.DocumentContext{})
	c.Assert(err, qt.IsNil)
	r, err := conv.Convert(converter.RenderContext{Src: []byte(`:toc: auto

== Some ` + "`code`" + ` in the title
`)})
	c.Assert(err, qt.IsNil)
	toc, ok := r.(converter.TableOfContentsProvider)
	c.Assert(ok, qt.Equals, true)
	c.Assert(toc.TableOfContents().HeadingsMap["_some_code_in_the_title"].Title, qt.Equals, "Some <code>code</code> in the title")
	c.Assert(string(r.Bytes()), qt.Not(qt.Contains), "<div id=\"toc\" class=\"toc\">")
}

func TestTableOfContentsPreserveTOC(t *testing.T) {
	if !asciidocext.Supports() {
		t.Skip("asciidoctor not installed")
	}
	c := qt.New(t)
	confStr := `
[markup]
[markup.asciidocExt]
preserveTOC = true
	`
	p := getProvider(c, confStr)

	conv, err := p.New(converter.DocumentContext{})
	c.Assert(err, qt.IsNil)
	r, err := conv.Convert(converter.RenderContext{Src: []byte(`:toc:
:idprefix:
:idseparator: -

== Some title
`)})
	c.Assert(err, qt.IsNil)
	toc, ok := r.(converter.TableOfContentsProvider)
	c.Assert(ok, qt.Equals, true)

	c.Assert(toc.TableOfContents().Identifiers, qt.DeepEquals, collections.SortedStringSlice{"some-title"})
	c.Assert(string(r.Bytes()), qt.Contains, "<div id=\"toc\" class=\"toc\">")
}
