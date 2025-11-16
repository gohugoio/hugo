// Copyright 2025 The Hugo Authors. All rights reserved.
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

package tplimplinit

import (
	// Init the template funcs namespaces
	"context"
	"html/template"

	"github.com/gohugoio/hugo/deps"
	_ "github.com/gohugoio/hugo/tpl/cast"
	_ "github.com/gohugoio/hugo/tpl/collections"
	_ "github.com/gohugoio/hugo/tpl/compare"
	_ "github.com/gohugoio/hugo/tpl/crypto"
	_ "github.com/gohugoio/hugo/tpl/css"
	_ "github.com/gohugoio/hugo/tpl/data"
	_ "github.com/gohugoio/hugo/tpl/debug"
	_ "github.com/gohugoio/hugo/tpl/diagrams"
	_ "github.com/gohugoio/hugo/tpl/encoding"
	_ "github.com/gohugoio/hugo/tpl/fmt"
	_ "github.com/gohugoio/hugo/tpl/hash"
	_ "github.com/gohugoio/hugo/tpl/hugo"
	_ "github.com/gohugoio/hugo/tpl/images"
	_ "github.com/gohugoio/hugo/tpl/inflect"
	"github.com/gohugoio/hugo/tpl/internal"
	_ "github.com/gohugoio/hugo/tpl/js"
	_ "github.com/gohugoio/hugo/tpl/lang"
	_ "github.com/gohugoio/hugo/tpl/math"
	_ "github.com/gohugoio/hugo/tpl/openapi/openapi3"
	_ "github.com/gohugoio/hugo/tpl/os"
	_ "github.com/gohugoio/hugo/tpl/page"
	_ "github.com/gohugoio/hugo/tpl/partials"
	_ "github.com/gohugoio/hugo/tpl/path"
	_ "github.com/gohugoio/hugo/tpl/reflect"
	_ "github.com/gohugoio/hugo/tpl/resources"
	_ "github.com/gohugoio/hugo/tpl/safe"
	_ "github.com/gohugoio/hugo/tpl/site"
	_ "github.com/gohugoio/hugo/tpl/strings"
	_ "github.com/gohugoio/hugo/tpl/templates"
	_ "github.com/gohugoio/hugo/tpl/time"
	_ "github.com/gohugoio/hugo/tpl/transform"
	_ "github.com/gohugoio/hugo/tpl/urls"
)

// CreateFuncMap creates a template.FuncMap with all of Hugo's template funcs,
// excluding the Go built-ins.
func CreateFuncMap(d *deps.Deps) map[string]any {
	funcMap := template.FuncMap{}
	nsMap := make(map[string]any)
	var onCreated []func(namespaces map[string]any)

	// Merge the namespace funcs
	for _, nsf := range internal.TemplateFuncsNamespaceRegistry {
		ns := nsf(d)
		if _, exists := funcMap[ns.Name]; exists {
			panic(ns.Name + " is a duplicate template func")
		}
		funcMap[ns.Name] = ns.Context
		contextV, err := ns.Context(context.Background())
		if err != nil {
			panic(err)
		}
		nsMap[ns.Name] = contextV
		for _, mm := range ns.MethodMappings {
			for _, alias := range mm.Aliases {
				if _, exists := funcMap[alias]; exists {
					panic(alias + " is a duplicate template func")
				}
				funcMap[alias] = mm.Method
			}
		}

		if ns.OnCreated != nil {
			onCreated = append(onCreated, ns.OnCreated)
		}
	}

	for _, f := range onCreated {
		f(nsMap)
	}

	return funcMap
}
