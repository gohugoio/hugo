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
	"time"

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

	writeRegularPagesForNodeAsPageTests(t)

	viper.Set("paginate", 1)
	viper.Set("title", "Hugo Rocks")
	viper.Set("rssURI", "customrss.xml")

	s := newSiteDefaultLang()

	if err := buildAndRenderSite(s); err != nil {
		t.Fatalf("Failed to build site: %s", err)
	}

	// date order: home, sect1, sect2, cat/hugo, cat/web, categories

	assertFileContent(t, filepath.Join("public", "index.html"), false,
		"Index Title: Home Sweet Home!",
		"Home <strong>Content!</strong>",
		"# Pages: 4",
		"Date: 2009-01-02",
		"Lastmod: 2009-01-03",
		"GetPage: Section1 ",
	)

	assertFileContent(t, filepath.Join("public", "sect1", "regular1", "index.html"), false, "Single Title: Page 01", "Content Page 01")

	h := s.owner
	nodes := h.findAllPagesByKindNotIn(KindPage)
	require.Len(t, nodes, 6)

	home := nodes[5] // oldest

	require.True(t, home.IsHome())
	require.True(t, home.IsNode())
	require.False(t, home.IsPage())
	require.True(t, home.Path() != "")

	section2 := nodes[3]
	require.Equal(t, "Section2", section2.Title)

	pages := h.findAllPagesByKind(KindPage)
	require.Len(t, pages, 4)

	first := pages[0]

	require.False(t, first.IsHome())
	require.False(t, first.IsNode())
	require.True(t, first.IsPage())

	// Check Home paginator
	assertFileContent(t, filepath.Join("public", "page", "2", "index.html"), false,
		"Pag: Page 02")

	// Check Sections
	assertFileContent(t, filepath.Join("public", "sect1", "index.html"), false,
		"Section Title: Section", "Section1 <strong>Content!</strong>",
		"Date: 2009-01-04",
		"Lastmod: 2009-01-05",
	)

	assertFileContent(t, filepath.Join("public", "sect2", "index.html"), false,
		"Section Title: Section", "Section2 <strong>Content!</strong>",
		"Date: 2009-01-06",
		"Lastmod: 2009-01-07",
	)

	// Check Sections paginator
	assertFileContent(t, filepath.Join("public", "sect1", "page", "2", "index.html"), false,
		"Pag: Page 02")

	sections := h.findAllPagesByKind(KindSection)

	require.Len(t, sections, 2)

	// Check taxonomy lists
	assertFileContent(t, filepath.Join("public", "categories", "hugo", "index.html"), false,
		"Taxonomy Title: Taxonomy Hugo", "Taxonomy Hugo <strong>Content!</strong>",
		"Date: 2009-01-08",
		"Lastmod: 2009-01-09",
	)

	assertFileContent(t, filepath.Join("public", "categories", "web", "index.html"), false,
		"Taxonomy Title: Taxonomy Web",
		"Taxonomy Web <strong>Content!</strong>",
		"Date: 2009-01-10",
		"Lastmod: 2009-01-11",
	)

	// Check taxonomy list paginator
	assertFileContent(t, filepath.Join("public", "categories", "hugo", "page", "2", "index.html"), false,
		"Taxonomy Title: Taxonomy Hugo",
		"Pag: Page 02")

	// Check taxonomy terms
	assertFileContent(t, filepath.Join("public", "categories", "index.html"), false,
		"Taxonomy Terms Title: Taxonomy Term Categories", "Taxonomy Term Categories <strong>Content!</strong>", "k/v: hugo",
		"Date: 2009-01-12",
		"Lastmod: 2009-01-13",
	)

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
	homePages := s.findPagesByKind(KindHome)
	require.Len(t, homePages, 1)

	homePage := homePages[0]
	require.Len(t, homePage.Data["Pages"], 4)
	require.Len(t, homePage.Pages, 4)
	require.True(t, homePage.Path() == "")

	assertFileContent(t, filepath.Join("public", "index.html"), false,
		"Index Title: Hugo Rocks!",
		"Date: 2010-06-12",
		"Lastmod: 2010-06-13",
	)

	// Taxonomy list
	assertFileContent(t, filepath.Join("public", "categories", "hugo", "index.html"), false,
		"Taxonomy Title: Hugo",
		"Date: 2010-06-12",
		"Lastmod: 2010-06-13",
	)

	// Taxonomy terms
	assertFileContent(t, filepath.Join("public", "categories", "index.html"), false,
		"Taxonomy Terms Title: Categories",
	)

	// Sections
	assertFileContent(t, filepath.Join("public", "sect1", "index.html"), false,
		"Section Title: Sect1s",
		"Date: 2010-06-12",
		"Lastmod: 2010-06-13",
	)

	assertFileContent(t, filepath.Join("public", "sect2", "index.html"), false,
		"Section Title: Sect2s",
		"Date: 2008-07-06",
		"Lastmod: 2008-07-09",
	)

	// RSS
	assertFileContent(t, filepath.Join("public", "customrss.xml"), false, "Hugo Rocks!", "<rss")
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
defaultContentLanguage = "nn"
defaultContentLanguageInSubdir = true


[languages]
[languages.nn]
languageName = "Nynorsk"
weight = 1
title = "Hugo på norsk"

[languages.en]
languageName = "English"
weight = 2
title = "Hugo in English"

[languages.de]
languageName = "Deutsch"
weight = 3
title = "Deutsche Hugo"
`)

	for _, lang := range []string{"nn", "en"} {
		writeRegularPagesForNodeAsPageTestsWithLang(t, lang)
	}

	// Only write node pages for the English and Deutsch
	writeNodePagesForNodeAsPageTests("en", t)
	writeNodePagesForNodeAsPageTests("de", t)

	if err := LoadGlobalConfig("", "config.toml"); err != nil {
		t.Fatalf("Failed to load config: %s", err)
	}

	sites, err := NewHugoSitesFromConfiguration()

	if err != nil {
		t.Fatalf("Failed to create sites: %s", err)
	}

	if len(sites.Sites) != 3 {
		t.Fatalf("Got %d sites", len(sites.Sites))
	}

	err = sites.Build(BuildCfg{})

	if err != nil {
		t.Fatalf("Failed to build sites: %s", err)
	}

	// The en and de language have content pages
	enHome := sites.Sites[1].getPage("home")
	require.NotNil(t, enHome)
	require.Equal(t, "en", enHome.Language().Lang)
	require.Contains(t, enHome.Content, "l-en")

	deHome := sites.Sites[2].getPage("home")
	require.NotNil(t, deHome)
	require.Equal(t, "de", deHome.Language().Lang)
	require.Contains(t, deHome.Content, "l-de")

	require.Len(t, deHome.Translations(), 2, deHome.Translations()[0].Language().Lang)
	require.Equal(t, "en", deHome.Translations()[1].Language().Lang)
	require.Equal(t, "nn", deHome.Translations()[0].Language().Lang)

	enSect := sites.Sites[1].getPage("section", "sect1")
	require.NotNil(t, enSect)
	require.Equal(t, "en", enSect.Language().Lang)
	require.Len(t, enSect.Translations(), 2, enSect.Translations()[0].Language().Lang)
	require.Equal(t, "de", enSect.Translations()[1].Language().Lang)
	require.Equal(t, "nn", enSect.Translations()[0].Language().Lang)

	assertFileContent(t, filepath.Join("public", "nn", "index.html"), true,
		"Index Title: Hugo på norsk")
	assertFileContent(t, filepath.Join("public", "en", "index.html"), true,
		"Index Title: Home Sweet Home!", "<strong>Content!</strong>")
	assertFileContent(t, filepath.Join("public", "de", "index.html"), true,
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
	assertFileContent(t, filepath.Join("public", "nn", "customrss.xml"), true, "Hugo på norsk", "<rss")
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

	assertFileContent(t, filepath.Join("public", "categories", "hugo", "index.html"), true, "Taxonomy Title: Hugo", "# Pages: 5")
	assertFileContent(t, filepath.Join("public", "categories", "home", "index.html"), true, "Taxonomy Title: Home", "# Pages: 1")

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
	viper.Set("baseURL", "http://base/")
	viper.Set("title", "Hugo Rocks!")

	s := newSiteDefaultLang()

	if err := buildAndRenderSite(s); err != nil {
		t.Fatalf("Failed to build site: %s", err)
	}

	assertFileContent(t, filepath.Join("public", "index.html"), true, "Home With Alias")
	assertFileContent(t, filepath.Join("public", "my", "new", "home.html"), true, "content=\"0; url=http://base/")

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

func TestNodesWithURLs(t *testing.T) {
	testCommonResetState()

	writeLayoutsForNodeAsPageTests(t)

	writeRegularPagesForNodeAsPageTests(t)

	writeSource(t, filepath.Join("content", "sect", "_index.md"), `---
title: MySection
url: foo.html
---
My Section Content
`)

	viper.Set("paginate", 1)
	viper.Set("title", "Hugo Rocks!")
	viper.Set("baseURL", "http://bep.is/base/")

	s := newSiteDefaultLang()

	if err := buildAndRenderSite(s); err != nil {
		t.Fatalf("Failed to build site: %s", err)
	}

	assertFileContent(t, filepath.Join("public", "sect", "index.html"), true, "My Section")

	p := s.RegularPages[0]

	require.Equal(t, "/base/sect1/regular1/", p.URL())

	// Section with front matter and url set (which should not be used)
	sect := s.getPage(KindSection, "sect")
	require.Equal(t, "/base/sect/", sect.URL())
	require.Equal(t, "http://bep.is/base/sect/", sect.Permalink())
	require.Equal(t, "/base/sect/", sect.RelPermalink())

	// Home page without front matter
	require.Equal(t, "/base/", s.getPage(KindHome).URL())

}

func writeRegularPagesForNodeAsPageTests(t *testing.T) {
	writeRegularPagesForNodeAsPageTestsWithLang(t, "")
}

func writeRegularPagesForNodeAsPageTestsWithLang(t *testing.T, lang string) {
	var langStr string

	if lang != "" {
		langStr = lang + "."
	}

	format := "2006-01-02"

	date, _ := time.Parse(format, "2010-06-15")

	for i := 1; i <= 4; i++ {
		sect := "sect1"
		if i > 2 {
			sect = "sect2"

			date, _ = time.Parse(format, "2008-07-15") // Nodes are placed in 2009

		}
		date = date.Add(-24 * time.Duration(i) * time.Hour)
		writeSource(t, filepath.Join("content", sect, fmt.Sprintf("regular%d.%smd", i, langStr)), fmt.Sprintf(`---
title: Page %02d
lastMod : %q
date : %q
categories:  [
        "Hugo",
		"Web"
]
---
Content Page %02d
`, i, date.Add(time.Duration(i)*-24*time.Hour).Format(time.RFC822), date.Add(time.Duration(i)*-2*24*time.Hour).Format(time.RFC822), i))
	}
}

func writeNodePagesForNodeAsPageTests(lang string, t *testing.T) {

	filename := "_index.md"

	if lang != "" {
		filename = fmt.Sprintf("_index.%s.md", lang)
	}

	format := "2006-01-02"

	date, _ := time.Parse(format, "2009-01-01")

	writeSource(t, filepath.Join("content", filename), fmt.Sprintf(`---
title: Home Sweet Home!
date : %q
lastMod : %q
---
l-%s Home **Content!**
`, date.Add(1*24*time.Hour).Format(time.RFC822), date.Add(2*24*time.Hour).Format(time.RFC822), lang))

	writeSource(t, filepath.Join("content", "sect1", filename), fmt.Sprintf(`---
title: Section1
date : %q
lastMod : %q
---
Section1 **Content!**
`, date.Add(3*24*time.Hour).Format(time.RFC822), date.Add(4*24*time.Hour).Format(time.RFC822)))
	writeSource(t, filepath.Join("content", "sect2", filename), fmt.Sprintf(`---
title: Section2
date : %q
lastMod : %q
---
Section2 **Content!**
`, date.Add(5*24*time.Hour).Format(time.RFC822), date.Add(6*24*time.Hour).Format(time.RFC822)))

	writeSource(t, filepath.Join("content", "categories", "hugo", filename), fmt.Sprintf(`---
title: Taxonomy Hugo
date : %q
lastMod : %q
---
Taxonomy Hugo **Content!**
`, date.Add(7*24*time.Hour).Format(time.RFC822), date.Add(8*24*time.Hour).Format(time.RFC822)))

	writeSource(t, filepath.Join("content", "categories", "web", filename), fmt.Sprintf(`---
title: Taxonomy Web
date : %q
lastMod : %q
---
Taxonomy Web **Content!**
`, date.Add(9*24*time.Hour).Format(time.RFC822), date.Add(10*24*time.Hour).Format(time.RFC822)))

	writeSource(t, filepath.Join("content", "categories", filename), fmt.Sprintf(`---
title: Taxonomy Term Categories
date : %q
lastMod : %q
---
Taxonomy Term Categories **Content!**
`, date.Add(11*24*time.Hour).Format(time.RFC822), date.Add(12*24*time.Hour).Format(time.RFC822)))
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
Date: {{ .Date.Format "2006-01-02" }}
Lastmod: {{ .Lastmod.Format "2006-01-02" }}
GetPage: {{ with .Site.GetPage "section" "sect1" }}{{ .Title }}{{ end }} 
`)

	writeSource(t, filepath.Join("layouts", "_default", "single.html"), `
Single Title: {{ .Title }}
Single Content: {{ .Content }}
Date: {{ .Date.Format "2006-01-02" }}
Lastmod: {{ .Lastmod.Format "2006-01-02" }}
`)

	writeSource(t, filepath.Join("layouts", "_default", "section.html"), `
Section Title: {{ .Title }}
Section Content: {{ .Content }}
# Pages: {{ len .Data.Pages }}
{{ range .Paginator.Pages }}
	Pag: {{ .Title }}
{{ end }}
Date: {{ .Date.Format "2006-01-02" }}
Lastmod: {{ .Lastmod.Format "2006-01-02" }}
`)

	// Taxonomy lists
	writeSource(t, filepath.Join("layouts", "_default", "taxonomy.html"), `
Taxonomy Title: {{ .Title }}
Taxonomy Content: {{ .Content }}
# Pages: {{ len .Data.Pages }}
{{ range .Paginator.Pages }}
	Pag: {{ .Title }}
{{ end }}
Date: {{ .Date.Format "2006-01-02" }}
Lastmod: {{ .Lastmod.Format "2006-01-02" }}
`)

	// Taxonomy terms
	writeSource(t, filepath.Join("layouts", "_default", "terms.html"), `
Taxonomy Terms Title: {{ .Title }}
Taxonomy Terms Content: {{ .Content }}
{{ range $key, $value := .Data.Terms }}
	k/v: {{ $key }} / {{ printf "%s" $value }}
{{ end }}
Date: {{ .Date.Format "2006-01-02" }}
Lastmod: {{ .Lastmod.Format "2006-01-02" }}
`)
}
