package helpers_test

import (
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/config/testconfig"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/spf13/afero"
)

func newTestPathSpecFromCfgAndLang(cfg config.Provider, lang string) *helpers.PathSpec {
	mfs := afero.NewMemMapFs()

	configs := testconfig.GetTestConfigs(mfs, cfg)
	var conf config.AllProvider
	if lang == "" {
		conf = configs.GetFirstLanguageConfig()
	} else {
		conf = configs.GetByLang(lang)
		if conf == nil {
			panic("no config for lang " + lang)
		}
	}
	fs := hugofs.NewFrom(mfs, conf.BaseConfig())
	ps, err := helpers.NewPathSpec(fs, conf, loggers.NewDefault())
	if err != nil {
		panic(err)
	}
	return ps
}

func newTestPathSpec(configKeyValues ...any) *helpers.PathSpec {
	cfg := config.New()
	for i := 0; i < len(configKeyValues); i += 2 {
		cfg.Set(configKeyValues[i].(string), configKeyValues[i+1])
	}
	return newTestPathSpecFromCfgAndLang(cfg, "")
}

func newTestContentSpec(cfg config.Provider) *helpers.ContentSpec {
	fs := afero.NewMemMapFs()
	conf := testconfig.GetTestConfig(fs, cfg)
	spec, err := helpers.NewContentSpec(conf, loggers.NewDefault(), fs, nil)
	if err != nil {
		panic(err)
	}
	return spec
}
