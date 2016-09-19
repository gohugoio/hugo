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
	"testing"

	"github.com/nicksnyder/go-i18n/i18n/bundle"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

type test struct {
	file    string
	content []byte
}

func doTestI18nTranslate(t *testing.T, data []test, lang, id string) string {
	i18nBundle := bundle.New()

	for _, r := range data {
		err := i18nBundle.ParseTranslationFileBytes(r.file, r.content)
		if err != nil {
			t.Errorf("Error parsing translation file: %s", err)
		}
	}

	SetI18nTfuncs(i18nBundle)
	SetTranslateLang(lang)

	translated, err := I18nTranslate(id, nil)
	if err != nil {
		t.Errorf("Error translating '%s': %s", id, err)
	}
	return translated
}

func TestI18nTranslate(t *testing.T) {
	var data []test
	var actual, expected string

	viper.SetDefault("DefaultContentLanguage", "en")

	// Test without and with placeholders
	for _, enablePlaceholders := range []bool{false, true} {
		viper.Set("EnableMissingTranslationPlaceholders", enablePlaceholders)

		// All translations present
		data = []test{
			{"en.yaml", []byte("- id: \"hello\"\n  translation: \"Hello, World!\"")},
			{"es.yaml", []byte("- id: \"hello\"\n  translation: \"¡Hola, Mundo!\"")},
		}
		expected = "¡Hola, Mundo!"
		actual = doTestI18nTranslate(t, data, "es", "hello")
		assert.Equal(t, expected, actual)

		// Translation missing in current language but present in default
		data = []test{
			{"en.yaml", []byte("- id: \"hello\"\n  translation: \"Hello, World!\"")},
			{"es.yaml", []byte("- id: \"goodbye\"\n  translation: \"¡Adiós, Mundo!\"")},
		}
		if enablePlaceholders {
			expected = "[i18n] hello"
		} else {
			expected = "Hello, World!"
		}
		actual = doTestI18nTranslate(t, data, "es", "hello")
		assert.Equal(t, expected, actual)

		// Translation missing in default language but present in current
		data = []test{
			{"en.yaml", []byte("- id: \"goodbye\"\n  translation: \"Goodbye, World!\"")},
			{"es.yaml", []byte("- id: \"hello\"\n  translation: \"¡Hola, Mundo!\"")},
		}
		expected = "¡Hola, Mundo!"
		actual = doTestI18nTranslate(t, data, "es", "hello")
		assert.Equal(t, expected, actual)

		// Translation missing in both default and current language
		data = []test{
			{"en.yaml", []byte("- id: \"goodbye\"\n  translation: \"Goodbye, World!\"")},
			{"es.yaml", []byte("- id: \"goodbye\"\n  translation: \"¡Adiós, Mundo!\"")},
		}
		if enablePlaceholders {
			expected = "[i18n] hello"
		} else {
			expected = ""
		}
		actual = doTestI18nTranslate(t, data, "es", "hello")
		assert.Equal(t, expected, actual)

		// Default translation file missing or empty
		data = []test{
			{"en.yaml", []byte("")},
		}
		actual = doTestI18nTranslate(t, data, "es", "hello")
		if enablePlaceholders {
			expected = "[i18n] hello"
		} else {
			expected = ""
		}
		assert.Equal(t, expected, actual)
	}
}
