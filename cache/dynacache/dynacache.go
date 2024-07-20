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

const minMaxSize = 10

type KeyIdentity struct {
	Key      any
	Identity identity.Identity
}

// New creates a new cache.
func New(opts Options) *Cache {
	if opts.CheckInterval == 0 {
		opts.CheckInterval = time.Second * 2
	}

	if opts.MaxSize == 0 {
		opts.MaxSize = 100000
	}
	if opts.Log == nil {
		panic("nil Log")
	}

	if opts.MinMaxSize == 0 {
		opts.MinMaxSize = 30
	}

	stats := &stats{
		opts:             opts,
		adjustmentFactor: 1.0,
		currentMaxSize:   opts.MaxSize,
		availableMemory:  config.GetMemoryLimit(),
	}

	infol := opts.Log.InfoCommand("dynacache")

	evictedIdentities := collections.NewStack[KeyIdentity]()

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
		stats:             stats,
		infol:             infol,
	}

	c.stop = c.start()

	return c
}

// Options for the cache.
type Options struct {
	Log           loggers.Logger
	CheckInterval time.Duration
	MaxSize       int
	MinMaxSize    int
	Watching      bool
}

// Options for a partition.
type OptionsPartition struct {
	// When to clear the this partition.
	ClearWhen ClearWhen

	// Weight is a number between 1 and 100 that indicates how, in general, how big this partition may get.
	Weight int
}

func (o OptionsPartition) WeightFraction() float64 {
	return float64(o.Weight) / 100
}

func (o OptionsPartition) CalculateMaxSize(maxSizePerPartition int) int {
	return int(math.Floor(float64(maxSizePerPartition) * o.WeightFraction()))
}

// A dynamic partitioned cache.
type Cache struct {
	mu sync.RWMutex

	partitions map[string]PartitionManager

	onEvict           func(k, v any)
	evictedIdentities *collections.Stack[KeyIdentity]

	opts  Options
	infol logg.LevelLogger

	stats    *stats
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
func (c *Cache) ClearOnRebuild(changeset ...identity.Identity) {
	g := rungroup.Run[PartitionManager](context.Background(), rungroup.Config[PartitionManager]{
		NumWorkers: len(c.partitions),
		Handle: func(ctx context.Context, partition PartitionManager) error {
			partition.clearOnRebuild(changeset...)
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

func calculateMaxSizePerPartition(maxItemsTotal, totalWeightQuantity, numPartitions int) int {
	if numPartitions == 0 {
		panic("numPartitions must be > 0")
	}
	if totalWeightQuantity == 0 {
		panic("totalWeightQuantity must be > 0")
	}

	avgWeight := float64(totalWeightQuantity) / float64(numPartitions)
	return int(math.Floor(float64(maxItemsTotal) / float64(numPartitions) * (100.0 / avgWeight)))
}

// Stop stops the cache.
func (c *Cache) Stop() {
	c.stopOnce.Do(func() {
		c.stop()
	})
}

func (c *Cache) adjustCurrentMaxSize() {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.partitions) == 0 {
		return
	}
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	s := c.stats
	s.memstatsCurrent = m
	// fmt.Printf("\n\nAvailable = %v\nAlloc = %v\nTotalAlloc = %v\nSys = %v\nNumGC = %v\nMaxSize = %d\nAdjustmentFactor=%f\n\n", helpers.FormatByteCount(s.availableMemory), helpers.FormatByteCount(m.Alloc), helpers.FormatByteCount(m.TotalAlloc), helpers.FormatByteCount(m.Sys), m.NumGC, c.stats.currentMaxSize, s.adjustmentFactor)

	if s.availableMemory >= s.memstatsCurrent.Alloc {
		if s.adjustmentFactor <= 1.0 {
			s.adjustmentFactor += 0.2
		}
	} else {
		// We're low on memory.
		s.adjustmentFactor -= 0.4
	}

	if s.adjustmentFactor <= 0 {
		s.adjustmentFactor = 0.05
	}

	if !s.adjustCurrentMaxSize() {
		return
	}

	totalWeight := 0
	for _, pm := range c.partitions {
		totalWeight += pm.getOptions().Weight
	}

	maxSizePerPartition := calculateMaxSizePerPartition(c.stats.currentMaxSize, totalWeight, len(c.partitions))

	evicted := 0
	for _, p := range c.partitions {
		evicted += p.adjustMaxSize(p.getOptions().CalculateMaxSize(maxSizePerPartition))
	}

	if evicted > 0 {
		c.infol.
			WithFields(
				logg.Fields{
					{Name: "evicted", Value: evicted},
					{Name: "numGC", Value: m.NumGC},
					{Name: "limit", Value: helpers.FormatByteCount(c.stats.availableMemory)},
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
				c.adjustCurrentMaxSize()
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

	if partitionNameRe.FindString(name) != name {
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

	// At this point, we don't know the number of partitions or their configuration, but
	// this will be re-adjusted later.
	const numberOfPartitionsEstimate = 10
	maxSize := opts.CalculateMaxSize(c.opts.MaxSize / numberOfPartitionsEstimate)

	onEvict := func(k K, v V) {
		c.onEvict(k, v)
	}

	// Create a new partition and cache it.
	partition := &Partition[K, V]{
		c:       lazycache.New(lazycache.Options[K, V]{MaxEntries: maxSize, OnEvict: onEvict}),
		maxSize: maxSize,
		trace:   c.opts.Log.Logger().WithLevel(logg.LevelTrace).WithField("partition", name),
		opts:    opts,
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

	maxSize int
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
		v, _, err := p.c.GetOrCreate(key, create)
		if err != nil {
			errch <- err
			return
		}
		resultch <- v
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

func (p *Partition[K, V]) clearOnRebuild(changeset ...identity.Identity) {
	opts := p.getOptions()
	if opts.ClearWhen == ClearNever {
		return
	}

	if opts.ClearWhen == ClearOnRebuild {
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
		if shouldDelete(key, v) {
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

// adjustMaxSize adjusts the max size of the and returns the number of items evicted.
func (p *Partition[K, V]) adjustMaxSize(newMaxSize int) int {
	if newMaxSize < minMaxSize {
		newMaxSize = minMaxSize
	}
	oldMaxSize := p.maxSize
	if newMaxSize == oldMaxSize {
		return 0
	}
	p.maxSize = newMaxSize
	// fmt.Println("Adjusting max size of partition from", oldMaxSize, "to", newMaxSize)
	return p.c.Resize(newMaxSize)
}

func (p *Partition[K, V]) getMaxSize() int {
	return p.maxSize
}

func (p *Partition[K, V]) getOptions() OptionsPartition {
	return p.opts
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
	adjustMaxSize(addend int) int
	getMaxSize() int
	getOptions() OptionsPartition
	clearOnRebuild(changeset ...identity.Identity)
	clearMatching(predicate func(k, v any) bool)
	clearStale()
}

const (
	ClearOnRebuild ClearWhen = iota + 1
	ClearOnChange
	ClearNever
)

type ClearWhen int

type stats struct {
	opts            Options
	memstatsCurrent runtime.MemStats
	currentMaxSize  int
	availableMemory uint64

	adjustmentFactor float64
}

func (s *stats) adjustCurrentMaxSize() bool {
	newCurrentMaxSize := int(math.Floor(float64(s.opts.MaxSize) * s.adjustmentFactor))

	if newCurrentMaxSize < s.opts.MinMaxSize {
		newCurrentMaxSize = int(s.opts.MinMaxSize)
	}
	changed := newCurrentMaxSize != s.currentMaxSize
	s.currentMaxSize = newCurrentMaxSize
	return changed
}

// CleanKey turns s into a format suitable for a cache key for this package.
// The key will be a Unix-styled path with a leading slash but no trailing slash.
func CleanKey(s string) string {
	return path.Clean(paths.ToSlashPreserveLeading(s))
}
