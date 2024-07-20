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

// Package output contains Output Format types and functions.
package output

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/gohugoio/hugo/media"
)

// Format represents an output representation, usually to a file on disk.
// <docsmeta>{ "name": "OutputFormat" }</docsmeta>
type Format struct {
	// The Name is used as an identifier. Internal output formats (i.e. html and rss)
	// can be overridden by providing a new definition for those types.
	// <docsmeta>{ "identifiers": ["html", "rss"] }</docsmeta>
	Name string `json:"-"`

	MediaType media.Type `json:"-"`

	// Must be set to a value when there are two or more conflicting mediatype for the same resource.
	Path string `json:"path"`

	// The base output file name used when not using "ugly URLs", defaults to "index".
	BaseName string `json:"baseName"`

	// The value to use for rel links.
	Rel string `json:"rel"`

	// The protocol to use, i.e. "webcal://". Defaults to the protocol of the baseURL.
	Protocol string `json:"protocol"`

	// IsPlainText decides whether to use text/template or html/template
	// as template parser.
	IsPlainText bool `json:"isPlainText"`

	// IsHTML returns whether this format is int the HTML family. This includes
	// HTML, AMP etc. This is used to decide when to create alias redirects etc.
	IsHTML bool `json:"isHTML"`

	// Enable to ignore the global uglyURLs setting.
	NoUgly bool `json:"noUgly"`

	// Enable to override the global uglyURLs setting.
	Ugly bool `json:"ugly"`

	// Enable if it doesn't make sense to include this format in an alternative
	// format listing, CSS being one good example.
	// Note that we use the term "alternative" and not "alternate" here, as it
	// does not necessarily replace the other format, it is an alternative representation.
	NotAlternative bool `json:"notAlternative"`

	// Eneable if this is a resource which path always starts at the root,
	// e.g. /robots.txt.
	Root bool `json:"root"`

	// Setting this will make this output format control the value of
	// .Permalink and .RelPermalink for a rendered Page.
	// If not set, these values will point to the main (first) output format
	// configured. That is probably the behavior you want in most situations,
	// as you probably don't want to link back to the RSS version of a page, as an
	// example. AMP would, however, be a good example of an output format where this
	// behavior is wanted.
	Permalinkable bool `json:"permalinkable"`

	// Setting this to a non-zero value will be used as the first sort criteria.
	Weight int `json:"weight"`
}

// Built-in output formats.
var (
	AMPFormat = Format{
		Name:          "amp",
		MediaType:     media.Builtin.HTMLType,
		BaseName:      "index",
		Path:          "amp",
		Rel:           "amphtml",
		IsHTML:        true,
		Permalinkable: true,
		// See https://www.ampproject.org/learn/overview/
	}

	CalendarFormat = Format{
		Name:        "calendar",
		MediaType:   media.Builtin.CalendarType,
		IsPlainText: true,
		Protocol:    "webcal://",
		BaseName:    "index",
		Rel:         "alternate",
	}

	CSSFormat = Format{
		Name:           "css",
		MediaType:      media.Builtin.CSSType,
		BaseName:       "styles",
		IsPlainText:    true,
		Rel:            "stylesheet",
		NotAlternative: true,
	}
	CSVFormat = Format{
		Name:        "csv",
		MediaType:   media.Builtin.CSVType,
		BaseName:    "index",
		IsPlainText: true,
		Rel:         "alternate",
	}

	HTMLFormat = Format{
		Name:          "html",
		MediaType:     media.Builtin.HTMLType,
		BaseName:      "index",
		Rel:           "canonical",
		IsHTML:        true,
		Permalinkable: true,

		// Weight will be used as first sort criteria. HTML will, by default,
		// be rendered first, but set it to 10 so it's easy to put one above it.
		Weight: 10,
	}

	MarkdownFormat = Format{
		Name:        "markdown",
		MediaType:   media.Builtin.MarkdownType,
		BaseName:    "index",
		Rel:         "alternate",
		IsPlainText: true,
	}

	JSONFormat = Format{
		Name:        "json",
		MediaType:   media.Builtin.JSONType,
		BaseName:    "index",
		IsPlainText: true,
		Rel:         "alternate",
	}

	WebAppManifestFormat = Format{
		Name:           "webappmanifest",
		MediaType:      media.Builtin.WebAppManifestType,
		BaseName:       "manifest",
		IsPlainText:    true,
		NotAlternative: true,
		Rel:            "manifest",
	}

	RobotsTxtFormat = Format{
		Name:        "robots",
		MediaType:   media.Builtin.TextType,
		BaseName:    "robots",
		IsPlainText: true,
		Root:        true,
		Rel:         "alternate",
	}

	RSSFormat = Format{
		Name:      "rss",
		MediaType: media.Builtin.RSSType,
		BaseName:  "index",
		NoUgly:    true,
		Rel:       "alternate",
	}

	SitemapFormat = Format{
		Name:      "sitemap",
		MediaType: media.Builtin.XMLType,
		BaseName:  "sitemap",
		Ugly:      true,
		Rel:       "sitemap",
	}

	SitemapIndexFormat = Format{
		Name:      "sitemapindex",
		MediaType: media.Builtin.XMLType,
		BaseName:  "sitemap",
		Ugly:      true,
		Root:      true,
		Rel:       "sitemap",
	}

	HTTPStatusHTMLFormat = Format{
		Name:           "httpstatus",
		MediaType:      media.Builtin.HTMLType,
		NotAlternative: true,
		Ugly:           true,
		IsHTML:         true,
		Permalinkable:  true,
	}
)

// DefaultFormats contains the default output formats supported by Hugo.
var DefaultFormats = Formats{
	AMPFormat,
	CalendarFormat,
	CSSFormat,
	CSVFormat,
	HTMLFormat,
	JSONFormat,
	MarkdownFormat,
	WebAppManifestFormat,
	RobotsTxtFormat,
	RSSFormat,
	SitemapFormat,
}

func init() {
	sort.Sort(DefaultFormats)
}

// Formats is a slice of Format.
// <docsmeta>{ "name": "OutputFormats" }</docsmeta>
type Formats []Format

func (formats Formats) Len() int      { return len(formats) }
func (formats Formats) Swap(i, j int) { formats[i], formats[j] = formats[j], formats[i] }
func (formats Formats) Less(i, j int) bool {
	fi, fj := formats[i], formats[j]
	if fi.Weight == fj.Weight {
		return fi.Name < fj.Name
	}

	if fj.Weight == 0 {
		return true
	}

	return fi.Weight > 0 && fi.Weight < fj.Weight
}

// GetBySuffix gets a output format given as suffix, e.g. "html".
// It will return false if no format could be found, or if the suffix given
// is ambiguous.
// The lookup is case insensitive.
func (formats Formats) GetBySuffix(suffix string) (f Format, found bool) {
	for _, ff := range formats {
		for _, suffix2 := range ff.MediaType.Suffixes() {
			if strings.EqualFold(suffix, suffix2) {
				if found {
					// ambiguous
					found = false
					return
				}
				f = ff
				found = true
			}
		}
	}
	return
}

// GetByName gets a format by its identifier name.
func (formats Formats) GetByName(name string) (f Format, found bool) {
	for _, ff := range formats {
		if strings.EqualFold(name, ff.Name) {
			f = ff
			found = true
			return
		}
	}
	return
}

// GetByNames gets a list of formats given a list of identifiers.
func (formats Formats) GetByNames(names ...string) (Formats, error) {
	var types []Format

	for _, name := range names {
		tpe, ok := formats.GetByName(name)
		if !ok {
			return types, fmt.Errorf("OutputFormat with key %q not found", name)
		}
		types = append(types, tpe)
	}
	return types, nil
}

// FromFilename gets a Format given a filename.
func (formats Formats) FromFilename(filename string) (f Format, found bool) {
	// mytemplate.amp.html
	// mytemplate.html
	// mytemplate
	var ext, outFormat string

	parts := strings.Split(filename, ".")
	if len(parts) > 2 {
		outFormat = parts[1]
		ext = parts[2]
	} else if len(parts) > 1 {
		ext = parts[1]
	}

	if outFormat != "" {
		return formats.GetByName(outFormat)
	}

	if ext != "" {
		f, found = formats.GetBySuffix(ext)
		if !found && len(parts) == 2 {
			// For extensionless output formats (e.g. Netlify's _redirects)
			// we must fall back to using the extension as format lookup.
			f, found = formats.GetByName(ext)
		}
	}
	return
}

// BaseFilename returns the base filename of f including an extension (ie.
// "index.xml").
func (f Format) BaseFilename() string {
	return f.BaseName + f.MediaType.FirstSuffix.FullSuffix
}

// IsZero returns true if f represents a zero value.
func (f Format) IsZero() bool {
	return f.Name == ""
}

// MarshalJSON returns the JSON encoding of f.
// For internal use only.
func (f Format) MarshalJSON() ([]byte, error) {
	type Alias Format
	return json.Marshal(&struct {
		MediaType string `json:"mediaType"`
		Alias
	}{
		MediaType: f.MediaType.String(),
		Alias:     (Alias)(f),
	})
}
