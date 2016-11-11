// Copyright 2016 The Hugo Authors. All rights reserved.
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
	"path/filepath"
	"testing"

	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

/*
	This file will test the "making everything a page" transition.

	See https://github.com/spf13/hugo/issues/2297

*/

func TestNodesAsPage(t *testing.T) {
	//jww.SetStdoutThreshold(jww.LevelDebug)
	jww.SetStdoutThreshold(jww.LevelFatal)

	/* Will have to decide what to name the node content files, but:

		Home page should have:
		Content, shortcode support
	   	Metadata (title, dates etc.)
		Params
	   	Taxonomies (categories, tags)

	*/

	testCommonResetState()

	writeLayoutsForNodeAsPageTests(t)
	writeNodePagesForNodeAsPageTests("", t)

	// Add some regular pages
	for i := 1; i <= 4; i++ {
		sect := "sect1"
		if i > 2 {
			sect = "sect2"
		}
		writeSource(t, filepath.Join("content", sect, fmt.Sprintf("regular%d.md", i)), fmt.Sprintf(`---
title: Page %02d
categories:  [
        "Hugo",
		"Web"
]
---
Content Page %02d
`, i, i))
	}

	viper.Set("paginate", 1)
	viper.Set("title", "Hugo Rocks")
	viper.Set("rssURI", "customrss.xml")

	s := newSiteDefaultLang()

	if err := buildAndRenderSite(s); err != nil {
		t.Fatalf("Failed to build site: %s", err)
	}

	assertFileContent(t, filepath.Join("public", "index.html"), false,
		"Index Title: Home Sweet Home!",
		"Home <strong>Content!</strong>",
		"# Pages: 9")

	assertFileContent(t, filepath.Join("public", "sect1", "regular1", "index.html"), false, "Single Title: Page 01", "Content Page 01")

	h := s.owner
	nodes := h.findAllPagesByNodeType(PageHome)
	require.Len(t, nodes, 1)

	home := nodes[0]

	require.True(t, home.IsHome())
	require.True(t, home.IsNode())
	require.False(t, home.IsPage())

	pages := h.findAllPagesByNodeType(PagePage)
	require.Len(t, pages, 4)

	first := pages[0]
	require.False(t, first.IsHome())
	require.False(t, first.IsNode())
	require.True(t, first.IsPage())

	first.Paginator()

	// Check Home paginator
	assertFileContent(t, filepath.Join("public", "page", "2", "index.html"), false,
		"Pag: Page 02")

	// Check Sections
	assertFileContent(t, filepath.Join("public", "sect1", "index.html"), false, "Section Title: Section", "Section1 <strong>Content!</strong>")
	assertFileContent(t, filepath.Join("public", "sect2", "index.html"), false, "Section Title: Section", "Section2 <strong>Content!</strong>")

	// Check Sections paginator
	assertFileContent(t, filepath.Join("public", "sect1", "page", "2", "index.html"), false,
		"Pag: Page 02")

	sections := h.findAllPagesByNodeType(PageSection)

	require.Len(t, sections, 2)

	// Check taxonomy lists
	assertFileContent(t, filepath.Join("public", "categories", "hugo", "index.html"), false,
		"Taxonomy Title: Taxonomy Hugo", "Taxonomy Hugo <strong>Content!</strong>")

	assertFileContent(t, filepath.Join("public", "categories", "web", "index.html"), false,
		"Taxonomy Title: Taxonomy Web", "Taxonomy Web <strong>Content!</strong>")

	// Check taxonomy list paginator
	assertFileContent(t, filepath.Join("public", "categories", "hugo", "page", "2", "index.html"), false,
		"Taxonomy Title: Taxonomy Hugo",
		"Pag: Page 02")

	// Check taxonomy terms
	assertFileContent(t, filepath.Join("public", "categories", "index.html"), false,
		"Taxonomy Terms Title: Taxonomy Term Categories", "Taxonomy Term Categories <strong>Content!</strong>", "k/v: hugo")

	// There are no pages to paginate over in the taxonomy terms.

	// RSS
	assertFileContent(t, filepath.Join("public", "customrss.xml"), false, "Recent content in Home Sweet Home! on Hugo Rocks", "<rss")
	assertFileContent(t, filepath.Join("public", "sect1", "customrss.xml"), false, "Recent content in Section1 on Hugo Rocks", "<rss")
	assertFileContent(t, filepath.Join("public", "sect2", "customrss.xml"), false, "Recent content in Section2 on Hugo Rocks", "<rss")
	assertFileContent(t, filepath.Join("public", "categories", "hugo", "customrss.xml"), false, "Recent content in Taxonomy Hugo on Hugo Rocks", "<rss")
	assertFileContent(t, filepath.Join("public", "categories", "web", "customrss.xml"), false, "Recent content in Taxonomy Web on Hugo Rocks", "<rss")

}

func TestNodesWithNoContentFile(t *testing.T) {
	//jww.SetStdoutThreshold(jww.LevelDebug)
	jww.SetStdoutThreshold(jww.LevelFatal)

	testCommonResetState()

	writeLayoutsForNodeAsPageTests(t)
	writeRegularPagesForNodeAsPageTests(t)

	viper.Set("paginate", 1)
	viper.Set("title", "Hugo Rocks!")
	viper.Set("rssURI", "customrss.xml")

	s := newSiteDefaultLang()

	if err := buildAndRenderSite(s); err != nil {
		t.Fatalf("Failed to build site: %s", err)
	}

	// Home page
	homePages := s.findIndexNodesByNodeType(PageHome)
	require.Len(t, homePages, 1)

	homePage := homePages[0]
	require.Len(t, homePage.Data["Pages"], 9)
	require.Len(t, homePage.Pages, 9) // Alias

	assertFileContent(t, filepath.Join("public", "index.html"), false,
		"Index Title: Hugo Rocks!")

	// Taxonomy list
	assertFileContent(t, filepath.Join("public", "categories", "hugo", "index.html"), false,
		"Taxonomy Title: Hugo")

	// Taxonomy terms
	assertFileContent(t, filepath.Join("public", "categories", "index.html"), false,
		"Taxonomy Terms Title: Categories")

	// Sections
	assertFileContent(t, filepath.Join("public", "sect1", "index.html"), false,
		"Section Title: Sect1s")
	assertFileContent(t, filepath.Join("public", "sect2", "index.html"), false,
		"Section Title: Sect2s")

	// RSS
	assertFileContent(t, filepath.Join("public", "customrss.xml"), false, "Recent content in Hugo Rocks! on Hugo Rocks!", "<rss")
	assertFileContent(t, filepath.Join("public", "sect1", "customrss.xml"), false, "Recent content in Sect1s on Hugo Rocks!", "<rss")
	assertFileContent(t, filepath.Join("public", "sect2", "customrss.xml"), false, "Recent content in Sect2s on Hugo Rocks!", "<rss")
	assertFileContent(t, filepath.Join("public", "categories", "hugo", "customrss.xml"), false, "Recent content in Hugo on Hugo Rocks!", "<rss")
	assertFileContent(t, filepath.Join("public", "categories", "web", "customrss.xml"), false, "Recent content in Web on Hugo Rocks!", "<rss")

}

func TestNodesAsPageMultilingual(t *testing.T) {

	testCommonResetState()

	writeLayoutsForNodeAsPageTests(t)

	writeSource(t, "config.toml",
		`
paginage = 1
title = "Hugo Multilingual Rocks!"
rssURI = "customrss.xml"

[languages]
[languages.nn]
languageName = "Nynorsk"
weight = 1
title = "Hugo på norsk"
defaultContentLanguage = "nn"

[languages.en]
languageName = "English"
weight = 2
title = "Hugo in English"
`)

	for _, lang := range []string{"nn", "en"} {
		writeRegularPagesForNodeAsPageTestsWithLang(t, lang)
	}

	// Only write node pages for the English side of the fence
	writeNodePagesForNodeAsPageTests("en", t)

	if err := LoadGlobalConfig("", "config.toml"); err != nil {
		t.Fatalf("Failed to load config: %s", err)
	}

	sites, err := NewHugoSitesFromConfiguration()

	if err != nil {
		t.Fatalf("Failed to create sites: %s", err)
	}

	if len(sites.Sites) != 2 {
		t.Fatalf("Got %d sites", len(sites.Sites))
	}

	err = sites.Build(BuildCfg{})

	if err != nil {
		t.Fatalf("Failed to build sites: %s", err)
	}

	// The en language has content pages

	// TODO(bep) np alias URL check

	assertFileContent(t, filepath.Join("public", "nn", "index.html"), true,
		"Index Title: Hugo på norsk")
	assertFileContent(t, filepath.Join("public", "en", "index.html"), true,
		"Index Title: Home Sweet Home!", "<strong>Content!</strong>")

	// Taxonomy list
	assertFileContent(t, filepath.Join("public", "nn", "categories", "hugo", "index.html"), true,
		"Taxonomy Title: Hugo")
	assertFileContent(t, filepath.Join("public", "en", "categories", "hugo", "index.html"), true,
		"Taxonomy Title: Taxonomy Hugo")

	// Taxonomy terms
	assertFileContent(t, filepath.Join("public", "nn", "categories", "index.html"), true,
		"Taxonomy Terms Title: Categories")
	assertFileContent(t, filepath.Join("public", "en", "categories", "index.html"), true,
		"Taxonomy Terms Title: Taxonomy Term Categories")

	// Sections
	assertFileContent(t, filepath.Join("public", "nn", "sect1", "index.html"), true,
		"Section Title: Sect1s")
	assertFileContent(t, filepath.Join("public", "nn", "sect2", "index.html"), true,
		"Section Title: Sect2s")
	assertFileContent(t, filepath.Join("public", "en", "sect1", "index.html"), true,
		"Section Title: Section1")
	assertFileContent(t, filepath.Join("public", "en", "sect2", "index.html"), true,
		"Section Title: Section2")

	// RSS
	assertFileContent(t, filepath.Join("public", "nn", "customrss.xml"), true, "Recent content in Hugo på norsk on Hugo på norsk", "<rss")
	assertFileContent(t, filepath.Join("public", "nn", "sect1", "customrss.xml"), true, "Recent content in Sect1s on Hugo på norsk", "<rss")
	assertFileContent(t, filepath.Join("public", "nn", "sect2", "customrss.xml"), true, "Recent content in Sect2s on Hugo på norsk", "<rss")
	assertFileContent(t, filepath.Join("public", "nn", "categories", "hugo", "customrss.xml"), true, "Recent content in Hugo on Hugo på norsk", "<rss")
	assertFileContent(t, filepath.Join("public", "nn", "categories", "web", "customrss.xml"), true, "Recent content in Web on Hugo på norsk", "<rss")

	assertFileContent(t, filepath.Join("public", "en", "customrss.xml"), true, "Recent content in Home Sweet Home! on Hugo in English", "<rss")
	assertFileContent(t, filepath.Join("public", "en", "sect1", "customrss.xml"), true, "Recent content in Section1 on Hugo in English", "<rss")
	assertFileContent(t, filepath.Join("public", "en", "sect2", "customrss.xml"), true, "Recent content in Section2 on Hugo in English", "<rss")
	assertFileContent(t, filepath.Join("public", "en", "categories", "hugo", "customrss.xml"), true, "Recent content in Taxonomy Hugo on Hugo in English", "<rss")
	assertFileContent(t, filepath.Join("public", "en", "categories", "web", "customrss.xml"), true, "Recent content in Taxonomy Web on Hugo in English", "<rss")

}

func TestNodesWithTaxonomies(t *testing.T) {
	testCommonResetState()

	writeLayoutsForNodeAsPageTests(t)
	writeRegularPagesForNodeAsPageTests(t)

	writeSource(t, filepath.Join("content", "_index.md"), `---
title: Home With Taxonomies
categories:  [
        "Hugo",
		"Home"
]
---
`)

	viper.Set("paginate", 1)
	viper.Set("title", "Hugo Rocks!")

	s := newSiteDefaultLang()

	if err := buildAndRenderSite(s); err != nil {
		t.Fatalf("Failed to build site: %s", err)
	}

	assertFileContent(t, filepath.Join("public", "categories", "hugo", "index.html"), true, "Taxonomy Title: Hugo", "# Pages: 5", "Pag: Home With Taxonomies")
	assertFileContent(t, filepath.Join("public", "categories", "home", "index.html"), true, "Taxonomy Title: Home", "# Pages: 1", "Pag: Home With Taxonomies")

}

func TestNodesWithMenu(t *testing.T) {
	testCommonResetState()

	writeLayoutsForNodeAsPageTests(t)
	writeRegularPagesForNodeAsPageTests(t)

	writeSource(t, filepath.Join("content", "_index.md"), `---
title: Home With Menu
menu:
  mymenu:
    name: "Go Home!"
---
`)

	viper.Set("paginate", 1)
	viper.Set("title", "Hugo Rocks!")

	s := newSiteDefaultLang()

	if err := buildAndRenderSite(s); err != nil {
		t.Fatalf("Failed to build site: %s", err)
	}

	assertFileContent(t, filepath.Join("public", "index.html"), true, "Home With Menu", "Menu Item: Go Home!")

}

func TestNodesWithAlias(t *testing.T) {
	testCommonResetState()

	writeLayoutsForNodeAsPageTests(t)
	writeRegularPagesForNodeAsPageTests(t)

	writeSource(t, filepath.Join("content", "_index.md"), `---
title: Home With Alias
aliases:
    - /my/new/home.html
---
`)

	viper.Set("paginate", 1)
	viper.Set("title", "Hugo Rocks!")

	s := newSiteDefaultLang()

	if err := buildAndRenderSite(s); err != nil {
		t.Fatalf("Failed to build site: %s", err)
	}

	assertFileContent(t, filepath.Join("public", "index.html"), true, "Home With Alias")
	assertFileContent(t, filepath.Join("public", "my", "new", "home.html"), true, "content=\"0; url=/")

}

func TestNodesWithSectionWithIndexPageOnly(t *testing.T) {
	testCommonResetState()

	writeLayoutsForNodeAsPageTests(t)

	writeSource(t, filepath.Join("content", "sect", "_index.md"), `---
title: MySection
---
My Section Content
`)

	viper.Set("paginate", 1)
	viper.Set("title", "Hugo Rocks!")

	s := newSiteDefaultLang()

	if err := buildAndRenderSite(s); err != nil {
		t.Fatalf("Failed to build site: %s", err)
	}

	assertFileContent(t, filepath.Join("public", "sect", "index.html"), true, "My Section")

}

func writeRegularPagesForNodeAsPageTests(t *testing.T) {
	writeRegularPagesForNodeAsPageTestsWithLang(t, "")
}

func writeRegularPagesForNodeAsPageTestsWithLang(t *testing.T, lang string) {
	var langStr string

	if lang != "" {
		langStr = lang + "."
	}

	for i := 1; i <= 4; i++ {
		sect := "sect1"
		if i > 2 {
			sect = "sect2"
		}
		writeSource(t, filepath.Join("content", sect, fmt.Sprintf("regular%d.%smd", i, langStr)), fmt.Sprintf(`---
title: Page %02d
categories:  [
        "Hugo",
		"Web"
]
---
Content Page %02d
`, i, i))
	}
}

func writeNodePagesForNodeAsPageTests(lang string, t *testing.T) {

	filename := "_index.md"

	if lang != "" {
		filename = fmt.Sprintf("_index.%s.md", lang)
	}

	writeSource(t, filepath.Join("content", filename), `---
title: Home Sweet Home!
---
Home **Content!**
`)

	writeSource(t, filepath.Join("content", "sect1", filename), `---
title: Section1
---
Section1 **Content!**
`)

	writeSource(t, filepath.Join("content", "sect2", filename), `---
title: Section2
---
Section2 **Content!**
`)

	writeSource(t, filepath.Join("content", "categories", "hugo", filename), `---
title: Taxonomy Hugo
---
Taxonomy Hugo **Content!**
`)

	writeSource(t, filepath.Join("content", "categories", "web", filename), `---
title: Taxonomy Web
---
Taxonomy Web **Content!**
`)

	writeSource(t, filepath.Join("content", "categories", filename), `---
title: Taxonomy Term Categories
---
Taxonomy Term Categories **Content!**
`)
}

func writeLayoutsForNodeAsPageTests(t *testing.T) {
	writeSource(t, filepath.Join("layouts", "index.html"), `
Index Title: {{ .Title }}
Index Content: {{ .Content }}
# Pages: {{ len .Data.Pages }}
{{ range .Paginator.Pages }}
	Pag: {{ .Title }}
{{ end }}
{{ with .Site.Menus.mymenu }}
{{ range . }}
Menu Item: {{ .Name }}
{{ end }}
{{ end }}
`)

	writeSource(t, filepath.Join("layouts", "_default", "single.html"), `
Single Title: {{ .Title }}
Single Content: {{ .Content }}
`)

	writeSource(t, filepath.Join("layouts", "_default", "section.html"), `
Section Title: {{ .Title }}
Section Content: {{ .Content }}
# Pages: {{ len .Data.Pages }}
{{ range .Paginator.Pages }}
	Pag: {{ .Title }}
{{ end }}
`)

	// Taxonomy lists
	writeSource(t, filepath.Join("layouts", "_default", "taxonomy.html"), `
Taxonomy Title: {{ .Title }}
Taxonomy Content: {{ .Content }}
# Pages: {{ len .Data.Pages }}
{{ range .Paginator.Pages }}
	Pag: {{ .Title }}
{{ end }}
`)

	// Taxonomy terms
	writeSource(t, filepath.Join("layouts", "_default", "terms.html"), `
Taxonomy Terms Title: {{ .Title }}
Taxonomy Terms Content: {{ .Content }}
{{ range $key, $value := .Data.Terms }}
	k/v: {{ $key }} / {{ printf "%=v" $value }}
{{ end }}
`)
}
