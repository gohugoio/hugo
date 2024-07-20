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

	"github.com/gohugoio/hugo/resources/kinds"
	"github.com/gohugoio/hugo/resources/page"
)

func newPagePaginator(source *pageState) *pagePaginator {
	return &pagePaginator{
		source:            source,
		pagePaginatorInit: &pagePaginatorInit{},
	}
}

type pagePaginator struct {
	*pagePaginatorInit
	source *pageState
}

type pagePaginatorInit struct {
	init    sync.Once
	current *page.Pager
}

// reset resets the paginator to allow for a rebuild.
func (p *pagePaginator) reset() {
	p.pagePaginatorInit = &pagePaginatorInit{}
}

func (p *pagePaginator) Paginate(seq any, options ...any) (*page.Pager, error) {
	var initErr error
	p.init.Do(func() {
		pagerSize, err := page.ResolvePagerSize(p.source.s.Conf, options...)
		if err != nil {
			initErr = err
			return
		}

		pd := p.source.targetPathDescriptor
		pd.Type = p.source.outputFormat()
		paginator, err := page.Paginate(pd, seq, pagerSize)
		if err != nil {
			initErr = err
			return
		}

		p.current = paginator.Pagers()[0]
	})

	if initErr != nil {
		return nil, initErr
	}

	return p.current, nil
}

func (p *pagePaginator) Paginator(options ...any) (*page.Pager, error) {
	var initErr error
	p.init.Do(func() {
		pagerSize, err := page.ResolvePagerSize(p.source.s.Conf, options...)
		if err != nil {
			initErr = err
			return
		}

		pd := p.source.targetPathDescriptor
		pd.Type = p.source.outputFormat()

		var pages page.Pages

		switch p.source.Kind() {
		case kinds.KindHome:
			// From Hugo 0.57 we made home.Pages() work like any other
			// section. To avoid the default paginators for the home page
			// changing in the wild, we make this a special case.
			pages = p.source.s.RegularPages()
		case kinds.KindTerm, kinds.KindTaxonomy:
			pages = p.source.Pages()
		default:
			pages = p.source.RegularPages()
		}

		paginator, err := page.Paginate(pd, pages, pagerSize)
		if err != nil {
			initErr = err
			return
		}

		p.current = paginator.Pagers()[0]
	})

	if initErr != nil {
		return nil, initErr
	}

	return p.current, nil
}
