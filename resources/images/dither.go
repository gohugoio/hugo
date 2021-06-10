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
	"image/color"
	"image/draw"
	"strings"

	"github.com/disintegration/gift"
)

var _ gift.Filter = (*ditherFilter)(nil)

var matrixs = map[string][][]float32{

	strings.ToLower("FloydSteinberg"):    {{0, 0, 7.0 / 16.0}, {3.0 / 16.0, 5.0 / 16.0, 1.0 / 16.0}},
	strings.ToLower("JarvisJudiceNinke"): {{0, 0, 0, 7.0 / 48.0, 5.0 / 48.0}, {3.0 / 48.0, 5.0 / 48.0, 7.0 / 48.0, 5.0 / 48.0, 3.0 / 48.0}, {1.0 / 48.0, 3.0 / 48.0, 5.0 / 48.0, 3.0 / 48.0, 1.0 / 48.0}},
	strings.ToLower("Stucki"):            {{0, 0, 0, 8.0 / 42.0, 4.0 / 42.0}, {2.0 / 42.0, 4.0 / 42.0, 8.0 / 42.0, 4.0 / 42.0, 2.0 / 42.0}, {1.0 / 42.0, 2.0 / 42.0, 4.0 / 42.0, 2.0 / 42.0, 1.0 / 42.0}},
	strings.ToLower("Atkinson"):          {{0, 0, 1.0 / 8.0, 1.0 / 8.0}, {1.0 / 8.0, 1.0 / 8.0, 1.0 / 8.0, 0}, {0, 1.0 / 8.0, 0, 0}},
	strings.ToLower("Burkes"):            {{0, 0, 0, 8.0 / 32.0, 4.0 / 32.0}, {2.0 / 32.0, 4.0 / 32.0, 8.0 / 32.0, 4.0 / 32.0, 2.0 / 32.0}},
	strings.ToLower("Sierra"):            {{0, 0, 0, 5.0 / 32.0, 3.0 / 32.0}, {2.0 / 32.0, 4.0 / 32.0, 5.0 / 32.0, 4.0 / 32.0, 2.0 / 32.0}, {0, 2.0 / 32.0, 3.0 / 32.0, 2.0 / 32.0, 0}},
	strings.ToLower("TwoRowSierra"):      {{0, 0, 0, 4.0 / 16.0, 3.0 / 16.0}, {1.0 / 32.0, 2.0 / 32.0, 3.0 / 32.0, 2.0 / 32.0, 1.0 / 32.0}},
	strings.ToLower("SierraLite"):        {{0, 0, 2.0 / 4.0}, {1.0 / 4.0, 1.0 / 4.0, 0}},
}

var (
	defaultMatrix  = [][]float32{{0, 0, 7.0 / 16.0}, {3.0 / 16.0, 5.0 / 16.0, 1.0 / 16.0}}
	defaultPalette = color.Palette{color.Black, color.White}
)

type ditherFilter struct {
	config string
}

func (f ditherFilter) Draw(dst draw.Image, src image.Image, options *gift.Options) {
	var (
		palette []color.Color
		matrix  [][]float32
	)
	parts := strings.Fields(f.config)
	for _, part := range parts {
		part = strings.ToLower(part)
		if m, ok := matrixs[part]; ok {
			matrix = m
		} else if part[0] == '#' {
			p, err := hexStringToColor(part[1:])
			palette = append(palette, p)
			if err != nil {
				return
			}
		}
	}

	// set to default matrix if matrix not found
	if len(matrix) <= 0 {
		matrix = defaultMatrix
	}

	// append palette if only has one color or less 
	if len(palette) <= 1 {
		palette = append(palette, defaultPalette...)
	}

	rect := dst.Bounds()
	errImg := NewErrorImage(rect)
	shift := findShift(matrix)

	animation := make(chan draw.Image)

	pixPerFrame := (rect.Dx() * rect.Dy()) / 1

	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			// using the closest color
			r, e, _ := findColor(errImg.PixelErrorAt(x, y), src.At(x, y), palette)
			dst.Set(x, y, r)
			errImg.SetPixelError(x, y, e)

			if (y != 0 && x != 0) && (((y*rect.Dy())+x)%pixPerFrame == 0) {
				animation <- dst
			}

			// diffusing the error using the diffusion matrix
			for i, v1 := range matrix {
				for j, v2 := range v1 {
					errImg.SetPixelError(x+j+shift, y+i,
						errImg.PixelErrorAt(x+j+shift, y+i).Add(errImg.PixelErrorAt(x, y).Mul(v2)))
				}
			}
		}
	}
}

func (f ditherFilter) Bounds(srcBounds image.Rectangle) (dstBounds image.Rectangle) {
	dstBounds = image.Rect(0, 0, srcBounds.Dx(), srcBounds.Dy())
	return
}

func findShift(matrix [][]float32) int {
	for _, v1 := range matrix {
		for j, v2 := range v1 {
			if v2 > 0.0 {
				return -j + 1
			}
		}
	}
	return 0
}

// findColor determines the closest color in a palette given the pixel color and the error
//
// It returns the closest color, the updated error and the distance between the error and the color
func findColor(err color.Color, pix color.Color, pal color.Palette) (color.RGBA, PixelError, uint16) {
	var errR, errG, errB,
		pixR, pixG, pixB,
		colR, colG, colB int16
	_errR, _errG, _errB, _ := err.RGBA()
	_pixR, _pixG, _pixB, _ := pix.RGBA()

	// Low-pass filter
	errR = int16(float32(int16(_errR)) * 0.75)
	errG = int16(float32(int16(_errG)) * 0.75)
	errB = int16(float32(int16(_errB)) * 0.75)

	pixR = int16(uint8(_pixR)) + errR
	pixG = int16(uint8(_pixG)) + errG
	pixB = int16(uint8(_pixB)) + errB

	var index int
	var minDiff uint16 = 1<<16 - 1

	for i, col := range pal {
		_colR, _colG, _colB, _ := col.RGBA()

		colR = int16(uint8(_colR))
		colG = int16(uint8(_colG))
		colB = int16(uint8(_colB))
		var distance = abs(pixR-colR) + abs(pixG-colG) + abs(pixB-colB)

		if distance < minDiff {
			index = i
			minDiff = distance
		}
	}

	_colR, _colG, _colB, _ := pal[index].RGBA()

	colR = int16(uint8(_colR))
	colG = int16(uint8(_colG))
	colB = int16(uint8(_colB))

	return color.RGBA{uint8(colR), uint8(colG), uint8(colB), 255},
		PixelError{float32(pixR - colR),
			float32(pixG - colG),
			float32(pixB - colB),
			1<<16 - 1},
		minDiff
}

// abs gives the absolute value of a signed integer
func abs(x int16) uint16 {
	if x < 0 {
		return uint16(-x)
	}
	return uint16(x)
}
