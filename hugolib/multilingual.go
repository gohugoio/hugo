// Copyright 2016-present The Hugo Authors. All rights reserved.
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

package hugolib

import (
	"sync"

	"github.com/gohugoio/hugo/common/maps"

	"sort"

	"errors"
	"fmt"

	"github.com/gohugoio/hugo/langs"

	"github.com/gohugoio/hugo/config"
	"github.com/spf13/cast"
)

// Multilingual manages the all languages used in a multilingual site.
type Multilingual struct {
	Languages langs.Languages

	DefaultLang *langs.Language

	langMap     map[string]*langs.Language
	langMapInit sync.Once
}

// Language returns the Language associated with the given string.
func (ml *Multilingual) Language(lang string) *langs.Language {
	ml.langMapInit.Do(func() {
		ml.langMap = make(map[string]*langs.Language)
		for _, l := range ml.Languages {
			ml.langMap[l.Lang] = l
		}
	})
	return ml.langMap[lang]
}

func getLanguages(cfg config.Provider) langs.Languages {
	if cfg.IsSet("languagesSorted") {
		return cfg.Get("languagesSorted").(langs.Languages)
	}

	return langs.Languages{langs.NewDefaultLanguage(cfg)}
}

func newMultiLingualFromSites(cfg config.Provider, sites ...*Site) (*Multilingual, error) {
	languages := make(langs.Languages, len(sites))

	for i, s := range sites {
		if s.Language == nil {
			return nil, errors.New("Missing language for site")
		}
		languages[i] = s.Language
	}

	defaultLang := cfg.GetString("defaultContentLanguage")

	if defaultLang == "" {
		defaultLang = "en"
	}

	return &Multilingual{Languages: languages, DefaultLang: langs.NewLanguage(defaultLang, cfg)}, nil

}

func newMultiLingualForLanguage(language *langs.Language) *Multilingual {
	languages := langs.Languages{language}
	return &Multilingual{Languages: languages, DefaultLang: language}
}
func (ml *Multilingual) enabled() bool {
	return len(ml.Languages) > 1
}

func (s *Site) multilingualEnabled() bool {
	if s.owner == nil {
		return false
	}
	return s.owner.multilingual != nil && s.owner.multilingual.enabled()
}

func toSortedLanguages(cfg config.Provider, l map[string]interface{}) (langs.Languages, error) {
	languages := make(langs.Languages, len(l))
	i := 0

	for lang, langConf := range l {
		langsMap, err := cast.ToStringMapE(langConf)

		if err != nil {
			return nil, fmt.Errorf("Language config is not a map: %T", langConf)
		}

		language := langs.NewLanguage(lang, cfg)

		for loki, v := range langsMap {
			switch loki {
			case "title":
				language.Title = cast.ToString(v)
			case "languagename":
				language.LanguageName = cast.ToString(v)
			case "weight":
				language.Weight = cast.ToInt(v)
			case "contentdir":
				language.ContentDir = cast.ToString(v)
			case "disabled":
				language.Disabled = cast.ToBool(v)
			case "params":
				m := cast.ToStringMap(v)
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
