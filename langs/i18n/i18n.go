// Copyright 2022 The Hugo Authors. All rights reserved.
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
	"reflect"
	"strings"

	"github.com/spf13/cast"

	"github.com/gohugoio/hugo/common/hreflect"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/langs"

	"github.com/gohugoio/go-i18n/v2/i18n"
)

type translator struct {
	translate      func(translationID string, templateData any) string
	hasTranslation func(translationID string) bool
}

var nopTranslator = translator{}

func (t translator) Translate(translationID string, templateData any) string {
	if t.translate == nil {
		return ""
	}
	return t.translate(translationID, templateData)
}

func (t translator) HasTranslation(translationID string) bool {
	if t.hasTranslation == nil {
		return false
	}
	return t.hasTranslation(translationID)
}

type translateFunc func(translationID string, templateData any) string

var i18nWarningLogger = helpers.NewDistinctErrorLogger()

// Translators handles i18n translations.
type Translators struct {
	translators map[string]langs.Translator
	cfg         config.Provider
	logger      loggers.Logger
}

// NewTranslator creates a new Translator for the given language bundle and configuration.
func NewTranslator(b *i18n.Bundle, cfg config.Provider, logger loggers.Logger) Translators {
	t := Translators{cfg: cfg, logger: logger, translators: make(map[string]langs.Translator)}
	t.initFuncs(b)
	return t
}

// Get gets the Translator for the given language, or for the default
// configured language if not found.
func (ts Translators) Get(lang string) langs.Translator {
	if t, ok := ts.translators[lang]; ok {
		return t
	}
	ts.logger.Infof("Translation func for language %v not found, use default.", lang)
	if tt, ok := ts.translators[ts.cfg.GetString("defaultContentLanguage")]; ok {
		return tt
	}

	ts.logger.Infoln("i18n not initialized; if you need string translations, check that you have a bundle in /i18n that matches the site language or the default language.")

	return nopTranslator
}

func (ts Translators) initFuncs(bndl *i18n.Bundle) {
	enableMissingTranslationPlaceholders := ts.cfg.GetBool("enableMissingTranslationPlaceholders")
	for _, lang := range bndl.LanguageTags() {
		currentLang := lang
		currentLangStr := currentLang.String()
		// This may be pt-BR; make it case insensitive.
		currentLangKey := strings.ToLower(strings.TrimPrefix(currentLangStr, artificialLangTagPrefix))
		localizer := i18n.NewLocalizer(bndl, currentLangStr)

		translate := func(translationID string, templateData any) (string, error) {
			pluralCount := getPluralCount(templateData)

			if templateData != nil {
				tp := reflect.TypeOf(templateData)
				if hreflect.IsInt(tp.Kind()) {
					// This was how go-i18n worked in v1,
					// and we keep it like this to avoid breaking
					// lots of sites in the wild.
					templateData = intCount(cast.ToInt(templateData))
				}
			}

			translated, translatedLang, err := localizer.LocalizeWithTag(&i18n.LocalizeConfig{
				MessageID:    translationID,
				TemplateData: templateData,
				PluralCount:  pluralCount,
			})

			sameLang := currentLang == translatedLang

			if err == nil && sameLang {
				return translated, nil
			}

			if err != nil && sameLang && translated != "" {
				// See #8492
				// TODO(bep) this needs to be improved/fixed upstream,
				// but currently we get an error even if the fallback to
				// "other" succeeds.
				if fmt.Sprintf("%T", err) == "i18n.pluralFormNotFoundError" {
					return translated, nil
				}
			}

			return translated, err

		}

		translateAndLogIfNeeded := func(translationID string, templateData any) string {
			translated, err := translate(translationID, templateData)
			if err == nil {
				return translated
			}

			if _, ok := err.(*i18n.MessageNotFoundErr); !ok {
				ts.logger.Warnf("Failed to get translated string for language %q and ID %q: %s", currentLangStr, translationID, err)
			}

			if ts.cfg.GetBool("logI18nWarnings") {
				i18nWarningLogger.Printf("i18n|MISSING_TRANSLATION|%s|%s", currentLangStr, translationID)
			}

			if enableMissingTranslationPlaceholders {
				return "[i18n] " + translationID
			}
			return translated
		}

		ts.translators[currentLangKey] = translator{
			translate: translateAndLogIfNeeded,
			hasTranslation: func(translationID string) bool {
				_, err := translate(translationID, nil)
				return err == nil
			},
		}
	}
}

// intCount wraps the Count method.
type intCount int

func (c intCount) Count() int {
	return int(c)
}

const countFieldName = "Count"

// getPluralCount gets the plural count as a string (floats) or an integer.
// If v is nil, nil is returned.
func getPluralCount(v any) any {
	if v == nil {
		// i18n called without any argument, make sure it does not
		// get any plural count.
		return nil
	}

	switch v := v.(type) {
	case map[string]any:
		for k, vv := range v {
			if strings.EqualFold(k, countFieldName) {
				return toPluralCountValue(vv)
			}
		}
	default:
		vv := reflect.Indirect(reflect.ValueOf(v))
		if vv.Kind() == reflect.Interface && !vv.IsNil() {
			vv = vv.Elem()
		}
		tp := vv.Type()

		if tp.Kind() == reflect.Struct {
			f := vv.FieldByName(countFieldName)
			if f.IsValid() {
				return toPluralCountValue(f.Interface())
			}
			m := hreflect.GetMethodByName(vv, countFieldName)
			if m.IsValid() && m.Type().NumIn() == 0 && m.Type().NumOut() == 1 {
				c := m.Call(nil)
				return toPluralCountValue(c[0].Interface())
			}
		}
	}

	return toPluralCountValue(v)
}

// go-i18n expects floats to be represented by string.
func toPluralCountValue(in any) any {
	k := reflect.TypeOf(in).Kind()
	switch {
	case hreflect.IsFloat(k):
		f := cast.ToString(in)
		if !strings.Contains(f, ".") {
			f += ".0"
		}
		return f
	case k == reflect.String:
		if _, err := cast.ToFloat64E(in); err == nil {
			return in
		}
		// A non-numeric value.
		return nil
	default:
		if i, err := cast.ToIntE(in); err == nil {
			return i
		}
		return nil
	}
}
