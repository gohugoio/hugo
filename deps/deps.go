package deps

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gohugoio/hugo/cache/filecache"
	"github.com/gohugoio/hugo/common/hexec"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/common/types"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/config/security"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/langs"
	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/resources/page"
	"github.com/gohugoio/hugo/resources/postpub"

	"github.com/gohugoio/hugo/metrics"
	"github.com/gohugoio/hugo/output"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/source"
	"github.com/gohugoio/hugo/tpl"
	"github.com/spf13/cast"
	jww "github.com/spf13/jwalterweatherman"
)

// Deps holds dependencies used by many.
// There will be normally only one instance of deps in play
// at a given time, i.e. one per Site built.
type Deps struct {

	// The logger to use.
	Log loggers.Logger `json:"-"`

	// Used to log errors that may repeat itself many times.
	LogDistinct loggers.Logger

	ExecHelper *hexec.Exec

	// The templates to use. This will usually implement the full tpl.TemplateManager.
	tmpl tpl.TemplateHandler

	// We use this to parse and execute ad-hoc text templates.
	textTmpl tpl.TemplateParseFinder

	// The file systems to use.
	Fs *hugofs.Fs `json:"-"`

	// The PathSpec to use
	*helpers.PathSpec `json:"-"`

	// The ContentSpec to use
	*helpers.ContentSpec `json:"-"`

	// The SourceSpec to use
	SourceSpec *source.SourceSpec `json:"-"`

	// The Resource Spec to use
	ResourceSpec *resources.Spec

	// The configuration to use
	Cfg config.Provider `json:"-"`

	// The file cache to use.
	FileCaches filecache.Caches

	// The translation func to use
	Translate func(ctx context.Context, translationID string, templateData any) string `json:"-"`

	// The language in use. TODO(bep) consolidate with site
	Language *langs.Language

	// The site building.
	Site page.Site

	// All the output formats available for the current site.
	OutputFormatsConfig output.Formats

	// FilenameHasPostProcessPrefix is a set of filenames in /public that
	// contains a post-processing prefix.
	FilenameHasPostProcessPrefix []string

	templateProvider ResourceProvider
	WithTemplate     func(templ tpl.TemplateManager) error `json:"-"`

	// Used in tests
	OverloadedTemplateFuncs map[string]any

	translationProvider ResourceProvider

	Metrics metrics.Provider

	// Timeout is configurable in site config.
	Timeout time.Duration

	// BuildStartListeners will be notified before a build starts.
	BuildStartListeners *Listeners

	// Resources that gets closed when the build is done or the server shuts down.
	BuildClosers *Closers

	// Atomic values set during a build.
	// This is common/global for all sites.
	BuildState *BuildState

	// Whether we are in running (server) mode
	Running bool

	*globalErrHandler
}

type globalErrHandler struct {
	// Channel for some "hard to get to" build errors
	buildErrors chan error
}

// SendErr sends the error on a channel to be handled later.
// This can be used in situations where returning and aborting the current
// operation isn't practical.
func (e *globalErrHandler) SendError(err error) {
	if e.buildErrors != nil {
		select {
		case e.buildErrors <- err:
		default:
		}
		return
	}

	jww.ERROR.Println(err)
}

func (e *globalErrHandler) StartErrorCollector() chan error {
	e.buildErrors = make(chan error, 10)
	return e.buildErrors
}

// Listeners represents an event listener.
type Listeners struct {
	sync.Mutex

	// A list of funcs to be notified about an event.
	listeners []func()
}

// Add adds a function to a Listeners instance.
func (b *Listeners) Add(f func()) {
	if b == nil {
		return
	}
	b.Lock()
	defer b.Unlock()
	b.listeners = append(b.listeners, f)
}

// Notify executes all listener functions.
func (b *Listeners) Notify() {
	b.Lock()
	defer b.Unlock()
	for _, notify := range b.listeners {
		notify()
	}
}

// ResourceProvider is used to create and refresh, and clone resources needed.
type ResourceProvider interface {
	Update(deps *Deps) error
	Clone(deps *Deps) error
}

func (d *Deps) Tmpl() tpl.TemplateHandler {
	return d.tmpl
}

func (d *Deps) TextTmpl() tpl.TemplateParseFinder {
	return d.textTmpl
}

func (d *Deps) SetTmpl(tmpl tpl.TemplateHandler) {
	d.tmpl = tmpl
}

func (d *Deps) SetTextTmpl(tmpl tpl.TemplateParseFinder) {
	d.textTmpl = tmpl
}

// LoadResources loads translations and templates.
func (d *Deps) LoadResources() error {
	// Note that the translations need to be loaded before the templates.
	if err := d.translationProvider.Update(d); err != nil {
		return fmt.Errorf("loading translations: %w", err)
	}

	if err := d.templateProvider.Update(d); err != nil {
		return fmt.Errorf("loading templates: %w", err)
	}

	return nil
}

// New initializes a Dep struct.
// Defaults are set for nil values,
// but TemplateProvider, TranslationProvider and Language are always required.
func New(cfg DepsCfg) (*Deps, error) {
	var (
		logger = cfg.Logger
		fs     = cfg.Fs
		d      *Deps
	)

	if cfg.TemplateProvider == nil {
		panic("Must have a TemplateProvider")
	}

	if cfg.TranslationProvider == nil {
		panic("Must have a TranslationProvider")
	}

	if cfg.Language == nil {
		panic("Must have a Language")
	}

	if logger == nil {
		logger = loggers.NewErrorLogger()
	}

	if fs == nil {
		// Default to the production file system.
		fs = hugofs.NewDefault(cfg.Language)
	}

	if cfg.MediaTypes == nil {
		cfg.MediaTypes = media.DefaultTypes
	}

	if cfg.OutputFormats == nil {
		cfg.OutputFormats = output.DefaultFormats
	}

	securityConfig, err := security.DecodeConfig(cfg.Cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create security config from configuration: %w", err)
	}
	execHelper := hexec.New(securityConfig)

	var filenameHasPostProcessPrefixMu sync.Mutex
	hashBytesReceiverFunc := func(name string, match bool) {
		if !match {
			return
		}
		filenameHasPostProcessPrefixMu.Lock()
		d.FilenameHasPostProcessPrefix = append(d.FilenameHasPostProcessPrefix, name)
		filenameHasPostProcessPrefixMu.Unlock()
	}

	// Skip binary files.
	hashBytesSHouldCheck := func(name string) bool {
		ext := strings.TrimPrefix(filepath.Ext(name), ".")
		mime, _, found := cfg.MediaTypes.GetBySuffix(ext)
		if !found {
			return false
		}
		switch mime.MainType {
		case "text", "application":
			return true
		default:
			return false
		}
	}
	fs.PublishDir = hugofs.NewHasBytesReceiver(fs.PublishDir, hashBytesSHouldCheck, hashBytesReceiverFunc, []byte(postpub.PostProcessPrefix))

	ps, err := helpers.NewPathSpec(fs, cfg.Language, logger)
	if err != nil {
		return nil, fmt.Errorf("create PathSpec: %w", err)
	}

	fileCaches, err := filecache.NewCaches(ps)
	if err != nil {
		return nil, fmt.Errorf("failed to create file caches from configuration: %w", err)
	}

	errorHandler := &globalErrHandler{}
	buildState := &BuildState{}

	resourceSpec, err := resources.NewSpec(ps, fileCaches, buildState, logger, errorHandler, execHelper, cfg.OutputFormats, cfg.MediaTypes)
	if err != nil {
		return nil, err
	}

	contentSpec, err := helpers.NewContentSpec(cfg.Language, logger, ps.BaseFs.Content.Fs, execHelper)
	if err != nil {
		return nil, err
	}

	sp := source.NewSourceSpec(ps, nil, fs.Source)

	timeout := 30 * time.Second
	if cfg.Cfg.IsSet("timeout") {
		v := cfg.Cfg.Get("timeout")
		d, err := types.ToDurationE(v)
		if err == nil {
			timeout = d
		}
	}
	ignoreErrors := cast.ToStringSlice(cfg.Cfg.Get("ignoreErrors"))
	ignorableLogger := loggers.NewIgnorableLogger(logger, ignoreErrors...)

	logDistinct := helpers.NewDistinctLogger(logger)

	d = &Deps{
		Fs:                      fs,
		Log:                     ignorableLogger,
		LogDistinct:             logDistinct,
		ExecHelper:              execHelper,
		templateProvider:        cfg.TemplateProvider,
		translationProvider:     cfg.TranslationProvider,
		WithTemplate:            cfg.WithTemplate,
		OverloadedTemplateFuncs: cfg.OverloadedTemplateFuncs,
		PathSpec:                ps,
		ContentSpec:             contentSpec,
		SourceSpec:              sp,
		ResourceSpec:            resourceSpec,
		Cfg:                     cfg.Language,
		Language:                cfg.Language,
		Site:                    cfg.Site,
		FileCaches:              fileCaches,
		BuildStartListeners:     &Listeners{},
		BuildClosers:            &Closers{},
		BuildState:              buildState,
		Running:                 cfg.Running,
		Timeout:                 timeout,
		globalErrHandler:        errorHandler,
	}

	if cfg.Cfg.GetBool("templateMetrics") {
		d.Metrics = metrics.NewProvider(cfg.Cfg.GetBool("templateMetricsHints"))
	}

	return d, nil
}

func (d *Deps) Close() error {
	return d.BuildClosers.Close()
}

// ForLanguage creates a copy of the Deps with the language dependent
// parts switched out.
func (d Deps) ForLanguage(cfg DepsCfg, onCreated func(d *Deps) error) (*Deps, error) {
	l := cfg.Language
	var err error

	d.PathSpec, err = helpers.NewPathSpecWithBaseBaseFsProvided(d.Fs, l, d.Log, d.BaseFs)
	if err != nil {
		return nil, err
	}

	d.ContentSpec, err = helpers.NewContentSpec(l, d.Log, d.BaseFs.Content.Fs, d.ExecHelper)
	if err != nil {
		return nil, err
	}

	d.Site = cfg.Site

	// These are common for all sites, so reuse.
	// TODO(bep) clean up these inits.
	resourceCache := d.ResourceSpec.ResourceCache
	postBuildAssets := d.ResourceSpec.PostBuildAssets
	d.ResourceSpec, err = resources.NewSpec(d.PathSpec, d.ResourceSpec.FileCaches, d.BuildState, d.Log, d.globalErrHandler, d.ExecHelper, cfg.OutputFormats, cfg.MediaTypes)
	if err != nil {
		return nil, err
	}
	d.ResourceSpec.ResourceCache = resourceCache
	d.ResourceSpec.PostBuildAssets = postBuildAssets

	d.Cfg = l
	d.Language = l

	if onCreated != nil {
		if err = onCreated(&d); err != nil {
			return nil, err
		}
	}

	if err := d.translationProvider.Clone(&d); err != nil {
		return nil, err
	}

	if err := d.templateProvider.Clone(&d); err != nil {
		return nil, err
	}

	d.BuildStartListeners = &Listeners{}

	return &d, nil
}

// DepsCfg contains configuration options that can be used to configure Hugo
// on a global level, i.e. logging etc.
// Nil values will be given default values.
type DepsCfg struct {

	// The Logger to use.
	Logger loggers.Logger

	// The file systems to use
	Fs *hugofs.Fs

	// The language to use.
	Language *langs.Language

	// The Site in use
	Site page.Site

	// The configuration to use.
	Cfg config.Provider

	// The media types configured.
	MediaTypes media.Types

	// The output formats configured.
	OutputFormats output.Formats

	// Template handling.
	TemplateProvider ResourceProvider
	WithTemplate     func(templ tpl.TemplateManager) error
	// Used in tests
	OverloadedTemplateFuncs map[string]any

	// i18n handling.
	TranslationProvider ResourceProvider

	// Whether we are in running (server) mode
	Running bool
}

// BuildState are flags that may be turned on during a build.
type BuildState struct {
	counter uint64
}

func (b *BuildState) Incr() int {
	return int(atomic.AddUint64(&b.counter, uint64(1)))
}

func NewBuildState() BuildState {
	return BuildState{}
}

type Closer interface {
	Close() error
}

type Closers struct {
	mu sync.Mutex
	cs []Closer
}

func (cs *Closers) Add(c Closer) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.cs = append(cs.cs, c)
}

func (cs *Closers) Close() error {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	for _, c := range cs.cs {
		c.Close()
	}

	cs.cs = cs.cs[:0]

	return nil
}
