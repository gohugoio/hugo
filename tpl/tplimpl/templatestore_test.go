package tplimpl_test

import (
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

func TestIssue13877(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['home','rss','section','sitemap','taxonomy','term']

[mediaTypes.'text/html']
suffixes = ['b','a','d','c']

[outputFormats.html]
mediaType = 'text/html'
-- content/p1.md --
---
title: p1
---
-- layouts/page.html.a --
{{ templates.Current.Name }}
-- layouts/page.html.b --
{{ templates.Current.Name }}
-- layouts/page.html.c --
{{ templates.Current.Name }}
-- layouts/page.html.d --
{{ templates.Current.Name }}
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/p1/index.b", "page.html.b")
}
