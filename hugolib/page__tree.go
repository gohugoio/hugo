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
	"path"
	"strings"

	"github.com/gohugoio/hugo/resources/page"
)

// pageTree holds the treen navigational method for a Page.
type pageTree struct {
	p *pageState
}

func (pt pageTree) IsAncestor(other interface{}) (bool, error) {
	tp, ok := other.(treeRefProvider)
	if !ok {
		return false, nil
	}

	ref1, ref2 := pt.p.getTreeRef(), tp.getTreeRef()

	if ref1.Key() == "" {
		return true, nil
	}

	if ref1.Key() == ref2.Key() {
		return true, nil
	}

	return strings.HasPrefix(ref2.Key(), ref1.Key()+"/"), nil
}

// 2 TODO1 create issue: CurrentSection should navigate sideways for all branch nodes.
func (pt pageTree) CurrentSection() page.Page {
	return pt.p.m.treeRef.GetBranch().n.p
}

func (pt pageTree) IsDescendant(other interface{}) (bool, error) {
	tp, ok := other.(treeRefProvider)
	if !ok {
		return false, nil
	}

	ref1, ref2 := pt.p.getTreeRef(), tp.getTreeRef()

	if ref2.Key() == "" {
		return true, nil
	}

	if ref1.Key() == ref2.Key() {
		return true, nil
	}

	return strings.HasPrefix(ref1.Key(), ref2.Key()+"/"), nil
}

func (pt pageTree) FirstSection() page.Page {
	ref := pt.p.getTreeRef()
	key := ref.Key()
	n := ref.GetNode()
	branch := ref.GetBranch()

	if branch.n != n {
		key = path.Dir(key)
	}
	_, b := pt.p.s.pageMap.GetFirstSection(key)
	if b == nil {
		return nil
	}
	return b.p
}

func (pt pageTree) InSection(other interface{}) (bool, error) {
	tp, ok := other.(treeRefProvider)
	if !ok {
		return false, nil
	}

	ref1, ref2 := pt.p.getTreeRef(), tp.getTreeRef()

	return ref1.GetBranch() == ref2.GetBranch(), nil
}

func (pt pageTree) Parent() page.Page {
	owner := pt.p.getTreeRef().GetContainerNode()
	if owner == nil {
		return nil
	}
	return owner.p
}

func (pt pageTree) Sections() page.Pages {
	return pt.p.bucket.getSections()
}

func (pt pageTree) Page() page.Page {
	return pt.p
}
