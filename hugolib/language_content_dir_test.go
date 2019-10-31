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

package hugolib

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cast"

	"github.com/gohugoio/hugo/resources/page"

	qt "github.com/frankban/quicktest"
)

/*

/en/p1.md
/nn/p1.md

.Readdir

- Name() => p1.en.md, p1.nn.md

.Stat(name)

.Open() --- real file name


*/

func TestLanguageContentRoot(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	config := `
baseURL = "https://example.org/"

defaultContentLanguage = "en"
defaultContentLanguageInSubdir = true

contentDir = "content/main"
workingDir = "/my/project"

[Languages]
[Languages.en]
weight = 10
title = "In English"
languageName = "English"

[Languages.nn]
weight = 20
title = "På Norsk"
languageName = "Norsk"
# This tells Hugo that all content in this directory is in the Norwegian language.
# It does not have to have the "my-page.nn.md" format. It can, but that is optional.
contentDir = "content/norsk"

[Languages.sv]
weight = 30
title = "På Svenska"
languageName = "Svensk"
contentDir = "content/svensk"
`

	pageTemplate := `
---
title: %s
slug: %s
weight: %d
---

Content.

SVP3-REF: {{< ref path="/sect/page3.md" lang="sv" >}}
SVP3-RELREF: {{< relref path="/sect/page3.md" lang="sv" >}}

`

	pageBundleTemplate := `
---
title: %s
weight: %d
---

Content.

`
	var contentFiles []string
	section := "sect"

	var contentRoot = func(lang string) string {
		switch lang {
		case "nn":
			return "content/norsk"
		case "sv":
			return "content/svensk"
		default:
			return "content/main"
		}

	}

	var contentSectionRoot = func(lang string) string {
		return contentRoot(lang) + "/" + section
	}

	for _, lang := range []string{"en", "nn", "sv"} {
		for j := 1; j <= 10; j++ {
			if (lang == "nn" || lang == "en") && j%4 == 0 {
				// Skip 4 and 8 for nn
				// We also skip it for en, but that is added to the Swedish directory below.
				continue
			}

			if lang == "sv" && j%5 == 0 {
				// Skip 5 and 10 for sv
				continue
			}

			base := fmt.Sprintf("p-%s-%d", lang, j)
			slug := base
			langID := ""

			if lang == "sv" && j%4 == 0 {
				// Put an English page in the Swedish content dir.
				langID = ".en"
			}

			if lang == "en" && j == 8 {
				// This should win over the sv variant above.
				langID = ".en"
			}

			slug += langID

			contentRoot := contentSectionRoot(lang)

			filename := filepath.Join(contentRoot, fmt.Sprintf("page%d%s.md", j, langID))
			contentFiles = append(contentFiles, filename, fmt.Sprintf(pageTemplate, slug, slug, j))
		}
	}

	// Put common translations in all of them
	for i, lang := range []string{"en", "nn", "sv"} {
		contentRoot := contentSectionRoot(lang)

		slug := fmt.Sprintf("common_%s", lang)

		filename := filepath.Join(contentRoot, "common.md")
		contentFiles = append(contentFiles, filename, fmt.Sprintf(pageTemplate, slug, slug, 100+i))

		for j, lang2 := range []string{"en", "nn", "sv"} {
			filename := filepath.Join(contentRoot, fmt.Sprintf("translated_all.%s.md", lang2))
			langSlug := slug + "_translated_all_" + lang2
			contentFiles = append(contentFiles, filename, fmt.Sprintf(pageTemplate, langSlug, langSlug, 200+i+j))
		}

		for j, lang2 := range []string{"sv", "nn"} {
			if lang == "en" {
				continue
			}
			filename := filepath.Join(contentRoot, fmt.Sprintf("translated_some.%s.md", lang2))
			langSlug := slug + "_translated_some_" + lang2
			contentFiles = append(contentFiles, filename, fmt.Sprintf(pageTemplate, langSlug, langSlug, 300+i+j))
		}
	}

	// Add a bundle with some images
	for i, lang := range []string{"en", "nn", "sv"} {
		contentRoot := contentSectionRoot(lang)
		slug := fmt.Sprintf("bundle_%s", lang)
		filename := filepath.Join(contentRoot, "mybundle", "index.md")
		contentFiles = append(contentFiles, filename, fmt.Sprintf(pageBundleTemplate, slug, 400+i))
		if lang == "en" {
			imageFilename := filepath.Join(contentRoot, "mybundle", "logo.png")
			contentFiles = append(contentFiles, imageFilename, "PNG Data")
		}
		imageFilename := filepath.Join(contentRoot, "mybundle", "featured.png")
		contentFiles = append(contentFiles, imageFilename, fmt.Sprintf("PNG Data for %s", lang))

		// Add some bundled pages
		contentFiles = append(contentFiles, filepath.Join(contentRoot, "mybundle", "p1.md"), fmt.Sprintf(pageBundleTemplate, slug, 401+i))
		contentFiles = append(contentFiles, filepath.Join(contentRoot, "mybundle", "sub", "p1.md"), fmt.Sprintf(pageBundleTemplate, slug, 402+i))

	}

	// Add some static files inside the content dir
	// https://github.com/gohugoio/hugo/issues/5759
	for _, lang := range []string{"en", "nn", "sv"} {
		contentRoot := contentRoot(lang)
		for i := 0; i < 2; i++ {
			filename := filepath.Join(contentRoot, "mystatic", fmt.Sprintf("file%d.yaml", i))
			contentFiles = append(contentFiles, filename, lang)
		}
	}

	b := newTestSitesBuilder(t)
	b.WithWorkingDir("/my/project").WithConfigFile("toml", config).WithContent(contentFiles...).CreateSites()

	_ = os.Stdout

	err := b.BuildE(BuildCfg{})

	//dumpPages(b.H.Sites[1].RegularPages()...)

	c.Assert(err, qt.IsNil)

	c.Assert(len(b.H.Sites), qt.Equals, 3)

	enSite := b.H.Sites[0]
	nnSite := b.H.Sites[1]
	svSite := b.H.Sites[2]

	b.AssertFileContent("/my/project/public/en/mystatic/file1.yaml", "en")
	b.AssertFileContent("/my/project/public/nn/mystatic/file1.yaml", "nn")

	//dumpPages(nnSite.RegularPages()...)

	c.Assert(len(nnSite.RegularPages()), qt.Equals, 12)
	c.Assert(len(enSite.RegularPages()), qt.Equals, 13)

	c.Assert(len(svSite.RegularPages()), qt.Equals, 10)

	svP2, err := svSite.getPageNew(nil, "/sect/page2.md")
	c.Assert(err, qt.IsNil)
	nnP2, err := nnSite.getPageNew(nil, "/sect/page2.md")
	c.Assert(err, qt.IsNil)

	enP2, err := enSite.getPageNew(nil, "/sect/page2.md")
	c.Assert(err, qt.IsNil)
	c.Assert(enP2.Language().Lang, qt.Equals, "en")
	c.Assert(svP2.Language().Lang, qt.Equals, "sv")
	c.Assert(nnP2.Language().Lang, qt.Equals, "nn")

	content, _ := nnP2.Content()
	contentStr := cast.ToString(content)
	c.Assert(contentStr, qt.Contains, "SVP3-REF: https://example.org/sv/sect/p-sv-3/")
	c.Assert(contentStr, qt.Contains, "SVP3-RELREF: /sv/sect/p-sv-3/")

	// Test RelRef with and without language indicator.
	nn3RefArgs := map[string]interface{}{
		"path": "/sect/page3.md",
		"lang": "nn",
	}
	nnP3RelRef, err := svP2.RelRef(
		nn3RefArgs,
	)
	c.Assert(err, qt.IsNil)
	c.Assert(nnP3RelRef, qt.Equals, "/nn/sect/p-nn-3/")
	nnP3Ref, err := svP2.Ref(
		nn3RefArgs,
	)
	c.Assert(err, qt.IsNil)
	c.Assert(nnP3Ref, qt.Equals, "https://example.org/nn/sect/p-nn-3/")

	for i, p := range enSite.RegularPages() {
		j := i + 1
		c.Assert(p.Language().Lang, qt.Equals, "en")
		c.Assert(p.Section(), qt.Equals, "sect")
		if j < 9 {
			if j%4 == 0 {
			} else {
				c.Assert(p.Title(), qt.Contains, "p-en")
			}
		}
	}

	for _, p := range nnSite.RegularPages() {
		c.Assert(p.Language().Lang, qt.Equals, "nn")
		c.Assert(p.Title(), qt.Contains, "nn")
	}

	for _, p := range svSite.RegularPages() {
		c.Assert(p.Language().Lang, qt.Equals, "sv")
		c.Assert(p.Title(), qt.Contains, "sv")
	}

	// Check bundles
	bundleEn := enSite.RegularPages()[len(enSite.RegularPages())-1]
	bundleNn := nnSite.RegularPages()[len(nnSite.RegularPages())-1]
	bundleSv := svSite.RegularPages()[len(svSite.RegularPages())-1]

	c.Assert(bundleEn.RelPermalink(), qt.Equals, "/en/sect/mybundle/")
	c.Assert(bundleSv.RelPermalink(), qt.Equals, "/sv/sect/mybundle/")

	c.Assert(len(bundleNn.Resources()), qt.Equals, 4)
	c.Assert(len(bundleSv.Resources()), qt.Equals, 4)
	c.Assert(len(bundleEn.Resources()), qt.Equals, 4)

	b.AssertFileContent("/my/project/public/en/sect/mybundle/index.html", "image/png: /en/sect/mybundle/logo.png")
	b.AssertFileContent("/my/project/public/nn/sect/mybundle/index.html", "image/png: /nn/sect/mybundle/logo.png")
	b.AssertFileContent("/my/project/public/sv/sect/mybundle/index.html", "image/png: /sv/sect/mybundle/logo.png")

	b.AssertFileContent("/my/project/public/sv/sect/mybundle/featured.png", "PNG Data for sv")
	b.AssertFileContent("/my/project/public/nn/sect/mybundle/featured.png", "PNG Data for nn")
	b.AssertFileContent("/my/project/public/en/sect/mybundle/featured.png", "PNG Data for en")
	b.AssertFileContent("/my/project/public/en/sect/mybundle/logo.png", "PNG Data")
	b.AssertFileContent("/my/project/public/sv/sect/mybundle/logo.png", "PNG Data")
	b.AssertFileContent("/my/project/public/nn/sect/mybundle/logo.png", "PNG Data")

	nnSect := nnSite.getPage(page.KindSection, "sect")
	c.Assert(nnSect, qt.Not(qt.IsNil))
	c.Assert(len(nnSect.Pages()), qt.Equals, 12)
	nnHome, _ := nnSite.Info.Home()
	c.Assert(nnHome.RelPermalink(), qt.Equals, "/nn/")

}

// https://github.com/gohugoio/hugo/issues/6463
func TestLanguageRootSectionsMismatch(t *testing.T) {
	t.Parallel()

	config := `
baseURL: "https://example.org/"
languageCode: "en-us"
title: "My New Hugo Site"
theme: "mytheme"

contentDir: "content/en"

languages:
    en:
        weight: 1
        languageName: "English"
        contentDir: content/en
    es:
        weight: 2
        languageName: "Español"
        contentDir: content/es
    fr:
        weight: 4
        languageName: "Française"
        contentDir: content/fr

        
`
	createPage := func(title string) string {
		return fmt.Sprintf(`---
title: %q
---

`, title)
	}

	b := newTestSitesBuilder(t)
	b.WithConfigFile("yaml", config)

	b.WithSourceFile("themes/mytheme/layouts/index.html", `MYTHEME`)
	b.WithTemplates("index.html", `
Lang: {{ .Lang }}
{{ range .Site.RegularPages }}
Page: {{ .RelPermalink }}|{{ .Title -}}
{{ end }}

`)
	b.WithSourceFile("static/hello.txt", `hello`)
	b.WithContent("en/_index.md", createPage("en home"))
	b.WithContent("es/_index.md", createPage("es home"))
	b.WithContent("fr/_index.md", createPage("fr home"))

	for i := 1; i < 3; i++ {
		b.WithContent(fmt.Sprintf("en/event/page%d.md", i), createPage(fmt.Sprintf("ev-en%d", i)))
		b.WithContent(fmt.Sprintf("es/event/page%d.md", i), createPage(fmt.Sprintf("ev-es%d", i)))
		b.WithContent(fmt.Sprintf("fr/event/page%d.md", i), createPage(fmt.Sprintf("ev-fr%d", i)))
		b.WithContent(fmt.Sprintf("en/blog/page%d.md", i), createPage(fmt.Sprintf("blog-en%d", i)))
		b.WithContent(fmt.Sprintf("es/blog/page%d.md", i), createPage(fmt.Sprintf("blog-es%d", i)))
		b.WithContent(fmt.Sprintf("fr/other/page%d.md", i), createPage(fmt.Sprintf("other-fr%d", i)))
	}

	b.Build(BuildCfg{})

	b.AssertFileContent("public/index.html", `
Lang: en
Page: /blog/page1/|blog-en1
Page: /blog/page2/|blog-en2
Page: /event/page1/|ev-en1
Page: /event/page2/|ev-en2
`)

	b.AssertFileContent("public/es/index.html", `
Lang: es
Page: /es/blog/page1/|blog-es1
Page: /es/blog/page2/|blog-es2
Page: /es/event/page1/|ev-es1
Page: /es/event/page2/|ev-es2
`)
	b.AssertFileContent("public/fr/index.html", `
Lang: fr
Page: /fr/event/page1/|ev-fr1
Page: /fr/event/page2/|ev-fr2
Page: /fr/other/page1/|other-fr1
Page: /fr/other/page2/|other-fr2`)

}
