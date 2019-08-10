// Copyright 2019 The Hugo Authors. All rights reserved.
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

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/resources/resource"
)

// TODO(bep) move and rewrite in resource/page.

func TestMergeLanguages(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	b := newTestSiteForLanguageMerge(t, 30)
	b.CreateSites()

	b.Build(BuildCfg{SkipRender: true})

	h := b.H

	enSite := h.Sites[0]
	frSite := h.Sites[1]
	nnSite := h.Sites[2]

	c.Assert(len(enSite.RegularPages()), qt.Equals, 31)
	c.Assert(len(frSite.RegularPages()), qt.Equals, 6)
	c.Assert(len(nnSite.RegularPages()), qt.Equals, 12)

	for i := 0; i < 2; i++ {
		mergedNN := nnSite.RegularPages().MergeByLanguage(enSite.RegularPages())
		c.Assert(len(mergedNN), qt.Equals, 31)
		for i := 1; i <= 31; i++ {
			expectedLang := "en"
			if i == 2 || i%3 == 0 || i == 31 {
				expectedLang = "nn"
			}
			p := mergedNN[i-1]
			c.Assert(p.Language().Lang, qt.Equals, expectedLang)
		}
	}

	mergedFR := frSite.RegularPages().MergeByLanguage(enSite.RegularPages())
	c.Assert(len(mergedFR), qt.Equals, 31)
	for i := 1; i <= 31; i++ {
		expectedLang := "en"
		if i%5 == 0 {
			expectedLang = "fr"
		}
		p := mergedFR[i-1]
		c.Assert(p.Language().Lang, qt.Equals, expectedLang)
	}

	firstNN := nnSite.RegularPages()[0]
	c.Assert(len(firstNN.Sites()), qt.Equals, 4)
	c.Assert(firstNN.Sites().First().Language().Lang, qt.Equals, "en")

	nnBundle := nnSite.getPage("page", "bundle")
	enBundle := enSite.getPage("page", "bundle")

	c.Assert(len(enBundle.Resources()), qt.Equals, 6)
	c.Assert(len(nnBundle.Resources()), qt.Equals, 2)

	var ri interface{} = nnBundle.Resources()

	// This looks less ugly in the templates ...
	mergedNNResources := ri.(resource.ResourcesLanguageMerger).MergeByLanguage(enBundle.Resources())
	c.Assert(len(mergedNNResources), qt.Equals, 6)

	unchanged, err := nnSite.RegularPages().MergeByLanguageInterface(nil)
	c.Assert(err, qt.IsNil)
	c.Assert(unchanged, deepEqualsPages, nnSite.RegularPages())

}

func TestMergeLanguagesTemplate(t *testing.T) {
	t.Parallel()

	b := newTestSiteForLanguageMerge(t, 15)
	b.WithTemplates("home.html", `
{{ $pages := .Site.RegularPages }}
{{ .Scratch.Set "pages" $pages }}
{{ if eq .Language.Lang "nn" }}:
{{ $enSite := index .Sites 0 }}
{{ $frSite := index .Sites 1 }}
{{ $nnBundle := .Site.GetPage "page" "bundle" }}
{{ $enBundle := $enSite.GetPage "page" "bundle" }}
{{ .Scratch.Set "pages" ($pages | lang.Merge $frSite.RegularPages| lang.Merge $enSite.RegularPages) }}
{{ .Scratch.Set "pages2" (sort ($nnBundle.Resources | lang.Merge $enBundle.Resources) "Title") }}
{{ end }}
{{ $pages := .Scratch.Get "pages" }}
{{ $pages2 := .Scratch.Get "pages2" }}
Pages1: {{ range $i, $p := $pages }}{{ add $i 1 }}: {{ .File.Path }} {{ .Language.Lang }} | {{ end }}
Pages2: {{ range $i, $p := $pages2 }}{{ add $i 1 }}: {{ .Title }} {{ .Language.Lang }} | {{ end }}

`,
		"shortcodes/shortcode.html", "MyShort",
		"shortcodes/lingo.html", "MyLingo",
	)

	b.CreateSites()
	b.Build(BuildCfg{})

	b.AssertFileContent("public/nn/index.html", "Pages1: 1: p1.md en | 2: p2.nn.md nn | 3: p3.nn.md nn | 4: p4.md en | 5: p5.fr.md fr | 6: p6.nn.md nn | 7: p7.md en | 8: p8.md en | 9: p9.nn.md nn | 10: p10.fr.md fr | 11: p11.md en | 12: p12.nn.md nn | 13: p13.md en | 14: p14.md en | 15: p15.nn.md nn")
	b.AssertFileContent("public/nn/index.html", "Pages2: 1: doc100 en | 2: doc101 nn | 3: doc102 nn | 4: doc103 en | 5: doc104 en | 6: doc105 en")
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

	// See https://github.com/gohugoio/hugo/issues/4644
	// Add a bundles
	j := 100
	contentPairs = append(contentPairs, []string{"bundle/index.md", fmt.Sprintf(contentTemplate, j, j)}...)
	for i := 0; i < 6; i++ {
		contentPairs = append(contentPairs, []string{fmt.Sprintf("bundle/pb%d.md", i), fmt.Sprintf(contentTemplate, i+j, i+j)}...)
	}
	contentPairs = append(contentPairs, []string{"bundle/index.nn.md", fmt.Sprintf(contentTemplate, j, j)}...)
	for i := 1; i < 3; i++ {
		contentPairs = append(contentPairs, []string{fmt.Sprintf("bundle/pb%d.nn.md", i), fmt.Sprintf(contentTemplate, i+j, i+j)}...)
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
		merged := nnSite.RegularPages().MergeByLanguage(enSite.RegularPages())
		if len(merged) != count {
			b.Fatal("Count mismatch")
		}
	}
}
