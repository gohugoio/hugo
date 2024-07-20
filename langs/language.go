// Copyright 2024 The Hugo Authors. All rights reserved.
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

// Package langs contains the language related types and function.
package langs

import (
	"fmt"
	"sync"
	"time"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"

	"github.com/gohugoio/hugo/common/htime"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/locales"
	translators "github.com/gohugoio/localescompressed"
)

type Language struct {
	// The language code, e.g. "en" or "no".
	// This is currently only settable as the key in the language map in the config.
	Lang string

	// Fields from the language config.
	LanguageConfig

	// Used for date formatting etc. We don't want these exported to the
	// templates.
	translator    locales.Translator
	timeFormatter htime.TimeFormatter
	tag           language.Tag
	// collator1 and collator2 are the same, we have 2 to prevent deadlocks.
	collator1 *Collator
	collator2 *Collator

	location *time.Location

	// This is just an alias of Site.Params.
	params maps.Params
}

// NewLanguage creates a new language.
func NewLanguage(lang, defaultContentLanguage, timeZone string, languageConfig LanguageConfig) (*Language, error) {
	translator := translators.GetTranslator(lang)
	if translator == nil {
		translator = translators.GetTranslator(defaultContentLanguage)
		if translator == nil {
			translator = translators.GetTranslator("en")
		}
	}

	var coll1, coll2 *Collator
	tag, err := language.Parse(lang)
	if err == nil {
		coll1 = &Collator{
			c: collate.New(tag),
		}
		coll2 = &Collator{
			c: collate.New(tag),
		}
	} else {
		coll1 = &Collator{
			c: collate.New(language.English),
		}
		coll2 = &Collator{
			c: collate.New(language.English),
		}
	}

	l := &Language{
		Lang:           lang,
		LanguageConfig: languageConfig,
		translator:     translator,
		timeFormatter:  htime.NewTimeFormatter(translator),
		tag:            tag,
		collator1:      coll1,
		collator2:      coll2,
	}

	return l, l.loadLocation(timeZone)
}

// This is injected from hugolib to avoid circular dependencies.
var DeprecationFunc = func(item, alternative string, err bool) {}

// Params returns the language params.
// Note that this is the same as the Site.Params, but we keep it here for legacy reasons.
// Deprecated: Use the site.Params instead.
func (l *Language) Params() maps.Params {
	// TODO(bep) Remove this for now as it created a little too much noise. Need to think about this.
	// See https://github.com/gohugoio/hugo/issues/11025
	// DeprecationFunc(".Language.Params", paramsDeprecationWarning, false)
	return l.params
}

func (l *Language) LanguageCode() string {
	if l.LanguageConfig.LanguageCode != "" {
		return l.LanguageConfig.LanguageCode
	}
	return l.Lang
}

func (l *Language) loadLocation(tzStr string) error {
	location, err := time.LoadLocation(tzStr)
	if err != nil {
		return fmt.Errorf("invalid timeZone for language %q: %w", l.Lang, err)
	}
	l.location = location

	return nil
}

func (l *Language) String() string {
	return l.Lang
}

// Languages is a sortable list of languages.
type Languages []*Language

func (l Languages) AsSet() map[string]bool {
	m := make(map[string]bool)
	for _, lang := range l {
		m[lang.Lang] = true
	}

	return m
}

// AsIndexSet returns a map with the language code as key and index in l as value.
func (l Languages) AsIndexSet() map[string]int {
	m := make(map[string]int)
	for i, lang := range l {
		m[lang.Lang] = i
	}

	return m
}

// Internal access to unexported Language fields.
// This construct is to prevent them from leaking to the templates.

func SetParams(l *Language, params maps.Params) {
	l.params = params
}

func GetTimeFormatter(l *Language) htime.TimeFormatter {
	return l.timeFormatter
}

func GetTranslator(l *Language) locales.Translator {
	return l.translator
}

func GetLocation(l *Language) *time.Location {
	return l.location
}

func GetCollator1(l *Language) *Collator {
	return l.collator1
}

func GetCollator2(l *Language) *Collator {
	return l.collator2
}

type Collator struct {
	sync.Mutex
	c *collate.Collator
}

// CompareStrings compares a and b.
// It returns -1 if a < b, 1 if a > b and 0 if a == b.
// Note that the Collator is not thread safe, so you may want
// to acquire a lock on it before calling this method.
func (c *Collator) CompareStrings(a, b string) int {
	return c.c.CompareString(a, b)
}
