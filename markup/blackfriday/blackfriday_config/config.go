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

package blackfriday_config

import (
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

// Default holds the default BlackFriday config.
// Do not change!
var Default = Config{
	Smartypants:           true,
	AngledQuotes:          false,
	SmartypantsQuotesNBSP: false,
	Fractions:             true,
	HrefTargetBlank:       false,
	NofollowLinks:         false,
	NoreferrerLinks:       false,
	SmartDashes:           true,
	LatexDashes:           true,
	PlainIDAnchors:        true,
	TaskLists:             true,
	SkipHTML:              false,
}

// Config holds configuration values for BlackFriday rendering.
// It is kept here because it's used in several packages.
type Config struct {
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

func UpdateConfig(b Config, m map[string]interface{}) (Config, error) {
	if err := mapstructure.Decode(m, &b); err != nil {
		return b, errors.WithMessage(err, "failed to decode rendering config")
	}
	return b, nil
}
