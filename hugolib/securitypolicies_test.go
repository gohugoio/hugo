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

package hugolib

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/markup/asciidocext"
	"github.com/gohugoio/hugo/markup/pandoc"
	"github.com/gohugoio/hugo/markup/rst"
	"github.com/gohugoio/hugo/resources/resource_transformers/tocss/dartsass"
)

func TestSecurityPolicies(t *testing.T) {
	c := qt.New(t)

	c.Run("os.GetEnv, denied", func(c *qt.C) {
		c.Parallel()
		files := `
-- hugo.toml --
baseURL = "https://example.org"
-- layouts/index.html --
{{ os.Getenv "FOOBAR" }}
`
		_, err := TestE(c, files)
		c.Assert(err, qt.IsNotNil)
		c.Assert(err, qt.ErrorMatches, `(?s).*"FOOBAR" is not whitelisted in policy "security\.funcs\.getenv".*`)
	})

	c.Run("os.GetEnv, OK", func(c *qt.C) {
		c.Parallel()
		files := `
-- hugo.toml --
baseURL = "https://example.org"
-- layouts/index.html --
		{{ os.Getenv "HUGO_FOO" }}
`
		Test(c, files)
	})

	c.Run("Asciidoc, denied", func(c *qt.C) {
		c.Parallel()
		if !asciidocext.Supports() {
			c.Skip()
		}

		files := `
-- hugo.toml --
baseURL = "https://example.org"
-- content/page.ad --
foo
`
		_, err := TestE(c, files)
		c.Assert(err, qt.IsNotNil)
		c.Assert(err, qt.ErrorMatches, `(?s).*"asciidoctor" is not whitelisted in policy "security\.exec\.allow".*`)
	})

	c.Run("RST, denied", func(c *qt.C) {
		c.Parallel()
		if !rst.Supports() {
			c.Skip()
		}

		files := `
-- hugo.toml --
baseURL = "https://example.org"
-- content/page.rst --
foo
`
		_, err := TestE(c, files)
		c.Assert(err, qt.IsNotNil)
		if runtime.GOOS == "windows" {
			c.Assert(err, qt.ErrorMatches, `(?s).*python(\.exe)?" is not whitelisted in policy "security\.exec\.allow".*`)
		} else {
			c.Assert(err, qt.ErrorMatches, `(?s).*"rst2html(\.py)?" is not whitelisted in policy "security\.exec\.allow".*`)
		}
	})

	c.Run("Pandoc, denied", func(c *qt.C) {
		c.Parallel()
		if !pandoc.Supports() {
			c.Skip()
		}

		files := `
-- hugo.toml --
baseURL = "https://example.org"
-- content/page.pdc --
foo
`
		_, err := TestE(c, files)
		c.Assert(err, qt.IsNotNil)
		c.Assert(err, qt.ErrorMatches, `(?s).*pandoc" is not whitelisted in policy "security\.exec\.allow".*`)
	})

	c.Run("Dart SASS, OK", func(c *qt.C) {
		c.Parallel()
		if !dartsass.Supports() {
			c.Skip()
		}
		files := `
-- hugo.toml --
baseURL = "https://example.org"
-- layouts/index.html --
{{ $scss := "body { color: #333; }" | resources.FromString "foo.scss"  | css.Sass (dict "transpiler" "dartsass") }}
`
		Test(c, files)
	})

	c.Run("Dart SASS, denied", func(c *qt.C) {
		c.Parallel()
		if !dartsass.Supports() {
			c.Skip()
		}
		files := `
-- hugo.toml --
baseURL = "https://example.org"
[security]
[security.exec]
allow="none"
-- layouts/index.html --
{{ $scss := "body { color: #333; }" | resources.FromString "foo.scss"  | css.Sass (dict "transpiler" "dartsass") }}
		`
		_, err := TestE(c, files)
		c.Assert(err, qt.IsNotNil)
		c.Assert(err, qt.ErrorMatches, `(?s).*sass(-embedded)?" is not whitelisted in policy "security\.exec\.allow".*`)
	})

	c.Run("resources.GetRemote, OK", func(c *qt.C) {
		c.Parallel()
		ts := httptest.NewServer(http.FileServer(http.Dir("testdata/")))
		c.Cleanup(func() {
			ts.Close()
		})
		files := fmt.Sprintf(`
-- hugo.toml --
baseURL = "https://example.org"
-- layouts/index.html --
{{ $json := resources.GetRemote "%s/fruits.json" }}{{ $json.Content }}
`, ts.URL)
		Test(c, files)
	})

	c.Run("resources.GetRemote, denied method", func(c *qt.C) {
		c.Parallel()
		ts := httptest.NewServer(http.FileServer(http.Dir("testdata/")))
		c.Cleanup(func() {
			ts.Close()
		})
		files := fmt.Sprintf(`
-- hugo.toml --
baseURL = "https://example.org"
-- layouts/index.html --
{{ $json := resources.GetRemote "%s/fruits.json" (dict "method" "DELETE" ) }}{{ $json.Content }}
`, ts.URL)
		_, err := TestE(c, files)
		c.Assert(err, qt.IsNotNil)
		c.Assert(err, qt.ErrorMatches, `(?s).*"DELETE" is not whitelisted in policy "security\.http\.method".*`)
	})

	c.Run("resources.GetRemote, denied URL", func(c *qt.C) {
		c.Parallel()
		ts := httptest.NewServer(http.FileServer(http.Dir("testdata/")))
		c.Cleanup(func() {
			ts.Close()
		})
		files := fmt.Sprintf(`
-- hugo.toml --
baseURL = "https://example.org"
[security]
[security.http]
urls="none"
-- layouts/index.html --
{{ $json := resources.GetRemote "%s/fruits.json" }}{{ $json.Content }}
`, ts.URL)
		_, err := TestE(c, files)
		c.Assert(err, qt.IsNotNil)
		c.Assert(err, qt.ErrorMatches, `(?s).*is not whitelisted in policy "security\.http\.urls".*`)
	})

	c.Run("resources.GetRemote, fake JSON", func(c *qt.C) {
		c.Parallel()
		ts := httptest.NewServer(http.FileServer(http.Dir("testdata/")))
		c.Cleanup(func() {
			ts.Close()
		})
		files := fmt.Sprintf(`
-- hugo.toml --
baseURL = "https://example.org"
[security]
-- layouts/index.html --
{{ $json := resources.GetRemote "%s/fakejson.json" }}{{ $json.Content }}
`, ts.URL)
		_, err := TestE(c, files)
		c.Assert(err, qt.IsNotNil)
		c.Assert(err, qt.ErrorMatches, `(?s).*failed to resolve media type.*`)
	})

	c.Run("resources.GetRemote, fake JSON whitelisted", func(c *qt.C) {
		c.Parallel()
		ts := httptest.NewServer(http.FileServer(http.Dir("testdata/")))
		c.Cleanup(func() {
			ts.Close()
		})
		files := fmt.Sprintf(`
-- hugo.toml --
baseURL = "https://example.org"
[security]
[security.http]
mediaTypes=["application/json"]
-- layouts/index.html --
{{ $json := resources.GetRemote "%s/fakejson.json" }}{{ $json.Content }}
`, ts.URL)
		Test(c, files)
	})
}
