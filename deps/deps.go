package deps

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/hugo/config"
	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/hugofs"
	"github.com/spf13/hugo/output"
	"github.com/spf13/hugo/tpl"
	jww "github.com/spf13/jwalterweatherman"
)

// Deps holds dependencies used by many.
// There will be normally be only one instance of deps in play
// at a given time, i.e. one per Site built.
type Deps struct {
	// The logger to use.
	Log *jww.Notepad `json:"-"`

	// The templates to use. This will usually implement the full tpl.TemplateHandler.
	Tmpl tpl.TemplateFinder `json:"-"`

	// The file systems to use.
	Fs *hugofs.Fs `json:"-"`

	// The PathSpec to use
	*helpers.PathSpec `json:"-"`

	// The ContentSpec to use
	*helpers.ContentSpec `json:"-"`

	// The configuration to use
	Cfg config.Provider `json:"-"`

	// The translation func to use
	Translate func(translationID string, args ...interface{}) string `json:"-"`

	Language *helpers.Language

	// All the output formats available for the current site.
	OutputFormatsConfig output.Formats

	templateProvider ResourceProvider
	WithTemplate     func(templ tpl.TemplateHandler) error `json:"-"`

	translationProvider ResourceProvider
}

// ResourceProvider is used to create and refresh, and clone resources needed.
type ResourceProvider interface {
	Update(deps *Deps) error
	Clone(deps *Deps) error
}

func (d *Deps) TemplateHandler() tpl.TemplateHandler {
	return d.Tmpl.(tpl.TemplateHandler)
}

func (d *Deps) LoadResources() error {
	// Note that the translations need to be loaded before the templates.
	if err := d.translationProvider.Update(d); err != nil {
		return err
	}

	if err := d.templateProvider.Update(d); err != nil {
		return err
	}

	if th, ok := d.Tmpl.(tpl.TemplateHandler); ok {
		th.PrintErrors()
	}

	return nil
}

func New(cfg DepsCfg) (*Deps, error) {
	var (
		logger = cfg.Logger
		fs     = cfg.Fs
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
		logger = jww.NewNotepad(jww.LevelError, jww.LevelError, os.Stdout, ioutil.Discard, "", log.Ldate|log.Ltime)
	}

	if fs == nil {
		// Default to the production file system.
		fs = hugofs.NewDefault(cfg.Language)
	}

	ps, err := helpers.NewPathSpec(fs, cfg.Language)

	if err != nil {
		return nil, err
	}

	d := &Deps{
		Fs:                  fs,
		Log:                 logger,
		templateProvider:    cfg.TemplateProvider,
		translationProvider: cfg.TranslationProvider,
		WithTemplate:        cfg.WithTemplate,
		PathSpec:            ps,
		ContentSpec:         helpers.NewContentSpec(cfg.Language),
		Cfg:                 cfg.Language,
		Language:            cfg.Language,
	}

	return d, nil
}

// ForLanguage creates a copy of the Deps with the language dependent
// parts switched out.
func (d Deps) ForLanguage(l *helpers.Language) (*Deps, error) {
	var err error

	d.PathSpec, err = helpers.NewPathSpec(d.Fs, l)
	if err != nil {
		return nil, err
	}

	d.ContentSpec = helpers.NewContentSpec(l)
	d.Cfg = l
	d.Language = l

	if err := d.translationProvider.Clone(&d); err != nil {
		return nil, err
	}

	if err := d.templateProvider.Clone(&d); err != nil {
		return nil, err
	}

	return &d, nil

}

// DepsCfg contains configuration options that can be used to configure Hugo
// on a global level, i.e. logging etc.
// Nil values will be given default values.
type DepsCfg struct {

	// The Logger to use.
	Logger *jww.Notepad

	// The file systems to use
	Fs *hugofs.Fs

	// The language to use.
	Language *helpers.Language

	// The configuration to use.
	Cfg config.Provider

	// Template handling.
	TemplateProvider ResourceProvider
	WithTemplate     func(templ tpl.TemplateHandler) error

	// i18n handling.
	TranslationProvider ResourceProvider
}
