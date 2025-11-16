// Copyright 2024 The Hugo Authors. All rights reserved.
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

package minifiers_test

import (
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/config/testconfig"
)

func TestConfig(t *testing.T) {
	c := qt.New(t)
	v := config.New()

	v.Set("minify", map[string]any{
		"disablexml": true,
		"tdewolff": map[string]any{
			"html": map[string]any{
				"keepwhitespace": false,
			},
		},
	})

	conf := testconfig.GetTestConfigs(nil, v).Base.Minify

	c.Assert(conf.MinifyOutput, qt.Equals, false)

	// explicitly set value
	c.Assert(conf.Tdewolff.HTML.KeepWhitespace, qt.Equals, false)
	// default value
	c.Assert(conf.Tdewolff.HTML.KeepEndTags, qt.Equals, true)
	c.Assert(conf.Tdewolff.CSS.Version, qt.Equals, 0)

	// `enable` flags
	c.Assert(conf.DisableHTML, qt.Equals, false)
	c.Assert(conf.DisableXML, qt.Equals, true)
}

func TestConfigDeprecations(t *testing.T) {
	c := qt.New(t)

	// Test default values of deprecated root keys.
	v := config.New()
	v.Set("minify", false)
	conf := testconfig.GetTestConfigs(nil, v).Base.Minify
	c.Assert(conf.MinifyOutput, qt.Equals, false)

	v = config.New()
	v.Set("minifyoutput", false)
	conf = testconfig.GetTestConfigs(nil, v).Base.Minify
	c.Assert(conf.MinifyOutput, qt.Equals, false)

	// Test non-default values of deprecated root keys.
	v = config.New()
	v.Set("minify", true)
	conf = testconfig.GetTestConfigs(nil, v).Base.Minify
	c.Assert(conf.MinifyOutput, qt.Equals, true)

	v = config.New()
	v.Set("minifyoutput", true)
	conf = testconfig.GetTestConfigs(nil, v).Base.Minify
	c.Assert(conf.MinifyOutput, qt.Equals, true)
}

func TestConfigUpstreamDeprecations(t *testing.T) {
	c := qt.New(t)

	// Test default values of deprecated keys.
	v := config.New()
	v.Set("minify", map[string]any{
		"tdewolff": map[string]any{
			"css": map[string]any{
				"decimals": 0,
				"keepcss2": true,
			},
			"html": map[string]any{
				"keepconditionalcomments": true,
			},
			"svg": map[string]any{
				"decimals": 0,
			},
		},
	})

	conf := testconfig.GetTestConfigs(nil, v).Base.Minify

	c.Assert(conf.Tdewolff.CSS.Precision, qt.Equals, 0)
	c.Assert(conf.Tdewolff.CSS.Version, qt.Equals, 2)
	c.Assert(conf.Tdewolff.HTML.KeepSpecialComments, qt.Equals, true)
	c.Assert(conf.Tdewolff.SVG.Precision, qt.Equals, 0)

	// Test non-default values of deprecated keys.
	v = config.New()
	v.Set("minify", map[string]any{
		"tdewolff": map[string]any{
			"css": map[string]any{
				"decimals": 6,
				"keepcss2": false,
			},
			"html": map[string]any{
				"keepconditionalcomments": false,
			},
			"svg": map[string]any{
				"decimals": 7,
			},
		},
	})

	conf = testconfig.GetTestConfigs(nil, v).Base.Minify

	c.Assert(conf.Tdewolff.CSS.Precision, qt.Equals, 6)
	c.Assert(conf.Tdewolff.CSS.Version, qt.Equals, 0)
	c.Assert(conf.Tdewolff.HTML.KeepSpecialComments, qt.Equals, false)
	c.Assert(conf.Tdewolff.SVG.Precision, qt.Equals, 7)
}
