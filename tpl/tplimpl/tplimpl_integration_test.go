// Copyright 2025 The Hugo Authors. All rights reserved.
//
// Portions Copyright The Go Authors.

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

package tplimpl_test

import (
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/hugolib"
)

// Verify that the new keywords in Go 1.18 is available.
func TestGo18Constructs(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = 'http://example.com/'
disableKinds = ["section", "home", "rss", "taxonomy",  "term", "rss"]
-- content/p1.md --
---
title: "P1"
---
-- layouts/_partials/counter.html --
{{ if .Scratch.Get "counter" }}{{ .Scratch.Add "counter" 1 }}{{ else }}{{ .Scratch.Set "counter" 1 }}{{ end }}{{ return true }}
-- layouts/single.html --
continue:{{ range seq 5 }}{{ if eq . 2 }}{{continue}}{{ end }}{{ . }}{{ end }}:END:
break:{{ range seq 5 }}{{ if eq . 2 }}{{break}}{{ end }}{{ . }}{{ end }}:END:
continue2:{{ range seq 5 }}{{ if eq . 2 }}{{ continue }}{{ end }}{{ . }}{{ end }}:END:
break2:{{ range seq 5 }}{{ if eq . 2 }}{{ break }}{{ end }}{{ . }}{{ end }}:END:

counter1: {{ partial "counter.html" . }}/{{ .Scratch.Get "counter" }}
and1: {{ if (and false (partial "counter.html" .)) }}true{{ else }}false{{ end }}
or1: {{ if (or true (partial "counter.html" .)) }}true{{ else }}false{{ end }}
and2: {{ if (and true (partial "counter.html" .)) }}true{{ else }}false{{ end }}
or2: {{ if (or false (partial "counter.html" .)) }}true{{ else }}false{{ end }}


counter2: {{ .Scratch.Get "counter" }}


	`

	b := hugolib.Test(t, files, hugolib.TestOptOsFs())

	b.AssertFileContent("public/p1/index.html", `
continue:1345:END:
break:1:END:
continue2:1345:END:
break2:1:END:
counter1: true/1
and1: false
or1: true
and2: true
or2: true
counter2: 3
`)
}

func TestGo23ElseWith(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
title = "Hugo"
-- layouts/home.html --
{{ with false }}{{ else with .Site }}{{ .Title }}{{ end }}|
`
	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html", "Hugo|")
}

// Issue 10495
func TestCommentsBeforeBlockDefinition(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = 'http://example.com/'
-- content/s1/p1.md --
---
title: "S1P1"
---
-- content/s2/p1.md --
---
title: "S2P1"
---
-- content/s3/p1.md --
---
title: "S3P1"
---
-- layouts/baseof.html --
{{ block "main" . }}{{ end }}
-- layouts/s1/single.html --
{{/* foo */}}
{{ define "main" }}{{ .Title }}{{ end }}
-- layouts/s2/single.html --
{{- /* foo */}}
{{ define "main" }}{{ .Title }}{{ end }}
-- layouts/s3/single.html --
{{- /* foo */ -}}
{{ define "main" }}{{ .Title }}{{ end }}
	`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/s1/p1/index.html", `S1P1`)
	b.AssertFileContent("public/s2/p1/index.html", `S2P1`)
	b.AssertFileContent("public/s3/p1/index.html", `S3P1`)
}

func TestGoTemplateBugs(t *testing.T) {
	t.Run("Issue 11112", func(t *testing.T) {
		t.Parallel()

		files := `
-- hugo.toml --
-- layouts/home.html --
{{ $m := dict "key" "value" }}
{{ $k := "" }}
{{ $v := "" }}
{{ range $k, $v = $m }}
{{ $k }} = {{ $v }}
{{ end }}
	`

		b := hugolib.Test(t, files)

		b.AssertFileContent("public/index.html", `key = value`)
	})
}

func TestSecurityAllowActionJSTmpl(t *testing.T) {
	filesTemplate := `
-- hugo.toml --
SECURITYCONFIG
-- layouts/home.html --
<script>
var a = §§{{.Title }}§§;
</script>
	`

	files := strings.ReplaceAll(filesTemplate, "SECURITYCONFIG", "")

	b, err := hugolib.TestE(t, files)

	// This used to fail, but not in >= Hugo 0.121.0.
	b.Assert(err, qt.IsNil)
}

func TestSitemap(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['home','rss','section','taxonomy','term']
[sitemap]
disable = true
-- content/p1.md --
---
title: p1
sitemap:
  p1_disable: foo
---
-- content/p2.md --
---
title: p2

---
-- layouts/single.html --
{{ .Title }}
`

	// Test A: Exclude all pages via project config.
	b := hugolib.Test(t, files)
	b.AssertFileContentExact("public/sitemap.xml",
		"<?xml version=\"1.0\" encoding=\"utf-8\" standalone=\"yes\"?>\n<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\"\n  xmlns:xhtml=\"http://www.w3.org/1999/xhtml\">\n  \n</urlset>\n",
	)

	// Test B: Include all pages via project config.
	files_b := strings.ReplaceAll(files, "disable = true", "disable = false")
	b = hugolib.Test(t, files_b)
	b.AssertFileContentExact("public/sitemap.xml",
		"<?xml version=\"1.0\" encoding=\"utf-8\" standalone=\"yes\"?>\n<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\"\n  xmlns:xhtml=\"http://www.w3.org/1999/xhtml\">\n  <url>\n    <loc>/p1/</loc>\n  </url><url>\n    <loc>/p2/</loc>\n  </url>\n</urlset>\n",
	)

	// Test C: Exclude all pages via project config, but include p1 via front matter.
	files_c := strings.ReplaceAll(files, "p1_disable: foo", "disable: false")
	b = hugolib.Test(t, files_c)
	b.AssertFileContentExact("public/sitemap.xml",
		"<?xml version=\"1.0\" encoding=\"utf-8\" standalone=\"yes\"?>\n<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\"\n  xmlns:xhtml=\"http://www.w3.org/1999/xhtml\">\n  <url>\n    <loc>/p1/</loc>\n  </url>\n</urlset>\n",
	)

	// Test D:  Include all pages via project config, but exclude p1 via front matter.
	files_d := strings.ReplaceAll(files_b, "p1_disable: foo", "disable: true")
	b = hugolib.Test(t, files_d)
	b.AssertFileContentExact("public/sitemap.xml",
		"<?xml version=\"1.0\" encoding=\"utf-8\" standalone=\"yes\"?>\n<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\"\n  xmlns:xhtml=\"http://www.w3.org/1999/xhtml\">\n  <url>\n    <loc>/p2/</loc>\n  </url>\n</urlset>\n",
	)
}

// Issue 12963
func TestEditBaseofParseAfterExecute(t *testing.T) {
	files := `
-- hugo.toml --
baseURL = "https://example.com"
disableLiveReload = true
disableKinds = ["taxonomy", "term", "rss", "404", "sitemap"]
[internal]
fastRenderMode = true
-- layouts/baseof.html --
Baseof!
{{ block "main" . }}default{{ end }}
{{ with (templates.Defer (dict "key" "global")) }}
Now. {{ now }}
{{ end }}
-- layouts/single.html --
{{ define "main" }}
Single.
{{ end }}
-- layouts/list.html --
{{ define "main" }}
List.
{{ .Content }}
{{ range .Pages }}{{ .Title }}{{ end }}|
{{ end }}
-- content/mybundle1/index.md --
---
title: "My Bundle 1"
---
-- content/mybundle2/index.md --
---
title: "My Bundle 2"
---
-- content/_index.md --
---
title: "Home"
---
Home!
`

	b := hugolib.TestRunning(t, files)
	b.AssertFileContent("public/index.html", "Home!")
	b.EditFileReplaceAll("layouts/baseof.html", "baseof", "Baseof!").Build()
	b.BuildPartial("/")
	b.AssertFileContent("public/index.html", "Baseof!")
	b.BuildPartial("/mybundle1/")
	b.AssertFileContent("public/mybundle1/index.html", "Baseof!")
}
