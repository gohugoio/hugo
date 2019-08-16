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

// Package helpers implements general utility functions that work with
// and on content.  The helper functions defined here lay down the
// foundation of how Hugo works with files and filepaths, and perform
// string operations on content.

package internal

import (
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/markup/converter"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

// BlackFriday holds configuration values for BlackFriday rendering.
// It is kept here because it's used in several packages.
type BlackFriday struct {
	Smartypants           bool
	SmartypantsQuotesNBSP bool
	AngledQuotes          bool
	Fractions             bool
	HrefTargetBlank       bool
	NofollowLinks         bool
	NoreferrerLinks       bool
	SmartDashes           bool
	LatexDashes           bool
	TaskLists             bool
	PlainIDAnchors        bool
	Extensions            []string
	ExtensionsMask        []string
	SkipHTML              bool

	FootnoteAnchorPrefix       string
	FootnoteReturnLinkContents string
}

func UpdateBlackFriday(old *BlackFriday, m map[string]interface{}) (*BlackFriday, error) {
	// Create a copy so we can modify it.
	bf := *old
	if err := mapstructure.Decode(m, &bf); err != nil {
		return nil, errors.WithMessage(err, "failed to decode rendering config")
	}
	return &bf, nil
}

// NewBlackfriday creates a new Blackfriday filled with site config or some sane defaults.
func NewBlackfriday(cfg converter.ProviderConfig) (*BlackFriday, error) {
	var siteConfig map[string]interface{}
	if cfg.Cfg != nil {
		siteConfig = cfg.Cfg.GetStringMap("blackfriday")
	}

	defaultParam := map[string]interface{}{
		"smartypants":           true,
		"angledQuotes":          false,
		"smartypantsQuotesNBSP": false,
		"fractions":             true,
		"hrefTargetBlank":       false,
		"nofollowLinks":         false,
		"noreferrerLinks":       false,
		"smartDashes":           true,
		"latexDashes":           true,
		"plainIDAnchors":        true,
		"taskLists":             true,
		"skipHTML":              false,
	}

	maps.ToLower(defaultParam)

	config := make(map[string]interface{})

	for k, v := range defaultParam {
		config[k] = v
	}

	for k, v := range siteConfig {
		config[k] = v
	}

	combinedConfig := &BlackFriday{}
	if err := mapstructure.Decode(config, combinedConfig); err != nil {
		return nil, errors.Errorf("failed to decode Blackfriday config: %s", err)
	}

	// TODO(bep) update/consolidate docs
	if combinedConfig.FootnoteAnchorPrefix == "" {
		combinedConfig.FootnoteAnchorPrefix = cfg.Cfg.GetString("footnoteAnchorPrefix")
	}

	if combinedConfig.FootnoteReturnLinkContents == "" {
		combinedConfig.FootnoteReturnLinkContents = cfg.Cfg.GetString("footnoteReturnLinkContents")
	}

	return combinedConfig, nil
}
