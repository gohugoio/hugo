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

package templates_test

import (
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/hugolib"
)

func TestDecoratorInnerNeverCalled(t *testing.T) {
	t.Parallel()

	filesTemplate := `
-- hugo.toml --
disableKinds = ["section", "taxonomy", "term", "sitemap", "RSS"]
-- content/p1.md --
---
title: "Page 1"
---
-- content/p2.md --
---
title: "Page 2"
---
-- layouts/_partials/cards.html --
Start:{{ range . }}{{ PLACEHOLDER . }}{{ end }}End$
-- layouts/home.html --
1:${{ with partial "cards.html" (site.RegularPages) }}{{ printf "Got %T" . }}|{{ end }}$
2:${{ with partial "cards.html" (site.RegularPages | first 0) }}{{ printf "Got %T" . }}|{{ end }}$
`

	for _, placeholder := range []string{"inner", "templates.Inner"} {
		files := strings.ReplaceAll(filesTemplate, "PLACEHOLDER", placeholder)
		b := hugolib.Test(t, files)

		b.AssertFileContent("public/index.html",
			"1:$Start:Got *hugolib.pageState|Got *hugolib.pageState|End$$",
			"2:$Start:End$$",
		)

	}
}

func TestDecoratorInlinePartialInnerNeverCalled(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ["section", "taxonomy", "term", "sitemap", "RSS"]
-- content/p1.md --
---
title: "Page 1"
---
-- content/p2.md --
---
title: "Page 2"
---
-- layouts/home.html --
${{ with partial "cards.html" (site.RegularPages) }}{{ printf "Got %T" . }}|{{ end }}$
{{ define "_partials/cards.html" }}Start:{{ range . }}{{ inner . }}{{ end }}End${{ end }}
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html",
		"$Start:Got *hugolib.pageState|Got *hugolib.pageState|End$$",
	)
}

func TestDecoratorInlinePartial(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ["section", "taxonomy", "term", "sitemap", "rss"]
-- layouts/home.html --
Home.
{{ with partial "decorate.html" "Important!" }}Notice: {{ . }}{{ end }}
{{ define "_partials/decorate.html" }}<b>{{ inner . }}</b>{{ end }}
`
	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html", "<b>Notice: Important!</b>")
}

func TestDecoratorNestedSimple(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ["section", "taxonomy", "term", "sitemap", "rss"]
-- layouts/home.html --
Home.
{{ with partial "a.html" "warning" }}{{ with partial "b.html" . }}{{ with partial "c.html" . }}{{ . }}{{ end }}{{ end }}{{ end }}
-- layouts/_partials/a.html --
<a>{{ inner . }}</b>
-- layouts/_partials/b.html --
<b>{{ inner . }}</b>
-- layouts/_partials/c.html --
<c>{{ inner . }}</c>
`
	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html", "<a><b><c>warning</c></b></b>")
}

// TODO1 detect inner usage inside with partial. That' a bite in the tail.
func TestDecoratorNested2(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ["section", "taxonomy", "term", "sitemap", "RSS"]
title = "Test title"
-- content/p1.md --
---
title: "Page 1"
---
-- content/p2.md --
---
title: "Page 2"
---
-- layouts/page.html --
{{ .Title }}
-- layouts/home.html --
{{ $pages := site.RegularPages }}
{{ with partial "ul.html" $pages }}<a href="{{ .RelPermalink }}">{{ with partial "bold.html" . }}<span>{{ .LinkTitle }}</span>{{ end }}</a>{{ end }}
-- layouts/_partials/ul.html --
<ul>
{{- range . }}
   <li>{{ inner . }}</li>
{{- end }}
</ul>
-- layouts/_partials/bold.html --
<b>{{ inner $ }}</b>
`

	b, err := hugolib.TestE(t, files)

	b.Assert(err, qt.IsNil)
	b.AssertFileContent("public/index.html", `
<ul>
   <li><a href="/p1/"><b><span>Page 1</span></b></a></li>
   <li><a href="/p2/"><b><span>Page 2</span></b></a></li>
</ul>
`)
}

func TestDecoratorMultiple(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
title = "Test title"
-- layouts/home.html --
{{ with partial "d1.html" . }}X2{{ . }}X4{{ end }}
-- layouts/_partials/d1.html --
X1{{ inner "X3" }}X5
{{ with partial "d2.html" . }}X7{{ . }}X9{{ end }}
{{ with partial "noinner.html" "N3" }}N1{{ . }}N5{{ end }}
X14{{ inner "X15" }}X16
-- layouts/_partials/d2.html --
X6{{ inner "X8" }}X10
{{ with partial "d3.html" . }}A1{{ . }}A2{{ end }}
X11{{ inner "X12" }}X13
-- layouts/_partials/d3.html --
A3{{ inner "A4" }}A5
X6{{ inner "A7" }}A8
-- layouts/_partials/noinner.html --
N2{{ . }}N4
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html",
		"X1X2X3X4X5",
		"X6X7X8X9X10",
		"X11X7X12X9X13",
		"X14X2X15X4X16",
		"A3A1A4A2A5",
		"N1N2N3N4N5", // partial with with, but no inner.
	)
}
