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

package pagemeta

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"

	"github.com/gohugoio/hugo/helpers"

	"github.com/gohugoio/hugo/config"
	"github.com/spf13/cast"
	jww "github.com/spf13/jwalterweatherman"
)

// TODO(bep) should probably make the date handling chain complete to give people the flexibility they want.

type FrontmatterHandler struct {
	// Ordered chain.
	dateHandlers frontMatterFieldHandler

	logger *jww.Notepad
}

type FrontMatterDescriptor struct {

	// This the Page's front matter.
	Frontmatter map[string]interface{}

	// This is the Page's base filename, e.g. page.md.
	BaseFilename string

	// The content file's mod time.
	ModTime time.Time

	// The below are pointers to values on Page and will be updated.

	// This is the Page's params.
	Params map[string]interface{}

	// This is the Page's dates.
	Dates *PageDates

	// This is the Page's Slug etc.
	PageURLs *URLPath
}

func (f FrontmatterHandler) handleDate(d FrontMatterDescriptor) error {
	_, err := f.dateHandlers(d)
	return err
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

func (f FrontmatterHandler) HandleDates(d FrontMatterDescriptor) error {
	if d.Dates == nil {
		panic("missing dates")
	}

	err := f.handleDate(d)
	if err != nil {
		return err
	}
	d.Dates.Lastmod = f.setParamsAndReturnFirstDate(d, lastModFrontMatterKeys)
	d.Dates.PublishDate = f.setParamsAndReturnFirstDate(d, publishDateFrontMatterKeys)
	d.Dates.ExpiryDate = f.setParamsAndReturnFirstDate(d, expiryDateFrontMatterKeys)

	// Hugo really needs a date!
	if d.Dates.Date.IsZero() {
		d.Dates.Date = d.Dates.PublishDate
	}

	if d.Dates.Lastmod.IsZero() {
		d.Dates.Lastmod = d.Dates.Date
	}

	// TODO(bep) date decide vs https://github.com/gohugoio/hugo/issues/3977
	if d.Dates.PublishDate.IsZero() {
		//d.dates.PublishDate = d.dates.Date
	}

	if d.Dates.Date.IsZero() {
		d.Dates.Date = d.Dates.Lastmod
	}

	f.setParamIfNotZero("date", d.Params, d.Dates.Date)
	f.setParamIfNotZero("lastmod", d.Params, d.Dates.Lastmod)
	f.setParamIfNotZero("publishdate", d.Params, d.Dates.PublishDate)
	f.setParamIfNotZero("expirydate", d.Params, d.Dates.ExpiryDate)

	return nil
}

func (f FrontmatterHandler) IsDateKey(key string) bool {
	return allDateFrontMatterKeys[key]
}

func (f FrontmatterHandler) setParamIfNotZero(name string, params map[string]interface{}, date time.Time) {
	if date.IsZero() {
		return
	}
	params[name] = date
}

func (f FrontmatterHandler) setParamsAndReturnFirstDate(d FrontMatterDescriptor, keys []string) time.Time {
	var date time.Time

	for _, key := range keys {
		v, found := d.Frontmatter[key]
		if found {
			currentDate, err := cast.ToTimeE(v)
			if err == nil {
				d.Params[key] = currentDate
				if date.IsZero() {
					date = currentDate
				}
			} else {
				d.Params[key] = v
			}
		}
	}

	return date
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

type frontMatterFieldHandler func(d FrontMatterDescriptor) (bool, error)

func (f FrontmatterHandler) newChainedFrontMatterFieldHandler(handlers ...frontMatterFieldHandler) frontMatterFieldHandler {
	return func(d FrontMatterDescriptor) (bool, error) {
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
	lastMod     []string
	publishDate []string
	expiryDate  []string
}

const (
	fmDate       = "date"
	fmPubDate    = "publishdate"
	fmLastMod    = "lastmod"
	fmExpiryDate = "expirydate"
)

func newDefaultFrontmatterConfig() frontmatterConfig {
	return frontmatterConfig{
		date:        []string{fmDate, fmPubDate, fmLastMod},
		lastMod:     []string{fmLastMod, fmDate},
		publishDate: []string{fmPubDate, fmDate},
		expiryDate:  []string{fmExpiryDate},
	}
}

func newFrontmatterConfig(cfg config.Provider) (frontmatterConfig, error) {
	c := newDefaultFrontmatterConfig()

	if cfg.IsSet("frontmatter") {
		fm := cfg.GetStringMap("frontmatter")
		if fm != nil {
			for k, v := range fm {
				loki := strings.ToLower(k)
				switch loki {
				case fmDate:
					c.date = toLowerSlice(v)
				case fmPubDate:
					c.publishDate = toLowerSlice(v)
				case fmLastMod:
					c.lastMod = toLowerSlice(v)
				case fmExpiryDate:
					c.expiryDate = toLowerSlice(v)
				}
			}
		}
		err := mapstructure.WeakDecode(fm, &c)
		return c, err
	}
	return c, nil
}

func toLowerSlice(in interface{}) []string {
	out := cast.ToStringSlice(in)
	for i := 0; i < len(out); i++ {
		out[i] = strings.ToLower(out[i])
	}

	return out
}

func NewFrontmatterHandler(logger *jww.Notepad, cfg config.Provider) (FrontmatterHandler, error) {

	if logger == nil {
		logger = jww.NewNotepad(jww.LevelWarn, jww.LevelWarn, os.Stdout, ioutil.Discard, "", log.Ldate|log.Ltime)
	}

	f := FrontmatterHandler{logger: logger}

	handlers := &frontmatterFieldHandlers{}

	/*

		[frontmatter]
		date = ["date", "publishDate", "lastMod"]
		lastMod = ["lastMod", "date"]
		publishDate  = ["publishDate", "date"]
		expiryDate = ["expiryDate"]

	*/

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
}

func (f *frontmatterFieldHandlers) defaultDateHandler(d FrontMatterDescriptor) (bool, error) {
	v, found := d.Frontmatter["date"]
	if !found {
		return false, nil
	}

	date, err := cast.ToTimeE(v)
	if err != nil {
		return false, err
	}

	d.Dates.Date = date

	return true, nil
}

func (f *frontmatterFieldHandlers) defaultDateFilenameHandler(d FrontMatterDescriptor) (bool, error) {
	date, slug := dateAndSlugFromBaseFilename(d.BaseFilename)
	if date.IsZero() {
		return false, nil
	}

	d.Dates.Date = date

	if _, found := d.Frontmatter["slug"]; !found {
		// Use slug from filename
		d.PageURLs.Slug = slug
	}

	return true, nil
}

func (f *frontmatterFieldHandlers) defaultDateModTimeHandler(d FrontMatterDescriptor) (bool, error) {
	if !d.ModTime.IsZero() {
		d.Dates.Date = d.ModTime
		return true, nil
	}
	return false, nil
}
