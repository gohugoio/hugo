// Copyright 2017 The Hugo Authors. All rights reserved.
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

package reflect

import (
	"fmt"
	"reflect"
)

// New returns a new instance of the reflect-namespaced template functions.
func New() *Namespace {
	return &Namespace{}
}

// Namespace provides template functions for the "reflect" namespace.
type Namespace struct{}

// KindIs reports whether v is of kind k.
func (ns *Namespace) KindIs(k string, v interface{}) bool {
	return ns.KindOf(v) == k
}

// KindOf reports v's kind.
func (ns *Namespace) KindOf(v interface{}) string {
	return reflect.ValueOf(v).Kind().String()
}

// TypeIs reports whether v is of type t.
func (ns *Namespace) TypeIs(t string, v interface{}) bool {
	return ns.TypeOf(v) == t
}

// TypeIsLike reports whether v is of type t or a pointer to type t.
func (ns *Namespace) TypeIsLike(t string, v interface{}) bool {
	s := ns.TypeOf(v)
	return s == t || s == "*"+t
}

// TypeOf reports v's type.
func (ns *Namespace) TypeOf(v interface{}) string {
	return fmt.Sprintf("%T", v)
}
