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
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/gohugoio/hugo/hugolib"
)

func TestRolesAndVersions(t *testing.T) {
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

[cascade.sites.matrix]
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
sites:
  matrix:
    languages: "**"
    roles: member
    versions: "v4.0.0"
---
Member content.
-- content/publicpost.md --
---
title: "Public"
sites:
  matrix:
    versions: ["v1.2.3", "v2.**", "! v2.1.*"]
  complements:
    versions: "v3**"
---
Users with guest role will see this.
-- content/v3publicpost.md --
---
title: "Public v3"
sites:
  matrix:
    versions: "v3**"
---
Users with guest role will see this.
-- layouts/all.html --
Role: {{ .Site.Role.Name }}|Version: {{ .Site.Version.Name }}|Lang: {{ .Site.Language.Lang }}|
RegularPages: {{ range .RegularPages }}{{ .RelPermalink }} r: {{ .Site.Language.Name }}  v: {{ .Site.Version.Name }} l: {{ .Site.Role.Name }}|{{ end }}$

`

	for range 3 {
		b := hugolib.Test(t, files)

		b.AssertPublishDir(
			"guest/v1.2.3/en/publicpost", "guest/v2.0.0/en/publicpost", "! guest/v2.1.0/en/publicpost",
			"member/v4.0.0/en/memberonlypost", "member/v4.0.0/nn/memberonlypost",
		)

		b.AssertFileContent("public/guest/v2.0.0/en/index.html",
			"Role: guest|Version: v2.0.0|",
		)

		b.AssertFileContent("public/guest/v3.0.0/en/index.html", "egularPages: /guest/v2.0.0/en/publicpost/ r: en  v: v2.0.0 l: guest|/guest/v3.0.0/en/v3publicpost/ r: en  v: v3.0.0 l: guest|$")

	}
}

func TestDefaultContentRoleDoesNotExist(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ["rss", "sitemap", "section", "taxonomy", "term"]
defaultContentLanguage = "en"
defaultContentLanguageInSubDir = true
defaultContentRole = "doesnotexist"
[roles]
[roles.guest]
weight = 300
[roles.member]
weight = 200
`
	b, err := hugolib.TestE(t, files)
	b.Assert(err, qt.ErrorMatches, `.*failed to decode "roles": the configured defaultContentRole "doesnotexist" does not exist`)
}
