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
	"testing"
)

func TestLanguageContentRoot(t *testing.T) {
	files := `
-- hugo.toml --
baseURL = "https://example.org/"
defaultContentLanguage = "en"
defaultContentLanguageInSubdir = true
[languages]
[languages.en]
weight = 10
contentDir = "content/en"
[languages.nn]
weight = 20
contentDir = "content/nn"
-- content/en/_index.md --
---
title: "Home"
---
-- content/nn/_index.md --
---
title: "Heim"
---
-- content/en/myfiles/file1.txt --
file 1 en
-- content/en/myfiles/file2.txt --
file 2 en
-- content/nn/myfiles/file1.txt --
file 1 nn
-- layouts/index.html --
Title: {{ .Title }}|
Len Resources: {{ len .Resources }}|
{{ range $i, $e := .Resources }}
{{ $i }}|{{ .RelPermalink }}|{{ .Content }}|
{{ end }}

`
	b := Test(t, files)
	b.AssertFileContent("public/en/index.html", "Home", "0|/en/myfiles/file1.txt|file 1 en|\n\n1|/en/myfiles/file2.txt|file 2 en|")
	b.AssertFileContent("public/nn/index.html", "Heim", "0|/nn/myfiles/file1.txt|file 1 nn|\n\n1|/en/myfiles/file2.txt|file 2 en|")
}

func TestContentMountMerge(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
baseURL = 'https://example.org/'
languageCode = 'en-us'
title = 'Hugo Forum Topic #37225'
theme = 'mytheme'

disableKinds = ['sitemap','RSS','taxonomy','term']
defaultContentLanguage = 'en'
defaultContentLanguageInSubdir = true

[languages.en]
languageName = 'English'
weight = 1
[languages.de]
languageName = 'Deutsch'
weight = 2
[languages.nl]
languageName = 'Nederlands'
weight = 3

# EN content
[[module.mounts]]
source = 'content/en'
target = 'content'
lang = 'en'

# DE content
[[module.mounts]]
source = 'content/de'
target = 'content'
lang = 'de'

# This fills in the gaps in DE content with EN content
[[module.mounts]]
source = 'content/en'
target = 'content'
lang = 'de'

# NL content
[[module.mounts]]
source = 'content/nl'
target = 'content'
lang = 'nl'

# This should fill in the gaps in NL content with EN content
[[module.mounts]]
source = 'content/en'
target = 'content'
lang = 'nl'

-- content/de/_index.md --
---
title: "home (de)"
---
-- content/de/p1.md --
---
title: "p1 (de)"
---
-- content/en/_index.md --
---
title: "home (en)"
---
-- content/en/p1.md --
---
title: "p1 (en)"
---
-- content/en/p2.md --
---
title: "p2 (en)"
---
-- content/en/p3.md --
---
title: "p3 (en)"
---
-- content/nl/_index.md --
---
title: "home (nl)"
---
-- content/nl/p1.md --
---
title: "p1 (nl)"
---
-- content/nl/p3.md --
---
title: "p3 (nl)"
---
-- layouts/home.html --
{{ .Title }}: {{ site.Language.Lang }}: {{ range site.RegularPages }}{{ .Title }}|{{ end }}:END
-- themes/mytheme/config.toml --
[[module.mounts]]
source = 'content/nlt'
target = 'content'
lang = 'nl'
-- themes/mytheme/content/nlt/p3.md --
---
title: "p3 theme (nl)"
---
-- themes/mytheme/content/nlt/p4.md --
---
title: "p4 theme (nl)"
---
`

	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,
		},
	).Build()

	b.AssertFileContent("public/nl/index.html", `home (nl): nl: p1 (nl)|p2 (en)|p3 (nl)|p4 theme (nl)|:END`)
	b.AssertFileContent("public/de/index.html", `home (de): de: p1 (de)|p2 (en)|p3 (en)|:END`)
	b.AssertFileContent("public/en/index.html", `home (en): en: p1 (en)|p2 (en)|p3 (en)|:END`)
}

// Issue 13993
func TestIssue13993(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['home','rss','section','sitemap','taxonomy','term']
printPathWarnings = true

defaultContentLanguage = 'en'
defaultContentLanguageInSubdir = true

[languages.en]
contentDir = "content/en"
weight = 1

[languages.es]
contentDir = "content/es"
weight = 2

# Default content mounts

[[module.mounts]]
source = "content/en"
target = "content"
lang = "en"

[[module.mounts]]
source = "content/es"
target = "content"
lang = "es"

# Populate the missing es content with en content

[[module.mounts]]
source = "content/en"
target = "content"
lang = "es"
-- layouts/all.html --
{{ .Title }}
-- content/en/p1.md --
---
title: p1 (en)
---
-- content/en/p2.md --
---
title: p2 (en)
---
-- content/es/p1.md --
---
title: p1 (es)
---
`

	b := Test(t, files, TestOptInfo())

	b.AssertFileExists("public/en/p1/index.html", true)
	b.AssertFileExists("public/en/p2/index.html", true)
	b.AssertFileExists("public/es/p1/index.html", true)
	b.AssertFileExists("public/es/p2/index.html", true)

	b.AssertLogContains("INFO  Duplicate")
	b.AssertLogContains("! WARN  Duplicate")
}
