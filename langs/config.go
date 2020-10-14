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
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gohugoio/hugo/common/maps"

	"github.com/spf13/cast"

	"github.com/pkg/errors"

	"github.com/gohugoio/hugo/config"
)

type LanguagesConfig struct {
	Languages                      Languages
	Multihost                      bool
	DefaultContentLanguageInSubdir bool
}

func LoadLanguageSettings(cfg config.Provider, oldLangs Languages) (c LanguagesConfig, err error) {

	defaultLang := strings.ToLower(cfg.GetString("defaultContentLanguage"))
	if defaultLang == "" {
		defaultLang = "en"
		cfg.Set("defaultContentLanguage", defaultLang)
	}

	var languages map[string]interface{}

	languagesFromConfig := cfg.GetStringMap("languages")
	disableLanguages := cfg.GetStringSlice("disableLanguages")

	if len(disableLanguages) == 0 {
		languages = languagesFromConfig
	} else {
		languages = make(map[string]interface{})
		for k, v := range languagesFromConfig {
			for _, disabled := range disableLanguages {
				if disabled == defaultLang {
					return c, fmt.Errorf("cannot disable default language %q", defaultLang)
				}

				if strings.EqualFold(k, disabled) {
					v.(map[string]interface{})["disabled"] = true
					break
				}
			}
			languages[k] = v
		}
	}

	var languages2 Languages

	if len(languages) == 0 {
		languages2 = append(languages2, NewDefaultLanguage(cfg))
	} else {
		languages2, err = toSortedLanguages(cfg, languages)
		if err != nil {
			return c, errors.Wrap(err, "Failed to parse multilingual config")
		}
	}

	if oldLangs != nil {
		// When in multihost mode, the languages are mapped to a server, so
		// some structural language changes will need a restart of the dev server.
		// The validation below isn't complete, but should cover the most
		// important cases.
		var invalid bool
		if languages2.IsMultihost() != oldLangs.IsMultihost() {
			invalid = true
		} else {
			if languages2.IsMultihost() && len(languages2) != len(oldLangs) {
				invalid = true
			}
		}

		if invalid {
			return c, errors.New("language change needing a server restart detected")
		}

		if languages2.IsMultihost() {
			// We need to transfer any server baseURL to the new language
			for i, ol := range oldLangs {
				nl := languages2[i]
				nl.Set("baseURL", ol.GetString("baseURL"))
			}
		}
	}

	// The defaultContentLanguage is something the user has to decide, but it needs
	// to match a language in the language definition list.
	langExists := false
	for _, lang := range languages2 {
		if lang.Lang == defaultLang {
			langExists = true
			break
		}
	}

	if !langExists {
		return c, fmt.Errorf("site config value %q for defaultContentLanguage does not match any language definition", defaultLang)
	}

	c.Languages = languages2
	c.Multihost = languages2.IsMultihost()
	c.DefaultContentLanguageInSubdir = c.Multihost

	sortedDefaultFirst := make(Languages, len(c.Languages))
	for i, v := range c.Languages {
		sortedDefaultFirst[i] = v
	}
	sort.Slice(sortedDefaultFirst, func(i, j int) bool {
		li, lj := sortedDefaultFirst[i], sortedDefaultFirst[j]
		if li.Lang == defaultLang {
			return true
		}

		if lj.Lang == defaultLang {
			return false
		}

		return i < j
	})

	cfg.Set("languagesSorted", c.Languages)
	cfg.Set("languagesSortedDefaultFirst", sortedDefaultFirst)
	cfg.Set("multilingual", len(languages2) > 1)

	multihost := c.Multihost

	if multihost {
		cfg.Set("defaultContentLanguageInSubdir", true)
		cfg.Set("multihost", true)
	}

	if multihost {
		// The baseURL may be provided at the language level. If that is true,
		// then every language must have a baseURL. In this case we always render
		// to a language sub folder, which is then stripped from all the Permalink URLs etc.
		for _, l := range languages2 {
			burl := l.GetLocal("baseURL")
			if burl == nil {
				return c, errors.New("baseURL must be set on all or none of the languages")
			}
		}

	}

	return c, nil
}

func toSortedLanguages(cfg config.Provider, l map[string]interface{}) (Languages, error) {
	languages := make(Languages, len(l))
	i := 0

	for lang, langConf := range l {
		langsMap, err := maps.ToStringMapE(langConf)

		if err != nil {
			return nil, fmt.Errorf("Language config is not a map: %T", langConf)
		}

		language := NewLanguage(lang, cfg)

		for loki, v := range langsMap {
			switch loki {
			case "title":
				language.Title = cast.ToString(v)
			case "languagename":
				language.LanguageName = cast.ToString(v)
			case "languagedirection":
				language.LanguageDirection = cast.ToString(v)
			case "weight":
				language.Weight = cast.ToInt(v)
			case "contentdir":
				language.ContentDir = filepath.Clean(cast.ToString(v))
			case "disabled":
				language.Disabled = cast.ToBool(v)
			case "params":
				m := maps.ToStringMap(v)
				// Needed for case insensitive fetching of params values
				maps.ToLower(m)
				for k, vv := range m {
					language.SetParam(k, vv)
				}
			}

			// Put all into the Params map
			language.SetParam(loki, v)

			// Also set it in the configuration map (for baseURL etc.)
			language.Set(loki, v)
		}

		languages[i] = language
		i++
	}

	sort.Sort(languages)

	return languages, nil
}
