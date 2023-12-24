package hugolib

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestMultihost(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
paginate = 1
defaultContentLanguage = "fr"
defaultContentLanguageInSubdir = false
staticDir = ["s1", "s2"]
enableRobotsTXT = true

[permalinks]
other = "/somewhere/else/:filename"

[taxonomies]
tag = "tags"

[languages]
[languages.en]
staticDir2 = ["staticen"]
baseURL = "https://example.com/docs"
weight = 10
title = "In English"
languageName = "English"
[languages.fr]
staticDir2 = ["staticfr"]
baseURL = "https://example.fr"
weight = 20
title = "Le Français"
languageName = "Français"
-- assets/css/main.css --
body { color: red; }
-- content/mysect/mybundle/index.md --
---
tags: [a, b]
title: "My Bundle fr"
---
My Bundle
-- content/mysect/mybundle/index.en.md --
---
tags: [c, d]
title: "My Bundle en"
---
My Bundle
-- content/mysect/mybundle/foo.txt --
Foo
-- layouts/_default/list.html --
List|{{ .Title }}|{{ .Lang }}|{{ .Permalink}}|{{ .RelPermalink }}|
-- layouts/_default/single.html --
Single|{{ .Title }}|{{ .Lang }}|{{ .Permalink}}|{{ .RelPermalink }}|
{{ $foo := .Resources.Get "foo.txt" | fingerprint }}
Foo: {{ $foo.Permalink }}|
{{ $css := resources.Get "css/main.css" | fingerprint }}
CSS: {{ $css.Permalink }}|{{ $css.RelPermalink }}|
-- layouts/robots.txt --
robots|{{ site.Language.Lang }}
-- layouts/404.html --
404|{{ site.Language.Lang }}


	
`

	b := Test(t, files)

	b.Assert(b.H.Conf.IsMultiLingual(), qt.Equals, true)
	b.Assert(b.H.Conf.IsMultihost(), qt.Equals, true)

	// helpers.PrintFs(b.H.Fs.PublishDir, "", os.Stdout)

	// Check regular pages.
	b.AssertFileContent("public/en/mysect/mybundle/index.html", "Single|My Bundle en|en|https://example.com/docs/mysect/mybundle/|")
	b.AssertFileContent("public/fr/mysect/mybundle/index.html", "Single|My Bundle fr|fr|https://example.fr/mysect/mybundle/|")

	// Check robots.txt
	b.AssertFileContent("public/en/robots.txt", "robots|en")
	b.AssertFileContent("public/fr/robots.txt", "robots|fr")

	// Check sitemap.xml
	b.AssertFileContent("public/en/sitemap.xml", "https://example.com/docs/mysect/mybundle/")
	b.AssertFileContent("public/fr/sitemap.xml", "https://example.fr/mysect/mybundle/")

	// Check 404
	b.AssertFileContent("public/en/404.html", "404|en")
	b.AssertFileContent("public/fr/404.html", "404|fr")

	// Check tags.
	b.AssertFileContent("public/en/tags/d/index.html", "List|D|en|https://example.com/docs/tags/d/")
	b.AssertFileContent("public/fr/tags/b/index.html", "List|B|fr|https://example.fr/tags/b/")
	b.AssertFileExists("public/en/tags/b/index.html", false)
	b.AssertFileExists("public/fr/tags/d/index.html", false)

	// en/mysect/mybundle/foo.txt fingerprinted
	b.AssertFileContent("public/en/mysect/mybundle/foo.1cbec737f863e4922cee63cc2ebbfaafcd1cff8b790d8cfd2e6a5d550b648afa.txt", "Foo")
	b.AssertFileContent("public/en/mysect/mybundle/index.html", "Foo: https://example.fr/mysect/mybundle/foo.1cbec737f863e4922cee63cc2ebbfaafcd1cff8b790d8cfd2e6a5d550b648afa.txt|")
	b.AssertFileContent("public/fr/mysect/mybundle/foo.1cbec737f863e4922cee63cc2ebbfaafcd1cff8b790d8cfd2e6a5d550b648afa.txt", "Foo")
	b.AssertFileContent("public/fr/mysect/mybundle/index.html", "Foo: https://example.fr/mysect/mybundle/foo.1cbec737f863e4922cee63cc2ebbfaafcd1cff8b790d8cfd2e6a5d550b648afa.txt|")

	// Assets CSS fingerprinted
	b.AssertFileContent("public/en/mysect/mybundle/index.html", "CSS: https://example.fr/css/main.5de625c36355cce7c1d5408826a0b21abfb49fb6c0e1f16c945a6f2aef38200c.css|")
	b.AssertFileContent("public/en/css/main.5de625c36355cce7c1d5408826a0b21abfb49fb6c0e1f16c945a6f2aef38200c.css", "body { color: red; }")
	b.AssertFileContent("public/fr/mysect/mybundle/index.html", "CSS: https://example.fr/css/main.5de625c36355cce7c1d5408826a0b21abfb49fb6c0e1f16c945a6f2aef38200c.css|")
	b.AssertFileContent("public/fr/css/main.5de625c36355cce7c1d5408826a0b21abfb49fb6c0e1f16c945a6f2aef38200c.css", "body { color: red; }")
}
