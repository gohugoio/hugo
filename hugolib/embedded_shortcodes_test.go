// Copyright 2019 The Hugo Authors. All rights reserved.
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

	"github.com/gohugoio/hugo/htesting"
)

func TestEmbeddedShortcodes(t *testing.T) {
	if !htesting.IsCI() {
		t.Skip("skip on non-CI for now")
	}

	t.Run("with theme", func(t *testing.T) {
		t.Parallel()

		files := `
-- hugo.toml --
baseURL = "https://example.com"
disableKinds = ["taxonomy", "term", "RSS", "sitemap", "robotsTXT", "page", "section"]
ignoreErrors = ["error-missing-instagram-accesstoken"]
[params]
foo = "bar"
-- content/_index.md --
---
title: "Home"
---

## Figure

{{< figure src="image.png" >}}

## Gist

{{< gist spf13 7896402 >}}

## Highlight

{{< highlight go >}}
package main
{{< /highlight >}}

## Instagram

{{< instagram BWNjjyYFxVx >}}

## Tweet

{{< tweet user="1626985695280603138" id="877500564405444608" >}}

## Vimeo

{{< vimeo 20097015 >}}

## YouTube

{{< youtube 0RKpf3rK57I >}}

## Param

Foo: {{< param foo >}}

-- layouts/index.html --
Content: {{ .Content }}|
`
		b := Test(t, files)

		b.AssertFileContent("public/index.html", `
<figure>
https://gist.github.com/spf13/7896402.js
<span style="color:#a6e22e">main</span></span>
https://t.co/X94FmYDEZJ
https://www.youtube.com/embed/0RKpf3rK57I
Foo: bar



`)
	})
}
