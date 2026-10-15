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

package hugolib

import (
	"fmt"
	"runtime"
	"strings"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
)

// Issue 15449.
func TestRenderPagesManyErrorsShouldNotHang(t *testing.T) {
	// These only change scheduling, but makes the race hit almost every time.
	t.Setenv("HUGO_NUMWORKERMULTIPLIER", "2")
	defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(1))

	var files strings.Builder
	files.WriteString(`
-- hugo.toml --
disableKinds = ["taxonomy", "term", "rss", "sitemap"]
-- layouts/single.html --
{{ (resources.Get "missing.svg").Content }}
-- layouts/list.html --
{{ (resources.Get "missing.svg").Content }}
`)
	for i := range 50 {
		fmt.Fprintf(&files, "-- content/p%d.md --\n---\ntitle: p%d\n---\n", i, i)
	}

	for range 10 {
		done := make(chan error, 1)
		go func() {
			_, err := TestE(t, files.String())
			done <- err
		}()
		select {
		case err := <-done:
			qt.Assert(t, err, qt.IsNotNil)
			qt.Assert(t, err.Error(), qt.Contains, "failed to render pages")
		case <-time.After(10 * time.Second):
			t.Fatal("build hung")
		}
	}
}
