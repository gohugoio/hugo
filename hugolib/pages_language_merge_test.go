// Copyright 2018 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hugolib

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMergeLanguages(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	b := newTestSiteForLanguageMerge(t, 30)
	b.CreateSites()

	b.Build(BuildCfg{SkipRender: true})

	h := b.H

	enSite := h.Sites[0]
	frSite := h.Sites[1]
	nnSite := h.Sites[2]

	assert.Equal(30, len(enSite.RegularPages))
	assert.Equal(6, len(frSite.RegularPages))
	assert.Equal(11, len(nnSite.RegularPages))

	for i := 0; i < 2; i++ {
		mergedNN := nnSite.RegularPages.MergeByLanguage(enSite.RegularPages)
		assert.Equal(30, len(mergedNN))
		for i := 1; i <= 30; i++ {
			expectedLang := "en"
			if i == 2 || i%3 == 0 {
				expectedLang = "nn"
			}
			p := mergedNN[i-1]
			assert.Equal(expectedLang, p.Lang(), fmt.Sprintf("Test %d", i))
		}
	}

	mergedFR := frSite.RegularPages.MergeByLanguage(enSite.RegularPages)
	assert.Equal(30, len(mergedFR))
	for i := 1; i <= 30; i++ {
		expectedLang := "en"
		if i%5 == 0 {
			expectedLang = "fr"
		}
		p := mergedFR[i-1]
		assert.Equal(expectedLang, p.Lang(), fmt.Sprintf("Test %d", i))
	}

	firstNN := nnSite.RegularPages[0]
	assert.Equal(4, len(firstNN.Sites()))
	assert.Equal("en", firstNN.Sites().First().Language.Lang)

}

func TestMergeLanguagesTemplate(t *testing.T) {
	t.Parallel()

	b := newTestSiteForLanguageMerge(t, 15)
	b.WithTemplates("home.html", `
{{ $pages := .Site.RegularPages }}
{{ .Scratch.Set "pages" $pages }}
{{ if eq .Lang "nn" }}:
{{ $enSite := index .Sites 0 }}
{{ $frSite := index .Sites 1 }}
{{ .Scratch.Set "pages" ($pages | lang.Merge $frSite.RegularPages| lang.Merge $enSite.RegularPages) }}
{{ end }}
{{ $pages := .Scratch.Get "pages" }}
{{ range $i, $p := $pages }}{{ add $i 1 }}: {{ .Path }} {{ .Lang }} | {{ end }}
`,
		"shortcodes/shortcode.html", "MyShort",
		"shortcodes/lingo.html", "MyLingo",
	)

	b.CreateSites()
	b.Build(BuildCfg{})

	b.AssertFileContent("public/nn/index.html", "p1.md en | 2: p2.nn.md nn | 3: p3.nn.md nn | 4: p4.md en | 5: p5.fr.md fr | 6: p6.nn.md nn | 7: p7.md en | 8: p8.md en | 9: p9.nn.md nn | 10: p10.fr.md fr | 11: p11.md en | 12: p12.nn.md nn | 13: p13.md en | 14: p14.md en | 15: p15.nn.md nn")
}

func newTestSiteForLanguageMerge(t testing.TB, count int) *sitesBuilder {
	contentTemplate := `---
title: doc%d
weight: %d
date: "2018-02-28"
---
# doc
*some "content"*

{{< shortcode >}}

{{< lingo >}}
`

	builder := newTestSitesBuilder(t).WithDefaultMultiSiteConfig()

	// We need some content with some missing translations.
	// "en" is the main language, so add some English content + some Norwegian (nn, nynorsk) content.
	var contentPairs []string
	for i := 1; i <= count; i++ {
		content := fmt.Sprintf(contentTemplate, i, i)
		contentPairs = append(contentPairs, []string{fmt.Sprintf("p%d.md", i), content}...)
		if i == 2 || i%3 == 0 {
			// Add page 2,3, 6, 9 ... to both languages
			contentPairs = append(contentPairs, []string{fmt.Sprintf("p%d.nn.md", i), content}...)
		}
		if i%5 == 0 {
			// Add some French content, too.
			contentPairs = append(contentPairs, []string{fmt.Sprintf("p%d.fr.md", i), content}...)
		}
	}

	builder.WithContent(contentPairs...)
	return builder
}

func BenchmarkMergeByLanguage(b *testing.B) {
	const count = 100

	builder := newTestSiteForLanguageMerge(b, count)
	builder.CreateSites()
	builder.Build(BuildCfg{SkipRender: true})
	h := builder.H

	enSite := h.Sites[0]
	nnSite := h.Sites[2]

	for i := 0; i < b.N; i++ {
		merged := nnSite.RegularPages.MergeByLanguage(enSite.RegularPages)
		if len(merged) != count {
			b.Fatal("Count mismatch")
		}
	}
}
