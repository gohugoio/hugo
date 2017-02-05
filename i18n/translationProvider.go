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
	"fmt"

	"github.com/nicksnyder/go-i18n/i18n/bundle"
	"github.com/spf13/hugo/deps"
	"github.com/spf13/hugo/source"
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
	dir := d.PathSpec.AbsPathify(d.Cfg.GetString("i18nDir"))
	sp := source.NewSourceSpec(d.Cfg, d.Fs)
	sources := []source.Input{sp.NewFilesystem(dir)}

	themeI18nDir, err := d.PathSpec.GetThemeI18nDirPath()

	if err == nil {
		sources = []source.Input{sp.NewFilesystem(themeI18nDir), sources[0]}
	}

	d.Log.DEBUG.Printf("Load I18n from %q", sources)

	i18nBundle := bundle.New()

	for _, currentSource := range sources {
		for _, r := range currentSource.Files() {
			err := i18nBundle.ParseTranslationFileBytes(r.LogicalName(), r.Bytes())
			if err != nil {
				return fmt.Errorf("Failed to load translations in file %q: %s", r.LogicalName(), err)
			}
		}
	}

	tp.t = NewTranslator(i18nBundle, d.Cfg, d.Log)

	d.Translate = tp.t.Func(d.Language.Lang)

	return nil

}

// Clone sets the language func for the new language.
func (tp *TranslationProvider) Clone(d *deps.Deps) error {
	d.Translate = tp.t.Func(d.Language.Lang)

	return nil
}
