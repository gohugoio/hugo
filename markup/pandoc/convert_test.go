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
	p, err := Provider.New(converter.ProviderConfig{Exec: hexec.New(sc), Logger: loggers.NewDefault()})
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

func runCiteprocTest(t *testing.T, content string, expected string) {
	c, conv, cfg := setupTestConverter(t)
	if !supportsCitations(cfg) {
		t.Skip("pandoc does not support citations")
	}
	output, err := conv.Convert(converter.RenderContext{Src: []byte(content)})
	c.Assert(err, qt.IsNil)
	c.Assert(string(output.Bytes()), qt.Equals, expected)
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
	v3 := pandocVersion{2, 2}
	v4 := pandocVersion{1, 2}
	v5 := pandocVersion{2, 11}

	// 1 >= 1 -> true
	c.Assert(v1.greaterThanOrEqual(v1), qt.IsTrue)

	// 1 >= 2 -> false, 2 >= 1 -> tru
	c.Assert(v1.greaterThanOrEqual(v2), qt.IsFalse)
	c.Assert(v2.greaterThanOrEqual(v1), qt.IsTrue)

	// 2.0 >= 2.2 -> false, 2.2 >= 2.0 -> true
	c.Assert(v2.greaterThanOrEqual(v3), qt.IsFalse)
	c.Assert(v3.greaterThanOrEqual(v2), qt.IsTrue)

	// 2.2 >= 1.2 -> true, 1.2 >= 2.2 -> false
	c.Assert(v3.greaterThanOrEqual(v4), qt.IsTrue)
	c.Assert(v4.greaterThanOrEqual(v3), qt.IsFalse)

	// 2.11 >= 2.2 -> true, 2.2 >= 2.11 -> false
	c.Assert(v5.greaterThanOrEqual(v3), qt.IsTrue)
	c.Assert(v3.greaterThanOrEqual(v5), qt.IsFalse)
}

func TestCiteprocWithHugoMeta(t *testing.T) {
	content := `
---
title: Test
published: 2022-05-30
---
testContent
`
	expected := "<p>testContent</p>\n"
	runCiteprocTest(t, content, expected)
}

func TestCiteprocWithPandocMeta(t *testing.T) {
	content := `
---
---
---
...
testContent
`
	expected := "<p>testContent</p>\n"
	runCiteprocTest(t, content, expected)
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
	expected := "<p>testContent</p>\n"
	runCiteprocTest(t, content, expected)
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
	expected := `<p><span class="citation" data-cites="Doe2022">Doe and Mustermann
(2022)</span></p>
<div id="refs" class="references csl-bib-body hanging-indent"
role="doc-bibliography">
<div id="ref-Doe2022" class="csl-entry" role="doc-biblioentry">
Doe, Jane, and Max Mustermann. 2022. <span>“A Treatise on Hugo
Tests.”</span> <em>Hugo Websites</em>.
</div>
</div>
`
	runCiteprocTest(t, content, expected)
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
	expected := `<div id="refs" class="references csl-bib-body hanging-indent"
role="doc-bibliography">
<div id="ref-Doe2022" class="csl-entry" role="doc-biblioentry">
Doe, Jane, and Max Mustermann. 2022. <span>“A Treatise on Hugo
Tests.”</span> <em>Hugo Websites</em>.
</div>
</div>
`
	runCiteprocTest(t, content, expected)
}
