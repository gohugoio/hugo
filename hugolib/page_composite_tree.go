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
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/resources/page"
)

type pageTreeProvider struct {
	p *pageState
}

func (pt pageTreeProvider) CurrentSection() page.Page {
	p := pt.p

	if p.IsHome() || p.IsSection() {
		return p
	}

	return p.Parent()
}

func (pt pageTreeProvider) FirstSection() page.Page {
	p := pt.p

	parent := p.Parent()

	if parent == nil || parent.IsHome() {
		return p
	}

	for {
		current := parent
		parent = parent.Parent()
		if parent == nil || parent.IsHome() {
			return current
		}
	}

}

func (pt pageTreeProvider) InSection(other interface{}) (bool, error) {
	if pt.p == nil || other == nil {
		return false, nil
	}

	pp, err := unwrapPage(other)
	if err != nil {
		return false, err
	}

	if pp == nil {
		return false, nil
	}

	return pp.CurrentSection().Eq(pt.p.CurrentSection()), nil

}

func (pt pageTreeProvider) IsAncestor(other interface{}) (bool, error) {
	if pt.p == nil {
		return false, nil
	}

	pp, err := unwrapPage(other)
	if err != nil || pp == nil {
		return false, err
	}

	if pt.p.Kind() == page.KindPage && len(pt.p.SectionsEntries()) == len(pp.SectionsEntries()) {
		// A regular page is never its section's ancestor.
		return false, nil
	}

	return helpers.HasStringsPrefix(pp.SectionsEntries(), pt.p.SectionsEntries()), nil
}

func (pt pageTreeProvider) IsDescendant(other interface{}) (bool, error) {
	if pt.p == nil {
		return false, nil
	}
	pp, err := unwrapPage(other)
	if err != nil || pp == nil {
		return false, err
	}

	if pp.Kind() == page.KindPage && len(pt.p.SectionsEntries()) == len(pp.SectionsEntries()) {
		// A regular page is never its section's descendant.
		return false, nil
	}
	return helpers.HasStringsPrefix(pt.p.SectionsEntries(), pp.SectionsEntries()), nil
}

// TODO(bep) page check
func (pt pageTreeProvider) Parent() page.Page {
	return pt.p.parent
}

func (pt pageTreeProvider) Sections() page.Pages {
	return pt.p.subSections
}
