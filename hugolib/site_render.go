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

	"github.com/gohugoio/hugo/output"

	"github.com/gohugoio/hugo/tpl"

	"github.com/gohugoio/hugo/config"

	"github.com/pkg/errors"

	"github.com/gohugoio/hugo/resources/page"
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
func (s siteRenderContext) shouldRenderSingletonPages() bool {
	if s.multihost {
		// 1 per site
		return s.outIdx == 0
	}

	// 1 for all sites
	return s.sitesOutIdx == 0
}

// renderPages renders this Site's pages for the output format defined in ctx.
func (s *Site) renderPages(ctx *siteRenderContext) error {
	numWorkers := config.GetNumWorkerMultiplier()

	results := make(chan error)
	pages := make(chan *pageState, numWorkers) // buffered for performance
	errs := make(chan error)

	go s.errorCollator(results, errs)

	wg := &sync.WaitGroup{}
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go s.renderPage(ctx, pages, results, wg)
	}

	cfg := ctx.cfg

	s.pageMap.WalkPagesAllPrefixSection("", nil, nil, func(np contentNodeProvider) bool {
		n := np.GetNode()
		if ctx.outIdx > 0 && n.p.getTreeRef().GetNode().IsStandalone() {
			// Only render the standalone pages (e.g. 404) once.
			return false
		}

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

func (s *Site) renderPage(
	ctx *siteRenderContext,
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
			s.logMissingLayout("", p.Layout(), p.Kind(), p.f.Name)
			continue
		}

		targetPath := p.targetPaths().TargetFilename

		var statCounter *uint64
		switch p.outputFormat().Name {
		case output.SitemapFormat.Name:
			statCounter = &s.PathSpec.ProcessingStats.Sitemaps
		default:
			statCounter = &s.PathSpec.ProcessingStats.Pages
		}

		if err := s.renderAndWritePage(statCounter, "page "+p.Title(), targetPath, p, templ); err != nil {
			results <- err
		}

		if p.paginator != nil && p.paginator.current != nil {
			if err := s.renderPaginator(p, templ); err != nil {
				results <- err
			}
		}
	}
}

func (s *Site) logMissingLayout(name, layout, kind, outputFormat string) {
	log := s.Log.Warn()
	if name != "" && infoOnMissingLayout[name] {
		log = s.Log.Info()
	}

	errMsg := "You should create a template file which matches Hugo Layouts Lookup Rules for this combination."
	var args []interface{}
	msg := "found no layout file for"
	if outputFormat != "" {
		msg += " %q"
		args = append(args, outputFormat)
	}

	if layout != "" {
		msg += " for layout %q"
		args = append(args, layout)
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

// renderAliases renders shell pages that simply have a redirect in the header.
func (s *Site) renderAliases() error {
	var err error

	s.pageMap.WalkPagesAllPrefixSection("", nil, contentTreeNoLinkFilter, func(np contentNodeProvider) bool {
		n := np.GetNode()
		p := n.p

		if len(p.Aliases()) == 0 {
			return false
		}
		pathSeen := make(map[string]bool)

		for _, of := range p.OutputFormats() {
			if !of.Format.IsHTML {
				continue
			}

			f := of.Format

			if pathSeen[f.Path] {
				continue
			}
			pathSeen[f.Path] = true

			plink := of.Permalink()

			for _, a := range p.Aliases() {
				isRelative := !strings.HasPrefix(a, "/")

				if isRelative {
					// Make alias relative, where "." will be on the
					// same directory level as the current page.
					basePath := path.Join(p.targetPaths().SubResourceBaseLink, "..")
					a = path.Join(basePath, a)

				} else {
					// Make sure AMP and similar doesn't clash with regular aliases.
					a = path.Join(f.Path, a)
				}

				if s.UglyURLs && !strings.HasSuffix(a, ".html") {
					a += ".html"
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
			mainLangURL := s.PathSpec.AbsURL(mainLang.Lang+"/", false)
			s.Log.Debugf("Write redirect to main language %s: %s", mainLang, mainLangURL)
			if err := s.publishDestAlias(true, "/", mainLangURL, html, nil); err != nil {
				return err
			}
		} else {
			mainLangURL := s.PathSpec.AbsURL("", false)
			s.Log.Debugf("Write redirect to main language %s: %s", mainLang, mainLangURL)
			if err := s.publishDestAlias(true, mainLang.Lang, mainLangURL, html, nil); err != nil {
				return err
			}
		}
	}

	return nil
}
