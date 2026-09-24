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
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/spf13/afero"
)

func TestUnlinkOnCreateFs(t *testing.T) {
	c := qt.New(t)
	dir := t.TempDir()
	src := filepath.Join(dir, "src.txt")
	dst := filepath.Join(dir, "dst.txt")
	c.Assert(os.WriteFile(src, []byte("source"), 0o644), qt.IsNil)
	if err := os.Link(src, dst); err != nil {
		t.Skipf("hard links not supported: %v", err)
	}

	fs := NewUnlinkOnCreateFs(Os)

	f, err := fs.Create(dst)
	c.Assert(err, qt.IsNil)
	_, err = f.WriteString("created")
	c.Assert(err, qt.IsNil)
	c.Assert(f.Close(), qt.IsNil)
	b, err := os.ReadFile(src)
	c.Assert(err, qt.IsNil)
	c.Assert(string(b), qt.Equals, "source")

	c.Assert(os.Remove(dst), qt.IsNil)
	c.Assert(os.Link(src, dst), qt.IsNil)
	f, err = fs.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	c.Assert(err, qt.IsNil)
	_, err = f.WriteString("opened")
	c.Assert(err, qt.IsNil)
	c.Assert(f.Close(), qt.IsNil)
	b, err = os.ReadFile(src)
	c.Assert(err, qt.IsNil)
	c.Assert(string(b), qt.Equals, "source")
	b, err = os.ReadFile(dst)
	c.Assert(err, qt.IsNil)
	c.Assert(string(b), qt.Equals, "opened")

	// New files.
	c.Assert(afero.WriteFile(fs, filepath.Join(dir, "new.txt"), []byte("new"), 0o644), qt.IsNil)
}
