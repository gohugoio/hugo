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

package hstore_test

import (
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

func TestScratchSetOnce(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
-- content/p1.md --
-- content/p2.md --
-- content/p3.md --
-- content/p4.md --
-- content/p5.md --
-- layouts/all.html --
all.html:
{{ if (hugo.Store.SetOnce "once" true) }}
ONCE
{{ end }}
-- layouts/all.rss --
all.rss:
{{ if (hugo.Store.SetOnce "once" true) }}
ONCE
{{ end }}

`

	b := hugolib.Test(t, files)

	b.AssertFileContentWalk("public", 1, "ONCE")
}
