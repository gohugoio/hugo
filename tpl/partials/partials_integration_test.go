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

	"github.com/gohugoio/hugo/htesting"
	"github.com/gohugoio/hugo/htesting/hqt"
	"github.com/gohugoio/hugo/hugolib"
)

func TestInclude(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = 'http://example.com/'
-- layouts/home.html --
partial: {{ partials.Include "foo.html" . }}
-- layouts/_partials/foo.html --
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
-- hugo.toml --
baseURL = 'http://example.com/'
-- layouts/home.html --
partialCached: {{ partials.IncludeCached "foo.html" . }}
partialCached: {{ partials.IncludeCached "foo.html" . }}
-- layouts/_partials/foo.html --
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
-- hugo.toml --
baseURL = 'http://example.com/'
-- layouts/home.html --
{{ partials.IncludeCached "p1.html" . }}
-- layouts/_partials/p1.html --
{{ partials.IncludeCached "p2.html" . }}
-- layouts/_partials/p2.html --
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
-- hugo.toml --
baseURL = 'http://example.com/'
-- content/_index.md --
---
title: "Index"
---
{{< short >}}
-- layouts/home.html --
{{ partials.IncludeCached "p1.html" . }}
-- layouts/_partials/p1.html --
{{ .Content }}
{{ partials.IncludeCached "p2.html" . }}
-- layouts/_partials/p2.html --
-- layouts/_shortcodes/short.html --
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
-- hugo.toml --
baseURL = 'http://example.com/'
templateMetrics=true
templateMetricsHints=true
disableKinds = ["page", "section", "taxonomy", "term", "sitemap"]
[outputs]
home = ["HTML"]
-- layouts/home.html --
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

-- layouts/_partials/static1.html --
P1
-- layouts/_partials/static2.html --
P2
-- layouts/_partials/dynamic1.html --
{{ math.Counter }}
-- layouts/_partials/halfdynamic1.html --
D1
{{ math.Counter }}


  `

	b := hugolib.Test(t, files)

	// fmt.Println(b.FileContent("public/index.html"))

	var buf bytes.Buffer
	b.H.Metrics.WriteMetrics(&buf)

	got := buf.String()

	// Get rid of all the durations, they are never the same.
	durationRe := regexp.MustCompile(`\b[\.\d]*(ms|ns|µs|s)\b`)

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
	0        0       0      1  home.html
	100        0       0      1  _partials/static2.html
	100       50       1      2  _partials/static1.html
	25       50       2      4  _partials/dynamic1.html
	66       33       1      3  _partials/halfdynamic1.html
	`

	b.Assert(got, hqt.IsSameString, expect)
}

// gobench --package ./tpl/partials
func BenchmarkIncludeCached(b *testing.B) {
	files := `
-- hugo.toml --
baseURL = 'http://example.com/'
-- layouts/home.html --
-- layouts/single.html --
{{ partialCached "heavy.html" "foo" }}
{{ partialCached "easy1.html" "bar" }}
{{ partialCached "easy1.html" "baz" }}
{{ partialCached "easy2.html" "baz" }}
-- layouts/_partials/easy1.html --
ABCD
-- layouts/_partials/easy2.html --
ABCDE
-- layouts/_partials/heavy.html --
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

	for b.Loop() {
		b.StopTimer()
		bb := hugolib.NewIntegrationTestBuilder(cfg)
		b.StartTimer()
		bb.Build()
	}
}

func TestIncludeTimeout(t *testing.T) {
	htesting.SkipSlowTestUnlessCI(t)
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = 'http://example.com/'
-- layouts/home.html --
{{ partials.Include "foo.html" . }}
-- layouts/_partials/foo.html --
{{ partial "foo.html" . }}
  `

	b, err := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
		},
	).BuildE()

	b.Assert(err, qt.Not(qt.IsNil))
	b.Assert(err.Error(), qt.Contains, "maximum template call stack size exceeded")
}

func TestIncludeCachedTimeout(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = 'http://example.com/'
timeout = '200ms'
-- layouts/home.html --
{{ partials.IncludeCached "foo.html" . }}
-- layouts/_partials/foo.html --
{{ partialCached "bar.html" . }}
-- layouts/_partials/bar.html --
{{ partialCached "foo.html" . }}
  `

	b, err := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
		},
	).BuildE()

	b.Assert(err, qt.Not(qt.IsNil))
	b.Assert(err.Error(), qt.Contains, `error calling partialCached: circular call stack detected in partial`)
}

// See Issue #13889
func TestIncludeCachedDifferentKey(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = 'http://example.com/'
timeout = '200ms'
-- layouts/home.html --
{{ partialCached "foo.html" "a" "a" }}
-- layouts/_partials/foo.html --
{{ if eq . "a" }}
{{ partialCached "bar.html" . }}
{{ else }}
DONE
{{ end }}
-- layouts/_partials/bar.html --
{{ partialCached "foo.html" "b" "b" }}
  `
	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html", `
DONE
`)
}

// See Issue #10789
func TestReturnExecuteFromTemplateInPartial(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = 'http://example.com/'
-- layouts/home.html --
{{ $r :=  partial "foo" }}
FOO:{{ $r.Content }}
-- layouts/_partials/foo.html --
{{ $r := §§{{ partial "bar" }}§§ | resources.FromString "bar.html" | resources.ExecuteAsTemplate "bar.html" . }}
{{ return $r }}
-- layouts/_partials/bar.html --
BAR
  `

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html", "OO:BAR")
}
