// Copyright 2018 The Hugo Authors. All rights reserved.
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

package services

import (
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/config"
	"github.com/spf13/viper"
)

func TestDecodeConfigFromTOML(t *testing.T) {
	c := qt.New(t)

	tomlConfig := `

someOtherValue = "foo"

[services]
[services.disqus]
shortname = "DS"
[services.googleAnalytics]
id = "ga_id"
[services.instagram]
disableInlineCSS = true
[services.twitter]
disableInlineCSS = true
`
	cfg, err := config.FromConfigString(tomlConfig, "toml")
	c.Assert(err, qt.IsNil)

	config, err := DecodeConfig(cfg)
	c.Assert(err, qt.IsNil)
	c.Assert(config, qt.Not(qt.IsNil))

	c.Assert(config.Disqus.Shortname, qt.Equals, "DS")
	c.Assert(config.GoogleAnalytics.ID, qt.Equals, "ga_id")

	c.Assert(config.Instagram.DisableInlineCSS, qt.Equals, true)
}

// Support old root-level GA settings etc.
func TestUseSettingsFromRootIfSet(t *testing.T) {
	c := qt.New(t)

	cfg := viper.New()
	cfg.Set("disqusShortname", "root_short")
	cfg.Set("googleAnalytics", "ga_root")

	config, err := DecodeConfig(cfg)
	c.Assert(err, qt.IsNil)
	c.Assert(config, qt.Not(qt.IsNil))

	c.Assert(config.Disqus.Shortname, qt.Equals, "root_short")
	c.Assert(config.GoogleAnalytics.ID, qt.Equals, "ga_root")

}
