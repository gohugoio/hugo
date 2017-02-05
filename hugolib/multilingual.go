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

	"sort"

	"errors"
	"fmt"

	"github.com/spf13/cast"
	"github.com/spf13/hugo/config"
	"github.com/spf13/hugo/helpers"
)

type Multilingual struct {
	Languages helpers.Languages

	DefaultLang *helpers.Language

	langMap     map[string]*helpers.Language
	langMapInit sync.Once
}

func (ml *Multilingual) Language(lang string) *helpers.Language {
	ml.langMapInit.Do(func() {
		ml.langMap = make(map[string]*helpers.Language)
		for _, l := range ml.Languages {
			ml.langMap[l.Lang] = l
		}
	})
	return ml.langMap[lang]
}

func newMultiLingualFromSites(cfg config.Provider, sites ...*Site) (*Multilingual, error) {
	languages := make(helpers.Languages, len(sites))

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

	return &Multilingual{Languages: languages, DefaultLang: helpers.NewLanguage(defaultLang, cfg)}, nil

}

func newMultiLingualForLanguage(language *helpers.Language) *Multilingual {
	languages := helpers.Languages{language}
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

func toSortedLanguages(cfg config.Provider, l map[string]interface{}) (helpers.Languages, error) {
	langs := make(helpers.Languages, len(l))
	i := 0

	for lang, langConf := range l {
		langsMap, err := cast.ToStringMapE(langConf)

		if err != nil {
			return nil, fmt.Errorf("Language config is not a map: %T", langConf)
		}

		language := helpers.NewLanguage(lang, cfg)

		for loki, v := range langsMap {
			switch loki {
			case "title":
				language.Title = cast.ToString(v)
			case "languagename":
				language.LanguageName = cast.ToString(v)
			case "weight":
				language.Weight = cast.ToInt(v)
			}

			// Put all into the Params map
			language.SetParam(loki, v)
		}

		langs[i] = language
		i++
	}

	sort.Sort(langs)

	return langs, nil
}
