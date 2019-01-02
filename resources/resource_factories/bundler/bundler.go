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

// Package bundler contains functions for concatenation etc. of Resource objects.
package bundler

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/resources/resource"
)

// Client contains methods perform concatenation and other bundling related
// tasks to Resource objects.
type Client struct {
	rs *resources.Spec
}

// New creates a new Client with the given specification.
func New(rs *resources.Spec) *Client {
	return &Client{rs: rs}
}

type multiReadSeekCloser struct {
	mr      io.Reader
	sources []hugio.ReadSeekCloser
}

func (r *multiReadSeekCloser) Read(p []byte) (n int, err error) {
	return r.mr.Read(p)
}

func (r *multiReadSeekCloser) Seek(offset int64, whence int) (newOffset int64, err error) {
	for _, s := range r.sources {
		newOffset, err = s.Seek(offset, whence)
		if err != nil {
			return
		}
	}
	return
}

func (r *multiReadSeekCloser) Close() error {
	for _, s := range r.sources {
		s.Close()
	}
	return nil
}

// Concat concatenates the list of Resource objects.
func (c *Client) Concat(targetPath string, r resource.Resources) (resource.Resource, error) {
	// The CACHE_OTHER will make sure this will be re-created and published on rebuilds.
	return c.rs.ResourceCache.GetOrCreate(resources.CACHE_OTHER, targetPath, func() (resource.Resource, error) {
		var resolvedm media.Type

		// The given set of resources must be of the same Media Type.
		// We may improve on that in the future, but then we need to know more.
		for i, r := range r {
			if i > 0 && r.MediaType().Type() != resolvedm.Type() {
				return nil, fmt.Errorf("resources in Concat must be of the same Media Type, got %q and %q", r.MediaType().Type(), resolvedm.Type())
			}
			resolvedm = r.MediaType()
		}

		concatr := func() (hugio.ReadSeekCloser, error) {
			var rcsources []hugio.ReadSeekCloser
			for _, s := range r {
				rcr, ok := s.(resource.ReadSeekCloserResource)
				if !ok {
					return nil, fmt.Errorf("resource %T does not implement resource.ReadSeekerCloserResource", s)
				}
				rc, err := rcr.ReadSeekCloser()
				if err != nil {
					// Close the already opened.
					for _, rcs := range rcsources {
						rcs.Close()
					}
					return nil, err
				}
				rcsources = append(rcsources, rc)
			}

			readers := make([]io.Reader, len(rcsources))
			for i := 0; i < len(rcsources); i++ {
				readers[i] = rcsources[i]
			}

			mr := io.MultiReader(readers...)

			return &multiReadSeekCloser{mr: mr, sources: rcsources}, nil
		}

		composite, err := c.rs.NewForFs(
			c.rs.FileCaches.AssetsCache().Fs,
			resources.ResourceSourceDescriptor{
				LazyPublish:        true,
				OpenReadSeekCloser: concatr,
				RelTargetFilename:  filepath.Clean(targetPath)})

		if err != nil {
			return nil, err
		}

		return composite, nil
	})

}
