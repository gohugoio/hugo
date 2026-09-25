// Copyright 2025 The Hugo Authors. All rights reserved.
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

package hugolib

import (
	"fmt"
	"html/template"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestPermalink(t *testing.T) {
	t.Parallel()

	tests := []struct {
		file         string
		base         template.URL
		slug         string
		url          string
		uglyURLs     bool
		canonifyURLs bool
		expectedAbs  string
		expectedRel  string
	}{
		{"x/y/z/boofar.md", "", "", "", false, false, "/x/y/z/boofar/", "/x/y/z/boofar/"},
		{"x/y/z/boofar.md", "", "", "", false, false, "/x/y/z/boofar/", "/x/y/z/boofar/"},
		// Issue #1174
		{"x/y/z/boofar.md", "http://gopher.com/", "", "", false, true, "http://gopher.com/x/y/z/boofar/", "/x/y/z/boofar/"},
		{"x/y/z/boofar.md", "http://gopher.com/", "", "", true, true, "http://gopher.com/x/y/z/boofar.html", "/x/y/z/boofar.html"},
		{"x/y/z/boofar.md", "", "boofar", "", false, false, "/x/y/z/boofar/", "/x/y/z/boofar/"},
		{"x/y/z/boofar.md", "http://barnew/", "", "", false, false, "http://barnew/x/y/z/boofar/", "/x/y/z/boofar/"},
		{"x/y/z/boofar.md", "http://barnew/", "boofar", "", false, false, "http://barnew/x/y/z/boofar/", "/x/y/z/boofar/"},
		{"x/y/z/boofar.md", "", "", "", true, false, "/x/y/z/boofar.html", "/x/y/z/boofar.html"},
		{"x/y/z/boofar.md", "", "", "", true, false, "/x/y/z/boofar.html", "/x/y/z/boofar.html"},
		{"x/y/z/boofar.md", "", "boofar", "", true, false, "/x/y/z/boofar.html", "/x/y/z/boofar.html"},
		{"x/y/z/boofar.md", "http://barnew/", "", "", true, false, "http://barnew/x/y/z/boofar.html", "/x/y/z/boofar.html"},
		{"x/y/z/boofar.md", "http://barnew/", "boofar", "", true, false, "http://barnew/x/y/z/boofar.html", "/x/y/z/boofar.html"},
		{"x/y/z/boofar.md", "http://barnew/boo/", "booslug", "", true, false, "http://barnew/boo/x/y/z/booslug.html", "/boo/x/y/z/booslug.html"},
		{"x/y/z/boofar.md", "http://barnew/boo/", "booslug", "", false, true, "http://barnew/boo/x/y/z/booslug/", "/x/y/z/booslug/"},
		{"x/y/z/boofar.md", "http://barnew/boo/", "booslug", "", false, false, "http://barnew/boo/x/y/z/booslug/", "/boo/x/y/z/booslug/"},
		{"x/y/z/boofar.md", "http://barnew/boo/", "booslug", "", true, true, "http://barnew/boo/x/y/z/booslug.html", "/x/y/z/booslug.html"},
		{"x/y/z/boofar.md", "http://barnew/boo", "booslug", "", true, true, "http://barnew/boo/x/y/z/booslug.html", "/x/y/z/booslug.html"},
		// Issue #4666
		{"x/y/z/boo-makeindex.md", "http://barnew/boo", "", "", true, true, "http://barnew/boo/x/y/z/boo-makeindex.html", "/x/y/z/boo-makeindex.html"},

		// test URL overrides
		{"x/y/z/boofar.md", "", "", "/z/y/q/", false, false, "/z/y/q/", "/z/y/q/"},
		// test URL override with expands
		{"x/y/z/boofar.md", "", "test", "/z/:slug/", false, false, "/z/test/", "/z/test/"},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%s-%d", test.file, i), func(t *testing.T) {
			t.Parallel()
			c := qt.New(t)

			files := fmt.Sprintf(`
-- hugo.toml --
baseURL = %q
uglyURLs = %t
canonifyURLs = %t
-- content/%s --
---
title: Page
slug: %q
url: %q
output: ["HTML"]
---
`, test.base, test.uglyURLs, test.canonifyURLs, test.file, test.slug, test.url)

			b := Test(t, files)
			s := b.H.Sites[0]
			c.Assert(len(s.RegularPages()), qt.Equals, 1)
			p := s.RegularPages()[0]
			u := p.Permalink()

			expected := test.expectedAbs
			if u != expected {
				t.Fatalf("[%d] Expected abs url: %s, got: %s", i, expected, u)
			}

			u = p.RelPermalink()

			expected = test.expectedRel
			if u != expected {
				t.Errorf("[%d] Expected rel url: %s, got: %s", i, expected, u)
			}
		})
	}
}

func TestRelativeURLInFrontMatter(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com"
defaultContentLanguage = "en"
defaultContentLanguageInSubdir = false

[Languages]
[Languages.en]
weight = 10
contentDir = "content/en"
[Languages.nn]
weight = 20
contentDir = "content/nn"
-- layouts/single.html --
Single: {{ .Title }}|Hello|{{ .Lang }}|RelPermalink: {{ .RelPermalink }}|Permalink: {{ .Permalink }}|
-- layouts/list.html --
List Page 1|{{ .Title }}|Hello|{{ .Permalink }}|
-- content/en/blog/page1.md --
---
title: "A page"
url: "myblog/p1/"
---

Some content.
-- content/en/blog/page2.md --
---
title: "A page"
url: "../../../../../myblog/p2/"
---

Some content.
-- content/en/blog/page3.md --
---
title: "A page"
url: "../myblog/../myblog/p3/"
---

Some content.
-- content/en/blog/_index.md --
---
title: "A page"
url: "this-is-my-english-blog"
---

Some content.
-- content/nn/blog/page1.md --
---
title: "A page"
url: "myblog/p1/"
---

Some content.
-- content/nn/blog/_index.md --
---
title: "A page"
url: "this-is-my-blog"
---

Some content.
`
	b := Test(t, files)

	b.AssertFileContent("public/nn/myblog/p1/index.html", "Single: A page|Hello|nn|RelPermalink: /nn/myblog/p1/|")
	b.AssertFileContent("public/nn/this-is-my-blog/index.html", "List Page 1|A page|Hello|https://example.com/nn/this-is-my-blog/|")
	b.AssertFileContent("public/this-is-my-english-blog/index.html", "List Page 1|A page|Hello|https://example.com/this-is-my-english-blog/|")
	b.AssertFileContent("public/myblog/p1/index.html", "Single: A page|Hello|en|RelPermalink: /myblog/p1/|Permalink: https://example.com/myblog/p1/|")
	b.AssertFileContent("public/myblog/p2/index.html", "Single: A page|Hello|en|RelPermalink: /myblog/p2/|Permalink: https://example.com/myblog/p2/|")
	b.AssertFileContent("public/myblog/p3/index.html", "Single: A page|Hello|en|RelPermalink: /myblog/p3/|Permalink: https://example.com/myblog/p3/|")
}

// Issue 13869
func TestPermalinksSliceFromModule13869(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['home','rss','section','sitemap','taxonomy','term']
theme = 'foo'
[permalinks]
_merge = 'deep'
-- content/books/my-book.md --
---
title: My Book
---
-- layouts/page.html --
RelPermalink: {{ .RelPermalink }}
-- themes/foo/hugo.toml --
[[permalinks]]
pattern = '/shelf/:slug/'
[permalinks.target]
path = '/books/**'
`

	b := Test(t, files)

	b.AssertFileContent("public/shelf/my-book/index.html", "RelPermalink: /shelf/my-book/")
}

// See issue 14352.
func TestSlugInBranchPagesCascadesToDescendants(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
defaultContentLanguage = "en"
[languages.en]
weight = 1
[languages.es]
weight = 2
[taxonomies]
tag = "tags"
-- content/help/_index.en.md --
---
title: Help
---
-- content/help/_index.es.md --
---
title: Ayuda
slug: ayuda
---
-- content/help/how-to-frob/index.en.md --
---
title: How to frob
tags: [red]
---
-- content/help/how-to-frob/index.es.md --
---
title: Cómo frobear
slug: como-frobear
tags: [red]
---
-- content/help/advanced/_index.en.md --
---
title: Advanced
---
-- content/help/advanced/_index.es.md --
---
title: Avanzado
slug: avanzado
---
-- content/help/advanced/p1.en.md --
---
title: P1
---
-- content/help/advanced/p1.es.md --
---
title: P1
---
-- content/help/deep/dir/p2.es.md --
---
title: P2
---
-- content/tags/_index.es.md --
---
title: Etiquetas
slug: etiquetas
---
-- content/tags/red/_index.es.md --
---
title: Rojo
slug: rojo
---
-- layouts/all.html --
{{ .Kind }}|{{ .Path }}|{{ .RelPermalink }}|{{ range .Pages }}{{ .RelPermalink }}|{{ end }}
`

	b := Test(t, files)

	b.AssertFileContent("public/help/index.html", "section|/help|/help/|/help/advanced/|/help/how-to-frob/|")
	b.AssertFileContent("public/help/advanced/p1/index.html", "page|/help/advanced/p1|/help/advanced/p1/|")
	b.AssertFileContent("public/tags/red/index.html", "term|/tags/red|/tags/red/|/help/how-to-frob/|")

	b.AssertFileContent("public/es/ayuda/index.html", "section|/help|/es/ayuda/|/es/ayuda/avanzado/|/es/ayuda/como-frobear/|")
	b.AssertFileContent("public/es/ayuda/index.xml", "<link>/es/ayuda/</link>")
	b.AssertFileContent("public/es/ayuda/como-frobear/index.html", "page|/help/how-to-frob|/es/ayuda/como-frobear/|")
	b.AssertFileContent("public/es/ayuda/avanzado/index.html", "section|/help/advanced|/es/ayuda/avanzado/|/es/ayuda/avanzado/p1/|")
	b.AssertFileContent("public/es/ayuda/avanzado/p1/index.html", "page|/help/advanced/p1|/es/ayuda/avanzado/p1/|")
	b.AssertFileContent("public/es/ayuda/deep/dir/p2/index.html", "page|/help/deep/dir/p2|/es/ayuda/deep/dir/p2/|")
	b.AssertFileContent("public/es/etiquetas/index.html", "taxonomy|/tags|/es/etiquetas/|/es/etiquetas/rojo/|")
	b.AssertFileContent("public/es/etiquetas/rojo/index.html", "term|/tags/red|/es/etiquetas/rojo/|/es/ayuda/como-frobear/|")
}

func TestSlugInBranchPagesDisablePathToLowerUglyURLs(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disablePathToLower = true
uglyURLs = true
-- content/help/_index.md --
---
title: Help
slug: Ayuda
---
-- content/help/p1.md --
---
title: P1
---
-- layouts/all.html --
{{ .Kind }}|{{ .RelPermalink }}|{{ range .Pages }}{{ .RelPermalink }}|{{ end }}
`

	b := Test(t, files)

	b.AssertFileContent("public/Ayuda/index.html", "section|/Ayuda/index.html|/Ayuda/p1.html|")
	b.AssertFileContent("public/Ayuda/p1.html", "page|/Ayuda/p1.html|")
}

func BenchmarkTargetPathDir(b *testing.B) {
	files := `
-- content/help/_index.md --
---
title: Help
slug: Ayuda
---
-- content/help/advanced/_index.md --
---
title: Advanced
slug: Avanzado
---
-- content/help/advanced/p1.md --
---
title: P1			
---
`

	bb := Test(b, files)
	p, err := bb.H.Sites[0].GetPage("help/advanced/p1")
	bb.Assert(err, qt.IsNil)
	bb.Assert(p, qt.IsNotNil)
	b.ResetTimer()
	for b.Loop() {
		_ = targetPathDir(p.(*pageState))
	}
}
