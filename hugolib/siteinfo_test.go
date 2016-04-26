// Copyright 2015 The Hugo Authors. All rights reserved.
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

	"github.com/spf13/viper"
)

const siteInfoParamTemplate = `{{ .Site.Params.MyGlobalParam }}`

func TestSiteInfoParams(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	viper.Set("Params", map[string]interface{}{"MyGlobalParam": "FOOBAR_PARAM"})
	s := &Site{}

	s.initialize()
	if s.Info.Params["MyGlobalParam"] != "FOOBAR_PARAM" {
		t.Errorf("Unable to set site.Info.Param")
	}

	s.prepTemplates("template", siteInfoParamTemplate)

	buf := new(bytes.Buffer)

	err := s.renderThing(s.newNode(), "template", buf)
	if err != nil {
		t.Errorf("Unable to render template: %s", err)
	}

	if buf.String() != "FOOBAR_PARAM" {
		t.Errorf("Expected FOOBAR_PARAM: got %s", buf.String())
	}
}

func TestSiteInfoPermalinks(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	viper.Set("Permalinks", map[string]interface{}{"section": "/:title"})
	s := &Site{}

	s.initialize()
	permalink := s.Info.Permalinks["section"]

	if permalink != "/:title" {
		t.Errorf("Could not set permalink (%#v)", permalink)
	}
}
