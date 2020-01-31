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
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/transform"
)

func TestLiveReloadInject(t *testing.T) {
	c := qt.New(t)

	expectBase := `<script src="/livereload.js?port=1313&mindelay=10&v=2" data-no-instant defer></script>`
	apply := func(s string) string {
		out := new(bytes.Buffer)
		in := strings.NewReader(s)

		tr := transform.New(New(1313))
		tr.Apply(out, in)

		return out.String()
	}

	c.Run("Head lower", func(c *qt.C) {
		c.Assert(apply("<html><head>foo"), qt.Equals, "<html><head>"+expectBase+"foo")
	})

	c.Run("Head upper", func(c *qt.C) {
		c.Assert(apply("<html><HEAD>foo"), qt.Equals, "<html><HEAD>"+expectBase+"foo")
	})

	c.Run("Body lower", func(c *qt.C) {
		c.Assert(apply("foo</body>"), qt.Equals, "foo"+expectBase+"</body>")
	})

	c.Run("Body upper", func(c *qt.C) {
		c.Assert(apply("foo</BODY>"), qt.Equals, "foo"+expectBase+"</BODY>")
	})

	c.Run("No match", func(c *qt.C) {
		c.Assert(apply("<h1>No match</h1>"), qt.Equals, "<h1>No match</h1>")
	})

}
