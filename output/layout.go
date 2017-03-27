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
	"sync"
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
type LayoutHandler struct {
	hasTheme bool

	mu    sync.RWMutex
	cache map[layoutCacheKey][]string
}

type layoutCacheKey struct {
	d              LayoutDescriptor
	layoutOverride string
	f              Format
}

func NewLayoutHandler(hasTheme bool) *LayoutHandler {
	return &LayoutHandler{hasTheme: hasTheme, cache: make(map[layoutCacheKey][]string)}
}

// RSS:
// Home:"rss.xml", "_default/rss.xml", "_internal/_default/rss.xml"
// Section: "section/" + section + ".rss.xml", "_default/rss.xml", "rss.xml", "_internal/_default/rss.xml"
// Taxonomy "taxonomy/" + singular + ".rss.xml", "_default/rss.xml", "rss.xml", "_internal/_default/rss.xml"
// Tax term: taxonomy/" + singular + ".terms.rss.xml", "_default/rss.xml", "rss.xml", "_internal/_default/rss.xml"

const (

	// The RSS templates doesn't map easily into the regular pages.
	layoutsRSSHome         = `NAME.SUFFIX _default/NAME.SUFFIX _internal/_default/rss.xml`
	layoutsRSSSection      = `section/SECTION.NAME.SUFFIX _default/NAME.SUFFIX NAME.SUFFIX _internal/_default/rss.xml`
	layoutsRSSTaxonomy     = `taxonomy/SECTION.NAME.SUFFIX _default/NAME.SUFFIX NAME.SUFFIX _internal/_default/rss.xml`
	layoutsRSSTaxonomyTerm = `taxonomy/SECTION.terms.NAME.SUFFIX _default/NAME.SUFFIX NAME.SUFFIX _internal/_default/rss.xml`

	layoutsHome    = "index.NAME.SUFFIX index.SUFFIX _default/list.NAME.SUFFIX _default/list.SUFFIX"
	layoutsSection = `
section/SECTION.NAME.SUFFIX section/SECTION.SUFFIX
SECTION/list.NAME.SUFFIX SECTION/list.SUFFIX
_default/section.NAME.SUFFIX _default/section.SUFFIX
_default/list.NAME.SUFFIX _default/list.SUFFIX
indexes/SECTION.NAME.SUFFIX indexes/SECTION.SUFFIX
_default/indexes.NAME.SUFFIX _default/indexes.SUFFIX
`
	layoutsTaxonomy = `
taxonomy/SECTION.NAME.SUFFIX taxonomy/SECTION.SUFFIX
indexes/SECTION.NAME.SUFFIX indexes/SECTION.SUFFIX 
_default/taxonomy.NAME.SUFFIX _default/taxonomy.SUFFIX
_default/list.NAME.SUFFIX _default/list.SUFFIX
`
	layoutsTaxonomyTerm = `
taxonomy/SECTION.terms.NAME.SUFFIX taxonomy/SECTION.terms.SUFFIX
_default/terms.NAME.SUFFIX _default/terms.SUFFIX
indexes/indexes.NAME.SUFFIX indexes/indexes.SUFFIX
`
)

func (l *LayoutHandler) For(d LayoutDescriptor, layoutOverride string, f Format) ([]string, error) {

	// We will get lots of requests for the same layouts, so avoid recalculations.
	key := layoutCacheKey{d, layoutOverride, f}
	l.mu.RLock()
	if cacheVal, found := l.cache[key]; found {
		l.mu.RUnlock()
		return cacheVal, nil
	}
	l.mu.RUnlock()

	var layouts []string

	if layoutOverride != "" && d.Kind != "page" {
		return layouts, fmt.Errorf("Custom layout (%q) only supported for regular pages, not kind %q", layoutOverride, d.Kind)
	}

	layout := d.Layout

	if layoutOverride != "" {
		layout = layoutOverride
	}

	isRSS := f.Name == RSSFormat.Name

	if d.Kind == "page" {
		if isRSS {
			return []string{}, nil
		}
		layouts = regularPageLayouts(d.Type, layout, f)
	} else {
		if isRSS {
			layouts = resolveListTemplate(d, f,
				layoutsRSSHome,
				layoutsRSSSection,
				layoutsRSSTaxonomy,
				layoutsRSSTaxonomyTerm)
		} else {
			layouts = resolveListTemplate(d, f,
				layoutsHome,
				layoutsSection,
				layoutsTaxonomy,
				layoutsTaxonomyTerm)
		}
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

		layouts = layoutsWithThemeLayouts
	}

	layouts = prependTextPrefixIfNeeded(f, layouts...)

	l.mu.Lock()
	l.cache[key] = layouts
	l.mu.Unlock()

	return layouts, nil
}

func resolveListTemplate(d LayoutDescriptor, f Format,
	homeLayouts,
	sectionLayouts,
	taxonomyLayouts,
	taxonomyTermLayouts string) []string {
	var layouts []string

	switch d.Kind {
	case "home":
		layouts = resolveTemplate(homeLayouts, d, f)
	case "section":
		layouts = resolveTemplate(sectionLayouts, d, f)
	case "taxonomy":
		layouts = resolveTemplate(taxonomyLayouts, d, f)
	case "taxonomyTerm":
		layouts = resolveTemplate(taxonomyTermLayouts, d, f)
	}

	return layouts
}

func resolveTemplate(templ string, d LayoutDescriptor, f Format) []string {
	layouts := strings.Fields(replaceKeyValues(templ,
		"SUFFIX", f.MediaType.Suffix,
		"NAME", strings.ToLower(f.Name),
		"SECTION", d.Section))

	return layouts
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

func regularPageLayouts(types string, layout string, f Format) []string {
	var layouts []string

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

	return layouts
}
