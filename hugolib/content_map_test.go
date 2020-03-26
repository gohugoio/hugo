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

	"github.com/gohugoio/hugo/helpers"

	"github.com/gohugoio/hugo/htesting/hqt"

	"github.com/gohugoio/hugo/hugofs/files"

	"github.com/gohugoio/hugo/hugofs"
	"github.com/spf13/afero"

	qt "github.com/frankban/quicktest"
)

func BenchmarkContentMap(b *testing.B) {
	writeFile := func(c *qt.C, fs afero.Fs, filename, content string) hugofs.FileMetaInfo {
		c.Helper()
		filename = filepath.FromSlash(filename)
		c.Assert(fs.MkdirAll(filepath.Dir(filename), 0777), qt.IsNil)
		c.Assert(afero.WriteFile(fs, filename, []byte(content), 0777), qt.IsNil)

		fi, err := fs.Stat(filename)
		c.Assert(err, qt.IsNil)

		mfi := fi.(hugofs.FileMetaInfo)
		return mfi

	}

	createFs := func(fs afero.Fs, lang string) afero.Fs {
		return hugofs.NewBaseFileDecorator(fs,
			func(fi hugofs.FileMetaInfo) {
				meta := fi.Meta()
				// We have a more elaborate filesystem setup in the
				// real flow, so simulate this here.
				meta["lang"] = lang
				meta["path"] = meta.Filename()
				meta["classifier"] = files.ClassifyContentFile(fi.Name(), meta.GetOpener())

			})
	}

	b.Run("CreateMissingNodes", func(b *testing.B) {
		c := qt.New(b)
		b.StopTimer()
		mps := make([]*contentMap, b.N)
		for i := 0; i < b.N; i++ {
			m := newContentMap(contentMapConfig{lang: "en"})
			mps[i] = m
			memfs := afero.NewMemMapFs()
			fs := createFs(memfs, "en")
			for i := 1; i <= 20; i++ {
				c.Assert(m.AddFilesBundle(writeFile(c, fs, fmt.Sprintf("sect%d/a/index.md", i), "page")), qt.IsNil)
				c.Assert(m.AddFilesBundle(writeFile(c, fs, fmt.Sprintf("sect2%d/%sindex.md", i, strings.Repeat("b/", i)), "page")), qt.IsNil)
			}

		}

		b.StartTimer()

		for i := 0; i < b.N; i++ {
			m := mps[i]
			c.Assert(m.CreateMissingNodes(), qt.IsNil)

			b.StopTimer()
			m.pages.DeletePrefix("/")
			m.sections.DeletePrefix("/")
			b.StartTimer()
		}
	})

}

func TestContentMap(t *testing.T) {
	c := qt.New(t)

	writeFile := func(c *qt.C, fs afero.Fs, filename, content string) hugofs.FileMetaInfo {
		c.Helper()
		filename = filepath.FromSlash(filename)
		c.Assert(fs.MkdirAll(filepath.Dir(filename), 0777), qt.IsNil)
		c.Assert(afero.WriteFile(fs, filename, []byte(content), 0777), qt.IsNil)

		fi, err := fs.Stat(filename)
		c.Assert(err, qt.IsNil)

		mfi := fi.(hugofs.FileMetaInfo)
		return mfi

	}

	createFs := func(fs afero.Fs, lang string) afero.Fs {
		return hugofs.NewBaseFileDecorator(fs,
			func(fi hugofs.FileMetaInfo) {
				meta := fi.Meta()
				// We have a more elaborate filesystem setup in the
				// real flow, so simulate this here.
				meta["lang"] = lang
				meta["path"] = meta.Filename()
				meta["classifier"] = files.ClassifyContentFile(fi.Name(), meta.GetOpener())
				meta["translationBaseName"] = helpers.Filename(fi.Name())

			})
	}

	c.Run("AddFiles", func(c *qt.C) {

		memfs := afero.NewMemMapFs()

		fsl := func(lang string) afero.Fs {
			return createFs(memfs, lang)
		}

		fs := fsl("en")

		header := writeFile(c, fs, "blog/a/index.md", "page")

		c.Assert(header.Meta().Lang(), qt.Equals, "en")

		resources := []hugofs.FileMetaInfo{
			writeFile(c, fs, "blog/a/b/data.json", "data"),
			writeFile(c, fs, "blog/a/logo.png", "image"),
		}

		m := newContentMap(contentMapConfig{lang: "en"})

		c.Assert(m.AddFilesBundle(header, resources...), qt.IsNil)

		c.Assert(m.AddFilesBundle(writeFile(c, fs, "blog/b/c/index.md", "page")), qt.IsNil)

		c.Assert(m.AddFilesBundle(
			writeFile(c, fs, "blog/_index.md", "section page"),
			writeFile(c, fs, "blog/sectiondata.json", "section resource"),
		), qt.IsNil)

		got := m.testDump()

		expect := `
          Tree 0:
              	/blog__hb_/a__hl_
              	/blog__hb_/b/c__hl_
              Tree 1:
              	/blog
              Tree 2:
              	/blog__hb_/a__hl_b/data.json
              	/blog__hb_/a__hl_logo.png
              	/blog__hl_sectiondata.json
              en/pages/blog__hb_/a__hl_|f:blog/a/index.md
              	 - R: blog/a/b/data.json
              	 - R: blog/a/logo.png
              en/pages/blog__hb_/b/c__hl_|f:blog/b/c/index.md
              en/sections/blog|f:blog/_index.md
              	 - P: blog/a/index.md
              	 - P: blog/b/c/index.md
              	 - R: blog/sectiondata.json
    
`

		c.Assert(got, hqt.IsSameString, expect, qt.Commentf(got))

		// Add a data file to the section bundle
		c.Assert(m.AddFiles(
			writeFile(c, fs, "blog/sectiondata2.json", "section resource"),
		), qt.IsNil)

		// And then one to the leaf bundles
		c.Assert(m.AddFiles(
			writeFile(c, fs, "blog/a/b/data2.json", "data2"),
		), qt.IsNil)

		c.Assert(m.AddFiles(
			writeFile(c, fs, "blog/b/c/d/data3.json", "data3"),
		), qt.IsNil)

		got = m.testDump()

		expect = `
			 Tree 0:
              	/blog__hb_/a__hl_
              	/blog__hb_/b/c__hl_
              Tree 1:
              	/blog
              Tree 2:
              	/blog__hb_/a__hl_b/data.json
              	/blog__hb_/a__hl_b/data2.json
              	/blog__hb_/a__hl_logo.png
              	/blog__hb_/b/c__hl_d/data3.json
              	/blog__hl_sectiondata.json
              	/blog__hl_sectiondata2.json
              en/pages/blog__hb_/a__hl_|f:blog/a/index.md
              	 - R: blog/a/b/data.json
              	 - R: blog/a/b/data2.json
              	 - R: blog/a/logo.png
              en/pages/blog__hb_/b/c__hl_|f:blog/b/c/index.md
              	 - R: blog/b/c/d/data3.json
              en/sections/blog|f:blog/_index.md
              	 - P: blog/a/index.md
              	 - P: blog/b/c/index.md
              	 - R: blog/sectiondata.json
              	 - R: blog/sectiondata2.json
             
`

		c.Assert(got, hqt.IsSameString, expect, qt.Commentf(got))

		// Add a regular page (i.e. not a bundle)
		c.Assert(m.AddFilesBundle(writeFile(c, fs, "blog/b.md", "page")), qt.IsNil)

		c.Assert(m.testDump(), hqt.IsSameString, `
		 Tree 0:
              	/blog__hb_/a__hl_
              	/blog__hb_/b/c__hl_
              	/blog__hb_/b__hl_
              Tree 1:
              	/blog
              Tree 2:
              	/blog__hb_/a__hl_b/data.json
              	/blog__hb_/a__hl_b/data2.json
              	/blog__hb_/a__hl_logo.png
              	/blog__hb_/b/c__hl_d/data3.json
              	/blog__hl_sectiondata.json
              	/blog__hl_sectiondata2.json
              en/pages/blog__hb_/a__hl_|f:blog/a/index.md
              	 - R: blog/a/b/data.json
              	 - R: blog/a/b/data2.json
              	 - R: blog/a/logo.png
              en/pages/blog__hb_/b/c__hl_|f:blog/b/c/index.md
              	 - R: blog/b/c/d/data3.json
              en/pages/blog__hb_/b__hl_|f:blog/b.md
              en/sections/blog|f:blog/_index.md
              	 - P: blog/a/index.md
              	 - P: blog/b/c/index.md
              	 - P: blog/b.md
              	 - R: blog/sectiondata.json
              	 - R: blog/sectiondata2.json
             
       
				`, qt.Commentf(m.testDump()))

	})

	c.Run("CreateMissingNodes", func(c *qt.C) {

		memfs := afero.NewMemMapFs()

		fsl := func(lang string) afero.Fs {
			return createFs(memfs, lang)
		}

		fs := fsl("en")

		m := newContentMap(contentMapConfig{lang: "en"})

		c.Assert(m.AddFilesBundle(writeFile(c, fs, "blog/page.md", "page")), qt.IsNil)
		c.Assert(m.AddFilesBundle(writeFile(c, fs, "blog/a/index.md", "page")), qt.IsNil)
		c.Assert(m.AddFilesBundle(writeFile(c, fs, "bundle/index.md", "page")), qt.IsNil)

		c.Assert(m.CreateMissingNodes(), qt.IsNil)

		got := m.testDump()

		c.Assert(got, hqt.IsSameString, `
			
			 Tree 0:
              	/__hb_/bundle__hl_
              	/blog__hb_/a__hl_
              	/blog__hb_/page__hl_
              Tree 1:
              	/
              	/blog
              Tree 2:
              en/pages/__hb_/bundle__hl_|f:bundle/index.md
              en/pages/blog__hb_/a__hl_|f:blog/a/index.md
              en/pages/blog__hb_/page__hl_|f:blog/page.md
              en/sections/
              	 - P: bundle/index.md
              en/sections/blog
              	 - P: blog/a/index.md
              	 - P: blog/page.md
            
			`, qt.Commentf(got))

	})

	c.Run("cleanKey", func(c *qt.C) {
		for _, test := range []struct {
			in       string
			expected string
		}{
			{"/a/b/", "/a/b"},
			{filepath.FromSlash("/a/b/"), "/a/b"},
			{"/a//b/", "/a/b"},
		} {

			c.Assert(cleanTreeKey(test.in), qt.Equals, test.expected)

		}
	})
}

func TestContentMapSite(t *testing.T) {

	b := newTestSitesBuilder(t)

	pageTempl := `
---
title: "Page %d"
date: "2019-06-0%d"	
lastMod: "2019-06-0%d"
categories: ["funny"]
---

Page content.
`
	createPage := func(i int) string {
		return fmt.Sprintf(pageTempl, i, i, i+1)
	}

	draftTemplate := `---
title: "Draft"
draft: true
---

`

	b.WithContent("_index.md", `
---
title: "Hugo Home"
cascade:
    description: "Common Description"
    
---

Home Content.
`)

	b.WithContent("blog/page1.md", createPage(1))
	b.WithContent("blog/page2.md", createPage(2))
	b.WithContent("blog/page3.md", createPage(3))
	b.WithContent("blog/bundle/index.md", createPage(12))
	b.WithContent("blog/bundle/data.json", "data")
	b.WithContent("blog/bundle/page.md", createPage(99))
	b.WithContent("blog/subsection/_index.md", createPage(3))
	b.WithContent("blog/subsection/subdata.json", "data")
	b.WithContent("blog/subsection/page4.md", createPage(8))
	b.WithContent("blog/subsection/page5.md", createPage(10))
	b.WithContent("blog/subsection/draft/index.md", draftTemplate)
	b.WithContent("blog/subsection/draft/data.json", "data")
	b.WithContent("blog/draftsection/_index.md", draftTemplate)
	b.WithContent("blog/draftsection/page/index.md", createPage(12))
	b.WithContent("blog/draftsection/page/folder/data.json", "data")
	b.WithContent("blog/draftsection/sub/_index.md", createPage(12))
	b.WithContent("blog/draftsection/sub/page.md", createPage(13))
	b.WithContent("docs/page6.md", createPage(11))
	b.WithContent("tags/_index.md", createPage(32))
	b.WithContent("overlap/_index.md", createPage(33))
	b.WithContent("overlap2/_index.md", createPage(34))

	b.WithTemplatesAdded("layouts/index.html", `
Num Regular: {{ len .Site.RegularPages }}
Main Sections: {{ .Site.Params.mainSections }}
Pag Num Pages: {{ len .Paginator.Pages }}
{{ $home := .Site.Home }}
{{ $blog := .Site.GetPage "blog" }}
{{ $categories := .Site.GetPage "categories" }}
{{ $funny := .Site.GetPage "categories/funny" }}
{{ $blogSub := .Site.GetPage "blog/subsection" }}
{{ $page := .Site.GetPage "blog/page1" }}
{{ $page2 := .Site.GetPage "blog/page2" }}
{{ $page4 := .Site.GetPage "blog/subsection/page4" }}
{{ $bundle := .Site.GetPage "blog/bundle" }}
{{ $overlap1 := .Site.GetPage "overlap" }}
{{ $overlap2 := .Site.GetPage "overlap2" }}

Home: {{ template "print-page" $home }}
Blog Section: {{ template "print-page" $blog }}
Blog Sub Section: {{ template "print-page" $blogSub }}
Page: {{ template "print-page" $page }}
Bundle: {{ template "print-page" $bundle }}
IsDescendant: true: {{ $page.IsDescendant $blog }} true: {{ $blogSub.IsDescendant $blog }} true: {{ $blog.IsDescendant $home }} false: {{ $home.IsDescendant $blog }}
IsAncestor: true: {{ $blog.IsAncestor $page }} true: {{ $home.IsAncestor $blog }} true: {{ $blog.IsAncestor $blogSub }} true: {{ $home.IsAncestor $page }} false: {{ $page.IsAncestor $blog }} false: {{ $blog.IsAncestor $home }}  false: {{ $blogSub.IsAncestor $blog }}
IsDescendant overlap1: false: {{ $overlap1.IsDescendant $overlap2 }}
IsDescendant overlap2: false: {{ $overlap2.IsDescendant $overlap1 }}
IsAncestor overlap1: false: {{ $overlap1.IsAncestor $overlap2 }}
IsAncestor overlap2: false: {{ $overlap2.IsAncestor $overlap1 }}
FirstSection: {{ $blogSub.FirstSection.RelPermalink }} {{ $blog.FirstSection.RelPermalink }} {{ $home.FirstSection.RelPermalink }} {{ $page.FirstSection.RelPermalink }}
InSection: true: {{ $page.InSection $blog }} false: {{ $page.InSection $blogSub }} 
Next: {{ $page2.Next.RelPermalink }}
NextInSection: {{ $page2.NextInSection.RelPermalink }}
Pages: {{ range $blog.Pages }}{{ .RelPermalink }}|{{ end }}
Sections: {{ range $home.Sections }}{{ .RelPermalink }}|{{ end }}
Categories: {{ range .Site.Taxonomies.categories }}{{ .Page.RelPermalink }}; {{ .Page.Title }}; {{ .Count }}|{{ end }}
Category Terms:  {{ $categories.Kind}}: {{ range $categories.Data.Terms.Alphabetical }}{{ .Page.RelPermalink }}; {{ .Page.Title }}; {{ .Count }}|{{ end }}
Category Funny:  {{ $funny.Kind}}; {{ $funny.Data.Term }}: {{ range $funny.Pages }}{{ .RelPermalink }};|{{ end }}
Pag Num Pages: {{ len .Paginator.Pages }}
Pag Blog Num Pages: {{ len $blog.Paginator.Pages }}
Blog Num RegularPages: {{ len $blog.RegularPages }}
Blog Num Pages: {{ len $blog.Pages }}

Draft1: {{ if (.Site.GetPage "blog/subsection/draft") }}FOUND{{ end }}|
Draft2: {{ if (.Site.GetPage "blog/draftsection") }}FOUND{{ end }}|
Draft3: {{ if (.Site.GetPage "blog/draftsection/page") }}FOUND{{ end }}|
Draft4: {{ if (.Site.GetPage "blog/draftsection/sub") }}FOUND{{ end }}|
Draft5: {{ if (.Site.GetPage "blog/draftsection/sub/page") }}FOUND{{ end }}|

{{ define "print-page" }}{{ .Title }}|{{ .RelPermalink }}|{{ .Date.Format "2006-01-02" }}|Current Section: {{ .CurrentSection.SectionsPath }}|Resources: {{ range .Resources }}{{ .ResourceType }}: {{ .RelPermalink }}|{{ end }}{{ end }}
`)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/index.html",

		`
	 Num Regular: 7
        Main Sections: [blog]
        Pag Num Pages: 7
        
      Home: Hugo Home|/|2019-06-08|Current Section: |Resources: 
        Blog Section: Blogs|/blog/|2019-06-08|Current Section: blog|Resources: 
        Blog Sub Section: Page 3|/blog/subsection/|2019-06-03|Current Section: blog/subsection|Resources: json: /blog/subsection/subdata.json|
        Page: Page 1|/blog/page1/|2019-06-01|Current Section: blog|Resources: 
        Bundle: Page 12|/blog/bundle/|0001-01-01|Current Section: blog|Resources: json: /blog/bundle/data.json|page: |
        IsDescendant: true: true true: true true: true false: false
        IsAncestor: true: true true: true true: true true: true false: false false: false  false: false
        IsDescendant overlap1: false: false
        IsDescendant overlap2: false: false
        IsAncestor overlap1: false: false
        IsAncestor overlap2: false: false
        FirstSection: /blog/ /blog/ / /blog/
        InSection: true: true false: false 
        Next: /blog/page3/
        NextInSection: /blog/page3/
        Pages: /blog/page3/|/blog/subsection/|/blog/page2/|/blog/page1/|/blog/bundle/|
        Sections: /blog/|/docs/|
        Categories: /categories/funny/; funny; 11|
        Category Terms:  taxonomyTerm: /categories/funny/; funny; 11|
 		Category Funny:  taxonomy; funny: /blog/subsection/page4/;|/blog/page3/;|/blog/subsection/;|/blog/page2/;|/blog/page1/;|/blog/subsection/page5/;|/docs/page6/;|/blog/bundle/;|;|
 		Pag Num Pages: 7
        Pag Blog Num Pages: 4
        Blog Num RegularPages: 4
        Blog Num Pages: 5
        
        Draft1: |
        Draft2: |
        Draft3: |
        Draft4: |
        Draft5: |
           
`)
}
