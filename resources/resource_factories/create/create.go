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

// Package create contains functions for to create Resource objects. This will
// typically non-files.
package create

import (
	"path"
	"path/filepath"
	"strings"

	"github.com/gohugoio/hugo/hugofs/glob"

	"github.com/gohugoio/hugo/hugofs"

	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/resources/resource"
)

// Client contains methods to create Resource objects.
// tasks to Resource objects.
type Client struct {
	rs *resources.Spec
}

// New creates a new Client with the given specification.
func New(rs *resources.Spec) *Client {
	return &Client{rs: rs}
}

// Get creates a new Resource by opening the given filename in the assets filesystem.
func (c *Client) Get(filename string) (resource.Resource, error) {
	filename = filepath.Clean(filename)
	return c.rs.ResourceCache.GetOrCreate(resources.ResourceCacheKey(filename), func() (resource.Resource, error) {
		return c.rs.New(resources.ResourceSourceDescriptor{
			Fs:             c.rs.BaseFs.Assets.Fs,
			LazyPublish:    true,
			SourceFilename: filename})
	})

}

// Match gets the resources matching the given pattern from the assets filesystem.
func (c *Client) Match(pattern string) (resource.Resources, error) {
	return c.match(pattern, false)
}

// GetMatch gets first resource matching the given pattern from the assets filesystem.
func (c *Client) GetMatch(pattern string) (resource.Resource, error) {
	res, err := c.match(pattern, true)
	if err != nil || len(res) == 0 {
		return nil, err
	}
	return res[0], err
}

func (c *Client) match(pattern string, firstOnly bool) (resource.Resources, error) {
	var name string
	if firstOnly {
		name = "__get-match"
	} else {
		name = "__match"
	}

	pattern = glob.NormalizePath(pattern)
	partitions := glob.FilterGlobParts(strings.Split(pattern, "/"))
	if len(partitions) == 0 {
		partitions = []string{resources.CACHE_OTHER}
	}
	key := path.Join(name, path.Join(partitions...))
	key = path.Join(key, pattern)

	return c.rs.ResourceCache.GetOrCreateResources(key, func() (resource.Resources, error) {
		var res resource.Resources

		handle := func(info hugofs.FileMetaInfo) (bool, error) {
			meta := info.Meta()
			r, err := c.rs.New(resources.ResourceSourceDescriptor{
				LazyPublish: true,
				FileInfo:    info,
				OpenReadSeekCloser: func() (hugio.ReadSeekCloser, error) {
					return meta.Open()
				},
				RelTargetFilename: meta.Path()})

			if err != nil {
				return true, err
			}

			res = append(res, r)

			return firstOnly, nil

		}

		if err := hugofs.Glob(c.rs.BaseFs.Assets.Fs, pattern, handle); err != nil {
			return nil, err
		}

		return res, nil

	})
}

// FromString creates a new Resource from a string with the given relative target path.
func (c *Client) FromString(targetPath, content string) (resource.Resource, error) {
	return c.rs.ResourceCache.GetOrCreate(path.Join(resources.CACHE_OTHER, targetPath), func() (resource.Resource, error) {
		return c.rs.New(
			resources.ResourceSourceDescriptor{
				Fs:          c.rs.FileCaches.AssetsCache().Fs,
				LazyPublish: true,
				OpenReadSeekCloser: func() (hugio.ReadSeekCloser, error) {
					return hugio.NewReadSeekerNoOpCloserFromString(content), nil
				},
				RelTargetFilename: filepath.Clean(targetPath)})

	})

}
