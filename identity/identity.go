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

// Package provides ways to identify values in Hugo. Used for dependency tracking etc.
package identity

import (
	"sync"
	"sync/atomic"

	"github.com/gohugoio/hugo/compare"
)

// NewIdentityManager creates a new Manager starting at id.
func NewManager(id Provider) Manager {
	return &identityManager{
		Provider: id,
		ids:      Identities{id.GetIdentity(): id},
	}
}

// Identities stores identity providers.
type Identities map[Identity]Provider

func (ids Identities) isNotNotDependent(p1, p2 Provider) (Provider, bool) {

	// Let both p1 and p2 vote on this:
	if !p1.IsNotDependent(p2) {
		return p1, true
	}

	if !p2.IsNotDependent(p1) {
		return p1, true
	}

	// There may false positives in the above, but we should know be
	// sure about p1 and p2 not being dependent on each other.
	return nil, false

}
func (ids Identities) search(depth int, probableMatch bool, id Provider) Provider {

	// TODO(bep): .Base()
	if v, found := ids[id.GetIdentity()]; found {
		return v
	}

	depth++

	// There may be infinite recursion in templates.
	if depth > 100 {
		// Bail out.
		return nil
	}

	for _, v := range ids {
		if probableMatch {
			if p, ok := ids.isNotNotDependent(v, id); ok {

				// It is not not dependent, which may be a false positive,
				// but that is OK in this case.
				return p

			}
		}

		switch t := v.(type) {
		case IdentitiesProvider:
			if nested := t.GetIdentities().search(depth, probableMatch, id); nested != nil {
				return nested
			}
		}
	}
	return nil
}

// IdentitiesProvider provides all Identities.
type IdentitiesProvider interface {
	GetIdentities() Identities
}

/*

- manager.Search(id/manager)

*/
type Provider interface {
	GetIdentity() Identity
	IsNotDependentProvider
}

// Identity represents an thing that can provide an identify. This can be
// any Go type, but the Identity returned by GetIdentify must be hashable.
type Identity interface {
	//Provider
	compare.ProbablyEqer
	compare.Eqer
	Base() interface{}
	Name() string
}

// Manager manages identities, and is itself a Provider of Identity.
type Manager interface {
	IdentitiesProvider
	Provider
	Add(ids ...Provider)
	Search(id Provider) Provider
	SearchProbablyDependent(id Provider) Provider
	Reset()
}

// A KeyValueIdentity a general purpose identity.
type KeyValueIdentity struct {
	Key   string
	Value string
}

// GetIdentity returns itself.
func (id KeyValueIdentity) GetIdentity() Identity {
	return id
}

func (id KeyValueIdentity) IsNotDependent(other Provider) bool {
	return id != other
}

func (id KeyValueIdentity) Base() interface{} {
	return id
}

// Name returns the Key.
func (id KeyValueIdentity) Name() string {
	return id.Key
}

func (id KeyValueIdentity) Eq(other interface{}) bool {
	return id == other
}

func (id KeyValueIdentity) ProbablyEq(other interface{}) bool {
	return id == other
}

// IsNotDependentProvider provides a method to determin if the other
// Provider is not dependent on this one.
type IsNotDependentProvider interface {
	IsNotDependent(other Provider) bool
}

type identityManager struct {
	sync.RWMutex
	Provider
	ids Identities
}

func (im *identityManager) Add(ids ...Provider) {
	im.Lock()
	for _, id := range ids {
		im.ids[id.GetIdentity()] = id
	}
	im.Unlock()
}

func (im *identityManager) Reset() {
	im.Lock()
	im.ids = Identities{im.GetIdentity(): im}
	im.Unlock()
}

// TODO(bep) these identities are currently only read on server reloads
// so there should be no concurrency issues, but that may change.
func (im *identityManager) GetIdentities() Identities {
	im.Lock()
	defer im.Unlock()
	return im.ids
}

func (im *identityManager) IsNotDependent(other Provider) bool {
	got := im.SearchProbablyDependent(other)
	return got == nil
}

func (im *identityManager) Search(id Provider) Provider {
	im.Lock()
	defer im.Unlock()
	return im.ids.search(0, false, id)
}

func (im *identityManager) SearchProbablyDependent(id Provider) Provider {
	im.Lock()
	defer im.Unlock()
	return im.ids.search(0, true, id)
}

// Incrementer increments and returns the value.
// Typically used for IDs.
type Incrementer interface {
	Incr() int
}

// IncrementByOne implements Incrementer adding 1 every time Incr is called.
type IncrementByOne struct {
	counter uint64
}

func (c *IncrementByOne) Incr() int {
	return int(atomic.AddUint64(&c.counter, uint64(1)))
}
