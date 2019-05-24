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

package services

import (
	"github.com/gohugoio/hugo/config"
	"github.com/mitchellh/mapstructure"
)

const (
	servicesConfigKey = "services"

	disqusShortnameKey = "disqusshortname"
	googleAnalyticsKey = "googleanalytics"
	rssLimitKey        = "rssLimit"
)

// Config is a privacy configuration for all the relevant services in Hugo.
type Config struct {
	Disqus          Disqus
	GoogleAnalytics GoogleAnalytics
	Instagram       Instagram
	Twitter         Twitter
	RSS             RSS
}

// Disqus holds the functional configuration settings related to the Disqus template.
type Disqus struct {
	// A Shortname is the unique identifier assigned to a Disqus site.
	Shortname string
}

// GoogleAnalytics holds the functional configuration settings related to the Google Analytics template.
type GoogleAnalytics struct {
	// The GA tracking ID.
	ID string
}

// Instagram holds the functional configuration settings related to the Instagram shortcodes.
type Instagram struct {
	// The Simple variant of the Instagram is decorated with Bootstrap 4 card classes.
	// This means that if you use Bootstrap 4 or want to provide your own CSS, you want
	// to disable the inline CSS provided by Hugo.
	DisableInlineCSS bool
}

// Twitter holds the functional configuration settings related to the Twitter shortcodes.
type Twitter struct {
	// The Simple variant of Twitter is decorated with a basic set of inline styles.
	// This means that if you want to provide your own CSS, you want
	// to disable the inline CSS provided by Hugo.
	DisableInlineCSS bool
}

// RSS holds the functional configuration settings related to the RSS feeds.
type RSS struct {
	// Limit the number of pages.
	Limit int
}

// DecodeConfig creates a services Config from a given Hugo configuration.
func DecodeConfig(cfg config.Provider) (c Config, err error) {
	m := cfg.GetStringMap(servicesConfigKey)

	err = mapstructure.WeakDecode(m, &c)

	// Keep backwards compatibility.
	if c.GoogleAnalytics.ID == "" {
		// Try the global config
		c.GoogleAnalytics.ID = cfg.GetString(googleAnalyticsKey)
	}
	if c.Disqus.Shortname == "" {
		c.Disqus.Shortname = cfg.GetString(disqusShortnameKey)
	}

	if c.RSS.Limit == 0 {
		c.RSS.Limit = cfg.GetInt(rssLimitKey)
	}

	return
}
