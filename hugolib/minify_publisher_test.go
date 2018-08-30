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

	"github.com/spf13/viper"

	"github.com/stretchr/testify/require"
)

func TestMinifyPublisher(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	v := viper.New()
	v.Set("minify", true)
	v.Set("baseURL", "https://example.org/")

	htmlTemplate := `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<title>HTML5 boilerplate – all you really need…</title>
	<link rel="stylesheet" href="css/style.css">
	<!--[if IE]>
		<script src="http://html5shiv.googlecode.com/svn/trunk/html5.js"></script>
	<![endif]-->
</head>

<body id="home">

	<h1>{{ .Page.Title }}</h1>

</body>
</html>
`

	b := newTestSitesBuilder(t)
	b.WithViper(v).WithContent("page.md", pageWithAlias)
	b.WithTemplates("_default/list.html", htmlTemplate, "_default/single.html", htmlTemplate, "alias.html", htmlTemplate)
	b.CreateSites().Build(BuildCfg{})

	assert.Equal(1, len(b.H.Sites))
	require.Len(t, b.H.Sites[0].RegularPages, 1)

	// Check minification
	// HTML
	b.AssertFileContent("public/page/index.html", "<!doctype html><html lang=en><head><meta charset=utf-8><title>HTML5 boilerplate – all you really need…</title><link rel=stylesheet href=css/style.css></head><body id=home><h1>Has Alias</h1></body></html>")
	// HTML alias. Note the custom template which does no redirect.
	b.AssertFileContent("public/foo/bar/index.html", "<!doctype html><html lang=en><head><meta charset=utf-8><title>HTML5 boilerplate ")

	// RSS
	b.AssertFileContent("public/index.xml", "<?xml version=\"1.0\" encoding=\"utf-8\" standalone=\"yes\"?><rss version=\"2.0\" xmlns:atom=\"http://www.w3.org/2005/Atom\"><channel><title/><link>https://example.org/</link>")

	// Sitemap
	b.AssertFileContent("public/sitemap.xml", "<?xml version=\"1.0\" encoding=\"utf-8\" standalone=\"yes\"?><urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\" xmlns:xhtml=\"http://www.w3.org/1999/xhtml\"><url><loc>https://example.org/</loc><priority>0</priority></url><url>")
}
