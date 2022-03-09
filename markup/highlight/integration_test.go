// Copyright 2022 The Hugo Authors. All rights reserved.
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

package highlight_test

import (
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

func TestHighlightInline(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
[markup]
[markup.highlight]
codeFences = true
noClasses = false
-- content/p1.md --
---
title: "p1"
---

Inline Classes:{{< highlight emacs "hl_inline=true" >}}abc{{< /highlight >}}:End.
Inline No Classes:{{< highlight emacs "hl_inline=true,noClasses=true" >}}abc{{< /highlight >}}:End.


-- layouts/_default/single.html --
{{ .Content }}
`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			NeedsOsFS:   false,
		},
	).Build()

	b.AssertFileContent("public/p1/index.html",
		"<p>Inline Classes:<span class=\"chroma inline\"><span class=\"cl\"><span class=\"nv\">abc</span></span></span>:End.",
		"Inline No Classes:<span style=\"\"><span>abc</span></span>:End.",
	)
}
