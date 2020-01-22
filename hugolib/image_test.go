// Copyright 2019 The Hugo Authors. All rights reserved.
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
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/gohugoio/hugo/htesting"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/spf13/viper"
)

// We have many tests for the different resize operations etc. in the resource package,
// this is an integration test.
func TestImageOps(t *testing.T) {
	c := qt.New(t)
	// Make this a real as possible.
	workDir, clean, err := htesting.CreateTempDir(hugofs.Os, "image-resize")
	c.Assert(err, qt.IsNil)
	defer clean()

	newBuilder := func(timeout interface{}) *sitesBuilder {

		v := viper.New()
		v.Set("workingDir", workDir)
		v.Set("baseURL", "https://example.org")
		v.Set("timeout", timeout)

		b := newTestSitesBuilder(t).WithWorkingDir(workDir)
		b.Fs = hugofs.NewDefault(v)
		b.WithViper(v)
		b.WithContent("mybundle/index.md", `
---
title: "My bundle"
---

{{< imgproc >}}

`)

		b.WithTemplatesAdded(
			"shortcodes/imgproc.html", `
{{ $img := resources.Get "images/sunset.jpg" }}
{{ $r := $img.Resize "129x239" }}
IMG SHORTCODE: {{ $r.RelPermalink }}/{{ $r.Width }}
`,
			"index.html", `
{{ $p := .Site.GetPage "mybundle" }}
{{ $img1 := resources.Get "images/sunset.jpg" }}
{{ $img2 := $p.Resources.GetMatch "sunset.jpg" }}
{{ $img3 := resources.GetMatch "images/*.jpg" }}
{{ $r := $img1.Resize "123x234" }}
{{ $r2 := $r.Resize "12x23" }}
{{ $b := $img2.Resize "345x678" }}
{{ $b2 := $b.Resize "34x67" }}
{{ $c := $img3.Resize "456x789" }}
{{ $fingerprinted := $img1.Resize "350x" | fingerprint }}

{{ $images := slice $r $r2 $b $b2 $c $fingerprinted }}

{{ range $i, $r := $images }}
{{ printf "Resized%d:" (add $i  1) }} {{ $r.Name }}|{{ $r.Width }}|{{ $r.Height }}|{{ $r.MediaType }}|{{ $r.RelPermalink }}|
{{ end }}

{{ $blurryGrayscale1 := $r | images.Filter images.Grayscale (images.GaussianBlur 8) }}
BG1: {{ $blurryGrayscale1.RelPermalink }}/{{ $blurryGrayscale1.Width }}
{{ $blurryGrayscale2 := $r.Filter images.Grayscale (images.GaussianBlur 8) }}
BG2: {{ $blurryGrayscale2.RelPermalink }}/{{ $blurryGrayscale2.Width }}
{{ $blurryGrayscale2_2 := $r.Filter images.Grayscale (images.GaussianBlur 8) }}
BG2_2: {{ $blurryGrayscale2_2.RelPermalink }}/{{ $blurryGrayscale2_2.Width }}

{{ $filters := slice images.Grayscale (images.GaussianBlur 9) }}
{{ $blurryGrayscale3 := $r | images.Filter $filters }}
BG3: {{ $blurryGrayscale3.RelPermalink }}/{{ $blurryGrayscale3.Width }}

{{ $blurryGrayscale4 := $r.Filter $filters }}
BG4: {{ $blurryGrayscale4.RelPermalink }}/{{ $blurryGrayscale4.Width }}

{{ $p.Content }}

`)

		return b
	}

	imageDir := filepath.Join(workDir, "assets", "images")
	bundleDir := filepath.Join(workDir, "content", "mybundle")

	c.Assert(os.MkdirAll(imageDir, 0777), qt.IsNil)
	c.Assert(os.MkdirAll(bundleDir, 0777), qt.IsNil)
	src, err := os.Open("testdata/sunset.jpg")
	c.Assert(err, qt.IsNil)
	out, err := os.Create(filepath.Join(imageDir, "sunset.jpg"))
	c.Assert(err, qt.IsNil)
	_, err = io.Copy(out, src)
	c.Assert(err, qt.IsNil)
	out.Close()

	src.Seek(0, 0)

	out, err = os.Create(filepath.Join(bundleDir, "sunset.jpg"))
	c.Assert(err, qt.IsNil)
	_, err = io.Copy(out, src)
	c.Assert(err, qt.IsNil)
	out.Close()
	src.Close()

	// First build it with a very short timeout to trigger errors.
	b := newBuilder("10ns")

	imgExpect := `
Resized1: images/sunset.jpg|123|234|image/jpeg|/images/sunset_hu59e56ffff1bc1d8d122b1403d34e039f_90587_123x234_resize_q75_box.jpg|
Resized2: images/sunset.jpg|12|23|image/jpeg|/images/sunset_hu59e56ffff1bc1d8d122b1403d34e039f_90587_ada4bb1a57f77a63306e3bd67286248e.jpg|
Resized3: sunset.jpg|345|678|image/jpeg|/mybundle/sunset_hu59e56ffff1bc1d8d122b1403d34e039f_90587_345x678_resize_q75_box.jpg|
Resized4: sunset.jpg|34|67|image/jpeg|/mybundle/sunset_hu59e56ffff1bc1d8d122b1403d34e039f_90587_44d8c928664d7c5a67377c6ec58425ce.jpg|
Resized5: images/sunset.jpg|456|789|image/jpeg|/images/sunset_hu59e56ffff1bc1d8d122b1403d34e039f_90587_456x789_resize_q75_box.jpg|
Resized6: images/sunset.jpg|350|219|image/jpeg|/images/sunset_hu59e56ffff1bc1d8d122b1403d34e039f_90587_350x0_resize_q75_box.a86fe88d894e5db613f6aa8a80538fefc25b20fa24ba0d782c057adcef616f56.jpg|
BG1: /images/sunset_hu59e56ffff1bc1d8d122b1403d34e039f_90587_2ae8bb993431ec1aec40fe59927b46b4.jpg/123
BG2: /images/sunset_hu59e56ffff1bc1d8d122b1403d34e039f_90587_2ae8bb993431ec1aec40fe59927b46b4.jpg/123
BG3: /images/sunset_hu59e56ffff1bc1d8d122b1403d34e039f_90587_ed7740a90b82802261c2fbdb98bc8082.jpg/123
BG4: /images/sunset_hu59e56ffff1bc1d8d122b1403d34e039f_90587_ed7740a90b82802261c2fbdb98bc8082.jpg/123
IMG SHORTCODE: /images/sunset_hu59e56ffff1bc1d8d122b1403d34e039f_90587_129x239_resize_q75_box.jpg/129
`

	assertImages := func() {
		b.Helper()
		b.AssertFileContent(filepath.Join(workDir, "public/index.html"), imgExpect)
		b.AssertImage(350, 219, "public/images/sunset_hu59e56ffff1bc1d8d122b1403d34e039f_90587_350x0_resize_q75_box.a86fe88d894e5db613f6aa8a80538fefc25b20fa24ba0d782c057adcef616f56.jpg")
		b.AssertImage(129, 239, "public/images/sunset_hu59e56ffff1bc1d8d122b1403d34e039f_90587_129x239_resize_q75_box.jpg")
	}

	err = b.BuildE(BuildCfg{})
	if runtime.GOOS != "windows" && !strings.Contains(runtime.GOARCH, "arm") {
		// TODO(bep)
		c.Assert(err, qt.Not(qt.IsNil))
	}

	b = newBuilder(29000)
	b.Build(BuildCfg{})

	assertImages()

	// Truncate one image.
	imgInCache := filepath.Join(workDir, "resources/_gen/images/sunset_hu59e56ffff1bc1d8d122b1403d34e039f_90587_ed7740a90b82802261c2fbdb98bc8082.jpg")
	f, err := os.Create(imgInCache)
	c.Assert(err, qt.IsNil)
	f.Close()

	// Build it again to make sure we read images from file cache.
	b = newBuilder("30s")
	b.Build(BuildCfg{})

	assertImages()

}

func TestImageResizeMultilingual(t *testing.T) {

	b := newTestSitesBuilder(t).WithConfigFile("toml", `
baseURL="https://example.org"
defaultContentLanguage = "en"

[languages]
[languages.en]
title = "Title in English"
languageName = "English"
weight = 1
[languages.nn]
languageName = "Nynorsk"
weight = 2
title = "Tittel på nynorsk"
[languages.nb]
languageName = "Bokmål"
weight = 3
title = "Tittel på bokmål"
[languages.fr]
languageName = "French"
weight = 4
title = "French Title"

`)

	pageContent := `---
title: "Page"
---
`

	b.WithContent("bundle/index.md", pageContent)
	b.WithContent("bundle/index.nn.md", pageContent)
	b.WithContent("bundle/index.fr.md", pageContent)
	b.WithSunset("content/bundle/sunset.jpg")
	b.WithSunset("assets/images/sunset.jpg")
	b.WithTemplates("index.html", `
{{ with (.Site.GetPage "bundle" ) }}
{{ $sunset := .Resources.GetMatch "sunset*" }}
{{ if $sunset }}
{{ $resized := $sunset.Resize "200x200" }}
SUNSET FOR: {{ $.Site.Language.Lang }}: {{ $resized.RelPermalink }}/{{ $resized.Width }}/Lat: {{ $resized.Exif.Lat }}
{{ end }}
{{ else }}
No bundle for {{ $.Site.Language.Lang }}
{{ end }}

{{ $sunset2 := resources.Get "images/sunset.jpg" }}
{{ $resized2 := $sunset2.Resize "123x234" }}
SUNSET2: {{ $resized2.RelPermalink }}/{{ $resized2.Width }}/Lat: {{ $resized2.Exif.Lat }}


`)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/index.html", "SUNSET FOR: en: /bundle/sunset_hu59e56ffff1bc1d8d122b1403d34e039f_90587_200x200_resize_q75_box.jpg/200/Lat: 36.59744166666667")
	b.AssertFileContent("public/fr/index.html", "SUNSET FOR: fr: /fr/bundle/sunset_hu59e56ffff1bc1d8d122b1403d34e039f_90587_200x200_resize_q75_box.jpg/200/Lat: 36.59744166666667")
	b.AssertFileContent("public/index.html", " SUNSET2: /images/sunset_hu59e56ffff1bc1d8d122b1403d34e039f_90587_123x234_resize_q75_box.jpg/123/Lat: 36.59744166666667")
	b.AssertFileContent("public/nn/index.html", " SUNSET2: /images/sunset_hu59e56ffff1bc1d8d122b1403d34e039f_90587_123x234_resize_q75_box.jpg/123/Lat: 36.59744166666667")

	b.AssertImage(200, 200, "public/fr/bundle/sunset_hu59e56ffff1bc1d8d122b1403d34e039f_90587_200x200_resize_q75_box.jpg")
	b.AssertImage(200, 200, "public/bundle/sunset_hu59e56ffff1bc1d8d122b1403d34e039f_90587_200x200_resize_q75_box.jpg")

	// Check the file cache
	b.AssertImage(200, 200, "resources/_gen/images/bundle/sunset_hu59e56ffff1bc1d8d122b1403d34e039f_90587_200x200_resize_q75_box.jpg")

	b.AssertFileContent("resources/_gen/images/bundle/sunset_7645215769587362592.json",
		"DateTimeDigitized|time.Time", "PENTAX")
	b.AssertImage(123, 234, "resources/_gen/images/sunset_hu59e56ffff1bc1d8d122b1403d34e039f_90587_123x234_resize_q75_box.jpg")
	b.AssertFileContent("resources/_gen/images/sunset_7645215769587362592.json",
		"DateTimeDigitized|time.Time", "PENTAX")

	// TODO(bep) add this as a default assertion after Build()?
	b.AssertNoDuplicateWrites()

}
