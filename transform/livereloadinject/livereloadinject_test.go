// Copyright 2018 The Hugo Authors. All rights reserved.
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

package livereloadinject

import (
	"bytes"
	"io"
	"net/url"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/transform"
)

func TestLiveReloadInject(t *testing.T) {
	c := qt.New(t)

	lrurl, err := url.Parse("http://localhost:1234/subpath")
	if err != nil {
		t.Errorf("Parsing test URL failed")
		return
	}
	expectBase := `<script src="/subpath/livereload.js?mindelay=10&amp;v=2&amp;port=1234&amp;path=subpath/livereload" data-no-instant defer></script>`
	apply := func(s string) string {
		out := new(bytes.Buffer)
		in := strings.NewReader(s)

		tr := transform.New(New(lrurl))
		tr.Apply(out, in)

		return out.String()
	}

	c.Run("Inject after head tag", func(c *qt.C) {
		c.Assert(apply("<!doctype html><html><head>after"), qt.Equals, "<!doctype html><html><head>"+expectBase+"after")
	})

	c.Run("Inject after head tag when doctype and html omitted", func(c *qt.C) {
		c.Assert(apply("<head>after"), qt.Equals, "<head>"+expectBase+"after")
	})

	c.Run("Inject after html when head omitted", func(c *qt.C) {
		c.Assert(apply("<html>after"), qt.Equals, "<html>"+expectBase+"after")
	})

	c.Run("Inject after doctype when head and html omitted", func(c *qt.C) {
		c.Assert(apply("<!doctype html>after"), qt.Equals, "<!doctype html>"+expectBase+"after")
	})

	c.Run("Inject before other elements if all else omitted", func(c *qt.C) {
		c.Assert(apply("<title>after</title>"), qt.Equals, expectBase+"<title>after</title>")
	})

	c.Run("Inject before text content if all else omitted", func(c *qt.C) {
		c.Assert(apply("after"), qt.Equals, expectBase+"after")
	})

	c.Run("Inject after HeAd tag MiXed CaSe", func(c *qt.C) {
		c.Assert(apply("<HeAd>AfTer"), qt.Equals, "<HeAd>"+expectBase+"AfTer")
	})

	c.Run("Inject after HtMl tag MiXed CaSe", func(c *qt.C) {
		c.Assert(apply("<HtMl>AfTer"), qt.Equals, "<HtMl>"+expectBase+"AfTer")
	})

	c.Run("Inject after doctype mixed case", func(c *qt.C) {
		c.Assert(apply("<!DocType HtMl>AfTer"), qt.Equals, "<!DocType HtMl>"+expectBase+"AfTer")
	})

	c.Run("Inject after html tag with attributes", func(c *qt.C) {
		c.Assert(apply(`<html lang="en">after`), qt.Equals, `<html lang="en">`+expectBase+"after")
	})

	c.Run("Inject after html tag with newline", func(c *qt.C) {
		c.Assert(apply("<html\n>after"), qt.Equals, "<html\n>"+expectBase+"after")
	})

	c.Run("Skip comments and whitespace", func(c *qt.C) {
		c.Assert(
			apply(" <!--x--> <!doctype html>\n<?xml instruction ?> <head>after"),
			qt.Equals,
			" <!--x--> <!doctype html>\n<?xml instruction ?> <head>"+expectBase+"after",
		)
	})

	c.Run("Do not search inside comment", func(c *qt.C) {
		c.Assert(apply("<html><!--<head>-->"), qt.Equals, "<html><!--<head>-->"+expectBase)
	})

	c.Run("Do not search inside scripts", func(c *qt.C) {
		c.Assert(apply("<html><script>`<head>`</script>"), qt.Equals, "<html>"+expectBase+"<script>`<head>`</script>")
	})

	c.Run("Do not search inside templates", func(c *qt.C) {
		c.Assert(apply("<html><template><head></template>"), qt.Not(qt.Equals), "<html><template><head>"+expectBase+"</template>")
	})

	c.Run("Search from the start of the input", func(c *qt.C) {
		c.Assert(apply("<head>after<head>"), qt.Equals, "<head>"+expectBase+"after<head>")
	})

	c.Run("Do not mistake header for head", func(c *qt.C) {
		c.Assert(apply("<html><header>"), qt.Equals, "<html>"+expectBase+"<header>")
	})

	c.Run("Do not mistake custom elements for head", func(c *qt.C) {
		c.Assert(apply("<html><head-custom>"), qt.Equals, "<html>"+expectBase+"<head-custom>")
	})
}

func BenchmarkLiveReloadInject(b *testing.B) {
	s := `
<html>
<head>
</head>
<body>
</body>
</html>	
`
	in := strings.NewReader(s)
	lrurl, err := url.Parse("http://localhost:1234/subpath")
	if err != nil {
		b.Fatalf("Parsing test URL failed")
	}
	tr := transform.New(New(lrurl))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		in.Seek(0, 0)
		tr.Apply(io.Discard, in)
	}
}
