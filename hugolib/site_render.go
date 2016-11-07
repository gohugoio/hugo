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
	"path"
	"path/filepath"
	"sync"

	"github.com/spf13/hugo/helpers"

	jww "github.com/spf13/jwalterweatherman"
)

// renderPages renders pages each corresponding to a markdown file.
// TODO(bep np doc
func (s *Site) renderPages() error {

	results := make(chan error)
	pages := make(chan *Page)
	errs := make(chan error)

	go errorCollator(results, errs)

	procs := getGoMaxProcs()

	wg := &sync.WaitGroup{}

	for i := 0; i < procs*4; i++ {
		wg.Add(1)
		go pageRenderer(s, pages, results, wg)
	}

	for _, page := range s.Nodes {
		pages <- page
	}

	close(pages)

	wg.Wait()

	close(results)

	err := <-errs
	if err != nil {
		return fmt.Errorf("Error(s) rendering pages: %s", err)
	}
	return nil
}

func pageRenderer(s *Site, pages <-chan *Page, results chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	for p := range pages {
		targetPath := p.TargetPath()
		layouts := p.layouts()
		jww.DEBUG.Printf("Render %s to %q with layouts %q", p.NodeType, targetPath, layouts)

		if err := s.renderAndWritePage("page "+p.FullFilePath(), targetPath, p, s.appendThemeTemplates(layouts)...); err != nil {
			results <- err
		}

		// Taxonomy terms have no page set to paginate, so skip that for now.
		if p.NodeType.IsNode() && p.NodeType != NodeTaxonomyTerms {
			if err := s.renderPaginator(p); err != nil {
				results <- err
			}
		}

		if err := s.renderRSS(p); err != nil {
			results <- err
		}
	}
}

// renderPaginator must be run after the owning Page has been rendered.
// TODO(bep) np
func (s *Site) renderPaginator(p *Page) error {
	if p.paginator != nil {
		jww.DEBUG.Printf("Render paginator for page %q", p.Path())
		paginatePath := helpers.Config().GetString("paginatePath")

		// write alias for page 1
		// TODO(bep) ml all of these n.addLang ... fix.
		// TODO(bep) np URL

		aliasPath := p.addLangPathPrefix(helpers.PaginateAliasPath(path.Join(p.sections...), 1))
		//TODO(bep) np node.permalink
		s.writeDestAlias(aliasPath, p.Node.Permalink(), nil)

		pagers := p.paginator.Pagers()

		for i, pager := range pagers {
			if i == 0 {
				// already created
				continue
			}

			pagerNode := p.copy()

			pagerNode.paginator = pager
			if pager.TotalPages() > 0 {
				first, _ := pager.page(0)
				pagerNode.Date = first.Date
				pagerNode.Lastmod = first.Lastmod
			}

			pageNumber := i + 1
			htmlBase := path.Join(p.URLPath.URL, fmt.Sprintf("/%s/%d", paginatePath, pageNumber))
			htmlBase = p.addLangPathPrefix(htmlBase)

			if err := s.renderAndWritePage(pagerNode.Title,
				filepath.FromSlash(htmlBase), pagerNode, p.layouts()...); err != nil {
				return err
			}

		}
	}
	return nil
}

func (s *Site) renderRSS(p *Page) error {
	layouts := p.rssLayouts()

	if layouts == nil {
		// No RSS for this NodeType
		return nil
	}

	// TODO(bep) np check RSS titles
	// TODO(bep) np check RSS page limit, 50?
	rssNode := p.copy()

	// TODO(bep) np todelido URL
	rssURI := s.Language.GetString("rssURI")
	rssNode.URLPath.URL = path.Join(rssNode.URLPath.URL, rssURI)

	if err := s.renderAndWriteXML(rssNode.Title, rssNode.addLangFilepathPrefix(rssNode.URLPath.URL), rssNode, s.appendThemeTemplates(layouts)...); err != nil {
		return err
	}

	return nil
}
