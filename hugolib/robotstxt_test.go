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

const ROBOTSTXT_TEMPLATE = `User-agent: Googlebot
  {{ range .Data.Pages }}
	Disallow: {{.RelPermalink}}
	{{ end }}
`

func TestRobotsTXTOutput(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	hugofs.DestinationFS = new(afero.MemMapFs)

	viper.Set("baseurl", "http://auth/bub/")

	s := &Site{
		Source: &source.InMemorySource{ByteSource: WEIGHTED_SOURCES},
	}

	s.initializeSiteInfo()

	s.prepTemplates("robots.txt", ROBOTSTXT_TEMPLATE)

	if err := s.CreatePages(); err != nil {
		t.Fatalf("Unable to create pages: %s", err)
	}

	if err := s.BuildSiteMeta(); err != nil {
		t.Fatalf("Unable to build site metadata: %s", err)
	}

	if err := s.RenderHomePage(); err != nil {
		t.Fatalf("Unable to RenderHomePage: %s", err)
	}

	if err := s.RenderSitemap(); err != nil {
		t.Fatalf("Unable to RenderSitemap: %s", err)
	}

	if err := s.RenderRobotsTXT(); err != nil {
		t.Fatalf("Unable to RenderRobotsTXT :%s", err)
	}

	robotsFile, err := hugofs.DestinationFS.Open("robots.txt")

	if err != nil {
		t.Fatalf("Unable to locate: robots.txt")
	}

	robots := helpers.ReaderToBytes(robotsFile)
	if !bytes.HasPrefix(robots, []byte("User-agent: Googlebot")) {
		t.Errorf("Robots file should start with 'User-agentL Googlebot'. %s", robots)
	}
}
