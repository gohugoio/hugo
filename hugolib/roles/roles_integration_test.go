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

package roles_test

import (
	"strings"
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

// TODO1 role hierarchy. Maybe.
// TODO1 #13663 negate.
// TODO1 test content adapter incl. cascade from config,
// TODO1 throw error when using any of the se new slices in cascade other than the config.
func TestRolesAndVersions(t *testing.T) {
	// TODO1 for resources, don't apply a default lang,role, etc. Insert with -1 as a null value.
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
---
Users with guest role will see this.
-- layouts/all.html --
Role: {{ .Site.Role.Name }}|Version: {{ .Site.Version.Name }}|Lang: {{ .Site.Language.Lang }}|
Roles: {{ range .Site.Roles }}Name: {{ .Name }} Site.Version: {{.Site.Version.Name }} Site.Language.Lang: {{ .Site.Language.Lang}}|{{ end }}$
Versions: {{ range site.Versions }}Name: {{ .Name }} Site.Role: {{ .Site.Role.Name }} Site.Language.Lang: {{ .Site.Language.Lang }}|{{ end }}$
RegularPages: {{ range .RegularPages }}{{ .RelPermalink }} r: {{ .Site.Language.Name }}  v: {{ .Site.Version.Name }} l: {{ .Site.Role.Name }}|{{ end }}$

`

	for range 3 {
		b := hugolib.Test(t, files)

		// TODO1 export Default?
		// /guest/v1.2.3/en/publicpost/index.html
		// TODO1 redirect Aliases.

		b.AssertPublishDir(
			"guest/v1.2.3/en/publicpost", "guest/v2.0.0/en/publicpost", "! guest/v2.1.0/en/publicpost",
			"member/v4.0.0/en/memberonlypost", "member/v4.0.0/nn/memberonlypost",
		)

		b.AssertFileContent("public/guest/v2.0.0/en/index.html",
			"Role: guest|Version: v2.0.0|",
			"Roles: Name: member Site.Version: v2.0.0 Site.Language.Lang: en|Name: guest Site.Version: v2.0.0 Site.Language.Lang: en|$",
			"Versions: Name: v4.0.0 Site.Role: guest Site.Language.Lang: en|Name: v3.0.0 Site.Role: guest Site.Language.Lang: en|Name: v2.1.0 Site.Role: guest Site.Language.Lang: en|Name: v2.0.0 Site.Role: guest Site.Language.Lang: en|Name: v1.2.3 Site.Role: guest Site.Language.Lang: en|$")

		b.AssertFileContent("public/guest/v3.0.0/en/index.html", "RegularPages: /guest/v2.0.0/en/publicpost/ r: en  v: v2.0.0 l: guest|/guest/v3.0.0/en/v3publicpost/ r: en  v: v3.0.0 l: guest|$")

	}
}

// TODO1 add dimensions to cascade config (instead of lang).

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
Rotate(version): {{/* with .Rotate "version" }}{{ range . }}{{ template "printp" . }}|{{ end }}{{ end */}}$
Rotate(role): {{/* with .Rotate "role" }}{{ range . }}{{ template "printp" . }}|{{ end }}{{ end */}}$
{{ define "printp" }}{{ .RelPermalink }}:{{ with .Site }}{{ template "prints" . }}{{ end }}{{ end }}
{{ define "prints" }}/l:{{ .Language.Name }}/v:{{ .Version.Name }}/r:{{ .Role.Name }}{{ end }}


`

	for range 3 {
		b := hugolib.Test(t, files)

		b.AssertFileContent("public/guest/v3.0.0/en/index.html",
			"Rotate(language): /guest/v3.0.0/en/:/l:en/v:v3.0.0/r:guest|/guest/v3.0.0/nn/:/l:nn/v:v3.0.0/r:guest|$",
			//"Rotate(version): /guest/v4.0.0/en/:/l:en/v:v4.0.0/r:guest|/guest/v3.0.0/en/:/l:en/v:v3.0.0/r:guest|/guest/v2.1.0/en/:/l:en/v:v2.1.0/r:guest|/guest/v2.0.0/en/:/l:en/v:v2.0.0/r:guest|/guest/v1.2.3/en/:/l:en/v:v1.2.3/r:guest",
			//"Rotate(role): /member/v3.0.0/en/:/l:en/v:v3.0.0/r:member|/guest/v3.0.0/en/:/l:en/v:v3.0.0/r:guest|$",
		)

	}
}

/*

Notes new API:

* Page.Rotate
* Site.Role
* Site.Version
* Site.Language
* Role.Name, Version.Name, Language.Name (new).


Remove: Site.Versions, Roles and (if new) Languages.


*/

// TODO1 check defaultCOntentVersionInSubDir = false vs language.
func TestDimensionsFileMount(t *testing.T) {
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
[module.mounts.sites]
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

	testOne := func(t *testing.T, files string) {
		t.Helper()
		b := hugolib.Test(t, files)

		// b.AssertPublishDir("asdf")
		// b.AssertFileContent("public/v1.2.3/nn/p1/index.html", "asdfasfd")
		// b.AssertFileContent("public/en/index.html", "Title English", "Text English")
		b.AssertFileContent("public/v2.0.0/nn/p1/index.html", "Tittel Nynorsk", "Tekst Nynorsk", "site.GetPage p1: Tittel Nynorsk|")
		b.AssertFileContent("public/v2.0.0/nn/p2/index.html", "p2 all||", "site.GetPage p2: p2 all", "site.GetPage p1: Tittel Nynorsk|")
		b.AssertFileContent("public/v2.0.0/en/p2/index.html", "p2 all||", "site.GetPage p1: $")
		b.AssertFileContent("public/v2.0.0/nn/p2/index.html", "p2 all||", "site.GetPage p1: Tittel Nynorsk|$")
		b.AssertFileContent("public/v1.2.3/en/p2/index.html", "p2 all||", "site.GetPage p2: p2 all")
	}

	// Format from v0.148.0:
	dims := `[module.mounts.sites]
languages = ["en"]
versions = ["v1**"]
`
	files := strings.Replace(filesTemplate, "DIMSEN", dims, 1)
	dims = strings.Replace(dims, `["en"]`, `["nn"]`, 1)
	dims = strings.Replace(dims, `["v1**"]`, `["v2**"]`, 1)
	files = strings.Replace(files, "DIMSNN", dims, 1)
	testOne(t, files)

	if true {
		return
	}

	// Old format:
	files = strings.Replace(filesTemplate, "DIMSEN", `lang = "en"`, 1)
	files = strings.Replace(files, "DIMSNN", `lang = "nn"`, 1)
	testOne(t, files)
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
[module.mounts.sites]
languages = ["nn"]
versions  = ["v1.**"]
[[module.mounts]]
source = 'content/en'
target = 'content'
[module.mounts.sites]
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

const filesVariationsSiteMatrix = `
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
[module.mounts.sites]
languages = ["nn"]
versions  = ["v1.**"]
[[module.mounts]]
source = 'content/en'
target = 'content'
[module.mounts.sites]
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

func TestFrontMatterSites(t *testing.T) {
	t.Parallel()

	files := filesVariationsSiteMatrix

	files += `
-- content/other/p2.md --
---
title: "NN p2"
sites:
   languages: ["nn"]
   versions: ["v1.2.3"]
---
`
	b := hugolib.Test(t, files)
	// b.AssertPublishDir("asdf")
	b.AssertFileContent("public/v2.0.0/nn/p2/index.html", "title: NN p2|")
}
