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

package asciidocext_test

import (
	"strings"
	"testing"

	"github.com/gohugoio/hugo/common/hexec"
	"github.com/gohugoio/hugo/hugolib"
	"github.com/gohugoio/hugo/markup/asciidocext"
)

func TestAsciiDocConverterTemplates(t *testing.T) {
	if !asciidocext.Supports() {
		t.Skip("asciidoctor not installed")
	}
	missingGems := listMissingConverterTemplateGems()
	if len(missingGems) > 0 {
		t.Skip("these ruby gems, required to use AsciiDoc converter templates, are not installed:", strings.Join(missingGems, ", "))
	}

	files := `
-- hugo.toml --
disableKinds = ['page','section','rss','section','sitemap','taxonomy','term']
[markup.asciidocext]
templateDirectories = ['a','b']
templateEngine = 'handlebars'
[security.exec]
allow = ['asciidoctor']
-- content/_index.adoc --
---
title: home
---
https://gohugo.io[This is a link,title="Hugo rocks!"]

image:a.jpg[alt=A kitten,title=This is my kitten!]
-- layouts/index.html --
{{ .Content }}
-- a/inline_anchor.html.handlebars --
inline_anchor_html_handlebars
-- b/inline_image.html.handlebars --
inline_image_html_handlebars
`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			NeedsOsFS:   true,
		}).Build()

	b.AssertFileContent("public/index.html",
		`<p>inline_anchor_html_handlebars</p>`,
		`<p>inline_image_html_handlebars</p>`,
	)
}

func TestAsciiDocConverterTemplatesWithDisallowedFile(t *testing.T) {
	if !asciidocext.Supports() {
		t.Skip("asciidoctor not installed")
	}
	missingGems := listMissingConverterTemplateGems()
	if len(missingGems) > 0 {
		t.Skip("these ruby gems, required to use AsciiDoc converter templates, are not installed:", strings.Join(missingGems, ", "))
	}

	files := `
-- hugo.toml --
disableKinds = ['page','section','rss','section','sitemap','taxonomy','term']
[markup.asciidocext]
templateDirectories = ['a','b']
templateEngine = 'handlebars'
[security.exec]
allow = ['asciidoctor']
-- content/_index.adoc --
---
title: home
---
https://gohugo.io[This is a link,title="Hugo rocks!"]

image:a.jpg[alt=A kitten,title=This is my kitten!]
-- layouts/index.html --
{{ .Content }}
-- a/inline_anchor.html.handlebars --
inline_anchor_html_handlebars
-- b/inline_image.html.handlebars --
inline_image_html_handlebars
-- b/helpers.js --
I have the potential to execute arbitrary code; skip this directory!
`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			NeedsOsFS:   true,
		}).Build()

	b.AssertFileContent("public/index.html",
		`<p>inline_anchor_html_handlebars</p>`,
		`<img src="a.jpg" alt="A kitten" title="This is my kitten!"/>`,
	)
}

// converterTemplateGems is a map of the ruby gems required to test AsciiDoc
// converter templates. The key is the gem name, while the value is the
// executable name.
var converterTemplateGems = map[string]string{
	"tilt":            "tilt",
	"tilt-handlebars": "handlebars",
}

// listMissingConverterTemplateGems returns a slice of missing (not installed)
// ruby gems that are required to test AsciiDoc converter templates.
func listMissingConverterTemplateGems() []string {
	var gems []string
	for name, exec := range converterTemplateGems {
		if !hexec.InPath(exec) {
			gems = append(gems, name)
		}
	}
	return gems
}
