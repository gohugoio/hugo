package commands

import (
	"sort"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

var langConfigs map[string]interface{}
var langConfigsList langConfigsSortable

func readMultilingualConfiguration() {
	multilingual := viper.GetStringMap("Multilingual")
	if len(multilingual) == 0 {
		langConfigsList = append(langConfigsList, "")
		return
	}

	langConfigs = make(map[string]interface{})
	for lang, config := range multilingual {
		langConfigs[lang] = config
		langConfigsList = append(langConfigsList, lang)
	}
	sort.Sort(langConfigsList)
}

type langConfigsSortable []string

func (p langConfigsSortable) Len() int           { return len(p) }
func (p langConfigsSortable) Less(i, j int) bool { return weightForLang(p[i]) < weightForLang(p[j]) }
func (p langConfigsSortable) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func weightForLang(lang string) int {
	conf := langConfigs[lang]
	if conf == nil {
		return 0
	}
	m := cast.ToStringMap(conf)
	return cast.ToInt(m["weight"])
}
