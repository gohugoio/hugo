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

package cssjs_test

import (
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/hugolib"
)

func TestInlineImports(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['page','rss','section','sitemap','taxonomy','term']
-- assets/css/styles.css --
@import "components/a.css";
@import "components/b.css";
h1 { color: red; }
-- assets/css/components/a.css --
.a { color: blue; }
-- assets/css/components/b.css --
@import "a.css";
.b { color: green; }
-- layouts/home.html --
{{ $css := resources.Get "css/styles.css" | css.InlineImports }}
CSS:{{ $css.Content | safeCSS }}:END
`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
		}).Build()

	b.AssertFileContent("public/index.html",
		".a { color: blue; }",
		".b { color: green; }",
		"h1 { color: red; }",
	)
}

func TestInlineImportsNotFound(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['page','rss','section','sitemap','taxonomy','term']
-- assets/css/styles.css --
@import "components/doesnotexist.css";
h1 { color: red; }
-- layouts/home.html --
{{ $css := resources.Get "css/styles.css" | css.InlineImports }}
CSS:{{ $css.Content | safeCSS }}:END
`

	_, err := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
		}).BuildE()

	b := qt.New(t)
	b.Assert(err, qt.IsNotNil)
	b.Assert(err.Error(), qt.Contains, `failed to resolve CSS @import`)
}

func TestInlineImportsSkipNotFound(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['page','rss','section','sitemap','taxonomy','term']
-- assets/css/styles.css --
@import "components/doesnotexist.css";
h1 { color: red; }
-- layouts/home.html --
{{ $opts := dict "skipInlineImportsNotFound" true }}
{{ $css := resources.Get "css/styles.css" | css.InlineImports $opts }}
RelPermalink: {{ $css.RelPermalink }}
`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
		}).Build()

	b.AssertFileContent("public/css/styles.css",
		`@import "components/doesnotexist.css";`,
		"h1 { color: red; }",
	)
}
