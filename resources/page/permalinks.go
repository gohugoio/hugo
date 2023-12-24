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

	"github.com/gohugoio/hugo/common/hstrings"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/resources/kinds"
)

// PermalinkExpander holds permalin mappings per section.
type PermalinkExpander struct {
	// knownPermalinkAttributes maps :tags in a permalink specification to a
	// function which, given a page and the tag, returns the resulting string
	// to be used to replace that tag.
	knownPermalinkAttributes map[string]pageToPermaAttribute

	expanders map[string]map[string]func(Page) (string, error)

	urlize func(uri string) string
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

	// Make sure this comes after all the other checks.
	if referenceTime.Format(attr) != attr {
		return p.pageToPermalinkDate, true
	}

	return nil, false
}

// NewPermalinkExpander creates a new PermalinkExpander configured by the given
// urlize func.
func NewPermalinkExpander(urlize func(uri string) string, patterns map[string]map[string]string) (PermalinkExpander, error) {
	p := PermalinkExpander{urlize: urlize}

	p.knownPermalinkAttributes = map[string]pageToPermaAttribute{
		"year":           p.pageToPermalinkDate,
		"month":          p.pageToPermalinkDate,
		"monthname":      p.pageToPermalinkDate,
		"day":            p.pageToPermalinkDate,
		"weekday":        p.pageToPermalinkDate,
		"weekdayname":    p.pageToPermalinkDate,
		"yearday":        p.pageToPermalinkDate,
		"section":        p.pageToPermalinkSection,
		"sections":       p.pageToPermalinkSections,
		"title":          p.pageToPermalinkTitle,
		"slug":           p.pageToPermalinkSlugElseTitle,
		"slugorfilename": p.pageToPermalinkSlugElseFilename,
		"filename":       p.pageToPermalinkFilename,
	}

	p.expanders = make(map[string]map[string]func(Page) (string, error))

	for kind, patterns := range patterns {
		e, err := p.parse(patterns)
		if err != nil {
			return p, err
		}
		p.expanders[kind] = e
	}

	return p, nil
}

// Expand expands the path in p according to the rules defined for the given key.
// If no rules are found for the given key, an empty string is returned.
func (l PermalinkExpander) Expand(key string, p Page) (string, error) {
	expanders, found := l.expanders[p.Kind()]

	if !found {
		return "", nil
	}

	expand, found := expanders[key]

	if !found {
		return "", nil
	}

	return expand(p)
}

// Allow " " and / to represent the root section.
var sectionCutSet = " /"

func init() {
	if string(os.PathSeparator) != "/" {
		sectionCutSet += string(os.PathSeparator)
	}
}

func (l PermalinkExpander) parse(patterns map[string]string) (map[string]func(Page) (string, error), error) {
	expanders := make(map[string]func(Page) (string, error))

	for k, pattern := range patterns {
		k = strings.Trim(k, sectionCutSet)

		if !l.validate(pattern) {
			return nil, &permalinkExpandError{pattern: pattern, err: errPermalinkIllFormed}
		}

		pattern := pattern
		matches := attributeRegexp.FindAllStringSubmatch(pattern, -1)

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

		expanders[k] = func(p Page) (string, error) {
			if matches == nil {
				return pattern, nil
			}

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

			return newField, nil
		}

	}

	return expanders, nil
}

// pageToPermaAttribute is the type of a function which, given a page and a tag
// can return a string to go in that position in the page (or an error)
type pageToPermaAttribute func(Page, string) (string, error)

var attributeRegexp = regexp.MustCompile(`:\w+(\[.+?\])?`)

// validate determines if a PathPattern is well-formed
func (l PermalinkExpander) validate(pp string) bool {
	if len(pp) == 0 {
		return false
	}
	fragments := strings.Split(pp[1:], "/")
	bail := false
	for i := range fragments {
		if bail {
			return false
		}
		if len(fragments[i]) == 0 {
			bail = true
			continue
		}

		matches := attributeRegexp.FindAllStringSubmatch(fragments[i], -1)
		if matches == nil {
			continue
		}

		for _, match := range matches {
			k := match[0][1:]
			if _, ok := l.callback(k); !ok {
				return false
			}
		}
	}
	return true
}

type permalinkExpandError struct {
	pattern string
	err     error
}

func (pee *permalinkExpandError) Error() string {
	return fmt.Sprintf("error expanding %q: %s", pee.pattern, pee.err)
}

var (
	errPermalinkIllFormed        = errors.New("permalink ill-formed")
	errPermalinkAttributeUnknown = errors.New("permalink attribute not recognised")
)

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

func (l PermalinkExpander) pageToPermalinkSections(p Page, _ string) (string, error) {
	return p.CurrentSection().SectionsPath(), nil
}

func (l PermalinkExpander) translationBaseName(p Page) string {
	if p.File() == nil {
		return ""
	}
	return p.File().TranslationBaseName()
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

var permalinksKindsSupport = []string{kinds.KindPage, kinds.KindSection, kinds.KindTaxonomy, kinds.KindTerm}

// DecodePermalinksConfig decodes the permalinks configuration in the given map
func DecodePermalinksConfig(m map[string]any) (map[string]map[string]string, error) {
	permalinksConfig := make(map[string]map[string]string)

	permalinksConfig[kinds.KindPage] = make(map[string]string)
	permalinksConfig[kinds.KindSection] = make(map[string]string)
	permalinksConfig[kinds.KindTaxonomy] = make(map[string]string)
	permalinksConfig[kinds.KindTerm] = make(map[string]string)

	config := maps.CleanConfigStringMap(m)
	for k, v := range config {
		switch v := v.(type) {
		case string:
			// [permalinks]
			//   key = '...'

			// To sucessfully be backward compatible, "default" patterns need to be set for both page and term
			permalinksConfig[kinds.KindPage][k] = v
			permalinksConfig[kinds.KindTerm][k] = v

		case maps.Params:
			// [permalinks.key]
			//   xyz = ???

			if hstrings.InSlice(permalinksKindsSupport, k) {
				// TODO: warn if we overwrite an already set value
				for k2, v2 := range v {
					switch v2 := v2.(type) {
					case string:
						permalinksConfig[k][k2] = v2

					default:
						return nil, fmt.Errorf("permalinks configuration invalid: unknown value %q for key %q for kind %q", v2, k2, k)
					}
				}
			} else {
				return nil, fmt.Errorf("permalinks configuration not supported for kind %q, supported kinds are %v", k, permalinksKindsSupport)
			}

		default:
			return nil, fmt.Errorf("permalinks configuration invalid: unknown value %q for key %q", v, k)
		}
	}
	return permalinksConfig, nil
}
