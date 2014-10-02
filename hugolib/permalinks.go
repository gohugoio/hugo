package hugolib

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/hugo/helpers"
)

// PathPattern represents a string which builds up a URL from attributes
type PathPattern string

// PageToPermaAttribute is the type of a function which, given a page and a tag
// can return a string to go in that position in the page (or an error)
type PageToPermaAttribute func(*Page, string) (string, error)

// PermalinkOverrides maps a section name to a PathPattern
type PermalinkOverrides map[string]PathPattern

// knownPermalinkAttributes maps :tags in a permalink specification to a
// function which, given a page and the tag, returns the resulting string
// to be used to replace that tag.
var knownPermalinkAttributes map[string]PageToPermaAttribute

// validate determines if a PathPattern is well-formed
func (pp PathPattern) validate() bool {
	fragments := strings.Split(string(pp[1:]), "/")
	var bail = false
	for i := range fragments {
		if bail {
			return false
		}
		if len(fragments[i]) == 0 {
			bail = true
			continue
		}
		if !strings.HasPrefix(fragments[i], ":") {
			continue
		}
		k := strings.ToLower(fragments[i][1:])
		if _, ok := knownPermalinkAttributes[k]; !ok {
			return false
		}
	}
	return true
}

type permalinkExpandError struct {
	pattern PathPattern
	section string
	err     error
}

func (pee *permalinkExpandError) Error() string {
	return fmt.Sprintf("error expanding %q section %q: %s", string(pee.pattern), pee.section, pee.err)
}

var (
	errPermalinkIllFormed        = errors.New("permalink ill-formed")
	errPermalinkAttributeUnknown = errors.New("permalink attribute not recognised")
)

// Expand on a PathPattern takes a Page and returns the fully expanded Permalink
// or an error explaining the failure.
func (pp PathPattern) Expand(p *Page) (string, error) {
	if !pp.validate() {
		return "", &permalinkExpandError{pattern: pp, section: "<all>", err: errPermalinkIllFormed}
	}
	sections := strings.Split(string(pp), "/")
	for i, field := range sections {
		if len(field) == 0 || field[0] != ':' {
			continue
		}
		attr := field[1:]
		callback, ok := knownPermalinkAttributes[attr]
		if !ok {
			return "", &permalinkExpandError{pattern: pp, section: strconv.Itoa(i), err: errPermalinkAttributeUnknown}
		}
		newField, err := callback(p, attr)
		if err != nil {
			return "", &permalinkExpandError{pattern: pp, section: strconv.Itoa(i), err: err}
		}
		sections[i] = newField
	}
	return strings.Join(sections, "/"), nil
}

func pageToPermalinkDate(p *Page, dateField string) (string, error) {
	// a Page contains a Node which provides a field Date, time.Time
	switch dateField {
	case "year":
		return strconv.Itoa(p.Date.Year()), nil
	case "month":
		return fmt.Sprintf("%02d", int(p.Date.Month())), nil
	case "monthname":
		return p.Date.Month().String(), nil
	case "day":
		return fmt.Sprintf("%02d", int(p.Date.Day())), nil
	case "weekday":
		return strconv.Itoa(int(p.Date.Weekday())), nil
	case "weekdayname":
		return p.Date.Weekday().String(), nil
	case "yearday":
		return strconv.Itoa(p.Date.YearDay()), nil
	}
	//TODO: support classic strftime escapes too
	// (and pass those through despite not being in the map)
	panic("coding error: should not be here")
}

// pageToPermalinkTitle returns the URL-safe form of the title
func pageToPermalinkTitle(p *Page, _ string) (string, error) {
	// Page contains Node which has Title
	// (also contains UrlPath which has Slug, sometimes)
	return helpers.Urlize(p.Title), nil
}

// pageToPermalinkFilename returns the URL-safe form of the filename
func pageToPermalinkFilename(p *Page, _ string) (string, error) {
	var extension = filepath.Ext(p.FileName)
	var name = p.FileName[0 : len(p.FileName)-len(extension)]
	return helpers.Urlize(name), nil
}

// if the page has a slug, return the slug, else return the title
func pageToPermalinkSlugElseTitle(p *Page, a string) (string, error) {
	if p.Slug != "" {
		// Don't start or end with a -
		if strings.HasPrefix(p.Slug, "-") {
			p.Slug = p.Slug[1:len(p.Slug)]
		}

		if strings.HasSuffix(p.Slug, "-") {
			p.Slug = p.Slug[0 : len(p.Slug)-1]
		}
		return p.Slug, nil
	}
	return pageToPermalinkTitle(p, a)
}

func pageToPermalinkSection(p *Page, _ string) (string, error) {
	// Page contains Node contains UrlPath which has Section
	return p.Section, nil
}

func init() {
	knownPermalinkAttributes = map[string]PageToPermaAttribute{
		"year":        pageToPermalinkDate,
		"month":       pageToPermalinkDate,
		"monthname":   pageToPermalinkDate,
		"day":         pageToPermalinkDate,
		"weekday":     pageToPermalinkDate,
		"weekdayname": pageToPermalinkDate,
		"yearday":     pageToPermalinkDate,
		"section":     pageToPermalinkSection,
		"title":       pageToPermalinkTitle,
		"slug":        pageToPermalinkSlugElseTitle,
		"filename":    pageToPermalinkFilename,
	}
}
