// Copyright 2016 The Hugo Authors. All rights reserved.
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
	"testing"

	"reflect"

	"github.com/spf13/hugo/deps"
	"github.com/spf13/hugo/tpl"
)

const sitemapTemplate = `<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  {{ range .Data.Pages }}
  <url>
    <loc>{{ .Permalink }}</loc>{{ if not .Lastmod.IsZero }}
    <lastmod>{{ safeHTML ( .Lastmod.Format "2006-01-02T15:04:05-07:00" ) }}</lastmod>{{ end }}{{ with .Sitemap.ChangeFreq }}
    <changefreq>{{ . }}</changefreq>{{ end }}{{ if ge .Sitemap.Priority 0.0 }}
    <priority>{{ .Sitemap.Priority }}</priority>{{ end }}
  </url>
  {{ end }}
</urlset>`

func TestSitemapOutput(t *testing.T) {
	t.Parallel()
	for _, internal := range []bool{false, true} {
		doTestSitemapOutput(t, internal)
	}
}

func doTestSitemapOutput(t *testing.T, internal bool) {

	cfg, fs := newTestCfg()
	cfg.Set("baseURL", "http://auth/bub/")

	depsCfg := deps.DepsCfg{Fs: fs, Cfg: cfg}

	if !internal {
		depsCfg.WithTemplate = func(templ tpl.TemplateHandler) error {
			templ.AddTemplate("sitemap.xml", sitemapTemplate)
			return nil
		}
	}

	writeSourcesToSource(t, "content", fs, weightedSources...)
	s := buildSingleSite(t, depsCfg, BuildCfg{})
	th := testHelper{s.Cfg, s.Fs, t}

	th.assertFileContent("public/sitemap.xml",
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

}

func TestParseSitemap(t *testing.T) {
	t.Parallel()
	expected := Sitemap{Priority: 3.0, Filename: "doo.xml", ChangeFreq: "3"}
	input := map[string]interface{}{
		"changefreq": "3",
		"priority":   3.0,
		"filename":   "doo.xml",
		"unknown":    "ignore",
	}
	result := parseSitemap(input)

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Got \n%v expected \n%v", result, expected)
	}

}
