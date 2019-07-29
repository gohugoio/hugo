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
	"fmt"
	"path"
	"strings"
	"sync"

	"github.com/gohugoio/hugo/config"

	"github.com/gohugoio/hugo/output"
	"github.com/pkg/errors"

	"github.com/gohugoio/hugo/resources/page"
	"github.com/gohugoio/hugo/resources/page/pagemeta"
)

type siteRenderContext struct {
	cfg *BuildCfg

	// Zero based index for all output formats combined.
	sitesOutIdx int

	// Zero based index of the output formats configured within a Site.
	// Note that these outputs are sorted.
	outIdx int

	multihost bool
}

// Whether to render 404.html, robotsTXT.txt which usually is rendered
// once only in the site root.
func (s siteRenderContext) renderSingletonPages() bool {
	if s.multihost {
		// 1 per site
		return s.outIdx == 0
	}

	// 1 for all sites
	return s.sitesOutIdx == 0

}

// renderPages renders pages each corresponding to a markdown file.
// TODO(bep np doc
func (s *Site) renderPages(ctx *siteRenderContext) error {

	numWorkers := config.GetNumWorkerMultiplier()

	results := make(chan error)
	pages := make(chan *pageState, numWorkers) // buffered for performance
	errs := make(chan error)

	go s.errorCollator(results, errs)

	wg := &sync.WaitGroup{}

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go pageRenderer(ctx, s, pages, results, wg)
	}

	cfg := ctx.cfg

	if !cfg.PartialReRender && ctx.outIdx == 0 && len(s.headlessPages) > 0 {
		wg.Add(1)
		go headlessPagesPublisher(s, wg)
	}

L:
	for _, page := range s.workAllPages {
		if cfg.shouldRender(page) {
			select {
			case <-s.h.Done():
				break L
			default:
				pages <- page
			}
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
	for _, p := range s.headlessPages {
		if err := p.renderResources(); err != nil {
			s.SendError(p.errorf(err, "failed to render page resources"))
		}
	}
}

func pageRenderer(
	ctx *siteRenderContext,
	s *Site,
	pages <-chan *pageState,
	results chan<- error,
	wg *sync.WaitGroup) {

	defer wg.Done()

	for p := range pages {
		f := p.outputFormat()

		// TODO(bep) get rid of this odd construct. RSS is an output format.
		if f.Name == "RSS" && !s.isEnabled(kindRSS) {
			continue
		}

		if err := p.renderResources(); err != nil {
			s.SendError(p.errorf(err, "failed to render page resources"))
			continue
		}

		layouts, err := p.getLayouts()
		if err != nil {
			s.Log.ERROR.Printf("Failed to resolve layout for output %q for page %q: %s", f.Name, p, err)
			continue
		}

		targetPath := p.targetPaths().TargetFilename

		if targetPath == "" {
			s.Log.ERROR.Printf("Failed to create target path for output %q for page %q: %s", f.Name, p, err)
			continue
		}

		if err := s.renderAndWritePage(&s.PathSpec.ProcessingStats.Pages, "page "+p.Title(), targetPath, p, layouts...); err != nil {
			results <- err
		}

		if p.paginator != nil && p.paginator.current != nil {
			if err := s.renderPaginator(p, layouts); err != nil {
				results <- err
			}
		}
	}
}

// renderPaginator must be run after the owning Page has been rendered.
func (s *Site) renderPaginator(p *pageState, layouts []string) error {

	paginatePath := s.Cfg.GetString("paginatePath")

	d := p.targetPathDescriptor
	f := p.s.rc.Format
	d.Type = f

	if p.paginator.current == nil || p.paginator.current != p.paginator.current.First() {
		panic(fmt.Sprintf("invalid paginator state for %q", p.pathOrTitle()))
	}

	// Write alias for page 1
	d.Addends = fmt.Sprintf("/%s/%d", paginatePath, 1)
	targetPaths := page.CreateTargetPaths(d)

	if err := s.writeDestAlias(targetPaths.TargetFilename, p.Permalink(), f, nil); err != nil {
		return err
	}

	// Render pages for the rest
	for current := p.paginator.current.Next(); current != nil; current = current.Next() {

		p.paginator.current = current
		d.Addends = fmt.Sprintf("/%s/%d", paginatePath, current.PageNumber())
		targetPaths := page.CreateTargetPaths(d)

		if err := s.renderAndWritePage(
			&s.PathSpec.ProcessingStats.PaginatorPages,
			p.Title(),
			targetPaths.TargetFilename, p, layouts...); err != nil {
			return err
		}

	}

	return nil
}

func (s *Site) render404() error {
	if !s.isEnabled(kind404) {
		return nil
	}

	p, err := newPageStandalone(&pageMeta{
		s:    s,
		kind: kind404,
		urlPaths: pagemeta.URLPath{
			URL: "404.html",
		},
	},
		output.HTMLFormat,
	)

	if err != nil {
		return err
	}

	nfLayouts := []string{"404.html"}

	targetPath := p.targetPaths().TargetFilename

	if targetPath == "" {
		return errors.New("failed to create targetPath for 404 page")
	}

	return s.renderAndWritePage(&s.PathSpec.ProcessingStats.Pages, "404 page", targetPath, p, nfLayouts...)
}

func (s *Site) renderSitemap() error {
	if !s.isEnabled(kindSitemap) {
		return nil
	}

	p, err := newPageStandalone(&pageMeta{
		s:    s,
		kind: kindSitemap,
		urlPaths: pagemeta.URLPath{
			URL: s.siteCfg.sitemap.Filename,
		}},
		output.HTMLFormat,
	)

	if err != nil {
		return err
	}

	targetPath := p.targetPaths().TargetFilename

	if targetPath == "" {
		return errors.New("failed to create targetPath for sitemap")
	}

	smLayouts := []string{"sitemap.xml", "_default/sitemap.xml", "_internal/_default/sitemap.xml"}

	return s.renderAndWriteXML(&s.PathSpec.ProcessingStats.Sitemaps, "sitemap", targetPath, p, smLayouts...)
}

func (s *Site) renderRobotsTXT() error {
	if !s.isEnabled(kindRobotsTXT) {
		return nil
	}

	if !s.Cfg.GetBool("enableRobotsTXT") {
		return nil
	}

	p, err := newPageStandalone(&pageMeta{
		s:    s,
		kind: kindRobotsTXT,
		urlPaths: pagemeta.URLPath{
			URL: "robots.txt",
		},
	},
		output.RobotsTxtFormat)

	if err != nil {
		return err
	}

	rLayouts := []string{"robots.txt", "_default/robots.txt", "_internal/_default/robots.txt"}

	return s.renderAndWritePage(&s.PathSpec.ProcessingStats.Pages, "Robots Txt", p.targetPaths().TargetFilename, p, rLayouts...)

}

// renderAliases renders shell pages that simply have a redirect in the header.
func (s *Site) renderAliases() error {
	for _, p := range s.workAllPages {

		if len(p.Aliases()) == 0 {
			continue
		}

		for _, of := range p.OutputFormats() {
			if !of.Format.IsHTML {
				continue
			}

			plink := of.Permalink()
			f := of.Format

			for _, a := range p.Aliases() {
				isRelative := !strings.HasPrefix(a, "/")

				if isRelative {
					// Make alias relative, where "." will be on the
					// same directory level as the current page.
					// TODO(bep) ugly URLs doesn't seem to be supported in
					// aliases, I'm not sure why not.
					basePath := of.RelPermalink()
					if strings.HasSuffix(basePath, "/") {
						basePath = path.Join(basePath, "..")
					}
					a = path.Join(basePath, a)

				} else if f.Path != "" {
					// Make sure AMP and similar doesn't clash with regular aliases.
					a = path.Join(f.Path, a)
				}

				lang := p.Language().Lang

				if s.h.multihost && !strings.HasPrefix(a, "/"+lang) {
					// These need to be in its language root.
					a = path.Join(lang, a)
				}

				if err := s.writeDestAlias(a, plink, f, p); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// renderMainLanguageRedirect creates a redirect to the main language home,
// depending on if it lives in sub folder (e.g. /en) or not.
func (s *Site) renderMainLanguageRedirect() error {

	if !s.h.multilingual.enabled() || s.h.IsMultihost() {
		// No need for a redirect
		return nil
	}

	if s.Cfg.GetBool("disableDefaultLanguageSubdirRedir") {
		// Do not generate the redirect to default language in subfolder.
		// todo( Write in docs how to use URL front-matter to place  the /index.html file.
		return nil
	}

	html, found := s.outputFormatsConfig.GetByName("HTML")
	if found {
		mainLang := s.h.multilingual.DefaultLang
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

	return nil
}
