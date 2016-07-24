package hugolib

import (
	"sync"

	"strings"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

type Language struct {
	Lang       string
	Title      string
	Weight     int
	params     map[string]interface{}
	paramsInit sync.Once
}

func NewLanguage(lang string) *Language {
	return &Language{Lang: lang, params: make(map[string]interface{})}
}

type Languages []*Language

func (l Languages) Len() int           { return len(l) }
func (l Languages) Less(i, j int) bool { return l[i].Weight < l[j].Weight }
func (l Languages) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }

type Multilingual struct {
	enabled bool

	Languages Languages
}

func (l *Language) Params() map[string]interface{} {
	l.paramsInit.Do(func() {
		// Merge with global config.
		// TODO(bep) consider making this part of a constructor func.
		globalParams := viper.GetStringMap("Params")
		for k, v := range globalParams {
			if _, ok := l.params[k]; !ok {
				l.params[k] = v
			}
		}
	})
	return l.params
}

func (l *Language) SetParam(k string, v interface{}) {
	l.params[k] = v
}

func (l *Language) GetString(key string) string { return cast.ToString(l.Get(key)) }
func (ml *Language) GetStringMap(key string) map[string]interface{} {
	return cast.ToStringMap(ml.Get(key))
}

func (l *Language) GetStringMapString(key string) map[string]string {
	return cast.ToStringMapString(l.Get(key))
}

func (l *Language) Get(key string) interface{} {
	key = strings.ToLower(key)
	if v, ok := l.params[key]; ok {
		return v
	}
	return viper.Get(key)
}

func (s *Site) SetMultilingualConfig(currentLang *Language, languages Languages) {

	// TODO(bep) multilingo evaluate
	viper.Set("CurrentLanguage", currentLang)
	ml := &Multilingual{
		enabled:   len(languages) > 0,
		Languages: languages,
	}
	viper.Set("Multilingual", ml.enabled)
	s.Multilingual = ml
}

func (s *Site) multilingualEnabled() bool {
	return s.Multilingual != nil && s.Multilingual.enabled
}

func currentLanguageString() string {
	return currentLanguage().Lang
}

func currentLanguage() *Language {
	l := viper.Get("CurrentLanguage")
	if l == nil {
		panic("CurrentLanguage not set")
	}
	return l.(*Language)
}
