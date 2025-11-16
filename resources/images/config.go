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
	"image/color"
	"strconv"
	"strings"

	"github.com/gohugoio/hugo/common/hashing"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/media"
	"github.com/mitchellh/mapstructure"

	"github.com/bep/gowebp/libwebp/webpoptions"

	"github.com/disintegration/gift"
)

const (
	ActionResize = "resize"
	ActionCrop   = "crop"
	ActionFit    = "fit"
	ActionFill   = "fill"
)

var Actions = map[string]bool{
	ActionResize: true,
	ActionCrop:   true,
	ActionFit:    true,
	ActionFill:   true,
}

var (
	imageFormats = map[string]Format{
		".jpg":  JPEG,
		".jpeg": JPEG,
		".jpe":  JPEG,
		".jif":  JPEG,
		".jfif": JPEG,
		".png":  PNG,
		".tif":  TIFF,
		".tiff": TIFF,
		".bmp":  BMP,
		".gif":  GIF,
		".webp": WEBP,
	}

	imageFormatsBySubType = map[string]Format{
		media.Builtin.JPEGType.SubType: JPEG,
		media.Builtin.PNGType.SubType:  PNG,
		media.Builtin.TIFFType.SubType: TIFF,
		media.Builtin.BMPType.SubType:  BMP,
		media.Builtin.GIFType.SubType:  GIF,
		media.Builtin.WEBPType.SubType: WEBP,
	}

	// Add or increment if changes to an image format's processing requires
	// re-generation.
	imageFormatsVersions = map[Format]int{
		PNG:  0,
		WEBP: 0,
		GIF:  0,
	}

	// Increment to mark all processed images as stale. Only use when absolutely needed.
	// See the finer grained smartCropVersionNumber and imageFormatsVersions.
	mainImageVersionNumber = 0
)

var anchorPositions = map[string]gift.Anchor{
	strings.ToLower("Center"):      gift.CenterAnchor,
	strings.ToLower("TopLeft"):     gift.TopLeftAnchor,
	strings.ToLower("Top"):         gift.TopAnchor,
	strings.ToLower("TopRight"):    gift.TopRightAnchor,
	strings.ToLower("Left"):        gift.LeftAnchor,
	strings.ToLower("Right"):       gift.RightAnchor,
	strings.ToLower("BottomLeft"):  gift.BottomLeftAnchor,
	strings.ToLower("Bottom"):      gift.BottomAnchor,
	strings.ToLower("BottomRight"): gift.BottomRightAnchor,
	smartCropIdentifier:            SmartCropAnchor,
}

// These encoding hints are currently only relevant for Webp.
var hints = map[string]webpoptions.EncodingPreset{
	"picture": webpoptions.EncodingPresetPicture,
	"photo":   webpoptions.EncodingPresetPhoto,
	"drawing": webpoptions.EncodingPresetDrawing,
	"icon":    webpoptions.EncodingPresetIcon,
	"text":    webpoptions.EncodingPresetText,
}

var imageFilters = map[string]gift.Resampling{
	strings.ToLower("NearestNeighbor"):   gift.NearestNeighborResampling,
	strings.ToLower("Box"):               gift.BoxResampling,
	strings.ToLower("Linear"):            gift.LinearResampling,
	strings.ToLower("Hermite"):           hermiteResampling,
	strings.ToLower("MitchellNetravali"): mitchellNetravaliResampling,
	strings.ToLower("CatmullRom"):        catmullRomResampling,
	strings.ToLower("BSpline"):           bSplineResampling,
	strings.ToLower("Gaussian"):          gaussianResampling,
	strings.ToLower("Lanczos"):           gift.LanczosResampling,
	strings.ToLower("Hann"):              hannResampling,
	strings.ToLower("Hamming"):           hammingResampling,
	strings.ToLower("Blackman"):          blackmanResampling,
	strings.ToLower("Bartlett"):          bartlettResampling,
	strings.ToLower("Welch"):             welchResampling,
	strings.ToLower("Cosine"):            cosineResampling,
}

func ImageFormatFromExt(ext string) (Format, bool) {
	f, found := imageFormats[ext]
	return f, found
}

func ImageFormatFromMediaSubType(sub string) (Format, bool) {
	f, found := imageFormatsBySubType[sub]
	return f, found
}

const (
	defaultJPEGQuality    = 75
	defaultResampleFilter = "box"
	defaultBgColor        = "#ffffff"
	defaultHint           = "photo"
)

var (
	defaultImaging = map[string]any{
		"resampleFilter": defaultResampleFilter,
		"bgColor":        defaultBgColor,
		"hint":           defaultHint,
		"quality":        defaultJPEGQuality,
	}

	defaultImageConfig *config.ConfigNamespace[ImagingConfig, ImagingConfigInternal]
)

func init() {
	var err error
	defaultImageConfig, err = DecodeConfig(defaultImaging)
	if err != nil {
		panic(err)
	}
}

func DecodeConfig(in map[string]any) (*config.ConfigNamespace[ImagingConfig, ImagingConfigInternal], error) {
	if in == nil {
		in = make(map[string]any)
	}

	buildConfig := func(in any) (ImagingConfigInternal, any, error) {
		m, err := maps.ToStringMapE(in)
		if err != nil {
			return ImagingConfigInternal{}, nil, err
		}
		// Merge in the defaults.
		maps.MergeShallow(m, defaultImaging)

		var i ImagingConfigInternal
		if err := mapstructure.Decode(m, &i.Imaging); err != nil {
			return i, nil, err
		}

		if err := i.Imaging.init(); err != nil {
			return i, nil, err
		}

		i.BgColor, err = hexStringToColorGo(i.Imaging.BgColor)
		if err != nil {
			return i, nil, err
		}

		if i.Imaging.Anchor != "" {
			anchor, found := anchorPositions[i.Imaging.Anchor]
			if !found {
				return i, nil, fmt.Errorf("invalid anchor value %q in imaging config", i.Anchor)
			}
			i.Anchor = anchor
		}

		filter, found := imageFilters[i.Imaging.ResampleFilter]
		if !found {
			return i, nil, fmt.Errorf("%q is not a valid resample filter", filter)
		}

		i.ResampleFilter = filter

		return i, nil, nil
	}

	ns, err := config.DecodeNamespace[ImagingConfig](in, buildConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to decode media types: %w", err)
	}
	return ns, nil
}

func DecodeImageConfig(options []string, defaults *config.ConfigNamespace[ImagingConfig, ImagingConfigInternal], sourceFormat Format) (ImageConfig, error) {
	var (
		c   ImageConfig = GetDefaultImageConfig(defaults)
		err error
	)

	// Make to lower case, trim space and remove any empty strings.
	n := 0
	for _, s := range options {
		s = strings.TrimSpace(s)
		if s != "" {
			options[n] = strings.ToLower(s)
			n++
		}
	}
	options = options[:n]

	for _, part := range options {
		if _, ok := Actions[part]; ok {
			c.Action = part
		} else if pos, ok := anchorPositions[part]; ok {
			c.Anchor = pos
		} else if filter, ok := imageFilters[part]; ok {
			c.Filter = filter
		} else if hint, ok := hints[part]; ok {
			c.Hint = hint
		} else if part[0] == '#' {
			c.BgColor, err = hexStringToColorGo(part[1:])
			if err != nil {
				return c, err
			}
		} else if part[0] == 'q' {
			c.Quality, err = strconv.Atoi(part[1:])
			if err != nil {
				return c, err
			}
			if c.Quality < 1 || c.Quality > 100 {
				return c, errors.New("quality ranges from 1 to 100 inclusive")
			}
			c.qualitySetForImage = true
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
		} else if f, ok := ImageFormatFromExt("." + part); ok {
			c.TargetFormat = f
		}
	}

	switch c.Action {
	case ActionCrop, ActionFill, ActionFit:
		if c.Width == 0 || c.Height == 0 {
			return c, errors.New("must provide Width and Height")
		}
	case ActionResize:
		if c.Width == 0 && c.Height == 0 {
			return c, errors.New("must provide Width or Height")
		}
	default:
		if c.Width != 0 || c.Height != 0 {
			return c, errors.New("width or height are not supported for this action")
		}
	}

	if c.Action != "" && c.Filter == nil {
		c.Filter = defaults.Config.ResampleFilter
	}

	if c.Hint == 0 {
		c.Hint = webpoptions.EncodingPresetPhoto
	}

	if c.Action != "" && c.Anchor == -1 {
		c.Anchor = defaults.Config.Anchor
	}

	// default to the source format
	if c.TargetFormat == 0 {
		c.TargetFormat = sourceFormat
	}

	if c.Quality <= 0 && c.TargetFormat.RequiresDefaultQuality() {
		// We need a quality setting for all JPEGs and WEBPs.
		c.Quality = defaults.Config.Imaging.Quality
	}

	if c.BgColor == nil && c.TargetFormat != sourceFormat {
		if sourceFormat.SupportsTransparency() && !c.TargetFormat.SupportsTransparency() {
			c.BgColor = defaults.Config.BgColor
		}
	}

	if mainImageVersionNumber > 0 {
		options = append(options, strconv.Itoa(mainImageVersionNumber))
	}

	if v, ok := imageFormatsVersions[sourceFormat]; ok && v > 0 {
		options = append(options, strconv.Itoa(v))
	}

	if smartCropVersionNumber > 0 && c.Anchor == SmartCropAnchor {
		options = append(options, strconv.Itoa(smartCropVersionNumber))
	}

	c.Key = hashing.HashStringHex(options)

	return c, nil
}

// ImageConfig holds configuration to create a new image from an existing one, resize etc.
type ImageConfig struct {
	// This defines the output format of the output image. It defaults to the source format.
	TargetFormat Format

	Action string

	// If set, this will be used as the key in filenames etc.
	Key string

	// Quality ranges from 1 to 100 inclusive, higher is better.
	// This is only relevant for JPEG and WEBP images.
	// Default is 75.
	Quality            int
	qualitySetForImage bool // Whether the above is set for this image.

	// Rotate rotates an image by the given angle counter-clockwise.
	// The rotation will be performed first.
	Rotate int

	// Used to fill any transparency.
	// When set in site config, it's used when converting to a format that does
	// not support transparency.
	// When set per image operation, it's used even for formats that does support
	// transparency.
	BgColor color.Color

	// Hint about what type of picture this is. Used to optimize encoding
	// when target is set to webp.
	Hint webpoptions.EncodingPreset

	Width  int
	Height int

	Filter gift.Resampling

	Anchor gift.Anchor
}

func (cfg ImageConfig) Reanchor(a gift.Anchor) ImageConfig {
	cfg.Anchor = a
	cfg.Key = hashing.HashStringHex(cfg.Key, "reanchor", a)
	return cfg
}

type ImagingConfigInternal struct {
	BgColor        color.Color
	Hint           webpoptions.EncodingPreset
	ResampleFilter gift.Resampling
	Anchor         gift.Anchor

	Imaging ImagingConfig
}

func (i *ImagingConfigInternal) Compile(externalCfg *ImagingConfig) error {
	var err error
	i.BgColor, err = hexStringToColorGo(externalCfg.BgColor)
	if err != nil {
		return err
	}

	if externalCfg.Anchor != "" {
		anchor, found := anchorPositions[externalCfg.Anchor]
		if !found {
			return fmt.Errorf("invalid anchor value %q in imaging config", i.Anchor)
		}
		i.Anchor = anchor
	}

	filter, found := imageFilters[externalCfg.ResampleFilter]
	if !found {
		return fmt.Errorf("%q is not a valid resample filter", filter)
	}
	i.ResampleFilter = filter

	return nil
}

// ImagingConfig contains default image processing configuration. This will be fetched
// from site (or language) config.
type ImagingConfig struct {
	// Default image quality setting (1-100). Only used for JPEG images.
	Quality int

	// Resample filter to use in resize operations.
	ResampleFilter string

	// Hint about what type of image this is.
	// Currently only used when encoding to Webp.
	// Default is "photo".
	// Valid values are "picture", "photo", "drawing", "icon", or "text".
	Hint string

	// The anchor to use in Fill. Default is "smart", i.e. Smart Crop.
	Anchor string

	// Default color used in fill operations (e.g. "fff" for white).
	BgColor string

	Exif ExifConfig
}

func (cfg *ImagingConfig) init() error {
	if cfg.Quality < 0 || cfg.Quality > 100 {
		return errors.New("image quality must be a number between 1 and 100")
	}

	cfg.BgColor = strings.ToLower(strings.TrimPrefix(cfg.BgColor, "#"))
	cfg.Anchor = strings.ToLower(cfg.Anchor)
	cfg.ResampleFilter = strings.ToLower(cfg.ResampleFilter)
	cfg.Hint = strings.ToLower(cfg.Hint)

	if cfg.Anchor == "" {
		cfg.Anchor = smartCropIdentifier
	}

	if strings.TrimSpace(cfg.Exif.IncludeFields) == "" && strings.TrimSpace(cfg.Exif.ExcludeFields) == "" {
		// Don't change this for no good reason. Please don't.
		cfg.Exif.ExcludeFields = "GPS|Exif|Exposure[M|P|B]|Contrast|Resolution|Sharp|JPEG|Metering|Sensing|Saturation|ColorSpace|Flash|WhiteBalance"
	}

	return nil
}

type ExifConfig struct {
	// Regexp matching the Exif fields you want from the (massive) set of Exif info
	// available. As we cache this info to disk, this is for performance and
	// disk space reasons more than anything.
	// If you want it all, put ".*" in this config setting.
	// Note that if neither this or ExcludeFields is set, Hugo will return a small
	// default set.
	IncludeFields string

	// Regexp matching the Exif fields you want to exclude. This may be easier to use
	// than IncludeFields above, depending on what you want.
	ExcludeFields string

	// Hugo extracts the "photo taken" date/time into .Date by default.
	// Set this to true to turn it off.
	DisableDate bool

	// Hugo extracts the "photo taken where" (GPS latitude and longitude) into
	// .Long and .Lat. Set this to true to turn it off.
	DisableLatLong bool
}
