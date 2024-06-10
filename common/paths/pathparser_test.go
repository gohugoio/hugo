// Copyright 2024 The Hugo Authors. All rights reserved.
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

package paths

import (
	"path/filepath"
	"testing"

	"github.com/gohugoio/hugo/hugofs/files"

	qt "github.com/frankban/quicktest"
)

var testParser = &PathParser{
	LanguageIndex: map[string]int{
		"no": 0,
		"en": 1,
	},
	IsContentExt: func(ext string) bool {
		return ext == "md"
	},
}

func TestParse(t *testing.T) {
	c := qt.New(t)

	tests := []struct {
		name   string
		path   string
		assert func(c *qt.C, p *Path)
	}{
		{
			"Basic text file",
			"/a/b.txt",
			func(c *qt.C, p *Path) {
				c.Assert(p.Name(), qt.Equals, "b.txt")
				c.Assert(p.Base(), qt.Equals, "/a/b.txt")
				c.Assert(p.Container(), qt.Equals, "a")
				c.Assert(p.Dir(), qt.Equals, "/a")
				c.Assert(p.Ext(), qt.Equals, "txt")
				c.Assert(p.IsContent(), qt.IsFalse)
			},
		},
		{
			"Basic text file, upper case",
			"/A/B.txt",
			func(c *qt.C, p *Path) {
				c.Assert(p.Name(), qt.Equals, "b.txt")
				c.Assert(p.NameNoExt(), qt.Equals, "b")
				c.Assert(p.NameNoIdentifier(), qt.Equals, "b")
				c.Assert(p.BaseNameNoIdentifier(), qt.Equals, "b")
				c.Assert(p.Base(), qt.Equals, "/a/b.txt")
				c.Assert(p.Ext(), qt.Equals, "txt")
			},
		},
		{
			"Basic text file, 1 space in dir",
			"/a b/c.txt",
			func(c *qt.C, p *Path) {
				c.Assert(p.Base(), qt.Equals, "/a-b/c.txt")
			},
		},
		{
			"Basic text file, 2 spaces in dir",
			"/a  b/c.txt",
			func(c *qt.C, p *Path) {
				c.Assert(p.Base(), qt.Equals, "/a--b/c.txt")
			},
		},
		{
			"Basic text file, 1 space in filename",
			"/a/b c.txt",
			func(c *qt.C, p *Path) {
				c.Assert(p.Base(), qt.Equals, "/a/b-c.txt")
			},
		},
		{
			"Basic text file, 2 spaces in filename",
			"/a/b  c.txt",
			func(c *qt.C, p *Path) {
				c.Assert(p.Base(), qt.Equals, "/a/b--c.txt")
			},
		},
		{
			"Basic text file, mixed case and spaces, unnormalized",
			"/a/Foo BAR.txt",
			func(c *qt.C, p *Path) {
				pp := p.Unnormalized()
				c.Assert(pp, qt.IsNotNil)
				c.Assert(pp.BaseNameNoIdentifier(), qt.Equals, "Foo BAR")
			},
		},
		{
			"Basic Markdown file",
			"/a/b/c.md",
			func(c *qt.C, p *Path) {
				c.Assert(p.IsContent(), qt.IsTrue)
				c.Assert(p.IsLeafBundle(), qt.IsFalse)
				c.Assert(p.Name(), qt.Equals, "c.md")
				c.Assert(p.Base(), qt.Equals, "/a/b/c")
				c.Assert(p.Section(), qt.Equals, "a")
				c.Assert(p.BaseNameNoIdentifier(), qt.Equals, "c")
				c.Assert(p.Path(), qt.Equals, "/a/b/c.md")
				c.Assert(p.Dir(), qt.Equals, "/a/b")
				c.Assert(p.Container(), qt.Equals, "b")
				c.Assert(p.ContainerDir(), qt.Equals, "/a/b")
				c.Assert(p.Ext(), qt.Equals, "md")
			},
		},
		{
			"Content resource",
			"/a/b.md",
			func(c *qt.C, p *Path) {
				c.Assert(p.Name(), qt.Equals, "b.md")
				c.Assert(p.Base(), qt.Equals, "/a/b")
				c.Assert(p.BaseNoLeadingSlash(), qt.Equals, "a/b")
				c.Assert(p.Section(), qt.Equals, "a")
				c.Assert(p.BaseNameNoIdentifier(), qt.Equals, "b")

				// Reclassify it as a content resource.
				ModifyPathBundleTypeResource(p)
				c.Assert(p.BundleType(), qt.Equals, PathTypeContentResource)
				c.Assert(p.IsContent(), qt.IsTrue)
				c.Assert(p.Name(), qt.Equals, "b.md")
				c.Assert(p.Base(), qt.Equals, "/a/b.md")
			},
		},
		{
			"No ext",
			"/a/b",
			func(c *qt.C, p *Path) {
				c.Assert(p.Name(), qt.Equals, "b")
				c.Assert(p.NameNoExt(), qt.Equals, "b")
				c.Assert(p.Base(), qt.Equals, "/a/b")
				c.Assert(p.Ext(), qt.Equals, "")
			},
		},
		{
			"No ext, trailing slash",
			"/a/b/",
			func(c *qt.C, p *Path) {
				c.Assert(p.Name(), qt.Equals, "b")
				c.Assert(p.Base(), qt.Equals, "/a/b")
				c.Assert(p.Ext(), qt.Equals, "")
			},
		},
		{
			"Identifiers",
			"/a/b.a.b.no.txt",
			func(c *qt.C, p *Path) {
				c.Assert(p.Name(), qt.Equals, "b.a.b.no.txt")
				c.Assert(p.NameNoIdentifier(), qt.Equals, "b.a.b")
				c.Assert(p.NameNoLang(), qt.Equals, "b.a.b.txt")
				c.Assert(p.Identifiers(), qt.DeepEquals, []string{"txt", "no"})
				c.Assert(p.Base(), qt.Equals, "/a/b.a.b.txt")
				c.Assert(p.BaseNoLeadingSlash(), qt.Equals, "a/b.a.b.txt")
				c.Assert(p.PathNoLang(), qt.Equals, "/a/b.a.b.txt")
				c.Assert(p.Ext(), qt.Equals, "txt")
				c.Assert(p.PathNoIdentifier(), qt.Equals, "/a/b.a.b")
			},
		},
		{
			"Home branch cundle",
			"/_index.md",
			func(c *qt.C, p *Path) {
				c.Assert(p.Base(), qt.Equals, "/")
				c.Assert(p.Path(), qt.Equals, "/_index.md")
				c.Assert(p.Container(), qt.Equals, "")
				c.Assert(p.ContainerDir(), qt.Equals, "/")
			},
		},
		{
			"Index content file in root",
			"/a/index.md",
			func(c *qt.C, p *Path) {
				c.Assert(p.Base(), qt.Equals, "/a")
				c.Assert(p.BaseNameNoIdentifier(), qt.Equals, "a")
				c.Assert(p.Container(), qt.Equals, "a")
				c.Assert(p.Container(), qt.Equals, "a")
				c.Assert(p.ContainerDir(), qt.Equals, "")
				c.Assert(p.Dir(), qt.Equals, "/a")
				c.Assert(p.Ext(), qt.Equals, "md")
				c.Assert(p.Identifiers(), qt.DeepEquals, []string{"md"})
				c.Assert(p.IsBranchBundle(), qt.IsFalse)
				c.Assert(p.IsBundle(), qt.IsTrue)
				c.Assert(p.IsLeafBundle(), qt.IsTrue)
				c.Assert(p.Lang(), qt.Equals, "")
				c.Assert(p.NameNoExt(), qt.Equals, "index")
				c.Assert(p.NameNoIdentifier(), qt.Equals, "index")
				c.Assert(p.NameNoLang(), qt.Equals, "index.md")
				c.Assert(p.Section(), qt.Equals, "")
			},
		},
		{
			"Index content file with lang",
			"/a/b/index.no.md",
			func(c *qt.C, p *Path) {
				c.Assert(p.Base(), qt.Equals, "/a/b")
				c.Assert(p.BaseNameNoIdentifier(), qt.Equals, "b")
				c.Assert(p.Container(), qt.Equals, "b")
				c.Assert(p.ContainerDir(), qt.Equals, "/a")
				c.Assert(p.Dir(), qt.Equals, "/a/b")
				c.Assert(p.Ext(), qt.Equals, "md")
				c.Assert(p.Identifiers(), qt.DeepEquals, []string{"md", "no"})
				c.Assert(p.IsBranchBundle(), qt.IsFalse)
				c.Assert(p.IsBundle(), qt.IsTrue)
				c.Assert(p.IsLeafBundle(), qt.IsTrue)
				c.Assert(p.Lang(), qt.Equals, "no")
				c.Assert(p.NameNoExt(), qt.Equals, "index.no")
				c.Assert(p.NameNoIdentifier(), qt.Equals, "index")
				c.Assert(p.NameNoLang(), qt.Equals, "index.md")
				c.Assert(p.PathNoLang(), qt.Equals, "/a/b/index.md")
				c.Assert(p.Section(), qt.Equals, "a")
			},
		},
		{
			"Index branch content file",
			"/a/b/_index.no.md",
			func(c *qt.C, p *Path) {
				c.Assert(p.Base(), qt.Equals, "/a/b")
				c.Assert(p.BaseNameNoIdentifier(), qt.Equals, "b")
				c.Assert(p.Container(), qt.Equals, "b")
				c.Assert(p.ContainerDir(), qt.Equals, "/a")
				c.Assert(p.Ext(), qt.Equals, "md")
				c.Assert(p.Identifiers(), qt.DeepEquals, []string{"md", "no"})
				c.Assert(p.IsBranchBundle(), qt.IsTrue)
				c.Assert(p.IsBundle(), qt.IsTrue)
				c.Assert(p.IsLeafBundle(), qt.IsFalse)
				c.Assert(p.NameNoExt(), qt.Equals, "_index.no")
				c.Assert(p.NameNoLang(), qt.Equals, "_index.md")
			},
		},
		{
			"Index root no slash",
			"_index.md",
			func(c *qt.C, p *Path) {
				c.Assert(p.Base(), qt.Equals, "/")
				c.Assert(p.Ext(), qt.Equals, "md")
				c.Assert(p.Name(), qt.Equals, "_index.md")
			},
		},
		{
			"Index root",
			"/_index.md",
			func(c *qt.C, p *Path) {
				c.Assert(p.Base(), qt.Equals, "/")
				c.Assert(p.Ext(), qt.Equals, "md")
				c.Assert(p.Name(), qt.Equals, "_index.md")
			},
		},
		{
			"Index first",
			"/a/_index.md",
			func(c *qt.C, p *Path) {
				c.Assert(p.Section(), qt.Equals, "a")
			},
		},
		{
			"Index text file",
			"/a/b/index.no.txt",
			func(c *qt.C, p *Path) {
				c.Assert(p.Base(), qt.Equals, "/a/b/index.txt")
				c.Assert(p.Ext(), qt.Equals, "txt")
				c.Assert(p.Identifiers(), qt.DeepEquals, []string{"txt", "no"})
				c.Assert(p.IsLeafBundle(), qt.IsFalse)
				c.Assert(p.PathNoIdentifier(), qt.Equals, "/a/b/index")
			},
		},
		{
			"Empty",
			"",
			func(c *qt.C, p *Path) {
				c.Assert(p.Base(), qt.Equals, "/")
				c.Assert(p.Ext(), qt.Equals, "")
				c.Assert(p.Name(), qt.Equals, "")
				c.Assert(p.Path(), qt.Equals, "/")
			},
		},
		{
			"Slash",
			"/",
			func(c *qt.C, p *Path) {
				c.Assert(p.Base(), qt.Equals, "/")
				c.Assert(p.Ext(), qt.Equals, "")
				c.Assert(p.Name(), qt.Equals, "")
			},
		},
		{
			"Trim Leading Slash bundle",
			"foo/bar/index.no.md",
			func(c *qt.C, p *Path) {
				c.Assert(p.Path(), qt.Equals, "/foo/bar/index.no.md")
				pp := p.TrimLeadingSlash()
				c.Assert(pp.Path(), qt.Equals, "foo/bar/index.no.md")
				c.Assert(pp.PathNoLang(), qt.Equals, "foo/bar/index.md")
				c.Assert(pp.Base(), qt.Equals, "foo/bar")
				c.Assert(pp.Dir(), qt.Equals, "foo/bar")
				c.Assert(pp.ContainerDir(), qt.Equals, "foo")
				c.Assert(pp.Container(), qt.Equals, "bar")
				c.Assert(pp.BaseNameNoIdentifier(), qt.Equals, "bar")
			},
		},
		{
			"Trim Leading Slash file",
			"foo/bar.txt",
			func(c *qt.C, p *Path) {
				c.Assert(p.Path(), qt.Equals, "/foo/bar.txt")
				pp := p.TrimLeadingSlash()
				c.Assert(pp.Path(), qt.Equals, "foo/bar.txt")
				c.Assert(pp.PathNoLang(), qt.Equals, "foo/bar.txt")
				c.Assert(pp.Base(), qt.Equals, "foo/bar.txt")
				c.Assert(pp.Dir(), qt.Equals, "foo")
				c.Assert(pp.ContainerDir(), qt.Equals, "foo")
				c.Assert(pp.Container(), qt.Equals, "foo")
				c.Assert(pp.BaseNameNoIdentifier(), qt.Equals, "bar")
			},
		},
		{
			"File separator",
			filepath.FromSlash("/a/b/c.txt"),
			func(c *qt.C, p *Path) {
				c.Assert(p.Base(), qt.Equals, "/a/b/c.txt")
				c.Assert(p.Ext(), qt.Equals, "txt")
				c.Assert(p.Name(), qt.Equals, "c.txt")
				c.Assert(p.Path(), qt.Equals, "/a/b/c.txt")
			},
		},
		{
			"Content data file gotmpl",
			"/a/b/_content.gotmpl",
			func(c *qt.C, p *Path) {
				c.Assert(p.Path(), qt.Equals, "/a/b/_content.gotmpl")
				c.Assert(p.Ext(), qt.Equals, "gotmpl")
				c.Assert(p.IsContentData(), qt.IsTrue)
			},
		},
		{
			"Content data file yaml",
			"/a/b/_content.yaml",
			func(c *qt.C, p *Path) {
				c.Assert(p.IsContentData(), qt.IsFalse)
			},
		},
	}
	for _, test := range tests {
		c.Run(test.name, func(c *qt.C) {
			test.assert(c, testParser.Parse(files.ComponentFolderContent, test.path))
		})
	}
}

func TestHasExt(t *testing.T) {
	c := qt.New(t)

	c.Assert(HasExt("/a/b/c.txt"), qt.IsTrue)
	c.Assert(HasExt("/a/b.c/d.txt"), qt.IsTrue)
	c.Assert(HasExt("/a/b/c"), qt.IsFalse)
	c.Assert(HasExt("/a/b.c/d"), qt.IsFalse)
}

func BenchmarkParseIdentity(b *testing.B) {
	for i := 0; i < b.N; i++ {
		testParser.ParseIdentity(files.ComponentFolderAssets, "/a/b.css")
	}
}
