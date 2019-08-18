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
	"image"
	"image/jpeg"
	"io"
	"sync"

	"github.com/disintegration/imaging"
	"github.com/gohugoio/hugo/common/hugio"
	"github.com/pkg/errors"
)

func NewImage(f imaging.Format, proc *ImageProcessor, img image.Image, s Spec) *Image {
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
	Format imaging.Format

	Proc *ImageProcessor

	Spec Spec

	*imageConfig
}

func (i *Image) EncodeTo(conf ImageConfig, img image.Image, w io.Writer) error {
	switch i.Format {
	case imaging.JPEG:

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
	default:
		return imaging.Encode(w, img, i.Format)
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

func (i *Image) initConfig() error {
	var err error
	i.configInit.Do(func() {
		if i.configLoaded {
			return
		}

		var (
			f      hugio.ReadSeekCloser
			config image.Config
		)

		f, err = i.Spec.ReadSeekCloser()
		if err != nil {
			return
		}
		defer f.Close()

		config, _, err = image.DecodeConfig(f)
		if err != nil {
			return
		}
		i.config = config
	})

	if err != nil {
		return errors.Wrap(err, "failed to load image config")
	}

	return nil
}

type ImageProcessor struct {
	Cfg Imaging
}

func (p *ImageProcessor) Fill(src image.Image, conf ImageConfig) (image.Image, error) {
	if conf.AnchorStr == SmartCropIdentifier {
		return smartCrop(src, conf.Width, conf.Height, conf.Anchor, conf.Filter)
	}
	return imaging.Fill(src, conf.Width, conf.Height, conf.Anchor, conf.Filter), nil
}

func (p *ImageProcessor) Fit(src image.Image, conf ImageConfig) (image.Image, error) {
	return imaging.Fit(src, conf.Width, conf.Height, conf.Filter), nil
}

func (p *ImageProcessor) Resize(src image.Image, conf ImageConfig) (image.Image, error) {
	return imaging.Resize(src, conf.Width, conf.Height, conf.Filter), nil
}

type Spec interface {
	// Loads the image source.
	ReadSeekCloser() (hugio.ReadSeekCloser, error)
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
