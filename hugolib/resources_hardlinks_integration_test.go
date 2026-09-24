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
	"os"
	"path/filepath"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/hugofs"
)

// See issue 14392.
func TestPublishResourcesHardlinks(t *testing.T) {
	files := `
-- hugo.toml --
disableKinds = ["taxonomy", "term", "rss", "sitemap"]
-- content/p1/index.md --
---
title: p1
---
-- content/p1/a.txt --
a
-- content/p1/sunset.jpg --
SUNSET_BASE64
-- content/_content.gotmpl --
{{ $.AddPage (dict "kind" "page" "path" "p2" "title" "p2") }}
{{ $.AddResource (dict "path" "p2/c.txt" "content" (dict "mediaType" "text/plain" "value" "inline")) }}
-- assets/b.txt --
b
-- layouts/page.html --
{{ with .Resources.Get "a.txt" }}{{ .RelPermalink }}|{{ end }}
{{ with .Resources.Get "c.txt" }}{{ .RelPermalink }}|{{ end }}
{{ (resources.Get "b.txt").RelPermalink }}|
{{ with .Resources.Get "sunset.jpg" }}{{ (.Resize "10x").RelPermalink }}|{{ end }}
`
	files = strings.ReplaceAll(files, "SUNSET_BASE64", getTestSunset(t))

	b := Test(t, files, TestOptWithOSFs(), TestOptRunning())

	assertSameFile := func(a, bb string, same bool) {
		b.Helper()
		fa, err := os.Stat(filepath.Join(b.Cfg.WorkingDir, filepath.FromSlash(a)))
		b.Assert(err, qt.IsNil)
		fb, err := os.Stat(filepath.Join(b.Cfg.WorkingDir, filepath.FromSlash(bb)))
		b.Assert(err, qt.IsNil)
		b.Assert(os.SameFile(fa, fb), qt.Equals, same, qt.Commentf("%s vs %s", a, bb))
	}

	b.AssertFileContent("public/p1/index.html", "/p1/a.txt|", "/b.txt|", "/p1/sunset_hu_")
	b.AssertFileContent("public/p2/index.html", "/p2/c.txt|")
	b.AssertFileContent("public/p1/a.txt", "a")
	b.AssertFileContent("public/b.txt", "b")
	b.AssertFileContent("public/p2/c.txt", "inline")
	assertSameFile("content/p1/a.txt", "public/p1/a.txt", true)
	assertSameFile("assets/b.txt", "public/b.txt", true)
	// Content from the content adapter must not link its template file.
	assertSameFile("content/_content.gotmpl", "public/p2/c.txt", false)

	// The processed image is linked from the file cache.
	cacheDir, ok := hugofs.RealFilename(b.H.ResourceSpec.FileCaches.ImageCache().Fs, "")
	b.Assert(ok, qt.IsTrue)
	published, err := filepath.Glob(filepath.Join(b.Cfg.WorkingDir, "public", "p1", "sunset_hu_*.jpg"))
	b.Assert(err, qt.IsNil)
	b.Assert(published, qt.HasLen, 1)
	cached := filepath.Join(cacheDir, "p1", filepath.Base(published[0]))
	fp, err := os.Stat(published[0])
	b.Assert(err, qt.IsNil)
	fc, err := os.Stat(cached)
	b.Assert(err, qt.IsNil)
	b.Assert(os.SameFile(fp, fc), qt.IsTrue)

	// Edit a source file and rebuild.
	b.EditFileReplaceAll("content/p1/a.txt", "a", "a edited").Build()
	b.AssertFileContent("public/p1/a.txt", "a edited")
	assertSameFile("content/p1/a.txt", "public/p1/a.txt", true)
}
