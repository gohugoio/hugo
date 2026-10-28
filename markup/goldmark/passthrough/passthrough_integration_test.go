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

package passthrough_test

import (
	"strings"
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

func TestPassthroughRenderHook(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
[markup.goldmark.extensions.passthrough]
enable = true
[markup.goldmark.extensions.passthrough.delimiters]
block = [['$$', '$$']]
inline = [['$', '$']]
-- content/p1.md --
---
title: "p1"
---
## LaTeX test

Some inline LaTeX 1: $a^*=x-b^*$.

Block equation that would be mangled by default parser:

$$a^*=x-b^*$$

Some inline LaTeX 2: $a^*=x-b^*$.

-- layouts/single.html --
{{ .Content }}
-- layouts/_markup/render-passthrough-block.html --
Passthrough block: {{ .Inner | safeHTML }}|{{ .Type }}|{{ .Ordinal }}:END
-- layouts/_markup/render-passthrough-inline.html --
Passthrough inline: {{ .Inner | safeHTML }}|{{ .Type }}|{{ .Ordinal }}:END

`

	b := hugolib.Test(t, files)
	b.AssertFileContent("public/p1/index.html", `
		Some inline LaTeX 1: Passthrough inline: a^*=x-b^*|inline|0:END
		Passthrough block: a^*=x-b^*|block|1:END
		Some inline LaTeX 2: Passthrough inline: a^*=x-b^*|inline|2:END

	`)
}

func TestPassthroughNoRenderHookEscapesHTML(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
[markup.goldmark.extensions.passthrough]
enable = true
[markup.goldmark.extensions.passthrough.delimiters]
block = [['$$', '$$']]
inline = [['\(', '\)']]
-- content/p1.md --
---
title: "p1"
---
Inline \(<img src=x onerror=alert(1)> a<b\).

$$
x <script>alert(2)</script>
$$
-- layouts/page.html --
{{ .Content }}
`

	t.Run("Default", func(t *testing.T) {
		t.Parallel()
		b := hugolib.Test(t, files)
		b.AssertFileContent("public/p1/index.html",
			"Inline \\(&lt;img src=x onerror=alert(1)&gt; a&lt;b\\).",
			"x &lt;script&gt;alert(2)&lt;/script&gt;",
			"! <img",
			"! <script>",
		)
	})

	t.Run("Unsafe", func(t *testing.T) {
		t.Parallel()
		b := hugolib.Test(t, strings.Replace(files, "[markup.goldmark.extensions.passthrough]", "[markup.goldmark.renderer]\nunsafe = true\n[markup.goldmark.extensions.passthrough]", 1))
		b.AssertFileContent("public/p1/index.html",
			"Inline \\(<img src=x onerror=alert(1)> a<b\\).",
			"x <script>alert(2)</script>",
		)
	})
}
