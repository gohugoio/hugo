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

package resources_test

import (
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/hugolib"
)

// Issue 8931
func TestImageCache(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
baseURL = "https://example.org"
-- content/mybundle/index.md --
---
title: "My Bundle"
---
-- content/mybundle/pixel.png --
iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==
-- content/mybundle/giphy.gif --
sourcefilename: testdata/giphy.gif
-- layouts/foo.html --
-- layouts/index.html --
{{ $p := site.GetPage "mybundle"}}
{{ $img := $p.Resources.Get "pixel.png" }}
{{ $giphy := $p.Resources.Get "giphy.gif" }}
{{ $gif := $img.Resize "1x2 gif" }}
{{ $bmp := $img.Resize "2x3 bmp" }}
{{ $anigif := $giphy.Resize "4x5" }}


gif: {{ $gif.RelPermalink }}|}|{{ $gif.Width }}|{{ $gif.Height }}|{{ $gif.MediaType }}|
bmp: {{ $bmp.RelPermalink }}|}|{{ $bmp.Width }}|{{ $bmp.Height }}|{{ $bmp.MediaType }}|
anigif: {{ $anigif.RelPermalink }}|{{ $anigif.Width }}|{{ $anigif.Height }}|{{ $anigif.MediaType }}|
`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			NeedsOsFS:   true,
			Running:     true,
		}).Build()

	assertImages := func() {
		b.AssertFileContent("public/index.html", `
		gif: /mybundle/pixel_hu8aa3346827e49d756ff4e630147c42b5_70_1x2_resize_box_3.gif|}|1|2|image/gif|
		bmp: /mybundle/pixel_hu8aa3346827e49d756ff4e630147c42b5_70_2x3_resize_box_3.bmp|}|2|3|image/bmp|
		anigif: /mybundle/giphy_hu3eafc418e52414ace6236bf1d31f82e1_52213_4x5_resize_box_1.gif|4|5|image/gif|
		`)
	}

	assertImages()

	b.EditFileReplaceFunc("content/mybundle/index.md", func(s string) string { return strings.ReplaceAll(s, "Bundle", "BUNDLE") })
	b.Build()

	assertImages()
}

func TestSVGError(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
-- assets/circle.svg --
<svg height="100" width="100"><circle cx="50" cy="50" r="40" stroke="black" stroke-width="3" fill="red" /></svg>
-- layouts/index.html --
{{ $svg := resources.Get "circle.svg" }}
Width: {{ $svg.Width }}
`

	b, err := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			NeedsOsFS:   true,
			Running:     true,
		}).BuildE()

	b.Assert(err, qt.IsNotNil)
	b.Assert(err.Error(), qt.Contains, `error calling Width: this method is only available for raster images. To determine if an image is SVG, you can do {{ if eq .MediaType.SubType "svg" }}{{ end }}`)
}

// Issue 10255.
func TestNoPublishOfUnusedProcessedImage(t *testing.T) {
	t.Parallel()

	workingDir := t.TempDir()

	files := `
-- assets/images/pixel.png --
iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==
-- layouts/index.html --
{{ $image := resources.Get "images/pixel.png" }}
{{ $image = $image.Resize "400x" }}
{{ $image = $image.Resize "300x" }}
{{ $image = $image.Resize "200x" }}
{{ $image = $image.Resize "100x" }}
{{ $image = $image.Crop "50x50" }}
{{ $image = $image.Filter (images.GaussianBlur 6) }}
{{ ($image | fingerprint).Permalink }}


`

	for i := 0; i < 3; i++ {

		b := hugolib.NewIntegrationTestBuilder(
			hugolib.IntegrationTestConfig{
				T:           t,
				TxtarString: files,
				NeedsOsFS:   true,
				WorkingDir:  workingDir,
			}).Build()

		b.AssertFileCount("resources/_gen/images", 6)
		b.AssertFileCount("public/images", 1)
		b.Build()
	}
}

func TestProcessFilter(t *testing.T) {
	t.Parallel()

	files := `
-- assets/images/pixel.png --
iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==
-- layouts/index.html --
{{ $pixel := resources.Get "images/pixel.png" }}
{{ $filters := slice (images.GaussianBlur 6) (images.Pixelate 8) (images.Process "jpg") }}
{{ $image := $pixel.Filter $filters }}
jpg|RelPermalink: {{ $image.RelPermalink }}|MediaType: {{ $image.MediaType }}|Width: {{ $image.Width }}|Height: {{ $image.Height }}|
{{ $filters := slice (images.GaussianBlur 6) (images.Pixelate 8) (images.Process "jpg resize 20x30") }}
{{ $image := $pixel.Filter $filters }}
resize 1|RelPermalink: {{ $image.RelPermalink }}|MediaType: {{ $image.MediaType }}|Width: {{ $image.Width }}|Height: {{ $image.Height }}|
{{ $image := $pixel.Filter $filters }}
resize 2|RelPermalink: {{ $image.RelPermalink }}|MediaType: {{ $image.MediaType }}|Width: {{ $image.Width }}|Height: {{ $image.Height }}|

`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html",
		"jpg|RelPermalink: /images/pixel_hu8aa3346827e49d756ff4e630147c42b5_70_filter_14347832029652749901.jpg|MediaType: image/jpeg|Width: 1|Height: 1|",
		"resize 1|RelPermalink: /images/pixel_hu8aa3346827e49d756ff4e630147c42b5_70_filter_14038289290831899195.jpg|MediaType: image/jpeg|Width: 20|Height: 30|",
		"resize 2|RelPermalink: /images/pixel_hu8aa3346827e49d756ff4e630147c42b5_70_filter_14038289290831899195.jpg|MediaType: image/jpeg|Width: 20|Height: 30|",
	)
}

// Issue #11563
func TestGroupByParamDate(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
disableKinds = ['section','rss','sitemap','taxonomy','term']
-- layouts/index.html --
{{- range site.RegularPages.GroupByParamDate "eventDate" "2006-01" }}
	{{- .Key }}|{{ range .Pages }}{{ .Title }}|{{ end }}
{{- end }}
-- content/p1.md --
+++
title = 'p1'
eventDate = 2023-09-01
+++
-- content/p2.md --
+++
title = 'p2'
eventDate = '2023-09-01'
+++
-- content/p3.md --
---
title: p3
eventDate: 2023-09-01
---
-- content/p4.md --
+++
title = 'p4'
eventDate = 2023-10-01T08:00:00
+++
-- content/p5.md --
+++
title = 'p5'
eventDate = '2023-10-01T08:00:00'
+++
-- content/p6.md --
---
title: p6
eventDate: 2023-10-01T08:00:00
---
-- content/p7.md --
+++
title = 'p7'
eventDate = 2023-11-01T07:00:00-08:00
+++
-- content/p8.md --
+++
title = 'p8'
eventDate = '2023-11-01T07:00:00-08:00'
+++
-- content/p9.md --
---
title: p9
eventDate: 2023-11-01T07:00:00-08:00
---
  `

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html", "2023-11|p9|p8|p7|2023-10|p6|p5|p4|2023-09|p3|p2|p1|")
}

// Issue 10412
func TestImageTransformThenCopy(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['page','rss','section','sitemap','taxonomy','term']
-- assets/pixel.png --
iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==
-- layouts/index.html --
{{- with resources.Get "pixel.png" }}
  {{- with .Resize "200x" | resources.Copy "pixel.png" }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}">|{{ .Key }}
  {{- end }}
{{- end }}
`

	b := hugolib.Test(t, files)

	b.AssertFileExists("public/pixel.png", true)
	b.AssertFileContent("public/index.html",
		`<img src="/pixel.png" width="200" height="200">|/pixel.png`,
	)
}

// Issue 12310
func TestUseDifferentCacheKeyForResourceCopy(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['page','section','rss','sitemap','taxonomy','term']
-- assets/a.txt --
This was assets/a.txt
-- layouts/index.html --
{{ $nilResource := resources.Get "/p1/b.txt" }}
{{ $r := resources.Get "a.txt" }}
{{ $r = resources.Copy "/p1/b.txt" $r }}
{{ $r.RelPermalink }}
`

	b, err := hugolib.TestE(t, files)

	b.Assert(err, qt.IsNil)
	b.AssertFileContent("public/p1/b.txt", "This was assets/a.txt")
}
