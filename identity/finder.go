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

package identity

import (
	"fmt"
	"sync"

	"github.com/gohugoio/hugo/compare"
)

// NewFinder creates a new Finder.
// This is a thread safe implementation with a cache.
func NewFinder(cfg FinderConfig) *Finder {
	return &Finder{cfg: cfg, answers: make(map[ManagerIdentity]FinderResult), seenFindOnce: make(map[Identity]bool)}
}

var searchIDPool = sync.Pool{
	New: func() interface{} {
		return &searchID{seen: make(map[Manager]bool)}
	},
}

func getSearchID() *searchID {
	return searchIDPool.Get().(*searchID)
}

func putSearchID(sid *searchID) {
	sid.id = nil
	sid.isDp = false
	sid.isPeq = false
	sid.hasEqer = false
	sid.maxDepth = 0
	sid.dp = nil
	sid.peq = nil
	sid.eqer = nil
	for k := range sid.seen {
		delete(sid.seen, k)
	}
	searchIDPool.Put(sid)
}

// GetSearchID returns a searchID from the pool.

// Finder finds identities inside another.
type Finder struct {
	cfg FinderConfig

	answers   map[ManagerIdentity]FinderResult
	muAnswers sync.RWMutex

	seenFindOnce   map[Identity]bool
	muSeenFindOnce sync.RWMutex
}

type FinderResult int

const (
	FinderNotFound FinderResult = iota
	FinderFoundOneOfManyRepetition
	FinderFoundOneOfMany
	FinderFound
)

// Contains returns whether in contains id.
func (f *Finder) Contains(id, in Identity, maxDepth int) FinderResult {
	if id == Anonymous || in == Anonymous {
		return FinderNotFound
	}

	if id == GenghisKhan && in == GenghisKhan {
		return FinderNotFound
	}

	if id == GenghisKhan {
		return FinderFound
	}

	if id == in {
		return FinderFound
	}

	if id == nil || in == nil {
		return FinderNotFound
	}

	var (
		isDp  bool
		isPeq bool

		dp  IsProbablyDependentProvider
		peq compare.ProbablyEqer
	)

	if !f.cfg.Exact {
		dp, isDp = id.(IsProbablyDependentProvider)
		peq, isPeq = id.(compare.ProbablyEqer)
	}

	eqer, hasEqer := id.(compare.Eqer)

	sid := getSearchID()
	sid.id = id
	sid.isDp = isDp
	sid.isPeq = isPeq
	sid.hasEqer = hasEqer
	sid.dp = dp
	sid.peq = peq
	sid.eqer = eqer
	sid.maxDepth = maxDepth

	defer putSearchID(sid)

	if r := f.checkOne(sid, in, 0); r > 0 {
		return r
	}

	m := GetDependencyManager(in)
	if m != nil {
		if r := f.checkManager(sid, m, 0); r > 0 {
			return r
		}
	}
	return FinderNotFound
}

func (f *Finder) checkMaxDepth(sid *searchID, level int) FinderResult {
	if sid.maxDepth >= 0 && level > sid.maxDepth {
		return FinderNotFound
	}
	if level > 100 {
		// This should never happen, but some false positives are probably better than a panic.
		if !f.cfg.Exact {
			return FinderFound
		}
		panic("too many levels")
	}
	return -1
}

func (f *Finder) checkManager(sid *searchID, m Manager, level int) FinderResult {
	if r := f.checkMaxDepth(sid, level); r >= 0 {
		return r
	}

	if m == nil {
		return FinderNotFound
	}
	if sid.seen[m] {
		return FinderNotFound
	}
	sid.seen[m] = true

	f.muAnswers.RLock()
	r, ok := f.answers[ManagerIdentity{Manager: m, Identity: sid.id}]
	f.muAnswers.RUnlock()
	if ok {
		return r
	}

	r = f.search(sid, m, level)

	if r == FinderFoundOneOfMany {
		// Don't cache this one.
		return r
	}

	f.muAnswers.Lock()
	f.answers[ManagerIdentity{Manager: m, Identity: sid.id}] = r
	f.muAnswers.Unlock()

	return r
}

func (f *Finder) checkOne(sid *searchID, v Identity, depth int) (r FinderResult) {
	if ff, ok := v.(FindFirstManagerIdentityProvider); ok {
		f.muSeenFindOnce.RLock()
		mi := ff.FindFirstManagerIdentity()
		seen := f.seenFindOnce[mi.Identity]
		f.muSeenFindOnce.RUnlock()
		if seen {
			return FinderFoundOneOfManyRepetition
		}

		r = f.doCheckOne(sid, mi.Identity, depth)
		if r == 0 {
			r = f.checkManager(sid, mi.Manager, depth)
		}

		if r > FinderFoundOneOfManyRepetition {
			f.muSeenFindOnce.Lock()
			// Double check.
			if f.seenFindOnce[mi.Identity] {
				f.muSeenFindOnce.Unlock()
				return FinderFoundOneOfManyRepetition
			}
			f.seenFindOnce[mi.Identity] = true
			f.muSeenFindOnce.Unlock()
			r = FinderFoundOneOfMany
		}
		return r
	} else {
		return f.doCheckOne(sid, v, depth)
	}
}

func (f *Finder) doCheckOne(sid *searchID, v Identity, depth int) FinderResult {
	id2 := Unwrap(v)
	if id2 == Anonymous {
		return FinderNotFound
	}
	id := sid.id
	if sid.hasEqer {
		if sid.eqer.Eq(id2) {
			return FinderFound
		}
	} else if id == id2 {
		return FinderFound
	}

	if f.cfg.Exact {
		return FinderNotFound
	}

	if id2 == nil {
		return FinderNotFound
	}

	if id2 == GenghisKhan {
		return FinderFound
	}

	if id.IdentifierBase() == id2.IdentifierBase() {
		return FinderFound
	}

	if sid.isDp && sid.dp.IsProbablyDependent(id2) {
		return FinderFound
	}

	if sid.isPeq && sid.peq.ProbablyEq(id2) {
		return FinderFound
	}

	if pdep, ok := id2.(IsProbablyDependencyProvider); ok && pdep.IsProbablyDependency(id) {
		return FinderFound
	}

	if peq, ok := id2.(compare.ProbablyEqer); ok && peq.ProbablyEq(id) {
		return FinderFound
	}

	return FinderNotFound
}

// search searches for id in ids.
func (f *Finder) search(sid *searchID, m Manager, depth int) FinderResult {
	id := sid.id

	if id == Anonymous {
		return FinderNotFound
	}

	if !f.cfg.Exact && id == GenghisKhan {
		return FinderNotFound
	}

	var r FinderResult
	m.forEeachIdentity(
		func(v Identity) bool {
			if r > 0 {
				panic("should be terminated")
			}
			r = f.checkOne(sid, v, depth)
			if r > 0 {
				return true
			}
			m := GetDependencyManager(v)
			if r = f.checkManager(sid, m, depth+1); r > 0 {
				return true
			}
			return false
		},
	)
	return r
}

// FinderConfig provides configuration for the Finder.
// Note that we by default will use a strategy where probable matches are
// good enough. The primary use case for this is to identity the change set
// for a given changed identity (e.g. a template), and we don't want to
// have any false negatives there, but some false positives are OK. Also, speed is important.
type FinderConfig struct {
	// Match exact matches only.
	Exact bool
}

// ManagerIdentity wraps a pair of Identity and Manager.
type ManagerIdentity struct {
	Identity
	Manager
}

func (p ManagerIdentity) String() string {
	return fmt.Sprintf("%s:%s", p.Identity.IdentifierBase(), p.Manager.IdentifierBase())
}

type searchID struct {
	id      Identity
	isDp    bool
	isPeq   bool
	hasEqer bool

	maxDepth int

	seen map[Manager]bool

	dp   IsProbablyDependentProvider
	peq  compare.ProbablyEqer
	eqer compare.Eqer
}
