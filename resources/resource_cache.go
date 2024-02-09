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

package resources

import (
	"context"
	"encoding/json"
	"io"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gohugoio/hugo/resources/resource"

	"github.com/gohugoio/hugo/cache/dynacache"
	"github.com/gohugoio/hugo/cache/filecache"
)

func newResourceCache(rs *Spec, memCache *dynacache.Cache) *ResourceCache {
	return &ResourceCache{
		fileCache: rs.FileCaches.AssetsCache(),
		cacheResource: dynacache.GetOrCreatePartition[string, resource.Resource](
			memCache,
			"/res1",
			dynacache.OptionsPartition{ClearWhen: dynacache.ClearOnChange, Weight: 40},
		),
		cacheResources: dynacache.GetOrCreatePartition[string, resource.Resources](
			memCache,
			"/ress",
			dynacache.OptionsPartition{ClearWhen: dynacache.ClearOnRebuild, Weight: 40},
		),
		cacheResourceTransformation: dynacache.GetOrCreatePartition[string, *resourceAdapterInner](
			memCache,
			"/res1/tra",
			dynacache.OptionsPartition{ClearWhen: dynacache.ClearOnChange, Weight: 40},
		),
	}
}

type ResourceCache struct {
	sync.RWMutex

	cacheResource               *dynacache.Partition[string, resource.Resource]
	cacheResources              *dynacache.Partition[string, resource.Resources]
	cacheResourceTransformation *dynacache.Partition[string, *resourceAdapterInner]

	fileCache *filecache.Cache
}

func (c *ResourceCache) cleanKey(key string) string {
	return strings.TrimPrefix(path.Clean(strings.ToLower(filepath.ToSlash(key))), "/")
}

func (c *ResourceCache) Get(ctx context.Context, key string) (resource.Resource, bool) {
	return c.cacheResource.Get(ctx, key)
}

func (c *ResourceCache) GetOrCreate(key string, f func() (resource.Resource, error)) (resource.Resource, error) {
	return c.cacheResource.GetOrCreate(key, func(key string) (resource.Resource, error) {
		return f()
	})
}

func (c *ResourceCache) GetOrCreateResources(key string, f func() (resource.Resources, error)) (resource.Resources, error) {
	return c.cacheResources.GetOrCreate(key, func(key string) (resource.Resources, error) {
		return f()
	})
}

func (c *ResourceCache) getFilenames(key string) (string, string) {
	filenameMeta := key + ".json"
	filenameContent := key + ".content"

	return filenameMeta, filenameContent
}

func (c *ResourceCache) getFromFile(key string) (filecache.ItemInfo, io.ReadCloser, transformedResourceMetadata, bool) {
	c.RLock()
	defer c.RUnlock()

	var meta transformedResourceMetadata
	filenameMeta, filenameContent := c.getFilenames(key)

	_, jsonContent, _ := c.fileCache.GetBytes(filenameMeta)
	if jsonContent == nil {
		return filecache.ItemInfo{}, nil, meta, false
	}

	if err := json.Unmarshal(jsonContent, &meta); err != nil {
		return filecache.ItemInfo{}, nil, meta, false
	}

	fi, rc, _ := c.fileCache.Get(filenameContent)

	return fi, rc, meta, rc != nil
}

// writeMeta writes the metadata to file and returns a writer for the content part.
func (c *ResourceCache) writeMeta(key string, meta transformedResourceMetadata) (filecache.ItemInfo, io.WriteCloser, error) {
	filenameMeta, filenameContent := c.getFilenames(key)
	raw, err := json.Marshal(meta)
	if err != nil {
		return filecache.ItemInfo{}, nil, err
	}

	_, fm, err := c.fileCache.WriteCloser(filenameMeta)
	if err != nil {
		return filecache.ItemInfo{}, nil, err
	}
	defer fm.Close()

	if _, err := fm.Write(raw); err != nil {
		return filecache.ItemInfo{}, nil, err
	}

	fi, fc, err := c.fileCache.WriteCloser(filenameContent)

	return fi, fc, err
}
