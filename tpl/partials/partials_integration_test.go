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

package partials_test

import (
	"bytes"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/gohugoio/hugo/htesting/hqt"
	"github.com/gohugoio/hugo/hugolib"
)

func TestInclude(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
baseURL = 'http://example.com/'
-- layouts/index.html --
partial: {{ partials.Include "foo.html" . }}
-- layouts/partials/foo.html --
foo
  `

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html", `
partial: foo
`)
}

func TestIncludeCached(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
baseURL = 'http://example.com/'
-- layouts/index.html --
partialCached: {{ partials.IncludeCached "foo.html" . }}
partialCached: {{ partials.IncludeCached "foo.html" . }}
-- layouts/partials/foo.html --
foo
  `

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html", `
partialCached: foo
partialCached: foo
`)
}

// Issue 9519
func TestIncludeCachedRecursion(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
baseURL = 'http://example.com/'
-- layouts/index.html --
{{ partials.IncludeCached "p1.html" . }}
-- layouts/partials/p1.html --
{{ partials.IncludeCached "p2.html" . }}
-- layouts/partials/p2.html --
P2

  `

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html", `
P2
`)
}

// Issue #588
func TestIncludeCachedRecursionShortcode(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
baseURL = 'http://example.com/'
-- content/_index.md --
---
title: "Index"
---
{{< short >}}
-- layouts/index.html --
{{ partials.IncludeCached "p1.html" . }}
-- layouts/partials/p1.html --
{{ .Content }}
{{ partials.IncludeCached "p2.html" . }}
-- layouts/partials/p2.html --
-- layouts/shortcodes/short.html --
SHORT
{{ partials.IncludeCached "p2.html" . }}
P2

  `

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html", `
SHORT
P2
`)
}

func TestIncludeCacheHints(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
baseURL = 'http://example.com/'
templateMetrics=true
templateMetricsHints=true
disableKinds = ["page", "section", "taxonomy", "term", "sitemap"]
[outputs]
home = ["HTML"]
-- layouts/index.html --
{{ partials.IncludeCached "static1.html" . }}
{{ partials.IncludeCached "static1.html" . }}
{{ partials.Include "static2.html" . }}

D1I: {{ partials.Include "dynamic1.html" . }}
D1C: {{ partials.IncludeCached "dynamic1.html" . }}
D1C: {{ partials.IncludeCached "dynamic1.html" . }}
D1C: {{ partials.IncludeCached "dynamic1.html" . }}
H1I: {{ partials.Include "halfdynamic1.html" . }}
H1C: {{ partials.IncludeCached "halfdynamic1.html" . }}
H1C: {{ partials.IncludeCached "halfdynamic1.html" . }}

-- layouts/partials/static1.html --
P1
-- layouts/partials/static2.html --
P2
-- layouts/partials/dynamic1.html --
{{ math.Counter }}
-- layouts/partials/halfdynamic1.html --
D1
{{ math.Counter }}


  `

	b := hugolib.Test(t, files)

	// fmt.Println(b.FileContent("public/index.html"))

	var buf bytes.Buffer
	b.H.Metrics.WriteMetrics(&buf)

	got := buf.String()

	// Get rid of all the durations, they are never the same.
	durationRe := regexp.MustCompile(`\b[\.\d]*(ms|µs|s)\b`)

	normalize := func(s string) string {
		s = durationRe.ReplaceAllString(s, "")
		linesIn := strings.Split(s, "\n")[3:]
		var lines []string
		for _, l := range linesIn {
			l = strings.TrimSpace(l)
			if l == "" {
				continue
			}
			lines = append(lines, l)
		}

		sort.Strings(lines)

		return strings.Join(lines, "\n")
	}

	got = normalize(got)

	expect := `
	0        0       0      1  index.html
	100        0       0      1  partials/static2.html
	100       50       1      2  partials/static1.html
	25       50       2      4  partials/dynamic1.html
	66       33       1      3  partials/halfdynamic1.html
	`

	b.Assert(got, hqt.IsSameString, expect)
}

// gobench --package ./tpl/partials
func BenchmarkIncludeCached(b *testing.B) {
	files := `
-- config.toml --
baseURL = 'http://example.com/'
-- layouts/index.html --
-- layouts/_default/single.html --
{{ partialCached "heavy.html" "foo" }}
{{ partialCached "easy1.html" "bar" }}
{{ partialCached "easy1.html" "baz" }}
{{ partialCached "easy2.html" "baz" }}
-- layouts/partials/easy1.html --
ABCD
-- layouts/partials/easy2.html --
ABCDE
-- layouts/partials/heavy.html --
{{ $result := slice }}
{{ range site.RegularPages }}
{{ $result = $result | append (dict "title" .Title "link" .RelPermalink "readingTime" .ReadingTime) }}
{{ end }}
{{ range $result }}
* {{ .title }} {{ .link }} {{ .readingTime }}
{{ end }}


`

	for i := 1; i < 100; i++ {
		files += fmt.Sprintf("\n-- content/p%d.md --\n---\ntitle: page\n---\n"+strings.Repeat("FOO ", i), i)
	}

	cfg := hugolib.IntegrationTestConfig{
		T:           b,
		TxtarString: files,
	}
	builders := make([]*hugolib.IntegrationTestBuilder, b.N)

	for i := range builders {
		builders[i] = hugolib.NewIntegrationTestBuilder(cfg)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		builders[i].Build()
	}
}

func TestIncludeTimeout(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
baseURL = 'http://example.com/'
timeout = '200ms'
-- layouts/index.html --
{{ partials.Include "foo.html" . }}
-- layouts/partials/foo.html --
{{ partial "foo.html" . }}
  `

	b, err := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
		},
	).BuildE()

	b.Assert(err, qt.Not(qt.IsNil))
	b.Assert(err.Error(), qt.Contains, "timed out")
}

func TestIncludeCachedTimeout(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
baseURL = 'http://example.com/'
timeout = '200ms'
-- layouts/index.html --
{{ partials.IncludeCached "foo.html" . }}
-- layouts/partials/foo.html --
{{ partialCached "foo.html" . }}
  `

	b, err := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
		},
	).BuildE()

	b.Assert(err, qt.Not(qt.IsNil))
	b.Assert(err.Error(), qt.Contains, "timed out")
}

// See Issue #10789
func TestReturnExecuteFromTemplateInPartial(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
baseURL = 'http://example.com/'
-- layouts/index.html --
{{ $r :=  partial "foo" }}
FOO:{{ $r.Content }}
-- layouts/partials/foo.html --
{{ $r := §§{{ partial "bar" }}§§ | resources.FromString "bar.html" | resources.ExecuteAsTemplate "bar.html" . }}
{{ return $r }}
-- layouts/partials/bar.html --
BAR
  `

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html", "OO:BAR")
}
