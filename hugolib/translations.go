// Copyright 2019 The Hugo Authors. All rights reserved.
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

func pagesToTranslationsMap(sites []*Site) map[string]page.Pages {
	out := make(map[string]page.Pages)

	for _, s := range sites {
		s.pageMap.pageTrees.Walk(func(ss string, n *contentNode) bool {
			p := n.p
			// TranslationKey is implemented for all page types.
			base := p.TranslationKey()

			pageTranslations, found := out[base]
			if !found {
				pageTranslations = make(page.Pages, 0)
			}

			pageTranslations = append(pageTranslations, p)
			out[base] = pageTranslations

			return false
		})
	}

	return out
}

func assignTranslationsToPages(allTranslations map[string]page.Pages, sites []*Site) {
	for _, s := range sites {
		s.pageMap.pageTrees.Walk(func(ss string, n *contentNode) bool {
			p := n.p
			base := p.TranslationKey()
			translations, found := allTranslations[base]
			if !found {
				return false
			}
			p.setTranslations(translations)
			return false
		})
	}
}
