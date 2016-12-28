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
	"time"

	bp "github.com/spf13/hugo/bufferpool"
	"github.com/spf13/hugo/helpers"
	"github.com/spf13/viper"

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

	for _, page := range s.Pages {
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
		jww.DEBUG.Printf("Render %s to %q with layouts %q", p.Kind, targetPath, layouts)

		if err := s.renderAndWritePage("page "+p.FullFilePath(), targetPath, p, s.appendThemeTemplates(layouts)...); err != nil {
			results <- err
		}

		// Taxonomy terms have no page set to paginate, so skip that for now.
		if p.IsNode() && p.Kind != KindTaxonomyTerm {
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
func (s *Site) renderPaginator(p *Page) error {
	if p.paginator != nil {
		jww.DEBUG.Printf("Render paginator for page %q", p.Path())
		paginatePath := helpers.Config().GetString("paginatePath")

		// write alias for page 1
		// TODO(bep) ml all of these n.addLang ... fix.

		aliasPath := p.addLangPathPrefix(helpers.PaginateAliasPath(path.Join(p.sections...), 1))
		link := p.Permalink()
		s.writeDestAlias(aliasPath, link, nil)

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
			htmlBase := path.Join(append(p.sections, fmt.Sprintf("/%s/%d", paginatePath, pageNumber))...)
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

	if viper.GetBool("disableRSS") {
		return nil
	}

	layouts := p.rssLayouts()

	if layouts == nil {
		// No RSS for this Kind of page.
		return nil
	}

	rssPage := p.copy()
	rssPage.Kind = kindRSS

	// TODO(bep) we zero the date here to get the number of diffs down in
	// testing. But this should be set back later; the RSS feed should
	// inherit the publish date from the node it represents.
	if p.Kind == KindTaxonomy {
		var zeroDate time.Time
		rssPage.Date = zeroDate
	}

	high := 50
	if len(rssPage.Pages) > high {
		rssPage.Pages = rssPage.Pages[:high]
		rssPage.Data["Pages"] = rssPage.Pages
	}
	rssURI := s.Language.GetString("rssURI")

	rssPath := path.Join(append(rssPage.sections, rssURI)...)
	s.setPageURLs(rssPage, rssPath)

	return s.renderAndWriteXML(rssPage.Title,
		rssPage.addLangFilepathPrefix(rssPath), rssPage, s.appendThemeTemplates(layouts)...)
}

func (s *Site) render404() error {
	if viper.GetBool("disable404") {
		return nil
	}

	p := s.newNodePage(kind404)
	p.Title = "404 Page not found"
	p.Data["Pages"] = s.Pages
	p.Pages = s.Pages
	s.setPageURLs(p, "404.html")

	nfLayouts := []string{"404.html"}

	return s.renderAndWritePage("404 page", "404.html", p, s.appendThemeTemplates(nfLayouts)...)

}

func (s *Site) renderSitemap() error {
	if viper.GetBool("disableSitemap") {
		return nil
	}

	sitemapDefault := parseSitemap(viper.GetStringMap("sitemap"))

	n := s.newNodePage(kindSitemap)

	// Include all pages (regular, home page, taxonomies etc.)
	pages := s.Pages

	page := s.newNodePage(kindSitemap)
	page.URLPath.URL = ""
	page.Sitemap.ChangeFreq = sitemapDefault.ChangeFreq
	page.Sitemap.Priority = sitemapDefault.Priority
	page.Sitemap.Filename = sitemapDefault.Filename

	n.Data["Pages"] = pages
	n.Pages = pages

	// TODO(bep) this should be done somewhere else
	for _, page := range pages {
		if page.Sitemap.ChangeFreq == "" {
			page.Sitemap.ChangeFreq = sitemapDefault.ChangeFreq
		}

		if page.Sitemap.Priority == -1 {
			page.Sitemap.Priority = sitemapDefault.Priority
		}

		if page.Sitemap.Filename == "" {
			page.Sitemap.Filename = sitemapDefault.Filename
		}
	}

	smLayouts := []string{"sitemap.xml", "_default/sitemap.xml", "_internal/_default/sitemap.xml"}
	addLanguagePrefix := n.Site.IsMultiLingual()

	return s.renderAndWriteXML("sitemap",
		n.addLangPathPrefixIfFlagSet(page.Sitemap.Filename, addLanguagePrefix), n, s.appendThemeTemplates(smLayouts)...)
}

func (s *Site) renderRobotsTXT() error {
	if !viper.GetBool("enableRobotsTXT") {
		return nil
	}

	n := s.newNodePage(kindRobotsTXT)
	n.Data["Pages"] = s.Pages
	n.Pages = s.Pages

	rLayouts := []string{"robots.txt", "_default/robots.txt", "_internal/_default/robots.txt"}
	outBuffer := bp.GetBuffer()
	defer bp.PutBuffer(outBuffer)
	err := s.renderForLayouts("robots", n, outBuffer, s.appendThemeTemplates(rLayouts)...)

	if err == nil {
		err = s.writeDestFile("robots.txt", outBuffer)
	}

	return err
}

// renderAliases renders shell pages that simply have a redirect in the header.
func (s *Site) renderAliases() error {
	for _, p := range s.Pages {
		if len(p.Aliases) == 0 {
			continue
		}

		plink := p.Permalink()

		for _, a := range p.Aliases {
			if err := s.writeDestAlias(a, plink, p); err != nil {
				return err
			}
		}
	}

	if s.owner.multilingual.enabled() {
		mainLang := s.owner.multilingual.DefaultLang.Lang
		if s.Info.defaultContentLanguageInSubdir {
			mainLangURL := s.Info.pathSpec.AbsURL(mainLang, false)
			jww.DEBUG.Printf("Write redirect to main language %s: %s", mainLang, mainLangURL)
			if err := s.publishDestAlias(s.languageAliasTarget(), "/", mainLangURL, nil); err != nil {
				return err
			}
		} else {
			mainLangURL := s.Info.pathSpec.AbsURL("", false)
			jww.DEBUG.Printf("Write redirect to main language %s: %s", mainLang, mainLangURL)
			if err := s.publishDestAlias(s.languageAliasTarget(), mainLang, mainLangURL, nil); err != nil {
				return err
			}
		}
	}

	return nil
}
