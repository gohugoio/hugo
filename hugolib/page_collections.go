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
)

// TODO(bep) np pages names
// TODO(bep) np this is a somewhat breaking change and should be doc. + release notes: See AllPages vs. "this language only". Looks like it is like this alread, check.
type PageCollections struct {
	// Includes only pages of NodePage type, and only pages in the current language.
	Pages Pages

	// Includes all pages in all languages, including the current one.
	// Only pages of NodePage type.
	AllPages Pages

	// Includes pages of all types, but only pages in the current language.
	Nodes Pages

	// Includes all pages in all languages, including the current one.
	// Includes pages of all types.
	AllNodes Pages

	// A convenience cache for the traditional node types, taxonomies, home page etc.
	// This is for the current language only.
	indexNodes Pages

	// Includes absolute all pages (of all types), including drafts etc.
	rawAllPages Pages
}

func (c *PageCollections) refreshPageCaches() {
	// All pages are stored in AllNodes and Nodes. Filter from those.
	c.Pages = c.findPagesByNodeTypeIn(NodePage, c.Nodes)
	c.indexNodes = c.findPagesByNodeTypeNotIn(NodePage, c.Nodes)
	c.AllPages = c.findPagesByNodeTypeIn(NodePage, c.AllNodes)

	for _, n := range c.Nodes {
		if n.NodeType == NodeUnknown {
			panic(fmt.Sprintf("Got unknown type %s", n.Title))
		}
	}
}

func newPageCollections() *PageCollections {
	return &PageCollections{}
}

func newPageCollectionsFromPages(pages Pages) *PageCollections {
	return &PageCollections{rawAllPages: pages}
}

// TODO(bep) np clean and remove finders

func (c *PageCollections) findPagesByNodeType(n NodeType) Pages {
	return c.findPagesByNodeTypeIn(n, c.Nodes)
}

func (c *PageCollections) getPage(n NodeType, path ...string) *Page {
	pages := c.findPagesByNodeTypeIn(n, c.Nodes)

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

func (c *PageCollections) findIndexNodesByNodeType(n NodeType) Pages {
	return c.findPagesByNodeTypeIn(n, c.indexNodes)
}

func (*PageCollections) findPagesByNodeTypeIn(n NodeType, inPages Pages) Pages {
	var pages Pages
	for _, p := range inPages {
		if p.NodeType == n {
			pages = append(pages, p)
		}
	}
	return pages
}

func (*PageCollections) findPagesByNodeTypeNotIn(n NodeType, inPages Pages) Pages {
	var pages Pages
	for _, p := range inPages {
		if p.NodeType != n {
			pages = append(pages, p)
		}
	}
	return pages
}

func (c *PageCollections) findAllPagesByNodeType(n NodeType) Pages {
	return c.findPagesByNodeTypeIn(n, c.rawAllPages)
}

func (c *PageCollections) findRawAllPagesByNodeType(n NodeType) Pages {
	return c.findPagesByNodeTypeIn(n, c.rawAllPages)
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

func (c *PageCollections) replacePage(page *Page) {
	// will find existing page that matches filepath and remove it
	c.removePage(page)
	c.addPage(page)
}
