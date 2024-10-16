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

package page_test

import (
	"testing"

	"github.com/bep/logg"
	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/htesting"
	"github.com/gohugoio/hugo/hugolib"
)

func TestPermalinks(t *testing.T) {
	t.Parallel()

	files := `
-- layouts/_default/list.html --
List|{{ .Kind }}|{{ .RelPermalink }}|
-- layouts/_default/single.html --
Single|{{ .Kind }}|{{ .RelPermalink }}|
-- hugo.toml --
[taxonomies]
tag = "tags"
[permalinks.page]
withpageslug = '/pageslug/:slug/'
withallbutlastsection = '/:sections[:last]/:slug/'
[permalinks.section]
withfilefilename = '/sectionwithfilefilename/:filename/'
withfilefiletitle = '/sectionwithfilefiletitle/:title/'
withfileslug = '/sectionwithfileslug/:slug/'
nofileslug = '/sectionnofileslug/:slug/'
nofilefilename = '/sectionnofilefilename/:filename/'
nofiletitle1 = '/sectionnofiletitle1/:title/'
nofiletitle2 = '/sectionnofiletitle2/:sections[:last]/'
[permalinks.term]
tags = '/tagsslug/tag/:slug/'
[permalinks.taxonomy]
tags = '/tagsslug/:slug/'
-- content/withpageslug/p1.md --
---
slug: "p1slugvalue"
tags: ["mytag"]
---
-- content/withfilefilename/_index.md --
-- content/withfileslug/_index.md --
---
slug: "withfileslugvalue"
---
-- content/nofileslug/p1.md --
-- content/nofilefilename/p1.md --
-- content/nofiletitle1/p1.md --
-- content/nofiletitle2/asdf/p1.md --
-- content/withallbutlastsection/subsection/p1.md --
-- content/tags/_index.md --
---
slug: "tagsslug"
---
-- content/tags/mytag/_index.md --
---
slug: "mytagslug"
---


`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			LogLevel:    logg.LevelWarn,
		}).Build()

	t.Log(b.LogString())
	// No .File.TranslationBaseName on zero object etc. warnings.
	b.Assert(b.H.Log.LoggCount(logg.LevelWarn), qt.Equals, 0)
	b.AssertFileContent("public/pageslug/p1slugvalue/index.html", "Single|page|/pageslug/p1slugvalue/|")
	b.AssertFileContent("public/sectionwithfilefilename/index.html", "List|section|/sectionwithfilefilename/|")
	b.AssertFileContent("public/sectionwithfileslug/withfileslugvalue/index.html", "List|section|/sectionwithfileslug/withfileslugvalue/|")
	b.AssertFileContent("public/sectionnofilefilename/index.html", "List|section|/sectionnofilefilename/|")
	b.AssertFileContent("public/sectionnofileslug/nofileslugs/index.html", "List|section|/sectionnofileslug/nofileslugs/|")
	b.AssertFileContent("public/sectionnofiletitle1/nofiletitle1s/index.html", "List|section|/sectionnofiletitle1/nofiletitle1s/|")
	b.AssertFileContent("public/sectionnofiletitle2/index.html", "List|section|/sectionnofiletitle2/|")

	b.AssertFileContent("public/tagsslug/tag/mytagslug/index.html", "List|term|/tagsslug/tag/mytagslug/|")
	b.AssertFileContent("public/tagsslug/tagsslug/index.html", "List|taxonomy|/tagsslug/tagsslug/|")

	permalinksConf := b.H.Configs.Base.Permalinks
	b.Assert(permalinksConf, qt.DeepEquals, map[string]map[string]string{
		"page":     {"withallbutlastsection": "/:sections[:last]/:slug/", "withpageslug": "/pageslug/:slug/"},
		"section":  {"nofilefilename": "/sectionnofilefilename/:filename/", "nofileslug": "/sectionnofileslug/:slug/", "nofiletitle1": "/sectionnofiletitle1/:title/", "nofiletitle2": "/sectionnofiletitle2/:sections[:last]/", "withfilefilename": "/sectionwithfilefilename/:filename/", "withfilefiletitle": "/sectionwithfilefiletitle/:title/", "withfileslug": "/sectionwithfileslug/:slug/"},
		"taxonomy": {"tags": "/tagsslug/:slug/"},
		"term":     {"tags": "/tagsslug/tag/:slug/"},
	})
}

func TestPermalinksOldSetup(t *testing.T) {
	t.Parallel()

	files := `
-- layouts/_default/list.html --
List|{{ .Kind }}|{{ .RelPermalink }}|
-- layouts/_default/single.html --
Single|{{ .Kind }}|{{ .RelPermalink }}|
-- hugo.toml --
[permalinks]
withpageslug = '/pageslug/:slug/'
-- content/withpageslug/p1.md --
---
slug: "p1slugvalue"
---




`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			LogLevel:    logg.LevelWarn,
		}).Build()

	t.Log(b.LogString())
	// No .File.TranslationBaseName on zero object etc. warnings.
	b.Assert(b.H.Log.LoggCount(logg.LevelWarn), qt.Equals, 0)
	b.AssertFileContent("public/pageslug/p1slugvalue/index.html", "Single|page|/pageslug/p1slugvalue/|")

	permalinksConf := b.H.Configs.Base.Permalinks
	b.Assert(permalinksConf, qt.DeepEquals, map[string]map[string]string{
		"page":     {"withpageslug": "/pageslug/:slug/"},
		"section":  {},
		"taxonomy": {},
		"term":     {"withpageslug": "/pageslug/:slug/"},
	})
}

func TestPermalinksNestedSections(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
[permalinks.page]
books = '/libros/:sections[1:]/:filename'

[permalinks.section]
books = '/libros/:sections[1:]'
-- content/books/_index.md --
---
title: Books
---
-- content/books/fiction/_index.md --
---
title: Fiction
---
-- content/books/fiction/2023/_index.md --
---
title: 2023
---
-- content/books/fiction/2023/book1/index.md --
---
title: Book 1
---
-- layouts/_default/single.html --
Single.
-- layouts/_default/list.html --
List.
`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			LogLevel:    logg.LevelWarn,
		}).Build()

	t.Log(b.LogString())
	// No .File.TranslationBaseName on zero object etc. warnings.
	b.Assert(b.H.Log.LoggCount(logg.LevelWarn), qt.Equals, 0)

	b.AssertFileContent("public/libros/index.html", "List.")
	b.AssertFileContent("public/libros/fiction/index.html", "List.")
	b.AssertFileContent("public/libros/fiction/2023/book1/index.html", "Single.")
}

func TestPermalinksUrlCascade(t *testing.T) {
	t.Parallel()

	files := `
-- layouts/_default/list.html --
List|{{ .Kind }}|{{ .RelPermalink }}|
-- layouts/_default/single.html --
Single|{{ .Kind }}|{{ .RelPermalink }}|
-- hugo.toml --
-- content/cooking/delicious-recipes/_index.md --
---
url: /delicious-recipe/
cascade:
  url: /delicious-recipe/:slug/
---
-- content/cooking/delicious-recipes/example1.md --
---
title: Recipe 1
---
-- content/cooking/delicious-recipes/example2.md --
---
title: Recipe 2
slug: custom-recipe-2
---
`
	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			LogLevel:    logg.LevelWarn,
		}).Build()

	t.Log(b.LogString())
	b.Assert(b.H.Log.LoggCount(logg.LevelWarn), qt.Equals, 0)
	b.AssertFileContent("public/delicious-recipe/index.html", "List|section|/delicious-recipe/")
	b.AssertFileContent("public/delicious-recipe/recipe-1/index.html", "Single|page|/delicious-recipe/recipe-1/")
	b.AssertFileContent("public/delicious-recipe/custom-recipe-2/index.html", "Single|page|/delicious-recipe/custom-recipe-2/")
}

// Issue 12948.
func TestPermalinksWithEscapedColons(t *testing.T) {
	t.Parallel()

	if htesting.IsWindows() {
		t.Skip("Windows does not support colons in paths")
	}

	files := `
-- hugo.toml --
disableKinds = ['home','rss','section','sitemap','taxonomy','term']
[permalinks.page]
s2 = "/c\\:d/:slug/"
-- content/s1/p1.md --
---
title: p1
url: "/a\\:b/:slug/"
---
-- content/s2/p2.md --
---
title: p2
---
-- layouts/_default/single.html --
{{ .Title }}
`

	b := hugolib.Test(t, files)

	b.AssertFileExists("public/a:b/p1/index.html", true)

	// The above URL comes from the URL front matter field where everything is allowed.
	// We strip colons from paths constructed by Hugo (they are not supported on Windows).
	b.AssertFileExists("public/cd/p2/index.html", true)
}
