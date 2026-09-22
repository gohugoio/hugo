// Copyright 2026 The Hugo Authors. All rights reserved.
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
	"image/gif"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/gift"
)

func TestSmartCropRotationPreservesSource(t *testing.T) {
	c := qt.New(t)
	p := &ImageProcessor{}
	palette := color.Palette{color.Black, color.White}
	frame := image.NewPaletted(image.Rect(0, 0, 80, 40), palette)
	frame.SetColorIndex(10, 20, 1)
	pixels := bytes.Clone(frame.Pix)
	anim := &giphy{
		Image: frame,
		gif: &gif.GIF{
			Image:  []*image.Paletted{frame, image.NewPaletted(frame.Bounds(), palette)},
			Config: image.Config{Width: 80, Height: 40},
		},
	}
	frames := anim.GetFrames()
	for _, src := range []image.Image{frame, anim} {
		for _, action := range []string{ActionCrop, ActionFill} {
			_, err := p.FiltersFromConfig(src, ImageConfig{
				Action: action, Width: 20, Height: 10, Rotate: 90,
				Anchor: SmartCropAnchor, Filter: gift.BoxResampling,
			})
			c.Assert(err, qt.IsNil)
			c.Assert(src.Bounds(), qt.Equals, image.Rect(0, 0, 80, 40))
			c.Assert(frame.Pix, qt.DeepEquals, pixels)
			c.Assert(anim.GetFrames(), qt.DeepEquals, frames)
			c.Assert(anim.GetImageConfig(), qt.Equals, image.Config{Width: 80, Height: 40})
		}
	}
}

func TestExpandCropRectToMinSize(t *testing.T) {
	c := qt.New(t)

	for _, test := range []struct {
		name         string
		r            image.Rectangle
		bounds       image.Rectangle
		width        int
		height       int
		expectedRect image.Rectangle
	}{
		{
			name:         "top left",
			r:            image.Rect(0, 0, 899, 560),
			bounds:       image.Rect(0, 0, 900, 562),
			width:        900,
			height:       561,
			expectedRect: image.Rect(0, 0, 900, 561),
		},
		{
			name:         "bottom right",
			r:            image.Rect(1, 2, 900, 562),
			bounds:       image.Rect(0, 0, 900, 562),
			width:        900,
			height:       561,
			expectedRect: image.Rect(0, 1, 900, 562),
		},
		{
			name:         "centered",
			r:            image.Rect(10, 10, 109, 109),
			bounds:       image.Rect(0, 0, 200, 200),
			width:        101,
			height:       101,
			expectedRect: image.Rect(9, 9, 110, 110),
		},
		{
			name:         "source too small",
			r:            image.Rect(0, 0, 899, 560),
			bounds:       image.Rect(0, 0, 899, 560),
			width:        900,
			height:       561,
			expectedRect: image.Rect(0, 0, 899, 560),
		},
	} {
		c.Run(test.name, func(c *qt.C) {
			got := expandCropRectToMinSize(test.r, test.bounds, test.width, test.height)
			c.Assert(got, qt.Equals, test.expectedRect)
		})
	}
}
