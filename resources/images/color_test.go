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
	"image/color"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/htesting/hqt"
)

func TestHexStringToColor(t *testing.T) {
	c := qt.New(t)

	for _, test := range []struct {
		arg    string
		expect any
	}{
		{"f", false},
		{"#f", false},
		{"#fffffff", false},
		{"fffffff", false},
		{"#fff", color.White},
		{"fff", color.White},
		{"FFF", color.White},
		{"FfF", color.White},
		{"#ffffff", color.White},
		{"ffffff", color.White},
		{"#000", color.Black},
		{"#4287f5", color.RGBA{R: 0x42, G: 0x87, B: 0xf5, A: 0xff}},
		{"777", color.RGBA{R: 0x77, G: 0x77, B: 0x77, A: 0xff}},
	} {

		test := test
		c.Run(test.arg, func(c *qt.C) {
			c.Parallel()

			result, err := hexStringToColorGo(test.arg)

			if b, ok := test.expect.(bool); ok && !b {
				c.Assert(err, qt.Not(qt.IsNil))
				return
			}

			c.Assert(err, qt.IsNil)
			c.Assert(result, qt.DeepEquals, test.expect)
		})

	}
}

func TestColorToHexString(t *testing.T) {
	c := qt.New(t)

	for _, test := range []struct {
		arg    color.Color
		expect string
	}{
		{color.White, "#ffffff"},
		{color.Black, "#000000"},
		{color.RGBA{R: 0x42, G: 0x87, B: 0xf5, A: 0xff}, "#4287f5"},

		// 50% opacity.
		// Note that the .Colors (dominant colors) received from the Image resource
		// will always have an alpha value of 0xff.
		{color.RGBA{R: 0x42, G: 0x87, B: 0xf5, A: 0x80}, "#4287f580"},
	} {

		test := test
		c.Run(test.expect, func(c *qt.C) {
			c.Parallel()

			result := ColorGoToHexString(test.arg)

			c.Assert(result, qt.Equals, test.expect)
		})

	}
}

func TestAddColorToPalette(t *testing.T) {
	c := qt.New(t)

	palette := color.Palette{color.White, color.Black}

	c.Assert(AddColorToPalette(color.White, palette), qt.HasLen, 2)

	blue1, _ := hexStringToColorGo("34c3eb")
	blue2, _ := hexStringToColorGo("34c3eb")
	white, _ := hexStringToColorGo("fff")

	c.Assert(AddColorToPalette(white, palette), qt.HasLen, 2)
	c.Assert(AddColorToPalette(blue1, palette), qt.HasLen, 3)
	c.Assert(AddColorToPalette(blue2, palette), qt.HasLen, 3)
}

func TestReplaceColorInPalette(t *testing.T) {
	c := qt.New(t)

	palette := color.Palette{color.White, color.Black}
	offWhite, _ := hexStringToColorGo("fcfcfc")

	ReplaceColorInPalette(offWhite, palette)

	c.Assert(palette, qt.HasLen, 2)
	c.Assert(palette[0], qt.Equals, offWhite)
}

func TestColorLuminance(t *testing.T) {
	c := qt.New(t)
	c.Assert(hexStringToColor("#000000").Luminance(), hqt.IsSameFloat64, 0.0)
	c.Assert(hexStringToColor("#768a9a").Luminance(), hqt.IsSameFloat64, 0.24361603589088263)
	c.Assert(hexStringToColor("#d5bc9f").Luminance(), hqt.IsSameFloat64, 0.5261577672685374)
	c.Assert(hexStringToColor("#ffffff").Luminance(), hqt.IsSameFloat64, 1.0)
}
