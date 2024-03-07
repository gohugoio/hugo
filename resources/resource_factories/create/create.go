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
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugofs/glob"
	"github.com/gohugoio/hugo/identity"

	"github.com/gohugoio/hugo/hugofs"

	"github.com/gohugoio/hugo/cache/dynacache"
	"github.com/gohugoio/hugo/cache/filecache"
	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/resources/resource"
)

// Client contains methods to create Resource objects.
// tasks to Resource objects.
type Client struct {
	rs               *resources.Spec
	httpClient       *http.Client
	cacheGetResource *filecache.Cache
}

// New creates a new Client with the given specification.
func New(rs *resources.Spec) *Client {
	return &Client{
		rs: rs,
		httpClient: &http.Client{
			Timeout: time.Minute,
		},
		cacheGetResource: rs.FileCaches.GetResourceCache(),
	}
}

// Copy copies r to the new targetPath.
func (c *Client) Copy(r resource.Resource, targetPath string) (resource.Resource, error) {
	key := dynacache.CleanKey(targetPath)
	return c.rs.ResourceCache.GetOrCreate(key, func() (resource.Resource, error) {
		return resources.Copy(r, targetPath), nil
	})
}

// Get creates a new Resource by opening the given pathname in the assets filesystem.
func (c *Client) Get(pathname string) (resource.Resource, error) {
	pathname = path.Clean(pathname)
	key := dynacache.CleanKey(pathname)

	return c.rs.ResourceCache.GetOrCreate(key, func() (resource.Resource, error) {
		// The resource file will not be read before it gets used (e.g. in .Content),
		// so we need to check that the file exists here.
		filename := filepath.FromSlash(pathname)
		fi, err := c.rs.BaseFs.Assets.Fs.Stat(filename)
		if err != nil {
			if os.IsNotExist(err) {
				return nil, nil
			}
			// A real error.
			return nil, err
		}

		pi := fi.(hugofs.FileMetaInfo).Meta().PathInfo

		return c.rs.NewResource(resources.ResourceSourceDescriptor{
			LazyPublish: true,
			OpenReadSeekCloser: func() (hugio.ReadSeekCloser, error) {
				return c.rs.BaseFs.Assets.Fs.Open(filename)
			},
			Path:          pi,
			GroupIdentity: pi,
			TargetPath:    pathname,
		})
	})
}

// Match gets the resources matching the given pattern from the assets filesystem.
func (c *Client) Match(pattern string) (resource.Resources, error) {
	return c.match("__match", pattern, nil, false)
}

func (c *Client) ByType(tp string) resource.Resources {
	res, err := c.match(path.Join("_byType", tp), "**", func(r resource.Resource) bool { return r.ResourceType() == tp }, false)
	if err != nil {
		panic(err)
	}
	return res
}

// GetMatch gets first resource matching the given pattern from the assets filesystem.
func (c *Client) GetMatch(pattern string) (resource.Resource, error) {
	res, err := c.match("__get-match", pattern, nil, true)
	if err != nil || len(res) == 0 {
		return nil, err
	}
	return res[0], err
}

func (c *Client) match(name, pattern string, matchFunc func(r resource.Resource) bool, firstOnly bool) (resource.Resources, error) {
	pattern = glob.NormalizePath(pattern)
	partitions := glob.FilterGlobParts(strings.Split(pattern, "/"))
	key := path.Join(name, path.Join(partitions...))
	key = path.Join(key, pattern)

	return c.rs.ResourceCache.GetOrCreateResources(key, func() (resource.Resources, error) {
		var res resource.Resources

		handle := func(info hugofs.FileMetaInfo) (bool, error) {
			meta := info.Meta()

			r, err := c.rs.NewResource(resources.ResourceSourceDescriptor{
				LazyPublish: true,
				OpenReadSeekCloser: func() (hugio.ReadSeekCloser, error) {
					return meta.Open()
				},
				NameNormalized: meta.PathInfo.Name(),
				NameOriginal:   meta.PathInfo.Unnormalized().Name(),
				GroupIdentity:  meta.PathInfo,
				TargetPath:     meta.PathInfo.Unnormalized().Path(),
			})
			if err != nil {
				return true, err
			}

			if matchFunc != nil && !matchFunc(r) {
				return false, nil
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
// TODO(bep) see #10912; we currently emit a warning for this config scenario.
func (c *Client) FromString(targetPath, content string) (resource.Resource, error) {
	targetPath = path.Clean(targetPath)
	key := dynacache.CleanKey(targetPath) + helpers.MD5String(content)
	r, err := c.rs.ResourceCache.GetOrCreate(key, func() (resource.Resource, error) {
		return c.rs.NewResource(
			resources.ResourceSourceDescriptor{
				LazyPublish:   true,
				GroupIdentity: identity.Anonymous, // All usage of this resource are tracked via its string content.
				OpenReadSeekCloser: func() (hugio.ReadSeekCloser, error) {
					return hugio.NewReadSeekerNoOpCloserFromString(content), nil
				},
				TargetPath: targetPath,
			})
	})

	return r, err
}
