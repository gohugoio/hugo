// Copyright 2020 The Hugo Authors. All rights reserved.
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

// Package memcache provides the core memory cache used in Hugo.
package memcache

import (
	"context"
	"math"
	"path"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/tpl"

	"github.com/gohugoio/hugo/media"

	"github.com/gohugoio/hugo/config"

	"github.com/gohugoio/hugo/resources/resource"

	"github.com/gohugoio/hugo/helpers"

	"github.com/BurntSushi/locker"
	"github.com/karlseguin/ccache/v2"
)

const (
	ClearOnRebuild ClearWhen = iota + 1
	ClearOnChange
	ClearNever
)

const (
	cacheVirtualRoot = "_root/"
)

var (

	// Consider a change in files matching this expression a "JS change".
	isJSFileRe = regexp.MustCompile(`\.(js|ts|jsx|tsx)`)

	// Consider a change in files matching this expression a "CSS change".
	isCSSFileRe = regexp.MustCompile(`\.(css|scss|sass)`)

	// These config files are tightly related to CSS editing, so consider
	// a change to any of them a "CSS change".
	isCSSConfigRe = regexp.MustCompile(`(postcss|tailwind)\.config\.js`)
)

const unknownExtension = "unkn"

// New creates a new cache.
func New(conf Config) *Cache {
	if conf.TTL == 0 {
		conf.TTL = time.Second * 33
	}
	if conf.CheckInterval == 0 {
		conf.CheckInterval = time.Second * 2
	}
	if conf.MaxSize == 0 {
		conf.MaxSize = 100000
	}
	if conf.ItemsToPrune == 0 {
		conf.ItemsToPrune = 200
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	stats := &stats{
		memstatsStart:   m,
		maxSize:         conf.MaxSize,
		availableMemory: config.GetMemoryLimit(),
	}

	conf.MaxSize = stats.adjustMaxSize(nil)

	c := &Cache{
		conf:    conf,
		cache:   ccache.Layered(ccache.Configure().MaxSize(conf.MaxSize).ItemsToPrune(conf.ItemsToPrune)),
		getters: make(map[string]*getter),
		ttl:     conf.TTL,
		stats:   stats,
		nlocker: locker.NewLocker(),
	}

	c.stop = c.start()

	return c
}

// CleanKey turns s into a format suitable for a cache key for this package.
// The key will be a Unix-styled path without any leading slash.
// If the input string does not contain any slash, a root will be prepended.
// If the input string does not contain any ".", a dummy file suffix will be appended.
// These are to make sure that they can effectively partake in the "cache cleaning"
// strategy used in server mode.
func CleanKey(s string) string {
	s = path.Clean(helpers.ToSlashTrimLeading(s))
	if !strings.ContainsRune(s, '/') {
		s = cacheVirtualRoot + s
	}
	if !strings.ContainsRune(s, '.') {
		s += "." + unknownExtension
	}

	return s
}

// InsertKeyPathElement inserts the given element after the first '/' in key.
func InsertKeyPathElements(key string, elements ...string) string {
	slashIdx := strings.Index(key, "/")
	return key[:slashIdx] + "/" + path.Join(elements...) + key[slashIdx:]
}

// Cache configures a cache.
type Cache struct {
	mu      sync.Mutex
	getters map[string]*getter

	conf  Config
	cache *ccache.LayeredCache

	ttl     time.Duration
	nlocker *locker.Locker

	stats    *stats
	stopOnce sync.Once
	stop     func()
}

// Clear clears the cache state.
// This method is not thread safe.
func (c *Cache) Clear() {
	c.nlocker = locker.NewLocker()
	for _, g := range c.getters {
		g.c.DeleteAll(g.partition)
	}
}

// ClearOn clears all the caches given a eviction strategy and (optional) a
// change set.
// This method is not thread safe.
func (c *Cache) ClearOn(when ClearWhen, changeset ...identity.Identity) {
	if when == 0 {
		panic("invalid ClearWhen")
	}

	// Fist pass.
	for _, g := range c.getters {
		if g.clearWhen == ClearNever {
			continue
		}

		if g.clearWhen == when {
			// Clear all.
			g.Clear()
			continue
		}

		shouldDelete := func(key string, e Entry) bool {
			// We always clear elements marked as stale.
			if resource.IsStaleAny(e, e.Value) {
				return true
			}

			if e.ClearWhen == ClearNever {
				return false
			}

			if e.ClearWhen == when && e.ClearWhen == ClearOnRebuild {
				return true
			}

			// Now check if this entry has changed based on the changeset
			// based on filesystem events.

			if len(changeset) == 0 {
				// Nothing changed.
				return false
			}

			var notNotDependent bool
			identity.WalkIdentities(e.Value, func(id2 identity.Identity) bool {
				for _, id := range changeset {
					if !identity.IsNotDependent(id2, id) {
						// It's probably dependent, evict from cache.
						notNotDependent = true
						return true
					}
				}
				return false
			})

			return notNotDependent
		}

		// Two passes, the last one to catch any leftover values marked stale in the first.
		g.c.cache.DeleteFunc(g.partition, func(key string, item *ccache.Item) bool {
			e := item.Value().(Entry)
			if shouldDelete(key, e) {
				resource.MarkStale(e.Value)
				return true
			}
			return false
		})

	}

	// Second pass: Clear all entries marked as stale in the first.
	for _, g := range c.getters {
		if g.clearWhen == ClearNever || g.clearWhen == when {
			continue
		}

		g.c.cache.DeleteFunc(g.partition, func(key string, item *ccache.Item) bool {
			e := item.Value().(Entry)
			return resource.IsStaleAny(e, e.Value)
		})
	}
}

type resourceTP interface {
	ResourceTarget() resource.Resource
}

/*
	assets: css/mystyles.scss
	content: blog/mybundle/data.json
*/
func shouldEvict(key string, e Entry, when ClearWhen, changeset ...identity.PathIdentity) bool {
	return false
}

func (c *Cache) DeleteAll(primary string) bool {
	return c.cache.DeleteAll(primary)
}

func (c *Cache) GetDropped() int {
	return c.cache.GetDropped()
}

func (c *Cache) GetOrCreatePartition(partition string, clearWhen ClearWhen) Getter {
	if clearWhen == 0 {
		panic("GetOrCreatePartition: invalid ClearWhen")
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	g, found := c.getters[partition]
	if found {
		if g.clearWhen != clearWhen {
			panic("GetOrCreatePartition called with the same partition but different clearing strategy.")
		}
		return g
	}

	g = &getter{
		partition: partition,
		c:         c,
		clearWhen: clearWhen,
	}

	c.getters[partition] = g

	return g
}

func (c *Cache) Stop() {
	c.stopOnce.Do(func() {
		c.stop()
		c.cache.Stop()
	})
}

func (c *Cache) start() func() {
	ticker := time.NewTicker(c.conf.CheckInterval)
	quit := make(chan struct{})

	checkAndAdjustMaxSize := func() {
		var m runtime.MemStats
		cacheDropped := c.GetDropped()
		c.stats.decr(cacheDropped)

		runtime.ReadMemStats(&m)
		c.stats.memstatsCurrent = m
		c.stats.adjustMaxSize(c.cache.SetMaxSize)

		// fmt.Printf("\n\nAlloc = %v\nTotalAlloc = %v\nSys = %v\nNumGC = %v\nMemCacheDropped = %d\n\n", helpers.FormatByteCount(m.Alloc), helpers.FormatByteCount(m.TotalAlloc), helpers.FormatByteCount(m.Sys), m.NumGC, cacheDropped)
	}
	go func() {
		for {
			select {
			case <-ticker.C:
				checkAndAdjustMaxSize()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	return func() {
		close(quit)
	}
}

// get tries to get the value with the given cache paths.
// It returns nil if not found
func (c *Cache) get(primary, secondary string) (interface{}, error) {
	if v := c.cache.Get(primary, secondary); v != nil {
		e := v.Value().(Entry)
		if !resource.IsStaleAny(e, e.Value) {
			return e.Value, e.Err
		}
	}
	return nil, nil
}

// getOrCreate tries to get the value with the given cache paths, if not found
// create will be called and the result cached.
//
// This method is thread safe.
func (c *Cache) getOrCreate(primary, secondary string, create func() Entry) (interface{}, error) {
	if v, err := c.get(primary, secondary); v != nil || err != nil {
		return v, err
	}

	// The provided create function may be a relatively time consuming operation,
	// and there will in the commmon case be concurrent requests for the same key'd
	// resource, so make sure we pause these until the result is ready.
	path := primary + secondary
	c.nlocker.Lock(path)
	defer c.nlocker.Unlock(path)

	// Try again.
	if v := c.cache.Get(primary, secondary); v != nil {
		e := v.Value().(Entry)
		if !resource.IsStaleAny(e, e.Value) {
			return e.Value, e.Err
		}
	}

	// Create it and store it in cache.
	entry := create()

	if entry.Err != nil {
		entry.ClearWhen = ClearOnRebuild
	} else if entry.ClearWhen == 0 {
		panic("entry: invalid ClearWhen")
	}

	entry.size = 1 // For now.

	c.cache.Set(primary, secondary, entry, c.ttl)
	c.stats.incr(1)

	return entry.Value, entry.Err
}

func (c *Cache) trackDependencyIfRunning(ctx context.Context, v interface{}) {
	if !c.conf.Running {
		return
	}

	tpl.AddIdentiesToDataContext(ctx, v)
}

type ClearWhen int

type Config struct {
	CheckInterval time.Duration
	MaxSize       int64
	ItemsToPrune  uint32
	TTL           time.Duration
	Running       bool
}

type Entry struct {
	Value     interface{}
	size      int64
	Err       error
	StaleFunc func() bool
	ClearWhen
}

func (e Entry) Size() int64 {
	return e.size
}

func (e Entry) IsStale() bool {
	return e.StaleFunc != nil && e.StaleFunc()
}

type Getter interface {
	Clear()
	Get(ctx context.Context, path string) (interface{}, error)
	GetOrCreate(ctx context.Context, path string, create func() Entry) (interface{}, error)
}

type getter struct {
	c         *Cache
	partition string

	clearWhen ClearWhen
}

func (g *getter) Clear() {
	g.c.DeleteAll(g.partition)
}

func (g *getter) Get(ctx context.Context, path string) (interface{}, error) {
	v, err := g.c.get(g.partition, path)
	if err != nil {
		return nil, err
	}

	g.c.trackDependencyIfRunning(ctx, v)

	return v, nil
}

func (g *getter) GetOrCreate(ctx context.Context, path string, create func() Entry) (interface{}, error) {
	v, err := g.c.getOrCreate(g.partition, path, create)
	if err != nil {
		return nil, err
	}

	g.c.trackDependencyIfRunning(ctx, v)

	return v, nil
}

type stats struct {
	memstatsStart   runtime.MemStats
	memstatsCurrent runtime.MemStats
	maxSize         int64
	availableMemory uint64
	numItems        uint64
}

func (s *stats) getNumItems() uint64 {
	return atomic.LoadUint64(&s.numItems)
}

func (s *stats) adjustMaxSize(setter func(size int64)) int64 {
	newSize := int64(float64(s.maxSize) * s.resizeFactor())
	if newSize != s.maxSize && setter != nil {
		setter(newSize)
	}
	return newSize
}

func (s *stats) decr(i int) {
	atomic.AddUint64(&s.numItems, ^uint64(i-1))
}

func (s *stats) incr(i int) {
	atomic.AddUint64(&s.numItems, uint64(i))
}

func (s *stats) resizeFactor() float64 {
	if s.memstatsCurrent.Alloc == 0 {
		return 1.0
	}
	return math.Floor(float64(s.availableMemory/s.memstatsCurrent.Alloc)*10) / 10
}

// Helpers to help eviction of related media types.
func isCSSType(m media.Type) bool {
	tp := m.Type()
	return tp == media.CSSType.Type() || tp == media.SASSType.Type() || tp == media.SCSSType.Type()
}

func isJSType(m media.Type) bool {
	tp := m.Type()
	return tp == media.JavascriptType.Type() || tp == media.TypeScriptType.Type() || tp == media.JSXType.Type() || tp == media.TSXType.Type()
}

func keyValid(s string) bool {
	if len(s) < 5 {
		return false
	}
	if strings.ContainsRune(s, '\\') {
		return false
	}
	if strings.HasPrefix(s, "/") {
		return false
	}
	if !strings.ContainsRune(s, '/') {
		return false
	}

	dotIdx := strings.Index(s, ".")
	if dotIdx == -1 || dotIdx == len(s)-1 {
		return false
	}

	return true
}

// This assumes a valid key path.
func splitBasePathAndExt(path string) (string, string) {
	dotIdx := strings.LastIndex(path, ".")
	ext := path[dotIdx+1:]
	slashIdx := strings.Index(path, "/")

	return path[:slashIdx], ext
}
