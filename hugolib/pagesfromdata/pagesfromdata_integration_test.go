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

package pagesfromdata_test

import (
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

// TODO1 add a section file in same file as the _content.json.
var pagesFromDataBasicJSON = `
-- hugo.toml --
disableKinds = ["home", "taxonomy", "term", "rss", "sitemap"]
baseURL = "https://example.com"
disableLiveReload = true
-- layouts/_default/single.html --
Single: {{ .Title }}|{{ .Content }}|
-- layouts/_default/list.html --
List: {{ .Title }}|{{ .Content }}|
-- content/docs/_index.md --
---
title: "My Docs"
---
-- content/docs/_content.jsonl --
{
    "path": "my-section",
    "title": "My Section",
    "kind": "section",
    "date": "2018-12-25",
    "lastMod": "2018-12-25",
    "resources": [
        {
            "path": "sunset.jpg",
            "mediaType": "image/jpg",
            "url": "file:///Users/bep/Downloads/IMG_20181225_123456.jpg"
        },
        {
            "path": "bootstrap.min.css",
            "mediaType": "txt/css",
            "url": "https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css"
        }
    ]
}
{
    "path": "my-section/my-page",
    "kind": "page",
    "title": "My Page",
    "date": "2018-12-25",
    "lastMod": "2018-12-25",
    "content": {
        "type": "text",
        "value": "My **Page** Content"
    }
}


`

// TODO1 yaml support, see https://github.com/go-yaml/yaml/pull/301
func TestPagesFromDataBasicJSON(t *testing.T) {
	t.Parallel()

	b := hugolib.Test(t, pagesFromDataBasicJSON)

	b.AssertFileContent("public/docs/my-section/my-page/index.html", "Single: My Page|<p>My <strong>Page</strong> Content</p>\n|")
	b.AssertFileContent("public/docs/my-section/index.html", "List: My Section||")
}

func TestPagesFromDataBasicYAML(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ["home", "taxonomy", "term", "rss", "sitemap"]
baseURL = "https://example.com"
disableLiveReload = true
-- layouts/_default/single.html --
Single: {{ .Title }}|{{ .Content }}|
-- layouts/_default/list.html --
List: {{ .Title }}|{{ .Content }}|
-- content/docs/_index.md --
---
title: "My Docs"
---
-- content/docs/_content.yaml --
path: "my-section"
title: "My Section"
kind: "section"
date: "2018-12-25"
---
path: "my-section/my-page"
kind: "page"
title: "My Page"
    `

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/docs/my-section/my-page/index.html", "Single: My Page||")
	b.AssertFileContent("public/docs/my-section/index.html", "List: My Section||")
}

func TestPagesFromDataRebuildEditPage(t *testing.T) {
	t.Parallel()

	b := hugolib.TestRunning(t, pagesFromDataBasicJSON)

	b.AssertPublishDir("docs/my-section/my-page/index.html")
	b.AssertFileContent("public/docs/my-section/my-page/index.html", "Single: My Page")
	b.AssertRenderCountPage(3)
	b.EditFileReplaceAll("content/docs/_content.jsonl", "My Page", "My Page edited").Build()
	b.AssertFileContent("public/docs/my-section/my-page/index.html", "Single: My Page edited")
	b.AssertRenderCountPage(1)
}
