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
	"bytes"
	"testing"

	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/hugofs"
	"github.com/spf13/hugo/source"
	"github.com/spf13/viper"
	"reflect"
)

const SITEMAP_TEMPLATE = `<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  {{ range .Data.Pages }}
  <url>
    <loc>{{ .Permalink }}</loc>
    <lastmod>{{ safeHTML ( .Date.Format "2006-01-02T15:04:05-07:00" ) }}</lastmod>{{ with .Sitemap.ChangeFreq }}
    <changefreq>{{ . }}</changefreq>{{ end }}{{ if ge .Sitemap.Priority 0.0 }}
    <priority>{{ .Sitemap.Priority }}</priority>{{ end }}
  </url>
  {{ end }}
</urlset>`

func TestSitemapOutput(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	hugofs.InitMemFs()

	viper.Set("baseurl", "http://auth/bub/")

	s := &Site{
		Source: &source.InMemorySource{ByteSource: weightedSources},
	}

	s.initializeSiteInfo()

	s.prepTemplates("sitemap.xml", SITEMAP_TEMPLATE)

	if err := s.createPages(); err != nil {
		t.Fatalf("Unable to create pages: %s", err)
	}

	if err := s.buildSiteMeta(); err != nil {
		t.Fatalf("Unable to build site metadata: %s", err)
	}

	if err := s.renderHomePage(); err != nil {
		t.Fatalf("Unable to RenderHomePage: %s", err)
	}

	if err := s.renderSitemap(); err != nil {
		t.Fatalf("Unable to RenderSitemap: %s", err)
	}

	if err := s.renderRobotsTXT(); err != nil {
		t.Fatalf("Unable to RenderRobotsTXT :%s", err)
	}

	sitemapFile, err := hugofs.Destination().Open("sitemap.xml")

	if err != nil {
		t.Fatalf("Unable to locate: sitemap.xml")
	}

	sitemap := helpers.ReaderToBytes(sitemapFile)
	if !bytes.HasPrefix(sitemap, []byte("<?xml")) {
		t.Errorf("Sitemap file should start with <?xml. %s", sitemap)
	}
}

func TestParseSitemap(t *testing.T) {
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
