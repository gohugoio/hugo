// Copyright 2018 The Hugo Authors. All rights reserved.
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

	"github.com/stretchr/testify/require"
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
	assert := require.New(t)

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
		contentRoot := "content/main"

		switch lang {
		case "nn":
			contentRoot = "content/norsk"
		case "sv":
			contentRoot = "content/svensk"
		}
		return contentRoot + "/" + section
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
			slug := fmt.Sprintf("%s", base)
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

			contentRoot := contentRoot(lang)

			filename := filepath.Join(contentRoot, fmt.Sprintf("page%d%s.md", j, langID))
			contentFiles = append(contentFiles, filename, fmt.Sprintf(pageTemplate, slug, slug, j))
		}
	}

	// Put common translations in all of them
	for i, lang := range []string{"en", "nn", "sv"} {
		contentRoot := contentRoot(lang)

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
		contentRoot := contentRoot(lang)
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

	b := newTestSitesBuilder(t)
	b.WithWorkingDir("/my/project").WithConfigFile("toml", config).WithContent(contentFiles...).CreateSites()

	_ = os.Stdout
	//printFs(b.H.BaseFs.ContentFs, "/", os.Stdout)

	b.Build(BuildCfg{})

	assert.Equal(3, len(b.H.Sites))

	enSite := b.H.Sites[0]
	nnSite := b.H.Sites[1]
	svSite := b.H.Sites[2]

	//dumpPages(nnSite.RegularPages...)
	assert.Equal(12, len(nnSite.RegularPages))
	assert.Equal(13, len(enSite.RegularPages))

	assert.Equal(10, len(svSite.RegularPages))

	svP2, err := svSite.getPageNew(nil, "/sect/page2.md")
	assert.NoError(err)
	nnP2, err := nnSite.getPageNew(nil, "/sect/page2.md")
	assert.NoError(err)

	enP2, err := enSite.getPageNew(nil, "/sect/page2.md")
	assert.NoError(err)
	assert.Equal("en", enP2.Lang())
	assert.Equal("sv", svP2.Lang())
	assert.Equal("nn", nnP2.Lang())

	content, _ := nnP2.Content()
	assert.Contains(content, "SVP3-REF: https://example.org/sv/sect/p-sv-3/")
	assert.Contains(content, "SVP3-RELREF: /sv/sect/p-sv-3/")

	// Test RelRef with and without language indicator.
	nn3RefArgs := map[string]interface{}{
		"path": "/sect/page3.md",
		"lang": "nn",
	}
	nnP3RelRef, err := svP2.RelRef(
		nn3RefArgs,
	)
	assert.NoError(err)
	assert.Equal("/nn/sect/p-nn-3/", nnP3RelRef)
	nnP3Ref, err := svP2.Ref(
		nn3RefArgs,
	)
	assert.NoError(err)
	assert.Equal("https://example.org/nn/sect/p-nn-3/", nnP3Ref)

	for i, p := range enSite.RegularPages {
		j := i + 1
		msg := fmt.Sprintf("Test %d", j)
		assert.Equal("en", p.Lang(), msg)
		assert.Equal("sect", p.Section())
		if j < 9 {
			if j%4 == 0 {
				assert.Contains(p.Title(), fmt.Sprintf("p-sv-%d.en", i+1), msg)
			} else {
				assert.Contains(p.Title(), "p-en", msg)
			}
		}
	}

	// Check bundles
	bundleEn := enSite.RegularPages[len(enSite.RegularPages)-1]
	bundleNn := nnSite.RegularPages[len(nnSite.RegularPages)-1]
	bundleSv := svSite.RegularPages[len(svSite.RegularPages)-1]

	assert.Equal("/en/sect/mybundle/", bundleEn.RelPermalink())
	assert.Equal("/sv/sect/mybundle/", bundleSv.RelPermalink())

	assert.Equal(4, len(bundleEn.Resources))
	assert.Equal(4, len(bundleNn.Resources))
	assert.Equal(4, len(bundleSv.Resources))

	assert.Equal("/en/sect/mybundle/logo.png", bundleEn.Resources.GetMatch("logo*").RelPermalink())
	assert.Equal("/nn/sect/mybundle/logo.png", bundleNn.Resources.GetMatch("logo*").RelPermalink())
	assert.Equal("/sv/sect/mybundle/logo.png", bundleSv.Resources.GetMatch("logo*").RelPermalink())

	b.AssertFileContent("/my/project/public/sv/sect/mybundle/featured.png", "PNG Data for sv")
	b.AssertFileContent("/my/project/public/nn/sect/mybundle/featured.png", "PNG Data for nn")
	b.AssertFileContent("/my/project/public/en/sect/mybundle/featured.png", "PNG Data for en")
	b.AssertFileContent("/my/project/public/en/sect/mybundle/logo.png", "PNG Data")
	b.AssertFileContent("/my/project/public/sv/sect/mybundle/logo.png", "PNG Data")
	b.AssertFileContent("/my/project/public/nn/sect/mybundle/logo.png", "PNG Data")

	nnSect := nnSite.getPage(KindSection, "sect")
	assert.NotNil(nnSect)
	assert.Equal(12, len(nnSect.Pages))
	nnHome, _ := nnSite.Info.Home()
	assert.Equal("/nn/", nnHome.RelPermalink())

}
