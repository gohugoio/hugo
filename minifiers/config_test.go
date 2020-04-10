// Copyright 2019 The Hugo Authors. All rights reserved.
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

package minifiers

import (
	"testing"

	"github.com/spf13/viper"

	qt "github.com/frankban/quicktest"
)

func TestConfig(t *testing.T) {
	c := qt.New(t)
	v := viper.New()

	v.Set("minify", map[string]interface{}{
		"disablexml": true,
		"tdewolff": map[string]interface{}{
			"html": map[string]interface{}{
				"keepwhitespace": false,
			},
		},
	})

	conf, err := decodeConfig(v)

	c.Assert(err, qt.IsNil)

	c.Assert(conf.MinifyOutput, qt.Equals, false)

	// explicitly set value
	c.Assert(conf.Tdewolff.HTML.KeepWhitespace, qt.Equals, false)
	// default value
	c.Assert(conf.Tdewolff.HTML.KeepEndTags, qt.Equals, true)
	c.Assert(conf.Tdewolff.CSS.KeepCSS2, qt.Equals, true)

	// `enable` flags
	c.Assert(conf.DisableHTML, qt.Equals, false)
	c.Assert(conf.DisableXML, qt.Equals, true)
}

func TestConfigLegacy(t *testing.T) {
	c := qt.New(t)
	v := viper.New()

	// This was a bool < Hugo v0.58.
	v.Set("minify", true)

	conf, err := decodeConfig(v)
	c.Assert(err, qt.IsNil)
	c.Assert(conf.MinifyOutput, qt.Equals, true)

}
