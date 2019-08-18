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
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gohugoio/hugo/resources/internal"

	"github.com/disintegration/imaging"
	"github.com/mitchellh/mapstructure"
)

const (
	defaultJPEGQuality    = 75
	defaultResampleFilter = "box"
)

var (
	imageFormats = map[string]imaging.Format{
		".jpg":  imaging.JPEG,
		".jpeg": imaging.JPEG,
		".png":  imaging.PNG,
		".tif":  imaging.TIFF,
		".tiff": imaging.TIFF,
		".bmp":  imaging.BMP,
		".gif":  imaging.GIF,
	}

	// Add or increment if changes to an image format's processing requires
	// re-generation.
	imageFormatsVersions = map[imaging.Format]int{
		imaging.PNG: 2, // Floyd Steinberg dithering
	}

	// Increment to mark all processed images as stale. Only use when absolutely needed.
	// See the finer grained smartCropVersionNumber and imageFormatsVersions.
	mainImageVersionNumber = 0

	// Increment to mark all traced SVGs as stale.
	traceVersionNumber = 0
)

var anchorPositions = map[string]imaging.Anchor{
	strings.ToLower("Center"):      imaging.Center,
	strings.ToLower("TopLeft"):     imaging.TopLeft,
	strings.ToLower("Top"):         imaging.Top,
	strings.ToLower("TopRight"):    imaging.TopRight,
	strings.ToLower("Left"):        imaging.Left,
	strings.ToLower("Right"):       imaging.Right,
	strings.ToLower("BottomLeft"):  imaging.BottomLeft,
	strings.ToLower("Bottom"):      imaging.Bottom,
	strings.ToLower("BottomRight"): imaging.BottomRight,
}

var imageFilters = map[string]imaging.ResampleFilter{
	strings.ToLower("NearestNeighbor"):   imaging.NearestNeighbor,
	strings.ToLower("Box"):               imaging.Box,
	strings.ToLower("Linear"):            imaging.Linear,
	strings.ToLower("Hermite"):           imaging.Hermite,
	strings.ToLower("MitchellNetravali"): imaging.MitchellNetravali,
	strings.ToLower("CatmullRom"):        imaging.CatmullRom,
	strings.ToLower("BSpline"):           imaging.BSpline,
	strings.ToLower("Gaussian"):          imaging.Gaussian,
	strings.ToLower("Lanczos"):           imaging.Lanczos,
	strings.ToLower("Hann"):              imaging.Hann,
	strings.ToLower("Hamming"):           imaging.Hamming,
	strings.ToLower("Blackman"):          imaging.Blackman,
	strings.ToLower("Bartlett"):          imaging.Bartlett,
	strings.ToLower("Welch"):             imaging.Welch,
	strings.ToLower("Cosine"):            imaging.Cosine,
}

func ImageFormatFromExt(ext string) (imaging.Format, bool) {
	f, found := imageFormats[ext]
	return f, found
}

func DecodeConfig(m map[string]interface{}) (Imaging, error) {
	var i Imaging
	if err := mapstructure.WeakDecode(m, &i); err != nil {
		return i, err
	}

	if i.Quality == 0 {
		i.Quality = defaultJPEGQuality
	} else if i.Quality < 0 || i.Quality > 100 {
		return i, errors.New("JPEG quality must be a number between 1 and 100")
	}

	if i.Anchor == "" || strings.EqualFold(i.Anchor, SmartCropIdentifier) {
		i.Anchor = SmartCropIdentifier
	} else {
		i.Anchor = strings.ToLower(i.Anchor)
		if _, found := anchorPositions[i.Anchor]; !found {
			return i, errors.New("invalid anchor value in imaging config")
		}
	}

	if i.ResampleFilter == "" {
		i.ResampleFilter = defaultResampleFilter
	} else {
		filter := strings.ToLower(i.ResampleFilter)
		_, found := imageFilters[filter]
		if !found {
			return i, fmt.Errorf("%q is not a valid resample filter", filter)
		}
		i.ResampleFilter = filter
	}

	return i, nil
}

func DecodeImageConfig(action, config string, defaults Imaging) (ImageConfig, error) {
	var (
		c   ImageConfig
		err error
	)

	c.Action = action

	if c.Action == "trace" {
		return c, nil
	}

	if config == "" {
		return c, errors.New("image config cannot be empty")
	}

	parts := strings.Fields(config)
	for _, part := range parts {
		part = strings.ToLower(part)

		if part == SmartCropIdentifier {
			c.AnchorStr = SmartCropIdentifier
		} else if pos, ok := anchorPositions[part]; ok {
			c.Anchor = pos
			c.AnchorStr = part
		} else if filter, ok := imageFilters[part]; ok {
			c.Filter = filter
			c.FilterStr = part
		} else if part[0] == 'q' {
			c.Quality, err = strconv.Atoi(part[1:])
			if err != nil {
				return c, err
			}
			if c.Quality < 1 || c.Quality > 100 {
				return c, errors.New("quality ranges from 1 to 100 inclusive")
			}
		} else if part[0] == 'r' {
			c.Rotate, err = strconv.Atoi(part[1:])
			if err != nil {
				return c, err
			}
		} else if strings.Contains(part, "x") {
			widthHeight := strings.Split(part, "x")
			if len(widthHeight) <= 2 {
				first := widthHeight[0]
				if first != "" {
					c.Width, err = strconv.Atoi(first)
					if err != nil {
						return c, err
					}
				}

				if len(widthHeight) == 2 {
					second := widthHeight[1]
					if second != "" {
						c.Height, err = strconv.Atoi(second)
						if err != nil {
							return c, err
						}
					}
				}
			} else {
				return c, errors.New("invalid image dimensions")
			}

		}
	}

	if c.Width == 0 && c.Height == 0 {
		return c, errors.New("must provide Width or Height")
	}

	if c.FilterStr == "" {
		c.FilterStr = defaults.ResampleFilter
		c.Filter = imageFilters[c.FilterStr]
	}

	if c.AnchorStr == "" {
		c.AnchorStr = defaults.Anchor
		if !strings.EqualFold(c.AnchorStr, SmartCropIdentifier) {
			c.Anchor = anchorPositions[c.AnchorStr]
		}
	}

	return c, nil
}

// ImageConfig holds configuration to create a new image from an existing one, resize etc.
type ImageConfig struct {
	Action string

	// Configures the trace action.
	TraceOptions TraceOptions

	// Quality ranges from 1 to 100 inclusive, higher is better.
	// This is only relevant for JPEG images.
	// Default is 75.
	Quality int

	// Rotate rotates an image by the given angle counter-clockwise.
	// The rotation will be performed first.
	Rotate int

	Width  int
	Height int

	Filter    imaging.ResampleFilter
	FilterStr string

	Anchor    imaging.Anchor
	AnchorStr string
}

func (i ImageConfig) Key(format imaging.Format) string {

	if i.Action == "trace" {
		key := internal.NewResourceTransformationKey(fmt.Sprintf("%s_%d", i.Action, traceVersionNumber), i.TraceOptions)
		return key.Value()
	}

	k := strconv.Itoa(i.Width) + "x" + strconv.Itoa(i.Height)
	if i.Action != "" {
		k += "_" + i.Action
	}
	if i.Quality > 0 {
		k += "_q" + strconv.Itoa(i.Quality)
	}
	if i.Rotate != 0 {
		k += "_r" + strconv.Itoa(i.Rotate)
	}
	anchor := i.AnchorStr
	if anchor == SmartCropIdentifier {
		anchor = anchor + strconv.Itoa(smartCropVersionNumber)
	}

	k += "_" + i.FilterStr

	if strings.EqualFold(i.Action, "fill") {
		k += "_" + anchor
	}

	if v, ok := imageFormatsVersions[format]; ok {
		k += "_" + strconv.Itoa(v)
	}

	if mainImageVersionNumber > 0 {
		k += "_" + strconv.Itoa(mainImageVersionNumber)
	}

	return k
}

// Imaging contains default image processing configuration. This will be fetched
// from site (or language) config.
type Imaging struct {
	// Default image quality setting (1-100). Only used for JPEG images.
	Quality int

	// Resample filter used. See https://github.com/disintegration/imaging
	ResampleFilter string

	// The anchor used in Fill. Default is "smart", i.e. Smart Crop.
	Anchor string
}
