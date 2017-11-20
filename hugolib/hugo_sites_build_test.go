package hugolib

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"html/template"
	"os"
	"path/filepath"
	"time"

	"github.com/fortytw2/leaktest"
	"github.com/fsnotify/fsnotify"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/source"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

type testSiteConfig struct {
	DefaultContentLanguage         string
	DefaultContentLanguageInSubdir bool
	Fs                             afero.Fs
}

func TestMultiSitesMainLangInRoot(t *testing.T) {
	t.Parallel()
	for _, b := range []bool{true, false} {
		doTestMultiSitesMainLangInRoot(t, b)
	}
}

func doTestMultiSitesMainLangInRoot(t *testing.T, defaultInSubDir bool) {

	siteConfig := testSiteConfig{Fs: afero.NewMemMapFs(), DefaultContentLanguage: "fr", DefaultContentLanguageInSubdir: defaultInSubDir}

	sites := createMultiTestSites(t, siteConfig, multiSiteTOMLConfigTemplate)
	fs := sites.Fs
	th := testHelper{sites.Cfg, fs, t}

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

	require.Equal(t, "/blog/en/foo", enSite.PathSpec.RelURL("foo", true))

	doc1en := enSite.RegularPages[0]
	doc1fr := frSite.RegularPages[0]

	enPerm := doc1en.Permalink()
	enRelPerm := doc1en.RelPermalink()
	require.Equal(t, "http://example.com/blog/en/sect/doc1-slug/", enPerm)
	require.Equal(t, "/blog/en/sect/doc1-slug/", enRelPerm)

	frPerm := doc1fr.Permalink()
	frRelPerm := doc1fr.RelPermalink()
	// Main language in root
	require.Equal(t, th.replaceDefaultContentLanguageValue("http://example.com/blog/fr/sect/doc1/"), frPerm)
	require.Equal(t, th.replaceDefaultContentLanguageValue("/blog/fr/sect/doc1/"), frRelPerm)

	th.assertFileContent("public/fr/sect/doc1/index.html", "Single", "Bonjour")
	th.assertFileContent("public/en/sect/doc1-slug/index.html", "Single", "Hello")

	// Check home
	if defaultInSubDir {
		// should have a redirect on top level.
		th.assertFileContentStraight("public/index.html", `<meta http-equiv="refresh" content="0; url=http://example.com/blog/fr" />`)
	} else {
		// should have redirect back to root
		th.assertFileContentStraight("public/fr/index.html", `<meta http-equiv="refresh" content="0; url=http://example.com/blog" />`)
	}
	th.assertFileContent("public/fr/index.html", "Home", "Bonjour")
	th.assertFileContent("public/en/index.html", "Home", "Hello")

	// Check list pages
	th.assertFileContent("public/fr/sect/index.html", "List", "Bonjour")
	th.assertFileContent("public/en/sect/index.html", "List", "Hello")
	th.assertFileContent("public/fr/plaques/frtag1/index.html", "List", "Bonjour")
	th.assertFileContent("public/en/tags/tag1/index.html", "List", "Hello")

	// Check sitemaps
	// Sitemaps behaves different: In a multilanguage setup there will always be a index file and
	// one sitemap in each lang folder.
	th.assertFileContentStraight("public/sitemap.xml",
		"<loc>http://example.com/blog/en/sitemap.xml</loc>",
		"<loc>http://example.com/blog/fr/sitemap.xml</loc>")

	if defaultInSubDir {
		th.assertFileContentStraight("public/fr/sitemap.xml", "<loc>http://example.com/blog/fr/</loc>")
	} else {
		th.assertFileContentStraight("public/fr/sitemap.xml", "<loc>http://example.com/blog/</loc>")
	}
	th.assertFileContent("public/en/sitemap.xml", "<loc>http://example.com/blog/en/</loc>")

	// Check rss
	th.assertFileContent("public/fr/index.xml", `<atom:link href="http://example.com/blog/fr/index.xml"`,
		`rel="self" type="application/rss+xml"`)
	th.assertFileContent("public/en/index.xml", `<atom:link href="http://example.com/blog/en/index.xml"`)
	th.assertFileContent("public/fr/sect/index.xml", `<atom:link href="http://example.com/blog/fr/sect/index.xml"`)
	th.assertFileContent("public/en/sect/index.xml", `<atom:link href="http://example.com/blog/en/sect/index.xml"`)
	th.assertFileContent("public/fr/plaques/frtag1/index.xml", `<atom:link href="http://example.com/blog/fr/plaques/frtag1/index.xml"`)
	th.assertFileContent("public/en/tags/tag1/index.xml", `<atom:link href="http://example.com/blog/en/tags/tag1/index.xml"`)

	// Check paginators
	th.assertFileContent("public/fr/page/1/index.html", `refresh" content="0; url=http://example.com/blog/fr/"`)
	th.assertFileContent("public/en/page/1/index.html", `refresh" content="0; url=http://example.com/blog/en/"`)
	th.assertFileContent("public/fr/page/2/index.html", "Home Page 2", "Bonjour", "http://example.com/blog/fr/")
	th.assertFileContent("public/en/page/2/index.html", "Home Page 2", "Hello", "http://example.com/blog/en/")
	th.assertFileContent("public/fr/sect/page/1/index.html", `refresh" content="0; url=http://example.com/blog/fr/sect/"`)
	th.assertFileContent("public/en/sect/page/1/index.html", `refresh" content="0; url=http://example.com/blog/en/sect/"`)
	th.assertFileContent("public/fr/sect/page/2/index.html", "List Page 2", "Bonjour", "http://example.com/blog/fr/sect/")
	th.assertFileContent("public/en/sect/page/2/index.html", "List Page 2", "Hello", "http://example.com/blog/en/sect/")
	th.assertFileContent("public/fr/plaques/frtag1/page/1/index.html", `refresh" content="0; url=http://example.com/blog/fr/plaques/frtag1/"`)
	th.assertFileContent("public/en/tags/tag1/page/1/index.html", `refresh" content="0; url=http://example.com/blog/en/tags/tag1/"`)
	th.assertFileContent("public/fr/plaques/frtag1/page/2/index.html", "List Page 2", "Bonjour", "http://example.com/blog/fr/plaques/frtag1/")
	th.assertFileContent("public/en/tags/tag1/page/2/index.html", "List Page 2", "Hello", "http://example.com/blog/en/tags/tag1/")
	// nn (Nynorsk) and nb (Bokmål) have custom pagePath: side ("page" in Norwegian)
	th.assertFileContent("public/nn/side/1/index.html", `refresh" content="0; url=http://example.com/blog/nn/"`)
	th.assertFileContent("public/nb/side/1/index.html", `refresh" content="0; url=http://example.com/blog/nb/"`)
}

func TestMultiSitesWithTwoLanguages(t *testing.T) {
	t.Parallel()
	mm := afero.NewMemMapFs()

	writeToFs(t, mm, "config.toml", `

defaultContentLanguage = "nn"

[languages]
[languages.nn]
languageName = "Nynorsk"
weight = 1
title = "Tittel på Nynorsk"

[languages.en]
title = "Title in English"
languageName = "English"
weight = 2
`,
	)

	cfg, err := LoadConfig(mm, "", "config.toml")
	require.NoError(t, err)

	fs := hugofs.NewFrom(mm, cfg)

	sites, err := NewHugoSites(deps.DepsCfg{Fs: fs, Cfg: cfg})

	if err != nil {
		t.Fatalf("Failed to create sites: %s", err)
	}

	writeSource(t, fs, filepath.Join("content", "foo.md"), "foo")

	// Add some data
	writeSource(t, fs, filepath.Join("data", "hugo.toml"), "slogan = \"Hugo Rocks!\"")

	require.NoError(t, sites.Build(BuildCfg{}))
	require.Len(t, sites.Sites, 2)

	nnSite := sites.Sites[0]
	nnSiteHome := nnSite.getPage(KindHome)
	require.Len(t, nnSiteHome.AllTranslations(), 2)
	require.Len(t, nnSiteHome.Translations(), 1)
	require.True(t, nnSiteHome.IsTranslated())

}

//
func TestMultiSitesBuild(t *testing.T) {
	t.Parallel()

	for _, config := range []struct {
		content string
		suffix  string
	}{
		{multiSiteTOMLConfigTemplate, "toml"},
		{multiSiteYAMLConfigTemplate, "yml"},
		{multiSiteJSONConfigTemplate, "json"},
	} {
		doTestMultiSitesBuild(t, config.content, config.suffix)
	}
}

func doTestMultiSitesBuild(t *testing.T, configTemplate, configSuffix string) {
	siteConfig := testSiteConfig{Fs: afero.NewMemMapFs(), DefaultContentLanguage: "fr", DefaultContentLanguageInSubdir: true}
	sites := createMultiTestSitesForConfig(t, siteConfig, configTemplate, configSuffix)

	require.Len(t, sites.Sites, 4)

	fs := sites.Fs
	th := testHelper{sites.Cfg, fs, t}
	err := sites.Build(BuildCfg{})

	if err != nil {
		t.Fatalf("Failed to build sites: %s", err)
	}

	// Check site config
	for _, s := range sites.Sites {
		require.True(t, s.Info.defaultContentLanguageInSubdir, s.Info.Title)
		require.NotNil(t, s.disabledKinds)
	}

	gp1 := sites.GetContentPage(filepath.FromSlash("content/sect/doc1.en.md"))
	require.NotNil(t, gp1)
	require.Equal(t, "doc1", gp1.Title)
	gp2 := sites.GetContentPage(filepath.FromSlash("content/sect/notfound.md"))
	require.Nil(t, gp2)

	enSite := sites.Sites[0]
	enSiteHome := enSite.getPage(KindHome)
	require.True(t, enSiteHome.IsTranslated())

	require.Equal(t, "en", enSite.Language.Lang)

	if len(enSite.RegularPages) != 4 {
		t.Fatal("Expected 4 english pages")
	}
	require.Len(t, enSite.Source.Files(), 14, "should have 13 source files")
	require.Len(t, enSite.AllPages, 28, "should have 28 total pages (including translations and index types)")

	doc1en := enSite.RegularPages[0]
	permalink := doc1en.Permalink()
	require.NoError(t, err, "permalink call failed")
	require.Equal(t, "http://example.com/blog/en/sect/doc1-slug/", permalink, "invalid doc1.en permalink")
	require.Len(t, doc1en.Translations(), 1, "doc1-en should have one translation, excluding itself")

	doc2 := enSite.RegularPages[1]
	permalink = doc2.Permalink()
	require.NoError(t, err, "permalink call failed")
	require.Equal(t, "http://example.com/blog/en/sect/doc2/", permalink, "invalid doc2 permalink")

	doc3 := enSite.RegularPages[2]
	permalink = doc3.Permalink()
	// Note that /superbob is a custom URL set in frontmatter.
	// We respect that URL literally (it can be /search.json)
	// and do no not do any language code prefixing.
	require.Equal(t, "http://example.com/blog/superbob/", permalink, "invalid doc3 permalink")

	require.Equal(t, "/superbob", doc3.URL(), "invalid url, was specified on doc3")
	th.assertFileContent("public/superbob/index.html", "doc3|Hello|en")
	require.Equal(t, doc2.Next, doc3, "doc3 should follow doc2, in .Next")

	doc1fr := doc1en.Translations()[0]
	permalink = doc1fr.Permalink()
	require.NoError(t, err, "permalink call failed")
	require.Equal(t, "http://example.com/blog/fr/sect/doc1/", permalink, "invalid doc1fr permalink")

	require.Equal(t, doc1en.Translations()[0], doc1fr, "doc1-en should have doc1-fr as translation")
	require.Equal(t, doc1fr.Translations()[0], doc1en, "doc1-fr should have doc1-en as translation")
	require.Equal(t, "fr", doc1fr.Language().Lang)

	doc4 := enSite.AllPages[4]
	permalink = doc4.Permalink()
	require.Equal(t, "http://example.com/blog/fr/sect/doc4/", permalink, "invalid doc4 permalink")
	require.Equal(t, "/blog/fr/sect/doc4/", doc4.URL())

	require.Len(t, doc4.Translations(), 0, "found translations for doc4")

	doc5 := enSite.AllPages[5]
	permalink = doc5.Permalink()
	require.Equal(t, "http://example.com/blog/fr/somewhere/else/doc5/", permalink, "invalid doc5 permalink")

	// Taxonomies and their URLs
	require.Len(t, enSite.Taxonomies, 1, "should have 1 taxonomy")
	tags := enSite.Taxonomies["tags"]
	require.Len(t, tags, 2, "should have 2 different tags")
	require.Equal(t, tags["tag1"][0].Page, doc1en, "first tag1 page should be doc1")

	frSite := sites.Sites[1]

	require.Equal(t, "fr", frSite.Language.Lang)
	require.Len(t, frSite.RegularPages, 3, "should have 3 pages")
	require.Len(t, frSite.AllPages, 28, "should have 28 total pages (including translations and nodes)")

	for _, frenchPage := range frSite.RegularPages {
		require.Equal(t, "fr", frenchPage.Lang())
	}

	// Check redirect to main language, French
	languageRedirect := readDestination(t, fs, "public/index.html")
	require.True(t, strings.Contains(languageRedirect, "0; url=http://example.com/blog/fr"), languageRedirect)

	// check home page content (including data files rendering)
	th.assertFileContent("public/en/index.html", "Default Home Page 1", "Hello", "Hugo Rocks!")
	th.assertFileContent("public/fr/index.html", "French Home Page 1", "Bonjour", "Hugo Rocks!")

	// check single page content
	th.assertFileContent("public/fr/sect/doc1/index.html", "Single", "Shortcode: Bonjour", "LingoFrench")
	th.assertFileContent("public/en/sect/doc1-slug/index.html", "Single", "Shortcode: Hello", "LingoDefault")

	// Check node translations
	homeEn := enSite.getPage(KindHome)
	require.NotNil(t, homeEn)
	require.Len(t, homeEn.Translations(), 3)
	require.Equal(t, "fr", homeEn.Translations()[0].Lang())
	require.Equal(t, "nn", homeEn.Translations()[1].Lang())
	require.Equal(t, "På nynorsk", homeEn.Translations()[1].Title)
	require.Equal(t, "nb", homeEn.Translations()[2].Lang())
	require.Equal(t, "På bokmål", homeEn.Translations()[2].Title, configSuffix)
	require.Equal(t, "Bokmål", homeEn.Translations()[2].Language().LanguageName, configSuffix)

	sectFr := frSite.getPage(KindSection, "sect")
	require.NotNil(t, sectFr)

	require.Equal(t, "fr", sectFr.Lang())
	require.Len(t, sectFr.Translations(), 1)
	require.Equal(t, "en", sectFr.Translations()[0].Lang())
	require.Equal(t, "Sects", sectFr.Translations()[0].Title)

	nnSite := sites.Sites[2]
	require.Equal(t, "nn", nnSite.Language.Lang)
	taxNn := nnSite.getPage(KindTaxonomyTerm, "lag")
	require.NotNil(t, taxNn)
	require.Len(t, taxNn.Translations(), 1)
	require.Equal(t, "nb", taxNn.Translations()[0].Lang())

	taxTermNn := nnSite.getPage(KindTaxonomy, "lag", "sogndal")
	require.NotNil(t, taxTermNn)
	require.Len(t, taxTermNn.Translations(), 1)
	require.Equal(t, "nb", taxTermNn.Translations()[0].Lang())

	// Check sitemap(s)
	sitemapIndex := readDestination(t, fs, "public/sitemap.xml")
	require.True(t, strings.Contains(sitemapIndex, "<loc>http://example.com/blog/en/sitemap.xml</loc>"), sitemapIndex)
	require.True(t, strings.Contains(sitemapIndex, "<loc>http://example.com/blog/fr/sitemap.xml</loc>"), sitemapIndex)
	sitemapEn := readDestination(t, fs, "public/en/sitemap.xml")
	sitemapFr := readDestination(t, fs, "public/fr/sitemap.xml")
	require.True(t, strings.Contains(sitemapEn, "http://example.com/blog/en/sect/doc2/"), sitemapEn)
	require.True(t, strings.Contains(sitemapFr, "http://example.com/blog/fr/sect/doc1/"), sitemapFr)

	// Check taxonomies
	enTags := enSite.Taxonomies["tags"]
	frTags := frSite.Taxonomies["plaques"]
	require.Len(t, enTags, 2, fmt.Sprintf("Tags in en: %v", enTags))
	require.Len(t, frTags, 2, fmt.Sprintf("Tags in fr: %v", frTags))
	require.NotNil(t, enTags["tag1"])
	require.NotNil(t, frTags["frtag1"])
	readDestination(t, fs, "public/fr/plaques/frtag1/index.html")
	readDestination(t, fs, "public/en/tags/tag1/index.html")

	// Check Blackfriday config
	require.True(t, strings.Contains(string(doc1fr.Content), "&laquo;"), string(doc1fr.Content))
	require.False(t, strings.Contains(string(doc1en.Content), "&laquo;"), string(doc1en.Content))
	require.True(t, strings.Contains(string(doc1en.Content), "&ldquo;"), string(doc1en.Content))

	// Check that the drafts etc. are not built/processed/rendered.
	assertShouldNotBuild(t, sites)

	// en and nn have custom site menus
	require.Len(t, frSite.Menus, 0, "fr: "+configSuffix)
	require.Len(t, enSite.Menus, 1, "en: "+configSuffix)
	require.Len(t, nnSite.Menus, 1, "nn: "+configSuffix)

	require.Equal(t, "Home", enSite.Menus["main"].ByName()[0].Name)
	require.Equal(t, "Heim", nnSite.Menus["main"].ByName()[0].Name)

	// Issue #1302
	require.Equal(t, template.URL(""), enSite.RegularPages[0].RSSLink())

	// Issue #3108
	next := enSite.RegularPages[0].Next
	require.NotNil(t, next)
	require.Equal(t, KindPage, next.Kind)

	for {
		if next == nil {
			break
		}
		require.Equal(t, KindPage, next.Kind)
		next = next.Next
	}

}

func TestMultiSitesRebuild(t *testing.T) {
	// t.Parallel() not supported, see https://github.com/fortytw2/leaktest/issues/4
	// This leaktest seems to be a little bit shaky on Travis.
	if !isCI() {
		defer leaktest.CheckTimeout(t, 30*time.Second)()
	}
	siteConfig := testSiteConfig{Fs: afero.NewMemMapFs(), DefaultContentLanguage: "fr", DefaultContentLanguageInSubdir: true}
	sites := createMultiTestSites(t, siteConfig, multiSiteTOMLConfigTemplate)
	fs := sites.Fs
	cfg := BuildCfg{Watching: true}
	th := testHelper{sites.Cfg, fs, t}

	err := sites.Build(cfg)

	if err != nil {
		t.Fatalf("Failed to build sites: %s", err)
	}

	_, err = fs.Destination.Open("public/en/sect/doc2/index.html")

	if err != nil {
		t.Fatalf("Unable to locate file")
	}

	enSite := sites.Sites[0]
	frSite := sites.Sites[1]

	require.Len(t, enSite.RegularPages, 4)
	require.Len(t, frSite.RegularPages, 3)

	// Verify translations
	th.assertFileContent("public/en/sect/doc1-slug/index.html", "Hello")
	th.assertFileContent("public/fr/sect/doc1/index.html", "Bonjour")

	// check single page content
	th.assertFileContent("public/fr/sect/doc1/index.html", "Single", "Shortcode: Bonjour")
	th.assertFileContent("public/en/sect/doc1-slug/index.html", "Single", "Shortcode: Hello")

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
				require.Len(t, enSite.RegularPages, 3, "1 en removed")

				// Check build stats
				require.Equal(t, 1, enSite.draftCount, "Draft")
				require.Equal(t, 1, enSite.futureCount, "Future")
				require.Equal(t, 1, enSite.expiredCount, "Expired")
				require.Equal(t, 0, frSite.draftCount, "Draft")
				require.Equal(t, 1, frSite.futureCount, "Future")
				require.Equal(t, 1, frSite.expiredCount, "Expired")
			},
		},
		{
			func(t *testing.T) {
				writeNewContentFile(t, fs, "new_en_1", "2016-07-31", "content/new1.en.md", -5)
				writeNewContentFile(t, fs, "new_en_2", "1989-07-30", "content/new2.en.md", -10)
				writeNewContentFile(t, fs, "new_fr_1", "2016-07-30", "content/new1.fr.md", 10)
			},
			[]fsnotify.Event{
				{Name: "content/new1.en.md", Op: fsnotify.Create},
				{Name: "content/new2.en.md", Op: fsnotify.Create},
				{Name: "content/new1.fr.md", Op: fsnotify.Create},
			},
			func(t *testing.T) {
				require.Len(t, enSite.RegularPages, 5)
				require.Len(t, enSite.AllPages, 30)
				require.Len(t, frSite.RegularPages, 4)
				require.Equal(t, "new_fr_1", frSite.RegularPages[3].Title)
				require.Equal(t, "new_en_2", enSite.RegularPages[0].Title)
				require.Equal(t, "new_en_1", enSite.RegularPages[1].Title)

				rendered := readDestination(t, fs, "public/en/new1/index.html")
				require.True(t, strings.Contains(rendered, "new_en_1"), rendered)
			},
		},
		{
			func(t *testing.T) {
				p := "content/sect/doc1.en.md"
				doc1 := readSource(t, fs, p)
				doc1 += "CHANGED"
				writeSource(t, fs, p, doc1)
			},
			[]fsnotify.Event{{Name: "content/sect/doc1.en.md", Op: fsnotify.Write}},
			func(t *testing.T) {
				require.Len(t, enSite.RegularPages, 5)
				doc1 := readDestination(t, fs, "public/en/sect/doc1-slug/index.html")
				require.True(t, strings.Contains(doc1, "CHANGED"), doc1)

			},
		},
		// Rename a file
		{
			func(t *testing.T) {
				if err := fs.Source.Rename("content/new1.en.md", "content/new1renamed.en.md"); err != nil {
					t.Fatalf("Rename failed: %s", err)
				}
			},
			[]fsnotify.Event{
				{Name: "content/new1renamed.en.md", Op: fsnotify.Rename},
				{Name: "content/new1.en.md", Op: fsnotify.Rename},
			},
			func(t *testing.T) {
				require.Len(t, enSite.RegularPages, 5, "Rename")
				require.Equal(t, "new_en_1", enSite.RegularPages[1].Title)
				rendered := readDestination(t, fs, "public/en/new1renamed/index.html")
				require.True(t, strings.Contains(rendered, "new_en_1"), rendered)
			}},
		{
			// Change a template
			func(t *testing.T) {
				template := "layouts/_default/single.html"
				templateContent := readSource(t, fs, template)
				templateContent += "{{ print \"Template Changed\"}}"
				writeSource(t, fs, template, templateContent)
			},
			[]fsnotify.Event{{Name: "layouts/_default/single.html", Op: fsnotify.Write}},
			func(t *testing.T) {
				require.Len(t, enSite.RegularPages, 5)
				require.Len(t, enSite.AllPages, 30)
				require.Len(t, frSite.RegularPages, 4)
				doc1 := readDestination(t, fs, "public/en/sect/doc1-slug/index.html")
				require.True(t, strings.Contains(doc1, "Template Changed"), doc1)
			},
		},
		{
			// Change a language file
			func(t *testing.T) {
				languageFile := "i18n/fr.yaml"
				langContent := readSource(t, fs, languageFile)
				langContent = strings.Replace(langContent, "Bonjour", "Salut", 1)
				writeSource(t, fs, languageFile, langContent)
			},
			[]fsnotify.Event{{Name: "i18n/fr.yaml", Op: fsnotify.Write}},
			func(t *testing.T) {
				require.Len(t, enSite.RegularPages, 5)
				require.Len(t, enSite.AllPages, 30)
				require.Len(t, frSite.RegularPages, 4)
				docEn := readDestination(t, fs, "public/en/sect/doc1-slug/index.html")
				require.True(t, strings.Contains(docEn, "Hello"), "No Hello")
				docFr := readDestination(t, fs, "public/fr/sect/doc1/index.html")
				require.True(t, strings.Contains(docFr, "Salut"), "No Salut")

				homeEn := enSite.getPage(KindHome)
				require.NotNil(t, homeEn)
				require.Len(t, homeEn.Translations(), 3)
				require.Equal(t, "fr", homeEn.Translations()[0].Lang())

			},
		},
		// Change a shortcode
		{
			func(t *testing.T) {
				writeSource(t, fs, "layouts/shortcodes/shortcode.html", "Modified Shortcode: {{ i18n \"hello\" }}")
			},
			[]fsnotify.Event{
				{Name: "layouts/shortcodes/shortcode.html", Op: fsnotify.Write},
			},
			func(t *testing.T) {
				require.Len(t, enSite.RegularPages, 5)
				require.Len(t, enSite.AllPages, 30)
				require.Len(t, frSite.RegularPages, 4)
				th.assertFileContent("public/fr/sect/doc1/index.html", "Single", "Modified Shortcode: Salut")
				th.assertFileContent("public/en/sect/doc1-slug/index.html", "Single", "Modified Shortcode: Hello")
			},
		},
	} {

		if this.preFunc != nil {
			this.preFunc(t)
		}

		err = sites.Build(cfg, this.events...)

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
		require.Equal(t, p.shouldBuild(), bytes.Contains(p.workContent, []byte("</")), p.BaseFileName()+": "+string(p.workContent))
		require.Equal(t, p.shouldBuild(), p.Content != "", p.BaseFileName())

		require.Equal(t, p.shouldBuild(), p.Content != "", p.BaseFileName())

	}
}

func TestAddNewLanguage(t *testing.T) {
	t.Parallel()
	siteConfig := testSiteConfig{Fs: afero.NewMemMapFs(), DefaultContentLanguage: "fr", DefaultContentLanguageInSubdir: true}

	sites := createMultiTestSites(t, siteConfig, multiSiteTOMLConfigTemplate)
	cfg := BuildCfg{}

	err := sites.Build(cfg)

	if err != nil {
		t.Fatalf("Failed to build sites: %s", err)
	}

	fs := sites.Fs

	newConfig := multiSiteTOMLConfigTemplate + `

[Languages.sv]
weight = 15
title = "Svenska"
`

	newConfig = createConfig(t, siteConfig, newConfig)

	writeNewContentFile(t, fs, "Swedish Contentfile", "2016-01-01", "content/sect/doc1.sv.md", 10)
	// replace the config
	writeSource(t, fs, "multilangconfig.toml", newConfig)

	// Watching does not work with in-memory fs, so we trigger a reload manually
	require.NoError(t, sites.Cfg.(*helpers.Language).Cfg.(*viper.Viper).ReadInConfig())
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

	homeEn := enSite.getPage(KindHome)
	require.NotNil(t, homeEn)
	require.Len(t, homeEn.Translations(), 4)
	require.Equal(t, "sv", homeEn.Translations()[0].Lang())

	require.Len(t, enSite.RegularPages, 4)
	require.Len(t, frSite.RegularPages, 3)

	// Veriy Swedish site
	require.Len(t, svSite.RegularPages, 1)
	svPage := svSite.RegularPages[0]
	require.Equal(t, "Swedish Contentfile", svPage.Title)
	require.Equal(t, "sv", svPage.Lang())
	require.Len(t, svPage.Translations(), 2)
	require.Len(t, svPage.AllTranslations(), 3)
	require.Equal(t, "en", svPage.Translations()[0].Lang())

	// Regular pages have no children
	require.Len(t, svPage.Pages, 0)
	require.Len(t, svPage.Data["Pages"], 0)

}

func TestChangeDefaultLanguage(t *testing.T) {
	t.Parallel()
	mf := afero.NewMemMapFs()

	sites := createMultiTestSites(t, testSiteConfig{Fs: mf, DefaultContentLanguage: "fr", DefaultContentLanguageInSubdir: false}, multiSiteTOMLConfigTemplate)

	require.Equal(t, mf, sites.Fs.Source)

	cfg := BuildCfg{}
	fs := sites.Fs
	th := testHelper{sites.Cfg, fs, t}

	err := sites.Build(cfg)

	if err != nil {
		t.Fatalf("Failed to build sites: %s", err)
	}

	th.assertFileContent("public/sect/doc1/index.html", "Single", "Bonjour")
	th.assertFileContent("public/en/sect/doc2/index.html", "Single", "Hello")

	newConfig := createConfig(t, testSiteConfig{Fs: mf, DefaultContentLanguage: "en", DefaultContentLanguageInSubdir: false}, multiSiteTOMLConfigTemplate)

	// replace the config
	writeSource(t, fs, "multilangconfig.toml", newConfig)

	// Watching does not work with in-memory fs, so we trigger a reload manually
	// This does not look pretty, so we should think of something else.
	require.NoError(t, th.Cfg.(*helpers.Language).Cfg.(*viper.Viper).ReadInConfig())
	err = sites.Build(BuildCfg{CreateSitesFromConfig: true})

	if err != nil {
		t.Fatalf("Failed to rebuild sites: %s", err)
	}

	// Default language is now en, so that should now be the "root" language
	th.assertFileContent("public/fr/sect/doc1/index.html", "Single", "Bonjour")
	th.assertFileContent("public/sect/doc2/index.html", "Single", "Hello")
}

func TestTableOfContentsInShortcodes(t *testing.T) {
	t.Parallel()
	mf := afero.NewMemMapFs()

	writeToFs(t, mf, "layouts/shortcodes/toc.html", tocShortcode)
	writeToFs(t, mf, "content/post/simple.en.md", tocPageSimple)
	writeToFs(t, mf, "content/post/withSCInHeading.en.md", tocPageWithShortcodesInHeadings)

	sites := createMultiTestSites(t, testSiteConfig{Fs: mf, DefaultContentLanguage: "en", DefaultContentLanguageInSubdir: true}, multiSiteTOMLConfigTemplate)

	cfg := BuildCfg{}

	err := sites.Build(cfg)

	if err != nil {
		t.Fatalf("Failed to build sites: %s", err)
	}

	fs := sites.Fs
	th := testHelper{sites.Cfg, fs, t}

	th.assertFileContent("public/en/post/simple/index.html", tocPageSimpleExpected)
	th.assertFileContent("public/en/post/withSCInHeading/index.html", tocPageWithShortcodesInHeadingsExpected)
}

var tocShortcode = `
{{ .Page.TableOfContents }}
`

var tocPageSimple = `---
title: tocTest
publishdate: "2000-01-01"
---

{{< toc >}}

# Heading 1 {#1}

Some text.

## Subheading 1.1 {#1-1}

Some more text.

# Heading 2 {#2}

Even more text.

## Subheading 2.1 {#2-1}

Lorem ipsum...
`

var tocPageSimpleExpected = `<nav id="TableOfContents">
<ul>
<li><a href="#1">Heading 1</a>
<ul>
<li><a href="#1-1">Subheading 1.1</a></li>
</ul></li>
<li><a href="#2">Heading 2</a>
<ul>
<li><a href="#2-1">Subheading 2.1</a></li>
</ul></li>
</ul>
</nav>`

var tocPageWithShortcodesInHeadings = `---
title: tocTest
publishdate: "2000-01-01"
---

{{< toc >}}

# Heading 1 {#1}

Some text.

## Subheading 1.1 {{< shortcode >}} {#1-1}

Some more text.

# Heading 2 {{% shortcode %}} {#2}

Even more text.

## Subheading 2.1 {#2-1}

Lorem ipsum...
`

var tocPageWithShortcodesInHeadingsExpected = `<nav id="TableOfContents">
<ul>
<li><a href="#1">Heading 1</a>
<ul>
<li><a href="#1-1">Subheading 1.1 Shortcode: Hello</a></li>
</ul></li>
<li><a href="#2">Heading 2 Shortcode: Hello</a>
<ul>
<li><a href="#2-1">Subheading 2.1</a></li>
</ul></li>
</ul>
</nav>`

var multiSiteTOMLConfigTemplate = `
baseURL = "http://example.com/blog"
disableSitemap = false
disableRSS = false
rssURI = "index.xml"

paginate = 1
disablePathToLower = true
defaultContentLanguage = "{{ .DefaultContentLanguage }}"
defaultContentLanguageInSubdir = {{ .DefaultContentLanguageInSubdir }}

[permalinks]
other = "/somewhere/else/:filename"

[blackfriday]
angledQuotes = true

[Taxonomies]
tag = "tags"

[Languages]
[Languages.en]
weight = 10
title = "In English"
languageName = "English"
[Languages.en.blackfriday]
angledQuotes = false
[[Languages.en.menu.main]]
url    = "/"
name   = "Home"
weight = 0

[Languages.fr]
weight = 20
title = "Le Français"
languageName = "Français"
[Languages.fr.Taxonomies]
plaque = "plaques"

[Languages.nn]
weight = 30
title = "På nynorsk"
languageName = "Nynorsk"
paginatePath = "side"
[Languages.nn.Taxonomies]
lag = "lag"
[[Languages.nn.menu.main]]
url    = "/"
name   = "Heim"
weight = 1

[Languages.nb]
weight = 40
title = "På bokmål"
languageName = "Bokmål"
paginatePath = "side"
[Languages.nb.Taxonomies]
lag = "lag"
`

var multiSiteYAMLConfigTemplate = `
baseURL: "http://example.com/blog"
disableSitemap: false
disableRSS: false
rssURI: "index.xml"

disablePathToLower: true
paginate: 1
defaultContentLanguage: "{{ .DefaultContentLanguage }}"
defaultContentLanguageInSubdir: {{ .DefaultContentLanguageInSubdir }}

permalinks:
    other: "/somewhere/else/:filename"

blackfriday:
    angledQuotes: true

Taxonomies:
    tag: "tags"

Languages:
    en:
        weight: 10
        title: "In English"
        languageName: "English"
        blackfriday:
            angledQuotes: false
        menu:
            main:
                - url: "/"
                  name: "Home"
                  weight: 0
    fr:
        weight: 20
        title: "Le Français"
        languageName: "Français"
        Taxonomies:
            plaque: "plaques"
    nn:
        weight: 30
        title: "På nynorsk"
        languageName: "Nynorsk"
        paginatePath: "side"
        Taxonomies:
            lag: "lag"
        menu:
            main:
                - url: "/"
                  name: "Heim"
                  weight: 1
    nb:
        weight: 40
        title: "På bokmål"
        languageName: "Bokmål"
        paginatePath: "side"
        Taxonomies:
            lag: "lag"

`

var multiSiteJSONConfigTemplate = `
{
  "baseURL": "http://example.com/blog",
  "disableSitemap": false,
  "disableRSS": false,
  "rssURI": "index.xml",
  "paginate": 1,
  "disablePathToLower": true,
  "defaultContentLanguage": "{{ .DefaultContentLanguage }}",
  "defaultContentLanguageInSubdir": true,
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
      "title": "In English",
      "languageName": "English",
      "blackfriday": {
        "angledQuotes": false
      },
	  "menu": {
        "main": [
			{
			"url": "/",
			"name": "Home",
			"weight": 0
			}
		]
      }
    },
    "fr": {
      "weight": 20,
      "title": "Le Français",
      "languageName": "Français",
      "Taxonomies": {
        "plaque": "plaques"
      }
    },
    "nn": {
      "weight": 30,
      "title": "På nynorsk",
      "paginatePath": "side",
      "languageName": "Nynorsk",
      "Taxonomies": {
        "lag": "lag"
      },
	  "menu": {
        "main": [
			{
        	"url": "/",
			"name": "Heim",
			"weight": 1
			}
      	]
      }
    },
    "nb": {
      "weight": 40,
      "title": "På bokmål",
      "paginatePath": "side",
      "languageName": "Bokmål",
      "Taxonomies": {
        "lag": "lag"
      }
    }
  }
}
`

func createMultiTestSites(t *testing.T, siteConfig testSiteConfig, tomlConfigTemplate string) *HugoSites {
	return createMultiTestSitesForConfig(t, siteConfig, tomlConfigTemplate, "toml")
}

func createMultiTestSitesForConfig(t *testing.T, siteConfig testSiteConfig, configTemplate, configSuffix string) *HugoSites {

	configContent := createConfig(t, siteConfig, configTemplate)

	mf := siteConfig.Fs

	// Add some layouts
	if err := afero.WriteFile(mf,
		filepath.Join("layouts", "_default/single.html"),
		[]byte("Single: {{ .Title }}|{{ i18n \"hello\" }}|{{.Lang}}|{{ .Content }}"),
		0755); err != nil {
		t.Fatalf("Failed to write layout file: %s", err)
	}

	if err := afero.WriteFile(mf,
		filepath.Join("layouts", "_default/list.html"),
		[]byte("{{ $p := .Paginator }}List Page {{ $p.PageNumber }}: {{ .Title }}|{{ i18n \"hello\" }}|{{ .Permalink }}|Pager: {{ template \"_internal/pagination.html\" . }}"),
		0755); err != nil {
		t.Fatalf("Failed to write layout file: %s", err)
	}

	if err := afero.WriteFile(mf,
		filepath.Join("layouts", "index.html"),
		[]byte("{{ $p := .Paginator }}Default Home Page {{ $p.PageNumber }}: {{ .Title }}|{{ .IsHome }}|{{ i18n \"hello\" }}|{{ .Permalink }}|{{  .Site.Data.hugo.slogan }}"),
		0755); err != nil {
		t.Fatalf("Failed to write layout file: %s", err)
	}

	if err := afero.WriteFile(mf,
		filepath.Join("layouts", "index.fr.html"),
		[]byte("{{ $p := .Paginator }}French Home Page {{ $p.PageNumber }}: {{ .Title }}|{{ .IsHome }}|{{ i18n \"hello\" }}|{{ .Permalink }}|{{  .Site.Data.hugo.slogan }}"),
		0755); err != nil {
		t.Fatalf("Failed to write layout file: %s", err)
	}

	// Add a shortcode
	if err := afero.WriteFile(mf,
		filepath.Join("layouts", "shortcodes", "shortcode.html"),
		[]byte("Shortcode: {{ i18n \"hello\" }}"),
		0755); err != nil {
		t.Fatalf("Failed to write layout file: %s", err)
	}

	// A shortcode in multiple languages
	if err := afero.WriteFile(mf,
		filepath.Join("layouts", "shortcodes", "lingo.html"),
		[]byte("LingoDefault"),
		0755); err != nil {
		t.Fatalf("Failed to write layout file: %s", err)
	}

	if err := afero.WriteFile(mf,
		filepath.Join("layouts", "shortcodes", "lingo.fr.html"),
		[]byte("LingoFrench"),
		0755); err != nil {
		t.Fatalf("Failed to write layout file: %s", err)
	}

	// Add some language files
	if err := afero.WriteFile(mf,
		filepath.Join("i18n", "en.yaml"),
		[]byte(`
hello:
  other: "Hello"
`),
		0755); err != nil {
		t.Fatalf("Failed to write language file: %s", err)
	}
	if err := afero.WriteFile(mf,
		filepath.Join("i18n", "fr.yaml"),
		[]byte(`
hello:
  other: "Bonjour"
`),
		0755); err != nil {
		t.Fatalf("Failed to write language file: %s", err)
	}

	// Sources
	sources := []source.ByteSource{
		{Name: filepath.FromSlash("root.en.md"), Content: []byte(`---
title: root
weight: 10000
slug: root
publishdate: "2000-01-01"
---
# root
`)},
		{Name: filepath.FromSlash("sect/doc1.en.md"), Content: []byte(`---
title: doc1
weight: 1
slug: doc1-slug
tags:
 - tag1
publishdate: "2000-01-01"
---
# doc1
*some "content"*

{{< shortcode >}}

{{< lingo >}}

NOTE: slug should be used as URL
`)},
		{Name: filepath.FromSlash("sect/doc1.fr.md"), Content: []byte(`---
title: doc1
weight: 1
plaques:
 - frtag1
 - frtag2
publishdate: "2000-01-04"
---
# doc1
*quelque "contenu"*

{{< shortcode >}}

{{< lingo >}}

NOTE: should be in the 'en' Page's 'Translations' field.
NOTE: date is after "doc3"
`)},
		{Name: filepath.FromSlash("sect/doc2.en.md"), Content: []byte(`---
title: doc2
weight: 2
publishdate: "2000-01-02"
---
# doc2
*some content*
NOTE: without slug, "doc2" should be used, without ".en" as URL
`)},
		{Name: filepath.FromSlash("sect/doc3.en.md"), Content: []byte(`---
title: doc3
weight: 3
publishdate: "2000-01-03"
aliases: [/en/al/alias1,/al/alias2/]
tags:
 - tag2
 - tag1
url: /superbob
---
# doc3
*some content*
NOTE: third 'en' doc, should trigger pagination on home page.
`)},
		{Name: filepath.FromSlash("sect/doc4.md"), Content: []byte(`---
title: doc4
weight: 4
plaques:
 - frtag1
publishdate: "2000-01-05"
---
# doc4
*du contenu francophone*
NOTE: should use the defaultContentLanguage and mark this doc as 'fr'.
NOTE: doesn't have any corresponding translation in 'en'
`)},
		{Name: filepath.FromSlash("other/doc5.fr.md"), Content: []byte(`---
title: doc5
weight: 5
publishdate: "2000-01-06"
---
# doc5
*autre contenu francophone*
NOTE: should use the "permalinks" configuration with :filename
`)},
		// Add some for the stats
		{Name: filepath.FromSlash("stats/expired.fr.md"), Content: []byte(`---
title: expired
publishdate: "2000-01-06"
expiryDate: "2001-01-06"
---
# Expired
`)},
		{Name: filepath.FromSlash("stats/future.fr.md"), Content: []byte(`---
title: future
weight: 6
publishdate: "2100-01-06"
---
# Future
`)},
		{Name: filepath.FromSlash("stats/expired.en.md"), Content: []byte(`---
title: expired
weight: 7
publishdate: "2000-01-06"
expiryDate: "2001-01-06"
---
# Expired
`)},
		{Name: filepath.FromSlash("stats/future.en.md"), Content: []byte(`---
title: future
weight: 6
publishdate: "2100-01-06"
---
# Future
`)},
		{Name: filepath.FromSlash("stats/draft.en.md"), Content: []byte(`---
title: expired
publishdate: "2000-01-06"
draft: true
---
# Draft
`)},
		{Name: filepath.FromSlash("stats/tax.nn.md"), Content: []byte(`---
title: Tax NN
weight: 8
publishdate: "2000-01-06"
weight: 1001
lag:
- Sogndal
---
# Tax NN
`)},
		{Name: filepath.FromSlash("stats/tax.nb.md"), Content: []byte(`---
title: Tax NB
weight: 8
publishdate: "2000-01-06"
weight: 1002
lag:
- Sogndal
---
# Tax NB
`)},
	}

	configFile := "multilangconfig." + configSuffix
	writeToFs(t, mf, configFile, configContent)

	cfg, err := LoadConfig(mf, "", configFile)
	require.NoError(t, err)

	fs := hugofs.NewFrom(mf, cfg)

	// Hugo support using ByteSource's directly (for testing),
	// but to make it more real, we write them to the mem file system.
	for _, s := range sources {
		if err := afero.WriteFile(mf, filepath.Join("content", s.Name), s.Content, 0755); err != nil {
			t.Fatalf("Failed to write file: %s", err)
		}
	}

	// Add some data
	writeSource(t, fs, "data/hugo.toml", "slogan = \"Hugo Rocks!\"")

	sites, err := NewHugoSites(deps.DepsCfg{Fs: fs, Cfg: cfg}) //, Logger: newDebugLogger()})

	if err != nil {
		t.Fatalf("Failed to create sites: %s", err)
	}

	if len(sites.Sites) == 0 {
		t.Fatalf("Got %d sites", len(sites.Sites))
	}

	if sites.Fs.Source != mf {
		t.Fatal("FS mismatch")
	}

	return sites
}

func writeSource(t testing.TB, fs *hugofs.Fs, filename, content string) {
	writeToFs(t, fs.Source, filename, content)
}

func writeToFs(t testing.TB, fs afero.Fs, filename, content string) {
	if err := afero.WriteFile(fs, filepath.FromSlash(filename), []byte(content), 0755); err != nil {
		t.Fatalf("Failed to write file: %s", err)
	}
}

func readDestination(t testing.TB, fs *hugofs.Fs, filename string) string {
	return readFileFromFs(t, fs.Destination, filename)
}

func destinationExists(fs *hugofs.Fs, filename string) bool {
	b, err := helpers.Exists(filename, fs.Destination)
	if err != nil {
		panic(err)
	}
	return b
}

func readSource(t *testing.T, fs *hugofs.Fs, filename string) string {
	return readFileFromFs(t, fs.Source, filename)
}

func readFileFromFs(t testing.TB, fs afero.Fs, filename string) string {
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

func writeNewContentFile(t *testing.T, fs *hugofs.Fs, title, date, filename string, weight int) {
	content := newTestPage(title, date, weight)
	writeSource(t, fs, filename, content)
}

func createConfig(t *testing.T, config testSiteConfig, configTemplate string) string {
	templ, err := template.New("test").Parse(configTemplate)
	if err != nil {
		t.Fatal("Template parse failed:", err)
	}
	var b bytes.Buffer
	templ.Execute(&b, config)
	return b.String()
}
