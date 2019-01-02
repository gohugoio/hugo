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

package resources

import (
	"image"

	"github.com/disintegration/imaging"
	"github.com/muesli/smartcrop"
)

const (
	// Do not change.
	smartCropIdentifier = "smart"

	// This is just a increment, starting on 1. If Smart Crop improves its cropping, we
	// need a way to trigger a re-generation of the crops in the wild, so increment this.
	smartCropVersionNumber = 1
)

// Needed by smartcrop
type imagingResizer struct {
	filter imaging.ResampleFilter
}

func (r imagingResizer) Resize(img image.Image, width, height uint) image.Image {
	return imaging.Resize(img, int(width), int(height), r.filter)
}

func newSmartCropAnalyzer(filter imaging.ResampleFilter) smartcrop.Analyzer {
	return smartcrop.NewAnalyzer(imagingResizer{filter: filter})
}

func smartCrop(img image.Image, width, height int, anchor imaging.Anchor, filter imaging.ResampleFilter) (*image.NRGBA, error) {

	if width <= 0 || height <= 0 {
		return &image.NRGBA{}, nil
	}

	srcBounds := img.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()

	if srcW <= 0 || srcH <= 0 {
		return &image.NRGBA{}, nil
	}

	if srcW == width && srcH == height {
		return imaging.Clone(img), nil
	}

	smart := newSmartCropAnalyzer(filter)

	rect, err := smart.FindBestCrop(img, width, height)

	if err != nil {
		return nil, err
	}

	b := img.Bounds().Intersect(rect)

	cropped, err := imaging.Crop(img, b), nil
	if err != nil {
		return nil, err
	}

	return imaging.Resize(cropped, width, height, filter), nil

}
