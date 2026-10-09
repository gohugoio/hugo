// Copyright 2026 The Hugo Authors. All rights reserved.
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
	"path/filepath"
	"runtime"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/spf13/afero"
)

func TestLinker(t *testing.T) {
	c := qt.New(t)
	dir := t.TempDir()
	src := filepath.Join(dir, "src.txt")
	c.Assert(os.WriteFile(src, []byte("src"), 0o644), qt.IsNil)
	if err := os.Link(src, filepath.Join(dir, "probe.txt")); err != nil {
		t.Skipf("hard links not supported: %v", err)
	}

	pubDir := filepath.Join(dir, "public")
	pubFs := NewBasePathFs(Os, pubDir)
	l := &Linker{}

	assertSameFile := func(a, b string, same bool) {
		c.Helper()
		fa, err := os.Stat(a)
		c.Assert(err, qt.IsNil)
		fb, err := os.Stat(b)
		c.Assert(err, qt.IsNil)
		c.Assert(os.SameFile(fa, fb), qt.Equals, same)
	}

	link := func(src string, dst string) bool {
		c.Helper()
		linked, err := l.Link(src, pubFs, dst)
		c.Assert(err, qt.IsNil)
		return linked
	}

	// Creates missing directories.
	c.Assert(link(src, "a/b.txt"), qt.IsTrue)
	assertSameFile(src, filepath.Join(pubDir, "a", "b.txt"), true)

	// Replaces existing files.
	c.Assert(afero.WriteFile(pubFs, "c.txt", []byte("old"), 0o644), qt.IsNil)
	c.Assert(link(src, "c.txt"), qt.IsTrue)
	assertSameFile(src, filepath.Join(pubDir, "c.txt"), true)

	// Replaces existing links to other files, leaving those intact.
	other := filepath.Join(dir, "other.txt")
	c.Assert(os.WriteFile(other, []byte("other"), 0o644), qt.IsNil)
	c.Assert(os.Link(other, filepath.Join(pubDir, "d.txt")), qt.IsNil)
	c.Assert(link(src, "d.txt"), qt.IsTrue)
	assertSameFile(src, filepath.Join(pubDir, "d.txt"), true)
	b, err := os.ReadFile(other)
	c.Assert(err, qt.IsNil)
	c.Assert(string(b), qt.Equals, "other")
	entries, err := os.ReadDir(pubDir)
	c.Assert(err, qt.IsNil)
	for _, e := range entries {
		c.Assert(e.Name(), qt.Not(qt.Contains), ".hugolink")
	}

	// Not an OS destination.
	linked, err := l.Link(src, NewBasePathFs(&afero.MemMapFs{}, "/public"), "d.txt")
	c.Assert(err, qt.IsNil)
	c.Assert(linked, qt.IsFalse)

	// Read-only source, e.g. from the module cache.
	ro := filepath.Join(dir, "ro.txt")
	c.Assert(os.WriteFile(ro, []byte("ro"), 0o444), qt.IsNil)
	c.Assert(link(ro, "ro.txt"), qt.IsTrue)
	assertSameFile(ro, filepath.Join(pubDir, "ro.txt"), true)

	// Symlink source.
	if runtime.GOOS != "windows" {
		symlink := filepath.Join(dir, "link.txt")
		c.Assert(os.Symlink(src, symlink), qt.IsNil)
		c.Assert(link(symlink, "link.txt"), qt.IsFalse)
	}

	// Missing source.
	c.Assert(link(filepath.Join(dir, "missing.txt"), "missing.txt"), qt.IsFalse)

	// Destination is a non-empty directory.
	c.Assert(pubFs.MkdirAll("dir/sub", 0o777), qt.IsNil)
	_, err = l.Link(src, pubFs, "dir")
	c.Assert(err, qt.IsNotNil)
}
