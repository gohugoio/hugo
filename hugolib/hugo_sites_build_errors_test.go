package hugolib

import (
	"fmt"
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
	t.assert.NotNil(ferr, fmt.Sprintf("[%s] got %T: %+v", t.name, err, err))
	return ferr
}

func (t testSiteBuildErrorAsserter) assertLineNumber(lineNumber int, err error) {
	fe := t.getFileError(err)
	t.assert.Equal(lineNumber, fe.LineNumber, fmt.Sprintf("[%s]  got => %s", t.name, fe))
}

func TestSiteBuildErrors(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	const (
		yamlcontent = "yamlcontent"
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
				a.assertLineNumber(2, err)
			},
		},
		{
			name:     "Base template execute failed",
			fileType: base,
			fileFixer: func(content string) string {
				return strings.Replace(content, ".Title", ".Titles", 1)
			},
			assertBuildError: func(a testSiteBuildErrorAsserter, err error) {
				a.assertLineNumber(2, err)
			},
		},
		{
			name:     "Single template parse failed",
			fileType: single,
			fileFixer: func(content string) string {
				return strings.Replace(content, ".Title }}", ".Title }", 1)
			},
			assertCreateError: func(a testSiteBuildErrorAsserter, err error) {
				a.assertLineNumber(3, err)
			},
		},
		{
			name:     "Single template execute failed",
			fileType: single,
			fileFixer: func(content string) string {
				return strings.Replace(content, ".Title", ".Titles", 1)
			},
			assertBuildError: func(a testSiteBuildErrorAsserter, err error) {
				a.assertLineNumber(3, err)
			},
		},
		{
			name:     "Shortcode parse failed",
			fileType: shortcode,
			fileFixer: func(content string) string {
				return strings.Replace(content, ".Title }}", ".Title }", 1)
			},
			assertCreateError: func(a testSiteBuildErrorAsserter, err error) {
				a.assertLineNumber(2, err)
			},
		},
		// TODO(bep) 2errors
		/*		{
				name:     "Shortode execute failed",
				fileType: shortcode,
				fileFixer: func(content string) string {
					return strings.Replace(content, ".Title", ".Titles", 1)
				},
				assertBuildError: func(a testSiteBuildErrorAsserter, err error) {
					a.assertLineNumber(2, err)
				},
			},*/

		{
			name:     "Panic in template Execute",
			fileType: single,
			fileFixer: func(content string) string {
				return strings.Replace(content, ".Title", ".Parent.Parent.Parent", 1)
			},
			assertBuildError: func(a testSiteBuildErrorAsserter, err error) {
				assert.Error(err)
				assert.Contains(err.Error(), "layouts/_default/single.html")
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
