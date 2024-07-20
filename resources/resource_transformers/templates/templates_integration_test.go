// Copyright 2021 The Hugo Authors. All rights reserved.
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

package templates_test

import (
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

func TestExecuteAsTemplateMultipleLanguages(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
baseURL = "http://example.com/blog"
defaultContentLanguage = "fr"
defaultContentLanguageInSubdir = true
[Languages]
[Languages.en]
weight = 10
title = "In English"
languageName = "English"
[Languages.fr]
weight = 20
title = "Le Français"
languageName = "Français"
-- i18n/en.toml --
[hello]
other = "Hello"
-- i18n/fr.toml --
[hello]
other = "Bonjour"
-- layouts/index.fr.html --
Lang: {{ site.Language.Lang }}
{{ $templ := "{{T \"hello\"}}" | resources.FromString "f1.html" }}
{{ $helloResource := $templ | resources.ExecuteAsTemplate (print "f%s.html" .Lang) . }}
Hello1: {{T "hello"}}
Hello2: {{ $helloResource.Content }}
LangURL: {{ relLangURL "foo" }}
-- layouts/index.html --
Lang: {{ site.Language.Lang }}
{{ $templ := "{{T \"hello\"}}" | resources.FromString "f1.html" }}
{{ $helloResource := $templ | resources.ExecuteAsTemplate (print "f%s.html" .Lang) . }}
Hello1: {{T "hello"}}
Hello2: {{ $helloResource.Content }}
LangURL: {{ relLangURL "foo" }}

	`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/en/index.html", `
		Hello1: Hello
		Hello2: Hello
		`)

	b.AssertFileContent("public/fr/index.html", `
		Hello1: Bonjour
		Hello2: Bonjour
		`)
}
