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
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/htesting"
	"github.com/spf13/afero"
)

func TestLanguageRootMapping(t *testing.T) {
	c := qt.New(t)
	v := viper.New()
	v.Set("contentDir", "content")

	fs := NewBaseFileDecorator(afero.NewMemMapFs())

	c.Assert(afero.WriteFile(fs, filepath.Join("content/sv/svdir", "main.txt"), []byte("main sv"), 0755), qt.IsNil)
	c.Assert(afero.WriteFile(fs, filepath.Join("themes/a/mysvblogcontent", "sv-f.txt"), []byte("some sv blog content"), 0755), qt.IsNil)
	c.Assert(afero.WriteFile(fs, filepath.Join("themes/a/myenblogcontent", "en-f.txt"), []byte("some en blog content in a"), 0755), qt.IsNil)
	c.Assert(afero.WriteFile(fs, filepath.Join("themes/a/myotherenblogcontent", "en-f2.txt"), []byte("some en content"), 0755), qt.IsNil)
	c.Assert(afero.WriteFile(fs, filepath.Join("themes/a/mysvdocs", "sv-docs.txt"), []byte("some sv docs content"), 0755), qt.IsNil)
	c.Assert(afero.WriteFile(fs, filepath.Join("themes/b/myenblogcontent", "en-b-f.txt"), []byte("some en content"), 0755), qt.IsNil)

	rfs, err := NewRootMappingFs(fs,
		RootMapping{
			From: "content/blog",             // Virtual path, first element is one of content, static, layouts etc.
			To:   "themes/a/mysvblogcontent", // Real path
			Meta: FileMeta{"lang": "sv"},
		},
		RootMapping{
			From: "content/blog",
			To:   "themes/a/myenblogcontent",
			Meta: FileMeta{"lang": "en"},
		},
		RootMapping{
			From: "content/blog",
			To:   "content/sv",
			Meta: FileMeta{"lang": "sv"},
		},
		RootMapping{
			From: "content/blog",
			To:   "themes/a/myotherenblogcontent",
			Meta: FileMeta{"lang": "en"},
		},
		RootMapping{
			From: "content/docs",
			To:   "themes/a/mysvdocs",
			Meta: FileMeta{"lang": "sv"},
		},
	)

	c.Assert(err, qt.IsNil)

	collected, err := collectFilenames(rfs, "content", "content")
	c.Assert(err, qt.IsNil)
	c.Assert(collected, qt.DeepEquals, []string{"blog/en-f.txt", "blog/en-f2.txt", "blog/sv-f.txt", "blog/svdir/main.txt", "docs/sv-docs.txt"})

	bfs := afero.NewBasePathFs(rfs, "content")
	collected, err = collectFilenames(bfs, "", "")
	c.Assert(err, qt.IsNil)
	c.Assert(collected, qt.DeepEquals, []string{"blog/en-f.txt", "blog/en-f2.txt", "blog/sv-f.txt", "blog/svdir/main.txt", "docs/sv-docs.txt"})

	dirs, err := rfs.Dirs(filepath.FromSlash("content/blog"))
	c.Assert(err, qt.IsNil)

	c.Assert(len(dirs), qt.Equals, 4)

	getDirnames := func(name string, rfs *RootMappingFs) []string {
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
		return rm.Meta.Lang() == "en"
	})

	c.Assert(getDirnames("content/blog", rfsEn), qt.DeepEquals, []string{"en-f.txt", "en-f2.txt"})

	rfsSv := rfs.Filter(func(rm RootMapping) bool {
		return rm.Meta.Lang() == "sv"
	})

	c.Assert(getDirnames("content/blog", rfsSv), qt.DeepEquals, []string{"sv-f.txt", "svdir"})

	// Make sure we have not messed with the original
	c.Assert(getDirnames("content/blog", rfs), qt.DeepEquals, []string{"sv-f.txt", "en-f.txt", "svdir", "en-f2.txt"})

	c.Assert(getDirnames("content", rfsSv), qt.DeepEquals, []string{"blog", "docs"})
	c.Assert(getDirnames("content", rfs), qt.DeepEquals, []string{"blog", "docs"})

}

func TestRootMappingFsDirnames(t *testing.T) {
	c := qt.New(t)
	fs := NewBaseFileDecorator(afero.NewMemMapFs())

	testfile := "myfile.txt"
	c.Assert(fs.Mkdir("f1t", 0755), qt.IsNil)
	c.Assert(fs.Mkdir("f2t", 0755), qt.IsNil)
	c.Assert(fs.Mkdir("f3t", 0755), qt.IsNil)
	c.Assert(afero.WriteFile(fs, filepath.Join("f2t", testfile), []byte("some content"), 0755), qt.IsNil)

	rfs, err := NewRootMappingFsFromFromTo(fs, "static/bf1", "f1t", "static/cf2", "f2t", "static/af3", "f3t")
	c.Assert(err, qt.IsNil)

	fif, err := rfs.Stat(filepath.Join("static/cf2", testfile))
	c.Assert(err, qt.IsNil)
	c.Assert(fif.Name(), qt.Equals, "myfile.txt")
	fifm := fif.(FileMetaInfo).Meta()
	c.Assert(fifm.Filename(), qt.Equals, filepath.FromSlash("f2t/myfile.txt"))

	root, err := rfs.Open(filepathSeparator)
	c.Assert(err, qt.IsNil)

	dirnames, err := root.Readdirnames(-1)
	c.Assert(err, qt.IsNil)
	c.Assert(dirnames, qt.DeepEquals, []string{"bf1", "cf2", "af3"})

}

func TestRootMappingFsFilename(t *testing.T) {
	c := qt.New(t)
	workDir, clean, err := htesting.CreateTempDir(Os, "hugo-root-filename")
	c.Assert(err, qt.IsNil)
	defer clean()
	fs := NewBaseFileDecorator(Os)

	testfilename := filepath.Join(workDir, "f1t/foo/file.txt")

	c.Assert(fs.MkdirAll(filepath.Join(workDir, "f1t/foo"), 0777), qt.IsNil)
	c.Assert(afero.WriteFile(fs, testfilename, []byte("content"), 0666), qt.IsNil)

	rfs, err := NewRootMappingFsFromFromTo(fs, "static/f1", filepath.Join(workDir, "f1t"), "static/f2", filepath.Join(workDir, "f2t"))
	c.Assert(err, qt.IsNil)

	fi, err := rfs.Stat(filepath.FromSlash("static/f1/foo/file.txt"))
	c.Assert(err, qt.IsNil)
	fim := fi.(FileMetaInfo)
	c.Assert(fim.Meta().Filename(), qt.Equals, testfilename)
	_, err = rfs.Stat(filepath.FromSlash("static/f1"))
	c.Assert(err, qt.IsNil)
}

func TestRootMappingFsMount(t *testing.T) {
	c := qt.New(t)
	fs := NewBaseFileDecorator(afero.NewMemMapFs())

	testfile := "test.txt"

	c.Assert(afero.WriteFile(fs, filepath.Join("themes/a/mynoblogcontent", testfile), []byte("some no content"), 0755), qt.IsNil)
	c.Assert(afero.WriteFile(fs, filepath.Join("themes/a/myenblogcontent", testfile), []byte("some en content"), 0755), qt.IsNil)
	c.Assert(afero.WriteFile(fs, filepath.Join("themes/a/mysvblogcontent", testfile), []byte("some sv content"), 0755), qt.IsNil)
	c.Assert(afero.WriteFile(fs, filepath.Join("themes/a/mysvblogcontent", "other.txt"), []byte("some sv content"), 0755), qt.IsNil)
	c.Assert(afero.WriteFile(fs, filepath.Join("themes/a/singlefiles", "no.txt"), []byte("no text"), 0755), qt.IsNil)
	c.Assert(afero.WriteFile(fs, filepath.Join("themes/a/singlefiles", "sv.txt"), []byte("sv text"), 0755), qt.IsNil)

	bfs := afero.NewBasePathFs(fs, "themes/a").(*afero.BasePathFs)
	rm := []RootMapping{
		// Directories
		RootMapping{
			From: "content/blog",
			To:   "mynoblogcontent",
			Meta: FileMeta{"lang": "no"},
		},
		RootMapping{
			From: "content/blog",
			To:   "myenblogcontent",
			Meta: FileMeta{"lang": "en"},
		},
		RootMapping{
			From: "content/blog",
			To:   "mysvblogcontent",
			Meta: FileMeta{"lang": "sv"},
		},
		// Files
		RootMapping{
			From:      "content/singles/p1.md",
			To:        "singlefiles/no.txt",
			ToBasedir: "singlefiles",
			Meta:      FileMeta{"lang": "no"},
		},
		RootMapping{
			From:      "content/singles/p1.md",
			To:        "singlefiles/sv.txt",
			ToBasedir: "singlefiles",
			Meta:      FileMeta{"lang": "sv"},
		},
	}

	rfs, err := NewRootMappingFs(bfs, rm...)
	c.Assert(err, qt.IsNil)

	blog, err := rfs.Stat(filepath.FromSlash("content/blog"))
	c.Assert(err, qt.IsNil)
	c.Assert(blog.IsDir(), qt.Equals, true)
	blogm := blog.(FileMetaInfo).Meta()
	c.Assert(blogm.Lang(), qt.Equals, "no") // First match

	f, err := blogm.Open()
	c.Assert(err, qt.IsNil)
	defer f.Close()
	dirs1, err := f.Readdirnames(-1)
	c.Assert(err, qt.IsNil)
	// Union with duplicate dir names filtered.
	c.Assert(dirs1, qt.DeepEquals, []string{"test.txt", "test.txt", "other.txt", "test.txt"})

	files, err := afero.ReadDir(rfs, filepath.FromSlash("content/blog"))
	c.Assert(err, qt.IsNil)
	c.Assert(len(files), qt.Equals, 4)

	testfilefi := files[1]
	c.Assert(testfilefi.Name(), qt.Equals, testfile)

	testfilem := testfilefi.(FileMetaInfo).Meta()
	c.Assert(testfilem.Filename(), qt.Equals, filepath.FromSlash("themes/a/mynoblogcontent/test.txt"))

	tf, err := testfilem.Open()
	c.Assert(err, qt.IsNil)
	defer tf.Close()
	b, err := ioutil.ReadAll(tf)
	c.Assert(err, qt.IsNil)
	c.Assert(string(b), qt.Equals, "some no content")

	// Check file mappings
	single, err := rfs.Stat(filepath.FromSlash("content/singles/p1.md"))
	c.Assert(err, qt.IsNil)
	c.Assert(single.IsDir(), qt.Equals, false)
	singlem := single.(FileMetaInfo).Meta()
	c.Assert(singlem.Lang(), qt.Equals, "no") // First match

	singlesDir, err := rfs.Open(filepath.FromSlash("content/singles"))
	c.Assert(err, qt.IsNil)
	defer singlesDir.Close()
	singles, err := singlesDir.Readdir(-1)
	c.Assert(err, qt.IsNil)
	c.Assert(singles, qt.HasLen, 2)
	for i, lang := range []string{"no", "sv"} {
		fi := singles[i].(FileMetaInfo)
		c.Assert(fi.Meta().PathFile(), qt.Equals, lang+".txt")
		c.Assert(fi.Meta().Lang(), qt.Equals, lang)
		c.Assert(fi.Name(), qt.Equals, "p1.md")
	}
}

func TestRootMappingFsMountOverlap(t *testing.T) {
	c := qt.New(t)
	fs := NewBaseFileDecorator(afero.NewMemMapFs())

	c.Assert(afero.WriteFile(fs, filepath.FromSlash("da/a.txt"), []byte("some no content"), 0755), qt.IsNil)
	c.Assert(afero.WriteFile(fs, filepath.FromSlash("db/b.txt"), []byte("some no content"), 0755), qt.IsNil)
	c.Assert(afero.WriteFile(fs, filepath.FromSlash("dc/c.txt"), []byte("some no content"), 0755), qt.IsNil)
	c.Assert(afero.WriteFile(fs, filepath.FromSlash("de/e.txt"), []byte("some no content"), 0755), qt.IsNil)

	rm := []RootMapping{
		RootMapping{
			From: "static",
			To:   "da",
		},
		RootMapping{
			From: "static/b",
			To:   "db",
		},
		RootMapping{
			From: "static/b/c",
			To:   "dc",
		},
		RootMapping{
			From: "/static/e/",
			To:   "de",
		},
	}

	rfs, err := NewRootMappingFs(fs, rm...)
	c.Assert(err, qt.IsNil)

	getDirnames := func(name string) []string {
		name = filepath.FromSlash(name)
		f, err := rfs.Open(name)
		c.Assert(err, qt.IsNil)
		defer f.Close()
		names, err := f.Readdirnames(-1)
		c.Assert(err, qt.IsNil)
		return names
	}

	c.Assert(getDirnames("static"), qt.DeepEquals, []string{"a.txt", "b", "e"})
	c.Assert(getDirnames("static/b"), qt.DeepEquals, []string{"b.txt", "c"})
	c.Assert(getDirnames("static/b/c"), qt.DeepEquals, []string{"c.txt"})

	fi, err := rfs.Stat(filepath.FromSlash("static/b/b.txt"))
	c.Assert(err, qt.IsNil)
	c.Assert(fi.Name(), qt.Equals, "b.txt")

}

func TestRootMappingFsOs(t *testing.T) {
	c := qt.New(t)
	fs := afero.NewOsFs()

	d, err := ioutil.TempDir("", "hugo-root-mapping")
	c.Assert(err, qt.IsNil)
	defer func() {
		os.RemoveAll(d)
	}()

	testfile := "myfile.txt"
	c.Assert(fs.Mkdir(filepath.Join(d, "f1t"), 0755), qt.IsNil)
	c.Assert(fs.Mkdir(filepath.Join(d, "f2t"), 0755), qt.IsNil)
	c.Assert(fs.Mkdir(filepath.Join(d, "f3t"), 0755), qt.IsNil)
	c.Assert(afero.WriteFile(fs, filepath.Join(d, "f2t", testfile), []byte("some content"), 0755), qt.IsNil)

	rfs, err := NewRootMappingFsFromFromTo(fs, "static/bf1", filepath.Join(d, "f1t"), "static/cf2", filepath.Join(d, "f2t"), "static/af3", filepath.Join(d, "f3t"))
	c.Assert(err, qt.IsNil)

	fif, err := rfs.Stat(filepath.Join("static/cf2", testfile))
	c.Assert(err, qt.IsNil)
	c.Assert(fif.Name(), qt.Equals, "myfile.txt")

	root, err := rfs.Open(filepathSeparator)
	c.Assert(err, qt.IsNil)

	dirnames, err := root.Readdirnames(-1)
	c.Assert(err, qt.IsNil)
	c.Assert(dirnames, qt.DeepEquals, []string{"bf1", "cf2", "af3"})

}
