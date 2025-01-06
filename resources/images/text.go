// Copyright 2021 The Hugo Authors. All rights reserved.
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
	"io"
	"strings"

	"github.com/disintegration/gift"
	"github.com/gohugoio/hugo/common/hugio"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

var _ gift.Filter = (*textFilter)(nil)

type textFilter struct {
	text        string
	color       color.Color
	x, y        int
	alignx      string
	size        float64
	linespacing int
	fontSource  hugio.ReadSeekCloserProvider
}

func (f textFilter) Draw(dst draw.Image, src image.Image, options *gift.Options) {
	// Load and parse font
	ttf := goregular.TTF
	if f.fontSource != nil {
		rs, err := f.fontSource.ReadSeekCloser()
		if err != nil {
			panic(err)
		}
		defer rs.Close()
		ttf, err = io.ReadAll(rs)
		if err != nil {
			panic(err)
		}
	}
	otf, err := opentype.Parse(ttf)
	if err != nil {
		panic(err)
	}

	// Set font options
	face, err := opentype.NewFace(otf, &opentype.FaceOptions{
		Size:    f.size,
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		panic(err)
	}

	d := font.Drawer{
		Dst:  dst,
		Src:  image.NewUniform(f.color),
		Face: face,
	}

	gift.New().Draw(dst, src)

	maxWidth := dst.Bounds().Dx() - 20

	var availableWidth int
	switch f.alignx {
	case "right":
		availableWidth = f.x
	case "center":
		availableWidth = min((maxWidth-f.x), f.x) * 2
	case "left":
		availableWidth = maxWidth - f.x
	}

	fontHeight := face.Metrics().Ascent.Ceil()

	// Calculate lines, consider and include linebreaks
	finalLines := []string{}
	f.text = strings.ReplaceAll(f.text, "\r", "")
	for _, line := range strings.Split(f.text, "\n") {
		currentLine := ""
		// Break each line at the maximum width.
		for _, str := range strings.Fields(line) {
			fieldStrWidth := font.MeasureString(face, str)
			currentLineStrWidth := font.MeasureString(face, currentLine)

			if (currentLineStrWidth.Ceil() + fieldStrWidth.Ceil()) >= availableWidth {
				finalLines = append(finalLines, currentLine)
				currentLine = ""
			}
			currentLine += str + " "
		}
		finalLines = append(finalLines, currentLine)
	}

	// Correct y position based on font and size
	f.y = f.y + fontHeight

	// Start position
	y := f.y

	// Draw text line by line
	for _, line := range finalLines {
		line = strings.TrimSpace(line)
		strWidth := font.MeasureString(face, line)
		var x int
		switch f.alignx {
		case "right":
			x = f.x - strWidth.Ceil()
		case "center":
			x = f.x - (strWidth.Ceil() / 2)

		case "left":
			x = f.x
		}
		d.Dot = fixed.P(x, y)
		d.DrawString(line)
		y = y + fontHeight + f.linespacing
	}
}

func (f textFilter) Bounds(srcBounds image.Rectangle) image.Rectangle {
	return image.Rect(0, 0, srcBounds.Dx(), srcBounds.Dy())
}
