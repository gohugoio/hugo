// Copyright 2023 The Hugo Authors. All rights reserved.
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
	collator      *Collator
	location      *time.Location
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

	var coll *Collator
	tag, err := language.Parse(lang)
	if err == nil {
		coll = &Collator{
			c: collate.New(tag),
		}
	} else {
		coll = &Collator{
			c: collate.New(language.English),
		}
	}

	l := &Language{
		Lang:           lang,
		LanguageConfig: languageConfig,
		translator:     translator,
		timeFormatter:  htime.NewTimeFormatter(translator),
		tag:            tag,
		collator:       coll,
	}

	return l, l.loadLocation(timeZone)

}

func (l *Language) loadLocation(tzStr string) error {
	location, err := time.LoadLocation(tzStr)
	if err != nil {
		return fmt.Errorf("invalid timeZone for language %q: %w", l.Lang, err)
	}
	l.location = location

	return nil
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

func (l Languages) AsOrdinalSet() map[string]int {
	m := make(map[string]int)
	for i, lang := range l {
		m[lang.Lang] = i
	}

	return m
}

// Internal access to unexported Language fields.
// This construct is to prevent them from leaking to the templates.

func GetTimeFormatter(l *Language) htime.TimeFormatter {
	return l.timeFormatter
}

func GetTranslator(l *Language) locales.Translator {
	return l.translator
}

func GetLocation(l *Language) *time.Location {
	return l.location
}

func GetCollator(l *Language) *Collator {
	return l.collator
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
