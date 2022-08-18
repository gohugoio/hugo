package helpers

import (
	"sort"

	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/gohugoio/hugo/docshelper"
)

// This is is just some helpers used to create some JSON used in the Hugo docs.
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

		return docshelper.DocProvider{"chroma": map[string]any{"lexers": chromaLexers}}
	}

	docshelper.AddDocProviderFunc(docsProvider)
}
