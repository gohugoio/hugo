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

// Package debug provides template functions to help debugging templates.
package debug

import (
	"reflect"
	"sort"

	"github.com/sanity-io/litter"
	"github.com/spf13/cast"
	"github.com/yuin/goldmark/util"

	"github.com/gohugoio/hugo/deps"
)

// New returns a new instance of the debug-namespaced template functions.
func New(d *deps.Deps) *Namespace {
	return &Namespace{}
}

// Namespace provides template functions for the "debug" namespace.
type Namespace struct {
}

// Dump returns a object dump of val as a string.
// Note that not every value passed to Dump will print so nicely, but
// we'll improve on that.
//
// We recommend using the "go" Chroma lexer to format the output
// nicely.
//
// Also note that the output from Dump may change from Hugo version to the next,
// so don't depend on a specific output.
func (ns *Namespace) Dump(val any) string {
	return litter.Sdump(val)
}

// List returns a slice of field names and method names of the struct/pointer or keys of the map
// This method scans the provided value shallow, non-recursively.
func (ns *Namespace) List(val any) []string {

	fields := make([]string, 0)
	value := reflect.ValueOf(val)

	if value.Kind() == reflect.Map {
		for _, key := range value.MapKeys() {
			fields = append(fields, key.String())
			sort.Strings(fields)
		}
	}

	// Dereference the pointer if needed
	if value.Kind() == reflect.Pointer {
		value = value.Elem()
	}

	if value.Kind() == reflect.Struct {
		// Iterate over the fields
		for i := 0; i < value.NumField(); i++ {
			field := value.Type().Field(i)

			// Only add exported fields
			if field.PkgPath == "" {
				fields = append(fields, field.Name)
			}
		}

		// Calling NumMethod() on the pointer type returns the number of methods
		// defined for the pointer type as well as the non pointer type.
		// Calling NumMethod() on the non pointer type returns on the other hand only the number of non pointer methods.
		pointerType := reflect.PointerTo(value.Type())

		for i := 0; i < pointerType.NumMethod(); i++ {
			method := pointerType.Method(i)
			fields = append(fields, method.Name)
		}
	}

	return fields
}

// VisualizeSpaces returns a string with spaces replaced by a visible string.
func (ns *Namespace) VisualizeSpaces(val any) string {
	s := cast.ToString(val)
	return string(util.VisualizeSpaces([]byte(s)))
}
