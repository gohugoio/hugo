// Copyright 2024 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless requiredF by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package debug_test

import (
	"testing"

	"github.com/bep/logg"
	"github.com/gohugoio/hugo/hugolib"
)

func TestTimer(t *testing.T) {
	files := `
-- hugo.toml --
baseURL = "https://example.org/"
disableKinds = ["taxonomy", "term"]
-- layouts/index.html --
{{ range seq 5 }}
{{ $t := debug.Timer "foo" }}
{{ seq 1 1000 }}
{{ $t.Stop }}
{{ end }}

`
	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			LogLevel:    logg.LevelInfo,
		},
	).Build()

	b.AssertLogContains("timer:  name foo count 5 duration")
}

func TestDebugDumpPage(t *testing.T) {
	files := `
-- hugo.toml --
baseURL = "https://example.org/"
disableLiveReload = true
[taxonomies]
tag = "tags"
-- content/_index.md --
---
title: "The Index"
date: 2012-03-15
---
-- content/p1.md --
---
title: "The First"
tags: ["a", "b"]
---
-- layouts/_default/list.html --
Dump: {{ debug.Dump . | safeHTML }}
Dump Site: {{ debug.Dump site }}
Dum site.Taxonomies: {{ debug.Dump site.Taxonomies | safeHTML }}
-- layouts/_default/single.html --
Dump: {{ debug.Dump . | safeHTML }}


`
	b := hugolib.TestRunning(t, files)
	b.AssertFileContent("public/index.html", "Dump: {\n  \"Date\": \"2012-03-15T00:00:00Z\"")
}

func TestDebugList(t *testing.T) {
	files := `
-- hugo.toml --
baseURL = "https://example.org/"
disableLiveReload = true
-- content/_index.md --
---
title: "The Index"
---
-- layouts/_default/list.html --
{{ $dict := dict "name" "Hugo" "version" "0.120" "type" "SSG" }}
Map Keys: {{ debug.List $dict }}

{{ $page := . }}
Page Fields: {{ debug.List $page }}

Nil: {{ debug.List nil }}
String: {{ debug.List "hello" }}
Number: {{ debug.List 42 }}
Slice: {{ debug.List (slice 1 2 3) }}


`
	b := hugolib.TestRunning(t, files)

	// Test map keys
	b.AssertFileContent("public/index.html", "Map Keys: [name type version]")

	// Test that page struct returns field names and methods (should include common page fields)
	b.AssertFileContent("public/index.html", "Page Fields:")
	b.AssertFileContent("public/index.html", "Title")
	b.AssertFileContent("public/index.html", "Content")

	// Test edge cases
	b.AssertFileContent("public/index.html", "Nil: []")
	b.AssertFileContent("public/index.html", "String: []")
	b.AssertFileContent("public/index.html", "Number: []")
	b.AssertFileContent("public/index.html", "Slice: []")
}
