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
	"fmt"
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

# All lowercase.
[[deployment.targets]]
name = "name0"
url = "url0"
cloudfrontdistributionid = "cdn0"

# All uppercase.
[[deployment.targets]]
NAME = "name1"
URL = "url1"
CLOUDFRONTDISTRIBUTIONID = "cdn1"

# Camelcase.
[[deployment.targets]]
name = "name2"
url = "url2"
cloudFrontDistributionID = "cdn2"

# All lowercase.
[[deployment.matchers]]
pattern = "^pattern0$"
cachecontrol = "cachecontrol0"
contentencoding = "contentencoding0"
contenttype = "contenttype0"

# All uppercase.
[[deployment.matchers]]
PATTERN = "^pattern1$"
CACHECONTROL = "cachecontrol1"
CONTENTENCODING = "contentencoding1"
CONTENTTYPE = "contenttype1"
GZIP = true
FORCE = true

# Camelcase.
[[deployment.matchers]]
pattern = "^pattern2$"
cacheControl = "cachecontrol2"
contentEncoding = "contentencoding2"
contentType = "contenttype2"
gzip = true
force = true
`
	cfg, err := config.FromConfigString(tomlConfig, "toml")
	assert.NoError(err)

	dcfg, err := decodeConfig(cfg)
	assert.NoError(err)

	// Order.
	assert.Equal(2, len(dcfg.Order))
	assert.Equal("o1", dcfg.Order[0])
	assert.Equal("o2", dcfg.Order[1])
	assert.Equal(2, len(dcfg.ordering))

	// Targets.
	assert.Equal(3, len(dcfg.Targets))
	for i := 0; i < 3; i++ {
		tgt := dcfg.Targets[i]
		assert.Equal(fmt.Sprintf("name%d", i), tgt.Name)
		assert.Equal(fmt.Sprintf("url%d", i), tgt.URL)
		assert.Equal(fmt.Sprintf("cdn%d", i), tgt.CloudFrontDistributionID)
	}

	// Matchers.
	assert.Equal(3, len(dcfg.Matchers))
	for i := 0; i < 3; i++ {
		m := dcfg.Matchers[i]
		assert.Equal(fmt.Sprintf("^pattern%d$", i), m.Pattern)
		assert.NotNil(m.re)
		assert.Equal(fmt.Sprintf("cachecontrol%d", i), m.CacheControl)
		assert.Equal(fmt.Sprintf("contentencoding%d", i), m.ContentEncoding)
		assert.Equal(fmt.Sprintf("contenttype%d", i), m.ContentType)
		assert.Equal(i != 0, m.Gzip)
		assert.Equal(i != 0, m.Force)
	}
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
