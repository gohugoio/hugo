package identity

import (
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
)

// NewIdentityManager creates a new Manager starting at id.
func NewManager(id Provider) Manager {
	return &identityManager{
		Provider: id,
		ids:      Identities{id.GetIdentity(): id},
	}
}

// NewPathIdentity creates a new Identity with the two identifiers
// type and path.
func NewPathIdentity(typ, pat string) PathIdentity {
	pat = strings.ToLower(strings.TrimPrefix(filepath.ToSlash(pat), "/"))
	return PathIdentity{Type: typ, Path: pat}
}

// Identities stores identity providers.
type Identities map[Identity]Provider

func (ids Identities) search(depth int, id Identity) Provider {

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
		switch t := v.(type) {
		case IdentitiesProvider:
			if nested := t.GetIdentities().search(depth, id); nested != nil {
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

// Identity represents an thing that can provide an identify. This can be
// any Go type, but the Identity returned by GetIdentify must be hashable.
type Identity interface {
	Provider
	Name() string
}

// Manager manages identities, and is itself a Provider of Identity.
type Manager interface {
	IdentitiesProvider
	Provider
	Add(ids ...Provider)
	Search(id Identity) Provider
	Reset()
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

// A KeyValueIdentity a general purpose identity.
type KeyValueIdentity struct {
	Key   string
	Value string
}

// GetIdentity returns itself.
func (id KeyValueIdentity) GetIdentity() Identity {
	return id
}

// Name returns the Key.
func (id KeyValueIdentity) Name() string {
	return id.Key
}

// Provider provides the hashable Identity.
type Provider interface {
	GetIdentity() Identity
}

type identityManager struct {
	sync.Mutex
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
	id := im.GetIdentity()
	im.ids = Identities{id.GetIdentity(): id}
	im.Unlock()
}

func (im *identityManager) GetIdentities() Identities {
	im.Lock()
	defer im.Unlock()
	return im.ids
}

func (im *identityManager) Search(id Identity) Provider {
	im.Lock()
	defer im.Unlock()
	return im.ids.search(0, id.GetIdentity())
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
