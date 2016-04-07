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
	viper.Reset()
	defer viper.Reset()

	rssURI := "customrss.xml"
	viper.Set("baseurl", "http://auth/bub/")
	viper.Set("RSSUri", rssURI)

	hugofs.InitMemFs()
	s := &Site{
		Source: &source.InMemorySource{ByteSource: weightedSources},
	}
	s.initializeSiteInfo()
	s.prepTemplates("rss.xml", rssTemplate)

	if err := s.createPages(); err != nil {
		t.Fatalf("Unable to create pages: %s", err)
	}

	if err := s.buildSiteMeta(); err != nil {
		t.Fatalf("Unable to build site metadata: %s", err)
	}

	if err := s.renderHomePage(); err != nil {
		t.Fatalf("Unable to RenderHomePage: %s", err)
	}

	file, err := hugofs.Destination().Open(rssURI)

	if err != nil {
		t.Fatalf("Unable to locate: %s", rssURI)
	}

	rss := helpers.ReaderToBytes(file)
	if !bytes.HasPrefix(rss, []byte("<?xml")) {
		t.Errorf("rss feed should start with <?xml. %s", rss)
	}
}
