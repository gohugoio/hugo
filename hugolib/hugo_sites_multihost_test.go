package hugolib

import (
	"testing"

	"github.com/gohugoio/hugo/resources/page"

	qt "github.com/frankban/quicktest"
)

func TestMultihosts(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

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

	s1h := s1.getPage(page.KindHome)
	c.Assert(s1h.IsTranslated(), qt.Equals, true)
	c.Assert(len(s1h.Translations()), qt.Equals, 2)
	c.Assert(s1h.Permalink(), qt.Equals, "https://example.com/docs/")

	// For “regular multilingual” we kept the aliases pages with url in front matter
	// as a literal value that we use as is.
	// There is an ambiguity in the guessing.
	// For multihost, we never want any content in the root.
	//
	// check url in front matter:
	pageWithURLInFrontMatter := s1.getPage(page.KindPage, "sect/doc3.en.md")
	c.Assert(pageWithURLInFrontMatter, qt.Not(qt.IsNil))
	c.Assert(pageWithURLInFrontMatter.RelPermalink(), qt.Equals, "/docs/superbob/")
	b.AssertFileContent("public/en/superbob/index.html", "doc3|Hello|en")

	// check alias:
	b.AssertFileContent("public/en/al/alias1/index.html", `content="0; url=https://example.com/docs/superbob/"`)
	b.AssertFileContent("public/en/al/alias2/index.html", `content="0; url=https://example.com/docs/superbob/"`)

	s2 := b.H.Sites[1]

	s2h := s2.getPage(page.KindHome)
	c.Assert(s2h.Permalink(), qt.Equals, "https://example.fr/")

	b.AssertFileContent("public/fr/index.html", "French Home Page", "String Resource: /docs/text/pipes.txt")
	b.AssertFileContent("public/fr/text/pipes.txt", "Hugo Pipes")
	b.AssertFileContent("public/en/index.html", "Default Home Page", "String Resource: /docs/text/pipes.txt")
	b.AssertFileContent("public/en/text/pipes.txt", "Hugo Pipes")

	// Check paginators
	b.AssertFileContent("public/en/page/1/index.html", `refresh" content="0; url=https://example.com/docs/"`)
	b.AssertFileContent("public/nn/page/1/index.html", `refresh" content="0; url=https://example.no/"`)
	b.AssertFileContent("public/en/sect/page/2/index.html", "List Page 2", "Hello", "https://example.com/docs/sect/", "\"/docs/sect/page/3/")
	b.AssertFileContent("public/fr/sect/page/2/index.html", "List Page 2", "Bonjour", "https://example.fr/sect/")

	// Check bundles

	bundleEn := s1.getPage(page.KindPage, "bundles/b1/index.en.md")
	c.Assert(bundleEn, qt.Not(qt.IsNil))
	c.Assert(bundleEn.RelPermalink(), qt.Equals, "/docs/bundles/b1/")
	c.Assert(len(bundleEn.Resources()), qt.Equals, 1)

	b.AssertFileContent("public/en/bundles/b1/logo.png", "PNG Data")
	b.AssertFileContent("public/en/bundles/b1/index.html", " image/png: /docs/bundles/b1/logo.png")

	bundleFr := s2.getPage(page.KindPage, "bundles/b1/index.md")
	c.Assert(bundleFr, qt.Not(qt.IsNil))
	c.Assert(bundleFr.RelPermalink(), qt.Equals, "/bundles/b1/")
	c.Assert(len(bundleFr.Resources()), qt.Equals, 1)
	b.AssertFileContent("public/fr/bundles/b1/logo.png", "PNG Data")
	b.AssertFileContent("public/fr/bundles/b1/index.html", " image/png: /bundles/b1/logo.png")

}
