// Copyright 2025 The Hugo Authors. All rights reserved.
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

package modules_test

import (
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/hugolib"
)

func TestMountsProjectDefaults(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
[languages]
[languages.en]
weight = 1
[languages.sv]
weight = 2

[module]
[[module.mounts]]
source = 'content'
target = 'content'
lang = 'en'
-- content/p1.md --

`
	b := hugolib.Test(t, files)

	b.Assert(len(b.H.Configs.Modules), qt.Equals, 1)
	projectMod := b.H.Configs.Modules[0]
	b.Assert(projectMod.Path(), qt.Equals, "project")
	mounts := projectMod.Mounts()
	b.Assert(len(mounts), qt.Equals, 7)
	contentMount := mounts[0]
	b.Assert(contentMount.Source, qt.Equals, "content")
	b.Assert(contentMount.Sites.Matrix.Languages, qt.DeepEquals, []string{"en"})
}

func TestMountsLangIsDeprecated(t *testing.T) {
	t.Parallel()
	files := `
-- hugo.toml --
[module]
[[module.mounts]]
source = 'content'
target = 'content'
lang = 'en'
-- layouts/all.html --
All.
`

	b := hugolib.Test(t, files, hugolib.TestOptInfo())
	b.AssertLogContains("deprecated")
}
