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
	"fmt"
	"testing"

	"github.com/bep/logg"
	"github.com/gohugoio/hugo/htesting"
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
	// Go changed their JPEG implementation in Go 1.26, so we cannot run this test on earlier versions. See https://go.dev/doc/go1.26#imagejpegpkgimagejpeg
	if htesting.GoMinorVersion() < 26 {
		t.Skip("Go 1.26+ required")
	}
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
		b.AssertFileContent("public/index.html", "Colors: [#2e2f33 #a69e94 #d29d59 #a26a3f #747c83 #7b848b]|")

	}
}

func TestImageCropSmartKeepsTargetSizeIssue13688(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['page','rss','section','sitemap','taxonomy','term']
-- assets/sunset.jpg --
sourcefilename: ../testdata/sunset.jpg
-- layouts/home.html --
{{ with resources.Get "sunset.jpg" }}
Original: {{ .Width }}x{{ .Height }}|
{{- with .Crop "900x561 TopLeft" -}}
CropTopLeft: {{ .Width }}x{{ .Height }}|
{{- end -}}
{{- with .Crop "900x561 Smart" -}}
CropSmart: {{ .Width }}x{{ .Height }}|
{{- end -}}
{{ end }}
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html", "Original: 900x562|CropTopLeft: 900x561|CropSmart: 900x561|")
}

// See issue 11266.
func TestImageSmartCropWithRotation(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['page','rss','section','sitemap','taxonomy','term']
-- assets/sunset.jpg --
sourcefilename: ../testdata/sunset.jpg
-- layouts/home.html --
{{ $landscape := (resources.Get "sunset.jpg").Process "png" }}
{{ $portrait := $landscape.Process "r90 png" }}
{{ range $name, $img := dict "landscape" $landscape "portrait" $portrait }}
  {{ range $angle := slice 0 90 180 270 45 }}
    {{ $rotated := $img.Process (printf "r%d png" $angle) }}
    {{ $size := "200x150 smart png" }}
    {{ if eq $angle 45 }}{{ $size = printf "%s nearestneighbor" $size }}{{ end }}
    {{ range $action := slice "crop" "fill" }}
      {{ $spec := printf "r%d %s" $angle $size }}
      {{ $got := $img.Crop $spec }}
      {{ $want := $rotated.Crop $size }}
      {{ if eq $action "fill" }}
        {{ $got = $img.Fill $spec }}
        {{ $want = $rotated.Fill $size }}
      {{ end }}
      {{ $filtered := $img.Filter (images.Process (printf "%s %s" $action $spec)) }}
      {{ $name }} {{ $action }} {{ $angle }}: {{ eq (sha256 $got.Content) (sha256 $want.Content) }}|
      {{ $name }} {{ $action }} {{ $angle }} filter: {{ eq (sha256 $filtered.Content) (sha256 $want.Content) }}|
    {{ end }}
  {{ end }}
{{ end }}
`

	b := hugolib.Test(t, files)
	for _, name := range []string{"landscape", "portrait"} {
		for _, angle := range []int{0, 90, 180, 270, 45} {
			for _, action := range []string{"crop", "fill"} {
				b.AssertFileContent("public/index.html",
					fmt.Sprintf("%s %s %d: true|", name, action, angle),
					fmt.Sprintf("%s %s %d filter: true|", name, action, angle),
				)
			}
		}
	}
}

func TestImagingGlobalsDeprecated(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
[imaging]
quality = 70
hint = "picture"
compression = "lossless"
-- layouts/home.html --
Home.
`

	b := hugolib.Test(t, files, hugolib.TestOptWithConfig(func(cfg *hugolib.IntegrationTestConfig) {
		cfg.LogLevel = logg.LevelInfo
	}))

	b.AssertLogContains(
		"project config key imaging.quality was deprecated in Hugo v0.163.0",
		"project config key imaging.hint was deprecated in Hugo v0.163.0",
		"project config key imaging.compression was deprecated in Hugo v0.163.0",
	)
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
