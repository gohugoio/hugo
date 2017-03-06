// Copyright 2017-present The Hugo Authors. All rights reserved.
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

package output

import (
	"fmt"
	"path"
	"strings"
)

// LayoutIdentifier is used to pick the correct layout for a piece of content.
type LayoutIdentifier interface {
	PageType() string
	PageSection() string // TODO(bep) name
	PageKind() string
	PageLayout() string
}

// Layout calculates the layout template to use to render a given output type.
// TODO(bep) output improve names
type LayoutHandler struct {
	hasTheme bool
}

func NewLayoutHandler(hasTheme bool) *LayoutHandler {
	return &LayoutHandler{hasTheme: hasTheme}
}

// TODO(bep) output theme layouts
var (
	layoutsHome        = "index.html _default/list.html"
	layoutsSection     = "section/SECTION.html  SECTION/list.html _default/section.html _default/list.html indexes/SECTION.html _default/indexes.html"
	layoutTaxonomy     = "taxonomy/SECTION.html indexes/SECTION.html _default/taxonomy.html _default/list.html"
	layoutTaxonomyTerm = "taxonomy/SECTION.terms.html _default/terms.html indexes/indexes.html"
)

func (l *LayoutHandler) For(id LayoutIdentifier, layoutOverride string, tp Type) []string {
	var layouts []string

	layout := id.PageLayout()

	if layoutOverride != "" {
		layout = layoutOverride
	}

	switch id.PageKind() {
	// TODO(bep) move the Kind constants some common place.
	case "home":
		layouts = strings.Fields(layoutsHome)
	case "section":
		layouts = strings.Fields(strings.Replace(layoutsSection, "SECTION", id.PageSection(), -1))
	case "taxonomy":
		layouts = strings.Fields(strings.Replace(layoutTaxonomy, "SECTION", id.PageSection(), -1))
	case "taxonomyTerm":
		layouts = strings.Fields(strings.Replace(layoutTaxonomyTerm, "SECTION", id.PageSection(), -1))
	case "page":
		layouts = regularPageLayouts(id.PageType(), layout)
	}

	if l.hasTheme {
		layoutsWithThemeLayouts := []string{}
		// First place all non internal templates
		for _, t := range layouts {
			if !strings.HasPrefix(t, "_internal/") {
				layoutsWithThemeLayouts = append(layoutsWithThemeLayouts, t)
			}
		}

		// Then place theme templates with the same names
		for _, t := range layouts {
			if !strings.HasPrefix(t, "_internal/") {
				layoutsWithThemeLayouts = append(layoutsWithThemeLayouts, "theme/"+t)
			}
		}

		// Lastly place internal templates
		for _, t := range layouts {
			if strings.HasPrefix(t, "_internal/") {
				layoutsWithThemeLayouts = append(layoutsWithThemeLayouts, t)
			}
		}

		return layoutsWithThemeLayouts
	}

	return layouts
}

func regularPageLayouts(types string, layout string) (layouts []string) {
	if layout == "" {
		layout = "single"
	}

	if types != "" {
		t := strings.Split(types, "/")

		// Add type/layout.html
		for i := range t {
			search := t[:len(t)-i]
			layouts = append(layouts, fmt.Sprintf("%s/%s.html", strings.ToLower(path.Join(search...)), layout))
		}
	}

	// Add _default/layout.html
	layouts = append(layouts, fmt.Sprintf("_default/%s.html", layout))

	return
}
