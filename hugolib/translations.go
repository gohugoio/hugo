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

package hugolib

// Translations represent the other translations for a given page. The
// string here is the language code, as affected by the `post.LANG.md`
// filename.
type Translations map[string]*Page

func pagesToTranslationsMap(pages []*Page) map[string]Translations {
	out := make(map[string]Translations)

	for _, page := range pages {
		base := page.TranslationBaseName()

		pageTranslation, present := out[base]
		if !present {
			pageTranslation = make(Translations)
		}

		pageLang := page.Lang()
		if pageLang == "" {
			continue
		}

		pageTranslation[pageLang] = page
		out[base] = pageTranslation
	}

	return out
}

func assignTranslationsToPages(allTranslations map[string]Translations, pages []*Page) {
	for _, page := range pages {
		base := page.TranslationBaseName()
		trans, exist := allTranslations[base]
		if !exist {
			continue
		}

		for lang, translatedPage := range trans {
			page.Translations[lang] = translatedPage
		}
	}
}
