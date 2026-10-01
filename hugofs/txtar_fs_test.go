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
	"path/filepath"
	"sort"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/afero"
)

func TestTxtarFs(t *testing.T) {
	c := qt.New(t)

	workingDir := filepath.FromSlash("/mysite")
	abs := func(s string) string { return filepath.Join(workingDir, filepath.FromSlash(s)) }
	txtarFilename := abs(TxtarFilename)

	base := afero.NewMemMapFs()
	writeTxtar := func(s string) {
		c.Assert(afero.WriteFile(base, txtarFilename, []byte(s), 0o666), qt.IsNil)
	}
	readFile := func(fs afero.Fs, name string) string {
		b, err := afero.ReadFile(fs, abs(name))
		c.Assert(err, qt.IsNil)
		return string(b)
	}
	eventStrings := func(evs []fsnotify.Event) []string {
		var ss []string
		for _, ev := range evs {
			ss = append(ss, ev.Op.String()+" "+filepath.ToSlash(strings.TrimPrefix(ev.Name, workingDir)))
		}
		sort.Strings(ss)
		return ss
	}

	fs, err := NewTxtarFsIfExists(base, workingDir)
	c.Assert(err, qt.IsNil)
	c.Assert(fs, qt.Equals, base)

	c.Assert(afero.WriteFile(base, abs("content/ondisk.md"), []byte("disk"), 0o666), qt.IsNil)
	c.Assert(afero.WriteFile(base, abs("content/both.md"), []byte("disk"), 0o666), qt.IsNil)

	writeTxtar(`
-- hugo.toml --
baseURL = "https://example.org/"
-- content/both.md --
txtar
-- content/p1.md --
p1
-- content/sub/p2.md --
p2
`)

	fs, err = NewTxtarFsIfExists(base, workingDir)
	c.Assert(err, qt.IsNil)
	tfs, ok := fs.(*TxtarFs)
	c.Assert(ok, qt.IsTrue)
	c.Assert(tfs.Filename(), qt.Equals, txtarFilename)

	c.Assert(readFile(fs, "hugo.toml"), qt.Equals, "baseURL = \"https://example.org/\"\n")
	c.Assert(readFile(fs, "content/p1.md"), qt.Equals, "p1\n")
	c.Assert(readFile(fs, "content/ondisk.md"), qt.Equals, "disk")
	// Files on disk take precedence.
	c.Assert(readFile(fs, "content/both.md"), qt.Equals, "disk")

	names, err := afero.ReadDir(fs, abs("content"))
	c.Assert(err, qt.IsNil)
	var entries []string
	for _, fi := range names {
		entries = append(entries, fi.Name())
	}
	sort.Strings(entries)
	c.Assert(entries, qt.DeepEquals, []string{"both.md", "ondisk.md", "p1.md", "sub"})

	// Writes go to base.
	c.Assert(afero.WriteFile(fs, abs("content/new.md"), []byte("new"), 0o666), qt.IsNil)
	c.Assert(readFile(base, "content/new.md"), qt.Equals, "new")

	// No changes.
	evs, err := tfs.Reload()
	c.Assert(err, qt.IsNil)
	c.Assert(evs, qt.HasLen, 0)

	// Change, add and remove; replace a directory with a file.
	writeTxtar(`
-- hugo.toml --
baseURL = "https://example.org/"
-- content/p1.md --
p1 changed
-- content/p3.md --
p3
-- content/sub --
sub
`)

	evs, err = tfs.Reload()
	c.Assert(err, qt.IsNil)
	c.Assert(eventStrings(evs), qt.DeepEquals, []string{
		"CREATE /content/p3.md",
		"CREATE /content/sub",
		"REMOVE /content/both.md",
		"REMOVE /content/sub/p2.md",
		"WRITE /content/p1.md",
	})
	c.Assert(readFile(fs, "content/p1.md"), qt.Equals, "p1 changed\n")
	c.Assert(readFile(fs, "content/sub"), qt.Equals, "sub\n")
	_, err = fs.Stat(abs("content/sub/p2.md"))
	c.Assert(err, qt.IsNotNil)

	// Archive removed.
	c.Assert(base.Remove(txtarFilename), qt.IsNil)
	evs, err = tfs.Reload()
	c.Assert(err, qt.IsNil)
	c.Assert(evs, qt.HasLen, 4)
	for _, ev := range evs {
		c.Assert(ev.Op, qt.Equals, fsnotify.Remove)
	}
	_, err = fs.Stat(abs("content/sub"))
	c.Assert(err, qt.IsNotNil)
	c.Assert(readFile(fs, "content/ondisk.md"), qt.Equals, "disk")

	// Invalid filename.
	writeTxtar("-- ../escape.md --\n")
	_, err = tfs.Reload()
	c.Assert(err, qt.ErrorMatches, `.*invalid filename "../escape.md"`)
}
