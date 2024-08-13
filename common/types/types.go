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
	"sync/atomic"

	"github.com/spf13/cast"
)

// RLocker represents the read locks in sync.RWMutex.
type RLocker interface {
	RLock()
	RUnlock()
}

// KeyValue is a interface{} tuple.
type KeyValue struct {
	Key   any
	Value any
}

// KeyValueStr is a string tuple.
type KeyValueStr struct {
	Key   string
	Value string
}

// KeyValues holds an key and a slice of values.
type KeyValues struct {
	Key    any
	Values []any
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
	iv := make([]any, len(values))
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
func IsNil(v any) bool {
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

// Unwrapper is implemented by types that can unwrap themselves.
type Unwrapper interface {
	// Unwrapv is for internal use only.
	// It got its slightly odd name to prevent collisions with user types.
	Unwrapv() any
}

// Unwrap returns the underlying value of v if it implements Unwrapper, otherwise v is returned.
func Unwrapv(v any) any {
	if u, ok := v.(Unwrapper); ok {
		return u.Unwrapv()
	}
	return v
}

// LowHigh represents a byte or slice boundary.
type LowHigh[S ~[]byte | string] struct {
	Low  int
	High int
}

func (l LowHigh[S]) IsZero() bool {
	return l.Low < 0 || (l.Low == 0 && l.High == 0)
}

func (l LowHigh[S]) Value(source S) S {
	return source[l.Low:l.High]
}

// This is only used for debugging purposes.
var InvocationCounter atomic.Int64

// NewTrue returns a pointer to b.
func NewBool(b bool) *bool {
	return &b
}

// PrintableValueProvider is implemented by types that can provide a printable value.
type PrintableValueProvider interface {
	PrintableValue() any
}

var _ PrintableValueProvider = Result[any]{}

// Result is a generic result type.
type Result[T any] struct {
	// The result value.
	Value T

	// The error value.
	Err error
}

// PrintableValue returns the value or panics if there is an error.
func (r Result[T]) PrintableValue() any {
	if r.Err != nil {
		panic(r.Err)
	}
	return r.Value
}
