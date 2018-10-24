package hugolib

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/stretchr/testify/require"
)

type testSiteBuildErrorAsserter struct {
	name   string
	assert *require.Assertions
}

func (t testSiteBuildErrorAsserter) getFileError(err error) *herrors.ErrorWithFileContext {
	t.assert.NotNil(err, t.name)
	ferr := herrors.UnwrapErrorWithFileContext(err)
	t.assert.NotNil(ferr, fmt.Sprintf("[%s] got %T: %+v\n%s", t.name, err, err, trace()))
	return ferr
}

func (t testSiteBuildErrorAsserter) assertLineNumber(lineNumber int, err error) {
	fe := t.getFileError(err)
	t.assert.Equal(lineNumber, fe.LineNumber, fmt.Sprintf("[%s]  got => %s\n%s", t.name, fe, trace()))
}

func (t testSiteBuildErrorAsserter) assertErrorMessage(e1, e2 string) {
	// The error message will contain filenames with OS slashes. Normalize before compare.
	e1, e2 = filepath.ToSlash(e1), filepath.ToSlash(e2)
	t.assert.Contains(e2, e1, trace())

}

func TestSiteBuildErrors(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	const (
		yamlcontent = "yamlcontent"
		tomlcontent = "tomlcontent"
		shortcode   = "shortcode"
		base        = "base"
		single      = "single"
	)

	// TODO(bep) add content tests after https://github.com/gohugoio/hugo/issues/5324
	// is implemented.

	tests := []struct {
		name              string
		fileType          string
		fileFixer         func(content string) string
		assertCreateError func(a testSiteBuildErrorAsserter, err error)
		assertBuildError  func(a testSiteBuildErrorAsserter, err error)
	}{

		{
			name:     "Base template parse failed",
			fileType: base,
			fileFixer: func(content string) string {
				return strings.Replace(content, ".Title }}", ".Title }", 1)
			},
			assertCreateError: func(a testSiteBuildErrorAsserter, err error) {
				a.assertLineNumber(4, err)
			},
		},
		{
			name:     "Base template execute failed",
			fileType: base,
			fileFixer: func(content string) string {
				return strings.Replace(content, ".Title", ".Titles", 1)
			},
			assertBuildError: func(a testSiteBuildErrorAsserter, err error) {
				a.assertLineNumber(4, err)
			},
		},
		{
			name:     "Single template parse failed",
			fileType: single,
			fileFixer: func(content string) string {
				return strings.Replace(content, ".Title }}", ".Title }", 1)
			},
			assertCreateError: func(a testSiteBuildErrorAsserter, err error) {
				fe := a.getFileError(err)
				assert.Equal(5, fe.LineNumber)
				assert.Equal(1, fe.ColumnNumber)
				assert.Equal("go-html-template", fe.ChromaLexer)
				a.assertErrorMessage("\"layouts/_default/single.html:5:1\": parse failed: template: _default/single.html:5: unexpected \"}\" in operand", fe.Error())

			},
		},
		{
			name:     "Single template execute failed",
			fileType: single,
			fileFixer: func(content string) string {
				return strings.Replace(content, ".Title", ".Titles", 1)
			},
			assertBuildError: func(a testSiteBuildErrorAsserter, err error) {
				fe := a.getFileError(err)
				assert.Equal(5, fe.LineNumber)
				assert.Equal(14, fe.ColumnNumber)
				assert.Equal("go-html-template", fe.ChromaLexer)
				a.assertErrorMessage("\"layouts/_default/single.html:5:14\": execute of template failed", fe.Error())

			},
		},
		{
			name:     "Single template execute failed, long keyword",
			fileType: single,
			fileFixer: func(content string) string {
				return strings.Replace(content, ".Title", ".ThisIsAVeryLongTitle", 1)
			},
			assertBuildError: func(a testSiteBuildErrorAsserter, err error) {
				fe := a.getFileError(err)
				assert.Equal(5, fe.LineNumber)
				assert.Equal(14, fe.ColumnNumber)
				assert.Equal("go-html-template", fe.ChromaLexer)
				a.assertErrorMessage("\"layouts/_default/single.html:5:14\": execute of template failed", fe.Error())

			},
		},
		{
			name:     "Shortcode parse failed",
			fileType: shortcode,
			fileFixer: func(content string) string {
				return strings.Replace(content, ".Title }}", ".Title }", 1)
			},
			assertCreateError: func(a testSiteBuildErrorAsserter, err error) {
				a.assertLineNumber(4, err)
			},
		},
		{
			name:     "Shortode execute failed",
			fileType: shortcode,
			fileFixer: func(content string) string {
				return strings.Replace(content, ".Title", ".Titles", 1)
			},
			assertBuildError: func(a testSiteBuildErrorAsserter, err error) {
				fe := a.getFileError(err)
				assert.Equal(7, fe.LineNumber)
				assert.Equal("md", fe.ChromaLexer)
				// Make sure that it contains both the content file and template
				a.assertErrorMessage(`content/myyaml.md:7:10": failed to render shortcode "sc"`, fe.Error())
				a.assertErrorMessage(`shortcodes/sc.html:4:22: executing "shortcodes/sc.html" at <.Page.Titles>: can't evaluate`, fe.Error())
			},
		},
		{
			name:     "Shortode does not exist",
			fileType: yamlcontent,
			fileFixer: func(content string) string {
				return strings.Replace(content, "{{< sc >}}", "{{< nono >}}", 1)
			},
			assertBuildError: func(a testSiteBuildErrorAsserter, err error) {
				fe := a.getFileError(err)
				assert.Equal(7, fe.LineNumber)
				assert.Equal(14, fe.ColumnNumber)
				assert.Equal("md", fe.ChromaLexer)
				a.assertErrorMessage("\"content/myyaml.md:7:14\": failed to extract shortcode: template for shortcode \"nono\" not found", fe.Error())
			},
		},
		{
			name:     "Invalid YAML front matter",
			fileType: yamlcontent,
			fileFixer: func(content string) string {
				return strings.Replace(content, "title:", "title: %foo", 1)
			},
			assertBuildError: func(a testSiteBuildErrorAsserter, err error) {
				a.assertLineNumber(2, err)
			},
		},
		{
			name:     "Invalid TOML front matter",
			fileType: tomlcontent,
			fileFixer: func(content string) string {
				return strings.Replace(content, "description = ", "description &", 1)
			},
			assertBuildError: func(a testSiteBuildErrorAsserter, err error) {
				fe := a.getFileError(err)
				assert.Equal(6, fe.LineNumber)
				assert.Equal("toml", fe.ErrorContext.ChromaLexer)

			},
		},
		{
			name:     "Invalid JSON front matter",
			fileType: tomlcontent,
			fileFixer: func(content string) string {
				return strings.Replace(content, "\"description\":", "\"description\"", 1)
			},
			assertBuildError: func(a testSiteBuildErrorAsserter, err error) {
				fe := a.getFileError(err)

				assert.Equal(3, fe.LineNumber)
				assert.Equal("json", fe.ErrorContext.ChromaLexer)

			},
		},
		{
			// See https://github.com/gohugoio/hugo/issues/5327
			name:     "Panic in template Execute",
			fileType: single,
			fileFixer: func(content string) string {
				return strings.Replace(content, ".Title", ".Parent.Parent.Parent", 1)
			},

			assertBuildError: func(a testSiteBuildErrorAsserter, err error) {
				assert.Error(err)
				// This is fixed in latest Go source
				if strings.Contains(runtime.Version(), "devel") {
					fe := a.getFileError(err)
					assert.Equal(5, fe.LineNumber)
					assert.Equal(21, fe.ColumnNumber)
				} else {
					assert.Contains(err.Error(), `execute of template failed: panic in Execute`)
				}
			},
		},
	}

	for _, test := range tests {

		errorAsserter := testSiteBuildErrorAsserter{
			assert: assert,
			name:   test.name,
		}

		b := newTestSitesBuilder(t).WithSimpleConfigFile()

		f := func(fileType, content string) string {
			if fileType != test.fileType {
				return content
			}
			return test.fileFixer(content)

		}

		b.WithTemplatesAdded("layouts/shortcodes/sc.html", f(shortcode, `SHORTCODE L1
SHORTCODE L2
SHORTCODE L3:
SHORTCODE L4: {{ .Page.Title }}
`))
		b.WithTemplatesAdded("layouts/_default/baseof.html", f(base, `BASEOF L1
BASEOF L2
BASEOF L3
BASEOF L4{{ if .Title }}{{ end }}
{{block "main" .}}This is the main content.{{end}}
BASEOF L6
`))

		b.WithTemplatesAdded("layouts/_default/single.html", f(single, `{{ define "main" }}
SINGLE L2:
SINGLE L3:
SINGLE L4:
SINGLE L5: {{ .Title }} {{ .Content }}
{{ end }}
`))

		b.WithContent("myyaml.md", f(yamlcontent, `---
title: "The YAML"
---

Some content.

         {{< sc >}}

Some more text.

The end.

`))

		b.WithContent("mytoml.md", f(tomlcontent, `+++
title = "The TOML"
p1 = "v"
p2 = "v"
p3 = "v"
description = "Descriptioon"
+++

Some content.


`))

		b.WithContent("myjson.md", f(tomlcontent, `{
	"title": "This is a title",
	"description": "This is a description."
}

Some content.


`))

		createErr := b.CreateSitesE()
		if test.assertCreateError != nil {
			test.assertCreateError(errorAsserter, createErr)
		} else {
			assert.NoError(createErr)
		}

		if createErr == nil {
			buildErr := b.BuildE(BuildCfg{})
			if test.assertBuildError != nil {
				test.assertBuildError(errorAsserter, buildErr)
			} else {
				assert.NoError(buildErr)
			}
		}
	}
}
