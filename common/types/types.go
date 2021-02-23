// Copyright 2019 The Hugo Authors. All rights reserved.
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

// Package types contains types shared between packages in Hugo.
package types

import (
	"fmt"
	"reflect"

	"github.com/spf13/cast"
)

// RLocker represents the read locks in sync.RWMutex.
type RLocker interface {
	RLock()
	RUnlock()
}

// KeyValueStr is a string tuple.
type KeyValueStr struct {
	Key   string
	Value string
}

// KeyValues holds an key and a slice of values.
type KeyValues struct {
	Key    interface{}
	Values []interface{}
}

// KeyString returns the key as a string, an empty string if conversion fails.
func (k KeyValues) KeyString() string {
	return cast.ToString(k.Key)
}

func (k KeyValues) String() string {
	return fmt.Sprintf("%v: %v", k.Key, k.Values)
}

// NewKeyValuesStrings takes a given key and slice of values and returns a new
// KeyValues struct.
func NewKeyValuesStrings(key string, values ...string) KeyValues {
	iv := make([]interface{}, len(values))
	for i := 0; i < len(values); i++ {
		iv[i] = values[i]
	}
	return KeyValues{Key: key, Values: iv}
}

// Zeroer, as implemented by time.Time, will be used by the truth template
// funcs in Hugo (if, with, not, and, or).
type Zeroer interface {
	IsZero() bool
}

// IsNil reports whether v is nil.
func IsNil(v interface{}) bool {
	if v == nil {
		return true
	}

	value := reflect.ValueOf(v)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return value.IsNil()
	}

	return false
}

// DevMarker is a marker interface for types that should only be used during
// development.
type DevMarker interface {
	DevOnly()
}
