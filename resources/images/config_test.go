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
	"fmt"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestDecodeConfig(t *testing.T) {
	c := qt.New(t)
	m := map[string]interface{}{
		"quality":        42,
		"resampleFilter": "NearestNeighbor",
		"anchor":         "topLeft",
	}

	imagingConfig, err := DecodeConfig(m)

	c.Assert(err, qt.IsNil)
	imaging := imagingConfig.Cfg
	c.Assert(imaging.Quality, qt.Equals, 42)
	c.Assert(imaging.ResampleFilter, qt.Equals, "nearestneighbor")
	c.Assert(imaging.Anchor, qt.Equals, "topleft")

	m = map[string]interface{}{}

	imagingConfig, err = DecodeConfig(m)
	c.Assert(err, qt.IsNil)
	imaging = imagingConfig.Cfg
	c.Assert(imaging.Quality, qt.Equals, defaultJPEGQuality)
	c.Assert(imaging.ResampleFilter, qt.Equals, "box")
	c.Assert(imaging.Anchor, qt.Equals, "smart")

	_, err = DecodeConfig(map[string]interface{}{
		"quality": 123,
	})
	c.Assert(err, qt.Not(qt.IsNil))

	_, err = DecodeConfig(map[string]interface{}{
		"resampleFilter": "asdf",
	})
	c.Assert(err, qt.Not(qt.IsNil))

	_, err = DecodeConfig(map[string]interface{}{
		"anchor": "asdf",
	})
	c.Assert(err, qt.Not(qt.IsNil))

	imagingConfig, err = DecodeConfig(map[string]interface{}{
		"anchor": "Smart",
	})
	imaging = imagingConfig.Cfg
	c.Assert(err, qt.IsNil)
	c.Assert(imaging.Anchor, qt.Equals, "smart")

	imagingConfig, err = DecodeConfig(map[string]interface{}{
		"exif": map[string]interface{}{
			"disableLatLong": true,
		},
	})
	c.Assert(err, qt.IsNil)
	imaging = imagingConfig.Cfg
	c.Assert(imaging.Exif.DisableLatLong, qt.Equals, true)
	c.Assert(imaging.Exif.ExcludeFields, qt.Equals, "GPS|Exif|Exposure[M|P|B]|Contrast|Resolution|Sharp|JPEG|Metering|Sensing|Saturation|ColorSpace|Flash|WhiteBalance")

}

func TestDecodeImageConfig(t *testing.T) {
	for i, this := range []struct {
		in     string
		expect interface{}
	}{
		{"300x400", newImageConfig(300, 400, 0, 0, "", "", "")},
		{"300x400 #fff", newImageConfig(300, 400, 0, 0, "", "", "fff")},
		{"100x200 bottomRight", newImageConfig(100, 200, 0, 0, "", "BottomRight", "")},
		{"10x20 topleft Lanczos", newImageConfig(10, 20, 0, 0, "Lanczos", "topleft", "")},
		{"linear left 10x r180", newImageConfig(10, 0, 0, 180, "linear", "left", "")},
		{"x20 riGht Cosine q95", newImageConfig(0, 20, 95, 0, "cosine", "right", "")},

		{"", false},
		{"foo", false},
	} {

		result, err := DecodeImageConfig("resize", this.in, Imaging{})
		if b, ok := this.expect.(bool); ok && !b {
			if err == nil {
				t.Errorf("[%d] parseImageConfig didn't return an expected error", i)
			}
		} else {
			if err != nil {
				t.Fatalf("[%d] err: %s", i, err)
			}
			if fmt.Sprint(result) != fmt.Sprint(this.expect) {
				t.Fatalf("[%d] got\n%v\n but expected\n%v", i, result, this.expect)
			}
		}
	}
}

func newImageConfig(width, height, quality, rotate int, filter, anchor, bgColor string) ImageConfig {
	var c ImageConfig
	c.Action = "resize"
	c.Width = width
	c.Height = height
	c.Quality = quality
	c.Rotate = rotate
	c.BgColorStr = bgColor
	c.BgColor, _ = hexStringToColor(bgColor)

	if filter != "" {
		filter = strings.ToLower(filter)
		if v, ok := imageFilters[filter]; ok {
			c.Filter = v
			c.FilterStr = filter
		}
	}

	if anchor != "" {
		anchor = strings.ToLower(anchor)
		if v, ok := anchorPositions[anchor]; ok {
			c.Anchor = v
			c.AnchorStr = anchor
		}
	}

	return c
}
