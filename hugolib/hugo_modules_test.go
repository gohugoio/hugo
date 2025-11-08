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

	qt "github.com/frankban/quicktest"
)

func TestModulesWithContent(t *testing.T) {
	t.Parallel()

	content := func(id string) string {
		return fmt.Sprintf(`---
title: Title %s
---
Content %s

`, id, id)
	}

	i18nContent := func(id, value string) string {
		return fmt.Sprintf(`
[%s]
other = %q
`, id, value)
	}

	files := `
-- hugo.toml --
baseURL="https://example.org"

defaultContentLanguage = "en"

[module]
[[module.imports]]
path="a"
[[module.imports.mounts]]
source="myacontent"
target="content/blog"
lang="en"
[[module.imports]]
path="b"
[[module.imports.mounts]]
source="mybcontent"
target="content/blog"
lang="nn"
[[module.imports]]
path="c"
[[module.imports]]
path="d"

[languages]

[languages.en]
title = "Title in English"
languageName = "English"
weight = 1
[languages.nn]
languageName = "Nynorsk"
weight = 2
title = "Tittel på nynorsk"
[languages.nb]
languageName = "Bokmål"
weight = 3
title = "Tittel på bokmål"
[languages.fr]
languageName = "French"
weight = 4
title = "French Title"

-- layouts/index.html --
{{ range .Site.RegularPages }}
|{{ .Title }}|{{ .RelPermalink }}|{{ .Plain }}
{{ end }}
{{ $data := .Site.Data }}
Data Common: {{ $data.common.value }}
Data C: {{ $data.c.value }}
Data D: {{ $data.d.value }}
All Data: {{ $data }}

i18n hello1: {{ i18n "hello1" . }}
i18n theme: {{ i18n "theme" . }}
i18n theme2: {{ i18n "theme2" . }}
-- themes/a/myacontent/page.md --
` + content("theme-a-en") + `
-- themes/b/mybcontent/page.md --
` + content("theme-b-nn") + `
-- themes/c/content/blog/c.md --
` + content("theme-c-nn") + `
-- data/common.toml --
value="Project"
-- themes/c/data/common.toml --
value="Theme C"
-- themes/c/data/c.toml --
value="Hugo Rocks!"
-- themes/d/data/c.toml --
value="Hugo Rodcks!"
-- themes/d/data/d.toml --
value="Hugo Rodks!"
-- i18n/en.toml --
` + i18nContent("hello1", "Project") + `
-- themes/c/i18n/en.toml --
[hello1]
other="Theme C Hello"
[theme]
other="Theme C"
-- themes/d/i18n/en.toml --
` + i18nContent("theme", "Theme D") + `
-- themes/d/i18n/en.toml --
` + i18nContent("theme2", "Theme2 D") + `
-- themes/c/static/hello.txt --
Hugo Rocks!"
`
	b := Test(t, files)

	b.AssertFileContent("public/index.html", "|Title theme-a-en|/blog/page/|Content theme-a-en")
	b.AssertFileContent("public/nn/index.html", "|Title theme-b-nn|/nn/blog/page/|Content theme-b-nn")

	// Data
	b.AssertFileContent("public/index.html",
		"Data Common: Project",
		"Data C: Hugo Rocks!",
		"Data D: Hugo Rodks!",
	)

	// i18n
	b.AssertFileContent("public/index.html",
		"i18n hello1: Project",
		"i18n theme: Theme C",
		"i18n theme2: Theme2 D",
	)
}

func TestModulesIgnoreConfig(t *testing.T) {
	files := `
-- hugo.toml --
baseURL="https://example.org"

[module]
[[module.imports]]
path="a"
ignoreConfig=true

-- themes/a/config.toml --
[params]
a = "Should Be Ignored!"
-- layouts/index.html --
Params: {{ .Site.Params }}
`
	Test(t, files).AssertFileContent("public/index.html", "! Ignored")
}

func TestModulesDisabled(t *testing.T) {
	files := `
-- hugo.toml --
baseURL="https://example.org"

[module]
[[module.imports]]
path="a"
[[module.imports]]
path="b"
disable=true

-- themes/a/config.toml --
[params]
a = "A param"
-- themes/b/config.toml --
[params]
b = "B param"
-- layouts/index.html --
Params: {{ .Site.Params }}
`
	Test(t, files).AssertFileContent("public/index.html", "A param", "! B param")
}

func TestModulesIncompatible(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL="https://example.org"

[module]
[[module.imports]]
path="ok"
[[module.imports]]
path="incompat1"
[[module.imports]]
path="incompat2"
[[module.imports]]
path="incompat3"

-- themes/ok/data/ok.toml --
title = "OK"
-- themes/incompat1/config.toml --

[module]
[module.hugoVersion]
min = "0.33.2"
max = "0.45.0"

-- themes/incompat2/theme.toml --
min_version = "5.0.0"

-- themes/incompat3/theme.toml --
min_version = 0.55.0

`
	b := Test(t, files, TestOptWarn())
	b.AssertLogContains("is not compatible with this Hugo version")
}

func TestMountsProject(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
baseURL="https://example.org"

[module]
[[module.mounts]]
source="mycontent"
target="content"
-- layouts/_default/single.html --
Permalink: {{ .Permalink }}|
-- mycontent/mypage.md --
---
title: "My Page"
---
`
	b := Test(t, files)

	b.AssertFileContent("public/mypage/index.html", "Permalink: https://example.org/mypage/|")
}

// https://github.com/gohugoio/hugo/issues/6684
func TestMountsContentFile(t *testing.T) {
	files := `
-- hugo.toml --
disableKinds = ["taxonomy", "term", "RSS", "sitemap", "robotsTXT", "page", "section"]
disableLiveReload = true
[module]
[[module.mounts]]
source = "README.md"
target = "content/_index.md"	
-- README.md --
# Hello World
-- layouts/index.html --
Home: {{ .Title }}|{{ .Content }}|
`
	b := Test(t, files)
	b.AssertFileContent("public/index.html", "Home: |<h1 id=\"hello-world\">Hello World</h1>\n|")
}

// https://github.com/gohugoio/hugo/issues/6299
func TestSiteWithGoModButNoModules(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.org"
-- go.mod --

`

	b, err := TestE(t, files, TestOptOsFs())
	b.Assert(err, qt.IsNil)
}

// Issue 9426
func TestMountSameSource(t *testing.T) {
	files := `
-- hugo.toml --
baseURL = 'https://example.org/'
languageCode = 'en-us'
title = 'Hugo GitHub Issue #9426'

disableKinds = ['RSS','sitemap','taxonomy','term']

[[module.mounts]]
source = "content"
target = "content"

[[module.mounts]]
source = "extra-content"
target = "content/resources-a"

[[module.mounts]]
source = "extra-content"
target = "content/resources-b"
-- layouts/_default/single.html --
Single
-- content/p1.md --
-- extra-content/_index.md --
-- extra-content/subdir/_index.md --
-- extra-content/subdir/about.md --
"
`
	b := Test(t, files)

	b.AssertFileContent("public/resources-a/subdir/about/index.html", "Single")
	b.AssertFileContent("public/resources-b/subdir/about/index.html", "Single")
}

func TestMountData(t *testing.T) {
	files := `
-- hugo.toml --
baseURL = 'https://example.org/'
disableKinds = ["taxonomy", "term", "RSS", "sitemap", "robotsTXT", "page", "section"]

[[module.mounts]]
source = "data"
target = "data"

[[module.mounts]]
source = "extra-data"
target = "data/extra"
-- extra-data/test.yaml --
message: Hugo Rocks
-- layouts/index.html --
{{ site.Data.extra.test.message }}
`

	b := Test(t, files)

	b.AssertFileContent("public/index.html", "Hugo Rocks")
}
