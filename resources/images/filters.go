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

// Package images provides template functions for manipulating images.
package images

import (
	"fmt"
	"image/color"

	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/resources/resource"

	"github.com/disintegration/gift"
	"github.com/spf13/cast"
)

// Increment for re-generation of images using these filters.
const filterAPIVersion = 0

type Filters struct{}

// Process creates a filter that processes an image using the given specification.
func (*Filters) Process(spec any) gift.Filter {
	return filter{
		Options: newFilterOpts(spec),
		Filter: processFilter{
			spec: cast.ToString(spec),
		},
	}
}

// Overlay creates a filter that overlays src at position x y.
func (*Filters) Overlay(src ImageSource, x, y any) gift.Filter {
	return filter{
		Options: newFilterOpts(src.Key(), x, y),
		Filter:  overlayFilter{src: src, x: cast.ToInt(x), y: cast.ToInt(y)},
	}
}

// Opacity creates a filter that changes the opacity of an image.
// The opacity parameter must be in range (0, 1).
func (*Filters) Opacity(opacity any) gift.Filter {
	return filter{
		Options: newFilterOpts(opacity),
		Filter:  opacityFilter{opacity: cast.ToFloat32(opacity)},
	}
}

// Text creates a filter that draws text with the given options.
func (*Filters) Text(text string, options ...any) gift.Filter {
	tf := textFilter{
		text:        text,
		color:       "#ffffff",
		size:        20,
		x:           10,
		y:           10,
		linespacing: 2,
	}

	var opt maps.Params
	if len(options) > 0 {
		opt = maps.MustToParamsAndPrepare(options[0])
		for option, v := range opt {
			switch option {
			case "color":
				tf.color = cast.ToString(v)
			case "size":
				tf.size = cast.ToFloat64(v)
			case "x":
				tf.x = cast.ToInt(v)
			case "y":
				tf.y = cast.ToInt(v)
			case "linespacing":
				tf.linespacing = cast.ToInt(v)
			case "font":
				if err, ok := v.(error); ok {
					panic(fmt.Sprintf("invalid font source: %s", err))
				}
				fontSource, ok1 := v.(hugio.ReadSeekCloserProvider)
				identifier, ok2 := v.(resource.Identifier)

				if !(ok1 && ok2) {
					panic(fmt.Sprintf("invalid text font source: %T", v))
				}

				tf.fontSource = fontSource

				// The input value isn't hashable and will not make a stable key.
				// Replace it with a string in the map used as basis for the
				// hash string.
				opt["font"] = identifier.Key()

			}
		}
	}

	return filter{
		Options: newFilterOpts(text, opt),
		Filter:  tf,
	}
}

// Padding creates a filter that resizes the image canvas without resizing the
// image. The last argument is the canvas color, expressed as an RGB or RGBA
// hexadecimal color. The default value is `ffffffff` (opaque white). The
// preceding arguments are the padding values, in pixels, using the CSS
// shorthand property syntax. Negative padding values will crop the image. The
// signature is images.Padding V1 [V2] [V3] [V4] [COLOR].
func (*Filters) Padding(args ...any) gift.Filter {
	if len(args) < 1 || len(args) > 5 {
		panic("the padding filter requires between 1 and 5 arguments")
	}

	var top, right, bottom, left int
	var ccolor color.Color = color.White // canvas color
	var err error

	_args := args // preserve original args for most stable hash

	if vcs, ok := (args[len(args)-1]).(string); ok {
		ccolor, err = hexStringToColor(vcs)
		if err != nil {
			panic("invalid canvas color: specify RGB or RGBA using hex notation")
		}
		args = args[:len(args)-1]
		if len(args) == 0 {
			panic("not enough arguments: provide one or more padding values using the CSS shorthand property syntax")
		}
	}

	var vals []int
	for _, v := range args {
		vi := cast.ToInt(v)
		if vi > 5000 {
			panic("padding values must not exceed 5000 pixels")
		}
		vals = append(vals, vi)
	}

	switch len(args) {
	case 1:
		top, right, bottom, left = vals[0], vals[0], vals[0], vals[0]
	case 2:
		top, right, bottom, left = vals[0], vals[1], vals[0], vals[1]
	case 3:
		top, right, bottom, left = vals[0], vals[1], vals[2], vals[1]
	case 4:
		top, right, bottom, left = vals[0], vals[1], vals[2], vals[3]
	default:
		panic(fmt.Sprintf("too many padding values: received %d, expected maximum of 4", len(args)))
	}

	return filter{
		Options: newFilterOpts(_args...),
		Filter: paddingFilter{
			top:    top,
			right:  right,
			bottom: bottom,
			left:   left,
			ccolor: ccolor,
		},
	}
}

// Brightness creates a filter that changes the brightness of an image.
// The percentage parameter must be in range (-100, 100).
func (*Filters) Brightness(percentage any) gift.Filter {
	return filter{
		Options: newFilterOpts(percentage),
		Filter:  gift.Brightness(cast.ToFloat32(percentage)),
	}
}

// ColorBalance creates a filter that changes the color balance of an image.
// The percentage parameters for each color channel (red, green, blue) must be in range (-100, 500).
func (*Filters) ColorBalance(percentageRed, percentageGreen, percentageBlue any) gift.Filter {
	return filter{
		Options: newFilterOpts(percentageRed, percentageGreen, percentageBlue),
		Filter:  gift.ColorBalance(cast.ToFloat32(percentageRed), cast.ToFloat32(percentageGreen), cast.ToFloat32(percentageBlue)),
	}
}

// Colorize creates a filter that produces a colorized version of an image.
// The hue parameter is the angle on the color wheel, typically in range (0, 360).
// The saturation parameter must be in range (0, 100).
// The percentage parameter specifies the strength of the effect, it must be in range (0, 100).
func (*Filters) Colorize(hue, saturation, percentage any) gift.Filter {
	return filter{
		Options: newFilterOpts(hue, saturation, percentage),
		Filter:  gift.Colorize(cast.ToFloat32(hue), cast.ToFloat32(saturation), cast.ToFloat32(percentage)),
	}
}

// Contrast creates a filter that changes the contrast of an image.
// The percentage parameter must be in range (-100, 100).
func (*Filters) Contrast(percentage any) gift.Filter {
	return filter{
		Options: newFilterOpts(percentage),
		Filter:  gift.Contrast(cast.ToFloat32(percentage)),
	}
}

// Gamma creates a filter that performs a gamma correction on an image.
// The gamma parameter must be positive. Gamma = 1 gives the original image.
// Gamma less than 1 darkens the image and gamma greater than 1 lightens it.
func (*Filters) Gamma(gamma any) gift.Filter {
	return filter{
		Options: newFilterOpts(gamma),
		Filter:  gift.Gamma(cast.ToFloat32(gamma)),
	}
}

// GaussianBlur creates a filter that applies a gaussian blur to an image.
func (*Filters) GaussianBlur(sigma any) gift.Filter {
	return filter{
		Options: newFilterOpts(sigma),
		Filter:  gift.GaussianBlur(cast.ToFloat32(sigma)),
	}
}

// Grayscale creates a filter that produces a grayscale version of an image.
func (*Filters) Grayscale() gift.Filter {
	return filter{
		Filter: gift.Grayscale(),
	}
}

// Hue creates a filter that rotates the hue of an image.
// The hue angle shift is typically in range -180 to 180.
func (*Filters) Hue(shift any) gift.Filter {
	return filter{
		Options: newFilterOpts(shift),
		Filter:  gift.Hue(cast.ToFloat32(shift)),
	}
}

// Invert creates a filter that negates the colors of an image.
func (*Filters) Invert() gift.Filter {
	return filter{
		Filter: gift.Invert(),
	}
}

// Pixelate creates a filter that applies a pixelation effect to an image.
func (*Filters) Pixelate(size any) gift.Filter {
	return filter{
		Options: newFilterOpts(size),
		Filter:  gift.Pixelate(cast.ToInt(size)),
	}
}

// Saturation creates a filter that changes the saturation of an image.
func (*Filters) Saturation(percentage any) gift.Filter {
	return filter{
		Options: newFilterOpts(percentage),
		Filter:  gift.Saturation(cast.ToFloat32(percentage)),
	}
}

// Sepia creates a filter that produces a sepia-toned version of an image.
func (*Filters) Sepia(percentage any) gift.Filter {
	return filter{
		Options: newFilterOpts(percentage),
		Filter:  gift.Sepia(cast.ToFloat32(percentage)),
	}
}

// Sigmoid creates a filter that changes the contrast of an image using a sigmoidal function and returns the adjusted image.
// It's a non-linear contrast change useful for photo adjustments as it preserves highlight and shadow detail.
func (*Filters) Sigmoid(midpoint, factor any) gift.Filter {
	return filter{
		Options: newFilterOpts(midpoint, factor),
		Filter:  gift.Sigmoid(cast.ToFloat32(midpoint), cast.ToFloat32(factor)),
	}
}

// UnsharpMask creates a filter that sharpens an image.
// The sigma parameter is used in a gaussian function and affects the radius of effect.
// Sigma must be positive. Sharpen radius roughly equals 3 * sigma.
// The amount parameter controls how much darker and how much lighter the edge borders become. Typically between 0.5 and 1.5.
// The threshold parameter controls the minimum brightness change that will be sharpened. Typically between 0 and 0.05.
func (*Filters) UnsharpMask(sigma, amount, threshold any) gift.Filter {
	return filter{
		Options: newFilterOpts(sigma, amount, threshold),
		Filter:  gift.UnsharpMask(cast.ToFloat32(sigma), cast.ToFloat32(amount), cast.ToFloat32(threshold)),
	}
}

type filter struct {
	Options filterOpts
	gift.Filter
}

// For cache-busting.
type filterOpts struct {
	Version int
	Vals    any
}

func newFilterOpts(vals ...any) filterOpts {
	return filterOpts{
		Version: filterAPIVersion,
		Vals:    vals,
	}
}
