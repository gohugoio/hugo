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

	"github.com/gohugoio/hugo/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestDecodeConfigFromTOML(t *testing.T) {
	assert := require.New(t)

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
	assert.NoError(err)

	config, err := DecodeConfig(cfg)
	assert.NoError(err)
	assert.NotNil(config)

	assert.Equal("DS", config.Disqus.Shortname)
	assert.Equal("ga_id", config.GoogleAnalytics.ID)

	assert.True(config.Instagram.DisableInlineCSS)
}

// Support old root-level GA settings etc.
func TestUseSettingsFromRootIfSet(t *testing.T) {
	assert := require.New(t)

	cfg := viper.New()
	cfg.Set("disqusShortname", "root_short")
	cfg.Set("googleAnalytics", "ga_root")

	config, err := DecodeConfig(cfg)
	assert.NoError(err)
	assert.NotNil(config)

	assert.Equal("root_short", config.Disqus.Shortname)
	assert.Equal("ga_root", config.GoogleAnalytics.ID)

}
