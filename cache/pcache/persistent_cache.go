// Copyright 2018 The Hugo Authors. All rights reserved.
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

package pcache

import (
	"strings"

	"github.com/mitchellh/mapstructure"

	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"sort"

	"sync"
)

var (
	_ OpenCloserCache = (*persistentCache)(nil)
	_ Cache           = (*noCache)(nil)
	_ Identifier      = (*VersionedID)(nil)
	// Used in unit tests.
	NoCache Cache = &noCache{}
)

type Identifier interface {
	_ID() string
}

// This is persisted to one cache file. The semantics are simple: If a cache request
// arrives with for a different version or there are no entries to be found,
// we start fresh.
// We sort this by ID before saving it to disk.
type cacheEntries struct {
	elemType reflect.Type
	Version  string
	Entries  []Identifier
}

func (c *cacheEntries) toCacheEntriesJSON() *cacheEntriesJSON {
	return &cacheEntriesJSON{Version: c.Version}
}

type cacheEntriesJSON struct {
	Version string
	Entries json.RawMessage
}

func (m *cacheEntries) UnmarshalJSON(value []byte) error {

	mj := m.toCacheEntriesJSON()

	if err := json.Unmarshal(value, mj); err != nil {
		return err
	}

	dec := json.NewDecoder(bytes.NewReader(mj.Entries))
	dec.UseNumber()

	// We need to massage this later, and this is the best we can do with
	// Go's JSON package.
	var ms []map[string]interface{}

	if err := dec.Decode(&ms); err != nil {
		return err
	}

	slice := reflect.MakeSlice(reflect.SliceOf(m.elemType), 0, 0)

	for _, msm := range ms {

		n := reflect.New(m.elemType.Elem())
		result := n.Interface()

		hook := func(t1, t2 reflect.Type, v interface{}) (interface{}, error) {
			// TODO(bep) defaultTypeConverters => struct
			vv, _, err := defaultTypeConverters.ConvertTypes(v, t1, t2)

			return vv, err

		}

		mdec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			DecodeHook:       hook,
			Metadata:         nil,
			Result:           result,
			WeaklyTypedInput: true,
		})
		if err != nil {
			return err
		}

		if err := mdec.Decode(msm); err != nil {
			return err
		}

		slice = reflect.Append(slice, reflect.ValueOf(result))

	}

	m.Version = mj.Version

	m.Entries = make([]Identifier, slice.Len())
	for i := 0; i < slice.Len(); i++ {
		e := slice.Index(i).Interface()
		id := e.(Identifier)
		m.Entries[i] = id
	}

	return nil
}

func (c *cacheEntries) toMap() *cacheEntriesMap {
	m := &cacheEntriesMap{Version: c.Version, Entries: make(map[string]*flaggedIdentifier), elemType: c.elemType}
	for _, e := range c.Entries {
		eid := e.(Identifier)
		m.Entries[eid._ID()] = newFlaggedIdentifier(eid, false)
	}
	return m
}

// We use this to track cache entries no longer in use.
type flaggedIdentifier struct {
	Identifier

	// Set if this item is used between Open and Close.
	IsTouched bool

	// Set for newly created items.
	IsNew bool
}

func newFlaggedIdentifier(id Identifier, isNew bool) *flaggedIdentifier {
	return &flaggedIdentifier{Identifier: id, IsTouched: isNew, IsNew: isNew}
}

// This is the variant used for lookups.
type cacheEntriesMap struct {
	elemType reflect.Type
	Version  string
	Entries  map[string]*flaggedIdentifier
}

// isStale returns whether this map has changed since open; changed, added or removed
// entries.
func (c *cacheEntriesMap) isStale() bool {
	for _, v := range c.Entries {
		if !v.IsTouched || v.IsNew {
			return true
		}
	}
	return false
}

func (c *cacheEntriesMap) getActiveEntriesSorted() *cacheEntries {
	ce := &cacheEntries{Version: c.Version, Entries: make([]Identifier, 0), elemType: c.elemType}
	for _, v := range c.Entries {
		if v.IsTouched {
			ce.Entries = append(ce.Entries, v.Identifier)
		}
	}

	sort.Slice(ce.Entries, func(i, j int) bool {
		return ce.Entries[i]._ID() < ce.Entries[j]._ID()
	})

	return ce
}

// Note that the caller must lock.
func (c *cacheEntriesMap) getItem(ID string) (Identifier, bool) {
	item, found := c.Entries[ID]
	if !found {
		return nil, false
	}
	item.IsTouched = true
	return item.Identifier, true
}

// VersionedID identifies an entity in the cache.
type VersionedID struct {
	Version string
	ID      string
}

func (vid VersionedID) _ID() string {
	return vid.ID
}

func NewVersionedID(version, ID string) VersionedID {
	return VersionedID{Version: version, ID: ID}
}

type Cache interface {
	GetOrCreate(ID VersionedID, create func() (Identifier, error)) (Identifier, error)
}

type OpenCloserCache interface {
	Cache
	Open() error
	Close() error
}

type persistentCache struct {
	elemType reflect.Type
	data     *cacheEntriesMap

	typeHandlers typeConverters

	openInit sync.Once
	openErr  error
	open     bool

	filename string
	sync.RWMutex
}

func New(filename string, entity interface{}) OpenCloserCache {
	return &persistentCache{filename: filename, elemType: reflect.TypeOf(entity), typeHandlers: defaultTypeConverters}
}

func (c *persistentCache) newCacheEntries() *cacheEntries {
	return &cacheEntries{elemType: c.elemType}
}

func (c *persistentCache) newCacheEntriesMap() *cacheEntriesMap {
	return &cacheEntriesMap{elemType: c.elemType, Entries: make(map[string]*flaggedIdentifier)}
}

func (c *persistentCache) Open() error {
	//  We delay the open until it gets used. Maybe never.
	return nil
}

func (c *persistentCache) lazyOpen() error {

	c.openInit.Do(func() {
		c.Lock()
		defer c.Unlock()

		b, e := ioutil.ReadFile(c.filename)
		if e != nil {
			if os.IsNotExist(e) {
				c.data = c.newCacheEntriesMap()
				c.open = true
				return
			} else {
				c.openErr = e
				return
			}
		}

		data := c.newCacheEntries()
		c.unmarshal(b, data)
		if c.openErr != nil {
			return
		}

		c.data = data.toMap()

		c.open = true

	})

	return c.openErr
}

func (c *persistentCache) unmarshal(b []byte, ce *cacheEntries) {
	c.openErr = json.Unmarshal(b, ce)
}

var prettyJSONReplacer = strings.NewReplacer(
	`"Entries":[`, "\"Entries\":[\n",
	"},", "},\n",
	"}]}", "}\n]}",
)

func (c *persistentCache) Close() error {
	c.Lock()
	defer c.Unlock()

	if !c.open || !c.data.isStale() {
		return nil
	}

	var (
		f   *os.File
		err error
	)

	// TODO(bep) mkdirall
	_, err = os.Stat(c.filename)

	if os.IsNotExist(err) {
		f, err = os.Create(c.filename)
	} else {
		f, err = os.OpenFile(c.filename, os.O_RDWR, 0755)
	}

	if err != nil {
		return err
	}
	defer f.Close()

	b, err := json.Marshal(c.data.getActiveEntriesSorted())
	if err != nil {
		return err
	}

	s := prettyJSONReplacer.Replace(string(b))

	_, err = f.WriteString(s)

	return err
}

// GetOrCreate fetches the value from the versioned data store. If not, found
// it is created with the supplied create func and put there.
func (c *persistentCache) GetOrCreate(ID VersionedID, create func() (Identifier, error)) (Identifier, error) {
	if err := c.lazyOpen(); err != nil {
		return nil, err
	}

	c.RLock()

	if c.data.Version != ID.Version {
		// Version upgrade.
		c.RUnlock()

		c.Lock()
		defer c.Unlock()
		c.data.Version = ID.Version
		c.data.Entries = make(map[string]*flaggedIdentifier)
		v, err := create()
		if err != nil {
			return nil, err
		}

		fi := newFlaggedIdentifier(v, true)
		fi.IsTouched = true
		c.data.Entries[ID.ID] = fi

		return v, nil

	}

	v, found := c.data.getItem(ID.ID)
	c.RUnlock()

	if found {
		return v, nil
	}

	c.Lock()
	defer c.Unlock()
	if v, found = c.data.getItem(ID.ID); found {
		return v, nil
	}

	vc, err := create()
	if err != nil {
		return nil, err
	}

	fi := newFlaggedIdentifier(vc, true)
	fi.IsTouched = true
	c.data.Entries[ID.ID] = fi

	return fi.Identifier, nil
}

type noCache struct {
}

func (c *noCache) GetOrCreate(ID VersionedID, create func() (Identifier, error)) (Identifier, error) {
	return create()
}
