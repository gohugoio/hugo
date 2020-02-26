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

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/config"
	"github.com/spf13/viper"
)

func TestDecodeConfigFromTOML(t *testing.T) {
	c := qt.New(t)

	tomlConfig := `

someOtherValue = "foo"

[deployment]

order = ["o1", "o2"]

# All lowercase.
[[deployment.targets]]
name = "name0"
url = "url0"
cloudfrontdistributionid = "cdn0"
include = "*.html"

# All uppercase.
[[deployment.targets]]
NAME = "name1"
URL = "url1"
CLOUDFRONTDISTRIBUTIONID = "cdn1"
INCLUDE = "*.jpg"

# Camelcase.
[[deployment.targets]]
name = "name2"
url = "url2"
cloudFrontDistributionID = "cdn2"
exclude = "*.png"

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
	c.Assert(err, qt.IsNil)

	dcfg, err := decodeConfig(cfg)
	c.Assert(err, qt.IsNil)

	// Order.
	c.Assert(len(dcfg.Order), qt.Equals, 2)
	c.Assert(dcfg.Order[0], qt.Equals, "o1")
	c.Assert(dcfg.Order[1], qt.Equals, "o2")
	c.Assert(len(dcfg.ordering), qt.Equals, 2)

	// Targets.
	c.Assert(len(dcfg.Targets), qt.Equals, 3)
	wantInclude := []string{"*.html", "*.jpg", ""}
	wantExclude := []string{"", "", "*.png"}
	for i := 0; i < 3; i++ {
		tgt := dcfg.Targets[i]
		c.Assert(tgt.Name, qt.Equals, fmt.Sprintf("name%d", i))
		c.Assert(tgt.URL, qt.Equals, fmt.Sprintf("url%d", i))
		c.Assert(tgt.CloudFrontDistributionID, qt.Equals, fmt.Sprintf("cdn%d", i))
		c.Assert(tgt.Include, qt.Equals, wantInclude[i])
		if wantInclude[i] != "" {
			c.Assert(tgt.includeGlob, qt.Not(qt.IsNil))
		}
		c.Assert(tgt.Exclude, qt.Equals, wantExclude[i])
		if wantExclude[i] != "" {
			c.Assert(tgt.excludeGlob, qt.Not(qt.IsNil))
		}
	}

	// Matchers.
	c.Assert(len(dcfg.Matchers), qt.Equals, 3)
	for i := 0; i < 3; i++ {
		m := dcfg.Matchers[i]
		c.Assert(m.Pattern, qt.Equals, fmt.Sprintf("^pattern%d$", i))
		c.Assert(m.re, qt.Not(qt.IsNil))
		c.Assert(m.CacheControl, qt.Equals, fmt.Sprintf("cachecontrol%d", i))
		c.Assert(m.ContentEncoding, qt.Equals, fmt.Sprintf("contentencoding%d", i))
		c.Assert(m.ContentType, qt.Equals, fmt.Sprintf("contenttype%d", i))
		c.Assert(m.Gzip, qt.Equals, i != 0)
		c.Assert(m.Force, qt.Equals, i != 0)
	}
}

func TestInvalidOrderingPattern(t *testing.T) {
	c := qt.New(t)

	tomlConfig := `

someOtherValue = "foo"

[deployment]
order = ["["]  # invalid regular expression
`
	cfg, err := config.FromConfigString(tomlConfig, "toml")
	c.Assert(err, qt.IsNil)

	_, err = decodeConfig(cfg)
	c.Assert(err, qt.Not(qt.IsNil))
}

func TestInvalidMatcherPattern(t *testing.T) {
	c := qt.New(t)

	tomlConfig := `

someOtherValue = "foo"

[deployment]
[[deployment.matchers]]
Pattern = "["  # invalid regular expression
`
	cfg, err := config.FromConfigString(tomlConfig, "toml")
	c.Assert(err, qt.IsNil)

	_, err = decodeConfig(cfg)
	c.Assert(err, qt.Not(qt.IsNil))
}

func TestDecodeConfigDefault(t *testing.T) {
	c := qt.New(t)

	dcfg, err := decodeConfig(viper.New())
	c.Assert(err, qt.IsNil)
	c.Assert(len(dcfg.Targets), qt.Equals, 0)
	c.Assert(len(dcfg.Matchers), qt.Equals, 0)
}
