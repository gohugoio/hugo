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

	"github.com/gohugoio/hugo/common/types"
	"github.com/gohugoio/hugo/resources/page"
)

type pageTree struct {
	p *pageState
}

func (pt pageTree) IsAncestor(other interface{}) (bool, error) {
	if pt.p == nil {
		return false, nil
	}

	tp, ok := other.(treeRefProvider)
	if !ok {
		return false, nil
	}

	ref1, ref2 := pt.p.getTreeRef(), tp.getTreeRef()

	if ref1 != nil && ref1.key == "/" {
		return true, nil
	}

	if ref1 == nil || ref2 == nil {
		if ref1 == nil {
			// A 404 or other similar standalone page.
			return false, nil
		}

		return ref1.n.p.IsHome(), nil
	}

	if !ref1.isSection() {
		return false, nil
	}

	if ref2.isSection() {
		return strings.HasPrefix(ref2.key, ref1.key+"/"), nil
	}

	return strings.HasPrefix(ref2.key, ref1.key+cmBranchSeparator), nil

}

func (pt pageTree) CurrentSection() page.Page {
	p := pt.p

	if p.IsHome() || p.IsSection() {
		return p
	}

	return p.Parent()
}

func (pt pageTree) IsDescendant(other interface{}) (bool, error) {
	if pt.p == nil {
		return false, nil
	}

	tp, ok := other.(treeRefProvider)
	if !ok {
		return false, nil
	}

	ref1, ref2 := pt.p.getTreeRef(), tp.getTreeRef()

	if ref2 != nil && ref2.key == "/" {
		return true, nil
	}

	if ref1 == nil || ref2 == nil {
		if ref2 == nil {
			// A 404 or other similar standalone page.
			return false, nil
		}

		return ref2.n.p.IsHome(), nil
	}

	if !ref2.isSection() {
		return false, nil
	}

	if ref1.isSection() {
		return strings.HasPrefix(ref1.key, ref2.key+"/"), nil
	}

	return strings.HasPrefix(ref1.key, ref2.key+cmBranchSeparator), nil

}

func (pt pageTree) FirstSection() page.Page {
	ref := pt.p.getTreeRef()
	if ref == nil {
		return pt.p.s.home
	}
	key := ref.key
	if !ref.isSection() {
		key = path.Dir(key)
	}
	_, b := ref.m.getFirstSection(key)
	if b == nil {
		return nil
	}
	return b.p
}

func (pt pageTree) InSection(other interface{}) (bool, error) {
	if pt.p == nil || types.IsNil(other) {
		return false, nil
	}

	tp, ok := other.(treeRefProvider)
	if !ok {
		return false, nil
	}

	ref1, ref2 := pt.p.getTreeRef(), tp.getTreeRef()

	if ref1 == nil || ref2 == nil {
		if ref1 == nil {
			// A 404 or other similar standalone page.
			return false, nil
		}
		return ref1.n.p.IsHome(), nil
	}

	s1, _ := ref1.getCurrentSection()
	s2, _ := ref2.getCurrentSection()

	return s1 == s2, nil

}

func (pt pageTree) Page() page.Page {
	return pt.p
}

func (pt pageTree) Parent() page.Page {
	p := pt.p

	if p.parent != nil {
		return p.parent
	}

	if pt.p.IsHome() {
		return nil
	}

	tree := p.getTreeRef()

	if tree == nil || pt.p.Kind() == page.KindTaxonomyTerm {
		return pt.p.s.home
	}

	_, b := tree.getSection()
	if b == nil {
		return nil
	}

	return b.p
}

func (pt pageTree) Sections() page.Pages {
	if pt.p.bucket == nil {
		return nil
	}

	return pt.p.bucket.getSections()
}
