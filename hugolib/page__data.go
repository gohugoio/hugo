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
	"sync"

	"github.com/gohugoio/hugo/resources/page"
)

type pageData struct {
	*pageState

	dataInit sync.Once
	data     page.Data
}

func (p *pageData) Data() interface{} {
	p.dataInit.Do(func() {
		p.data = make(page.Data)

		if p.Kind() == page.KindPage {
			return
		}

		switch p.Kind() {
		case page.KindTaxonomy:
			termInfo := p.getTaxonomyNodeInfo()
			pluralInfo := termInfo.parent

			singular := pluralInfo.singular
			plural := pluralInfo.plural
			term := termInfo.term
			taxonomy := p.s.Taxonomies[plural].Get(termInfo.termKey)

			p.data[singular] = taxonomy
			p.data["Singular"] = singular
			p.data["Plural"] = plural
			p.data["Term"] = term
		case page.KindTaxonomyTerm:
			info := p.getTaxonomyNodeInfo()
			plural := info.plural
			singular := info.singular

			p.data["Singular"] = singular
			p.data["Plural"] = plural
			p.data["Terms"] = p.s.Taxonomies[plural]
			// keep the following just for legacy reasons
			p.data["OrderedIndex"] = p.data["Terms"]
			p.data["Index"] = p.data["Terms"]
		}

		// Assign the function to the map to make sure it is lazily initialized
		p.data["pages"] = p.Pages

	})

	return p.data
}
