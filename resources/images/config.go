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
	"maps"
	"slices"
	"strconv"
	"strings"

	"github.com/gohugoio/hugo/common/hashing"
	"github.com/gohugoio/hugo/common/hmaps"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/media"
	"github.com/mitchellh/mapstructure"

	"github.com/gohugoio/gift"
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
		".avif": AVIF,
		".heif": HEIF,
		".heic": HEIC,
	}

	// These are the image types we can process.
	processableImageSubTypes = map[string]Format{
		media.Builtin.JPEGType.SubType: JPEG,
		media.Builtin.PNGType.SubType:  PNG,
		media.Builtin.TIFFType.SubType: TIFF,
		media.Builtin.BMPType.SubType:  BMP,
		media.Builtin.GIFType.SubType:  GIF,
		media.Builtin.WEBPType.SubType: WEBP,
		media.Builtin.AVIFType.SubType: AVIF,
	}

	// We cannot process these formats, but we can provide metadata support for them (including width/height).
	metaOnlyImageSubTypes = map[string]Format{
		media.Builtin.AVIFType.SubType: AVIF,
		media.Builtin.HEIFType.SubType: HEIF,
		media.Builtin.HEICType.SubType: HEIC,
	}

	// Increment to mark all processed images as stale. Only use when absolutely needed.
	// See the finer grained smartCropVersionNumber.
	mainImageVersionNumber = 1

	// Increment a format's version number to mark all processed images targeting
	// that format as stale, e.g. after a change to its encoder. This is finer
	// grained than mainImageVersionNumber, which invalidates every format.
	formatVersionNumbers = map[Format]int{
		AVIF: 1,
	}
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

var compressionMethods = map[string]bool{
	"lossy":    true,
	"lossless": true,
}

// These encoding hints are used by Webp (preset) and Avif (chroma subsampling).
var hints = map[string]bool{
	"picture": true,
	"photo":   true,
	"drawing": true,
	"icon":    true,
	"text":    true,
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

type ImageResourceType int

const (
	// ImageResourceTypeNone means that the resource is not an image, and thus does not support any image operations.
	ImageResourceTypeNone ImageResourceType = iota
	// This is an image, but with no support for any image operations.
	ImageResourceTypeBasic
	// ImageResourceTypeMetaOnly means that only metadata operations (e.g. getting width/height and other metadata) are supported for this format.
	ImageResourceTypeMetaOnly
	// ImageResourceTypeProcessable means that all image operations (resizing, cropping, etc.) are supported for this format.
	ImageResourceTypeProcessable
)

// ImageFormatFromMediaSubType returns the image format for the given media subtype, and how much image processing operations are supported for this format.
func ImageFormatFromMediaSubType(sub string) (Format, ImageResourceType) {
	f, found := processableImageSubTypes[sub]
	if found {
		return f, ImageResourceTypeProcessable
	}
	if f, found = metaOnlyImageSubTypes[sub]; found {
		return f, ImageResourceTypeMetaOnly
	}
	return f, ImageResourceTypeBasic
}

const (
	defaultResampleFilter   = "box"
	defaultBgColor          = "#ffffff"
	defaultHint             = "photo"
	defaultCompression      = "lossy"
	defaultWebpUseSharpYuv  = false
	defaultWebpMethod       = 2
	defaultAvifEncoderSpeed = 10
)

func DecodeConfig(in map[string]any) (*config.ConfigNamespace[ImagingConfig, ImagingConfigInternal], error) {
	if in == nil {
		in = make(map[string]any)
	}

	buildConfig := func(in any) (ImagingConfigInternal, any, error) {
		m, err := hmaps.ToStringMapE(in)
		if err != nil {
			return ImagingConfigInternal{}, nil, err
		}

		i := ImagingConfigInternal{
			Imaging: ImagingConfig{
				BgColor:        defaultBgColor,
				ResampleFilter: defaultResampleFilter,
				Webp: WebpConfig{
					UseSharpYuv: defaultWebpUseSharpYuv,
					Method:      defaultWebpMethod,
				},
				Avif: AvifConfig{
					EncoderSpeed: defaultAvifEncoderSpeed,
				},
			},
		}

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
				return i, nil, fmt.Errorf("invalid anchor value %q in imaging config", i.Imaging.Anchor)
			}
			i.Anchor = anchor
		}

		filter, found := imageFilters[i.Imaging.ResampleFilter]
		if !found {
			return i, nil, fmt.Errorf("%q is not a valid resample filter", filter)
		}

		i.ResampleFilter = filter

		return i, i.Imaging, nil
	}

	ns, err := config.DecodeNamespace[ImagingConfig](in, buildConfig)
	if err != nil {
		return nil, err
	}
	return ns, nil
}

func DecodeImageConfig(options []string, defaults *config.ConfigNamespace[ImagingConfig, ImagingConfigInternal], sourceFormat Format) (ImageConfig, error) {
	var (
		c   ImageConfig = newImageConfig()
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
		} else if _, ok := hints[part]; ok {
			c.Hint = part
		} else if _, ok := compressionMethods[part]; ok {
			c.Compression = part
		} else if f, ok := ImageFormatFromExt("." + part); ok {
			c.TargetFormat = f
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

	if err := c.init(defaults.Config, sourceFormat); err != nil {
		return c, err
	}

	if mainImageVersionNumber > 0 {
		options = append(options, strconv.Itoa(mainImageVersionNumber))
	}

	if v := formatVersionNumbers[c.TargetFormat]; v > 0 {
		options = append(options, "tfv"+strconv.Itoa(v))
	}

	usesSmartCrop := c.Anchor == SmartCropAnchor && (c.Action == ActionCrop || c.Action == ActionFill)
	if smartCropVersionNumber > 0 && usesSmartCrop {
		options = append(options, strconv.Itoa(smartCropVersionNumber))
	}

	c.Key = hashing.HashStringHex(options)

	return c, nil
}

// ImageConfig holds configuration to create a new image from an existing one, resize etc.
type ImageConfig struct {
	// This defines the output format of the output image. It defaults to the source format.
	TargetFormat Format

	// Optional image processing action to perform ("resize", "crop", "fit", "fill"). When empty, the operation may still perform e.g. format conversion.
	Action string

	// If set, this will be used as the key in filenames etc.
	Key string

	// Quality ranges from 1 to 100 inclusive, higher is better.
	// This is only relevant for JPEG and WEBP images.
	// For WEBP it's only relevant for lossy encoding.
	// Default is 75.
	Quality int

	// Rotate rotates an image by the given angle counter-clockwise.
	// The rotation will be performed first.
	Rotate int

	// Used to fill any transparency.
	// When set in project config, it's used when converting to a format that does
	// not support transparency.
	// When set per image operation, it's used even for formats that does support
	// transparency.
	BgColor color.Color

	Filter gift.Resampling
	Anchor gift.Anchor

	// Hint about what type of picture this is. Used to optimize encoding
	// when target is webp (preset) or avif (chroma subsampling).
	Hint string

	Compression string

	// WebP-specific options.
	UseSharpYuv bool
	Method      int

	// AVIF-specific options.
	EncoderSpeed int

	Width  int
	Height int
}

func (c *ImageConfig) init(defaults ImagingConfigInternal, sourceFormat Format) error {
	if c.Action != "" && c.Anchor == -1 {
		c.Anchor = defaults.Anchor
	}
	if c.TargetFormat == 0 {
		c.TargetFormat = sourceFormat
	}
	if c.Filter == nil {
		c.Filter = defaults.ResampleFilter
	}
	if c.Anchor == -1 {
		c.Anchor = defaults.Anchor
	}
	if c.Hint == "" {
		c.Hint = defaults.Imaging.hintFor(c.TargetFormat)
	}

	if c.Quality == 0 {
		c.Quality = defaults.Imaging.qualityFor(c.TargetFormat)
	}

	if c.Compression == "" {
		c.Compression = defaults.Imaging.compressionFor(c.TargetFormat)
	}

	if c.TargetFormat == AVIF {
		c.EncoderSpeed = defaults.Imaging.Avif.EncoderSpeed
	}

	if c.TargetFormat == WEBP {
		c.Method = defaults.Imaging.Webp.Method
		c.UseSharpYuv = defaults.Imaging.Webp.UseSharpYuv
	}

	if c.BgColor == nil && c.TargetFormat != sourceFormat {
		if sourceFormat.SupportsTransparency() && !c.TargetFormat.SupportsTransparency() {
			c.BgColor = defaults.BgColor
		}
	}
	return nil
}

func (c ImageConfig) Reanchor(a gift.Anchor) ImageConfig {
	c.Anchor = a
	c.Key = hashing.HashStringHex(c.Key, "reanchor", a)
	return c
}

type ImagingConfigInternal struct {
	BgColor        color.Color
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
			return fmt.Errorf("invalid anchor value %q in imaging config", externalCfg.Anchor)
		}
		i.Anchor = anchor
	}

	filter, found := imageFilters[externalCfg.ResampleFilter]
	if !found {
		return fmt.Errorf("%q is not a valid resample filter", externalCfg.ResampleFilter)
	}
	i.ResampleFilter = filter

	return nil
}

// ImagingConfig contains default image processing configuration. This will be fetched
// from site (or language) config.
type ImagingConfig struct {
	// Default image quality setting (1-100). Used as the fallback for JPEG,
	// WebP and AVIF when no per-format quality is set. When left unset, JPEG
	// and WebP default to 75 and AVIF to 60 (its scale differs perceptually).
	// Deprecated in v0.163.0: set the quality per format instead, see
	// imaging.jpeg.quality, imaging.webp.quality and imaging.avif.quality.
	Quality int `json:"-"`

	// Compression method to use.
	// One of "lossy" or "lossless".
	// Note that lossless is currently only supported for WebP and AVIF.
	// Deprecated in v0.163.0: set the compression method per format instead, see imaging.webp.compression and imaging.avif.compression.
	Compression string `json:"-"`

	// Resample filter to use in resize operations.
	ResampleFilter string

	// Hint about what type of image this is.
	// Used when encoding to Webp (preset) and Avif (chroma subsampling).
	// Default is "photo".
	// Valid values are "picture", "photo", "drawing", "icon", or "text".
	// Moved to WebpConfig in v0.155.0, but kept here for backwards compatibility.
	Hint string `json:"-"`

	// The anchor to use in Fill. Default is "smart", i.e. Smart Crop.
	Anchor string

	// Default color used in fill operations (e.g. "fff" for white).
	BgColor string

	Exif ExifConfig
	Meta MetaConfig
	Jpeg JpegConfig
	Webp WebpConfig
	Avif AvifConfig
}

func (cfg *ImagingConfig) qualityFor(f Format) int {
	switch f {
	case JPEG:
		return cfg.Jpeg.Quality
	case WEBP:
		return cfg.Webp.Quality
	case AVIF:
		return cfg.Avif.Quality
	}
	return 0
}

func (cfg *ImagingConfig) compressionFor(f Format) string {
	switch f {
	case WEBP:
		return cfg.Webp.Compression
	case AVIF:
		return cfg.Avif.Compression
	}
	return ""
}

// hintFor returns the configured hint for the given target format.
func (cfg *ImagingConfig) hintFor(f Format) string {
	switch f {
	case WEBP:
		return cfg.Webp.Hint
	case AVIF:
		return cfg.Avif.Hint
	}
	return ""
}

var validMetaSources = map[string]bool{
	"exif": true,
	"iptc": true,
	"xmp":  true,
}

func (cfg *ImagingConfig) init() error {
	cfg.BgColor = strings.ToLower(strings.TrimPrefix(cfg.BgColor, "#"))
	cfg.Anchor = strings.ToLower(cfg.Anchor)
	cfg.ResampleFilter = strings.ToLower(cfg.ResampleFilter)
	cfg.Hint = strings.ToLower(cfg.Hint)
	cfg.Compression = strings.ToLower(cfg.Compression)
	if err := cfg.Jpeg.init(cfg); err != nil {
		return fmt.Errorf("invalid jpeg config: %w", err)
	}
	if err := cfg.Webp.init(cfg); err != nil {
		return fmt.Errorf("invalid webp config: %w", err)
	}
	if err := cfg.Avif.init(cfg); err != nil {
		return fmt.Errorf("invalid avif config: %w", err)
	}
	if cfg.Quality < 0 || cfg.Quality > 100 {
		return fmt.Errorf("imaging.quality must be between 1 and 100 inclusive, got %d", cfg.Quality)
	}

	if cfg.Anchor == "" {
		cfg.Anchor = smartCropIdentifier
	}

	if strings.TrimSpace(cfg.Exif.IncludeFields) == "" && strings.TrimSpace(cfg.Exif.ExcludeFields) == "" {
		// Don't change this for no good reason. Please don't.
		cfg.Exif.ExcludeFields = "GPS|Exif|Exposure[M|P|B]|Contrast|Resolution|Sharp|JPEG|Metering|Sensing|Saturation|ColorSpace|Flash|WhiteBalance"
	}

	if len(cfg.Meta.Fields) == 0 {
		// Default: include all fields except technical metadata.
		// Don't change this for no good reason. Please don't.
		cfg.Meta.Fields = []string{
			"! *{GPS,Exif,Exposure[MPB],Contrast,Resolution,Sharp,JPEG,Metering,Sensing,Saturation,ColorSpace,Flash,WhiteBalance}*",
		}
	}

	if len(cfg.Meta.Sources) == 0 {
		// Default to EXIF and IPTC (XMP is slower to decode).
		cfg.Meta.Sources = []string{"exif", "iptc"}
	} else {
		// Normalize to lowercase.
		for i, s := range cfg.Meta.Sources {
			cfg.Meta.Sources[i] = strings.ToLower(s)
			if !validMetaSources[cfg.Meta.Sources[i]] {
				keys := slices.Collect(maps.Keys(validMetaSources))
				slices.Sort(keys)
				return fmt.Errorf("invalid metadata source %q in imaging.meta.sources config; must be one of %s", s, keys)
			}
		}
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

type MetaConfig struct {
	// Glob patterns for which metadata fields to include.
	// Use "! " prefix to exclude patterns (e.g., "! *GPS*" excludes GPS fields).
	// Patterns are OR'd together for inclusion, AND'd for exclusion.
	// If empty, a default set excluding technical metadata is used.
	// Use ["**"] to include all fields.
	Fields []string

	// Which metadata sources to include.
	// Valid values are "exif", "iptc", "xmp".
	// Default is ["exif", "iptc"] (XMP is excluded for performance reasons).
	Sources []string
}

// JpegConfig holds JPEG-specific encoding configuration.
type JpegConfig struct {
	// Quality setting (1-100). Falls back to the global imaging.quality if unset.
	Quality int
}

func (c *JpegConfig) init(ic *ImagingConfig) error {
	if c.Quality == 0 {
		c.Quality = ic.Quality
	}
	if c.Quality == 0 {
		c.Quality = 75
	}
	if c.Quality < 1 || c.Quality > 100 {
		return fmt.Errorf("imaging.jpeg.quality must be between 1 and 100 inclusive, got %d", c.Quality)
	}
	return nil
}

// AvifConfig holds AVIF-specific encoding configuration.
type AvifConfig struct {
	// Quality setting (1-100). Falls back to the global imaging.quality if unset.
	Quality int

	// Compression method to use.
	// One of "lossy" or "lossless".
	Compression string

	// Hint about what type of image this is. Used for chroma subsampling.
	// Valid values are "picture", "photo", "drawing", "icon", or "text".
	// Default is "photo".
	Hint string

	// Encoder quality/speed trade-off, 1 (slowest, best quality / smallest
	// files) to 10 (fastest). Default is 10 — fast enough for incremental
	// builds with quality indistinguishable from slower settings at typical
	// web thumbnail sizes. Lower values reduce file size at the cost of
	// build time.
	// We recommend sticking with the default of 10 unless you have a specific reason to change it,
	// and to stay above 5 to avoid very long build times and timeouts.
	EncoderSpeed int
}

func (c *AvifConfig) init(ic *ImagingConfig) error {
	if c.Hint == "" {
		c.Hint = ic.Hint
	}
	if c.Compression == "" {
		c.Compression = ic.Compression
	}
	if c.Quality == 0 {
		c.Quality = ic.Quality
	}
	if c.Hint == "" {
		c.Hint = defaultHint
	}
	if c.Compression == "" {
		c.Compression = defaultCompression
	}
	if c.Quality == 0 {
		c.Quality = 60
	}

	c.Hint = strings.ToLower(c.Hint)
	c.Compression = strings.ToLower(c.Compression)

	if c.Hint != "" && !hints[c.Hint] {
		return fmt.Errorf("imaging.avif.hint must be one of picture, photo, drawing, icon, or text, got %q", c.Hint)
	}
	if c.EncoderSpeed < 1 || c.EncoderSpeed > 10 {
		return fmt.Errorf("imaging.avif.encoderSpeed must be between 1 and 10, got %d", c.EncoderSpeed)
	}

	if c.Compression != "" && !compressionMethods[c.Compression] {
		return fmt.Errorf("imaging.avif.compression must be one of lossy or lossless, got %q", c.Compression)
	}
	if c.Quality < 1 || c.Quality > 100 {
		return fmt.Errorf("imaging.avif.quality must be between 1 and 100 inclusive, got %d", c.Quality)
	}

	return nil
}

// WebpConfig holds WebP-specific encoding configuration.
type WebpConfig struct {
	// Quality setting (1-100). Falls back to the global imaging.quality if unset.
	// Only relevant for lossy encoding.
	Quality int

	// Compression method to use.
	// One of "lossy" or "lossless".
	Compression string

	// Hint about what type of image this is.
	// Valid values are "picture", "photo", "drawing", "icon", or "text".
	// Default is "photo".
	Hint string

	// Use sharp (and slow) RGB->YUV conversion.
	// Default is true.
	UseSharpYuv bool

	// Quality/speed trade-off (0=fast, 6=slower-better).
	// Default is 2.
	Method int
}

func (c *WebpConfig) init(ic *ImagingConfig) error {
	if c.Hint == "" {
		c.Hint = ic.Hint
	}
	if c.Compression == "" {
		c.Compression = ic.Compression
	}
	if c.Quality == 0 {
		c.Quality = ic.Quality
	}
	if c.Hint == "" {
		c.Hint = defaultHint
	}
	if c.Compression == "" {
		c.Compression = defaultCompression
	}
	if c.Quality == 0 {
		c.Quality = 75
	}

	c.Hint = strings.ToLower(c.Hint)
	c.Compression = strings.ToLower(c.Compression)

	if c.Hint != "" && !hints[c.Hint] {
		return fmt.Errorf("imaging.webp.hint must be one of picture, photo, drawing, icon, or text, got %q", c.Hint)
	}
	if c.Method < 0 || c.Method > 6 {
		return fmt.Errorf("imaging.webp.method must be between 0 and 6, got %d", c.Method)
	}
	if c.Compression != "" && !compressionMethods[c.Compression] {
		return fmt.Errorf("imaging.webp.compression must be one of lossy or lossless, got %q", c.Compression)
	}
	if c.Quality < 1 || c.Quality > 100 {
		return fmt.Errorf("imaging.webp.quality must be between 1 and 100 inclusive, got %d", c.Quality)
	}
	return nil
}
