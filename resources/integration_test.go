// Copyright 2022 The Hugo Authors. All rights reserved.
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
-- layouts/foo.html --
-- layouts/index.html --
{{ $p := site.GetPage "mybundle"}}
{{ $img := $p.Resources.Get "pixel.png" }}
{{ $gif := $img.Resize "1x1 gif" }}
{{ $bmp := $img.Resize "1x1 bmp" }}

gif: {{ $gif.RelPermalink }}|{{ $gif.MediaType }}|
bmp: {{ $bmp.RelPermalink }}|{{ $bmp.MediaType }}|	
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
		gif: /mybundle/pixel_hu8aa3346827e49d756ff4e630147c42b5_70_1x1_resize_box_3.gif|image/gif|
		bmp: /mybundle/pixel_hu8aa3346827e49d756ff4e630147c42b5_70_1x1_resize_box_3.bmp|image/bmp|
		
		`)
	}

	assertImages()

	b.EditFileReplace("content/mybundle/index.md", func(s string) string { return strings.ReplaceAll(s, "Bundle", "BUNDLE") })
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
