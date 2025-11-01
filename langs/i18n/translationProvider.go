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
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/parser/metadecoders"

	"github.com/gohugoio/hugo/common/herrors"
	"golang.org/x/text/language"

	"github.com/gohugoio/hugo/helpers"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	toml "github.com/pelletier/go-toml/v2"

	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/source"
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
func (tp *TranslationProvider) NewResource(dst *deps.Deps) error {
	defaultLangTag, err := language.Parse(dst.Conf.DefaultContentLanguage())
	if err != nil {
		defaultLangTag = language.English
	}
	builder := newBundleBuilder(defaultLangTag)

	w := hugofs.NewWalkway(
		hugofs.WalkwayConfig{
			Fs:         dst.BaseFs.I18n.Fs,
			IgnoreFile: dst.SourceSpec.IgnoreFile,
			PathParser: dst.SourceSpec.Cfg.PathParser(),
			WalkFn: func(path string, info hugofs.FileMetaInfo) error {
				if info.IsDir() {
					return nil
				}
				return builder.addTranslationFile(source.NewFileInfo(info))
			},
		})

	if err := w.Walk(); err != nil {
		return err
	}

	bundle, err := builder.Build()
	if err != nil {
		return err
	}

	tp.t = newTranslator(bundle, dst.Conf, dst.Log)

	dst.Translate = tp.t.Func(dst.Conf.Language().Lang)

	return nil
}

func newBundleBuilder(defaultLangTag language.Tag) *bundleBuilder {
	b := i18n.NewBundle(defaultLangTag)

	b.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	b.RegisterUnmarshalFunc("yaml", metadecoders.UnmarshalYaml)
	b.RegisterUnmarshalFunc("yml", metadecoders.UnmarshalYaml)
	b.RegisterUnmarshalFunc("json", json.Unmarshal)

	bb := &bundle{
		b:              b,
		definedLangs:   make(map[language.Tag]bool),
		undefinedLangs: make(map[language.Tag]string),
	}

	return &bundleBuilder{b: bb}
}

type bundleBuilder struct {
	b *bundle

	// The Go i18n library we use does not support artificial language tags.
	// Store them away and add them later using available real language tags.
	undefinedLangs []*source.File
}

type bundle struct {
	b *i18n.Bundle

	// Bundled languages.
	definedLangs map[language.Tag]bool

	// Maps an arbitrary but real language tag to a Hugo artificial language key.
	undefinedLangs map[language.Tag]string
}

var errUndefinedLang = fmt.Errorf("undefined language")

func (b *bundleBuilder) Build() (*bundle, error) {
	const retries = 10
	for range retries {
		if len(b.undefinedLangs) == 0 {
			break
		}
		var undefinedLangs []*source.File
		for _, r := range b.undefinedLangs {
			name := r.LogicalName()
			lang := paths.Filename(name)
			var tag language.Tag
			// Find an unused language tag.
			for _, t := range languageTags {
				if !b.b.definedLangs[t] {
					tag = t
					break
				}
			}
			if tag == language.Und {
				return nil, fmt.Errorf("failed to resolve language for file %q", r.LogicalName())
			}
			ext := paths.Ext(name)
			name = tag.String() + ext
			if err := b.doAddTranslationFile(r, tag, name); err != nil {
				if err == errUndefinedLang {
					undefinedLangs = append(undefinedLangs, r)
					continue
				}
				return nil, err
			}
			b.b.undefinedLangs[tag] = lang
		}
		b.undefinedLangs = undefinedLangs
	}

	if len(b.undefinedLangs) != 0 {
		return nil, fmt.Errorf("failed to resolve languages for some translation files")
	}

	return b.b, nil
}

func (bb *bundleBuilder) addTranslationFile(r *source.File) error {
	name := r.LogicalName()
	lang := paths.Filename(name)
	const artificialLangTagPrefix = "art-x-"
	isArtificial := strings.HasPrefix(lang, artificialLangTagPrefix)
	var tag language.Tag
	if !isArtificial {
		tag = language.Make(lang)
	}
	if isArtificial || tag == language.Und {
		if len(strings.TrimPrefix(lang, artificialLangTagPrefix)) > 8 {
			// The upstream language.Matcher does not support private use subtags.
			// Prior to v0.152.0 we maintained a fork of go-i18n that worked around this,
			// but to get rid of that fork, we reworked this to use real language tags.
			// But we still want to preserve the tag validation.
			return fmt.Errorf("%q: language: tag is not well-formed", lang)
		}
		bb.undefinedLangs = append(bb.undefinedLangs, r)
		return nil
	}
	err := bb.doAddTranslationFile(r, tag, name)

	if err == errUndefinedLang {
		bb.undefinedLangs = append(bb.undefinedLangs, r)
		return nil
	}

	return err
}

// Note that name must include the file extension.
func (bb *bundleBuilder) doAddTranslationFile(r *source.File, tag language.Tag, name string) error {
	f, err := r.FileInfo().Meta().Open()
	if err != nil {
		return fmt.Errorf("failed to open translations file %q:: %w", r.LogicalName(), err)
	}

	b := helpers.ReaderToBytes(f)
	f.Close()

	_, err = bb.b.b.ParseMessageFileBytes(b, name)
	if err != nil {
		if strings.Contains(err.Error(), "no plural rule") {
			// https://github.com/gohugoio/hugo/issues/7798
			return errUndefinedLang
		}
		var guidance string
		if strings.Contains(err.Error(), "mixed with unreserved keys") {
			guidance = ": see the lang.Translate documentation for a list of reserved keys"
		}
		return errWithFileContext(fmt.Errorf("failed to load translations: %w%s", err, guidance), r)
	}

	bb.b.definedLangs[tag] = true

	return nil
}

// CloneResource sets the language func for the new language.
func (tp *TranslationProvider) CloneResource(dst, src *deps.Deps) error {
	dst.Translate = tp.t.Func(dst.Conf.Language().Lang)
	return nil
}

func errWithFileContext(inerr error, r *source.File) error {
	meta := r.FileInfo().Meta()
	realFilename := meta.Filename
	f, err := meta.Open()
	if err != nil {
		return inerr
	}
	defer f.Close()

	return herrors.NewFileErrorFromName(inerr, realFilename).UpdateContent(f, nil)
}

// A list of languages in no particular order.
var languageTags = []language.Tag{
	language.Georgian,
	language.Urdu,
	language.Vietnamese,
	language.Catalan,
	language.Swedish,
	language.Filipino,
	language.Icelandic,
	language.Punjabi,
	language.Persian,
	language.EuropeanPortuguese,
	language.BritishEnglish,
	language.Kannada,
	language.Ukrainian,
	language.EuropeanSpanish,
	language.Arabic,
	language.Kazakh,
	language.Hebrew,
	language.Danish,
	language.Serbian,
	language.SimplifiedChinese,
	language.Lithuanian,
	language.Gujarati,
	language.Italian,
	language.Russian,
	language.Macedonian,
	language.Burmese,
	language.Portuguese,
	language.Bengali,
	language.Swahili,
	language.Tamil,
	language.Zulu,
	language.Croatian,
	language.Dutch,
	language.Khmer,
	language.LatinAmericanSpanish,
	language.Japanese,
	language.AmericanEnglish,
	language.Azerbaijani,
	language.Turkish,
	language.Norwegian,
	language.TraditionalChinese,
	language.Hungarian,
	language.Finnish,
	language.Estonian,
	language.Lao,
	language.Marathi,
	language.Greek,
	language.Korean,
	language.Uzbek,
	language.Latvian,
	language.Nepali,
	language.Albanian,
	language.SerbianLatin,
	language.BrazilianPortuguese,
	language.Romanian,
	language.Chinese,
	language.Amharic,
	language.English,
	language.French,
	language.CanadianFrench,
	language.Indonesian,
	language.Malayalam,
	language.Slovak,
	language.Slovenian,
	language.Telugu,
	language.Thai,
	language.Sinhala,
	language.Armenian,
	language.Czech,
	language.German,
	language.Polish,
	language.Spanish,
	language.Malay,
	language.Mongolian,
	language.Afrikaans,
	language.ModernStandardArabic,
	language.Bulgarian,
	language.Hindi,
	language.Kirghiz,
}
