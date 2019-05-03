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

package deploy

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

[deployment]

order = ["o1", "o2"]

[[deployment.targets]]
Name = "name1"
URL = "url1"
CloudFrontDistributionID = "cdn1"

[[deployment.targets]]
name = "name2"
url = "url2"
cloudfrontdistributionid = "cdn2"

[[deployment.matchers]]
Pattern = "^pattern1$"
Cache-Control = "cachecontrol1"
Content-Encoding = "contentencoding1"
Content-Type = "contenttype1"
Gzip = true
Force = true

[[deployment.matchers]]
pattern = "^pattern2$"
cache-control = "cachecontrol2"
content-encoding = "contentencoding2"
content-type = "contenttype2"
`
	cfg, err := config.FromConfigString(tomlConfig, "toml")
	assert.NoError(err)

	dcfg, err := decodeConfig(cfg)
	assert.NoError(err)

	assert.Equal(2, len(dcfg.Order))
	assert.Equal("o1", dcfg.Order[0])
	assert.Equal("o2", dcfg.Order[1])
	assert.Equal(2, len(dcfg.ordering))

	assert.Equal(2, len(dcfg.Targets))
	assert.Equal("name1", dcfg.Targets[0].Name)
	assert.Equal("url1", dcfg.Targets[0].URL)
	assert.Equal("cdn1", dcfg.Targets[0].CloudFrontDistributionID)
	assert.Equal("name2", dcfg.Targets[1].Name)
	assert.Equal("url2", dcfg.Targets[1].URL)
	assert.Equal("cdn2", dcfg.Targets[1].CloudFrontDistributionID)

	assert.Equal(2, len(dcfg.Matchers))
	assert.Equal("^pattern1$", dcfg.Matchers[0].Pattern)
	assert.NotNil(dcfg.Matchers[0].re)
	assert.Equal("cachecontrol1", dcfg.Matchers[0].CacheControl)
	assert.Equal("contentencoding1", dcfg.Matchers[0].ContentEncoding)
	assert.Equal("contenttype1", dcfg.Matchers[0].ContentType)
	assert.True(dcfg.Matchers[0].Gzip)
	assert.True(dcfg.Matchers[0].Force)
	assert.Equal("^pattern2$", dcfg.Matchers[1].Pattern)
	assert.NotNil(dcfg.Matchers[1].re)
	assert.Equal("cachecontrol2", dcfg.Matchers[1].CacheControl)
	assert.Equal("contentencoding2", dcfg.Matchers[1].ContentEncoding)
	assert.Equal("contenttype2", dcfg.Matchers[1].ContentType)
	assert.False(dcfg.Matchers[1].Gzip)
	assert.False(dcfg.Matchers[1].Force)
}

func TestInvalidOrderingPattern(t *testing.T) {
	assert := require.New(t)

	tomlConfig := `

someOtherValue = "foo"

[deployment]
order = ["["]  # invalid regular expression
`
	cfg, err := config.FromConfigString(tomlConfig, "toml")
	assert.NoError(err)

	_, err = decodeConfig(cfg)
	assert.Error(err)
}

func TestInvalidMatcherPattern(t *testing.T) {
	assert := require.New(t)

	tomlConfig := `

someOtherValue = "foo"

[deployment]
[[deployment.matchers]]
Pattern = "["  # invalid regular expression
`
	cfg, err := config.FromConfigString(tomlConfig, "toml")
	assert.NoError(err)

	_, err = decodeConfig(cfg)
	assert.Error(err)
}

func TestDecodeConfigDefault(t *testing.T) {
	assert := require.New(t)

	dcfg, err := decodeConfig(viper.New())
	assert.NoError(err)
	assert.Equal(0, len(dcfg.Targets))
	assert.Equal(0, len(dcfg.Matchers))
}
