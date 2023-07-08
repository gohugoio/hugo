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

package langs

import (
	"errors"

	"github.com/gohugoio/hugo/common/maps"
	"github.com/mitchellh/mapstructure"
)

// LanguageConfig holds the configuration for a single language.
// This is what is read from the config file.
type LanguageConfig struct {
	// The language name, e.g. "English".
	LanguageName string

	// The language code, e.g. "en-US".
	LanguageCode string

	// The language title. When set, this will
	// override site.Title for this language.
	Title string

	// The language direction, e.g. "ltr" or "rtl".
	LanguageDirection string

	// The language weight. When set to a non-zero value, this will
	// be the main sort criteria for the language.
	Weight int

	// Set to true to disable this language.
	Disabled bool
}

func DecodeConfig(m map[string]any) (map[string]LanguageConfig, error) {
	m = maps.CleanConfigStringMap(m)
	var langs map[string]LanguageConfig

	if err := mapstructure.WeakDecode(m, &langs); err != nil {
		return nil, err
	}
	if len(langs) == 0 {
		return nil, errors.New("no languages configured")
	}
	return langs, nil
}
