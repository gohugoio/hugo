// Copyright 2024 The Hugo Authors. All rights reserved.
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

package dynacache

import (
	"context"
	"fmt"
	"math"
	"path"
	"regexp"
	"runtime"
	"sync"
	"time"

	"github.com/bep/lazycache"
	"github.com/bep/logg"
	"github.com/gohugoio/hugo/common/collections"
	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/common/rungroup"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/resources/resource"
)

const (
	minMaxSize          = 1
	maxEntriesPerWeight = 10000
)

type KeyIdentity struct {
	Key      any
	Identity identity.Identity
}

// New creates a new cache.
func New(opts Options) *Cache {
	if opts.CheckInterval == 0 {
		opts.CheckInterval = time.Second
	}

	if opts.Log == nil {
		panic("nil Log")
	}

	infol := opts.Log.InfoCommand("dynacache")

	evictedIdentities := collections.NewStackThreadSafe[KeyIdentity]()

	onEvict := func(k, v any) {
		if !opts.Watching {
			return
		}
		identity.WalkIdentitiesShallow(v, func(level int, id identity.Identity) bool {
			evictedIdentities.Push(KeyIdentity{Key: k, Identity: id})
			return false
		})
		resource.MarkStale(v)
	}

	c := &Cache{
		partitions:        make(map[string]PartitionManager),
		onEvict:           onEvict,
		evictedIdentities: evictedIdentities,
		opts:              opts,
		availableMemory:   config.GetMemoryLimit(),
		infol:             infol,
	}

	c.stop = c.start()

	return c
}

// Options for the cache.
type Options struct {
	Log           loggers.Logger
	CheckInterval time.Duration
	Watching      bool
}

// Options for a partition.
type OptionsPartition struct {
	// When to clear this partition.
	ClearWhen ClearWhen

	// Weight is a number between 1 and 100 that indicates how important/large this partition may get.
	Weight int
}

// A dynamic partitioned cache.
type Cache struct {
	mu sync.RWMutex

	partitions map[string]PartitionManager

	onEvict           func(k, v any)
	evictedIdentities *collections.StackThreadSafe[KeyIdentity]

	opts            Options
	availableMemory uint64
	infol           logg.LevelLogger

	stopOnce sync.Once
	stop     func()
}

// DrainEvictedIdentities drains the evicted identities from the cache.
func (c *Cache) DrainEvictedIdentities() []KeyIdentity {
	return c.evictedIdentities.Drain()
}

// DrainEvictedIdentitiesMatching drains the evicted identities from the cache that match the given predicate.
func (c *Cache) DrainEvictedIdentitiesMatching(predicate func(KeyIdentity) bool) []KeyIdentity {
	return c.evictedIdentities.DrainMatching(predicate)
}

// ClearMatching clears all partition for which the predicate returns true.
func (c *Cache) ClearMatching(predicatePartition func(k string, p PartitionManager) bool, predicateValue func(k, v any) bool) {
	if predicatePartition == nil {
		predicatePartition = func(k string, p PartitionManager) bool { return true }
	}
	if predicateValue == nil {
		panic("nil predicateValue")
	}
	g := rungroup.Run[PartitionManager](context.Background(), rungroup.Config[PartitionManager]{
		NumWorkers: len(c.partitions),
		Handle: func(ctx context.Context, partition PartitionManager) error {
			partition.clearMatching(predicateValue)
			return nil
		},
	})

	for k, p := range c.partitions {
		if !predicatePartition(k, p) {
			continue
		}
		g.Enqueue(p)
	}

	g.Wait()
}

// ClearOnRebuild prepares the cache for a new rebuild taking the given changeset into account.
// predicate is optional and will clear any entry for which it returns true.
func (c *Cache) ClearOnRebuild(predicate func(k, v any) bool, changeset ...identity.Identity) {
	g := rungroup.Run[PartitionManager](context.Background(), rungroup.Config[PartitionManager]{
		NumWorkers: len(c.partitions),
		Handle: func(ctx context.Context, partition PartitionManager) error {
			partition.clearOnRebuild(predicate, changeset...)
			return nil
		},
	})

	for _, p := range c.partitions {
		g.Enqueue(p)
	}

	g.Wait()

	// Clear any entries marked as stale above.
	g = rungroup.Run[PartitionManager](context.Background(), rungroup.Config[PartitionManager]{
		NumWorkers: len(c.partitions),
		Handle: func(ctx context.Context, partition PartitionManager) error {
			partition.clearStale()
			return nil
		},
	})

	for _, p := range c.partitions {
		g.Enqueue(p)
	}

	g.Wait()
}

type keysProvider interface {
	Keys() []string
}

// Keys returns a list of keys in all partitions.
func (c *Cache) Keys(predicate func(s string) bool) []string {
	if predicate == nil {
		predicate = func(s string) bool { return true }
	}
	var keys []string
	for pn, g := range c.partitions {
		pkeys := g.(keysProvider).Keys()
		for _, k := range pkeys {
			p := path.Join(pn, k)
			if predicate(p) {
				keys = append(keys, p)
			}
		}

	}
	return keys
}

// Stop stops the cache.
func (c *Cache) Stop() {
	c.stopOnce.Do(func() {
		c.stop()
	})
}

func (c *Cache) checkMemory() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	c.adjustMaxSize(m)
}

func (c *Cache) adjustMaxSize(m runtime.MemStats) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.partitions) == 0 {
		return
	}

	limit := c.availableMemory
	// Hysteresis: shrink above limit, grow below 80% of limit, hold steady in between.
	growThreshold := limit * 4 / 5

	evicted := 0
	if m.Alloc > limit {
		for _, p := range c.partitions {
			evicted += p.shrink()
		}
	} else if m.Alloc < growThreshold {
		for _, p := range c.partitions {
			p.grow()
		}
	}

	if evicted > 0 {
		c.infol.
			WithFields(
				logg.Fields{
					{Name: "evicted", Value: evicted},
					{Name: "numGC", Value: m.NumGC},
					{Name: "limit", Value: helpers.FormatByteCount(limit)},
					{Name: "alloc", Value: helpers.FormatByteCount(m.Alloc)},
					{Name: "totalAlloc", Value: helpers.FormatByteCount(m.TotalAlloc)},
				},
			).Logf("adjusted partitions' max size")
	}
}

func (c *Cache) start() func() {
	ticker := time.NewTicker(c.opts.CheckInterval)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				c.checkMemory()
				// Reset the ticker to avoid drift.
				ticker.Reset(c.opts.CheckInterval)
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

var partitionNameRe = regexp.MustCompile(`^\/[a-zA-Z0-9]{4}(\/[a-zA-Z0-9]+)?(\/[a-zA-Z0-9]+)?`)

// GetOrCreatePartition gets or creates a partition with the given name.
func GetOrCreatePartition[K comparable, V any](c *Cache, name string, opts OptionsPartition) *Partition[K, V] {
	if c == nil {
		panic("nil Cache")
	}
	if opts.Weight < 1 || opts.Weight > 100 {
		panic("invalid Weight, must be between 1 and 100")
	}

	if !partitionNameRe.MatchString(name) {
		panic(fmt.Sprintf("invalid partition name %q", name))
	}

	c.mu.RLock()
	p, found := c.partitions[name]
	c.mu.RUnlock()
	if found {
		return p.(*Partition[K, V])
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Double check.
	p, found = c.partitions[name]
	if found {
		return p.(*Partition[K, V])
	}

	defaultMaxSize := opts.Weight * maxEntriesPerWeight

	onEvict := func(k K, v V) {
		c.onEvict(k, v)
	}

	// Create a new partition and cache it.
	partition := &Partition[K, V]{
		c:              lazycache.New(lazycache.Options[K, V]{MaxEntries: defaultMaxSize, OnEvict: onEvict}),
		defaultMaxSize: defaultMaxSize,
		maxSize:        defaultMaxSize,
		trace:          c.opts.Log.Logger().WithLevel(logg.LevelTrace).WithField("partition", name),
		opts:           opts,
	}

	c.partitions[name] = partition

	return partition
}

// Partition is a partition in the cache.
type Partition[K comparable, V any] struct {
	c *lazycache.Cache[K, V]

	zero V

	trace logg.LevelLogger
	opts  OptionsPartition

	defaultMaxSize int
	maxSize        int
}

// GetOrCreate gets or creates a value for the given key.
func (p *Partition[K, V]) GetOrCreate(key K, create func(key K) (V, error)) (V, error) {
	v, err := p.doGetOrCreate(key, create)
	if err != nil {
		return p.zero, err
	}
	if resource.StaleVersion(v) > 0 {
		p.c.Delete(key)
		return p.doGetOrCreate(key, create)
	}
	return v, err
}

func (p *Partition[K, V]) doGetOrCreate(key K, create func(key K) (V, error)) (V, error) {
	v, _, err := p.c.GetOrCreate(key, create)
	return v, err
}

func (p *Partition[K, V]) GetOrCreateWitTimeout(key K, duration time.Duration, create func(key K) (V, error)) (V, error) {
	v, err := p.doGetOrCreateWitTimeout(key, duration, create)
	if err != nil {
		return p.zero, err
	}
	if resource.StaleVersion(v) > 0 {
		p.c.Delete(key)
		return p.doGetOrCreateWitTimeout(key, duration, create)
	}
	return v, err
}

// GetOrCreateWitTimeout gets or creates a value for the given key and times out if the create function
// takes too long.
func (p *Partition[K, V]) doGetOrCreateWitTimeout(key K, duration time.Duration, create func(key K) (V, error)) (V, error) {
	resultch := make(chan V, 1)
	errch := make(chan error, 1)

	go func() {
		var (
			v   V
			err error
		)
		defer func() {
			if r := recover(); r != nil {
				if rerr, ok := r.(error); ok {
					err = rerr
				} else {
					err = fmt.Errorf("panic: %v", r)
				}
			}
			if err != nil {
				errch <- err
			} else {
				resultch <- v
			}
		}()
		v, _, err = p.c.GetOrCreate(key, create)
	}()

	select {
	case v := <-resultch:
		return v, nil
	case err := <-errch:
		return p.zero, err
	case <-time.After(duration):
		return p.zero, &herrors.TimeoutError{
			Duration: duration,
		}
	}
}

func (p *Partition[K, V]) clearMatching(predicate func(k, v any) bool) {
	p.c.DeleteFunc(func(key K, v V) bool {
		if predicate(key, v) {
			p.trace.Log(
				logg.StringFunc(
					func() string {
						return fmt.Sprintf("clearing cache key %v", key)
					},
				),
			)
			return true
		}
		return false
	})
}

func (p *Partition[K, V]) clearOnRebuild(predicate func(k, v any) bool, changeset ...identity.Identity) {
	if predicate == nil {
		predicate = func(k, v any) bool {
			return false
		}
	}
	if p.opts.ClearWhen == ClearNever {
		return
	}

	if p.opts.ClearWhen == ClearOnRebuild {
		// Clear all.
		p.Clear()
		return
	}

	depsFinder := identity.NewFinder(identity.FinderConfig{})

	shouldDelete := func(key K, v V) bool {
		// We always clear elements marked as stale.
		if resource.StaleVersion(v) > 0 {
			return true
		}

		// Now check if this entry has changed based on the changeset
		// based on filesystem events.
		if len(changeset) == 0 {
			// Nothing changed.
			return false
		}

		var probablyDependent bool
		identity.WalkIdentitiesShallow(v, func(level int, id2 identity.Identity) bool {
			for _, id := range changeset {
				if r := depsFinder.Contains(id, id2, -1); r > 0 {
					// It's probably dependent, evict from cache.
					probablyDependent = true
					return true
				}
			}
			return false
		})

		return probablyDependent
	}

	// First pass.
	// Second pass needs to be done in a separate loop to catch any
	// elements marked as stale in the other partitions.
	p.c.DeleteFunc(func(key K, v V) bool {
		if predicate(key, v) || shouldDelete(key, v) {
			p.trace.Log(
				logg.StringFunc(
					func() string {
						return fmt.Sprintf("first pass: clearing cache key %v", key)
					},
				),
			)
			return true
		}
		return false
	})
}

func (p *Partition[K, V]) Keys() []K {
	var keys []K
	p.c.DeleteFunc(func(key K, v V) bool {
		keys = append(keys, key)
		return false
	})
	return keys
}

func (p *Partition[K, V]) clearStale() {
	p.c.DeleteFunc(func(key K, v V) bool {
		staleVersion := resource.StaleVersion(v)
		if staleVersion > 0 {
			p.trace.Log(
				logg.StringFunc(
					func() string {
						return fmt.Sprintf("second pass: clearing cache key %v", key)
					},
				),
			)
		}

		return staleVersion > 0
	})
}

func (p *Partition[K, V]) shrink() int {
	n := p.c.Len()
	if n == 0 {
		return 0
	}
	// Weight in [1, 100] maps to retain factor in [0.255, 0.75].
	retain := 0.25 + 0.5*(float64(p.opts.Weight)/100.0)
	base := min(p.maxSize, n)
	newMaxSize := max(minMaxSize, int(math.Floor(float64(base)*retain)))
	return p.adjustMaxSize(newMaxSize)
}

func (p *Partition[K, V]) grow() {
	if p.maxSize >= p.defaultMaxSize {
		return
	}
	newMaxSize := min(p.defaultMaxSize, int(math.Ceil(float64(p.maxSize)*1.5)))
	p.adjustMaxSize(newMaxSize)
}

// adjustMaxSize adjusts the max size of the partition and returns the number of items evicted.
func (p *Partition[K, V]) adjustMaxSize(newMaxSize int) int {
	if newMaxSize < minMaxSize {
		newMaxSize = minMaxSize
	}
	if newMaxSize == p.maxSize {
		return 0
	}
	p.maxSize = newMaxSize
	return p.c.Resize(newMaxSize)
}

func (p *Partition[K, V]) getMaxSize() int {
	return p.maxSize
}

func (p *Partition[K, V]) Clear() {
	p.c.DeleteFunc(func(key K, v V) bool {
		return true
	})
}

func (p *Partition[K, V]) Get(ctx context.Context, key K) (V, bool) {
	return p.c.Get(key)
}

type PartitionManager interface {
	shrink() int
	grow()
	getMaxSize() int
	clearOnRebuild(predicate func(k, v any) bool, changeset ...identity.Identity)
	clearMatching(predicate func(k, v any) bool)
	clearStale()
}

const (
	ClearOnRebuild ClearWhen = iota + 1
	ClearOnChange
	ClearNever
)

type ClearWhen int

// CleanKey turns s into a format suitable for a cache key for this package.
// The key will be a Unix-styled path with a leading slash but no trailing slash.
func CleanKey(s string) string {
	return path.Clean(paths.ToSlashPreserveLeading(s))
}
