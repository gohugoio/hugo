package deps

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/hugofs"
	"github.com/spf13/hugo/tplapi"
	jww "github.com/spf13/jwalterweatherman"
)

// Deps holds dependencies used by many.
// There will be normally be only one instance of deps in play
// at a given time, i.e. one per Site built.
type Deps struct {
	// The logger to use.
	Log *jww.Notepad `json:"-"`

	// The templates to use.
	Tmpl tplapi.Template `json:"-"`

	// The file systems to use.
	Fs *hugofs.Fs `json:"-"`

	// The PathSpec to use
	*helpers.PathSpec `json:"-"`

	templateProvider TemplateProvider
	WithTemplate     func(templ tplapi.Template) error

	// TODO(bep) globals next in line: Viper

}

// Used to create and refresh, and clone the template.
type TemplateProvider interface {
	Update(deps *Deps) error
	Clone(deps *Deps) error
}

func (d *Deps) LoadTemplates() error {
	if err := d.templateProvider.Update(d); err != nil {
		return err
	}
	d.Tmpl.PrintErrors()
	return nil
}

func New(cfg DepsCfg) *Deps {
	var (
		logger = cfg.Logger
		fs     = cfg.Fs
	)

	if cfg.TemplateProvider == nil {
		panic("Must have a TemplateProvider")
	}

	if cfg.Language == nil {
		panic("Must have a Language")
	}

	if logger == nil {
		logger = jww.NewNotepad(jww.LevelError, jww.LevelError, os.Stdout, ioutil.Discard, "", log.Ldate|log.Ltime)
	}

	if fs == nil {
		// Default to the most used file systems.
		fs = hugofs.NewMem()
	}

	d := &Deps{
		Fs:               fs,
		Log:              logger,
		templateProvider: cfg.TemplateProvider,
		WithTemplate:     cfg.WithTemplate,
		PathSpec:         helpers.NewPathSpec(fs, cfg.Language),
	}

	return d
}

// ForLanguage creates a copy of the Deps with the language dependent
// parts switched out.
func (d Deps) ForLanguage(l *helpers.Language) (*Deps, error) {

	d.PathSpec = helpers.NewPathSpec(d.Fs, l)
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

	// Template handling.
	TemplateProvider TemplateProvider
	WithTemplate     func(templ tplapi.Template) error
}
