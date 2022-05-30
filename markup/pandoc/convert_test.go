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

package pandoc

import (
	"testing"

	"github.com/gohugoio/hugo/common/hexec"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/config/security"

	"github.com/gohugoio/hugo/markup/converter"

	qt "github.com/frankban/quicktest"
)

func setupTestConverter(t *testing.T) (*qt.C, converter.Converter, converter.ProviderConfig) {
	if !Supports() {
		t.Skip("pandoc not installed")
	}
	c := qt.New(t)
	sc := security.DefaultConfig
	var err error
	sc.Exec.Allow, err = security.NewWhitelist("pandoc")
	c.Assert(err, qt.IsNil)
	cfg := converter.ProviderConfig{Exec: hexec.New(sc, "", loggers.NewDefault()), Logger: loggers.NewDefault()}
	p, err := Provider.New(cfg)
	c.Assert(err, qt.IsNil)
	conv, err := p.New(converter.DocumentContext{})
	c.Assert(err, qt.IsNil)
	return c, conv, cfg
}

func TestConvert(t *testing.T) {
	c, conv, _ := setupTestConverter(t)
	output, err := conv.Convert(converter.RenderContext{Src: []byte("testContent")})
	c.Assert(err, qt.IsNil)
	c.Assert(string(output.Bytes()), qt.Equals, "<p>testContent</p>\n")
}

func runCiteprocTest(t *testing.T, content string, expectContained []string, expectNotContained []string) {
	c, conv, cfg := setupTestConverter(t)
	if !supportsCitations(cfg) {
		t.Skip("pandoc does not support citations")
	}
	output, err := conv.Convert(converter.RenderContext{Src: []byte(content)})
	c.Assert(err, qt.IsNil)
	for _, expected := range expectContained {
		c.Assert(string(output.Bytes()), qt.Contains, expected)
	}
	for _, notExpected := range expectNotContained {
		c.Assert(string(output.Bytes()), qt.Not(qt.Contains), notExpected)
	}
}

func TestGetPandocSupportsCiteprocCallTwice(t *testing.T) {
	c, _, cfg := setupTestConverter(t)

	supports1, err1 := getPandocSupportsCiteproc(cfg)
	supports2, err2 := getPandocSupportsCiteproc(cfg)
	c.Assert(supports1, qt.Equals, supports2)
	c.Assert(err1, qt.IsNil)
	c.Assert(err2, qt.IsNil)
}

func TestCiteprocWithHugoMeta(t *testing.T) {
	content := `
---
title: Test
published: 2022-05-30
---
testContent
`
	expected := []string{"testContent"}
	unexpected := []string{"Doe", "Mustermann", "2022", "Treatise"}
	runCiteprocTest(t, content, expected, unexpected)
}

func TestCiteprocWithPandocMeta(t *testing.T) {
	content := `
---
---
---
...
testContent
`
	expected := []string{"testContent"}
	unexpected := []string{"Doe", "Mustermann", "2022", "Treatise"}
	runCiteprocTest(t, content, expected, unexpected)
}

func TestCiteprocWithBibliography(t *testing.T) {
	content := `
---
---
---
bibliography: testdata/bibliography.bib
...
testContent
`
	expected := []string{"testContent"}
	unexpected := []string{"Doe", "Mustermann", "2022", "Treatise"}
	runCiteprocTest(t, content, expected, unexpected)
}

func TestCiteprocWithExplicitCitation(t *testing.T) {
	content := `
---
---
---
bibliography: testdata/bibliography.bib
...
@Doe2022
`
	expected := []string{"Doe", "Mustermann", "2022", "Treatise"}
	runCiteprocTest(t, content, expected, []string{})
}

func TestCiteprocWithNocite(t *testing.T) {
	content := `
---
---
---
bibliography: testdata/bibliography.bib
nocite: |
  @*
...
`
	expected := []string{"Doe", "Mustermann", "2022", "Treatise"}
	runCiteprocTest(t, content, expected, []string{})
}
