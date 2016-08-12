package hugolib

import (
	"bytes"
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
	jww.SetStdoutThreshold(jww.LevelCritical)
}

func testCommonResetState() {
	hugofs.InitMemFs()
	viper.Reset()
	viper.SetFs(hugofs.Source())
	loadDefaultSettings()

	// Default is false, but true is easier to use as default in tests
	viper.Set("DefaultContentLanguageInSubdir", true)

	if err := hugofs.Source().Mkdir("content", 0755); err != nil {
		panic("Content folder creation failed.")
	}

}

func TestMultiSitesMainLangInRoot(t *testing.T) {
	for _, b := range []bool{false, true} {
		doTestMultiSitesMainLangInRoot(t, b)
	}
}

func doTestMultiSitesMainLangInRoot(t *testing.T, defaultInSubDir bool) {
	testCommonResetState()
	viper.Set("DefaultContentLanguageInSubdir", defaultInSubDir)

	sites := createMultiTestSites(t, multiSiteTomlConfig)

	err := sites.Build(BuildCfg{})

	if err != nil {
		t.Fatalf("Failed to build sites: %s", err)
	}

	require.Len(t, sites.Sites, 4)

	enSite := sites.Sites[0]
	frSite := sites.Sites[1]

	require.Equal(t, "/en", enSite.Info.LanguagePrefix)

	if defaultInSubDir {
		require.Equal(t, "/fr", frSite.Info.LanguagePrefix)
	} else {
		require.Equal(t, "", frSite.Info.LanguagePrefix)
	}

	doc1en := enSite.Pages[0]
	doc1fr := frSite.Pages[0]

	enPerm, _ := doc1en.Permalink()
	enRelPerm, _ := doc1en.RelPermalink()
	require.Equal(t, "http://example.com/blog/en/sect/doc1-slug/", enPerm)
	require.Equal(t, "/blog/en/sect/doc1-slug/", enRelPerm)

	frPerm, _ := doc1fr.Permalink()
	frRelPerm, _ := doc1fr.RelPermalink()
	// Main language in root
	require.Equal(t, replaceDefaultContentLanguageValue("http://example.com/blog/fr/sect/doc1/", defaultInSubDir), frPerm)
	require.Equal(t, replaceDefaultContentLanguageValue("/blog/fr/sect/doc1/", defaultInSubDir), frRelPerm)

	assertFileContent(t, "public/fr/sect/doc1/index.html", defaultInSubDir, "Single", "Bonjour")
	assertFileContent(t, "public/en/sect/doc1-slug/index.html", defaultInSubDir, "Single", "Hello")

	// Check home
	if defaultInSubDir {
		// should have a redirect on top level.
		assertFileContent(t, "public/index.html", true, `<meta http-equiv="refresh" content="0; url=http://example.com/blog/fr" />`)
	}
	assertFileContent(t, "public/fr/index.html", defaultInSubDir, "Home", "Bonjour")
	assertFileContent(t, "public/en/index.html", defaultInSubDir, "Home", "Hello")

	// Check list pages
	assertFileContent(t, "public/fr/sect/index.html", defaultInSubDir, "List", "Bonjour")
	assertFileContent(t, "public/en/sect/index.html", defaultInSubDir, "List", "Hello")
	assertFileContent(t, "public/fr/plaques/frtag1/index.html", defaultInSubDir, "List", "Bonjour")
	assertFileContent(t, "public/en/tags/tag1/index.html", defaultInSubDir, "List", "Hello")

	// Check sitemaps
	// Sitemaps behaves different: In a multilanguage setup there will always be a index file and
	// one sitemap in each lang folder.
	assertFileContent(t, "public/sitemap.xml", true,
		"<loc>http:/example.com/blog/en/sitemap.xml</loc>",
		"<loc>http:/example.com/blog/fr/sitemap.xml</loc>")

	if defaultInSubDir {
		assertFileContent(t, "public/fr/sitemap.xml", true, "<loc>http://example.com/blog/fr/</loc>")
	} else {
		assertFileContent(t, "public/fr/sitemap.xml", true, "<loc>http://example.com/blog/</loc>")
	}
	assertFileContent(t, "public/en/sitemap.xml", true, "<loc>http://example.com/blog/en/</loc>")

	// Check rss
	assertFileContent(t, "public/fr/index.xml", defaultInSubDir, `<atom:link href="http://example.com/blog/fr/index.xml"`)
	assertFileContent(t, "public/en/index.xml", defaultInSubDir, `<atom:link href="http://example.com/blog/en/index.xml"`)
	assertFileContent(t, "public/fr/sect/index.xml", defaultInSubDir, `<atom:link href="http://example.com/blog/fr/sect/index.xml"`)
	assertFileContent(t, "public/en/sect/index.xml", defaultInSubDir, `<atom:link href="http://example.com/blog/en/sect/index.xml"`)
	assertFileContent(t, "public/fr/plaques/frtag1/index.xml", defaultInSubDir, `<atom:link href="http://example.com/blog/fr/plaques/frtag1/index.xml"`)
	assertFileContent(t, "public/en/tags/tag1/index.xml", defaultInSubDir, `<atom:link href="http://example.com/blog/en/tags/tag1/index.xml"`)

	// Check paginators
	assertFileContent(t, "public/fr/page/1/index.html", defaultInSubDir, `refresh" content="0; url=http://example.com/blog/fr/"`)
	assertFileContent(t, "public/en/page/1/index.html", defaultInSubDir, `refresh" content="0; url=http://example.com/blog/en/"`)
	assertFileContent(t, "public/fr/page/2/index.html", defaultInSubDir, "Home Page 2", "Bonjour", "http://example.com/blog/fr/")
	assertFileContent(t, "public/en/page/2/index.html", defaultInSubDir, "Home Page 2", "Hello", "http://example.com/blog/en/")
	assertFileContent(t, "public/fr/sect/page/1/index.html", defaultInSubDir, `refresh" content="0; url=http://example.com/blog/fr/sect/"`)
	assertFileContent(t, "public/en/sect/page/1/index.html", defaultInSubDir, `refresh" content="0; url=http://example.com/blog/en/sect/"`)
	assertFileContent(t, "public/fr/sect/page/2/index.html", defaultInSubDir, "List Page 2", "Bonjour", "http://example.com/blog/fr/sect/")
	assertFileContent(t, "public/en/sect/page/2/index.html", defaultInSubDir, "List Page 2", "Hello", "http://example.com/blog/en/sect/")
	assertFileContent(t, "public/fr/plaques/frtag1/page/1/index.html", defaultInSubDir, `refresh" content="0; url=http://example.com/blog/fr/plaques/frtag1/"`)
	assertFileContent(t, "public/en/tags/tag1/page/1/index.html", defaultInSubDir, `refresh" content="0; url=http://example.com/blog/en/tags/tag1/"`)
	assertFileContent(t, "public/fr/plaques/frtag1/page/2/index.html", defaultInSubDir, "List Page 2", "Bonjour", "http://example.com/blog/fr/plaques/frtag1/")
	assertFileContent(t, "public/en/tags/tag1/page/2/index.html", defaultInSubDir, "List Page 2", "Hello", "http://example.com/blog/en/tags/tag1/")

}

func replaceDefaultContentLanguageValue(value string, defaultInSubDir bool) string {
	replace := viper.GetString("DefaultContentLanguage") + "/"
	if !defaultInSubDir {
		value = strings.Replace(value, replace, "", 1)

	}
	return value

}

func assertFileContent(t *testing.T, filename string, defaultInSubDir bool, matches ...string) {
	filename = replaceDefaultContentLanguageValue(filename, defaultInSubDir)
	content := readDestination(t, filename)
	for _, match := range matches {
		match = replaceDefaultContentLanguageValue(match, defaultInSubDir)
		require.True(t, strings.Contains(content, match), fmt.Sprintf("File no match for %q in %q: %s", match, filename, content))
	}
}

//
func TestMultiSitesBuild(t *testing.T) {
	for _, config := range []struct {
		content string
		suffix  string
	}{
		{multiSiteTomlConfig, "toml"},
		{multiSiteYAMLConfig, "yml"},
		{multiSiteJSONConfig, "json"},
	} {
		doTestMultiSitesBuild(t, config.content, config.suffix)
	}
}

func doTestMultiSitesBuild(t *testing.T, configContent, configSuffix string) {
	testCommonResetState()
	sites := createMultiTestSitesForConfig(t, configContent, configSuffix)

	err := sites.Build(BuildCfg{})

	if err != nil {
		t.Fatalf("Failed to build sites: %s", err)
	}

	enSite := sites.Sites[0]

	assert.Equal(t, "en", enSite.Language.Lang)

	if len(enSite.Pages) != 3 {
		t.Fatal("Expected 3 english pages")
	}
	assert.Len(t, enSite.Source.Files(), 13, "should have 13 source files")
	assert.Len(t, enSite.AllPages, 8, "should have 8 total pages (including translations)")

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

	assert.Equal(t, "/en/superbob", doc3.URL(), "invalid url, was specified on doc3")

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
	assert.Len(t, frSite.AllPages, 8, "should have 8 total pages (including translations)")

	for _, frenchPage := range frSite.Pages {
		assert.Equal(t, "fr", frenchPage.Lang())
	}

	// Check redirect to main language, French
	languageRedirect := readDestination(t, "public/index.html")
	require.True(t, strings.Contains(languageRedirect, "0; url=http://example.com/blog/fr"), languageRedirect)

	// check home page content (including data files rendering)
	assertFileContent(t, "public/en/index.html", true, "Home Page 1", "Hello", "Hugo Rocks!")
	assertFileContent(t, "public/fr/index.html", true, "Home Page 1", "Bonjour", "Hugo Rocks!")

	// check single page content
	assertFileContent(t, "public/fr/sect/doc1/index.html", true, "Single", "Shortcode: Bonjour")
	assertFileContent(t, "public/en/sect/doc1-slug/index.html", true, "Single", "Shortcode: Hello")

	// Check node translations
	homeEn := enSite.getNode("home-0")
	require.NotNil(t, homeEn)
	require.Len(t, homeEn.Translations(), 3)
	require.Equal(t, "fr", homeEn.Translations()[0].Lang())
	require.Equal(t, "nn", homeEn.Translations()[1].Lang())
	require.Equal(t, "Nynorsk", homeEn.Translations()[1].Title)
	require.Equal(t, "nb", homeEn.Translations()[2].Lang())
	require.Equal(t, "Bokmål", homeEn.Translations()[2].Title)

	sectFr := frSite.getNode("sect-sect-0")
	require.NotNil(t, sectFr)

	require.Equal(t, "fr", sectFr.Lang())
	require.Len(t, sectFr.Translations(), 1)
	require.Equal(t, "en", sectFr.Translations()[0].Lang())
	require.Equal(t, "Sects", sectFr.Translations()[0].Title)

	nnSite := sites.Sites[2]
	require.Equal(t, "nn", nnSite.Language.Lang)
	taxNn := nnSite.getNode("taxlist-lag-0")
	require.NotNil(t, taxNn)
	require.Len(t, taxNn.Translations(), 1)
	require.Equal(t, "nb", taxNn.Translations()[0].Lang())

	taxTermNn := nnSite.getNode("tax-lag-sogndal-0")
	require.NotNil(t, taxTermNn)
	require.Len(t, taxTermNn.Translations(), 1)
	require.Equal(t, "nb", taxTermNn.Translations()[0].Lang())

	// Check sitemap(s)
	sitemapIndex := readDestination(t, "public/sitemap.xml")
	require.True(t, strings.Contains(sitemapIndex, "<loc>http:/example.com/blog/en/sitemap.xml</loc>"), sitemapIndex)
	require.True(t, strings.Contains(sitemapIndex, "<loc>http:/example.com/blog/fr/sitemap.xml</loc>"), sitemapIndex)
	sitemapEn := readDestination(t, "public/en/sitemap.xml")
	sitemapFr := readDestination(t, "public/fr/sitemap.xml")
	require.True(t, strings.Contains(sitemapEn, "http://example.com/blog/en/sect/doc2/"), sitemapEn)
	require.True(t, strings.Contains(sitemapFr, "http://example.com/blog/fr/sect/doc1/"), sitemapFr)

	// Check taxonomies
	enTags := enSite.Taxonomies["tags"]
	frTags := frSite.Taxonomies["plaques"]
	require.Len(t, enTags, 2, fmt.Sprintf("Tags in en: %v", enTags))
	require.Len(t, frTags, 2, fmt.Sprintf("Tags in fr: %v", frTags))
	require.NotNil(t, enTags["tag1"])
	require.NotNil(t, frTags["frtag1"])
	readDestination(t, "public/fr/plaques/frtag1/index.html")
	readDestination(t, "public/en/tags/tag1/index.html")

	// Check Blackfriday config
	assert.True(t, strings.Contains(string(doc1fr.Content), "&laquo;"), string(doc1fr.Content))
	assert.False(t, strings.Contains(string(doc1en.Content), "&laquo;"), string(doc1en.Content))
	assert.True(t, strings.Contains(string(doc1en.Content), "&ldquo;"), string(doc1en.Content))

	// Check that the drafts etc. are not built/processed/rendered.
	assertShouldNotBuild(t, sites)

}

func TestMultiSitesRebuild(t *testing.T) {
	testCommonResetState()
	sites := createMultiTestSites(t, multiSiteTomlConfig)
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
				assert.Len(t, enSite.AllPages, 10)
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
				assert.Len(t, enSite.AllPages, 10)
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
				assert.Len(t, enSite.AllPages, 10)
				assert.Len(t, frSite.Pages, 4)
				docEn := readDestination(t, "public/en/sect/doc1-slug/index.html")
				assert.True(t, strings.Contains(docEn, "Hello"), "No Hello")
				docFr := readDestination(t, "public/fr/sect/doc1/index.html")
				assert.True(t, strings.Contains(docFr, "Salut"), "No Salut")

				homeEn := enSite.getNode("home-0")
				require.NotNil(t, homeEn)
				require.Len(t, homeEn.Translations(), 3)
				require.Equal(t, "fr", homeEn.Translations()[0].Lang())

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

	// Check that the drafts etc. are not built/processed/rendered.
	assertShouldNotBuild(t, sites)

}

func assertShouldNotBuild(t *testing.T, sites *HugoSites) {
	s := sites.Sites[0]

	for _, p := range s.rawAllPages {
		// No HTML when not processed
		require.Equal(t, p.shouldBuild(), bytes.Contains(p.rawContent, []byte("</")), p.BaseFileName()+": "+string(p.rawContent))
		require.Equal(t, p.shouldBuild(), p.Content != "", p.BaseFileName())

		require.Equal(t, p.shouldBuild(), p.Content != "", p.BaseFileName())

		filename := filepath.Join("public", p.TargetPath())
		if strings.HasSuffix(filename, ".html") {
			// TODO(bep) the end result is correct, but it is weird that we cannot use targetPath directly here.
			filename = strings.Replace(filename, ".html", "/index.html", 1)
		}

		require.Equal(t, p.shouldBuild(), destinationExists(filename), filename)
	}
}

func TestAddNewLanguage(t *testing.T) {
	testCommonResetState()

	sites := createMultiTestSites(t, multiSiteTomlConfig)
	cfg := BuildCfg{}

	err := sites.Build(cfg)

	if err != nil {
		t.Fatalf("Failed to build sites: %s", err)
	}

	newConfig := multiSiteTomlConfig + `

[Languages.sv]
weight = 15
title = "Svenska"
`

	writeNewContentFile(t, "Swedish Contentfile", "2016-01-01", "content/sect/doc1.sv.md", 10)
	// replace the config
	writeSource(t, "multilangconfig.toml", newConfig)

	// Watching does not work with in-memory fs, so we trigger a reload manually
	require.NoError(t, viper.ReadInConfig())
	err = sites.Build(BuildCfg{CreateSitesFromConfig: true})

	if err != nil {
		t.Fatalf("Failed to rebuild sites: %s", err)
	}

	require.Len(t, sites.Sites, 5, fmt.Sprintf("Len %d", len(sites.Sites)))

	// The Swedish site should be put in the middle (language weight=15)
	enSite := sites.Sites[0]
	svSite := sites.Sites[1]
	frSite := sites.Sites[2]
	require.True(t, enSite.Language.Lang == "en", enSite.Language.Lang)
	require.True(t, svSite.Language.Lang == "sv", svSite.Language.Lang)
	require.True(t, frSite.Language.Lang == "fr", frSite.Language.Lang)

	homeEn := enSite.getNode("home-0")
	require.NotNil(t, homeEn)
	require.Len(t, homeEn.Translations(), 4)
	require.Equal(t, "sv", homeEn.Translations()[0].Lang())

	require.Len(t, enSite.Pages, 3)
	require.Len(t, frSite.Pages, 3)

	// Veriy Swedish site
	require.Len(t, svSite.Pages, 1)
	svPage := svSite.Pages[0]
	require.Equal(t, "Swedish Contentfile", svPage.Title)
	require.Equal(t, "sv", svPage.Lang())
	require.Len(t, svPage.Translations(), 2)
	require.Len(t, svPage.AllTranslations(), 3)
	require.Equal(t, "en", svPage.Translations()[0].Lang())

}

var multiSiteTomlConfig = `
DefaultExtension = "html"
baseurl = "http://example.com/blog"
DisableSitemap = false
DisableRSS = false
RSSUri = "index.xml"

paginate = 1
DefaultContentLanguage = "fr"

[permalinks]
other = "/somewhere/else/:filename"

[blackfriday]
angledQuotes = true

[Taxonomies]
tag = "tags"

[Languages]
[Languages.en]
weight = 10
title = "English"
[Languages.en.blackfriday]
angledQuotes = false

[Languages.fr]
weight = 20
title = "Français"
[Languages.fr.Taxonomies]
plaque = "plaques"

[Languages.nn]
weight = 30
title = "Nynorsk"
[Languages.nn.Taxonomies]
lag = "lag"

[Languages.nb]
weight = 40
title = "Bokmål"
[Languages.nb.Taxonomies]
lag = "lag"
`

var multiSiteYAMLConfig = `
DefaultExtension: "html"
baseurl: "http://example.com/blog"
DisableSitemap: false
DisableRSS: false
RSSUri: "index.xml"

paginate: 1
DefaultContentLanguage: "fr"

permalinks:
    other: "/somewhere/else/:filename"

blackfriday:
    angledQuotes: true

Taxonomies:
    tag: "tags"

Languages:
    en:
        weight: 10
        title: "English"
        blackfriday:
            angledQuotes: false
    fr:
        weight: 20
        title: "Français"
        Taxonomies:
            plaque: "plaques"
    nn:
        weight: 30
        title: "Nynorsk"
        Taxonomies:
            lag: "lag"
    nb:
        weight: 40
        title: "Bokmål"
        Taxonomies:
            lag: "lag"

`

var multiSiteJSONConfig = `
{
  "DefaultExtension": "html",
  "baseurl": "http://example.com/blog",
  "DisableSitemap": false,
  "DisableRSS": false,
  "RSSUri": "index.xml",
  "paginate": 1,
  "DefaultContentLanguage": "fr",
  "permalinks": {
    "other": "/somewhere/else/:filename"
  },
  "blackfriday": {
    "angledQuotes": true
  },
  "Taxonomies": {
    "tag": "tags"
  },
  "Languages": {
    "en": {
      "weight": 10,
      "title": "English",
      "blackfriday": {
        "angledQuotes": false
      }
    },
    "fr": {
      "weight": 20,
      "title": "Français",
      "Taxonomies": {
        "plaque": "plaques"
      }
    },
    "nn": {
      "weight": 30,
      "title": "Nynorsk",
      "Taxonomies": {
        "lag": "lag"
      }
    },
    "nb": {
      "weight": 40,
      "title": "Bokmål",
      "Taxonomies": {
        "lag": "lag"
      }
    }
  }
}
`

func createMultiTestSites(t *testing.T, tomlConfig string) *HugoSites {
	return createMultiTestSitesForConfig(t, tomlConfig, "toml")
}

func createMultiTestSitesForConfig(t *testing.T, configContent, configSuffix string) *HugoSites {

	// Add some layouts
	if err := afero.WriteFile(hugofs.Source(),
		filepath.Join("layouts", "_default/single.html"),
		[]byte("Single: {{ .Title }}|{{ i18n \"hello\" }}|{{.Lang}}|{{ .Content }}"),
		0755); err != nil {
		t.Fatalf("Failed to write layout file: %s", err)
	}

	if err := afero.WriteFile(hugofs.Source(),
		filepath.Join("layouts", "_default/list.html"),
		[]byte("{{ $p := .Paginator }}List Page {{ $p.PageNumber }}: {{ .Title }}|{{ i18n \"hello\" }}|{{ .Permalink }}"),
		0755); err != nil {
		t.Fatalf("Failed to write layout file: %s", err)
	}

	if err := afero.WriteFile(hugofs.Source(),
		filepath.Join("layouts", "index.html"),
		[]byte("{{ $p := .Paginator }}Home Page {{ $p.PageNumber }}: {{ .Title }}|{{ .IsHome }}|{{ i18n \"hello\" }}|{{ .Permalink }}|{{  .Site.Data.hugo.slogan }}"),
		0755); err != nil {
		t.Fatalf("Failed to write layout file: %s", err)
	}

	// Add a shortcode
	if err := afero.WriteFile(hugofs.Source(),
		filepath.Join("layouts", "shortcodes", "shortcode.html"),
		[]byte("Shortcode: {{ i18n \"hello\" }}"),
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
*some "content"*

{{< shortcode >}}

NOTE: slug should be used as URL
`)},
		{filepath.FromSlash("sect/doc1.fr.md"), []byte(`---
title: doc1
plaques:
 - frtag1
 - frtag2
publishdate: "2000-01-04"
---
# doc1
*quelque "contenu"*

{{< shortcode >}}

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
 - tag1
url: /superbob
---
# doc3
*some content*
NOTE: third 'en' doc, should trigger pagination on home page.
`)},
		{filepath.FromSlash("sect/doc4.md"), []byte(`---
title: doc4
plaques:
 - frtag1
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
		{filepath.FromSlash("stats/tax.nn.md"), []byte(`---
title: Tax NN
publishdate: "2000-01-06"
weight: 1001
lag:
- Sogndal
---
# Tax NN
`)},
		{filepath.FromSlash("stats/tax.nb.md"), []byte(`---
title: Tax NB
publishdate: "2000-01-06"
weight: 1002
lag:
- Sogndal
---
# Tax NB
`)},
	}

	configFile := "multilangconfig." + configSuffix
	writeSource(t, configFile, configContent)
	if err := LoadGlobalConfig("", configFile); err != nil {
		t.Fatalf("Failed to load config: %s", err)
	}

	// Hugo support using ByteSource's directly (for testing),
	// but to make it more real, we write them to the mem file system.
	for _, s := range sources {
		if err := afero.WriteFile(hugofs.Source(), filepath.Join("content", s.Name), s.Content, 0755); err != nil {
			t.Fatalf("Failed to write file: %s", err)
		}
	}

	// Add some data
	writeSource(t, "data/hugo.toml", "slogan = \"Hugo Rocks!\"")

	sites, err := NewHugoSitesFromConfiguration()

	if err != nil {
		t.Fatalf("Failed to create sites: %s", err)
	}

	if len(sites.Sites) != 4 {
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

func destinationExists(filename string) bool {
	b, err := helpers.Exists(filename, hugofs.Destination())
	if err != nil {
		panic(err)
	}
	return b
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
