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

package hugofs

import (
	"fmt"
	"path/filepath"
	"sort"
	"testing"

	iofs "io/fs"

	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/hugofs/glob"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/htesting"
	"github.com/spf13/afero"
)

func TestLanguageRootMapping(t *testing.T) {
	c := qt.New(t)
	v := config.New()
	v.Set("contentDir", "content")

	fs := NewBaseFileDecorator(afero.NewMemMapFs())

	c.Assert(afero.WriteFile(fs, filepath.Join("content/sv/svdir", "main.txt"), []byte("main sv"), 0o755), qt.IsNil)

	c.Assert(afero.WriteFile(fs, filepath.Join("themes/a/mysvblogcontent", "sv-f.txt"), []byte("some sv blog content"), 0o755), qt.IsNil)
	c.Assert(afero.WriteFile(fs, filepath.Join("themes/a/myenblogcontent", "en-f.txt"), []byte("some en blog content in a"), 0o755), qt.IsNil)
	c.Assert(afero.WriteFile(fs, filepath.Join("themes/a/mysvblogcontent/d1", "sv-d1-f.txt"), []byte("some sv blog content"), 0o755), qt.IsNil)
	c.Assert(afero.WriteFile(fs, filepath.Join("themes/a/myenblogcontent/d1", "en-d1-f.txt"), []byte("some en blog content in a"), 0o755), qt.IsNil)

	c.Assert(afero.WriteFile(fs, filepath.Join("themes/a/myotherenblogcontent", "en-f2.txt"), []byte("some en content"), 0o755), qt.IsNil)
	c.Assert(afero.WriteFile(fs, filepath.Join("themes/a/mysvdocs", "sv-docs.txt"), []byte("some sv docs content"), 0o755), qt.IsNil)
	c.Assert(afero.WriteFile(fs, filepath.Join("themes/b/myenblogcontent", "en-b-f.txt"), []byte("some en content"), 0o755), qt.IsNil)

	rfs, err := NewRootMappingFs(fs,
		RootMapping{
			From: "content/blog",             // Virtual path, first element is one of content, static, layouts etc.
			To:   "themes/a/mysvblogcontent", // Real path
			Meta: &FileMeta{Lang: "sv"},
		},
		RootMapping{
			From: "content/blog",
			To:   "themes/a/myenblogcontent",
			Meta: &FileMeta{Lang: "en"},
		},
		RootMapping{
			From: "content/blog",
			To:   "content/sv",
			Meta: &FileMeta{Lang: "sv"},
		},
		RootMapping{
			From: "content/blog",
			To:   "themes/a/myotherenblogcontent",
			Meta: &FileMeta{Lang: "en"},
		},
		RootMapping{
			From: "content/docs",
			To:   "themes/a/mysvdocs",
			Meta: &FileMeta{Lang: "sv"},
		},
	)

	c.Assert(err, qt.IsNil)

	collected, err := collectPaths(rfs, "content")
	c.Assert(err, qt.IsNil)
	c.Assert(collected, qt.DeepEquals,
		[]string{"/blog/d1/en-d1-f.txt", "/blog/d1/sv-d1-f.txt", "/blog/en-f.txt", "/blog/en-f2.txt", "/blog/sv-f.txt", "/blog/svdir/main.txt", "/docs/sv-docs.txt"}, qt.Commentf("%#v", collected))

	dirs, err := rfs.Mounts(filepath.FromSlash("content/blog"))
	c.Assert(err, qt.IsNil)
	c.Assert(len(dirs), qt.Equals, 4)
	for _, dir := range dirs {
		f, err := dir.Meta().Open()
		c.Assert(err, qt.IsNil)
		f.Close()
	}

	blog, err := rfs.Open(filepath.FromSlash("content/blog"))
	c.Assert(err, qt.IsNil)
	fis, err := blog.(iofs.ReadDirFile).ReadDir(-1)
	c.Assert(err, qt.IsNil)
	for _, fi := range fis {
		f, err := fi.(FileMetaInfo).Meta().Open()
		c.Assert(err, qt.IsNil)
		f.Close()
	}
	blog.Close()

	getDirnames := func(name string, rfs *RootMappingFs) []string {
		c.Helper()
		filename := filepath.FromSlash(name)
		f, err := rfs.Open(filename)
		c.Assert(err, qt.IsNil)
		names, err := f.Readdirnames(-1)

		f.Close()
		c.Assert(err, qt.IsNil)

		info, err := rfs.Stat(filename)
		c.Assert(err, qt.IsNil)
		f2, err := info.(FileMetaInfo).Meta().Open()
		c.Assert(err, qt.IsNil)
		names2, err := f2.Readdirnames(-1)
		c.Assert(err, qt.IsNil)
		c.Assert(names2, qt.DeepEquals, names)
		f2.Close()

		return names
	}

	rfsEn := rfs.Filter(func(rm RootMapping) bool {
		return rm.Meta.Lang == "en"
	})

	c.Assert(getDirnames("content/blog", rfsEn), qt.DeepEquals, []string{"d1", "en-f.txt", "en-f2.txt"})

	rfsSv := rfs.Filter(func(rm RootMapping) bool {
		return rm.Meta.Lang == "sv"
	})

	c.Assert(getDirnames("content/blog", rfsSv), qt.DeepEquals, []string{"d1", "sv-f.txt", "svdir"})

	// Make sure we have not messed with the original
	c.Assert(getDirnames("content/blog", rfs), qt.DeepEquals, []string{"d1", "sv-f.txt", "en-f.txt", "svdir", "en-f2.txt"})

	c.Assert(getDirnames("content", rfsSv), qt.DeepEquals, []string{"blog", "docs"})
	c.Assert(getDirnames("content", rfs), qt.DeepEquals, []string{"blog", "docs"})
}

func TestRootMappingFsDirnames(t *testing.T) {
	c := qt.New(t)
	fs := NewBaseFileDecorator(afero.NewMemMapFs())

	testfile := "myfile.txt"
	c.Assert(fs.Mkdir("f1t", 0o755), qt.IsNil)
	c.Assert(fs.Mkdir("f2t", 0o755), qt.IsNil)
	c.Assert(fs.Mkdir("f3t", 0o755), qt.IsNil)
	c.Assert(afero.WriteFile(fs, filepath.Join("f2t", testfile), []byte("some content"), 0o755), qt.IsNil)

	rfs, err := newRootMappingFsFromFromTo("", fs, "static/bf1", "f1t", "static/cf2", "f2t", "static/af3", "f3t")
	c.Assert(err, qt.IsNil)

	fif, err := rfs.Stat(filepath.Join("static/cf2", testfile))
	c.Assert(err, qt.IsNil)
	c.Assert(fif.Name(), qt.Equals, "myfile.txt")
	fifm := fif.(FileMetaInfo).Meta()
	c.Assert(fifm.Filename, qt.Equals, filepath.FromSlash("f2t/myfile.txt"))

	root, err := rfs.Open("static")
	c.Assert(err, qt.IsNil)

	dirnames, err := root.Readdirnames(-1)
	c.Assert(err, qt.IsNil)
	c.Assert(dirnames, qt.DeepEquals, []string{"af3", "bf1", "cf2"})
}

func TestRootMappingFsFilename(t *testing.T) {
	c := qt.New(t)
	workDir, clean, err := htesting.CreateTempDir(Os, "hugo-root-filename")
	c.Assert(err, qt.IsNil)
	defer clean()
	fs := NewBaseFileDecorator(Os)

	testfilename := filepath.Join(workDir, "f1t/foo/file.txt")

	c.Assert(fs.MkdirAll(filepath.Join(workDir, "f1t/foo"), 0o777), qt.IsNil)
	c.Assert(afero.WriteFile(fs, testfilename, []byte("content"), 0o666), qt.IsNil)

	rfs, err := newRootMappingFsFromFromTo(workDir, fs, "static/f1", filepath.Join(workDir, "f1t"), "static/f2", filepath.Join(workDir, "f2t"))
	c.Assert(err, qt.IsNil)

	fi, err := rfs.Stat(filepath.FromSlash("static/f1/foo/file.txt"))
	c.Assert(err, qt.IsNil)
	fim := fi.(FileMetaInfo)
	c.Assert(fim.Meta().Filename, qt.Equals, testfilename)
	_, err = rfs.Stat(filepath.FromSlash("static/f1"))
	c.Assert(err, qt.IsNil)
}

func TestRootMappingFsMount(t *testing.T) {
	c := qt.New(t)
	fs := NewBaseFileDecorator(afero.NewMemMapFs())

	testfile := "test.txt"

	c.Assert(afero.WriteFile(fs, filepath.Join("themes/a/mynoblogcontent", testfile), []byte("some no content"), 0o755), qt.IsNil)
	c.Assert(afero.WriteFile(fs, filepath.Join("themes/a/myenblogcontent", testfile), []byte("some en content"), 0o755), qt.IsNil)
	c.Assert(afero.WriteFile(fs, filepath.Join("themes/a/mysvblogcontent", testfile), []byte("some sv content"), 0o755), qt.IsNil)
	c.Assert(afero.WriteFile(fs, filepath.Join("themes/a/mysvblogcontent", "other.txt"), []byte("some sv content"), 0o755), qt.IsNil)
	c.Assert(afero.WriteFile(fs, filepath.Join("themes/a/singlefiles", "no.txt"), []byte("no text"), 0o755), qt.IsNil)
	c.Assert(afero.WriteFile(fs, filepath.Join("themes/a/singlefiles", "sv.txt"), []byte("sv text"), 0o755), qt.IsNil)

	bfs := NewBasePathFs(fs, "themes/a")
	rm := []RootMapping{
		// Directories
		{
			From: "content/blog",
			To:   "mynoblogcontent",
			Meta: &FileMeta{Lang: "no"},
		},
		{
			From: "content/blog",
			To:   "myenblogcontent",
			Meta: &FileMeta{Lang: "en"},
		},
		{
			From: "content/blog",
			To:   "mysvblogcontent",
			Meta: &FileMeta{Lang: "sv"},
		},
		// Files
		{
			From:   "content/singles/p1.md",
			To:     "singlefiles/no.txt",
			ToBase: "singlefiles",
			Meta:   &FileMeta{Lang: "no"},
		},
		{
			From:   "content/singles/p1.md",
			To:     "singlefiles/sv.txt",
			ToBase: "singlefiles",
			Meta:   &FileMeta{Lang: "sv"},
		},
	}

	rfs, err := NewRootMappingFs(bfs, rm...)
	c.Assert(err, qt.IsNil)

	blog, err := rfs.Stat(filepath.FromSlash("content/blog"))
	c.Assert(err, qt.IsNil)
	c.Assert(blog.IsDir(), qt.Equals, true)
	blogm := blog.(FileMetaInfo).Meta()
	c.Assert(blogm.Lang, qt.Equals, "no") // First match

	f, err := blogm.Open()
	c.Assert(err, qt.IsNil)
	defer f.Close()
	dirs1, err := f.Readdirnames(-1)
	c.Assert(err, qt.IsNil)
	// Union with duplicate dir names filtered.
	c.Assert(dirs1, qt.DeepEquals, []string{"test.txt", "test.txt", "other.txt", "test.txt"})

	d, err := rfs.Open(filepath.FromSlash("content/blog"))
	c.Assert(err, qt.IsNil)
	files, err := d.(iofs.ReadDirFile).ReadDir(-1)
	c.Assert(err, qt.IsNil)
	c.Assert(len(files), qt.Equals, 4)

	singlesDir, err := rfs.Open(filepath.FromSlash("content/singles"))
	c.Assert(err, qt.IsNil)
	defer singlesDir.Close()
	singles, err := singlesDir.(iofs.ReadDirFile).ReadDir(-1)
	c.Assert(err, qt.IsNil)
	c.Assert(singles, qt.HasLen, 2)
	for i, lang := range []string{"no", "sv"} {
		fi := singles[i].(FileMetaInfo)
		c.Assert(fi.Meta().Lang, qt.Equals, lang)
		c.Assert(fi.Name(), qt.Equals, "p1.md")
	}

	// Test ReverseLookup.
	// Single file mounts.
	cps, err := rfs.ReverseLookup(filepath.FromSlash("singlefiles/no.txt"))
	c.Assert(err, qt.IsNil)
	c.Assert(cps, qt.DeepEquals, []ComponentPath{
		{Component: "content", Path: "singles/p1.md", Lang: "no"},
	})

	cps, err = rfs.ReverseLookup(filepath.FromSlash("singlefiles/sv.txt"))
	c.Assert(err, qt.IsNil)
	c.Assert(cps, qt.DeepEquals, []ComponentPath{
		{Component: "content", Path: "singles/p1.md", Lang: "sv"},
	})

	// File inside directory mount.
	cps, err = rfs.ReverseLookup(filepath.FromSlash("mynoblogcontent/test.txt"))
	c.Assert(err, qt.IsNil)
	c.Assert(cps, qt.DeepEquals, []ComponentPath{
		{Component: "content", Path: "blog/test.txt", Lang: "no"},
	})
}

func TestRootMappingFsMountOverlap(t *testing.T) {
	c := qt.New(t)
	fs := NewBaseFileDecorator(afero.NewMemMapFs())

	c.Assert(afero.WriteFile(fs, filepath.FromSlash("da/a.txt"), []byte("some no content"), 0o755), qt.IsNil)
	c.Assert(afero.WriteFile(fs, filepath.FromSlash("db/b.txt"), []byte("some no content"), 0o755), qt.IsNil)
	c.Assert(afero.WriteFile(fs, filepath.FromSlash("dc/c.txt"), []byte("some no content"), 0o755), qt.IsNil)
	c.Assert(afero.WriteFile(fs, filepath.FromSlash("de/e.txt"), []byte("some no content"), 0o755), qt.IsNil)

	rm := []RootMapping{
		{
			From: "static",
			To:   "da",
		},
		{
			From: "static/b",
			To:   "db",
		},
		{
			From: "static/b/c",
			To:   "dc",
		},
		{
			From: "/static/e/",
			To:   "de",
		},
	}

	rfs, err := NewRootMappingFs(fs, rm...)
	c.Assert(err, qt.IsNil)

	checkDirnames := func(name string, expect []string) {
		c.Helper()
		name = filepath.FromSlash(name)
		f, err := rfs.Open(name)
		c.Assert(err, qt.IsNil)
		defer f.Close()
		names, err := f.Readdirnames(-1)
		c.Assert(err, qt.IsNil)
		c.Assert(names, qt.DeepEquals, expect, qt.Commentf(fmt.Sprintf("%#v", names)))
	}

	checkDirnames("static", []string{"a.txt", "b", "e"})
	checkDirnames("static/b", []string{"b.txt", "c"})
	checkDirnames("static/b/c", []string{"c.txt"})

	fi, err := rfs.Stat(filepath.FromSlash("static/b/b.txt"))
	c.Assert(err, qt.IsNil)
	c.Assert(fi.Name(), qt.Equals, "b.txt")
}

func TestRootMappingFsOs(t *testing.T) {
	c := qt.New(t)
	fs := NewBaseFileDecorator(afero.NewOsFs())

	d, clean, err := htesting.CreateTempDir(fs, "hugo-root-mapping-os")
	c.Assert(err, qt.IsNil)
	defer clean()

	testfile := "myfile.txt"
	c.Assert(fs.Mkdir(filepath.Join(d, "f1t"), 0o755), qt.IsNil)
	c.Assert(fs.Mkdir(filepath.Join(d, "f2t"), 0o755), qt.IsNil)
	c.Assert(fs.Mkdir(filepath.Join(d, "f3t"), 0o755), qt.IsNil)

	// Deep structure
	deepDir := filepath.Join(d, "d1", "d2", "d3", "d4", "d5")
	c.Assert(fs.MkdirAll(deepDir, 0o755), qt.IsNil)
	for i := 1; i <= 3; i++ {
		c.Assert(fs.MkdirAll(filepath.Join(d, "d1", "d2", "d3", "d4", fmt.Sprintf("d4-%d", i)), 0o755), qt.IsNil)
		c.Assert(afero.WriteFile(fs, filepath.Join(d, "d1", "d2", "d3", fmt.Sprintf("f-%d.txt", i)), []byte("some content"), 0o755), qt.IsNil)
	}

	c.Assert(afero.WriteFile(fs, filepath.Join(d, "f2t", testfile), []byte("some content"), 0o755), qt.IsNil)

	// https://github.com/gohugoio/hugo/issues/6854
	mystaticDir := filepath.Join(d, "mystatic", "a", "b", "c")
	c.Assert(fs.MkdirAll(mystaticDir, 0o755), qt.IsNil)
	c.Assert(afero.WriteFile(fs, filepath.Join(mystaticDir, "ms-1.txt"), []byte("some content"), 0o755), qt.IsNil)

	rfs, err := newRootMappingFsFromFromTo(
		d,
		fs,
		"static/bf1", filepath.Join(d, "f1t"),
		"static/cf2", filepath.Join(d, "f2t"),
		"static/af3", filepath.Join(d, "f3t"),
		"static", filepath.Join(d, "mystatic"),
		"static/a/b/c", filepath.Join(d, "d1", "d2", "d3"),
		"layouts", filepath.Join(d, "d1"),
	)

	c.Assert(err, qt.IsNil)

	fif, err := rfs.Stat(filepath.Join("static/cf2", testfile))
	c.Assert(err, qt.IsNil)
	c.Assert(fif.Name(), qt.Equals, "myfile.txt")

	root, err := rfs.Open("static")
	c.Assert(err, qt.IsNil)

	dirnames, err := root.Readdirnames(-1)
	c.Assert(err, qt.IsNil)
	c.Assert(dirnames, qt.DeepEquals, []string{"a", "af3", "bf1", "cf2"}, qt.Commentf(fmt.Sprintf("%#v", dirnames)))

	getDirnames := func(dirname string) []string {
		dirname = filepath.FromSlash(dirname)
		f, err := rfs.Open(dirname)
		c.Assert(err, qt.IsNil)
		defer f.Close()
		dirnames, err := f.Readdirnames(-1)
		c.Assert(err, qt.IsNil)
		sort.Strings(dirnames)
		return dirnames
	}

	c.Assert(getDirnames("static/a/b"), qt.DeepEquals, []string{"c"})
	c.Assert(getDirnames("static/a/b/c"), qt.DeepEquals, []string{"d4", "f-1.txt", "f-2.txt", "f-3.txt", "ms-1.txt"})
	c.Assert(getDirnames("static/a/b/c/d4"), qt.DeepEquals, []string{"d4-1", "d4-2", "d4-3", "d5"})

	all, err := collectPaths(rfs, "static")
	c.Assert(err, qt.IsNil)

	c.Assert(all, qt.DeepEquals, []string{"/a/b/c/f-1.txt", "/a/b/c/f-2.txt", "/a/b/c/f-3.txt", "/a/b/c/ms-1.txt", "/cf2/myfile.txt"})

	fis, err := collectFileinfos(rfs, "static")
	c.Assert(err, qt.IsNil)

	dirc := fis[3].Meta()

	f, err := dirc.Open()
	c.Assert(err, qt.IsNil)
	defer f.Close()
	dirEntries, err := f.(iofs.ReadDirFile).ReadDir(-1)
	c.Assert(err, qt.IsNil)
	sortDirEntries(dirEntries)
	i := 0
	for _, fi := range dirEntries {
		if fi.IsDir() || fi.Name() == "ms-1.txt" {
			continue
		}
		i++
		meta := fi.(FileMetaInfo).Meta()
		c.Assert(meta.Filename, qt.Equals, filepath.Join(d, fmt.Sprintf("/d1/d2/d3/f-%d.txt", i)))
	}

	_, err = rfs.Stat(filepath.FromSlash("layouts/d2/d3/f-1.txt"))
	c.Assert(err, qt.IsNil)
	_, err = rfs.Stat(filepath.FromSlash("layouts/d2/d3"))
	c.Assert(err, qt.IsNil)
}

func TestRootMappingFsOsBase(t *testing.T) {
	c := qt.New(t)
	fs := NewBaseFileDecorator(afero.NewOsFs())

	d, clean, err := htesting.CreateTempDir(fs, "hugo-root-mapping-os-base")
	c.Assert(err, qt.IsNil)
	defer clean()

	// Deep structure
	deepDir := filepath.Join(d, "d1", "d2", "d3", "d4", "d5")
	c.Assert(fs.MkdirAll(deepDir, 0o755), qt.IsNil)
	for i := 1; i <= 3; i++ {
		c.Assert(fs.MkdirAll(filepath.Join(d, "d1", "d2", "d3", "d4", fmt.Sprintf("d4-%d", i)), 0o755), qt.IsNil)
		c.Assert(afero.WriteFile(fs, filepath.Join(d, "d1", "d2", "d3", fmt.Sprintf("f-%d.txt", i)), []byte("some content"), 0o755), qt.IsNil)
	}

	mystaticDir := filepath.Join(d, "mystatic", "a", "b", "c")
	c.Assert(fs.MkdirAll(mystaticDir, 0o755), qt.IsNil)
	c.Assert(afero.WriteFile(fs, filepath.Join(mystaticDir, "ms-1.txt"), []byte("some content"), 0o755), qt.IsNil)

	bfs := NewBasePathFs(fs, d)

	rfs, err := newRootMappingFsFromFromTo(
		"",
		bfs,
		"static", "mystatic",
		"static/a/b/c", filepath.Join("d1", "d2", "d3"),
	)
	c.Assert(err, qt.IsNil)

	getDirnames := func(dirname string) []string {
		dirname = filepath.FromSlash(dirname)
		f, err := rfs.Open(dirname)
		c.Assert(err, qt.IsNil)
		defer f.Close()
		dirnames, err := f.Readdirnames(-1)
		c.Assert(err, qt.IsNil)
		sort.Strings(dirnames)
		return dirnames
	}

	c.Assert(getDirnames("static/a/b/c"), qt.DeepEquals, []string{"d4", "f-1.txt", "f-2.txt", "f-3.txt", "ms-1.txt"})
}

func TestRootMappingFileFilter(t *testing.T) {
	c := qt.New(t)
	fs := NewBaseFileDecorator(afero.NewMemMapFs())

	for _, lang := range []string{"no", "en", "fr"} {
		for i := 1; i <= 3; i++ {
			c.Assert(afero.WriteFile(fs, filepath.Join(lang, fmt.Sprintf("my%s%d.txt", lang, i)), []byte("some text file for"+lang), 0o755), qt.IsNil)
		}
	}

	for _, lang := range []string{"no", "en", "fr"} {
		for i := 1; i <= 3; i++ {
			c.Assert(afero.WriteFile(fs, filepath.Join(lang, "sub", fmt.Sprintf("mysub%s%d.txt", lang, i)), []byte("some text file for"+lang), 0o755), qt.IsNil)
		}
	}

	rm := []RootMapping{
		{
			From: "content",
			To:   "no",
			Meta: &FileMeta{Lang: "no", InclusionFilter: glob.MustNewFilenameFilter(nil, []string{"**.txt"})},
		},
		{
			From: "content",
			To:   "en",
			Meta: &FileMeta{Lang: "en"},
		},
		{
			From: "content",
			To:   "fr",
			Meta: &FileMeta{Lang: "fr", InclusionFilter: glob.MustNewFilenameFilter(nil, []string{"**.txt"})},
		},
	}

	rfs, err := NewRootMappingFs(fs, rm...)
	c.Assert(err, qt.IsNil)

	assertExists := func(filename string, shouldExist bool) {
		c.Helper()
		filename = filepath.Clean(filename)
		_, err1 := rfs.Stat(filename)
		f, err2 := rfs.Open(filename)
		if shouldExist {
			c.Assert(err1, qt.IsNil)
			c.Assert(err2, qt.IsNil)
			c.Assert(f.Close(), qt.IsNil)
		} else {
			c.Assert(err1, qt.Not(qt.IsNil))
			c.Assert(err2, qt.Not(qt.IsNil))
		}
	}

	assertExists("content/myno1.txt", false)
	assertExists("content/myen1.txt", true)
	assertExists("content/myfr1.txt", false)

	dirEntriesSub, err := afero.ReadDir(rfs, filepath.Join("content", "sub"))
	c.Assert(err, qt.IsNil)
	c.Assert(len(dirEntriesSub), qt.Equals, 3)

	f, err := rfs.Open("content")
	c.Assert(err, qt.IsNil)
	defer f.Close()
	dirEntries, err := f.(iofs.ReadDirFile).ReadDir(-1)

	c.Assert(err, qt.IsNil)
	c.Assert(len(dirEntries), qt.Equals, 4)
}
