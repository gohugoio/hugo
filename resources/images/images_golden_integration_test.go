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
	"image"
	"image/gif"
	_ "image/jpeg"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/disintegration/gift"
	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/common/hashing"
	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/htesting"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/hugolib"
	"github.com/google/go-cmp/cmp"
)

var eq = qt.CmpEquals(
	cmp.Comparer(func(p1, p2 os.FileInfo) bool {
		return p1.Name() == p2.Name() && p1.Size() == p2.Size() && p1.IsDir() == p2.IsDir()
	}),
	cmp.Comparer(func(d1, d2 fs.DirEntry) bool {
		p1, err1 := d1.Info()
		p2, err2 := d2.Info()
		if err1 != nil || err2 != nil {
			return false
		}
		return p1.Name() == p2.Name() && p1.Size() == p2.Size() && p1.IsDir() == p2.IsDir()
	}),
)

var goldenOpts = struct {
	// Toggle this to write golden files to disk.
	// Note: Remember to set this to false before committing.
	writeGoldenFiles bool

	// This will skip any assertions. Useful when adding new golden variants to a test.
	devMode bool
}{
	writeGoldenFiles: false,
	devMode:          false,
}

// Note, if you're enabling writeGoldenFiles on a MacOS ARM 64 you need to run the test with GOARCH=amd64, e.g.
// GOARCH=amd64 go test -count 1 -timeout 30s -run "^TestGolden" ./resources/images
func TestGoldenFiltersMisc(t *testing.T) {
	t.Parallel()

	if skipGolden {
		t.Skip("Skip golden test on this architecture")
	}

	// Will be used to generate golden files.
	name := "filters_misc"

	files := `
-- hugo.toml --
-- assets/rotate270.jpg --
sourcefilename: ../testdata/exif/orientation6.jpg
-- assets/sunset.jpg --
sourcefilename: ../testdata/sunset.jpg
-- assets/gopher.png --
sourcefilename: ../testdata/gopher-hero8.png
-- layouts/index.html --
Home.
{{ $sunset := resources.Get "sunset.jpg" }}
{{ $sunsetGrayscale := $sunset.Filter (images.Grayscale) }}
{{ $gopher := resources.Get "gopher.png" }}
{{ $overlayFilter := images.Overlay $gopher 20 20 }}

{{ $textOpts := dict
  "color" "#fbfaf5"
  "linespacing" 8
  "size" 40
  "x" 25
  "y" 190
}}

{{/* These are sorted. */}}
{{ template "filters" (dict "name" "brightness-40.jpg" "img" $sunset "filters" (images.Brightness 40)) }}
{{ template "filters" (dict "name" "contrast-50.jpg" "img" $sunset "filters" (images.Contrast 50)) }}
{{ template "filters" (dict "name" "dither-default.jpg" "img" $sunset  "filters" (images.Dither)) }}
{{ template "filters" (dict "name" "gamma-1.667.jpg" "img" $sunset  "filters" (images.Gamma 1.667)) }}
{{ template "filters" (dict "name" "gaussianblur-5.jpg" "img" $sunset  "filters" (images.GaussianBlur 5)) }}
{{ template "filters" (dict "name" "grayscale.jpg" "img" $sunset  "filters" (images.Grayscale)) }}
{{ template "filters" (dict "name" "grayscale+colorize-180-50-20.jpg" "img" $sunset "filters" (slice images.Grayscale (images.Colorize 180 50 20))) }}
{{ template "filters" (dict "name" "colorbalance-180-50-20.jpg" "img" $sunset "filters"  (images.ColorBalance 180 50 20)) }}
{{ template "filters" (dict "name" "hue--15.jpg" "img" $sunset  "filters" (images.Hue -15)) }}
{{ template "filters" (dict "name" "invert.jpg" "img" $sunset  "filters" (images.Invert)) }}
{{ template "filters" (dict "name" "opacity-0.65.jpg" "img" $sunset  "filters" (images.Opacity 0.65)) }}
{{ template "filters" (dict "name" "overlay-20-20.jpg" "img" $sunset  "filters" ($overlayFilter)) }}
{{ template "filters" (dict "name" "padding-20-40-#976941.jpg" "img" $sunset  "filters" (images.Padding 20 40 "#976941" )) }}
{{ template "filters" (dict "name" "pixelate-10.jpg" "img" $sunset  "filters" (images.Pixelate 10)) }}
{{ template "filters" (dict "name" "rotate270.jpg" "img" (resources.Get "rotate270.jpg") "filters" images.AutoOrient) }}
{{ template "filters" (dict "name" "saturation-65.jpg" "img" $sunset  "filters" (images.Saturation 65)) }}
{{ template "filters" (dict "name" "sepia-80.jpg" "img" $sunsetGrayscale  "filters" (images.Sepia 80)) }}
{{ template "filters" (dict "name" "sigmoid-0.6--4.jpg" "img" $sunset  "filters" (images.Sigmoid 0.6 -4 )) }}
{{ template "filters" (dict "name" "text.jpg" "img" $sunset  "filters" (images.Text "Hugo Rocks!" $textOpts )) }}
{{ template "filters" (dict "name" "unsharpmask.jpg" "img" $sunset  "filters" (images.UnsharpMask 10 0.4 0.03)) }}


{{ define "filters"}}
{{ if lt (len (path.Ext .name)) 4 }}
	{{ errorf "No extension in %q" .name }}
{{ end }}
{{ $img := .img.Filter .filters }}
{{ $name := printf "images/%s" .name  }}
{{ with $img | resources.Copy $name }}
{{ .Publish }}
{{ end }}
{{ end }}
`

	runGolden(t, name, files)
}

func TestGoldenFiltersMask(t *testing.T) {
	t.Parallel()

	if skipGolden {
		t.Skip("Skip golden test on this architecture")
	}

	// Will be used to generate golden files.
	name := "filters_mask"

	files := `
-- hugo.toml --
[imaging]
  bgColor = '#ebcc34'
  hint = 'photo'
  quality = 75
  resampleFilter = 'Lanczos'
-- assets/sunset.jpg --
sourcefilename: ../testdata/sunset.jpg
-- assets/mask.png --
sourcefilename: ../testdata/mask.png

-- layouts/index.html --
Home.
{{ $sunset := resources.Get "sunset.jpg" }}
{{ $mask := resources.Get "mask.png" }}

{{ template "mask" (dict "name" "transparant.png" "base" $sunset  "mask" $mask) }}
{{ template "mask" (dict "name" "yellow.jpg" "base" $sunset  "mask" $mask) }}
{{ template "mask" (dict "name" "wide.jpg" "base" $sunset "mask" $mask "spec" "resize 600x200") }}


{{ define "mask"}}
{{ $ext := path.Ext .name }}
{{ if lt (len (path.Ext .name)) 4 }}
	{{ errorf "No extension in %q" .name }}
{{ end }}
{{ $format := strings.TrimPrefix "." $ext }}
{{ $spec := .spec | default (printf "resize 300x300 %s" $format) }}
{{ $filters := slice (images.Process $spec) (images.Mask .mask) }}
{{ $name := printf "images/%s" .name  }}
{{ $img := .base.Filter $filters }}
{{ with $img | resources.Copy $name }}
{{ .Publish }}
{{ end }}
{{ end }}
`

	runGolden(t, name, files)
}

func TestGoldenFiltersText(t *testing.T) {
	t.Parallel()

	if skipGolden {
		t.Skip("Skip golden test on this architecture")
	}

	// Will be used to generate golden files.
	name := "filters_text"

	files := `
-- hugo.toml --
-- assets/sunset.jpg --
sourcefilename: ../testdata/sunset.jpg

-- layouts/index.html --
Home.
{{ $sunset := resources.Get "sunset.jpg" }}
{{ $textOpts := dict
  "color" "#fbfaf5"
  "linespacing" 8
  "size" 28
  "x" (div $sunset.Width 2 | int)
  "alignx" "center"
  "y" 190
}}

{{ $text := "Pariatur deserunt sunt nisi sunt tempor quis eu. Sint et nulla enim officia sunt cupidatat. Eu amet ipsum qui velit cillum cillum ad Lorem in non ad aute." }}
{{ template "filters" (dict "name" "text_alignx-center.jpg" "img" $sunset  "filters" (images.Text $text $textOpts )) }}
{{ $textOpts = (dict "alignx" "right") | merge $textOpts }}
{{ template "filters" (dict "name" "text_alignx-right.jpg" "img" $sunset  "filters" (images.Text $text $textOpts )) }}
{{ $textOpts = (dict "alignx" "left") | merge $textOpts }}
{{ template "filters" (dict "name" "text_alignx-left.jpg" "img" $sunset  "filters" (images.Text $text $textOpts )) }}

{{ define "filters"}}
{{ if lt (len (path.Ext .name)) 4 }}
	{{ errorf "No extension in %q" .name }}
{{ end }}
{{ $img := .img.Filter .filters }}
{{ $name := printf "images/%s" .name  }}
{{ with $img | resources.Copy $name }}
{{ .Publish }}
{{ end }}
{{ end }}
`

	runGolden(t, name, files)
}

func TestGoldenProcessMisc(t *testing.T) {
	t.Parallel()

	if skipGolden {
		t.Skip("Skip golden test on this architecture")
	}

	// Will be used to generate golden files.
	name := "process_misc"

	files := `
-- hugo.toml --
-- assets/giphy.gif --
sourcefilename: ../testdata/giphy.gif
-- assets/sunset.jpg --
sourcefilename: ../testdata/sunset.jpg
-- assets/gopher.png --
sourcefilename: ../testdata/gopher-hero8.png
-- layouts/index.html --
Home.
{{ $sunset := resources.Get "sunset.jpg" }}
{{ $sunsetGrayscale := $sunset.Filter (images.Grayscale) }}
{{ $gopher := resources.Get "gopher.png" }}
{{ $giphy := resources.Get "giphy.gif" }}


{{/* These are sorted. The end file name will be created from the spec + extension, so make sure these are unique. */}}
{{ template "process" (dict "spec" "crop 500x200 smart" "img" $sunset) }}
{{ template "process" (dict "spec" "fill 500x200 smart" "img" $sunset) }}
{{ template "process" (dict "spec" "fit 500x200 smart" "img" $sunset) }}
{{ template "process" (dict "spec" "resize 100x100 gif" "img" $giphy) }}
{{ template "process" (dict "spec" "resize 100x100 r180" "img" $gopher) }}
{{ template "process" (dict "spec" "resize 300x300 jpg #b31280" "img" $gopher) }}

{{ define "process"}}
{{ $img := .img.Process .spec }}
{{ $ext := path.Ext $img.RelPermalink }}
{{ $name := printf "images/%s%s" (.spec | anchorize) $ext  }}
{{ with $img | resources.Copy $name }}
{{ .Publish }}
{{ end }}
{{ end }}
`

	runGolden(t, name, files)
}

func TestGoldenFuncs(t *testing.T) {
	t.Parallel()

	if skipGolden {
		t.Skip("Skip golden test on this architecture")
	}

	// Will be used to generate golden files.
	name := "funcs"

	files := `
-- hugo.toml --
-- assets/sunset.jpg --
sourcefilename: ../testdata/sunset.jpg

-- layouts/index.html --
Home.

{{ template "copy" (dict "name" "qr-default.png" "img" (images.QR "https://gohugo.io"))  }}
{{ template "copy" (dict "name" "qr-level-high_scale-6.png" "img" (images.QR "https://gohugo.io" (dict "level" "high" "scale" 6)))  }}

{{ define "copy"}}
{{ if lt (len (path.Ext .name)) 4 }}
	{{ errorf "No extension in %q" .name }}
{{ end }}
{{ $img := .img }}
{{ $name := printf "images/%s" .name  }}
{{ with $img | resources.Copy $name }}
{{ .Publish }}
{{ end }}
{{ end }}
`

	runGolden(t, name, files)
}

func runGolden(t testing.TB, name, files string) *hugolib.IntegrationTestBuilder {
	t.Helper()

	c := hugolib.Test(t, files, hugolib.TestOptWithOSFs()) // hugolib.TestOptWithPrintAndKeepTempDir(true))
	c.AssertFileContent("public/index.html", "Home.")

	outputDir := filepath.Join(c.H.Conf.WorkingDir(), "public", "images")
	goldenBaseDir := filepath.Join("testdata", "images_golden")
	goldenDir := filepath.Join(goldenBaseDir, name)
	if goldenOpts.writeGoldenFiles {
		c.Assert(htesting.IsRealCI(), qt.IsFalse)
		c.Assert(os.MkdirAll(goldenBaseDir, 0o777), qt.IsNil)
		c.Assert(os.RemoveAll(goldenDir), qt.IsNil)
		c.Assert(hugio.CopyDir(hugofs.Os, outputDir, goldenDir, nil), qt.IsNil)
		return c
	}

	if goldenOpts.devMode {
		c.Assert(htesting.IsRealCI(), qt.IsFalse)
		return c
	}

	decodeAll := func(f *os.File) []image.Image {
		c.Helper()

		var images []image.Image

		if strings.HasSuffix(f.Name(), ".gif") {
			gif, err := gif.DecodeAll(f)
			c.Assert(err, qt.IsNil, qt.Commentf(f.Name()))
			images = make([]image.Image, len(gif.Image))
			for i, img := range gif.Image {
				images[i] = img
			}
		} else {
			img, _, err := image.Decode(f)
			c.Assert(err, qt.IsNil, qt.Commentf(f.Name()))
			images = append(images, img)
		}
		return images
	}

	entries1, err := os.ReadDir(outputDir)
	c.Assert(err, qt.IsNil)
	entries2, err := os.ReadDir(goldenDir)
	c.Assert(err, qt.IsNil)
	c.Assert(len(entries1), qt.Equals, len(entries2))
	for i, e1 := range entries1 {
		c.Assert(filepath.Ext(e1.Name()), qt.Not(qt.Equals), "")
		func() {
			e2 := entries2[i]

			f1, err := os.Open(filepath.Join(outputDir, e1.Name()))
			c.Assert(err, qt.IsNil)
			defer f1.Close()

			f2, err := os.Open(filepath.Join(goldenDir, e2.Name()))
			c.Assert(err, qt.IsNil)
			defer f2.Close()

			imgs2 := decodeAll(f2)
			imgs1 := decodeAll(f1)
			c.Assert(len(imgs1), qt.Equals, len(imgs2))

			if !usesFMA {
				c.Assert(e1, eq, e2)
				_, err = f1.Seek(0, 0)
				c.Assert(err, qt.IsNil)
				_, err = f2.Seek(0, 0)
				c.Assert(err, qt.IsNil)

				hash1, _, err := hashing.XXHashFromReader(f1)
				c.Assert(err, qt.IsNil)
				hash2, _, err := hashing.XXHashFromReader(f2)
				c.Assert(err, qt.IsNil)

				c.Assert(hash1, qt.Equals, hash2)
			}

			for i, img1 := range imgs1 {
				img2 := imgs2[i]
				nrgba1 := image.NewNRGBA(img1.Bounds())
				gift.New().Draw(nrgba1, img1)
				nrgba2 := image.NewNRGBA(img2.Bounds())
				gift.New().Draw(nrgba2, img2)
				c.Assert(goldenEqual(nrgba1, nrgba2), qt.Equals, true, qt.Commentf(e1.Name()))
			}
		}()
	}
	return c
}

// goldenEqual compares two NRGBA images.  It is used in golden tests only.
// A small tolerance is allowed on architectures using "fused multiply and add"
// (FMA) instruction to accommodate for floating-point rounding differences
// with control golden images that were generated on amd64 architecture.
// See https://golang.org/ref/spec#Floating_point_operators
// and https://github.com/gohugoio/hugo/issues/6387 for more information.
//
// Based on https://github.com/disintegration/gift/blob/a999ff8d5226e5ab14b64a94fca07c4ac3f357cf/gift_test.go#L598-L625
// Copyright (c) 2014-2019 Grigory Dryapak
// Licensed under the MIT License.
func goldenEqual(img1, img2 *image.NRGBA) bool {
	maxDiff := 0
	if runtime.GOARCH != "amd64" {
		// The golden files are created using the AMD64 architecture.
		// Be lenient on other platforms due to floaging point and dithering differences.
		maxDiff = 15
	}
	if !img1.Rect.Eq(img2.Rect) {
		return false
	}
	if len(img1.Pix) != len(img2.Pix) {
		return false
	}
	for i := 0; i < len(img1.Pix); i++ {
		diff := int(img1.Pix[i]) - int(img2.Pix[i])
		if diff < 0 {
			diff = -diff
		}
		if diff > maxDiff {
			return false
		}
	}
	return true
}

// We don't have a CI test environment for these, and there are known dithering issues that makes these time consuming to maintain.
var skipGolden = runtime.GOARCH == "ppc64" || runtime.GOARCH == "ppc64le" || runtime.GOARCH == "s390x"

// usesFMA indicates whether "fused multiply and add" (FMA) instruction is
// used.  The command "grep FMADD go/test/codegen/floats.go" can help keep
// the FMA-using architecture list updated.
var usesFMA = runtime.GOARCH == "s390x" ||
	runtime.GOARCH == "ppc64" ||
	runtime.GOARCH == "ppc64le" ||
	runtime.GOARCH == "arm64" ||
	runtime.GOARCH == "riscv64"
