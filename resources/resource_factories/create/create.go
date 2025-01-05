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

	"github.com/bep/logg"
	"github.com/gohugoio/httpcache"
	hhttpcache "github.com/gohugoio/hugo/cache/httpcache"
	"github.com/gohugoio/hugo/hugofs/glob"
	"github.com/gohugoio/hugo/identity"

	"github.com/gohugoio/hugo/hugofs"

	"github.com/gohugoio/hugo/cache/dynacache"
	"github.com/gohugoio/hugo/cache/filecache"
	"github.com/gohugoio/hugo/common/hashing"
	"github.com/gohugoio/hugo/common/hcontext"
	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/common/tasks"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/resources/resource"
)

// Client contains methods to create Resource objects.
// tasks to Resource objects.
type Client struct {
	rs                   *resources.Spec
	httpClient           *http.Client
	httpCacheConfig      hhttpcache.ConfigCompiled
	cacheGetResource     *filecache.Cache
	resourceIDDispatcher hcontext.ContextDispatcher[string]

	// Set when watching.
	remoteResourceChecker *tasks.RunEvery
	remoteResourceLogger  logg.LevelLogger
}

type contextKey string

// New creates a new Client with the given specification.
func New(rs *resources.Spec) *Client {
	fileCache := rs.FileCaches.GetResourceCache()
	resourceIDDispatcher := hcontext.NewContextDispatcher[string](contextKey("resourceID"))
	httpCacheConfig := rs.Cfg.GetConfigSection("httpCacheCompiled").(hhttpcache.ConfigCompiled)
	var remoteResourceChecker *tasks.RunEvery
	if rs.Cfg.Watching() && !httpCacheConfig.IsPollingDisabled() {
		remoteResourceChecker = &tasks.RunEvery{
			HandleError: func(name string, err error) {
				rs.Logger.Warnf("Failed to check remote resource: %s", err)
			},
			RunImmediately: false,
		}

		if err := remoteResourceChecker.Start(); err != nil {
			panic(err)
		}

		rs.BuildClosers.Add(remoteResourceChecker)
	}

	httpTimeout := 2 * time.Minute // Need to cover retries.
	if httpTimeout < (rs.Cfg.Timeout() + 30*time.Second) {
		httpTimeout = rs.Cfg.Timeout() + 30*time.Second
	}

	return &Client{
		rs:                    rs,
		httpCacheConfig:       httpCacheConfig,
		resourceIDDispatcher:  resourceIDDispatcher,
		remoteResourceChecker: remoteResourceChecker,
		remoteResourceLogger:  rs.Logger.InfoCommand("remote"),
		httpClient: &http.Client{
			Timeout: httpTimeout,
			Transport: &httpcache.Transport{
				Cache: fileCache.AsHTTPCache(),
				CacheKey: func(req *http.Request) string {
					return resourceIDDispatcher.Get(req.Context())
				},
				Around: func(req *http.Request, key string) func() {
					return fileCache.NamedLock(key)
				},
				AlwaysUseCachedResponse: func(req *http.Request, key string) bool {
					return !httpCacheConfig.For(req.URL.String())
				},
				ShouldCache: func(req *http.Request, resp *http.Response, key string) bool {
					return shouldCache(resp.StatusCode)
				},
				MarkCachedResponses: true,
				EnableETagPair:      true,
				Transport: &transport{
					Cfg:    rs.Cfg,
					Logger: rs.Logger,
				},
			},
		},
		cacheGetResource: fileCache,
	}
}

// Copy copies r to the new targetPath.
func (c *Client) Copy(r resource.Resource, targetPath string) (resource.Resource, error) {
	key := dynacache.CleanKey(targetPath) + "__copy"
	return c.rs.ResourceCache.GetOrCreate(key, func() (resource.Resource, error) {
		return resources.Copy(r, targetPath), nil
	})
}

// Get creates a new Resource by opening the given pathname in the assets filesystem.
func (c *Client) Get(pathname string) (resource.Resource, error) {
	pathname = path.Clean(pathname)
	key := dynacache.CleanKey(pathname) + "__get"

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

		return c.getOrCreateFileResource(fi.(hugofs.FileMetaInfo))
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

func (c *Client) getOrCreateFileResource(info hugofs.FileMetaInfo) (resource.Resource, error) {
	meta := info.Meta()
	return c.rs.ResourceCache.GetOrCreateFile(filepath.ToSlash(meta.Filename), func() (resource.Resource, error) {
		return c.rs.NewResource(resources.ResourceSourceDescriptor{
			LazyPublish: true,
			OpenReadSeekCloser: func() (hugio.ReadSeekCloser, error) {
				return meta.Open()
			},
			NameNormalized:       meta.PathInfo.Path(),
			NameOriginal:         meta.PathInfo.Unnormalized().Path(),
			GroupIdentity:        meta.PathInfo,
			TargetPath:           meta.PathInfo.Unnormalized().Path(),
			SourceFilenameOrPath: meta.Filename,
		})
	})
}

func (c *Client) match(name, pattern string, matchFunc func(r resource.Resource) bool, firstOnly bool) (resource.Resources, error) {
	pattern = glob.NormalizePath(pattern)
	partitions := glob.FilterGlobParts(strings.Split(pattern, "/"))
	key := path.Join(name, path.Join(partitions...))
	key = path.Join(key, pattern)

	return c.rs.ResourceCache.GetOrCreateResources(key, func() (resource.Resources, error) {
		var res resource.Resources

		handle := func(info hugofs.FileMetaInfo) (bool, error) {
			r, err := c.getOrCreateFileResource(info)
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

type Options struct {
	// The target path relative to the publish directory.
	// Unix style path, i.e. "images/logo.png".
	TargetPath string

	// Whether the TargetPath has a hash in it which will change if the resource changes.
	// If not, we will calculate a hash from the content.
	TargetPathHasHash bool

	// The content to create the Resource from.
	CreateContent func() (func() (hugio.ReadSeekCloser, error), error)
}

// FromOpts creates a new Resource from the given Options.
// Make sure to set optis.TargetPathHasHash if the TargetPath already contains a hash,
// as this avoids the need to calculate it.
// To create a new ReadSeekCloser from a string, use hugio.NewReadSeekerNoOpCloserFromString,
// or hugio.NewReadSeekerNoOpCloserFromBytes for a byte slice.
// See FromString.
func (c *Client) FromOpts(opts Options) (resource.Resource, error) {
	opts.TargetPath = path.Clean(opts.TargetPath)
	var hash string
	var newReadSeeker func() (hugio.ReadSeekCloser, error) = nil
	if !opts.TargetPathHasHash {
		var err error
		newReadSeeker, err = opts.CreateContent()
		if err != nil {
			return nil, err
		}
		if err := func() error {
			r, err := newReadSeeker()
			if err != nil {
				return err
			}
			defer r.Close()

			hash, err = hashing.XxHashFromReaderHexEncoded(r)
			if err != nil {
				return err
			}
			return nil
		}(); err != nil {
			return nil, err
		}
	}

	key := dynacache.CleanKey(opts.TargetPath) + hash
	r, err := c.rs.ResourceCache.GetOrCreate(key, func() (resource.Resource, error) {
		if newReadSeeker == nil {
			var err error
			newReadSeeker, err = opts.CreateContent()
			if err != nil {
				return nil, err
			}
		}
		return c.rs.NewResource(
			resources.ResourceSourceDescriptor{
				LazyPublish:   true,
				GroupIdentity: identity.Anonymous, // All usage of this resource are tracked via its string content.
				OpenReadSeekCloser: func() (hugio.ReadSeekCloser, error) {
					return newReadSeeker()
				},
				TargetPath: opts.TargetPath,
			})
	})

	return r, err
}

// FromString creates a new Resource from a string with the given relative target path.
func (c *Client) FromString(targetPath, content string) (resource.Resource, error) {
	return c.FromOpts(Options{
		TargetPath: targetPath,
		CreateContent: func() (func() (hugio.ReadSeekCloser, error), error) {
			return func() (hugio.ReadSeekCloser, error) {
				return hugio.NewReadSeekerNoOpCloserFromString(content), nil
			}, nil
		},
	})
}
