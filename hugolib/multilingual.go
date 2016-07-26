package hugolib

import (
	"sync"

	"sort"
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

// TODO(bep) multilingo
func newDefaultLanguage() *Language {
	return NewLanguage("en")
}

type Languages []*Language

func NewLanguages(l ...*Language) Languages {
	languages := make(Languages, len(l))
	for i := 0; i < len(l); i++ {
		languages[i] = l[i]
	}
	sort.Sort(languages)
	return languages
}

func (l Languages) Len() int           { return len(l) }
func (l Languages) Less(i, j int) bool { return l[i].Weight < l[j].Weight }
func (l Languages) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }

type Multilingual struct {
	Languages Languages

	langMap     map[string]*Language
	langMapInit sync.Once
}

func (ml *Multilingual) Language(lang string) *Language {
	ml.langMapInit.Do(func() {
		ml.langMap = make(map[string]*Language)
		for _, l := range ml.Languages {
			ml.langMap[l.Lang] = l
		}
	})
	return ml.langMap[lang]
}

func (ml *Multilingual) enabled() bool {
	return len(ml.Languages) > 0
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

// TODO(bep) multilingo move this to a constructor.
func (s *Site) SetMultilingualConfig(currentLang *Language, languages Languages) {

	ml := &Multilingual{
		Languages: languages,
	}
	viper.Set("Multilingual", ml.enabled())
	s.Multilingual = ml
}

func (s *Site) multilingualEnabled() bool {
	return s.Multilingual != nil && s.Multilingual.enabled()
}

// TODO(bep) multilingo remove these
func (s *Site) currentLanguageString() string {
	return s.currentLanguage().Lang
}

func (s *Site) currentLanguage() *Language {
	return s.Lang
}
