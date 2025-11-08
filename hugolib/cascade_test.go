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
	"testing"

	"github.com/gohugoio/hugo/hugolib/sitesmatrix"

	qt "github.com/frankban/quicktest"
)

func BenchmarkCascadeTarget(b *testing.B) {
	files := `
-- content/_index.md --
background = 'yosemite.jpg'
[cascade._target]
kind = '{section,term}'
-- content/posts/_index.md --
-- content/posts/funny/_index.md --
`

	for i := 1; i < 100; i++ {
		files += fmt.Sprintf("\n-- content/posts/p%d.md --\n", i+1)
	}

	for i := 1; i < 100; i++ {
		files += fmt.Sprintf("\n-- content/posts/funny/pf%d.md --\n", i+1)
	}

	b.Run("Kind", func(b *testing.B) {
		cfg := IntegrationTestConfig{
			T:           b,
			TxtarString: files,
		}
		b.ResetTimer()

		for b.Loop() {
			b.StopTimer()
			builder := NewIntegrationTestBuilder(cfg)
			b.StartTimer()
			builder.Build()
		}
	})
}

func TestCascadeBuildOptionsTaxonomies(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL="https://example.org"
[taxonomies]
tag = "tags"

[[cascade]]

[cascade._build]
render = "never"
list = "never"
publishResources = false

[cascade._target]
path = '/hidden/**'
-- content/p1.md --
---
title: P1
---
-- content/hidden/p2.md --
---
title: P2
tags: [t1, t2]
---
-- layouts/_default/list.html --
List: {{ len .Pages }}|
-- layouts/_default/single.html --
Single: Tags: {{ site.Taxonomies.tags }}|
`

	b := Test(t, files)

	b.AssertFileContent("public/p1/index.html", "Single: Tags: map[]|")
	b.AssertFileContent("public/tags/index.html", "List: 0|")
	b.AssertFileExists("public/hidden/p2/index.html", false)
	b.AssertFileExists("public/tags/t2/index.html", false)
}

func TestCascadeEditIssue12449(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com"
disableKinds = ['sitemap','rss', 'home', 'taxonomy','term']
disableLiveReload = true
-- layouts/_default/list.html --
Title: {{ .Title }}|{{ .Content }}|cascadeparam: {{ .Params.cascadeparam }}|
-- layouts/_default/single.html --
Title: {{ .Title }}|{{ .Content }}|cascadeparam: {{ .Params.cascadeparam }}|
-- content/mysect/_index.md --
---
title: mysect
cascade:
  description: descriptionvalue
  params:
    cascadeparam: cascadeparamvalue
---
mysect-content|
-- content/mysect/p1/index.md --
---
slug: p1
---
p1-content|
-- content/mysect/subsect/_index.md --
---
slug: subsect
---
subsect-content|
`

	b := TestRunning(t, files)

	// Make the cascade set the title.
	b.EditFileReplaceAll("content/mysect/_index.md", "description: descriptionvalue", "title: cascadetitle").Build()
	b.AssertFileContent("public/mysect/subsect/index.html", "Title: cascadetitle|")

	// Edit cascade title.
	b.EditFileReplaceAll("content/mysect/_index.md", "title: cascadetitle", "title: cascadetitle-edit").Build()
	b.AssertFileContent("public/mysect/subsect/index.html", "Title: cascadetitle-edit|")

	// Revert title change.
	// The step below failed in #12449.
	b.EditFileReplaceAll("content/mysect/_index.md", "title: cascadetitle-edit", "description: descriptionvalue").Build()
	b.AssertFileContent("public/mysect/subsect/index.html", "Title: |")
}

func TestCascadeIssue12172(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['rss','sitemap','taxonomy','term']
[[cascade]]
headless = true
[cascade._target]
path = '/s1**'
-- content/s1/p1.md --
---
title: p1
---
-- layouts/_default/single.html --
{{ .Title }}|
-- layouts/_default/list.html --
{{ .Title }}|
  `
	b := Test(t, files)

	b.AssertFileExists("public/index.html", true)
	b.AssertFileExists("public/s1/index.html", false)
	b.AssertFileExists("public/s1/p1/index.html", false)
}

// Issue 12594.
func TestCascadeOrder(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['rss','sitemap','taxonomy','term', 'home']
-- content/_index.md --
---
title: Home
cascade:
- _target:
    path: "**"
  params:
    background: yosemite.jpg
- _target:
  params:
    background: goldenbridge.jpg
---
-- content/p1.md --
---
title: p1
---
-- layouts/_default/single.html --
Background: {{ .Params.background }}|
-- layouts/_default/list.html --
{{ .Title }}|
  `

	for range 10 {
		b := Test(t, files)
		b.AssertFileContent("public/p1/index.html", "Background: yosemite.jpg")
	}
}

// Issue #12465.
func TestCascadeOverlap(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['home','rss','sitemap','taxonomy','term']
-- layouts/_default/list.html --
{{ .Title }}
-- layouts/_default/single.html --
{{ .Title }}
-- content/s/_index.md --
---
title: s
cascade:
  build:
    render: never
---
-- content/s/p1.md --
---
title: p1
---
-- content/sx/_index.md --
---
title: sx
---
-- content/sx/p2.md --
---
title: p2
---
`

	b := Test(t, files)

	b.AssertFileExists("public/s/index.html", false)
	b.AssertFileExists("public/s/p1/index.html", false)

	b.AssertFileExists("public/sx/index.html", true)    // failing
	b.AssertFileExists("public/sx/p2/index.html", true) // failing
}

func TestCascadeGotmplIssue13743(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['home','rss','section','sitemap','taxonomy','term']
[cascade.params]
foo = 'bar'
[cascade.target]
path = '/p1'
-- content/_content.gotmpl --
{{ .AddPage (dict "title" "p1" "path" "p1") }}
-- layouts/all.html --
{{ .Title }}|{{ .Params.foo }}
`

	b := Test(t, files)

	b.AssertFileContent("public/p1/index.html", "p1|bar") // actual content is "p1|"
}

func TestSitesMatrixCascadeConfig(t *testing.T) {
	files := `
-- hugo.toml --
disableKinds = ["taxonomy", "term", "rss", "sitemap"]
[languages]
[languages.en]
weight = 1
[languages.nn]
weight = 2
[languages.sv]
weight = 3

[versions]
[versions."v1.0.0"]
[versions."v2.0.0"]
[versions."v2.1.0"]

[roles]
[roles.guest]
[roles.member]
[cascade]
[cascade.sites.matrix]
languages = ["en"]
versions = ["v2**"]
roles = ["member"]
[cascade.sites.complements]
languages = ["nn"]
versions = ["v1.0.*"]
roles = ["guest"]
-- content/_index.md --
---
title: "Home"
sites:
  matrix:
    roles: ["guest"]
---
-- layouts/all.html --
All.

`

	b := Test(t, files)

	s0 := b.H.sitesVersionsRolesMap[sitesmatrix.Vector{0, 0, 0}] // en, v2.1.0, guest
	b.Assert(s0.home, qt.IsNotNil)
	b.Assert(s0.home.File(), qt.IsNotNil)
	b.Assert(s0.language.Name(), qt.Equals, "en")
	b.Assert(s0.version.Name(), qt.Equals, "v2.1.0")
	b.Assert(s0.role.Name(), qt.Equals, "guest")
	s0Pconfig := s0.Home().(*pageState).m.pageConfigSource
	b.Assert(s0Pconfig.SitesMatrix.Vectors(), qt.DeepEquals, []sitesmatrix.Vector{{0, 0, 0}, {0, 1, 0}}) // en, v2.1.0, guest + en, v2.0.0, guest

	s1 := b.H.sitesVersionsRolesMap[sitesmatrix.Vector{1, 2, 0}]
	b.Assert(s1.home, qt.IsNotNil)
	b.Assert(s1.home.File(), qt.IsNil)
	b.Assert(s1.language.Name(), qt.Equals, "nn")
	b.Assert(s1.version.Name(), qt.Equals, "v1.0.0")
	b.Assert(s1.role.Name(), qt.Equals, "guest")
	s1Pconfig := s1.Home().(*pageState).m.pageConfigSource
	b.Assert(s1Pconfig.SitesMatrix.HasVector(sitesmatrix.Vector{1, 2, 0}), qt.IsTrue) // nn, v1.0.0, guest
	// Every site needs a home page. This matrix adds the missing ones, (3 * 3 * 2) - 2 = 16
	b.Assert(s1Pconfig.SitesMatrix.LenVectors(), qt.Equals, 16)
}

func TestCascadeBundledPage(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.org"
-- content/_index.md --
---
title: Home
cascade:
  params:
    p1: v1
---
-- content/b1/index.md --
---
title: b1
---
-- content/b1/p2.md --
---
title: p2
---
-- layouts/all.html --
Title: {{ .Title }}|p1: {{ .Params.p1 }}|
{{ range .Resources }}
Resource: {{ .Name }}|p1: {{ .Params.p1 }}|
{{ end }}
`

	b := Test(t, files)

	b.AssertFileContent("public/b1/index.html", "Title: b1|p1: v1|", "Resource: p2.md|p1: v1|")
}
