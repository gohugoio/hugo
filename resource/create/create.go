// Copyright 2018 The Hugo Authors. All rights reserved.
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
	"path/filepath"

	"github.com/spf13/afero"

	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/resource"
)

// Client contains methods to create Resource objects.
// tasks to Resource objects.
type Client struct {
	rs *resource.Spec
}

// New creates a new Client with the given specification.
func New(rs *resource.Spec) *Client {
	return &Client{rs: rs}
}

// Get creates a new Resource by opening the given filename in the given filesystem.
func (c *Client) Get(fs afero.Fs, filename string) (resource.Resource, error) {
	filename = filepath.Clean(filename)
	return c.rs.ResourceCache.GetOrCreate(resource.ResourceKeyPartition(filename), filename, func() (resource.Resource, error) {
		return c.rs.NewForFs(fs,
			resource.ResourceSourceDescriptor{
				LazyPublish:    true,
				SourceFilename: filename})
	})

}

// FromString creates a new Resource from a string with the given relative target path.
func (c *Client) FromString(targetPath, content string) (resource.Resource, error) {
	return c.rs.ResourceCache.GetOrCreate(resource.CACHE_OTHER, targetPath, func() (resource.Resource, error) {
		return c.rs.NewForFs(
			c.rs.BaseFs.Resources.Fs,
			resource.ResourceSourceDescriptor{
				LazyPublish: true,
				OpenReadSeekCloser: func() (hugio.ReadSeekCloser, error) {
					return hugio.NewReadSeekerNoOpCloserFromString(content), nil
				},
				RelTargetFilename: filepath.Clean(targetPath)})

	})

}
