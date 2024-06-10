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

package goldmark_test

import (
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

func TestTableOfContents(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['home','rss','section','sitemap','taxonomy','term']
enableEmoji = false

[markup.tableOfContents]
startLevel = 2
endLevel = 4
ordered = false

[markup.goldmark.extensions]
strikethrough = false

[markup.goldmark.extensions.typographer]
disable = true

[markup.goldmark.parser]
autoHeadingID = false
autoHeadingIDType = 'github'

[markup.goldmark.renderer]
unsafe = false
xhtml = false
-- layouts/_default/single.html --
{{ .TableOfContents }}
-- content/p1.md --
---
title: p1 (basic)
---
# Title
## Section 1
### Section 1.1
### Section 1.2
#### Section 1.2.1
##### Section 1.2.1.1
-- content/p2.md --
---
title: p2 (markdown)
---
## Some *emphasized* text
## Some ` + "`" + `inline` + "`" + ` code
## Something to escape A < B && C > B
---
-- content/p3.md --
---
title: p3 (image)
---
## An image ![kitten](a.jpg)
-- content/p4.md --
---
title: p4 (raw html)
---
## Some <span>raw</span> HTML
-- content/p5.md --
---
title: p5 (typographer)
---
## Some "typographer" markup
-- content/p6.md --
---
title: p6 (strikethrough)
---
## Some ~~deleted~~ text
-- content/p7.md --
---
title: p7 (emoji)
---
## A :snake: emoji
`

	b := hugolib.Test(t, files)

	// basic
	b.AssertFileContentExact("public/p1/index.html", `<nav id="TableOfContents">
  <ul>
    <li><a href="#">Section 1</a>
      <ul>
        <li><a href="#">Section 1.1</a></li>
        <li><a href="#">Section 1.2</a>
          <ul>
            <li><a href="#">Section 1.2.1</a></li>
          </ul>
        </li>
      </ul>
    </li>
  </ul>
</nav>`)

	// markdown
	b.AssertFileContent("public/p2/index.html", `<nav id="TableOfContents">
<li><a href="#">Some <em>emphasized</em> text</a></li>
<li><a href="#">Some <code>inline</code> code</a></li>
<li><a href="#">Something to escape A &lt; B &amp;&amp; C &gt; B</a></li>
`)

	// image
	b.AssertFileContent("public/p3/index.html", `
<li><a href="#">An image <img src="a.jpg" alt="kitten"></a></li>
`)

	// raw html
	b.AssertFileContent("public/p4/index.html", `
<li><a href="#">Some <!-- raw HTML omitted -->raw<!-- raw HTML omitted --> HTML</a></li>
`)

	// typographer
	b.AssertFileContent("public/p5/index.html", `
<li><a href="#">Some &quot;typographer&quot; markup</a></li>
`)

	// strikethrough
	b.AssertFileContent("public/p6/index.html", `
<li><a href="#">Some ~~deleted~~ text</a></li>
	`)

	// emoji
	b.AssertFileContent("public/p7/index.html", `
<li><a href="#">A :snake: emoji</a></li>
		`)
}

func TestTableOfContentsAdvanced(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['home','rss','section','sitemap','taxonomy','term']
enableEmoji = true

[markup.tableOfContents]
startLevel = 2
endLevel = 3
ordered = true

[markup.goldmark.extensions]
strikethrough = true

[markup.goldmark.extensions.typographer]
disable = false

[markup.goldmark.parser]
autoHeadingID = true
autoHeadingIDType = 'github'

[markup.goldmark.renderer]
unsafe = true
xhtml = true
-- layouts/_default/single.html --
{{ .TableOfContents }}
-- content/p1.md --
---
title: p1 (basic)
---
# Title
## Section 1
### Section 1.1
### Section 1.2
#### Section 1.2.1
##### Section 1.2.1.1
-- content/p2.md --
---
title: p2 (markdown)
---
## Some *emphasized* text
## Some ` + "`" + `inline` + "`" + ` code
## Something to escape A < B && C > B
---
-- content/p3.md --
---
title: p3 (image)
---
## An image ![kitten](a.jpg)
-- content/p4.md --
---
title: p4 (raw html)
---
## Some <span>raw</span> HTML
-- content/p5.md --
---
title: p5 (typographer)
---
## Some "typographer" markup
-- content/p6.md --
---
title: p6 (strikethrough)
---
## Some ~~deleted~~ text
-- content/p7.md --
---
title: p7 (emoji)
---
## A :snake: emoji
`

	b := hugolib.Test(t, files)

	// basic
	b.AssertFileContentExact("public/p1/index.html", `<nav id="TableOfContents">
  <ol>
    <li><a href="#section-1">Section 1</a>
      <ol>
        <li><a href="#section-11">Section 1.1</a></li>
        <li><a href="#section-12">Section 1.2</a></li>
      </ol>
    </li>
  </ol>
</nav>`)

	// markdown
	b.AssertFileContent("public/p2/index.html", `<nav id="TableOfContents">
<li><a href="#some-emphasized-text">Some <em>emphasized</em> text</a></li>
<li><a href="#some-inline-code">Some <code>inline</code> code</a></li>
<li><a href="#something-to-escape-a--b--c--b">Something to escape A &lt; B &amp;&amp; C &gt; B</a></li>
`)

	// image
	b.AssertFileContent("public/p3/index.html", `
<li><a href="#an-image-kittenajpg">An image <img src="a.jpg" alt="kitten" /></a></li>
`)

	// raw html
	b.AssertFileContent("public/p4/index.html", `
<li><a href="#some-spanrawspan-html">Some <span>raw</span> HTML</a></li>
`)

	// typographer
	b.AssertFileContent("public/p5/index.html", `
<li><a href="#some-typographer-markup">Some &ldquo;typographer&rdquo; markup</a></li>
`)

	// strikethrough
	b.AssertFileContent("public/p6/index.html", `
<li><a href="#some-deleted-text">Some <del>deleted</del> text</a></li>
`)

	// emoji
	b.AssertFileContent("public/p7/index.html", `
<li><a href="#a-snake-emoji">A &#x1f40d; emoji</a></li>
`)
}
