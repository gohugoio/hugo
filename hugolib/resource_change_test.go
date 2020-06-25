// Copyright 2020 The Hugo Authors. All rights reserved.
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
	"strings"
	"testing"
)

func TestResourceEditMetadata(t *testing.T) {
	b := newTestSitesBuilder(t).Running()

	content := `+++
title = "My Bundle With TOML Meta"

[[resources]]
src = "**.toml"
title = "My TOML :counter"
+++

Content.
`

	b.WithContent(
		"bundle/index.md", content,
		"bundle/my1.toml", `a = 1`,
		"bundle/my2.toml", `a = 2`)

	b.WithTemplatesAdded("index.html", `
{{ $bundle := site.GetPage "bundle" }}
{{ $toml := $bundle.Resources.GetMatch "*.toml"  }}
TOML: {{ $toml.Title }}
	
`)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/index.html", "TOML: My TOML 1")

	b.EditFiles("content/bundle/index.md", strings.ReplaceAll(content, "My TOML", "My Changed TOML 1"))

	b.Build(BuildCfg{})

	b.AssertFileContent("public/index.html", "TOML: My Changed TOML")

}

func TestResourceCacheSimpleTest(t *testing.T) {
	conf := `
baseURL = "https://example.org"

defaultContentLanguage = "en"

[module]
[[module.mounts]]
source = "content/cen"
target = "content"
lang="en"
[[module.mounts]]
source = "content/cno"
target = "content"
lang="no"
[[module.mounts]]
source = "assets"
target = "assets"
[[module.mounts]]
source = "assets_common"
target = "assets/aen"
[[module.mounts]]
source = "assets_common"
target = "assets/ano"

[languages]
[languages.en]
weight = 1

[languages.no]
weight = 2

`
	b := newTestSitesBuilder(t).WithConfigFile("toml", conf).Running()

	b.WithSourceFile(
		"content/cen/bundle/index.md", "---\ntitle: En Bundle\n---",
		"content/cen/bundle/data1.json", `{ "data1": "en" }`,
		"content/cen/bundle/data2.json", `{ "data2": "en" }`,
		"content/cno/bundle/index.md", "---\ntitle: No Bundle\n---",
		"content/cno/bundle/data1.json", `{ "data1": "no" }`,
		"content/cno/bundle/data3.json", `{ "data3": "no" }`,
	)

	b.WithSourceFile("assets_common/data/common.json", `{
    "Hugo": "Rocks!",
 }`)

	b.WithSourceFile("assets/data/mydata.json", `{
    "a": 32,
 }`)

	b.WithTemplatesAdded("index.html", `
{{ $data := resources.Get "data/mydata.json" }}
{{ template "print-resource" ( dict "title" "data" "r" $data ) }}
{{ $dataMinified := $data | minify }}
{{ template "print-resource" ( dict "title" "data-minified" "r" $dataMinified ) }}
{{ $dataUnmarshaled := $dataMinified | transform.Unmarshal }}
Data Unmarshaled: {{ $dataUnmarshaled }}
{{ $bundle := site.GetPage "bundle" }}
{{ range (seq 3) }}
{{ $i := . }}
{{ with $bundle.Resources.GetMatch (printf "data%d.json" . ) }}
{{ $minified := . | minify }}
{{ template "print-resource" ( dict "title" (printf "bundle data %d" $i) "r" . ) }}
{{ template "print-resource" ( dict "title" (printf "bundle data %d min" $i) "r" $minified ) }}
{{ end }}
{{ end }}
{{ $common1 := resources.Get "aen/data/common.json" }}
{{ $common2 := resources.Get "ano/data/common.json" }}
{{ template "print-resource" ( dict "title" "common1" "r" $common1 ) }}
{{ template "print-resource" ( dict "title" "common2" "r" $common2 ) }}
{{ define "print-resource" }}{{ .title }}|{{ .r.RelPermalink }}|{{ .r.Key }}|{{ .r.Content | safeHTML }}|{{ end }}



`)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/index.html", `
data-minified|/data/mydata.min.json|data/mydata_d3f53f09220d597dac26fe7840c31fc9.json|{"a":32}|
Data Unmarshaled: map[a:32]
bundle data 1|/bundle/data1.json|bundle/content/en/data1.json|{ "data1": "en" }|
bundle data 1 min|/bundle/data1.min.json|bundle/content/en/data1_d3f53f09220d597dac26fe7840c31fc9.json|{"data1":"en"}|
bundle data 3|/bundle/data3.json|bundle/content/en/data3.json|{ "data3": "no" }|
bundle data 3 min|/bundle/data3.min.json|bundle/content/en/data3_d3f53f09220d597dac26fe7840c31fc9.json|{"data3":"no"}|
common1|/aen/data/common.json|aen/data/common.json|
common2|/ano/data/common.json|ano/data/common.json|
`)

	b.AssertFileContent("public/no/index.html", `
data-minified|/data/mydata.min.json|data/mydata_d3f53f09220d597dac26fe7840c31fc9.json|{"a":32}|
bundle data 1|/no/bundle/data1.json|bundle/content/no/data1.json|{ "data1": "no" }|
bundle data 2|/no/bundle/data2.json|bundle/content/no/data2.json|{ "data2": "en" }|
 bundle data 3|/no/bundle/data3.json|bundle/content/no/data3.json|{ "data3": "no" }|
bundle data 3 min|/no/bundle/data3.min.json|bundle/content/no/data3_d3f53f09220d597dac26fe7840c31fc9.json|{"data3":"no"}|
common1|/aen/data/common.json|aen/data/common.json
`)

	b.EditFiles("assets/data/mydata.json", `{ "a": 42 }`)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/index.html", `
data|/data/mydata.json|data/mydata.json|{ "a": 42 }|
data-minified|/data/mydata.min.json|data/mydata_d3f53f09220d597dac26fe7840c31fc9.json|{"a":42}|
Data Unmarshaled: map[a:42]
`)

}

func TestResourceCacheMultihost(t *testing.T) {
	toLang := func(format, lang string) string {
		return fmt.Sprintf(format, lang)
	}

	addContent := func(b *sitesBuilder, contentPath func(path, lang string) string) {
		b.WithNoContentAdded()
		for _, lang := range []string{"en", "fr"} {
			b.WithSourceFile(
				contentPath("b1/index.md", lang), toLang("---\ntitle: Bundle 1 %s\n---", lang),
				contentPath("b1/styles/style11.css", lang), toLang(".%s1: { color: blue };", lang),
				contentPath("b1/styles/style12.css", lang), toLang(".%s2: { color: red };", lang),
			)
			b.WithSourceFile(
				contentPath("b2/index.md", lang), toLang("---\ntitle: Bundle 2 %s\n---", lang),
				contentPath("b2/styles/style21.css", lang), toLang(".%s21: { color: green };", lang),
				contentPath("b2/styles/style22.css", lang), toLang(".%s22: { color: orange };", lang),
			)
		}
	}

	addTemplates := func(b *sitesBuilder) {
		b.WithTemplates("_default/single.html", `
{{ template "print-page" (dict "page" . "title" "Self") }}
{{ $other := site.Sites.First.GetPage "b1" }}
{{ template "print-page" (dict "page" $other "title" "Other") }}
	

{{ define "print-page" }}
{{ $p := .page }}
{{ $title := .title }}
{{ $styles := $p.Resources.Match "**/style1*.css" }}
{{ if $styles }}
{{ $firststyle := index $styles 0 }}
{{ $mystyles := $styles | resources.Concat "mystyles.css" }}
{{ $title }} Mystyles First CSS: {{ $firststyle.RelPermalink }}|Key: {{ $firststyle.Key }}|{{ $firststyle.Content }}|
{{ $title }} Mystyles CSS: {{ $mystyles.RelPermalink }}|Key: {{ $mystyles.Key }}|{{ $mystyles.Content }}|
{{ end }}
{{ $title }} Bundle: {{ $p.Permalink }}
{{ $style := $p.Resources.GetMatch "**.css" }}
{{ $title }} CSS: {{ $style.RelPermalink }}|Key: {{ $style.Key }}|{{ $style.Content }}|
{{ $minified := $style | minify }}
{{ $title }} Minified CSS: {{ $minified.RelPermalink }}|Key: {{ $minified.Key }}|{{ $minified.Content }}
{{ end }}

`)
	}

	assertContent := func(b *sitesBuilder) {

		otherAssert := `Other Mystyles First CSS: /b1/styles/style11.css|Key: b1/content/en/styles/style11.css|.en1: { color: blue };|
Other Bundle: https://example.com/b1/
Other CSS: /b1/styles/style11.css|Key: b1/content/en/styles/style11.css|.en1: { color: blue };
Other Minified CSS: /b1/styles/style11.min.css|Key: b1/content/en/styles/style11_d3f53f09220d597dac26fe7840c31fc9.css|.en1:{color:blue}`

		b.AssertFileContent("public/fr/b1/index.html", `
Self Mystyles CSS: /mystyles.css|Key: _root/mystyles.css|.en1: { color: blue };.en2: { color: red };
Self Bundle: https://example.fr/b1/
Self CSS: /b1/styles/style11.css|Key: b1/content/fr/styles/style11.css|.fr1: { color: blue };
Self Minified CSS: /b1/styles/style11.min.css|Key: b1/content/fr/styles/style11_d3f53f09220d597dac26fe7840c31fc9.css|.fr1:{color:blue}
`,
			otherAssert)
		b.AssertFileContent("public/en/b1/index.html", `
Self Mystyles First CSS: /b1/styles/style11.css|Key: b1/content/en/styles/style11.css|.en1: { color: blue };|
Self Mystyles CSS: /mystyles.css|Key: _root/mystyles.css|.en1: { color: blue };.en2: { color: red };|
Self Bundle: https://example.com/b1/
Self CSS: /b1/styles/style11.css|Key: b1/content/en/styles/style11.css|.en1: { color: blue };|
Self Minified CSS: /b1/styles/style11.min.css|Key: b1/content/en/styles/style11_d3f53f09220d597dac26fe7840c31fc9.css|.en1:{color:blue}
`, otherAssert)

		b.AssertFileContent("public/fr/b2/index.html", `
Self Bundle: https://example.fr/b2/
Self CSS: /b2/styles/style21.css|Key: b2/content/fr/styles/style21.css|.fr21: { color: green };|
Self Minified CSS: /b2/styles/style21.min.css|Key: b2/content/fr/styles/style21_d3f53f09220d597dac26fe7840c31fc9.css|.fr21:{color:green}
`, otherAssert)

		b.AssertFileContent("public/en/b2/index.html", `
Self Bundle: https://example.com/b2/
Self CSS: /b2/styles/style21.css|Key: b2/content/en/styles/style21.css|.en21: { color: green };|
Self Minified CSS: /b2/styles/style21.min.css|Key: b2/content/en/styles/style21_d3f53f09220d597dac26fe7840c31fc9.css|.en21:{color:green}
`, otherAssert)

	}

	t.Run("Default content", func(t *testing.T) {

		var configTemplate = `
paginate = 1
defaultContentLanguage = "fr"
defaultContentLanguageInSubdir = false
contentDir = "content"

[Languages]
[Languages.en]
baseURL = "https://example.com/"
weight = 10
languageName = "English"

[Languages.fr]
baseURL = "https://example.fr"
weight = 20
languageName = "Français"

`

		b := newTestSitesBuilder(t).WithConfigFile("toml", configTemplate)
		fmt.Println(b.workingDir)
		addContent(b, func(path, lang string) string {
			path = strings.Replace(path, ".", "."+lang+".", 1)
			path = "content/" + path
			return path
		})
		addTemplates(b)
		b.Build(BuildCfg{})
		assertContent(b)

	})

	t.Run("Content dir per language", func(t *testing.T) {

		var configTemplate = `
paginate = 1
defaultContentLanguage = "fr"
defaultContentLanguageInSubdir = false

[Languages]
[Languages.en]
contentDir = "content_en"
baseURL = "https://example.com/"
weight = 10
languageName = "English"

[Languages.fr]
contentDir = "content_fr"
baseURL = "https://example.fr"
weight = 20
languageName = "Français"

`

		b := newTestSitesBuilder(t).WithConfigFile("toml", configTemplate)
		addContent(b, func(path, lang string) string {
			return "content_" + lang + "/" + path
		})
		addTemplates(b)
		b.Build(BuildCfg{})
		assertContent(b)

	})

}
