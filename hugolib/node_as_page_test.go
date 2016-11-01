// Copyright 2016 The Hugo Authors. All rights reserved.
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
	"path/filepath"
	"testing"

	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

/*
	This file will test the "making everything a page" transition.

	See https://github.com/spf13/hugo/issues/2297

*/

func TestNodesAsPage(t *testing.T) {
	//jww.SetStdoutThreshold(jww.LevelDebug)
	jww.SetStdoutThreshold(jww.LevelFatal)

	nodePageFeatureFlag = true
	defer toggleNodePageFeatureFlag()

	/* Will have to decide what to name the node content files, but:

		Home page should have:
		Content, shortcode support
	   	Metadata (title, dates etc.)
		Params
	   	Taxonomies (categories, tags)

	*/

	testCommonResetState()

	writeLayoutsForNodeAsPageTests(t)

	writeSource(t, filepath.Join("content", "_node.md"), `---
title: Home Sweet Home!
---
Home **Content!**
`)

	writeSource(t, filepath.Join("content", "sect1", "_node.md"), `---
title: Section1
---
Section1 **Content!**
`)

	writeSource(t, filepath.Join("content", "sect2", "_node.md"), `---
title: Section2
---
Section2 **Content!**
`)

	writeSource(t, filepath.Join("content", "categories", "hugo", "_node.md"), `---
title: Taxonomy Hugo
---
Taxonomy Hugo **Content!**
`)

	writeSource(t, filepath.Join("content", "categories", "web", "_node.md"), `---
title: Taxonomy Web
---
Taxonomy Web **Content!**
`)

	writeSource(t, filepath.Join("content", "categories", "_node.md"), `---
title: Taxonomy Term Categories
---
Taxonomy Term Categories **Content!**
`)

	// Add some regular pages
	for i := 1; i <= 4; i++ {
		sect := "sect1"
		if i > 2 {
			sect = "sect2"
		}
		writeSource(t, filepath.Join("content", sect, fmt.Sprintf("regular%d.md", i)), fmt.Sprintf(`---
title: Page %02d
categories:  [
        "Hugo",
		"Web"
]
---
Content Page %02d
`, i, i))
	}

	viper.Set("paginate", 1)

	s := newSiteDefaultLang()

	if err := buildAndRenderSite(s); err != nil {
		t.Fatalf("Failed to build site: %s", err)
	}

	assertFileContent(t, filepath.Join("public", "index.html"), false,
		"Index Title: Home Sweet Home!",
		"Home <strong>Content!</strong>",
		"# Pages: 4")

	assertFileContent(t, filepath.Join("public", "sect1", "regular1", "index.html"), false, "Single Title: Page 01", "Content Page 01")

	h := s.owner
	nodes := h.findAllPagesByNodeType(NodeHome)
	require.Len(t, nodes, 1)

	home := nodes[0]

	require.True(t, home.IsHome())
	require.True(t, home.IsNode())
	require.False(t, home.IsPage())

	pages := h.findAllPagesByNodeType(NodePage)
	require.Len(t, pages, 4)

	first := pages[0]
	require.False(t, first.IsHome())
	require.False(t, first.IsNode())
	require.True(t, first.IsPage())

	first.Paginator()

	// Check Home paginator
	assertFileContent(t, filepath.Join("public", "page", "2", "index.html"), false,
		"Pag: Page 02")

	// Check Sections
	assertFileContent(t, filepath.Join("public", "sect1", "index.html"), false, "Section Title: Section", "Section1 <strong>Content!</strong>")
	assertFileContent(t, filepath.Join("public", "sect2", "index.html"), false, "Section Title: Section", "Section2 <strong>Content!</strong>")

	// Check Sections paginator
	assertFileContent(t, filepath.Join("public", "sect1", "page", "2", "index.html"), false,
		"Pag: Page 02")

	sections := h.findAllPagesByNodeType(NodeSection)
	require.Len(t, sections, 2)

	// Check taxonomy lists
	assertFileContent(t, filepath.Join("public", "categories", "hugo", "index.html"), false,
		"Taxonomy Title: Taxonomy Hugo", "Taxonomy Hugo <strong>Content!</strong>")

	assertFileContent(t, filepath.Join("public", "categories", "web", "index.html"), false,
		"Taxonomy Title: Taxonomy Web", "Taxonomy Web <strong>Content!</strong>")

	// Check taxonomy list paginator
	assertFileContent(t, filepath.Join("public", "categories", "hugo", "page", "2", "index.html"), false,
		"Taxonomy Title: Taxonomy Hugo",
		"Pag: Page 02")

	// Check taxonomy terms
	assertFileContent(t, filepath.Join("public", "categories", "index.html"), false,
		"Taxonomy Terms Title: Taxonomy Term Categories", "Taxonomy Term Categories <strong>Content!</strong>", "k/v: hugo")

	// There are no pages to paginate over in the taxonomy terms.

}

func TestNodesWithNoContentFile(t *testing.T) {
	//jww.SetStdoutThreshold(jww.LevelDebug)
	jww.SetStdoutThreshold(jww.LevelFatal)

	nodePageFeatureFlag = true
	defer toggleNodePageFeatureFlag()

	testCommonResetState()

	writeLayoutsForNodeAsPageTests(t)

	for i := 1; i <= 4; i++ {
		sect := "sect1"
		if i > 2 {
			sect = "sect2"
		}
		writeSource(t, filepath.Join("content", sect, fmt.Sprintf("regular%d.md", i)), fmt.Sprintf(`---
title: Page %02d
categories:  [
        "Hugo",
		"Web"
]
---
Content Page %02d
`, i, i))
	}

	viper.Set("paginate", 1)
	viper.Set("title", "Hugo Rocks!")

	s := newSiteDefaultLang()

	if err := buildAndRenderSite(s); err != nil {
		t.Fatalf("Failed to build site: %s", err)
	}

	// Home page
	homePages := s.findPagesByNodeType(NodeHome)
	require.Len(t, homePages, 1)

	homePage := homePages[0]
	require.Len(t, homePage.Data["Pages"], 4)

	assertFileContent(t, filepath.Join("public", "index.html"), false,
		"Index Title: Hugo Rocks!")

	// Taxonomy list
	assertFileContent(t, filepath.Join("public", "categories", "hugo", "index.html"), false,
		"Taxonomy Title: Hugo")

	// Taxonomy terms
	assertFileContent(t, filepath.Join("public", "categories", "index.html"), false,
		"Taxonomy Terms Title: Categories")

	// Sections
	assertFileContent(t, filepath.Join("public", "sect1", "index.html"), false,
		"Section Title: Sect1s")
	assertFileContent(t, filepath.Join("public", "sect2", "index.html"), false,
		"Section Title: Sect2s")

}

func writeLayoutsForNodeAsPageTests(t *testing.T) {
	writeSource(t, filepath.Join("layouts", "index.html"), `
Index Title: {{ .Title }}
Index Content: {{ .Content }}
# Pages: {{ len .Data.Pages }}
{{ range .Paginator.Pages }}
	Pag: {{ .Title }}
{{ end }}
`)

	writeSource(t, filepath.Join("layouts", "_default", "single.html"), `
Single Title: {{ .Title }}
Single Content: {{ .Content }}
`)

	writeSource(t, filepath.Join("layouts", "_default", "section.html"), `
Section Title: {{ .Title }}
Section Content: {{ .Content }}
# Pages: {{ len .Data.Pages }}
{{ range .Paginator.Pages }}
	Pag: {{ .Title }}
{{ end }}
`)

	// Taxonomy lists
	writeSource(t, filepath.Join("layouts", "_default", "taxonomy.html"), `
Taxonomy Title: {{ .Title }}
Taxonomy Content: {{ .Content }}
# Pages: {{ len .Data.Pages }}
{{ range .Paginator.Pages }}
	Pag: {{ .Title }}
{{ end }}
`)

	// Taxonomy terms
	writeSource(t, filepath.Join("layouts", "_default", "terms.html"), `
Taxonomy Terms Title: {{ .Title }}
Taxonomy Terms Content: {{ .Content }}
{{ range $key, $value := .Data.Terms }}
	k/v: {{ $key }} / {{ printf "%=v" $value }}
{{ end }}
`)
}
