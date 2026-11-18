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

package commands

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
	"github.com/fsnotify/fsnotify"
	"github.com/gohugoio/hugo/internal/warpc"
)

func TestErrStateShouldRecover(t *testing.T) {
	c := qt.New(t)

	broken := fmt.Errorf("render failed: %w", warpc.ErrBroken)
	other := errors.New("template error")

	var e hugoBuilderErrState

	steps := []struct {
		err  error
		want bool
	}{
		// timeout, ok, ok, timeout.
		{broken, true},
		{nil, false},
		{nil, false},
		{broken, true},
		// The recovery rebuild fails the same way: don't loop.
		{broken, false},
		// But recover again on the next, unrelated failure.
		{broken, true},
		// Any other build result resets.
		{other, false},
		{broken, true},
	}

	for i, step := range steps {
		c.Assert(e.shouldRecover(step.err), qt.Equals, step.want, qt.Commentf("step %d", i))
	}
}

// See issue 15579.
func TestDetectRenamedDirectoriesInWatchBatch(t *testing.T) {
	c := qt.New(t)

	mkdirs := func(names ...string) {
		for _, name := range names {
			c.Assert(os.MkdirAll(name, 0o755), qt.IsNil)
		}
	}

	// A rename is reported twice for the old path: once by the parent watch
	// and once by the renamed directory's own watch.
	base := t.TempDir()
	oldDir := filepath.Join(base, "sub")
	newDir := filepath.Join(base, "r1")
	mkdirs(filepath.Join(oldDir, "deep"))
	c.Assert(os.Rename(oldDir, newDir), qt.IsNil)
	c.Assert(renamedDirsInBatch([]fsnotify.Event{
		{Name: oldDir, Op: fsnotify.Rename},
		{Name: newDir, Op: fsnotify.Create},
		{Name: oldDir, Op: fsnotify.Rename},
	}), qt.DeepEquals, map[string]string{newDir: oldDir})

	// Two renames sharing a parent in one batch are ambiguous.
	base = t.TempDir()
	a1, a2 := filepath.Join(base, "a1"), filepath.Join(base, "a2")
	b1, b2 := filepath.Join(base, "b1"), filepath.Join(base, "b2")
	mkdirs(a1, a2)
	c.Assert(os.Rename(a1, b1), qt.IsNil)
	c.Assert(os.Rename(a2, b2), qt.IsNil)
	c.Assert(renamedDirsInBatch([]fsnotify.Event{
		{Name: a1, Op: fsnotify.Rename},
		{Name: b1, Op: fsnotify.Create},
		{Name: a2, Op: fsnotify.Rename},
		{Name: b2, Op: fsnotify.Create},
	}), qt.DeepEquals, map[string]string{})

	// A single pair across parents is unambiguous.
	base = t.TempDir()
	p1, p2 := filepath.Join(base, "p1"), filepath.Join(base, "p2")
	oldDir = filepath.Join(p1, "a")
	newDir = filepath.Join(p2, "a")
	mkdirs(p1, p2, oldDir)
	c.Assert(os.Rename(oldDir, newDir), qt.IsNil)
	c.Assert(renamedDirsInBatch([]fsnotify.Event{
		{Name: oldDir, Op: fsnotify.Rename},
		{Name: newDir, Op: fsnotify.Create},
	}), qt.DeepEquals, map[string]string{newDir: oldDir})

	// A created file does not count as a moved-in directory.
	base = t.TempDir()
	oldDir = filepath.Join(base, "sub")
	newDir = filepath.Join(base, "r1")
	mkdirs(oldDir)
	c.Assert(os.Rename(oldDir, newDir), qt.IsNil)
	f := filepath.Join(base, "top.md")
	c.Assert(os.WriteFile(f, []byte("x"), 0o644), qt.IsNil)
	c.Assert(renamedDirsInBatch([]fsnotify.Event{
		{Name: oldDir, Op: fsnotify.Rename},
		{Name: newDir, Op: fsnotify.Create},
		{Name: f, Op: fsnotify.Create},
	}), qt.DeepEquals, map[string]string{newDir: oldDir})

	// A Rename event for a path that still exists is not a rename.
	base = t.TempDir()
	oldDir = filepath.Join(base, "sub")
	newDir = filepath.Join(base, "r1")
	mkdirs(oldDir, newDir)
	c.Assert(renamedDirsInBatch([]fsnotify.Event{
		{Name: oldDir, Op: fsnotify.Rename},
		{Name: newDir, Op: fsnotify.Create},
	}), qt.DeepEquals, map[string]string{})
}

// See issue 15579.
func TestStaleWatchPathAfterRename(t *testing.T) {
	c := qt.New(t)
	renamed := map[string]string{
		filepath.Join("x", "r1"): filepath.Join("x", "sub"),
	}
	c.Assert(staleWatchPath(renamed, filepath.Join("x", "r1")), qt.Equals, filepath.Join("x", "sub"))
	c.Assert(staleWatchPath(renamed, filepath.Join("x", "r1", "deep")), qt.Equals, filepath.Join("x", "sub", "deep"))
	c.Assert(staleWatchPath(renamed, filepath.Join("x", "r1", "deep", "deeper")), qt.Equals, filepath.Join("x", "sub", "deep", "deeper"))
	c.Assert(staleWatchPath(renamed, filepath.Join("x", "other")), qt.Equals, "")
	c.Assert(staleWatchPath(renamed, filepath.Join("x")), qt.Equals, "")
	c.Assert(staleWatchPath(renamed, filepath.Join("x", "r1x")), qt.Equals, "")
	c.Assert(staleWatchPath(map[string]string{}, filepath.Join("x", "r1")), qt.Equals, "")
}

// See issue 15579.
func TestWatchRepairAfterDirectoryRename(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("the stale watch left behind by a rename is inotify-specific")
	}
	c := qt.New(t)
	base := t.TempDir()
	oldDir := filepath.Join(base, "a")
	newDir := filepath.Join(base, "b")
	oldSub := filepath.Join(oldDir, "sub")
	newSub := filepath.Join(newDir, "sub")
	c.Assert(os.MkdirAll(oldSub, 0o755), qt.IsNil)
	c.Assert(os.WriteFile(filepath.Join(oldSub, "file.txt"), []byte("v1"), 0o644), qt.IsNil)

	w, err := fsnotify.NewWatcher()
	c.Assert(err, qt.IsNil)
	defer w.Close()
	go func() {
		for range w.Errors {
		}
	}()
	for _, p := range []string{base, oldDir, oldSub} {
		c.Assert(w.Add(p), qt.IsNil)
	}

	c.Assert(os.Rename(oldDir, newDir), qt.IsNil)
	var evs []fsnotify.Event
	sawRename, sawCreate := false, false
	deadline := time.After(3 * time.Second)
collect:
	for !(sawRename && sawCreate) {
		select {
		case ev := <-w.Events:
			evs = append(evs, ev)
			sawRename = sawRename || ev.Has(fsnotify.Rename)
			sawCreate = sawCreate || ev.Has(fsnotify.Create)
		case <-deadline:
			break collect
		}
	}
	renamed := renamedDirsInBatch(evs)
	c.Assert(renamed, qt.DeepEquals, map[string]string{newDir: oldDir})

	// Re-add the renamed tree the way handleEvents does: drop the stale
	// watch first, or the Add is a no-op and the old path sticks.
	for _, p := range []string{newDir, newSub} {
		if old := staleWatchPath(renamed, p); old != "" {
			_ = w.Remove(old)
		}
		c.Assert(w.Add(p), qt.IsNil)
	}
	watched := make(map[string]bool)
	for _, p := range w.WatchList() {
		watched[p] = true
	}
	c.Assert(watched[newSub], qt.IsTrue)
	c.Assert(watched[oldSub], qt.IsFalse)

	newFile := filepath.Join(newSub, "file.txt")
	oldFile := filepath.Join(oldSub, "file.txt")
	c.Assert(os.WriteFile(newFile, []byte("v2"), 0o644), qt.IsNil)
	sawNew := false
	timeout := time.After(5 * time.Second)
events:
	for {
		select {
		case ev := <-w.Events:
			if ev.Has(fsnotify.Write) && ev.Name == newFile {
				sawNew = true
				break events
			}
			if ev.Has(fsnotify.Write) && ev.Name == oldFile {
				break events
			}
		case <-timeout:
			break events
		}
	}
	c.Assert(sawNew, qt.IsTrue)
}
