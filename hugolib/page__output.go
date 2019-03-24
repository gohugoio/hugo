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
	"github.com/gohugoio/hugo/output"
	"github.com/gohugoio/hugo/resources/page"
	"github.com/gohugoio/hugo/resources/resource"
)

func newPageOutput(
	cp *pageContentOutput, // may be nil
	ps *pageState,
	pp pagePaths,
	f output.Format,
	render bool) *pageOutput {

	var targetPathsProvider targetPathsHolder
	var linksProvider resource.ResourceLinksProvider

	ft, found := pp.targetPaths[f.Name]
	if !found {
		// Link to the main output format
		ft = pp.targetPaths[pp.OutputFormats()[0].Format.Name]
	}
	targetPathsProvider = ft
	linksProvider = ft

	var paginatorProvider page.PaginatorProvider = page.NopPage
	var pag *pagePaginator

	if render && ps.IsNode() {
		pag = &pagePaginator{source: ps}
		paginatorProvider = pag
	}

	var contentProvider page.ContentProvider = page.NopPage
	var tableOfContentsProvider page.TableOfContentsProvider = page.NopPage

	if cp != nil {
		contentProvider = cp
		tableOfContentsProvider = cp
	}

	providers := struct {
		page.ContentProvider
		page.TableOfContentsProvider
		page.PaginatorProvider
		resource.ResourceLinksProvider
		targetPather
	}{
		contentProvider,
		tableOfContentsProvider,
		paginatorProvider,
		linksProvider,
		targetPathsProvider,
	}

	po := &pageOutput{
		f:                      f,
		cp:                     cp,
		pagePerOutputProviders: providers,
		render:                 render,
		paginator:              pag,
	}

	return po

}

// We create a pageOutput for every output format combination, even if this
// particular page isn't configured to be rendered to that format.
type pageOutput struct {
	// Set if this page isn't configured to be rendered to this format.
	render bool

	f output.Format

	// Only set if render is set.
	// Note that this will be lazily initialized, so only used if actually
	// used in template(s).
	paginator *pagePaginator

	// This interface provides the functionality that is specific for this
	// output format.
	pagePerOutputProviders

	// This may be nil.
	cp *pageContentOutput
}

func (p *pageOutput) enablePlaceholders() {
	if p.cp != nil {
		p.cp.enablePlaceholders()
	}
}
