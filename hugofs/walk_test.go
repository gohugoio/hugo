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
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/gohugoio/hugo/common/hugo"

	"github.com/pkg/errors"

	"github.com/gohugoio/hugo/htesting"

	"github.com/spf13/afero"

	qt "github.com/frankban/quicktest"
)

func TestWalk(t *testing.T) {
	c := qt.New(t)

	fs := NewBaseFileDecorator(afero.NewMemMapFs())

	afero.WriteFile(fs, "b.txt", []byte("content"), 0777)
	afero.WriteFile(fs, "c.txt", []byte("content"), 0777)
	afero.WriteFile(fs, "a.txt", []byte("content"), 0777)

	names, err := collectFilenames(fs, "", "")

	c.Assert(err, qt.IsNil)
	c.Assert(names, qt.DeepEquals, []string{"a.txt", "b.txt", "c.txt"})
}

func TestWalkRootMappingFs(t *testing.T) {
	c := qt.New(t)
	fs := NewBaseFileDecorator(afero.NewMemMapFs())

	testfile := "test.txt"

	c.Assert(afero.WriteFile(fs, filepath.Join("a/b", testfile), []byte("some content"), 0755), qt.IsNil)
	c.Assert(afero.WriteFile(fs, filepath.Join("c/d", testfile), []byte("some content"), 0755), qt.IsNil)
	c.Assert(afero.WriteFile(fs, filepath.Join("e/f", testfile), []byte("some content"), 0755), qt.IsNil)

	rm := []RootMapping{
		RootMapping{
			From: "static/b",
			To:   "e/f",
		},
		RootMapping{
			From: "static/a",
			To:   "c/d",
		},

		RootMapping{
			From: "static/c",
			To:   "a/b",
		},
	}

	rfs, err := NewRootMappingFs(fs, rm...)
	c.Assert(err, qt.IsNil)
	bfs := afero.NewBasePathFs(rfs, "static")

	names, err := collectFilenames(bfs, "", "")

	c.Assert(err, qt.IsNil)
	c.Assert(names, qt.DeepEquals, []string{"a/test.txt", "b/test.txt", "c/test.txt"})

}

func skipSymlink() bool {
	return runtime.GOOS == "windows" && os.Getenv("CI") == ""
}

func TestWalkSymbolicLink(t *testing.T) {
	if skipSymlink() {
		t.Skip("Skip; os.Symlink needs administrator rights on Windows")
	}
	c := qt.New(t)
	workDir, clean, err := htesting.CreateTempDir(Os, "hugo-walk-sym")
	c.Assert(err, qt.IsNil)
	defer clean()
	wd, _ := os.Getwd()
	defer func() {
		os.Chdir(wd)
	}()

	fs := NewBaseFileDecorator(Os)

	blogDir := filepath.Join(workDir, "blog")
	docsDir := filepath.Join(workDir, "docs")
	blogReal := filepath.Join(blogDir, "real")
	blogRealSub := filepath.Join(blogReal, "sub")
	c.Assert(os.MkdirAll(blogRealSub, 0777), qt.IsNil)
	c.Assert(os.MkdirAll(docsDir, 0777), qt.IsNil)
	afero.WriteFile(fs, filepath.Join(blogRealSub, "a.txt"), []byte("content"), 0777)
	afero.WriteFile(fs, filepath.Join(docsDir, "b.txt"), []byte("content"), 0777)

	os.Chdir(blogDir)
	c.Assert(os.Symlink("real", "symlinked"), qt.IsNil)
	os.Chdir(blogReal)
	c.Assert(os.Symlink("../real", "cyclic"), qt.IsNil)
	os.Chdir(docsDir)
	c.Assert(os.Symlink("../blog/real/cyclic", "docsreal"), qt.IsNil)

	t.Run("OS Fs", func(t *testing.T) {
		c := qt.New(t)

		names, err := collectFilenames(fs, workDir, workDir)
		c.Assert(err, qt.IsNil)

		c.Assert(names, qt.DeepEquals, []string{"blog/real/sub/a.txt", "docs/b.txt"})
	})

	t.Run("BasePath Fs", func(t *testing.T) {
		if hugo.GoMinorVersion() < 12 {
			// https://github.com/golang/go/issues/30520
			// This is fixed in Go 1.13 and in the latest Go 1.12
			t.Skip("skip this for Go <= 1.11 due to a bug in Go's stdlib")

		}
		c := qt.New(t)

		docsFs := afero.NewBasePathFs(fs, docsDir)

		names, err := collectFilenames(docsFs, "", "")
		c.Assert(err, qt.IsNil)

		// Note: the docsreal folder is considered cyclic when walking from the root, but this works.
		c.Assert(names, qt.DeepEquals, []string{"b.txt", "docsreal/sub/a.txt"})
	})

}

func collectFilenames(fs afero.Fs, base, root string) ([]string, error) {
	var names []string

	walkFn := func(path string, info FileMetaInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		filename := info.Meta().Path()
		filename = filepath.ToSlash(filename)

		names = append(names, filename)

		return nil
	}

	w := NewWalkway(WalkwayConfig{Fs: fs, BasePath: base, Root: root, WalkFn: walkFn})

	err := w.Walk()

	return names, err

}

func collectFileinfos(fs afero.Fs, base, root string) ([]FileMetaInfo, error) {
	var fis []FileMetaInfo

	walkFn := func(path string, info FileMetaInfo, err error) error {
		if err != nil {
			return err
		}

		fis = append(fis, info)

		return nil
	}

	w := NewWalkway(WalkwayConfig{Fs: fs, BasePath: base, Root: root, WalkFn: walkFn})

	err := w.Walk()

	return fis, err

}

func BenchmarkWalk(b *testing.B) {
	c := qt.New(b)
	fs := NewBaseFileDecorator(afero.NewMemMapFs())

	writeFiles := func(dir string, numfiles int) {
		for i := 0; i < numfiles; i++ {
			filename := filepath.Join(dir, fmt.Sprintf("file%d.txt", i))
			c.Assert(afero.WriteFile(fs, filename, []byte("content"), 0777), qt.IsNil)
		}
	}

	const numFilesPerDir = 20

	writeFiles("root", numFilesPerDir)
	writeFiles("root/l1_1", numFilesPerDir)
	writeFiles("root/l1_1/l2_1", numFilesPerDir)
	writeFiles("root/l1_1/l2_2", numFilesPerDir)
	writeFiles("root/l1_2", numFilesPerDir)
	writeFiles("root/l1_2/l2_1", numFilesPerDir)
	writeFiles("root/l1_3", numFilesPerDir)

	walkFn := func(path string, info FileMetaInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		filename := info.Meta().Filename()
		if !strings.HasPrefix(filename, "root") {
			return errors.New(filename)
		}

		return nil
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := NewWalkway(WalkwayConfig{Fs: fs, Root: "root", WalkFn: walkFn})

		if err := w.Walk(); err != nil {
			b.Fatal(err)
		}
	}

}
