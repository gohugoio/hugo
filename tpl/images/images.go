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
	"path"
	"sync"

	"github.com/bep/overlayfs"
	"github.com/gohugoio/hugo/common/hashing"
	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/resources/images"
	"github.com/gohugoio/hugo/resources/resource_factories/create"
	"github.com/mitchellh/mapstructure"
	"rsc.io/qr"

	// Importing image codecs for image.DecodeConfig
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	// Import webp codec
	_ "golang.org/x/image/webp"

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
		cache:        map[string]image.Config{},
		deps:         d,
		createClient: create.New(d.ResourceSpec),
	}
}

// Namespace provides template functions for the "images" namespace.
type Namespace struct {
	*images.Filters
	readFileFs   afero.Fs
	cacheMu      sync.RWMutex
	cache        map[string]image.Config
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

	// Check cache for image config.
	ns.cacheMu.RLock()
	config, ok := ns.cache[filename]
	ns.cacheMu.RUnlock()

	if ok {
		return config, nil
	}

	f, err := ns.readFileFs.Open(filename)
	if err != nil {
		return image.Config{}, err
	}
	defer f.Close()

	config, _, err = image.DecodeConfig(f)
	if err != nil {
		return config, err
	}

	ns.cacheMu.Lock()
	ns.cache[filename] = config
	ns.cacheMu.Unlock()

	return config, nil
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
func (ns *Namespace) QR(options any) (images.ImageResource, error) {
	const (
		qrDefaultErrorCorrectionLevel = "medium"
		qrDefaultScale                = 4
	)

	opts := struct {
		Text      string // text to encode
		Level     string // error correction level; one of low, medium, quartile, or high
		Scale     int    // number of image pixels per QR code module
		TargetDir string // target directory relative to publishDir
	}{
		Level: qrDefaultErrorCorrectionLevel,
		Scale: qrDefaultScale,
	}

	err := mapstructure.WeakDecode(options, &opts)
	if err != nil {
		return nil, err
	}

	if opts.Text == "" {
		return nil, errors.New("cannot encode an empty string")
	}

	level, ok := qrErrorCorrectionLevels[opts.Level]
	if !ok {
		return nil, errors.New("error correction level must be one of low, medium, quartile, or high")
	}

	if opts.Scale < 2 {
		return nil, errors.New("scale must be an integer greater than or equal to 2")
	}

	targetPath := path.Join(opts.TargetDir, fmt.Sprintf("qr_%s.png", hashing.HashStringHex(opts)))

	r, err := ns.createClient.FromOpts(
		create.Options{
			TargetPath:        targetPath,
			TargetPathHasHash: true,
			CreateContent: func() (func() (hugio.ReadSeekCloser, error), error) {
				code, err := qr.Encode(opts.Text, level)
				if err != nil {
					return nil, err
				}
				code.Scale = opts.Scale
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
