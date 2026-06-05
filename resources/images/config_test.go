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
	c.Assert(err, qt.Not(qt.IsNil))

	_, err = DecodeConfig(map[string]any{
		"resampleFilter": "asdf",
	})
	c.Assert(err, qt.Not(qt.IsNil))

	_, err = DecodeConfig(map[string]any{
		"anchor": "asdf",
	})
	c.Assert(err, qt.Not(qt.IsNil))

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
	c.Assert(conf.Imaging.Exif.DisableLatLong, qt.Equals, true)
	c.Assert(conf.Imaging.Exif.ExcludeFields, qt.Equals, "GPS|Exif|Exposure[M|P|B]|Contrast|Resolution|Sharp|JPEG|Metering|Sensing|Saturation|ColorSpace|Flash|WhiteBalance")

	// AVIF: default is speed 10.
	imagingConfig, err = DecodeConfig(map[string]any{})
	c.Assert(err, qt.IsNil)
	c.Assert(imagingConfig.Config.Imaging.Avif.EncoderSpeed, qt.Equals, defaultAvifEncoderSpeed)

	// AVIF: override via config.
	imagingConfig, err = DecodeConfig(map[string]any{
		"avif": map[string]any{"encoderSpeed": 5},
	})
	c.Assert(err, qt.IsNil)
	c.Assert(imagingConfig.Config.Imaging.Avif.EncoderSpeed, qt.Equals, 5)

	// AVIF: out-of-range rejected.
	_, err = DecodeConfig(map[string]any{
		"avif": map[string]any{"encoderSpeed": 11},
	})
	c.Assert(err, qt.ErrorMatches, ".*encoderSpeed must be between.*")

	imagingConfig, err = DecodeConfig(map[string]any{
		"avif": map[string]any{"encoderSpeed": 1},
	})
	c.Assert(err, qt.IsNil)
	c.Assert(imagingConfig.Config.Imaging.Avif.EncoderSpeed, qt.Equals, 1)
}

func TestImageConfigPerFormat(t *testing.T) {
	c := qt.New(t)

	cfg, err := DecodeConfig(map[string]any{
		"quality":     80,
		"compression": "lossy",
		"hint":        "text",
		"jpeg":        map[string]any{"quality": 80},
		"webp":        map[string]any{"quality": 70, "hint": "picture", "compression": "lossless"},
		"avif":        map[string]any{"quality": 55},
	})
	c.Assert(err, qt.IsNil)

	conf := func(opts ...string) ImageConfig {
		conf, err := DecodeImageConfig(append([]string{"resize", "100x"}, opts...), cfg, JPEG)
		c.Assert(err, qt.IsNil)
		return conf
	}

	c.Assert(conf("jpg").Quality, qt.Equals, 80)
	c.Assert(conf("avif").Quality, qt.Equals, 55)
	c.Assert(conf("webp").Quality, qt.Equals, 70)

	c.Assert(conf("webp", "q33").Quality, qt.Equals, 33)
	c.Assert(conf("avif", "q33").Quality, qt.Equals, 33)

	c.Assert(conf("webp").Hint, qt.Equals, "picture")
	c.Assert(conf("avif").Hint, qt.Equals, "text")
	c.Assert(conf("jpeg").Hint, qt.Equals, "")

	c.Assert(conf("webp").Compression, qt.Equals, "lossless")
	c.Assert(conf("avif").Compression, qt.Equals, "lossy")
	c.Assert(conf("jpeg").Compression, qt.Equals, "")
}

func TestDecodeImageConfig(t *testing.T) {
	c := qt.New(t)

	for i, this := range []struct {
		action string
		in     string
		expect any
	}{
		{"resize", "300x400", newTestImageConfig("resize", 300, 400, 75, 0, "box", "smart", "")},
		{"resize", "300x400 #fff", newTestImageConfig("resize", 300, 400, 75, 0, "box", "smart", "fff")},
		{"resize", "100x200 bottomRight", newTestImageConfig("resize", 100, 200, 75, 0, "box", "BottomRight", "")},
		{"resize", "10x20 topleft Lanczos", newTestImageConfig("resize", 10, 20, 75, 0, "Lanczos", "topleft", "")},
		{"resize", "linear left 10x r180", newTestImageConfig("resize", 10, 0, 75, 180, "linear", "left", "")},
		{"resize", "x20 riGht Cosine q95", newTestImageConfig("resize", 0, 20, 95, 0, "cosine", "right", "")},
		{"crop", "300x400", newTestImageConfig("crop", 300, 400, 75, 0, "box", "smart", "")},
		{"fill", "300x400", newTestImageConfig("fill", 300, 400, 75, 0, "box", "smart", "")},
		{"fit", "300x400", newTestImageConfig("fit", 300, 400, 75, 0, "box", "smart", "")},

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
		options := append([]string{this.action}, strings.Fields(this.in)...)
		result, err := DecodeImageConfig(options, cfg, WEBP)
		if b, ok := this.expect.(bool); ok && !b {
			if err == nil {
				t.Errorf("[%d] parseImageConfig didn't return an expected error", i)
			}
		} else {
			if err != nil {
				t.Fatalf("[%d] err: %s", i, err)
			}
			expect := this.expect.(ImageConfig)
			result.Key = ""

			c.Assert(fmt.Sprint(result), qt.Equals, fmt.Sprint(expect))

		}
	}
}

func newTestImageConfig(action string, width, height, quality, rotate int, filter, anchor, bgColor string) ImageConfig {
	var c ImageConfig = newImageConfig()
	c.Action = action
	c.TargetFormat = WEBP
	c.Hint = defaultHint
	c.Compression = defaultCompression
	c.Width = width
	c.Height = height
	c.Quality = quality
	c.Rotate = rotate
	c.BgColor, _ = hexStringToColorGo(bgColor)
	c.Anchor = SmartCropAnchor
	c.Method = defaultWebpMethod

	if filter != "" {
		filter = strings.ToLower(filter)
		if v, ok := imageFilters[filter]; ok {
			c.Filter = v
		}
	}

	if anchor != "" {
		anchor = strings.ToLower(anchor)
		if v, ok := anchorPositions[anchor]; ok {
			c.Anchor = v
		}
	}

	return c
}
