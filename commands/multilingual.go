package commands

import (
	"fmt"
	"sort"

	"strings"

	"github.com/spf13/cast"
	"github.com/spf13/hugo/hugolib"
	"github.com/spf13/viper"
)

var langConfigsList hugolib.Languages

func readMultilingualConfiguration() error {
	multilingual := viper.GetStringMap("Multilingual")
	if len(multilingual) == 0 {
		// TODO(bep) multilingo langConfigsList = append(langConfigsList, hugolib.NewLanguage("en"))
		return nil
	}

	var err error

	langConfigsList, err = toSortedLanguages(multilingual)

	if err != nil {
		return fmt.Errorf("Failed to parse multilingual config: %s", err)
	}

	return nil
}

func toSortedLanguages(l map[string]interface{}) (hugolib.Languages, error) {
	langs := make(hugolib.Languages, len(l))

	for lang, langConf := range l {
		langsMap, ok := langConf.(map[string]interface{})

		if !ok {
			return nil, fmt.Errorf("Language config is not a map: %v", langsMap)
		}

		language := hugolib.NewLanguage(lang)

		for k, v := range langsMap {
			loki := strings.ToLower(k)
			switch loki {
			case "title":
				language.Title = cast.ToString(v)
			case "weight":
				language.Weight = cast.ToInt(v)
			}

			// Put all into the Params map
			// TODO(bep) reconsile with the type handling etc. from other params handlers.
			language.SetParam(loki, v)
		}

		langs = append(langs, language)
	}

	sort.Sort(langs)

	return langs, nil
}
