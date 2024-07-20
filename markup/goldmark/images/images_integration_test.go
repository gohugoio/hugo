package images_test

import (
	"strings"
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

func TestDisableWrapStandAloneImageWithinParagraph(t *testing.T) {
	t.Parallel()

	filesTemplate := `
-- config.toml --
[markup.goldmark.renderer]
	unsafe = false
[markup.goldmark.parser]
wrapStandAloneImageWithinParagraph = CONFIG_VALUE
[markup.goldmark.parser.attribute]
	block = true
	title = true
-- content/p1.md --
---
title: "p1"
---

This is an inline image: ![Inline Image](/inline.jpg). Some more text.

![Block Image](/block.jpg)
{.b}


-- layouts/_default/single.html --
{{ .Content }}
`

	t.Run("With Hook, no wrap", func(t *testing.T) {
		files := strings.ReplaceAll(filesTemplate, "CONFIG_VALUE", "false")
		files = files + `-- layouts/_default/_markup/render-image.html --
{{ if .IsBlock }}
<figure class="{{ .Attributes.class }}">
	<img src="{{ .Destination | safeURL }}" alt="{{ .Text }}|{{ .Ordinal }}" />
</figure>
{{ else }}
	<img src="{{ .Destination | safeURL }}" alt="{{ .Text }}|{{ .Ordinal }}" />
{{ end }}
`
		b := hugolib.Test(t, files)

		b.AssertFileContent("public/p1/index.html",
			"This is an inline image: \n\t<img src=\"/inline.jpg\" alt=\"Inline Image|0\" />\n. Some more text.</p>",
			"<figure class=\"b\">\n\t<img src=\"/block.jpg\" alt=\"Block Image|1\" />",
		)
	})

	t.Run("With Hook, wrap", func(t *testing.T) {
		files := strings.ReplaceAll(filesTemplate, "CONFIG_VALUE", "true")
		files = files + `-- layouts/_default/_markup/render-image.html --
{{ if .IsBlock }}
<figure class="{{ .Attributes.class }}">
	<img src="{{ .Destination | safeURL }}" alt="{{ .Text }}" />
</figure>
{{ else }}
	<img src="{{ .Destination | safeURL }}" alt="{{ .Text }}" />
{{ end }}
`
		b := hugolib.Test(t, files)

		b.AssertFileContent("public/p1/index.html",
			"This is an inline image: \n\t<img src=\"/inline.jpg\" alt=\"Inline Image\" />\n. Some more text.</p>",
			"<p class=\"b\">\n\t<img src=\"/block.jpg\" alt=\"Block Image\" />\n</p>",
		)
	})

	t.Run("No Hook, no wrap", func(t *testing.T) {
		files := strings.ReplaceAll(filesTemplate, "CONFIG_VALUE", "false")
		b := hugolib.Test(t, files)

		b.AssertFileContent("public/p1/index.html", "<p>This is an inline image: <img src=\"/inline.jpg\" alt=\"Inline Image\">. Some more text.</p>\n<img src=\"/block.jpg\" alt=\"Block Image\" class=\"b\">")
	})

	t.Run("No Hook, wrap", func(t *testing.T) {
		files := strings.ReplaceAll(filesTemplate, "CONFIG_VALUE", "true")
		b := hugolib.Test(t, files)

		b.AssertFileContent("public/p1/index.html", "<p class=\"b\"><img src=\"/block.jpg\" alt=\"Block Image\"></p>")
	})
}
