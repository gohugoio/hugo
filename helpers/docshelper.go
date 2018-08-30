package helpers

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/alecthomas/chroma/lexers"
	"github.com/gohugoio/hugo/docshelper"
)

// This is is just some helpers used to create some JSON used in the Hugo docs.
func init() {

	docsProvider := func() map[string]interface{} {
		docs := make(map[string]interface{})

		var chromaLexers []interface{}

		sort.Sort(lexers.Registry.Lexers)

		for _, l := range lexers.Registry.Lexers {

			config := l.Config()

			var filenames []string
			filenames = append(filenames, config.Filenames...)
			filenames = append(filenames, config.AliasFilenames...)

			aliases := config.Aliases

			for _, filename := range filenames {
				alias := strings.TrimSpace(strings.TrimPrefix(filepath.Ext(filename), "."))
				if alias != "" {
					aliases = append(aliases, alias)
				}
			}

			sort.Strings(aliases)
			aliases = UniqueStrings(aliases)

			lexerEntry := struct {
				Name    string
				Aliases []string
			}{
				config.Name,
				aliases,
			}

			chromaLexers = append(chromaLexers, lexerEntry)

			docs["lexers"] = chromaLexers
		}
		return docs

	}

	docshelper.AddDocProvider("chroma", docsProvider)
}
