package hugolib

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/common/paths"
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

func TestBenchmarkAssembleDeepSiteWithManySections(t *testing.T) {
	t.Parallel()

	b := createBenchmarkAssembleDeepSiteWithManySectionsBuilder(t, false, 2, 3, 4).Build()
	b.AssertFileContent("public/index.html", "Num regular pages recursive: 48|")
}

func createBenchmarkAssembleDeepSiteWithManySectionsBuilder(t testing.TB, skipRender bool, sectionDepth, sectionsPerLevel, pagesPerSection int) *IntegrationTestBuilder {
	t.Helper()

	const contentTemplate = `---
title: P%d
---

## A title

Some content with a shortcode: {{< foo >}}.

Some more content and then another shortcode: {{< foo >}}.

Some final content.
`

	const filesTemplate = `
-- hugo.toml --
baseURL = "http://example.org/"
disableKinds = ["taxonomy", "term", "rss", "sitemap", "robotsTXT", "404"]
-- layouts/all.html --
All.{{ .Title }}|{{ .Content }}|Num pages: {{ len .Pages }}|Num sections: {{ len .Sections }}|Num regular pages recursive: {{ len .RegularPagesRecursive }}|
Sections: {{ range .Sections }}{{ .Title }}|{{ end }}|
RegularPagesRecursive: {{ range .RegularPagesRecursive }}{{ .RelPermalink }}|{{ end }}|
-- layouts/_shortcodes/foo.html --
`
	page := func(section string, i int) string {
		return fmt.Sprintf("\n-- content/%s/p%d.md --\n"+contentTemplate, section, i, i)
	}

	section := func(section string, i int) string {
		if section != "" {
			section = paths.AddTrailingSlash(section)
		}
		return fmt.Sprintf("\n-- content/%ss%d/_index.md --\n"+contentTemplate, section, i, i)
	}

	var sb strings.Builder

	// s0
	// s0/s0
	// s0/s1
	// etc.
	var (
		pageCount    int
		sectionCount int
	)
	var createSections func(currentSection string, currentDepth int)
	createSections = func(currentSection string, currentDepth int) {
		if currentDepth > sectionDepth {
			return
		}

		for i := 0; i < sectionsPerLevel; i++ {
			sectionCount++
			sectionName := fmt.Sprintf("s%d", i)
			sectionPath := sectionName
			if currentSection != "" {
				sectionPath = currentSection + "/" + sectionName
			}
			sb.WriteString(section(currentSection, i))

			// Pages in this section
			for j := 0; j < pagesPerSection; j++ {
				pageCount++
				sb.WriteString(page(sectionPath, j))
			}

			// Recurse
			createSections(sectionPath, currentDepth+1)
		}
	}

	createSections("", 1)

	sb.WriteString(filesTemplate)

	files := sb.String()

	return NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			BuildCfg: BuildCfg{
				SkipRender: skipRender,
			},
		},
	)
}

func BenchmarkAssembleDeepSiteWithManySections(b *testing.B) {
	runOne := func(sectionDepth, sectionsPerLevel, pagesPerSection int) {
		name := fmt.Sprintf("depth=%d/sectionsPerLevel=%d/pagesPerSection=%d", sectionDepth, sectionsPerLevel, pagesPerSection)
		b.Run(name, func(b *testing.B) {
			for b.Loop() {
				b.StopTimer()
				bt := createBenchmarkAssembleDeepSiteWithManySectionsBuilder(b, true, sectionDepth, sectionsPerLevel, pagesPerSection)
				b.StartTimer()
				bt.Build()
			}
		})
	}

	runOne(1, 1, 50)
	runOne(1, 6, 100)
	runOne(1, 6, 500)
	runOne(2, 6, 100)
	runOne(4, 2, 100)
}
