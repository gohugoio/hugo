// Copyright 2017-present The Hugo Authors. All rights reserved.
//
// Portions Copyright The Go Authors.

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

package internal

import (
	"encoding/json"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/spf13/hugo/deps"
)

var TemplateFuncsNamespaceRegistry []func(d *deps.Deps) *TemplateFuncsNamespace

func AddTemplateFuncsNamespace(ns func(d *deps.Deps) *TemplateFuncsNamespace) {
	TemplateFuncsNamespaceRegistry = append(TemplateFuncsNamespaceRegistry, ns)
}

type TemplateFuncsNamespace struct {
	// The namespace name, "strings", "lang", etc.
	Name string

	// This is the method receiver.
	Context interface{}

	// Additional info, aliases and examples, per method name.
	MethodMappings map[string]TemplateFuncMethodMapping
}

func (t *TemplateFuncsNamespace) AddMethodMapping(m interface{}, aliases []string, examples [][2]string) {
	if t.MethodMappings == nil {
		t.MethodMappings = make(map[string]TemplateFuncMethodMapping)
	}

	name := methodToName(m)

	// sanity check
	for _, e := range examples {
		if e[0] == "" {
			panic(t.Name + ": Empty example for " + name)
		}
	}
	for _, a := range aliases {
		if a == "" {
			panic(t.Name + ": Empty alias for " + name)
		}
	}

	t.MethodMappings[name] = TemplateFuncMethodMapping{
		Method:   m,
		Aliases:  aliases,
		Examples: examples,
	}

}

type TemplateFuncMethodMapping struct {
	Method interface{}

	// Any template funcs aliases. This is mainly motivated by keeping
	// backwards compability, but some new template funcs may also make
	// sense to give short and snappy aliases.
	// Note that these aliases are global and will be merged, so the last
	// key will win.
	Aliases []string

	// A slice of input/expected examples.
	// We keep it a the namespace level for now, but may find a way to keep track
	// of the single template func, for documentation purposes.
	// Some of these, hopefully just a few, may depend on some test data to run.
	Examples [][2]string
}

func methodToName(m interface{}) string {
	name := runtime.FuncForPC(reflect.ValueOf(m).Pointer()).Name()
	name = filepath.Ext(name)
	name = strings.TrimPrefix(name, ".")
	name = strings.TrimSuffix(name, "-fm")
	return name
}

func (t *TemplateFuncsNamespace) MarshalJSON() ([]byte, error) {
	type Func struct {
		Name        string
		Description string // TODO(bep)
		Aliases     []string
		Examples    [][2]string
	}
	// TODO(bep) map/lookup from docs template Namespace + Func name.
	var funcs []Func

	ctx := t.Context.(func() interface{})()
	ctxType := reflect.TypeOf(ctx)
	for i := 0; i < ctxType.NumMethod(); i++ {
		method := ctxType.Method(i)
		f := Func{
			Name: method.Name,
		}
		if mapping, ok := t.MethodMappings[method.Name]; ok {
			f.Aliases = mapping.Aliases
			f.Examples = mapping.Examples
		}
		funcs = append(funcs, f)
	}

	return json.Marshal(&struct {
		Name  string
		Funcs []Func
	}{
		Name:  t.Name,
		Funcs: funcs,
	})
}
