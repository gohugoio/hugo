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

package page

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gohugoio/hugo/common/hmaps"
	"github.com/gohugoio/hugo/common/hstrings"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugolib/sitesmatrix"
	"github.com/gohugoio/hugo/resources/kinds"
	"github.com/mitchellh/mapstructure"
)

// PermalinkConfig holds a single permalink rule with a target matcher and a pattern.
type PermalinkConfig struct {
	Target  PageMatcher
	Pattern string
}

// PermalinksConfig is an ordered slice of permalink rules.
// For any given Page, the first matching rule wins.
type PermalinksConfig []PermalinkConfig

// InitConfig compiles the sites matrix for each permalink target.
func (c PermalinksConfig) InitConfig(logger loggers.Logger, defaultSitesMatrix sitesmatrix.VectorStore, configuredDimensions *sitesmatrix.ConfiguredDimensions) error {
	for i := range c {
		if err := c[i].Target.compileSitesMatrix(nil, configuredDimensions); err != nil {
			return fmt.Errorf("failed to compile permalink target %d: %w", i, err)
		}
	}
	return nil
}

// PermalinkExpander holds permalink mappings.
type PermalinkExpander struct {
	// knownPermalinkAttributes maps :tags in a permalink specification to a
	// function which, given a page and the tag, returns the resulting string
	// to be used to replace that tag.
	knownPermalinkAttributes map[string]pageToPermaAttribute

	configs PermalinksConfig

	urlize func(uri string) string

	patternCache *hmaps.Cache[string, func(Page) (string, error)]
}

// Time for checking date formats. Every field is different than the
// Go reference time for date formatting. This ensures that formatting this date
// with a Go time format always has a different output than the format itself.
var referenceTime = time.Date(2019, time.November, 9, 23, 1, 42, 1, time.UTC)

// Return the callback for the given permalink attribute and a boolean indicating if the attribute is valid or not.
func (p PermalinkExpander) callback(attr string) (pageToPermaAttribute, bool) {
	if callback, ok := p.knownPermalinkAttributes[attr]; ok {
		return callback, true
	}

	if strings.HasPrefix(attr, "sections[") {
		fn := p.toSliceFunc(strings.TrimPrefix(attr, "sections"))
		return func(p Page, s string) (string, error) {
			return path.Join(fn(p.CurrentSection().SectionsEntries())...), nil
		}, true
	}

	if strings.HasPrefix(attr, "sectionslugs[") {
		fn := p.toSliceFunc(strings.TrimPrefix(attr, "sectionslugs"))
		sectionSlugsFunc := p.withSectionPagesFunc(p.pageToPermalinkSlugElseTitle, func(s ...string) string {
			return path.Join(fn(s)...)
		})
		return sectionSlugsFunc, true
	}

	// Make sure this comes after all the other checks.
	if referenceTime.Format(attr) != attr {
		return p.pageToPermalinkDate, true
	}

	return nil, false
}

// NewPermalinkExpander creates a new PermalinkExpander configured by the given
// urlize func.
func NewPermalinkExpander(urlize func(uri string) string, configs PermalinksConfig) (PermalinkExpander, error) {
	p := PermalinkExpander{
		urlize:       urlize,
		configs:      configs,
		patternCache: hmaps.NewCache[string, func(Page) (string, error)](),
	}

	p.knownPermalinkAttributes = map[string]pageToPermaAttribute{
		"year":                  p.pageToPermalinkDate,
		"month":                 p.pageToPermalinkDate,
		"monthname":             p.pageToPermalinkDate,
		"day":                   p.pageToPermalinkDate,
		"weekday":               p.pageToPermalinkDate,
		"weekdayname":           p.pageToPermalinkDate,
		"yearday":               p.pageToPermalinkDate,
		"section":               p.pageToPermalinkSection,
		"sectionslug":           p.pageToPermalinkSectionSlug,
		"sections":              p.pageToPermalinkSections,
		"sectionslugs":          p.pageToPermalinkSectionSlugs,
		"title":                 p.pageToPermalinkTitle,
		"slug":                  p.pageToPermalinkSlugElseTitle,
		"slugorfilename":        p.pageToPermalinkSlugElseFilename,
		"filename":              p.pageToPermalinkFilename,
		"contentbasename":       p.pageToPermalinkContentBaseName,
		"slugorcontentbasename": p.pageToPermalinkSlugOrContentBaseName,
	}

	// Validate all patterns at init time.
	for _, cfg := range configs {
		if _, err := p.getOrParsePattern(cfg.Pattern); err != nil {
			return p, err
		}
	}

	return p, nil
}

// Escape sequence for colons in permalink patterns.
const escapePlaceholderColon = "\x00"

func (l PermalinkExpander) normalizeEscapeSequencesIn(s string) (string, bool) {
	s2 := strings.ReplaceAll(s, "\\:", escapePlaceholderColon)
	return s2, s2 != s
}

func (l PermalinkExpander) normalizeEscapeSequencesOut(result string) string {
	return strings.ReplaceAll(result, escapePlaceholderColon, ":")
}

// ExpandPattern expands the path in p with the specified expand pattern.
func (l PermalinkExpander) ExpandPattern(pattern string, p Page) (string, error) {
	expand, err := l.getOrParsePattern(pattern)
	if err != nil {
		return "", err
	}

	return expand(p)
}

// Expand expands the path in p according to the first matching permalink rule.
// If no rules match, an empty string is returned.
func (l PermalinkExpander) Expand(p Page) (string, error) {
	kind := p.Kind()
	if !hstrings.InSlice(permalinksKindsSupport, kind) {
		return "", nil
	}
	var siteVector sitesmatrix.VectorProvider
	if sv := GetSiteVector(p); sv != (sitesmatrix.Vector{}) {
		siteVector = sv
	}
	for _, cfg := range l.configs {
		if cfg.Target.Match(kind, p.Path(), "", siteVector) {
			return l.ExpandPattern(cfg.Pattern, p)
		}
	}
	return "", nil
}

// Allow " " and / to represent the root section.
var sectionCutSet = " /"

func init() {
	if string(os.PathSeparator) != "/" {
		sectionCutSet += string(os.PathSeparator)
	}
}

func (l PermalinkExpander) getOrParsePattern(pattern string) (func(Page) (string, error), error) {
	return l.patternCache.GetOrCreate(pattern, func() (func(Page) (string, error), error) {
		var normalized bool
		pattern, normalized = l.normalizeEscapeSequencesIn(pattern)

		matches := attributeRegexp.FindAllStringSubmatch(pattern, -1)
		if matches == nil {
			result := pattern
			if normalized {
				result = l.normalizeEscapeSequencesOut(result)
			}
			return func(p Page) (string, error) {
				return result, nil
			}, nil
		}

		callbacks := make([]pageToPermaAttribute, len(matches))
		replacements := make([]string, len(matches))
		for i, m := range matches {
			replacement := m[0]
			attr := replacement[1:]
			replacements[i] = replacement
			callback, ok := l.callback(attr)

			if !ok {
				return nil, &permalinkExpandError{pattern: pattern, err: errPermalinkAttributeUnknown}
			}

			callbacks[i] = callback
		}

		return func(p Page) (string, error) {
			newField := pattern

			for i, replacement := range replacements {
				attr := replacement[1:]
				callback := callbacks[i]
				newAttr, err := callback(p, attr)
				if err != nil {
					return "", &permalinkExpandError{pattern: pattern, err: err}
				}

				newField = strings.Replace(newField, replacement, newAttr, 1)
			}

			if normalized {
				newField = l.normalizeEscapeSequencesOut(newField)
			}

			return newField, nil
		}, nil
	})
}

// pageToPermaAttribute is the type of a function which, given a page and a tag
// can return a string to go in that position in the page (or an error)
type pageToPermaAttribute func(Page, string) (string, error)

var attributeRegexp = regexp.MustCompile(`:\w+(\[.+?\])?`)

type permalinkExpandError struct {
	pattern string
	err     error
}

func (pee *permalinkExpandError) Error() string {
	return fmt.Sprintf("error expanding %q: %s", pee.pattern, pee.err)
}

var errPermalinkAttributeUnknown = errors.New("permalink attribute not recognised")

func (l PermalinkExpander) pageToPermalinkDate(p Page, dateField string) (string, error) {
	// a Page contains a Node which provides a field Date, time.Time
	switch dateField {
	case "year":
		return strconv.Itoa(p.Date().Year()), nil
	case "month":
		return fmt.Sprintf("%02d", int(p.Date().Month())), nil
	case "monthname":
		return p.Date().Month().String(), nil
	case "day":
		return fmt.Sprintf("%02d", p.Date().Day()), nil
	case "weekday":
		return strconv.Itoa(int(p.Date().Weekday())), nil
	case "weekdayname":
		return p.Date().Weekday().String(), nil
	case "yearday":
		return strconv.Itoa(p.Date().YearDay()), nil
	}

	return p.Date().Format(dateField), nil
}

// pageToPermalinkTitle returns the URL-safe form of the title
func (l PermalinkExpander) pageToPermalinkTitle(p Page, _ string) (string, error) {
	return l.urlize(p.Title()), nil
}

// pageToPermalinkFilename returns the URL-safe form of the filename
func (l PermalinkExpander) pageToPermalinkFilename(p Page, _ string) (string, error) {
	name := l.translationBaseName(p)
	if name == "index" {
		// Page bundles; the directory name will hopefully have a better name.
		dir := strings.TrimSuffix(p.File().Dir(), helpers.FilePathSeparator)
		_, name = filepath.Split(dir)
	} else if name == "_index" {
		return "", nil
	}

	return l.urlize(name), nil
}

// if the page has a slug, return the slug, else return the title
func (l PermalinkExpander) pageToPermalinkSlugElseTitle(p Page, a string) (string, error) {
	if p.Slug() != "" {
		return l.urlize(p.Slug()), nil
	}
	return l.pageToPermalinkTitle(p, a)
}

// if the page has a slug, return the slug, else return the filename
func (l PermalinkExpander) pageToPermalinkSlugElseFilename(p Page, a string) (string, error) {
	if p.Slug() != "" {
		return l.urlize(p.Slug()), nil
	}
	return l.pageToPermalinkFilename(p, a)
}

func (l PermalinkExpander) pageToPermalinkSection(p Page, _ string) (string, error) {
	return p.Section(), nil
}

// pageToPermalinkSectionSlug returns the URL-safe form of the first section's slug or title
func (l PermalinkExpander) pageToPermalinkSectionSlug(p Page, attr string) (string, error) {
	sectionPage := p.FirstSection()
	if sectionPage == nil || sectionPage.IsHome() {
		return "", nil
	}
	return l.pageToPermalinkSlugElseTitle(sectionPage, attr)
}

func (l PermalinkExpander) pageToPermalinkSections(p Page, _ string) (string, error) {
	return p.CurrentSection().SectionsPath(), nil
}

// pageToPermalinkSectionSlugs returns a path built from all ancestor sections using their slugs or titles
func (l PermalinkExpander) pageToPermalinkSectionSlugs(p Page, attr string) (string, error) {
	sectionSlugsFunc := l.withSectionPagesFunc(l.pageToPermalinkSlugElseTitle, path.Join)
	return sectionSlugsFunc(p, attr)
}

// pageToPermalinkContentBaseName returns the URL-safe form of the content base name.
func (l PermalinkExpander) pageToPermalinkContentBaseName(p Page, _ string) (string, error) {
	return l.urlize(p.PathInfo().Unnormalized().BaseNameNoIdentifier()), nil
}

// pageToPermalinkSlugOrContentBaseName returns the URL-safe form of the slug, content base name.
func (l PermalinkExpander) pageToPermalinkSlugOrContentBaseName(p Page, a string) (string, error) {
	if p.Slug() != "" {
		return l.urlize(p.Slug()), nil
	}
	name, err := l.pageToPermalinkContentBaseName(p, a)
	if err != nil {
		return "", nil
	}
	return name, nil
}

func (l PermalinkExpander) translationBaseName(p Page) string {
	if p.File() == nil {
		return ""
	}
	return p.File().TranslationBaseName()
}

// withSectionPagesFunc returns a function that builds permalink attributes from section pages.
// It applies the transformation function f to each ancestor section (Page), then joins the results with the join function.
//
// Current use is create section-based hierarchical paths using section slugs.
func (l PermalinkExpander) withSectionPagesFunc(f func(Page, string) (string, error), join func(...string) string) func(p Page, s string) (string, error) {
	return func(p Page, s string) (string, error) {
		var entries []string
		currentSection := p.CurrentSection()

		// Build section hierarchy: ancestors (reversed to root-first) + current section
		sections := currentSection.Ancestors().Reverse()
		sections = append(sections, currentSection)

		for _, section := range sections {
			if section.IsHome() {
				continue
			}
			entry, err := f(section, s)
			if err != nil {
				return "", err
			}
			entries = append(entries, entry)
		}

		return join(entries...), nil
	}
}

var (
	nilSliceFunc = func(s []string) []string {
		return nil
	}
	allSliceFunc = func(s []string) []string {
		return s
	}
)

// toSliceFunc returns a slice func that slices s according to the cut spec.
// The cut spec must be on form [low:high] (one or both can be omitted),
// also allowing single slice indices (e.g. [2]) and the special [last] keyword
// giving the last element of the slice.
// The returned function will be lenient and not panic in out of bounds situation.
//
// The current use case for this is to use parts of the sections path in permalinks.
func (l PermalinkExpander) toSliceFunc(cut string) func(s []string) []string {
	cut = strings.ToLower(strings.TrimSpace(cut))
	if cut == "" {
		return allSliceFunc
	}

	if len(cut) < 3 || (cut[0] != '[' || cut[len(cut)-1] != ']') {
		return nilSliceFunc
	}

	toNFunc := func(s string, low bool) func(ss []string) int {
		if s == "" {
			if low {
				return func(ss []string) int {
					return 0
				}
			} else {
				return func(ss []string) int {
					return len(ss)
				}
			}
		}

		if s == "last" {
			return func(ss []string) int {
				return len(ss) - 1
			}
		}

		n, _ := strconv.Atoi(s)
		if n < 0 {
			n = 0
		}
		return func(ss []string) int {
			// Prevent out of bound situations. It would not make
			// much sense to panic here.
			if n >= len(ss) {
				if low {
					return -1
				}
				return len(ss)
			}
			return n
		}
	}

	opsStr := cut[1 : len(cut)-1]
	opts := strings.Split(opsStr, ":")

	if !strings.Contains(opsStr, ":") {
		toN := toNFunc(opts[0], true)
		return func(s []string) []string {
			if len(s) == 0 {
				return nil
			}
			n := toN(s)
			if n < 0 {
				return []string{}
			}
			v := s[n]
			if v == "" {
				return nil
			}
			return []string{v}
		}
	}

	toN1, toN2 := toNFunc(opts[0], true), toNFunc(opts[1], false)

	return func(s []string) []string {
		if len(s) == 0 {
			return nil
		}
		n1, n2 := toN1(s), toN2(s)
		if n1 < 0 || n2 < 0 {
			return []string{}
		}
		return s[n1:n2]
	}
}

var permalinksKindsSupport = []string{kinds.KindPage, kinds.KindHome, kinds.KindSection, kinds.KindTaxonomy, kinds.KindTerm}

func sectionToPathGlob(section string) string {
	section = strings.Trim(section, sectionCutSet)
	if section == "" {
		return "/*"
	}
	return "/{" + section + "," + section + "/**}"
}

// DecodePermalinksConfig decodes the permalinks configuration.
// It supports both the new slice-based format and the legacy map-based formats.
func DecodePermalinksConfig(in any) (PermalinksConfig, error) {
	if in == nil {
		return nil, nil
	}

	// Try the new slice-based format first.
	if configs, err := decodePermalinksSlice(in); err == nil && configs != nil {
		return configs, nil
	}

	// Fall back to legacy map-based formats.
	switch v := in.(type) {
	case map[string]any:
		return decodePermalinksMap(v)
	case hmaps.Params:
		return decodePermalinksMap(v)
	default:
		return nil, fmt.Errorf("permalinks: unsupported config type %T", in)
	}
}

func decodePermalinksSlice(in any) (PermalinksConfig, error) {
	ms, err := hmaps.ToSliceStringMap(in)
	if err != nil {
		return nil, err
	}

	var configs PermalinksConfig
	for _, m := range ms {
		m = hmaps.CleanConfigStringMap(m)
		var cfg PermalinkConfig

		if targetVal, ok := m["target"]; ok {
			if err := mapstructure.WeakDecode(targetVal, &cfg.Target); err != nil {
				return nil, fmt.Errorf("permalinks: failed to decode target: %w", err)
			}
			cfg.Target.Kind = strings.ToLower(cfg.Target.Kind)
			cfg.Target.Path = filepath.ToSlash(strings.ToLower(cfg.Target.Path))
		}

		if patternVal, ok := m["pattern"]; ok {
			cfg.Pattern, ok = patternVal.(string)
			if !ok {
				return nil, fmt.Errorf("permalinks: pattern must be a string, got %T", patternVal)
			}
		} else {
			return nil, fmt.Errorf("permalinks: missing pattern")
		}

		configs = append(configs, cfg)
	}

	return configs, nil
}

func decodePermalinksMap(m map[string]any) (PermalinksConfig, error) {
	var configs PermalinksConfig
	config := hmaps.CleanConfigStringMap(m)

	for k, v := range config {
		switch v := v.(type) {
		case string:
			// [permalinks]
			//   key = '...'
			// Backward compat: set for both page and term.
			configs = append(configs,
				PermalinkConfig{Target: PageMatcher{Kind: kinds.KindPage, Path: sectionToPathGlob(k)}, Pattern: v},
				PermalinkConfig{Target: PageMatcher{Kind: kinds.KindTerm, Path: sectionToPathGlob(k)}, Pattern: v},
			)

		case hmaps.Params:
			// [permalinks.kind]
			//   section = '...'
			if !hstrings.InSlice(permalinksKindsSupport, k) {
				return nil, fmt.Errorf("permalinks configuration not supported for kind %q, supported kinds are %v", k, permalinksKindsSupport)
			}
			for k2, v2 := range v {
				switch v2 := v2.(type) {
				case string:
					configs = append(configs,
						PermalinkConfig{Target: PageMatcher{Kind: k, Path: sectionToPathGlob(k2)}, Pattern: v2},
					)
				default:
					return nil, fmt.Errorf("permalinks configuration invalid: unknown value %q for key %q for kind %q", v2, k2, k)
				}
			}

		default:
			return nil, fmt.Errorf("permalinks configuration invalid: unknown value %q for key %q", v, k)
		}
	}
	return configs, nil
}
