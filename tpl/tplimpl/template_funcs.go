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

package tplimpl

import (
	"html/template"

	"github.com/gohugoio/hugo/deps"

	"github.com/gohugoio/hugo/tpl/internal"

	// Init the namespaces
	_ "github.com/gohugoio/hugo/tpl/cast"
	_ "github.com/gohugoio/hugo/tpl/collections"
	_ "github.com/gohugoio/hugo/tpl/compare"
	_ "github.com/gohugoio/hugo/tpl/crypto"
	_ "github.com/gohugoio/hugo/tpl/data"
	_ "github.com/gohugoio/hugo/tpl/encoding"
	_ "github.com/gohugoio/hugo/tpl/fmt"
	_ "github.com/gohugoio/hugo/tpl/hugo"
	_ "github.com/gohugoio/hugo/tpl/images"
	_ "github.com/gohugoio/hugo/tpl/inflect"
	_ "github.com/gohugoio/hugo/tpl/lang"
	_ "github.com/gohugoio/hugo/tpl/math"
	_ "github.com/gohugoio/hugo/tpl/os"
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

func createFuncMap(d *deps.Deps) map[string]interface{} {
	funcMap := template.FuncMap{}

	// Merge the namespace funcs
	for _, nsf := range internal.TemplateFuncsNamespaceRegistry {
		ns := nsf(d)
		if _, exists := funcMap[ns.Name]; exists {
			panic(ns.Name + " is a duplicate template func")
		}
		funcMap[ns.Name] = ns.Context
		for _, mm := range ns.MethodMappings {
			for _, alias := range mm.Aliases {
				if _, exists := funcMap[alias]; exists {
					panic(alias + " is a duplicate template func")
				}
				funcMap[alias] = mm.Method
			}

		}

	}

	return funcMap

}
func (t *templateFuncster) initFuncMap(funcMap template.FuncMap) {
	t.funcMap = funcMap
	t.Tmpl.(*templateHandler).setFuncs(funcMap)
}
