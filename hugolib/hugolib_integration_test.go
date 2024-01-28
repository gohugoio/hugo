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

package hugolib_test

import (
	"strings"
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

// Issue 9073
func TestPageTranslationsMap(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
baseURL = 'https://example.org/'
title = 'Issue-9073'
defaultContentLanguageInSubdir = true

[taxonomies]
tag = 'tags'

[languages.en]
contentDir = 'content/en'
weight = 1
disableKinds = ['RSS','sitemap']

[languages.de]
contentDir = 'content/de'
weight = 2
disableKinds = ['home', 'page', 'section', 'taxonomy', 'term','RSS','sitemap']
-- content/de/posts/p1.md --
---
title: P1
tags: ['T1']
---
-- content/en/posts/p1.md --
---
title: P1
tags: ['T1']
---
-- layouts/_default/single.html --
<ul>{{ range .AllTranslations }}<li>{{ .Title }}-{{ .Lang }}</li>{{ end }}</ul>
-- layouts/_default/list.html --
<ul>{{ range .AllTranslations }}<li>{{ .Title }}-{{ .Lang }}</li>{{ end }}</ul>
	`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
		},
	)
	b.Build()

	// Kind home
	b.AssertFileContent("public/en/index.html",
		"<ul><li>Issue-9073-en</li></ul>",
	)
	// Kind section
	b.AssertFileContent("public/en/posts/index.html",
		"<ul><li>Posts-en</li></ul>",
	)
	// Kind page
	b.AssertFileContent("public/en/posts/p1/index.html",
		"<ul><li>P1-en</li></ul>",
	)
	// Kind taxonomy
	b.AssertFileContent("public/en/tags/index.html",
		"<ul><li>Tags-en</li></ul>",
	)
	// Kind term
	b.AssertFileContent("public/en/tags/t1/index.html",
		"<ul><li>T1-en</li></ul>",
	)
}

// Issue #11538
func TestRenderStringBadMarkupOpt(t *testing.T) {
	t.Parallel()

	files := `
-- layouts/index.html --
{{ $opts := dict "markup" "foo" }}
{{ "something" | .RenderString $opts }}
	`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
		},
	)

	_, err := b.BuildE()

	want := `no content renderer found for markup "foo"`
	if !strings.Contains(err.Error(), want) {
		t.Errorf("error msg must contain %q, error msg actually contains %q", want, err.Error())
	}
}

// Issue #11547
func TestTitleCaseStyleWithAutomaticSectionPages(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
titleCaseStyle = 'none'
-- content/books/book-1.md --
---
title: Book 1
tags: [fiction]
---
-- content/films/_index.md --
---
title: Films
---
-- layouts/index.html --
{{ (site.GetPage "/tags").Title }}
{{ (site.GetPage "/tags/fiction").Title }}
{{ (site.GetPage "/books").Title }}
{{ (site.GetPage "/films").Title }}
	`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
		},
	)
	b.Build()
	b.AssertFileContent("public/index.html", "tags\nfiction\nbooks\nFilms")
}
