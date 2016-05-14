package hugolib

import (
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

type Multilingual struct {
	enabled bool
	config  *viper.Viper

	Languages []string
}

func (ml *Multilingual) GetString(key string) string { return cast.ToString(ml.Get(key)) }
func (ml *Multilingual) GetStringMap(key string) map[string]interface{} {
	return cast.ToStringMap(ml.Get(key))
}

func (ml *Multilingual) GetStringMapString(key string) map[string]string {
	return cast.ToStringMapString(ml.Get(key))
}

func (ml *Multilingual) Get(key string) interface{} {
	if ml != nil && ml.config != nil && ml.config.IsSet(key) {
		return ml.config.Get(key)
	}
	return viper.Get(key)
}

func (s *Site) SetMultilingualConfig(currentLang string, orderedLanguages []string, langConfigs map[string]interface{}) {
	conf := viper.New()
	for k, val := range cast.ToStringMap(langConfigs[currentLang]) {
		conf.Set(k, val)
	}
	conf.Set("CurrentLanguage", currentLang)
	ml := &Multilingual{
		enabled:   len(langConfigs) > 0,
		config:    conf,
		Languages: orderedLanguages,
	}
	viper.Set("Multilingual", ml.enabled)
	s.Multilingual = ml
}

func (s *Site) multilingualEnabled() bool {
	return s.Multilingual != nil && s.Multilingual.enabled
}
