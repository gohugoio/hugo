// Copyright 2017-present The Hugo Authors. All rights reserved.
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
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"testing"

	"io"

	"github.com/spf13/afero"

	"github.com/gohugoio/hugo/media"

	"path/filepath"

	"fmt"

	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/resource"
	"github.com/spf13/viper"

	"github.com/stretchr/testify/require"
)

func TestPageBundlerSite(t *testing.T) {
	t.Parallel()

	for _, ugly := range []bool{false, true} {
		t.Run(fmt.Sprintf("ugly=%t", ugly),
			func(t *testing.T) {

				assert := require.New(t)
				cfg, fs := newTestBundleSources(t)

				cfg.Set("permalinks", map[string]string{
					"a": ":sections/:filename",
					"b": ":year/:slug/",
					"c": ":sections/:slug",
					"":  ":filename/",
				})

				cfg.Set("outputFormats", map[string]interface{}{
					"CUSTOMO": map[string]interface{}{
						"mediaType": media.HTMLType,
						"baseName":  "cindex",
						"path":      "cpath",
					},
				})

				cfg.Set("outputs", map[string]interface{}{
					"home":    []string{"HTML", "CUSTOMO"},
					"page":    []string{"HTML", "CUSTOMO"},
					"section": []string{"HTML", "CUSTOMO"},
				})

				cfg.Set("uglyURLs", ugly)

				s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{})

				th := testHelper{s.Cfg, s.Fs, t}

				assert.Len(s.RegularPages, 8)

				singlePage := s.getPage(KindPage, "a/1.md")

				assert.NotNil(singlePage)
				assert.Equal(singlePage, s.getPage("page", "a/1"))
				assert.Equal(singlePage, s.getPage("page", "1"))

				assert.Contains(singlePage.Content, "TheContent")

				if ugly {
					assert.Equal("/a/1.html", singlePage.RelPermalink())
					th.assertFileContent(filepath.FromSlash("/work/public/a/1.html"), "TheContent")

				} else {
					assert.Equal("/a/1/", singlePage.RelPermalink())
					th.assertFileContent(filepath.FromSlash("/work/public/a/1/index.html"), "TheContent")
				}

				th.assertFileContent(filepath.FromSlash("/work/public/images/hugo-logo.png"), "content")

				// This should be just copied to destination.
				th.assertFileContent(filepath.FromSlash("/work/public/assets/pic1.png"), "content")

				leafBundle1 := s.getPage(KindPage, "b/my-bundle/index.md")
				assert.NotNil(leafBundle1)
				assert.Equal("b", leafBundle1.Section())
				assert.NotNil(s.getPage(KindSection, "b"))

				// This is a root bundle and should live in the "home section"
				// See https://github.com/gohugoio/hugo/issues/4332
				rootBundle := s.getPage(KindPage, "root")
				assert.NotNil(rootBundle)
				assert.True(rootBundle.Parent().IsHome())
				if ugly {
					assert.Equal("/root.html", rootBundle.RelPermalink())
				} else {
					assert.Equal("/root/", rootBundle.RelPermalink())
				}

				leafBundle2 := s.getPage(KindPage, "a/b/index.md")
				assert.NotNil(leafBundle2)
				unicodeBundle := s.getPage(KindPage, "c/bundle/index.md")
				assert.NotNil(unicodeBundle)

				pageResources := leafBundle1.Resources.ByType(pageResourceType)
				assert.Len(pageResources, 2)
				firstPage := pageResources[0].(*Page)
				secondPage := pageResources[1].(*Page)
				assert.Equal(filepath.FromSlash("b/my-bundle/1.md"), firstPage.pathOrTitle(), secondPage.pathOrTitle())
				assert.Contains(firstPage.Content, "TheContent")
				assert.Equal(6, len(leafBundle1.Resources))

				assert.Equal(firstPage, pageResources.GetByPrefix("1"))
				assert.Equal(secondPage, pageResources.GetByPrefix("2"))
				assert.Nil(pageResources.GetByPrefix("doesnotexist"))

				imageResources := leafBundle1.Resources.ByType("image")
				assert.Equal(3, len(imageResources))
				image := imageResources[0]

				altFormat := leafBundle1.OutputFormats().Get("CUSTOMO")
				assert.NotNil(altFormat)

				assert.Equal(filepath.FromSlash("/work/base/b/my-bundle/c/logo.png"), image.(resource.Source).AbsSourceFilename())
				assert.Equal("https://example.com/2017/pageslug/c/logo.png", image.Permalink())
				th.assertFileContent(filepath.FromSlash("/work/public/2017/pageslug/c/logo.png"), "content")
				th.assertFileContent(filepath.FromSlash("/work/public/cpath/2017/pageslug/c/logo.png"), "content")

				// Custom media type defined in site config.
				assert.Len(leafBundle1.Resources.ByType("bepsays"), 1)

				if ugly {
					assert.Equal("/2017/pageslug.html", leafBundle1.RelPermalink())
					th.assertFileContent(filepath.FromSlash("/work/public/2017/pageslug.html"),
						"TheContent",
						"Sunset RelPermalink: /2017/pageslug/sunset1.jpg",
						"Thumb Width: 123",
						"Thumb Name: my-sunset-1",
						"Short Sunset RelPermalink: /2017/pageslug/sunset2.jpg",
						"Short Thumb Width: 56",
						"1: Image Title: Sunset Galore 1",
						"1: Image Params: map[myparam:My Sunny Param]",
						"2: Image Title: Sunset Galore 2",
						"2: Image Params: map[myparam:My Sunny Param]",
						"1: Image myParam: Lower: My Sunny Param Caps: My Sunny Param",
					)
					th.assertFileContent(filepath.FromSlash("/work/public/cpath/2017/pageslug.html"), "TheContent")

					assert.Equal("/a/b.html", leafBundle2.RelPermalink())

					// 은행
					assert.Equal("/c/%EC%9D%80%ED%96%89.html", unicodeBundle.RelPermalink())
					th.assertFileContent(filepath.FromSlash("/work/public/c/은행.html"), "Content for 은행")
					th.assertFileContent(filepath.FromSlash("/work/public/c/은행/logo-은행.png"), "은행 PNG")

				} else {
					assert.Equal("/2017/pageslug/", leafBundle1.RelPermalink())
					th.assertFileContent(filepath.FromSlash("/work/public/2017/pageslug/index.html"), "TheContent")
					th.assertFileContent(filepath.FromSlash("/work/public/cpath/2017/pageslug/cindex.html"), "TheContent")
					th.assertFileContent(filepath.FromSlash("/work/public/2017/pageslug/index.html"), "Single Title")
					th.assertFileContent(filepath.FromSlash("/work/public/root/index.html"), "Single Title")

					assert.Equal("/a/b/", leafBundle2.RelPermalink())

				}

			})
	}

}

func TestPageBundlerSiteMultilingual(t *testing.T) {
	t.Parallel()

	for _, ugly := range []bool{false, true} {
		t.Run(fmt.Sprintf("ugly=%t", ugly),
			func(t *testing.T) {

				assert := require.New(t)
				cfg, fs := newTestBundleSourcesMultilingual(t)

				cfg.Set("uglyURLs", ugly)

				assert.NoError(loadDefaultSettingsFor(cfg))
				sites, err := NewHugoSites(deps.DepsCfg{Fs: fs, Cfg: cfg})
				assert.NoError(err)
				assert.Equal(2, len(sites.Sites))

				assert.NoError(sites.Build(BuildCfg{}))

				s := sites.Sites[0]

				assert.Equal(8, len(s.RegularPages))
				assert.Equal(16, len(s.Pages))
				assert.Equal(31, len(s.AllPages))

				bundleWithSubPath := s.getPage(KindPage, "lb/index")
				assert.NotNil(bundleWithSubPath)

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
				bfBundle := s.getPage(KindPage, "bf/my-bf-bundle/index")
				assert.NotNil(bfBundle)
				assert.Equal("en", bfBundle.Lang())
				assert.Equal(bfBundle, s.getPage(KindPage, "bf/my-bf-bundle/index.md"))
				assert.Equal(bfBundle, s.getPage(KindPage, "bf/my-bf-bundle"))
				assert.Equal(bfBundle, s.getPage(KindPage, "my-bf-bundle"))

				nnSite := sites.Sites[1]
				assert.Equal(7, len(nnSite.RegularPages))

				bfBundleNN := nnSite.getPage(KindPage, "bf/my-bf-bundle/index")
				assert.NotNil(bfBundleNN)
				assert.Equal("nn", bfBundleNN.Lang())
				assert.Equal(bfBundleNN, nnSite.getPage(KindPage, "bf/my-bf-bundle/index.nn.md"))
				assert.Equal(bfBundleNN, nnSite.getPage(KindPage, "bf/my-bf-bundle"))
				assert.Equal(bfBundleNN, nnSite.getPage(KindPage, "my-bf-bundle"))

				// See https://github.com/gohugoio/hugo/issues/4295
				// Every resource should have its Name prefixed with its base folder.
				cBundleResources := bundleWithSubPath.Resources.ByPrefix("c/")
				assert.Equal(4, len(cBundleResources))
				bundlePage := bundleWithSubPath.Resources.GetByPrefix("c/page")
				assert.NotNil(bundlePage)
				assert.IsType(&Page{}, bundlePage)

			})
	}
}

func TestMultilingualDisableDefaultLanguage(t *testing.T) {
	t.Parallel()

	assert := require.New(t)
	cfg, _ := newTestBundleSourcesMultilingual(t)

	cfg.Set("disableLanguages", []string{"en"})

	err := loadDefaultSettingsFor(cfg)
	assert.Error(err)
	assert.Contains(err.Error(), "cannot disable default language")
}

func TestMultilingualDisableLanguage(t *testing.T) {
	t.Parallel()

	assert := require.New(t)
	cfg, fs := newTestBundleSourcesMultilingual(t)
	cfg.Set("disableLanguages", []string{"nn"})

	assert.NoError(loadDefaultSettingsFor(cfg))
	sites, err := NewHugoSites(deps.DepsCfg{Fs: fs, Cfg: cfg})
	assert.NoError(err)
	assert.Equal(1, len(sites.Sites))

	assert.NoError(sites.Build(BuildCfg{}))

	s := sites.Sites[0]

	assert.Equal(8, len(s.RegularPages))
	assert.Equal(16, len(s.Pages))
	// No nn pages
	assert.Equal(16, len(s.AllPages))
	for _, p := range s.rawAllPages {
		assert.True(p.Lang() != "nn")
	}
	for _, p := range s.AllPages {
		assert.True(p.Lang() != "nn")
	}

}

func TestPageBundlerSiteWitSymbolicLinksInContent(t *testing.T) {
	assert := require.New(t)
	cfg, fs, workDir := newTestBundleSymbolicSources(t)

	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg, Logger: newWarningLogger()}, BuildCfg{})

	th := testHelper{s.Cfg, s.Fs, t}

	assert.Equal(7, len(s.RegularPages))
	a1Bundle := s.getPage(KindPage, "symbolic2/a1/index.md")
	assert.NotNil(a1Bundle)
	assert.Equal(2, len(a1Bundle.Resources))
	assert.Equal(1, len(a1Bundle.Resources.ByType(pageResourceType)))

	th.assertFileContent(filepath.FromSlash(workDir+"/public/a/page/index.html"), "TheContent")
	th.assertFileContent(filepath.FromSlash(workDir+"/public/symbolic1/s1/index.html"), "TheContent")
	th.assertFileContent(filepath.FromSlash(workDir+"/public/symbolic2/a1/index.html"), "TheContent")

}

func TestPageBundlerHeadless(t *testing.T) {
	t.Parallel()

	cfg, fs := newTestCfg()
	assert := require.New(t)

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

	assert.Equal(1, len(s.RegularPages))
	assert.Equal(1, len(s.headlessPages))

	regular := s.getPage(KindPage, "a/index")
	assert.Equal("/a/s1/", regular.RelPermalink())

	headless := s.getPage(KindPage, "b/index")
	assert.NotNil(headless)
	assert.True(headless.headless)
	assert.Equal("Headless Bundle in Topless Bar", headless.Title())
	assert.Equal("", headless.RelPermalink())
	assert.Equal("", headless.Permalink())
	assert.Contains(headless.Content, "HEADLESS SHORTCODE")

	headlessResources := headless.Resources
	assert.Equal(3, len(headlessResources))
	assert.Equal(2, len(headlessResources.Match("l*")))
	pageResource := headlessResources.GetMatch("p*")
	assert.NotNil(pageResource)
	assert.IsType(&Page{}, pageResource)
	p := pageResource.(*Page)
	assert.Contains(p.Content, "SHORTCODE")
	assert.Equal("p1.md", p.Name())

	th := testHelper{s.Cfg, s.Fs, t}

	th.assertFileContent(filepath.FromSlash(workDir+"/public/a/s1/index.html"), "TheContent")
	th.assertFileContent(filepath.FromSlash(workDir+"/public/a/s1/l1.png"), "PNG")

	th.assertFileNotExist(workDir + "/public/b/s2/index.html")
	// But the bundled resources needs to be published
	th.assertFileContent(filepath.FromSlash(workDir+"/public/b/s2/l1.png"), "PNG")

}

func newTestBundleSources(t *testing.T) (*viper.Viper, *hugofs.Fs) {
	cfg, fs := newTestCfg()
	assert := require.New(t)

	workDir := "/work"
	cfg.Set("workingDir", workDir)
	cfg.Set("contentDir", "base")
	cfg.Set("baseURL", "https://example.com")
	cfg.Set("mediaTypes", map[string]interface{}{
		"text/bepsays": map[string]interface{}{
			"suffix": "bep",
		},
	})

	pageContent := `---
title: "Bundle Galore"
slug: pageslug
date: 2017-10-09
---

TheContent.
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
Content: {{ .Content }}
{{ $sunset := .Resources.GetByPrefix "my-sunset-1" }}
{{ with $sunset }}
Sunset RelPermalink: {{ .RelPermalink }}
{{ $thumb := .Fill "123x123" }}
Thumb Width: {{ $thumb.Width }}
Thumb Name: {{ $thumb.Name }}
Thumb Title: {{ $thumb.Title }}
Thumb RelPermalink: {{ $thumb.RelPermalink }}
{{ end }}
{{ range $i, $e := .Resources.ByType "image" }}
{{ $i }}: Image Title: {{ .Title }}
{{ $i }}: Image Name: {{ .Name }}
{{ $i }}: Image Params: {{ printf "%v" .Params }}
{{ $i }}: Image myParam: Lower: {{ .Params.myparam }} Caps: {{ .Params.MYPARAM }}
{{ end }}
`

	myShort := `
{{ $sunset := .Page.Resources.GetByPrefix "my-sunset-2" }}
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
	writeSource(t, fs, filepath.Join(workDir, "base", "b", "my-bundle", "2.md"), pageContent)
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
	assert.NoError(err)

	// We need 2 to test https://github.com/gohugoio/hugo/issues/4202
	out, err := fs.Source.Create(filepath.Join(workDir, "base", "b", "my-bundle", "sunset1.jpg"))
	assert.NoError(err)
	out2, err := fs.Source.Create(filepath.Join(workDir, "base", "b", "my-bundle", "sunset2.jpg"))
	assert.NoError(err)

	_, err = io.Copy(out, src)
	out.Close()
	src.Seek(0, 0)
	_, err = io.Copy(out2, src)
	out2.Close()
	src.Close()
	assert.NoError(err)

	return cfg, fs
}

func newTestBundleSourcesMultilingual(t *testing.T) (*viper.Viper, *hugofs.Fs) {
	cfg, fs := newTestCfg()

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
	writeSource(t, fs, filepath.Join(workDir, "base", "bc", "page.md"), pageContent)
	writeSource(t, fs, filepath.Join(workDir, "base", "bc", "logo-bc.png"), pageContent)
	writeSource(t, fs, filepath.Join(workDir, "base", "bc", "page.nn.md"), pageContent)

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

	return cfg, fs
}

func newTestBundleSymbolicSources(t *testing.T) (*viper.Viper, *hugofs.Fs, string) {
	assert := require.New(t)
	// We need to use the OS fs for this.
	cfg := viper.New()
	fs := hugofs.NewFrom(hugofs.Os, cfg)
	fs.Destination = &afero.MemMapFs{}
	loadDefaultSettingsFor(cfg)

	workDir, err := ioutil.TempDir("", "hugosym")

	if runtime.GOOS == "darwin" && !strings.HasPrefix(workDir, "/private") {
		// To get the entry folder in line with the rest. This its a little bit
		// mysterious, but so be it.
		workDir = "/private" + workDir
	}

	contentDir := "base"
	cfg.Set("workingDir", workDir)
	cfg.Set("contentDir", contentDir)
	cfg.Set("baseURL", "https://example.com")

	layout := `{{ .Title }}|{{ .Content }}`
	pageContent := `---
slug: %s
date: 2017-10-09
---

TheContent.
`

	fs.Source.MkdirAll(filepath.Join(workDir, "layouts", "_default"), 0777)
	fs.Source.MkdirAll(filepath.Join(workDir, contentDir), 0777)
	fs.Source.MkdirAll(filepath.Join(workDir, contentDir, "a"), 0777)
	for i := 1; i <= 3; i++ {
		fs.Source.MkdirAll(filepath.Join(workDir, fmt.Sprintf("symcontent%d", i)), 0777)

	}
	fs.Source.MkdirAll(filepath.Join(workDir, "symcontent2", "a1"), 0777)

	writeSource(t, fs, filepath.Join(workDir, "layouts", "_default", "single.html"), layout)
	writeSource(t, fs, filepath.Join(workDir, "layouts", "_default", "list.html"), layout)

	writeSource(t, fs, filepath.Join(workDir, contentDir, "a", "regular.md"), fmt.Sprintf(pageContent, "a1"))

	// Regular files inside symlinked folder.
	writeSource(t, fs, filepath.Join(workDir, "symcontent1", "s1.md"), fmt.Sprintf(pageContent, "s1"))
	writeSource(t, fs, filepath.Join(workDir, "symcontent1", "s2.md"), fmt.Sprintf(pageContent, "s2"))

	// A bundle
	writeSource(t, fs, filepath.Join(workDir, "symcontent2", "a1", "index.md"), fmt.Sprintf(pageContent, ""))
	writeSource(t, fs, filepath.Join(workDir, "symcontent2", "a1", "page.md"), fmt.Sprintf(pageContent, "page"))
	writeSource(t, fs, filepath.Join(workDir, "symcontent2", "a1", "logo.png"), "image")

	// Assets
	writeSource(t, fs, filepath.Join(workDir, "symcontent3", "s1.png"), "image")
	writeSource(t, fs, filepath.Join(workDir, "symcontent3", "s2.png"), "image")

	wd, _ := os.Getwd()
	defer func() {
		os.Chdir(wd)
	}()
	// Symlinked sections inside content.
	os.Chdir(filepath.Join(workDir, contentDir))
	for i := 1; i <= 3; i++ {
		assert.NoError(os.Symlink(filepath.FromSlash(fmt.Sprintf(("../symcontent%d"), i)), fmt.Sprintf("symbolic%d", i)))
	}

	os.Chdir(filepath.Join(workDir, contentDir, "a"))

	// Create a symlink to one single content file
	assert.NoError(os.Symlink(filepath.FromSlash("../../symcontent2/a1/page.md"), "page_s.md"))

	os.Chdir(filepath.FromSlash("../../symcontent3"))

	// Create a circular symlink. Will print some warnings.
	assert.NoError(os.Symlink(filepath.Join("..", contentDir), filepath.FromSlash("circus")))

	os.Chdir(workDir)
	assert.NoError(err)

	return cfg, fs, workDir
}
