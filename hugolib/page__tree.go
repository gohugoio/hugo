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
	"strings"

	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/common/types"
	"github.com/gohugoio/hugo/hugolib/doctree"
	"github.com/gohugoio/hugo/resources/kinds"
	"github.com/gohugoio/hugo/resources/page"
)

// pageTree holds the treen navigational method for a Page.
type pageTree struct {
	p *pageState
}

func (pt pageTree) IsAncestor(other any) bool {
	n, ok := other.(contentNode)
	if !ok {
		return false
	}

	if n.Path() == pt.p.Path() {
		return false
	}

	return strings.HasPrefix(n.Path(), paths.AddTrailingSlash(pt.p.Path()))
}

func (pt pageTree) IsDescendant(other any) bool {
	n, ok := other.(contentNode)
	if !ok {
		return false
	}

	if n.Path() == pt.p.Path() {
		return false
	}

	return strings.HasPrefix(pt.p.Path(), paths.AddTrailingSlash(n.Path()))
}

func (pt pageTree) CurrentSection() page.Page {
	if kinds.IsBranch(pt.p.Kind()) {
		return pt.p
	}

	dir := pt.p.m.pathInfo.Dir()

	if dir == "/" {
		if pt.p.s.home == nil {
			panic(fmt.Sprintf("[%v] home page is nil for %q", pt.p.s.siteVector, pt.p.Path()))
		}
		return pt.p.s.home
	}

	_, n := pt.p.s.pageMap.treePages.LongestPrefix(dir, false, func(n contentNode) bool {
		return cnh.isBranchNode(n)
	})

	if n != nil {
		return n.(page.Page)
	}

	panic(fmt.Sprintf("CurrentSection not found for %q for %s", pt.p.Path(), pt.p.s.resolveDimensionNames()))
}

func (pt pageTree) FirstSection() page.Page {
	s := pt.p.m.pathInfo.Dir()
	if s == "/" {
		return pt.p.s.home
	}

	for {
		k, n := pt.p.s.pageMap.treePages.LongestPrefix(s, false, func(n contentNode) bool { return cnh.isBranchNode(n) })
		if n == nil {
			return nil
		}

		// /blog
		if strings.Count(k, "/") < 2 {
			return n.(page.Page)
		}

		if s == "" {
			return nil
		}

		s = paths.Dir(s)

	}
}

func (pt pageTree) InSection(other any) bool {
	if pt.p == nil || types.IsNil(other) {
		return false
	}

	p, ok := other.(page.Page)
	if !ok {
		return false
	}

	return pt.CurrentSection() == p.CurrentSection()
}

func (pt pageTree) Parent() page.Page {
	if pt.p.IsHome() {
		return nil
	}

	dir := pt.p.m.pathInfo.ContainerDir()

	if dir == "" {
		return pt.p.s.home
	}

	for {
		_, n := pt.p.s.pageMap.treePages.LongestPrefix(dir, false, nil)
		if n == nil {
			return pt.p.s.home
		}
		if pt.p.m.bundled || cnh.isBranchNode(n) {
			return n.(page.Page)
		}
		dir = paths.Dir(dir)
	}
}

func (pt pageTree) Ancestors() page.Pages {
	var ancestors page.Pages
	parent := pt.Parent()
	for parent != nil {
		ancestors = append(ancestors, parent)
		parent = parent.Parent()

	}
	return ancestors
}

func (pt pageTree) Sections() page.Pages {
	var (
		pages               page.Pages
		currentBranchPrefix string
		s                   = pt.p.Path()
		prefix              = paths.AddTrailingSlash(s)
		tree                = pt.p.s.pageMap.treePages
	)

	w := &doctree.NodeShiftTreeWalker[contentNode]{
		Tree:   tree,
		Prefix: prefix,
	}
	w.Handle = func(ss string, n contentNode) (bool, error) {
		if !cnh.isBranchNode(n) {
			return false, nil
		}
		if currentBranchPrefix == "" || !strings.HasPrefix(ss, currentBranchPrefix) {
			if p, ok := n.(*pageState); ok && p.IsSection() && p.m.shouldList(false) && p.Parent() == pt.p {
				pages = append(pages, p)
			} else {
				w.SkipPrefix(ss + "/")
			}
		}
		currentBranchPrefix = ss + "/"
		return false, nil
	}

	if err := w.Walk(context.Background()); err != nil {
		panic(err)
	}

	page.SortByDefault(pages)
	return pages
}

func (pt pageTree) Page() page.Page {
	return pt.p
}

func (p pageTree) SectionsEntries() []string {
	sp := p.SectionsPath()
	if sp == "/" {
		return nil
	}
	entries := strings.Split(sp[1:], "/")
	if len(entries) == 0 {
		return nil
	}
	return entries
}

func (p pageTree) SectionsPath() string {
	return p.CurrentSection().Path()
}
