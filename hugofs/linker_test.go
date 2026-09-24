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

	// Creates missing directories.
	c.Assert(l.Link(src, pubFs, "a/b.txt"), qt.IsTrue)
	assertSameFile(src, filepath.Join(pubDir, "a", "b.txt"), true)

	// Replaces existing files.
	c.Assert(afero.WriteFile(pubFs, "c.txt", []byte("old"), 0o644), qt.IsNil)
	c.Assert(l.Link(src, pubFs, "c.txt"), qt.IsTrue)
	assertSameFile(src, filepath.Join(pubDir, "c.txt"), true)

	// Not an OS destination.
	c.Assert(l.Link(src, NewBasePathFs(&afero.MemMapFs{}, "/public"), "d.txt"), qt.IsFalse)

	// Read-only source, e.g. from the module cache.
	ro := filepath.Join(dir, "ro.txt")
	c.Assert(os.WriteFile(ro, []byte("ro"), 0o444), qt.IsNil)
	c.Assert(l.Link(ro, pubFs, "ro.txt"), qt.IsTrue)
	assertSameFile(ro, filepath.Join(pubDir, "ro.txt"), true)

	// Symlink source.
	if runtime.GOOS != "windows" {
		link := filepath.Join(dir, "link.txt")
		c.Assert(os.Symlink(src, link), qt.IsNil)
		c.Assert(l.Link(link, pubFs, "link.txt"), qt.IsFalse)
	}

	// Missing source.
	c.Assert(l.Link(filepath.Join(dir, "missing.txt"), pubFs, "missing.txt"), qt.IsFalse)
}
