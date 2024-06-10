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
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gohugoio/hugo/common/para"

	"github.com/spf13/afero"

	qt "github.com/frankban/quicktest"
)

func TestWalk(t *testing.T) {
	c := qt.New(t)

	fs := NewBaseFileDecorator(afero.NewMemMapFs())

	afero.WriteFile(fs, "b.txt", []byte("content"), 0o777)
	afero.WriteFile(fs, "c.txt", []byte("content"), 0o777)
	afero.WriteFile(fs, "a.txt", []byte("content"), 0o777)

	names, err := collectPaths(fs, "")

	c.Assert(err, qt.IsNil)
	c.Assert(names, qt.DeepEquals, []string{"/a.txt", "/b.txt", "/c.txt"})
}

func TestWalkRootMappingFs(t *testing.T) {
	c := qt.New(t)

	prepare := func(c *qt.C) afero.Fs {
		fs := NewBaseFileDecorator(afero.NewMemMapFs())

		testfile := "test.txt"

		c.Assert(afero.WriteFile(fs, filepath.Join("a/b", testfile), []byte("some content"), 0o755), qt.IsNil)
		c.Assert(afero.WriteFile(fs, filepath.Join("c/d", testfile), []byte("some content"), 0o755), qt.IsNil)
		c.Assert(afero.WriteFile(fs, filepath.Join("e/f", testfile), []byte("some content"), 0o755), qt.IsNil)

		rm := []RootMapping{
			{
				From: "static/b",
				To:   "e/f",
			},
			{
				From: "static/a",
				To:   "c/d",
			},

			{
				From: "static/c",
				To:   "a/b",
			},
		}

		rfs, err := NewRootMappingFs(fs, rm...)
		c.Assert(err, qt.IsNil)
		return NewBasePathFs(rfs, "static")
	}

	c.Run("Basic", func(c *qt.C) {
		bfs := prepare(c)

		names, err := collectPaths(bfs, "")

		c.Assert(err, qt.IsNil)
		c.Assert(names, qt.DeepEquals, []string{"/a/test.txt", "/b/test.txt", "/c/test.txt"})
	})

	c.Run("Para", func(c *qt.C) {
		bfs := prepare(c)

		p := para.New(4)
		r, _ := p.Start(context.Background())

		for i := 0; i < 8; i++ {
			r.Run(func() error {
				_, err := collectPaths(bfs, "")
				if err != nil {
					return err
				}
				fi, err := bfs.Stat("b/test.txt")
				if err != nil {
					return err
				}
				meta := fi.(FileMetaInfo).Meta()
				if meta.Filename == "" {
					return errors.New("fail")
				}
				return nil
			})
		}

		c.Assert(r.Wait(), qt.IsNil)
	})
}

func collectPaths(fs afero.Fs, root string) ([]string, error) {
	var names []string

	walkFn := func(path string, info FileMetaInfo) error {
		if info.IsDir() {
			return nil
		}
		names = append(names, info.Meta().PathInfo.Path())

		return nil
	}

	w := NewWalkway(WalkwayConfig{Fs: fs, Root: root, WalkFn: walkFn, SortDirEntries: true, FailOnNotExist: true})

	err := w.Walk()

	return names, err
}

func collectFileinfos(fs afero.Fs, root string) ([]FileMetaInfo, error) {
	var fis []FileMetaInfo

	walkFn := func(path string, info FileMetaInfo) error {
		fis = append(fis, info)

		return nil
	}

	w := NewWalkway(WalkwayConfig{Fs: fs, Root: root, WalkFn: walkFn, SortDirEntries: true, FailOnNotExist: true})

	err := w.Walk()

	return fis, err
}

func BenchmarkWalk(b *testing.B) {
	c := qt.New(b)
	fs := NewBaseFileDecorator(afero.NewMemMapFs())

	writeFiles := func(dir string, numfiles int) {
		for i := 0; i < numfiles; i++ {
			filename := filepath.Join(dir, fmt.Sprintf("file%d.txt", i))
			c.Assert(afero.WriteFile(fs, filename, []byte("content"), 0o777), qt.IsNil)
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

	walkFn := func(path string, info FileMetaInfo) error {
		if info.IsDir() {
			return nil
		}

		filename := info.Meta().Filename
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
