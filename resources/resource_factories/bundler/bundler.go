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
	"path"
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

func toReaders(sources []hugio.ReadSeekCloser) []io.Reader {
	readers := make([]io.Reader, len(sources))
	for i, r := range sources {
		readers[i] = r
	}
	return readers
}

func newMultiReadSeekCloser(sources ...hugio.ReadSeekCloser) *multiReadSeekCloser {
	mr := io.MultiReader(toReaders(sources)...)
	return &multiReadSeekCloser{mr, sources}
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

	r.mr = io.MultiReader(toReaders(r.sources)...)

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
	return c.rs.ResourceCache.GetOrCreate(path.Join(resources.CACHE_OTHER, targetPath), func() (resource.Resource, error) {
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

			// Arbitrary JavaScript files require a barrier between them to be safely concatenated together.
			// Without this, the last line of one file can affect the first line of the next file and change how both files are interpreted.
			if resolvedm.MainType == media.JavascriptType.MainType && resolvedm.SubType == media.JavascriptType.SubType {
				readers := make([]hugio.ReadSeekCloser, 2*len(rcsources)-1)
				j := 0
				for i := 0; i < len(rcsources); i++ {
					if i > 0 {
						readers[j] = hugio.NewReadSeekerNoOpCloserFromString("\n;\n")
						j++
					}
					readers[j] = rcsources[i]
					j++
				}
				return newMultiReadSeekCloser(readers...), nil
			}

			return newMultiReadSeekCloser(rcsources...), nil

		}

		composite, err := c.rs.New(
			resources.ResourceSourceDescriptor{
				Fs:                 c.rs.FileCaches.AssetsCache().Fs,
				LazyPublish:        true,
				OpenReadSeekCloser: concatr,
				RelTargetFilename:  filepath.Clean(targetPath)})

		if err != nil {
			return nil, err
		}

		return composite, nil
	})

}
