// Copyright 2018 The Hugo Authors. All rights reserved.
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

package hugolib

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gohugoio/hugo/helpers"

	"github.com/gohugoio/hugo/config"
	"github.com/spf13/cast"
	jww "github.com/spf13/jwalterweatherman"
)

// TODO(bep) should probably make the date handling chain complete to give people the flexibility they want.

type frontmatterHandler struct {
	// Ordered chain.
	dateHandlers frontMatterFieldHandler

	logger *jww.Notepad
}

type frontMatterDescriptor struct {

	// This the Page's front matter.
	frontmatter map[string]interface{}

	// This is the Page's base filename, e.g. page.md.
	baseFilename string

	// The content file's mod time.
	modTime time.Time

	// The below are pointers to values on Page and will be updated.

	// This is the Page's params.
	params map[string]interface{}

	// This is the Page's dates.
	dates *PageDates
}

func (f frontmatterHandler) handleDate(d frontMatterDescriptor) (time.Time, error) {
	v, err := f.dateHandlers(d)
	if err == nil {
		if d, ok := v.(time.Time); ok {
			return d, nil
		}
	}

	return time.Time{}, err
}

var (
	lastModFrontMatterKeys     = []string{"lastmod", "modified"}
	publishDateFrontMatterKeys = []string{"publishdate", "pubdate", "published"}
	expiryDateFrontMatterKeys  = []string{"expirydate", "unpublishdate"}
	allDateFrontMatterKeys     = make(map[string]bool)
)

func init() {
	for _, key := range lastModFrontMatterKeys {
		allDateFrontMatterKeys[key] = true
	}
	for _, key := range publishDateFrontMatterKeys {
		allDateFrontMatterKeys[key] = true
	}
	for _, key := range expiryDateFrontMatterKeys {
		allDateFrontMatterKeys[key] = true
	}

	allDateFrontMatterKeys["date"] = true
}

func (f frontmatterHandler) handleDates(d frontMatterDescriptor) error {

	date, err := f.handleDate(d)
	if err != nil {
		return err
	}
	d.dates.Date = date
	d.dates.Lastmod = f.setParamsAndReturnFirstDate(d, lastModFrontMatterKeys)
	d.dates.PublishDate = f.setParamsAndReturnFirstDate(d, publishDateFrontMatterKeys)
	d.dates.ExpiryDate = f.setParamsAndReturnFirstDate(d, expiryDateFrontMatterKeys)

	// Hugo really needs a date!
	if d.dates.Date.IsZero() {
		d.dates.Date = d.dates.PublishDate
	}

	if d.dates.PublishDate.IsZero() {
		d.dates.PublishDate = d.dates.Date
	}

	if d.dates.Lastmod.IsZero() {
		d.dates.Lastmod = d.dates.Date
	}

	f.setParamIfNotZero("date", d.params, d.dates.Date)
	f.setParamIfNotZero("lastmod", d.params, d.dates.Lastmod)
	f.setParamIfNotZero("publishdate", d.params, d.dates.PublishDate)
	f.setParamIfNotZero("expirydate", d.params, d.dates.ExpiryDate)

	return nil
}

func (f frontmatterHandler) isDateKey(key string) bool {
	return allDateFrontMatterKeys[key]
}

func (f frontmatterHandler) setParamIfNotZero(name string, params map[string]interface{}, date time.Time) {
	if date.IsZero() {
		return
	}
	params[name] = date
}

func (f frontmatterHandler) setParamsAndReturnFirstDate(d frontMatterDescriptor, keys []string) time.Time {
	var date time.Time

	for _, key := range keys {
		v, found := d.frontmatter[key]
		if found {
			currentDate, err := cast.ToTimeE(v)
			if err == nil {
				d.params[key] = currentDate
				if date.IsZero() {
					date = currentDate
				}
			} else {
				d.params[key] = v
			}
		}
	}

	return date
}

// A Zero date is a signal that the name can not be parsed.
// This follows the format as outlined in Jekyll, https://jekyllrb.com/docs/posts/:
// "Where YEAR is a four-digit number, MONTH and DAY are both two-digit numbers"
func (f frontmatterHandler) dateAndSlugFromBaseFilename(name string) (time.Time, string) {
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

type frontMatterFieldHandler func(d frontMatterDescriptor) (interface{}, error)

func (f frontmatterHandler) newChainedFrontMatterFieldHandler(handlers ...frontMatterFieldHandler) frontMatterFieldHandler {
	return func(d frontMatterDescriptor) (interface{}, error) {
		for _, h := range handlers {
			// First non-nil value wins.
			val, err := h(d)
			if err != nil {
				f.logger.ERROR.Println(err)
			} else if val != nil {
				return val, nil
			}
		}
		return nil, nil
	}

}

func newFrontmatterConfig(logger *jww.Notepad, cfg config.Provider) (frontmatterHandler, error) {

	if logger == nil {
		logger = jww.NewNotepad(jww.LevelWarn, jww.LevelWarn, os.Stdout, ioutil.Discard, "", log.Ldate|log.Ltime)
	}

	f := frontmatterHandler{logger: logger}

	handlers := &frontmatterFieldHandlers{logger: logger}

	dateHandlers := []frontMatterFieldHandler{handlers.defaultDateHandler}

	defaultDate := cfg.Get("frontmatter.defaultdate")

	if defaultDate != nil {
		slice, err := cast.ToStringSliceE(defaultDate)
		if err != nil {
			return f, fmt.Errorf("invalid value for defaultDate, expeced a string slice, got %T", defaultDate)
		}

		for _, v := range slice {
			if strings.EqualFold(v, "filename") {
				dateHandlers = append(dateHandlers, handlers.defaultDateFilenameHandler)
				// No more for now.
				break
			}
		}
	}

	// This is deprecated
	if cfg.GetBool("useModTimeAsFallback") {
		dateHandlers = append(dateHandlers, handlers.defaultDateModTimeHandler)
	}

	f.dateHandlers = f.newChainedFrontMatterFieldHandler(dateHandlers...)

	return f, nil
}

type frontmatterFieldHandlers struct {
	logger *jww.Notepad
}

func (f *frontmatterFieldHandlers) defaultDateHandler(d frontMatterDescriptor) (interface{}, error) {
	v, found := d.frontmatter["date"]
	if !found {
		return nil, nil
	}

	date, err := cast.ToTimeE(v)
	if err != nil {
		return nil, err
	}

	return date, nil
}

func (f *frontmatterFieldHandlers) defaultDateFilenameHandler(d frontMatterDescriptor) (interface{}, error) {
	return true, nil
}

func (f *frontmatterFieldHandlers) defaultDateModTimeHandler(d frontMatterDescriptor) (interface{}, error) {
	if !d.modTime.IsZero() {
		return d.modTime, nil
	}
	return nil, nil
}
