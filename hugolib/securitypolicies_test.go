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
	"sync/atomic"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/markup/asciidocext"
	"github.com/gohugoio/hugo/markup/pandoc"
	"github.com/gohugoio/hugo/markup/rst"
	"github.com/gohugoio/hugo/resources/resource_transformers/tocss/dartsass"
)

func TestSecurityPolicies(t *testing.T) {
	c := qt.New(t)

	c.Run("HTML content, denied by default", func(c *qt.C) {
		c.Parallel()
		files := `
-- hugo.toml --
baseURL = "https://example.org"
-- content/page.html --
---
title: "Untrusted"
---
<script>alert(1)</script>
-- layouts/single.html --
{{ .Content }}
`
		_, err := TestE(c, files)
		c.Assert(err, qt.IsNotNil)
		c.Assert(err, qt.ErrorMatches, `(?s).*"text/html" is not whitelisted in policy "security\.allowContent".*`)
	})

	c.Run("HTML content, allowed via override", func(c *qt.C) {
		c.Parallel()
		files := `
-- hugo.toml --
baseURL = "https://example.org"
[security]
allowContent = ['.*']
-- content/page.html --
---
title: "Trusted"
---
<p>hello</p>
-- layouts/single.html --
{{ .Content }}
`
		b := Test(c, files)
		b.AssertFileContent("public/page/index.html", "<p>hello</p>")
	})

	c.Run("HTML content from content adapter, denied by default", func(c *qt.C) {
		c.Parallel()
		files := `
-- hugo.toml --
baseURL = "https://example.org"
-- content/_content.gotmpl --
{{ .AddPage (dict "path" "p1" "title" "Untrusted" "content" (dict "value" "<script>alert(1)</script>" "mediaType" "text/html")) }}
-- layouts/single.html --
{{ .Content }}
`
		_, err := TestE(c, files)
		c.Assert(err, qt.IsNotNil)
		c.Assert(err, qt.ErrorMatches, `(?s).*"text/html" is not whitelisted in policy "security\.allowContent".*`)
	})

	c.Run("HTML content from content adapter, allowed via override", func(c *qt.C) {
		c.Parallel()
		files := `
-- hugo.toml --
baseURL = "https://example.org"
[security]
allowContent = ['.*']
-- content/_content.gotmpl --
{{ .AddPage (dict "path" "p1" "title" "Trusted" "content" (dict "value" "<p>hello</p>" "mediaType" "text/html")) }}
-- layouts/single.html --
{{ .Content }}
`
		b := Test(c, files)
		b.AssertFileContent("public/p1/index.html", "<p>hello</p>")
	})

	c.Run("Org content, denied by default", func(c *qt.C) {
		c.Parallel()
		files := `
-- hugo.toml --
baseURL = "https://example.org"
-- content/page.org --
---
title: "Untrusted"
---
@@html:<script>alert(1)</script>@@
-- layouts/single.html --
{{ .Content }}
`
		_, err := TestE(c, files)
		c.Assert(err, qt.IsNotNil)
		c.Assert(err, qt.ErrorMatches, `(?s).*"text/org" is not whitelisted in policy "security\.allowContent".*`)
	})

	c.Run("Org content, allowed via override", func(c *qt.C) {
		c.Parallel()
		files := `
-- hugo.toml --
baseURL = "https://example.org"
[security]
allowContent = ['.*']
-- content/page.org --
---
title: "Trusted"
---
hello
-- layouts/single.html --
{{ .Content }}
`
		b := Test(c, files)
		b.AssertFileContent("public/page/index.html", "hello")
	})

	c.Run("Org content from content adapter, denied by default", func(c *qt.C) {
		c.Parallel()
		files := `
-- hugo.toml --
baseURL = "https://example.org"
-- content/_content.gotmpl --
{{ .AddPage (dict "path" "p1" "title" "Untrusted" "content" (dict "value" "@@html:<script>alert(1)</script>@@" "mediaType" "text/org")) }}
-- layouts/single.html --
{{ .Content }}
`
		_, err := TestE(c, files)
		c.Assert(err, qt.IsNotNil)
		c.Assert(err, qt.ErrorMatches, `(?s).*"text/org" is not whitelisted in policy "security\.allowContent".*`)
	})

	c.Run("os.GetEnv, denied", func(c *qt.C) {
		c.Parallel()
		files := `
-- hugo.toml --
baseURL = "https://example.org"
-- layouts/home.html --
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
-- layouts/home.html --
		{{ os.Getenv "HUGO_FOO" }}
`
		Test(c, files)
	})

	c.Run("AsciiDoc, denied", func(c *qt.C) {
		c.Parallel()
		if ok, err := asciidocext.Supports(); !ok {
			c.Skip(err)
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
-- layouts/home.html --
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
-- layouts/home.html --
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
[security]
[security.http]
urls = ['.*']
-- layouts/home.html --
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
[security]
[security.http]
urls = ['.*']
-- layouts/home.html --
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
-- layouts/home.html --
{{ $json := resources.GetRemote "%s/fruits.json" }}{{ $json.Content }}
`, ts.URL)
		_, err := TestE(c, files)
		c.Assert(err, qt.IsNotNil)
		c.Assert(err, qt.ErrorMatches, `(?s).*is not whitelisted in policy "security\.http\.urls".*`)
	})

	c.Run("resources.GetRemote, denied loopback URL by default", func(c *qt.C) {
		c.Parallel()
		ts := httptest.NewServer(http.FileServer(http.Dir("testdata/")))
		c.Cleanup(func() {
			ts.Close()
		})
		files := fmt.Sprintf(`
-- hugo.toml --
baseURL = "https://example.org"
-- layouts/home.html --
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
[security.http]
urls = ['.*']
-- layouts/home.html --
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
urls = ['.*']
mediaTypes=["application/json"]
-- layouts/home.html --
{{ $json := resources.GetRemote "%s/fakejson.json" }}{{ $json.Content }}
`, ts.URL)
		Test(c, files)
	})
}

// See issue 15302.
func TestProxyFromEnvironment(t *testing.T) {
	c := qt.New(t)

	var proxyHit atomic.Bool
	proxySrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxyHit.Store(true)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("from-proxy"))
	}))
	t.Cleanup(proxySrv.Close)

	t.Setenv("HTTP_PROXY", proxySrv.URL)

	// Case 1: Default config (proxyFromEnvironment = false).
	filesDefault := `
-- hugo.toml --
baseURL = "https://example.org"
[security.http]
urls = ['.*']
-- layouts/home.html --
{{ $res := resources.GetRemote "http://example.com/test.txt" }}
`
	proxyHit.Store(false)
	TestE(c, filesDefault)
	c.Assert(proxyHit.Load(), qt.IsFalse)

	// Case 2: Explicit proxyFromEnvironment = true.
	filesProxied := `
-- hugo.toml --
baseURL = "https://example.org"
[security.http]
urls = ['.*']
proxyFromEnvironment = true
-- layouts/home.html --
{{ $res := resources.GetRemote "http://example.com/test.txt" }}{{ $res.Content }}
`
	proxyHit.Store(false)
	b := Test(c, filesProxied)
	c.Assert(proxyHit.Load(), qt.IsTrue)
	b.AssertFileContent("public/index.html", "from-proxy")
}
