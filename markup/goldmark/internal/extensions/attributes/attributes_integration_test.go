package attributes_test

import (
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

func TestDescriptionListAutoID(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
[markup.goldmark.parser]
autoHeadingID = true
autoDefinitionTermID = true
autoIDType = 'github-ascii'
-- content/p1.md --
---
title: "Title"
---

## Title with id set {#title-with-id}

## Title with id set duplicate {#title-with-id}

## My Title

Base Name
: Base name of the file.

Base Name
: Duplicate term name.

My Title
: Term with same name as title.

Foo@Bar
: The foo bar.

foo [something](/a/b/) bar
: A foo bar.

良善天父
: The good father.

Ā ā Ă ă Ą ą Ć ć Ĉ ĉ Ċ ċ Č č Ď
: Testing accents.

Multiline set text header
Second line
---------------

## Example [hyperlink](https://example.com/) in a header

-- layouts/_default/single.html --
{{ .Content }}|Identifiers: {{ .Fragments.Identifiers }}|
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/p1/index.html",
		`<dt id="base-name">Base Name</dt>`,
		`<dt id="base-name-1">Base Name</dt>`,
		`<dt id="foobar">Foo@Bar</dt>`,
		`<h2 id="my-title">My Title</h2>`,
		`<dt id="foo-something-bar">foo <a href="/a/b/">something</a> bar</dt>`,
		`<h2 id="title-with-id">Title with id set</h2>`,
		`<h2 id="title-with-id">Title with id set duplicate</h2>`,
		`<dt id="my-title-1">My Title</dt>`,
		`<dt id="term">良善天父</dt>`,
		`<dt id="a-a-a-a-a-a-c-c-c-c-c-c-c-c-d">Ā ā Ă ă Ą ą Ć ć Ĉ ĉ Ċ ċ Č č Ď</dt>`,
		`<h2 id="second-line">`,
		`<h2 id="example-hyperlink-in-a-header">`,
		"|Identifiers: [a-a-a-a-a-a-c-c-c-c-c-c-c-c-d base-name base-name-1 example-hyperlink-in-a-header foo-something-bar foobar my-title my-title-1 second-line term title-with-id title-with-id]|",
	)
}
