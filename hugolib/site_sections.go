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

	radix "github.com/hashicorp/go-immutable-radix"
	"github.com/spf13/hugo/helpers"
)

// Deprecated: Use .Site.Home.Sections.
// TODO(bep) Hugo 0.23 = Reuse as an alias for home's sections.
func (s *SiteInfo) Sections() Taxonomy {

	helpText := `In Hugo 0.22 we introduced nested sections, so this method now returns an empty taxonomy.
	
To list sections with its pages, you can do something like this:

{{ range $.Site.Home.Sections }}
  Section: {{ .Title }}
  {{ range .Pages }}
     Section Page: {{ .Title }}
  {{ end }}
{{ end }}

To get a specific section, you can do this:

{{ $section := $.Site.GetPage "section" "blog" }}
`

	helpers.Deprecated("Site", "Sections", helpText, true)

	return Taxonomy{}
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

// current returns the page's current section.
// Note that this will return nil for pages that is not regular, home or section pages.
// Note that for paginated sections and home pages, this will return the original page pointer.
func (p *Page) current() *Page {
	v := p
	if v.origOnCopy != nil {
		v = v.origOnCopy
	}
	if v.IsHome() || v.IsSection() {
		return v
	}

	return v.parent
}

// InSection returns whether the given page is in the current section.
// Note that this will always return false for pages that are
// not either regular, home or section pages.
func (p *Page) InSection(other interface{}) (bool, error) {
	if p == nil || other == nil {
		return false, nil
	}

	if po, ok := other.(*PageOutput); ok {
		other = po.Page
	}

	pp, ok := other.(*Page)
	if !ok {
		return false, fmt.Errorf("%T not supported in InSection", other)
	}

	if pp == nil {
		return false, nil
	}

	return pp.current() == p.current(), nil
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
	sectionPages := make(map[string]*Page)

	// The sections with content files will already have been created.
	for _, sect := range s.findPagesByKind(KindSection) {
		sectionPages[path.Join(sect.sections...)] = sect
	}

	const (
		sectKey     = "__hs"
		sectSectKey = "_a" + sectKey
		sectPageKey = "_b" + sectKey
	)

	var (
		home       *Page
		inPages    = radix.New().Txn()
		inSections = radix.New().Txn()
		undecided  Pages
	)

	for i, p := range s.Pages {
		if p.Kind != KindPage {
			if p.Kind == KindHome {
				home = p
			}
			continue
		}

		if len(p.sections) == 0 {
			// Root level pages. These will have the home page as their Parent.
			p.parent = home
			continue
		}

		sectionKey := path.Join(p.sections...)
		sect, found := sectionPages[sectionKey]

		if !found && len(p.sections) == 1 {
			// We only create content-file-less sections for the root sections.
			sect = s.newSectionPage(p.sections[0])
			sectionPages[sectionKey] = sect
			newPages = append(newPages, sect)
		} else if !found {
			// We don't know what to do with this section yet.
			undecided = append(undecided, p)
		}

		pagePath := path.Join(sectionKey, sectPageKey, strconv.Itoa(i))
		inPages.Insert([]byte(pagePath), p)
	}

	// Create any missing sections in the tree.
	// A sub-section needs a content file, but to create a navigational tree,
	// given a content file in /content/a/b/c/_index.md, we cannot create just
	// the c section.
	for _, sect := range sectionPages {
		for i := len(sect.sections); i > 0; i-- {
			sectionPath := sect.sections[:i]
			sectionKey := path.Join(sectionPath...)
			sect, found := sectionPages[sectionKey]
			if !found {
				sect = s.newSectionPage(sectionPath[len(sectionPath)-1])
				sect.sections = sectionPath
				sectionPages[sectionKey] = sect
				newPages = append(newPages, sect)
			}
		}
	}

	// Create any missing root sections.
	for _, p := range undecided {
		sectionKey := p.sections[0]
		sect, found := sectionPages[sectionKey]
		if !found {
			sect = s.newSectionPage(sectionKey)
			sectionPages[sectionKey] = sect
			newPages = append(newPages, sect)
		}
	}

	for k, sect := range sectionPages {
		inPages.Insert([]byte(path.Join(k, sectSectKey)), sect)
		inSections.Insert([]byte(k), sect)
	}

	var (
		currentSection *Page
		children       Pages
		rootPages      = inPages.Commit().Root()
		rootSections   = inSections.Commit().Root()
	)

	rootPages.Walk(func(path []byte, v interface{}) bool {
		p := v.(*Page)

		if p.Kind == KindSection {
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
		if len(sect.sections) == 1 {
			sect.parent = home
		} else {
			parentSearchKey := path.Join(sect.sections[:len(sect.sections)-1]...)
			_, v, _ := rootSections.LongestPrefix([]byte(parentSearchKey))
			p := v.(*Page)
			sect.parent = p
		}

		if sect.parent != nil {
			sect.parent.subSections = append(sect.parent.subSections, sect)
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
		if sect.parent != nil {
			sect.parent.subSections.Sort()
		}

		for i, p := range sect.Pages {
			if i > 0 {
				p.NextInSection = sect.Pages[i-1]
			}
			if i < len(sect.Pages)-1 {
				p.PrevInSection = sect.Pages[i+1]
			}
		}

		if !mainSectionsFound {
			weight := len(sect.Pages) + (len(sect.Sections()) * 5)
			if weight >= maxSectionWeight {
				mainSections = []string{sect.Section()}
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
	pages.Sort()
	p.Pages = pages
	p.Data = make(map[string]interface{})
	p.Data["Pages"] = pages
}
