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

package warpc_test

import (
	"path/filepath"
	"runtime"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/hugolib"
)

func TestWebPMisc(t *testing.T) {
	files := `
-- assets/sunrise.webp --
sourcefilename: ../../resources/testdata/sunrise.webp
-- layouts/home.html --
{{ $image := resources.Get "sunrise.webp" }}
{{ $resized := $image.Resize "123x456" }}
Original Width/Height: {{ $image.Width }}/{{ $image.Height }}|
Resized Width:/Height {{ $resized.Width }}/{{ $resized.Height }}|
Resized RelPermalink: {{ $resized.RelPermalink }}|
{{ $ic := images.Config "/assets/sunrise.webp" }}
ImageConfig: {{ printf "%d/%d" $ic.Width $ic.Height }}|
`

	b := hugolib.Test(t, files)

	b.ImageHelper("public/sunrise_hu_8dd1706a77fb35ce.webp").AssertFormat("webp").AssertIsAnimated(false)

	b.AssertFileContent("public/index.html",
		"Original Width/Height: 1024/640|",
		"Resized Width:/Height 123/456|",
		"Resized RelPermalink: /sunrise_hu_8dd1706a77fb35ce.webp|",
		"ImageConfig: 1024/640|",
	)
}

func TestWebPEncodeGrayscale(t *testing.T) {
	files := `
-- assets/gopher.png --
sourcefilename: ../../resources/testdata/bw-gopher.png
-- layouts/home.html --
{{ $image := resources.Get "gopher.png" }}
{{ $resized := $image.Resize "123x456 webp" }}
Resized RelPermalink: {{ $resized.RelPermalink }}|
`

	b := hugolib.Test(t, files)

	b.ImageHelper("public/gopher_hu_cc98ebaf742cba8e.webp").AssertFormat("webp")
}

func TestWebPInvalid(t *testing.T) {
	files := `
-- assets/invalid.webp --
sourcefilename: ../../resources/testdata/webp/invalid.webp
-- layouts/home.html --
{{ $image := resources.Get "invalid.webp" }}
{{ $resized := $image.Resize "123x456 webp" }}
Resized RelPermalink: {{ $resized.RelPermalink }}|
`
	tempDir := t.TempDir()

	b, err := hugolib.TestE(t, files, hugolib.TestOptWithConfig(func(cfg *hugolib.IntegrationTestConfig) {
		cfg.NeedsOsFS = true
		cfg.WorkingDir = tempDir
	}))
	b.Assert(err, qt.IsNotNil)

	if runtime.GOOS != "windows" {
		// Make sure the full image filename is in the error message.
		filename := filepath.Join(tempDir, "assets/invalid.webp")
		b.Assert(err.Error(), qt.Contains, filename)
	}
}

// This test isn't great, but we have golden tests to verify the output itself.
func TestWebPAnimation(t *testing.T) {
	files := `
-- hugo.toml --
disableKinds = ["page", "section", "taxonomy", "term", "sitemap", "robotsTXT", "404"]
-- assets/anim.webp --
sourcefilename: ../../resources/testdata/webp/anim.webp
-- assets/giphy.gif --
sourcefilename: ../../resources/testdata/giphy.gif
-- layouts/home.html --
{{ $webpAnim := resources.Get "anim.webp" }}
{{ $gifAnim := resources.Get "giphy.gif" }}
{{ ($webpAnim.Resize "100x100 webp").Publish }}
{{ ($webpAnim.Resize "100x100 gif").Publish }}
{{ ($gifAnim.Resize "100x100 gif").Publish }}
{{ ($gifAnim.Resize "100x100 webp").Publish }}

`

	b := hugolib.Test(t, files)

	// Source animated gif:
	// Frame durations in ms.
	giphyFrameDurations := []int{200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200}

	b.ImageHelper("public/giphy_hu_bb052284cc220165.webp").AssertFormat("webp").AssertIsAnimated(true).AssertLoopCount(0).AssertFrameDurations(giphyFrameDurations)
	b.ImageHelper("public/giphy_hu_c6b8060edf0363b1.gif").AssertFormat("gif").AssertIsAnimated(true).AssertLoopCount(0).AssertFrameDurations(giphyFrameDurations)

	// Source animated webp:
	animFrameDurations := []int{80, 80, 80, 80, 80, 80, 80, 80, 80, 80, 80, 80, 80, 80, 80, 80, 80}
	b.ImageHelper("public/anim_hu_edc2f24aaad2cee6.webp").AssertFormat("webp").AssertIsAnimated(true).AssertLoopCount(0).AssertFrameDurations(animFrameDurations)
	b.ImageHelper("public/anim_hu_58eb49733894e7ce.gif").AssertFormat("gif").AssertIsAnimated(true).AssertLoopCount(0).AssertFrameDurations(animFrameDurations)
}

func BenchmarkWebp(b *testing.B) {
	files := `
-- content/p1/sunrise.webp --
sourcefilename: ../../resources/testdata/sunrise.webp
-- content/p1/index.md --
-- content/p2/sunrise.webp --
sourcefilename: ../../resources/testdata/sunrise.webp
-- content/p2/index.md --
-- content/p3/sunrise.webp --
sourcefilename: ../../resources/testdata/sunrise.webp
-- content/p3/index.md --
-- content/p4/sunrise.webp --
sourcefilename: ../../resources/testdata/sunrise.webp
-- content/p4/index.md --
-- content/p5/sunrise.webp --
sourcefilename: ../../resources/testdata/sunrise.webp
-- content/p5/index.md --
-- content/p6/sunrise.webp --
sourcefilename: ../../resources/testdata/sunrise.webp
-- content/p6/index.md --
-- content/p7/sunrise.webp --
sourcefilename: ../../resources/testdata/sunrise.webp
-- content/p7/index.md --
-- content/p8/sunrise.webp --
sourcefilename: ../../resources/testdata/sunrise.webp
-- content/p8/index.md --
-- content/p9/sunrise.webp --
sourcefilename: ../../resources/testdata/sunrise.webp
-- content/p9/index.md --
-- content/p10/sunrise.webp --
sourcefilename: ../../resources/testdata/sunrise.webp
-- content/p10/index.md --
-- layouts/home.html --
Home.
-- layouts/page.html --
{{ $image := .Resources.Get "sunrise.webp" }}
{{ $resized := $image.Fit "400x300 webp" }}
Resized RelPermalink: {{ $resized.RelPermalink }}|
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
