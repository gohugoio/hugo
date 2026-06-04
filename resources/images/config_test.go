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

	"github.com/gohugoio/hugo/config"

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

	// AVIF: minimum is 1; 0 is treated as unset and falls back to the default.
	imagingConfig, err = DecodeConfig(map[string]any{
		"avif": map[string]any{"encoderSpeed": 0},
	})
	c.Assert(err, qt.IsNil)
	c.Assert(imagingConfig.Config.Imaging.Avif.EncoderSpeed, qt.Equals, defaultAvifEncoderSpeed)

	imagingConfig, err = DecodeConfig(map[string]any{
		"avif": map[string]any{"encoderSpeed": 1},
	})
	c.Assert(err, qt.IsNil)
	c.Assert(imagingConfig.Config.Imaging.Avif.EncoderSpeed, qt.Equals, 1)
}

// See issue 14992.
func TestImageConfigHintPerFormat(t *testing.T) {
	c := qt.New(t)

	// Default: both WebP and AVIF hint default to "photo".
	cfg, err := DecodeConfig(map[string]any{})
	c.Assert(err, qt.IsNil)
	c.Assert(cfg.Config.Imaging.Webp.Hint, qt.Equals, "photo")
	c.Assert(cfg.Config.Imaging.Avif.Hint, qt.Equals, "photo")

	// Per-format hint from config.
	cfg, err = DecodeConfig(map[string]any{
		"webp": map[string]any{"hint": "drawing"},
		"avif": map[string]any{"hint": "icon"},
	})
	c.Assert(err, qt.IsNil)
	c.Assert(cfg.Config.Imaging.Webp.Hint, qt.Equals, "drawing")
	c.Assert(cfg.Config.Imaging.Avif.Hint, qt.Equals, "icon")

	// Root-level hint applies to both formats for backwards compatibility.
	cfg, err = DecodeConfig(map[string]any{"hint": "text"})
	c.Assert(err, qt.IsNil)
	c.Assert(cfg.Config.Imaging.Webp.Hint, qt.Equals, "text")
	c.Assert(cfg.Config.Imaging.Avif.Hint, qt.Equals, "text")

	// Invalid AVIF hint is rejected.
	_, err = DecodeConfig(map[string]any{"avif": map[string]any{"hint": "nope"}})
	c.Assert(err, qt.ErrorMatches, ".*invalid avif hint.*")

	// Per-image config picks the hint for the target format.
	cfg, err = DecodeConfig(map[string]any{
		"webp": map[string]any{"hint": "drawing"},
		"avif": map[string]any{"hint": "icon"},
	})
	c.Assert(err, qt.IsNil)
	hint := func(f Format, opts ...string) string {
		conf, err := DecodeImageConfig(append([]string{"resize", "100x"}, opts...), cfg, f)
		c.Assert(err, qt.IsNil)
		return conf.Hint
	}
	c.Assert(hint(WEBP), qt.Equals, "drawing")
	c.Assert(hint(AVIF), qt.Equals, "icon")
	// A per-image hint always wins.
	c.Assert(hint(AVIF, "photo"), qt.Equals, "photo")
}

// See issue 14957.
func TestImageConfigQualityPerFormat(t *testing.T) {
	c := qt.New(t)

	cfg, err := DecodeConfig(map[string]any{
		"quality": 80,
		"webp":    map[string]any{"quality": 70},
		"avif":    map[string]any{"quality": 55},
	})
	c.Assert(err, qt.IsNil)

	quality := func(opts ...string) int {
		conf, err := DecodeImageConfig(append([]string{"resize", "100x"}, opts...), cfg, JPEG)
		c.Assert(err, qt.IsNil)
		return conf.Quality
	}

	// Per-format quality from config.
	c.Assert(quality("webp"), qt.Equals, 70)
	c.Assert(quality("avif"), qt.Equals, 55)
	// JPEG has no per-format override, so it falls back to the global quality.
	c.Assert(quality("jpg"), qt.Equals, 80)

	// A per-image quality always wins.
	c.Assert(quality("webp", "q33"), qt.Equals, 33)
	c.Assert(quality("avif", "q33"), qt.Equals, 33)

	// Without per-format config, all formats fall back to the global quality.
	cfg2, err := DecodeConfig(map[string]any{"quality": 80})
	c.Assert(err, qt.IsNil)
	for _, f := range []string{"jpg", "webp", "avif"} {
		conf, err := DecodeImageConfig([]string{"resize", "100x", f}, cfg2, JPEG)
		c.Assert(err, qt.IsNil)
		c.Assert(conf.Quality, qt.Equals, 80)
	}

	// Out-of-range per-format quality is rejected.
	_, err = DecodeConfig(map[string]any{"avif": map[string]any{"quality": 123}})
	c.Assert(err, qt.ErrorMatches, ".*quality must be.*")
}

// See issue 14979.
func TestImageConfigAvifDefaultQuality(t *testing.T) {
	c := qt.New(t)

	quality := func(cfg *config.ConfigNamespace[ImagingConfig, ImagingConfigInternal], opts ...string) int {
		conf, err := DecodeImageConfig(append([]string{"resize", "100x"}, opts...), cfg, JPEG)
		c.Assert(err, qt.IsNil)
		return conf.Quality
	}

	// With nothing configured, AVIF defaults to 60 while JPEG/WebP stay at 75.
	cfg, err := DecodeConfig(map[string]any{})
	c.Assert(err, qt.IsNil)
	c.Assert(quality(cfg, "avif"), qt.Equals, 60)
	c.Assert(quality(cfg, "jpg"), qt.Equals, 75)
	c.Assert(quality(cfg, "webp"), qt.Equals, 75)

	// An explicit global quality applies to AVIF too.
	cfg, err = DecodeConfig(map[string]any{"quality": 80})
	c.Assert(err, qt.IsNil)
	c.Assert(quality(cfg, "avif"), qt.Equals, 80)

	// An explicit per-format AVIF quality wins over the default.
	cfg, err = DecodeConfig(map[string]any{"avif": map[string]any{"quality": 90}})
	c.Assert(err, qt.IsNil)
	c.Assert(quality(cfg, "avif"), qt.Equals, 90)

	// A per-image quality still wins.
	c.Assert(quality(cfg, "avif", "q33"), qt.Equals, 33)
}

// See issue 14990.
func TestImageConfigFormatVersionNumber(t *testing.T) {
	c := qt.New(t)

	cfg, err := DecodeConfig(map[string]any{})
	c.Assert(err, qt.IsNil)

	key := func(f string) string {
		conf, err := DecodeImageConfig([]string{"resize", "100x", f}, cfg, JPEG)
		c.Assert(err, qt.IsNil)
		return conf.Key
	}

	avifBefore, jpgBefore := key("avif"), key("jpg")

	orig := formatVersionNumbers[AVIF]
	defer func() { formatVersionNumbers[AVIF] = orig }()
	formatVersionNumbers[AVIF] = orig + 1

	// Bumping the AVIF version invalidates AVIF images only.
	c.Assert(key("avif"), qt.Not(qt.Equals), avifBefore)
	c.Assert(key("jpg"), qt.Equals, jpgBefore)
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
		options := append([]string{this.action}, strings.Fields(this.in)...)
		result, err := DecodeImageConfig(options, cfg, PNG)
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

			if fmt.Sprint(result) != fmt.Sprint(expect) {
				t.Fatalf("[%d] got\n%v\n but expected\n%v", i, result, expect)
			}
		}
	}
}

func newImageConfig(action string, width, height, quality, rotate int, filter, anchor, bgColor string) ImageConfig {
	var c ImageConfig = GetDefaultImageConfig(nil)
	c.Action = action
	c.TargetFormat = PNG
	c.Hint = defaultHint // Resolved per target format in DecodeImageConfig.
	c.Width = width
	c.Height = height
	c.Quality = quality
	c.Rotate = rotate
	c.BgColor, _ = hexStringToColorGo(bgColor)
	c.Anchor = SmartCropAnchor

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
