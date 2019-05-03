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

package hugolib

import (
	"sync"

	"errors"

	"github.com/gohugoio/hugo/langs"

	"github.com/gohugoio/hugo/config"
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
		if s.language == nil {
			return nil, errors.New("missing language for site")
		}
		languages[i] = s.language
	}

	defaultLang := cfg.GetString("defaultContentLanguage")

	if defaultLang == "" {
		defaultLang = "en"
	}

	return &Multilingual{Languages: languages, DefaultLang: langs.NewLanguage(defaultLang, cfg)}, nil

}

func (ml *Multilingual) enabled() bool {
	return len(ml.Languages) > 1
}

func (s *Site) multilingualEnabled() bool {
	if s.h == nil {
		return false
	}
	return s.h.multilingual != nil && s.h.multilingual.enabled()
}
