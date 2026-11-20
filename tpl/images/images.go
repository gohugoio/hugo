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

// Package images provides template functions for manipulating images.
package images

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"path"
	"path/filepath"

	"github.com/bep/helpers/maphelpers"
	"github.com/bep/overlayfs"
	"github.com/gohugoio/hugo/common/hashing"
	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/resources/images"
	"github.com/gohugoio/hugo/resources/resource_factories/create"
	"github.com/gohugoio/qr"
	"github.com/mitchellh/mapstructure"

	"github.com/gohugoio/hugo/deps"
	"github.com/spf13/afero"
	"github.com/spf13/cast"
)

// New returns a new instance of the images-namespaced template functions.
func New(d *deps.Deps) *Namespace {
	var readFileFs afero.Fs

	// The docshelper script does not have or need all the dependencies set up.
	if d.PathSpec != nil {
		readFileFs = overlayfs.New(overlayfs.Options{
			Fss: []afero.Fs{
				d.PathSpec.BaseFs.Work,
				d.PathSpec.BaseFs.Content.Fs,
			},
		})
	}

	return &Namespace{
		readFileFs:   readFileFs,
		Filters:      &images.Filters{},
		cache:        maphelpers.NewConcurrentMap[string, image.Config](),
		deps:         d,
		createClient: create.New(d.ResourceSpec),
	}
}

// Namespace provides template functions for the "images" namespace.
type Namespace struct {
	*images.Filters
	readFileFs   afero.Fs
	cache        *maphelpers.ConcurrentMap[string, image.Config]
	deps         *deps.Deps
	createClient *create.Client
}

// Config returns the image.Config for the specified path relative to the
// working directory.
func (ns *Namespace) Config(path any) (image.Config, error) {
	filename, err := cast.ToStringE(path)
	if err != nil {
		return image.Config{}, err
	}

	if filename == "" {
		return image.Config{}, errors.New("config needs a filename")
	}

	return ns.cache.GetOrCreate(filename, func() (image.Config, error) {
		f, err := ns.readFileFs.Open(filename)
		if err != nil {
			return image.Config{}, err
		}
		defer f.Close()

		ext := filepath.Ext(filename)
		format, _ := images.ImageFormatFromExt(ext)

		config, _, err := ns.deps.ResourceSpec.Imaging.Codec.DecodeConfig(format, f)
		return config, err
	})
}

// Filter applies the given filters to the image given as the last element in args.
func (ns *Namespace) Filter(args ...any) (images.ImageResource, error) {
	if len(args) < 2 {
		return nil, errors.New("must provide an image and one or more filters")
	}

	img := args[len(args)-1].(images.ImageResource)
	filtersv := args[:len(args)-1]

	return img.Filter(filtersv...)
}

var qrErrorCorrectionLevels = map[string]qr.Level{
	"low":      qr.L,
	"medium":   qr.M,
	"quartile": qr.Q,
	"high":     qr.H,
}

// QR encodes the given text into a QR code using the specified options,
// returning an image resource.
func (ns *Namespace) QR(args ...any) (images.ImageResource, error) {
	const (
		qrDefaultErrorCorrectionLevel = "medium"
		qrDefaultScale                = 4
	)

	opts := struct {
		Level     string // error correction level; one of low, medium, quartile, or high
		Scale     int    // number of image pixels per QR code module
		TargetDir string // target directory relative to publishDir
		FgColor   any    // foreground color as a hex string or a Color; nil or empty means black
		BgColor   any    // background color as a hex string or a Color; nil or empty means white
	}{
		Level: qrDefaultErrorCorrectionLevel,
		Scale: qrDefaultScale,
	}

	if len(args) == 0 || len(args) > 2 {
		return nil, errors.New("requires 1 or 2 arguments")
	}

	text, err := cast.ToStringE(args[0])
	if err != nil {
		return nil, err
	}

	if text == "" {
		return nil, errors.New("cannot encode an empty string")
	}

	// The hash source is the text and the options as provided by the user.
	var optsm any
	if len(args) == 2 {
		optsm = args[1]
		err := mapstructure.WeakDecode(optsm, &opts)
		if err != nil {
			return nil, err
		}
	}

	level, ok := qrErrorCorrectionLevels[opts.Level]
	if !ok {
		return nil, errors.New("error correction level must be one of low, medium, quartile, or high")
	}

	if opts.Scale < 2 {
		return nil, errors.New("scale must be an integer greater than or equal to 2")
	}

	fgColor, err := qrColor("fgColor", opts.FgColor)
	if err != nil {
		return nil, err
	}
	bgColor, err := qrColor("bgColor", opts.BgColor)
	if err != nil {
		return nil, err
	}

	targetPath := path.Join(opts.TargetDir, fmt.Sprintf("qr_%s.png", hashing.HashStringHex(text, optsm)))

	r, err := ns.createClient.FromOpts(
		create.Options{
			TargetPath:        targetPath,
			TargetPathHasHash: true,
			CreateContent: func() (func() (hugio.ReadSeekCloser, error), error) {
				code, err := qr.Encode(text, qr.Options{Level: level, Scale: opts.Scale, FgColor: fgColor, BgColor: bgColor})
				if err != nil {
					return nil, err
				}
				png := code.PNG()
				return func() (hugio.ReadSeekCloser, error) {
					return hugio.NewReadSeekerNoOpCloserFromBytes(png), nil
				}, nil
			},
		},
	)
	if err != nil {
		return nil, err
	}

	ir, ok := r.(images.ImageResource)
	if !ok {
		panic("bug: resource is not an image resource")
	}

	return ir, nil
}

// qrColor converts v, a hex color string or a Color, to a color.Color.
// It returns nil if v is nil or an empty string.
func qrColor(name string, v any) (color.Color, error) {
	if v == nil || v == "" {
		return nil, nil
	}
	c, ok, err := images.ToColorGo(v)
	if err != nil {
		return nil, fmt.Errorf("invalid %s: %w", name, err)
	}
	if !ok {
		return nil, fmt.Errorf("invalid %s: unsupported type %T", name, v)
	}
	return c, nil
}
