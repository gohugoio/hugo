package hugolib

import (
	"bytes"
	"testing"
)

const SITE_INFO_PARAM_TEMPLATE = `{{ .Site.Params.MyGlobalParam }}`

func TestSiteInfoParams(t *testing.T) {
	s := &Site{
		Config: Config{Params: map[string]interface{}{"MyGlobalParam": "FOOBAR_PARAM"}},
	}

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
