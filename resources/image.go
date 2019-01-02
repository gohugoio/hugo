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

package resources

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/gohugoio/hugo/resources/resource"

	_errors "github.com/pkg/errors"

	"github.com/disintegration/imaging"
	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/helpers"
	"github.com/mitchellh/mapstructure"

	// Blind import for image.Decode
	_ "image/gif"
	_ "image/png"

	// Blind import for image.Decode
	_ "golang.org/x/image/webp"
)

var (
	_ resource.Resource = (*Image)(nil)
	_ resource.Source   = (*Image)(nil)
	_ resource.Cloner   = (*Image)(nil)
)

// Imaging contains default image processing configuration. This will be fetched
// from site (or language) config.
type Imaging struct {
	// Default image quality setting (1-100). Only used for JPEG images.
	Quality int

	// Resample filter used. See https://github.com/disintegration/imaging
	ResampleFilter string

	// The anchor used in Fill. Default is "smart", i.e. Smart Crop.
	Anchor string
}

const (
	defaultJPEGQuality    = 75
	defaultResampleFilter = "box"
)

var (
	imageFormats = map[string]imaging.Format{
		".jpg":  imaging.JPEG,
		".jpeg": imaging.JPEG,
		".png":  imaging.PNG,
		".tif":  imaging.TIFF,
		".tiff": imaging.TIFF,
		".bmp":  imaging.BMP,
		".gif":  imaging.GIF,
	}

	// Add or increment if changes to an image format's processing requires
	// re-generation.
	imageFormatsVersions = map[imaging.Format]int{
		imaging.PNG: 2, // Floyd Steinberg dithering
	}

	// Increment to mark all processed images as stale. Only use when absolutely needed.
	// See the finer grained smartCropVersionNumber and imageFormatsVersions.
	mainImageVersionNumber = 0
)

var anchorPositions = map[string]imaging.Anchor{
	strings.ToLower("Center"):      imaging.Center,
	strings.ToLower("TopLeft"):     imaging.TopLeft,
	strings.ToLower("Top"):         imaging.Top,
	strings.ToLower("TopRight"):    imaging.TopRight,
	strings.ToLower("Left"):        imaging.Left,
	strings.ToLower("Right"):       imaging.Right,
	strings.ToLower("BottomLeft"):  imaging.BottomLeft,
	strings.ToLower("Bottom"):      imaging.Bottom,
	strings.ToLower("BottomRight"): imaging.BottomRight,
}

var imageFilters = map[string]imaging.ResampleFilter{
	strings.ToLower("NearestNeighbor"):   imaging.NearestNeighbor,
	strings.ToLower("Box"):               imaging.Box,
	strings.ToLower("Linear"):            imaging.Linear,
	strings.ToLower("Hermite"):           imaging.Hermite,
	strings.ToLower("MitchellNetravali"): imaging.MitchellNetravali,
	strings.ToLower("CatmullRom"):        imaging.CatmullRom,
	strings.ToLower("BSpline"):           imaging.BSpline,
	strings.ToLower("Gaussian"):          imaging.Gaussian,
	strings.ToLower("Lanczos"):           imaging.Lanczos,
	strings.ToLower("Hann"):              imaging.Hann,
	strings.ToLower("Hamming"):           imaging.Hamming,
	strings.ToLower("Blackman"):          imaging.Blackman,
	strings.ToLower("Bartlett"):          imaging.Bartlett,
	strings.ToLower("Welch"):             imaging.Welch,
	strings.ToLower("Cosine"):            imaging.Cosine,
}

// Image represents an image resource.
type Image struct {
	config       image.Config
	configInit   sync.Once
	configLoaded bool

	copyToDestinationInit sync.Once

	imaging *Imaging

	format imaging.Format

	*genericResource
}

// Width returns i's width.
func (i *Image) Width() int {
	i.initConfig()
	return i.config.Width
}

// Height returns i's height.
func (i *Image) Height() int {
	i.initConfig()
	return i.config.Height
}

// WithNewBase implements the Cloner interface.
func (i *Image) WithNewBase(base string) resource.Resource {
	return &Image{
		imaging:         i.imaging,
		format:          i.format,
		genericResource: i.genericResource.WithNewBase(base).(*genericResource)}
}

// Resize resizes the image to the specified width and height using the specified resampling
// filter and returns the transformed image. If one of width or height is 0, the image aspect
// ratio is preserved.
func (i *Image) Resize(spec string) (*Image, error) {
	return i.doWithImageConfig("resize", spec, func(src image.Image, conf imageConfig) (image.Image, error) {
		return imaging.Resize(src, conf.Width, conf.Height, conf.Filter), nil
	})
}

// Fit scales down the image using the specified resample filter to fit the specified
// maximum width and height.
func (i *Image) Fit(spec string) (*Image, error) {
	return i.doWithImageConfig("fit", spec, func(src image.Image, conf imageConfig) (image.Image, error) {
		return imaging.Fit(src, conf.Width, conf.Height, conf.Filter), nil
	})
}

// Fill scales the image to the smallest possible size that will cover the specified dimensions,
// crops the resized image to the specified dimensions using the given anchor point.
// Space delimited config: 200x300 TopLeft
func (i *Image) Fill(spec string) (*Image, error) {
	return i.doWithImageConfig("fill", spec, func(src image.Image, conf imageConfig) (image.Image, error) {
		if conf.AnchorStr == smartCropIdentifier {
			return smartCrop(src, conf.Width, conf.Height, conf.Anchor, conf.Filter)
		}
		return imaging.Fill(src, conf.Width, conf.Height, conf.Anchor, conf.Filter), nil
	})
}

// Holds configuration to create a new image from an existing one, resize etc.
type imageConfig struct {
	Action string

	// Quality ranges from 1 to 100 inclusive, higher is better.
	// This is only relevant for JPEG images.
	// Default is 75.
	Quality int

	// Rotate rotates an image by the given angle counter-clockwise.
	// The rotation will be performed first.
	Rotate int

	Width  int
	Height int

	Filter    imaging.ResampleFilter
	FilterStr string

	Anchor    imaging.Anchor
	AnchorStr string
}

func (i *Image) isJPEG() bool {
	name := strings.ToLower(i.relTargetDirFile.file)
	return strings.HasSuffix(name, ".jpg") || strings.HasSuffix(name, ".jpeg")
}

// Serialize image processing. The imaging library spins up its own set of Go routines,
// so there is not much to gain from adding more load to the mix. That
// can even have negative effect in low resource scenarios.
// Note that this only effects the non-cached scenario. Once the processed
// image is written to disk, everything is fast, fast fast.
const imageProcWorkers = 1

var imageProcSem = make(chan bool, imageProcWorkers)

func (i *Image) doWithImageConfig(action, spec string, f func(src image.Image, conf imageConfig) (image.Image, error)) (*Image, error) {
	conf, err := parseImageConfig(spec)
	if err != nil {
		return nil, err
	}
	conf.Action = action

	if conf.Quality <= 0 && i.isJPEG() {
		// We need a quality setting for all JPEGs
		conf.Quality = i.imaging.Quality
	}

	if conf.FilterStr == "" {
		conf.FilterStr = i.imaging.ResampleFilter
		conf.Filter = imageFilters[conf.FilterStr]
	}

	if conf.AnchorStr == "" {
		conf.AnchorStr = i.imaging.Anchor
		if !strings.EqualFold(conf.AnchorStr, smartCropIdentifier) {
			conf.Anchor = anchorPositions[conf.AnchorStr]
		}
	}

	return i.spec.imageCache.getOrCreate(i, conf, func() (*Image, image.Image, error) {
		imageProcSem <- true
		defer func() {
			<-imageProcSem
		}()

		ci := i.clone()

		errOp := action
		errPath := i.sourceFilename

		ci.setBasePath(conf)

		src, err := i.decodeSource()
		if err != nil {
			return nil, nil, &os.PathError{Op: errOp, Path: errPath, Err: err}
		}

		if conf.Rotate != 0 {
			// Rotate it before any scaling to get the dimensions correct.
			src = imaging.Rotate(src, float64(conf.Rotate), color.Transparent)
		}

		converted, err := f(src, conf)
		if err != nil {
			return ci, nil, &os.PathError{Op: errOp, Path: errPath, Err: err}
		}

		if i.format == imaging.PNG {
			// Apply the colour palette from the source
			if paletted, ok := src.(*image.Paletted); ok {
				tmp := image.NewPaletted(converted.Bounds(), paletted.Palette)
				draw.FloydSteinberg.Draw(tmp, tmp.Bounds(), converted, converted.Bounds().Min)
				converted = tmp
			}
		}

		b := converted.Bounds()
		ci.config = image.Config{Width: b.Max.X, Height: b.Max.Y}
		ci.configLoaded = true

		return ci, converted, nil
	})

}

func (i imageConfig) key(format imaging.Format) string {
	k := strconv.Itoa(i.Width) + "x" + strconv.Itoa(i.Height)
	if i.Action != "" {
		k += "_" + i.Action
	}
	if i.Quality > 0 {
		k += "_q" + strconv.Itoa(i.Quality)
	}
	if i.Rotate != 0 {
		k += "_r" + strconv.Itoa(i.Rotate)
	}
	anchor := i.AnchorStr
	if anchor == smartCropIdentifier {
		anchor = anchor + strconv.Itoa(smartCropVersionNumber)
	}

	k += "_" + i.FilterStr

	if strings.EqualFold(i.Action, "fill") {
		k += "_" + anchor
	}

	if v, ok := imageFormatsVersions[format]; ok {
		k += "_" + strconv.Itoa(v)
	}

	if mainImageVersionNumber > 0 {
		k += "_" + strconv.Itoa(mainImageVersionNumber)
	}

	return k
}

func newImageConfig(width, height, quality, rotate int, filter, anchor string) imageConfig {
	var c imageConfig

	c.Width = width
	c.Height = height
	c.Quality = quality
	c.Rotate = rotate

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

func parseImageConfig(config string) (imageConfig, error) {
	var (
		c   imageConfig
		err error
	)

	if config == "" {
		return c, errors.New("image config cannot be empty")
	}

	parts := strings.Fields(config)
	for _, part := range parts {
		part = strings.ToLower(part)

		if part == smartCropIdentifier {
			c.AnchorStr = smartCropIdentifier
		} else if pos, ok := anchorPositions[part]; ok {
			c.Anchor = pos
			c.AnchorStr = part
		} else if filter, ok := imageFilters[part]; ok {
			c.Filter = filter
			c.FilterStr = part
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

	if c.Width == 0 && c.Height == 0 {
		return c, errors.New("must provide Width or Height")
	}

	return c, nil
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

		f, err = i.ReadSeekCloser()
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
		return _errors.Wrap(err, "failed to load image config")
	}

	return nil
}

func (i *Image) decodeSource() (image.Image, error) {
	f, err := i.ReadSeekCloser()
	if err != nil {
		return nil, _errors.Wrap(err, "failed to open image for decode")
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	return img, err
}

func (i *Image) openDestinationsForWriting() (io.WriteCloser, error) {
	targetFilenames := i.targetFilenames()
	var changedFilenames []string

	// Fast path:
	// This is a processed version of the original.
	// If it exists on destination with the same filename and file size, it is
	// the same file, so no need to transfer it again.
	for _, targetFilename := range targetFilenames {
		if fi, err := i.spec.BaseFs.PublishFs.Stat(targetFilename); err == nil && fi.Size() == i.osFileInfo.Size() {
			continue
		}
		changedFilenames = append(changedFilenames, targetFilename)
	}

	if len(changedFilenames) == 0 {
		return struct {
			io.Writer
			io.Closer
		}{
			ioutil.Discard,
			ioutil.NopCloser(nil),
		}, nil

	}

	return helpers.OpenFilesForWriting(i.spec.BaseFs.PublishFs, changedFilenames...)

}

func (i *Image) encodeTo(conf imageConfig, img image.Image, w io.Writer) error {
	switch i.format {
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
		return imaging.Encode(w, img, i.format)
	}
}

func (i *Image) clone() *Image {
	g := *i.genericResource
	g.resourceContent = &resourceContent{}

	return &Image{
		imaging:         i.imaging,
		format:          i.format,
		genericResource: &g}
}

func (i *Image) setBasePath(conf imageConfig) {
	i.relTargetDirFile = i.relTargetPathFromConfig(conf)
}

func (i *Image) relTargetPathFromConfig(conf imageConfig) dirFile {
	p1, p2 := helpers.FileAndExt(i.relTargetDirFile.file)

	idStr := fmt.Sprintf("_hu%s_%d", i.hash, i.osFileInfo.Size())

	// Do not change for no good reason.
	const md5Threshold = 100

	key := conf.key(i.format)

	// It is useful to have the key in clear text, but when nesting transforms, it
	// can easily be too long to read, and maybe even too long
	// for the different OSes to handle.
	if len(p1)+len(idStr)+len(p2) > md5Threshold {
		key = helpers.MD5String(p1 + key + p2)
		huIdx := strings.Index(p1, "_hu")
		if huIdx != -1 {
			p1 = p1[:huIdx]
		} else {
			// This started out as a very long file name. Making it even longer
			// could melt ice in the Arctic.
			p1 = ""
		}
	} else if strings.Contains(p1, idStr) {
		// On scaling an already scaled image, we get the file info from the original.
		// Repeating the same info in the filename makes it stuttery for no good reason.
		idStr = ""
	}

	return dirFile{
		dir:  i.relTargetDirFile.dir,
		file: fmt.Sprintf("%s%s_%s%s", p1, idStr, key, p2),
	}

}

func decodeImaging(m map[string]interface{}) (Imaging, error) {
	var i Imaging
	if err := mapstructure.WeakDecode(m, &i); err != nil {
		return i, err
	}

	if i.Quality == 0 {
		i.Quality = defaultJPEGQuality
	} else if i.Quality < 0 || i.Quality > 100 {
		return i, errors.New("JPEG quality must be a number between 1 and 100")
	}

	if i.Anchor == "" || strings.EqualFold(i.Anchor, smartCropIdentifier) {
		i.Anchor = smartCropIdentifier
	} else {
		i.Anchor = strings.ToLower(i.Anchor)
		if _, found := anchorPositions[i.Anchor]; !found {
			return i, errors.New("invalid anchor value in imaging config")
		}
	}

	if i.ResampleFilter == "" {
		i.ResampleFilter = defaultResampleFilter
	} else {
		filter := strings.ToLower(i.ResampleFilter)
		_, found := imageFilters[filter]
		if !found {
			return i, fmt.Errorf("%q is not a valid resample filter", filter)
		}
		i.ResampleFilter = filter
	}

	return i, nil
}
