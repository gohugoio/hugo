// Copyright 2017-present The Hugo Authors. All rights reserved.
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

	"github.com/gohugoio/hugo/deps"
	"github.com/stretchr/testify/require"
)

func TestNestedSections(t *testing.T) {
	t.Parallel()

	var (
		assert  = require.New(t)
		cfg, fs = newTestCfg()
		th      = testHelper{cfg, fs, t}
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
{{ $sect := (.Site.GetPage "section" "l1" "l2") }}
<html>List|{{ .Title }}|L1/l2-IsActive: {{ .InSection $sect }}
{{ range .Paginator.Pages }}
PAG|{{ .Title }}|{{ $sect.InSection . }}
{{ end }}
</html>`)

	cfg.Set("paginate", 2)

	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{})
	require.Len(t, s.RegularPages, 21)

	tests := []struct {
		sections string
		verify   func(p *Page)
	}{
		{"elsewhere", func(p *Page) {
			assert.Len(p.Pages, 1)
			for _, p := range p.Pages {
				assert.Equal([]string{"elsewhere"}, p.sections)
			}
		}},
		{"post", func(p *Page) {
			assert.Len(p.Pages, 2)
			for _, p := range p.Pages {
				assert.Equal("post", p.Section())
			}
		}},
		{"empty1", func(p *Page) {
			// > b,c
			assert.NotNil(p.s.getPage(KindSection, "empty1", "b"))
			assert.NotNil(p.s.getPage(KindSection, "empty1", "b", "c"))

		}},
		{"empty2", func(p *Page) {
			// > b,c,d where b and d have content files.
			b := p.s.getPage(KindSection, "empty2", "b")
			assert.NotNil(b)
			assert.Equal("T40_-1", b.Title)
			c := p.s.getPage(KindSection, "empty2", "b", "c")
			assert.NotNil(c)
			assert.Equal("Cs", c.Title)
			d := p.s.getPage(KindSection, "empty2", "b", "c", "d")
			assert.NotNil(d)
			assert.Equal("T41_-1", d.Title)

			assert.False(c.Eq(d))
			assert.True(c.Eq(c))
			assert.False(c.Eq("asdf"))

		}},
		{"empty3", func(p *Page) {
			// b,c,d with regular page in b
			b := p.s.getPage(KindSection, "empty3", "b")
			assert.NotNil(b)
			assert.Len(b.Pages, 1)
			assert.Equal("empty3.md", b.Pages[0].File.LogicalName())

		}},
		{"top", func(p *Page) {
			assert.Equal("Tops", p.Title)
			assert.Len(p.Pages, 2)
			assert.Equal("mypage2.md", p.Pages[0].LogicalName())
			assert.Equal("mypage3.md", p.Pages[1].LogicalName())
			home := p.Parent()
			assert.True(home.IsHome())
			assert.Len(p.Sections(), 0)
			assert.Equal(home, home.CurrentSection())
			active, err := home.InSection(home)
			assert.NoError(err)
			assert.True(active)
		}},
		{"l1", func(p *Page) {
			assert.Equal("L1s", p.Title)
			assert.Len(p.Pages, 2)
			assert.True(p.Parent().IsHome())
			assert.Len(p.Sections(), 2)
		}},
		{"l1,l2", func(p *Page) {
			assert.Equal("T2_-1", p.Title)
			assert.Len(p.Pages, 3)
			assert.Equal(p, p.Pages[0].Parent())
			assert.Equal("L1s", p.Parent().Title)
			assert.Equal("/l1/l2/", p.URLPath.URL)
			assert.Equal("/l1/l2/", p.RelPermalink())
			assert.Len(p.Sections(), 1)

			for _, child := range p.Pages {
				assert.Equal(p, child.CurrentSection())
				active, err := child.InSection(p)
				assert.NoError(err)
				assert.True(active)
				active, err = p.InSection(child)
				assert.NoError(err)
				assert.True(active)
				active, err = p.InSection(p.s.getPage(KindHome))
				assert.NoError(err)
				assert.False(active)

				isAncestor, err := p.IsAncestor(child)
				assert.NoError(err)
				assert.True(isAncestor)
				isAncestor, err = child.IsAncestor(p)
				assert.NoError(err)
				assert.False(isAncestor)

				isDescendant, err := p.IsDescendant(child)
				assert.NoError(err)
				assert.False(isDescendant)
				isDescendant, err = child.IsDescendant(p)
				assert.NoError(err)
				assert.True(isDescendant)
			}

			assert.Equal(p, p.CurrentSection())

		}},
		{"l1,l2_2", func(p *Page) {
			assert.Equal("T22_-1", p.Title)
			assert.Len(p.Pages, 2)
			assert.Equal(filepath.FromSlash("l1/l2_2/page_2_2_1.md"), p.Pages[0].Path())
			assert.Equal("L1s", p.Parent().Title)
			assert.Len(p.Sections(), 0)
		}},
		{"l1,l2,l3", func(p *Page) {
			assert.Equal("T3_-1", p.Title)
			assert.Len(p.Pages, 2)
			assert.Equal("T2_-1", p.Parent().Title)
			assert.Len(p.Sections(), 0)

			l1 := p.s.getPage(KindSection, "l1")
			isDescendant, err := l1.IsDescendant(p)
			assert.NoError(err)
			assert.False(isDescendant)
			isDescendant, err = p.IsDescendant(l1)
			assert.NoError(err)
			assert.True(isDescendant)

			isAncestor, err := l1.IsAncestor(p)
			assert.NoError(err)
			assert.True(isAncestor)
			isAncestor, err = p.IsAncestor(l1)
			assert.NoError(err)
			assert.False(isAncestor)

		}},
		{"perm a,link", func(p *Page) {
			assert.Equal("T9_-1", p.Title)
			assert.Equal("/perm-a/link/", p.RelPermalink())
			assert.Len(p.Pages, 4)
			first := p.Pages[0]
			assert.Equal("/perm-a/link/t1_1/", first.RelPermalink())
			th.assertFileContent("public/perm-a/link/t1_1/index.html", "Single|T1_1")

			last := p.Pages[3]
			assert.Equal("/perm-a/link/t1_5/", last.RelPermalink())

		}},
	}

	for _, test := range tests {
		sections := strings.Split(test.sections, ",")
		p := s.getPage(KindSection, sections...)
		assert.NotNil(p, fmt.Sprint(sections))

		if p.Pages != nil {
			assert.Equal(p.Pages, p.Data["Pages"])
		}
		assert.NotNil(p.Parent(), fmt.Sprintf("Parent nil: %q", test.sections))
		test.verify(p)
	}

	home := s.getPage(KindHome)

	assert.NotNil(home)

	assert.Len(home.Sections(), 9)
	assert.Equal(home.Sections(), s.Info.Sections())

	rootPage := s.getPage(KindPage, "mypage.md")
	assert.NotNil(rootPage)
	assert.True(rootPage.Parent().IsHome())

	// Add a odd test for this as this looks a little bit off, but I'm not in the mood
	// to think too hard a out this right now. It works, but people will have to spell
	// out the directory name as is.
	// If we later decide to do something about this, we will have to do some normalization in
	// getPage.
	// TODO(bep)
	sectionWithSpace := s.getPage(KindSection, "Spaces in Section")
	require.NotNil(t, sectionWithSpace)
	require.Equal(t, "/spaces-in-section/", sectionWithSpace.RelPermalink())

	th.assertFileContent("public/l1/l2/page/2/index.html", "L1/l2-IsActive: true", "PAG|T2_3|true")

}
