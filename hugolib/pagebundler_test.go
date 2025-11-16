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
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestPageBundlerBasic(t *testing.T) {
	files := `
-- hugo.toml --
-- content/mybundle/index.md --
---
title: "My Bundle"
---
-- content/mybundle/p1.md --
---
title: "P1"
---
-- content/mybundle/p1.html --
---
title: "P1 HTML"
---
-- content/mybundle/data.txt --
Data txt.
-- content/mybundle/sub/data.txt --
Data sub txt.
-- content/mybundle/sub/index.md --
-- content/mysection/_index.md --
---
title: "My Section"
---
-- content/mysection/data.txt --
Data section txt.
-- layouts/all.html --
All. {{ .Title }}|{{ .RelPermalink }}|{{ .Kind }}|{{ .BundleType }}|
`

	b := Test(t, files)

	b.AssertFileContent("public/mybundle/index.html", "My Bundle|/mybundle/|page|leaf|")
}

func TestPageBundlerBundleInRoot(t *testing.T) {
	t.Parallel()
	files := `
-- hugo.toml --
baseURL = "https://example.com"
disableKinds = ["taxonomy", "term"]
-- content/root/index.md --
---
title: "Root"
---
-- layouts/_default/single.html --
Basic: {{ .Title }}|{{ .Kind }}|{{ .BundleType }}|{{ .RelPermalink }}|
Tree: Section: {{ .Section }}|CurrentSection: {{ .CurrentSection.RelPermalink }}|Parent: {{ .Parent.RelPermalink }}|FirstSection: {{ .FirstSection.RelPermalink }}
`
	b := Test(t, files)

	b.AssertFileContent("public/root/index.html",
		"Basic: Root|page|leaf|/root/|",
		"Tree: Section: |CurrentSection: /|Parent: /|FirstSection: /",
	)
}

func TestPageBundlerShortcodeInBundledPage(t *testing.T) {
	t.Parallel()
	files := `
-- hugo.toml --
baseURL = "https://example.com"
disableKinds = ["taxonomy", "term"]
-- content/section/mybundle/index.md --
---
title: "Mybundle"
---
-- content/section/mybundle/p1.md --
---
title: "P1"
---

P1 content.

{{< myShort >}}

-- layouts/_default/single.html --
Bundled page: {{ .RelPermalink}}|{{ with .Resources.Get "p1.md" }}Title: {{ .Title }}|Content: {{ .Content }}{{ end }}|
-- layouts/shortcodes/myShort.html --
MyShort.

`
	b := Test(t, files)

	b.AssertFileContent("public/section/mybundle/index.html",
		"Bundled page: /section/mybundle/|Title: P1|Content: <p>P1 content.</p>\nMyShort.",
	)
}

func TestPageBundlerResourceMultipleOutputFormatsWithDifferentPaths(t *testing.T) {
	t.Parallel()
	files := `
-- hugo.toml --
baseURL = "https://example.com"
disableKinds = ["taxonomy", "term"]
[outputformats]
[outputformats.cpath]
mediaType = "text/html"
path = "cpath"
-- content/section/mybundle/index.md --
---
title: "My Bundle"
outputs: ["html", "cpath"]
---
-- content/section/mybundle/hello.txt --
Hello.
-- content/section/mybundle/p1.md --
---
title: "P1"
---
P1.

{{< hello >}}

-- layouts/shortcodes/hello.html --
Hello HTML.
-- layouts/_default/single.html --
Basic: {{ .Title }}|{{ .Kind }}|{{ .BundleType }}|{{ .RelPermalink }}|
Resources: {{ range .Resources }}RelPermalink: {{ .RelPermalink }}|Content: {{ .Content }}|{{ end }}|
-- layouts/shortcodes/hello.cpath --
Hello CPATH.
-- layouts/_default/single.cpath --
Basic: {{ .Title }}|{{ .Kind }}|{{ .BundleType }}|{{ .RelPermalink }}|
Resources: {{ range .Resources }}RelPermalink: {{ .RelPermalink }}|Content: {{ .Content }}|{{ end }}|
`

	b := Test(t, files)

	b.AssertFileContent("public/section/mybundle/index.html",
		"Basic: My Bundle|page|leaf|/section/mybundle/|",
		"Resources: RelPermalink: |Content: <p>P1.</p>\nHello HTML.\n|RelPermalink: /section/mybundle/hello.txt|Content: Hello.||",
	)

	b.AssertFileContent("public/cpath/section/mybundle/index.html", "Basic: My Bundle|page|leaf|/section/mybundle/|\nResources: RelPermalink: |Content: <p>P1.</p>\nHello CPATH.\n|RelPermalink: /section/mybundle/hello.txt|Content: Hello.||")
}

func TestPageBundlerMultilingualTextResource(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com"
disableKinds = ["taxonomy", "term"]
defaultContentLanguage = "en"
defaultContentLanguageInSubdir = true
[languages]
[languages.en]
weight = 1
[languages.en.permalinks]
"/" = "/enpages/:slug/"
[languages.nn]
weight = 2
-- content/mybundle/index.md --
---
title: "My Bundle"
---
-- content/mybundle/index.nn.md --
---
title: "My Bundle NN"
---
-- content/mybundle/f1.txt --
F1
-- content/mybundle/f2.txt --
F2
-- content/mybundle/f2.nn.txt --
F2 nn.
-- layouts/_default/single.html --
{{ .Title }}|{{ .RelPermalink }}|{{ .Lang }}|
Resources: {{ range .Resources }}RelPermalink: {{ .RelPermalink }}|Content: {{ .Content }}|{{ end }}|

`
	b := Test(t, files)

	b.AssertFileContent("public/en/enpages/my-bundle/index.html", "My Bundle|/en/enpages/my-bundle/|en|\nResources: RelPermalink: /en/enpages/my-bundle/f1.txt|Content: F1|RelPermalink: /en/enpages/my-bundle/f2.txt|Content: F2||")
	b.AssertFileContent("public/nn/mybundle/index.html", "My Bundle NN|/nn/mybundle/|nn|\nResources: RelPermalink: /en/enpages/my-bundle/f1.txt|Content: F1|RelPermalink: /nn/mybundle/f2.nn.txt|Content: F2 nn.||")
}

func TestMultilingualDisableLanguage(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com"
disableKinds = ["taxonomy", "term"]
defaultContentLanguage = "en"
defaultContentLanguageInSubdir = true
[languages]
[languages.en]
weight = 1
[languages.nn]
weight = 2
disabled = true
-- content/mysect/_index.md --
---
title: "My Sect En"
---
-- content/mysect/p1/index.md --
---
title: "P1"
---
P1
-- content/mysect/_index.nn.md --
---
title: "My Sect Nn"
---
-- content/mysect/p1/index.nn.md --
---
title: "P1nn"
---
P1nn
-- layouts/index.html --
Len RegularPages: {{ len .Site.RegularPages }}|RegularPages: {{ range site.RegularPages }}{{ .RelPermalink }}: {{ .Title }}|{{ end }}|
Len Pages: {{ len .Site.Pages }}|
Len Sites: {{ len .Site.Sites }}|
-- layouts/_default/single.html --
{{ .Title }}|{{ .Content }}|{{ .Lang }}|

`
	b := Test(t, files)

	b.AssertFileContent("public/en/index.html", "Len RegularPages: 1|")
	b.AssertFileContent("public/en/mysect/p1/index.html", "P1|<p>P1</p>\n|en|")
	b.AssertFileExists("public/public/nn/mysect/p1/index.html", false)
	b.Assert(len(b.H.Sites), qt.Equals, 1)
}

func TestMultilingualDisableLanguageMounts(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com"
disableKinds = ["taxonomy", "term"]
defaultContentLanguage = "en"
defaultContentLanguageInSubdir = true
[[module.mounts]]
source = 'content/nn'
target = 'content'
[module.mounts.sites.matrix]
languages = "nn"
[[module.mounts]]
source = 'content/en'
target = 'content'
[module.mounts.sites.matrix]
languages = "en"

[languages]
[languages.en]
weight = 1
[languages.nn]
weight = 2
disabled = true
-- content/en/mysect/_index.md --
---
title: "My Sect En"
---
-- content/en/mysect/p1/index.md --
---
title: "P1"
---
P1
-- content/nn/mysect/_index.md --
---
title: "My Sect Nn"
---
-- content/nn/mysect/p1/index.md --
---
title: "P1nn"
---
P1nn
-- layouts/index.html --
Len RegularPages: {{ len .Site.RegularPages }}|RegularPages: {{ range site.RegularPages }}{{ .RelPermalink }}: {{ .Title }}|{{ end }}|
Len Pages: {{ len .Site.Pages }}|
Len Sites: {{ len .Site.Sites }}|
-- layouts/_default/single.html --
{{ .Title }}|{{ .Content }}|{{ .Lang }}|

`
	b := Test(t, files)

	b.AssertFileContent("public/en/index.html", "Len RegularPages: 1|")
	b.AssertFileContent("public/en/mysect/p1/index.html", "P1|<p>P1</p>\n|en|")
	b.AssertFileExists("public/public/nn/mysect/p1/index.html", false)
	b.Assert(len(b.H.Sites), qt.Equals, 1)
}

func TestPageBundlerHeadlessIssue6552(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "http://example.com/"
-- content/headless/h1/index.md --
---
title: My Headless Bundle1
headless: true
---
-- content/headless/h1/p1.md --
---
title: P1
---
-- content/headless/h2/index.md --
---
title: My Headless Bundle2
headless: true
---
-- layouts/index.html --
{{ $headless1 := .Site.GetPage "headless/h1" }}
{{ $headless2 := .Site.GetPage "headless/h2" }}

HEADLESS1: {{ $headless1.Title }}|{{ $headless1.RelPermalink }}|{{ len $headless1.Resources }}|
HEADLESS2: {{ $headless2.Title }}{{ $headless2.RelPermalink }}|{{ len $headless2.Resources }}|
`
	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			BuildCfg:    BuildCfg{},
		},
	).Build()

	b.AssertFileContent("public/index.html", `
HEADLESS1: My Headless Bundle1||1|
HEADLESS2: My Headless Bundle2|0|
`)
}

func TestMultiSiteBundles(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "http://example.com/"
defaultContentLanguage = "en"

[languages]
[languages.en]
weight = 10
contentDir = "content/en"
[languages.nn]
weight = 20
contentDir = "content/nn"

-- content/en/mybundle/index.md --
---
headless: true
---
-- content/nn/mybundle/index.md --
---
headless: true
---
-- content/en/mybundle/data.yaml --
data en
-- content/en/mybundle/forms.yaml --
forms en
-- content/nn/mybundle/data.yaml --
data nn
-- content/en/_index.md --
---
Title: Home
---

Home content.
-- content/en/section-not-bundle/_index.md --
---
Title: Section Page
---

Section content.
-- content/en/section-not-bundle/single.md --
---
Title: Section Single
Date: 2018-02-01
---

Single content.
-- layouts/_default/single.html --
{{ .Title }}|{{ .Content }}
-- layouts/_default/list.html --
{{ .Title }}|{{ .Content }}
`
	b := Test(t, files)

	b.AssertFileContent("public/nn/mybundle/data.yaml", "data nn")
	b.AssertFileContent("public/mybundle/data.yaml", "data en")
	b.AssertFileContent("public/mybundle/forms.yaml", "forms en")

	b.AssertFileExists("public/nn/nn/mybundle/data.yaml", false)
	b.AssertFileExists("public/en/mybundle/data.yaml", false)

	homeEn := b.H.Sites[0].home
	b.Assert(homeEn, qt.Not(qt.IsNil))
	b.Assert(homeEn.Date().Year(), qt.Equals, 2018)

	b.AssertFileContent("public/section-not-bundle/index.html", "Section Page|<p>Section content.</p>")
	b.AssertFileContent("public/section-not-bundle/single/index.html", "Section Single|<p>Single content.</p>")
}

func TestBundledResourcesMultilingualDuplicateResourceFiles(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com/"
[markup]
[markup.goldmark]
duplicateResourceFiles = true
[languages]
[languages.en]
weight = 1
[languages.en.permalinks]
"/" = "/enpages/:slug/"
[languages.nn]
weight = 2
[languages.nn.permalinks]
"/" = "/nnpages/:slug/"
-- content/mybundle/index.md --
---
title: "My Bundle"
---
{{< getresource "f1.txt" >}}
{{< getresource "f2.txt" >}}
-- content/mybundle/index.nn.md --
---
title: "My Bundle NN"
---
{{< getresource "f1.txt" >}}
f2.nn.txt is the original name.
{{< getresource "f2.nn.txt" >}}
{{< getresource "f2.txt" >}}
{{< getresource "sub/f3.txt" >}}
-- content/mybundle/f1.txt --
F1 en.
-- content/mybundle/sub/f3.txt --
F1 en.
-- content/mybundle/f2.txt --
F2 en.
-- content/mybundle/f2.nn.txt --
F2 nn.
-- layouts/shortcodes/getresource.html --
{{ $r := .Page.Resources.Get (.Get 0)}}
Resource: {{ (.Get 0) }}|{{ with $r }}{{ .RelPermalink }}|{{ .Content }}|{{ else }}Not found.{{ end}}
-- layouts/_default/single.html --
{{ .Title }}|{{ .RelPermalink }}|{{ .Lang }}|{{ .Content }}|
`
	b := Test(t, files)

	// helpers.PrintFs(b.H.Fs.PublishDir, "", os.Stdout)
	b.AssertFileContent("public/nn/nnpages/my-bundle-nn/index.html", `
My Bundle NN
Resource: f1.txt|/nn/nnpages/my-bundle-nn/f1.txt|
Resource: f2.txt|/nn/nnpages/my-bundle-nn/f2.nn.txt|F2 nn.|
Resource: f2.nn.txt|/nn/nnpages/my-bundle-nn/f2.nn.txt|F2 nn.|
Resource: sub/f3.txt|/nn/nnpages/my-bundle-nn/sub/f3.txt|F1 en.|
`)

	b.AssertFileContent("public/enpages/my-bundle/f2.txt", "F2 en.")
	b.AssertFileContent("public/nn/nnpages/my-bundle-nn/f2.nn.txt", "F2 nn")

	b.AssertFileContent("public/enpages/my-bundle/index.html", `
Resource: f1.txt|/enpages/my-bundle/f1.txt|F1 en.|
Resource: f2.txt|/enpages/my-bundle/f2.txt|F2 en.|
`)
	b.AssertFileContent("public/enpages/my-bundle/f1.txt", "F1 en.")

	// Should be duplicated to the nn bundle.
	b.AssertFileContent("public/nn/nnpages/my-bundle-nn/f1.txt", "F1 en.")
}

// https://github.com/gohugoio/hugo/issues/5858
func TestBundledResourcesWhenMultipleOutputFormats(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.org"
disableKinds = ["taxonomy", "term"]
disableLiveReload = true
[outputs]
# This looks odd, but it triggers the behavior in #5858
# The total output formats list gets sorted, so CSS before HTML.
home = [ "CSS" ]
-- content/mybundle/index.md --
---
title: Page
---
-- content/mybundle/data.json --
MyData
-- layouts/_default/single.html --
{{ range .Resources }}
{{ .ResourceType }}|{{ .Title }}|
{{ end }}
`

	b := TestRunning(t, files)

	b.AssertFileContent("public/mybundle/data.json", "MyData")

	b.EditFileReplaceAll("content/mybundle/data.json", "MyData", "My changed data").Build()

	b.AssertFileContent("public/mybundle/data.json", "My changed data")
}

// See #11663
func TestPageBundlerPartialTranslations(t *testing.T) {
	t.Parallel()
	files := `
-- hugo.toml --
baseURL = "https://example.com"
disableKinds = ["taxonomy", "term"]
defaultContentLanguage = "en"
defaultContentLanguageInSubDir = true
[languages]
[languages.nn]
weight = 2
[languages.en]
weight = 1
-- content/section/mybundle/index.md --
---
title: "Mybundle"
---
-- content/section/mybundle/bundledpage.md --
---
title: "Bundled page en"
---
-- content/section/mybundle/bundledpage.nn.md --
---
title: "Bundled page nn"
---

-- layouts/_default/single.html --
Bundled page: {{ .RelPermalink}}|Len resources: {{ len .Resources }}|


`
	b := Test(t, files)

	b.AssertFileContent("public/en/section/mybundle/index.html",
		"Bundled page: /en/section/mybundle/|Len resources: 1|",
	)

	b.AssertFileExists("public/nn/section/mybundle/index.html", false)
}

// #6208
func TestBundleIndexInSubFolder(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.com"
-- layouts/_default/single.html --
{{ range .Resources }}
{{ .ResourceType }}|{{ .Title }}|
{{ end }}
-- content/bundle/index.md --
---
title: "bundle index"
---
-- content/bundle/p1.md --
---
title: "bundle p1"
---
-- content/bundle/sub/p2.md --
---
title: "bundle sub p2"
---
-- content/bundle/sub/index.md --
---
title: "bundle sub index"
---
-- content/bundle/sub/data.json --
data
`
	b := Test(t, files)

	b.AssertFileContent("public/bundle/index.html", `
        application|sub/data.json|
        page|bundle p1|
        page|bundle sub index|
        page|bundle sub p2|
`)
}

func TestPageBundlerHome(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "http://example.com/"
-- content/_index.md --
---
title: Home
---
![Alt text](image.jpg)
-- content/data.json --
DATA
-- layouts/index.html --
Title: {{ .Title }}|First Resource: {{ index .Resources 0 }}|Content: {{ .Content }}
-- layouts/_default/_markup/render-image.html --
Hook Len Page Resources {{ len .Page.Resources }}
`
	b := Test(t, files)

	b.AssertFileContent("public/index.html", `
Title: Home|First Resource: data.json|Content: <p>Hook Len Page Resources 1</p>
`)
}

func TestHTMLFilesIsue11999(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ["taxonomy", "term", "rss", "sitemap", "robotsTXT", "404"]
[permalinks]
posts = "/myposts/:slugorfilename"
-- content/posts/markdown-without-frontmatter.md --
-- content/posts/html-without-frontmatter.html --
<html>hello</html>
-- content/posts/html-with-frontmatter.html --
---
title: "HTML with frontmatter"
---
<html>hello</html>
-- content/posts/html-with-commented-out-frontmatter.html --
<!--
---
title: "HTML with commented out frontmatter"
---
-->
<html>hello</html>
-- content/posts/markdown-with-frontmatter.md --
---
title: "Markdown"
---
-- content/posts/mybundle/index.md --
---
title: My Bundle
---
-- content/posts/mybundle/data.txt --
Data.txt
-- content/posts/mybundle/html-in-bundle-without-frontmatter.html --
<html>hell</html>
-- content/posts/mybundle/html-in-bundle-with-frontmatter.html --
---
title: Hello
---
<html>hello</html>
-- content/posts/mybundle/html-in-bundle-with-commented-out-frontmatter.html --
<!--
---
title: "HTML with commented out frontmatter"
---
-->
<html>hello</html>
-- layouts/index.html --
{{ range site.RegularPages }}{{ .RelPermalink }}|{{ end }}$
-- layouts/_default/single.html --
{{ .Title }}|{{ .RelPermalink }}Resources: {{ range .Resources }}{{ .Name }}|{{ end }}$

`
	b := Test(t, files)

	b.AssertFileContent("public/index.html", "/myposts/html-with-commented-out-frontmatter/|/myposts/html-without-frontmatter/|/myposts/markdown-without-frontmatter/|/myposts/html-with-frontmatter/|/myposts/markdown-with-frontmatter/|/myposts/mybundle/|$")

	b.AssertFileContent("public/myposts/mybundle/index.html",
		"My Bundle|/myposts/mybundle/Resources: html-in-bundle-with-commented-out-frontmatter.html|html-in-bundle-without-frontmatter.html|html-in-bundle-with-frontmatter.html|data.txt|$")

	b.AssertPublishDir(`
index.html
myposts/html-with-commented-out-frontmatter
myposts/html-with-commented-out-frontmatter/index.html
myposts/html-with-frontmatter
myposts/html-with-frontmatter/index.html
myposts/html-without-frontmatter
myposts/html-without-frontmatter/index.html
myposts/markdown-with-frontmatter
myposts/markdown-with-frontmatter/index.html
myposts/markdown-without-frontmatter
myposts/markdown-without-frontmatter/index.html
myposts/mybundle/data.txt
myposts/mybundle/index.html
! myposts/mybundle/html-in-bundle-with-frontmatter.html
`)
}

func TestBundleDuplicatePagesAndResources(t *testing.T) {
	files := `
-- hugo.toml --
baseURL = "https://example.com"
disableKinds = ["taxonomy", "term"]
-- content/mysection/mybundle/index.md --
-- content/mysection/mybundle/index.html --
-- content/mysection/mybundle/p1.md --
-- content/mysection/mybundle/p1.html --
-- content/mysection/mybundle/foo/p1.html --
-- content/mysection/mybundle/data.txt --
Data txt.
-- content/mysection/mybundle/data.en.txt --
Data en txt.
-- content/mysection/mybundle/data.json --
Data JSON.
-- content/mysection/_index.md --
-- content/mysection/_index.html --
-- content/mysection/sectiondata.json --
Secion data JSON.
-- content/mysection/sectiondata.txt --
Section data TXT.
-- content/mysection/p2.md --
-- content/mysection/p2.html --
-- content/mysection/foo/p2.md --
-- layouts/_default/single.html --
Single:{{ .Title }}|{{ .Path }}|File LogicalName: {{ with .File }}{{ .LogicalName }}{{ end }}||{{ .RelPermalink }}|{{ .Kind }}|Resources: {{ range .Resources}}{{ .Name }}: {{ .Content }}|{{ end }}$
-- layouts/_default/list.html --
List: {{ .Title }}|{{ .Path }}|File LogicalName: {{ with .File }}{{ .LogicalName }}{{ end }}|{{ .RelPermalink }}|{{ .Kind }}|Resources: {{ range .Resources}}{{ .Name }}: {{ .Content }}|{{ end }}$
RegularPages: {{ range .RegularPages }}{{ .RelPermalink }}|File LogicalName: {{ with .File }}{{ .LogicalName }}|{{ end }}{{ end }}$
`

	b := Test(t, files)

	// Note that the sort order gives us the most specific data file for the en language (the data.en.json).
	b.AssertFileContent("public/mysection/mybundle/index.html", `Single:|/mysection/mybundle|File LogicalName: index.md||/mysection/mybundle/|page|Resources: data.en.txt: Data en txt.|data.json: Data JSON.|foo/p1.html: |p1.html: |p1.md: |$`)
	b.AssertFileContent("public/mysection/index.html",
		"List: |/mysection|File LogicalName: _index.md|/mysection/|section|Resources: sectiondata.json: Secion data JSON.|sectiondata.txt: Section data TXT.|$",
		"RegularPages: /mysection/foo/p2/|File LogicalName: p2.md|/mysection/mybundle/|File LogicalName: index.md|/mysection/p2/|File LogicalName: p2.md|$")
}

func TestBundleResourcesGetMatchOriginalName(t *testing.T) {
	files := `
-- hugo.toml --
baseURL = "https://example.com"
-- content/mybundle/index.md --
-- content/mybundle/f1.en.txt --
F1.
-- layouts/_default/single.html --
GetMatch: {{ with .Resources.GetMatch "f1.en.*" }}{{ .Name }}: {{ .Content }}|{{ end }}
Match: {{ range .Resources.Match "f1.En.*" }}{{ .Name }}: {{ .Content }}|{{ end }}
`

	b := Test(t, files)

	b.AssertFileContent("public/mybundle/index.html", "GetMatch: f1.en.txt: F1.|", "Match: f1.en.txt: F1.|")
}

func TestBundleResourcesWhenLanguageVariantIsDraft(t *testing.T) {
	files := `
-- hugo.toml --
baseURL = "https://example.com"
defaultContentLanguage = "en"
[languages]
[languages.en]
weight = 1
[languages.nn]
weight = 2
-- content/mybundle/index.en.md --
-- content/mybundle/index.nn.md --
---
draft: true
---
-- content/mybundle/f1.en.txt --
F1.
-- layouts/_default/single.html --
GetMatch: {{ with .Resources.GetMatch "f1.*" }}{{ .Name }}: {{ .Content }}|{{ end }}$
`

	b := Test(t, files)

	b.AssertFileContent("public/mybundle/index.html", "GetMatch: f1.en.txt: F1.|")
}

func TestBundleBranchIssue12320(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['rss','sitemap','taxonomy','term']
defaultContentLanguage = 'en'
defaultContentLanguageInSubdir = true
[languages.en]
baseURL = "https://en.example.org/"
contentDir = "content/en"
[languages.fr]
baseURL = "https://fr.example.org/"
contentDir = "content/fr"
-- content/en/s1/p1.md --
---
title: p1
---
-- content/en/s1/p1.txt --
---
p1.txt
---
-- layouts/_default/single.html --
{{ .Title }}|
-- layouts/_default/list.html --
{{ .Title }}|
`

	b := Test(t, files)

	b.AssertFileExists("public/en/s1/index.html", true)
	b.AssertFileExists("public/en/s1/p1/index.html", true)
	b.AssertFileExists("public/en/s1/p1.txt", true)

	b.AssertFileExists("public/fr/s1/index.html", false)
	b.AssertFileExists("public/fr/s1/p1/index.html", false)
	b.AssertFileExists("public/fr/s1/p1.txt", false) // failing test
}
