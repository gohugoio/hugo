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
	"testing"
)

// Issue 15521.
func TestBothIndexFilesInContentDir(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['taxonomy', 'term', 'rss', 'sitemap']
-- layouts/all.html --
{{ .Title }}|{{ .Kind }}
-- content/inhalt/_index.md --
---
title: Section
---
-- content/inhalt/index.md --
---
title: Leaf
---
-- content/inhalt/foo.md --
---
title: Foo
---
`

	b := Test(t, files, TestOptWarn())

	b.AssertFileContent("public/inhalt/index.html", "Leaf|page")
	// The directory is treated as a leaf bundle, so neither the section page
	// nor the regular pages in it are rendered.
	b.AssertFileExists("public/inhalt/foo/index.html", false)
	b.AssertLogContains(`Directory "/inhalt" contains both "inhalt/index.md" and "inhalt/_index.md"`)
}

func TestBothIndexFilesInContentDirIgnoreLogs(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['taxonomy', 'term', 'rss', 'sitemap']
ignoreLogs = ['warning-both-index-files-in-content-dir']
-- layouts/all.html --
{{ .Title }}|{{ .Kind }}
-- content/inhalt/_index.md --
---
title: Section
---
-- content/inhalt/index.md --
---
title: Leaf
---
`

	b := Test(t, files, TestOptWarn())

	b.AssertFileContent("public/inhalt/index.html", "Leaf|page")
	b.AssertLogContains(`! Directory "/inhalt" contains both`)
}

// The warning is per language: an index file in one language does not conflict
// with an _index file in another, but two index files for the same language do.
func TestBothIndexFilesInContentDirDifferentLanguages(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['taxonomy', 'term', 'rss', 'sitemap']
defaultContentLanguage = 'de'
[languages]
[languages.de]
weight = 1
[languages.en]
weight = 2
-- layouts/all.html --
{{ .Title }}|{{ .Kind }}
-- content/inhalt/_index.de.md --
---
title: Sektion
---
-- content/inhalt/index.en.md --
---
title: Leaf
---
-- content/andere/_index.de.md --
---
title: Andere Sektion
---
-- content/andere/index.de.md --
---
title: Anderes Leaf
---
`

	b := Test(t, files, TestOptWarn())

	b.AssertLogContains(`! Directory "/inhalt" contains both`)
	b.AssertLogContains(`Directory "/andere" contains both "andere/index.de.md" and "andere/_index.de.md"`)
}

// The home page conflict has its own warning, see TestHomePageIsLeafBundle.
func TestBothIndexFilesInContentDirHomePage(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['taxonomy', 'term', 'rss', 'sitemap']
-- layouts/all.html --
{{ .Title }}|{{ .Kind }}
-- content/_index.md --
---
title: Section
---
-- content/index.md --
---
title: Leaf
---
`

	b := Test(t, files, TestOptWarn())

	b.AssertLogContains(`! Directory "/" contains both`)
	b.AssertLogContains("Using index.md in your content's root directory is usually incorrect for your home page.")
}
