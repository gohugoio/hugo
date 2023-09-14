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
	m := map[string]any{
		"quality":        42,
		"resampleFilter": "NearestNeighbor",
		"anchor":         "topLeft",
	}

	imagingConfig, err := DecodeConfig(m)

	c.Assert(err, qt.IsNil)
	conf := imagingConfig.Config
	c.Assert(conf.Imaging.Quality, qt.Equals, 42)
	c.Assert(conf.Imaging.ResampleFilter, qt.Equals, "nearestneighbor")
	c.Assert(conf.Imaging.Anchor, qt.Equals, "topleft")

	m = map[string]any{}

	imagingConfig, err = DecodeConfig(m)
	c.Assert(err, qt.IsNil)
	conf = imagingConfig.Config
	c.Assert(conf.Imaging.ResampleFilter, qt.Equals, "box")
	c.Assert(conf.Imaging.Anchor, qt.Equals, "smart")

	_, err = DecodeConfig(map[string]any{
		"quality": 123,
	})
	c.Assert(err, qt.IsNotNil)

	_, err = DecodeConfig(map[string]any{
		"resampleFilter": "asdf",
	})
	c.Assert(err, qt.IsNotNil)

	_, err = DecodeConfig(map[string]any{
		"anchor": "asdf",
	})
	c.Assert(err, qt.IsNotNil)

	imagingConfig, err = DecodeConfig(map[string]any{
		"anchor": "Smart",
	})
	conf = imagingConfig.Config
	c.Assert(err, qt.IsNil)
	c.Assert(conf.Imaging.Anchor, qt.Equals, "smart")

	imagingConfig, err = DecodeConfig(map[string]any{
		"exif": map[string]any{
			"disableLatLong": true,
		},
	})
	c.Assert(err, qt.IsNil)
	conf = imagingConfig.Config
	c.Assert(conf.Imaging.Exif.DisableLatLong, qt.IsTrue)
	c.Assert(conf.Imaging.Exif.ExcludeFields, qt.Equals, "GPS|Exif|Exposure[M|P|B]|Contrast|Resolution|Sharp|JPEG|Metering|Sensing|Saturation|ColorSpace|Flash|WhiteBalance")
}

func TestDecodeImageConfig(t *testing.T) {
	for i, this := range []struct {
		action string
		in     string
		expect any
	}{
		{"resize", "300x400", newImageConfig("resize", 300, 400, 75, 0, "box", "smart", "")},
		{"resize", "300x400 #fff", newImageConfig("resize", 300, 400, 75, 0, "box", "smart", "fff")},
		{"resize", "100x200 bottomRight", newImageConfig("resize", 100, 200, 75, 0, "box", "BottomRight", "")},
		{"resize", "10x20 topleft Lanczos", newImageConfig("resize", 10, 20, 75, 0, "Lanczos", "topleft", "")},
		{"resize", "linear left 10x r180", newImageConfig("resize", 10, 0, 75, 180, "linear", "left", "")},
		{"resize", "x20 riGht Cosine q95", newImageConfig("resize", 0, 20, 95, 0, "cosine", "right", "")},
		{"crop", "300x400", newImageConfig("crop", 300, 400, 75, 0, "box", "smart", "")},
		{"fill", "300x400", newImageConfig("fill", 300, 400, 75, 0, "box", "smart", "")},
		{"fit", "300x400", newImageConfig("fit", 300, 400, 75, 0, "box", "smart", "")},

		{"resize", "", false},
		{"resize", "foo", false},
		{"crop", "100x", false},
		{"fill", "100x", false},
		{"fit", "100x", false},
		{"foo", "100x", false},
	} {

		cfg, err := DecodeConfig(nil)
		if err != nil {
			t.Fatal(err)
		}
		result, err := DecodeImageConfig(this.action, this.in, cfg, PNG)
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

func newImageConfig(action string, width, height, quality, rotate int, filter, anchor, bgColor string) ImageConfig {
	var c ImageConfig = GetDefaultImageConfig(action, nil)
	c.TargetFormat = PNG
	c.Hint = 2
	c.Width = width
	c.Height = height
	c.Quality = quality
	c.qualitySetForImage = quality != 75
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
		if anchor == smartCropIdentifier {
			c.AnchorStr = anchor
		} else {
			anchor = strings.ToLower(anchor)
			if v, ok := anchorPositions[anchor]; ok {
				c.Anchor = v
				c.AnchorStr = anchor
			}
		}
	}

	return c
}
