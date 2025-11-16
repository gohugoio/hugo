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

package sitesmatrix_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/hugolib"
)

func TestPageRotate(t *testing.T) {
	t.Parallel()
	files := `
-- hugo.toml --
baseURL = "https://example.org/"
defaultContentVersion = "v4.0.0"
defaultContentVersionInSubdir = true
defaultContentRoleInSubdir = true
defaultContentRole = "guest"
defaultContentLanguage = "en"
defaultContentLanguageInSubdir = true
disableKinds = ["taxonomy", "term", "rss", "sitemap"]

[cascade]
versions = ["v2**"]

[languages]
[languages.en]
weight = 1
[languages.nn]
weight = 2
[roles]
[roles.guest]
weight = 300
[roles.member]
weight = 200
[versions]
[versions."v2.0.0"]
[versions."v1.2.3"]
[versions."v2.1.0"]
[versions."v3.0.0"]
[versions."v4.0.0"]
-- content/_index.en.md --
---
title: "Home"
roles: ["**"]
versions: ["**"]
---
-- content/_index.nn.md --
---
title: "Heim"
roles: ["**"]
versions: ["**"]
---
-- content/memberonlypost.md --
---
title: "Member Only"
roles: ["member"]
languages: ["**"]
---
Member content.
-- content/publicpost.md --
---
title: "Public"
versions: ["v1.2.3", "v2.**", "! v2.1.*"]
versionDelegees: ["v3**"]
---
Users with guest role will see this.
-- content/v3publicpost.md --
---
title: "Public v3"
versions: ["v3**"]
languages: ["**"]
---
Users with guest role will see this.
-- layouts/all.html --
Rotate(language): {{ with .Rotate "language" }}{{ range . }}{{ template "printp" . }}|{{ end }}{{ end }}$
Rotate(version): {{ with .Rotate "version" }}{{ range . }}{{ template "printp" . }}|{{ end }}{{ end }}$
Rotate(role): {{ with .Rotate "role" }}{{ range . }}{{ template "printp" . }}|{{ end }}{{ end }}$
.Site.Dimension.language {{ (.Site.Dimension "language").Name }}|
.Site.Dimension.version: {{ (.Site.Dimension "version").Name }}|
.Site.Dimension.role: {{ (.Site.Dimension "role").Name }}|
{{ define "printp" }}{{ .RelPermalink }}:{{ with .Site }}{{ template "prints" . }}{{ end }}{{ end }}
{{ define "prints" }}/l:{{ .Language.Name }}/v:{{ .Version.Name }}/r:{{ .Role.Name }}{{ end }}


`

	for range 3 {
		b := hugolib.Test(t, files)

		b.AssertFileContent("public/guest/v3.0.0/en/index.html",
			"Rotate(language): /guest/v3.0.0/en/:/l:en/v:v3.0.0/r:guest|/guest/v3.0.0/nn/:/l:nn/v:v3.0.0/r:guest|$",
			"Rotate(version): /guest/v4.0.0/en/:/l:en/v:v4.0.0/r:guest|/guest/v3.0.0/en/:/l:en/v:v3.0.0/r:guest|/guest/v2.1.0/en/:/l:en/v:v2.1.0/r:guest|/guest/v2.0.0/en/:/l:en/v:v2.0.0/r:guest|/guest/v1.2.3/en/:/l:en/v:v1.2.3/r:guest",
			"Rotate(role): /member/v3.0.0/en/:/l:en/v:v3.0.0/r:member|/guest/v3.0.0/en/:/l:en/v:v3.0.0/r:guest|$",
			".Site.Dimension.language en|",
			".Site.Dimension.version: v3.0.0|",
			".Site.Dimension.role: guest|",
		)

	}
}

func TestFileMountSitesMatrix(t *testing.T) {
	filesTemplate := `
-- hugo.toml --
disableKinds = ["taxonomy", "term", "rss", "sitemap", "section"]
defaultContentLanguage = "en"
defaultContentLanguageInSubDir = true
defaultCOntentVersionInSubDir = true
[versions]
[versions."v1.2.3"]
[versions."v2.0.0"]
[languages]
[languages.en]
weight = 1
[languages.nn]
weight = 2
[module]
[[module.mounts]]
source = 'content/en'
target = 'content'
DIMSEN
[[module.mounts]]
source = 'content/nn'
target = 'content'
DIMSNN
[[module.mounts]]
source = 'content/all'
target = 'content'
[module.mounts.sites.matrix]
languages = ["**"]
versions = ["**"]
-- content/en/p1/index.md --
---
title: "Title English"
---
-- content/en/p1/mytext.txt --
Text English
-- content/nn/p1/index.md --
---
title: "Tittel Nynorsk"
---
-- content/all/p2/index.md --
---
title: "p2 all"
---
-- content/nn/p1/mytext.txt --
Tekst Nynorsk
-- layouts/all.html --
{{ $mytext := .Resources.Get "mytext.txt" }}
{{ .Title }}|{{ with  $mytext }}{{ .Content | safeHTML }}{{ end }}|
site.GetPage p2: {{ with .Site.GetPage "p2" }}{{ .Title }}|{{ end }}$
site.GetPage p1: {{ with .Site.GetPage "p1" }}{{ .Title }}|{{ end }}$

`

	testOne := func(t *testing.T, name, files string) {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			b := hugolib.Test(t, files)

			nn20 := b.SiteHelper("nn", "v2.0.0", "")
			b.Assert(nn20.PageHelper("/p2").MatrixFromFile()["matrix"],
				qt.DeepEquals,
				map[string][]string{"languages": {"en", "nn"}, "roles": {"guest"}, "versions": {"v2.0.0", "v1.2.3"}})

			b.AssertFileContent("public/v2.0.0/nn/p1/index.html", "Tittel Nynorsk", "Tekst Nynorsk", "site.GetPage p1: Tittel Nynorsk|")
			b.AssertFileContent("public/v2.0.0/nn/p2/index.html", "p2 all||", "site.GetPage p2: p2 all", "site.GetPage p1: Tittel Nynorsk|")
			b.AssertFileContent("public/v2.0.0/nn/p2/index.html", "p2 all||", "site.GetPage p1: Tittel Nynorsk|$")
			b.AssertFileContent("public/v1.2.3/en/p2/index.html", "p2 all||", "site.GetPage p2: p2 all")
		})
	}

	// Format from v0.148.0:
	dims := `[module.mounts.sites.matrix]
languages = ["en"]
versions = ["v1**"]
`
	files := strings.Replace(filesTemplate, "DIMSEN", dims, 1)
	dims = strings.Replace(dims, `["en"]`, `["nn"]`, 1)
	dims = strings.Replace(dims, `["v1**"]`, `["v2**"]`, 1)
	files = strings.Replace(files, "DIMSNN", dims, 1)
	testOne(t, "new", files)

	// Old format:
	files = strings.Replace(filesTemplate, "DIMSEN", `lang = "en"`, 1)
	files = strings.Replace(files, "DIMSNN", `lang = "nn"`, 1)
	testOne(t, "old", files)
}

func TestSpecificMountShouldAlwaysWin(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ["rss", "sitemap", "section"]
defaultContentLanguage = "en"
defaultContentLanguageInSubDir = true
defaultCOntentVersionInSubDir = true
defaultContentVersion = "v2.0.0"
[taxonomies]
tag = "tags"
[languages]
[languages.en]
weight = 1
[languages.nn]
weight = 2
[versions]
[versions."v1.2.3"]
[versions."v2.0.0"]
[module]
[[module.mounts]]
source = 'content/nn'
target = 'content'
[module.mounts.sites.matrix]
languages = ["nn"]
versions  = ["v1.**"]
[[module.mounts]]
source = 'content/en'
target = 'content'
[module.mounts.sites.matrix]
languages = ["**"]
versions  = ["**"]
-- content/en/_index.md --
---
title: "English Home"
tags: ["tag1"]
---
-- content/en/p1.md --
---
title: "English p1"
---
-- content/nn/_index.md --
---
title: "Nynorsk Heim"
tags: ["tag2"]
---
-- layouts/all.html --
title: {{ .Title }}|
tags: {{ range $term, $taxonomy := .Site.Taxonomies.tags }}{{ $term }}: {{ range $taxonomy.Pages }}{{ .Title }}: {{ .RelPermalink}}|{{ end }}{{ end }}$
`

	for range 2 {
		b := hugolib.Test(t, files)

		// b.AssertPublishDir("asdf")
		b.AssertFileContent("public/v1.2.3/nn/index.html", "title: Nynorsk Heim|", "tags: tag2: Nynorsk Heim: /v1.2.3/nn/|$")
		b.AssertFileContent("public/v2.0.0/en/index.html", "title: English Home|", "tags: tag1: English Home: /v2.0.0/en/|$")
		b.AssertFileContent("public/v2.0.0/nn/index.html", "title: English Home|", "tags: tag1: English Home: /v2.0.0/nn/|$") // v2.0.0 is only in English.
		b.AssertFileContent("public/v1.2.3/en/index.html", "title: English Home|")
	}
}

const filesVariationsSitesMatrixBase = `
-- hugo.toml --
baseURL = "https://example.org/"
disableKinds = ["rss", "sitemap", "section"]
defaultContentLanguage = "en"
defaultContentLanguageInSubDir = true
defaultCOntentVersionInSubDir = true
defaultContentVersion = "v2.0.0"
defaultContentRole = "guest"
defaultContentRoleInSubDir = true
[taxonomies]
tag = "tags"
[languages]
[languages.en]
weight = 1
[languages.nn]
weight = 2
[versions]
[versions."v1.2.3"]
[versions."v1.4.0"]
[versions."v2.0.0"]
[roles]
[roles.guest]
weight = 300
[roles.member]
weight = 200
[module]
[[module.mounts]]
source = 'content/nn'
target = 'content'
[module.mounts.sites.matrix]
languages = ["nn"]
versions  = ["v1.2.*"]
[[module.mounts]]
source = 'content/en'
target = 'content'
[module.mounts.sites.matrix]
languages = ["**"]
versions  = ["**"]
[[module.mounts]]
source = 'content/other'
target = 'content'
-- content/en/_index.md --
---
title: "English Home"
tags: ["tag1"]
---

Ref home: {{< ref "/" >}}|
-- content/en/p1.md --
---
title: "English p1"
---
-- content/nn/_index.md --
---
title: "Nynorsk Heim"
tags: ["tag2"]
---

Ref home: {{< ref "/" >}}|
`

const filesVariationsSitesMatrix = filesVariationsSitesMatrixBase + `
-- layouts/all.html --
title: {{ .Title }}|
tags: {{ range $term, $taxonomy := .Site.Taxonomies.tags }}{{ $term }}: {{ range $taxonomy.Pages }}{{ .Title }}: {{ .RelPermalink}}|{{ end }}{{ end }}$
.Language.IsDefault: {{ with .Rotate "language" }}{{ range . }}{{ .RelPermalink }}: {{ with .Site.Language }}{{ .Name}}: {{ .IsDefault }}|{{ end }}{{ end }}{{ end }}$
.Version.IsDefault: {{ with .Rotate "version" }}{{ range . }}{{ .RelPermalink }}: {{ with .Site.Version }}{{ .Name}}: {{ .IsDefault }}|{{ end }}{{ end }}{{ end }}$
.Role.IsDefault: {{ with .Rotate "role" }}{{ range . }}{{ .RelPermalink }}: {{ with .Site.Role }}{{ .Name}}: {{ .IsDefault }}|{{ end }}{{ end }}{{ end }}$
`

func TestFrontMatterSitesMatrix(t *testing.T) {
	t.Parallel()

	files := filesVariationsSitesMatrix

	files += `
-- content/other/p2.md --
---
title: "NN p2"
sites:
  matrix:
    languages:
      - nn
    versions:
      - v1.2.3
---
`
	b := hugolib.Test(t, files)
	b.AssertFileContent("public/guest/v1.2.3/nn/p2/index.html", "title: NN p2|")
}

func TestFrontMatterSitesMatrixShouldWin(t *testing.T) {
	t.Parallel()

	files := filesVariationsSitesMatrix

	// nn mount config is nn, v1.2.*.
	files += `
-- content/nn/p2.md --
---
title: "EN p2"
sites:
  matrix:
    languages:
      - en
    versions:
      - v1.4.*
---
`
	b := hugolib.Test(t, files)
	b.AssertFileContent("public/guest/v1.4.0/en/p2/index.html", "title: EN p2|")
}

func TestGetPageAndRef(t *testing.T) {
	t.Parallel()

	files := filesVariationsSitesMatrixBase + `
-- layouts/all.html --
Title: {{ .Title }}|
Home: {{ with .Site.GetPage "/" }}{{ with .Site }}Language: {{ .Language.Name }}|Version: {{ .Version.Name }}|Role: {{ .Role.Name }}{{ end }}|{{ end }}$
Content: {{ .Content }}$
`

	b := hugolib.Test(t, files)
	b.AssertFileContent(
		"public/guest/v2.0.0/en/index.html", "Home: Language: en|Version: v2.0.0|Role: guest|$",
		"Ref home: https://example.org/guest/v2.0.0/en/|",
	)
	b.AssertFileContent(
		"public/member/v1.4.0/nn/index.html", "Home: Language: nn|Version: v1.4.0|Role: member|$",
	)

	b.AssertFileContent(
		"public/guest/v1.2.3/nn/index.html", "Home: Language: nn|Version: v1.2.3|Role: guest|$",
		"Ref home: https://example.org/guest/v1.2.3/nn/|",
	)
}

func TestFrontMatterSitesMatrixShouldBeMergedWithMount(t *testing.T) {
	t.Parallel()

	files := filesVariationsSitesMatrix

	// nn mount config is nn, v1.2.*.
	// This changes only the language, not the version.
	files += `
-- content/nn/p2.md --
---
title: "EN p2"
sites:
  matrix:
    languages:
      - en
---
`
	b := hugolib.Test(t, files)
	b.AssertFileContent("public/guest/v1.2.3/en/p2/index.html", "title: EN p2|")
}

func TestCascadeMatrixInFrontMatter(t *testing.T) {
	t.Parallel()
	files := `
-- hugo.toml --
disableKinds = ["rss", "sitemap", "taxonomy", "term"]
defaultContentLanguage = "en"
defaultContentLanguageInSubDir = true
[languages]
[languages.en]
weight = 1
[languages.nn]
weight = 2
[languages.sv]
weight = 3
-- content/_index.md --
+++
title = "English home"
[cascade]
[cascade.params]
p1 = "p1cascade"
[cascade.target.sites.matrix]
languages = ["en"]
+++
-- content/_index.nn.md --
+++
title = "Scandinavian home"
[sites.matrix]
languages = "{nn,sv}"
[cascade]
[cascade.params]
p1 = "p1cascadescandinavian"
[cascade.sites.matrix]
languages = "{nn,sv}"
[cascade.target]
path = "**scandinavian**"
+++
-- content/mysection/_index.md --
+++
title = "My section"
[cascade.target.sites.matrix]
languages = "**"
+++
-- content/mysection/p1.md --
+++
title = "English p1"
+++
-- content/mysection/scandinavianpages/p1.md --
+++
title = "Scandinavian p1"
+++
-- layouts/all.html --
{{ .Title }}|{{ .Site.Language.Name }}|{{ .Site.Version.Name }}|p1: {{ .Params.p1 }}|
`
	b := hugolib.Test(t, files)
	b.AssertFileContent("public/en/mysection/p1/index.html", "English p1|en|v1.0.0|p1: p1cascade|")
	b.AssertFileContent("public/sv/mysection/scandinavianpages/p1/index.html", "Scandinavian p1|sv|v1.0.0|p1: p1cascadescandinavian|")
}

func TestCascadeMatrixConfigPerLanguage(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ["rss", "sitemap", "taxonomy", "term"]
defaultContentLanguage = "en"
defaultContentLanguageInSubDir = true
[cascade.params]
p1 = "p1cascadeall"
p2 = "p2cascadeall"
[languages]
[languages.en]
weight = 1
[languages.en.cascade.params]
p1 = "p1cascadeen"
[languages.nn]
weight = 2
[languages.nn.cascade.params]
p1 = "p1cascadenn"
-- content/_index.md --
---
title: "English home"
---
-- content/mysection/p1.md --
---
title: "English p1"
---
-- content/mysection/p1.nn.md --
---
title: "Nynorsk p1"
---
-- layouts/all.html --
{{ .Title }}|{{ .Site.Language.Name }}|{{ .Site.Version.Name }}|p1: {{ .Params.p1 }}|p2: {{ .Params.p2 }}|

`
	b := hugolib.Test(t, files)
	b.AssertFileContent("public/en/mysection/p1/index.html", "English p1|en|v1.0.0|p1: p1cascadeen|p2: p2cascadeall|")
	b.AssertFileContent("public/nn/mysection/p1/index.html", "Nynorsk p1|nn|v1.0.0|p1: p1cascadenn|p2: p2cascadeall|")
}

func TestCascadeMatrixNoHomeContent(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ["rss", "sitemap", "taxonomy", "term"]
defaultContentLanguage = "en"
defaultContentLanguageInSubDir = true
[cascade.params]
p1 = "p1cascadeall"
p2 = "p2cascadeall"
-- content/mysection/p1.md --
-- layouts/all.html --
{{ .Title }}|{{ .Kind }}|{{ .Site.Language.Name }}|{{ .Site.Version.Name }}|p1: {{ .Params.p1 }}|p2: {{ .Params.p2 }}|
`

	b := hugolib.Test(t, files)
	b.AssertFileContent("public/en/index.html", "|home|en|v1.0.0|p1: p1cascadeall|p2: p2cascadeall|")
	b.AssertFileContent("public/en/mysection/index.html", "Mysections|section|en|v1.0.0|p1: p1cascadeall|p2: p2cascadeall|")
	b.AssertFileContent("public/en/mysection/p1/index.html", "|page|en|v1.0.0|p1: p1cascadeall|p2: p2cascadeall|")
}

func TestMountCascadeFrontMatterSitesMatrixAndComplementsShouldBeMerged(t *testing.T) {
	t.Parallel()

	// Pick language from mount, role from cascade and version from front matter.
	files := `
-- hugo.toml --
disableKinds = ["rss", "sitemap", "section", "taxonomy", "term"]
defaultContentLanguage = "en"
defaultContentLanguageInSubDir = true
defaultCOntentVersionInSubDir = true
defaultContentVersion = "v1.2.3"
defaultContentRole = "guest"
defaultContentRoleInSubDir = true

[languages]
[languages.en]
weight = 1
[languages.nn]
weight = 2

[versions]
[versions."v1.2.3"]
[versions."v1.4.0"]
[versions."v2.0.0"]

[roles]
[roles.guest]
weight = 300
[roles.member]
weight = 200

[[module.mounts]]
source = 'content/other'
target = 'content'
[module.mounts.sites.matrix]
languages = ["nn"] # used.
versions = ["v1.2.*"] # not used.
roles = ["guest"] # not used.
[module.mounts.sites.complements]
languages = ["en"] # used.
versions = ["v1.4.*"] # not used.
roles = ["member"] # not used.

[cascade.sites.matrix]
roles = ["member"] # used
versions = ["v1.2.*"] # not used.
[cascade.sites.complements]
roles = ["guest"] # used
versions = ["v2**"] # not used.

-- content/other/p2.md --
+++
title = "NN p2"
[sites.matrix]
versions = ["v1.2.*","v1.4.*"]
[sites.complements]
versions = ["v2.*.*"]
+++
-- content/other/p3.md --
+++
title = "NN p3"
[sites.matrix]
versions = "v1.4.*"
+++
-- layouts/all.html --
All. {{ site.Language.Name }}|{{ site.Version.Name }}|{{ site.Role.Name }}

`

	b := hugolib.Test(t, files)
	b.AssertFileContent("public/member/v1.4.0/nn/p2/index.html", "All.")
	s := b.SiteHelper("nn", "v1.4.0", "member")
	p2 := s.PageHelper("/p2")
	s.Assert(p2.MatrixFromPageConfig(), qt.DeepEquals,
		map[string]map[string][]string{
			"complements": {
				"languages": {"en"},
				"roles":     {"guest"},
				"versions":  {"v2.0.0"},
			},
			"matrix": {
				"languages": {"nn"},
				"roles":     {"member"},
				"versions":  {"v1.4.0", "v1.2.3"},
			},
		},
	)

	p3 := s.PageHelper("/p3")
	s.Assert(p3.MatrixFromPageConfig(), qt.DeepEquals, map[string]map[string][]string{
		"complements": {
			"languages": {"en"},
			"roles":     {"guest"},
			"versions":  {"v2.0.0"},
		},
		"matrix": {
			"languages": {"nn"},
			"roles":     {"member"},
			"versions":  {"v1.4.0"},
		},
	})
}

func TestLanguageVersionRoleIsDefault(t *testing.T) {
	files := filesVariationsSitesMatrix

	b := hugolib.Test(t, files)
	b.AssertFileContent("public/guest/v2.0.0/en/index.html",
		".Language.IsDefault: /guest/v2.0.0/en/: en: true|/guest/v2.0.0/nn/: nn: false|$",
		".Version.IsDefault: /guest/v2.0.0/en/: v2.0.0: true|/guest/v1.4.0/en/: v1.4.0: false|/guest/v1.2.3/en/: v1.2.3: false|$",
		".Role.IsDefault: /member/v2.0.0/en/: member: false|/guest/v2.0.0/en/: guest: true|$",
	)
}

func TestModuleMountsLanguageOverlap(t *testing.T) {
	files := `
-- hugo.toml --
disableKinds = ["rss", "sitemap", "section", "taxonomy", "term"]
theme = "mytheme"
defaultContentLanguage = "en"
defaultContentLanguageInSubDir = true
[languages]
[languages.en]
weight = 1
[languages.nn]
weight = 2
[[module.mounts]]
source = 'content'
target = 'content'
[module.mounts.sites.matrix]
languages = "en"
-- content/p1.md --
---
title: "Project p1"
---
-- themes/mytheme/hugo.toml --

[[module.mounts]]
source = 'content'
target = 'content'
[module.mounts.sites.matrix]
languages = "**"
-- themes/mytheme/content/p1.md --
---
title: "Theme p1"
---
-- themes/mytheme/content/p2.md --
---
title: "Theme p1"
---
-- layouts/all.html --
{{ .Title }}|
`

	b := hugolib.Test(t, files)

	sEn := b.SiteHelper("en", "", "")
	sNN := b.SiteHelper("nn", "", "guest")

	b.Assert(sEn.PageHelper("/p1").MatrixFromPageConfig()["matrix"], qt.DeepEquals, map[string][]string{"languages": {"en"}, "roles": {"guest"}, "versions": {"v1.0.0"}})
	b.Assert(sEn.PageHelper("/p2").MatrixFromPageConfig()["matrix"], qt.DeepEquals, map[string][]string{
		"languages": {"en", "nn"},
		"roles":     {"guest"},
		"versions":  {"v1.0.0"},
	})

	p1NN := sNN.PageHelper("/p1")
	p2NN := sNN.PageHelper("/p2")

	b.Assert(p2NN.MatrixFromPageConfig()["matrix"], qt.DeepEquals, map[string][]string{
		"languages": {"en", "nn"},
		"roles":     {"guest"},
		"versions":  {"v1.0.0"},
	})

	b.Assert(p1NN.MatrixFromFile()["matrix"], qt.DeepEquals, map[string][]string{
		"languages": {"nn"}, // It's defined as ** in the mount.
		"roles":     {"guest"},
		"versions":  {"v1.0.0"},
	})
}

func TestMountLanguageComplements(t *testing.T) {
	files := `
-- hugo.toml --
disableKinds = ["rss", "sitemap", "section", "taxonomy", "term"]
defaultContentLanguage = "en"
defaultContentLanguageInSubDir = true
[languages]
[languages.en]
weight = 1
[languages.sv]
weight = 2
[languages.da]
weight = 3
[[module.mounts]]
source = 'content/en'
target = 'content'
[module.mounts.sites.matrix]
languages = "en"
[module.mounts.sites.complements]
languages = "{sv,da}"
[[module.mounts]]
source = 'content/sv'
target = 'content'
[module.mounts.sites.matrix]
languages = "sv"
-- content/en/p1.md --
---
title: "English p1"
---
-- content/en/p2.md --
---
title: "English p2"
---
-- content/sv/p1.md --
---
title: "Swedish p1"
---
-- layouts/home.html --
RegularPagesRecursive: {{ range .RegularPagesRecursive }}{{ .Title }}|{{ .RelPermalink }}{{ end }}$
-- layouts/all.html --
{{ .Title }}|{{ .Site.Language.Name }}|
`

	b := hugolib.Test(t, files)
	b.AssertFileContent("public/en/index.html", "RegularPagesRecursive: English p1|/en/p1/English p2|/en/p2/$")
	b.AssertFileContent("public/sv/index.html", "RegularPagesRecursive: English p2|/en/p2/Swedish p1|/sv/p1/$")
	b.AssertFileContent("public/da/index.html", "RegularPagesRecursive: English p1|/en/p1/English p2|/en/p2/$")
}

func TestContentAdapterSitesMatrixResources(t *testing.T) {
	t.Parallel()
	files := `
-- hugo.toml --
disableKinds = ["rss", "sitemap", "section", "taxonomy", "term"]
defaultContentLanguage = "en"
defaultContentLanguageInSubDir = true
[languages]
[languages.en]
weight = 1
[languages.nn]
weight = 2
[languages.sv]
weight = 3
[languages.da]
weight = 4
-- content/_content.gotmpl --
{{ $scandinavia := dict "languages" "{nn,sv}" }}
{{ $en := dict "languages" "en" }}
{{ $da := dict "languages" "da" }}
{{ $contentMarkdown := dict "value" "**Hello World**"  "mediaType" "text/markdown" }}
{{ $contentTextEnglish := dict "value" "Hello World"  "mediaType" "text/plain" }}
{{ $contentTextNorsk := dict "value" "Hallo verd"  "mediaType" "text/plain" }}
{{ .AddPage (dict "path" "p1" "title" "P1 en" "content" $contentMarkdown "sites" (dict "matrix"  $en )) }}
{{ .AddPage (dict "path" "p1" "title" "P1 scandinavia" "content" $contentMarkdown "sites" (dict "matrix"  $scandinavia )) }}
{{ .AddPage (dict "path" "p1" "title" "P1 da" "content" $contentMarkdown "sites" (dict "matrix"  $da )) }}
{{ .AddResource (dict "path" "p1/hello.txt" "title" "Hello en" "content" $contentTextEnglish "sites" (dict "matrix"  $en )) }}
{{ .AddResource (dict "path" "p1/hello.txt" "title" "Hello scandinavia" "content" $contentTextNorsk "sites" (dict "matrix"  $scandinavia "complements" $da )) }}
-- layouts/all.html --
len .Resources: {{ len .Resources}}|
{{ $hello := .Resources.Get "hello.txt" }}
All. {{ .Title }}|Hello: {{ with $hello }}{{ .RelPermalink }}|{{ .Content | safeHTML }}{{ end }}|
`

	b := hugolib.Test(t, files)
	b.AssertFileContent("public/en/p1/index.html", "Hello: /en/p1/hello.txt|Hello World|")
	b.AssertFileContent("public/nn/p1/index.html", "P1 scandinavia|", "/nn/p1/hello.txt|Hallo verd|")
	b.AssertFileContent("public/sv/p1/index.html", "P1 scandinavia|", "/sv/p1/hello.txt|Hallo verd|")
	b.AssertFileContent("public/da/p1/index.html", "P1 da|", "/sv/p1/hello.txt|Hallo verd|") // Because it's closest of the Scandinavian complements.
}

func TestContentAdapterSitesMatrixContentPageSwitchLanguage(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ["rss", "sitemap", "section", "taxonomy", "term"]
defaultContentLanguage = "en"
defaultContentLanguageInSubDir = true
disableLiveReload = true
[languages]
[languages.en]
weight = 1
[languages.sv]
weight = 2
[languages.da]
weight = 3
[languages.no]
weight = 4
-- content/_content.gotmpl --
{{ $contentMarkdown := dict "value" "**Hello World**"  "mediaType" "text/markdown" }}
{{ $en := dict "languages" "en" }}
{{ $sv := dict "languages" "sv" }}
{{ $da := dict "languages" "da" }}
{{ $no := dict "languages" "no" }}
{{ .AddPage (dict "path" "p1" "title" "P1-1" "content" $contentMarkdown "sites" (dict "matrix"  $en )) }}
{{ .AddPage (dict "path" "p1" "title" "P1-2" "content" $contentMarkdown "sites" (dict "matrix"  $sv )) }}
{{ .AddPage (dict "path" "p1" "title" "P1-3" "content" $contentMarkdown "sites" (dict "matrix"  $da )) }}
-- layouts/all.html --
All. {{ .Title }}|{{ .Site.Language.Name }}|
`

	b := hugolib.TestRunning(t, files)

	b.AssertFileExists("public/no/p1/index.html", false)
	b.AssertFileExists("public/sv/p1/index.html", true)

	b.RemovePublishDir()
	b.AssertFileExists("public/sv/p1/index.html", false)

	b.EditFileReplaceAll("content/_content.gotmpl", `"sv"`, `"no"`).Build()

	b.AssertFileExists("public/no/p1/index.html", true)
	b.AssertFileExists("public/sv/p1/index.html", false)
}

func TestContentAdapterSitesMatrixContentResourceSwitchLanguage(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ["rss", "sitemap", "section", "taxonomy", "term"]
disableLiveReload = true
defaultContentLanguage = "en"
defaultContentLanguageInSubDir = true
[languages]
[languages.en]
weight = 1
[languages.sv]
weight = 2
[languages.da]
weight = 3
[languages.no]
weight = 4
-- content/_content.gotmpl --
{{ $contentTextEnglish := dict "value" "Hello World"  "mediaType" "text/plain" }}
{{ $contentTextSwedish := dict "value" "Hei världen"  "mediaType" "text/plain" }}
{{ $contentTextDanish := dict "value" "Hej verden"  "mediaType" "text/plain" }}

{{ $en := dict "languages" "en" }}
{{ $sv := dict "languages" "sv" }}
{{ $da := dict "languages" "da" }}
{{ .AddResource (dict "path" "hello.txt" "title" "Hello" "content" $contentTextEnglish "sites" (dict "matrix"  $en )) }}
{{ .AddResource (dict "path" "hello.txt" "title" "Hello" "content" $contentTextSwedish "sites" (dict "matrix"  $sv )) }}
{{ .AddResource (dict "path" "hello.txt" "title" "Hello" "content" $contentTextDanish "sites" (dict "matrix"  $da )) }}
 
-- layouts/all.html --
{{ $hello := .Resources.Get "hello.txt" }}
All. {{ .Title }}|{{ .Site.Language.Name }}|hello: {{ with $hello }}{{ .RelPermalink }}: {{ .Content }}|{{ end }}|
`

	b := hugolib.TestRunning(t, files)

	b.AssertFileContent("public/en/index.html", "en|hello: /en/hello.txt: Hello World|")
	b.AssertFileContent("public/sv/index.html", "sv|hello: /sv/hello.txt: Hei världen|")
	b.AssertFileContent("public/no/index.html", "no|hello: |")

	b.EditFileReplaceAll("content/_content.gotmpl", `"sv"`, `"no"`).Build()

	b.AssertFileContent("public/no/index.html", "no|hello: /no/hello.txt: Hei världen|")
	b.AssertFileContent("public/en/index.html", "en|hello: /en/hello.txt: Hello World|")
	b.AssertFileContent("public/sv/index.html", "|sv|hello: |")
}

const filesContentAdapterSitesMatrixFromConfig = `
-- hugo.toml --
disableKinds = ["rss", "sitemap", "section", "taxonomy", "term"]
defaultContentLanguage = "en"
defaultContentLanguageInSubDir = true
defaultContentRole = "enrole"
defaultContentRoleInSubDir = true
[languages]
[languages.en]
weight = 1
[languages.en.cascade.sites.matrix]
roles = ["enrole"]
[languages.nn]
weight = 2
[languages.nn.cascade.sites.matrix]
roles = ["nnrole"]
[roles]
[roles.enrole]
weight = 200
[roles.nnrole]
weight = 100
-- layouts/all.html --
All. {{ .Title }}|{{ .Site.Language.Name }}|{{ .Site.Role.Name }}|
Resources: {{ range .Resources }}{{ .Title }}|{{ end }}$

`

func TestContentAdapterSitesMatrixSitesMatrixFromConfig(t *testing.T) {
	t.Parallel()

	files := filesContentAdapterSitesMatrixFromConfig + `
-- content/_content.gotmpl --
{{ $title := printf "P1 %s:%s" .Site.Language.Name .Site.Role.Name }}
{{ .AddPage (dict "path" "p1" "title" $title ) }}
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/enrole/en/p1/index.html", "P1 en:enrole|en|enrole|")
	b.AssertPublishDir("! nn/p1")
}

func TestContentAdapterSitesMatrixSitesMatrixFromConfigEnableAllLanguages(t *testing.T) {
	t.Parallel()

	files := filesContentAdapterSitesMatrixFromConfig + `
-- content/_content.gotmpl --
{{ .EnableAllLanguages }}
{{ $title := printf "P1 %s:%s" .Site.Language.Name .Site.Role.Name }}
{{ .AddPage (dict "path" "p1" "title" $title ) }}
{{ $.AddResource (dict "path" "p1/mytext.txt" "title" $title "content" (dict "value" $title) )}}
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/enrole/en/p1/index.html", "P1 en:enrole|en|enrole|", "Resources: P1 en:enrole|$")
	b.AssertFileContent("public/nnrole/nn/p1/index.html", "P1 nn:nnrole|nn|nnrole|", "Resources: P1 nn:nnrole|$")
}

func TestSitesMatrixCustomContentFilenameIdentifier(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ["rss", "sitemap", "section", "taxonomy", "term"]
defaultContentLanguage = "en"
defaultContentLanguageInSubDir = true
[languages]
[languages.en]
weight = 1
[languages.nn]
weight = 2
[languages.sv]
weight = 3
[languages.da]
weight = 4
[languages.de]
weight = 5
-- content/p1.en.md --
---
title: "P1 en"
---
-- content/p1._scandinavian_.md --
---
title: "P1 scandinavian"
sites:
  matrix:
    languages: "{nn,sv,da}"
---
-- content/p1.de.md --
---
title: "P1 de"
---
-- layouts/all.html --
All.{{ .Title }}|
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/en/p1/index.html", "All.P1 en|")
	b.AssertFileContent("public/nn/p1/index.html", "All.P1 scandinavian|")
	b.AssertFileContent("public/sv/p1/index.html", "All.P1 scandinavian|")
	b.AssertFileContent("public/da/p1/index.html", "All.P1 scandinavian|")
	b.AssertFileContent("public/de/p1/index.html", "All.P1 de|")
}

func TestSitesMatrixDefaultValues(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ["rss", "sitemap", "section", "taxonomy", "term"]
-- layouts/all.html --
All. {{ .Title }}|Lang: {{ .Site.Language.Name }}|Ver: {{ .Site.Version.Name }}|Role: {{ .Site.Role.Name }}|
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html", "All. |Lang: en|Ver: v1.0.0|Role: guest|")
}

func newSitesMatrixContentBenchmarkBuilder(t testing.TB, numPages int, skipRender, multipleDimensions bool) *hugolib.IntegrationTestBuilder {
	files := `
-- hugo.toml --
disableKinds = ["rss", "sitemap", "section", "taxonomy", "term"]
defaultContentLanguage = "en"
defaultContentLanguageInSubDir = true
defaultCOntentVersionInSubDir = true
defaultContentVersion = "v1.0.0"
defaultContentRole = "guest"
defaultContentRoleInSubDir = true
title = "Benchmark"

[languages]
[languages.en]
weight = 1


`

	if multipleDimensions {
		files += `
[languages.nn]
weight = 2
[languages.sv]
weight = 3

[versions]
[versions."v1.0.0"]
[versions."v2.0.0"]

[roles]
[roles.guest]
weight = 300
[roles.member]
weight = 200
`
	}

	files += `

[[module.mounts]]
source = 'content'
target = 'content'
[module.mounts.sites.matrix]
languages = ["**"]
versions  = ["**"]
roles = ["**"]
-- layouts/all.html --
All. {{ .Title }}|
`
	if multipleDimensions {
		files += `
Rotate(language): {{ with .Rotate "language" }}{{ range . }}{{ template "printp" . }}|{{ end }}{{ end }}$
Rotate(version): {{ with .Rotate "version" }}{{ range . }}{{ template "printp" . }}|{{ end }}{{ end }}$
Rotate(role): {{ with .Rotate "role" }}{{ range . }}{{ template "printp" . }}|{{ end }}{{ end }}$
{{ define "printp" }}{{ .RelPermalink }}:{{ with .Site }}{{ template "prints" . }}{{ end }}{{ end }}
{{ define "prints" }}/l:{{ .Language.Name }}/v:{{ .Version.Name }}/r:{{ .Role.Name }}{{ end }}
 
`
	}

	for i := range numPages {
		files += fmt.Sprintf(`
-- content/p%d/index.md --
---
title: "P%d"
---
`, i+1, i+1)
	}

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			BuildCfg: hugolib.BuildCfg{
				SkipRender: skipRender,
			},
		})
	return b
}

// See #14132. We recently reworked the config structs for languages, versions, and roles,
// which made them incomplete when generating the docshelper YAML file.
// Add a test here to ensure we don't regress.
func TestUnmarshalSitesMatrixConfig(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
defaultContentLanguage = "en"
defaultContentLanguageInSubDir = true
defaultCOntentVersionInSubDir = true
defaultContentVersion = "v1.0.0"
defaultContentRole = "guest"
defaultContentRoleInSubDir = true

[moule.mounts]
source = 'content'
target = 'content'


[languages]
[languages.en]

[versions]
[versions."v1.0.0"]

[roles]
[roles.guest]

`

	b := hugolib.Test(t, files)

	toJSONAndMap := func(v any) map[string]any {
		bb, err := json.Marshal(v)
		b.Assert(err, qt.IsNil)
		var m map[string]any
		err = json.Unmarshal(bb, &m)
		b.Assert(err, qt.IsNil)
		return m
	}

	conf := b.H.Configs.Base

	b.Assert(toJSONAndMap(conf.Languages), qt.DeepEquals,
		map[string]any{
			"en": map[string]any{
				"Disabled":          bool(false),
				"LanguageCode":      "",
				"LanguageDirection": "",
				"LanguageName":      "",
				"Title":             "",
				"Weight":            float64(0),
			},
		})

	b.Assert(toJSONAndMap(conf.Versions), qt.DeepEquals, map[string]any{
		"v1.0.0": map[string]any{
			"Weight": float64(0),
		},
	})

	b.Assert(toJSONAndMap(conf.Roles), qt.DeepEquals, map[string]any{
		"guest": map[string]any{
			"Weight": float64(0),
		},
	})

	firstMount := conf.Module.Mounts[0]
	b.Assert(toJSONAndMap(firstMount.Sites.Matrix), qt.DeepEquals, map[string]any{
		"languages": nil,
		"versions":  nil,
		"roles":     nil,
	})
	b.Assert(toJSONAndMap(firstMount.Sites.Complements), qt.DeepEquals, map[string]any{
		"languages": nil,
		"versions":  nil,
		"roles":     nil,
	})
}

func TestSitesMatrixContentBenchmark(t *testing.T) {
	const numPages = 3
	b := newSitesMatrixContentBenchmarkBuilder(t, numPages, false, true)

	b.Build()

	for _, lang := range []string{"en", "nn", "sv"} {
		for _, ver := range []string{"v1.0.0", "v2.0.0"} {
			for _, role := range []string{"guest", "member"} {
				base := "public/" + role + "/" + ver + "/" + lang
				b.AssertFileContent(base+"/index.html", "All. Benchmark|")
				for i := range numPages {
					b.AssertFileContent(fmt.Sprintf("%s/p%d/index.html", base, i+1), fmt.Sprintf("All. P%d|", i+1))
				}

			}
		}
	}
}

func BenchmarkSitesMatrixContent(b *testing.B) {
	for _, numPages := range []int{10, 100} {
		for _, multipleDimensions := range []bool{false, true} {
			b.Run(fmt.Sprintf("n%d/md%t", numPages, multipleDimensions), func(b *testing.B) {
				for b.Loop() {
					b.StopTimer()
					bb := newSitesMatrixContentBenchmarkBuilder(b, numPages, true, multipleDimensions)
					b.StartTimer()
					bb.Build()
				}
			})
		}
	}
}

// Just a test to test the development of concurrency in page assembly.
func TestCreateAllPagesPartitionSections(t *testing.T) {
	files := `
-- hugo.toml --
baseURL = "https://example.org/"
disableKinds = ["rss", "sitemap", "taxonomy", "term"]
-- content/_index.md --
-- content/s1/_index.md --
-- content/s1/p1.md --
-- content/s1/p2.md --
-- content/s2/_index.md --
-- content/s2/p1.md --
-- content/s2_1/_index.md --
-- content/s2_1/p1.md --
-- layouts/all.html --
{{ .Kind }}|{{ .RelPermalink }}|
`

	for range 3 {
		b := hugolib.Test(t, files)
		b.AssertFileContent("public/s1/index.html", "section|/s1/|")
	}
}
