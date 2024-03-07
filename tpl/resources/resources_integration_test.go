// Copyright 2022s The Hugo Authors. All rights reserved.
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

package resources_test

import (
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/hugolib"
)

func TestCopy(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
baseURL = "http://example.com/blog"
-- assets/images/pixel.png --
iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==
-- layouts/index.html --
{{/* Image resources */}}
{{ $img := resources.Get "images/pixel.png" }}
{{ $imgCopy1 := $img | resources.Copy "images/copy.png"  }}
{{ $imgCopy1 = $imgCopy1.Resize "3x4"}}
{{ $imgCopy2 := $imgCopy1 | resources.Copy "images/copy2.png" }}
{{ $imgCopy3 := $imgCopy1 | resources.Copy "images/copy3.png" }}
Image Orig:  {{ $img.RelPermalink}}|{{ $img.MediaType }}|{{ $img.Width }}|{{ $img.Height }}|
Image Copy1:  {{ $imgCopy1.RelPermalink}}|{{ $imgCopy1.MediaType }}|{{ $imgCopy1.Width }}|{{ $imgCopy1.Height }}|
Image Copy2:  {{ $imgCopy2.RelPermalink}}|{{ $imgCopy2.MediaType }}|{{ $imgCopy2.Width }}|{{ $imgCopy2.Height }}|
Image Copy3:  {{ $imgCopy3.MediaType }}|{{ $imgCopy3.Width }}|{{ $imgCopy3.Height }}|

{{/* Generic resources */}}
{{ $targetPath := "js/vars.js" }}
{{ $orig := "let foo;" | resources.FromString "js/foo.js" }}
{{ $copy1 := $orig | resources.Copy "js/copies/bar.js" }}
{{ $copy2 := $orig | resources.Copy "js/copies/baz.js" | fingerprint "md5" }}
{{ $copy3 := $copy2 | resources.Copy "js/copies/moo.js" | minify }}

Orig: {{ $orig.RelPermalink}}|{{ $orig.MediaType }}|{{ $orig.Content | safeJS }}|
Copy1: {{ $copy1.RelPermalink}}|{{ $copy1.MediaType }}|{{ $copy1.Content | safeJS }}|
Copy2: {{ $copy2.RelPermalink}}|{{ $copy2.MediaType }}|{{ $copy2.Content | safeJS }}|
Copy3: {{ $copy3.RelPermalink}}|{{ $copy3.MediaType }}|{{ $copy3.Content | safeJS }}|

	`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html", `
Image Orig:  /blog/images/pixel.png|image/png|1|1|
Image Copy1:  /blog/images/copy_hu8aa3346827e49d756ff4e630147c42b5_70_3x4_resize_box_3.png|image/png|3|4|
Image Copy2:  /blog/images/copy2.png|image/png|3|4
Image Copy3:  image/png|3|4|
Orig: /blog/js/foo.js|text/javascript|let foo;|
Copy1: /blog/js/copies/bar.js|text/javascript|let foo;|
Copy2: /blog/js/copies/baz.a677329fc6c4ad947e0c7116d91f37a2.js|text/javascript|let foo;|
Copy3: /blog/js/copies/moo.a677329fc6c4ad947e0c7116d91f37a2.min.js|text/javascript|let foo|

		`)

	b.AssertFileExists("public/images/copy2.png", true)
	// No permalink used.
	b.AssertFileExists("public/images/copy3.png", false)
}

func TestCopyPageShouldFail(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
-- layouts/index.html --
{{/* This is currently not supported. */}}
{{ $copy := .Copy "copy.md" }}

	`

	b, err := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
		}).BuildE()

	b.Assert(err, qt.IsNotNil)
}

func TestGet(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
baseURL = "http://example.com/blog"
-- assets/images/pixel.png --
iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==
-- layouts/index.html --
{{ with resources.Get "images/pixel.png" }}Image OK{{ else }}Image not found{{ end }}
{{ with resources.Get "" }}Failed{{ else }}Empty string not found{{ end }}


	`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html", `
Image OK
Empty string not found

		`)
}

func TestResourcesGettersShouldNotNormalizePermalinks(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
baseURL = "http://example.com/"
-- assets/401K Prospectus.txt --
Prospectus.
-- layouts/index.html --
{{ $name := "401K Prospectus.txt" }}
Get: {{ with resources.Get $name }}{{ .RelPermalink }}|{{ .Permalink }}|{{ end }}
GetMatch: {{ with resources.GetMatch $name }}{{ .RelPermalink }}|{{ .Permalink }}|{{ end }}
Match: {{ with (index (resources.Match $name) 0) }}{{ .RelPermalink }}|{{ .Permalink }}|{{ end }}
ByType: {{ with (index (resources.ByType "text") 0) }}{{ .RelPermalink }}|{{ .Permalink }}|{{ end }}
	`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html", `
Get: /401K%20Prospectus.txt|http://example.com/401K%20Prospectus.txt|
GetMatch: /401K%20Prospectus.txt|http://example.com/401K%20Prospectus.txt|
Match: /401K%20Prospectus.txt|http://example.com/401K%20Prospectus.txt|
ByType: /401K%20Prospectus.txt|http://example.com/401K%20Prospectus.txt|

		`)
}

func TestGlobalResourcesNotPublishedRegressionIssue12190(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['page','rss','section','sitemap','taxonomy','term']
-- assets/a.txt --
I am a.txt
-- assets/b.txt --
I am b.txt
-- layouts/index.html --
Home.
{{ with resources.ByType "text" }}
  {{ with .Get "a.txt" }}
    {{ .Publish }}
  {{ end }}
  {{ with .GetMatch "*b*" }}
    {{ .Publish }}
  {{ end }}
{{ end }}
`

	b := hugolib.Test(t, files)

	b.AssertFileExists("public/a.txt", true) // failing test
	b.AssertFileExists("public/b.txt", true) // failing test
}
