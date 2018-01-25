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
	"strings"
	"sync"

	"github.com/gohugoio/hugo/helpers"
)

// These may be used as content sections with potential conflicts. Avoid that.
var reservedSections = map[string]bool{
	"shortcodes": true,
	"partials":   true,
}

// LayoutDescriptor describes how a layout should be chosen. This is
// typically built from a Page.
type LayoutDescriptor struct {
	Type    string
	Section string
	Kind    string
	Lang    string
	Layout  string
	// LayoutOverride indicates what we should only look for the above layout.
	LayoutOverride bool
}

// LayoutHandler calculates the layout template to use to render a given output type.
type LayoutHandler struct {
	hasTheme bool

	mu    sync.RWMutex
	cache map[layoutCacheKey][]string
}

type layoutCacheKey struct {
	d LayoutDescriptor
	f Format
}

// NewLayoutHandler creates a new LayoutHandler.
func NewLayoutHandler(hasTheme bool) *LayoutHandler {
	return &LayoutHandler{hasTheme: hasTheme, cache: make(map[layoutCacheKey][]string)}
}

// For returns a layout for the given LayoutDescriptor and options.
// Layouts are rendered and cached internally.
func (l *LayoutHandler) For(d LayoutDescriptor, f Format) ([]string, error) {

	// We will get lots of requests for the same layouts, so avoid recalculations.
	key := layoutCacheKey{d, f}
	l.mu.RLock()
	if cacheVal, found := l.cache[key]; found {
		l.mu.RUnlock()
		return cacheVal, nil
	}
	l.mu.RUnlock()

	layouts := resolvePageTemplate(d, f)

	if l.hasTheme {
		// From Hugo 0.33 we interleave the project/theme templates. This was kind of a fundamental change, but the
		// previous behaviour was surprising.
		// As an example, an `index.html` in theme for the home page will now win over a `_default/list.html` in the project.
		layoutsWithThemeLayouts := []string{}

		// First place all non internal templates
		for _, t := range layouts {
			if !strings.HasPrefix(t, "_internal/") {
				layoutsWithThemeLayouts = append(layoutsWithThemeLayouts, t)
				layoutsWithThemeLayouts = append(layoutsWithThemeLayouts, "theme/"+t)
			}
		}

		// Lastly place internal templates
		for _, t := range layouts {
			if strings.HasPrefix(t, "_internal/") {
				layoutsWithThemeLayouts = append(layoutsWithThemeLayouts, t)
			}
		}

		layouts = layoutsWithThemeLayouts
	}

	layouts = prependTextPrefixIfNeeded(f, layouts...)
	layouts = helpers.UniqueStrings(layouts)

	l.mu.Lock()
	l.cache[key] = layouts
	l.mu.Unlock()

	return layouts, nil
}

type layoutBuilder struct {
	layoutVariations []string
	typeVariations   []string
	d                LayoutDescriptor
	f                Format
}

func (l *layoutBuilder) addLayoutVariations(vars ...string) {
	for _, layoutVar := range vars {
		if l.d.LayoutOverride && layoutVar != l.d.Layout {
			continue
		}
		l.layoutVariations = append(l.layoutVariations, layoutVar)
	}
}

func (l *layoutBuilder) addTypeVariations(vars ...string) {
	for _, typeVar := range vars {
		if !reservedSections[typeVar] {
			l.typeVariations = append(l.typeVariations, typeVar)
		}
	}
}

func (l *layoutBuilder) addSectionType() {
	if l.d.Section != "" {
		l.addTypeVariations(l.d.Section)
	}
}

func (l *layoutBuilder) addKind() {
	l.addLayoutVariations(l.d.Kind)
	l.addTypeVariations(l.d.Kind)
}

func resolvePageTemplate(d LayoutDescriptor, f Format) []string {

	b := &layoutBuilder{d: d, f: f}

	if d.Layout != "" {
		b.addLayoutVariations(d.Layout)
	}

	if d.Type != "" {
		b.addTypeVariations(d.Type)
	}

	switch d.Kind {
	case "page":
		b.addLayoutVariations("single")
		b.addSectionType()
	case "home":
		b.addLayoutVariations("index", "home")
		// Also look in the root
		b.addTypeVariations("")
	case "section":
		if d.Section != "" {
			b.addLayoutVariations(d.Section)
		}
		b.addSectionType()
		b.addKind()
	case "taxonomy":
		if d.Section != "" {
			b.addLayoutVariations(d.Section)
		}
		b.addKind()
		b.addSectionType()

	case "taxonomyTerm":
		if d.Section != "" {
			b.addLayoutVariations(d.Section + ".terms")
		}
		b.addTypeVariations("taxonomy")
		b.addSectionType()
		b.addLayoutVariations("terms")

	}

	isRSS := f.Name == RSSFormat.Name
	if isRSS {
		// The historic and common rss.xml case
		b.addLayoutVariations("")
	}

	// All have _default in their lookup path
	b.addTypeVariations("_default")

	if d.Kind != "page" {
		// Add the common list type
		b.addLayoutVariations("list")
	}

	layouts := b.resolveVariations()

	if isRSS {
		layouts = append(layouts, "_internal/_default/rss.xml")
	}

	return layouts

}

func (l *layoutBuilder) resolveVariations() []string {

	var layouts []string

	var variations []string
	name := strings.ToLower(l.f.Name)

	if l.d.Lang != "" {
		// We prefer the most specific type before language.
		variations = append(variations, []string{fmt.Sprintf("%s.%s", l.d.Lang, name), name, l.d.Lang}...)
	} else {
		variations = append(variations, name)
	}

	variations = append(variations, "")

	for _, typeVar := range l.typeVariations {
		for _, variation := range variations {
			for _, layoutVar := range l.layoutVariations {
				if variation == "" && layoutVar == "" {
					continue
				}
				template := layoutTemplate(typeVar, layoutVar)
				layouts = append(layouts, replaceKeyValues(template,
					"TYPE", typeVar,
					"LAYOUT", layoutVar,
					"VARIATIONS", variation,
					"EXTENSION", l.f.MediaType.Suffix,
				))
			}
		}

	}

	return filterDotLess(layouts)
}

func layoutTemplate(typeVar, layoutVar string) string {

	var l string

	if typeVar != "" {
		l = "TYPE/"
	}

	if layoutVar != "" {
		l += "LAYOUT.VARIATIONS.EXTENSION"
	} else {
		l += "VARIATIONS.EXTENSION"
	}

	return l
}

func filterDotLess(layouts []string) []string {
	var filteredLayouts []string

	for _, l := range layouts {
		l = strings.Replace(l, "..", ".", -1)
		l = strings.Trim(l, ".")
		// If media type has no suffix, we have "index" type of layouts in this list, which
		// doesn't make much sense.
		if strings.Contains(l, ".") {
			filteredLayouts = append(filteredLayouts, l)
		}
	}

	return filteredLayouts
}

func prependTextPrefixIfNeeded(f Format, layouts ...string) []string {
	if !f.IsPlainText {
		return layouts
	}

	newLayouts := make([]string, len(layouts))

	for i, l := range layouts {
		newLayouts[i] = "_text/" + l
	}

	return newLayouts
}

func replaceKeyValues(s string, oldNew ...string) string {
	replacer := strings.NewReplacer(oldNew...)
	return replacer.Replace(s)
}
