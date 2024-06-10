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

package filecache_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/bep/logg"
	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/htesting"
	"github.com/gohugoio/hugo/hugolib"
)

// See issue #10781. That issue wouldn't have been triggered if we kept
// the empty root directories (e.g. _resources/gen/images).
// It's still an upstream Go issue that we also need to handle, but
// this is a test for the first part.
func TestPruneShouldPreserveEmptyCacheRoots(t *testing.T) {
	files := `
-- hugo.toml --
baseURL = "https://example.com"
-- content/_index.md --
---
title: "Home"
---

`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{T: t, TxtarString: files, RunGC: true, NeedsOsFS: true},
	).Build()

	_, err := b.H.BaseFs.ResourcesCache.Stat(filepath.Join("_gen", "images"))

	b.Assert(err, qt.IsNil)
}

func TestPruneImages(t *testing.T) {
	if htesting.IsCI() {
		// TODO(bep)
		t.Skip("skip flaky test on CI server")
	}
	t.Skip("skip flaky test")
	files := `
-- hugo.toml --
baseURL = "https://example.com"
[caches]
[caches.images]
maxAge = "200ms"
dir = ":resourceDir/_gen"
-- content/_index.md --
---
title: "Home"
---
-- assets/a/pixel.png --
iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==
-- layouts/index.html --
{{ warnf "HOME!" }}
{{ $img := resources.GetMatch "**.png" }}
{{ $img = $img.Resize "3x3" }}
{{ $img.RelPermalink }}



`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{T: t, TxtarString: files, Running: true, RunGC: true, NeedsOsFS: true, LogLevel: logg.LevelInfo},
	).Build()

	b.Assert(b.GCCount, qt.Equals, 0)
	b.Assert(b.H, qt.IsNotNil)

	imagesCacheDir := filepath.Join("_gen", "images")
	_, err := b.H.BaseFs.ResourcesCache.Stat(imagesCacheDir)

	b.Assert(err, qt.IsNil)

	// TODO(bep) we need a way to test full rebuilds.
	// For now, just sleep a little so the cache elements expires.
	time.Sleep(500 * time.Millisecond)

	b.RenameFile("assets/a/pixel.png", "assets/b/pixel2.png").Build()

	b.Assert(b.GCCount, qt.Equals, 1)
	// Build it again to GC the empty a dir.
	b.Build()

	_, err = b.H.BaseFs.ResourcesCache.Stat(filepath.Join(imagesCacheDir, "a"))
	b.Assert(err, qt.Not(qt.IsNil))
	_, err = b.H.BaseFs.ResourcesCache.Stat(imagesCacheDir)
	b.Assert(err, qt.IsNil)
}
