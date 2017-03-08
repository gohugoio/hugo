// Copyright 2016 The Hugo Authors. All rights reserved.
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

import (
	"fmt"
)

// Translations represent the other translations for a given page. The
// string here is the language code, as affected by the `post.LANG.md`
// filename.
type Translations map[string]*Page

func pagesToTranslationsMap(pages []*Page) map[string]Translations {
	out := make(map[string]Translations)

	for _, page := range pages {
		base := createTranslationKey(page)

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

func createTranslationKey(p *Page) string {
	base := p.TranslationBaseName()

	if p.IsNode() {
		// TODO(bep) see https://github.com/spf13/hugo/issues/2699
		// Must prepend the section and kind to the key to make it unique
		base = fmt.Sprintf("%s/%s/%s", p.Kind, p.sections, base)
	}

	return base
}

func assignTranslationsToPages(allTranslations map[string]Translations, pages []*Page) {
	for _, page := range pages {
		page.translations = page.translations[:0]
		base := createTranslationKey(page)
		trans, exist := allTranslations[base]
		if !exist {
			continue
		}

		for _, translatedPage := range trans {
			page.translations = append(page.translations, translatedPage)
		}

		pageBy(languagePageSort).Sort(page.translations)
	}
}
