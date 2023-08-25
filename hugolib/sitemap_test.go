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
	"reflect"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/deps"
)

const sitemapTemplate = `<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  {{ range .Data.Pages }}
    {{- if .Permalink -}}
  <url>
    <loc>{{ .Permalink }}</loc>{{ if not .Lastmod.IsZero }}
    <lastmod>{{ safeHTML ( .Lastmod.Format "2006-01-02T15:04:05-07:00" ) }}</lastmod>{{ end }}{{ with .Sitemap.ChangeFreq }}
    <changefreq>{{ . }}</changefreq>{{ end }}{{ if ge .Sitemap.Priority 0.0 }}
    <priority>{{ .Sitemap.Priority }}</priority>{{ end }}
  </url>
    {{- end -}}
  {{ end }}
</urlset>`

func TestSitemapOutput(t *testing.T) {
	t.Parallel()
	for _, internal := range []bool{false, true} {
		doTestSitemapOutput(t, internal)
	}
}

func doTestSitemapOutput(t *testing.T, internal bool) {
	c := qt.New(t)
	cfg, fs := newTestCfg()
	cfg.Set("baseURL", "http://auth/bub/")
	cfg.Set("defaultContentLanguageInSubdir", false)
	configs, err := loadTestConfigFromProvider(cfg)
	c.Assert(err, qt.IsNil)
	writeSource(t, fs, "layouts/sitemap.xml", sitemapTemplate)
	// We want to check that the 404 page is not included in the sitemap
	// output. This template should have no effect either way, but include
	// it for the clarity.
	writeSource(t, fs, "layouts/404.html", "Not found")

	depsCfg := deps.DepsCfg{Fs: fs, Configs: configs}

	writeSourcesToSource(t, "content", fs, weightedSources...)
	s := buildSingleSite(t, depsCfg, BuildCfg{})
	th := newTestHelper(s.conf, s.Fs, t)
	outputSitemap := "public/sitemap.xml"

	th.assertFileContent(outputSitemap,
		// Regular page
		" <loc>http://auth/bub/sect/doc1/</loc>",
		// Home page
		"<loc>http://auth/bub/</loc>",
		// Section
		"<loc>http://auth/bub/sect/</loc>",
		// Tax terms
		"<loc>http://auth/bub/categories/</loc>",
		// Tax list
		"<loc>http://auth/bub/categories/hugo/</loc>",
	)

	content := readWorkingDir(th, th.Fs, outputSitemap)
	c.Assert(content, qt.Not(qt.Contains), "404")
	c.Assert(content, qt.Not(qt.Contains), "<loc></loc>")
}

func TestParseSitemap(t *testing.T) {
	t.Parallel()
	expected := config.SitemapConfig{Priority: 3.0, Filename: "doo.xml", ChangeFreq: "3"}
	input := map[string]any{
		"changefreq": "3",
		"priority":   3.0,
		"filename":   "doo.xml",
		"unknown":    "ignore",
	}
	result, err := config.DecodeSitemap(config.SitemapConfig{}, input)
	if err != nil {
		t.Fatalf("Failed to parse sitemap: %s", err)
	}

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Got \n%v expected \n%v", result, expected)
	}
}

// https://github.com/gohugoio/hugo/issues/5910
func TestSitemapOutputFormats(t *testing.T) {
	b := newTestSitesBuilder(t).WithSimpleConfigFile()

	b.WithContent("blog/html-amp.md", `
---
Title: AMP and HTML
outputs: [ "html", "amp" ]
---

`)

	b.Build(BuildCfg{})

	// Should link to the HTML version.
	b.AssertFileContent("public/sitemap.xml", " <loc>http://example.com/blog/html-amp/</loc>")
}
