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

func TestImageResizeMultilingual(t *testing.T) {
	b := newTestSitesBuilder(t).WithConfigFile("toml", `
baseURL="https://example.org"
defaultContentLanguage = "en"

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

`)

	pageContent := `---
title: "Page"
---
`

	b.WithContent("bundle/index.md", pageContent)
	b.WithContent("bundle/index.nn.md", pageContent)
	b.WithContent("bundle/index.fr.md", pageContent)
	b.WithSunset("content/bundle/sunset.jpg")
	b.WithSunset("assets/images/sunset.jpg")
	b.WithTemplates("index.html", `
{{ with (.Site.GetPage "bundle" ) }}
{{ $sunset := .Resources.GetMatch "sunset*" }}
{{ if $sunset }}
{{ $resized := $sunset.Resize "200x200" }}
SUNSET FOR: {{ $.Site.Language.Lang }}: {{ $resized.RelPermalink }}/{{ $resized.Width }}/Lat: {{ $resized.Exif.Lat }}
{{ end }}
{{ else }}
No bundle for {{ $.Site.Language.Lang }}
{{ end }}

{{ $sunset2 := resources.Get "images/sunset.jpg" }}
{{ $resized2 := $sunset2.Resize "123x234" }}
SUNSET2: {{ $resized2.RelPermalink }}/{{ $resized2.Width }}/Lat: {{ $resized2.Exif.Lat }}


`)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/index.html", "SUNSET FOR: en: /bundle/sunset_hu13235715490294913361.jpg/200/Lat: 36.59744166666667")
	b.AssertFileContent("public/fr/index.html", "SUNSET FOR: fr: /bundle/sunset_hu13235715490294913361.jpg/200/Lat: 36.59744166666667")
	b.AssertFileContent("public/index.html", " SUNSET2: /images/sunset_hu1573057890424052540.jpg/123/Lat: 36.59744166666667")
	b.AssertFileContent("public/nn/index.html", " SUNSET2: /images/sunset_hu1573057890424052540.jpg/123/Lat: 36.59744166666667")

	b.AssertImage(200, 200, "public/bundle/sunset_hu13235715490294913361.jpg")

	// Check the file cache
	b.AssertImage(200, 200, "resources/_gen/images/bundle/sunset_hu13235715490294913361.jpg")

	b.AssertFileContent("resources/_gen/images/bundle/sunset_17710516992648092201.json",
		"FocalLengthIn35mmFormat|uint16", "PENTAX")

	b.AssertFileContent("resources/_gen/images/images/sunset_17710516992648092201.json",
		"FocalLengthIn35mmFormat|uint16", "PENTAX")

	b.AssertNoDuplicateWrites()
}
