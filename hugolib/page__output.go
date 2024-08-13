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

	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/output"
	"github.com/gohugoio/hugo/resources/page"
	"github.com/gohugoio/hugo/resources/resource"
)

func newPageOutput(
	ps *pageState,
	pp pagePaths,
	f output.Format,
	render bool,
) *pageOutput {
	var targetPathsProvider targetPathsHolder
	var linksProvider resource.ResourceLinksProvider

	ft, found := pp.targetPaths[f.Name]
	if !found {
		// Link to the main output format
		ft = pp.targetPaths[pp.firstOutputFormat.Format.Name]
	}
	targetPathsProvider = ft
	linksProvider = ft

	var paginatorProvider page.PaginatorProvider
	var pag *pagePaginator

	if render && ps.IsNode() {
		pag = newPagePaginator(ps)
		paginatorProvider = pag
	} else {
		paginatorProvider = page.PaginatorNotSupportedFunc(func() error {
			return fmt.Errorf("pagination not supported for this page: %s", ps.getPageInfoForError())
		})
	}

	providers := struct {
		page.PaginatorProvider
		resource.ResourceLinksProvider
		targetPather
	}{
		paginatorProvider,
		linksProvider,
		targetPathsProvider,
	}

	po := &pageOutput{
		p:                       ps,
		f:                       f,
		pagePerOutputProviders:  providers,
		MarkupProvider:          page.NopPage,
		ContentProvider:         page.NopPage,
		PageRenderProvider:      page.NopPage,
		TableOfContentsProvider: page.NopPage,
		render:                  render,
		paginator:               pag,
		dependencyManagerOutput: ps.s.Conf.NewIdentityManager((ps.Path() + "/" + f.Name)),
	}

	return po
}

// We create a pageOutput for every output format combination, even if this
// particular page isn't configured to be rendered to that format.
type pageOutput struct {
	p *pageState

	// Set if this page isn't configured to be rendered to this format.
	render bool

	f output.Format

	// Only set if render is set.
	// Note that this will be lazily initialized, so only used if actually
	// used in template(s).
	paginator *pagePaginator

	// These interface provides the functionality that is specific for this
	// output format.
	contentRenderer page.ContentRenderer
	pagePerOutputProviders
	page.MarkupProvider
	page.ContentProvider
	page.PageRenderProvider
	page.TableOfContentsProvider
	page.RenderShortcodesProvider

	// May be nil.
	pco *pageContentOutput

	dependencyManagerOutput identity.Manager

	renderState int  // Reset when it needs to be rendered again.
	renderOnce  bool // To make sure we at least try to render it once.
}

func (po *pageOutput) incrRenderState() {
	po.renderState++
	po.renderOnce = true
}

// isRendered reports whether this output format or its content has been rendered.
func (po *pageOutput) isRendered() bool {
	if po.renderState > 0 {
		return true
	}
	if po.pco != nil && po.pco.contentRendered.Load() {
		return true
	}
	return false
}

func (po *pageOutput) IdentifierBase() string {
	return po.p.f.Name
}

func (po *pageOutput) GetDependencyManager() identity.Manager {
	return po.dependencyManagerOutput
}

func (p *pageOutput) setContentProvider(cp *pageContentOutput) {
	if cp == nil {
		return
	}
	p.contentRenderer = cp
	p.ContentProvider = cp
	p.MarkupProvider = cp
	p.PageRenderProvider = cp
	p.TableOfContentsProvider = cp
	p.RenderShortcodesProvider = cp
	p.pco = cp
}
