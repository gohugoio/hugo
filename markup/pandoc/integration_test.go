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

package pandoc_test

import (
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

func TestBasicConversion(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
-- content/p1.md --
testContent
-- layouts/_default/single.html --
{{ .Content }}
`
	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			NeedsOsFS:   true,
		},
	).Build()

	b.AssertFileContent("public/p1/index.html", `<p>testContent</p>`)
}

func TestConversionWithHeader(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
-- content/p1.md --
# testContent
-- layouts/_default/single.html --
{{ .Content }}
`
	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			NeedsOsFS:   true,
		},
	).Build()

	b.AssertFileContent("public/p1/index.html", `<h1 id="testcontent">testContent</h1>`)
}

func TestConversionWithExtractedToc(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
-- content/p1.md --
# title 1
## title 2
-- layouts/_default/single.html --
{{ .TableOfContents }}
{{ .Content }}
`
	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			NeedsOsFS:   true,
		},
	).Build()

	b.AssertFileContent("public/p1/index.html", "<nav id=\"TableOfContents\">\n  <ul>\n    <li><a href=\"#title-2\">title 2</a></li>\n  </ul>\n</nav>\n<h1 id=\"title-1\">title 1</h1>\n<h2 id=\"title-2\">title 2</h2>")
}
