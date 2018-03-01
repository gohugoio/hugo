// Copyright 2017 The Hugo Authors. All rights reserved.
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

package i18n

import (
	"errors"
	"fmt"

	"github.com/gohugoio/hugo/helpers"

	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/source"
	"github.com/nicksnyder/go-i18n/i18n/bundle"
	"github.com/nicksnyder/go-i18n/i18n/language"
)

// TranslationProvider provides translation handling, i.e. loading
// of bundles etc.
type TranslationProvider struct {
	t Translator
}

// NewTranslationProvider creates a new translation provider.
func NewTranslationProvider() *TranslationProvider {
	return &TranslationProvider{}
}

// Update updates the i18n func in the provided Deps.
func (tp *TranslationProvider) Update(d *deps.Deps) error {
	sp := source.NewSourceSpec(d.PathSpec, d.BaseFs.SourceFilesystems.I18n.Fs)
	src := sp.NewFilesystem("")

	i18nBundle := bundle.New()

	en := language.GetPluralSpec("en")
	if en == nil {
		return errors.New("The English language has vanished like an old oak table!")
	}
	var newLangs []string

	for _, r := range src.Files() {
		currentSpec := language.GetPluralSpec(r.BaseFileName())
		if currentSpec == nil {
			// This may is a language code not supported by go-i18n, it may be
			// Klingon or ... not even a fake language. Make sure it works.
			newLangs = append(newLangs, r.BaseFileName())
		}
	}

	if len(newLangs) > 0 {
		language.RegisterPluralSpec(newLangs, en)
	}

	// The source files are ordered so the most important comes first. Since this is a
	// last key win situation, we have to reverse the iteration order.
	files := src.Files()
	for i := len(files) - 1; i >= 0; i-- {
		if err := addTranslationFile(i18nBundle, files[i]); err != nil {
			return err
		}
	}

	tp.t = NewTranslator(i18nBundle, d.Cfg, d.Log)

	d.Translate = tp.t.Func(d.Language.Lang)

	return nil

}

func addTranslationFile(bundle *bundle.Bundle, r source.ReadableFile) error {
	f, err := r.Open()
	if err != nil {
		return fmt.Errorf("Failed to open translations file %q: %s", r.LogicalName(), err)
	}
	defer f.Close()
	err = bundle.ParseTranslationFileBytes(r.LogicalName(), helpers.ReaderToBytes(f))
	if err != nil {
		return fmt.Errorf("Failed to load translations in file %q: %s", r.LogicalName(), err)
	}
	return nil
}

// Clone sets the language func for the new language.
func (tp *TranslationProvider) Clone(d *deps.Deps) error {
	d.Translate = tp.t.Func(d.Language.Lang)

	return nil
}
