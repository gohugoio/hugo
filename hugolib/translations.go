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
	"github.com/gohugoio/hugo/resources/page"
)

// Translations represent the other translations for a given page. The
// string here is the language code, as affected by the `post.LANG.md`
// filename.
type Translations map[string]page.Page

func pagesToTranslationsMap(pages Pages) map[string]Translations {
	out := make(map[string]Translations)

	for _, page := range pages {
		pagep := page.(*Page)
		base := pagep.TranslationKey()

		pageTranslation, present := out[base]
		if !present {
			pageTranslation = make(Translations)
		}

		pageLang := pagep.Lang()
		if pageLang == "" {
			continue
		}

		pageTranslation[pageLang] = page
		out[base] = pageTranslation
	}

	return out
}

func assignTranslationsToPages(allTranslations map[string]Translations, pages Pages) {
	for _, page := range pages {
		pagep := page.(*Page)
		pagep.translations = pagep.translations[:0]
		base := pagep.TranslationKey()
		trans, exist := allTranslations[base]
		if !exist {
			continue
		}

		for _, translatedPage := range trans {
			pagep.translations = append(pagep.translations, translatedPage)
		}

		pageBy(languagePageSort).Sort(pagep.translations)
	}
}
