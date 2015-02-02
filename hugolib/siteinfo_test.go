package hugolib

import (
	"bytes"
	"testing"

	"github.com/spf13/viper"
)

const SITE_INFO_PARAM_TEMPLATE = `{{ .Site.Params.MyGlobalParam }}`

func TestSiteInfoParams(t *testing.T) {
	viper.Set("Params", map[string]interface{}{"MyGlobalParam": "FOOBAR_PARAM"})
	s := &Site{}

	s.initialize()
	if s.Info.Params["MyGlobalParam"] != "FOOBAR_PARAM" {
		t.Errorf("Unable to set site.Info.Param")
	}
	s.prepTemplates()
	s.addTemplate("template", SITE_INFO_PARAM_TEMPLATE)
	buf := new(bytes.Buffer)

	err := s.renderThing(s.NewNode(), "template", buf)
	if err != nil {
		t.Errorf("Unable to render template: %s", err)
	}

	if buf.String() != "FOOBAR_PARAM" {
		t.Errorf("Expected FOOBAR_PARAM: got %s", buf.String())
	}
}

func TestSiteInfoPermalinks(t *testing.T) {
	viper.Set("Permalinks", map[string]interface{}{"section": "/:title"})
	s := &Site{}

	s.initialize()
	permalink := s.Info.Permalinks["section"]

	if permalink != "/:title" {
		t.Errorf("Could not set permalink (%#v)", permalink)
	}
}
