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

package hugolib

import (
	"testing"
)

func TestInternalTemplatesImage(t *testing.T) {
	files := `
-- hugo.toml --
baseURL = "https://example.org"

[params]
images=["siteimg1.jpg", "siteimg2.jpg"]
-- content/mybundle/index.md --
---
title: My Bundle
date: 2021-02-26T18:02:00-01:00
lastmod: 2021-05-22T19:25:00-01:00
---
-- content/mypage/index.md --
---
title: My Page
images: ["pageimg1.jpg", "pageimg2.jpg", "https://example.local/logo.png", "sample.jpg"]
date: 2021-02-26T18:02:00+01:00
lastmod: 2021-05-22T19:25:00+01:00
---
-- content/mysite.md --
---
title: My Site
---
-- layouts/_default/single.html --

{{ template "_internal/twitter_cards.html" . }}
{{ template "_internal/opengraph.html" . }}
{{ template "_internal/schema.html" . }}

-- content/mybundle/featured-sunset.jpg --
/9j/4AAQSkZJRgABAQEAYABgAAD/2wBDAAgGBgcGBQgHBwcJCQgKDBQNDAsLDBkSEw8UHRofHh0aHBwgJC4nICIsIxwcKDcpLDAxNDQ0Hyc5PTgyPC4zNDL/2wBDAQkJCQwLDBgNDRgyIRwhMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjL/wAARCAABAAEDASIAAhEBAxEB/8QAFQABAQAAAAAAAAAAAAAAAAAAAAD/xAAUEAEAAAAAAAAAAAAAAAAAAAAA/8QAFAEBAAAAAAAAAAAAAAAAAAAAAP/EABQRAQAAAAAAAAAAAAAAAAAAAAD/2gAMAwEAAhEDEQA/AAT8AAAAA//Z
-- content/mypage/sample.jpg --
/9j/4AAQSkZJRgABAQEAYABgAAD/2wBDAAgGBgcGBQgHBwcJCQgKDBQNDAsLDBkSEw8UHRofHh0aHBwgJC4nICIsIxwcKDcpLDAxNDQ0Hyc5PTgyPC4zNDL/2wBDAQkJCQwLDBgNDRgyIRwhMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjL/wAARCAABAAEDASIAAhEBAxEB/8QAFQABAQAAAAAAAAAAAAAAAAAAAAD/xAAUEAEAAAAAAAAAAAAAAAAAAAAA/8QAFAEBAAAAAAAAAAAAAAAAAAAAAP/EABQRAQAAAAAAAAAAAAAAAAAAAAD/2gAMAwEAAhEDEQA/AAT8AAAAA//Z
`

	b := Test(t, files)

	b.AssertFileContent("public/mybundle/index.html", `
<meta name="twitter:image" content="https://example.org/mybundle/featured-sunset.jpg">
<meta name="twitter:title" content="My Bundle">
<meta property="og:title" content="My Bundle">
<meta property="og:url" content="https://example.org/mybundle/">
<meta property="og:image" content="https://example.org/mybundle/featured-sunset.jpg">
<meta property="article:published_time" content="2021-02-26T18:02:00-01:00">
<meta property="article:modified_time" content="2021-05-22T19:25:00-01:00">
<meta itemprop="name" content="My Bundle">
<meta itemprop="image" content="https://example.org/mybundle/featured-sunset.jpg">
<meta itemprop="datePublished" content="2021-02-26T18:02:00-01:00">
<meta itemprop="dateModified" content="2021-05-22T19:25:00-01:00">

`)
	b.AssertFileContent("public/mypage/index.html", `
<meta name="twitter:image" content="https://example.org/pageimg1.jpg">
<meta property="og:image" content="https://example.org/pageimg1.jpg">
<meta property="og:image" content="https://example.org/pageimg2.jpg">
<meta property="og:image" content="https://example.local/logo.png">
<meta property="og:image" content="https://example.org/mypage/sample.jpg">
<meta property="article:published_time" content="2021-02-26T18:02:00+01:00">
<meta property="article:modified_time" content="2021-05-22T19:25:00+01:00">
<meta itemprop="image" content="https://example.org/pageimg1.jpg">
<meta itemprop="image" content="https://example.org/pageimg2.jpg">
<meta itemprop="image" content="https://example.local/logo.png">
<meta itemprop="image" content="https://example.org/mypage/sample.jpg">
<meta itemprop="datePublished" content="2021-02-26T18:02:00+01:00">
<meta itemprop="dateModified" content="2021-05-22T19:25:00+01:00">
`)
	b.AssertFileContent("public/mysite/index.html", `
<meta name="twitter:image" content="https://example.org/siteimg1.jpg">
<meta property="og:image" content="https://example.org/siteimg1.jpg">
<meta itemprop="image" content="https://example.org/siteimg1.jpg">
`)
}

func TestEmbeddedPaginationTemplate(t *testing.T) {
	t.Parallel()

	test := func(variant string, expectedOutput string) {
		files := `
-- hugo.toml --
pagination.pagerSize = 1
-- content/s1/p01.md --
---
title: p01
---
-- content/s1/p02.md --
---
title: p02
---
-- content/s1/p03.md --
---
title: p03
---
-- content/s1/p04.md --
---
title: p04
---
-- content/s1/p05.md --
---
title: p05
---
-- content/s1/p06.md --
---
title: p06
---
-- content/s1/p07.md --
---
title: p07
---
-- content/s1/p08.md --
---
title: p08
---
-- content/s1/p09.md --
---
title: p09
---
-- content/s1/p10.md --
---
title: p10
---
-- layouts/index.html --
{{ .Paginate (where site.RegularPages "Section" "s1") }}` + variant + `
`
		b := Test(t, files)
		b.AssertFileContent("public/index.html", expectedOutput)
	}

	expectedOutputDefaultFormat := "Pager 1\n    <ul class=\"pagination pagination-default\">\n      <li class=\"page-item disabled\">\n        <a aria-disabled=\"true\" aria-label=\"First\" class=\"page-link\" role=\"button\" tabindex=\"-1\"><span aria-hidden=\"true\">&laquo;&laquo;</span></a>\n      </li>\n      <li class=\"page-item disabled\">\n        <a aria-disabled=\"true\" aria-label=\"Previous\" class=\"page-link\" role=\"button\" tabindex=\"-1\"><span aria-hidden=\"true\">&laquo;</span></a>\n      </li>\n      <li class=\"page-item active\">\n        <a aria-current=\"page\" aria-label=\"Page 1\" class=\"page-link\" role=\"button\">1</a>\n      </li>\n      <li class=\"page-item\">\n        <a href=\"/page/2/\" aria-label=\"Page 2\" class=\"page-link\" role=\"button\">2</a>\n      </li>\n      <li class=\"page-item\">\n        <a href=\"/page/3/\" aria-label=\"Page 3\" class=\"page-link\" role=\"button\">3</a>\n      </li>\n      <li class=\"page-item\">\n        <a href=\"/page/4/\" aria-label=\"Page 4\" class=\"page-link\" role=\"button\">4</a>\n      </li>\n      <li class=\"page-item\">\n        <a href=\"/page/5/\" aria-label=\"Page 5\" class=\"page-link\" role=\"button\">5</a>\n      </li>\n      <li class=\"page-item\">\n        <a href=\"/page/2/\" aria-label=\"Next\" class=\"page-link\" role=\"button\"><span aria-hidden=\"true\">&raquo;</span></a>\n      </li>\n      <li class=\"page-item\">\n        <a href=\"/page/10/\" aria-label=\"Last\" class=\"page-link\" role=\"button\"><span aria-hidden=\"true\">&raquo;&raquo;</span></a>\n      </li>\n    </ul>"
	expectedOutputTerseFormat := "Pager 1\n    <ul class=\"pagination pagination-terse\">\n      <li class=\"page-item active\">\n        <a aria-current=\"page\" aria-label=\"Page 1\" class=\"page-link\" role=\"button\">1</a>\n      </li>\n      <li class=\"page-item\">\n        <a href=\"/page/2/\" aria-label=\"Page 2\" class=\"page-link\" role=\"button\">2</a>\n      </li>\n      <li class=\"page-item\">\n        <a href=\"/page/3/\" aria-label=\"Page 3\" class=\"page-link\" role=\"button\">3</a>\n      </li>\n      <li class=\"page-item\">\n        <a href=\"/page/2/\" aria-label=\"Next\" class=\"page-link\" role=\"button\"><span aria-hidden=\"true\">&raquo;</span></a>\n      </li>\n      <li class=\"page-item\">\n        <a href=\"/page/10/\" aria-label=\"Last\" class=\"page-link\" role=\"button\"><span aria-hidden=\"true\">&raquo;&raquo;</span></a>\n      </li>\n    </ul>"

	variant := `{{ template "_internal/pagination.html" . }}`
	test(variant, expectedOutputDefaultFormat)

	variant = `{{ template "_internal/pagination.html" (dict "page" .) }}`
	test(variant, expectedOutputDefaultFormat)

	variant = `{{ template "_internal/pagination.html" (dict "page" . "format" "default") }}`
	test(variant, expectedOutputDefaultFormat)

	variant = `{{ template "_internal/pagination.html" (dict "page" . "format" "terse") }}`
	test(variant, expectedOutputTerseFormat)
}
