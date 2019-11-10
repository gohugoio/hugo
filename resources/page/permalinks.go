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
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/gohugoio/hugo/helpers"
)

// PermalinkExpander holds permalin mappings per section.
type PermalinkExpander struct {
	// knownPermalinkAttributes maps :tags in a permalink specification to a
	// function which, given a page and the tag, returns the resulting string
	// to be used to replace that tag.
	knownPermalinkAttributes map[string]pageToPermaAttribute

	expanders map[string]func(Page) (string, error)

	ps *helpers.PathSpec
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

	if referenceTime.Format(attr) != attr {
		return p.pageToPermalinkDate, true
	}

	return nil, false
}

// NewPermalinkExpander creates a new PermalinkExpander configured by the given
// PathSpec.
func NewPermalinkExpander(ps *helpers.PathSpec) (PermalinkExpander, error) {

	p := PermalinkExpander{ps: ps}

	p.knownPermalinkAttributes = map[string]pageToPermaAttribute{
		"year":        p.pageToPermalinkDate,
		"month":       p.pageToPermalinkDate,
		"monthname":   p.pageToPermalinkDate,
		"day":         p.pageToPermalinkDate,
		"weekday":     p.pageToPermalinkDate,
		"weekdayname": p.pageToPermalinkDate,
		"yearday":     p.pageToPermalinkDate,
		"section":     p.pageToPermalinkSection,
		"sections":    p.pageToPermalinkSections,
		"title":       p.pageToPermalinkTitle,
		"slug":        p.pageToPermalinkSlugElseTitle,
		"filename":    p.pageToPermalinkFilename,
	}

	patterns := ps.Cfg.GetStringMapString("permalinks")
	if patterns == nil {
		return p, nil
	}

	e, err := p.parse(patterns)
	if err != nil {
		return p, err
	}

	p.expanders = e

	return p, nil
}

// Expand expands the path in p according to the rules defined for the given key.
// If no rules are found for the given key, an empty string is returned.
func (l PermalinkExpander) Expand(key string, p Page) (string, error) {
	expand, found := l.expanders[key]

	if !found {
		return "", nil
	}

	return expand(p)

}

func (l PermalinkExpander) parse(patterns map[string]string) (map[string]func(Page) (string, error), error) {

	expanders := make(map[string]func(Page) (string, error))

	// Allow " " and / to represent the root section.
	const sectionCutSet = " /" + string(os.PathSeparator)

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

var attributeRegexp = regexp.MustCompile(`:\w+`)

// validate determines if a PathPattern is well-formed
func (l PermalinkExpander) validate(pp string) bool {
	fragments := strings.Split(pp[1:], "/")
	var bail = false
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
	return l.ps.URLize(p.Title()), nil
}

// pageToPermalinkFilename returns the URL-safe form of the filename
func (l PermalinkExpander) pageToPermalinkFilename(p Page, _ string) (string, error) {
	name := p.File().TranslationBaseName()
	if name == "index" {
		// Page bundles; the directory name will hopefully have a better name.
		dir := strings.TrimSuffix(p.File().Dir(), helpers.FilePathSeparator)
		_, name = filepath.Split(dir)
	}

	return l.ps.URLize(name), nil
}

// if the page has a slug, return the slug, else return the title
func (l PermalinkExpander) pageToPermalinkSlugElseTitle(p Page, a string) (string, error) {
	if p.Slug() != "" {
		return l.ps.URLize(p.Slug()), nil
	}
	return l.pageToPermalinkTitle(p, a)
}

func (l PermalinkExpander) pageToPermalinkSection(p Page, _ string) (string, error) {
	return p.Section(), nil
}

func (l PermalinkExpander) pageToPermalinkSections(p Page, _ string) (string, error) {
	return p.CurrentSection().SectionsPath(), nil
}
