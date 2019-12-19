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
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"sync"

	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/resources/images/exif"

	"github.com/disintegration/gift"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"

	"github.com/gohugoio/hugo/common/hugio"
	"github.com/pkg/errors"
)

func NewImage(f Format, proc *ImageProcessor, img image.Image, s Spec) *Image {
	if img != nil {
		return &Image{
			Format: f,
			Proc:   proc,
			Spec:   s,
			imageConfig: &imageConfig{
				config:       imageConfigFromImage(img),
				configLoaded: true,
			},
		}
	}
	return &Image{Format: f, Proc: proc, Spec: s, imageConfig: &imageConfig{}}
}

type Image struct {
	Format Format
	Proc   *ImageProcessor
	Spec   Spec
	*imageConfig
}

func (i *Image) EncodeTo(conf ImageConfig, img image.Image, w io.Writer) error {
	switch conf.TargetFormat {
	case JPEG:

		var rgba *image.RGBA
		quality := conf.Quality

		if nrgba, ok := img.(*image.NRGBA); ok {
			if nrgba.Opaque() {
				rgba = &image.RGBA{
					Pix:    nrgba.Pix,
					Stride: nrgba.Stride,
					Rect:   nrgba.Rect,
				}
			}
		}
		if rgba != nil {
			return jpeg.Encode(w, rgba, &jpeg.Options{Quality: quality})
		}
		return jpeg.Encode(w, img, &jpeg.Options{Quality: quality})
	case PNG:
		encoder := png.Encoder{CompressionLevel: png.DefaultCompression}
		return encoder.Encode(w, img)

	case GIF:
		return gif.Encode(w, img, &gif.Options{
			NumColors: 256,
		})
	case TIFF:
		return tiff.Encode(w, img, &tiff.Options{Compression: tiff.Deflate, Predictor: true})

	case BMP:
		return bmp.Encode(w, img)
	default:
		return errors.New("format not supported")
	}

}

// Height returns i's height.
func (i *Image) Height() int {
	i.initConfig()
	return i.config.Height
}

// Width returns i's width.
func (i *Image) Width() int {
	i.initConfig()
	return i.config.Width
}

func (i Image) WithImage(img image.Image) *Image {
	i.Spec = nil
	i.imageConfig = &imageConfig{
		config:       imageConfigFromImage(img),
		configLoaded: true,
	}

	return &i
}

func (i Image) WithSpec(s Spec) *Image {
	i.Spec = s
	i.imageConfig = &imageConfig{}
	return &i
}

// InitConfig reads the image config from the given reader.
func (i *Image) InitConfig(r io.Reader) error {
	var err error
	i.configInit.Do(func() {
		i.config, _, err = image.DecodeConfig(r)
	})
	return err
}

func (i *Image) initConfig() error {
	var err error
	i.configInit.Do(func() {
		if i.configLoaded {
			return
		}

		var f hugio.ReadSeekCloser

		f, err = i.Spec.ReadSeekCloser()
		if err != nil {
			return
		}
		defer f.Close()

		i.config, _, err = image.DecodeConfig(f)
	})

	if err != nil {
		return errors.Wrap(err, "failed to load image config")
	}

	return nil
}

func NewImageProcessor(cfg ImagingConfig) (*ImageProcessor, error) {
	e := cfg.Cfg.Exif
	exifDecoder, err := exif.NewDecoder(
		exif.WithDateDisabled(e.DisableDate),
		exif.WithLatLongDisabled(e.DisableLatLong),
		exif.ExcludeFields(e.ExcludeFields),
		exif.IncludeFields(e.IncludeFields),
	)

	if err != nil {
		return nil, err
	}

	return &ImageProcessor{
		Cfg:         cfg,
		exifDecoder: exifDecoder,
	}, nil

}

type ImageProcessor struct {
	Cfg         ImagingConfig
	exifDecoder *exif.Decoder
}

func (p *ImageProcessor) DecodeExif(r io.Reader) (*exif.Exif, error) {
	return p.exifDecoder.Decode(r)
}

func (p *ImageProcessor) ApplyFiltersFromConfig(src image.Image, conf ImageConfig) (image.Image, error) {
	var filters []gift.Filter

	if conf.Rotate != 0 {
		// Apply any rotation before any resize.
		filters = append(filters, gift.Rotate(float32(conf.Rotate), color.Transparent, gift.NearestNeighborInterpolation))
	}

	switch conf.Action {
	case "resize":
		filters = append(filters, gift.Resize(conf.Width, conf.Height, conf.Filter))
	case "fill":
		if conf.AnchorStr == smartCropIdentifier {
			bounds, err := p.smartCrop(src, conf.Width, conf.Height, conf.Filter)
			if err != nil {
				return nil, err
			}

			// First crop it, then resize it.
			filters = append(filters, gift.Crop(bounds))
			filters = append(filters, gift.Resize(conf.Width, conf.Height, conf.Filter))

		} else {
			filters = append(filters, gift.ResizeToFill(conf.Width, conf.Height, conf.Filter, conf.Anchor))
		}
	case "fit":
		filters = append(filters, gift.ResizeToFit(conf.Width, conf.Height, conf.Filter))
	default:
		return nil, errors.Errorf("unsupported action: %q", conf.Action)
	}

	img, err := p.Filter(src, filters...)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func (p *ImageProcessor) Filter(src image.Image, filters ...gift.Filter) (image.Image, error) {
	g := gift.New(filters...)
	dst := image.NewRGBA(g.Bounds(src.Bounds()))
	g.Draw(dst, src)
	return dst, nil
}

func (p *ImageProcessor) GetDefaultImageConfig(action string) ImageConfig {
	return ImageConfig{
		Action:  action,
		Quality: p.Cfg.Cfg.Quality,
	}
}

type Spec interface {
	// Loads the image source.
	ReadSeekCloser() (hugio.ReadSeekCloser, error)
}

// Format is an image file format.
type Format int

const (
	JPEG Format = iota + 1
	PNG
	GIF
	TIFF
	BMP
)

// RequiresDefaultQuality returns if the default quality needs to be applied to images of this format
func (f Format) RequiresDefaultQuality() bool {
	return f == JPEG
}

// SupportsTransparency reports whether it supports transparency in any form.
func (f Format) SupportsTransparency() bool {
	return f != JPEG
}

// DefaultExtension returns the default file extension of this format, starting with a dot.
// For example: .jpg for JPEG
func (f Format) DefaultExtension() string {
	return f.MediaType().FullSuffix()
}

// MediaType returns the media type of this image, e.g. image/jpeg for JPEG
func (f Format) MediaType() media.Type {
	switch f {
	case JPEG:
		return media.JPEGType
	case PNG:
		return media.PNGType
	case GIF:
		return media.GIFType
	case TIFF:
		return media.TIFFType
	case BMP:
		return media.BMPType
	default:
		panic(fmt.Sprintf("%d is not a valid image format", f))
	}
}

type imageConfig struct {
	config       image.Config
	configInit   sync.Once
	configLoaded bool
}

func imageConfigFromImage(img image.Image) image.Config {
	b := img.Bounds()
	return image.Config{Width: b.Max.X, Height: b.Max.Y}
}

func ToFilters(in interface{}) []gift.Filter {
	switch v := in.(type) {
	case []gift.Filter:
		return v
	case []filter:
		vv := make([]gift.Filter, len(v))
		for i, f := range v {
			vv[i] = f
		}
		return vv
	case gift.Filter:
		return []gift.Filter{v}
	default:
		panic(fmt.Sprintf("%T is not an image filter", in))
	}
}

// IsOpaque returns false if the image has alpha channel and there is at least 1
// pixel that is not (fully) opaque.
func IsOpaque(img image.Image) bool {
	if oim, ok := img.(interface {
		Opaque() bool
	}); ok {
		return oim.Opaque()
	}

	return false
}
