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

// PageCollections contains the page collections for a site.
type PageCollections struct {
	// Includes only pages of all types, and only pages in the current language.
	Pages Pages

	// Includes all pages in all languages, including the current one.
	// Includes pages of all types.
	AllPages Pages

	// A convenience cache for the traditional index types, taxonomies, home page etc.
	// This is for the current language only.
	indexPages Pages

	// A convenience cache for the regular pages.
	// This is for the current language only.
	RegularPages Pages

	// A convenience cache for the all the regular pages.
	AllRegularPages Pages

	// Includes absolute all pages (of all types), including drafts etc.
	rawAllPages Pages
}

func (c *PageCollections) refreshPageCaches() {
	c.indexPages = c.findPagesByKindNotIn(KindPage, c.Pages)
	c.RegularPages = c.findPagesByKindIn(KindPage, c.Pages)
	c.AllRegularPages = c.findPagesByKindIn(KindPage, c.AllPages)
}

func newPageCollections() *PageCollections {
	return &PageCollections{}
}

func newPageCollectionsFromPages(pages Pages) *PageCollections {
	return &PageCollections{rawAllPages: pages}
}

func (c *PageCollections) getPage(typ string, path ...string) *Page {
	pages := c.findPagesByKindIn(typ, c.Pages)

	if len(pages) == 0 {
		return nil
	}

	if len(path) == 0 && len(pages) == 1 {
		return pages[0]
	}

	for _, p := range pages {
		match := false
		for i := 0; i < len(path); i++ {
			if len(p.sections) > i && path[i] == p.sections[i] {
				match = true
			} else {
				match = false
				break
			}
		}
		if match {
			return p
		}
	}

	return nil
}

func (*PageCollections) findPagesByKindIn(kind string, inPages Pages) Pages {
	var pages Pages
	for _, p := range inPages {
		if p.Kind == kind {
			pages = append(pages, p)
		}
	}
	return pages
}

func (*PageCollections) findPagesByKindNotIn(kind string, inPages Pages) Pages {
	var pages Pages
	for _, p := range inPages {
		if p.Kind != kind {
			pages = append(pages, p)
		}
	}
	return pages
}

func (c *PageCollections) findPagesByKind(kind string) Pages {
	return c.findPagesByKindIn(kind, c.Pages)
}

func (c *PageCollections) addPage(page *Page) {
	c.rawAllPages = append(c.rawAllPages, page)
}

func (c *PageCollections) removePageByPath(path string) {
	if i := c.rawAllPages.FindPagePosByFilePath(path); i >= 0 {
		c.rawAllPages = append(c.rawAllPages[:i], c.rawAllPages[i+1:]...)
	}
}

func (c *PageCollections) removePage(page *Page) {
	if i := c.rawAllPages.FindPagePos(page); i >= 0 {
		c.rawAllPages = append(c.rawAllPages[:i], c.rawAllPages[i+1:]...)
	}
}

func (c *PageCollections) findPagesByShortcode(shortcode string) Pages {
	var pages Pages

	for _, p := range c.rawAllPages {
		if p.shortcodeState != nil {
			if _, ok := p.shortcodeState.nameSet[shortcode]; ok {
				pages = append(pages, p)
			}
		}
	}
	return pages
}

func (c *PageCollections) replacePage(page *Page) {
	// will find existing page that matches filepath and remove it
	c.removePage(page)
	c.addPage(page)
}
