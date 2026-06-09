package helpers

import (
	"sort"

	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/gohugoio/hugo/docshelper"
)

// This is just a helper used to create some JSON used in the Hugo docs.
func init() {
	docsProvider := func() docshelper.DocProvider {
		var chromaLexers []any

		sort.Sort(lexers.GlobalLexerRegistry.Lexers)

		for _, l := range lexers.GlobalLexerRegistry.Lexers {

			config := l.Config()

			lexerEntry := struct {
				Name    string
				Aliases []string
			}{
				config.Name,
				config.Aliases,
			}

			chromaLexers = append(chromaLexers, lexerEntry)

		}

		var styleInfos []styleInfo

		for name := range styles.Registry {
			style := styles.Registry[name]

			styleInfos = append(styleInfos, styleInfo{
				Name:        name,
				Counterpart: style.Counterpart,
				Mode:        style.Mode().String(),
			})
		}

		sort.Slice(styleInfos, func(i, j int) bool {
			return styleInfos[i].Name < styleInfos[j].Name
		})

		return docshelper.DocProvider{"chroma": map[string]any{
			"lexers": chromaLexers,
			"styles": styleInfos,
		}}
	}

	docshelper.AddDocProviderFunc(docsProvider)
}

type styleInfo struct {
	Name        string `json:"name"`
	Mode        string `json:"mode"`
	Counterpart string `json:"counterpart"`
}
