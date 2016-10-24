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
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

const rssTemplate = `<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
  <channel>
    <title>{{ .Title }} on {{ .Site.Title }} </title>
    <link>{{ .Permalink }}</link>
    <language>en-us</language>
    <author>Steve Francia</author>
    <rights>Francia; all rights reserved.</rights>
    <updated>{{ .Date }}</updated>
    {{ range .Data.Pages }}
    <item>
      <title>{{ .Title }}</title>
      <link>{{ .Permalink }}</link>
      <pubDate>{{ .Date.Format "Mon, 02 Jan 2006 15:04:05 MST" }}</pubDate>
      <author>Steve Francia</author>
      <guid>{{ .Permalink }}</guid>
      <description>{{ .Content | html }}</description>
    </item>
    {{ end }}
  </channel>
</rss>`

func TestRSSOutput(t *testing.T) {
	testCommonResetState()

	rssURI := "public/customrss.xml"
	viper.Set("baseURL", "http://auth/bub/")
	viper.Set("rssURI", rssURI)

	for _, s := range weightedSources {
		writeSource(t, filepath.Join("content", s.Name), string(s.Content))
	}

	writeSource(t, filepath.Join("layouts", "rss.xml"), rssTemplate)

	if err := buildAndRenderSite(newSiteDefaultLang()); err != nil {
		t.Fatalf("Failed to build site: %s", err)
	}

	assertFileContent(t, filepath.Join("public", rssURI), true, "<?xml", "rss version")

}
