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
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/helpers"
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

// Get creates a new Resource by opening the given filename in the given filesystem.
func (c *Client) Get(fs afero.Fs, filename string) (resource.Resource, error) {
	u, err := url.Parse(filename)
	if err != nil {
		return nil, err
	}

	isLocalResource := u.Scheme == ""
	if isLocalResource {
		filename = filepath.Clean(filename)
		return c.rs.ResourceCache.GetOrCreate(resources.ResourceKeyPartition(filename), filename, func() (resource.Resource, error) {
			return c.rs.NewForFs(fs,
				resources.ResourceSourceDescriptor{
					LazyPublish:    true,
					SourceFilename: filename})
		})

	}

	resourceID := helpers.MD5String(filename)
	return c.rs.ResourceCache.GetOrCreate(resources.CACHE_OTHER, resourceID, func() (resource.Resource, error) {
		res, err := http.Get(filename)
		if err != nil {
			return nil, err
		}

		if res.StatusCode < 200 || res.StatusCode > 299 {
			return nil, errors.Errorf("Failed to retrieve remote file: %s", http.StatusText(res.StatusCode))
		}

		defer res.Body.Close()
		bodyBuff, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		contentType := ""
		if arr, _ := mime.ExtensionsByType(res.Header.Get("Content-Type")); arr != nil {
			contentType = arr[0]
		}

		if contentType == "" {
			if ct := http.DetectContentType(bodyBuff); ct != "application/octet-stream" {
				if arr, _ := mime.ExtensionsByType(ct); arr != nil {
					contentType = arr[0]
				}

			}
		}

		if contentType == "" {
			// maybe there's a useful suffix in the url path?
			if ext := filepath.Ext(u.Path); ext != "" {
				contentType = ext
			}
		}

		resourceID = resourceID + contentType

		return c.rs.NewForFs(c.rs.FileCaches.AssetsCache().Fs,
			resources.ResourceSourceDescriptor{
				OpenReadSeekCloser: func() (hugio.ReadSeekCloser, error) {
					return hugio.NewReadSeekerNoOpCloserFromString(string(bodyBuff)), nil
				},
				LazyPublish:       true,
				RelTargetFilename: resourceID})
	})
}

// FromString creates a new Resource from a string with the given relative target path.
func (c *Client) FromString(targetPath, content string) (resource.Resource, error) {
	return c.rs.ResourceCache.GetOrCreate(resources.CACHE_OTHER, targetPath, func() (resource.Resource, error) {
		return c.rs.NewForFs(
			c.rs.FileCaches.AssetsCache().Fs,
			resources.ResourceSourceDescriptor{
				LazyPublish: true,
				OpenReadSeekCloser: func() (hugio.ReadSeekCloser, error) {
					return hugio.NewReadSeekerNoOpCloserFromString(content), nil
				},
				RelTargetFilename: filepath.Clean(targetPath)})

	})

}
