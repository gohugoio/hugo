// Copyright 2018 The Hugo Authors. All rights reserved.
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

	qt "github.com/frankban/quicktest"
)

// Just some simple test of the embedded templates to avoid
// https://github.com/gohugoio/hugo/issues/4757 and similar.
// TODO(bep) fix me https://github.com/gohugoio/hugo/issues/5926
func _TestEmbeddedTemplates(t *testing.T) {
	t.Parallel()

	c := qt.New(t)
	c.Assert(true, qt.Equals, true)

	home := []string{"index.html", `
GA:
{{ template "_internal/google_analytics.html" . }}

GA async:

{{ template "_internal/google_analytics_async.html" . }}

Disqus:

{{ template "_internal/disqus.html" . }}

`}

	b := newTestSitesBuilder(t)
	b.WithSimpleConfigFile().WithTemplatesAdded(home...)

	b.Build(BuildCfg{})

	// Gheck GA regular and async
	b.AssertFileContent("public/index.html",
		"'anonymizeIp', true",
		"'script','https://www.google-analytics.com/analytics.js','ga');\n\tga('create', 'ga_id', 'auto')",
		"<script async src='https://www.google-analytics.com/analytics.js'>")

	// Disqus
	b.AssertFileContent("public/index.html", "\"disqus_shortname\" + '.disqus.com/embed.js';")
}

func TestInternalTemplatesImage(t *testing.T) {
	config := `
baseURL = "https://example.org"

[params]
images=["siteimg1.jpg", "siteimg2.jpg"]

`
	b := newTestSitesBuilder(t).WithConfigFile("toml", config)

	b.WithContent("mybundle/index.md", `---
title: My Bundle
---
`)

	b.WithContent("mypage.md", `---
title: My Page
images: ["pageimg1.jpg", "pageimg2.jpg"]
---
`)

	b.WithContent("mysite.md", `---
title: My Site
---
`)

	b.WithTemplatesAdded("_default/single.html", `

{{ template "_internal/twitter_cards.html" . }}
{{ template "_internal/opengraph.html" . }}
{{ template "_internal/schema.html" . }}

`)

	b.WithSunset("content/mybundle/featured-sunset.jpg")
	b.Build(BuildCfg{})

	b.AssertFileContent("public/mybundle/index.html", `
<meta name="twitter:image" content="https://example.org/mybundle/featured-sunset.jpg"/>
<meta name="twitter:title" content="My Bundle"/>
<meta property="og:title" content="My Bundle" />
<meta property="og:url" content="https://example.org/mybundle/" />
<meta property="og:image" content="https://example.org/mybundle/featured-sunset.jpg"/>
<meta itemprop="name" content="My Bundle">
<meta itemprop="image" content="https://example.org/mybundle/featured-sunset.jpg">

`)
	b.AssertFileContent("public/mypage/index.html", `
<meta name="twitter:image" content="https://example.org/pageimg1.jpg"/>
<meta property="og:image" content="https://example.org/pageimg1.jpg" />
<meta property="og:image" content="https://example.org/pageimg2.jpg" />
<meta itemprop="image" content="https://example.org/pageimg1.jpg">
<meta itemprop="image" content="https://example.org/pageimg2.jpg">        
`)
	b.AssertFileContent("public/mysite/index.html", `
<meta name="twitter:image" content="https://example.org/siteimg1.jpg"/>
<meta property="og:image" content="https://example.org/siteimg1.jpg"/>
<meta itemprop="image" content="https://example.org/siteimg1.jpg"/>
`)

}
