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

// LayoutDescriptor describes how a layout should be chosen. This is
// typically built from a Page.
type LayoutDescriptor struct {
	Type    string
	Section string
	Kind    string
	Layout  string
}

// Layout calculates the layout template to use to render a given output type.
// TODO(bep) output improve names
type LayoutHandler struct {
	hasTheme bool
}

func NewLayoutHandler(hasTheme bool) *LayoutHandler {
	return &LayoutHandler{hasTheme: hasTheme}
}

const (
	layoutsHome    = "index.NAME.SUFFIX index.SUFFIX _default/list.NAME.SUFFIX _default/list.SUFFIX"
	layoutsSection = `
section/SECTION.NAME.SUFFIX section/SECTION.SUFFIX
SECTION/list.NAME.SUFFIX SECTION/list.SUFFIX
_default/section.NAME.SUFFIX _default/section.SUFFIX
_default/list.NAME.SUFFIX _default/list.SUFFIX
indexes/SECTION.NAME.SUFFIX indexes/SECTION.SUFFIX
_default/indexes.NAME.SUFFIX _default/indexes.SUFFIX
`
	layoutTaxonomy = `
taxonomy/SECTION.NAME.SUFFIX taxonomy/SECTION.SUFFIX
indexes/SECTION.NAME.SUFFIX indexes/SECTION.SUFFIX 
_default/taxonomy.NAME.SUFFIX _default/taxonomy.SUFFIX
_default/list.NAME.SUFFIX _default/list.SUFFIX
`
	layoutTaxonomyTerm = `
taxonomy/SECTION.terms.NAME.SUFFIX taxonomy/SECTION.terms.SUFFIX
_default/terms.NAME.SUFFIX _default/terms.SUFFIX
indexes/indexes.NAME.SUFFIX indexes/indexes.SUFFIX
`
)

func (l *LayoutHandler) For(d LayoutDescriptor, layoutOverride string, f Format) []string {
	var layouts []string

	layout := d.Layout

	if layoutOverride != "" {
		layout = layoutOverride
	}

	switch d.Kind {
	// TODO(bep) move the Kind constants some common place.
	case "home":
		layouts = resolveTemplate(layoutsHome, d, f)
	case "section":
		layouts = resolveTemplate(layoutsSection, d, f)
	case "taxonomy":
		layouts = resolveTemplate(layoutTaxonomy, d, f)
	case "taxonomyTerm":
		layouts = resolveTemplate(layoutTaxonomyTerm, d, f)
	case "page":
		layouts = regularPageLayouts(d.Type, layout, f)
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

func resolveTemplate(templ string, d LayoutDescriptor, f Format) []string {
	return strings.Fields(replaceKeyValues(templ,
		"SUFFIX", f.MediaType.Suffix,
		"NAME", strings.ToLower(f.Name),
		"SECTION", d.Section))
}

func replaceKeyValues(s string, oldNew ...string) string {
	replacer := strings.NewReplacer(oldNew...)
	return replacer.Replace(s)
}

func regularPageLayouts(types string, layout string, f Format) (layouts []string) {
	if layout == "" {
		layout = "single"
	}

	suffix := f.MediaType.Suffix
	name := strings.ToLower(f.Name)

	if types != "" {
		t := strings.Split(types, "/")

		// Add type/layout.html
		for i := range t {
			search := t[:len(t)-i]
			layouts = append(layouts, fmt.Sprintf("%s/%s.%s.%s", strings.ToLower(path.Join(search...)), layout, name, suffix))
			layouts = append(layouts, fmt.Sprintf("%s/%s.%s", strings.ToLower(path.Join(search...)), layout, suffix))

		}
	}

	// Add _default/layout.html
	layouts = append(layouts, fmt.Sprintf("_default/%s.%s.%s", layout, name, suffix))
	layouts = append(layouts, fmt.Sprintf("_default/%s.%s", layout, suffix))

	return
}
