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
	"strings"
	"sync"

	"github.com/gohugoio/hugo/resources/kinds"
	"github.com/gohugoio/hugo/resources/page"
)

type pageData struct {
	*pageState

	dataInit sync.Once
	data     page.Data
}

func (p *pageData) Data() any {
	p.dataInit.Do(func() {
		p.data = make(page.Data)

		if p.Kind() == kinds.KindPage {
			return
		}

		switch p.Kind() {
		case kinds.KindTerm:
			path := p.Path()
			name := p.s.pageMap.cfg.getTaxonomyConfig(path)
			term := p.s.Taxonomies()[name.plural].Get(strings.TrimPrefix(path, name.pluralTreeKey))
			p.data[name.singular] = term
			p.data["Singular"] = name.singular
			p.data["Plural"] = name.plural
			p.data["Term"] = p.m.term
		case kinds.KindTaxonomy:
			viewCfg := p.s.pageMap.cfg.getTaxonomyConfig(p.Path())
			p.data["Singular"] = viewCfg.singular
			p.data["Plural"] = viewCfg.plural
			p.data["Terms"] = p.s.Taxonomies()[viewCfg.plural]
			// keep the following just for legacy reasons
			p.data["OrderedIndex"] = p.data["Terms"]
			p.data["Index"] = p.data["Terms"]
		}

		// Assign the function to the map to make sure it is lazily initialized
		p.data["pages"] = p.Pages
	})

	return p.data
}
