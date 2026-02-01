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

package hugofs

import (
	"os"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
	"github.com/spf13/afero"
)

func TestMtimeSyncer(t *testing.T) {
	c := qt.New(t)

	srcFs := afero.NewMemMapFs()
	dstFs := afero.NewMemMapFs()

	// Create source files
	c.Assert(afero.WriteFile(srcFs, "/static/file1.txt", []byte("content1"), 0644), qt.IsNil)
	c.Assert(afero.WriteFile(srcFs, "/static/subdir/file2.txt", []byte("content2"), 0644), qt.IsNil)
	c.Assert(afero.WriteFile(srcFs, "/static/file3.txt", []byte("content3"), 0644), qt.IsNil)

	syncer := &MtimeSyncer{
		SrcFs:  srcFs,
		DestFs: dstFs,
	}

	// First sync - all files should be copied
	err := syncer.Sync("/public", "/static")
	c.Assert(err, qt.IsNil)

	// Verify files exist in destination
	content, err := afero.ReadFile(dstFs, "/public/file1.txt")
	c.Assert(err, qt.IsNil)
	c.Assert(string(content), qt.Equals, "content1")

	content, err = afero.ReadFile(dstFs, "/public/subdir/file2.txt")
	c.Assert(err, qt.IsNil)
	c.Assert(string(content), qt.Equals, "content2")
}

func TestMtimeSyncerSkipsUnchanged(t *testing.T) {
	c := qt.New(t)

	srcFs := afero.NewMemMapFs()
	dstFs := afero.NewMemMapFs()

	// Create source file
	c.Assert(afero.WriteFile(srcFs, "/static/file.txt", []byte("content"), 0644), qt.IsNil)

	syncer := &MtimeSyncer{
		SrcFs:  srcFs,
		DestFs: dstFs,
	}

	// First sync
	err := syncer.Sync("/public", "/static")
	c.Assert(err, qt.IsNil)

	// Get mtime of dst file
	dstInfo1, err := dstFs.Stat("/public/file.txt")
	c.Assert(err, qt.IsNil)

	// Wait a bit
	time.Sleep(10 * time.Millisecond)

	// Sync again without changes
	err = syncer.Sync("/public", "/static")
	c.Assert(err, qt.IsNil)

	// Dst file should not be modified (mtime unchanged)
	dstInfo2, err := dstFs.Stat("/public/file.txt")
	c.Assert(err, qt.IsNil)
	c.Assert(dstInfo2.ModTime(), qt.Equals, dstInfo1.ModTime())
}

func TestMtimeSyncerCopiesNewer(t *testing.T) {
	c := qt.New(t)

	srcFs := afero.NewMemMapFs()
	dstFs := afero.NewMemMapFs()

	// Create files with different times
	oldTime := time.Now().Add(-1 * time.Hour)
	newTime := time.Now()

	// Create dst file first (older)
	c.Assert(afero.WriteFile(dstFs, "/public/file.txt", []byte("old"), 0644), qt.IsNil)
	c.Assert(dstFs.Chtimes("/public/file.txt", oldTime, oldTime), qt.IsNil)

	// Create src file (newer)
	c.Assert(afero.WriteFile(srcFs, "/static/file.txt", []byte("new"), 0644), qt.IsNil)
	c.Assert(srcFs.Chtimes("/static/file.txt", newTime, newTime), qt.IsNil)

	syncer := &MtimeSyncer{
		SrcFs:  srcFs,
		DestFs: dstFs,
	}

	err := syncer.Sync("/public", "/static")
	c.Assert(err, qt.IsNil)

	// Verify content was updated
	content, err := afero.ReadFile(dstFs, "/public/file.txt")
	c.Assert(err, qt.IsNil)
	c.Assert(string(content), qt.Equals, "new")
}

func TestMtimeSyncerSkipsOlder(t *testing.T) {
	c := qt.New(t)

	srcFs := afero.NewMemMapFs()
	dstFs := afero.NewMemMapFs()

	oldTime := time.Now().Add(-1 * time.Hour)
	newTime := time.Now()

	// Create dst file first (newer)
	c.Assert(afero.WriteFile(dstFs, "/public/file.txt", []byte("dst"), 0644), qt.IsNil)
	c.Assert(dstFs.Chtimes("/public/file.txt", newTime, newTime), qt.IsNil)

	// Create src file (older, same size)
	c.Assert(afero.WriteFile(srcFs, "/static/file.txt", []byte("src"), 0644), qt.IsNil)
	c.Assert(srcFs.Chtimes("/static/file.txt", oldTime, oldTime), qt.IsNil)

	syncer := &MtimeSyncer{
		SrcFs:  srcFs,
		DestFs: dstFs,
	}

	err := syncer.Sync("/public", "/static")
	c.Assert(err, qt.IsNil)

	// Verify content was NOT updated (dst is newer)
	content, err := afero.ReadFile(dstFs, "/public/file.txt")
	c.Assert(err, qt.IsNil)
	c.Assert(string(content), qt.Equals, "dst")
}

func TestMtimeSyncerDelete(t *testing.T) {
	c := qt.New(t)

	srcFs := afero.NewMemMapFs()
	dstFs := afero.NewMemMapFs()

	// Create source file
	c.Assert(afero.WriteFile(srcFs, "/static/keep.txt", []byte("keep"), 0644), qt.IsNil)

	// Create orphan file in destination
	c.Assert(afero.WriteFile(dstFs, "/public/orphan.txt", []byte("orphan"), 0644), qt.IsNil)

	syncer := &MtimeSyncer{
		SrcFs:  srcFs,
		DestFs: dstFs,
		Delete: true,
	}

	err := syncer.Sync("/public", "/static")
	c.Assert(err, qt.IsNil)

	// keep.txt should exist
	_, err = dstFs.Stat("/public/keep.txt")
	c.Assert(err, qt.IsNil)

	// orphan.txt should be deleted
	_, err = dstFs.Stat("/public/orphan.txt")
	c.Assert(os.IsNotExist(err), qt.IsTrue)
}

func TestMtimeSyncerDeleteFilter(t *testing.T) {
	c := qt.New(t)

	srcFs := afero.NewMemMapFs()
	dstFs := afero.NewMemMapFs()

	// Create source file
	c.Assert(afero.WriteFile(srcFs, "/static/keep.txt", []byte("keep"), 0644), qt.IsNil)

	// Create files that should and shouldn't be deleted
	c.Assert(afero.WriteFile(dstFs, "/public/orphan.txt", []byte("orphan"), 0644), qt.IsNil)
	c.Assert(afero.WriteFile(dstFs, "/public/.gitignore", []byte("ignore"), 0644), qt.IsNil)

	syncer := &MtimeSyncer{
		SrcFs:  srcFs,
		DestFs: dstFs,
		Delete: true,
		DeleteFilter: func(f FileInfo) bool {
			return f.Name() == ".gitignore"
		},
	}

	err := syncer.Sync("/public", "/static")
	c.Assert(err, qt.IsNil)

	// orphan.txt should be deleted
	_, err = dstFs.Stat("/public/orphan.txt")
	c.Assert(os.IsNotExist(err), qt.IsTrue)

	// .gitignore should be kept
	_, err = dstFs.Stat("/public/.gitignore")
	c.Assert(err, qt.IsNil)
}

func TestMtimeSyncerDifferentSize(t *testing.T) {
	c := qt.New(t)

	srcFs := afero.NewMemMapFs()
	dstFs := afero.NewMemMapFs()

	sameTime := time.Now()

	// Create dst file
	c.Assert(afero.WriteFile(dstFs, "/public/file.txt", []byte("short"), 0644), qt.IsNil)
	c.Assert(dstFs.Chtimes("/public/file.txt", sameTime, sameTime), qt.IsNil)

	// Create src file with different size but same mtime
	c.Assert(afero.WriteFile(srcFs, "/static/file.txt", []byte("much longer content"), 0644), qt.IsNil)
	c.Assert(srcFs.Chtimes("/static/file.txt", sameTime, sameTime), qt.IsNil)

	syncer := &MtimeSyncer{
		SrcFs:  srcFs,
		DestFs: dstFs,
	}

	err := syncer.Sync("/public", "/static")
	c.Assert(err, qt.IsNil)

	// Verify content was updated (different size triggers copy)
	content, err := afero.ReadFile(dstFs, "/public/file.txt")
	c.Assert(err, qt.IsNil)
	c.Assert(string(content), qt.Equals, "much longer content")
}
