package hugolib

import (
	"fmt"
	"strings"
	"testing"

	"path/filepath"

	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/afero"
	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/hugofs"
	"github.com/spf13/hugo/source"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	testCommonResetState()
	jww.SetStdoutThreshold(jww.LevelError)

}

func testCommonResetState() {
	hugofs.InitMemFs()
	viper.Reset()
	viper.SetFs(hugofs.Source())
	loadDefaultSettings()

	if err := hugofs.Source().Mkdir("content", 0755); err != nil {
		panic("Content folder creation failed.")
	}

}

func TestMultiSites(t *testing.T) {

	sites := createMultiTestSites(t)

	err := sites.Build(BuildCfg{})

	if err != nil {
		t.Fatalf("Failed to build sites: %s", err)
	}

	enSite := sites.Sites[0]

	assert.Equal(t, "en", enSite.Language.Lang)

	if len(enSite.Pages) != 3 {
		t.Fatal("Expected 3 english pages")
	}
	assert.Len(t, enSite.Source.Files(), 11, "should have 11 source files")
	assert.Len(t, enSite.AllPages, 6, "should have 6 total pages (including translations)")

	doc1en := enSite.Pages[0]
	permalink, err := doc1en.Permalink()
	assert.NoError(t, err, "permalink call failed")
	assert.Equal(t, "http://example.com/blog/en/sect/doc1-slug/", permalink, "invalid doc1.en permalink")
	assert.Len(t, doc1en.Translations(), 1, "doc1-en should have one translation, excluding itself")

	doc2 := enSite.Pages[1]
	permalink, err = doc2.Permalink()
	assert.NoError(t, err, "permalink call failed")
	assert.Equal(t, "http://example.com/blog/en/sect/doc2/", permalink, "invalid doc2 permalink")

	doc3 := enSite.Pages[2]
	permalink, err = doc3.Permalink()
	assert.NoError(t, err, "permalink call failed")
	assert.Equal(t, "http://example.com/blog/superbob", permalink, "invalid doc3 permalink")

	// TODO(bep) multilingo. Check this case. This has url set in frontmatter, but we must split into lang folders
	// The assertion below was missing the /en prefix.
	assert.Equal(t, "/en/superbob", doc3.URL(), "invalid url, was specified on doc3 TODO(bep)")

	assert.Equal(t, doc2.Next, doc3, "doc3 should follow doc2, in .Next")

	doc1fr := doc1en.Translations()[0]
	permalink, err = doc1fr.Permalink()
	assert.NoError(t, err, "permalink call failed")
	assert.Equal(t, "http://example.com/blog/fr/sect/doc1/", permalink, "invalid doc1fr permalink")

	assert.Equal(t, doc1en.Translations()[0], doc1fr, "doc1-en should have doc1-fr as translation")
	assert.Equal(t, doc1fr.Translations()[0], doc1en, "doc1-fr should have doc1-en as translation")
	assert.Equal(t, "fr", doc1fr.Language().Lang)

	doc4 := enSite.AllPages[4]
	permalink, err = doc4.Permalink()
	assert.NoError(t, err, "permalink call failed")
	assert.Equal(t, "http://example.com/blog/fr/sect/doc4/", permalink, "invalid doc4 permalink")
	assert.Len(t, doc4.Translations(), 0, "found translations for doc4")

	doc5 := enSite.AllPages[5]
	permalink, err = doc5.Permalink()
	assert.NoError(t, err, "permalink call failed")
	assert.Equal(t, "http://example.com/blog/fr/somewhere/else/doc5", permalink, "invalid doc5 permalink")

	// Taxonomies and their URLs
	assert.Len(t, enSite.Taxonomies, 1, "should have 1 taxonomy")
	tags := enSite.Taxonomies["tags"]
	assert.Len(t, tags, 2, "should have 2 different tags")
	assert.Equal(t, tags["tag1"][0].Page, doc1en, "first tag1 page should be doc1")

	frSite := sites.Sites[1]

	assert.Equal(t, "fr", frSite.Language.Lang)
	assert.Len(t, frSite.Pages, 3, "should have 3 pages")
	assert.Len(t, frSite.AllPages, 6, "should have 6 total pages (including translations)")

	for _, frenchPage := range frSite.Pages {
		assert.Equal(t, "fr", frenchPage.Lang())
	}

	// Check redirect to main language, French
	languageRedirect := readDestination(t, "public/index.html")
	require.True(t, strings.Contains(languageRedirect, "0; url=http://example.com/blog/fr"), languageRedirect)

	// Check sitemap(s)
	sitemapIndex := readDestination(t, "public/sitemap.xml")
	require.True(t, strings.Contains(sitemapIndex, "<loc>http:/example.com/blog/en/sitemap.xml</loc>"), sitemapIndex)
	require.True(t, strings.Contains(sitemapIndex, "<loc>http:/example.com/blog/fr/sitemap.xml</loc>"), sitemapIndex)
	sitemapEn := readDestination(t, "public/en/sitemap.xml")
	sitemapFr := readDestination(t, "public/fr/sitemap.xml")
	require.True(t, strings.Contains(sitemapEn, "http://example.com/blog/en/sect/doc2/"), sitemapEn)
	require.True(t, strings.Contains(sitemapFr, "http://example.com/blog/fr/sect/doc1/"), sitemapFr)
}

func TestMultiSitesRebuild(t *testing.T) {

	sites := createMultiTestSites(t)
	cfg := BuildCfg{}

	err := sites.Build(cfg)

	if err != nil {
		t.Fatalf("Failed to build sites: %s", err)
	}

	_, err = hugofs.Destination().Open("public/en/sect/doc2/index.html")

	if err != nil {
		t.Fatalf("Unable to locate file")
	}

	enSite := sites.Sites[0]
	frSite := sites.Sites[1]

	assert.Len(t, enSite.Pages, 3)
	assert.Len(t, frSite.Pages, 3)

	// Verify translations
	docEn := readDestination(t, "public/en/sect/doc1-slug/index.html")
	assert.True(t, strings.Contains(docEn, "Hello"), "No Hello")
	docFr := readDestination(t, "public/fr/sect/doc1/index.html")
	assert.True(t, strings.Contains(docFr, "Bonjour"), "No Bonjour")

	for i, this := range []struct {
		preFunc    func(t *testing.T)
		events     []fsnotify.Event
		assertFunc func(t *testing.T)
	}{
		// * Remove doc
		// * Add docs existing languages
		// (Add doc new language: TODO(bep) we should load config.toml as part of these so we can add languages).
		// * Rename file
		// * Change doc
		// * Change a template
		// * Change language file
		{
			nil,
			[]fsnotify.Event{{Name: "content/sect/doc2.en.md", Op: fsnotify.Remove}},
			func(t *testing.T) {
				assert.Len(t, enSite.Pages, 2, "1 en removed")

				// Check build stats
				assert.Equal(t, 1, enSite.draftCount, "Draft")
				assert.Equal(t, 1, enSite.futureCount, "Future")
				assert.Equal(t, 1, enSite.expiredCount, "Expired")
				assert.Equal(t, 0, frSite.draftCount, "Draft")
				assert.Equal(t, 1, frSite.futureCount, "Future")
				assert.Equal(t, 1, frSite.expiredCount, "Expired")
			},
		},
		{
			func(t *testing.T) {
				writeNewContentFile(t, "new_en_1", "2016-07-31", "content/new1.en.md", -5)
				writeNewContentFile(t, "new_en_2", "1989-07-30", "content/new2.en.md", -10)
				writeNewContentFile(t, "new_fr_1", "2016-07-30", "content/new1.fr.md", 10)
			},
			[]fsnotify.Event{
				{Name: "content/new1.en.md", Op: fsnotify.Create},
				{Name: "content/new2.en.md", Op: fsnotify.Create},
				{Name: "content/new1.fr.md", Op: fsnotify.Create},
			},
			func(t *testing.T) {
				assert.Len(t, enSite.Pages, 4)
				assert.Len(t, enSite.AllPages, 8)
				assert.Len(t, frSite.Pages, 4)
				assert.Equal(t, "new_fr_1", frSite.Pages[3].Title)
				assert.Equal(t, "new_en_2", enSite.Pages[0].Title)
				assert.Equal(t, "new_en_1", enSite.Pages[1].Title)

				rendered := readDestination(t, "public/en/new1/index.html")
				assert.True(t, strings.Contains(rendered, "new_en_1"), rendered)
			},
		},
		{
			func(t *testing.T) {
				p := "content/sect/doc1.en.md"
				doc1 := readSource(t, p)
				doc1 += "CHANGED"
				writeSource(t, p, doc1)
			},
			[]fsnotify.Event{{Name: "content/sect/doc1.en.md", Op: fsnotify.Write}},
			func(t *testing.T) {
				assert.Len(t, enSite.Pages, 4)
				doc1 := readDestination(t, "public/en/sect/doc1-slug/index.html")
				assert.True(t, strings.Contains(doc1, "CHANGED"), doc1)

			},
		},
		// Rename a file
		{
			func(t *testing.T) {
				if err := hugofs.Source().Rename("content/new1.en.md", "content/new1renamed.en.md"); err != nil {
					t.Fatalf("Rename failed: %s", err)
				}
			},
			[]fsnotify.Event{
				{Name: "content/new1renamed.en.md", Op: fsnotify.Rename},
				{Name: "content/new1.en.md", Op: fsnotify.Rename},
			},
			func(t *testing.T) {
				assert.Len(t, enSite.Pages, 4, "Rename")
				assert.Equal(t, "new_en_1", enSite.Pages[1].Title)
				rendered := readDestination(t, "public/en/new1renamed/index.html")
				assert.True(t, strings.Contains(rendered, "new_en_1"), rendered)
			}},
		{
			// Change a template
			func(t *testing.T) {
				template := "layouts/_default/single.html"
				templateContent := readSource(t, template)
				templateContent += "{{ print \"Template Changed\"}}"
				writeSource(t, template, templateContent)
			},
			[]fsnotify.Event{{Name: "layouts/_default/single.html", Op: fsnotify.Write}},
			func(t *testing.T) {
				assert.Len(t, enSite.Pages, 4)
				assert.Len(t, enSite.AllPages, 8)
				assert.Len(t, frSite.Pages, 4)
				doc1 := readDestination(t, "public/en/sect/doc1-slug/index.html")
				assert.True(t, strings.Contains(doc1, "Template Changed"), doc1)
			},
		},
		{
			// Change a language file
			func(t *testing.T) {
				languageFile := "i18n/fr.yaml"
				langContent := readSource(t, languageFile)
				langContent = strings.Replace(langContent, "Bonjour", "Salut", 1)
				writeSource(t, languageFile, langContent)
			},
			[]fsnotify.Event{{Name: "i18n/fr.yaml", Op: fsnotify.Write}},
			func(t *testing.T) {
				assert.Len(t, enSite.Pages, 4)
				assert.Len(t, enSite.AllPages, 8)
				assert.Len(t, frSite.Pages, 4)
				docEn := readDestination(t, "public/en/sect/doc1-slug/index.html")
				assert.True(t, strings.Contains(docEn, "Hello"), "No Hello")
				docFr := readDestination(t, "public/fr/sect/doc1/index.html")
				assert.True(t, strings.Contains(docFr, "Salut"), "No Salut")
			},
		},
	} {

		if this.preFunc != nil {
			this.preFunc(t)
		}
		err = sites.Rebuild(cfg, this.events...)

		if err != nil {
			t.Fatalf("[%d] Failed to rebuild sites: %s", i, err)
		}

		this.assertFunc(t)
	}

}

func createMultiTestSites(t *testing.T) *HugoSites {
	// General settings
	hugofs.InitMemFs()

	viper.Set("DefaultExtension", "html")
	viper.Set("baseurl", "http://example.com/blog")
	viper.Set("DisableSitemap", false)
	viper.Set("DisableRSS", false)
	viper.Set("RSSUri", "index.xml")
	viper.Set("Taxonomies", map[string]string{"tag": "tags"})
	viper.Set("Permalinks", map[string]string{"other": "/somewhere/else/:filename"})

	// Add some layouts
	if err := afero.WriteFile(hugofs.Source(),
		filepath.Join("layouts", "_default/single.html"),
		[]byte("Single: {{ .Title }}|{{ i18n \"hello\" }} {{ .Content }}"),
		0755); err != nil {
		t.Fatalf("Failed to write layout file: %s", err)
	}

	if err := afero.WriteFile(hugofs.Source(),
		filepath.Join("layouts", "_default/list.html"),
		[]byte("List: {{ .Title }}"),
		0755); err != nil {
		t.Fatalf("Failed to write layout file: %s", err)
	}

	if err := afero.WriteFile(hugofs.Source(),
		filepath.Join("layouts", "index.html"),
		[]byte("Home: {{ .Title }}|{{ .IsHome }}"),
		0755); err != nil {
		t.Fatalf("Failed to write layout file: %s", err)
	}

	// Add some language files
	if err := afero.WriteFile(hugofs.Source(),
		filepath.Join("i18n", "en.yaml"),
		[]byte(`
- id: hello
  translation: "Hello"
`),
		0755); err != nil {
		t.Fatalf("Failed to write language file: %s", err)
	}
	if err := afero.WriteFile(hugofs.Source(),
		filepath.Join("i18n", "fr.yaml"),
		[]byte(`
- id: hello
  translation: "Bonjour"
`),
		0755); err != nil {
		t.Fatalf("Failed to write language file: %s", err)
	}

	// Sources
	sources := []source.ByteSource{
		{filepath.FromSlash("sect/doc1.en.md"), []byte(`---
title: doc1
slug: doc1-slug
tags:
 - tag1
publishdate: "2000-01-01"
---
# doc1
*some content*
NOTE: slug should be used as URL
`)},
		{filepath.FromSlash("sect/doc1.fr.md"), []byte(`---
title: doc1
tags:
 - tag1
 - tag2
publishdate: "2000-01-04"
---
# doc1
*quelque contenu*
NOTE: should be in the 'en' Page's 'Translations' field.
NOTE: date is after "doc3"
`)},
		{filepath.FromSlash("sect/doc2.en.md"), []byte(`---
title: doc2
publishdate: "2000-01-02"
---
# doc2
*some content*
NOTE: without slug, "doc2" should be used, without ".en" as URL
`)},
		{filepath.FromSlash("sect/doc3.en.md"), []byte(`---
title: doc3
publishdate: "2000-01-03"
tags:
 - tag2
url: /superbob
---
# doc3
*some content*
NOTE: third 'en' doc, should trigger pagination on home page.
`)},
		{filepath.FromSlash("sect/doc4.md"), []byte(`---
title: doc4
tags:
 - tag1
publishdate: "2000-01-05"
---
# doc4
*du contenu francophone*
NOTE: should use the DefaultContentLanguage and mark this doc as 'fr'.
NOTE: doesn't have any corresponding translation in 'en'
`)},
		{filepath.FromSlash("other/doc5.fr.md"), []byte(`---
title: doc5
publishdate: "2000-01-06"
---
# doc5
*autre contenu francophone*
NOTE: should use the "permalinks" configuration with :filename
`)},
		// Add some for the stats
		{filepath.FromSlash("stats/expired.fr.md"), []byte(`---
title: expired
publishdate: "2000-01-06"
expiryDate: "2001-01-06"
---
# Expired
`)},
		{filepath.FromSlash("stats/future.fr.md"), []byte(`---
title: future
publishdate: "2100-01-06"
---
# Future
`)},
		{filepath.FromSlash("stats/expired.en.md"), []byte(`---
title: expired
publishdate: "2000-01-06"
expiryDate: "2001-01-06"
---
# Expired
`)},
		{filepath.FromSlash("stats/future.en.md"), []byte(`---
title: future
publishdate: "2100-01-06"
---
# Future
`)},
		{filepath.FromSlash("stats/draft.en.md"), []byte(`---
title: expired
publishdate: "2000-01-06"
draft: true
---
# Draft
`)},
	}

	// Multilingual settings
	viper.Set("Multilingual", true)
	en := NewLanguage("en")
	viper.Set("DefaultContentLanguage", "fr")
	viper.Set("paginate", "2")

	languages := NewLanguages(en, NewLanguage("fr"))

	// Hugo support using ByteSource's directly (for testing),
	// but to make it more real, we write them to the mem file system.
	for _, s := range sources {
		if err := afero.WriteFile(hugofs.Source(), filepath.Join("content", s.Name), s.Content, 0755); err != nil {
			t.Fatalf("Failed to write file: %s", err)
		}
	}
	_, err := hugofs.Source().Open("content/other/doc5.fr.md")

	if err != nil {
		t.Fatalf("Unable to locate file")
	}
	sites, err := newHugoSitesFromLanguages(languages)

	if err != nil {
		t.Fatalf("Failed to create sites: %s", err)
	}

	if len(sites.Sites) != 2 {
		t.Fatalf("Got %d sites", len(sites.Sites))
	}

	return sites
}

func writeSource(t *testing.T, filename, content string) {
	if err := afero.WriteFile(hugofs.Source(), filepath.FromSlash(filename), []byte(content), 0755); err != nil {
		t.Fatalf("Failed to write file: %s", err)
	}
}

func readDestination(t *testing.T, filename string) string {
	return readFileFromFs(t, hugofs.Destination(), filename)
}

func readSource(t *testing.T, filename string) string {
	return readFileFromFs(t, hugofs.Source(), filename)
}

func readFileFromFs(t *testing.T, fs afero.Fs, filename string) string {
	filename = filepath.FromSlash(filename)
	b, err := afero.ReadFile(fs, filename)
	if err != nil {
		// Print some debug info
		root := strings.Split(filename, helpers.FilePathSeparator)[0]
		afero.Walk(fs, root, func(path string, info os.FileInfo, err error) error {
			if info != nil && !info.IsDir() {
				fmt.Println("    ", path)
			}

			return nil
		})
		t.Fatalf("Failed to read file: %s", err)
	}
	return string(b)
}

const testPageTemplate = `---
title: "%s"
publishdate: "%s"
weight: %d
---
# Doc %s
`

func newTestPage(title, date string, weight int) string {
	return fmt.Sprintf(testPageTemplate, title, date, weight, title)
}

func writeNewContentFile(t *testing.T, title, date, filename string, weight int) {
	content := newTestPage(title, date, weight)
	writeSource(t, filename, content)
}
