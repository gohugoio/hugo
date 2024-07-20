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

package images

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/disintegration/gift"
)

var _ gift.Filter = (*opacityFilter)(nil)

type opacityFilter struct {
	opacity float32
}

func (f opacityFilter) Draw(dst draw.Image, src image.Image, options *gift.Options) {
	// 0 is fully transparent and 255 is opaque.
	alpha := uint8(f.opacity * 255)
	mask := image.NewUniform(color.Alpha{alpha})
	draw.DrawMask(dst, dst.Bounds(), src, image.Point{}, mask, image.Point{}, draw.Over)
}

func (f opacityFilter) Bounds(srcBounds image.Rectangle) image.Rectangle {
	return image.Rect(0, 0, srcBounds.Dx(), srcBounds.Dy())
}
