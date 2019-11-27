package identity

import (
	"path/filepath"
	"strings"
	"sync"
)

// NewIdentityManager creates a new Manager starting at root.
func NewIdentityManager(root Provider) Manager {
	return &identityManager{
		Provider: root,
		children: make(Identities, 0),
	}
}

// NewPathIdentity creates a new Identity with the two identifiers
// type and path.
func NewPathIdentity(typ, path string) PathIdentity {
	path = strings.TrimPrefix(filepath.ToSlash(path), "/")
	return PathIdentity{Type: typ, Path: path}
}

// Identities stores identity providers.
type Identities []Provider

// A set of identities.
type IdentitiesSet map[Identity]bool

// ToIdentitySet creates a set of these Identities.
func (ids Identities) ToIdentitySet() map[Identity]bool {
	m := make(map[Identity]bool)
	for _, id := range ids {
		m[id.GetIdentity()] = true
	}
	return m
}

func (ids Identities) search(id Identity) Provider {
	for _, v := range ids {
		vid := v.GetIdentity()

		if vid == id {
			return v
		}

		if idsp, ok := v.(ChildIdentitiesProvider); ok {
			if nested := idsp.GetChildIdentities().search(id); nested != nil {
				return nested
			}
		}
	}

	return nil
}

// IdentitiesProvider provides Identities as a set.
type IdentitiesProvider interface {
	GetIdentities() IdentitiesSet
}

// ChildIdentitiesProvider provides child Identities.
type ChildIdentitiesProvider interface {
	GetChildIdentities() Identities
}

// Identity represents an thing that can provide an identify. This can be
// any Go type, but the Identity returned by GetIdentify must be hashable.
type Identity interface {
	Provider
	Name() string
}

// Manager manages identities, and is itself a Provider of Identity.
type Manager interface {
	ChildIdentitiesProvider
	Provider
	Add(ids ...Provider)
	Search(id Identity) Provider
}

// A PathIdentity is a common identity identified by a type and a path, e.g. "layouts" and "_default/single.html".
type PathIdentity struct {
	Type string
	Path string
}

// GetIdentity returns itself.
func (id PathIdentity) GetIdentity() Identity {
	return id
}

// Name returns the Path.
func (id PathIdentity) Name() string {
	return id.Path
}

// Provider provides the hashable Identity.
type Provider interface {
	GetIdentity() Identity
}

type identityManager struct {
	sync.RWMutex
	Provider
	children Identities
}

func (im *identityManager) Add(ids ...Provider) {
	im.Lock()
	im.children = append(im.children, ids...)
	im.Unlock()
}

func (im *identityManager) GetChildIdentities() Identities {
	return im.children
}

func (im *identityManager) Search(id Identity) Provider {
	im.RLock()
	defer im.RUnlock()
	if id == im.GetIdentity() {
		return im
	}
	v := im.children.search(id)

	return v
}
