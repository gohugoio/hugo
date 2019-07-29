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

	"github.com/gohugoio/hugo/common/herrors"

	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/source"
	"github.com/nicksnyder/go-i18n/i18n/bundle"
	"github.com/nicksnyder/go-i18n/i18n/language"
	_errors "github.com/pkg/errors"
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
	spec := source.NewSourceSpec(d.PathSpec, nil)

	i18nBundle := bundle.New()

	en := language.GetPluralSpec("en")
	if en == nil {
		return errors.New("the English language has vanished like an old oak table")
	}
	var newLangs []string

	// The source dirs are ordered so the most important comes first. Since this is a
	// last key win situation, we have to reverse the iteration order.
	dirs := d.BaseFs.I18n.Dirs
	for i := len(dirs) - 1; i >= 0; i-- {
		dir := dirs[i]
		src := spec.NewFilesystemFromFileMetaInfo(dir)

		files, err := src.Files()
		if err != nil {
			return err
		}

		for _, r := range files {
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

		for _, file := range files {
			if err := addTranslationFile(i18nBundle, file); err != nil {
				return err
			}
		}
	}

	tp.t = NewTranslator(i18nBundle, d.Cfg, d.Log)

	d.Translate = tp.t.Func(d.Language.Lang)

	return nil

}

func addTranslationFile(bundle *bundle.Bundle, r source.File) error {
	f, err := r.FileInfo().Meta().Open()
	if err != nil {
		return _errors.Wrapf(err, "failed to open translations file %q:", r.LogicalName())
	}
	err = bundle.ParseTranslationFileBytes(r.LogicalName(), helpers.ReaderToBytes(f))
	f.Close()
	if err != nil {
		return errWithFileContext(_errors.Wrapf(err, "failed to load translations"), r)
	}
	return nil
}

// Clone sets the language func for the new language.
func (tp *TranslationProvider) Clone(d *deps.Deps) error {
	d.Translate = tp.t.Func(d.Language.Lang)

	return nil
}

func errWithFileContext(inerr error, r source.File) error {
	fim, ok := r.FileInfo().(hugofs.FileMetaInfo)
	if !ok {
		return inerr
	}

	meta := fim.Meta()
	realFilename := meta.Filename()
	f, err := meta.Open()
	if err != nil {
		return inerr
	}
	defer f.Close()

	err, _ = herrors.WithFileContext(
		inerr,
		realFilename,
		f,
		herrors.SimpleLineMatcher)

	return err

}
