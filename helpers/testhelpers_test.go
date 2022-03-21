package helpers

import (
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/config"
	"github.com/spf13/afero"

	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/langs"
	"github.com/gohugoio/hugo/modules"
)

func newTestPathSpec(fs *hugofs.Fs, v config.Provider) *PathSpec {
	l := langs.NewDefaultLanguage(v)
	ps, _ := NewPathSpec(fs, l, nil)
	return ps
}

func newTestDefaultPathSpec(configKeyValues ...any) *PathSpec {
	cfg := newTestCfg()
	fs := hugofs.NewMem(cfg)

	for i := 0; i < len(configKeyValues); i += 2 {
		cfg.Set(configKeyValues[i].(string), configKeyValues[i+1])
	}
	return newTestPathSpec(fs, cfg)
}

func newTestCfg() config.Provider {
	v := config.NewWithTestDefaults()
	langs.LoadLanguageSettings(v, nil)
	langs.LoadLanguageSettings(v, nil)
	mod, err := modules.CreateProjectModule(v)
	if err != nil {
		panic(err)
	}
	v.Set("allModules", modules.Modules{mod})

	return v
}

func newTestContentSpec() *ContentSpec {
	v := config.NewWithTestDefaults()
	spec, err := NewContentSpec(v, loggers.NewErrorLogger(), afero.NewMemMapFs(), nil)
	if err != nil {
		panic(err)
	}
	return spec
}
