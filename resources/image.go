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
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io"
	"os"
	"strings"
	"sync"

	color_extractor "github.com/marekm4/color-extractor"

	"github.com/gohugoio/hugo/cache/filecache"
	"github.com/gohugoio/hugo/common/hashing"
	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/config"

	"github.com/gohugoio/gift"

	"github.com/gohugoio/hugo/resources/images/meta"
	"github.com/gohugoio/hugo/resources/internal"

	"github.com/gohugoio/hugo/resources/resource"

	"github.com/gohugoio/hugo/resources/images"
)

var (
	_ images.ImageResource            = (*imageResource)(nil)
	_ resource.Source                 = (*imageResource)(nil)
	_ resource.Cloner                 = (*imageResource)(nil)
	_ resource.NameNormalizedProvider = (*imageResource)(nil)
	_ targetPathProvider              = (*imageResource)(nil)
)

// imageResource represents an image resource.
type imageResource struct {
	*images.Image

	// When a image is processed in a chain, this holds the reference to the
	// original (first).
	root *imageResource

	// These 2 are only set for the root imageResource.
	exifInfoFn func() (*meta.ExifInfo, error)
	metaInfoFn func() (*meta.MetaInfo, error)

	colorsFn func() ([]images.Color, error)

	baseResource
}

func newImageResource(img *images.Image, base baseResource) *imageResource {
	ir := &imageResource{
		Image:        img,
		baseResource: base,
	}
	ir.root = ir
	ir.exifInfoFn = ir.newExifInfoFn()
	ir.metaInfoFn = ir.newMetaInfoFn()
	ir.colorsFn = ir.newColorsFn()
	return ir
}

func (i *imageResource) newExifInfoFn() func() (*meta.ExifInfo, error) {
	return sync.OnceValues(func() (*meta.ExifInfo, error) {
		hugo.Deprecate("Image.Exif", "Use Image.Meta, see https://gohugo.io/content-management/image-processing/#meta", "v0.155.0")
		mf := i.Format.ToImageMetaImageFormatFormat()
		if mf == -1 {
			return nil, nil
		}

		var result *meta.ExifInfo
		key := i.getImageMetaCacheTargetPath()

		read := func(info filecache.ItemInfo, r io.ReadSeeker) error {
			data, err := io.ReadAll(r)
			if err != nil {
				return err
			}
			return json.Unmarshal(data, &result)
		}

		create := func(info filecache.ItemInfo, w io.WriteCloser) (err error) {
			defer w.Close()
			f, err := i.ReadSeekCloser()
			if err != nil {
				return err
			}
			defer f.Close()

			filename := i.getResourcePaths().Path()
			result, err = i.getSpec().Imaging.DecodeExif(filename, mf, f)
			if err != nil {
				i.getSpec().Logger.Warnf("Unable to decode Exif metadata from image: %s", i.Key())
				return nil
			}

			enc := json.NewEncoder(w)
			return enc.Encode(result)
		}

		_, err := i.getSpec().ImageCache.fcache.ReadOrCreate(key, read, create)
		return result, err
	})
}

func (i *imageResource) newMetaInfoFn() func() (*meta.MetaInfo, error) {
	return sync.OnceValues(func() (*meta.MetaInfo, error) {
		mf := i.Format.ToImageMetaImageFormatFormat()
		if mf == -1 {
			return nil, nil
		}

		var result *meta.MetaInfo
		key := i.getImageMetaInfoCacheTargetPath()

		read := func(info filecache.ItemInfo, r io.ReadSeeker) error {
			data, err := io.ReadAll(r)
			if err != nil {
				return err
			}
			return json.Unmarshal(data, &result)
		}

		create := func(info filecache.ItemInfo, w io.WriteCloser) (err error) {
			defer w.Close()
			f, err := i.ReadSeekCloser()
			if err != nil {
				return err
			}
			defer f.Close()

			filename := i.getResourcePaths().Path()
			result, err = i.getSpec().Imaging.DecodeMeta(filename, mf, f)
			if err != nil {
				i.getSpec().Logger.Warnf("Unable to decode metadata from image: %s", i.Key())
				return nil
			}

			enc := json.NewEncoder(w)
			return enc.Encode(result)
		}

		_, err := i.getSpec().ImageCache.fcache.ReadOrCreate(key, read, create)
		return result, err
	})
}

func (i *imageResource) newColorsFn() func() ([]images.Color, error) {
	return sync.OnceValues(func() ([]images.Color, error) {
		img, err := i.DecodeImage()
		if err != nil {
			return nil, err
		}
		colors := color_extractor.ExtractColors(img)
		result := make([]images.Color, len(colors))
		for j, c := range colors {
			result[j] = images.ColorGoToColor(c)
		}
		return result, nil
	})
}

func (i *imageResource) Exif() *meta.ExifInfo {
	x, err := i.root.exifInfoFn()
	if err != nil {
		panic(fmt.Sprintf("exif init failed: %s", err))
	}
	return x
}

func (i *imageResource) Meta() *meta.MetaInfo {
	m, err := i.root.metaInfoFn()
	if err != nil {
		panic(fmt.Sprintf("meta init failed: %s", err))
	}
	return m
}

func (i *imageResource) getImageMetaInfoCacheTargetPath() string {
	// Increment to invalidate the meta cache
	const imageMetaInfoVersionNumber = 1

	cfgHash := i.getSpec().Imaging.Cfg.SourceHash
	df := i.getResourcePaths()
	p1, _ := paths.FileAndExt(df.File)
	h := i.hash()
	idStr := hashing.HashStringHex(h, i.size(), imageMetaInfoVersionNumber, cfgHash)
	df.File = fmt.Sprintf("%s_%s_meta.json", p1, idStr)
	return df.TargetPath()
}

// Colors returns a slice of the most dominant colors in an image
// using a simple histogram method.
func (i *imageResource) Colors() ([]images.Color, error) {
	return i.colorsFn()
}

func (i *imageResource) targetPath() string {
	return i.TargetPath()
}

// Clone is for internal use.
func (i *imageResource) Clone() resource.Resource {
	gr := i.baseResource.Clone().(baseResource)
	ir := &imageResource{
		root:         i.root,
		Image:        i.WithSpec(gr),
		baseResource: gr,
	}
	ir.colorsFn = ir.newColorsFn()
	return ir
}

func (i *imageResource) cloneTo(targetPath string) resource.Resource {
	gr := i.baseResource.cloneTo(targetPath).(baseResource)
	ir := &imageResource{
		root:         i.root,
		Image:        i.WithSpec(gr),
		baseResource: gr,
	}
	ir.colorsFn = ir.newColorsFn()
	return ir
}

func (i *imageResource) cloneWithUpdates(u *transformationUpdate) (baseResource, error) {
	base, err := i.baseResource.cloneWithUpdates(u)
	if err != nil {
		return nil, err
	}

	var img *images.Image

	if u.isContentChanged() {
		img = i.WithSpec(base)
	} else {
		img = i.Image
	}

	ir := &imageResource{
		root:         i.root,
		Image:        img,
		baseResource: base,
	}
	ir.colorsFn = ir.newColorsFn()
	return ir, nil
}

// Process processes the image with the given spec.
// The spec can contain an optional action, one of "resize", "crop", "fit" or "fill".
// This makes this method a more flexible version that covers all of Resize, Crop, Fit and Fill,
// but it also supports e.g. format conversions without any resize action.
func (i *imageResource) Process(spec string) (images.ImageResource, error) {
	return i.processActionSpec("", spec)
}

// Resize resizes the image to the specified width and height using the specified resampling
// filter and returns the transformed image. If one of width or height is 0, the image aspect
// ratio is preserved.
func (i *imageResource) Resize(spec string) (images.ImageResource, error) {
	return i.processActionSpec(images.ActionResize, spec)
}

// Crop the image to the specified dimensions without resizing using the given anchor point.
// Space delimited config, e.g. `200x300 TopLeft`.
func (i *imageResource) Crop(spec string) (images.ImageResource, error) {
	return i.processActionSpec(images.ActionCrop, spec)
}

// Fit scales down the image using the specified resample filter to fit the specified
// maximum width and height.
func (i *imageResource) Fit(spec string) (images.ImageResource, error) {
	return i.processActionSpec(images.ActionFit, spec)
}

// Fill scales the image to the smallest possible size that will cover the specified dimensions,
// crops the resized image to the specified dimensions using the given anchor point.
// Space delimited config, e.g. `200x300 TopLeft`.
func (i *imageResource) Fill(spec string) (images.ImageResource, error) {
	return i.processActionSpec(images.ActionFill, spec)
}

func (i *imageResource) Filter(filters ...any) (images.ImageResource, error) {
	var confMain images.ImageConfig

	var gfilters []gift.Filter

	for _, f := range filters {
		gfilters = append(gfilters, images.ToFilters(f)...)
	}

	var options []string

	for _, f := range gfilters {
		f = images.UnwrapFilter(f)
		if specProvider, ok := f.(images.ImageProcessSpecProvider); ok {
			options = append(options, strings.Fields(specProvider.ImageProcessSpec())...)
		}
	}

	confMain, err := images.DecodeImageConfig(options, i.Proc.Cfg, i.Format)
	if err != nil {
		return nil, err
	}

	confMain.Action = "filter"
	confMain.Key = hashing.HashString(gfilters)

	return i.doWithImageConfig(confMain, func(src image.Image) (image.Image, error) {
		var filters []gift.Filter
		for _, f := range gfilters {
			f = images.UnwrapFilter(f)
			if specProvider, ok := f.(images.ImageProcessSpecProvider); ok {
				options := strings.Fields(specProvider.ImageProcessSpec())
				conf, err := images.DecodeImageConfig(options, i.Proc.Cfg, i.Format)
				if err != nil {
					return nil, err
				}
				pFilters, err := i.Proc.FiltersFromConfig(src, conf)
				if err != nil {
					return nil, err
				}
				filters = append(filters, pFilters...)
			} else if orientationProvider, ok := f.(images.ImageFilterFromOrientationProvider); ok {
				var orientation int
				if meta := i.Meta(); meta != nil {
					orientation = meta.Orientation
				}
				if tf := orientationProvider.AutoOrient(orientation); tf != nil {
					filters = append(filters, tf)
				}
			} else {
				filters = append(filters, f)
			}
		}
		return i.Proc.Filter(src, filters...)
	})
}

func (i *imageResource) processActionSpec(action, spec string) (images.ImageResource, error) {
	options := append([]string{action}, strings.Fields(strings.ToLower(spec))...)
	ir, err := i.processOptions(options)
	if err != nil {
		if sourcePath := i.sourcePath(); sourcePath != "" {
			err = fmt.Errorf("failed to %s image %q: %w", action, sourcePath, err)
		}
		return nil, err
	}
	return ir, nil
}

func (i *imageResource) processOptions(options []string) (images.ImageResource, error) {
	conf, err := images.DecodeImageConfig(options, i.Proc.Cfg, i.Format)
	if err != nil {
		return nil, err
	}

	img, err := i.doWithImageConfig(conf, func(src image.Image) (image.Image, error) {
		return i.Proc.ApplyFiltersFromConfig(src, conf)
	})
	if err != nil {
		return nil, err
	}

	if conf.Action == images.ActionFill {
		if conf.Anchor == images.SmartCropAnchor && img.Width() == 0 || img.Height() == 0 {
			// See https://github.com/gohugoio/hugo/issues/7955
			// Smartcrop fails silently in some rare cases.
			// Fall back to a center fill.
			conf = conf.Reanchor(gift.CenterAnchor)
			return i.doWithImageConfig(conf, func(src image.Image) (image.Image, error) {
				return i.Proc.ApplyFiltersFromConfig(src, conf)
			})
		}
	}

	return img, nil
}

var imageProcSem chan bool

func init() {
	imageProcWorkers := 1
	if n := config.GetNumWorkerMultiplier(); n > 4 {
		imageProcWorkers = 2
	}
	imageProcSem = make(chan bool, imageProcWorkers)
}

func (i *imageResource) doWithImageConfig(conf images.ImageConfig, f func(src image.Image) (image.Image, error)) (images.ImageResource, error) {
	img, err := i.getSpec().ImageCache.getOrCreate(i, conf, func() (*imageResource, image.Image, error) {
		imageProcSem <- true
		defer func() {
			<-imageProcSem
		}()

		src, err := i.DecodeImage()
		if err != nil {
			return nil, nil, &os.PathError{Op: conf.Action, Path: i.TargetPath(), Err: err}
		}

		converted, err := f(src)
		if err != nil {
			return nil, nil, &os.PathError{Op: conf.Action, Path: i.TargetPath(), Err: err}
		}

		hasAlpha := !images.IsOpaque(converted)
		shouldFill := conf.BgColor != nil && hasAlpha
		shouldFill = shouldFill || (!conf.TargetFormat.SupportsTransparency() && hasAlpha)
		var bgColor color.Color

		if shouldFill {
			bgColor = conf.BgColor
			if bgColor == nil {
				bgColor = i.Proc.Cfg.Config.BgColor
			}
			tmp := image.NewRGBA(converted.Bounds())
			draw.Draw(tmp, tmp.Bounds(), image.NewUniform(bgColor), image.Point{}, draw.Src)
			draw.Draw(tmp, tmp.Bounds(), converted, converted.Bounds().Min, draw.Over)
			converted = tmp
		}

		if conf.TargetFormat == images.PNG {
			// Apply the colour palette from the source
			if paletted, ok := src.(*image.Paletted); ok {
				palette := paletted.Palette
				if bgColor != nil && len(palette) < 256 {
					palette = images.AddColorToPalette(bgColor, palette)
				} else if bgColor != nil {
					images.ReplaceColorInPalette(bgColor, palette)
				}
				tmp := image.NewPaletted(converted.Bounds(), palette)
				draw.FloydSteinberg.Draw(tmp, tmp.Bounds(), converted, converted.Bounds().Min)
				converted = tmp
			}
		}

		ci := i.clone(converted)
		targetPath := i.relTargetPathFromConfig(conf, i.getSpec().Imaging.Cfg.SourceHash)
		ci.setTargetPath(targetPath)
		ci.Format = conf.TargetFormat
		ci.setMediaType(conf.TargetFormat.MediaType())

		return ci, converted, nil
	})
	if err != nil {
		return nil, err
	}
	return img, nil
}

// DecodeImage decodes the image source into an Image.
// This for internal use only.
func (i *imageResource) DecodeImage() (image.Image, error) {
	f, err := i.ReadSeekCloser()
	if err != nil {
		return nil, fmt.Errorf("failed to open image for decode: %w", err)
	}
	defer f.Close()

	return i.getSpec().Imaging.Codec.DecodeFormat(i.Format, f)
}

func (i *imageResource) clone(img image.Image) *imageResource {
	spec := i.baseResource.Clone().(baseResource)

	var image *images.Image
	if img != nil {
		image = i.WithImage(img)
	} else {
		image = i.WithSpec(spec)
	}

	ir := &imageResource{
		Image:        image,
		root:         i.root,
		baseResource: spec,
	}
	ir.colorsFn = ir.newColorsFn()
	return ir
}

func (i *imageResource) getImageMetaCacheTargetPath() string {
	// Increment to invalidate the meta cache
	// Last increment: v0.130.0 when change to the new imagemeta library for Exif.
	const imageMetaVersionNumber = 2

	cfgHash := i.getSpec().Imaging.Cfg.SourceHash
	df := i.getResourcePaths()
	p1, _ := paths.FileAndExt(df.File)
	h := i.hash()
	idStr := hashing.HashStringHex(h, i.size(), imageMetaVersionNumber, cfgHash)
	df.File = fmt.Sprintf("%s_%s.json", p1, idStr)
	return df.TargetPath()
}

func (i *imageResource) relTargetPathFromConfig(conf images.ImageConfig, imagingConfigSourceHash string) internal.ResourcePaths {
	p1, p2 := paths.FileAndExt(i.getResourcePaths().File)
	if conf.TargetFormat != i.Format {
		p2 = conf.TargetFormat.DefaultExtension()
	}

	// Do not change.
	const imageHashPrefix = "_hu_"

	huIdx := strings.LastIndex(p1, imageHashPrefix)
	incomingID := ""
	if huIdx > -1 {
		incomingID = p1[huIdx+len(imageHashPrefix):]
		p1 = p1[:huIdx]
	}

	hash := hashing.HashStringHex(incomingID, i.hash(), conf.Key, imagingConfigSourceHash)
	rp := i.getResourcePaths()
	rp.File = fmt.Sprintf("%s%s%s%s", p1, imageHashPrefix, hash, p2)

	return rp
}
