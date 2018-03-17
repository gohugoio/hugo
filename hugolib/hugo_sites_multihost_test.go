package hugolib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMultihosts(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	var configTemplate = `
paginate = 1
disablePathToLower = true
defaultContentLanguage = "fr"
defaultContentLanguageInSubdir = false
staticDir = ["s1", "s2"]

[permalinks]
other = "/somewhere/else/:filename"

[Taxonomies]
tag = "tags"

[Languages]
[Languages.en]
staticDir2 = ["ens1", "ens2"]
baseURL = "https://example.com/docs"
weight = 10
title = "In English"
languageName = "English"

[Languages.fr]
staticDir2 = ["frs1", "frs2"]
baseURL = "https://example.fr"
weight = 20
title = "Le Français"
languageName = "Français"

[Languages.nn]
staticDir2 = ["nns1", "nns2"]
baseURL = "https://example.no"
weight = 30
title = "På nynorsk"
languageName = "Nynorsk"

`

	b := newMultiSiteTestDefaultBuilder(t).WithConfigFile("toml", configTemplate)
	b.CreateSites().Build(BuildCfg{})

	b.AssertFileContent("public/en/sect/doc1-slug/index.html", "Hello")

	s1 := b.H.Sites[0]

	assert.Equal([]string{"s1", "s2", "ens1", "ens2"}, s1.StaticDirs())

	s1h := s1.getPage(KindHome)
	assert.True(s1h.IsTranslated())
	assert.Len(s1h.Translations(), 2)
	assert.Equal("https://example.com/docs/", s1h.Permalink())

	// For “regular multilingual” we kept the aliases pages with url in front matter
	// as a literal value that we use as is.
	// There is an ambiguity in the guessing.
	// For multihost, we never want any content in the root.
	//
	// check url in front matter:
	pageWithURLInFrontMatter := s1.getPage(KindPage, "sect/doc3.en.md")
	assert.NotNil(pageWithURLInFrontMatter)
	assert.Equal("/superbob", pageWithURLInFrontMatter.URL())
	assert.Equal("/docs/superbob/", pageWithURLInFrontMatter.RelPermalink())
	b.AssertFileContent("public/en/superbob/index.html", "doc3|Hello|en")

	// check alias:
	b.AssertFileContent("public/en/al/alias1/index.html", `content="0; url=https://example.com/docs/superbob/"`)
	b.AssertFileContent("public/en/al/alias2/index.html", `content="0; url=https://example.com/docs/superbob/"`)

	s2 := b.H.Sites[1]
	assert.Equal([]string{"s1", "s2", "frs1", "frs2"}, s2.StaticDirs())

	s2h := s2.getPage(KindHome)
	assert.Equal("https://example.fr/", s2h.Permalink())

	b.AssertFileContent("public/fr/index.html", "French Home Page")
	b.AssertFileContent("public/en/index.html", "Default Home Page")

	// Check paginators
	b.AssertFileContent("public/en/page/1/index.html", `refresh" content="0; url=https://example.com/docs/"`)
	b.AssertFileContent("public/nn/page/1/index.html", `refresh" content="0; url=https://example.no/"`)
	b.AssertFileContent("public/en/sect/page/2/index.html", "List Page 2", "Hello", "https://example.com/docs/sect/", "\"/docs/sect/page/3/")
	b.AssertFileContent("public/fr/sect/page/2/index.html", "List Page 2", "Bonjour", "https://example.fr/sect/")

	// Check bundles

	bundleEn := s1.getPage(KindPage, "bundles/b1/index.en.md")
	require.NotNil(t, bundleEn)
	require.Equal(t, "/docs/bundles/b1/", bundleEn.RelPermalink())
	require.Equal(t, 1, len(bundleEn.Resources))
	logoEn := bundleEn.Resources.GetByPrefix("logo")
	require.NotNil(t, logoEn)
	require.Equal(t, "/docs/bundles/b1/logo.png", logoEn.RelPermalink())
	b.AssertFileContent("public/en/bundles/b1/logo.png", "PNG Data")

	bundleFr := s2.getPage(KindPage, "bundles/b1/index.md")
	require.NotNil(t, bundleFr)
	require.Equal(t, "/bundles/b1/", bundleFr.RelPermalink())
	require.Equal(t, 1, len(bundleFr.Resources))
	logoFr := bundleFr.Resources.GetByPrefix("logo")
	require.NotNil(t, logoFr)
	require.Equal(t, "/bundles/b1/logo.png", logoFr.RelPermalink())
	b.AssertFileContent("public/fr/bundles/b1/logo.png", "PNG Data")

}
