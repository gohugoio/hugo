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

	"github.com/gohugoio/hugo/helpers"

	"github.com/gohugoio/hugo/output"
	"github.com/pkg/errors"

	"github.com/gohugoio/hugo/resources/page"
	"github.com/gohugoio/hugo/resources/page/pagemeta"
)

// renderPages renders pages each corresponding to a markdown file.
// TODO(bep np doc
func (s *Site) renderPages(cfg *BuildCfg) error {

	results := make(chan error)
	pages := make(chan *pageState)
	errs := make(chan error)

	go s.errorCollator(results, errs)

	numWorkers := getGoMaxProcs() * 4

	wg := &sync.WaitGroup{}

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go pageRenderer(s, pages, results, wg)
	}

	// TODO(bep) page
	/*	if !cfg.PartialReRender && len(s.headlessPages) > 0 {
		wg.Add(1)
		go headlessPagesPublisher(s, wg)
	}*/

	for _, page := range s.workAllPages {
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

// TODO(bep) page fixme
func headlessPagesPublisher(s *Site, wg *sync.WaitGroup) {
	if true {
		return
	}
	defer wg.Done()
	for _, page := range s.headlessPages {
		// TODO(bep) page
		outFormat := page.currentOutputFormat // There is only one
		if outFormat.Name != s.rc.Format.Name {
			// Avoid double work.
			continue
		}

		// TODO(bep) page
		//if err == nil {
		//page.p.mainPageOutput = pageOutput
		//err = pageOutput.renderResources()
		//}

		//if err != nil {
		//	s.Log.ERROR.Printf("Failed to render resources for headless page %q: %s", page, err)
		//}
	}
}

func pageRenderer(s *Site, pages <-chan *pageState, results chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	for p := range pages {
		for i, f := range p.m.outputFormats() {
			// TODO(bep) get rid of this odd construct. RSS is an output format.
			if f.Name == "RSS" && !s.isEnabled(kindRSS) {
				continue
			}

			if f.Name != s.rc.Format.Name {
				// Rendered later.
				continue
			}

			if i == 0 {
				// This will publish its resources to all output formats paths.
				if err := p.renderResources(); err != nil {
					s.SendError(p.errorf(err, "failed to render page resources"))
					continue
				}
			}

			layouts, err := p.getLayouts(f)
			if err != nil {
				s.Log.ERROR.Printf("Failed to resolve layout for output %q for page %q: %s", f.Name, p, err)
				continue
			}

			targetPath := p.targetPath()

			if targetPath == "" {
				s.Log.ERROR.Printf("Failed to create target path for output %q for page %q: %s", f.Name, p, err)
				continue
			}

			if err := s.renderAndWritePage(&s.PathSpec.ProcessingStats.Pages, "page "+p.Title(), targetPath, p, layouts...); err != nil {
				results <- err
			}

			// Only render paginators for the main output format
			if i == 0 && p.paginator != nil && p.paginator.current != nil {
				if err := s.renderPaginator(p, layouts); err != nil {
					results <- err
				}
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

	// Rewind
	p.paginator.current = p.paginator.current.First()

	// Write alias for page 1
	d.Addends = fmt.Sprintf("/%s/%d", paginatePath, 1)
	targetPath := page.CreateTargetPath(d)

	if err := s.writeDestAlias(targetPath, p.Permalink(), f, nil); err != nil {
		return err
	}

	// Render pages for the rest
	for current := p.paginator.current.Next(); current != nil; current = current.Next() {

		p.paginator.current = current
		d.Addends = fmt.Sprintf("/%s/%d", paginatePath, current.PageNumber())
		targetPath := page.CreateTargetPath(d)

		if err := s.renderAndWritePage(
			&s.PathSpec.ProcessingStats.PaginatorPages,
			p.Title(),
			targetPath, p, layouts...); err != nil {
			return err
		}

	}

	return nil
}

func (s *Site) render404() error {
	if !s.isEnabled(kind404) {
		return nil
	}

	p, err := newStandalonePage(&pageMeta{
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

	// TODO(bep) page
	p.initOutputFormat(output.HTMLFormat, true)

	nfLayouts := []string{"404.html"}

	targetPath := p.targetPath()

	if targetPath == "" {
		return errors.New("failed to create targetPath for 404 page")
	}

	return s.renderAndWritePage(&s.PathSpec.ProcessingStats.Pages, "404 page", targetPath, p, nfLayouts...)
}

func (s *Site) renderSitemap() error {
	if !s.isEnabled(kindSitemap) {
		return nil
	}

	p, err := newStandalonePage(&pageMeta{
		s:    s,
		kind: kindSitemap,
		urlPaths: pagemeta.URLPath{
			URL: s.siteConfigHolder.sitemap.Filename,
		}},
		output.HTMLFormat,
	)

	if err != nil {
		return err
	}

	// TODO(bep) page
	p.initOutputFormat(output.HTMLFormat, true)

	targetPath := p.targetPath()

	if targetPath == "" {
		return errors.New("failed to create targetPath for sitemap")
	}

	// TODO(bep) page check/consolidate
	if s.Info.IsMultiLingual() {
		targetPath = helpers.FilePathSeparator + s.Language().Lang + helpers.FilePathSeparator + targetPath
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

	p, err := newStandalonePage(&pageMeta{
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

	// TODO(bep) page
	p.initOutputFormat(output.RobotsTxtFormat, true)

	rLayouts := []string{"robots.txt", "_default/robots.txt", "_internal/_default/robots.txt"}

	return s.renderAndWritePage(&s.PathSpec.ProcessingStats.Pages, "Robots Txt", p.targetPath(), p, rLayouts...)

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
				if f.Path != "" {
					// Make sure AMP and similar doesn't clash with regular aliases.
					a = path.Join(a, f.Path)
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

	if s.h.multilingual.enabled() && !s.h.IsMultihost() {
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
	}

	return nil
}
