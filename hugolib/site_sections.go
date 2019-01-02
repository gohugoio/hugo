// Copyright 2017 The Hugo Authors. All rights reserved.
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
	"strconv"
	"strings"

	"github.com/gohugoio/hugo/resources/page"

	"github.com/gohugoio/hugo/helpers"

	radix "github.com/hashicorp/go-immutable-radix"
)

// Sections returns the top level sections.
func (s *SiteInfo) Sections() Pages {
	home, err := s.Home()
	if err == nil {
		return home.Sections()
	}
	return nil
}

// Home is a shortcut to the home page, equivalent to .Site.GetPage "home".
func (s *SiteInfo) Home() (*Page, error) {
	return s.GetPage(KindHome)
}

// Parent returns a section's parent section or a page's section.
// To get a section's subsections, see Page's Sections method.
func (p *Page) Parent() *Page {
	return p.parent
}

// CurrentSection returns the page's current section or the page itself if home or a section.
// Note that this will return nil for pages that is not regular, home or section pages.
func (p *Page) CurrentSection() *Page {
	v := p
	if v.origOnCopy != nil {
		v = v.origOnCopy
	}
	if v.IsHome() || v.IsSection() {
		return v
	}

	return v.parent
}

// FirstSection returns the section on level 1 below home, e.g. "/docs".
// For the home page, this will return itself.
func (p *Page) FirstSection() *Page {
	v := p
	if v.origOnCopy != nil {
		v = v.origOnCopy
	}

	if v.parent == nil || v.parent.IsHome() {
		return v
	}

	parent := v.parent
	for {
		current := parent
		parent = parent.parent
		if parent == nil || parent.IsHome() {
			return current
		}
	}

}

// InSection returns whether the given page is in the current section.
// Note that this will always return false for pages that are
// not either regular, home or section pages.
func (p *Page) InSection(other interface{}) (bool, error) {
	if p == nil || other == nil {
		return false, nil
	}

	pp, err := unwrapPage(other)
	if err != nil {
		return false, err
	}

	if pp == nil {
		return false, nil
	}

	return pp.CurrentSection() == p.CurrentSection(), nil
}

// IsDescendant returns whether the current page is a descendant of the given page.
// Note that this method is not relevant for taxonomy lists and taxonomy terms pages.
func (p *Page) IsDescendant(other interface{}) (bool, error) {
	if p == nil {
		return false, nil
	}
	pp, err := unwrapPage(other)
	if err != nil || pp == nil {
		return false, err
	}

	if pp.Kind() == KindPage && len(p.sections) == len(pp.sections) {
		// A regular page is never its section's descendant.
		return false, nil
	}
	return helpers.HasStringsPrefix(p.sections, pp.sections), nil
}

// IsAncestor returns whether the current page is an ancestor of the given page.
// Note that this method is not relevant for taxonomy lists and taxonomy terms pages.
func (p *Page) IsAncestor(other interface{}) (bool, error) {
	if p == nil {
		return false, nil
	}

	pp, err := unwrapPage(other)
	if err != nil || pp == nil {
		return false, err
	}

	if p.Kind() == KindPage && len(p.sections) == len(pp.sections) {
		// A regular page is never its section's ancestor.
		return false, nil
	}

	return helpers.HasStringsPrefix(pp.sections, p.sections), nil
}

// Eq returns whether the current page equals the given page.
// Note that this is more accurate than doing `{{ if eq $page $otherPage }}`
// since a Page can be embedded in another type.
func (p *Page) Eq(other interface{}) bool {
	pp, err := unwrapPage(other)
	if err != nil {
		return false
	}

	return p == pp
}

func unwrapPage(in interface{}) (*Page, error) {
	switch v := in.(type) {
	case *Page:
		return v, nil
	case *PageOutput:
		return v.Page, nil
	case *PageWithoutContent:
		return v.Page, nil
	case nil:
		return nil, nil
	default:
		return nil, fmt.Errorf("%T not supported", in)
	}
}

// Sections returns this section's subsections, if any.
// Note that for non-sections, this method will always return an empty list.
func (p *Page) Sections() Pages {
	return p.subSections
}

func (s *Site) assembleSections() Pages {
	var newPages Pages

	if !s.isEnabled(KindSection) {
		return newPages
	}

	// Maps section kind pages to their path, i.e. "my/section"
	sectionPages := make(map[string]page.Page)

	// The sections with content files will already have been created.
	for _, sect := range s.findPagesByKind(KindSection) {
		sectp := sect.(*Page)
		sectionPages[path.Join(sectp.sections...)] = sect

	}

	const (
		sectKey     = "__hs"
		sectSectKey = "_a" + sectKey
		sectPageKey = "_b" + sectKey
	)

	var (
		inPages    = radix.New().Txn()
		inSections = radix.New().Txn()
		undecided  Pages
	)

	home := s.findFirstPageByKindIn(KindHome, s.Pages)

	for i, p := range s.Pages {
		if p.Kind() != KindPage {
			continue
		}

		pp := p.(*Page)

		if len(pp.sections) == 0 {
			// Root level pages. These will have the home page as their Parent.
			pp.parent = home
			continue
		}

		sectionKey := path.Join(pp.sections...)
		sect, found := sectionPages[sectionKey]

		if !found && len(pp.sections) == 1 {
			// We only create content-file-less sections for the root sections.
			sect = s.newSectionPage(pp.sections[0])
			sectionPages[sectionKey] = sect
			newPages = append(newPages, sect)
			found = true
		}

		if len(pp.sections) > 1 {
			// Create the root section if not found.
			_, rootFound := sectionPages[pp.sections[0]]
			if !rootFound {
				sect = s.newSectionPage(pp.sections[0])
				sectionPages[pp.sections[0]] = sect
				newPages = append(newPages, sect)
			}
		}

		if found {
			pagePath := path.Join(sectionKey, sectPageKey, strconv.Itoa(i))
			inPages.Insert([]byte(pagePath), p)
		} else {
			undecided = append(undecided, p)
		}
	}

	// Create any missing sections in the tree.
	// A sub-section needs a content file, but to create a navigational tree,
	// given a content file in /content/a/b/c/_index.md, we cannot create just
	// the c section.
	for _, sect := range sectionPages {
		sectp := sect.(*Page)
		for i := len(sectp.sections); i > 0; i-- {
			sectionPath := sectp.sections[:i]
			sectionKey := path.Join(sectionPath...)
			_, found := sectionPages[sectionKey]
			if !found {
				sectp = s.newSectionPage(sectionPath[len(sectionPath)-1])
				sectp.sections = sectionPath
				sectionPages[sectionKey] = sectp
				newPages = append(newPages, sectp)
			}
		}
	}

	for k, sect := range sectionPages {
		inPages.Insert([]byte(path.Join(k, sectSectKey)), sect)
		inSections.Insert([]byte(k), sect)
	}

	var (
		currentSection *Page
		children       Pages
		rootSections   = inSections.Commit().Root()
	)

	for i, p := range undecided {
		pp := p.(*Page)
		// Now we can decide where to put this page into the tree.
		sectionKey := path.Join(pp.sections...)

		_, v, _ := rootSections.LongestPrefix([]byte(sectionKey))
		sect := v.(*Page)
		pagePath := path.Join(path.Join(sect.sections...), sectSectKey, "u", strconv.Itoa(i))
		inPages.Insert([]byte(pagePath), p)
	}

	var rootPages = inPages.Commit().Root()

	rootPages.Walk(func(path []byte, v interface{}) bool {
		p := v.(*Page)

		if p.Kind() == KindSection {
			if currentSection != nil {
				// A new section
				currentSection.setPagePages(children)
			}

			currentSection = p
			children = make(Pages, 0)

			return false

		}

		// Regular page
		p.parent = currentSection
		children = append(children, p)
		return false
	})

	if currentSection != nil {
		currentSection.setPagePages(children)
	}

	// Build the sections hierarchy
	for _, sect := range sectionPages {
		sectp := sect.(*Page)
		if len(sectp.sections) == 1 {
			sectp.parent = home
		} else {
			parentSearchKey := path.Join(sectp.sections[:len(sectp.sections)-1]...)
			_, v, _ := rootSections.LongestPrefix([]byte(parentSearchKey))
			p := v.(*Page)
			sectp.parent = p
		}

		if sectp.parent != nil {
			sectp.parent.subSections = append(sectp.parent.subSections, sect)
		}
	}

	var (
		sectionsParamId      = "mainSections"
		sectionsParamIdLower = strings.ToLower(sectionsParamId)
		mainSections         interface{}
		mainSectionsFound    bool
		maxSectionWeight     int
	)

	mainSections, mainSectionsFound = s.Info.Params[sectionsParamIdLower]

	for _, sect := range sectionPages {
		sectp := sect.(*Page)
		if sectp.parent != nil {
			sectp.parent.subSections.sort()
		}

		for i, p := range sectp.Pages {
			pp := p.(*Page)
			if i > 0 {
				pp.NextInSection = sectp.Pages[i-1]
			}
			if i < len(sectp.Pages)-1 {
				pp.PrevInSection = sectp.Pages[i+1]
			}
		}

		if !mainSectionsFound {
			weight := len(sectp.Pages) + (len(sectp.Sections()) * 5)
			if weight >= maxSectionWeight {
				mainSections = []string{sectp.Section()}
				maxSectionWeight = weight
			}
		}
	}

	// Try to make this as backwards compatible as possible.
	s.Info.Params[sectionsParamId] = mainSections
	s.Info.Params[sectionsParamIdLower] = mainSections

	return newPages

}

func (p *Page) setPagePages(pages Pages) {
	pages.sort()
	p.Pages = pages
	p.data = make(map[string]interface{})
	p.data["Pages"] = pages
}
