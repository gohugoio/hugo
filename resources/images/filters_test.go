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

package images

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/common/hashing"
)

func TestFilterHash(t *testing.T) {
	c := qt.New(t)

	f := &Filters{}

	c.Assert(hashing.HashString(f.Grayscale()), qt.Equals, hashing.HashString(f.Grayscale()))
	c.Assert(hashing.HashString(f.Grayscale()), qt.Not(qt.Equals), hashing.HashString(f.Invert()))
	c.Assert(hashing.HashString(f.Gamma(32)), qt.Not(qt.Equals), hashing.HashString(f.Gamma(33)))
	c.Assert(hashing.HashString(f.Gamma(32)), qt.Equals, hashing.HashString(f.Gamma(32)))
}

// Issue 12536. A transparent padding color must survive the trip through an
// alpha-premultiplied destination image. The destination type is derived from
// the source image type, so the padding color came out wrong whenever the
// source decoded to an alpha-premultiplied *image.RGBA (e.g. an opaque PNG
// written by a previous "resize x200 png" step), while the same filter applied
// to a JPEG source looked correct.
func TestPaddingTransparentColor(t *testing.T) {
	c := qt.New(t)

	ccolor, err := hexStringToColorGo("#00f7")
	c.Assert(err, qt.IsNil)

	pad := paddingFilter{
		top: 2, right: 2, bottom: 2, left: 2,
		ccolor: ccolor,
	}

	for _, test := range []struct {
		name string
		src  image.Image
	}{
		{"YCbCr source (e.g. JPEG)", image.NewYCbCr(image.Rect(0, 0, 10, 10), image.YCbCrSubsampleRatio420)},
		{"RGBA source (e.g. opaque PNG)", image.NewRGBA(image.Rect(0, 0, 10, 10))},
	} {
		c.Run(test.name, func(c *qt.C) {
			img, err := (&ImageProcessor{}).Filter(test.src, pad)
			c.Assert(err, qt.IsNil)

			// Encode and decode as PNG so the padding pixel is read back the
			// same way a browser would see it.
			var buf bytes.Buffer
			c.Assert(png.Encode(&buf, img), qt.IsNil)
			decoded, err := png.Decode(&buf)
			c.Assert(err, qt.IsNil)

			got := color.NRGBAModel.Convert(decoded.At(1, 1)).(color.NRGBA)
			c.Assert(got, qt.Equals, color.NRGBA{R: 0x00, G: 0x00, B: 0xff, A: 0x77})
		})
	}
}
