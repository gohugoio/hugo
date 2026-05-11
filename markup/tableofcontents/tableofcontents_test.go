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
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/common/collections"
)

var newTestTocBuilder = func() Builder {
	var b Builder
	b.AddAt(&Heading{Title: "Heading 1", ID: "h1-1"}, 0, 0)
	b.AddAt(&Heading{Title: "1-H2-1", ID: "1-h2-1"}, 0, 1)
	b.AddAt(&Heading{Title: "1-H2-2", ID: "1-h2-2"}, 0, 1)
	b.AddAt(&Heading{Title: "1-H3-1", ID: "1-h2-2"}, 0, 2)
	b.AddAt(&Heading{Title: "Heading 2", ID: "h1-2"}, 1, 0)
	return b
}

var newTestToc = func() *Fragments {
	return newTestTocBuilder().Build()
}

func TestToc(t *testing.T) {
	c := qt.New(t)

	toc := &Fragments{}

	toc.addAt(&Heading{Title: "Heading 1", ID: "h1-1"}, 0, 0)
	toc.addAt(&Heading{Title: "1-H2-1", ID: "1-h2-1"}, 0, 1)
	toc.addAt(&Heading{Title: "1-H2-2", ID: "1-h2-2"}, 0, 1)
	toc.addAt(&Heading{Title: "1-H3-1", ID: "1-h2-2"}, 0, 2)
	toc.addAt(&Heading{Title: "Heading 2", ID: "h1-2"}, 1, 0)

	tocHTML, _ := toc.ToHTML(1, -1, false)
	got := string(tocHTML)
	c.Assert(got, qt.Equals, `<nav id="TableOfContents">
  <ul>
    <li><a href="#h1-1">Heading 1</a>
      <ul>
        <li><a href="#1-h2-1">1-H2-1</a></li>
        <li><a href="#1-h2-2">1-H2-2</a>
          <ul>
            <li><a href="#1-h2-2">1-H3-1</a></li>
          </ul>
        </li>
      </ul>
    </li>
    <li><a href="#h1-2">Heading 2</a></li>
  </ul>
</nav>`, qt.Commentf(got))

	tocHTML, _ = toc.ToHTML(1, 1, false)
	got = string(tocHTML)
	c.Assert(got, qt.Equals, `<nav id="TableOfContents">
  <ul>
    <li><a href="#h1-1">Heading 1</a></li>
    <li><a href="#h1-2">Heading 2</a></li>
  </ul>
</nav>`, qt.Commentf(got))

	tocHTML, _ = toc.ToHTML(1, 2, false)
	got = string(tocHTML)
	c.Assert(got, qt.Equals, `<nav id="TableOfContents">
  <ul>
    <li><a href="#h1-1">Heading 1</a>
      <ul>
        <li><a href="#1-h2-1">1-H2-1</a></li>
        <li><a href="#1-h2-2">1-H2-2</a></li>
      </ul>
    </li>
    <li><a href="#h1-2">Heading 2</a></li>
  </ul>
</nav>`, qt.Commentf(got))

	tocHTML, _ = toc.ToHTML(2, 2, false)
	got = string(tocHTML)
	c.Assert(got, qt.Equals, `<nav id="TableOfContents">
  <ul>
    <li><a href="#1-h2-1">1-H2-1</a></li>
    <li><a href="#1-h2-2">1-H2-2</a></li>
  </ul>
</nav>`, qt.Commentf(got))

	tocHTML, _ = toc.ToHTML(1, -1, true)
	got = string(tocHTML)
	c.Assert(got, qt.Equals, `<nav id="TableOfContents">
  <ol>
    <li><a href="#h1-1">Heading 1</a>
      <ol>
        <li><a href="#1-h2-1">1-H2-1</a></li>
        <li><a href="#1-h2-2">1-H2-2</a>
          <ol>
            <li><a href="#1-h2-2">1-H3-1</a></li>
          </ol>
        </li>
      </ol>
    </li>
    <li><a href="#h1-2">Heading 2</a></li>
  </ol>
</nav>`, qt.Commentf(got))
}

func TestTocMissingParent(t *testing.T) {
	c := qt.New(t)

	toc := &Fragments{}

	toc.addAt(&Heading{Title: "H2", ID: "h2"}, 0, 1)
	toc.addAt(&Heading{Title: "H3", ID: "h3"}, 1, 2)
	toc.addAt(&Heading{Title: "H3", ID: "h3"}, 1, 2)

	tocHTML, _ := toc.ToHTML(1, -1, false)
	got := string(tocHTML)
	c.Assert(got, qt.Equals, `<nav id="TableOfContents">
  <ul>
    <li><a href="#h2">H2</a></li>
    <li><a href="#h3">H3</a></li>
    <li><a href="#h3">H3</a></li>
  </ul>
</nav>`, qt.Commentf(got))

	tocHTML, _ = toc.ToHTML(3, 3, false)
	got = string(tocHTML)
	c.Assert(got, qt.Equals, `<nav id="TableOfContents">
  <ul>
    <li><a href="#h3">H3</a></li>
    <li><a href="#h3">H3</a></li>
  </ul>
</nav>`, qt.Commentf(got))

	tocHTML, _ = toc.ToHTML(1, -1, true)
	got = string(tocHTML)
	c.Assert(got, qt.Equals, `<nav id="TableOfContents">
  <ol>
    <li><a href="#h2">H2</a></li>
    <li><a href="#h3">H3</a></li>
    <li><a href="#h3">H3</a></li>
  </ol>
</nav>`, qt.Commentf(got))
}

func TestTocMissingIntermediateLevels(t *testing.T) {
	c := qt.New(t)

	type item struct {
		title string
		id    string
		row   int
		level int
	}

	for _, test := range []struct {
		name     string
		items    []item
		expected string
	}{
		{
			name: "h2 to h4",
			items: []item{
				{title: "H2", id: "h2", level: 1},
				{title: "H4", id: "h4", level: 3},
			},
			expected: `<nav id="TableOfContents">
  <ul>
    <li><a href="#h2">H2</a>
      <ul>
        <li><a href="#h4">H4</a></li>
      </ul>
    </li>
  </ul>
</nav>`,
		},
		{
			name: "h2 to h5",
			items: []item{
				{title: "H2", id: "h2", level: 1},
				{title: "H5", id: "h5", level: 4},
			},
			expected: `<nav id="TableOfContents">
  <ul>
    <li><a href="#h2">H2</a>
      <ul>
        <li><a href="#h5">H5</a></li>
      </ul>
    </li>
  </ul>
</nav>`,
		},
		{
			name: "h2 to h6",
			items: []item{
				{title: "H2", id: "h2", level: 1},
				{title: "H6", id: "h6", level: 5},
			},
			expected: `<nav id="TableOfContents">
  <ul>
    <li><a href="#h2">H2</a>
      <ul>
        <li><a href="#h6">H6</a></li>
      </ul>
    </li>
  </ul>
</nav>`,
		},
		{
			name: "h3 to h5",
			items: []item{
				{title: "H3", id: "h3", level: 2},
				{title: "H5", id: "h5", level: 4},
			},
			expected: `<nav id="TableOfContents">
  <ul>
    <li><a href="#h3">H3</a>
      <ul>
        <li><a href="#h5">H5</a></li>
      </ul>
    </li>
  </ul>
</nav>`,
		},
		{
			name: "starts at h4",
			items: []item{
				{title: "H4", id: "h4", level: 3},
			},
			expected: `<nav id="TableOfContents">
  <ul>
    <li><a href="#h4">H4</a></li>
  </ul>
</nav>`,
		},
		{
			name: "starts at h6",
			items: []item{
				{title: "H6", id: "h6", level: 5},
			},
			expected: `<nav id="TableOfContents">
  <ul>
    <li><a href="#h6">H6</a></li>
  </ul>
</nav>`,
		},
	} {
		c.Run(test.name, func(c *qt.C) {
			toc := &Fragments{}
			for _, item := range test.items {
				toc.addAt(&Heading{Title: item.title, ID: item.id}, item.row, item.level)
			}

			tocHTML, err := toc.ToHTML(2, -1, false)
			c.Assert(err, qt.IsNil)
			got := string(tocHTML)
			c.Assert(got, qt.Equals, test.expected, qt.Commentf(got))
			c.Assert(got, qt.Not(qt.Contains), "<li>\n")
			c.Assert(hasListItemWithoutAnchor(got), qt.Equals, false)
		})
	}

	toc := &Fragments{}
	toc.addAt(&Heading{Title: "H2", ID: "h2"}, 0, 1)
	toc.addAt(&Heading{Title: "H4", ID: "h4"}, 0, 3)

	tocHTML, err := toc.ToHTML(2, 3, false)
	c.Assert(err, qt.IsNil)
	got := string(tocHTML)
	c.Assert(got, qt.Equals, `<nav id="TableOfContents">
  <ul>
    <li><a href="#h2">H2</a></li>
  </ul>
</nav>`, qt.Commentf(got))
	c.Assert(got, qt.Not(qt.Contains), "<li>\n")
	c.Assert(hasListItemWithoutAnchor(got), qt.Equals, false)
}

func hasListItemWithoutAnchor(s string) bool {
	for {
		i := strings.Index(s, "<li>")
		if i == -1 {
			return false
		}
		s = s[i+len("<li>"):]
		if !strings.HasPrefix(strings.TrimLeft(s, " \n\t\r"), "<a ") {
			return true
		}
	}
}

func TestTocMisc(t *testing.T) {
	c := qt.New(t)

	c.Run("Identifiers", func(c *qt.C) {
		toc := newTestToc()
		c.Assert(toc.Identifiers, qt.DeepEquals, collections.SortedStringSlice{"1-h2-1", "1-h2-2", "1-h2-2", "h1-1", "h1-2"})
	})

	c.Run("HeadingsMap", func(c *qt.C) {
		toc := newTestToc()
		m := toc.HeadingsMap
		c.Assert(m["h1-1"].Title, qt.Equals, "Heading 1")
		c.Assert(m["doesnot exist"], qt.IsNil)
	})
}

// Note that some of these cannot use b.Loop() because of golang/go#27217.
func BenchmarkToc(b *testing.B) {
	newTocs := func(n int) []*Fragments {
		var tocs []*Fragments
		for range n {
			tocs = append(tocs, newTestToc())
		}
		return tocs
	}

	b.Run("Build", func(b *testing.B) {
		var builders []Builder
		for i := 0; i < b.N; i++ {
			builders = append(builders, newTestTocBuilder())
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			b := builders[i]
			b.Build()
		}
	})

	b.Run("ToHTML", func(b *testing.B) {
		const size = 1000
		tocs := newTocs(size)
		for i := 0; b.Loop(); i++ {
			toc := tocs[i%size]
			toc.ToHTML(1, -1, false)
		}
	})
}
