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
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/masyagin1998/segmentator/segmentator"

	"github.com/dennwc/gotrace"
	"github.com/mitchellh/mapstructure"

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

// TODO1
var _ = &gotrace.Params{
	TurdSize:     2,
	TurnPolicy:   gotrace.TurnMinority,
	AlphaMax:     1.0,
	OptiCurve:    true,
	OptTolerance: 0.2,
}

var DefaultTraceOptions = TraceOptions{
	Color:  "#fff",
	Filter: "luma", // Threshold function; luma, sobel
	Low:    40,     // Luma range between 0-100.
	High:   99,
	TraceParams: gotrace.Params{
		TurdSize:     10, // Suppress speckles of up to this size
		TurnPolicy:   gotrace.TurnMinority,
		AlphaMax:     1.0,  //  Corner threshold parameter
		OptiCurve:    true, // Curve optimization
		OptTolerance: 0.2,  // Curve optimization tolerance
	},
}

type TraceOptions struct {
	Color       string
	Filter      string
	Low         int
	High        int
	TraceParams gotrace.Params `mapstructure:",squash"`
}

func (p *ImageProcessor) DecodeTraceOptions(m map[string]interface{}) (TraceOptions, error) {
	opts := DefaultTraceOptions
	if m == nil {
		return opts, nil
	}

	err := mapstructure.WeakDecode(m, &opts)

	opts.Filter = strings.ToLower(opts.Filter)

	// Luminance range is between 0-100.
	if opts.Low < 0 {
		opts.Low = 0
	}
	if opts.Low > 100 {
		opts.Low = 100
	}
	if opts.High < opts.Low {
		opts.High = opts.Low
	}
	if opts.High > 100 {
		opts.High = 100
	}

	return opts, err

}
func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

type segmentFunc struct {
	filter func(segments segmentator.Image) error
	value  func(p segmentator.Pixel) int
}

func (p *ImageProcessor) Trace(src image.Image, opts TraceOptions) (string, error) {
	defer timeTrack(time.Now(), "Trace")

	segments := createSegmentMap(src)

	pickFirst := func(p segmentator.Pixel) int {
		v := p.R
		return (v * 100) / 255
	}

	pickAll := func(p segmentator.Pixel) int {
		v := p.R + p.G + p.B
		return (v * 100) / 255
	}

	// Pixel Classification/threshold functions.
	tm := map[string]func(segments segmentator.Image) int{
		"otsu": func(segments segmentator.Image) int {
			return segmentator.FGPCOtsuThresholding2(segments)
		},
	}

	// Edge detection functions.
	sm := map[string]segmentFunc{
		"luma": segmentFunc{
			filter: func(segments segmentator.Image) error {
				segmentator.GSLuma(segments)
				return nil
			},
			value: pickAll,
		},
		"average": segmentFunc{
			filter: func(segments segmentator.Image) error {
				segmentator.GSAveraging(segments)
				return nil
			},
			value: pickAll,
		},
		"sobel": segmentFunc{ // Remove
			filter: func(segments segmentator.Image) error {
				return segmentator.FGEDSobel(segments, segmentator.GXGY)
			},
			value: pickFirst,
		},
		"previtt": segmentFunc{
			filter: func(segments segmentator.Image) error {
				return segmentator.FGEDPrevitt(segments, segmentator.GXGY)
			},
			value: pickFirst,
		},
		"roberts": segmentFunc{ // Keep
			filter: func(segments segmentator.Image) error {
				return segmentator.FGEDRoberts(segments, segmentator.GXGY)
			},
			value: pickFirst,
		},
		"scharr": segmentFunc{ // Remove
			filter: func(segments segmentator.Image) error {
				return segmentator.FGEDScharr(segments, segmentator.GXGY)
			},
			value: pickFirst,
		},
		"laplacian": segmentFunc{ // Remove
			filter: func(segments segmentator.Image) error {
				return segmentator.FGEDLaplacian4(segments)
			},
			value: pickFirst,
		},
	}

	var thresholdFunc func(x, y int, c color.Color) bool

	thresholdf, found := tm[opts.Filter]
	if found {
		threshold := thresholdf(segments)
		thresholdFunc = func(x, y int, c color.Color) bool {
			p := segments.Pixels[y][x]
			v := p.R
			return v >= threshold
		}

	} else {

		segmentf, found := sm[opts.Filter]
		if !found {
			segmentf = sm["luma"]
		}

		if err := segmentf.filter(segments); err != nil {
			return "", nil
		}

		thresholdFunc = func(x, y int, c color.Color) bool {
			p := segments.Pixels[y][x]
			v := segmentf.value(p)

			return v >= opts.Low && v <= opts.High
		}
	}

	bm := gotrace.NewBitmapFromImage(src, thresholdFunc)

	paths, err := gotrace.Trace(bm, &opts.TraceParams)
	if err != nil {
		return "", err
	}

	var b bytes.Buffer

	if err := gotrace.WriteSvg(&b, src.Bounds(), paths, opts.Color); err != nil {
		return "", err
	}

	return b.String(), nil
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

func createSegmentMap(src image.Image) segmentator.Image {
	var img segmentator.Image

	bounds := src.Bounds()
	img.Width = bounds.Max.X
	img.Height = bounds.Max.Y

	for x := 0; x < img.Height; x++ {
		var row []segmentator.Pixel
		for y := 0; y < img.Width; y++ {
			r, g, b, a := src.At(y, x).RGBA()
			p := segmentator.Pixel{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}
			row = append(row, p)
		}
		img.Pixels = append(img.Pixels, row)
	}

	return img
}
