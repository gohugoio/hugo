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
	"fmt"
	"image"
	"image/draw"

	"github.com/disintegration/gift"
)

var _ gift.Filter = (*overlayFilter)(nil)

type overlayFilter struct {
	src  ImageSource
	x, y int
}

func (f overlayFilter) Draw(dst draw.Image, src image.Image, options *gift.Options) {
	overlaySrc, err := f.src.DecodeImage()
	if err != nil {
		panic(fmt.Sprintf("failed to decode image: %s", err))
	}

	gift.New().Draw(dst, src)
	gift.New().DrawAt(dst, overlaySrc, image.Pt(f.x, f.y), gift.OverOperator)
}

func (f overlayFilter) Bounds(srcBounds image.Rectangle) image.Rectangle {
	return image.Rect(0, 0, srcBounds.Dx(), srcBounds.Dy())
}
