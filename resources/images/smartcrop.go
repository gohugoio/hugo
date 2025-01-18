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
	"image"
	"math"

	"github.com/disintegration/gift"

	"github.com/muesli/smartcrop"
)

const (
	// Do not change.
	smartCropIdentifier = "smart"
	SmartCropAnchor     = 1000
	// This is just a increment, starting on 0. If Smart Crop improves its cropping, we
	// need a way to trigger a re-generation of the crops in the wild, so increment this.
	smartCropVersionNumber = 0
)

func (p *ImageProcessor) newSmartCropAnalyzer(filter gift.Resampling) smartcrop.Analyzer {
	return smartcrop.NewAnalyzer(imagingResizer{p: p, filter: filter})
}

// Needed by smartcrop
type imagingResizer struct {
	p      *ImageProcessor
	filter gift.Resampling
}

func (r imagingResizer) Resize(img image.Image, width, height uint) image.Image {
	// See https://github.com/gohugoio/hugo/issues/7955#issuecomment-861710681
	scaleX, scaleY := calcFactorsNfnt(width, height, float64(img.Bounds().Dx()), float64(img.Bounds().Dy()))
	if width == 0 {
		width = uint(math.Ceil(float64(img.Bounds().Dx()) / scaleX))
	}
	if height == 0 {
		height = uint(math.Ceil(float64(img.Bounds().Dy()) / scaleY))
	}
	result, _ := r.p.Filter(img, gift.Resize(int(width), int(height), r.filter))
	return result
}

func (p *ImageProcessor) smartCrop(img image.Image, width, height int, filter gift.Resampling) (image.Rectangle, error) {
	if width <= 0 || height <= 0 {
		return image.Rectangle{}, nil
	}

	srcBounds := img.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()

	if srcW <= 0 || srcH <= 0 {
		return image.Rectangle{}, nil
	}

	if srcW == width && srcH == height {
		return srcBounds, nil
	}

	smart := p.newSmartCropAnalyzer(filter)

	rect, err := smart.FindBestCrop(img, width, height)
	if err != nil {
		return image.Rectangle{}, err
	}

	return img.Bounds().Intersect(rect), nil
}

// Calculates scaling factors using old and new image dimensions.
// Code borrowed from https://github.com/nfnt/resize/blob/83c6a9932646f83e3267f353373d47347b6036b2/resize.go#L593
func calcFactorsNfnt(width, height uint, oldWidth, oldHeight float64) (scaleX, scaleY float64) {
	if width == 0 {
		if height == 0 {
			scaleX = 1.0
			scaleY = 1.0
		} else {
			scaleY = oldHeight / float64(height)
			scaleX = scaleY
		}
	} else {
		scaleX = oldWidth / float64(width)
		if height == 0 {
			scaleY = scaleX
		} else {
			scaleY = oldHeight / float64(height)
		}
	}
	return
}
