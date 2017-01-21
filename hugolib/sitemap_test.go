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

	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/source"
	"github.com/spf13/viper"
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
	for _, internal := range []bool{false, true} {
		doTestSitemapOutput(t, internal)
	}
}

func doTestSitemapOutput(t *testing.T, internal bool) {
	testCommonResetState()

	viper.Set("baseURL", "http://auth/bub/")

	s := &Site{
		deps:     newDeps(DepsCfg{}),
		Source:   &source.InMemorySource{ByteSource: weightedSources},
		Language: helpers.NewDefaultLanguage(),
	}

	if internal {
		if err := buildAndRenderSite(s); err != nil {
			t.Fatalf("Failed to build site: %s", err)
		}

	} else {
		if err := buildAndRenderSite(s, "sitemap.xml", sitemapTemplate); err != nil {
			t.Fatalf("Failed to build site: %s", err)
		}
	}

	assertFileContent(t, "public/sitemap.xml", true,
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

func TestParseSitemapUnknownAndFilename(t *testing.T) {
	expected := Sitemap{Filename: "doo.xml", Priority: 0.5}
	input := map[string]interface{}{
		"filename": "doo.xml",
		"unknown":  "ignore",
	}
	result := parseSitemap(input)

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Got \n%v expected \n%v", result, expected)
	}
}

func TestParseSitemapChangefreq(t *testing.T) {
	var sitemaps = []struct {
		expected Sitemap
		result   Sitemap
	}{
		{
			Sitemap{Priority: 0.5, Filename: "sitemap.xml", ChangeFreq: "always"},
			parseSitemap(map[string]interface{}{"changefreq": "always"}),
		},
		{
			Sitemap{Priority: 0.5, Filename: "sitemap.xml", ChangeFreq: "hourly"},
			parseSitemap(map[string]interface{}{"changefreq": "hourly"}),
		},
		{
			Sitemap{Priority: 0.5, Filename: "sitemap.xml", ChangeFreq: "daily"},
			parseSitemap(map[string]interface{}{"changefreq": "daily"}),
		},
		{
			Sitemap{Priority: 0.5, Filename: "sitemap.xml", ChangeFreq: "weekly"},
			parseSitemap(map[string]interface{}{"changefreq": "weekly"}),
		},
		{
			Sitemap{Priority: 0.5, Filename: "sitemap.xml", ChangeFreq: "monthly"},
			parseSitemap(map[string]interface{}{"changefreq": "monthly"}),
		},
		{
			Sitemap{Priority: 0.5, Filename: "sitemap.xml", ChangeFreq: "yearly"},
			parseSitemap(map[string]interface{}{"changefreq": "yearly"}),
		},
		{
			Sitemap{Priority: 0.5, Filename: "sitemap.xml", ChangeFreq: "never"},
			parseSitemap(map[string]interface{}{"changefreq": "never"}),
		},
		{
			Sitemap{Priority: 0.5, Filename: "sitemap.xml", ChangeFreq: "never"},
			parseSitemap(map[string]interface{}{"changefreq": "NEVER"}),
		},
		{
			Sitemap{Priority: 0.5, Filename: "sitemap.xml"},
			parseSitemap(map[string]interface{}{"changefreq": "invalid"}),
		},
	}

	for _, s := range sitemaps {
		if !reflect.DeepEqual(s.expected, s.result) {
			t.Errorf("Got \n%v expected \n%v", s.result, s.expected)
		}
	}
}

func TestParseSitemapPriority(t *testing.T) {
	var sitemaps = []struct {
		expected Sitemap
		result   Sitemap
	}{
		{
			Sitemap{Priority: 0.5, Filename: "sitemap.xml"},
			parseSitemap(map[string]interface{}{"priority": -1.0}),
		},
		{
			Sitemap{Priority: 0, Filename: "sitemap.xml"},
			parseSitemap(map[string]interface{}{"priority": 0}),
		},
		{
			Sitemap{Priority: 0.1, Filename: "sitemap.xml"},
			parseSitemap(map[string]interface{}{"priority": 0.1}),
		},
		{
			Sitemap{Priority: 0.5, Filename: "sitemap.xml"},
			parseSitemap(map[string]interface{}{"priority": 0.5}),
		},
		{
			Sitemap{Priority: 1.0, Filename: "sitemap.xml"},
			parseSitemap(map[string]interface{}{"priority": 1.0}),
		},
		{
			Sitemap{Priority: 0.5, Filename: "sitemap.xml"},
			parseSitemap(map[string]interface{}{"priority": 1.1}),
		},
		{
			Sitemap{Priority: 0.5, Filename: "sitemap.xml"},
			parseSitemap(map[string]interface{}{"priority": 0.01}),
		},
		{
			Sitemap{Priority: 0.5, Filename: "sitemap.xml"},
			parseSitemap(map[string]interface{}{"priority": 1.01}),
		},
		{
			Sitemap{Priority: 0.5, Filename: "sitemap.xml"},
			parseSitemap(map[string]interface{}{"priority": "text"}),
		},
	}

	for _, s := range sitemaps {
		if !reflect.DeepEqual(s.expected, s.result) {
			t.Errorf("Got \n%v expected \n%v", s.result, s.expected)
		}
	}
}
