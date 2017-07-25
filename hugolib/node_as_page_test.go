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
	"strings"
	"testing"

	"time"

	"github.com/spf13/afero"

	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/stretchr/testify/require"
)

/*
	This file will test the "making everything a page" transition.

	See https://github.com/gohugoio/hugo/issues/2297

*/

func TestNodesAsPage(t *testing.T) {
	t.Parallel()
	for _, preserveTaxonomyNames := range []bool{false, true} {
		for _, ugly := range []bool{true, false} {
			doTestNodeAsPage(t, ugly, preserveTaxonomyNames)
		}
	}
}

func doTestNodeAsPage(t *testing.T, ugly, preserveTaxonomyNames bool) {

	/* Will have to decide what to name the node content files, but:

		Home page should have:
		Content, shortcode support
	   	Metadata (title, dates etc.)
		Params
	   	Taxonomies (categories, tags)

	*/

	var (
		cfg, fs = newTestCfg()
		th      = testHelper{cfg, fs, t}
	)

	cfg.Set("uglyURLs", ugly)
	cfg.Set("preserveTaxonomyNames", preserveTaxonomyNames)

	cfg.Set("paginate", 1)
	cfg.Set("title", "Hugo Rocks")
	cfg.Set("rssURI", "customrss.xml")

	writeLayoutsForNodeAsPageTests(t, fs)
	writeNodePagesForNodeAsPageTests(t, fs, "")

	writeRegularPagesForNodeAsPageTests(t, fs)

	sites, err := NewHugoSites(deps.DepsCfg{Fs: fs, Cfg: cfg})

	require.NoError(t, err)

	require.NoError(t, sites.Build(BuildCfg{}))

	// date order: home, sect1, sect2, cat/hugo, cat/web, categories

	th.assertFileContent(filepath.Join("public", "index.html"),
		"Index Title: Home Sweet Home!",
		"Home <strong>Content!</strong>",
		"# Pages: 4",
		"Date: 2009-01-02",
		"Lastmod: 2009-01-03",
		"GetPage: Section1 ",
	)

	th.assertFileContent(expectedFilePath(ugly, "public", "sect1", "regular1"), "Single Title: Page 01", "Content Page 01")

	nodes := sites.findAllPagesByKindNotIn(KindPage)

	require.Len(t, nodes, 8)

	home := nodes[7] // oldest

	require.True(t, home.IsHome())
	require.True(t, home.IsNode())
	require.False(t, home.IsPage())
	require.True(t, home.Path() != "")

	section2 := nodes[5]
	require.Equal(t, "Section2", section2.Title)

	pages := sites.findAllPagesByKind(KindPage)
	require.Len(t, pages, 4)

	first := pages[0]

	require.False(t, first.IsHome())
	require.False(t, first.IsNode())
	require.True(t, first.IsPage())

	// Check Home paginator
	th.assertFileContent(expectedFilePath(ugly, "public", "page", "2"),
		"Pag: Page 02")

	// Check Sections
	th.assertFileContent(expectedFilePath(ugly, "public", "sect1"),
		"Section Title: Section", "Section1 <strong>Content!</strong>",
		"Date: 2009-01-04",
		"Lastmod: 2009-01-05",
	)

	th.assertFileContent(expectedFilePath(ugly, "public", "sect2"),
		"Section Title: Section", "Section2 <strong>Content!</strong>",
		"Date: 2009-01-06",
		"Lastmod: 2009-01-07",
	)

	// Check Sections paginator
	th.assertFileContent(expectedFilePath(ugly, "public", "sect1", "page", "2"),
		"Pag: Page 02")

	sections := sites.findAllPagesByKind(KindSection)

	require.Len(t, sections, 2)

	// Check taxonomy lists
	th.assertFileContent(expectedFilePath(ugly, "public", "categories", "hugo"),
		"Taxonomy Title: Taxonomy Hugo", "Taxonomy Hugo <strong>Content!</strong>",
		"Date: 2009-01-08",
		"Lastmod: 2009-01-09",
	)

	th.assertFileContent(expectedFilePath(ugly, "public", "categories", "hugo-rocks"),
		"Taxonomy Title: Taxonomy Hugo Rocks",
	)

	s := sites.Sites[0]

	web := s.getPage(KindTaxonomy, "categories", "web")
	require.NotNil(t, web)
	require.Len(t, web.Data["Pages"].(Pages), 4)

	th.assertFileContent(expectedFilePath(ugly, "public", "categories", "web"),
		"Taxonomy Title: Taxonomy Web",
		"Taxonomy Web <strong>Content!</strong>",
		"Date: 2009-01-10",
		"Lastmod: 2009-01-11",
	)

	// Check taxonomy list paginator
	th.assertFileContent(expectedFilePath(ugly, "public", "categories", "hugo", "page", "2"),
		"Taxonomy Title: Taxonomy Hugo",
		"Pag: Page 02")

	// Check taxonomy terms
	th.assertFileContent(expectedFilePath(ugly, "public", "categories"),
		"Taxonomy Terms Title: Taxonomy Term Categories", "Taxonomy Term Categories <strong>Content!</strong>", "k/v: hugo",
		"Date: 2009-01-14",
		"Lastmod: 2009-01-15",
	)

	// Check taxonomy terms paginator
	th.assertFileContent(expectedFilePath(ugly, "public", "categories", "page", "2"),
		"Taxonomy Terms Title: Taxonomy Term Categories",
		"Pag: Taxonomy Web")

	// RSS
	th.assertFileContent(filepath.Join("public", "customrss.xml"), "Recent content in Home Sweet Home! on Hugo Rocks", "<rss")
	th.assertFileContent(filepath.Join("public", "sect1", "customrss.xml"), "Recent content in Section1 on Hugo Rocks", "<rss")
	th.assertFileContent(filepath.Join("public", "sect2", "customrss.xml"), "Recent content in Section2 on Hugo Rocks", "<rss")
	th.assertFileContent(filepath.Join("public", "categories", "hugo", "customrss.xml"), "Recent content in Taxonomy Hugo on Hugo Rocks", "<rss")
	th.assertFileContent(filepath.Join("public", "categories", "web", "customrss.xml"), "Recent content in Taxonomy Web on Hugo Rocks", "<rss")
	th.assertFileContent(filepath.Join("public", "categories", "customrss.xml"), "Recent content in Taxonomy Term Categories on Hugo Rocks", "<rss")

}

func TestNodesWithNoContentFile(t *testing.T) {
	t.Parallel()
	for _, ugly := range []bool{false, true} {
		doTestNodesWithNoContentFile(t, ugly)
	}
}

func doTestNodesWithNoContentFile(t *testing.T, ugly bool) {

	var (
		cfg, fs = newTestCfg()
		th      = testHelper{cfg, fs, t}
	)

	cfg.Set("uglyURLs", ugly)
	cfg.Set("paginate", 1)
	cfg.Set("title", "Hugo Rocks!")
	cfg.Set("rssURI", "customrss.xml")

	writeLayoutsForNodeAsPageTests(t, fs)
	writeRegularPagesForNodeAsPageTests(t, fs)

	sites, err := NewHugoSites(deps.DepsCfg{Fs: fs, Cfg: cfg})

	require.NoError(t, err)

	require.NoError(t, sites.Build(BuildCfg{}))

	s := sites.Sites[0]

	// Home page
	homePages := s.findPagesByKind(KindHome)
	require.Len(t, homePages, 1)

	homePage := homePages[0]
	require.Len(t, homePage.Data["Pages"], 4)
	require.Len(t, homePage.Pages, 4)
	require.True(t, homePage.Path() == "")

	th.assertFileContent(filepath.Join("public", "index.html"),
		"Index Title: Hugo Rocks!",
		"Date: 2010-06-12",
		"Lastmod: 2010-06-13",
	)

	// Taxonomy list
	th.assertFileContent(expectedFilePath(ugly, "public", "categories", "hugo"),
		"Taxonomy Title: Hugo",
		"Date: 2010-06-12",
		"Lastmod: 2010-06-13",
	)

	// Taxonomy terms
	th.assertFileContent(expectedFilePath(ugly, "public", "categories"),
		"Taxonomy Terms Title: Categories",
	)

	pages := s.findPagesByKind(KindTaxonomyTerm)
	for _, p := range pages {
		var want string
		if ugly {
			want = "/" + p.s.PathSpec.URLize(p.Title) + ".html"
		} else {
			want = "/" + p.s.PathSpec.URLize(p.Title) + "/"
		}
		if p.URL() != want {
			t.Errorf("Taxonomy term URL mismatch: want %q, got %q", want, p.URL())
		}
	}

	// Sections
	th.assertFileContent(expectedFilePath(ugly, "public", "sect1"),
		"Section Title: Sect1s",
		"Date: 2010-06-12",
		"Lastmod: 2010-06-13",
	)

	th.assertFileContent(expectedFilePath(ugly, "public", "sect2"),
		"Section Title: Sect2s",
		"Date: 2008-07-06",
		"Lastmod: 2008-07-09",
	)

	// RSS
	th.assertFileContent(filepath.Join("public", "customrss.xml"), "Hugo Rocks!", "<rss")
	th.assertFileContent(filepath.Join("public", "sect1", "customrss.xml"), "Recent content in Sect1s on Hugo Rocks!", "<rss")
	th.assertFileContent(filepath.Join("public", "sect2", "customrss.xml"), "Recent content in Sect2s on Hugo Rocks!", "<rss")
	th.assertFileContent(filepath.Join("public", "categories", "hugo", "customrss.xml"), "Recent content in Hugo on Hugo Rocks!", "<rss")
	th.assertFileContent(filepath.Join("public", "categories", "web", "customrss.xml"), "Recent content in Web on Hugo Rocks!", "<rss")

}

func TestNodesAsPageMultilingual(t *testing.T) {
	t.Parallel()
	for _, ugly := range []bool{false, true} {
		t.Run(fmt.Sprintf("ugly=%t", ugly), func(t *testing.T) {
			doTestNodesAsPageMultilingual(t, ugly)
		})
	}
}

func doTestNodesAsPageMultilingual(t *testing.T, ugly bool) {

	mf := afero.NewMemMapFs()

	writeToFs(t, mf, "config.toml",
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

	cfg, err := LoadConfig(mf, "", "config.toml")
	require.NoError(t, err)

	cfg.Set("uglyURLs", ugly)

	fs := hugofs.NewFrom(mf, cfg)

	writeLayoutsForNodeAsPageTests(t, fs)

	for _, lang := range []string{"nn", "en"} {
		writeRegularPagesForNodeAsPageTestsWithLang(t, fs, lang)
	}

	th := testHelper{cfg, fs, t}

	sites, err := NewHugoSites(deps.DepsCfg{Fs: fs, Cfg: cfg})

	if err != nil {
		t.Fatalf("Failed to create sites: %s", err)
	}

	if len(sites.Sites) != 3 {
		t.Fatalf("Got %d sites", len(sites.Sites))
	}

	// Only write node pages for the English and Deutsch
	writeNodePagesForNodeAsPageTests(t, fs, "en")
	writeNodePagesForNodeAsPageTests(t, fs, "de")

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
	// See issue #3179
	require.Equal(t, expetedPermalink(false, "/de/"), deHome.Permalink())

	enSect := sites.Sites[1].getPage("section", "sect1")
	require.NotNil(t, enSect)
	require.Equal(t, "en", enSect.Language().Lang)
	require.Len(t, enSect.Translations(), 2, enSect.Translations()[0].Language().Lang)
	require.Equal(t, "de", enSect.Translations()[1].Language().Lang)
	require.Equal(t, "nn", enSect.Translations()[0].Language().Lang)

	require.Equal(t, expetedPermalink(ugly, "/en/sect1/"), enSect.Permalink())

	th.assertFileContent(filepath.Join("public", "nn", "index.html"),
		"Index Title: Hugo på norsk")
	th.assertFileContent(filepath.Join("public", "en", "index.html"),
		"Index Title: Home Sweet Home!", "<strong>Content!</strong>")
	th.assertFileContent(filepath.Join("public", "de", "index.html"),
		"Index Title: Home Sweet Home!", "<strong>Content!</strong>")

	// Taxonomy list
	th.assertFileContent(expectedFilePath(ugly, "public", "nn", "categories", "hugo"),
		"Taxonomy Title: Hugo")
	th.assertFileContent(expectedFilePath(ugly, "public", "en", "categories", "hugo"),
		"Taxonomy Title: Taxonomy Hugo")

	// Taxonomy terms
	th.assertFileContent(expectedFilePath(ugly, "public", "nn", "categories"),
		"Taxonomy Terms Title: Categories")
	th.assertFileContent(expectedFilePath(ugly, "public", "en", "categories"),
		"Taxonomy Terms Title: Taxonomy Term Categories")

	// Sections
	th.assertFileContent(expectedFilePath(ugly, "public", "nn", "sect1"),
		"Section Title: Sect1s")
	th.assertFileContent(expectedFilePath(ugly, "public", "nn", "sect2"),
		"Section Title: Sect2s")
	th.assertFileContent(expectedFilePath(ugly, "public", "en", "sect1"),
		"Section Title: Section1")
	th.assertFileContent(expectedFilePath(ugly, "public", "en", "sect2"),
		"Section Title: Section2")

	// Regular pages
	th.assertFileContent(expectedFilePath(ugly, "public", "en", "sect1", "regular1"),
		"Single Title: Page 01")
	th.assertFileContent(expectedFilePath(ugly, "public", "nn", "sect1", "regular2"),
		"Single Title: Page 02")

	// RSS
	th.assertFileContent(filepath.Join("public", "nn", "customrss.xml"), "Hugo på norsk", "<rss")
	th.assertFileContent(filepath.Join("public", "nn", "sect1", "customrss.xml"), "Recent content in Sect1s on Hugo på norsk", "<rss")
	th.assertFileContent(filepath.Join("public", "nn", "sect2", "customrss.xml"), "Recent content in Sect2s on Hugo på norsk", "<rss")
	th.assertFileContent(filepath.Join("public", "nn", "categories", "hugo", "customrss.xml"), "Recent content in Hugo on Hugo på norsk", "<rss")
	th.assertFileContent(filepath.Join("public", "nn", "categories", "web", "customrss.xml"), "Recent content in Web on Hugo på norsk", "<rss")

	th.assertFileContent(filepath.Join("public", "en", "customrss.xml"), "Recent content in Home Sweet Home! on Hugo in English", "<rss")
	th.assertFileContent(filepath.Join("public", "en", "sect1", "customrss.xml"), "Recent content in Section1 on Hugo in English", "<rss")
	th.assertFileContent(filepath.Join("public", "en", "sect2", "customrss.xml"), "Recent content in Section2 on Hugo in English", "<rss")
	th.assertFileContent(filepath.Join("public", "en", "categories", "hugo", "customrss.xml"), "Recent content in Taxonomy Hugo on Hugo in English", "<rss")
	th.assertFileContent(filepath.Join("public", "en", "categories", "web", "customrss.xml"), "Recent content in Taxonomy Web on Hugo in English", "<rss")

}

func TestNodesWithTaxonomies(t *testing.T) {
	t.Parallel()
	var (
		cfg, fs = newTestCfg()
		th      = testHelper{cfg, fs, t}
	)

	cfg.Set("paginate", 1)
	cfg.Set("title", "Hugo Rocks!")

	writeLayoutsForNodeAsPageTests(t, fs)
	writeRegularPagesForNodeAsPageTests(t, fs)

	writeSource(t, fs, filepath.Join("content", "_index.md"), `---
title: Home With Taxonomies
categories:  [
        "Hugo",	
		"Home"
]
---
`)

	h, err := NewHugoSites(deps.DepsCfg{Fs: fs, Cfg: cfg})

	require.NoError(t, err)

	require.NoError(t, h.Build(BuildCfg{}))

	th.assertFileContent(filepath.Join("public", "categories", "hugo", "index.html"), "Taxonomy Title: Hugo", "# Pages: 5")
	th.assertFileContent(filepath.Join("public", "categories", "home", "index.html"), "Taxonomy Title: Home", "# Pages: 1")

}

func TestNodesWithMenu(t *testing.T) {
	t.Parallel()
	var (
		cfg, fs = newTestCfg()
		th      = testHelper{cfg, fs, t}
	)

	cfg.Set("paginate", 1)
	cfg.Set("title", "Hugo Rocks!")

	writeLayoutsForNodeAsPageTests(t, fs)
	writeRegularPagesForNodeAsPageTests(t, fs)

	writeSource(t, fs, filepath.Join("content", "_index.md"), `---
title: Home With Menu
menu:
  mymenu:
    name: "Go Home!"
---
`)

	writeSource(t, fs, filepath.Join("content", "sect1", "_index.md"), `---
title: Sect1 With Menu
menu:
  mymenu:
    name: "Go Sect1!"
---
`)

	writeSource(t, fs, filepath.Join("content", "categories", "hugo", "_index.md"), `---
title: Taxonomy With Menu
menu:
  mymenu:
    name: "Go Tax Hugo!"
---
`)

	h, err := NewHugoSites(deps.DepsCfg{Fs: fs, Cfg: cfg})

	require.NoError(t, err)

	require.NoError(t, h.Build(BuildCfg{}))

	th.assertFileContent(filepath.Join("public", "index.html"), "Home With Menu", "Home Menu Item: Go Home!: /")
	th.assertFileContent(filepath.Join("public", "sect1", "index.html"), "Sect1 With Menu", "Section Menu Item: Go Sect1!: /sect1/")
	th.assertFileContent(filepath.Join("public", "categories", "hugo", "index.html"), "Taxonomy With Menu", "Taxonomy Menu Item: Go Tax Hugo!: /categories/hugo/")

}

func TestNodesWithAlias(t *testing.T) {
	t.Parallel()
	var (
		cfg, fs = newTestCfg()
		th      = testHelper{cfg, fs, t}
	)

	cfg.Set("paginate", 1)
	cfg.Set("baseURL", "http://base/")
	cfg.Set("title", "Hugo Rocks!")

	writeLayoutsForNodeAsPageTests(t, fs)
	writeRegularPagesForNodeAsPageTests(t, fs)

	writeSource(t, fs, filepath.Join("content", "_index.md"), `---
title: Home With Alias
aliases:
    - /my/new/home.html
---
`)

	h, err := NewHugoSites(deps.DepsCfg{Fs: fs, Cfg: cfg})

	require.NoError(t, err)

	require.NoError(t, h.Build(BuildCfg{}))

	th.assertFileContent(filepath.Join("public", "index.html"), "Home With Alias")
	th.assertFileContent(filepath.Join("public", "my", "new", "home.html"), "content=\"0; url=http://base/")

}

func TestNodesWithSectionWithIndexPageOnly(t *testing.T) {
	t.Parallel()
	var (
		cfg, fs = newTestCfg()
		th      = testHelper{cfg, fs, t}
	)

	cfg.Set("paginate", 1)
	cfg.Set("title", "Hugo Rocks!")

	writeLayoutsForNodeAsPageTests(t, fs)

	writeSource(t, fs, filepath.Join("content", "sect", "_index.md"), `---
title: MySection
---
My Section Content
`)

	h, err := NewHugoSites(deps.DepsCfg{Fs: fs, Cfg: cfg})

	require.NoError(t, err)

	require.NoError(t, h.Build(BuildCfg{}))

	th.assertFileContent(filepath.Join("public", "sect", "index.html"), "My Section")

}

func TestNodesWithURLs(t *testing.T) {
	t.Parallel()
	var (
		cfg, fs = newTestCfg()
		th      = testHelper{cfg, fs, t}
	)

	cfg.Set("paginate", 1)
	cfg.Set("title", "Hugo Rocks!")
	cfg.Set("baseURL", "http://bep.is/base/")

	writeLayoutsForNodeAsPageTests(t, fs)
	writeRegularPagesForNodeAsPageTests(t, fs)

	writeSource(t, fs, filepath.Join("content", "sect", "_index.md"), `---
title: MySection
url: foo.html
---
My Section Content
`)

	h, err := NewHugoSites(deps.DepsCfg{Fs: fs, Cfg: cfg})

	require.NoError(t, err)

	require.NoError(t, h.Build(BuildCfg{}))

	th.assertFileContent(filepath.Join("public", "sect", "index.html"), "My Section")

	s := h.Sites[0]

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

func writeRegularPagesForNodeAsPageTests(t *testing.T, fs *hugofs.Fs) {
	writeRegularPagesForNodeAsPageTestsWithLang(t, fs, "")
}

func writeRegularPagesForNodeAsPageTestsWithLang(t *testing.T, fs *hugofs.Fs, lang string) {
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
		writeSource(t, fs, filepath.Join("content", sect, fmt.Sprintf("regular%d.%smd", i, langStr)), fmt.Sprintf(`---
title: Page %02d
lastMod : %q
date : %q
categories:  [
        "Hugo",
		"Web",
		"Hugo Rocks!"
]
---
Content Page %02d
`, i, date.Add(time.Duration(i)*-24*time.Hour).Format(time.RFC822), date.Add(time.Duration(i)*-2*24*time.Hour).Format(time.RFC822), i))
	}
}

func writeNodePagesForNodeAsPageTests(t *testing.T, fs *hugofs.Fs, lang string) {

	filename := "_index.md"

	if lang != "" {
		filename = fmt.Sprintf("_index.%s.md", lang)
	}

	format := "2006-01-02"

	date, _ := time.Parse(format, "2009-01-01")

	writeSource(t, fs, filepath.Join("content", filename), fmt.Sprintf(`---
title: Home Sweet Home!
date : %q
lastMod : %q
---
l-%s Home **Content!**
`, date.Add(1*24*time.Hour).Format(time.RFC822), date.Add(2*24*time.Hour).Format(time.RFC822), lang))

	writeSource(t, fs, filepath.Join("content", "sect1", filename), fmt.Sprintf(`---
title: Section1
date : %q
lastMod : %q
---
Section1 **Content!**
`, date.Add(3*24*time.Hour).Format(time.RFC822), date.Add(4*24*time.Hour).Format(time.RFC822)))
	writeSource(t, fs, filepath.Join("content", "sect2", filename), fmt.Sprintf(`---
title: Section2
date : %q
lastMod : %q
---
Section2 **Content!**
`, date.Add(5*24*time.Hour).Format(time.RFC822), date.Add(6*24*time.Hour).Format(time.RFC822)))

	writeSource(t, fs, filepath.Join("content", "categories", "hugo", filename), fmt.Sprintf(`---
title: Taxonomy Hugo
date : %q
lastMod : %q
---
Taxonomy Hugo **Content!**
`, date.Add(7*24*time.Hour).Format(time.RFC822), date.Add(8*24*time.Hour).Format(time.RFC822)))

	writeSource(t, fs, filepath.Join("content", "categories", "web", filename), fmt.Sprintf(`---
title: Taxonomy Web
date : %q
lastMod : %q
---
Taxonomy Web **Content!**
`, date.Add(9*24*time.Hour).Format(time.RFC822), date.Add(10*24*time.Hour).Format(time.RFC822)))

	writeSource(t, fs, filepath.Join("content", "categories", "hugo-rocks", filename), fmt.Sprintf(`---
title: Taxonomy Hugo Rocks
date : %q
lastMod : %q
---
Taxonomy Hugo Rocks **Content!**
`, date.Add(11*24*time.Hour).Format(time.RFC822), date.Add(12*24*time.Hour).Format(time.RFC822)))

	writeSource(t, fs, filepath.Join("content", "categories", filename), fmt.Sprintf(`---
title: Taxonomy Term Categories
date : %q
lastMod : %q
---
Taxonomy Term Categories **Content!**
`, date.Add(13*24*time.Hour).Format(time.RFC822), date.Add(14*24*time.Hour).Format(time.RFC822)))

	writeSource(t, fs, filepath.Join("content", "tags", filename), fmt.Sprintf(`---
title: Taxonomy Term Tags
date : %q
lastMod : %q
---
Taxonomy Term Tags **Content!**
`, date.Add(15*24*time.Hour).Format(time.RFC822), date.Add(16*24*time.Hour).Format(time.RFC822)))

}

func writeLayoutsForNodeAsPageTests(t *testing.T, fs *hugofs.Fs) {
	writeSource(t, fs, filepath.Join("layouts", "index.html"), `
Index Title: {{ .Title }}
Index Content: {{ .Content }}
# Pages: {{ len .Data.Pages }}
{{ range .Paginator.Pages }}
	Pag: {{ .Title }}
{{ end }}
{{ with .Site.Menus.mymenu }}
{{ range . }}
Home Menu Item: {{ .Name }}: {{ .URL }}
{{ end }}
{{ end }}
Date: {{ .Date.Format "2006-01-02" }}
Lastmod: {{ .Lastmod.Format "2006-01-02" }}
GetPage: {{ with .Site.GetPage "section" "sect1" }}{{ .Title }}{{ end }} 
`)

	writeSource(t, fs, filepath.Join("layouts", "_default", "single.html"), `
Single Title: {{ .Title }}
Single Content: {{ .Content }}
Date: {{ .Date.Format "2006-01-02" }}
Lastmod: {{ .Lastmod.Format "2006-01-02" }}
`)

	writeSource(t, fs, filepath.Join("layouts", "_default", "section.html"), `
Section Title: {{ .Title }}
Section Content: {{ .Content }}
# Pages: {{ len .Data.Pages }}
{{ range .Paginator.Pages }}
	Pag: {{ .Title }}
{{ end }}
{{ with .Site.Menus.mymenu }}
{{ range . }}
Section Menu Item: {{ .Name }}: {{ .URL }}
{{ end }}
{{ end }}
Date: {{ .Date.Format "2006-01-02" }}
Lastmod: {{ .Lastmod.Format "2006-01-02" }}
`)

	// Taxonomy lists
	writeSource(t, fs, filepath.Join("layouts", "_default", "taxonomy.html"), `
Taxonomy Title: {{ .Title }}
Taxonomy Content: {{ .Content }}
# Pages: {{ len .Data.Pages }}
{{ range .Paginator.Pages }}
	Pag: {{ .Title }}
{{ end }}
{{ with .Site.Menus.mymenu }}
{{ range . }}
Taxonomy Menu Item: {{ .Name }}: {{ .URL }}
{{ end }}
{{ end }}
Date: {{ .Date.Format "2006-01-02" }}
Lastmod: {{ .Lastmod.Format "2006-01-02" }}
`)

	// Taxonomy terms
	writeSource(t, fs, filepath.Join("layouts", "_default", "terms.html"), `
Taxonomy Terms Title: {{ .Title }}
Taxonomy Terms Content: {{ .Content }}
# Pages: {{ len .Data.Pages }}
{{ range .Paginator.Pages }}
	Pag: {{ .Title }}
{{ end }}
{{ range $key, $value := .Data.Terms }}
	k/v: {{ $key | lower }} / {{ printf "%s" $value }}
{{ end }}
{{ with .Site.Menus.mymenu }}
{{ range . }}
Taxonomy Terms Menu Item: {{ .Name }}: {{ .URL }}
{{ end }}
{{ end }}
Date: {{ .Date.Format "2006-01-02" }}
Lastmod: {{ .Lastmod.Format "2006-01-02" }}
`)
}

func expectedFilePath(ugly bool, path ...string) string {
	if ugly {
		return filepath.Join(append(path[0:len(path)-1], path[len(path)-1]+".html")...)
	}
	return filepath.Join(append(path, "index.html")...)
}

func expetedPermalink(ugly bool, path string) string {
	if ugly {
		return strings.TrimSuffix(path, "/") + ".html"
	}
	return path
}
