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
	"reflect"
	"strings"

	"github.com/gohugoio/hugo/tpl"

	"github.com/gohugoio/hugo/common/maps"

	template "github.com/gohugoio/hugo/tpl/internal/go_templates/htmltemplate"
	texttemplate "github.com/gohugoio/hugo/tpl/internal/go_templates/texttemplate"

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

var _ texttemplate.ExecHelper = (*templateExecHelper)(nil)
var zero reflect.Value

type templateExecHelper struct {
	funcs map[string]reflect.Value
}

func (t *templateExecHelper) GetFunc(tmpl texttemplate.Preparer, name string) (reflect.Value, bool) {
	if fn, found := t.funcs[name]; found {
		return fn, true
	}
	return zero, false
}

func (t *templateExecHelper) GetMapValue(tmpl texttemplate.Preparer, receiver, key reflect.Value) (reflect.Value, bool) {
	if params, ok := receiver.Interface().(maps.Params); ok {
		// Case insensitive.
		keystr := strings.ToLower(key.String())
		v, found := params[keystr]
		if !found {
			return zero, false
		}
		return reflect.ValueOf(v), true
	}

	v := receiver.MapIndex(key)

	return v, v.IsValid()
}

func (t *templateExecHelper) GetMethod(tmpl texttemplate.Preparer, receiver reflect.Value, name string) (method reflect.Value, firstArg reflect.Value) {
	// This is a hot path and receiver.MethodByName really shows up in the benchmarks.
	// Page.Render is the only method with a WithTemplateInfo as of now, so let's just
	// check that for now.
	// TODO(bep) find a more flexible, but still fast, way.
	if name == "Render" {
		if info, ok := tmpl.(tpl.Info); ok {
			if m := receiver.MethodByName(name + "WithTemplateInfo"); m.IsValid() {
				return m, reflect.ValueOf(info)
			}
		}
	}

	return receiver.MethodByName(name), zero
}

func newTemplateExecuter(d *deps.Deps) (texttemplate.Executer, map[string]reflect.Value) {
	funcs := createFuncMap(d)
	funcsv := make(map[string]reflect.Value)

	for k, v := range funcs {
		vv := reflect.ValueOf(v)
		funcsv[k] = vv
	}

	// Duplicate Go's internal funcs here for faster lookups.
	for k, v := range template.GoFuncs {
		if _, exists := funcsv[k]; !exists {
			vv, ok := v.(reflect.Value)
			if !ok {
				vv = reflect.ValueOf(v)
			}
			funcsv[k] = vv
		}
	}

	for k, v := range texttemplate.GoFuncs {
		if _, exists := funcsv[k]; !exists {
			funcsv[k] = v
		}
	}

	exeHelper := &templateExecHelper{
		funcs: funcsv,
	}

	return texttemplate.NewExecuter(
		exeHelper,
	), funcsv
}

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

	if d.OverloadedTemplateFuncs != nil {
		for k, v := range d.OverloadedTemplateFuncs {
			funcMap[k] = v
		}
	}

	return funcMap

}
