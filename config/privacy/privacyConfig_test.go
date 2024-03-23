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

package privacy

import (
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/config"
)

func TestDecodeConfigFromTOML(t *testing.T) {
	c := qt.New(t)

	tomlConfig := `

someOtherValue = "foo"

[privacy]
[privacy.disqus]
disable = true
[privacy.googleAnalytics]
disable = true
respectDoNotTrack = true
[privacy.instagram]
disable = true
simple = true
[privacy.twitter]
disable = true
enableDNT = true
simple = true
[privacy.vimeo]
disable = true
enableDNT = true
simple = true
[privacy.youtube]
disable = true
privacyEnhanced = true
simple = true
`
	cfg, err := config.FromConfigString(tomlConfig, "toml")
	c.Assert(err, qt.IsNil)

	pc, err := DecodeConfig(cfg)
	c.Assert(err, qt.IsNil)
	c.Assert(pc, qt.Not(qt.IsNil))

	got := []bool{
		pc.Disqus.Disable, pc.GoogleAnalytics.Disable,
		pc.GoogleAnalytics.RespectDoNotTrack, pc.Instagram.Disable,
		pc.Instagram.Simple, pc.Twitter.Disable, pc.Twitter.EnableDNT,
		pc.Twitter.Simple, pc.Vimeo.Disable, pc.Vimeo.EnableDNT, pc.Vimeo.Simple,
		pc.YouTube.PrivacyEnhanced, pc.YouTube.Disable,
	}

	c.Assert(got, qt.All(qt.Equals), true)
}

func TestDecodeConfigFromTOMLCaseInsensitive(t *testing.T) {
	c := qt.New(t)

	tomlConfig := `

someOtherValue = "foo"

[Privacy]
[Privacy.YouTube]
PrivacyENhanced = true
`
	cfg, err := config.FromConfigString(tomlConfig, "toml")
	c.Assert(err, qt.IsNil)

	pc, err := DecodeConfig(cfg)
	c.Assert(err, qt.IsNil)
	c.Assert(pc, qt.Not(qt.IsNil))
	c.Assert(pc.YouTube.PrivacyEnhanced, qt.Equals, true)
}

func TestDecodeConfigDefault(t *testing.T) {
	c := qt.New(t)

	pc, err := DecodeConfig(config.New())
	c.Assert(err, qt.IsNil)
	c.Assert(pc, qt.Not(qt.IsNil))
	c.Assert(pc.YouTube.PrivacyEnhanced, qt.Equals, false)
}
