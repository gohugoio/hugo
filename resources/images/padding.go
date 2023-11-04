// Copyright 2023 The Hugo Authors. All rights reserved.
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
	"image/color"
	"image/draw"

	"github.com/disintegration/gift"
)

var _ gift.Filter = (*paddingFilter)(nil)

type paddingFilter struct {
	top, right, bottom, left int
	ccolor                   color.Color // canvas color
}

func (f paddingFilter) Draw(dst draw.Image, src image.Image, options *gift.Options) {
	w := src.Bounds().Dx() + f.left + f.right
	h := src.Bounds().Dy() + f.top + f.bottom

	if w < 1 {
		panic("final image width will be less than 1 pixel: check padding values")
	}
	if h < 1 {
		panic("final image height will be less than 1 pixel: check padding values")
	}

	i := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(i, i.Bounds(), image.NewUniform(f.ccolor), image.Point{}, draw.Src)
	gift.New().Draw(dst, i)
	gift.New().DrawAt(dst, src, image.Pt(f.left, f.top), gift.OverOperator)
}

func (f paddingFilter) Bounds(srcBounds image.Rectangle) image.Rectangle {
	return image.Rect(0, 0, srcBounds.Dx()+f.left+f.right, srcBounds.Dy()+f.top+f.bottom)
}
