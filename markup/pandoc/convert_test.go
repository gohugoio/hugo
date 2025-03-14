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

func TestGetPandocVersionCallTwice(t *testing.T) {
	c, _, cfg := setupTestConverter(t)

	version1, err1 := getPandocVersion(cfg)
	version2, err2 := getPandocVersion(cfg)
	c.Assert(version1, qt.Equals, version2)
	c.Assert(err1, qt.IsNil)
	c.Assert(err2, qt.IsNil)
}

func TestPandocVersionEquality(t *testing.T) {
	c := qt.New(t)
	v1 := pandocVersion{1, 0}
	v2 := pandocVersion{2, 0}
	v2_2 := pandocVersion{2, 2}
	v1_2 := pandocVersion{1, 2}
	v2_11 := pandocVersion{2, 11}
	v3_9 := pandocVersion{3, 9}
	v1_15 := pandocVersion{1, 15}

	c.Assert(v1.greaterThanOrEqual(v1), qt.IsTrue)

	c.Assert(v1.greaterThanOrEqual(v2), qt.IsFalse)
	c.Assert(v2.greaterThanOrEqual(v1), qt.IsTrue)

	c.Assert(v2.greaterThanOrEqual(v2_2), qt.IsFalse)
	c.Assert(v2_2.greaterThanOrEqual(v2), qt.IsTrue)

	c.Assert(v2_2.greaterThanOrEqual(v1_2), qt.IsTrue)
	c.Assert(v1_2.greaterThanOrEqual(v2_2), qt.IsFalse)

	c.Assert(v2_11.greaterThanOrEqual(v2_2), qt.IsTrue)
	c.Assert(v2_2.greaterThanOrEqual(v2_11), qt.IsFalse)

	c.Assert(v3_9.greaterThanOrEqual(v2_11), qt.IsTrue)
	c.Assert(v2_11.greaterThanOrEqual(v3_9), qt.IsFalse)

	c.Assert(v2_11.greaterThanOrEqual(v1_15), qt.IsTrue)
	c.Assert(v1_15.greaterThanOrEqual(v2_11), qt.IsFalse)
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
