package hugolib

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestMultihost(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
defaultContentLanguage = "fr"
defaultContentLanguageInSubdir = false
staticDir = ["s1", "s2"]
enableRobotsTXT = true

[pagination]
pagerSize = 1

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

	b.Assert(b.H.Conf.IsMultilingual(), qt.Equals, true)
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
	b.AssertFileContent("public/en/mysect/mybundle/index.html", "Foo: https://example.com/docs/mysect/mybundle/foo.1cbec737f863e4922cee63cc2ebbfaafcd1cff8b790d8cfd2e6a5d550b648afa.txt|")
	b.AssertFileContent("public/fr/mysect/mybundle/foo.1cbec737f863e4922cee63cc2ebbfaafcd1cff8b790d8cfd2e6a5d550b648afa.txt", "Foo")
	b.AssertFileContent("public/fr/mysect/mybundle/index.html", "Foo: https://example.fr/mysect/mybundle/foo.1cbec737f863e4922cee63cc2ebbfaafcd1cff8b790d8cfd2e6a5d550b648afa.txt|")

	// Assets CSS fingerprinted
	b.AssertFileContent("public/en/mysect/mybundle/index.html", "CSS: https://example.fr/css/main.5de625c36355cce7c1d5408826a0b21abfb49fb6c0e1f16c945a6f2aef38200c.css|")
	b.AssertFileContent("public/en/css/main.5de625c36355cce7c1d5408826a0b21abfb49fb6c0e1f16c945a6f2aef38200c.css", "body { color: red; }")
	b.AssertFileContent("public/fr/mysect/mybundle/index.html", "CSS: https://example.fr/css/main.5de625c36355cce7c1d5408826a0b21abfb49fb6c0e1f16c945a6f2aef38200c.css|")
	b.AssertFileContent("public/fr/css/main.5de625c36355cce7c1d5408826a0b21abfb49fb6c0e1f16c945a6f2aef38200c.css", "body { color: red; }")
}

func TestMultihostResourcePerLanguageMultihostMinify(t *testing.T) {
	t.Parallel()
	files := `
-- hugo.toml --
disableKinds = ["taxonomy", "term"]
defaultContentLanguage = "en"
defaultContentLanguageInSubDir = true
[languages]
[languages.en]
baseURL = "https://example.en"
weight = 1
contentDir = "content/en"
[languages.fr]
baseURL = "https://example.fr"
weight = 2
contentDir = "content/fr"
-- content/en/section/mybundle/index.md --
---
title: "Mybundle en"
---
-- content/fr/section/mybundle/index.md --
---
title: "Mybundle fr"
---
-- content/en/section/mybundle/styles.css --
.body {
	color: english;
}
-- content/fr/section/mybundle/styles.css --
.body {
	color: french;
}
-- layouts/_default/single.html --
{{ $data := .Resources.GetMatch "styles*" | minify }}
{{ .Lang }}: {{ $data.Content}}|{{ $data.RelPermalink }}|

`
	b := Test(t, files)

	b.AssertFileContent("public/fr/section/mybundle/index.html",
		"fr: .body{color:french}|/section/mybundle/styles.min.css|",
	)

	b.AssertFileContent("public/en/section/mybundle/index.html",
		"en: .body{color:english}|/section/mybundle/styles.min.css|",
	)

	b.AssertFileContent("public/en/section/mybundle/styles.min.css", ".body{color:english}")
	b.AssertFileContent("public/fr/section/mybundle/styles.min.css", ".body{color:french}")
}

func TestResourcePerLanguageIssue12163(t *testing.T) {
	files := `
-- hugo.toml --
defaultContentLanguage = 'de'
disableKinds = ['rss','sitemap','taxonomy','term']

[languages.de]
baseURL = 'https://de.example.org/'
contentDir = 'content/de'
weight = 1

[languages.en]
baseURL = 'https://en.example.org/'
contentDir = 'content/en'
weight = 2
-- content/de/mybundle/index.md --
---
title: mybundle-de
---
-- content/de/mybundle/pixel.png --
iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==
-- content/en/mybundle/index.md --
---
title: mybundle-en
---
-- layouts/_default/single.html --
{{ with .Resources.Get "pixel.png" }}
  {{ with .Resize "2x2" }}
    {{ .RelPermalink }}|
  {{ end }}
{{ end }}
`

	b := Test(t, files)

	b.AssertFileExists("public/de/mybundle/index.html", true)
	b.AssertFileExists("public/en/mybundle/index.html", true)

	b.AssertFileExists("public/de/mybundle/pixel.png", true)
	b.AssertFileExists("public/en/mybundle/pixel.png", true)

	b.AssertFileExists("public/de/mybundle/pixel_hu8581513846771248023.png", true)
	// failing test below
	b.AssertFileExists("public/en/mybundle/pixel_hu8581513846771248023.png", true)
}

func TestMultihostResourceOneBaseURLWithSuPath(t *testing.T) {
	files := `
-- hugo.toml --
defaultContentLanguage = "en"
[languages]
[languages.en]
baseURL = "https://example.com/docs"
weight = 1
contentDir = "content/en"
[languages.en.permalinks]
section = "/enpages/:slug/"
[languages.fr]
baseURL = "https://example.fr"
contentDir = "content/fr"
-- content/en/section/mybundle/index.md --
---
title: "Mybundle en"
---
-- content/fr/section/mybundle/index.md --
---
title: "Mybundle fr"
---
-- content/fr/section/mybundle/file1.txt --
File 1 fr.
-- content/en/section/mybundle/file1.txt --
File 1 en.
-- content/en/section/mybundle/file2.txt --
File 2 en.
-- layouts/_default/single.html --
{{ $files := .Resources.Match "file*" }}
Files: {{ range $files }}{{ .Permalink }}|{{ end }}$

`

	b := Test(t, files)

	b.AssertFileContent("public/en/enpages/mybundle-en/index.html", "Files: https://example.com/docs/enpages/mybundle-en/file1.txt|https://example.com/docs/enpages/mybundle-en/file2.txt|$")
	b.AssertFileContent("public/fr/section/mybundle/index.html", "Files: https://example.fr/section/mybundle/file1.txt|https://example.fr/section/mybundle/file2.txt|$")

	b.AssertFileContent("public/en/enpages/mybundle-en/file1.txt", "File 1 en.")
	b.AssertFileContent("public/fr/section/mybundle/file1.txt", "File 1 fr.")
	b.AssertFileContent("public/en/enpages/mybundle-en/file2.txt", "File 2 en.")
	b.AssertFileContent("public/fr/section/mybundle/file2.txt", "File 2 en.")
}

func TestMultihostAllButOneLanguageDisabledIssue12288(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
defaultContentLanguage = "en"
disableLanguages = ["fr"]
#baseURL = "https://example.com"
[languages]
[languages.en]
baseURL = "https://example.en"
weight = 1
[languages.fr]
baseURL = "https://example.fr"
weight = 2
--  assets/css/main.css --
body { color: red; }
-- layouts/index.html --
{{ $css := resources.Get "css/main.css" | minify }}
CSS: {{ $css.Permalink }}|{{ $css.RelPermalink }}|
`

	b := Test(t, files)

	b.AssertFileContent("public/css/main.min.css", "body{color:red}")
	b.AssertFileContent("public/index.html", "CSS: https://example.en/css/main.min.css|/css/main.min.css|")
}
