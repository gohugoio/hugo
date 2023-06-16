package deps

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/bep/logg"
	"github.com/gohugoio/hugo/common/hexec"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/config/allconfig"
	"github.com/gohugoio/hugo/config/security"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/resources/page"
	"github.com/gohugoio/hugo/resources/postpub"

	"github.com/gohugoio/hugo/metrics"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/source"
	"github.com/gohugoio/hugo/tpl"
	"github.com/spf13/afero"
)

// Deps holds dependencies used by many.
// There will be normally only one instance of deps in play
// at a given time, i.e. one per Site built.
type Deps struct {

	// The logger to use.
	Log loggers.Logger `json:"-"`

	ExecHelper *hexec.Exec

	// The templates to use. This will usually implement the full tpl.TemplateManager.
	tmplHandlers *tpl.TemplateHandlers

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
	Conf config.AllProvider `json:"-"`

	// The translation func to use
	Translate func(ctx context.Context, translationID string, templateData any) string `json:"-"`

	// The site building.
	Site page.Site

	TemplateProvider ResourceProvider
	// Used in tests
	OverloadedTemplateFuncs map[string]any

	TranslationProvider ResourceProvider

	Metrics metrics.Provider

	// BuildStartListeners will be notified before a build starts.
	BuildStartListeners *Listeners

	// Resources that gets closed when the build is done or the server shuts down.
	BuildClosers *Closers

	// This is common/global for all sites.
	BuildState *BuildState

	*globalErrHandler
}

func (d Deps) Clone(s page.Site, conf config.AllProvider) (*Deps, error) {
	d.Conf = conf
	d.Site = s
	d.ExecHelper = nil
	d.ContentSpec = nil

	if err := d.Init(); err != nil {
		return nil, err
	}

	return &d, nil

}

func (d *Deps) SetTempl(t *tpl.TemplateHandlers) {
	d.tmplHandlers = t
}

func (d *Deps) Init() error {
	if d.Conf == nil {
		panic("conf is nil")
	}

	if d.Fs == nil {
		// For tests.
		d.Fs = hugofs.NewFrom(afero.NewMemMapFs(), d.Conf.BaseConfig())
	}

	if d.Log == nil {
		d.Log = loggers.NewDefault()
	}

	if d.globalErrHandler == nil {
		d.globalErrHandler = &globalErrHandler{
			logger: d.Log,
		}
	}

	if d.BuildState == nil {
		d.BuildState = &BuildState{}
	}

	if d.BuildStartListeners == nil {
		d.BuildStartListeners = &Listeners{}
	}

	if d.BuildClosers == nil {
		d.BuildClosers = &Closers{}
	}

	if d.Metrics == nil && d.Conf.TemplateMetrics() {
		d.Metrics = metrics.NewProvider(d.Conf.TemplateMetricsHints())
	}

	if d.ExecHelper == nil {
		d.ExecHelper = hexec.New(d.Conf.GetConfigSection("security").(security.Config))
	}

	if d.PathSpec == nil {
		hashBytesReceiverFunc := func(name string, match bool) {
			if !match {
				return
			}
			d.BuildState.AddFilenameWithPostPrefix(name)
		}

		// Skip binary files.
		mediaTypes := d.Conf.GetConfigSection("mediaTypes").(media.Types)
		hashBytesSHouldCheck := func(name string) bool {
			ext := strings.TrimPrefix(filepath.Ext(name), ".")
			return mediaTypes.IsTextSuffix(ext)
		}
		d.Fs.PublishDir = hugofs.NewHasBytesReceiver(d.Fs.PublishDir, hashBytesSHouldCheck, hashBytesReceiverFunc, []byte(postpub.PostProcessPrefix))
		pathSpec, err := helpers.NewPathSpec(d.Fs, d.Conf, d.Log)
		if err != nil {
			return err
		}
		d.PathSpec = pathSpec
	} else {
		var err error
		d.PathSpec, err = helpers.NewPathSpecWithBaseBaseFsProvided(d.Fs, d.Conf, d.Log, d.PathSpec.BaseFs)
		if err != nil {
			return err
		}
	}

	if d.ContentSpec == nil {
		contentSpec, err := helpers.NewContentSpec(d.Conf, d.Log, d.Content.Fs, d.ExecHelper)
		if err != nil {
			return err
		}
		d.ContentSpec = contentSpec
	}

	if d.SourceSpec == nil {
		d.SourceSpec = source.NewSourceSpec(d.PathSpec, nil, d.Fs.Source)
	}

	var common *resources.SpecCommon
	var imageCache *resources.ImageCache
	if d.ResourceSpec != nil {
		common = d.ResourceSpec.SpecCommon
		imageCache = d.ResourceSpec.ImageCache
	}

	resourceSpec, err := resources.NewSpec(d.PathSpec, common, imageCache, d.BuildState, d.Log, d, d.ExecHelper)
	if err != nil {
		return fmt.Errorf("failed to create resource spec: %w", err)
	}
	d.ResourceSpec = resourceSpec

	return nil
}

func (d *Deps) Compile(prototype *Deps) error {
	var err error
	if prototype == nil {
		if err = d.TemplateProvider.NewResource(d); err != nil {
			return err
		}
		if err = d.TranslationProvider.NewResource(d); err != nil {
			return err
		}
		return nil
	}

	if err = d.TemplateProvider.CloneResource(d, prototype); err != nil {
		return err
	}

	if err = d.TranslationProvider.CloneResource(d, prototype); err != nil {
		return err
	}

	return nil
}

type globalErrHandler struct {
	logger loggers.Logger

	// Channel for some "hard to get to" build errors
	buildErrors chan error
	// Used to signal that the build is done.
	quit chan struct{}
}

// SendError sends the error on a channel to be handled later.
// This can be used in situations where returning and aborting the current
// operation isn't practical.
func (e *globalErrHandler) SendError(err error) {
	if e.buildErrors != nil {
		select {
		case <-e.quit:
		case e.buildErrors <- err:
		default:
		}
		return
	}
	e.logger.Errorln(err)
}

func (e *globalErrHandler) StartErrorCollector() chan error {
	e.quit = make(chan struct{})
	e.buildErrors = make(chan error, 10)
	return e.buildErrors
}

func (e *globalErrHandler) StopErrorCollector() {
	if e.buildErrors != nil {
		close(e.quit)
		close(e.buildErrors)
	}
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
	NewResource(dst *Deps) error
	CloneResource(dst, src *Deps) error
}

func (d *Deps) Tmpl() tpl.TemplateHandler {
	return d.tmplHandlers.Tmpl
}

func (d *Deps) TextTmpl() tpl.TemplateParseFinder {
	return d.tmplHandlers.TxtTmpl
}

func (d *Deps) Close() error {
	return d.BuildClosers.Close()
}

// DepsCfg contains configuration options that can be used to configure Hugo
// on a global level, i.e. logging etc.
// Nil values will be given default values.
type DepsCfg struct {
	// The logger to use. Only set in some tests.
	// TODO(bep) get rid of this.
	TestLogger loggers.Logger

	// The logging level to use.
	LogLevel logg.Level

	// Where to write the logs.
	// Currently we typically write everything to stdout.
	LogOut io.Writer

	// The file systems to use
	Fs *hugofs.Fs

	// The Site in use
	Site page.Site

	Configs *allconfig.Configs

	// Template handling.
	TemplateProvider ResourceProvider

	// i18n handling.
	TranslationProvider ResourceProvider
}

// BuildState are state used during a build.
type BuildState struct {
	counter uint64

	mu sync.Mutex // protects state below.

	// A set of ilenames in /public that
	// contains a post-processing prefix.
	filenamesWithPostPrefix map[string]bool
}

func (b *BuildState) AddFilenameWithPostPrefix(filename string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.filenamesWithPostPrefix == nil {
		b.filenamesWithPostPrefix = make(map[string]bool)
	}
	b.filenamesWithPostPrefix[filename] = true
}

func (b *BuildState) GetFilenamesWithPostPrefix() []string {
	b.mu.Lock()
	defer b.mu.Unlock()
	var filenames []string
	for filename := range b.filenamesWithPostPrefix {
		filenames = append(filenames, filename)
	}
	sort.Strings(filenames)
	return filenames
}

func (b *BuildState) Incr() int {
	return int(atomic.AddUint64(&b.counter, uint64(1)))
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
