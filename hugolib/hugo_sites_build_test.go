package hugolib

import (
	"path/filepath"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/resources/kinds"

	"github.com/spf13/afero"
)

func TestMultiSitesMainLangInRoot(t *testing.T) {
	files := `
-- hugo.toml --
defaultContentLanguage = "fr"
defaultContentLanguageInSubdir = false
disableKinds = ["taxonomy", "term"]
[languages]
[languages.en]
weight = 1
[languages.fr]
weight = 2
-- content/sect/doc1.en.md --
---
title: doc1 en
---
-- content/sect/doc1.fr.md --
---
title: doc1 fr
slug: doc1-fr
---
-- layouts/_default/single.html --
Single: {{ .Title }}|{{ .Lang }}|{{ .RelPermalink }}|

`
	b := Test(t, files)
	b.AssertFileContent("public/sect/doc1-fr/index.html", "Single: doc1 fr|fr|/sect/doc1-fr/|")
	b.AssertFileContent("public/en/sect/doc1/index.html", "Single: doc1 en|en|/en/sect/doc1/|")
}

func TestMultiSitesWithTwoLanguages(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
defaultContentLanguage = "nn"

[languages]
[languages.nn]
languageName = "Nynorsk"
weight = 1
title = "Tittel på Nynorsk"
[languages.nn.params]
p1 = "p1nn"

[languages.en]
title = "Title in English"
languageName = "English"
weight = 2
[languages.en.params]
p1 = "p1en"
`
	b := Test(t, files, TestOptSkipRender())
	sites := b.H.Sites

	c := qt.New(t)
	c.Assert(len(sites), qt.Equals, 2)

	nnSite := sites[0]
	nnHome := nnSite.getPageOldVersion(kinds.KindHome)
	c.Assert(len(nnHome.AllTranslations()), qt.Equals, 2)
	c.Assert(len(nnHome.Translations()), qt.Equals, 1)
	c.Assert(nnHome.IsTranslated(), qt.Equals, true)

	enHome := sites[1].getPageOldVersion(kinds.KindHome)

	p1, err := enHome.Param("p1")
	c.Assert(err, qt.IsNil)
	c.Assert(p1, qt.Equals, "p1en")

	p1, err = nnHome.Param("p1")
	c.Assert(err, qt.IsNil)
	c.Assert(p1, qt.Equals, "p1nn")
}

func writeToFs(t testing.TB, fs afero.Fs, filename, content string) {
	t.Helper()
	if err := afero.WriteFile(fs, filepath.FromSlash(filename), []byte(content), 0o755); err != nil {
		t.Fatalf("Failed to write file: %s", err)
	}
}
