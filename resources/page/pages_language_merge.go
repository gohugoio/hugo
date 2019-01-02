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

package page

import (
	"fmt"
)

var (
	_ pagesLanguageMerger = (*Pages)(nil)
)

type pagesLanguageMerger interface {
	MergeByLanguage(other Pages) Pages
	// Needed for integration with the tpl package.
	MergeByLanguageInterface(other interface{}) (interface{}, error)
}

// MergeByLanguage supplies missing translations in p1 with values from p2.
// The result is sorted by the default sort order for pages.
func (p1 Pages) MergeByLanguage(p2 Pages) Pages {
	merge := func(pages *Pages) {
		m := make(map[string]bool)
		for _, p := range *pages {
			m[p.TranslationKey()] = true
		}

		for _, p := range p2 {
			if _, found := m[p.TranslationKey()]; !found {
				*pages = append(*pages, p)
			}
		}

		SortByDefault(*pages)
	}

	out, _ := spc.getP("pages.MergeByLanguage", merge, p1, p2)

	return out
}

// MergeByLanguageInterface is the generic version of MergeByLanguage. It
// is here just so it can be called from the tpl package.
func (p1 Pages) MergeByLanguageInterface(in interface{}) (interface{}, error) {
	if in == nil {
		return p1, nil
	}
	p2, ok := in.(Pages)
	if !ok {
		return nil, fmt.Errorf("%T cannot be merged by language", in)
	}
	return p1.MergeByLanguage(p2), nil
}
