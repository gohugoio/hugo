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
	"strings"
	"sync"

	"github.com/pkg/errors"

	"github.com/gohugoio/hugo/output"
)

// renderPages renders pages each corresponding to a markdown file.
// TODO(bep np doc
func (s *Site) renderPages(cfg *BuildCfg) error {

	results := make(chan error)
	pages := make(chan *Page)
	errs := make(chan error)

	go s.errorCollator(results, errs)

	numWorkers := getGoMaxProcs() * 4

	wg := &sync.WaitGroup{}

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go pageRenderer(s, pages, results, wg)
	}

	if !cfg.PartialReRender && len(s.headlessPages) > 0 {
		wg.Add(1)
		go headlessPagesPublisher(s, wg)
	}

	for _, page := range s.Pages {
		if cfg.shouldRender(page) {
			pages <- page
		}
	}

	close(pages)

	wg.Wait()

	close(results)

	err := <-errs
	if err != nil {
		return errors.Wrap(err, "failed to render pages")
	}
	return nil
}

func headlessPagesPublisher(s *Site, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, page := range s.headlessPages {
		outFormat := page.outputFormats[0] // There is only one
		if outFormat.Name != s.rc.Format.Name {
			// Avoid double work.
			continue
		}
		pageOutput, err := newPageOutput(page, false, false, outFormat)
		if err == nil {
			page.mainPageOutput = pageOutput
			err = pageOutput.renderResources()
		}

		if err != nil {
			s.Log.ERROR.Printf("Failed to render resources for headless page %q: %s", page, err)
		}
	}
}

func pageRenderer(s *Site, pages <-chan *Page, results chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	for page := range pages {

		for i, outFormat := range page.outputFormats {

			if outFormat.Name != page.s.rc.Format.Name {
				// Will be rendered  ... later.
				continue
			}

			var (
				pageOutput *PageOutput
				err        error
			)

			if i == 0 {
				pageOutput = page.mainPageOutput
			} else {
				pageOutput, err = page.mainPageOutput.copyWithFormat(outFormat, true)
			}

			if err != nil {
				s.Log.ERROR.Printf("Failed to create output page for type %q for page %q: %s", outFormat.Name, page, err)
				continue
			}

			if pageOutput == nil {
				panic("no pageOutput")
			}

			// We only need to re-publish the resources if the output format is different
			// from all of the previous (e.g. the "amp" use case).
			shouldRender := i == 0
			if i > 0 {
				for j := i; j >= 0; j-- {
					if outFormat.Path != page.outputFormats[j].Path {
						shouldRender = true
					} else {
						shouldRender = false
					}
				}
			}

			if shouldRender {
				if err := pageOutput.renderResources(); err != nil {
					s.SendError(page.errorf(err, "failed to render page resources"))
					continue
				}
			}

			var layouts []string

			if page.selfLayout != "" {
				layouts = []string{page.selfLayout}
			} else {
				layouts, err = s.layouts(pageOutput)
				if err != nil {
					s.Log.ERROR.Printf("Failed to resolve layout for output %q for page %q: %s", outFormat.Name, page, err)
					continue
				}
			}

			switch pageOutput.outputFormat.Name {

			case "RSS":
				if err := s.renderRSS(pageOutput); err != nil {
					results <- err
				}
			default:
				targetPath, err := pageOutput.targetPath()
				if err != nil {
					s.Log.ERROR.Printf("Failed to create target path for output %q for page %q: %s", outFormat.Name, page, err)
					continue
				}

				s.Log.DEBUG.Printf("Render %s to %q with layouts %q", pageOutput.Kind, targetPath, layouts)

				if err := s.renderAndWritePage(&s.PathSpec.ProcessingStats.Pages, "page "+pageOutput.FullFilePath(), targetPath, pageOutput, layouts...); err != nil {
					results <- err
				}

				// Only render paginators for the main output format
				if i == 0 && pageOutput.IsNode() {
					if err := s.renderPaginator(pageOutput); err != nil {
						results <- err
					}
				}
			}

		}
	}
}

// renderPaginator must be run after the owning Page has been rendered.
func (s *Site) renderPaginator(p *PageOutput) error {
	if p.paginator != nil {
		s.Log.DEBUG.Printf("Render paginator for page %q", p.Path())
		paginatePath := s.Cfg.GetString("paginatePath")

		// write alias for page 1
		addend := fmt.Sprintf("/%s/%d", paginatePath, 1)
		target, err := p.createTargetPath(p.outputFormat, false, addend)
		if err != nil {
			return err
		}

		// TODO(bep) do better
		link := newOutputFormat(p.Page, p.outputFormat).Permalink()
		if err := s.writeDestAlias(target, link, p.outputFormat, nil); err != nil {
			return err
		}

		pagers := p.paginator.Pagers()

		for i, pager := range pagers {
			if i == 0 {
				// already created
				continue
			}

			pagerNode, err := p.copy()
			if err != nil {
				return err
			}

			pagerNode.origOnCopy = p.Page

			pagerNode.paginator = pager
			if pager.TotalPages() > 0 {
				first, _ := pager.page(0)
				pagerNode.Date = first.Date
				pagerNode.Lastmod = first.Lastmod
			}

			pageNumber := i + 1
			addend := fmt.Sprintf("/%s/%d", paginatePath, pageNumber)
			targetPath, _ := p.targetPath(addend)
			layouts, err := p.layouts()

			if err != nil {
				return err
			}

			if err := s.renderAndWritePage(
				&s.PathSpec.ProcessingStats.PaginatorPages,
				pagerNode.title,
				targetPath, pagerNode, layouts...); err != nil {
				return err
			}

		}
	}
	return nil
}

func (s *Site) renderRSS(p *PageOutput) error {

	if !s.isEnabled(kindRSS) {
		return nil
	}

	limit := s.Cfg.GetInt("rssLimit")
	if limit >= 0 && len(p.Pages) > limit {
		p.Pages = p.Pages[:limit]
		p.data["Pages"] = p.Pages
	}

	layouts, err := s.layoutHandler.For(
		p.layoutDescriptor,
		p.outputFormat)
	if err != nil {
		return err
	}

	targetPath, err := p.targetPath()
	if err != nil {
		return err
	}

	return s.renderAndWriteXML(&s.PathSpec.ProcessingStats.Pages, p.title,
		targetPath, p, layouts...)
}

func (s *Site) render404() error {
	if !s.isEnabled(kind404) {
		return nil
	}

	p := s.newNodePage(kind404)

	p.title = "404 Page not found"
	p.data["Pages"] = s.Pages
	p.Pages = s.Pages
	p.URLPath.URL = "404.html"

	if err := p.initTargetPathDescriptor(); err != nil {
		return err
	}

	nfLayouts := []string{"404.html"}

	htmlOut := output.HTMLFormat
	htmlOut.BaseName = "404"

	pageOutput, err := newPageOutput(p, false, false, htmlOut)
	if err != nil {
		return err
	}

	targetPath, err := pageOutput.targetPath()
	if err != nil {
		s.Log.ERROR.Printf("Failed to create target path for page %q: %s", p, err)
	}

	return s.renderAndWritePage(&s.PathSpec.ProcessingStats.Pages, "404 page", targetPath, pageOutput, nfLayouts...)
}

func (s *Site) renderSitemap() error {
	if !s.isEnabled(kindSitemap) {
		return nil
	}

	sitemapDefault := parseSitemap(s.Cfg.GetStringMap("sitemap"))

	n := s.newNodePage(kindSitemap)

	// Include all pages (regular, home page, taxonomies etc.)
	pages := s.Pages

	page := s.newNodePage(kindSitemap)
	page.URLPath.URL = ""
	if err := page.initTargetPathDescriptor(); err != nil {
		return err
	}
	page.Sitemap.ChangeFreq = sitemapDefault.ChangeFreq
	page.Sitemap.Priority = sitemapDefault.Priority
	page.Sitemap.Filename = sitemapDefault.Filename

	n.data["Pages"] = pages
	n.Pages = pages

	// TODO(bep) we have several of these
	if err := page.initTargetPathDescriptor(); err != nil {
		return err
	}

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

	return s.renderAndWriteXML(&s.PathSpec.ProcessingStats.Sitemaps, "sitemap",
		n.addLangPathPrefixIfFlagSet(page.Sitemap.Filename, addLanguagePrefix), n, smLayouts...)
}

func (s *Site) renderRobotsTXT() error {
	if !s.isEnabled(kindRobotsTXT) {
		return nil
	}

	if !s.Cfg.GetBool("enableRobotsTXT") {
		return nil
	}

	p := s.newNodePage(kindRobotsTXT)
	if err := p.initTargetPathDescriptor(); err != nil {
		return err
	}
	p.data["Pages"] = s.Pages
	p.Pages = s.Pages

	rLayouts := []string{"robots.txt", "_default/robots.txt", "_internal/_default/robots.txt"}

	pageOutput, err := newPageOutput(p, false, false, output.RobotsTxtFormat)
	if err != nil {
		return err
	}

	targetPath, err := pageOutput.targetPath()
	if err != nil {
		s.Log.ERROR.Printf("Failed to create target path for page %q: %s", p, err)
	}

	return s.renderAndWritePage(&s.PathSpec.ProcessingStats.Pages, "Robots Txt", targetPath, pageOutput, rLayouts...)

}

// renderAliases renders shell pages that simply have a redirect in the header.
func (s *Site) renderAliases() error {
	for _, p := range s.Pages {
		if len(p.Aliases) == 0 {
			continue
		}

		for _, f := range p.outputFormats {
			if !f.IsHTML {
				continue
			}

			o := newOutputFormat(p, f)
			plink := o.Permalink()

			for _, a := range p.Aliases {
				if f.Path != "" {
					// Make sure AMP and similar doesn't clash with regular aliases.
					a = path.Join(a, f.Path)
				}

				lang := p.Lang()

				if s.owner.multihost && !strings.HasPrefix(a, "/"+lang) {
					// These need to be in its language root.
					a = path.Join(lang, a)
				}

				if err := s.writeDestAlias(a, plink, f, p); err != nil {
					return err
				}
			}
		}
	}

	if s.owner.multilingual.enabled() && !s.owner.IsMultihost() {
		html, found := s.outputFormatsConfig.GetByName("HTML")
		if found {
			mainLang := s.owner.multilingual.DefaultLang
			if s.Info.defaultContentLanguageInSubdir {
				mainLangURL := s.PathSpec.AbsURL(mainLang.Lang, false)
				s.Log.DEBUG.Printf("Write redirect to main language %s: %s", mainLang, mainLangURL)
				if err := s.publishDestAlias(true, "/", mainLangURL, html, nil); err != nil {
					return err
				}
			} else {
				mainLangURL := s.PathSpec.AbsURL("", false)
				s.Log.DEBUG.Printf("Write redirect to main language %s: %s", mainLang, mainLangURL)
				if err := s.publishDestAlias(true, mainLang.Lang, mainLangURL, html, nil); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
