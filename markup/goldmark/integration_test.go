// Copyright 2021 The Hugo Authors. All rights reserved.
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

func TestAttributeExclusion(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
[markup.goldmark.renderer]
	unsafe = false
[markup.goldmark.parser.attribute]
	block = true
	title = true
-- content/p1.md --
---
title: "p1"
---
## Heading {class="a" onclick="alert('heading')" linenos="inline"}

> Blockquote
{class="b" ondblclick="alert('blockquote')" LINENOS="inline"}

~~~bash {id="c" onmouseover="alert('code fence')"}
foo
~~~
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

	b.AssertFileContent("public/p1/index.html", `
<h2 class="a" id="heading">
<blockquote class="b">
<div class="highlight" id="c">
	`)
}
