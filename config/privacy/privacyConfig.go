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
	"github.com/gohugoio/hugo/config"
	"github.com/mitchellh/mapstructure"
)

const privacyConfigKey = "privacy"

// Service is the common values for a service in a policy definition.
type Service struct {
	Disable bool
}

// Config is a privacy configuration for all the relevant services in Hugo.
type Config struct {
	Disqus          Disqus
	GoogleAnalytics GoogleAnalytics
	Instagram       Instagram
	SpeakerDeck     SpeakerDeck
	Tweet           Tweet
	Vimeo           Vimeo
	YouTube         YouTube
}

// Disqus holds the privacy configuration settings related to the Disqus template.
type Disqus struct {
	Service `mapstructure:",squash"`
}

// GoogleAnalytics holds the privacy configuration settings related to the Google Analytics template.
type GoogleAnalytics struct {
	Service `mapstructure:",squash"`

	// Enabling this will make the GA templates respect the
	// "Do Not Track" HTTP header. See  https://www.paulfurley.com/google-analytics-dnt/.
	RespectDoNotTrack bool
}

// Instagram holds the privacy configuration settings related to the Instagram shortcode.
type Instagram struct {
	Service `mapstructure:",squash"`

	// If simple mode is enabled, a static and no-JS version of the Instagram
	// image card will be built.
	Simple bool
}

// SpeakerDeck holds the privacy configuration settings related to the SpeakerDeck shortcode.
type SpeakerDeck struct {
	Service `mapstructure:",squash"`
}

// Tweet holds the privacy configuration settingsrelated to the Tweet shortcode.
type Tweet struct {
	Service `mapstructure:",squash"`
}

// Vimeo holds the privacy configuration settingsrelated to the Vimeo shortcode.
type Vimeo struct {
	Service `mapstructure:",squash"`
}

// YouTube holds the privacy configuration settingsrelated to the YouTube shortcode.
type YouTube struct {
	Service `mapstructure:",squash"`

	// When you turn on privacy-enhanced mode,
	// YouTube won’t store information about visitors on your website
	// unless the user plays the embedded video.
	PrivacyEnhanced bool

	// If simple mode is enabled, only a thumbnail is fetched from ytimg.com and
	// shown with a play button overlaid. If a user clicks the button, he/she will
	// be taken to the video page on youtube.com in a new browser tab.
	Simple bool
}

func DecodeConfig(cfg config.Provider) (pc Config, err error) {
	if !cfg.IsSet(privacyConfigKey) {
		return
	}

	m := cfg.GetStringMap(privacyConfigKey)

	err = mapstructure.WeakDecode(m, &pc)

	return
}
