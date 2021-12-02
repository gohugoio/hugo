// Copyright 2020 The Hugo Authors. All rights reserved.
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
	"image"
	"image/draw"

	"github.com/disintegration/gift"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goitalic"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

var _ gift.Filter = (*textFilter)(nil)

type textFilter struct {
	text, color string
	x, y        int
	size        float64
}

func (f textFilter) Draw(dst draw.Image, src image.Image, options *gift.Options) {
	if f.color == "" {
		f.color = "#000000"
	}

	color, err := hexStringToColor(f.color)
	if err != nil {
		panic(err)
	}

	otf, err := opentype.Parse(goitalic.TTF)
	if err != nil {
		panic(err)
	}
	face, err := opentype.NewFace(otf, &opentype.FaceOptions{
		Size:    f.size,
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		panic(err)
	}

	d := font.Drawer{
		Dst:  dst,
		Src:  image.NewUniform(color),
		Face: face,
		Dot:  fixed.P(f.x, f.y),
	}

	gift.New().Draw(dst, src)
	d.DrawString(f.text)

}

func (f textFilter) Bounds(srcBounds image.Rectangle) image.Rectangle {
	return image.Rect(0, 0, srcBounds.Dx(), srcBounds.Dy())
}
