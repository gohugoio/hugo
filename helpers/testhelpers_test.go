package helpers

import (
	"github.com/spf13/viper"

	"github.com/gohugoio/hugo/hugofs"
)

func newTestPathSpec(fs *hugofs.Fs, v *viper.Viper) *PathSpec {
	l := NewDefaultLanguage(v)
	ps, _ := NewPathSpec(fs, l)
	return ps
}

func newTestDefaultPathSpec(configKeyValues ...interface{}) *PathSpec {
	v := viper.New()
	fs := hugofs.NewMem(v)
	cfg := newTestCfg(fs)

	for i := 0; i < len(configKeyValues); i += 2 {
		cfg.Set(configKeyValues[i].(string), configKeyValues[i+1])
	}
	return newTestPathSpec(fs, cfg)
}

func newTestCfg(fs *hugofs.Fs) *viper.Viper {
	v := viper.New()

	v.SetFs(fs.Source)

	return v

}

func newTestContentSpec() *ContentSpec {
	v := viper.New()
	spec, err := NewContentSpec(v)
	if err != nil {
		panic(err)
	}
	return spec
}
