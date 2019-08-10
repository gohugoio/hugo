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
	"path/filepath"
	"testing"

	"github.com/gohugoio/hugo/modules"

	"github.com/gohugoio/hugo/tpl/tplimpl"

	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/langs"
	"github.com/gohugoio/hugo/resources/page"
	"github.com/spf13/afero"
	"github.com/spf13/viper"

	"github.com/gohugoio/hugo/deps"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/hugofs"
)

var logger = loggers.NewErrorLogger()

type i18nTest struct {
	name                             string
	data                             map[string][]byte
	args                             interface{}
	lang, id, expected, expectedFlag string
}

var i18nTests = []i18nTest{
	// All translations present
	{
		name: "all-present",
		data: map[string][]byte{
			"en.toml": []byte("[hello]\nother = \"Hello, World!\""),
			"es.toml": []byte("[hello]\nother = \"¡Hola, Mundo!\""),
		},
		args:         nil,
		lang:         "es",
		id:           "hello",
		expected:     "¡Hola, Mundo!",
		expectedFlag: "¡Hola, Mundo!",
	},
	// Translation missing in current language but present in default
	{
		name: "present-in-default",
		data: map[string][]byte{
			"en.toml": []byte("[hello]\nother = \"Hello, World!\""),
			"es.toml": []byte("[goodbye]\nother = \"¡Adiós, Mundo!\""),
		},
		args:         nil,
		lang:         "es",
		id:           "hello",
		expected:     "Hello, World!",
		expectedFlag: "[i18n] hello",
	},
	// Translation missing in default language but present in current
	{
		name: "present-in-current",
		data: map[string][]byte{
			"en.toml": []byte("[goodbye]\nother = \"Goodbye, World!\""),
			"es.toml": []byte("[hello]\nother = \"¡Hola, Mundo!\""),
		},
		args:         nil,
		lang:         "es",
		id:           "hello",
		expected:     "¡Hola, Mundo!",
		expectedFlag: "¡Hola, Mundo!",
	},
	// Translation missing in both default and current language
	{
		name: "missing",
		data: map[string][]byte{
			"en.toml": []byte("[goodbye]\nother = \"Goodbye, World!\""),
			"es.toml": []byte("[goodbye]\nother = \"¡Adiós, Mundo!\""),
		},
		args:         nil,
		lang:         "es",
		id:           "hello",
		expected:     "",
		expectedFlag: "[i18n] hello",
	},
	// Default translation file missing or empty
	{
		name: "file-missing",
		data: map[string][]byte{
			"en.toml": []byte(""),
		},
		args:         nil,
		lang:         "es",
		id:           "hello",
		expected:     "",
		expectedFlag: "[i18n] hello",
	},
	// Context provided
	{
		name: "context-provided",
		data: map[string][]byte{
			"en.toml": []byte("[wordCount]\nother = \"Hello, {{.WordCount}} people!\""),
			"es.toml": []byte("[wordCount]\nother = \"¡Hola, {{.WordCount}} gente!\""),
		},
		args: struct {
			WordCount int
		}{
			50,
		},
		lang:         "es",
		id:           "wordCount",
		expected:     "¡Hola, 50 gente!",
		expectedFlag: "¡Hola, 50 gente!",
	},
	// Same id and translation in current language
	// https://github.com/gohugoio/hugo/issues/2607
	{
		name: "same-id-and-translation",
		data: map[string][]byte{
			"es.toml": []byte("[hello]\nother = \"hello\""),
			"en.toml": []byte("[hello]\nother = \"hi\""),
		},
		args:         nil,
		lang:         "es",
		id:           "hello",
		expected:     "hello",
		expectedFlag: "hello",
	},
	// Translation missing in current language, but same id and translation in default
	{
		name: "same-id-and-translation-default",
		data: map[string][]byte{
			"es.toml": []byte("[bye]\nother = \"bye\""),
			"en.toml": []byte("[hello]\nother = \"hello\""),
		},
		args:         nil,
		lang:         "es",
		id:           "hello",
		expected:     "hello",
		expectedFlag: "[i18n] hello",
	},
	// Unknown language code should get its plural spec from en
	{
		name: "unknown-language-code",
		data: map[string][]byte{
			"en.toml": []byte(`[readingTime]
one ="one minute read"
other = "{{.Count}} minutes read"`),
			"klingon.toml": []byte(`[readingTime]
one =  "eitt minutt med lesing"
other = "{{ .Count }} minuttar lesing"`),
		},
		args:         3,
		lang:         "klingon",
		id:           "readingTime",
		expected:     "3 minuttar lesing",
		expectedFlag: "3 minuttar lesing",
	},
}

func doTestI18nTranslate(t testing.TB, test i18nTest, cfg config.Provider) string {
	tp := prepareTranslationProvider(t, test, cfg)
	f := tp.t.Func(test.lang)
	return f(test.id, test.args)

}

func prepareTranslationProvider(t testing.TB, test i18nTest, cfg config.Provider) *TranslationProvider {
	c := qt.New(t)
	fs := hugofs.NewMem(cfg)

	for file, content := range test.data {
		err := afero.WriteFile(fs.Source, filepath.Join("i18n", file), []byte(content), 0755)
		c.Assert(err, qt.IsNil)
	}

	tp := NewTranslationProvider()
	depsCfg := newDepsConfig(tp, cfg, fs)
	d, err := deps.New(depsCfg)
	c.Assert(err, qt.IsNil)
	c.Assert(d.LoadResources(), qt.IsNil)

	return tp
}

func newDepsConfig(tp *TranslationProvider, cfg config.Provider, fs *hugofs.Fs) deps.DepsCfg {
	l := langs.NewLanguage("en", cfg)
	l.Set("i18nDir", "i18n")
	return deps.DepsCfg{
		Language:            l,
		Site:                page.NewDummyHugoSite(cfg),
		Cfg:                 cfg,
		Fs:                  fs,
		Logger:              logger,
		TemplateProvider:    tplimpl.DefaultTemplateProvider,
		TranslationProvider: tp,
	}
}

func getConfig() *viper.Viper {
	v := viper.New()
	v.SetDefault("defaultContentLanguage", "en")
	v.Set("contentDir", "content")
	v.Set("dataDir", "data")
	v.Set("i18nDir", "i18n")
	v.Set("layoutDir", "layouts")
	v.Set("archetypeDir", "archetypes")
	v.Set("assetDir", "assets")
	v.Set("resourceDir", "resources")
	v.Set("publishDir", "public")
	langs.LoadLanguageSettings(v, nil)
	mod, err := modules.CreateProjectModule(v)
	if err != nil {
		panic(err)
	}
	v.Set("allModules", modules.Modules{mod})

	return v

}

func TestI18nTranslate(t *testing.T) {
	c := qt.New(t)
	var actual, expected string
	v := getConfig()

	// Test without and with placeholders
	for _, enablePlaceholders := range []bool{false, true} {
		v.Set("enableMissingTranslationPlaceholders", enablePlaceholders)

		for _, test := range i18nTests {
			if enablePlaceholders {
				expected = test.expectedFlag
			} else {
				expected = test.expected
			}
			actual = doTestI18nTranslate(t, test, v)
			c.Assert(actual, qt.Equals, expected)
		}
	}
}

func BenchmarkI18nTranslate(b *testing.B) {
	v := getConfig()
	for _, test := range i18nTests {
		b.Run(test.name, func(b *testing.B) {
			tp := prepareTranslationProvider(b, test, v)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				f := tp.t.Func(test.lang)
				actual := f(test.id, test.args)
				if actual != test.expected {
					b.Fatalf("expected %v got %v", test.expected, actual)
				}
			}
		})
	}

}
