package output_test

import (
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

func TestCanonical(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
[outputs]
 home = ["notcanonical", "html", "rss"]

[outputFormats]
[outputFormats.notcanonical]
mediaType = 'text/html'
path = 'not'
isHTML = true
-- layouts/all.html --
All. Canonical: {{ .OutputFormats.Canonical.RelPermalink }}.
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/not/index.html", "All. Canonical: /.")
	b.AssertFileContent("public/index.html", "All. Canonical: /.")
}
