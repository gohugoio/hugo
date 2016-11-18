// Copyright 2015 The Hugo Authors. All rights reserved.
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

package tpl

import (
	"github.com/nicksnyder/go-i18n/i18n/bundle"
	"github.com/spf13/hugo/helpers"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

var (
	// Logi18nWarnings set to true to print warnings about missing language strings
	Logi18nWarnings   bool
	i18nWarningLogger = helpers.NewDistinctFeedbackLogger()
	currentLanguage   *helpers.Language
)

type translate struct {
	translateFuncs map[string]bundle.TranslateFunc

	current bundle.TranslateFunc
}

var translator *translate

// SetTranslateLang sets the translations language to use during template processing.
// This construction is unfortunate, but the template system is currently global.
func SetTranslateLang(language *helpers.Language) error {
	currentLanguage = language
	if f, ok := translator.translateFuncs[language.Lang]; ok {
		translator.current = f
	} else {
		jww.WARN.Printf("Translation func for language %v not found, use default.", language.Lang)
		translator.current = translator.translateFuncs[viper.GetString("defaultContentLanguage")]
	}
	return nil
}

// SetI18nTfuncs sets the language bundle to be used for i18n.
func SetI18nTfuncs(bndl *bundle.Bundle) {
	translator = &translate{translateFuncs: make(map[string]bundle.TranslateFunc)}
	defaultContentLanguage := viper.GetString("defaultContentLanguage")
	var (
		defaultT bundle.TranslateFunc
		err      error
	)

	defaultT, err = bndl.Tfunc(defaultContentLanguage)

	if err != nil {
		jww.WARN.Printf("No translation bundle found for default language %q", defaultContentLanguage)
	}

	enableMissingTranslationPlaceholders := viper.GetBool("enableMissingTranslationPlaceholders")
	for _, lang := range bndl.LanguageTags() {
		currentLang := lang

		translator.translateFuncs[currentLang] = func(translationID string, args ...interface{}) string {
			tFunc, err := bndl.Tfunc(currentLang)
			if err != nil {
				jww.WARN.Printf("could not load translations for language %q (%s), will use default content language.\n", lang, err)
			} else if translated := tFunc(translationID, args...); translated != translationID {
				return translated
			}
			if Logi18nWarnings {
				i18nWarningLogger.Printf("i18n|MISSING_TRANSLATION|%s|%s", currentLang, translationID)
			}
			if enableMissingTranslationPlaceholders {
				return "[i18n] " + translationID
			}
			if defaultT != nil {
				if translated := defaultT(translationID, args...); translated != translationID {
					return translated
				}
			}
			return ""
		}
	}
}

func i18nTranslate(id string, args ...interface{}) (string, error) {
	if translator == nil || translator.current == nil {
		helpers.DistinctErrorLog.Printf("i18n not initialized, check that you have language file (in i18n) that matches the site language or the default language.")
		return "", nil
	}
	return translator.current(id, args...), nil
}
