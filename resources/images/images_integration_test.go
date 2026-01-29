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

package images_test

import (
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

func TestAutoOrient(t *testing.T) {
	files := `
-- hugo.toml --
-- assets/rotate270.jpg --
sourcefilename: ../testdata/exif/orientation6.jpg
-- layouts/home.html --
{{ $img := resources.Get "rotate270.jpg" }}
W/H original: {{ $img.Width }}/{{ $img.Height }}
{{ $rotated := $img.Filter images.AutoOrient }}
W/H rotated: {{ $rotated.Width }}/{{ $rotated.Height }}
`

	b := hugolib.Test(t, files)
	b.AssertFileContent("public/index.html", "W/H original: 80/40\n\nW/H rotated: 40/80")
}

// Issue 12733.
func TestOrientationEq(t *testing.T) {
	files := `
-- hugo.toml --
-- assets/rotate270.jpg --
sourcefilename: ../testdata/exif/orientation6.jpg
-- layouts/home.html --
{{ $img := resources.Get "rotate270.jpg" }}
{{ $orientation := $img.Exif.Tags.Orientation }}
Orientation: {{ $orientation }}|eq 6: {{ eq $orientation 6 }}|Type: {{ printf "%T" $orientation }}|
`

	b := hugolib.Test(t, files)
	b.AssertFileContent("public/index.html", "Orientation: 6|eq 6: true|")
}

func TestColorsIssue14453(t *testing.T) {
	files := `
-- hugo.toml --
-- assets/sunset.jpg --
sourcefilename: ../testdata/sunset.jpg
-- layouts/home.html --
{{ $img := resources.Get "sunset.jpg" }}
{{ $img := $img.Fit "100x100" }}
{{ $img := $img.Filter (slice images.AutoOrient (images.Process "fit 100x100 webp")) -}}
{{ $colors := $img.Colors }}
Colors: {{ $colors }}|
`
	tempDir := t.TempDir()
	for range 2 {
		b := hugolib.Test(t, files, hugolib.TestOptWithConfig(func(cfg *hugolib.IntegrationTestConfig) {
			cfg.NeedsOsFS = true
			cfg.WorkingDir = tempDir
		}))
		b.AssertFileContent("public/index.html", "Colors: [#2e2f34 #a39e94 #d39e57 #a96b3a #747b84 #7c838a]|")

	}
}

func BenchmarkImageResize(b *testing.B) {
	files := `
-- content/p1/sunrise.jpg --
sourcefilename: ../../resources/testdata/sunrise.jpg
-- content/p1/index.md --
-- content/p2/sunrise.jpg --
sourcefilename: ../../resources/testdata/sunrise.jpg
-- content/p2/index.md --
-- content/p3/sunrise.jpg --
sourcefilename: ../../resources/testdata/sunrise.jpg
-- content/p3/index.md --
-- content/p4/sunrise.jpg --
sourcefilename: ../../resources/testdata/sunrise.jpg
-- content/p4/index.md --
-- content/p5/sunrise.jpg --
sourcefilename: ../../resources/testdata/sunrise.jpg
-- content/p5/index.md --
-- content/p6/sunrise.jpg --
sourcefilename: ../../resources/testdata/sunrise.jpg
-- content/p6/index.md --
-- content/p7/sunrise.jpg --
sourcefilename: ../../resources/testdata/sunrise.jpg
-- content/p7/index.md --
-- content/p8/sunrise.jpg --
sourcefilename: ../../resources/testdata/sunrise.jpg
-- content/p8/index.md --
-- content/p9/sunrise.jpg --
sourcefilename: ../../resources/testdata/sunrise.jpg
-- content/p9/index.md --
-- content/p10/sunrise.jpg --
sourcefilename: ../../resources/testdata/sunrise.jpg
-- content/p10/index.md --
-- layouts/home.html --
Home.
-- layouts/page.html --
Page.
{{ $image := .Resources.Get "sunrise.jpg" }}
{{ ($image.Process "resize 200x200").Publish }}

`

	cfg := hugolib.IntegrationTestConfig{
		T:           b,
		TxtarString: files,
	}

	for b.Loop() {
		b.StopTimer()
		builder := hugolib.NewIntegrationTestBuilder(cfg).Init()
		b.StartTimer()
		builder.Build()
	}
}
