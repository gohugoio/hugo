package tplimpl_test

import (
	"path/filepath"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/hugolib"
	"github.com/gohugoio/hugo/tpl"
)

func TestPrintUnusedTemplates(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
baseURL = 'http://example.com/'
printUnusedTemplates=true
-- content/p1.md --
---
title: "P1"
---
{{< usedshortcode >}}
-- layouts/baseof.html --
{{ block "main" . }}{{ end }}
-- layouts/baseof.json --
{{ block "main" . }}{{ end }}
-- layouts/index.html --
{{ define "main" }}FOO{{ end }}
-- layouts/_default/single.json --
-- layouts/_default/single.html --
{{ define "main" }}MAIN{{ end }}
-- layouts/post/single.html --
{{ define "main" }}MAIN{{ end }}
-- layouts/partials/usedpartial.html --
-- layouts/partials/unusedpartial.html --
-- layouts/shortcodes/usedshortcode.html --
{{ partial "usedpartial.html" }}
-- layouts/shortcodes/unusedshortcode.html --

	`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			NeedsOsFS:   true,
		},
	)
	b.Build()

	unused := b.H.Tmpl().(tpl.UnusedTemplatesProvider).UnusedTemplates()

	var names []string
	for _, tmpl := range unused {
		names = append(names, tmpl.Name())
	}

	b.Assert(names, qt.DeepEquals, []string{"_default/single.json", "baseof.json", "partials/unusedpartial.html", "post/single.html", "shortcodes/unusedshortcode.html"})
	b.Assert(unused[0].Filename(), qt.Equals, filepath.Join(b.Cfg.WorkingDir, "layouts/_default/single.json"))
}
