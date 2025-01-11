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
	_ "image/jpeg"
	"testing"

	"github.com/gohugoio/hugo/resources/images/imagetesting"
)

// Note, if you're enabling writeGoldenFiles on a MacOS ARM 64 you need to run the test with GOARCH=amd64, e.g.
func TestImagesGoldenFiltersMisc(t *testing.T) {
	t.Parallel()

	if imagetesting.SkipGoldenTests {
		t.Skip("Skip golden test on this architecture")
	}

	// Will be used as the base folder for generated images.
	name := "filters/misc"

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
{{ $sunset := (resources.Get "sunset.jpg").Resize "x300" }}
{{ $sunsetGrayscale := $sunset.Filter (images.Grayscale) }}
{{ $gopher := (resources.Get "gopher.png").Resize "x80" }}
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

	opts := imagetesting.DefaultGoldenOpts
	opts.T = t
	opts.Name = name
	opts.Files = files

	imagetesting.RunGolden(opts)
}

func TestImagesGoldenFiltersMask(t *testing.T) {
	t.Parallel()

	if imagetesting.SkipGoldenTests {
		t.Skip("Skip golden test on this architecture")
	}

	// Will be used as the base folder for generated images.
	name := "filters/mask"

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
{{/* This looks a little odd, but is correct and the recommended way to do this.
This will 1. Scale the image to x300, 2. Apply the mask, 3. Create the final image with background color #323ea. 
It's possible to have multiple images.Process filters in the chain, but for the options for the final image (target format, bgGolor etc.),
the last entry will win.
*/}}
{{ template "mask" (dict "name" "blue.jpg" "base" $sunset "mask" $mask "spec" "resize x300 #323ea8") }}

{{ define "mask"}}
{{ $ext := path.Ext .name }}
{{ if lt (len (path.Ext .name)) 4 }}
	{{ errorf "No extension in %q" .name }}
{{ end }}
{{ $format := strings.TrimPrefix "." $ext }}
{{ $spec := .spec | default (printf "resize x300 %s" $format) }}
{{ $filters := slice (images.Process $spec) (images.Mask .mask) }}
{{ $name := printf "images/%s" .name  }}
{{ $img := .base.Filter $filters }}
{{ with $img | resources.Copy $name }}
{{ .Publish }}
{{ end }}
{{ end }}
`

	opts := imagetesting.DefaultGoldenOpts
	opts.T = t
	opts.Name = name
	opts.Files = files

	imagetesting.RunGolden(opts)
}

func TestImagesGoldenFiltersText(t *testing.T) {
	t.Parallel()

	if imagetesting.SkipGoldenTests {
		t.Skip("Skip golden test on this architecture")
	}

	// Will be used as the base folder for generated images.
	name := "filters/text"

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

	opts := imagetesting.DefaultGoldenOpts
	opts.T = t
	opts.Name = name
	opts.Files = files

	imagetesting.RunGolden(opts)
}

func TestImagesGoldenProcessMisc(t *testing.T) {
	t.Parallel()

	if imagetesting.SkipGoldenTests {
		t.Skip("Skip golden test on this architecture")
	}

	// Will be used as the base folder for generated images.
	name := "process/misc"

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

	opts := imagetesting.DefaultGoldenOpts
	opts.T = t
	opts.Name = name
	opts.Files = files

	imagetesting.RunGolden(opts)
}
