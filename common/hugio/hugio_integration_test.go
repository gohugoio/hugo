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

package hugio_test

import (
	"testing"

	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/hugolib"
)

func TestTrimBOM(t *testing.T) {
	files := `
-- hugo.toml --
baseURL = "http://example.org/"
-- i18n/en.yaml --
hello: "Hello"
-- data/mydata.yaml --
foo: bar
-- assets/mycss.css --
body { background: #fff; }
-- assets/myjs.js --
console.log("Hello, world!");
-- assets/myyaml.yaml --
foo: bar
-- content/page.md --
---
title: "Page with BOM"
---
-- content/mysection/_content.gotmpl --
{{ .AddPage (dict "path" "p2" "title" "p2") }}
-- layouts/all.html --
Title: {{ .Title }}|
Unmarshal resource to map: {{ resources.Get "myyaml.yaml" | transform.Unmarshal }}|
Unmarshal resource.Content to map: {{ (resources.Get "myyaml.yaml").Content | transform.Unmarshal }}|
{{ $myyamlFronOsReadFile := os.ReadFile "assets/myyaml.yaml" }}
Unmarshal os.ReadFile result to map: {{ $myyamlFronOsReadFile | transform.Unmarshal }}|
Data: {{ .Site.Data.mydata.foo }}|
I18n: {{ i18n "hello" }}|
CSS1: {{ resources.Get "mycss.css" | minify | resources.Publish }}|
CSS2: {{ resources.Get "mycss.css" | css.Build | resources.Publish }}|
JS: {{ resources.Get "myjs.js" | js.Build | resources.Publish }}|

`
	b := hugolib.Test(t, files, hugolib.TestOptWithConfig(func(c *hugolib.IntegrationTestConfig) {
		c.FileContentPrefix = hugio.BOM
		c.NeedsOsFS = true
	}))

	// No BOM in the published text files.
	b.AssertStringPublished("! " + hugio.BOM)

	b.AssertFileContent(
		"public/page/index.html",
		"Title: Page with BOM|",
		"Unmarshal resource to map: map[foo:bar]|",
		"Unmarshal resource.Content to map: map[foo:bar]|",
		"Unmarshal os.ReadFile result to map: map[foo:bar]|",
		"I18n: Hello|",
	)
	b.AssertFileContent("public/mysection/p2/index.html", "Title: p2|")
}
