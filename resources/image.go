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
	_ "image/gif"
	_ "image/png"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/disintegration/gift"

	"github.com/gohugoio/hugo/cache/filecache"
	"github.com/gohugoio/hugo/resources/images/exif"

	"github.com/gohugoio/hugo/resources/resource"

	"github.com/pkg/errors"
	_errors "github.com/pkg/errors"

	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/resources/images"

	// Blind import for image.Decode
	_ "golang.org/x/image/webp"
)

var (
	_ resource.Image  = (*imageResource)(nil)
	_ resource.Source = (*imageResource)(nil)
	_ resource.Cloner = (*imageResource)(nil)
)

// ImageResource represents an image resource.
type imageResource struct {
	*images.Image

	// When a image is processed in a chain, this holds the reference to the
	// original (first).
	root *imageResource

	metaInit    sync.Once
	metaInitErr error
	meta        *imageMeta

	baseResource
}

type imageMeta struct {
	Exif *exif.Exif
}

func (i *imageResource) Exif() (*exif.Exif, error) {
	return i.root.getExif()
}

func (i *imageResource) getExif() (*exif.Exif, error) {

	i.metaInit.Do(func() {

		supportsExif := i.Format == images.JPEG || i.Format == images.TIFF
		if !supportsExif {
			return

		}

		key := i.getImageMetaCacheTargetPath()

		read := func(info filecache.ItemInfo, r io.ReadSeeker) error {
			meta := &imageMeta{}
			data, err := ioutil.ReadAll(r)
			if err != nil {
				return err
			}

			if err = json.Unmarshal(data, &meta); err != nil {
				return err
			}

			i.meta = meta

			return nil
		}

		create := func(info filecache.ItemInfo, w io.WriteCloser) (err error) {

			f, err := i.root.ReadSeekCloser()
			if err != nil {
				i.metaInitErr = err
				return
			}
			defer f.Close()

			x, err := i.getSpec().imaging.DecodeExif(f)
			if err != nil {
				i.metaInitErr = err
				return
			}

			i.meta = &imageMeta{Exif: x}

			// Also write it to cache
			enc := json.NewEncoder(w)
			return enc.Encode(i.meta)

		}

		_, i.metaInitErr = i.getSpec().imageCache.fileCache.ReadOrCreate(key, read, create)

	})

	if i.metaInitErr != nil {
		return nil, i.metaInitErr
	}

	return i.meta.Exif, nil
}

func (i *imageResource) Clone() resource.Resource {
	gr := i.baseResource.Clone().(baseResource)
	return &imageResource{
		root:         i.root,
		Image:        i.WithSpec(gr),
		baseResource: gr,
	}
}

func (i *imageResource) cloneWithUpdates(u *transformationUpdate) (baseResource, error) {
	base, err := i.baseResource.cloneWithUpdates(u)
	if err != nil {
		return nil, err
	}

	var img *images.Image

	if u.isContenChanged() {
		img = i.WithSpec(base)
	} else {
		img = i.Image
	}

	return &imageResource{
		root:         i.root,
		Image:        img,
		baseResource: base,
	}, nil
}

// Resize resizes the image to the specified width and height using the specified resampling
// filter and returns the transformed image. If one of width or height is 0, the image aspect
// ratio is preserved.
func (i *imageResource) Resize(spec string) (resource.Image, error) {
	conf, err := i.decodeImageConfig("resize", spec)
	if err != nil {
		return nil, err
	}

	return i.doWithImageConfig(conf, func(src image.Image) (image.Image, error) {
		return i.Proc.ApplyFiltersFromConfig(src, conf)
	})
}

// Fit scales down the image using the specified resample filter to fit the specified
// maximum width and height.
func (i *imageResource) Fit(spec string) (resource.Image, error) {
	conf, err := i.decodeImageConfig("fit", spec)
	if err != nil {
		return nil, err
	}

	return i.doWithImageConfig(conf, func(src image.Image) (image.Image, error) {
		return i.Proc.ApplyFiltersFromConfig(src, conf)
	})
}

// Fill scales the image to the smallest possible size that will cover the specified dimensions,
// crops the resized image to the specified dimensions using the given anchor point.
// Space delimited config: 200x300 TopLeft
func (i *imageResource) Fill(spec string) (resource.Image, error) {
	conf, err := i.decodeImageConfig("fill", spec)
	if err != nil {
		return nil, err
	}

	return i.doWithImageConfig(conf, func(src image.Image) (image.Image, error) {
		return i.Proc.ApplyFiltersFromConfig(src, conf)
	})
}

func (i *imageResource) Filter(filters ...interface{}) (resource.Image, error) {
	conf := i.Proc.GetDefaultImageConfig("filter")

	var gfilters []gift.Filter

	for _, f := range filters {
		gfilters = append(gfilters, images.ToFilters(f)...)
	}

	conf.Key = helpers.HashString(gfilters)
	conf.TargetFormat = i.Format

	return i.doWithImageConfig(conf, func(src image.Image) (image.Image, error) {
		return i.Proc.Filter(src, gfilters...)
	})
}

// Serialize image processing. The imaging library spins up its own set of Go routines,
// so there is not much to gain from adding more load to the mix. That
// can even have negative effect in low resource scenarios.
// Note that this only effects the non-cached scenario. Once the processed
// image is written to disk, everything is fast, fast fast.
const imageProcWorkers = 1

var imageProcSem = make(chan bool, imageProcWorkers)

func (i *imageResource) doWithImageConfig(conf images.ImageConfig, f func(src image.Image) (image.Image, error)) (resource.Image, error) {
	img, err := i.getSpec().imageCache.getOrCreate(i, conf, func() (*imageResource, image.Image, error) {
		imageProcSem <- true
		defer func() {
			<-imageProcSem
		}()

		errOp := conf.Action
		errPath := i.getSourceFilename()

		src, err := i.decodeSource()
		if err != nil {
			return nil, nil, &os.PathError{Op: errOp, Path: errPath, Err: err}
		}

		converted, err := f(src)
		if err != nil {
			return nil, nil, &os.PathError{Op: errOp, Path: errPath, Err: err}
		}

		hasAlpha := !images.IsOpaque(converted)
		shouldFill := conf.BgColor != nil && hasAlpha
		shouldFill = shouldFill || (!conf.TargetFormat.SupportsTransparency() && hasAlpha)
		var bgColor color.Color

		if shouldFill {
			bgColor = conf.BgColor
			if bgColor == nil {
				bgColor = i.Proc.Cfg.BgColor
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
		ci.setBasePath(conf)
		ci.Format = conf.TargetFormat
		ci.setMediaType(conf.TargetFormat.MediaType())

		return ci, converted, nil
	})

	if err != nil {
		if i.root != nil && i.root.getFileInfo() != nil {
			return nil, errors.Wrapf(err, "image %q", i.root.getFileInfo().Meta().Filename())
		}
	}
	return img, nil
}

func (i *imageResource) decodeImageConfig(action, spec string) (images.ImageConfig, error) {
	conf, err := images.DecodeImageConfig(action, spec, i.Proc.Cfg.Cfg)
	if err != nil {
		return conf, err
	}

	// default to the source format
	if conf.TargetFormat == 0 {
		conf.TargetFormat = i.Format
	}

	if conf.Quality <= 0 && conf.TargetFormat.RequiresDefaultQuality() {
		// We need a quality setting for all JPEGs
		conf.Quality = i.Proc.Cfg.Cfg.Quality
	}

	if conf.BgColor == nil && conf.TargetFormat != i.Format {
		if i.Format.SupportsTransparency() && !conf.TargetFormat.SupportsTransparency() {
			conf.BgColor = i.Proc.Cfg.BgColor
			conf.BgColorStr = i.Proc.Cfg.Cfg.BgColor
		}
	}

	return conf, nil
}

func (i *imageResource) decodeSource() (image.Image, error) {
	f, err := i.ReadSeekCloser()
	if err != nil {
		return nil, _errors.Wrap(err, "failed to open image for decode")
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	return img, err
}

func (i *imageResource) clone(img image.Image) *imageResource {
	spec := i.baseResource.Clone().(baseResource)

	var image *images.Image
	if img != nil {
		image = i.WithImage(img)
	} else {
		image = i.WithSpec(spec)
	}

	return &imageResource{
		Image:        image,
		root:         i.root,
		baseResource: spec,
	}
}

func (i *imageResource) setBasePath(conf images.ImageConfig) {
	i.getResourcePaths().relTargetDirFile = i.relTargetPathFromConfig(conf)
}

func (i *imageResource) getImageMetaCacheTargetPath() string {
	const imageMetaVersionNumber = 1 // Increment to invalidate the meta cache

	cfg := i.getSpec().imaging.Cfg.Cfg
	df := i.getResourcePaths().relTargetDirFile
	if fi := i.getFileInfo(); fi != nil {
		df.dir = filepath.Dir(fi.Meta().Path())
	}
	p1, _ := helpers.FileAndExt(df.file)
	h, _ := i.hash()
	idStr := helpers.HashString(h, i.size(), imageMetaVersionNumber, cfg)
	return path.Join(df.dir, fmt.Sprintf("%s_%s.json", p1, idStr))
}

func (i *imageResource) relTargetPathFromConfig(conf images.ImageConfig) dirFile {
	p1, p2 := helpers.FileAndExt(i.getResourcePaths().relTargetDirFile.file)
	if conf.TargetFormat != i.Format {
		p2 = conf.TargetFormat.DefaultExtension()
	}

	h, _ := i.hash()
	idStr := fmt.Sprintf("_hu%s_%d", h, i.size())

	// Do not change for no good reason.
	const md5Threshold = 100

	key := conf.GetKey(i.Format)

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
		dir:  i.getResourcePaths().relTargetDirFile.dir,
		file: fmt.Sprintf("%s%s_%s%s", p1, idStr, key, p2),
	}
}
