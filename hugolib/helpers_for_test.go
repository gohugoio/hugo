package hugolib

import (
	"bytes"
	"testing"

	"github.com/spf13/afero"
	"github.com/spf13/hugo/hugofs"
	"github.com/spf13/hugo/source"
	"github.com/spf13/hugo/tpl"
	"github.com/spf13/viper"
)

func siteFromByteSources(bs []source.ByteSource) *Site {

	viper.SetDefault("baseurl", "http://auth/bub")
	hugofs.DestinationFS = new(afero.MemMapFs)

	s := &Site{
		Source: &source.InMemorySource{ByteSource: bs},
	}
	s.initializeSiteInfo()
	return s
}

func buildSiteFromByteSources(bs []source.ByteSource, t *testing.T) *Site {
	s := siteFromByteSources(bs)

	if err := s.CreatePages(); err != nil {
		t.Fatalf("Unable to create pages: %s", err)
	}

	if err := s.BuildSiteMeta(); err != nil {
		t.Fatalf("Unable to build site metadata: %s", err)
	}

	return s
}

func setHugoDefaultTaxonomies() {
	taxonomies := make(map[string]string)
	taxonomies["tag"] = "tags"
	taxonomies["category"] = "categories"

	viper.Set("taxonomies", taxonomies)
}

func createAndRenderPages(t *testing.T, s *Site) {
	if err := s.CreatePages(); err != nil {
		t.Fatalf("Unable to create pages: %s", err)
	}

	if err := s.BuildSiteMeta(); err != nil {
		t.Fatalf("Unable to build site metadata: %s", err)
	}

	if err := s.RenderPages(); err != nil {
		t.Fatalf("Unable to render pages. %s", err)
	}
}

func testRenderPages(t *testing.T, s *Site) {
	if err := s.RenderPages(); err != nil {
		t.Fatalf("Unable to render pages. %s", err)
	}
}

func templatePrep(s *Site) {
	s.Tmpl = tpl.New()
	s.Tmpl.LoadTemplates(s.absLayoutDir())
	if s.hasTheme() {
		s.Tmpl.LoadTemplatesWithPrefix(s.absThemeDir()+"/layouts", "theme")
	}
}

func pageMust(p *Page, err error) *Page {
	if err != nil {
		panic(err)
	}
	return p
}

func matchRender(t *testing.T, s *Site, p *Page, tmplName string, expected string) {
	content := new(bytes.Buffer)
	err := s.renderThing(p, tmplName, NopCloser(content))
	if err != nil {
		t.Fatalf("Unable to render template.")
	}

	if string(content.Bytes()) != expected {
		t.Fatalf("Content did not match expected: %s. got: %s", expected, content)
	}
}
