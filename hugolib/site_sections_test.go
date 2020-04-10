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
	"path/filepath"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/resources/page"
)

func TestNestedSections(t *testing.T) {

	var (
		c       = qt.New(t)
		cfg, fs = newTestCfg()
		th      = newTestHelper(cfg, fs, t)
	)

	cfg.Set("permalinks", map[string]string{
		"perm a": ":sections/:title",
	})

	pageTemplate := `---
title: T%d_%d
---
Content
`

	// Home page
	writeSource(t, fs, filepath.Join("content", "_index.md"), fmt.Sprintf(pageTemplate, -1, -1))

	// Top level content page
	writeSource(t, fs, filepath.Join("content", "mypage.md"), fmt.Sprintf(pageTemplate, 1234, 5))

	// Top level section without index content page
	writeSource(t, fs, filepath.Join("content", "top", "mypage2.md"), fmt.Sprintf(pageTemplate, 12345, 6))
	// Just a page in a subfolder, i.e. not a section.
	writeSource(t, fs, filepath.Join("content", "top", "folder", "mypage3.md"), fmt.Sprintf(pageTemplate, 12345, 67))

	for level1 := 1; level1 < 3; level1++ {
		writeSource(t, fs, filepath.Join("content", "l1", fmt.Sprintf("page_1_%d.md", level1)),
			fmt.Sprintf(pageTemplate, 1, level1))
	}

	// Issue #3586
	writeSource(t, fs, filepath.Join("content", "post", "0000.md"), fmt.Sprintf(pageTemplate, 1, 2))
	writeSource(t, fs, filepath.Join("content", "post", "0000", "0001.md"), fmt.Sprintf(pageTemplate, 1, 3))
	writeSource(t, fs, filepath.Join("content", "elsewhere", "0003.md"), fmt.Sprintf(pageTemplate, 1, 4))

	// Empty nested section, i.e. no regular content pages.
	writeSource(t, fs, filepath.Join("content", "empty1", "b", "c", "_index.md"), fmt.Sprintf(pageTemplate, 33, -1))
	// Index content file a the end and in the middle.
	writeSource(t, fs, filepath.Join("content", "empty2", "b", "_index.md"), fmt.Sprintf(pageTemplate, 40, -1))
	writeSource(t, fs, filepath.Join("content", "empty2", "b", "c", "d", "_index.md"), fmt.Sprintf(pageTemplate, 41, -1))

	// Empty with content file in the middle.
	writeSource(t, fs, filepath.Join("content", "empty3", "b", "c", "d", "_index.md"), fmt.Sprintf(pageTemplate, 41, -1))
	writeSource(t, fs, filepath.Join("content", "empty3", "b", "empty3.md"), fmt.Sprintf(pageTemplate, 3, -1))

	// Section with permalink config
	writeSource(t, fs, filepath.Join("content", "perm a", "link", "_index.md"), fmt.Sprintf(pageTemplate, 9, -1))
	for i := 1; i < 4; i++ {
		writeSource(t, fs, filepath.Join("content", "perm a", "link", fmt.Sprintf("page_%d.md", i)),
			fmt.Sprintf(pageTemplate, 1, i))
	}
	writeSource(t, fs, filepath.Join("content", "perm a", "link", "regular", fmt.Sprintf("page_%d.md", 5)),
		fmt.Sprintf(pageTemplate, 1, 5))

	writeSource(t, fs, filepath.Join("content", "l1", "l2", "_index.md"), fmt.Sprintf(pageTemplate, 2, -1))
	writeSource(t, fs, filepath.Join("content", "l1", "l2_2", "_index.md"), fmt.Sprintf(pageTemplate, 22, -1))
	writeSource(t, fs, filepath.Join("content", "l1", "l2", "l3", "_index.md"), fmt.Sprintf(pageTemplate, 3, -1))

	for level2 := 1; level2 < 4; level2++ {
		writeSource(t, fs, filepath.Join("content", "l1", "l2", fmt.Sprintf("page_2_%d.md", level2)),
			fmt.Sprintf(pageTemplate, 2, level2))
	}
	for level2 := 1; level2 < 3; level2++ {
		writeSource(t, fs, filepath.Join("content", "l1", "l2_2", fmt.Sprintf("page_2_2_%d.md", level2)),
			fmt.Sprintf(pageTemplate, 2, level2))
	}
	for level3 := 1; level3 < 3; level3++ {
		writeSource(t, fs, filepath.Join("content", "l1", "l2", "l3", fmt.Sprintf("page_3_%d.md", level3)),
			fmt.Sprintf(pageTemplate, 3, level3))
	}

	writeSource(t, fs, filepath.Join("content", "Spaces in Section", "page100.md"), fmt.Sprintf(pageTemplate, 10, 0))

	writeSource(t, fs, filepath.Join("layouts", "_default", "single.html"), "<html>Single|{{ .Title }}</html>")
	writeSource(t, fs, filepath.Join("layouts", "_default", "list.html"),
		`
{{ $sect := (.Site.GetPage "l1/l2") }}
<html>List|{{ .Title }}|L1/l2-IsActive: {{ .InSection $sect }}
{{ range .Paginator.Pages }}
PAG|{{ .Title }}|{{ $sect.InSection . }}
{{ end }}
{{/* https://github.com/gohugoio/hugo/issues/4989 */}}
{{ $sections := (.Site.GetPage "section" .Section).Sections.ByWeight }}
</html>`)

	cfg.Set("paginate", 2)

	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{})

	c.Assert(len(s.RegularPages()), qt.Equals, 21)

	tests := []struct {
		sections string
		verify   func(c *qt.C, p page.Page)
	}{
		{"elsewhere", func(c *qt.C, p page.Page) {
			c.Assert(len(p.Pages()), qt.Equals, 1)
			for _, p := range p.Pages() {
				c.Assert(p.SectionsPath(), qt.Equals, "elsewhere")
			}
		}},
		{"post", func(c *qt.C, p page.Page) {
			c.Assert(len(p.Pages()), qt.Equals, 2)
			for _, p := range p.Pages() {
				c.Assert(p.Section(), qt.Equals, "post")
			}
		}},
		{"empty1", func(c *qt.C, p page.Page) {
			// > b,c
			c.Assert(getPage(p, "/empty1/b"), qt.IsNil) // No _index.md page.
			c.Assert(getPage(p, "/empty1/b/c"), qt.Not(qt.IsNil))

		}},
		{"empty2", func(c *qt.C, p page.Page) {
			// > b,c,d where b and d have _index.md files.
			b := getPage(p, "/empty2/b")
			c.Assert(b, qt.Not(qt.IsNil))
			c.Assert(b.Title(), qt.Equals, "T40_-1")

			cp := getPage(p, "/empty2/b/c")
			c.Assert(cp, qt.IsNil) // No _index.md

			d := getPage(p, "/empty2/b/c/d")
			c.Assert(d, qt.Not(qt.IsNil))
			c.Assert(d.Title(), qt.Equals, "T41_-1")

			c.Assert(cp.Eq(d), qt.Equals, false)
			c.Assert(cp.Eq(cp), qt.Equals, true)
			c.Assert(cp.Eq("asdf"), qt.Equals, false)

		}},
		{"empty3", func(c *qt.C, p page.Page) {
			// b,c,d with regular page in b
			b := getPage(p, "/empty3/b")
			c.Assert(b, qt.IsNil) // No _index.md
			e3 := getPage(p, "/empty3/b/empty3")
			c.Assert(e3, qt.Not(qt.IsNil))
			c.Assert(e3.File().LogicalName(), qt.Equals, "empty3.md")

		}},
		{"empty3", func(c *qt.C, p page.Page) {
			xxx := getPage(p, "/empty3/nil")
			c.Assert(xxx, qt.IsNil)
		}},
		{"top", func(c *qt.C, p page.Page) {
			c.Assert(p.Title(), qt.Equals, "Tops")
			c.Assert(len(p.Pages()), qt.Equals, 2)
			c.Assert(p.Pages()[0].File().LogicalName(), qt.Equals, "mypage2.md")
			c.Assert(p.Pages()[1].File().LogicalName(), qt.Equals, "mypage3.md")
			home := p.Parent()
			c.Assert(home.IsHome(), qt.Equals, true)
			c.Assert(len(p.Sections()), qt.Equals, 0)
			c.Assert(home.CurrentSection(), qt.Equals, home)
			active, err := home.InSection(home)
			c.Assert(err, qt.IsNil)
			c.Assert(active, qt.Equals, true)
			c.Assert(p.FirstSection(), qt.Equals, p)
		}},
		{"l1", func(c *qt.C, p page.Page) {
			c.Assert(p.Title(), qt.Equals, "L1s")
			c.Assert(len(p.Pages()), qt.Equals, 4) // 2 pages + 2 sections
			c.Assert(p.Parent().IsHome(), qt.Equals, true)
			c.Assert(len(p.Sections()), qt.Equals, 2)
		}},
		{"l1,l2", func(c *qt.C, p page.Page) {
			c.Assert(p.Title(), qt.Equals, "T2_-1")
			c.Assert(len(p.Pages()), qt.Equals, 4) // 3 pages + 1 section
			c.Assert(p.Pages()[0].Parent(), qt.Equals, p)
			c.Assert(p.Parent().Title(), qt.Equals, "L1s")
			c.Assert(p.RelPermalink(), qt.Equals, "/l1/l2/")
			c.Assert(len(p.Sections()), qt.Equals, 1)

			for _, child := range p.Pages() {
				if child.IsSection() {
					c.Assert(child.CurrentSection(), qt.Equals, child)
					continue
				}

				c.Assert(child.CurrentSection(), qt.Equals, p)
				active, err := child.InSection(p)
				c.Assert(err, qt.IsNil)

				c.Assert(active, qt.Equals, true)
				active, err = p.InSection(child)
				c.Assert(err, qt.IsNil)
				c.Assert(active, qt.Equals, true)
				active, err = p.InSection(getPage(p, "/"))
				c.Assert(err, qt.IsNil)
				c.Assert(active, qt.Equals, false)

				isAncestor, err := p.IsAncestor(child)
				c.Assert(err, qt.IsNil)
				c.Assert(isAncestor, qt.Equals, true)
				isAncestor, err = child.IsAncestor(p)
				c.Assert(err, qt.IsNil)
				c.Assert(isAncestor, qt.Equals, false)

				isDescendant, err := p.IsDescendant(child)
				c.Assert(err, qt.IsNil)
				c.Assert(isDescendant, qt.Equals, false)
				isDescendant, err = child.IsDescendant(p)
				c.Assert(err, qt.IsNil)
				c.Assert(isDescendant, qt.Equals, true)
			}

			c.Assert(p.Eq(p.CurrentSection()), qt.Equals, true)

		}},
		{"l1,l2_2", func(c *qt.C, p page.Page) {
			c.Assert(p.Title(), qt.Equals, "T22_-1")
			c.Assert(len(p.Pages()), qt.Equals, 2)
			c.Assert(p.Pages()[0].File().Path(), qt.Equals, filepath.FromSlash("l1/l2_2/page_2_2_1.md"))
			c.Assert(p.Parent().Title(), qt.Equals, "L1s")
			c.Assert(len(p.Sections()), qt.Equals, 0)
		}},
		{"l1,l2,l3", func(c *qt.C, p page.Page) {
			nilp, _ := p.GetPage("this/does/not/exist")

			c.Assert(p.Title(), qt.Equals, "T3_-1")
			c.Assert(len(p.Pages()), qt.Equals, 2)
			c.Assert(p.Parent().Title(), qt.Equals, "T2_-1")
			c.Assert(len(p.Sections()), qt.Equals, 0)

			l1 := getPage(p, "/l1")
			isDescendant, err := l1.IsDescendant(p)
			c.Assert(err, qt.IsNil)
			c.Assert(isDescendant, qt.Equals, false)
			isDescendant, err = l1.IsDescendant(nil)
			c.Assert(err, qt.IsNil)
			c.Assert(isDescendant, qt.Equals, false)
			isDescendant, err = nilp.IsDescendant(p)
			c.Assert(err, qt.IsNil)
			c.Assert(isDescendant, qt.Equals, false)
			isDescendant, err = p.IsDescendant(l1)
			c.Assert(err, qt.IsNil)
			c.Assert(isDescendant, qt.Equals, true)

			isAncestor, err := l1.IsAncestor(p)
			c.Assert(err, qt.IsNil)
			c.Assert(isAncestor, qt.Equals, true)
			isAncestor, err = p.IsAncestor(l1)
			c.Assert(err, qt.IsNil)
			c.Assert(isAncestor, qt.Equals, false)
			c.Assert(p.FirstSection(), qt.Equals, l1)
			isAncestor, err = p.IsAncestor(nil)
			c.Assert(err, qt.IsNil)
			c.Assert(isAncestor, qt.Equals, false)
			isAncestor, err = nilp.IsAncestor(l1)
			c.Assert(err, qt.IsNil)
			c.Assert(isAncestor, qt.Equals, false)

		}},
		{"perm a,link", func(c *qt.C, p page.Page) {
			c.Assert(p.Title(), qt.Equals, "T9_-1")
			c.Assert(p.RelPermalink(), qt.Equals, "/perm-a/link/")
			c.Assert(len(p.Pages()), qt.Equals, 4)
			first := p.Pages()[0]
			c.Assert(first.RelPermalink(), qt.Equals, "/perm-a/link/t1_1/")
			th.assertFileContent("public/perm-a/link/t1_1/index.html", "Single|T1_1")

			last := p.Pages()[3]
			c.Assert(last.RelPermalink(), qt.Equals, "/perm-a/link/t1_5/")

		}},
	}

	home := s.getPage(page.KindHome)

	for _, test := range tests {
		test := test
		t.Run(fmt.Sprintf("sections %s", test.sections), func(t *testing.T) {
			t.Parallel()
			c := qt.New(t)
			sections := strings.Split(test.sections, ",")
			p := s.getPage(page.KindSection, sections...)
			c.Assert(p, qt.Not(qt.IsNil), qt.Commentf(fmt.Sprint(sections)))

			if p.Pages() != nil {
				c.Assert(p.Data().(page.Data).Pages(), deepEqualsPages, p.Pages())
			}
			c.Assert(p.Parent(), qt.Not(qt.IsNil))
			test.verify(c, p)
		})
	}

	c.Assert(home, qt.Not(qt.IsNil))

	c.Assert(len(home.Sections()), qt.Equals, 9)
	c.Assert(s.Info.Sections(), deepEqualsPages, home.Sections())

	rootPage := s.getPage(page.KindPage, "mypage.md")
	c.Assert(rootPage, qt.Not(qt.IsNil))
	c.Assert(rootPage.Parent().IsHome(), qt.Equals, true)
	// https://github.com/gohugoio/hugo/issues/6365
	c.Assert(rootPage.Sections(), qt.HasLen, 0)

	// Add a odd test for this as this looks a little bit off, but I'm not in the mood
	// to think too hard a out this right now. It works, but people will have to spell
	// out the directory name as is.
	// If we later decide to do something about this, we will have to do some normalization in
	// getPage.
	// TODO(bep)
	sectionWithSpace := s.getPage(page.KindSection, "Spaces in Section")
	c.Assert(sectionWithSpace, qt.Not(qt.IsNil))
	c.Assert(sectionWithSpace.RelPermalink(), qt.Equals, "/spaces-in-section/")

	th.assertFileContent("public/l1/l2/page/2/index.html", "L1/l2-IsActive: true", "PAG|T2_3|true")

}

func TestNextInSectionNested(t *testing.T) {
	t.Parallel()

	pageContent := `---
title: "The Page"
weight: %d
---
Some content.
`
	createPageContent := func(weight int) string {
		return fmt.Sprintf(pageContent, weight)
	}

	b := newTestSitesBuilder(t)
	b.WithSimpleConfigFile()
	b.WithTemplates("_default/single.html", `
Prev: {{ with .PrevInSection }}{{ .RelPermalink }}{{ end }}|
Next: {{ with .NextInSection }}{{ .RelPermalink }}{{ end }}|
`)

	b.WithContent("blog/page1.md", createPageContent(1))
	b.WithContent("blog/page2.md", createPageContent(2))
	b.WithContent("blog/cool/_index.md", createPageContent(1))
	b.WithContent("blog/cool/cool1.md", createPageContent(1))
	b.WithContent("blog/cool/cool2.md", createPageContent(2))
	b.WithContent("root1.md", createPageContent(1))
	b.WithContent("root2.md", createPageContent(2))

	b.Build(BuildCfg{})

	b.AssertFileContent("public/root1/index.html",
		"Prev: /root2/|", "Next: |")
	b.AssertFileContent("public/root2/index.html",
		"Prev: |", "Next: /root1/|")
	b.AssertFileContent("public/blog/page1/index.html",
		"Prev: /blog/page2/|", "Next: |")
	b.AssertFileContent("public/blog/page2/index.html",
		"Prev: |", "Next: /blog/page1/|")
	b.AssertFileContent("public/blog/cool/cool1/index.html",
		"Prev: /blog/cool/cool2/|", "Next: |")
	b.AssertFileContent("public/blog/cool/cool2/index.html",
		"Prev: |", "Next: /blog/cool/cool1/|")

}
