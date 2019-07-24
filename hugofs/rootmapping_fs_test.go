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

	"github.com/gohugoio/hugo/htesting"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestLanguageRootMapping(t *testing.T) {
	assert := require.New(t)
	v := viper.New()
	v.Set("contentDir", "content")

	fs := NewBaseFileDecorator(afero.NewMemMapFs())

	assert.NoError(afero.WriteFile(fs, filepath.Join("content/sv/svdir", "main.txt"), []byte("main sv"), 0755))
	assert.NoError(afero.WriteFile(fs, filepath.Join("themes/a/mysvblogcontent", "sv-f.txt"), []byte("some sv blog content"), 0755))
	assert.NoError(afero.WriteFile(fs, filepath.Join("themes/a/myenblogcontent", "en-f.txt"), []byte("some en blog content in a"), 0755))
	assert.NoError(afero.WriteFile(fs, filepath.Join("themes/a/myotherenblogcontent", "en-f2.txt"), []byte("some en content"), 0755))
	assert.NoError(afero.WriteFile(fs, filepath.Join("themes/a/mysvdocs", "sv-docs.txt"), []byte("some sv docs content"), 0755))
	assert.NoError(afero.WriteFile(fs, filepath.Join("themes/b/myenblogcontent", "en-b-f.txt"), []byte("some en content"), 0755))

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

	assert.NoError(err)

	collected, err := collectFilenames(rfs, "content", "content")
	assert.NoError(err)
	assert.Equal([]string{"blog/en-f.txt", "blog/en-f2.txt", "blog/sv-f.txt", "blog/svdir/main.txt", "docs/sv-docs.txt"}, collected)

	bfs := afero.NewBasePathFs(rfs, "content")
	collected, err = collectFilenames(bfs, "", "")
	assert.NoError(err)
	assert.Equal([]string{"blog/en-f.txt", "blog/en-f2.txt", "blog/sv-f.txt", "blog/svdir/main.txt", "docs/sv-docs.txt"}, collected)

	dirs, err := rfs.Dirs(filepath.FromSlash("content/blog"))
	assert.NoError(err)

	assert.Equal(4, len(dirs))

	getDirnames := func(name string, rfs *RootMappingFs) []string {
		filename := filepath.FromSlash(name)
		f, err := rfs.Open(filename)
		assert.NoError(err)
		names, err := f.Readdirnames(-1)

		f.Close()
		assert.NoError(err)

		info, err := rfs.Stat(filename)
		assert.NoError(err)
		f2, err := info.(FileMetaInfo).Meta().Open()
		assert.NoError(err)
		names2, err := f2.Readdirnames(-1)
		assert.NoError(err)
		assert.Equal(names, names2)
		f2.Close()

		return names
	}

	rfsEn := rfs.Filter(func(rm RootMapping) bool {
		return rm.Meta.Lang() == "en"
	})

	assert.Equal([]string{"en-f.txt", "en-f2.txt"}, getDirnames("content/blog", rfsEn))

	rfsSv := rfs.Filter(func(rm RootMapping) bool {
		return rm.Meta.Lang() == "sv"
	})

	assert.Equal([]string{"sv-f.txt", "svdir"}, getDirnames("content/blog", rfsSv))

	// Make sure we have not messed with the original
	assert.Equal([]string{"sv-f.txt", "en-f.txt", "svdir", "en-f2.txt"}, getDirnames("content/blog", rfs))

	assert.Equal([]string{"blog", "docs"}, getDirnames("content", rfsSv))
	assert.Equal([]string{"blog", "docs"}, getDirnames("content", rfs))
}

func TestRootMappingFsDirnames(t *testing.T) {
	assert := require.New(t)
	fs := NewBaseFileDecorator(afero.NewMemMapFs())

	testfile := "myfile.txt"
	assert.NoError(fs.Mkdir("f1t", 0755))
	assert.NoError(fs.Mkdir("f2t", 0755))
	assert.NoError(fs.Mkdir("f3t", 0755))
	assert.NoError(afero.WriteFile(fs, filepath.Join("f2t", testfile), []byte("some content"), 0755))

	rfs, err := NewRootMappingFsFromFromTo(fs, "static/bf1", "f1t", "static/cf2", "f2t", "static/af3", "f3t")
	assert.NoError(err)

	fif, err := rfs.Stat(filepath.Join("static/cf2", testfile))
	assert.NoError(err)
	assert.Equal("myfile.txt", fif.Name())
	fifm := fif.(FileMetaInfo).Meta()
	assert.Equal(filepath.FromSlash("f2t/myfile.txt"), fifm.Filename())

	root, err := rfs.Open(filepathSeparator)
	assert.NoError(err)

	dirnames, err := root.Readdirnames(-1)
	assert.NoError(err)
	assert.Equal([]string{"bf1", "cf2", "af3"}, dirnames)
}

func TestRootMappingFsFilename(t *testing.T) {
	assert := require.New(t)
	workDir, clean, err := htesting.CreateTempDir(Os, "hugo-root-filename")
	assert.NoError(err)
	defer clean()
	fs := NewBaseFileDecorator(Os)

	testfilename := filepath.Join(workDir, "f1t/foo/file.txt")

	assert.NoError(fs.MkdirAll(filepath.Join(workDir, "f1t/foo"), 0777))
	assert.NoError(afero.WriteFile(fs, testfilename, []byte("content"), 0666))

	rfs, err := NewRootMappingFsFromFromTo(fs, "static/f1", filepath.Join(workDir, "f1t"), "static/f2", filepath.Join(workDir, "f2t"))
	assert.NoError(err)

	fi, err := rfs.Stat(filepath.FromSlash("static/f1/foo/file.txt"))
	assert.NoError(err)
	fim := fi.(FileMetaInfo)
	assert.Equal(testfilename, fim.Meta().Filename())
	_, err = rfs.Stat(filepath.FromSlash("static/f1"))
	assert.NoError(err)
}

func TestRootMappingFsMount(t *testing.T) {
	assert := require.New(t)
	fs := NewBaseFileDecorator(afero.NewMemMapFs())

	testfile := "test.txt"

	assert.NoError(afero.WriteFile(fs, filepath.Join("themes/a/mynoblogcontent", testfile), []byte("some no content"), 0755))
	assert.NoError(afero.WriteFile(fs, filepath.Join("themes/a/myenblogcontent", testfile), []byte("some en content"), 0755))
	assert.NoError(afero.WriteFile(fs, filepath.Join("themes/a/mysvblogcontent", testfile), []byte("some sv content"), 0755))
	assert.NoError(afero.WriteFile(fs, filepath.Join("themes/a/mysvblogcontent", "other.txt"), []byte("some sv content"), 0755))

	bfs := afero.NewBasePathFs(fs, "themes/a").(*afero.BasePathFs)
	rm := []RootMapping{
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
	}

	rfs, err := NewRootMappingFs(bfs, rm...)
	assert.NoError(err)

	blog, err := rfs.Stat(filepath.FromSlash("content/blog"))
	assert.NoError(err)
	blogm := blog.(FileMetaInfo).Meta()
	assert.Equal("sv", blogm.Lang()) // Last match

	f, err := blogm.Open()
	assert.NoError(err)
	defer f.Close()
	dirs1, err := f.Readdirnames(-1)
	assert.NoError(err)
	// Union with duplicate dir names filtered.
	assert.Equal([]string{"test.txt", "test.txt", "other.txt", "test.txt"}, dirs1)

	files, err := afero.ReadDir(rfs, filepath.FromSlash("content/blog"))
	assert.NoError(err)
	assert.Equal(4, len(files))

	testfilefi := files[1]
	assert.Equal(testfile, testfilefi.Name())

	testfilem := testfilefi.(FileMetaInfo).Meta()
	assert.Equal(filepath.FromSlash("themes/a/mynoblogcontent/test.txt"), testfilem.Filename())

	tf, err := testfilem.Open()
	assert.NoError(err)
	defer tf.Close()
	c, err := ioutil.ReadAll(tf)
	assert.NoError(err)
	assert.Equal("some no content", string(c))
}

func TestRootMappingFsOs(t *testing.T) {
	assert := require.New(t)
	fs := afero.NewOsFs()

	d, err := ioutil.TempDir("", "hugo-root-mapping")
	assert.NoError(err)
	defer func() {
		os.RemoveAll(d)
	}()

	testfile := "myfile.txt"
	assert.NoError(fs.Mkdir(filepath.Join(d, "f1t"), 0755))
	assert.NoError(fs.Mkdir(filepath.Join(d, "f2t"), 0755))
	assert.NoError(fs.Mkdir(filepath.Join(d, "f3t"), 0755))
	assert.NoError(afero.WriteFile(fs, filepath.Join(d, "f2t", testfile), []byte("some content"), 0755))

	rfs, err := NewRootMappingFsFromFromTo(fs, "static/bf1", filepath.Join(d, "f1t"), "static/cf2", filepath.Join(d, "f2t"), "static/af3", filepath.Join(d, "f3t"))
	assert.NoError(err)

	fif, err := rfs.Stat(filepath.Join("static/cf2", testfile))
	assert.NoError(err)
	assert.Equal("myfile.txt", fif.Name())

	root, err := rfs.Open(filepathSeparator)
	assert.NoError(err)

	dirnames, err := root.Readdirnames(-1)
	assert.NoError(err)
	assert.Equal([]string{"bf1", "cf2", "af3"}, dirnames)
}
