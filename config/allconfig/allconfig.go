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

// Package allconfig contains the full configuration for Hugo.
// <docsmeta>{ "name": "Configuration", "description": "This section holds all configuration options in Hugo." }</docsmeta>
package allconfig

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gohugoio/hugo/cache/filecache"
	"github.com/gohugoio/hugo/cache/httpcache"
	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/common/types"
	"github.com/gohugoio/hugo/common/urls"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/config/privacy"
	"github.com/gohugoio/hugo/config/security"
	"github.com/gohugoio/hugo/config/services"
	"github.com/gohugoio/hugo/deploy/deployconfig"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugolib/segments"
	"github.com/gohugoio/hugo/langs"
	"github.com/gohugoio/hugo/markup/markup_config"
	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/minifiers"
	"github.com/gohugoio/hugo/modules"
	"github.com/gohugoio/hugo/navigation"
	"github.com/gohugoio/hugo/output"
	"github.com/gohugoio/hugo/related"
	"github.com/gohugoio/hugo/resources/images"
	"github.com/gohugoio/hugo/resources/kinds"
	"github.com/gohugoio/hugo/resources/page"
	"github.com/gohugoio/hugo/resources/page/pagemeta"
	"github.com/spf13/afero"

	xmaps "golang.org/x/exp/maps"
)

// InternalConfig is the internal configuration for Hugo, not read from any user provided config file.
type InternalConfig struct {
	// Server mode?
	Running bool

	Quiet          bool
	Verbose        bool
	Clock          string
	Watch          bool
	FastRenderMode bool
	LiveReloadPort int
}

// All non-params config keys for language.
var configLanguageKeys map[string]bool

func init() {
	skip := map[string]bool{
		"internal":   true,
		"c":          true,
		"rootconfig": true,
	}
	configLanguageKeys = make(map[string]bool)
	addKeys := func(v reflect.Value) {
		for i := 0; i < v.NumField(); i++ {
			name := strings.ToLower(v.Type().Field(i).Name)
			if skip[name] {
				continue
			}
			configLanguageKeys[name] = true
		}
	}
	addKeys(reflect.ValueOf(Config{}))
	addKeys(reflect.ValueOf(RootConfig{}))
	addKeys(reflect.ValueOf(config.CommonDirs{}))
	addKeys(reflect.ValueOf(langs.LanguageConfig{}))
}

type Config struct {
	// For internal use only.
	Internal InternalConfig `mapstructure:"-" json:"-"`
	// For internal use only.
	C *ConfigCompiled `mapstructure:"-" json:"-"`

	RootConfig

	// Author information.
	// Deprecated: Use taxonomies instead.
	Author map[string]any

	// Social links.
	// Deprecated: Use .Site.Params instead.
	Social map[string]string

	// The build configuration section contains build-related configuration options.
	// <docsmeta>{"identifiers": ["build"] }</docsmeta>
	Build config.BuildConfig `mapstructure:"-"`

	// The caches configuration section contains cache-related configuration options.
	// <docsmeta>{"identifiers": ["caches"] }</docsmeta>
	Caches filecache.Configs `mapstructure:"-"`

	// The httpcache configuration section contains HTTP-cache-related configuration options.
	// <docsmeta>{"identifiers": ["httpcache"] }</docsmeta>
	HTTPCache httpcache.Config `mapstructure:"-"`

	// The markup configuration section contains markup-related configuration options.
	// <docsmeta>{"identifiers": ["markup"] }</docsmeta>
	Markup markup_config.Config `mapstructure:"-"`

	// The mediatypes configuration section maps the MIME type (a string) to a configuration object for that type.
	// <docsmeta>{"identifiers": ["mediatypes"], "refs": ["types:media:type"] }</docsmeta>
	MediaTypes *config.ConfigNamespace[map[string]media.MediaTypeConfig, media.Types] `mapstructure:"-"`

	Imaging *config.ConfigNamespace[images.ImagingConfig, images.ImagingConfigInternal] `mapstructure:"-"`

	// The outputformats configuration sections maps a format name (a string) to a configuration object for that format.
	OutputFormats *config.ConfigNamespace[map[string]output.OutputFormatConfig, output.Formats] `mapstructure:"-"`

	// The outputs configuration section maps a Page Kind (a string) to a slice of output formats.
	// This can be overridden in the front matter.
	Outputs map[string][]string `mapstructure:"-"`

	// The cascade configuration section contains the top level front matter cascade configuration options,
	// a slice of page matcher and params to apply to those pages.
	Cascade *config.ConfigNamespace[[]page.PageMatcherParamsConfig, map[page.PageMatcher]maps.Params] `mapstructure:"-"`

	// The segments defines segments for the site. Used for partial/segmented builds.
	Segments *config.ConfigNamespace[map[string]segments.SegmentConfig, segments.Segments] `mapstructure:"-"`

	// Menu configuration.
	// <docsmeta>{"refs": ["config:languages:menus"] }</docsmeta>
	Menus *config.ConfigNamespace[map[string]navigation.MenuConfig, navigation.Menus] `mapstructure:"-"`

	// The deployment configuration section contains for hugo deployconfig.
	Deployment deployconfig.DeployConfig `mapstructure:"-"`

	// Module configuration.
	Module modules.Config `mapstructure:"-"`

	// Front matter configuration.
	Frontmatter pagemeta.FrontmatterConfig `mapstructure:"-"`

	// Minification configuration.
	Minify minifiers.MinifyConfig `mapstructure:"-"`

	// Permalink configuration.
	Permalinks map[string]map[string]string `mapstructure:"-"`

	// Taxonomy configuration.
	Taxonomies map[string]string `mapstructure:"-"`

	// Sitemap configuration.
	Sitemap config.SitemapConfig `mapstructure:"-"`

	// Related content configuration.
	Related related.Config `mapstructure:"-"`

	// Server configuration.
	Server config.Server `mapstructure:"-"`

	// Pagination configuration.
	Pagination config.Pagination `mapstructure:"-"`

	// Page configuration.
	Page config.PageConfig `mapstructure:"-"`

	// Privacy configuration.
	Privacy privacy.Config `mapstructure:"-"`

	// Security configuration.
	Security security.Config `mapstructure:"-"`

	// Services configuration.
	Services services.Config `mapstructure:"-"`

	// User provided parameters.
	// <docsmeta>{"refs": ["config:languages:params"] }</docsmeta>
	Params maps.Params `mapstructure:"-"`

	// The languages configuration sections maps a language code (a string) to a configuration object for that language.
	Languages map[string]langs.LanguageConfig `mapstructure:"-"`

	// UglyURLs configuration. Either a boolean or a sections map.
	UglyURLs any `mapstructure:"-"`
}

type configCompiler interface {
	CompileConfig(logger loggers.Logger) error
}

func (c Config) cloneForLang() *Config {
	x := c
	x.C = nil
	copyStringSlice := func(in []string) []string {
		if in == nil {
			return nil
		}
		out := make([]string, len(in))
		copy(out, in)
		return out
	}

	// Copy all the slices to avoid sharing.
	x.DisableKinds = copyStringSlice(x.DisableKinds)
	x.DisableLanguages = copyStringSlice(x.DisableLanguages)
	x.MainSections = copyStringSlice(x.MainSections)
	x.IgnoreLogs = copyStringSlice(x.IgnoreLogs)
	x.IgnoreFiles = copyStringSlice(x.IgnoreFiles)
	x.Theme = copyStringSlice(x.Theme)

	// Collapse all static dirs to one.
	x.StaticDir = x.staticDirs()
	// These will go away soon ...
	x.StaticDir0 = nil
	x.StaticDir1 = nil
	x.StaticDir2 = nil
	x.StaticDir3 = nil
	x.StaticDir4 = nil
	x.StaticDir5 = nil
	x.StaticDir6 = nil
	x.StaticDir7 = nil
	x.StaticDir8 = nil
	x.StaticDir9 = nil
	x.StaticDir10 = nil

	return &x
}

func (c *Config) CompileConfig(logger loggers.Logger) error {
	var transientErr error
	s := c.Timeout
	if _, err := strconv.Atoi(s); err == nil {
		// A number, assume seconds.
		s = s + "s"
	}
	timeout, err := time.ParseDuration(s)
	if err != nil {
		return fmt.Errorf("failed to parse timeout: %s", err)
	}
	disabledKinds := make(map[string]bool)
	for _, kind := range c.DisableKinds {
		kind = strings.ToLower(kind)
		if newKind := kinds.IsDeprecatedAndReplacedWith(kind); newKind != "" {
			logger.Deprecatef(false, "Kind %q used in disableKinds is deprecated, use %q instead.", kind, newKind)
			// Legacy config.
			kind = newKind
		}
		if kinds.GetKindAny(kind) == "" {
			logger.Warnf("Unknown kind %q in disableKinds configuration.", kind)
			continue
		}
		disabledKinds[kind] = true
	}
	kindOutputFormats := make(map[string]output.Formats)
	isRssDisabled := disabledKinds["rss"]
	outputFormats := c.OutputFormats.Config
	for kind, formats := range c.Outputs {
		if newKind := kinds.IsDeprecatedAndReplacedWith(kind); newKind != "" {
			logger.Deprecatef(false, "Kind %q used in outputs configuration is deprecated, use %q instead.", kind, newKind)
			kind = newKind
		}
		if disabledKinds[kind] {
			continue
		}
		if kinds.GetKindAny(kind) == "" {
			logger.Warnf("Unknown kind %q in outputs configuration.", kind)
			continue
		}
		for _, format := range formats {
			if isRssDisabled && format == "rss" {
				// Legacy config.
				continue
			}
			f, found := outputFormats.GetByName(format)
			if !found {
				transientErr = fmt.Errorf("unknown output format %q for kind %q", format, kind)
				continue
			}
			kindOutputFormats[kind] = append(kindOutputFormats[kind], f)
		}
	}

	disabledLangs := make(map[string]bool)
	for _, lang := range c.DisableLanguages {
		disabledLangs[lang] = true
	}
	for lang, language := range c.Languages {
		if !language.Disabled && disabledLangs[lang] {
			language.Disabled = true
			c.Languages[lang] = language
		}
		if language.Disabled {
			disabledLangs[lang] = true
			if lang == c.DefaultContentLanguage {
				return fmt.Errorf("cannot disable default content language %q", lang)
			}
		}
	}

	for i, s := range c.IgnoreLogs {
		c.IgnoreLogs[i] = strings.ToLower(s)
	}

	ignoredLogIDs := make(map[string]bool)
	for _, err := range c.IgnoreLogs {
		ignoredLogIDs[err] = true
	}

	baseURL, err := urls.NewBaseURLFromString(c.BaseURL)
	if err != nil {
		return err
	}

	isUglyURL := func(section string) bool {
		switch v := c.UglyURLs.(type) {
		case bool:
			return v
		case map[string]bool:
			return v[section]
		default:
			return false
		}
	}

	ignoreFile := func(s string) bool {
		return false
	}
	if len(c.IgnoreFiles) > 0 {
		regexps := make([]*regexp.Regexp, len(c.IgnoreFiles))
		for i, pattern := range c.IgnoreFiles {
			var err error
			regexps[i], err = regexp.Compile(pattern)
			if err != nil {
				return fmt.Errorf("failed to compile ignoreFiles pattern %q: %s", pattern, err)
			}
		}
		ignoreFile = func(s string) bool {
			for _, r := range regexps {
				if r.MatchString(s) {
					return true
				}
			}
			return false
		}
	}

	var clock time.Time
	if c.Internal.Clock != "" {
		var err error
		clock, err = time.Parse(time.RFC3339, c.Internal.Clock)
		if err != nil {
			return fmt.Errorf("failed to parse clock: %s", err)
		}
	}

	httpCache, err := c.HTTPCache.Compile()
	if err != nil {
		return err
	}

	// Legacy paginate values.
	if c.Paginate != 0 {
		hugo.Deprecate("site config key paginate", "Use pagination.pagerSize instead.", "v0.128.0")
		c.Pagination.PagerSize = c.Paginate
	}

	if c.PaginatePath != "" {
		hugo.Deprecate("site config key paginatePath", "Use pagination.path instead.", "v0.128.0")
		c.Pagination.Path = c.PaginatePath
	}

	c.C = &ConfigCompiled{
		Timeout:           timeout,
		BaseURL:           baseURL,
		BaseURLLiveReload: baseURL,
		DisabledKinds:     disabledKinds,
		DisabledLanguages: disabledLangs,
		IgnoredLogs:       ignoredLogIDs,
		KindOutputFormats: kindOutputFormats,
		ContentTypes:      media.DefaultContentTypes.FromTypes(c.MediaTypes.Config),
		CreateTitle:       helpers.GetTitleFunc(c.TitleCaseStyle),
		IsUglyURLSection:  isUglyURL,
		IgnoreFile:        ignoreFile,
		SegmentFilter:     c.Segments.Config.Get(func(s string) { logger.Warnf("Render segment %q not found in configuration", s) }, c.RootConfig.RenderSegments...),
		MainSections:      c.MainSections,
		Clock:             clock,
		HTTPCache:         httpCache,
		transientErr:      transientErr,
	}

	for _, s := range allDecoderSetups {
		if getCompiler := s.getCompiler; getCompiler != nil {
			if err := getCompiler(c).CompileConfig(logger); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Config) IsKindEnabled(kind string) bool {
	return !c.C.DisabledKinds[kind]
}

func (c *Config) IsLangDisabled(lang string) bool {
	return c.C.DisabledLanguages[lang]
}

// ConfigCompiled holds values and functions that are derived from the config.
type ConfigCompiled struct {
	Timeout           time.Duration
	BaseURL           urls.BaseURL
	BaseURLLiveReload urls.BaseURL
	ServerInterface   string
	KindOutputFormats map[string]output.Formats
	ContentTypes      media.ContentTypes
	DisabledKinds     map[string]bool
	DisabledLanguages map[string]bool
	IgnoredLogs       map[string]bool
	CreateTitle       func(s string) string
	IsUglyURLSection  func(section string) bool
	IgnoreFile        func(filename string) bool
	SegmentFilter     segments.SegmentFilter
	MainSections      []string
	Clock             time.Time
	HTTPCache         httpcache.ConfigCompiled

	// This is set to the last transient error found during config compilation.
	// With themes/modules we compute the configuration in multiple passes, and
	// errors with missing output format definitions may resolve itself.
	transientErr error

	mu sync.Mutex
}

// This may be set after the config is compiled.
func (c *ConfigCompiled) SetMainSections(sections []string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.MainSections = sections
}

// IsMainSectionsSet returns whether the main sections have been set.
func (c *ConfigCompiled) IsMainSectionsSet() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.MainSections != nil
}

// This is set after the config is compiled by the server command.
func (c *ConfigCompiled) SetServerInfo(baseURL, baseURLLiveReload urls.BaseURL, serverInterface string) {
	c.BaseURL = baseURL
	c.BaseURLLiveReload = baseURLLiveReload
	c.ServerInterface = serverInterface
}

// RootConfig holds all the top-level configuration options in Hugo
type RootConfig struct {
	// The base URL of the site.
	// Note that the default value is empty, but Hugo requires a valid URL (e.g. "https://example.com/") to work properly.
	// <docsmeta>{"identifiers": ["URL"] }</docsmeta>
	BaseURL string

	// Whether to build content marked as draft.X
	// <docsmeta>{"identifiers": ["draft"] }</docsmeta>
	BuildDrafts bool

	// Whether to build content with expiryDate in the past.
	// <docsmeta>{"identifiers": ["expiryDate"] }</docsmeta>
	BuildExpired bool

	// Whether to build content with publishDate in the future.
	// <docsmeta>{"identifiers": ["publishDate"] }</docsmeta>
	BuildFuture bool

	// Copyright information.
	Copyright string

	// The language to apply to content without any language indicator.
	DefaultContentLanguage string

	// By default, we put the default content language in the root and the others below their language ID, e.g. /no/.
	// Set this to true to put all languages below their language ID.
	DefaultContentLanguageInSubdir bool

	// Disable creation of alias redirect pages.
	DisableAliases bool

	// Disable lower casing of path segments.
	DisablePathToLower bool

	// Disable page kinds from build.
	DisableKinds []string

	// A list of languages to disable.
	DisableLanguages []string

	// The named segments to render.
	// This needs to match the name of the segment in the segments configuration.
	RenderSegments []string

	// Disable the injection of the Hugo generator tag on the home page.
	DisableHugoGeneratorInject bool

	// Disable live reloading in server mode.
	DisableLiveReload bool

	// Enable replacement in Pages' Content of Emoji shortcodes with their equivalent Unicode characters.
	// <docsmeta>{"identifiers": ["Content", "Unicode"] }</docsmeta>
	EnableEmoji bool

	// THe main section(s) of the site.
	// If not set, Hugo will try to guess this from the content.
	MainSections []string

	// Enable robots.txt generation.
	EnableRobotsTXT bool

	// When enabled, Hugo will apply Git version information to each Page if possible, which
	// can be used to keep lastUpdated in synch and to print version information.
	// <docsmeta>{"identifiers": ["Page"] }</docsmeta>
	EnableGitInfo bool

	// Enable to track, calculate and print metrics.
	TemplateMetrics bool

	// Enable to track, print and calculate metric hints.
	TemplateMetricsHints bool

	// Enable to disable the build lock file.
	NoBuildLock bool

	// A list of log IDs to ignore.
	IgnoreLogs []string

	// A list of regexps that match paths to ignore.
	// Deprecated: Use the settings on module imports.
	IgnoreFiles []string

	// Ignore cache.
	IgnoreCache bool

	// Enable to print greppable placeholders (on the form "[i18n] TRANSLATIONID") for missing translation strings.
	EnableMissingTranslationPlaceholders bool

	// Enable to panic on warning log entries. This may make it easier to detect the source.
	PanicOnWarning bool

	// The configured environment. Default is "development" for server and "production" for build.
	Environment string

	// The default language code.
	LanguageCode string

	// Enable if the site content has CJK language (Chinese, Japanese, or Korean). This affects how Hugo counts words.
	HasCJKLanguage bool

	// The default number of pages per page when paginating.
	// Deprecated: Use the Pagination struct.
	Paginate int

	// The path to use when creating pagination URLs, e.g. "page" in /page/2/.
	// Deprecated: Use the Pagination struct.
	PaginatePath string

	// Whether to pluralize default list titles.
	// Note that this currently only works for English, but you can provide your own title in the content file's front matter.
	PluralizeListTitles bool

	// Whether to capitalize automatic page titles, applicable to section, taxonomy, and term pages.
	CapitalizeListTitles bool

	// Make all relative URLs absolute using the baseURL.
	// <docsmeta>{"identifiers": ["baseURL"] }</docsmeta>
	CanonifyURLs bool

	// Enable this to make all relative URLs relative to content root. Note that this does not affect absolute URLs.
	RelativeURLs bool

	// Removes non-spacing marks from composite characters in content paths.
	RemovePathAccents bool

	// Whether to track and print unused templates during the build.
	PrintUnusedTemplates bool

	// Enable to print warnings for missing translation strings.
	PrintI18nWarnings bool

	// ENable to print warnings for multiple files published to the same destination.
	PrintPathWarnings bool

	// URL to be used as a placeholder when a page reference cannot be found in ref or relref. Is used as-is.
	RefLinksNotFoundURL string

	// When using ref or relref to resolve page links and a link cannot be resolved, it will be logged with this log level.
	// Valid values are ERROR (default) or WARNING. Any ERROR will fail the build (exit -1).
	RefLinksErrorLevel string

	// This will create a menu with all the sections as menu items and all the sections’ pages as “shadow-members”.
	SectionPagesMenu string

	// The length of text in words to show in a .Summary.
	SummaryLength int

	// The site title.
	Title string

	// The theme(s) to use.
	// See Modules for more a more flexible way to load themes.
	Theme []string

	// Timeout for generating page contents, specified as a duration or in seconds.
	Timeout string

	// The time zone (or location), e.g. Europe/Oslo, used to parse front matter dates without such information and in the time function.
	TimeZone string

	// Set titleCaseStyle to specify the title style used by the title template function and the automatic section titles in Hugo.
	// It defaults to AP Stylebook for title casing, but you can also set it to Chicago or Go (every word starts with a capital letter).
	TitleCaseStyle string

	// The editor used for opening up new content.
	NewContentEditor string

	// Don't sync modification time of files for the static mounts.
	NoTimes bool

	// Don't sync modification time of files for the static mounts.
	NoChmod bool

	// Clean the destination folder before a new build.
	// This currently only handles static files.
	CleanDestinationDir bool

	// A Glob pattern of module paths to ignore in the _vendor folder.
	IgnoreVendorPaths string

	config.CommonDirs `mapstructure:",squash"`

	// The odd constructs below are kept for backwards compatibility.
	// Deprecated: Use module mount config instead.
	StaticDir []string
	// Deprecated: Use module mount config instead.
	StaticDir0 []string
	// Deprecated: Use module mount config instead.
	StaticDir1 []string
	// Deprecated: Use module mount config instead.
	StaticDir2 []string
	// Deprecated: Use module mount config instead.
	StaticDir3 []string
	// Deprecated: Use module mount config instead.
	StaticDir4 []string
	// Deprecated: Use module mount config instead.
	StaticDir5 []string
	// Deprecated: Use module mount config instead.
	StaticDir6 []string
	// Deprecated: Use module mount config instead.
	StaticDir7 []string
	// Deprecated: Use module mount config instead.
	StaticDir8 []string
	// Deprecated: Use module mount config instead.
	StaticDir9 []string
	// Deprecated: Use module mount config instead.
	StaticDir10 []string
}

func (c RootConfig) staticDirs() []string {
	var dirs []string
	dirs = append(dirs, c.StaticDir...)
	dirs = append(dirs, c.StaticDir0...)
	dirs = append(dirs, c.StaticDir1...)
	dirs = append(dirs, c.StaticDir2...)
	dirs = append(dirs, c.StaticDir3...)
	dirs = append(dirs, c.StaticDir4...)
	dirs = append(dirs, c.StaticDir5...)
	dirs = append(dirs, c.StaticDir6...)
	dirs = append(dirs, c.StaticDir7...)
	dirs = append(dirs, c.StaticDir8...)
	dirs = append(dirs, c.StaticDir9...)
	dirs = append(dirs, c.StaticDir10...)
	return helpers.UniqueStringsReuse(dirs)
}

type Configs struct {
	Base                *Config
	LoadingInfo         config.LoadConfigResult
	LanguageConfigMap   map[string]*Config
	LanguageConfigSlice []*Config

	IsMultihost bool

	Modules       modules.Modules
	ModulesClient *modules.Client

	// All below is set in Init.
	Languages             langs.Languages
	LanguagesDefaultFirst langs.Languages
	ContentPathParser     *paths.PathParser

	configLangs []config.AllProvider
}

func (c *Configs) Validate(logger loggers.Logger) error {
	for p := range c.Base.Cascade.Config {
		page.CheckCascadePattern(logger, p)
	}
	return nil
}

// transientErr returns the last transient error found during config compilation.
func (c *Configs) transientErr() error {
	for _, l := range c.LanguageConfigSlice {
		if l.C.transientErr != nil {
			return l.C.transientErr
		}
	}
	return nil
}

func (c *Configs) IsZero() bool {
	// A config always has at least one language.
	return c == nil || len(c.Languages) == 0
}

func (c *Configs) Init() error {
	var languages langs.Languages
	defaultContentLanguage := c.Base.DefaultContentLanguage
	for k, v := range c.LanguageConfigMap {
		c.LanguageConfigSlice = append(c.LanguageConfigSlice, v)
		languageConf := v.Languages[k]
		language, err := langs.NewLanguage(k, defaultContentLanguage, v.TimeZone, languageConf)
		if err != nil {
			return err
		}
		languages = append(languages, language)
	}

	// Sort the sites by language weight (if set) or lang.
	sort.Slice(languages, func(i, j int) bool {
		li := languages[i]
		lj := languages[j]
		if li.Weight != lj.Weight {
			return li.Weight < lj.Weight
		}
		return li.Lang < lj.Lang
	})

	for _, l := range languages {
		c.LanguageConfigSlice = append(c.LanguageConfigSlice, c.LanguageConfigMap[l.Lang])
	}

	// Filter out disabled languages.
	var n int
	for _, l := range languages {
		if !l.Disabled {
			languages[n] = l
			n++
		}
	}
	languages = languages[:n]

	var languagesDefaultFirst langs.Languages
	for _, l := range languages {
		if l.Lang == defaultContentLanguage {
			languagesDefaultFirst = append(languagesDefaultFirst, l)
		}
	}
	for _, l := range languages {
		if l.Lang != defaultContentLanguage {
			languagesDefaultFirst = append(languagesDefaultFirst, l)
		}
	}

	c.Languages = languages
	c.LanguagesDefaultFirst = languagesDefaultFirst

	c.ContentPathParser = &paths.PathParser{LanguageIndex: languagesDefaultFirst.AsIndexSet(), IsLangDisabled: c.Base.IsLangDisabled, IsContentExt: c.Base.C.ContentTypes.IsContentSuffix}

	c.configLangs = make([]config.AllProvider, len(c.Languages))
	for i, l := range c.LanguagesDefaultFirst {
		c.configLangs[i] = ConfigLanguage{
			m:          c,
			config:     c.LanguageConfigMap[l.Lang],
			baseConfig: c.LoadingInfo.BaseConfig,
			language:   l,
		}
	}

	if len(c.Modules) == 0 {
		return errors.New("no modules loaded (need at least the main module)")
	}

	// Apply default project mounts.
	if err := modules.ApplyProjectConfigDefaults(c.Modules[0], c.configLangs...); err != nil {
		return err
	}

	// We should consolidate this, but to get a full view of the mounts in e.g. "hugo config" we need to
	// transfer any default mounts added above to the config used to print the config.
	for _, m := range c.Modules[0].Mounts() {
		var found bool
		for _, cm := range c.Base.Module.Mounts {
			if cm.Source == m.Source && cm.Target == m.Target && cm.Lang == m.Lang {
				found = true
				break
			}
		}
		if !found {
			c.Base.Module.Mounts = append(c.Base.Module.Mounts, m)
		}
	}

	// Transfer the changed mounts to the language versions (all share the same mount set, but can be displayed in different languages).
	for _, l := range c.LanguageConfigSlice {
		l.Module.Mounts = c.Base.Module.Mounts
	}

	return nil
}

func (c Configs) ConfigLangs() []config.AllProvider {
	return c.configLangs
}

func (c Configs) GetFirstLanguageConfig() config.AllProvider {
	return c.configLangs[0]
}

func (c Configs) GetByLang(lang string) config.AllProvider {
	for _, l := range c.configLangs {
		if l.Language().Lang == lang {
			return l
		}
	}
	return nil
}

// fromLoadConfigResult creates a new Config from res.
func fromLoadConfigResult(fs afero.Fs, logger loggers.Logger, res config.LoadConfigResult) (*Configs, error) {
	if !res.Cfg.IsSet("languages") {
		// We need at least one
		lang := res.Cfg.GetString("defaultContentLanguage")
		res.Cfg.Set("languages", maps.Params{lang: maps.Params{}})
	}
	bcfg := res.BaseConfig
	cfg := res.Cfg

	all := &Config{}

	err := decodeConfigFromParams(fs, logger, bcfg, cfg, all, nil)
	if err != nil {
		return nil, err
	}

	langConfigMap := make(map[string]*Config)

	languagesConfig := cfg.GetStringMap("languages")
	var isMultihost bool

	if err := all.CompileConfig(logger); err != nil {
		return nil, err
	}

	for k, v := range languagesConfig {
		mergedConfig := config.New()
		var differentRootKeys []string
		switch x := v.(type) {
		case maps.Params:
			var params maps.Params
			pv, found := x["params"]
			if found {
				params = pv.(maps.Params)
			} else {
				params = maps.Params{
					maps.MergeStrategyKey: maps.ParamsMergeStrategyDeep,
				}
				x["params"] = params
			}

			for kk, vv := range x {
				if kk == "_merge" {
					continue
				}
				if kk != maps.MergeStrategyKey && !configLanguageKeys[kk] {
					// This should have been placed below params.
					// We accidentally allowed it in the past, so we need to support it a little longer,
					// But log a warning.
					if _, found := params[kk]; !found {
						hugo.Deprecate(fmt.Sprintf("config: languages.%s.%s: custom params on the language top level", k, kk), fmt.Sprintf("Put the value below [languages.%s.params]. See https://gohugo.io/content-management/multilingual/#changes-in-hugo-01120", k), "v0.112.0")
						params[kk] = vv
					}
				}
				if kk == "baseurl" {
					// baseURL configure don the language level is a multihost setup.
					isMultihost = true
				}
				mergedConfig.Set(kk, vv)
				rootv := cfg.Get(kk)
				if rootv != nil && cfg.IsSet(kk) {
					// This overrides a root key and potentially needs a merge.
					if !reflect.DeepEqual(rootv, vv) {
						switch vvv := vv.(type) {
						case maps.Params:
							differentRootKeys = append(differentRootKeys, kk)

							// Use the language value as base.
							mergedConfigEntry := xmaps.Clone(vvv)
							// Merge in the root value.
							maps.MergeParams(mergedConfigEntry, rootv.(maps.Params))

							mergedConfig.Set(kk, mergedConfigEntry)
						default:
							// Apply new values to the root.
							differentRootKeys = append(differentRootKeys, "")
						}
					}
				} else {
					switch vv.(type) {
					case maps.Params:
						differentRootKeys = append(differentRootKeys, kk)
					default:
						// Apply new values to the root.
						differentRootKeys = append(differentRootKeys, "")
					}
				}
			}
			differentRootKeys = helpers.UniqueStringsSorted(differentRootKeys)

			if len(differentRootKeys) == 0 {
				langConfigMap[k] = all
				continue
			}

			// Create a copy of the complete config and replace the root keys with the language specific ones.
			clone := all.cloneForLang()

			if err := decodeConfigFromParams(fs, logger, bcfg, mergedConfig, clone, differentRootKeys); err != nil {
				return nil, fmt.Errorf("failed to decode config for language %q: %w", k, err)
			}
			if err := clone.CompileConfig(logger); err != nil {
				return nil, err
			}

			// Adjust Goldmark config defaults for multilingual, single-host sites.
			if len(languagesConfig) > 1 && !isMultihost && !clone.Markup.Goldmark.DuplicateResourceFiles {
				if !clone.Markup.Goldmark.DuplicateResourceFiles {
					if clone.Markup.Goldmark.RenderHooks.Link.EnableDefault == nil {
						clone.Markup.Goldmark.RenderHooks.Link.EnableDefault = types.NewBool(true)
					}
					if clone.Markup.Goldmark.RenderHooks.Image.EnableDefault == nil {
						clone.Markup.Goldmark.RenderHooks.Image.EnableDefault = types.NewBool(true)
					}
				}
			}

			langConfigMap[k] = clone
		case maps.ParamsMergeStrategy:
		default:
			panic(fmt.Sprintf("unknown type in languages config: %T", v))

		}
	}

	bcfg.PublishDir = all.PublishDir
	res.BaseConfig = bcfg
	all.CommonDirs.CacheDir = bcfg.CacheDir
	for _, l := range langConfigMap {
		l.CommonDirs.CacheDir = bcfg.CacheDir
	}

	cm := &Configs{
		Base:              all,
		LanguageConfigMap: langConfigMap,
		LoadingInfo:       res,
		IsMultihost:       isMultihost,
	}

	return cm, nil
}

func decodeConfigFromParams(fs afero.Fs, logger loggers.Logger, bcfg config.BaseConfig, p config.Provider, target *Config, keys []string) error {
	var decoderSetups []decodeWeight

	if len(keys) == 0 {
		for _, v := range allDecoderSetups {
			decoderSetups = append(decoderSetups, v)
		}
	} else {
		for _, key := range keys {
			if v, found := allDecoderSetups[key]; found {
				decoderSetups = append(decoderSetups, v)
			} else {
				logger.Warnf("Skip unknown config key %q", key)
			}
		}
	}

	// Sort them to get the dependency order right.
	sort.Slice(decoderSetups, func(i, j int) bool {
		ki, kj := decoderSetups[i], decoderSetups[j]
		if ki.weight == kj.weight {
			return ki.key < kj.key
		}
		return ki.weight < kj.weight
	})

	for _, v := range decoderSetups {
		p := decodeConfig{p: p, c: target, fs: fs, bcfg: bcfg}
		if err := v.decode(v, p); err != nil {
			return fmt.Errorf("failed to decode %q: %w", v.key, err)
		}
	}

	return nil
}

func createDefaultOutputFormats(allFormats output.Formats) map[string][]string {
	if len(allFormats) == 0 {
		panic("no output formats")
	}
	rssOut, rssFound := allFormats.GetByName(output.RSSFormat.Name)
	htmlOut, _ := allFormats.GetByName(output.HTMLFormat.Name)

	defaultListTypes := []string{htmlOut.Name}
	if rssFound {
		defaultListTypes = append(defaultListTypes, rssOut.Name)
	}

	m := map[string][]string{
		kinds.KindPage:     {htmlOut.Name},
		kinds.KindHome:     defaultListTypes,
		kinds.KindSection:  defaultListTypes,
		kinds.KindTerm:     defaultListTypes,
		kinds.KindTaxonomy: defaultListTypes,
	}

	// May be disabled
	if rssFound {
		m["rss"] = []string{rssOut.Name}
	}

	return m
}
