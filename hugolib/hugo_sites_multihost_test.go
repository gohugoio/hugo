package hugolib

import (
	"testing"

	"github.com/gohugoio/hugo/resources/page/pagekinds"

	qt "github.com/frankban/quicktest"
)

func TestMultihosts(t *testing.T) {
	c := qt.New(t)

	files := `
-- config.toml --
paginate = 1
disablePathToLower = true
defaultContentLanguage = "fr"
defaultContentLanguageInSubdir = false
staticDir = ["s1", "s2"]
enableRobotsTXT = true

[permalinks]
other = "/somewhere/else/:filename"

[Taxonomies]
tag = "tags"

[languages]
[languages.en]
staticDir2 = ["ens1", "ens2"]
baseURL = "https://example.com/docs"
weight = 10
title = "In English"
languageName = "English"

[languages.fr]
staticDir2 = ["frs1", "frs2"]
baseURL = "https://example.fr"
weight = 20
title = "Le Français"
languageName = "Français"

[languages.nn]
staticDir2 = ["nns1", "nns2"]
baseURL = "https://example.no"
weight = 30
title = "På nynorsk"
languageName = "Nynorsk"
-- content/bundles/b1/index.en.md --
---
title: Bundle EN
publishdate: "2000-01-06"
weight: 2001
---
# Bundle Content EN
-- content/bundles/b1/index.md --
---
title: Bundle Default
publishdate: "2000-01-06"
weight: 2002
---
# Bundle Content Default
-- content/bundles/b1/logo.png --
PNG Data
-- content/other/doc5.fr.md --
---
title: doc5
weight: 5
publishdate: "2000-01-06"
---
# doc5
*autre contenu francophone*
NOTE: should use the "permalinks" configuration with :filename
-- content/root.en.md --
---
title: root
weight: 10000
slug: root
publishdate: "2000-01-01"
---
# root
-- content/sect/doc1.en.md --
---
title: doc1
weight: 1
slug: doc1-slug
tags:
  - tag1
publishdate: "2000-01-01"
---
# doc1
*some "content"*
-- content/sect/doc1.fr.md --
---
title: doc1
weight: 1
plaques:
  - FRtag1
  - FRtag2
publishdate: "2000-01-04"
---
# doc1
*quelque "contenu"*
NOTE: date is after "doc3"
-- content/sect/doc2.en.md --
---
title: doc2
weight: 2
publishdate: "2000-01-02"
---
# doc2
*some content*
NOTE: without slug, "doc2" should be used, without ".en" as URL
-- content/sect/doc3.en.md --
---
title: doc3
weight: 3
publishdate: "2000-01-03"
aliases: [/en/al/alias1,/al/alias2/]
tags:
  - tag2
  - tag1
url: /superbob/
---
# doc3
*some content*
NOTE: third 'en' doc, should trigger pagination on home page.
-- content/sect/doc4.md --
---
title: doc4
weight: 4
plaques:
  - FRtag1
publishdate: "2000-01-05"
---
# doc4
*du contenu francophone*
-- i18n/en.toml --
[hello]
other = "Hello"
-- i18n/en.yaml --
hello:
  other: "Hello"
-- i18n/fr.toml --
[hello]
other = "Bonjour"
-- i18n/fr.yaml --
hello:
  other: "Bonjour"
-- i18n/nb.toml --
[hello]
other = "Hallo"
-- i18n/nn.toml --
[hello]
other = "Hallo"
-- layouts/_default/list.html --
List Page {{ $p := .Paginator }}{{ $p.PageNumber }}|{{ .Title }}|{{ i18n "hello" }}|{{ .Permalink }}|Pager: {{ template "_internal/pagination.html" . }}|Kind: {{ .Kind }}|Content: {{ .Content }}|Len Pages: {{ len .Pages }}|Len RegularPages: {{ len .RegularPages }}| HasParent: {{ if .Parent }}YES{{ else }}NO{{ end }}
-- layouts/_default/single.html --
Single: {{ .Title }}|{{ i18n "hello" }}|{{.Language.Lang}}|RelPermalink: {{ .RelPermalink }}|Permalink: {{ .Permalink }}|{{ .Content }}|Resources: {{ range .Resources }}{{ .MediaType }}: {{ .RelPermalink}} -- {{ end }}|Summary: {{ .Summary }}|Truncated: {{ .Truncated }}|Parent: {{ .Parent.Title }}
-- layouts/_default/taxonomy.html --
-- layouts/index.fr.html --
{{ $p := .Paginator }}French Home Page {{ $p.PageNumber }}: {{ .Title }}|{{ .IsHome }}|{{ i18n "hello" }}|{{ .Permalink }}|{{  .Site.Data.hugo.slogan }}|String Resource: {{ ( "Hugo Pipes" | resources.FromString "text/pipes.txt").RelPermalink  }}
-- layouts/index.html --
{{ $p := .Paginator }}Default Home Page {{ $p.PageNumber }}: {{ .Title }}|{{ .IsHome }}|{{ i18n "hello" }}|{{ .Permalink }}|{{  .Site.Data.hugo.slogan }}|String Resource: {{ ( "Hugo Pipes" | resources.FromString "text/pipes.txt").RelPermalink  }}
-- layouts/robots.txt --
robots|{{ .Lang }}|{{ .Title }}
	`

	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:               c,
			NeedsOsFS:       false,
			NeedsNpmInstall: false,
			TxtarString:     files,
		}).Build()

	b.AssertFileContent("public/en/sect/doc1-slug/index.html", "Hello")

	s1 := b.H.Sites[0]

	s1h := s1.getPage(pagekinds.Home)
	c.Assert(s1h.IsTranslated(), qt.Equals, true)
	c.Assert(len(s1h.Translations()), qt.Equals, 2)
	c.Assert(s1h.Permalink(), qt.Equals, "https://example.com/docs/")

	// For “regular multilingual” we kept the aliases pages with url in front matter
	// as a literal value that we use as is.
	// There is an ambiguity in the guessing.
	// For multihost, we never want any content in the root.
	//
	// check url in front matter:
	pageWithURLInFrontMatter := s1.getPage(pagekinds.Page, "sect/doc3.en.md")
	c.Assert(pageWithURLInFrontMatter, qt.Not(qt.IsNil))
	c.Assert(pageWithURLInFrontMatter.RelPermalink(), qt.Equals, "/docs/superbob/")
	b.AssertFileContent("public/en/superbob/index.html", "doc3|Hello|en")

	// the domain root is the language directory for each language, so the robots.txt is created in the language directories
	b.AssertFileContent("public/en/robots.txt", "robots|en")
	b.AssertFileContent("public/fr/robots.txt", "robots|fr")
	b.AssertFileContent("public/nn/robots.txt", "robots|nn")
	b.AssertDestinationExists("public/robots.txt", false)

	// check alias:
	b.AssertFileContent("public/en/al/alias1/index.html", `content="0; url=https://example.com/docs/superbob/"`)
	b.AssertFileContent("public/en/al/alias2/index.html", `content="0; url=https://example.com/docs/superbob/"`)

	s2 := b.H.Sites[1]

	s2h := s2.getPage(pagekinds.Home)
	c.Assert(s2h.Permalink(), qt.Equals, "https://example.fr/")

	b.AssertFileContent("public/fr/index.html", "French Home Page", "String Resource: /text/pipes.txt")
	b.AssertFileContent("public/fr/text/pipes.txt", "Hugo Pipes")
	b.AssertFileContent("public/en/index.html", "Default Home Page", "String Resource: /docs/text/pipes.txt")
	b.AssertFileContent("public/en/text/pipes.txt", "Hugo Pipes")

	// Check paginators
	b.AssertFileContent("public/en/page/1/index.html", `refresh" content="0; url=https://example.com/docs/"`)
	b.AssertFileContent("public/nn/page/1/index.html", `refresh" content="0; url=https://example.no/"`)
	b.AssertFileContent("public/en/sect/page/2/index.html", "List Page 2", "Hello", "https://example.com/docs/sect/", "\"/docs/sect/page/3/")
	b.AssertFileContent("public/fr/sect/page/2/index.html", "List Page 2", "Bonjour", "https://example.fr/sect/")

	// Check bundles

	bundleEn := s1.getPage(pagekinds.Page, "bundles/b1/index.en.md")
	c.Assert(bundleEn, qt.Not(qt.IsNil))
	c.Assert(bundleEn.RelPermalink(), qt.Equals, "/docs/bundles/b1/")
	c.Assert(len(bundleEn.Resources()), qt.Equals, 1)

	b.AssertFileContent("public/en/bundles/b1/logo.png", "PNG Data")
	b.AssertFileContent("public/en/bundles/b1/index.html", " image/png: /docs/bundles/b1/logo.png")

	bundleFr := s2.getPage(pagekinds.Page, "bundles/b1/index.md")
	c.Assert(bundleFr, qt.Not(qt.IsNil))
	c.Assert(bundleFr.RelPermalink(), qt.Equals, "/bundles/b1/")
	c.Assert(len(bundleFr.Resources()), qt.Equals, 1)
	b.AssertFileContent("public/fr/bundles/b1/logo.png", "PNG Data")
	b.AssertFileContent("public/fr/bundles/b1/index.html", " image/png: /bundles/b1/logo.png")
}
