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
	"image/draw"

	"github.com/disintegration/gift"
	"github.com/gohugoio/hugo/resources/images/exif"
)

var _ gift.Filter = (*autoOrientFilter)(nil)

var transformationFilters = map[int]gift.Filter{
	2: gift.FlipHorizontal(),
	3: gift.Rotate180(),
	4: gift.FlipVertical(),
	5: gift.Transpose(),
	6: gift.Rotate270(),
	7: gift.Transverse(),
	8: gift.Rotate90(),
}

type autoOrientFilter struct{}

type ImageFilterFromOrientationProvider interface {
	AutoOrient(exifInfo *exif.ExifInfo) gift.Filter
}

func (f autoOrientFilter) Draw(dst draw.Image, src image.Image, options *gift.Options) {
	panic("not supported")
}

func (f autoOrientFilter) Bounds(srcBounds image.Rectangle) image.Rectangle {
	panic("not supported")
}

func (f autoOrientFilter) AutoOrient(exifInfo *exif.ExifInfo) gift.Filter {
	if exifInfo != nil {
		if orientation, ok := exifInfo.Tags["Orientation"].(int); ok {
			if filter, ok := transformationFilters[orientation]; ok {
				return filter
			}
		}
	}

	return nil
}
