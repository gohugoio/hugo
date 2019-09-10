// Copyright 2019 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hugolib

import (
	"os"
	"path"
	"regexp"
	"strings"
	"testing"

	"github.com/gohugoio/hugo/hugofs/files"

	"github.com/gohugoio/hugo/helpers"

	"github.com/gohugoio/hugo/hugofs"

	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/resources/page"

	"io"

	"github.com/gohugoio/hugo/htesting"

	"github.com/gohugoio/hugo/media"

	"path/filepath"

	"fmt"

	"github.com/gohugoio/hugo/deps"
	"github.com/spf13/viper"

	qt "github.com/frankban/quicktest"
)

func TestPageBundlerSiteRegular(t *testing.T) {
	c := qt.New(t)
	baseBaseURL := "https://example.com"

	for _, baseURLPath := range []string{"", "/hugo"} {
		for _, canonify := range []bool{false, true} {
			for _, ugly := range []bool{false, true} {
				baseURLPathId := baseURLPath
				if baseURLPathId == "" {
					baseURLPathId = "NONE"
				}
				ugly := ugly
				canonify := canonify
				c.Run(fmt.Sprintf("ugly=%t,canonify=%t,path=%s", ugly, canonify, baseURLPathId),
					func(c *qt.C) {
						c.Parallel()
						baseURL := baseBaseURL + baseURLPath
						relURLBase := baseURLPath
						if canonify {
							relURLBase = ""
						}
						fs, cfg := newTestBundleSources(c)
						cfg.Set("baseURL", baseURL)
						cfg.Set("canonifyURLs", canonify)

						cfg.Set("permalinks", map[string]string{
							"a": ":sections/:filename",
							"b": ":year/:slug/",
							"c": ":sections/:slug",
							"/": ":filename/",
						})

						cfg.Set("outputFormats", map[string]interface{}{
							"CUSTOMO": map[string]interface{}{
								"mediaType":     media.HTMLType,
								"baseName":      "cindex",
								"path":          "cpath",
								"permalinkable": true,
							},
						})

						cfg.Set("outputs", map[string]interface{}{
							"home":    []string{"HTML", "CUSTOMO"},
							"page":    []string{"HTML", "CUSTOMO"},
							"section": []string{"HTML", "CUSTOMO"},
						})

						cfg.Set("uglyURLs", ugly)

						b := newTestSitesBuilderFromDepsCfg(c, deps.DepsCfg{Logger: loggers.NewErrorLogger(), Fs: fs, Cfg: cfg}).WithNothingAdded()

						b.Build(BuildCfg{})

						s := b.H.Sites[0]

						c.Assert(len(s.RegularPages()), qt.Equals, 8)

						singlePage := s.getPage(page.KindPage, "a/1.md")
						c.Assert(singlePage.BundleType(), qt.Equals, files.ContentClass(""))

						c.Assert(singlePage, qt.Not(qt.IsNil))
						c.Assert(s.getPage("page", "a/1"), qt.Equals, singlePage)
						c.Assert(s.getPage("page", "1"), qt.Equals, singlePage)

						c.Assert(content(singlePage), qt.Contains, "TheContent")

						relFilename := func(basePath, outBase string) (string, string) {
							rel := basePath
							if ugly {
								rel = strings.TrimSuffix(basePath, "/") + ".html"
							}

							var filename string
							if !ugly {
								filename = path.Join(basePath, outBase)
							} else {
								filename = rel
							}

							rel = fmt.Sprintf("%s%s", relURLBase, rel)

							return rel, filename
						}

						// Check both output formats
						rel, filename := relFilename("/a/1/", "index.html")
						b.AssertFileContent(filepath.Join("/work/public", filename),
							"TheContent",
							"Single RelPermalink: "+rel,
						)

						rel, filename = relFilename("/cpath/a/1/", "cindex.html")

						b.AssertFileContent(filepath.Join("/work/public", filename),
							"TheContent",
							"Single RelPermalink: "+rel,
						)

						b.AssertFileContent(filepath.FromSlash("/work/public/images/hugo-logo.png"), "content")

						// This should be just copied to destination.
						b.AssertFileContent(filepath.FromSlash("/work/public/assets/pic1.png"), "content")

						leafBundle1 := s.getPage(page.KindPage, "b/my-bundle/index.md")
						c.Assert(leafBundle1, qt.Not(qt.IsNil))
						c.Assert(leafBundle1.BundleType(), qt.Equals, files.ContentClassLeaf)
						c.Assert(leafBundle1.Section(), qt.Equals, "b")
						sectionB := s.getPage(page.KindSection, "b")
						c.Assert(sectionB, qt.Not(qt.IsNil))
						home, _ := s.Info.Home()
						c.Assert(home.BundleType(), qt.Equals, files.ContentClassBranch)

						// This is a root bundle and should live in the "home section"
						// See https://github.com/gohugoio/hugo/issues/4332
						rootBundle := s.getPage(page.KindPage, "root")
						c.Assert(rootBundle, qt.Not(qt.IsNil))
						c.Assert(rootBundle.Parent().IsHome(), qt.Equals, true)
						if !ugly {
							b.AssertFileContent(filepath.FromSlash("/work/public/root/index.html"), "Single RelPermalink: "+relURLBase+"/root/")
							b.AssertFileContent(filepath.FromSlash("/work/public/cpath/root/cindex.html"), "Single RelPermalink: "+relURLBase+"/cpath/root/")
						}

						leafBundle2 := s.getPage(page.KindPage, "a/b/index.md")
						c.Assert(leafBundle2, qt.Not(qt.IsNil))
						unicodeBundle := s.getPage(page.KindPage, "c/bundle/index.md")
						c.Assert(unicodeBundle, qt.Not(qt.IsNil))

						pageResources := leafBundle1.Resources().ByType(pageResourceType)
						c.Assert(len(pageResources), qt.Equals, 2)
						firstPage := pageResources[0].(page.Page)
						secondPage := pageResources[1].(page.Page)

						c.Assert(firstPage.File().Filename(), qt.Equals, filepath.FromSlash("/work/base/b/my-bundle/1.md"))
						c.Assert(content(firstPage), qt.Contains, "TheContent")
						c.Assert(len(leafBundle1.Resources()), qt.Equals, 6)

						// Verify shortcode in bundled page
						c.Assert(content(secondPage), qt.Contains, filepath.FromSlash("MyShort in b/my-bundle/2.md"))

						// https://github.com/gohugoio/hugo/issues/4582
						c.Assert(firstPage.Parent(), qt.Equals, leafBundle1)
						c.Assert(secondPage.Parent(), qt.Equals, leafBundle1)

						c.Assert(pageResources.GetMatch("1*"), qt.Equals, firstPage)
						c.Assert(pageResources.GetMatch("2*"), qt.Equals, secondPage)
						c.Assert(pageResources.GetMatch("doesnotexist*"), qt.IsNil)

						imageResources := leafBundle1.Resources().ByType("image")
						c.Assert(len(imageResources), qt.Equals, 3)

						c.Assert(leafBundle1.OutputFormats().Get("CUSTOMO"), qt.Not(qt.IsNil))

						relPermalinker := func(s string) string {
							return fmt.Sprintf(s, relURLBase)
						}

						permalinker := func(s string) string {
							return fmt.Sprintf(s, baseURL)
						}

						if ugly {
							b.AssertFileContent("/work/public/2017/pageslug.html",
								relPermalinker("Single RelPermalink: %s/2017/pageslug.html"),
								permalinker("Single Permalink: %s/2017/pageslug.html"),
								relPermalinker("Sunset RelPermalink: %s/2017/pageslug/sunset1.jpg"),
								permalinker("Sunset Permalink: %s/2017/pageslug/sunset1.jpg"))
						} else {
							b.AssertFileContent("/work/public/2017/pageslug/index.html",
								relPermalinker("Sunset RelPermalink: %s/2017/pageslug/sunset1.jpg"),
								permalinker("Sunset Permalink: %s/2017/pageslug/sunset1.jpg"))

							b.AssertFileContent("/work/public/cpath/2017/pageslug/cindex.html",
								relPermalinker("Single RelPermalink: %s/cpath/2017/pageslug/"),
								relPermalinker("Short Sunset RelPermalink: %s/cpath/2017/pageslug/sunset2.jpg"),
								relPermalinker("Sunset RelPermalink: %s/cpath/2017/pageslug/sunset1.jpg"),
								permalinker("Sunset Permalink: %s/cpath/2017/pageslug/sunset1.jpg"),
							)
						}

						b.AssertFileContent(filepath.FromSlash("/work/public/2017/pageslug/c/logo.png"), "content")
						b.AssertFileContent(filepath.FromSlash("/work/public/cpath/2017/pageslug/c/logo.png"), "content")
						c.Assert(b.CheckExists("/work/public/cpath/cpath/2017/pageslug/c/logo.png"), qt.Equals, false)

						// Custom media type defined in site config.
						c.Assert(len(leafBundle1.Resources().ByType("bepsays")), qt.Equals, 1)

						if ugly {
							b.AssertFileContent(filepath.FromSlash("/work/public/2017/pageslug.html"),
								"TheContent",
								relPermalinker("Sunset RelPermalink: %s/2017/pageslug/sunset1.jpg"),
								permalinker("Sunset Permalink: %s/2017/pageslug/sunset1.jpg"),
								"Thumb Width: 123",
								"Thumb Name: my-sunset-1",
								relPermalinker("Short Sunset RelPermalink: %s/2017/pageslug/sunset2.jpg"),
								"Short Thumb Width: 56",
								"1: Image Title: Sunset Galore 1",
								"1: Image Params: map[myparam:My Sunny Param]",
								relPermalinker("1: Image RelPermalink: %s/2017/pageslug/sunset1.jpg"),
								"2: Image Title: Sunset Galore 2",
								"2: Image Params: map[myparam:My Sunny Param]",
								"1: Image myParam: Lower: My Sunny Param Caps: My Sunny Param",
								"0: Page Title: Bundle Galore",
							)

							// https://github.com/gohugoio/hugo/issues/5882
							b.AssertFileContent(
								filepath.FromSlash("/work/public/2017/pageslug.html"), "0: Page RelPermalink: |")

							b.AssertFileContent(filepath.FromSlash("/work/public/cpath/2017/pageslug.html"), "TheContent")

							// 은행
							b.AssertFileContent(filepath.FromSlash("/work/public/c/은행/logo-은행.png"), "은행 PNG")

						} else {
							b.AssertFileContent(filepath.FromSlash("/work/public/2017/pageslug/index.html"), "TheContent")
							b.AssertFileContent(filepath.FromSlash("/work/public/cpath/2017/pageslug/cindex.html"), "TheContent")
							b.AssertFileContent(filepath.FromSlash("/work/public/2017/pageslug/index.html"), "Single Title")
							b.AssertFileContent(filepath.FromSlash("/work/public/root/index.html"), "Single Title")

						}

					})
			}
		}
	}

}

func TestPageBundlerSiteMultilingual(t *testing.T) {
	t.Parallel()

	for _, ugly := range []bool{false, true} {
		ugly := ugly
		t.Run(fmt.Sprintf("ugly=%t", ugly),
			func(t *testing.T) {
				t.Parallel()
				c := qt.New(t)
				fs, cfg := newTestBundleSourcesMultilingual(t)
				cfg.Set("uglyURLs", ugly)

				b := newTestSitesBuilderFromDepsCfg(t, deps.DepsCfg{Fs: fs, Cfg: cfg}).WithNothingAdded()
				b.Build(BuildCfg{})

				sites := b.H

				c.Assert(len(sites.Sites), qt.Equals, 2)

				s := sites.Sites[0]

				c.Assert(len(s.RegularPages()), qt.Equals, 8)
				c.Assert(len(s.Pages()), qt.Equals, 16)
				//dumpPages(s.AllPages()...)
				c.Assert(len(s.AllPages()), qt.Equals, 31)

				bundleWithSubPath := s.getPage(page.KindPage, "lb/index")
				c.Assert(bundleWithSubPath, qt.Not(qt.IsNil))

				// See https://github.com/gohugoio/hugo/issues/4312
				// Before that issue:
				// A bundle in a/b/index.en.md
				// a/b/index.en.md => OK
				// a/b/index => OK
				// index.en.md => ambigous, but OK.
				// With bundles, the file name has little meaning, the folder it lives in does. So this should also work:
				// a/b
				// and probably also just b (aka "my-bundle")
				// These may also be translated, so we also need to test that.
				//  "bf", "my-bf-bundle", "index.md + nn
				bfBundle := s.getPage(page.KindPage, "bf/my-bf-bundle/index")
				c.Assert(bfBundle, qt.Not(qt.IsNil))
				c.Assert(bfBundle.Language().Lang, qt.Equals, "en")
				c.Assert(s.getPage(page.KindPage, "bf/my-bf-bundle/index.md"), qt.Equals, bfBundle)
				c.Assert(s.getPage(page.KindPage, "bf/my-bf-bundle"), qt.Equals, bfBundle)
				c.Assert(s.getPage(page.KindPage, "my-bf-bundle"), qt.Equals, bfBundle)

				nnSite := sites.Sites[1]
				c.Assert(len(nnSite.RegularPages()), qt.Equals, 7)

				bfBundleNN := nnSite.getPage(page.KindPage, "bf/my-bf-bundle/index")
				c.Assert(bfBundleNN, qt.Not(qt.IsNil))
				c.Assert(bfBundleNN.Language().Lang, qt.Equals, "nn")
				c.Assert(nnSite.getPage(page.KindPage, "bf/my-bf-bundle/index.nn.md"), qt.Equals, bfBundleNN)
				c.Assert(nnSite.getPage(page.KindPage, "bf/my-bf-bundle"), qt.Equals, bfBundleNN)
				c.Assert(nnSite.getPage(page.KindPage, "my-bf-bundle"), qt.Equals, bfBundleNN)

				// See https://github.com/gohugoio/hugo/issues/4295
				// Every resource should have its Name prefixed with its base folder.
				cBundleResources := bundleWithSubPath.Resources().Match("c/**")
				c.Assert(len(cBundleResources), qt.Equals, 4)
				bundlePage := bundleWithSubPath.Resources().GetMatch("c/page*")
				c.Assert(bundlePage, qt.Not(qt.IsNil))

				bcBundleNN, _ := nnSite.getPageNew(nil, "bc")
				c.Assert(bcBundleNN, qt.Not(qt.IsNil))
				bcBundleEN, _ := s.getPageNew(nil, "bc")
				c.Assert(bcBundleNN.Language().Lang, qt.Equals, "nn")
				c.Assert(bcBundleEN.Language().Lang, qt.Equals, "en")
				c.Assert(len(bcBundleNN.Resources()), qt.Equals, 3)
				c.Assert(len(bcBundleEN.Resources()), qt.Equals, 3)
				b.AssertFileContent("public/en/bc/data1.json", "data1")
				b.AssertFileContent("public/en/bc/data2.json", "data2")
				b.AssertFileContent("public/en/bc/logo-bc.png", "logo")
				b.AssertFileContent("public/nn/bc/data1.nn.json", "data1.nn")
				b.AssertFileContent("public/nn/bc/data2.json", "data2")
				b.AssertFileContent("public/nn/bc/logo-bc.png", "logo")

			})
	}
}

func TestMultilingualDisableDefaultLanguage(t *testing.T) {
	t.Parallel()

	c := qt.New(t)
	_, cfg := newTestBundleSourcesMultilingual(t)

	cfg.Set("disableLanguages", []string{"en"})

	err := loadDefaultSettingsFor(cfg)
	c.Assert(err, qt.IsNil)
	err = loadLanguageSettings(cfg, nil)
	c.Assert(err, qt.Not(qt.IsNil))
	c.Assert(err.Error(), qt.Contains, "cannot disable default language")
}

func TestMultilingualDisableLanguage(t *testing.T) {
	t.Parallel()

	c := qt.New(t)
	fs, cfg := newTestBundleSourcesMultilingual(t)
	cfg.Set("disableLanguages", []string{"nn"})

	b := newTestSitesBuilderFromDepsCfg(t, deps.DepsCfg{Fs: fs, Cfg: cfg}).WithNothingAdded()
	b.Build(BuildCfg{})
	sites := b.H

	c.Assert(len(sites.Sites), qt.Equals, 1)

	s := sites.Sites[0]

	c.Assert(len(s.RegularPages()), qt.Equals, 8)
	c.Assert(len(s.Pages()), qt.Equals, 16)
	// No nn pages
	c.Assert(len(s.AllPages()), qt.Equals, 16)
	s.pageMap.withEveryBundlePage(func(p *pageState) bool {
		c.Assert(p.Language().Lang != "nn", qt.Equals, true)
		return false
	})

}

func TestPageBundlerSiteWitSymbolicLinksInContent(t *testing.T) {
	skipSymlink(t)

	wd, _ := os.Getwd()
	defer func() {
		os.Chdir(wd)
	}()

	c := qt.New(t)
	// We need to use the OS fs for this.
	cfg := viper.New()
	fs := hugofs.NewFrom(hugofs.Os, cfg)

	workDir, clean, err := htesting.CreateTempDir(hugofs.Os, "hugosym")
	c.Assert(err, qt.IsNil)

	contentDirName := "content"

	contentDir := filepath.Join(workDir, contentDirName)
	c.Assert(os.MkdirAll(filepath.Join(contentDir, "a"), 0777), qt.IsNil)

	for i := 1; i <= 3; i++ {
		c.Assert(os.MkdirAll(filepath.Join(workDir, fmt.Sprintf("symcontent%d", i)), 0777), qt.IsNil)
	}

	c.Assert(os.MkdirAll(filepath.Join(workDir, "symcontent2", "a1"), 0777), qt.IsNil)

	// Symlinked sections inside content.
	os.Chdir(contentDir)
	for i := 1; i <= 3; i++ {
		c.Assert(os.Symlink(filepath.FromSlash(fmt.Sprintf(("../symcontent%d"), i)), fmt.Sprintf("symbolic%d", i)), qt.IsNil)
	}

	c.Assert(os.Chdir(filepath.Join(contentDir, "a")), qt.IsNil)

	// Create a symlink to one single content file
	c.Assert(os.Symlink(filepath.FromSlash("../../symcontent2/a1/page.md"), "page_s.md"), qt.IsNil)

	c.Assert(os.Chdir(filepath.FromSlash("../../symcontent3")), qt.IsNil)

	// Create a circular symlink. Will print some warnings.
	c.Assert(os.Symlink(filepath.Join("..", contentDirName), filepath.FromSlash("circus")), qt.IsNil)

	c.Assert(os.Chdir(workDir), qt.IsNil)

	defer clean()

	cfg.Set("workingDir", workDir)
	cfg.Set("contentDir", contentDirName)
	cfg.Set("baseURL", "https://example.com")

	layout := `{{ .Title }}|{{ .Content }}`
	pageContent := `---
slug: %s
date: 2017-10-09
---

TheContent.
`

	b := newTestSitesBuilderFromDepsCfg(t, deps.DepsCfg{
		Fs:  fs,
		Cfg: cfg,
	})

	b.WithTemplates(
		"_default/single.html", layout,
		"_default/list.html", layout,
	)

	b.WithContent(
		"a/regular.md", fmt.Sprintf(pageContent, "a1"),
	)

	b.WithSourceFile(
		"symcontent1/s1.md", fmt.Sprintf(pageContent, "s1"),
		"symcontent1/s2.md", fmt.Sprintf(pageContent, "s2"),
		// Regular files inside symlinked folder.
		"symcontent1/s1.md", fmt.Sprintf(pageContent, "s1"),
		"symcontent1/s2.md", fmt.Sprintf(pageContent, "s2"),

		// A bundle
		"symcontent2/a1/index.md", fmt.Sprintf(pageContent, ""),
		"symcontent2/a1/page.md", fmt.Sprintf(pageContent, "page"),
		"symcontent2/a1/logo.png", "image",

		// Assets
		"symcontent3/s1.png", "image",
		"symcontent3/s2.png", "image",
	)

	b.Build(BuildCfg{})
	s := b.H.Sites[0]

	c.Assert(len(s.RegularPages()), qt.Equals, 7)
	a1Bundle := s.getPage(page.KindPage, "symbolic2/a1/index.md")
	c.Assert(a1Bundle, qt.Not(qt.IsNil))
	c.Assert(len(a1Bundle.Resources()), qt.Equals, 2)
	c.Assert(len(a1Bundle.Resources().ByType(pageResourceType)), qt.Equals, 1)

	b.AssertFileContent(filepath.FromSlash(workDir+"/public/a/page/index.html"), "TheContent")
	b.AssertFileContent(filepath.FromSlash(workDir+"/public/symbolic1/s1/index.html"), "TheContent")
	b.AssertFileContent(filepath.FromSlash(workDir+"/public/symbolic2/a1/index.html"), "TheContent")

}

func TestPageBundlerHeadless(t *testing.T) {
	t.Parallel()

	cfg, fs := newTestCfg()
	c := qt.New(t)

	workDir := "/work"
	cfg.Set("workingDir", workDir)
	cfg.Set("contentDir", "base")
	cfg.Set("baseURL", "https://example.com")

	pageContent := `---
title: "Bundle Galore"
slug: s1
date: 2017-01-23
---

TheContent.

{{< myShort >}}
`

	writeSource(t, fs, filepath.Join(workDir, "layouts", "_default", "single.html"), "single {{ .Content }}")
	writeSource(t, fs, filepath.Join(workDir, "layouts", "_default", "list.html"), "list")
	writeSource(t, fs, filepath.Join(workDir, "layouts", "shortcodes", "myShort.html"), "SHORTCODE")

	writeSource(t, fs, filepath.Join(workDir, "base", "a", "index.md"), pageContent)
	writeSource(t, fs, filepath.Join(workDir, "base", "a", "l1.png"), "PNG image")
	writeSource(t, fs, filepath.Join(workDir, "base", "a", "l2.png"), "PNG image")

	writeSource(t, fs, filepath.Join(workDir, "base", "b", "index.md"), `---
title: "Headless Bundle in Topless Bar"
slug: s2
headless: true
date: 2017-01-23
---

TheContent.
HEADLESS {{< myShort >}}
`)
	writeSource(t, fs, filepath.Join(workDir, "base", "b", "l1.png"), "PNG image")
	writeSource(t, fs, filepath.Join(workDir, "base", "b", "l2.png"), "PNG image")
	writeSource(t, fs, filepath.Join(workDir, "base", "b", "p1.md"), pageContent)

	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{})

	c.Assert(len(s.RegularPages()), qt.Equals, 1)

	regular := s.getPage(page.KindPage, "a/index")
	c.Assert(regular.RelPermalink(), qt.Equals, "/s1/")

	headless := s.getPage(page.KindPage, "b/index")
	c.Assert(headless, qt.Not(qt.IsNil))
	c.Assert(headless.Title(), qt.Equals, "Headless Bundle in Topless Bar")
	c.Assert(headless.RelPermalink(), qt.Equals, "")
	c.Assert(headless.Permalink(), qt.Equals, "")
	c.Assert(content(headless), qt.Contains, "HEADLESS SHORTCODE")

	headlessResources := headless.Resources()
	c.Assert(len(headlessResources), qt.Equals, 3)
	c.Assert(len(headlessResources.Match("l*")), qt.Equals, 2)
	pageResource := headlessResources.GetMatch("p*")
	c.Assert(pageResource, qt.Not(qt.IsNil))
	p := pageResource.(page.Page)
	c.Assert(content(p), qt.Contains, "SHORTCODE")
	c.Assert(p.Name(), qt.Equals, "p1.md")

	th := newTestHelper(s.Cfg, s.Fs, t)

	th.assertFileContent(filepath.FromSlash(workDir+"/public/s1/index.html"), "TheContent")
	th.assertFileContent(filepath.FromSlash(workDir+"/public/s1/l1.png"), "PNG")

	th.assertFileNotExist(workDir + "/public/s2/index.html")
	// But the bundled resources needs to be published
	th.assertFileContent(filepath.FromSlash(workDir+"/public/s2/l1.png"), "PNG")

	// No headless bundles here, please.
	// https://github.com/gohugoio/hugo/issues/6492
	c.Assert(s.RegularPages(), qt.HasLen, 1)
	c.Assert(s.home.RegularPages(), qt.HasLen, 1)
	c.Assert(s.home.Pages(), qt.HasLen, 1)

}

func TestPageBundlerHeadlessIssue6552(t *testing.T) {
	t.Parallel()

	b := newTestSitesBuilder(t)
	b.WithContent("headless/h1/index.md", `
---
title: My Headless Bundle1
headless: true
---
`, "headless/h1/p1.md", `
---
title: P1
---
`, "headless/h2/index.md", `
---
title: My Headless Bundle2
headless: true
---
`)

	b.WithTemplatesAdded("index.html", `
{{ $headless1 := .Site.GetPage "headless/h1" }}
{{ $headless2 := .Site.GetPage "headless/h2" }}

HEADLESS1: {{ $headless1.Title }}|{{ $headless1.RelPermalink }}|{{ len $headless1.Resources }}|
HEADLESS2: {{ $headless2.Title }}{{ $headless2.RelPermalink }}|{{ len $headless2.Resources }}|

`)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/index.html", `
HEADLESS1: My Headless Bundle1||1|
HEADLESS2: My Headless Bundle2|0|
`)
}

func TestMultiSiteBundles(t *testing.T) {
	c := qt.New(t)
	b := newTestSitesBuilder(t)
	b.WithConfigFile("toml", `

baseURL = "http://example.com/"

defaultContentLanguage = "en"

[languages]
[languages.en]
weight = 10
contentDir = "content/en"
[languages.nn]
weight = 20
contentDir = "content/nn"


`)

	b.WithContent("en/mybundle/index.md", `
---
headless: true
---

`)

	b.WithContent("nn/mybundle/index.md", `
---
headless: true
---

`)

	b.WithContent("en/mybundle/data.yaml", `data en`)
	b.WithContent("en/mybundle/forms.yaml", `forms en`)
	b.WithContent("nn/mybundle/data.yaml", `data nn`)

	b.WithContent("en/_index.md", `
---
Title: Home
---

Home content.

`)

	b.WithContent("en/section-not-bundle/_index.md", `
---
Title: Section Page
---

Section content.

`)

	b.WithContent("en/section-not-bundle/single.md", `
---
Title: Section Single
Date: 2018-02-01
---

Single content.

`)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/nn/mybundle/data.yaml", "data nn")
	b.AssertFileContent("public/nn/mybundle/forms.yaml", "forms en")
	b.AssertFileContent("public/mybundle/data.yaml", "data en")
	b.AssertFileContent("public/mybundle/forms.yaml", "forms en")

	c.Assert(b.CheckExists("public/nn/nn/mybundle/data.yaml"), qt.Equals, false)
	c.Assert(b.CheckExists("public/en/mybundle/data.yaml"), qt.Equals, false)

	homeEn := b.H.Sites[0].home
	c.Assert(homeEn, qt.Not(qt.IsNil))
	c.Assert(homeEn.Date().Year(), qt.Equals, 2018)

	b.AssertFileContent("public/section-not-bundle/index.html", "Section Page", "Content: <p>Section content.</p>")
	b.AssertFileContent("public/section-not-bundle/single/index.html", "Section Single", "|<p>Single content.</p>")

}

func newTestBundleSources(t testing.TB) (*hugofs.Fs, *viper.Viper) {
	cfg, fs := newTestCfgBasic()
	c := qt.New(t)

	workDir := "/work"
	cfg.Set("workingDir", workDir)
	cfg.Set("contentDir", "base")
	cfg.Set("baseURL", "https://example.com")
	cfg.Set("mediaTypes", map[string]interface{}{
		"text/bepsays": map[string]interface{}{
			"suffixes": []string{"bep"},
		},
	})

	pageContent := `---
title: "Bundle Galore"
slug: pageslug
date: 2017-10-09
---

TheContent.
`

	pageContentShortcode := `---
title: "Bundle Galore"
slug: pageslug
date: 2017-10-09
---

TheContent.

{{< myShort >}}
`

	pageWithImageShortcodeAndResourceMetadataContent := `---
title: "Bundle Galore"
slug: pageslug
date: 2017-10-09
resources:
- src: "*.jpg"
  name: "my-sunset-:counter"
  title: "Sunset Galore :counter"
  params:
    myParam: "My Sunny Param"
---

TheContent.

{{< myShort >}}
`

	pageContentNoSlug := `---
title: "Bundle Galore #2"
date: 2017-10-09
---

TheContent.
`

	singleLayout := `
Single Title: {{ .Title }}
Single RelPermalink: {{ .RelPermalink }}
Single Permalink: {{ .Permalink }}
Content: {{ .Content }}
{{ $sunset := .Resources.GetMatch "my-sunset-1*" }}
{{ with $sunset }}
Sunset RelPermalink: {{ .RelPermalink }}
Sunset Permalink: {{ .Permalink }}
{{ $thumb := .Fill "123x123" }}
Thumb Width: {{ $thumb.Width }}
Thumb Name: {{ $thumb.Name }}
Thumb Title: {{ $thumb.Title }}
Thumb RelPermalink: {{ $thumb.RelPermalink }}
{{ end }}
{{ $types := slice "image" "page" }}
{{ range $types }}
{{ $typeTitle := . | title }}
{{ range $i, $e := $.Resources.ByType . }}
{{ $i }}: {{ $typeTitle }} Title: {{ .Title }}
{{ $i }}: {{ $typeTitle }} Name: {{ .Name }}
{{ $i }}: {{ $typeTitle }} RelPermalink: {{ .RelPermalink }}|
{{ $i }}: {{ $typeTitle }} Params: {{ printf "%v" .Params }}
{{ $i }}: {{ $typeTitle }} myParam: Lower: {{ .Params.myparam }} Caps: {{ .Params.MYPARAM }}
{{ end }}
{{ end }}
`

	myShort := `
MyShort in {{ .Page.File.Path }}:
{{ $sunset := .Page.Resources.GetMatch "my-sunset-2*" }}
{{ with $sunset }}
Short Sunset RelPermalink: {{ .RelPermalink }}
{{ $thumb := .Fill "56x56" }}
Short Thumb Width: {{ $thumb.Width }}
{{ end }}
`

	listLayout := `{{ .Title }}|{{ .Content }}`

	writeSource(t, fs, filepath.Join(workDir, "layouts", "_default", "single.html"), singleLayout)
	writeSource(t, fs, filepath.Join(workDir, "layouts", "_default", "list.html"), listLayout)
	writeSource(t, fs, filepath.Join(workDir, "layouts", "shortcodes", "myShort.html"), myShort)
	writeSource(t, fs, filepath.Join(workDir, "layouts", "shortcodes", "myShort.customo"), myShort)

	writeSource(t, fs, filepath.Join(workDir, "base", "_index.md"), pageContent)
	writeSource(t, fs, filepath.Join(workDir, "base", "_1.md"), pageContent)
	writeSource(t, fs, filepath.Join(workDir, "base", "_1.png"), pageContent)

	writeSource(t, fs, filepath.Join(workDir, "base", "images", "hugo-logo.png"), "content")
	writeSource(t, fs, filepath.Join(workDir, "base", "a", "2.md"), pageContent)
	writeSource(t, fs, filepath.Join(workDir, "base", "a", "1.md"), pageContent)

	writeSource(t, fs, filepath.Join(workDir, "base", "a", "b", "index.md"), pageContentNoSlug)
	writeSource(t, fs, filepath.Join(workDir, "base", "a", "b", "ab1.md"), pageContentNoSlug)

	// Mostly plain static assets in a folder with a page in a sub folder thrown in.
	writeSource(t, fs, filepath.Join(workDir, "base", "assets", "pic1.png"), "content")
	writeSource(t, fs, filepath.Join(workDir, "base", "assets", "pic2.png"), "content")
	writeSource(t, fs, filepath.Join(workDir, "base", "assets", "pages", "mypage.md"), pageContent)

	// Bundle
	writeSource(t, fs, filepath.Join(workDir, "base", "b", "my-bundle", "index.md"), pageWithImageShortcodeAndResourceMetadataContent)
	writeSource(t, fs, filepath.Join(workDir, "base", "b", "my-bundle", "1.md"), pageContent)
	writeSource(t, fs, filepath.Join(workDir, "base", "b", "my-bundle", "2.md"), pageContentShortcode)
	writeSource(t, fs, filepath.Join(workDir, "base", "b", "my-bundle", "custom-mime.bep"), "bepsays")
	writeSource(t, fs, filepath.Join(workDir, "base", "b", "my-bundle", "c", "logo.png"), "content")

	// Bundle with 은행 slug
	// See https://github.com/gohugoio/hugo/issues/4241
	writeSource(t, fs, filepath.Join(workDir, "base", "c", "bundle", "index.md"), `---
title: "은행 은행"
slug: 은행
date: 2017-10-09
---

Content for 은행.
`)

	// Bundle in root
	writeSource(t, fs, filepath.Join(workDir, "base", "root", "index.md"), pageWithImageShortcodeAndResourceMetadataContent)
	writeSource(t, fs, filepath.Join(workDir, "base", "root", "1.md"), pageContent)
	writeSource(t, fs, filepath.Join(workDir, "base", "root", "c", "logo.png"), "content")

	writeSource(t, fs, filepath.Join(workDir, "base", "c", "bundle", "logo-은행.png"), "은행 PNG")

	// Write a real image into one of the bundle above.
	src, err := os.Open("testdata/sunset.jpg")
	c.Assert(err, qt.IsNil)

	// We need 2 to test https://github.com/gohugoio/hugo/issues/4202
	out, err := fs.Source.Create(filepath.Join(workDir, "base", "b", "my-bundle", "sunset1.jpg"))
	c.Assert(err, qt.IsNil)
	out2, err := fs.Source.Create(filepath.Join(workDir, "base", "b", "my-bundle", "sunset2.jpg"))
	c.Assert(err, qt.IsNil)

	_, err = io.Copy(out, src)
	c.Assert(err, qt.IsNil)
	out.Close()
	src.Seek(0, 0)
	_, err = io.Copy(out2, src)
	out2.Close()
	src.Close()
	c.Assert(err, qt.IsNil)

	return fs, cfg

}

func newTestBundleSourcesMultilingual(t *testing.T) (*hugofs.Fs, *viper.Viper) {
	cfg, fs := newTestCfgBasic()

	workDir := "/work"
	cfg.Set("workingDir", workDir)
	cfg.Set("contentDir", "base")
	cfg.Set("baseURL", "https://example.com")
	cfg.Set("defaultContentLanguage", "en")

	langConfig := map[string]interface{}{
		"en": map[string]interface{}{
			"weight":       1,
			"languageName": "English",
		},
		"nn": map[string]interface{}{
			"weight":       2,
			"languageName": "Nynorsk",
		},
	}

	cfg.Set("languages", langConfig)

	pageContent := `---
slug: pageslug
date: 2017-10-09
---

TheContent.
`

	layout := `{{ .Title }}|{{ .Content }}|Lang: {{ .Site.Language.Lang }}`

	writeSource(t, fs, filepath.Join(workDir, "layouts", "_default", "single.html"), layout)
	writeSource(t, fs, filepath.Join(workDir, "layouts", "_default", "list.html"), layout)

	writeSource(t, fs, filepath.Join(workDir, "base", "1s", "mypage.md"), pageContent)
	writeSource(t, fs, filepath.Join(workDir, "base", "1s", "mypage.nn.md"), pageContent)
	writeSource(t, fs, filepath.Join(workDir, "base", "1s", "mylogo.png"), "content")

	writeSource(t, fs, filepath.Join(workDir, "base", "bb", "_index.md"), pageContent)
	writeSource(t, fs, filepath.Join(workDir, "base", "bb", "_index.nn.md"), pageContent)
	writeSource(t, fs, filepath.Join(workDir, "base", "bb", "en.md"), pageContent)
	writeSource(t, fs, filepath.Join(workDir, "base", "bb", "_1.md"), pageContent)
	writeSource(t, fs, filepath.Join(workDir, "base", "bb", "_1.nn.md"), pageContent)
	writeSource(t, fs, filepath.Join(workDir, "base", "bb", "a.png"), "content")
	writeSource(t, fs, filepath.Join(workDir, "base", "bb", "b.png"), "content")
	writeSource(t, fs, filepath.Join(workDir, "base", "bb", "b.nn.png"), "content")
	writeSource(t, fs, filepath.Join(workDir, "base", "bb", "c.nn.png"), "content")
	writeSource(t, fs, filepath.Join(workDir, "base", "bb", "b", "d.nn.png"), "content")

	writeSource(t, fs, filepath.Join(workDir, "base", "bc", "_index.md"), pageContent)
	writeSource(t, fs, filepath.Join(workDir, "base", "bc", "_index.nn.md"), pageContent)
	writeSource(t, fs, filepath.Join(workDir, "base", "bc", "page.md"), pageContent)
	writeSource(t, fs, filepath.Join(workDir, "base", "bc", "logo-bc.png"), "logo")
	writeSource(t, fs, filepath.Join(workDir, "base", "bc", "page.nn.md"), pageContent)
	writeSource(t, fs, filepath.Join(workDir, "base", "bc", "data1.json"), "data1")
	writeSource(t, fs, filepath.Join(workDir, "base", "bc", "data2.json"), "data2")
	writeSource(t, fs, filepath.Join(workDir, "base", "bc", "data1.nn.json"), "data1.nn")

	writeSource(t, fs, filepath.Join(workDir, "base", "bd", "index.md"), pageContent)
	writeSource(t, fs, filepath.Join(workDir, "base", "bd", "page.md"), pageContent)
	writeSource(t, fs, filepath.Join(workDir, "base", "bd", "page.nn.md"), pageContent)

	writeSource(t, fs, filepath.Join(workDir, "base", "be", "_index.md"), pageContent)
	writeSource(t, fs, filepath.Join(workDir, "base", "be", "page.md"), pageContent)
	writeSource(t, fs, filepath.Join(workDir, "base", "be", "page.nn.md"), pageContent)

	// Bundle leaf,  multilingual
	writeSource(t, fs, filepath.Join(workDir, "base", "lb", "index.md"), pageContent)
	writeSource(t, fs, filepath.Join(workDir, "base", "lb", "index.nn.md"), pageContent)
	writeSource(t, fs, filepath.Join(workDir, "base", "lb", "1.md"), pageContent)
	writeSource(t, fs, filepath.Join(workDir, "base", "lb", "2.md"), pageContent)
	writeSource(t, fs, filepath.Join(workDir, "base", "lb", "2.nn.md"), pageContent)
	writeSource(t, fs, filepath.Join(workDir, "base", "lb", "c", "page.md"), pageContent)
	writeSource(t, fs, filepath.Join(workDir, "base", "lb", "c", "logo.png"), "content")
	writeSource(t, fs, filepath.Join(workDir, "base", "lb", "c", "logo.nn.png"), "content")
	writeSource(t, fs, filepath.Join(workDir, "base", "lb", "c", "one.png"), "content")
	writeSource(t, fs, filepath.Join(workDir, "base", "lb", "c", "d", "deep.png"), "content")

	//Translated bundle in some sensible sub path.
	writeSource(t, fs, filepath.Join(workDir, "base", "bf", "my-bf-bundle", "index.md"), pageContent)
	writeSource(t, fs, filepath.Join(workDir, "base", "bf", "my-bf-bundle", "index.nn.md"), pageContent)
	writeSource(t, fs, filepath.Join(workDir, "base", "bf", "my-bf-bundle", "page.md"), pageContent)

	return fs, cfg
}

// https://github.com/gohugoio/hugo/issues/5858
func TestBundledResourcesWhenMultipleOutputFormats(t *testing.T) {
	t.Parallel()

	b := newTestSitesBuilder(t).Running().WithConfigFile("toml", `
baseURL = "https://example.org"
[outputs]
  # This looks odd, but it triggers the behaviour in #5858
  # The total output formats list gets sorted, so CSS before HTML.
  home = [ "CSS" ]

`)
	b.WithContent("mybundle/index.md", `
---
title: Page
date: 2017-01-15
---
`,
		"mybundle/data.json", "MyData",
	)

	b.CreateSites().Build(BuildCfg{})

	b.AssertFileContent("public/mybundle/data.json", "MyData")

	// Change the bundled JSON file and make sure it gets republished.
	b.EditFiles("content/mybundle/data.json", "My changed data")

	b.Build(BuildCfg{})

	b.AssertFileContent("public/mybundle/data.json", "My changed data")

}

// https://github.com/gohugoio/hugo/issues/4870
func TestBundleSlug(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	const pageTemplate = `---
title: Title
slug: %s
---
`

	b := newTestSitesBuilder(t)

	b.WithTemplatesAdded("index.html", `{{ range .Site.RegularPages }}|{{ .RelPermalink }}{{ end }}|`)
	b.WithSimpleConfigFile().
		WithContent("about/services1/misc.md", fmt.Sprintf(pageTemplate, "this-is-the-slug")).
		WithContent("about/services2/misc/index.md", fmt.Sprintf(pageTemplate, "this-is-another-slug"))

	b.CreateSites().Build(BuildCfg{})

	b.AssertHome(
		"|/about/services1/this-is-the-slug/|/",
		"|/about/services2/this-is-another-slug/|")

	c.Assert(b.CheckExists("public/about/services1/this-is-the-slug/index.html"), qt.Equals, true)
	c.Assert(b.CheckExists("public/about/services2/this-is-another-slug/index.html"), qt.Equals, true)

}

func TestBundleMisc(t *testing.T) {
	config := `
baseURL = "https://example.com"
defaultContentLanguage = "en"
defaultContentLanguageInSubdir = true
ignoreFiles = ["README\\.md", "content/en/ignore"]

[Languages]
[Languages.en]
weight = 99999
contentDir = "content/en"
[Languages.nn]
weight = 20
contentDir = "content/nn"
[Languages.sv]
weight = 30
contentDir = "content/sv"
[Languages.nb]
weight = 40
contentDir = "content/nb"

`

	const pageContent = `---
title: %q
---
`
	createPage := func(s string) string {
		return fmt.Sprintf(pageContent, s)
	}

	b := newTestSitesBuilder(t).WithConfigFile("toml", config)
	b.WithLogger(loggers.NewWarningLogger())

	b.WithTemplates("_default/list.html", `{{ range .Site.Pages }}
{{ .Kind }}|{{ .Path }}|{{ with .CurrentSection }}CurrentSection: {{ .Path }}{{ end }}|{{ .RelPermalink }}{{ end }}
`)

	b.WithTemplates("_default/single.html", `Single: {{ .Title }}`)

	b.WithContent("en/sect1/sect2/_index.md", createPage("en: Sect 2"))
	b.WithContent("en/sect1/sect2/page.md", createPage("en: Page"))
	b.WithContent("en/sect1/sect2/data-branch.json", "mydata")
	b.WithContent("nn/sect1/sect2/page.md", createPage("nn: Page"))
	b.WithContent("nn/sect1/sect2/data-branch.json", "my nn data")

	// En only
	b.WithContent("en/enonly/myen.md", createPage("en: Page"))
	b.WithContent("en/enonly/myendata.json", "mydata")

	// Leaf

	b.WithContent("nn/b1/index.md", createPage("nn: leaf"))
	b.WithContent("en/b1/index.md", createPage("en: leaf"))
	b.WithContent("sv/b1/index.md", createPage("sv: leaf"))
	b.WithContent("nb/b1/index.md", createPage("nb: leaf"))

	// Should be ignored
	b.WithContent("en/ignore/page.md", createPage("en: ignore"))
	b.WithContent("en/README.md", createPage("en: ignore"))

	// Both leaf and branch bundle in same dir
	b.WithContent("en/b2/index.md", `---
slug: leaf
---
`)
	b.WithContent("en/b2/_index.md", createPage("en: branch"))

	b.WithContent("en/b1/data1.json", "en: data")
	b.WithContent("sv/b1/data1.json", "sv: data")
	b.WithContent("sv/b1/data2.json", "sv: data2")
	b.WithContent("nb/b1/data2.json", "nb: data2")

	b.WithContent("en/b3/_index.md", createPage("en: branch"))
	b.WithContent("en/b3/p1.md", createPage("en: page"))
	b.WithContent("en/b3/data1.json", "en: data")

	b.Build(BuildCfg{})

	b.AssertFileContent("public/en/index.html",
		filepath.FromSlash("section|sect1/sect2/_index.md|CurrentSection: sect1/sect2/_index.md"),
		"myen.md|CurrentSection: enonly")

	b.AssertFileContentFn("public/en/index.html", func(s string) bool {
		// Check ignored files
		return !regexp.MustCompile("README|ignore").MatchString(s)

	})

	b.AssertFileContent("public/nn/index.html", filepath.FromSlash("page|sect1/sect2/page.md|CurrentSection: sect1"))
	b.AssertFileContentFn("public/nn/index.html", func(s string) bool {
		return !strings.Contains(s, "enonly")

	})

	// Check order of inherited data file
	b.AssertFileContent("public/nb/b1/data1.json", "en: data") // Default content
	b.AssertFileContent("public/nn/b1/data2.json", "sv: data") // First match

	b.AssertFileContent("public/en/enonly/myen/index.html", "Single: en: Page")
	b.AssertFileContent("public/en/enonly/myendata.json", "mydata")

	c := qt.New(t)
	c.Assert(b.CheckExists("public/sv/enonly/myen/index.html"), qt.Equals, false)

	// Both leaf and branch bundle in same dir
	// We log a warning about it, but we keep both.
	b.AssertFileContent("public/en/b2/index.html",
		"/en/b2/leaf/",
		filepath.FromSlash("section|sect1/sect2/_index.md|CurrentSection: sect1/sect2/_index.md"))

}

// Issue 6136
func TestPageBundlerPartialTranslations(t *testing.T) {
	config := `
baseURL = "https://example.org"
defaultContentLanguage = "en"
defaultContentLanguageInSubDir = true
disableKinds = ["taxonomyTerm", "taxonomy"]
[languages]
[languages.nn]
languageName = "Nynorsk"
weight = 2
title = "Tittel på Nynorsk"
[languages.en]
title = "Title in English"
languageName = "English"
weight = 1
`

	pageContent := func(id string) string {
		return fmt.Sprintf(`
---
title: %q
---
`, id)
	}

	dataContent := func(id string) string {
		return id
	}

	b := newTestSitesBuilder(t).WithConfigFile("toml", config)

	b.WithContent("blog/sect1/_index.nn.md", pageContent("s1.nn"))
	b.WithContent("blog/sect1/data.json", dataContent("s1.data"))

	b.WithContent("blog/sect1/b1/index.nn.md", pageContent("s1.b1.nn"))
	b.WithContent("blog/sect1/b1/data.json", dataContent("s1.b1.data"))

	b.WithContent("blog/sect2/_index.md", pageContent("s2"))
	b.WithContent("blog/sect2/data.json", dataContent("s2.data"))

	b.WithContent("blog/sect2/b1/index.md", pageContent("s2.b1"))
	b.WithContent("blog/sect2/b1/data.json", dataContent("s2.b1.data"))

	b.WithContent("blog/sect2/b2/index.md", pageContent("s2.b2"))
	b.WithContent("blog/sect2/b2/bp.md", pageContent("s2.b2.bundlecontent"))

	b.WithContent("blog/sect2/b3/index.md", pageContent("s2.b3"))
	b.WithContent("blog/sect2/b3/bp.nn.md", pageContent("s2.b3.bundlecontent.nn"))

	b.WithContent("blog/sect2/b4/index.nn.md", pageContent("s2.b4"))
	b.WithContent("blog/sect2/b4/bp.nn.md", pageContent("s2.b4.bundlecontent.nn"))

	b.WithTemplates("index.html", `
Num Pages: {{ len .Site.Pages }}
{{ range .Site.Pages }}
{{ .Kind }}|{{ .RelPermalink }}|Content: {{ .Title }}|Resources: {{ range .Resources }}R: {{ .Title }}|{{ .Content }}|{{ end -}}
{{ end }}
`)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/nn/index.html",
		"Num Pages: 6",
		"page|/nn/blog/sect1/b1/|Content: s1.b1.nn|Resources: R: data.json|s1.b1.data|",
		"page|/nn/blog/sect2/b3/|Content: s2.b3|Resources: R: s2.b3.bundlecontent.nn|",
		"page|/nn/blog/sect2/b4/|Content: s2.b4|Resources: R: s2.b4.bundlecontent.nn",
	)

	b.AssertFileContent("public/en/index.html",
		"Num Pages: 6",
		"section|/en/blog/sect2/|Content: s2|Resources: R: data.json|s2.data|",
		"page|/en/blog/sect2/b1/|Content: s2.b1|Resources: R: data.json|s2.b1.data|",
		"page|/en/blog/sect2/b2/|Content: s2.b2|Resources: R: s2.b2.bundlecontent|",
	)

}

// #6208
func TestBundleIndexInSubFolder(t *testing.T) {
	config := `
baseURL = "https://example.com"

`

	const pageContent = `---
title: %q
---
`
	createPage := func(s string) string {
		return fmt.Sprintf(pageContent, s)
	}

	b := newTestSitesBuilder(t).WithConfigFile("toml", config)
	b.WithLogger(loggers.NewWarningLogger())

	b.WithTemplates("_default/single.html", `{{ range .Resources }}
{{ .ResourceType }}|{{ .Title }}|
{{ end }}


`)

	b.WithContent("bundle/index.md", createPage("bundle index"))
	b.WithContent("bundle/p1.md", createPage("bundle p1"))
	b.WithContent("bundle/sub/p2.md", createPage("bundle sub p2"))
	b.WithContent("bundle/sub/index.md", createPage("bundle sub index"))
	b.WithContent("bundle/sub/data.json", "data")

	b.Build(BuildCfg{})

	b.AssertFileContent("public/bundle/index.html", `
        json|sub/data.json|
        page|bundle p1|
        page|bundle sub index|
        page|bundle sub p2|
`)

}

func TestBundleTransformMany(t *testing.T) {

	b := newTestSitesBuilder(t).WithSimpleConfigFile().Running()

	for i := 1; i <= 50; i++ {
		b.WithContent(fmt.Sprintf("bundle%d/index.md", i), fmt.Sprintf(`
---
title: "Page"
weight: %d
---
		
`, i))
		b.WithSourceFile(fmt.Sprintf("content/bundle%d/data.yaml", i), fmt.Sprintf(`data: v%d`, i))
		b.WithSourceFile(fmt.Sprintf("content/bundle%d/data.json", i), fmt.Sprintf(`{ "data": "v%d" }`, i))
		b.WithSourceFile(fmt.Sprintf("assets/data%d/data.yaml", i), fmt.Sprintf(`vdata: v%d`, i))

	}

	b.WithTemplatesAdded("_default/single.html", `
{{ $bundleYaml := .Resources.GetMatch "*.yaml" }}
{{ $bundleJSON := .Resources.GetMatch "*.json" }}
{{ $assetsYaml := resources.GetMatch (printf "data%d/*.yaml" .Weight) }}
{{ $data1 := $bundleYaml | transform.Unmarshal }}
{{ $data2 := $assetsYaml | transform.Unmarshal }}
{{ $bundleFingerprinted := $bundleYaml | fingerprint "md5" }}
{{ $assetsFingerprinted := $assetsYaml | fingerprint "md5" }}
{{ $jsonMin := $bundleJSON | minify }}
{{ $jsonMinMin := $jsonMin | minify }}
{{ $jsonMinMinMin := $jsonMinMin | minify }}

data content unmarshaled: {{ $data1.data }}
data assets content unmarshaled: {{ $data2.vdata }}
bundle fingerprinted: {{ $bundleFingerprinted.RelPermalink }}
assets fingerprinted: {{ $assetsFingerprinted.RelPermalink }}

bundle min min min: {{ $jsonMinMinMin.RelPermalink }}
bundle min min key: {{ $jsonMinMin.Key }}

`)

	for i := 0; i < 3; i++ {

		b.Build(BuildCfg{})

		for i := 1; i <= 50; i++ {
			index := fmt.Sprintf("public/bundle%d/index.html", i)
			b.AssertFileContent(fmt.Sprintf("public/bundle%d/data.yaml", i), fmt.Sprintf("data: v%d", i))
			b.AssertFileContent(index, fmt.Sprintf("data content unmarshaled: v%d", i))
			b.AssertFileContent(index, fmt.Sprintf("data assets content unmarshaled: v%d", i))

			md5Asset := helpers.MD5String(fmt.Sprintf(`vdata: v%d`, i))
			b.AssertFileContent(index, fmt.Sprintf("assets fingerprinted: /data%d/data.%s.yaml", i, md5Asset))

			// The original is not used, make sure it's not published.
			b.Assert(b.CheckExists(fmt.Sprintf("public/data%d/data.yaml", i)), qt.Equals, false)

			md5Bundle := helpers.MD5String(fmt.Sprintf(`data: v%d`, i))
			b.AssertFileContent(index, fmt.Sprintf("bundle fingerprinted: /bundle%d/data.%s.yaml", i, md5Bundle))

			b.AssertFileContent(index,
				fmt.Sprintf("bundle min min min: /bundle%d/data.min.min.min.json", i),
				fmt.Sprintf("bundle min min key: /bundle%d/data.min.min.json", i),
			)
			b.Assert(b.CheckExists(fmt.Sprintf("public/bundle%d/data.min.min.min.json", i)), qt.Equals, true)
			b.Assert(b.CheckExists(fmt.Sprintf("public/bundle%d/data.min.json", i)), qt.Equals, false)
			b.Assert(b.CheckExists(fmt.Sprintf("public/bundle%d/data.min.min.json", i)), qt.Equals, false)

		}

		b.EditFiles("assets/data/foo.yaml", "FOO")

	}
}

func TestPageBundlerHome(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	workDir, clean, err := htesting.CreateTempDir(hugofs.Os, "hugo-bundler-home")
	c.Assert(err, qt.IsNil)

	cfg := viper.New()
	cfg.Set("workingDir", workDir)
	fs := hugofs.NewFrom(hugofs.Os, cfg)

	os.MkdirAll(filepath.Join(workDir, "content"), 0777)

	defer clean()

	b := newTestSitesBuilder(t)
	b.Fs = fs

	b.WithWorkingDir(workDir).WithViper(cfg)

	b.WithContent("_index.md", "---\ntitle: Home\n---\n![Alt text](image.jpg)")
	b.WithSourceFile("content/data.json", "DATA")

	b.WithTemplates("index.html", `Title: {{ .Title }}|First Resource: {{ index .Resources 0 }}|Content: {{ .Content }}`)
	b.WithTemplates("_default/_markup/render-image.html", `Hook Len Page Resources {{ len .Page.Resources }}`)

	b.Build(BuildCfg{})
	b.AssertFileContent("public/index.html", `
Title: Home|First Resource: data.json|Content: <p>Hook Len Page Resources 1</p>
`)
}
