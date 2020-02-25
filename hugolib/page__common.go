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
	"sync"

	"github.com/bep/gitmap"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/compare"
	"github.com/gohugoio/hugo/lazy"
	"github.com/gohugoio/hugo/navigation"
	"github.com/gohugoio/hugo/output"
	"github.com/gohugoio/hugo/resources/page"
	"github.com/gohugoio/hugo/resources/resource"
)

type treeRefProvider interface {
	getTreeRef() *contentTreeRef
}

func (p *pageCommon) getTreeRef() *contentTreeRef {
	return p.treeRef
}

type nextPrevProvider interface {
	getNextPrev() *nextPrev
}

func (p *pageCommon) getNextPrev() *nextPrev {
	return p.posNextPrev
}

type nextPrevInSectionProvider interface {
	getNextPrevInSection() *nextPrev
}

func (p *pageCommon) getNextPrevInSection() *nextPrev {
	return p.posNextPrevSection
}

type pageCommon struct {
	s *Site
	m *pageMeta

	bucket  *pagesMapBucket
	treeRef *contentTreeRef

	// Laziliy initialized dependencies.
	init *lazy.Init

	// All of these represents the common parts of a page.Page
	maps.Scratcher
	navigation.PageMenusProvider
	page.AuthorProvider
	page.PageRenderProvider
	page.AlternativeOutputFormatsProvider
	page.ChildCareProvider
	page.FileProvider
	page.GetPageProvider
	page.GitInfoProvider
	page.InSectionPositioner
	page.OutputFormatsProvider
	page.PageMetaProvider
	page.Positioner
	page.RawContentProvider
	page.RelatedKeywordsProvider
	page.RefProvider
	page.ShortcodeInfoProvider
	page.SitesProvider
	page.DeprecatedWarningPageMethods
	page.TranslationsProvider
	page.TreeProvider
	resource.LanguageProvider
	resource.ResourceDataProvider
	resource.ResourceMetaProvider
	resource.ResourceParamsProvider
	resource.ResourceTypeProvider
	resource.MediaTypeProvider
	resource.TranslationKeyProvider
	compare.Eqer

	// Describes how paths and URLs for this page and its descendants
	// should look like.
	targetPathDescriptor page.TargetPathDescriptor

	layoutDescriptor     output.LayoutDescriptor
	layoutDescriptorInit sync.Once

	// The parsed page content.
	pageContent

	// Set if feature enabled and this is in a Git repo.
	gitInfo *gitmap.GitInfo

	// Positional navigation
	posNextPrev        *nextPrev
	posNextPrevSection *nextPrev

	// Menus
	pageMenus *pageMenus

	// Internal use
	page.InternalDependencies

	// The children. Regular pages will have none.
	*pagePages

	// Any bundled resources
	resources            resource.Resources
	resourcesInit        sync.Once
	resourcesPublishInit sync.Once

	translations    page.Pages
	allTranslations page.Pages

	// Calculated an cached translation mapping key
	translationKey     string
	translationKeyInit sync.Once

	// Will only be set for bundled pages.
	parent *pageState

	// Set in fast render mode to force render a given page.
	forceRender bool
}

type pagePages struct {
	pagesInit sync.Once
	pages     page.Pages

	regularPagesInit          sync.Once
	regularPages              page.Pages
	regularPagesRecursiveInit sync.Once
	regularPagesRecursive     page.Pages
}
