package hugolib

import (
	"bytes"
	"testing"

	"github.com/spf13/afero"
	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/hugofs"
	"github.com/spf13/hugo/source"
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
	viper.Reset()
	defer viper.Reset()

	rssUri := "customrss.xml"
	viper.Set("baseurl", "http://auth/bub/")
	viper.Set("RSSUri", rssUri)

	hugofs.DestinationFS = new(afero.MemMapFs)
	s := &Site{
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

	file, err := hugofs.DestinationFS.Open(rssUri)

	if err != nil {
		t.Fatalf("Unable to locate: %s", rssUri)
	}

	rss := helpers.ReaderToBytes(file)
	if !bytes.HasPrefix(rss, []byte("<?xml")) {
		t.Errorf("rss feed should start with <?xml. %s", rss)
	}
}
