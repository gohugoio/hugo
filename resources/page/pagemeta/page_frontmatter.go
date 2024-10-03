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
	"errors"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/gohugoio/hugo/common/hreflect"
	"github.com/gohugoio/hugo/common/htime"
	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/hugofs/files"
	"github.com/gohugoio/hugo/markup"
	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/resources/kinds"
	"github.com/gohugoio/hugo/resources/page"
	"github.com/gohugoio/hugo/resources/resource"

	"github.com/gohugoio/hugo/helpers"

	"github.com/gohugoio/hugo/config"
	"github.com/spf13/cast"
)

type DatesStrings struct {
	Date        string `json:"date"`
	Lastmod     string `json:"lastMod"`
	PublishDate string `json:"publishDate"`
	ExpiryDate  string `json:"expiryDate"`
}

type Dates struct {
	Date        time.Time
	Lastmod     time.Time
	PublishDate time.Time
	ExpiryDate  time.Time
}

func (d Dates) IsDateOrLastModAfter(in Dates) bool {
	return d.Date.After(in.Date) || d.Lastmod.After(in.Lastmod)
}

func (d *Dates) UpdateDateAndLastmodAndPublishDateIfAfter(in Dates) {
	if in.Date.After(d.Date) {
		d.Date = in.Date
	}
	if in.Lastmod.After(d.Lastmod) {
		d.Lastmod = in.Lastmod
	}

	if in.PublishDate.After(d.PublishDate) && in.PublishDate.Before(htime.Now()) {
		d.PublishDate = in.PublishDate
	}
}

func (d Dates) IsAllDatesZero() bool {
	return d.Date.IsZero() && d.Lastmod.IsZero() && d.PublishDate.IsZero() && d.ExpiryDate.IsZero()
}

// PageConfig configures a Page, typically from front matter.
// Note that all the top level fields are reserved Hugo keywords.
// Any custom configuration needs to be set in the Params map.
type PageConfig struct {
	Dates Dates `json:"-"` // Dates holds the four core dates for this page.
	DatesStrings
	Title          string   // The title of the page.
	LinkTitle      string   // The link title of the page.
	Type           string   // The content type of the page.
	Layout         string   // The layout to use for to render this page.
	Weight         int      // The weight of the page, used in sorting if set to a non-zero value.
	Kind           string   // The kind of page, e.g. "page", "section", "home" etc. This is usually derived from the content path.
	Path           string   // The canonical path to the page, e.g. /sect/mypage. Note: Leading slash, no trailing slash, no extensions or language identifiers.
	Lang           string   // The language code for this page. This is usually derived from the module mount or filename.
	URL            string   // The URL to the rendered page, e.g. /sect/mypage.html.
	Slug           string   // The slug for this page.
	Description    string   // The description for this page.
	Summary        string   // The summary for this page.
	Draft          bool     // Whether or not the content is a draft.
	Headless       bool     `json:"-"` // Whether or not the page should be rendered.
	IsCJKLanguage  bool     // Whether or not the content is in a CJK language.
	TranslationKey string   // The translation key for this page.
	Keywords       []string // The keywords for this page.
	Aliases        []string // The aliases for this page.
	Outputs        []string // The output formats to render this page in. If not set, the site's configured output formats for this page kind will be used.

	FrontMatterOnlyValues `mapstructure:"-" json:"-"`

	Cascade []map[string]any
	Sitemap config.SitemapConfig
	Build   BuildConfig
	Menus   []string

	// User defined params.
	Params maps.Params

	// Content holds the content for this page.
	Content Source

	// Compiled values.
	CascadeCompiled      map[page.PageMatcher]maps.Params
	ContentMediaType     media.Type `mapstructure:"-" json:"-"`
	IsFromContentAdapter bool       `mapstructure:"-" json:"-"`
}

var DefaultPageConfig = PageConfig{
	Build: DefaultBuildConfig,
}

func (p *PageConfig) Validate(pagesFromData bool) error {
	if pagesFromData {
		if p.Path == "" {
			return errors.New("path must be set")
		}
		if strings.HasPrefix(p.Path, "/") {
			return fmt.Errorf("path %q must not start with a /", p.Path)
		}
		if p.Lang != "" {
			return errors.New("lang must not be set")
		}

		if p.Content.Markup != "" {
			return errors.New("markup must not be set, use mediaType")
		}
	}

	if p.Cascade != nil {
		if !kinds.IsBranch(p.Kind) {
			return errors.New("cascade is only supported for branch nodes")
		}
	}

	return nil
}

// Compile sets up the page configuration after all fields have been set.
func (p *PageConfig) Compile(basePath string, pagesFromData bool, ext string, logger loggers.Logger, mediaTypes media.Types) error {
	// In content adapters, we always get relative paths.
	if basePath != "" {
		p.Path = path.Join(basePath, p.Path)
	}

	if p.Params == nil {
		p.Params = make(maps.Params)
	}
	maps.PrepareParams(p.Params)

	if p.Content.Markup == "" && p.Content.MediaType == "" {
		if ext == "" {
			ext = "md"
		}
		p.ContentMediaType = MarkupToMediaType(ext, mediaTypes)
		if p.ContentMediaType.IsZero() {
			return fmt.Errorf("failed to resolve media type for suffix %q", ext)
		}
	}

	var s string
	if p.ContentMediaType.IsZero() {
		if p.Content.MediaType != "" {
			s = p.Content.MediaType
			p.ContentMediaType, _ = mediaTypes.GetByType(s)
		}

		if p.ContentMediaType.IsZero() && p.Content.Markup != "" {
			s = p.Content.Markup
			p.ContentMediaType = MarkupToMediaType(s, mediaTypes)
		}
	}

	if p.ContentMediaType.IsZero() {
		return fmt.Errorf("failed to resolve media type for %q", s)
	}

	if p.Content.Markup == "" {
		p.Content.Markup = p.ContentMediaType.SubType
	}

	if pagesFromData {
		if p.Kind == "" {
			p.Kind = kinds.KindPage
		}

		// Note that NormalizePathStringBasic will make sure that we don't preserve the unnormalized path.
		// We do that when we create pages from the file system; mostly for backward compatibility,
		// but also because people tend to use use the filename to name their resources (with spaces and all),
		// and this isn't relevant when creating resources from an API where it's easy to add textual meta data.
		p.Path = paths.NormalizePathStringBasic(p.Path)
	}

	if p.Cascade != nil {
		cascade, err := page.DecodeCascade(logger, p.Cascade)
		if err != nil {
			return fmt.Errorf("failed to decode cascade: %w", err)
		}
		p.CascadeCompiled = cascade
	}

	return nil
}

// MarkupToMediaType converts a markup string to a media type.
func MarkupToMediaType(s string, mediaTypes media.Types) media.Type {
	s = strings.ToLower(s)
	mt, _ := mediaTypes.GetBestMatch(markup.ResolveMarkup(s))
	return mt
}

type ResourceConfig struct {
	Path    string
	Name    string
	Title   string
	Params  maps.Params
	Content Source

	// Compiled values.
	PathInfo         *paths.Path `mapstructure:"-" json:"-"`
	ContentMediaType media.Type
}

func (rc *ResourceConfig) Validate() error {
	if rc.Path == "" {
		return errors.New("path must be set")
	}
	if rc.Content.Markup != "" {
		return errors.New("markup must not be set, use mediaType")
	}
	return nil
}

func (rc *ResourceConfig) Compile(basePath string, pathParser *paths.PathParser, mediaTypes media.Types) error {
	if rc.Params != nil {
		maps.PrepareParams(rc.Params)
	}

	// Note that NormalizePathStringBasic will make sure that we don't preserve the unnormalized path.
	// We do that when we create resources from the file system; mostly for backward compatibility,
	// but also because people tend to use use the filename to name their resources (with spaces and all),
	// and this isn't relevant when creating resources from an API where it's easy to add textual meta data.
	rc.Path = paths.NormalizePathStringBasic(path.Join(basePath, rc.Path))
	rc.PathInfo = pathParser.Parse(files.ComponentFolderContent, rc.Path)
	if rc.Content.MediaType != "" {
		var found bool
		rc.ContentMediaType, found = mediaTypes.GetByType(rc.Content.MediaType)
		if !found {
			return fmt.Errorf("media type %q not found", rc.Content.MediaType)
		}
	}
	return nil
}

type Source struct {
	// MediaType is the media type of the content.
	MediaType string

	// The markup used in Value. Only used in front matter.
	Markup string

	// The content.
	Value any
}

func (s Source) IsZero() bool {
	return !hreflect.IsTruthful(s.Value)
}

func (s Source) IsResourceValue() bool {
	_, ok := s.Value.(resource.Resource)
	return ok
}

func (s Source) ValueAsString() string {
	if s.Value == nil {
		return ""
	}
	ss, err := cast.ToStringE(s.Value)
	if err != nil {
		panic(fmt.Errorf("content source: failed to convert %T to string: %s", s.Value, err))
	}
	return ss
}

func (s Source) ValueAsOpenReadSeekCloser() hugio.OpenReadSeekCloser {
	return hugio.NewOpenReadSeekCloser(hugio.NewReadSeekerNoOpCloserFromString(s.ValueAsString()))
}

// FrontMatterOnlyValues holds values that can only be set via front matter.
type FrontMatterOnlyValues struct {
	ResourcesMeta []map[string]any
}

// FrontMatterHandler maps front matter into Page fields and .Params.
// Note that we currently have only extracted the date logic.
type FrontMatterHandler struct {
	fmConfig FrontmatterConfig

	contentAdapterDatesHandler func(d *FrontMatterDescriptor) error

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

	// The Page's path if the page is backed by a file, else its title.
	PathOrTitle string

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

	if d.PageConfig.IsFromContentAdapter {
		if f.contentAdapterDatesHandler == nil {
			panic("missing content adapter date handler")
		}
		return f.contentAdapterDatesHandler(d)
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

	if f.contentAdapterDatesHandler, err = f.createContentAdapterDatesHandler(f.fmConfig); err != nil {
		return err
	}

	if f.dateHandler, err = f.createDateHandler(f.fmConfig.Date,
		func(d *FrontMatterDescriptor, t time.Time) {
			d.PageConfig.Dates.Date = t
			setParamIfNotSet(fmDate, t, d)
		}); err != nil {
		return err
	}

	if f.lastModHandler, err = f.createDateHandler(f.fmConfig.Lastmod,
		func(d *FrontMatterDescriptor, t time.Time) {
			setParamIfNotSet(fmLastmod, t, d)
			d.PageConfig.Dates.Lastmod = t
		}); err != nil {
		return err
	}

	if f.publishDateHandler, err = f.createDateHandler(f.fmConfig.PublishDate,
		func(d *FrontMatterDescriptor, t time.Time) {
			setParamIfNotSet(fmPubDate, t, d)
			d.PageConfig.Dates.PublishDate = t
		}); err != nil {
		return err
	}

	if f.expiryDateHandler, err = f.createDateHandler(f.fmConfig.ExpiryDate,
		func(d *FrontMatterDescriptor, t time.Time) {
			setParamIfNotSet(fmExpiryDate, t, d)
			d.PageConfig.Dates.ExpiryDate = t
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

func (f FrontMatterHandler) createContentAdapterDatesHandler(fmcfg FrontmatterConfig) (func(d *FrontMatterDescriptor) error, error) {
	setTime := func(key string, value time.Time, in *PageConfig) {
		switch key {
		case fmDate:
			in.Dates.Date = value
		case fmLastmod:
			in.Dates.Lastmod = value
		case fmPubDate:
			in.Dates.PublishDate = value
		case fmExpiryDate:
			in.Dates.ExpiryDate = value
		}
	}

	getTime := func(key string, in *PageConfig) time.Time {
		switch key {
		case fmDate:
			return in.Dates.Date
		case fmLastmod:
			return in.Dates.Lastmod
		case fmPubDate:
			return in.Dates.PublishDate
		case fmExpiryDate:
			return in.Dates.ExpiryDate
		}
		return time.Time{}
	}

	createSetter := func(identifiers []string, date string) func(pcfg *PageConfig) {
		var getTimes []func(in *PageConfig) time.Time
		for _, identifier := range identifiers {
			if strings.HasPrefix(identifier, ":") {
				continue
			}
			switch identifier {
			case fmDate:
				getTimes = append(getTimes, func(in *PageConfig) time.Time {
					return getTime(fmDate, in)
				})
			case fmLastmod:
				getTimes = append(getTimes, func(in *PageConfig) time.Time {
					return getTime(fmLastmod, in)
				})
			case fmPubDate:
				getTimes = append(getTimes, func(in *PageConfig) time.Time {
					return getTime(fmPubDate, in)
				})
			case fmExpiryDate:
				getTimes = append(getTimes, func(in *PageConfig) time.Time {
					return getTime(fmExpiryDate, in)
				})
			}
		}

		return func(pcfg *PageConfig) {
			for _, get := range getTimes {
				if t := get(pcfg); !t.IsZero() {
					setTime(date, t, pcfg)
					return
				}
			}
		}
	}

	setDate := createSetter(fmcfg.Date, fmDate)
	setLastmod := createSetter(fmcfg.Lastmod, fmLastmod)
	setPublishDate := createSetter(fmcfg.PublishDate, fmPubDate)
	setExpiryDate := createSetter(fmcfg.ExpiryDate, fmExpiryDate)

	fn := func(d *FrontMatterDescriptor) error {
		pcfg := d.PageConfig
		setDate(pcfg)
		setLastmod(pcfg)
		setPublishDate(pcfg)
		setExpiryDate(pcfg)
		return nil
	}
	return fn, nil
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

		var date time.Time
		if vt, ok := v.(time.Time); ok && vt.Location() == d.Location {
			date = vt
		} else {
			var err error
			date, err = htime.ToTimeInDefaultLocationE(v, d.Location)
			if err != nil {
				return false, fmt.Errorf("invalid front matter: %s: %s: see %s", key, v, d.PathOrTitle)
			}
			d.PageConfig.Params[key] = date
		}

		// We map several date keys to one, so, for example,
		// "expirydate", "unpublishdate" will all set .ExpiryDate (first found).
		setter(d, date)

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
