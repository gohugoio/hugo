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
	"path/filepath"
	"testing"

	"github.com/bep/logg"
	"github.com/gohugoio/hugo/common/types"
	"github.com/gohugoio/hugo/config/testconfig"

	"github.com/gohugoio/hugo/resources/page"
	"github.com/spf13/afero"

	"github.com/gohugoio/hugo/deps"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/config"
)

type i18nTest struct {
	name                             string
	data                             map[string][]byte
	args                             any
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
	// https://github.com/gohugoio/hugo/issues/7787
	{
		name: "readingTime-one",
		data: map[string][]byte{
			"en.toml": []byte(`[readingTime]
one = "One minute to read"
other = "{{ .Count }} minutes to read"
`),
		},
		args:         1,
		lang:         "en",
		id:           "readingTime",
		expected:     "One minute to read",
		expectedFlag: "One minute to read",
	},
	{
		name: "readingTime-many-dot",
		data: map[string][]byte{
			"en.toml": []byte(`[readingTime]
one = "One minute to read"
other = "{{ . }} minutes to read"
`),
		},
		args:         21,
		lang:         "en",
		id:           "readingTime",
		expected:     "21 minutes to read",
		expectedFlag: "21 minutes to read",
	},
	{
		name: "readingTime-many",
		data: map[string][]byte{
			"en.toml": []byte(`[readingTime]
one = "One minute to read"
other = "{{ .Count }} minutes to read"
`),
		},
		args:         21,
		lang:         "en",
		id:           "readingTime",
		expected:     "21 minutes to read",
		expectedFlag: "21 minutes to read",
	},
	// Issue #8454
	{
		name: "readingTime-map-one",
		data: map[string][]byte{
			"en.toml": []byte(`[readingTime]
one = "One minute to read"
other = "{{ .Count }} minutes to read"
`),
		},
		args:         map[string]any{"Count": 1},
		lang:         "en",
		id:           "readingTime",
		expected:     "One minute to read",
		expectedFlag: "One minute to read",
	},
	{
		name: "readingTime-string-one",
		data: map[string][]byte{
			"en.toml": []byte(`[readingTime]
one = "One minute to read"
other = "{{ . }} minutes to read"
`),
		},
		args:         "1",
		lang:         "en",
		id:           "readingTime",
		expected:     "One minute to read",
		expectedFlag: "One minute to read",
	},
	{
		name: "readingTime-map-many",
		data: map[string][]byte{
			"en.toml": []byte(`[readingTime]
one = "One minute to read"
other = "{{ .Count }} minutes to read"
`),
		},
		args:         map[string]any{"Count": 21},
		lang:         "en",
		id:           "readingTime",
		expected:     "21 minutes to read",
		expectedFlag: "21 minutes to read",
	},
	{
		name: "argument-float",
		data: map[string][]byte{
			"en.toml": []byte(`[float]
other = "Number is {{ . }}"
`),
		},
		args:         22.5,
		lang:         "en",
		id:           "float",
		expected:     "Number is 22.5",
		expectedFlag: "Number is 22.5",
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
	// Issue #7838
	{
		name: "unknown-language-codes",
		data: map[string][]byte{
			"en.toml": []byte(`[readingTime]
one ="en one"
other = "en count {{.Count}}"`),
			"a1.toml": []byte(`[readingTime]
one =  "a1 one"
other = "a1 count {{ .Count }}"`),
			"a2.toml": []byte(`[readingTime]
one =  "a2 one"
other = "a2 count {{ .Count }}"`),
		},
		args:         3,
		lang:         "a2",
		id:           "readingTime",
		expected:     "a2 count 3",
		expectedFlag: "a2 count 3",
	},
	// https://github.com/gohugoio/hugo/issues/7798
	{
		name: "known-language-missing-plural",
		data: map[string][]byte{
			"oc.toml": []byte(`[oc]
one =  "abc"`),
		},
		args:         1,
		lang:         "oc",
		id:           "oc",
		expected:     "abc",
		expectedFlag: "abc",
	},
	// https://github.com/gohugoio/hugo/issues/7794
	{
		name: "dotted-bare-key",
		data: map[string][]byte{
			"en.toml": []byte(`"shop_nextPage.one" = "Show Me The Money"
`),
		},
		args:         nil,
		lang:         "en",
		id:           "shop_nextPage.one",
		expected:     "Show Me The Money",
		expectedFlag: "Show Me The Money",
	},
	// https: //github.com/gohugoio/hugo/issues/7804
	{
		name: "lang-with-hyphen",
		data: map[string][]byte{
			"pt-br.toml": []byte(`foo.one =  "abc"`),
		},
		args:         1,
		lang:         "pt-br",
		id:           "foo",
		expected:     "abc",
		expectedFlag: "abc",
	},
}

func TestPlural(t *testing.T) {
	c := qt.New(t)

	for _, test := range []struct {
		name     string
		lang     string
		id       string
		templ    string
		variants []types.KeyValue
	}{
		{
			name: "English",
			lang: "en",
			id:   "hour",
			templ: `
[hour]
one = "{{ . }} hour"
other = "{{ . }} hours"`,
			variants: []types.KeyValue{
				{Key: 1, Value: "1 hour"},
				{Key: "1", Value: "1 hour"},
				{Key: 1.5, Value: "1.5 hours"},
				{Key: "1.5", Value: "1.5 hours"},
				{Key: 2, Value: "2 hours"},
				{Key: "2", Value: "2 hours"},
			},
		},
		{
			name: "Other only",
			lang: "en",
			id:   "hour",
			templ: `
[hour]
other = "{{ with . }}{{ . }}{{ end }} hours"`,
			variants: []types.KeyValue{
				{Key: 1, Value: "1 hours"},
				{Key: "1", Value: "1 hours"},
				{Key: 2, Value: "2 hours"},
				{Key: nil, Value: " hours"},
			},
		},
		{
			name: "Polish",
			lang: "pl",
			id:   "day",
			templ: `
[day]
one = "{{ . }} miesiąc"
few = "{{ . }} miesiące"
many = "{{ . }} miesięcy"
other = "{{ . }} miesiąca"
`,
			variants: []types.KeyValue{
				{Key: 1, Value: "1 miesiąc"},
				{Key: 2, Value: "2 miesiące"},
				{Key: 100, Value: "100 miesięcy"},
				{Key: "100.0", Value: "100.0 miesiąca"},
				{Key: 100.0, Value: "100 miesiąca"},
			},
		},
	} {
		c.Run(test.name, func(c *qt.C) {
			cfg := config.New()
			cfg.Set("enableMissingTranslationPlaceholders", true)
			cfg.Set("publishDir", "public")
			afs := afero.NewMemMapFs()

			err := afero.WriteFile(afs, filepath.Join("i18n", test.lang+".toml"), []byte(test.templ), 0o755)
			c.Assert(err, qt.IsNil)

			d, tp := prepareDeps(afs, cfg)

			f := tp.t.Func(test.lang)
			ctx := context.Background()

			for _, variant := range test.variants {
				c.Assert(f(ctx, test.id, variant.Key), qt.Equals, variant.Value, qt.Commentf("input: %v", variant.Key))
				c.Assert(d.Log.LoggCount(logg.LevelWarn), qt.Equals, 0)
			}
		})
	}
}

func doTestI18nTranslate(t testing.TB, test i18nTest, cfg config.Provider) string {
	tp := prepareTranslationProvider(t, test, cfg)
	f := tp.t.Func(test.lang)
	return f(context.Background(), test.id, test.args)
}

type countField struct {
	Count any
}

type noCountField struct {
	Counts int
}

type countMethod struct{}

func (c countMethod) Count() any {
	return 32.5
}

func TestGetPluralCount(t *testing.T) {
	c := qt.New(t)

	c.Assert(getPluralCount(map[string]any{"Count": 32}), qt.Equals, 32)
	c.Assert(getPluralCount(map[string]any{"Count": 1}), qt.Equals, 1)
	c.Assert(getPluralCount(map[string]any{"Count": 1.5}), qt.Equals, "1.5")
	c.Assert(getPluralCount(map[string]any{"Count": "32"}), qt.Equals, "32")
	c.Assert(getPluralCount(map[string]any{"Count": "32.5"}), qt.Equals, "32.5")
	c.Assert(getPluralCount(map[string]any{"count": 32}), qt.Equals, 32)
	c.Assert(getPluralCount(map[string]any{"Count": "32"}), qt.Equals, "32")
	c.Assert(getPluralCount(map[string]any{"Counts": 32}), qt.Equals, nil)
	c.Assert(getPluralCount("foo"), qt.Equals, nil)
	c.Assert(getPluralCount(countField{Count: 22}), qt.Equals, 22)
	c.Assert(getPluralCount(countField{Count: 1.5}), qt.Equals, "1.5")
	c.Assert(getPluralCount(&countField{Count: 22}), qt.Equals, 22)
	c.Assert(getPluralCount(noCountField{Counts: 23}), qt.Equals, nil)
	c.Assert(getPluralCount(countMethod{}), qt.Equals, "32.5")
	c.Assert(getPluralCount(&countMethod{}), qt.Equals, "32.5")

	c.Assert(getPluralCount(1234), qt.Equals, 1234)
	c.Assert(getPluralCount(1234.4), qt.Equals, "1234.4")
	c.Assert(getPluralCount(1234.0), qt.Equals, "1234.0")
	c.Assert(getPluralCount("1234"), qt.Equals, "1234")
	c.Assert(getPluralCount("0.5"), qt.Equals, "0.5")
	c.Assert(getPluralCount(nil), qt.Equals, nil)
}

func prepareTranslationProvider(t testing.TB, test i18nTest, cfg config.Provider) *TranslationProvider {
	c := qt.New(t)
	afs := afero.NewMemMapFs()

	for file, content := range test.data {
		err := afero.WriteFile(afs, filepath.Join("i18n", file), []byte(content), 0o755)
		c.Assert(err, qt.IsNil)
	}

	_, tp := prepareDeps(afs, cfg)
	return tp
}

func prepareDeps(afs afero.Fs, cfg config.Provider) (*deps.Deps, *TranslationProvider) {
	d := testconfig.GetTestDeps(afs, cfg)
	translationProvider := NewTranslationProvider()
	d.TranslationProvider = translationProvider
	d.Site = page.NewDummyHugoSite(d.Conf)
	if err := d.Compile(nil); err != nil {
		panic(err)
	}
	return d, translationProvider
}

func TestI18nTranslate(t *testing.T) {
	c := qt.New(t)
	var actual, expected string
	v := config.New()

	// Test without and with placeholders
	for _, enablePlaceholders := range []bool{false, true} {
		v.Set("enableMissingTranslationPlaceholders", enablePlaceholders)

		for _, test := range i18nTests {
			c.Run(fmt.Sprintf("%s-%t", test.name, enablePlaceholders), func(c *qt.C) {
				if enablePlaceholders {
					expected = test.expectedFlag
				} else {
					expected = test.expected
				}
				actual = doTestI18nTranslate(c, test, v)
				c.Assert(actual, qt.Equals, expected)
			})
		}
	}
}

// TestI18nLanguageCode: when languageCode is set, use it for translation file lookup.
func TestI18nLanguageCode(t *testing.T) {
	c := qt.New(t)

	afs := afero.NewMemMapFs()

	// Create only pt-br.toml (no pt.toml)
	err := afero.WriteFile(afs, filepath.Join("i18n", "pt-br.toml"), []byte("hello = \"Olá\""), 0o755)
	c.Assert(err, qt.IsNil)

	// Set up config with pt language and languageCode pt-BR
	cfg := config.New()
	cfg.Set("defaultContentLanguage", "pt")
	cfg.Set("languages", map[string]any{
		"pt": map[string]any{
			"languageCode": "pt-BR",
		},
	})

	d, _ := prepareDeps(afs, cfg)

	ctx := context.Background()
	actual := d.Translate(ctx, "hello", nil)

	c.Assert(actual, qt.Equals, "Olá", qt.Commentf("Expected translation from pt-br.toml when languageCode is pt-BR"))
}

// TestI18nLanguageCodePreference: when both Lang.toml and languageCode.toml exist, prefer languageCode.
func TestI18nLanguageCodePreference(t *testing.T) {
	c := qt.New(t)

	afs := afero.NewMemMapFs()

	// Create both pt.toml and pt-br.toml with different values
	err := afero.WriteFile(afs, filepath.Join("i18n", "pt.toml"), []byte("hello = \"Olá (PT)\""), 0o755)
	c.Assert(err, qt.IsNil)
	err = afero.WriteFile(afs, filepath.Join("i18n", "pt-br.toml"), []byte("hello = \"Olá (PT-BR)\""), 0o755)
	c.Assert(err, qt.IsNil)

	// Set up config with pt language and languageCode pt-BR
	cfg := config.New()
	cfg.Set("defaultContentLanguage", "pt")
	cfg.Set("languages", map[string]any{
		"pt": map[string]any{
			"languageCode": "pt-BR",
		},
	})

	d, _ := prepareDeps(afs, cfg)

	// Should use pt-br.toml (languageCode) not pt.toml (Lang)
	ctx := context.Background()
	actual := d.Translate(ctx, "hello", nil)

	c.Assert(actual, qt.Equals, "Olá (PT-BR)", qt.Commentf("Expected translation from pt-br.toml when both pt.toml and pt-br.toml exist"))
}

// TestI18nLanguageCodeEqualsLang: when languageCode equals Lang (case-insensitive), use Lang file.
func TestI18nLanguageCodeEqualsLang(t *testing.T) {
	c := qt.New(t)

	afs := afero.NewMemMapFs()

	err := afero.WriteFile(afs, filepath.Join("i18n", "pt.toml"), []byte("hello = \"Olá (PT)\""), 0o755)
	c.Assert(err, qt.IsNil)
	err = afero.WriteFile(afs, filepath.Join("i18n", "pt-br.toml"), []byte("hello = \"Olá\""), 0o755)
	c.Assert(err, qt.IsNil)

	// Set up config with pt-br language and languageCode pt-BR (same as Lang)
	cfg := config.New()
	cfg.Set("defaultContentLanguage", "pt-br")
	cfg.Set("languages", map[string]any{
		"pt-br": map[string]any{
			"languageCode": "pt-BR", // Same as Lang
		},
	})

	d, _ := prepareDeps(afs, cfg)

	// Should use pt-br.toml (Lang) since languageCode equals Lang
	ctx := context.Background()
	actual := d.Translate(ctx, "hello", nil)

	c.Assert(actual, qt.Equals, "Olá", qt.Commentf("Expected translation from pt-br.toml when languageCode equals Lang"))
}

// TestI18nLanguageCodeFallback: when languageCode file doesn't exist, should fall back to default language.
func TestI18nLanguageCodeFallback(t *testing.T) {
	c := qt.New(t)

	afs := afero.NewMemMapFs()

	// Create only en.toml (default language), no pt-br.toml
	err := afero.WriteFile(afs, filepath.Join("i18n", "en.toml"), []byte("hello = \"Hello\""), 0o755)
	c.Assert(err, qt.IsNil)

	// Set up config with pt language and languageCode pt-BR, but default is en
	// Note: defaultContentLanguage must be in the languages map
	cfg := config.New()
	cfg.Set("defaultContentLanguage", "en")
	cfg.Set("languages", map[string]any{
		"en": map[string]any{}, // Default language must be in languages map
		"pt": map[string]any{
			"languageCode": "pt-BR",
		},
	})

	d, _ := prepareDeps(afs, cfg)

	// Should fall back to default language (en) since pt-br.toml doesn't exist
	ctx := context.Background()
	actual := d.Translate(ctx, "hello", nil)

	c.Assert(actual, qt.Equals, "Hello", qt.Commentf("Expected fallback to default language (en) when pt-br.toml doesn't exist"))
}

// TestI18nLanguageCodeMultipleLanguages: test multiple languages with different languageCode settings.
func TestI18nLanguageCodeMultipleLanguages(t *testing.T) {
	c := qt.New(t)

	afs := afero.NewMemMapFs()

	// Create translation files for different language codes
	err := afero.WriteFile(afs, filepath.Join("i18n", "en.toml"), []byte("hello = \"Hello\""), 0o755)
	c.Assert(err, qt.IsNil)
	err = afero.WriteFile(afs, filepath.Join("i18n", "en-us.toml"), []byte("hello = \"Hello (US)\""), 0o755)
	c.Assert(err, qt.IsNil)
	err = afero.WriteFile(afs, filepath.Join("i18n", "zh-cn.toml"), []byte("hello = \"你好\""), 0o755)
	c.Assert(err, qt.IsNil)

	// Set up config with multiple languages
	cfg := config.New()
	cfg.Set("defaultContentLanguage", "en")
	cfg.Set("languages", map[string]any{
		"en": map[string]any{
			"languageCode": "en-US",
		},
		"zh": map[string]any{
			"languageCode": "zh-CN",
		},
	})

	// Test English with languageCode en-US
	dEn, _ := prepareDeps(afs, cfg)
	ctx := context.Background()
	actualEn := dEn.Translate(ctx, "hello", nil)
	c.Assert(actualEn, qt.Equals, "Hello (US)", qt.Commentf("Expected translation from en-us.toml"))

	// Test Chinese with languageCode zh-CN
	// We need to create a new Deps for zh language
	cfgZh := config.New()
	cfgZh.Set("defaultContentLanguage", "zh")
	cfgZh.Set("languages", map[string]any{
		"zh": map[string]any{
			"languageCode": "zh-CN",
		},
	})
	dZh, _ := prepareDeps(afs, cfgZh)
	actualZh := dZh.Translate(ctx, "hello", nil)
	c.Assert(actualZh, qt.Equals, "你好", qt.Commentf("Expected translation from zh-cn.toml"))
}

func BenchmarkI18nTranslate(b *testing.B) {
	v := config.New()
	for _, test := range i18nTests {
		b.Run(test.name, func(b *testing.B) {
			tp := prepareTranslationProvider(b, test, v)
			b.ResetTimer()
			for b.Loop() {
				f := tp.t.Func(test.lang)
				actual := f(context.Background(), test.id, test.args)
				if actual != test.expected {
					b.Fatalf("expected %v got %v", test.expected, actual)
				}
			}
		})
	}
}
