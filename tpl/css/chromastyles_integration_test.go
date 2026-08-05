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

package css_test

import (
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/hugolib"
	"github.com/gohugoio/hugo/tpl/css"
)

// See issue 15112.
func TestChromaStyles(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ["taxonomy", "term", "rss", "sitemap", "section", "page"]
[markup.highlight]
style = "monokai"
-- layouts/home.html --
{{ $light := css.ChromaStyles (dict "targetPath" "css/chroma-light.css" "style" "github" "mode" "light" "modeSelector" true) }}
{{ $dark := css.ChromaStyles (dict "targetPath" "css/chroma-dark.css" "style" "github" "mode" "dark" "modeSelector" true) }}
{{ $default := css.ChromaStyles (dict "targetPath" "css/chroma.css") }}
Light: {{ $light.RelPermalink }}|Dark: {{ $dark.RelPermalink }}|Default: {{ $default.RelPermalink }}|
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html", "Light: /css/chroma-light.css|Dark: /css/chroma-dark.css|Default: /css/chroma.css|")
	b.AssertFileContent("public/css/chroma-light.css", ".light .chroma ", ".light .bg {")
	b.AssertFileContent("public/css/chroma-dark.css", ".dark .chroma ", ".dark .bg {")
	// Style defaults to markup.highlight.style, monokai background.
	b.AssertFileContent("public/css/chroma.css", "/* Background */ .bg {", "background-color:#272822")
}

// See issue 15167.
func TestChromaStylesCustomModeClasses(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ["taxonomy", "term", "rss", "sitemap", "section", "page"]
-- layouts/home.html --
{{ $light := css.ChromaStyles (dict "targetPath" "css/chroma-light.css" "style" "github" "mode" "light" "modeSelector" true "classLight" "theme-light") }}
{{ $dark := css.ChromaStyles (dict "targetPath" "css/chroma-dark.css" "style" "github" "mode" "dark" "modeSelector" true "classDark" "theme-dark") }}
Light: {{ $light.RelPermalink }}|Dark: {{ $dark.RelPermalink }}|
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/css/chroma-light.css", ".theme-light .chroma ", ".theme-light .bg {")
	b.AssertFileContent("public/css/chroma-dark.css", ".theme-dark .chroma ", ".theme-dark .bg {")
}

// Chroma's minifier drops token rules whose color equals the style's default
// foreground (e.g. .nx in github-dark). In a paired light/dark setup the light
// sheet's explicit rule would then leak into dark mode, so these must be re-emitted.
// See issue 15161.
func TestChromaStylesKeepDefaultForegroundTokens(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ["taxonomy", "term", "rss", "sitemap", "section", "page"]
-- layouts/home.html --
{{ $light := css.ChromaStyles (dict "targetPath" "css/chroma-light.css" "style" "github" "mode" "light") }}
{{ $dark := css.ChromaStyles (dict "targetPath" "css/chroma-dark.css" "style" "github" "mode" "dark" "modeSelector" true) }}
Light: {{ $light.RelPermalink }}|Dark: {{ $dark.RelPermalink }}|
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/css/chroma-light.css", "/* NameOther */ .chroma .nx { color:#1f2328 }")
	b.AssertFileContent("public/css/chroma-dark.css", "/* NameOther */ .dark .chroma .nx { color:#e6edf3 }")
}

func TestChromaStylesErrors(t *testing.T) {
	t.Parallel()

	filesTemplate := `
-- hugo.toml --
disableKinds = ["taxonomy", "term", "rss", "sitemap", "section", "page"]
-- layouts/home.html --
{{ css.ChromaStyles OPTS }}
`

	for _, test := range []struct {
		opts  string
		match string
	}{
		{`(dict "style" "monokai")`, `.*targetPath cannot be empty.*`},
		{`(dict "targetPath" "css/chroma.css" "style" "doesnotexist")`, `.*invalid style: doesnotexist.*`},
		{`(dict "targetPath" "css/chroma.css" "mode" "blue")`, `.*invalid mode: blue.*`},
		{`(dict "targetPath" "css/chroma.css" "style" "friendly" "mode" "dark")`, `.*style "friendly" does not have a "dark" mode.*`},
	} {
		files := strings.ReplaceAll(filesTemplate, "OPTS", test.opts)
		b, err := hugolib.TestE(t, files)
		b.Assert(err, qt.ErrorMatches, test.match)
	}
}

func BenchmarkChromaStyles(b *testing.B) {
	files := `
-- hugo.toml --
disableKinds = ["taxonomy", "term", "rss", "sitemap", "section", "page"]
-- layouts/all.html --
All.
`

	bb := hugolib.Test(b, files)

	bb.Assert(true, qt.IsTrue)
	s := bb.H.Sites[0]
	ns := s.TemplateStore.GetTemplateFuncsNamespace("css").(*css.Namespace)

	opts := map[string]any{"targetPath": "css/chroma.css", "style": "monokai", "mode": "light", "modeSelector": true}

	b.ResetTimer()

	for b.Loop() {
		_, err := ns.ChromaStyles(opts)
		bb.Assert(err, qt.IsNil)
	}
}
