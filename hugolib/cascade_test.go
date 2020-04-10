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
	"bytes"
	"fmt"
	"path"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/parser"
	"github.com/gohugoio/hugo/parser/metadecoders"
)

func BenchmarkCascade(b *testing.B) {
	allLangs := []string{"en", "nn", "nb", "sv", "ab", "aa", "af", "sq", "kw", "da"}

	for i := 1; i <= len(allLangs); i += 2 {
		langs := allLangs[0:i]
		b.Run(fmt.Sprintf("langs-%d", len(langs)), func(b *testing.B) {
			c := qt.New(b)
			b.StopTimer()
			builders := make([]*sitesBuilder, b.N)
			for i := 0; i < b.N; i++ {
				builders[i] = newCascadeTestBuilder(b, langs)
			}
			b.StartTimer()

			for i := 0; i < b.N; i++ {
				builder := builders[i]
				err := builder.BuildE(BuildCfg{})
				c.Assert(err, qt.IsNil)
				first := builder.H.Sites[0]
				c.Assert(first, qt.Not(qt.IsNil))
			}
		})
	}
}

func TestCascade(t *testing.T) {

	allLangs := []string{"en", "nn", "nb", "sv"}

	langs := allLangs[:3]

	t.Run(fmt.Sprintf("langs-%d", len(langs)), func(t *testing.T) {
		b := newCascadeTestBuilder(t, langs)
		b.Build(BuildCfg{})

		b.AssertFileContent("public/index.html", `
12|taxonomy|categories/cool/_index.md|Cascade Category|cat.png|categories|HTML-|
12|taxonomy|categories/catsect1|catsect1|cat.png|categories|HTML-|
12|taxonomy|categories/funny|funny|cat.png|categories|HTML-|
12|taxonomyTerm|categories/_index.md|My Categories|cat.png|categories|HTML-|
32|taxonomy|categories/sad/_index.md|Cascade Category|sad.png|categories|HTML-|
42|taxonomy|tags/blue|blue|home.png|tags|HTML-|
42|taxonomyTerm|tags|Cascade Home|home.png|tags|HTML-|
42|section|sectnocontent|Cascade Home|home.png|sectnocontent|HTML-|
42|section|sect3|Cascade Home|home.png|sect3|HTML-|
42|page|bundle1/index.md|Cascade Home|home.png|page|HTML-|
42|page|p2.md|Cascade Home|home.png|page|HTML-|
42|page|sect2/p2.md|Cascade Home|home.png|sect2|HTML-|
42|page|sect3/nofrontmatter.md|Cascade Home|home.png|sect3|HTML-|
42|page|sect3/p1.md|Cascade Home|home.png|sect3|HTML-|
42|page|sectnocontent/p1.md|Cascade Home|home.png|sectnocontent|HTML-|
42|section|sectnofrontmatter/_index.md|Cascade Home|home.png|sectnofrontmatter|HTML-|
42|taxonomy|tags/green|green|home.png|tags|HTML-|
42|home|_index.md|Home|home.png|page|HTML-|
42|page|p1.md|p1|home.png|page|HTML-|
42|section|sect1/_index.md|Sect1|sect1.png|stype|HTML-|
42|section|sect1/s1_2/_index.md|Sect1_2|sect1.png|stype|HTML-|
42|page|sect1/s1_2/p1.md|Sect1_2_p1|sect1.png|stype|HTML-|
42|page|sect1/s1_2/p2.md|Sect1_2_p2|sect1.png|stype|HTML-|
42|section|sect2/_index.md|Sect2|home.png|sect2|HTML-|
42|page|sect2/p1.md|Sect2_p1|home.png|sect2|HTML-|
52|page|sect4/p1.md|Cascade Home|home.png|sect4|RSS-|
52|section|sect4/_index.md|Sect4|home.png|sect4|RSS-|
`)

		// Check that type set in cascade gets the correct layout.
		b.AssertFileContent("public/sect1/index.html", `stype list: Sect1`)
		b.AssertFileContent("public/sect1/s1_2/p2/index.html", `stype single: Sect1_2_p2`)

		// Check output formats set in cascade
		b.AssertFileContent("public/sect4/index.xml", `<link>https://example.org/sect4/index.xml</link>`)
		b.AssertFileContent("public/sect4/p1/index.xml", `<link>https://example.org/sect4/p1/index.xml</link>`)
		b.C.Assert(b.CheckExists("public/sect2/index.xml"), qt.Equals, false)

		// Check cascade into bundled page
		b.AssertFileContent("public/bundle1/index.html", `Resources: bp1.md|home.png|`)

	})

}

func TestCascadeEdit(t *testing.T) {
	p1Content := `---
title: P1
---
`

	indexContentNoCascade := `
---
title: Home
---
`

	indexContentCascade := `
---
title: Section
cascade:
  banner: post.jpg
  layout: postlayout
  type: posttype
---
`

	layout := `Banner: {{ .Params.banner }}|Layout: {{ .Layout }}|Type: {{ .Type }}|Content: {{ .Content }}`

	newSite := func(t *testing.T, cascade bool) *sitesBuilder {
		b := newTestSitesBuilder(t).Running()
		b.WithTemplates("_default/single.html", layout)
		b.WithTemplates("_default/list.html", layout)
		if cascade {
			b.WithContent("post/_index.md", indexContentCascade)
		} else {
			b.WithContent("post/_index.md", indexContentNoCascade)
		}
		b.WithContent("post/dir/p1.md", p1Content)

		return b
	}

	t.Run("Edit descendant", func(t *testing.T) {
		t.Parallel()

		b := newSite(t, true)
		b.Build(BuildCfg{})

		assert := func() {
			b.Helper()
			b.AssertFileContent("public/post/dir/p1/index.html",
				`Banner: post.jpg|`,
				`Layout: postlayout`,
				`Type: posttype`,
			)
		}

		assert()

		b.EditFiles("content/post/dir/p1.md", p1Content+"\ncontent edit")
		b.Build(BuildCfg{})

		assert()
		b.AssertFileContent("public/post/dir/p1/index.html",
			`content edit
Banner: post.jpg`,
		)
	})

	t.Run("Edit ancestor", func(t *testing.T) {
		t.Parallel()

		b := newSite(t, true)
		b.Build(BuildCfg{})

		b.AssertFileContent("public/post/dir/p1/index.html", `Banner: post.jpg|Layout: postlayout|Type: posttype|Content:`)

		b.EditFiles("content/post/_index.md", strings.Replace(indexContentCascade, "post.jpg", "edit.jpg", 1))

		b.Build(BuildCfg{})

		b.AssertFileContent("public/post/index.html", `Banner: edit.jpg|Layout: postlayout|Type: posttype|`)
		b.AssertFileContent("public/post/dir/p1/index.html", `Banner: edit.jpg|Layout: postlayout|Type: posttype|`)
	})

	t.Run("Edit ancestor, add cascade", func(t *testing.T) {
		t.Parallel()

		b := newSite(t, true)
		b.Build(BuildCfg{})

		b.AssertFileContent("public/post/dir/p1/index.html", `Banner: post.jpg`)

		b.EditFiles("content/post/_index.md", indexContentCascade)

		b.Build(BuildCfg{})

		b.AssertFileContent("public/post/index.html", `Banner: post.jpg|Layout: postlayout|Type: posttype|`)
		b.AssertFileContent("public/post/dir/p1/index.html", `Banner: post.jpg|Layout: postlayout|`)
	})

	t.Run("Edit ancestor, remove cascade", func(t *testing.T) {
		t.Parallel()

		b := newSite(t, false)
		b.Build(BuildCfg{})

		b.AssertFileContent("public/post/dir/p1/index.html", `Banner: |Layout: |`)

		b.EditFiles("content/post/_index.md", indexContentNoCascade)

		b.Build(BuildCfg{})

		b.AssertFileContent("public/post/index.html", `Banner: |Layout: |Type: post|`)
		b.AssertFileContent("public/post/dir/p1/index.html", `Banner: |Layout: |`)
	})

	t.Run("Edit ancestor, content only", func(t *testing.T) {
		t.Parallel()

		b := newSite(t, true)
		b.Build(BuildCfg{})

		b.EditFiles("content/post/_index.md", indexContentCascade+"\ncontent edit")

		counters := &testCounters{}
		b.Build(BuildCfg{testCounters: counters})
		// As we only changed the content, not the cascade front matter, make
		// only the home page is re-rendered.
		b.Assert(int(counters.contentRenderCounter), qt.Equals, 1)

		b.AssertFileContent("public/post/index.html", `Banner: post.jpg|Layout: postlayout|Type: posttype|Content: <p>content edit</p>`)
		b.AssertFileContent("public/post/dir/p1/index.html", `Banner: post.jpg|Layout: postlayout|`)
	})
}

func newCascadeTestBuilder(t testing.TB, langs []string) *sitesBuilder {
	p := func(m map[string]interface{}) string {
		var yamlStr string

		if len(m) > 0 {
			var b bytes.Buffer

			parser.InterfaceToConfig(m, metadecoders.YAML, &b)
			yamlStr = b.String()
		}

		metaStr := "---\n" + yamlStr + "\n---"

		return metaStr

	}

	createLangConfig := func(lang string) string {
		const langEntry = `
[languages.%s]
`
		return fmt.Sprintf(langEntry, lang)
	}

	createMount := func(lang string) string {
		const mountsTempl = `
[[module.mounts]]
source="content/%s"
target="content"
lang="%s"
`
		return fmt.Sprintf(mountsTempl, lang, lang)
	}

	config := `
baseURL = "https://example.org"
defaultContentLanguage = "en"
defaultContentLanguageInSubDir = false

[languages]`
	for _, lang := range langs {
		config += createLangConfig(lang)
	}

	config += "\n\n[module]\n"
	for _, lang := range langs {
		config += createMount(lang)
	}

	b := newTestSitesBuilder(t).WithConfigFile("toml", config)

	createContentFiles := func(lang string) {

		withContent := func(filenameContent ...string) {
			for i := 0; i < len(filenameContent); i += 2 {
				b.WithContent(path.Join(lang, filenameContent[i]), filenameContent[i+1])
			}
		}

		withContent(
			"_index.md", p(map[string]interface{}{
				"title": "Home",
				"cascade": map[string]interface{}{
					"title":   "Cascade Home",
					"ICoN":    "home.png",
					"outputs": []string{"HTML"},
					"weight":  42,
				},
			}),
			"p1.md", p(map[string]interface{}{
				"title": "p1",
			}),
			"p2.md", p(map[string]interface{}{}),
			"sect1/_index.md", p(map[string]interface{}{
				"title": "Sect1",
				"type":  "stype",
				"cascade": map[string]interface{}{
					"title":      "Cascade Sect1",
					"icon":       "sect1.png",
					"type":       "stype",
					"categories": []string{"catsect1"},
				},
			}),
			"sect1/s1_2/_index.md", p(map[string]interface{}{
				"title": "Sect1_2",
			}),
			"sect1/s1_2/p1.md", p(map[string]interface{}{
				"title": "Sect1_2_p1",
			}),
			"sect1/s1_2/p2.md", p(map[string]interface{}{
				"title": "Sect1_2_p2",
			}),
			"sect2/_index.md", p(map[string]interface{}{
				"title": "Sect2",
			}),
			"sect2/p1.md", p(map[string]interface{}{
				"title":      "Sect2_p1",
				"categories": []string{"cool", "funny", "sad"},
				"tags":       []string{"blue", "green"},
			}),
			"sect2/p2.md", p(map[string]interface{}{}),
			"sect3/p1.md", p(map[string]interface{}{}),

			// No front matter, see #6855
			"sect3/nofrontmatter.md", `**Hello**`,
			"sectnocontent/p1.md", `**Hello**`,
			"sectnofrontmatter/_index.md", `**Hello**`,

			"sect4/_index.md", p(map[string]interface{}{
				"title": "Sect4",
				"cascade": map[string]interface{}{
					"weight":  52,
					"outputs": []string{"RSS"},
				},
			}),
			"sect4/p1.md", p(map[string]interface{}{}),
			"p2.md", p(map[string]interface{}{}),
			"bundle1/index.md", p(map[string]interface{}{}),
			"bundle1/bp1.md", p(map[string]interface{}{}),
			"categories/_index.md", p(map[string]interface{}{
				"title": "My Categories",
				"cascade": map[string]interface{}{
					"title":  "Cascade Category",
					"icoN":   "cat.png",
					"weight": 12,
				},
			}),
			"categories/cool/_index.md", p(map[string]interface{}{}),
			"categories/sad/_index.md", p(map[string]interface{}{
				"cascade": map[string]interface{}{
					"icon":   "sad.png",
					"weight": 32,
				},
			}),
		)
	}

	createContentFiles("en")

	b.WithTemplates("index.html", `
	
{{ range .Site.Pages }}
{{- .Weight }}|{{ .Kind }}|{{ path.Join .Path }}|{{ .Title }}|{{ .Params.icon }}|{{ .Type }}|{{ range .OutputFormats }}{{ .Name }}-{{ end }}|
{{ end }}
`,

		"_default/single.html", "default single: {{ .Title }}|{{ .RelPermalink }}|{{ .Content }}|Resources: {{ range .Resources }}{{ .Name }}|{{ .Params.icon }}|{{ .Content }}{{ end }}",
		"_default/list.html", "default list: {{ .Title }}",
		"stype/single.html", "stype single: {{ .Title }}|{{ .RelPermalink }}|{{ .Content }}",
		"stype/list.html", "stype list: {{ .Title }}",
	)

	return b
}
