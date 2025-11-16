// Copyright 2025 The Hugo Authors. All rights reserved.
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
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/resources/resource"
	"github.com/google/go-cmp/cmp"
)

var deepEqualsPages = qt.CmpEquals(cmp.Comparer(func(p1, p2 *pageState) bool { return p1 == p2 }))

func TestMergeLanguages(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	files := generateLanguageMergeTxtar(30)
	b := Test(t, files)

	h := b.H

	enSite := h.Sites[0]
	frSite := h.Sites[1]
	nnSite := h.Sites[2]

	c.Assert(len(enSite.RegularPages()), qt.Equals, 31)
	c.Assert(len(frSite.RegularPages()), qt.Equals, 6)
	c.Assert(len(nnSite.RegularPages()), qt.Equals, 12)

	for range 2 {
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
	c.Assert(firstNN.Sites().Default().Language().Lang, qt.Equals, "en")

	nnBundle := nnSite.getPageOldVersion("page", "bundle")
	enBundle := enSite.getPageOldVersion("page", "bundle")

	c.Assert(len(enBundle.Resources()), qt.Equals, 6)
	c.Assert(len(nnBundle.Resources()), qt.Equals, 2)

	var ri any = nnBundle.Resources()

	// This looks less ugly in the templates ...
	mergedNNResources := ri.(resource.ResourcesLanguageMerger).MergeByLanguage(enBundle.Resources())
	c.Assert(len(mergedNNResources), qt.Equals, 6)

	unchanged, err := nnSite.RegularPages().MergeByLanguageInterface(nil)
	c.Assert(err, qt.IsNil)
	c.Assert(unchanged, deepEqualsPages, nnSite.RegularPages())
}

func TestMergeLanguagesTemplate(t *testing.T) {
	t.Parallel()

	files := generateLanguageMergeTxtar(15) + `
-- layouts/home.html --
{{ $pages := .Site.RegularPages }}
{{ .Scratch.Set "pages" $pages }}
{{ $enSite := index .Sites 0 }}
{{ $frSite := index .Sites 1 }}
{{ if eq .Language.Lang "nn" }}:
{{ $nnBundle := .Site.GetPage "page" "bundle" }}
{{ $enBundle := $enSite.GetPage "page" "bundle" }}
{{ .Scratch.Set "pages" ($pages | lang.Merge $frSite.RegularPages| lang.Merge $enSite.RegularPages) }}
{{ .Scratch.Set "pages2" (sort ($nnBundle.Resources | lang.Merge $enBundle.Resources) "Title") }}
{{ end }}
{{ $pages := .Scratch.Get "pages" }}
{{ $pages2 := .Scratch.Get "pages2" }}
Pages1: {{ range $i, $p := $pages }}{{ add $i 1 }}: {{ .File.Path }} {{ .Language.Lang }} | {{ end }}
Pages2: {{ range $i, $p := $pages2 }}{{ add $i 1 }}: {{ .Title }} {{ .Language.Lang }} | {{ end }}
{{ $nil := resources.Get "asdfasdfasdf" }}
Pages3: {{ $frSite.RegularPages | lang.Merge  $nil }}
Pages4: {{  $nil | lang.Merge $frSite.RegularPages }}
-- layouts/shortcodes/shortcode.html --
MyShort
-- layouts/shortcodes/lingo.html --
MyLingo
`
	b := Test(t, files)

	b.AssertFileContent("public/nn/index.html", "Pages1: 1: p1.md en | 2: p2.nn.md nn | 3: p3.nn.md nn | 4: p4.md en | 5: p5.fr.md fr | 6: p6.nn.md nn | 7: p7.md en | 8: p8.md en | 9: p9.nn.md nn | 10: p10.fr.md fr | 11: p11.md en | 12: p12.nn.md nn | 13: p13.md en | 14: p14.md en | 15: p15.nn.md nn")
	b.AssertFileContent("public/nn/index.html", "Pages2: 1: doc100 en | 2: doc101 nn | 3: doc102 nn | 4: doc103 en | 5: doc104 en | 6: doc105 en")
	b.AssertFileContent("public/nn/index.html", `
Pages3: Pages(3)
Pages4: Pages(3)
	`)
}

func generateLanguageMergeTxtar(count int) string {
	var b strings.Builder

	// hugo.toml for multisite configuration
	b.WriteString(`
-- hugo.toml --
baseURL = "https://example.org"
defaultContentLanguage = "en"

[languages]
[languages.en]
title = "English"
weight = 1
[languages.fr]
title = "French"
weight = 2
[languages.nn]
title = "Nynorsk"
weight = 3
[languages.no]
title = "Norwegian"
weight = 4

`)

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

	// Generate content files
	for i := 1; i <= count; i++ {
		content := fmt.Sprintf(contentTemplate, i, i)
		b.WriteString(fmt.Sprintf("\n-- content/p%d.md --\n%s", i, content))
		if i == 2 || i%3 == 0 {
			b.WriteString(fmt.Sprintf("\n-- content/p%d.nn.md --\n%s", i, content))
		}
		if i%5 == 0 {
			b.WriteString(fmt.Sprintf("\n-- content/p%d.fr.md --\n%s", i, content))
		}
	}

	// Add a bundles
	j := 100
	b.WriteString(fmt.Sprintf("\n-- content/bundle/index.md --\n%s", fmt.Sprintf(contentTemplate, j, j)))
	for i := range 6 {
		b.WriteString(fmt.Sprintf("\n-- content/bundle/pb%d.md --\n%s", i, fmt.Sprintf(contentTemplate, i+j, i+j)))
	}
	b.WriteString(fmt.Sprintf("\n-- content/bundle/index.nn.md --\n%s", fmt.Sprintf(contentTemplate, j, j)))
	for i := 1; i < 3; i++ {
		b.WriteString(fmt.Sprintf("\n-- content/bundle/pb%d.nn.md --\n%s", i, fmt.Sprintf(contentTemplate, i+j, i+j)))
	}

	// Add shortcode templates
	b.WriteString(`
-- layouts/shortcodes/shortcode.html --
MyShort
-- layouts/shortcodes/lingo.html --
MyLingo
`)

	return b.String()
}

func BenchmarkMergeByLanguage(b *testing.B) {
	const count = 100

	files := generateLanguageMergeTxtar(count - 1)
	builder := Test(b, files)
	h := builder.H

	enSite := h.Sites[0]
	nnSite := h.Sites[2]

	for b.Loop() {
		merged := nnSite.RegularPages().MergeByLanguage(enSite.RegularPages())
		if len(merged) != count {
			b.Fatal("Count mismatch")
		}
	}
}
