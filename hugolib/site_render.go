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

	"github.com/gohugoio/hugo/tpl"

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

	s.pageMap.pageTrees.Walk(func(ss string, n *contentNode) bool {
		if cfg.shouldRender(n.p) {
			select {
			case <-s.h.Done():
				return true
			default:
				pages <- n.p
			}
		}
		return false
	})

	close(pages)

	wg.Wait()

	close(results)

	err := <-errs
	if err != nil {
		return errors.Wrap(err, "failed to render pages")
	}
	return nil
}

func pageRenderer(
	ctx *siteRenderContext,
	s *Site,
	pages <-chan *pageState,
	results chan<- error,
	wg *sync.WaitGroup) {

	defer wg.Done()

	for p := range pages {
		if p.m.buildConfig.PublishResources {
			if err := p.renderResources(); err != nil {
				s.SendError(p.errorf(err, "failed to render page resources"))
				continue
			}
		}

		if !p.render {
			// Nothing more to do for this page.
			continue
		}

		templ, found, err := p.resolveTemplate()
		if err != nil {
			s.SendError(p.errorf(err, "failed to resolve template"))
			continue
		}

		if !found {
			s.logMissingLayout("", p.Kind(), p.f.Name)
			continue
		}

		targetPath := p.targetPaths().TargetFilename

		if err := s.renderAndWritePage(&s.PathSpec.ProcessingStats.Pages, "page "+p.Title(), targetPath, p, templ); err != nil {
			results <- err
		}

		if p.paginator != nil && p.paginator.current != nil {
			if err := s.renderPaginator(p, templ); err != nil {
				results <- err
			}
		}
	}
}

func (s *Site) logMissingLayout(name, kind, outputFormat string) {
	log := s.Log.WARN
	if name != "" && infoOnMissingLayout[name] {
		log = s.Log.INFO
	}

	errMsg := "You should create a template file which matches Hugo Layouts Lookup Rules for this combination."
	var args []interface{}
	msg := "found no layout file for"
	if outputFormat != "" {
		msg += " %q"
		args = append(args, outputFormat)
	}

	if kind != "" {
		msg += " for kind %q"
		args = append(args, kind)
	}

	if name != "" {
		msg += " for %q"
		args = append(args, name)
	}

	msg += ": " + errMsg

	log.Printf(msg, args...)
}

// renderPaginator must be run after the owning Page has been rendered.
func (s *Site) renderPaginator(p *pageState, templ tpl.Template) error {

	paginatePath := s.Cfg.GetString("paginatePath")

	d := p.targetPathDescriptor
	f := p.s.rc.Format
	d.Type = f

	if p.paginator.current == nil || p.paginator.current != p.paginator.current.First() {
		panic(fmt.Sprintf("invalid paginator state for %q", p.pathOrTitle()))
	}

	if f.IsHTML {
		// Write alias for page 1
		d.Addends = fmt.Sprintf("/%s/%d", paginatePath, 1)
		targetPaths := page.CreateTargetPaths(d)

		if err := s.writeDestAlias(targetPaths.TargetFilename, p.Permalink(), f, nil); err != nil {
			return err
		}
	}

	// Render pages for the rest
	for current := p.paginator.current.Next(); current != nil; current = current.Next() {

		p.paginator.current = current
		d.Addends = fmt.Sprintf("/%s/%d", paginatePath, current.PageNumber())
		targetPaths := page.CreateTargetPaths(d)

		if err := s.renderAndWritePage(
			&s.PathSpec.ProcessingStats.PaginatorPages,
			p.Title(),
			targetPaths.TargetFilename, p, templ); err != nil {
			return err
		}

	}

	return nil
}

func (s *Site) render404() error {
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

	if !p.render {
		return nil
	}

	var d output.LayoutDescriptor
	d.Kind = kind404

	templ, found, err := s.Tmpl().LookupLayout(d, output.HTMLFormat)
	if err != nil {
		return err
	}
	if !found {
		return nil
	}

	targetPath := p.targetPaths().TargetFilename

	if targetPath == "" {
		return errors.New("failed to create targetPath for 404 page")
	}

	return s.renderAndWritePage(&s.PathSpec.ProcessingStats.Pages, "404 page", targetPath, p, templ)
}

func (s *Site) renderSitemap() error {
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

	if !p.render {
		return nil
	}

	targetPath := p.targetPaths().TargetFilename

	if targetPath == "" {
		return errors.New("failed to create targetPath for sitemap")
	}

	templ := s.lookupLayouts("sitemap.xml", "_default/sitemap.xml", "_internal/_default/sitemap.xml")

	return s.renderAndWriteXML(&s.PathSpec.ProcessingStats.Sitemaps, "sitemap", targetPath, p, templ)
}

func (s *Site) renderRobotsTXT() error {
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

	if !p.render {
		return nil
	}

	templ := s.lookupLayouts("robots.txt", "_default/robots.txt", "_internal/_default/robots.txt")

	return s.renderAndWritePage(&s.PathSpec.ProcessingStats.Pages, "Robots Txt", p.targetPaths().TargetFilename, p, templ)

}

// renderAliases renders shell pages that simply have a redirect in the header.
func (s *Site) renderAliases() error {
	var err error
	s.pageMap.pageTrees.WalkRenderable(func(ss string, n *contentNode) bool {
		p := n.p
		if len(p.Aliases()) == 0 {
			return false
		}

		for _, of := range p.OutputFormats() {
			if !of.Format.IsHTML {
				return false
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

				err = s.writeDestAlias(a, plink, f, p)
				if err != nil {
					return true
				}
			}
		}
		return false
	})

	return err
}

// renderMainLanguageRedirect creates a redirect to the main language home,
// depending on if it lives in sub folder (e.g. /en) or not.
func (s *Site) renderMainLanguageRedirect() error {

	if !s.h.multilingual.enabled() || s.h.IsMultihost() {
		// No need for a redirect
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
