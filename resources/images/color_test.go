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
)

func TestHexStringToColor(t *testing.T) {
	c := qt.New(t)

	for _, test := range []struct {
		arg    string
		expect interface{}
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

			result, err := hexStringToColor(test.arg)

			if b, ok := test.expect.(bool); ok && !b {
				c.Assert(err, qt.Not(qt.IsNil))
				return
			}

			c.Assert(err, qt.IsNil)
			c.Assert(result, qt.DeepEquals, test.expect)
		})

	}
}

func TestAddColorToPalette(t *testing.T) {
	c := qt.New(t)

	palette := color.Palette{color.White, color.Black}

	c.Assert(AddColorToPalette(color.White, palette), qt.HasLen, 2)

	blue1, _ := hexStringToColor("34c3eb")
	blue2, _ := hexStringToColor("34c3eb")
	white, _ := hexStringToColor("fff")

	c.Assert(AddColorToPalette(white, palette), qt.HasLen, 2)
	c.Assert(AddColorToPalette(blue1, palette), qt.HasLen, 3)
	c.Assert(AddColorToPalette(blue2, palette), qt.HasLen, 3)

}

func TestReplaceColorInPalette(t *testing.T) {
	c := qt.New(t)

	palette := color.Palette{color.White, color.Black}
	offWhite, _ := hexStringToColor("fcfcfc")

	ReplaceColorInPalette(offWhite, palette)

	c.Assert(palette, qt.HasLen, 2)
	c.Assert(palette[0], qt.Equals, offWhite)
}
