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
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"
	"strings"

	"github.com/spf13/cast"

	"github.com/gohugoio/hugo/media"

	"github.com/gohugoio/hugo/resources/resource"

	_errors "github.com/pkg/errors"

	"github.com/disintegration/imaging"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/resources/images"

	// Blind import for image.Decode
	_ "image/gif"
	_ "image/png"

	// Blind import for image.Decode
	_ "golang.org/x/image/webp"
)

var (
	_ resource.Image  = (*imageResource)(nil)
	_ resource.Source = (*imageResource)(nil)
	_ resource.Cloner = (*imageResource)(nil)
	_ Transformer     = (*imageResource)(nil)
)

// ImageResource represents an image resource.
type imageResource struct {
	*images.Image

	baseResource
}

func (i *imageResource) Transform(t ResourceTransformation) (resource.Resource, error) {
	ic := i.Clone()

	return ic, nil
}

// CloneWithNewBase implements the Cloner interface.
func (i *imageResource) CloneWithNewBase(base string) resource.Resource {
	gr := i.baseResource.CloneWithNewBase(base).(baseResource)
	return &imageResource{
		Image:        i.WithSpec(gr),
		baseResource: gr}
}

func (i *imageResource) Clone() resource.Resource {
	gr := i.baseResource.Clone().(baseResource)
	return &imageResource{
		Image:        i.WithSpec(gr),
		baseResource: gr}
}

// Resize resizes the image to the specified width and height using the specified resampling
// filter and returns the transformed image. If one of width or height is 0, the image aspect
// ratio is preserved.
func (i *imageResource) Resize(spec string) (resource.Image, error) {
	return i.doWithImageConfig("resize", spec, func(src image.Image, conf images.ImageConfig) (image.Image, error) {
		return i.Proc.Resize(src, conf)
	})
}

// Fit scales down the image using the specified resample filter to fit the specified
// maximum width and height.
func (i *imageResource) Fit(spec string) (resource.Image, error) {
	return i.doWithImageConfig("fit", spec, func(src image.Image, conf images.ImageConfig) (image.Image, error) {
		return i.Proc.Fit(src, conf)
	})
}

// Fill scales the image to the smallest possible size that will cover the specified dimensions,
// crops the resized image to the specified dimensions using the given anchor point.
// Space delimited config: 200x300 TopLeft
func (i *imageResource) Fill(spec string) (resource.Image, error) {
	return i.doWithImageConfig("fill", spec, func(src image.Image, conf images.ImageConfig) (image.Image, error) {
		return i.Proc.Fill(src, conf)
	})
}

func (i *imageResource) Trace(opts ...interface{}) (resource.Resource, error) {
	var optsm map[string]interface{}
	if len(opts) > 0 {
		optsm = cast.ToStringMap(opts[0])
	}
	o, err := i.Proc.DecodeTraceOptions(optsm)
	if err != nil {
		return nil, err
	}

	conf := images.ImageConfig{
		Action:       "trace",
		TraceOptions: o,
	}

	return i.doWithImageConfigBase(conf, func(src image.Image, conf images.ImageConfig) (interface{}, error) {
		s, err := i.Proc.Trace(src, o)
		return s, err
	})

}

func (i *imageResource) isJPEG() bool {
	name := strings.ToLower(i.getResourcePaths().relTargetDirFile.file)
	return strings.HasSuffix(name, ".jpg") || strings.HasSuffix(name, ".jpeg")
}

// Serialize image processing. The imaging library spins up its own set of Go routines,
// so there is not much to gain from adding more load to the mix. That
// can even have negative effect in low resource scenarios.
// Note that this only effects the non-cached scenario. Once the processed
// image is written to disk, everything is fast, fast fast.
const imageProcWorkers = 1

var imageProcSem = make(chan bool, imageProcWorkers)

func (i *imageResource) doWithImageConfig(action, spec string, f func(src image.Image, conf images.ImageConfig) (image.Image, error)) (resource.Image, error) {
	conf, err := i.decodeImageConfig(action, spec)
	if err != nil {
		return nil, err
	}
	r, err := i.doWithImageConfigBase(conf, func(src image.Image, conf images.ImageConfig) (interface{}, error) {
		v, err := f(src, conf)
		if err != nil {
			return nil, err
		}
		return v.(image.Image), nil
	})
	if err != nil {
		return nil, err
	}
	return r.(resource.Image), nil
}

func (i *imageResource) decodeImageConfig(action, spec string) (images.ImageConfig, error) {
	conf, err := images.DecodeImageConfig(action, spec, i.Proc.Cfg)
	if err != nil {
		return conf, err
	}

	iconf := i.Proc.Cfg

	if conf.Quality <= 0 && i.isJPEG() {
		// We need a quality setting for all JPEGs
		conf.Quality = iconf.Quality
	}

	return conf, nil
}

func (i *imageResource) doWithImageConfigBase(conf images.ImageConfig, f func(src image.Image, conf images.ImageConfig) (interface{}, error)) (resource.Resource, error) {

	return i.getSpec().imageCache.getOrCreate(i, conf, func() (baseResource, interface{}, error) {
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

		if conf.Rotate != 0 {
			// Rotate it before any scaling to get the dimensions correct.
			src = imaging.Rotate(src, float64(conf.Rotate), color.Transparent)
		}

		convertedv, err := f(src, conf)
		if err != nil {
			return nil, nil, &os.PathError{Op: errOp, Path: errPath, Err: err}
		}

		if convertedv == nil {
			panic("converted is nil")
		}

		if s, ok := convertedv.(string); ok {
			// SVG TODO1
			ci := i.baseResource.Clone().(baseResource)
			ci.setMediaType(media.SVGType)
			/*ci.setOpenReadSeekerCloser(
				func() (hugio.ReadSeekCloser, error) {
					return hugio.NewReadSeekerNoOpCloserFromString(s), nil
				},
			)*/
			//oldname := ci.getResourcePaths().relTargetDirFile.file
			//ci.getResourcePaths().relTargetDirFile.file = oldname + ".svg"
			return ci, s, nil
		}

		converted := convertedv.(image.Image)

		if i.Format == imaging.PNG {
			// Apply the colour palette from the source
			if paletted, ok := src.(*image.Paletted); ok {
				tmp := image.NewPaletted(converted.Bounds(), paletted.Palette)
				draw.FloydSteinberg.Draw(tmp, tmp.Bounds(), converted, converted.Bounds().Min)
				converted = tmp
			}
		}

		ci := i.clone(converted)
		ci.setBasePath(conf)

		return ci, converted, nil
	})

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
		baseResource: spec,
	}
}

func (i *imageResource) setBasePath(conf images.ImageConfig) {
	i.getResourcePaths().relTargetDirFile = i.relTargetPathFromConfig(conf)
}

func (i *imageResource) relTargetPathFromConfig(conf images.ImageConfig) dirFile {
	p1, p2 := helpers.FileAndExt(i.getResourcePaths().relTargetDirFile.file)
	if conf.Action == "trace" {
		p2 = ".svg"
	}

	h, _ := i.hash()
	idStr := fmt.Sprintf("_hu%s_%d", h, i.size())

	// Do not change for no good reason.
	const md5Threshold = 100

	key := conf.Key(i.Format)

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
