// Copyright 2024 The Hugo Authors. All rights reserved.
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
	"path/filepath"
	"strings"

	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/hugofs/files"

	"github.com/gohugoio/hugo/common/paths"

	"github.com/gohugoio/hugo/resources/kinds"
	"github.com/gohugoio/hugo/resources/page"
)

// pageFinder provides ways to find a Page in a Site.
type pageFinder struct {
	pageMap *pageMap
}

func newPageFinder(m *pageMap) *pageFinder {
	if m == nil {
		panic("must provide a pageMap")
	}
	c := &pageFinder{pageMap: m}
	return c
}

// getPageRef resolves a Page from ref/relRef, with a slightly more comprehensive
// search path than getPage.
func (c *pageFinder) getPageRef(context page.Page, ref string) (page.Page, error) {
	n, err := c.getContentNode(context, true, ref)
	if err != nil {
		return nil, err
	}

	if p, ok := n.(page.Page); ok {
		return p, nil
	}
	return nil, nil
}

func (c *pageFinder) getPage(context page.Page, ref string) (page.Page, error) {
	n, err := c.getContentNode(context, false, ref)
	if err != nil {
		return nil, err
	}
	if p, ok := n.(page.Page); ok {
		return p, nil
	}
	return nil, nil
}

// Only used in tests.
func (c *pageFinder) getPageOldVersion(kind string, sections ...string) page.Page {
	refs := append([]string{kind}, path.Join(sections...))
	p, _ := c.getPageForRefs(refs...)
	return p
}

// This is an adapter func for the old API with Kind as first argument.
// This is invoked when you do .Site.GetPage. We drop the Kind and fails
// if there are more than 2 arguments, which would be ambiguous.
func (c *pageFinder) getPageForRefs(ref ...string) (page.Page, error) {
	var refs []string
	for _, r := range ref {
		// A common construct in the wild is
		// .Site.GetPage "home" "" or
		// .Site.GetPage "home" "/"
		if r != "" && r != "/" {
			refs = append(refs, r)
		}
	}

	var key string

	if len(refs) > 2 {
		// This was allowed in Hugo <= 0.44, but we cannot support this with the
		// new API. This should be the most unusual case.
		return nil, fmt.Errorf(`too many arguments to .Site.GetPage: %v. Use lookups on the form {{ .Site.GetPage "/posts/mypage-md" }}`, ref)
	}

	if len(refs) == 0 || refs[0] == kinds.KindHome {
		key = "/"
	} else if len(refs) == 1 {
		if len(ref) == 2 && refs[0] == kinds.KindSection {
			// This is an old style reference to the "Home Page section".
			// Typically fetched via {{ .Site.GetPage "section" .Section }}
			// See https://github.com/gohugoio/hugo/issues/4989
			key = "/"
		} else {
			key = refs[0]
		}
	} else {
		key = refs[1]
	}

	return c.getPage(nil, key)
}

const defaultContentExt = ".md"

func (c *pageFinder) getContentNode(context page.Page, isReflink bool, ref string) (contentNodeI, error) {
	ref = paths.ToSlashTrimTrailing(ref)
	inRef := ref
	if ref == "" {
		ref = "/"
	}

	if paths.HasExt(ref) {
		return c.getContentNodeForRef(context, isReflink, true, inRef, ref)
	}

	// We are always looking for a content file and having an extension greatly simplifies the code that follows,
	// even in the case where the extension does not match this one.
	if ref == "/" {
		if n, err := c.getContentNodeForRef(context, isReflink, false, inRef, "/_index"+defaultContentExt); n != nil || err != nil {
			return n, err
		}
	} else if strings.HasSuffix(ref, "/index") {
		if n, err := c.getContentNodeForRef(context, isReflink, false, inRef, ref+"/index"+defaultContentExt); n != nil || err != nil {
			return n, err
		}
		if n, err := c.getContentNodeForRef(context, isReflink, false, inRef, ref+defaultContentExt); n != nil || err != nil {
			return n, err
		}
	} else {
		if n, err := c.getContentNodeForRef(context, isReflink, false, inRef, ref+defaultContentExt); n != nil || err != nil {
			return n, err
		}
	}

	return nil, nil
}

func (c *pageFinder) getContentNodeForRef(context page.Page, isReflink, hadExtension bool, inRef, ref string) (contentNodeI, error) {
	s := c.pageMap.s
	contentPathParser := s.Conf.PathParser()

	if context != nil && !strings.HasPrefix(ref, "/") {
		// Try the page-relative path first.
		// Branch pages: /mysection, "./mypage" => /mysection/mypage
		// Regular pages: /mysection/mypage.md, Path=/mysection/mypage, "./someotherpage" => /mysection/mypage/../someotherpage
		// Regular leaf bundles: /mysection/mypage/index.md, Path=/mysection/mypage, "./someotherpage" => /mysection/mypage/../someotherpage
		// Given the above, for regular pages we use the containing folder.
		var baseDir string
		if pi := context.PathInfo(); pi != nil {
			if pi.IsBranchBundle() || (hadExtension && strings.HasPrefix(ref, "../")) {
				baseDir = pi.Dir()
			} else {
				baseDir = pi.ContainerDir()
			}
		}

		rel := path.Join(baseDir, ref)

		relPath, _ := contentPathParser.ParseBaseAndBaseNameNoIdentifier(files.ComponentFolderContent, rel)

		n, err := c.getContentNodeFromPath(relPath, ref)
		if n != nil || err != nil {
			return n, err
		}

		if hadExtension && context.File() != nil {
			if n, err := c.getContentNodeFromRefReverseLookup(inRef, context.File().FileInfo()); n != nil || err != nil {
				return n, err
			}
		}

	}

	if strings.HasPrefix(ref, ".") {
		// Page relative, no need to look further.
		return nil, nil
	}

	relPath, nameNoIdentifier := contentPathParser.ParseBaseAndBaseNameNoIdentifier(files.ComponentFolderContent, ref)

	n, err := c.getContentNodeFromPath(relPath, ref)

	if n != nil || err != nil {
		return n, err
	}

	if hadExtension && s.home != nil && s.home.File() != nil {
		if n, err := c.getContentNodeFromRefReverseLookup(inRef, s.home.File().FileInfo()); n != nil || err != nil {
			return n, err
		}
	}

	var doSimpleLookup bool
	if isReflink || context == nil {
		slashCount := strings.Count(inRef, "/")
		doSimpleLookup = slashCount == 0
	}

	if !doSimpleLookup {
		return nil, nil
	}

	n = c.pageMap.pageReverseIndex.Get(nameNoIdentifier)
	if n == ambiguousContentNode {
		return nil, fmt.Errorf("page reference %q is ambiguous", inRef)
	}

	return n, nil
}

func (c *pageFinder) getContentNodeFromRefReverseLookup(ref string, fi hugofs.FileMetaInfo) (contentNodeI, error) {
	s := c.pageMap.s
	meta := fi.Meta()
	dir := meta.Filename
	if !fi.IsDir() {
		dir = filepath.Dir(meta.Filename)
	}

	realFilename := filepath.Join(dir, ref)

	pcs, err := s.BaseFs.Content.ReverseLookup(realFilename, true)
	if err != nil {
		return nil, err
	}

	// There may be multiple matches, but we will only use the first one.
	for _, pc := range pcs {
		pi := s.Conf.PathParser().Parse(pc.Component, pc.Path)
		if n := c.pageMap.treePages.Get(pi.Base()); n != nil {
			return n, nil
		}
	}
	return nil, nil
}

func (c *pageFinder) getContentNodeFromPath(s string, ref string) (contentNodeI, error) {
	n := c.pageMap.treePages.Get(s)
	if n != nil {
		return n, nil
	}

	return nil, nil
}
