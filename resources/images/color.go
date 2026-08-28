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
	"encoding/hex"
	"fmt"
	"image/color"
	"math"
	"slices"
	"strings"

	"github.com/gohugoio/hugo/common/hstrings"
)

type colorGoProvider interface {
	ColorGo() color.Color
}

type Color struct {
	// The color.
	color color.Color

	// The color prefixed with a #.
	hex string

	// The relative luminance of the color.
	luminance float64
}

// Luminance as defined by w3.org.
// See https://www.w3.org/TR/WCAG21/#dfn-relative-luminance
func (c Color) Luminance() float64 {
	return c.luminance
}

// ColorGo returns the color as a color.Color.
// For internal use only.
func (c Color) ColorGo() color.Color {
	return c.color
}

// ColorHex returns the color as a hex string  prefixed with a #.
func (c Color) ColorHex() string {
	return c.hex
}

// String returns the color as a hex string prefixed with a #.
func (c Color) String() string {
	return c.hex
}

// For hashstructure. This struct is used in template func options
// that needs to be able to hash a Color.
// For internal use only.
func (c Color) Key() string {
	return c.hex
}

func (c *Color) init() error {
	// Normalize to the canonical non-premultiplied form so that equal colors
	// compare equal no matter where they came from (a hex string or the
	// pixels of an image).
	nrgba := color.NRGBAModel.Convert(c.color).(color.NRGBA)
	c.color = nrgba
	c.hex = ColorGoToHexString(nrgba)
	c.luminance = 0.2126*c.toSRGB(nrgba.R) + 0.7152*c.toSRGB(nrgba.G) + 0.0722*c.toSRGB(nrgba.B)
	return nil
}

func (c Color) toSRGB(i uint8) float64 {
	v := float64(i) / 255
	if v <= 0.04045 {
		return v / 12.92
	} else {
		return math.Pow((v+0.055)/1.055, 2.4)
	}
}

// AddColorToPalette adds c as the first color in p if not already there.
// Note that it does no additional checks, so callers must make sure
// that the palette is valid for the relevant format.
func AddColorToPalette(c color.Color, p color.Palette) color.Palette {
	// Compare color values rather than the interface values, so that e.g. a
	// color.NRGBA is found even if the palette holds it as a color.RGBA.
	cr, cg, cb, ca := c.RGBA()
	found := slices.ContainsFunc(p, func(v color.Color) bool {
		vr, vg, vb, va := v.RGBA()
		return cr == vr && cg == vg && cb == vb && ca == va
	})

	if !found {
		p = append(color.Palette{c}, p...)
	}

	return p
}

// ReplaceColorInPalette will replace the color in palette p closest to c in Euclidean
// R,G,B,A space with c.
func ReplaceColorInPalette(c color.Color, p color.Palette) {
	p[p.Index(c)] = c
}

// ColorGoToHexString converts a color.Color to a hex string.
func ColorGoToHexString(c color.Color) string {
	// Hex notation is not alpha-premultiplied, so convert via NRGBA to
	// un-premultiply the color components of any transparent color.
	nrgba := color.NRGBAModel.Convert(c).(color.NRGBA)
	if nrgba.A == 0xff {
		return fmt.Sprintf("#%.2x%.2x%.2x", nrgba.R, nrgba.G, nrgba.B)
	}
	return fmt.Sprintf("#%.2x%.2x%.2x%.2x", nrgba.R, nrgba.G, nrgba.B, nrgba.A)
}

// ColorGoToColor converts a color.Color to a Color.
func ColorGoToColor(c color.Color) Color {
	cc := Color{color: c}
	if err := cc.init(); err != nil {
		panic(err)
	}
	return cc
}

func hexStringToColor(s string) Color {
	c, err := hexStringToColorGo(s)
	if err != nil {
		panic(err)
	}
	return ColorGoToColor(c)
}

// HexStringsToColors converts a slice of hex strings to a slice of Colors.
func HexStringsToColors(s ...string) []Color {
	var colors []Color
	for _, v := range s {
		colors = append(colors, hexStringToColor(v))
	}
	return colors
}

func toColorGo(v any) (color.Color, bool, error) {
	switch vv := v.(type) {
	case colorGoProvider:
		return vv.ColorGo(), true, nil
	default:
		s, ok := hstrings.ToString(v)
		if !ok {
			return nil, false, nil
		}
		c, err := hexStringToColorGo(s)
		if err != nil {
			return nil, false, err
		}
		return c, true, nil
	}
}

func hexStringToColorGo(s string) (color.Color, error) {
	s = strings.TrimPrefix(s, "#")

	if len(s) != 3 && len(s) != 4 && len(s) != 6 && len(s) != 8 {
		return nil, fmt.Errorf("invalid color code: %q", s)
	}

	s = strings.ToLower(s)

	if len(s) == 3 || len(s) == 4 {
		var v strings.Builder
		for _, r := range s {
			v.WriteString(string(r) + string(r))
		}
		s = v.String()
	}

	// Standard colors.
	if s == "ffffff" {
		return color.White, nil
	}

	if s == "000000" {
		return color.Black, nil
	}

	// Set Alfa to white.
	if len(s) == 6 {
		s += "ff"
	}

	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}

	// Hex notation is not alpha-premultiplied, so this must be a color.NRGBA.
	// Returning the components as a color.RGBA (which is alpha-premultiplied
	// by definition) produces an invalid color whenever the alpha component
	// is less than the color components, and drawing it corrupts the color
	// on the alpha-premultiplied code paths. See issue #12536.
	return color.NRGBA{b[0], b[1], b[2], b[3]}, nil
}
