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
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bep/mapstructure"

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
	// NoCache is used in unit tests.
	NoCache Cache = &noCache{}
)

// Identifier is the interface all cached items must implement.
type Identifier interface {
	_ID() ID
}

// This is persisted to one cache file. The semantics are simple: If a cache request
// arrives with for a different version or there are no entries to be found,
// we start fresh.
// We sort this by ID before saving it to disk.
type cacheEntries struct {
	c *persistentCache

	FieldConverters map[string]string

	Version string
	Entries []Identifier
}

func (c *cacheEntries) toCacheEntriesJSON() *cacheEntriesJSON {
	return &cacheEntriesJSON{Version: c.Version, FieldConverters: c.FieldConverters}
}

type cacheEntriesJSON struct {
	Version         string
	FieldConverters map[string]string

	Entries json.RawMessage
}

func (c *cacheEntries) UnmarshalJSON(value []byte) error {

	mj := c.toCacheEntriesJSON()

	if err := json.Unmarshal(value, mj); err != nil {
		return err
	}

	c.Version = mj.Version
	c.FieldConverters = mj.FieldConverters

	dec := json.NewDecoder(bytes.NewReader(mj.Entries))
	dec.UseNumber()

	// We need to massage this later, and this is the best we can do with
	// Go's JSON package.
	var ms []map[string]interface{}

	if err := dec.Decode(&ms); err != nil {
		return err
	}

	slice := reflect.MakeSlice(reflect.SliceOf(c.c.elemType), 0, 0)

	for _, msm := range ms {
		elemType := c.c.elemType
		if c.c.isPtr {
			elemType = elemType.Elem()
		}
		v := reflect.New(elemType)
		resultp := v.Interface()

		hook := func(t1, t2 reflect.Type, v interface{}, ctx *mapstructure.DecodeContext) (interface{}, error) {
			if ctx.IsKey {
				return v, nil
			}

			fmt.Println(ctx.Name, ": >>> t1:", t1, "t2:", t2, "v:", v)

			converter := c.c.typem.converters.GetByTypes(t1, t2)
			if converter == nil {
				convertorName, found := c.FieldConverters[ctx.Name]
				if found {
					converter = c.c.typem.converters.GetByName(convertorName)
				}
			}

			if converter != nil {
				return converter.Convert(v)
			}

			return v, nil
		}

		mdec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			DecodeHook: hook,
			Metadata:   nil,
			Result:     resultp,
			//WeaklyTypedInput: true,
			//ZeroFields:       true,
		})
		if err != nil {
			return err
		}

		if err := mdec.Decode(msm); err != nil {
			return err
		}

		var result reflect.Value
		if c.c.isPtr {
			result = reflect.ValueOf(resultp)
		} else {
			result = reflect.Indirect(reflect.ValueOf(resultp))
		}

		slice = reflect.Append(slice, result)
	}

	c.Entries = make([]Identifier, slice.Len())
	for i := 0; i < slice.Len(); i++ {
		e := slice.Index(i).Interface()
		id := e.(Identifier)
		c.Entries[i] = id
	}

	return nil
}

func (c *cacheEntries) toMap() *cacheEntriesMap {
	m := &cacheEntriesMap{Version: c.Version, FieldConverters: c.FieldConverters, Entries: make(map[ID]*flaggedIdentifier), c: c.c}
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
	c *persistentCache

	FieldConverters map[string]string

	Version string
	Entries map[ID]*flaggedIdentifier
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
	ce := &cacheEntries{Version: c.Version, FieldConverters: c.FieldConverters, Entries: make([]Identifier, 0), c: c.c}
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
func (c *cacheEntriesMap) getItem(ID ID) (Identifier, bool) {
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
	ID
}

type ID string

func (ID ID) _ID() ID {
	return ID
}

// NewVersionedID creates a new versioned ID.
func NewVersionedID(version, id string) VersionedID {
	return VersionedID{Version: version, ID: ID(id)}
}

// Cache interface.
type Cache interface {
	GetOrCreate(ID VersionedID, create func() (Identifier, error)) (Identifier, error)
}

// OpenCloserCache is a cache with a Open/Close lifecycle.
type OpenCloserCache interface {
	Cache
	Open() error
	Close() error
}

type persistentCache struct {
	elemType reflect.Type
	isPtr    bool

	data *cacheEntriesMap

	typem *typeMapper

	openInit sync.Once
	openErr  error
	open     bool

	filename string
	sync.RWMutex
}

// New creates a new Cache with the provided absolute filename and entity. Note
// that the entity used must match exactly the type used when doing the actual saves.
func New(filename string, entity interface{}) OpenCloserCache {
	elementType := reflect.TypeOf(entity)
	isPtr := false
	if elementType.Kind() == reflect.Ptr {
		isPtr = true
	}

	return &persistentCache{filename: filename, elemType: elementType, isPtr: isPtr, typem: newDefaultTypeMapper()}
}

func (c *persistentCache) newCacheEntries() *cacheEntries {
	return &cacheEntries{c: c}
}

func (c *persistentCache) newCacheEntriesMap() *cacheEntriesMap {
	return &cacheEntriesMap{c: c, FieldConverters: make(map[string]string), Entries: make(map[ID]*flaggedIdentifier)}
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
			}
			c.openErr = e
			return
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

// We want the array entries separated nicely on each line to get easy diffing.
// Go's JSON marshaler doesn't seem to support this.
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

	if err = os.MkdirAll(filepath.Dir(c.filename), os.FileMode(0755)); err != nil {
		return err
	}

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

	fmt.Println(">> S", s)
	_, err = f.WriteString(s)

	return err
}

// GetOrCreate fetches the value from the versioned data store. If not, found
// it is created with the supplied create func and put there.
func (c *persistentCache) GetOrCreate(vID VersionedID, create func() (Identifier, error)) (Identifier, error) {
	if err := c.lazyOpen(); err != nil {
		return nil, err
	}

	c.RLock()

	if c.data.Version != vID.Version {
		// Version upgrade.
		c.RUnlock()

		c.Lock()
		defer c.Unlock()
		c.data.Version = vID.Version
		c.data.Entries = make(map[ID]*flaggedIdentifier)
		v, err := c.createIdentifier(create)
		if err != nil {
			return nil, err
		}

		fi := newFlaggedIdentifier(v, true)
		fi.IsTouched = true
		c.data.Entries[vID.ID] = fi

		return v, nil

	}

	v, found := c.data.getItem(vID.ID)
	c.RUnlock()

	if found {
		return v, nil
	}

	c.Lock()
	defer c.Unlock()
	if v, found = c.data.getItem(vID.ID); found {
		return v, nil
	}

	vc, err := c.createIdentifier(create)
	if err != nil {
		return nil, err
	}

	fi := newFlaggedIdentifier(vc, true)
	fi.IsTouched = true
	c.data.Entries[vID.ID] = fi

	return fi.Identifier, nil
}

type mapTypeWalker struct {
	typeConverter typeConverters

	keys []string

	currFieldname    string
	currMapFieldname string

	// The result goes here.
	namedConverters map[string]string
}

func (c *persistentCache) createIdentifier(f func() (Identifier, error)) (Identifier, error) {
	identifier, err := f()
	if err != nil {
		return nil, err
	}

	c.typem.mapTypes(identifier)
	for k, v := range c.typem.namedConverters {
		c.data.FieldConverters[k] = v
	}

	return identifier, nil
}

type noCache struct {
}

func (c *noCache) GetOrCreate(ID VersionedID, create func() (Identifier, error)) (Identifier, error) {
	return create()
}
