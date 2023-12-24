// Copyright 2024 The Hugo Authors. All rights reserved.
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

func TestSitesAndLanguageOrder(t *testing.T) {
	files := `
-- hugo.toml --
defaultContentLanguage = "fr"
defaultContentLanguageInSubdir = true
[languages]
[languages.en]
weight = 1
[languages.fr]
weight = 2
[languages.de]
weight = 3
-- layouts/index.html --
{{ $bundle := site.GetPage "bundle" }}
Bundle all translations: {{ range $bundle.AllTranslations }}{{ .Lang }}|{{ end }}$
Bundle translations: {{ range $bundle.Translations }}{{ .Lang }}|{{ end }}$
Site languages: {{ range site.Languages }}{{ .Lang }}|{{ end }}$
Sites: {{ range site.Sites }}{{ .Language.Lang }}|{{ end }}$
-- content/bundle/index.fr.md --
---
title: "Bundle Fr"
---
-- content/bundle/index.en.md --
---
title: "Bundle En"
---
-- content/bundle/index.de.md --
---
title: "Bundle De"
---
	
	`
	b := Test(t, files)

	b.AssertFileContent("public/en/index.html",
		"Bundle all translations: en|fr|de|$",
		"Bundle translations: fr|de|$",
		"Site languages: en|fr|de|$",
		"Sites: fr|en|de|$",
	)
}
