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

package tableofcontents

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestToc(t *testing.T) {
	c := qt.New(t)

	toc := &Root{}

	toc.AddAt(Header{Text: "Header 1", ID: "h1-1"}, 0, 0)
	toc.AddAt(Header{Text: "1-H2-1", ID: "1-h2-1"}, 0, 1)
	toc.AddAt(Header{Text: "1-H2-2", ID: "1-h2-2"}, 0, 1)
	toc.AddAt(Header{Text: "1-H3-1", ID: "1-h2-2"}, 0, 2)
	toc.AddAt(Header{Text: "Header 2", ID: "h1-2"}, 1, 0)

	got := toc.ToHTML(1, -1, false)
	c.Assert(got, qt.Equals, `<nav id="TableOfContents">
  <ul>
    <li><a href="#h1-1">Header 1</a>
      <ul>
        <li><a href="#1-h2-1">1-H2-1</a></li>
        <li><a href="#1-h2-2">1-H2-2</a>
          <ul>
            <li><a href="#1-h2-2">1-H3-1</a></li>
          </ul>
        </li>
      </ul>
    </li>
    <li><a href="#h1-2">Header 2</a></li>
  </ul>
</nav>`, qt.Commentf(got))

	got = toc.ToHTML(1, 1, false)
	c.Assert(got, qt.Equals, `<nav id="TableOfContents">
  <ul>
    <li><a href="#h1-1">Header 1</a></li>
    <li><a href="#h1-2">Header 2</a></li>
  </ul>
</nav>`, qt.Commentf(got))

	got = toc.ToHTML(1, 2, false)
	c.Assert(got, qt.Equals, `<nav id="TableOfContents">
  <ul>
    <li><a href="#h1-1">Header 1</a>
      <ul>
        <li><a href="#1-h2-1">1-H2-1</a></li>
        <li><a href="#1-h2-2">1-H2-2</a></li>
      </ul>
    </li>
    <li><a href="#h1-2">Header 2</a></li>
  </ul>
</nav>`, qt.Commentf(got))

	got = toc.ToHTML(2, 2, false)
	c.Assert(got, qt.Equals, `<nav id="TableOfContents">
  <ul>
    <li><a href="#1-h2-1">1-H2-1</a></li>
    <li><a href="#1-h2-2">1-H2-2</a></li>
  </ul>
</nav>`, qt.Commentf(got))

	got = toc.ToHTML(1, -1, true)
	c.Assert(got, qt.Equals, `<nav id="TableOfContents">
  <ol>
    <li><a href="#h1-1">Header 1</a>
      <ol>
        <li><a href="#1-h2-1">1-H2-1</a></li>
        <li><a href="#1-h2-2">1-H2-2</a>
          <ol>
            <li><a href="#1-h2-2">1-H3-1</a></li>
          </ol>
        </li>
      </ol>
    </li>
    <li><a href="#h1-2">Header 2</a></li>
  </ol>
</nav>`, qt.Commentf(got))
}

func TestTocMissingParent(t *testing.T) {
	c := qt.New(t)

	toc := &Root{}

	toc.AddAt(Header{Text: "H2", ID: "h2"}, 0, 1)
	toc.AddAt(Header{Text: "H3", ID: "h3"}, 1, 2)
	toc.AddAt(Header{Text: "H3", ID: "h3"}, 1, 2)

	got := toc.ToHTML(1, -1, false)
	c.Assert(got, qt.Equals, `<nav id="TableOfContents">
  <ul>
    <li>
      <ul>
        <li><a href="#h2">H2</a></li>
      </ul>
    </li>
    <li>
      <ul>
        <li>
          <ul>
            <li><a href="#h3">H3</a></li>
            <li><a href="#h3">H3</a></li>
          </ul>
        </li>
      </ul>
    </li>
  </ul>
</nav>`, qt.Commentf(got))

	got = toc.ToHTML(3, 3, false)
	c.Assert(got, qt.Equals, `<nav id="TableOfContents">
  <ul>
    <li><a href="#h3">H3</a></li>
    <li><a href="#h3">H3</a></li>
  </ul>
</nav>`, qt.Commentf(got))

	got = toc.ToHTML(1, -1, true)
	c.Assert(got, qt.Equals, `<nav id="TableOfContents">
  <ol>
    <li>
      <ol>
        <li><a href="#h2">H2</a></li>
      </ol>
    </li>
    <li>
      <ol>
        <li>
          <ol>
            <li><a href="#h3">H3</a></li>
            <li><a href="#h3">H3</a></li>
          </ol>
        </li>
      </ol>
    </li>
  </ol>
</nav>`, qt.Commentf(got))

}
