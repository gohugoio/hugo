// Copyright 2024 The Hugo Authors. All rights reserved.
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

package pagemeta

import (
	"strings"
	"time"

	"github.com/gohugoio/hugo/common/htime"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/resources/page"

	"github.com/gohugoio/hugo/helpers"

	"github.com/gohugoio/hugo/config"
	"github.com/spf13/cast"
)

type Dates struct {
	Date        time.Time
	Lastmod     time.Time
	PublishDate time.Time
	ExpiryDate  time.Time
}

func (d Dates) IsDateOrLastModAfter(in Dates) bool {
	return d.Date.After(in.Date) || d.Lastmod.After(in.Lastmod)
}

func (d *Dates) UpdateDateAndLastmodIfAfter(in Dates) {
	if in.Date.After(d.Date) {
		d.Date = in.Date
	}
	if in.Lastmod.After(d.Lastmod) {
		d.Lastmod = in.Lastmod
	}
}

func (d Dates) IsAllDatesZero() bool {
	return d.Date.IsZero() && d.Lastmod.IsZero() && d.PublishDate.IsZero() && d.ExpiryDate.IsZero()
}

// PageConfig configures a Page, typically from front matter.
// Note that all the top level fields are reserved Hugo keywords.
// Any custom configuration needs to be set in the Params map.
type PageConfig struct {
	Dates                   // Dates holds the fource core dates for this page.
	Title          string   // The title of the page.
	LinkTitle      string   // The link title of the page.
	Type           string   // The content type of the page.
	Layout         string   // The layout to use for to render this page.
	Markup         string   // The markup used in the content file.
	Weight         int      // The weight of the page, used in sorting if set to a non-zero value.
	Kind           string   // The kind of page, e.g. "page", "section", "home" etc. This is usually derived from the content path.
	Path           string   // The canonical path to the page, e.g. /sect/mypage. Note: Leading slash, no trailing slash, no extensions or language identifiers.
	Lang           string   // The language code for this page. This is usually derived from the module mount or filename.
	Slug           string   // The slug for this page.
	Description    string   // The description for this page.
	Summary        string   // The summary for this page.
	Draft          bool     // Whether or not the content is a draft.
	Headless       bool     // Whether or not the page should be rendered.
	IsCJKLanguage  bool     // Whether or not the content is in a CJK language.
	TranslationKey string   // The translation key for this page.
	Keywords       []string // The keywords for this page.
	Aliases        []string // The aliases for this page.
	Outputs        []string // The output formats to render this page in. If not set, the site's configured output formats for this page kind will be used.

	// These build options are set in the front matter,
	// but not passed on to .Params.
	Resources []map[string]any
	Cascade   map[page.PageMatcher]maps.Params // Only relevant for branch nodes.
	Sitemap   config.SitemapConfig
	Build     BuildConfig

	// User defined params.
	Params maps.Params
}

// FrontMatterHandler maps front matter into Page fields and .Params.
// Note that we currently have only extracted the date logic.
type FrontMatterHandler struct {
	fmConfig FrontmatterConfig

	dateHandler        frontMatterFieldHandler
	lastModHandler     frontMatterFieldHandler
	publishDateHandler frontMatterFieldHandler
	expiryDateHandler  frontMatterFieldHandler

	// A map of all date keys configured, including any custom.
	allDateKeys map[string]bool

	logger loggers.Logger
}

// FrontMatterDescriptor describes how to handle front matter for a given Page.
// It has pointers to values in the receiving page which gets updated.
type FrontMatterDescriptor struct {
	// This is the Page's base filename (BaseFilename), e.g. page.md., or
	// if page is a leaf bundle, the bundle folder name (ContentBaseName).
	BaseFilename string

	// The content file's mod time.
	ModTime time.Time

	// May be set from the author date in Git.
	GitAuthorDate time.Time

	// The below will be modified.
	PageConfig *PageConfig

	// The Location to use to parse dates without time zone info.
	Location *time.Location
}

var dateFieldAliases = map[string][]string{
	fmDate:       {},
	fmLastmod:    {"modified"},
	fmPubDate:    {"pubdate", "published"},
	fmExpiryDate: {"unpublishdate"},
}

// HandleDates updates all the dates given the current configuration and the
// supplied front matter params. Note that this requires all lower-case keys
// in the params map.
func (f FrontMatterHandler) HandleDates(d *FrontMatterDescriptor) error {
	if d.PageConfig == nil {
		panic("missing pageConfig")
	}

	if f.dateHandler == nil {
		panic("missing date handler")
	}

	if _, err := f.dateHandler(d); err != nil {
		return err
	}

	if _, err := f.lastModHandler(d); err != nil {
		return err
	}

	if _, err := f.publishDateHandler(d); err != nil {
		return err
	}

	if _, err := f.expiryDateHandler(d); err != nil {
		return err
	}

	return nil
}

// IsDateKey returns whether the given front matter key is considered a date by the current
// configuration.
func (f FrontMatterHandler) IsDateKey(key string) bool {
	return f.allDateKeys[key]
}

// A Zero date is a signal that the name can not be parsed.
// This follows the format as outlined in Jekyll, https://jekyllrb.com/docs/posts/:
// "Where YEAR is a four-digit number, MONTH and DAY are both two-digit numbers"
func dateAndSlugFromBaseFilename(location *time.Location, name string) (time.Time, string) {
	withoutExt, _ := paths.FileAndExt(name)

	if len(withoutExt) < 10 {
		// This can not be a date.
		return time.Time{}, ""
	}

	d, err := htime.ToTimeInDefaultLocationE(withoutExt[:10], location)
	if err != nil {
		return time.Time{}, ""
	}

	// Be a little lenient with the format here.
	slug := strings.Trim(withoutExt[10:], " -_")

	return d, slug
}

type frontMatterFieldHandler func(d *FrontMatterDescriptor) (bool, error)

func (f FrontMatterHandler) newChainedFrontMatterFieldHandler(handlers ...frontMatterFieldHandler) frontMatterFieldHandler {
	return func(d *FrontMatterDescriptor) (bool, error) {
		for _, h := range handlers {
			// First successful handler wins.
			success, err := h(d)
			if err != nil {
				f.logger.Errorln(err)
			} else if success {
				return true, nil
			}
		}
		return false, nil
	}
}

type FrontmatterConfig struct {
	// Controls how the Date is set from front matter.
	Date []string
	// Controls how the Lastmod is set from front matter.
	Lastmod []string
	// Controls how the PublishDate is set from front matter.
	PublishDate []string
	// Controls how the ExpiryDate is set from front matter.
	ExpiryDate []string
}

const (
	// These are all the date handler identifiers
	// All identifiers not starting with a ":" maps to a front matter parameter.
	fmDate       = "date"
	fmPubDate    = "publishdate"
	fmLastmod    = "lastmod"
	fmExpiryDate = "expirydate"

	// Gets date from filename, e.g 218-02-22-mypage.md
	fmFilename = ":filename"

	// Gets date from file OS mod time.
	fmModTime = ":filemodtime"

	// Gets date from Git
	fmGitAuthorDate = ":git"
)

// This is the config you get when doing nothing.
func newDefaultFrontmatterConfig() FrontmatterConfig {
	return FrontmatterConfig{
		Date:        []string{fmDate, fmPubDate, fmLastmod},
		Lastmod:     []string{fmGitAuthorDate, fmLastmod, fmDate, fmPubDate},
		PublishDate: []string{fmPubDate, fmDate},
		ExpiryDate:  []string{fmExpiryDate},
	}
}

func DecodeFrontMatterConfig(cfg config.Provider) (FrontmatterConfig, error) {
	c := newDefaultFrontmatterConfig()
	defaultConfig := c

	if cfg.IsSet("frontmatter") {
		fm := cfg.GetStringMap("frontmatter")
		for k, v := range fm {
			loki := strings.ToLower(k)
			switch loki {
			case fmDate:
				c.Date = toLowerSlice(v)
			case fmPubDate:
				c.PublishDate = toLowerSlice(v)
			case fmLastmod:
				c.Lastmod = toLowerSlice(v)
			case fmExpiryDate:
				c.ExpiryDate = toLowerSlice(v)
			}
		}
	}

	expander := func(c, d []string) []string {
		out := expandDefaultValues(c, d)
		out = addDateFieldAliases(out)
		return out
	}

	c.Date = expander(c.Date, defaultConfig.Date)
	c.PublishDate = expander(c.PublishDate, defaultConfig.PublishDate)
	c.Lastmod = expander(c.Lastmod, defaultConfig.Lastmod)
	c.ExpiryDate = expander(c.ExpiryDate, defaultConfig.ExpiryDate)

	return c, nil
}

func addDateFieldAliases(values []string) []string {
	var complete []string

	for _, v := range values {
		complete = append(complete, v)
		if aliases, found := dateFieldAliases[v]; found {
			complete = append(complete, aliases...)
		}
	}
	return helpers.UniqueStringsReuse(complete)
}

func expandDefaultValues(values []string, defaults []string) []string {
	var out []string
	for _, v := range values {
		if v == ":default" {
			out = append(out, defaults...)
		} else {
			out = append(out, v)
		}
	}
	return out
}

func toLowerSlice(in any) []string {
	out := cast.ToStringSlice(in)
	for i := 0; i < len(out); i++ {
		out[i] = strings.ToLower(out[i])
	}

	return out
}

// NewFrontmatterHandler creates a new FrontMatterHandler with the given logger and configuration.
// If no logger is provided, one will be created.
func NewFrontmatterHandler(logger loggers.Logger, frontMatterConfig FrontmatterConfig) (FrontMatterHandler, error) {
	if logger == nil {
		logger = loggers.NewDefault()
	}

	allDateKeys := make(map[string]bool)
	addKeys := func(vals []string) {
		for _, k := range vals {
			if !strings.HasPrefix(k, ":") {
				allDateKeys[k] = true
			}
		}
	}

	addKeys(frontMatterConfig.Date)
	addKeys(frontMatterConfig.ExpiryDate)
	addKeys(frontMatterConfig.Lastmod)
	addKeys(frontMatterConfig.PublishDate)

	f := FrontMatterHandler{logger: logger, fmConfig: frontMatterConfig, allDateKeys: allDateKeys}

	if err := f.createHandlers(); err != nil {
		return f, err
	}

	return f, nil
}

func (f *FrontMatterHandler) createHandlers() error {
	var err error

	if f.dateHandler, err = f.createDateHandler(f.fmConfig.Date,
		func(d *FrontMatterDescriptor, t time.Time) {
			d.PageConfig.Date = t
			setParamIfNotSet(fmDate, t, d)
		}); err != nil {
		return err
	}

	if f.lastModHandler, err = f.createDateHandler(f.fmConfig.Lastmod,
		func(d *FrontMatterDescriptor, t time.Time) {
			setParamIfNotSet(fmLastmod, t, d)
			d.PageConfig.Lastmod = t
		}); err != nil {
		return err
	}

	if f.publishDateHandler, err = f.createDateHandler(f.fmConfig.PublishDate,
		func(d *FrontMatterDescriptor, t time.Time) {
			setParamIfNotSet(fmPubDate, t, d)
			d.PageConfig.PublishDate = t
		}); err != nil {
		return err
	}

	if f.expiryDateHandler, err = f.createDateHandler(f.fmConfig.ExpiryDate,
		func(d *FrontMatterDescriptor, t time.Time) {
			setParamIfNotSet(fmExpiryDate, t, d)
			d.PageConfig.ExpiryDate = t
		}); err != nil {
		return err
	}

	return nil
}

func setParamIfNotSet(key string, value any, d *FrontMatterDescriptor) {
	if _, found := d.PageConfig.Params[key]; found {
		return
	}
	d.PageConfig.Params[key] = value
}

func (f FrontMatterHandler) createDateHandler(identifiers []string, setter func(d *FrontMatterDescriptor, t time.Time)) (frontMatterFieldHandler, error) {
	var h *frontmatterFieldHandlers
	var handlers []frontMatterFieldHandler

	for _, identifier := range identifiers {
		switch identifier {
		case fmFilename:
			handlers = append(handlers, h.newDateFilenameHandler(setter))
		case fmModTime:
			handlers = append(handlers, h.newDateModTimeHandler(setter))
		case fmGitAuthorDate:
			handlers = append(handlers, h.newDateGitAuthorDateHandler(setter))
		default:
			handlers = append(handlers, h.newDateFieldHandler(identifier, setter))
		}
	}

	return f.newChainedFrontMatterFieldHandler(handlers...), nil
}

type frontmatterFieldHandlers int

func (f *frontmatterFieldHandlers) newDateFieldHandler(key string, setter func(d *FrontMatterDescriptor, t time.Time)) frontMatterFieldHandler {
	return func(d *FrontMatterDescriptor) (bool, error) {
		v, found := d.PageConfig.Params[key]

		if !found {
			return false, nil
		}

		date, err := htime.ToTimeInDefaultLocationE(v, d.Location)
		if err != nil {
			return false, nil
		}

		// We map several date keys to one, so, for example,
		// "expirydate", "unpublishdate" will all set .ExpiryDate (first found).
		setter(d, date)

		// This is the params key as set in front matter.
		d.PageConfig.Params[key] = date

		return true, nil
	}
}

func (f *frontmatterFieldHandlers) newDateFilenameHandler(setter func(d *FrontMatterDescriptor, t time.Time)) frontMatterFieldHandler {
	return func(d *FrontMatterDescriptor) (bool, error) {
		date, slug := dateAndSlugFromBaseFilename(d.Location, d.BaseFilename)
		if date.IsZero() {
			return false, nil
		}

		setter(d, date)

		if _, found := d.PageConfig.Params["slug"]; !found {
			// Use slug from filename
			d.PageConfig.Slug = slug
		}

		return true, nil
	}
}

func (f *frontmatterFieldHandlers) newDateModTimeHandler(setter func(d *FrontMatterDescriptor, t time.Time)) frontMatterFieldHandler {
	return func(d *FrontMatterDescriptor) (bool, error) {
		if d.ModTime.IsZero() {
			return false, nil
		}
		setter(d, d.ModTime)
		return true, nil
	}
}

func (f *frontmatterFieldHandlers) newDateGitAuthorDateHandler(setter func(d *FrontMatterDescriptor, t time.Time)) frontMatterFieldHandler {
	return func(d *FrontMatterDescriptor) (bool, error) {
		if d.GitAuthorDate.IsZero() {
			return false, nil
		}
		setter(d, d.GitAuthorDate)
		return true, nil
	}
}
