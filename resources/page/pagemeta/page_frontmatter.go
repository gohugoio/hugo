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

package pagemeta

import (
	"strings"
	"time"

	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/resources/resource"

	"github.com/gohugoio/hugo/config"
	"github.com/spf13/cast"
)

// FrontMatterHandler maps front matter into Page fields and .Params.
// Note that we currently have only extracted the date logic.
type FrontMatterHandler struct {
	fmConfig frontmatterConfig

	dateHandler        frontMatterFieldHandler
	lastModHandler     frontMatterFieldHandler
	publishDateHandler frontMatterFieldHandler
	expiryDateHandler  frontMatterFieldHandler

	// A map of all date keys configured, including any custom.
	allDateKeys map[string]bool

	logger *loggers.Logger
}

// FrontMatterDescriptor describes how to handle front matter for a given Page.
// It has pointers to values in the receiving page which gets updated.
type FrontMatterDescriptor struct {

	// This the Page's front matter.
	Frontmatter map[string]interface{}

	// This is the Page's base filename (BaseFilename), e.g. page.md., or
	// if page is a leaf bundle, the bundle folder name (ContentBaseName).
	BaseFilename string

	// The content file's mod time.
	ModTime time.Time

	// May be set from the author date in Git.
	GitAuthorDate time.Time

	// The below are pointers to values on Page and will be modified.

	// This is the Page's params.
	Params map[string]interface{}

	// This is the Page's dates.
	Dates *resource.Dates

	// This is the Page's Slug etc.
	PageURLs *URLPath
}

var (
	dateFieldAliases = map[string][]string{
		fmDate:       {},
		fmLastmod:    {"modified"},
		fmPubDate:    {"pubdate", "published"},
		fmExpiryDate: {"unpublishdate"},
	}
)

// HandleDates updates all the dates given the current configuration and the
// supplied front matter params. Note that this requires all lower-case keys
// in the params map.
func (f FrontMatterHandler) HandleDates(d *FrontMatterDescriptor) error {
	if d.Dates == nil {
		panic("missing dates")
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
func dateAndSlugFromBaseFilename(name string) (time.Time, string) {
	withoutExt, _ := helpers.FileAndExt(name)

	if len(withoutExt) < 10 {
		// This can not be a date.
		return time.Time{}, ""
	}

	// Note: Hugo currently have no custom timezone support.
	// We will have to revisit this when that is in place.
	d, err := time.Parse("2006-01-02", withoutExt[:10])
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
				f.logger.ERROR.Println(err)
			} else if success {
				return true, nil
			}
		}
		return false, nil
	}
}

type frontmatterConfig struct {
	date        []string
	lastmod     []string
	publishDate []string
	expiryDate  []string
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
func newDefaultFrontmatterConfig() frontmatterConfig {
	return frontmatterConfig{
		date:        []string{fmDate, fmPubDate, fmLastmod},
		lastmod:     []string{fmGitAuthorDate, fmLastmod, fmDate, fmPubDate},
		publishDate: []string{fmPubDate, fmDate},
		expiryDate:  []string{fmExpiryDate},
	}
}

func newFrontmatterConfig(cfg config.Provider) (frontmatterConfig, error) {
	c := newDefaultFrontmatterConfig()
	defaultConfig := c

	if cfg.IsSet("frontmatter") {
		fm := cfg.GetStringMap("frontmatter")
		for k, v := range fm {
			loki := strings.ToLower(k)
			switch loki {
			case fmDate:
				c.date = toLowerSlice(v)
			case fmPubDate:
				c.publishDate = toLowerSlice(v)
			case fmLastmod:
				c.lastmod = toLowerSlice(v)
			case fmExpiryDate:
				c.expiryDate = toLowerSlice(v)
			}
		}
	}

	expander := func(c, d []string) []string {
		out := expandDefaultValues(c, d)
		out = addDateFieldAliases(out)
		return out
	}

	c.date = expander(c.date, defaultConfig.date)
	c.publishDate = expander(c.publishDate, defaultConfig.publishDate)
	c.lastmod = expander(c.lastmod, defaultConfig.lastmod)
	c.expiryDate = expander(c.expiryDate, defaultConfig.expiryDate)

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

func toLowerSlice(in interface{}) []string {
	out := cast.ToStringSlice(in)
	for i := 0; i < len(out); i++ {
		out[i] = strings.ToLower(out[i])
	}

	return out
}

// NewFrontmatterHandler creates a new FrontMatterHandler with the given logger and configuration.
// If no logger is provided, one will be created.
func NewFrontmatterHandler(logger *loggers.Logger, cfg config.Provider) (FrontMatterHandler, error) {

	if logger == nil {
		logger = loggers.NewErrorLogger()
	}

	frontMatterConfig, err := newFrontmatterConfig(cfg)
	if err != nil {
		return FrontMatterHandler{}, err
	}

	allDateKeys := make(map[string]bool)
	addKeys := func(vals []string) {
		for _, k := range vals {
			if !strings.HasPrefix(k, ":") {
				allDateKeys[k] = true
			}
		}
	}

	addKeys(frontMatterConfig.date)
	addKeys(frontMatterConfig.expiryDate)
	addKeys(frontMatterConfig.lastmod)
	addKeys(frontMatterConfig.publishDate)

	f := FrontMatterHandler{logger: logger, fmConfig: frontMatterConfig, allDateKeys: allDateKeys}

	if err := f.createHandlers(); err != nil {
		return f, err
	}

	return f, nil
}

func (f *FrontMatterHandler) createHandlers() error {
	var err error

	if f.dateHandler, err = f.createDateHandler(f.fmConfig.date,
		func(d *FrontMatterDescriptor, t time.Time) {
			d.Dates.FDate = t
			setParamIfNotSet(fmDate, t, d)
		}); err != nil {
		return err
	}

	if f.lastModHandler, err = f.createDateHandler(f.fmConfig.lastmod,
		func(d *FrontMatterDescriptor, t time.Time) {
			setParamIfNotSet(fmLastmod, t, d)
			d.Dates.FLastmod = t
		}); err != nil {
		return err
	}

	if f.publishDateHandler, err = f.createDateHandler(f.fmConfig.publishDate,
		func(d *FrontMatterDescriptor, t time.Time) {
			setParamIfNotSet(fmPubDate, t, d)
			d.Dates.FPublishDate = t
		}); err != nil {
		return err
	}

	if f.expiryDateHandler, err = f.createDateHandler(f.fmConfig.expiryDate,
		func(d *FrontMatterDescriptor, t time.Time) {
			setParamIfNotSet(fmExpiryDate, t, d)
			d.Dates.FExpiryDate = t
		}); err != nil {
		return err
	}

	return nil
}

func setParamIfNotSet(key string, value interface{}, d *FrontMatterDescriptor) {
	if _, found := d.Params[key]; found {
		return
	}
	d.Params[key] = value
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
		v, found := d.Frontmatter[key]

		if !found {
			return false, nil
		}

		date, err := cast.ToTimeE(v)
		if err != nil {
			return false, nil
		}

		// We map several date keys to one, so, for example,
		// "expirydate", "unpublishdate" will all set .ExpiryDate (first found).
		setter(d, date)

		// This is the params key as set in front matter.
		d.Params[key] = date

		return true, nil
	}
}

func (f *frontmatterFieldHandlers) newDateFilenameHandler(setter func(d *FrontMatterDescriptor, t time.Time)) frontMatterFieldHandler {
	return func(d *FrontMatterDescriptor) (bool, error) {
		date, slug := dateAndSlugFromBaseFilename(d.BaseFilename)
		if date.IsZero() {
			return false, nil
		}

		setter(d, date)

		if _, found := d.Frontmatter["slug"]; !found {
			// Use slug from filename
			d.PageURLs.Slug = slug
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
