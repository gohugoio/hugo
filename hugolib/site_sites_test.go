// Copyright 2026 The Hugo Authors. All rights reserved.
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

import "testing"

func TestSiteIsDefault(t *testing.T) {
	files := `
-- hugo.toml --
disableKinds = ['rss','sitemap','taxonomy','term']
defaultContentLanguage = 'fr'
defaultContentLanguageInSubdir = true
defaultContentVersion = "v2.0.0"
defaultContentVersionInSubdir = true
defaultContentRoleInSubdir = true
[languages]
[languages.en]
weight = 1
title = 'English'
[languages.fr]
weight = 2
title = 'French'
[languages.de]
weight = 3
title = 'German'
[roles]
[roles.guest]
weight = 1
[roles.member]
weight = 2
[versions]
[versions.'v1.0.0']
weight = 1
[versions.'v2.0.0']
weight = 2
-- content/p1.en.md --
---
title: Page 1 EN
---
-- content/p1.fr.md --
---
title: Page 1 FR
---
-- content/p1.de.md --
---
title: Page 1 DE
---
-- layouts/_default/single.html --
Current site is default: {{ .Site.IsDefault }}
{{ with hugo.Sites.Default }}
Default site: {{ .Language.Name }}-{{ .Role.Name }}-{{ .Version.Name }}: IsDefault={{ .IsDefault }}
{{ end }}
{{ range hugo.Sites -}}
{{ .Language.Name }}-{{ .Role.Name }}-{{ .Version.Name }}: IsDefault={{ .IsDefault }}
{{ end }}
`

	b := Test(t, files)
	b.AssertFileContent("public/guest/v2.0.0/en/p1/index.html",
		`
Current site is default: false
Default site: fr-guest-v2.0.0: IsDefault=true
en-guest-v1.0.0: IsDefault=false
en-member-v1.0.0: IsDefault=false
en-guest-v2.0.0: IsDefault=false
en-member-v2.0.0: IsDefault=false
fr-guest-v1.0.0: IsDefault=false
fr-member-v1.0.0: IsDefault=false
fr-guest-v2.0.0: IsDefault=true
fr-member-v2.0.0: IsDefault=false
de-guest-v1.0.0: IsDefault=false
de-member-v1.0.0: IsDefault=false
de-guest-v2.0.0: IsDefault=false
de-member-v2.0.0: IsDefault=false		
`,
	)
	b.AssertFileContent("public/guest/v2.0.0/fr/p1/index.html",
		`
Current site is default: true
Default site: fr-guest-v2.0.0: IsDefault=true
en-guest-v1.0.0: IsDefault=false
en-member-v1.0.0: IsDefault=false
en-guest-v2.0.0: IsDefault=false
en-member-v2.0.0: IsDefault=false
fr-guest-v1.0.0: IsDefault=false
fr-member-v1.0.0: IsDefault=false
fr-guest-v2.0.0: IsDefault=true
fr-member-v2.0.0: IsDefault=false
de-guest-v1.0.0: IsDefault=false
de-member-v1.0.0: IsDefault=false
de-guest-v2.0.0: IsDefault=false
de-member-v2.0.0: IsDefault=false
`,
	)
	b.AssertFileContent("public/guest/v2.0.0/de/p1/index.html",
		`
Current site is default: false
Default site: fr-guest-v2.0.0: IsDefault=true
en-guest-v1.0.0: IsDefault=false
en-member-v1.0.0: IsDefault=false
en-guest-v2.0.0: IsDefault=false
en-member-v2.0.0: IsDefault=false
fr-guest-v1.0.0: IsDefault=false
fr-member-v1.0.0: IsDefault=false
fr-guest-v2.0.0: IsDefault=true
fr-member-v2.0.0: IsDefault=false
de-guest-v1.0.0: IsDefault=false
de-member-v1.0.0: IsDefault=false
de-guest-v2.0.0: IsDefault=false
de-member-v2.0.0: IsDefault=false
`,
	)
}

func TestSiteSites(t *testing.T) {
	files := `
-- hugo.toml --
disableKinds = ['page','rss','section','sitemap','taxonomy','term']
defaultContentLanguage = 'fr'
defaultContentLanguageInSubdir = true
defaultContentVersionInSubdir = true
defaultContentRoleInSubdir = true
[languages]
[languages.en]
weight = 1
[languages.fr]
weight = 2
[languages.de]
weight = 3
[roles]
[roles.guest]
weight = 1
[roles.member]
weight = 2
[versions]
[versions.'v1.0.0']
weight = 1
[versions.'v2.0.0']
weight = 2
-- layouts/home.html --
{{ range .Site.Sites }}{{ .Language.Name }}-{{ .Role.Name }}-{{ .Version.Name }}|{{ end }}
	`
	b := Test(t, files, TestOptInfo())

	b.AssertFileContent("public/guest/v1.0.0/fr/index.html", "en-guest-v1.0.0|en-member-v1.0.0|en-guest-v2.0.0|en-member-v2.0.0|fr-guest-v1.0.0|fr-member-v1.0.0|fr-guest-v2.0.0|fr-member-v2.0.0|de-guest-v1.0.0|de-member-v1.0.0|de-guest-v2.0.0|de-member-v2.0.0|")
	b.AssertLogContains(".Site.Sites was deprecated")
}
