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

package hugolib

import (
	"testing"
	"time"

	"github.com/gohugoio/hugo/common/hmaps"
	"github.com/gohugoio/hugo/config"
)

// TestSkipUnchangedByMtime tests the --skipUnchanged optimization which
// skips rendering pages when the output file is newer than the source file.
func TestSkipUnchangedByMtime(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.org/"
disableKinds = ["taxonomy", "term", "RSS", "sitemap"]
-- content/posts/p1.md --
---
title: Post 1
---
Content for post 1.
-- content/posts/p2.md --
---
title: Post 2
---
Content for post 2.
-- layouts/_default/single.html --
Single: {{ .Title }}|{{ .Content }}
-- layouts/_default/list.html --
List: {{ .Title }}|{{ range .Pages }}{{ .Title }}|{{ end }}
-- layouts/index.html --
Home: {{ .Title }}|
`

	// Create config with skipUnchanged enabled
	cfg := config.New()
	cfg.Set("internal", hmaps.Params{
		"skipUnchanged": true,
	})

	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			BaseCfg:     cfg,
			NeedsOsFS:   true, // Need real FS for mtime comparison
		},
	).Build()

	// First build should render all pages
	b.AssertFileContent("public/posts/p1/index.html", "Single: Post 1")
	b.AssertFileContent("public/posts/p2/index.html", "Single: Post 2")
	b.AssertFileContent("public/index.html", "Home:")

	// Wait a moment to ensure mtime difference
	time.Sleep(100 * time.Millisecond)

	// Get the initial render count
	initialCount := b.counters.pageRenderCounter.Load()

	// Second build with skipUnchanged should skip pages with source files
	b.Build()

	// The render count should be less than double because some pages were skipped
	// Pages with source files (p1, p2) should be skipped
	// Pages without source files (home, list) may still render
	secondCount := b.counters.pageRenderCounter.Load()

	// Verify that files still exist and have correct content
	b.AssertFileContent("public/posts/p1/index.html", "Single: Post 1")
	b.AssertFileContent("public/posts/p2/index.html", "Single: Post 2")

	// The second build should have rendered fewer pages than the first
	// because pages with unchanged source files should be skipped
	if secondCount >= initialCount*2 {
		t.Errorf("Expected fewer renders on second build with skipUnchanged, got initial=%d, after second build=%d", initialCount, secondCount)
	}
}

// TestSkipUnchangedByMtimeRebuildsModified tests that modified pages are
// correctly rebuilt even when skipUnchanged is enabled.
func TestSkipUnchangedByMtimeRebuildsModified(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "https://example.org/"
disableKinds = ["taxonomy", "term", "RSS", "sitemap"]
-- content/posts/p1.md --
---
title: Post 1
---
Original content.
-- layouts/_default/single.html --
Single: {{ .Title }}|{{ .Content }}
-- layouts/_default/list.html --
List.
-- layouts/index.html --
Home.
`

	cfg := config.New()
	cfg.Set("internal", hmaps.Params{
		"skipUnchanged": true,
	})

	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			BaseCfg:     cfg,
			NeedsOsFS:   true,
		},
	).Build()

	b.AssertFileContent("public/posts/p1/index.html", "Original content")

	// Wait to ensure mtime changes
	time.Sleep(100 * time.Millisecond)

	// Modify the source file
	b.EditFiles("content/posts/p1.md", `---
title: Post 1
---
Modified content.
`)

	// Rebuild - the modified page should be re-rendered
	b.Build()

	// Verify the content was updated
	b.AssertFileContent("public/posts/p1/index.html", "Modified content")
}
