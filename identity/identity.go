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

// Package provides ways to identify values in Hugo. Used for dependency tracking etc.
package identity

import (
	"fmt"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/gohugoio/hugo/common/types"
	"github.com/gohugoio/hugo/compare"
)

const (
	// Anonymous is an Identity that can be used when identity doesn't matter.
	Anonymous = StringIdentity("__anonymous")

	// GenghisKhan is an Identity everyone relates to.
	GenghisKhan = StringIdentity("__genghiskhan")
)

var NopManager = new(nopManager)

// NewIdentityManager creates a new Manager.
func NewManager(name string, opts ...ManagerOption) Manager {
	idm := &identityManager{
		Identity: Anonymous,
		name:     name,
		ids:      Identities{},
	}

	for _, o := range opts {
		o(idm)
	}

	return idm
}

// CleanString cleans s to be suitable as an identifier.
func CleanString(s string) string {
	s = strings.ToLower(s)
	s = strings.Trim(filepath.ToSlash(s), "/")
	return "/" + path.Clean(s)
}

// CleanStringIdentity cleans s to be suitable as an identifier and wraps it in a StringIdentity.
func CleanStringIdentity(s string) StringIdentity {
	return StringIdentity(CleanString(s))
}

// GetDependencyManager returns the DependencyManager from v or nil if none found.
func GetDependencyManager(v any) Manager {
	switch vv := v.(type) {
	case Manager:
		return vv
	case types.Unwrapper:
		return GetDependencyManager(vv.Unwrapv())
	case DependencyManagerProvider:
		return vv.GetDependencyManager()
	}
	return nil
}

// FirstIdentity returns the first Identity in v, Anonymous if none found
func FirstIdentity(v any) Identity {
	var result Identity = Anonymous
	WalkIdentitiesShallow(v, func(level int, id Identity) bool {
		result = id
		return true
	})

	return result
}

// PrintIdentityInfo is used for debugging/tests only.
func PrintIdentityInfo(v any) {
	WalkIdentitiesDeep(v, func(level int, id Identity) bool {
		var s string
		if idm, ok := id.(*identityManager); ok {
			s = " " + idm.name
		}
		fmt.Printf("%s%s (%T)%s\n", strings.Repeat("  ", level), id.IdentifierBase(), id, s)
		return false
	})
}

func Unwrap(id Identity) Identity {
	switch t := id.(type) {
	case IdentityProvider:
		return t.GetIdentity()
	default:
		return id
	}
}

// WalkIdentitiesDeep walks identities in v and applies cb to every identity found.
// Return true from cb to terminate.
// If deep is true, it will also walk nested Identities in any Manager found.
func WalkIdentitiesDeep(v any, cb func(level int, id Identity) bool) {
	seen := make(map[Identity]bool)
	walkIdentities(v, 0, true, seen, cb)
}

// WalkIdentitiesShallow will not walk into a Manager's Identities.
// See WalkIdentitiesDeep.
// cb is called for every Identity found and returns whether to terminate the walk.
func WalkIdentitiesShallow(v any, cb func(level int, id Identity) bool) {
	walkIdentitiesShallow(v, 0, cb)
}

// WithOnAddIdentity sets a callback that will be invoked when an identity is added to the manager.
func WithOnAddIdentity(f func(id Identity)) ManagerOption {
	return func(m *identityManager) {
		m.onAddIdentity = f
	}
}

// DependencyManagerProvider provides a manager for dependencies.
type DependencyManagerProvider interface {
	GetDependencyManager() Manager
}

// DependencyManagerProviderFunc is a function that implements the DependencyManagerProvider interface.
type DependencyManagerProviderFunc func() Manager

func (d DependencyManagerProviderFunc) GetDependencyManager() Manager {
	return d()
}

// DependencyManagerScopedProvider provides a manager for dependencies with a given scope.
type DependencyManagerScopedProvider interface {
	GetDependencyManagerForScope(scope int) Manager
}

// ForEeachIdentityProvider provides a way iterate over identities.
type ForEeachIdentityProvider interface {
	// ForEeachIdentityProvider calls cb for each Identity.
	// If cb returns true, the iteration is terminated.
	// The return value is whether the iteration was terminated.
	ForEeachIdentity(cb func(id Identity) bool) bool
}

// ForEeachIdentityProviderFunc is a function that implements the ForEeachIdentityProvider interface.
type ForEeachIdentityProviderFunc func(func(id Identity) bool) bool

func (f ForEeachIdentityProviderFunc) ForEeachIdentity(cb func(id Identity) bool) bool {
	return f(cb)
}

// ForEeachIdentityByNameProvider provides a way to look up identities by name.
type ForEeachIdentityByNameProvider interface {
	// ForEeachIdentityByName calls cb for each Identity that relates to name.
	// If cb returns true, the iteration is terminated.
	ForEeachIdentityByName(name string, cb func(id Identity) bool)
}

type FindFirstManagerIdentityProvider interface {
	Identity
	FindFirstManagerIdentity() ManagerIdentity
}

func NewFindFirstManagerIdentityProvider(m Manager, id Identity) FindFirstManagerIdentityProvider {
	return findFirstManagerIdentity{
		Identity: Anonymous,
		ManagerIdentity: ManagerIdentity{
			Manager: m, Identity: id,
		},
	}
}

type findFirstManagerIdentity struct {
	Identity
	ManagerIdentity
}

func (f findFirstManagerIdentity) FindFirstManagerIdentity() ManagerIdentity {
	return f.ManagerIdentity
}

// Identities stores identity providers.
type Identities map[Identity]bool

func (ids Identities) AsSlice() []Identity {
	s := make([]Identity, len(ids))
	i := 0
	for v := range ids {
		s[i] = v
		i++
	}
	sort.Slice(s, func(i, j int) bool {
		return s[i].IdentifierBase() < s[j].IdentifierBase()
	})

	return s
}

func (ids Identities) String() string {
	var sb strings.Builder
	i := 0
	for id := range ids {
		sb.WriteString(fmt.Sprintf("[%s]", id.IdentifierBase()))
		if i < len(ids)-1 {
			sb.WriteString(", ")
		}
		i++
	}
	return sb.String()
}

// Identity represents a thing in Hugo (a Page, a template etc.)
// Any implementation must be comparable/hashable.
type Identity interface {
	IdentifierBase() string
}

// IdentityGroupProvider can be implemented by tightly connected types.
// Current use case is Resource transformation via Hugo Pipes.
type IdentityGroupProvider interface {
	GetIdentityGroup() Identity
}

// IdentityProvider can be implemented by types that isn't itself and Identity,
// usually because they're not comparable/hashable.
type IdentityProvider interface {
	GetIdentity() Identity
}

// SignalRebuilder is an optional interface for types that can signal a rebuild.
type SignalRebuilder interface {
	SignalRebuild(ids ...Identity)
}

// IncrementByOne implements Incrementer adding 1 every time Incr is called.
type IncrementByOne struct {
	counter uint64
}

func (c *IncrementByOne) Incr() int {
	return int(atomic.AddUint64(&c.counter, uint64(1)))
}

// Incrementer increments and returns the value.
// Typically used for IDs.
type Incrementer interface {
	Incr() int
}

// IsProbablyDependentProvider is an optional interface for Identity.
type IsProbablyDependentProvider interface {
	IsProbablyDependent(other Identity) bool
}

// IsProbablyDependencyProvider is an optional interface for Identity.
type IsProbablyDependencyProvider interface {
	IsProbablyDependency(other Identity) bool
}

// Manager  is an Identity that also manages identities, typically dependencies.
type Manager interface {
	Identity
	AddIdentity(ids ...Identity)
	AddIdentityForEach(ids ...ForEeachIdentityProvider)
	GetIdentity() Identity
	Reset()
	forEeachIdentity(func(id Identity) bool) bool
}

type ManagerOption func(m *identityManager)

// StringIdentity is an Identity that wraps a string.
type StringIdentity string

func (s StringIdentity) IdentifierBase() string {
	return string(s)
}

type identityManager struct {
	Identity

	// Only used for debugging.
	name string

	// mu protects _changes_ to this manager,
	// reads currently assumes no concurrent writes.
	mu         sync.RWMutex
	ids        Identities
	forEachIds []ForEeachIdentityProvider

	// Hooks used in debugging.
	onAddIdentity func(id Identity)
}

func (im *identityManager) AddIdentity(ids ...Identity) {
	im.mu.Lock()

	for _, id := range ids {
		if id == nil || id == Anonymous {
			continue
		}
		if _, found := im.ids[id]; !found {
			if im.onAddIdentity != nil {
				im.onAddIdentity(id)
			}
			im.ids[id] = true
		}
	}
	im.mu.Unlock()
}

func (im *identityManager) AddIdentityForEach(ids ...ForEeachIdentityProvider) {
	im.mu.Lock()
	im.forEachIds = append(im.forEachIds, ids...)
	im.mu.Unlock()
}

func (im *identityManager) ContainsIdentity(id Identity) FinderResult {
	if im.Identity != Anonymous && id == im.Identity {
		return FinderFound
	}

	f := NewFinder(FinderConfig{Exact: true})
	r := f.Contains(id, im, -1)

	return r
}

// Managers are always anonymous.
func (im *identityManager) GetIdentity() Identity {
	return im.Identity
}

func (im *identityManager) Reset() {
	im.mu.Lock()
	im.ids = Identities{}
	im.mu.Unlock()
}

func (im *identityManager) GetDependencyManagerForScope(int) Manager {
	return im
}

func (im *identityManager) String() string {
	return fmt.Sprintf("IdentityManager(%s)", im.name)
}

func (im *identityManager) forEeachIdentity(fn func(id Identity) bool) bool {
	// The absence of a lock here is deliberate. This is currently only used on server reloads
	// in a single-threaded context.
	for id := range im.ids {
		if fn(id) {
			return true
		}
	}
	for _, fe := range im.forEachIds {
		if fe.ForEeachIdentity(fn) {
			return true
		}
	}
	return false
}

type nopManager int

func (m *nopManager) AddIdentity(ids ...Identity) {
}

func (m *nopManager) AddIdentityForEach(ids ...ForEeachIdentityProvider) {
}

func (m *nopManager) IdentifierBase() string {
	return ""
}

func (m *nopManager) GetIdentity() Identity {
	return Anonymous
}

func (m *nopManager) Reset() {
}

func (m *nopManager) forEeachIdentity(func(id Identity) bool) bool {
	return false
}

// returns whether further walking should be terminated.
func walkIdentities(v any, level int, deep bool, seen map[Identity]bool, cb func(level int, id Identity) bool) {
	if level > 20 {
		panic("too deep")
	}
	var cbRecursive func(level int, id Identity) bool
	cbRecursive = func(level int, id Identity) bool {
		if id == nil {
			return false
		}
		if deep && seen[id] {
			return false
		}
		seen[id] = true
		if cb(level, id) {
			return true
		}

		if deep {
			if m := GetDependencyManager(id); m != nil {
				m.forEeachIdentity(func(id2 Identity) bool {
					return walkIdentitiesShallow(id2, level+1, cbRecursive)
				})
			}
		}
		return false
	}
	walkIdentitiesShallow(v, level, cbRecursive)
}

// returns whether further walking should be terminated.
// Anonymous identities are skipped.
func walkIdentitiesShallow(v any, level int, cb func(level int, id Identity) bool) bool {
	cb2 := func(level int, id Identity) bool {
		if id == Anonymous {
			return false
		}
		if id == nil {
			return false
		}
		return cb(level, id)
	}

	if id, ok := v.(Identity); ok {
		if cb2(level, id) {
			return true
		}
	}

	if ipd, ok := v.(IdentityProvider); ok {
		if cb2(level, ipd.GetIdentity()) {
			return true
		}
	}

	if ipdgp, ok := v.(IdentityGroupProvider); ok {
		if cb2(level, ipdgp.GetIdentityGroup()) {
			return true
		}
	}

	return false
}

var (
	_ Identity             = (*orIdentity)(nil)
	_ compare.ProbablyEqer = (*orIdentity)(nil)
)

func Or(a, b Identity) Identity {
	return orIdentity{a: a, b: b}
}

type orIdentity struct {
	a, b Identity
}

func (o orIdentity) IdentifierBase() string {
	return o.a.IdentifierBase()
}

func (o orIdentity) ProbablyEq(other any) bool {
	otherID, ok := other.(Identity)
	if !ok {
		return false
	}

	return probablyEq(o.a, otherID) || probablyEq(o.b, otherID)
}

func probablyEq(a, b Identity) bool {
	if a == b {
		return true
	}

	if a == Anonymous || b == Anonymous {
		return false
	}

	if a.IdentifierBase() == b.IdentifierBase() {
		return true
	}

	if a2, ok := a.(IsProbablyDependentProvider); ok {
		return a2.IsProbablyDependent(b)
	}

	return false
}
