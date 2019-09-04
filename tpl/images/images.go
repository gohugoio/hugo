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
	"image"
	"sync"

	"github.com/pkg/errors"

	"github.com/gohugoio/hugo/resources/images"
	"github.com/gohugoio/hugo/resources/resource"

	// Importing image codecs for image.DecodeConfig
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	// Import webp codec
	_ "golang.org/x/image/webp"

	"github.com/gohugoio/hugo/deps"
	"github.com/spf13/cast"
)

// New returns a new instance of the images-namespaced template functions.
func New(deps *deps.Deps) *Namespace {
	return &Namespace{
		Filters: &images.Filters{},
		cache:   map[string]image.Config{},
		deps:    deps,
	}
}

// Namespace provides template functions for the "images" namespace.
type Namespace struct {
	*images.Filters
	cacheMu sync.RWMutex
	cache   map[string]image.Config

	deps *deps.Deps
}

// Config returns the image.Config for the specified path relative to the
// working directory.
func (ns *Namespace) Config(path interface{}) (image.Config, error) {
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

	f, err := ns.deps.Fs.WorkingDir.Open(filename)
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

func (ns *Namespace) Filter(args ...interface{}) (resource.Image, error) {
	if len(args) < 2 {
		return nil, errors.New("must provide an image and one or more filters")
	}

	img := args[len(args)-1].(resource.Image)
	filtersv := args[:len(args)-1]

	return img.Filter(filtersv...)
}
