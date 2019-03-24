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
	"strconv"
	"strings"

	"github.com/gohugoio/hugo/resources/page"
	"github.com/gohugoio/hugo/resources/resource"

	radix "github.com/hashicorp/go-immutable-radix"
)

// Sections returns the top level sections.
func (s *SiteInfo) Sections() page.Pages {
	home, err := s.Home()
	if err == nil {
		return home.Sections()
	}
	return nil
}

// Home is a shortcut to the home page, equivalent to .Site.GetPage "home".
func (s *SiteInfo) Home() (page.Page, error) {
	return s.s.home, nil
}

func (s *Site) assembleSections() pageStatePages {
	var newPages pageStatePages

	if !s.isEnabled(page.KindSection) {
		return newPages
	}

	// Maps section kind pages to their path, i.e. "my/section"
	sectionPages := make(map[string]*pageState)

	// The sections with content files will already have been created.
	for _, sect := range s.findWorkPagesByKind(page.KindSection) {
		sectionPages[sect.SectionsPath()] = sect
	}

	const (
		sectKey     = "__hs"
		sectSectKey = "_a" + sectKey
		sectPageKey = "_b" + sectKey
	)

	var (
		inPages    = radix.New().Txn()
		inSections = radix.New().Txn()
		undecided  pageStatePages
	)

	home := s.findFirstWorkPageByKindIn(page.KindHome)

	for i, p := range s.workAllPages {

		if p.Kind() != page.KindPage {
			continue
		}

		sections := p.SectionsEntries()

		if len(sections) == 0 {
			// Root level pages. These will have the home page as their Parent.
			p.parent = home
			continue
		}

		sectionKey := p.SectionsPath()
		_, found := sectionPages[sectionKey]

		if !found && len(sections) == 1 {

			// We only create content-file-less sections for the root sections.
			n := s.newPage(page.KindSection, sections[0])

			sectionPages[sectionKey] = n
			newPages = append(newPages, n)
			found = true
		}

		if len(sections) > 1 {
			// Create the root section if not found.
			_, rootFound := sectionPages[sections[0]]
			if !rootFound {
				sect := s.newPage(page.KindSection, sections[0])
				sectionPages[sections[0]] = sect
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
		sections := sect.SectionsEntries()
		for i := len(sections); i > 0; i-- {
			sectionPath := sections[:i]
			sectionKey := path.Join(sectionPath...)
			_, found := sectionPages[sectionKey]
			if !found {
				sect = s.newPage(page.KindSection, sectionPath[len(sectionPath)-1])
				sect.m.sections = sectionPath
				sectionPages[sectionKey] = sect
				newPages = append(newPages, sect)
			}
		}
	}

	for k, sect := range sectionPages {
		inPages.Insert([]byte(path.Join(k, sectSectKey)), sect)
		inSections.Insert([]byte(k), sect)
	}

	var (
		currentSection *pageState
		children       page.Pages
		dates          *resource.Dates
		rootSections   = inSections.Commit().Root()
	)

	for i, p := range undecided {
		// Now we can decide where to put this page into the tree.
		sectionKey := p.SectionsPath()

		_, v, _ := rootSections.LongestPrefix([]byte(sectionKey))
		sect := v.(*pageState)
		pagePath := path.Join(path.Join(sect.SectionsEntries()...), sectSectKey, "u", strconv.Itoa(i))
		inPages.Insert([]byte(pagePath), p)
	}

	var rootPages = inPages.Commit().Root()

	rootPages.Walk(func(path []byte, v interface{}) bool {
		p := v.(*pageState)

		if p.Kind() == page.KindSection {
			if currentSection != nil {
				// A new section
				currentSection.setPages(children)
			}

			currentSection = p
			children = make(page.Pages, 0)
			dates = &resource.Dates{}

			return false

		}

		// Regular page
		p.parent = currentSection
		children = append(children, p)
		dates.UpdateDateAndLastmodIfAfter(p)
		return false
	})

	if currentSection != nil {
		currentSection.setPages(children)
		currentSection.m.Dates = *dates

	}

	// Build the sections hierarchy
	for _, sect := range sectionPages {
		sections := sect.SectionsEntries()
		if len(sections) == 1 {
			if home != nil {
				sect.parent = home
			}
		} else {
			parentSearchKey := path.Join(sect.SectionsEntries()[:len(sections)-1]...)
			_, v, _ := rootSections.LongestPrefix([]byte(parentSearchKey))
			p := v.(*pageState)
			sect.parent = p
		}

		sect.addSectionToParent()
	}

	var (
		sectionsParamId      = "mainSections"
		sectionsParamIdLower = strings.ToLower(sectionsParamId)
		mainSections         interface{}
		mainSectionsFound    bool
		maxSectionWeight     int
	)

	mainSections, mainSectionsFound = s.Info.Params()[sectionsParamIdLower]

	for _, sect := range sectionPages {
		sect.sortParentSections()

		if !mainSectionsFound {
			weight := len(sect.Pages()) + (len(sect.Sections()) * 5)
			if weight >= maxSectionWeight {
				mainSections = []string{sect.Section()}
				maxSectionWeight = weight
			}
		}
	}

	// Try to make this as backwards compatible as possible.
	s.Info.Params()[sectionsParamId] = mainSections
	s.Info.Params()[sectionsParamIdLower] = mainSections

	return newPages

}
