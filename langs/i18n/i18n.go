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
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/cast"

	"github.com/gohugoio/hugo/common/hreflect"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/resources/page"

	"github.com/gohugoio/go-i18n/v2/i18n"
)

type translateFunc func(ctx context.Context, translationID string, templateData any) string

// Translator handles i18n translations.
type Translator struct {
	translateFuncs map[string]translateFunc
	cfg            config.AllProvider
	logger         loggers.Logger
}

// NewTranslator creates a new Translator for the given language bundle and configuration.
func NewTranslator(b *i18n.Bundle, cfg config.AllProvider, logger loggers.Logger) Translator {
	t := Translator{cfg: cfg, logger: logger, translateFuncs: make(map[string]translateFunc)}
	t.initFuncs(b)
	return t
}

// Func gets the translate func for the given language, or for the default
// configured language if not found.
func (t Translator) Func(lang string) translateFunc {
	if f, ok := t.translateFuncs[lang]; ok {
		return f
	}
	t.logger.Infof("Translation func for language %v not found, use default.", lang)
	if f, ok := t.translateFuncs[t.cfg.DefaultContentLanguage()]; ok {
		return f
	}

	t.logger.Infoln("i18n not initialized; if you need string translations, check that you have a bundle in /i18n that matches the site language or the default language.")
	return func(ctx context.Context, translationID string, args any) string {
		return ""
	}
}

func (t Translator) initFuncs(bndl *i18n.Bundle) {
	enableMissingTranslationPlaceholders := t.cfg.EnableMissingTranslationPlaceholders()
	for _, lang := range bndl.LanguageTags() {
		currentLang := lang
		currentLangStr := currentLang.String()
		// This may be pt-BR; make it case insensitive.
		currentLangKey := strings.ToLower(strings.TrimPrefix(currentLangStr, artificialLangTagPrefix))
		localizer := i18n.NewLocalizer(bndl, currentLangStr)
		t.translateFuncs[currentLangKey] = func(ctx context.Context, translationID string, templateData any) string {
			pluralCount := getPluralCount(templateData)

			if templateData != nil {
				tp := reflect.TypeOf(templateData)
				if hreflect.IsInt(tp.Kind()) {
					// This was how go-i18n worked in v1,
					// and we keep it like this to avoid breaking
					// lots of sites in the wild.
					templateData = intCount(cast.ToInt(templateData))
				} else {
					if p, ok := templateData.(page.Page); ok {
						// See issue 10782.
						// The i18n has its own template handling and does not know about
						// the context.Context.
						// A common pattern is to pass Page to i18n, and use .ReadingTime etc.
						// We need to improve this, but that requires some upstream changes.
						// For now, just create a wrapper.
						templateData = page.PageWithContext{Page: p, Ctx: ctx}
					}
				}
			}

			translated, translatedLang, err := localizer.LocalizeWithTag(&i18n.LocalizeConfig{
				MessageID:    translationID,
				TemplateData: templateData,
				PluralCount:  pluralCount,
			})

			sameLang := currentLang == translatedLang

			if err == nil && sameLang {
				return translated
			}

			if err != nil && sameLang && translated != "" {
				// See #8492
				// TODO(bep) this needs to be improved/fixed upstream,
				// but currently we get an error even if the fallback to
				// "other" succeeds.
				if fmt.Sprintf("%T", err) == "i18n.pluralFormNotFoundError" {
					return translated
				}
			}

			if _, ok := err.(*i18n.MessageNotFoundErr); !ok {
				t.logger.Warnf("Failed to get translated string for language %q and ID %q: %s", currentLangStr, translationID, err)
			}

			if t.cfg.PrintI18nWarnings() {
				t.logger.Warnf("i18n|MISSING_TRANSLATION|%s|%s", currentLangStr, translationID)
			}

			if enableMissingTranslationPlaceholders {
				return "[i18n] " + translationID
			}

			return translated
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
