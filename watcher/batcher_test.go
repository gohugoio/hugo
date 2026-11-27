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

package watcher

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
)

// Issue 15579.
func TestBatcherRemoveBelowRenamedDir(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("inotify only")
	}
	c := qt.New(t)

	dir := t.TempDir()
	oldDir := filepath.Join(dir, "sub")
	newDir := filepath.Join(dir, "r1")
	c.Assert(os.MkdirAll(filepath.Join(oldDir, "deep"), 0o755), qt.IsNil)

	b, err := New(10*time.Millisecond, 100*time.Millisecond, false)
	c.Assert(err, qt.IsNil)
	defer b.Close()

	c.Assert(b.Add(oldDir), qt.IsNil)
	c.Assert(b.Add(filepath.Join(oldDir, "deep")), qt.IsNil)

	c.Assert(os.Rename(oldDir, newDir), qt.IsNil)
	b.RemoveBelow(oldDir)
	c.Assert(b.Add(newDir), qt.IsNil)
	c.Assert(b.Add(filepath.Join(newDir, "deep")), qt.IsNil)

	filename := filepath.Join(newDir, "deep", "f.md")
	c.Assert(os.WriteFile(filename, []byte("foo"), 0o644), qt.IsNil)

	timeout := time.After(5 * time.Second)
	for {
		select {
		case evs := <-b.Events:
			for _, ev := range evs {
				// Without RemoveBelow, the event is reported under the old path.
				c.Assert(ev.Name, qt.Not(qt.Equals), filepath.Join(oldDir, "deep", "f.md"))
				if ev.Name == filename {
					return
				}
			}
		case <-timeout:
			c.Fatal("no event for " + filename)
		}
	}
}
