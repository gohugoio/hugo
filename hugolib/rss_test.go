package hugolib

import (
	"bytes"
	"testing"

	"github.com/spf13/hugo/source"
	"github.com/spf13/hugo/target"
	"github.com/spf13/viper"
)

const RSS_TEMPLATE = `<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
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
	files := make(map[string][]byte)
	target := &target.InMemoryTarget{Files: files}
	viper.Set("baseurl", "http://auth/bub/")
	s := &Site{
		Target: target,
		Source: &source.InMemorySource{ByteSource: WEIGHTED_SOURCES},
	}
	s.initializeSiteInfo()
	s.prepTemplates()
	//  Add an rss.xml template to invoke the rss build.
	s.addTemplate("rss.xml", RSS_TEMPLATE)

	if err := s.CreatePages(); err != nil {
		t.Fatalf("Unable to create pages: %s", err)
	}

	if err := s.BuildSiteMeta(); err != nil {
		t.Fatalf("Unable to build site metadata: %s", err)
	}

	if err := s.RenderHomePage(); err != nil {
		t.Fatalf("Unable to RenderHomePage: %s", err)
	}

	if _, ok := files[".xml"]; !ok {
		t.Errorf("Unable to locate: %s", ".xml")
		t.Logf("%q", files)
	}

	rss, _ := files[".xml"]
	if !bytes.HasPrefix(rss, []byte("<?xml")) {
		t.Errorf("rss feed should start with <?xml. %s", rss)
	}
}
