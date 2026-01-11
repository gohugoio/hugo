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
	"context"
	"fmt"
	"path"
	"strings"
	"sync"

	"github.com/bep/logg"
	"github.com/gohugoio/go-radix"
	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/hugolib/doctree"
	"github.com/gohugoio/hugo/tpl/tplimpl"

	"github.com/gohugoio/hugo/config"

	"github.com/gohugoio/hugo/resources/kinds"
	"github.com/gohugoio/hugo/resources/page"
)

type siteRenderContext struct {
	cfg *BuildCfg

	infol logg.LevelLogger

	// languageIdx is the zero based index of the site.
	languageIdx int

	// Zero based index for all output formats combined.
	sitesOutIdx int

	// Zero based index of the output formats configured within a Site.
	// Note that these outputs are sorted.
	outIdx int

	multihost bool
}

// Whether to render 404.html, robotsTXT.txt and similar.
// These are usually rendered once in the root of public.
func (s siteRenderContext) shouldRenderStandalonePage(kind string) bool {
	if s.multihost || kind == kinds.KindSitemap {
		// 1 per site
		return s.outIdx == 0
	}

	if kind == kinds.KindTemporary || kind == kinds.KindStatus404 {
		// 1 for all output formats
		return s.outIdx == 0
	}

	// 1 for all sites and output formats.
	return s.languageIdx == 0 && s.outIdx == 0
}

// renderPages renders pages concurrently.
func (s *Site) renderPages(ctx *siteRenderContext) error {
	numWorkers := config.GetNumWorkerMultiplier()

	results := make(chan error)
	pages := make(chan *pageState, numWorkers) // buffered for performance
	errs := make(chan error)

	go s.errorCollator(results, errs)

	wg := &sync.WaitGroup{}

	for range numWorkers {
		wg.Add(1)
		go pageRenderer(ctx, s, pages, results, wg)
	}

	cfg := ctx.cfg

	w := &doctree.NodeShiftTreeWalker[contentNode]{
		Tree: s.pageMap.treePages,
		Handle: func(key string, n contentNode) (radix.WalkFlag, error) {
			if p, ok := n.(*pageState); ok {
				if cfg.shouldRender(ctx.infol, p) {
					select {
					case <-s.h.Done():
						return radix.WalkStop, nil
					default:
						pages <- p
					}
				}
			}
			return radix.WalkContinue, nil
		},
	}

	if err := w.Walk(context.Background()); err != nil {
		return err
	}

	close(pages)

	wg.Wait()

	close(results)

	err := <-errs
	if err != nil {
		return fmt.Errorf("%v failed to render pages: %w", s.resolveDimensionNames(), err)
	}
	return nil
}

func pageRenderer(
	ctx *siteRenderContext,
	s *Site,
	pages <-chan *pageState,
	results chan<- error,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	sendErr := func(err error) bool {
		select {
		case results <- err:
			return true
		case <-s.h.Done():
			return false
		}
	}

	for p := range pages {

		if p.m.isStandalone() && !ctx.shouldRenderStandalonePage(p.Kind()) {
			continue
		}

		if p.m.pageConfig.Build.PublishResources {
			if err := p.renderResources(); err != nil {
				if sendErr(p.errorf(err, "failed to render resources")) {
					continue
				} else {
					return
				}
			}
		}

		if !p.render {
			// Nothing more to do for this page.
			continue
		}

		templ, found, err := p.resolveTemplate()
		if err != nil {
			if sendErr(p.errorf(err, "failed to resolve template")) {
				continue
			} else {
				return
			}
		}

		if !found {
			s.Log.Trace(
				func() string {
					return fmt.Sprintf("no layout for kind %q found", p.Kind())
				},
			)
			// Don't emit warning for missing 404 etc. pages.
			if !p.m.isStandalone() {
				s.logMissingLayout("", p.Layout(), p.Kind(), p.f.Name)
			}
			continue
		}

		targetPath := p.targetPaths().TargetFilename

		s.Log.Trace(
			func() string {
				return fmt.Sprintf("rendering outputFormat %q kind %q using layout %q to %q", p.pageOutput.f.Name, p.Kind(), templ.Name(), targetPath)
			},
		)

		var d any = p
		switch p.Kind() {
		case kinds.KindSitemapIndex:
			d = s.h.Sites
		}

		if err := s.renderAndWritePage(&s.PathSpec.ProcessingStats.Pages, targetPath, p, d, templ); err != nil {
			if sendErr(err) {
				continue
			} else {
				return
			}
		}

		if p.IsHome() && p.outputFormat().IsHTML && s.isDefault() {
			if err = s.renderDefaultSiteRedirect(p); err != nil {
				if sendErr(err) {
					continue
				} else {
					return
				}
			}
		}

		if p.paginator != nil && p.paginator.current != nil {
			if err := s.renderPaginator(p, templ); err != nil {
				if sendErr(err) {
					continue
				} else {
					return
				}
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
	var args []any
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

	log.Logf(msg, args...)
}

// renderPaginator must be run after the owning Page has been rendered.
func (s *Site) renderPaginator(p *pageState, templ *tplimpl.TemplInfo) error {
	paginatePath := s.Conf.Pagination().Path

	d := p.targetPathDescriptor
	f := p.outputFormat()
	d.Type = f

	if p.paginator.current == nil || p.paginator.current != p.paginator.current.First() {
		panic(fmt.Sprintf("invalid paginator state for %q", p.pathOrTitle()))
	}

	if f.IsHTML && !s.Conf.Pagination().DisableAliases {
		// Write alias for page 1
		d.Addends = fmt.Sprintf("/%s/%d", paginatePath, 1)
		targetPaths := page.CreateTargetPaths(d)

		if err := s.writeDestAlias(targetPaths.TargetFilename, p.Permalink(), f, p); err != nil {
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
			targetPaths.TargetFilename, p, p, templ); err != nil {
			return err
		}

	}

	return nil
}

// renderAliases renders shell pages that simply have a redirect in the header.
func (s *Site) renderAliases() error {
	w := &doctree.NodeShiftTreeWalker[contentNode]{
		Tree: s.pageMap.treePages,
		Handle: func(key string, n contentNode) (radix.WalkFlag, error) {
			p := n.(*pageState)

			// We cannot alias a page that's not rendered.
			if p.m.noLink() || p.skipRender() {
				return radix.WalkContinue, nil
			}

			if len(p.Aliases()) == 0 {
				return radix.WalkContinue, nil
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

					if s.conf.C.IsUglyURLSection(p.Section()) && !strings.HasSuffix(a, ".html") {
						a += ".html"
					}

					lang := p.Language().Lang

					if s.h.Configs.IsMultihost && !strings.HasPrefix(a, "/"+lang) {
						// These need to be in its language root.
						a = path.Join(lang, a)
					}

					err := s.writeDestAlias(a, plink, f, p)
					if err != nil {
						return radix.WalkStop, err
					}
				}
			}
			return radix.WalkContinue, nil
		},
	}
	return w.Walk(context.TODO())
}

// renderDefaultSiteRedirect creates a redirect to the default site's home,
// depending on if it lives in sub folder (e.g. /en) or not.
// The default site is the site is the combination of defaultContentLanguage,
// defaultContentVersion and defaultContentRole.
func (s *Site) renderDefaultSiteRedirect(home *pageState) error {
	if s.conf.DisableDefaultLanguageRedirect || s.conf.DisableDefaultSiteRedirect {
		return nil
	}

	addRedirectInRoot := s.conf.DefaultContentLanguageInSubdir && !s.Conf.IsMultihost()
	addRedirectInRoot = addRedirectInRoot || s.conf.DefaultContentVersionInSubdir || s.conf.DefaultContentRoleInSubdir

	homeLink := home.pageOutput.targetPaths().Link // This doesn't have any baseURL paths in it.

	of := home.outputFormat()
	homePermalink := home.Permalink()

	var ps []string
	if of.Path != "" {
		if addRedirectInRoot {
			// For OutputFormats with a path, creating more than one alias will easily create path clashes without much value.
			ps = []string{paths.AddLeadingAndTrailingSlash(of.Path)}
		}
	} else {
		if addRedirectInRoot {
			ps = append(ps, "/")
		}

		if s.Conf.IsMultilingual() && !s.conf.DefaultContentLanguageInSubdir && !s.Conf.IsMultihost() {
			// Create redirect from e.g. /en => /
			ps = append(ps, homeLink+s.Lang()+"/")
		}

		// /guest/v1.0.0/en/
		//    /,/guest/,/guest/v1.0.0  = /guest/v1.0.0/en/
		// /guest/v1.0.0/
		//    /,/guest/ =  /guest/v1.0.0/
		parts := strings.Split(strings.Trim(homeLink, "/"), "/")

		for i := 0; i < len(parts)-1; i++ {
			ps = append(ps, paths.AddLeadingAndTrailingSlash(strings.Join(parts[0:i+1], "/")))
		}
	}

	if s.h.Configs.IsMultihost {
		prefix := "/" + s.LanguagePrefix()
		for i, p := range ps {
			ps[i] = path.Join(prefix, p)
		}
	}

	for _, p := range ps {
		if err := s.publishDestAlias(true, p, homePermalink, of, home); err != nil {
			return err
		}
	}

	return nil
}
