// Copyright 2019 The Hugo Authors. All rights reserved.
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

	testVariant := func(c *qt.C, withBuilder func(b *sitesBuilder), expectErr string) {
		c.Helper()
		b := newTestSitesBuilder(c)
		withBuilder(b)

		if expectErr != "" {
			err := b.BuildE(BuildCfg{})
			b.Assert(err, qt.IsNotNil)
			b.Assert(err, qt.ErrorMatches, expectErr)
		} else {
			b.Build(BuildCfg{})
		}
	}

	httpTestVariant := func(c *qt.C, templ, expectErr string, withBuilder func(b *sitesBuilder)) {
		ts := httptest.NewServer(http.FileServer(http.Dir("testdata/")))
		c.Cleanup(func() {
			ts.Close()
		})
		cb := func(b *sitesBuilder) {
			b.WithTemplatesAdded("index.html", fmt.Sprintf(templ, ts.URL))
			if withBuilder != nil {
				withBuilder(b)
			}
		}
		testVariant(c, cb, expectErr)
	}

	c.Run("os.GetEnv, denied", func(c *qt.C) {
		c.Parallel()
		cb := func(b *sitesBuilder) {
			b.WithTemplatesAdded("index.html", `{{ os.Getenv "FOOBAR" }}`)
		}
		testVariant(c, cb, `(?s).*"FOOBAR" is not whitelisted in policy "security\.funcs\.getenv".*`)
	})

	c.Run("os.GetEnv, OK", func(c *qt.C) {
		c.Parallel()
		cb := func(b *sitesBuilder) {
			b.WithTemplatesAdded("index.html", `{{ os.Getenv "HUGO_FOO" }}`)
		}
		testVariant(c, cb, "")
	})

	c.Run("Asciidoc, denied", func(c *qt.C) {
		c.Parallel()
		if !asciidocext.Supports() {
			c.Skip()
		}

		cb := func(b *sitesBuilder) {
			b.WithContent("page.ad", "foo")
		}

		testVariant(c, cb, `(?s).*"asciidoctor" is not whitelisted in policy "security\.exec\.allow".*`)
	})

	c.Run("RST, denied", func(c *qt.C) {
		c.Parallel()
		if !rst.Supports() {
			c.Skip()
		}

		cb := func(b *sitesBuilder) {
			b.WithContent("page.rst", "foo")
		}

		if runtime.GOOS == "windows" {
			testVariant(c, cb, `(?s).*python(\.exe)?" is not whitelisted in policy "security\.exec\.allow".*`)
		} else {
			testVariant(c, cb, `(?s).*"rst2html(\.py)?" is not whitelisted in policy "security\.exec\.allow".*`)
		}
	})

	c.Run("Pandoc, denied", func(c *qt.C) {
		c.Parallel()
		if !pandoc.Supports() {
			c.Skip()
		}

		cb := func(b *sitesBuilder) {
			b.WithContent("page.pdc", "foo")
		}

		testVariant(c, cb, `(?s).*pandoc" is not whitelisted in policy "security\.exec\.allow".*`)
	})

	c.Run("Dart SASS, OK", func(c *qt.C) {
		c.Parallel()
		if !dartsass.Supports() {
			c.Skip()
		}
		cb := func(b *sitesBuilder) {
			b.WithTemplatesAdded("index.html", `{{ $scss := "body { color: #333; }" | resources.FromString "foo.scss"  | css.Sass (dict "transpiler" "dartsass") }}`)
		}
		testVariant(c, cb, "")
	})

	c.Run("Dart SASS, denied", func(c *qt.C) {
		c.Parallel()
		if !dartsass.Supports() {
			c.Skip()
		}
		cb := func(b *sitesBuilder) {
			b.WithConfigFile("toml", `
[security]
[security.exec]
allow="none"

			`)
			b.WithTemplatesAdded("index.html", `{{ $scss := "body { color: #333; }" | resources.FromString "foo.scss"  | css.Sass (dict "transpiler" "dartsass") }}`)
		}
		testVariant(c, cb, `(?s).*sass(-embedded)?" is not whitelisted in policy "security\.exec\.allow".*`)
	})

	c.Run("resources.GetRemote, OK", func(c *qt.C) {
		c.Parallel()
		httpTestVariant(c, `{{ $json := resources.GetRemote "%[1]s/fruits.json" }}{{ $json.Content }}`, "", nil)
	})

	c.Run("resources.GetRemote, denied method", func(c *qt.C) {
		c.Parallel()
		httpTestVariant(c, `{{ $json := resources.GetRemote "%[1]s/fruits.json" (dict "method" "DELETE" ) }}{{ $json.Content }}`, `(?s).*"DELETE" is not whitelisted in policy "security\.http\.method".*`, nil)
	})

	c.Run("resources.GetRemote, denied URL", func(c *qt.C) {
		c.Parallel()
		httpTestVariant(c, `{{ $json := resources.GetRemote "%[1]s/fruits.json" }}{{ $json.Content }}`, `(?s).*is not whitelisted in policy "security\.http\.urls".*`,
			func(b *sitesBuilder) {
				b.WithConfigFile("toml", `
[security]
[security.http]
urls="none"
`)
			})
	})

	c.Run("resources.GetRemote, fake JSON", func(c *qt.C) {
		c.Parallel()
		httpTestVariant(c, `{{ $json := resources.GetRemote "%[1]s/fakejson.json" }}{{ $json.Content }}`, `(?s).*failed to resolve media type.*`,
			func(b *sitesBuilder) {
				b.WithConfigFile("toml", `
`)
			})
	})

	c.Run("resources.GetRemote, fake JSON whitelisted", func(c *qt.C) {
		c.Parallel()
		httpTestVariant(c, `{{ $json := resources.GetRemote "%[1]s/fakejson.json" }}{{ $json.Content }}`, ``,
			func(b *sitesBuilder) {
				b.WithConfigFile("toml", `
[security]
[security.http]
mediaTypes=["application/json"]

`)
			})
	})
}
