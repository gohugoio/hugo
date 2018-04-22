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

package hugolib

import (
	"fmt"
	"sync"
)

type orderedMap struct {
	sync.RWMutex
	keys []interface{}
	m    map[interface{}]interface{}
}

func newOrderedMap() *orderedMap {
	return &orderedMap{m: make(map[interface{}]interface{})}
}

func newOrderedMapFromStringMapString(m map[string]string) *orderedMap {
	om := newOrderedMap()
	for k, v := range m {
		om.Add(k, v)
	}
	return om
}

func (m *orderedMap) Add(k, v interface{}) {
	m.Lock()
	defer m.Unlock()
	_, found := m.m[k]
	if found {
		panic(fmt.Sprintf("%v already added", v))
	}
	m.m[k] = v
	m.keys = append(m.keys, k)
}

func (m *orderedMap) Get(k interface{}) (interface{}, bool) {
	m.RLock()
	defer m.RUnlock()
	v, found := m.m[k]
	return v, found
}

func (m *orderedMap) Contains(k interface{}) bool {
	m.RLock()
	defer m.RUnlock()
	_, found := m.m[k]
	return found
}

func (m *orderedMap) Keys() []interface{} {
	m.RLock()
	defer m.RUnlock()
	return m.keys
}

func (m *orderedMap) Len() int {
	m.RLock()
	defer m.RUnlock()
	return len(m.keys)
}

// Some shortcuts for known types.
func (m *orderedMap) getShortcode(k interface{}) *shortcode {
	v, found := m.Get(k)
	if !found {
		return nil
	}
	return v.(*shortcode)
}

func (m *orderedMap) getShortcodeRenderer(k interface{}) func() (string, error) {
	v, found := m.Get(k)
	if !found {
		return nil
	}
	return v.(func() (string, error))
}

func (m *orderedMap) getString(k interface{}) string {
	v, found := m.Get(k)
	if !found {
		return ""
	}
	return v.(string)
}
