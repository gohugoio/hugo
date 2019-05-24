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
	"bytes"
	"encoding/json"
	"fmt"
	"go/doc"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync"

	"github.com/gohugoio/hugo/deps"
)

// TemplateFuncsNamespaceRegistry describes a registry of functions that provide
// namespaces.
var TemplateFuncsNamespaceRegistry []func(d *deps.Deps) *TemplateFuncsNamespace

// AddTemplateFuncsNamespace adds a given function to a registry.
func AddTemplateFuncsNamespace(ns func(d *deps.Deps) *TemplateFuncsNamespace) {
	TemplateFuncsNamespaceRegistry = append(TemplateFuncsNamespaceRegistry, ns)
}

// TemplateFuncsNamespace represents a template function namespace.
type TemplateFuncsNamespace struct {
	// The namespace name, "strings", "lang", etc.
	Name string

	// This is the method receiver.
	Context func(v ...interface{}) interface{}

	// Additional info, aliases and examples, per method name.
	MethodMappings map[string]TemplateFuncMethodMapping
}

// TemplateFuncsNamespaces is a slice of TemplateFuncsNamespace.
type TemplateFuncsNamespaces []*TemplateFuncsNamespace

// AddMethodMapping adds a method to a template function namespace.
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

// TemplateFuncMethodMapping represents a mapping of functions to methods for a
// given namespace.
type TemplateFuncMethodMapping struct {
	Method interface{}

	// Any template funcs aliases. This is mainly motivated by keeping
	// backwards compatibility, but some new template funcs may also make
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

type goDocFunc struct {
	Name        string
	Description string
	Args        []string
	Aliases     []string
	Examples    [][2]string
}

func (t goDocFunc) toJSON() ([]byte, error) {
	args, err := json.Marshal(t.Args)
	if err != nil {
		return nil, err
	}
	aliases, err := json.Marshal(t.Aliases)
	if err != nil {
		return nil, err
	}
	examples, err := json.Marshal(t.Examples)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf(`%q:
    { "Description": %q, "Args": %s, "Aliases": %s, "Examples": %s }	
`, t.Name, t.Description, args, aliases, examples))

	return buf.Bytes(), nil
}

// MarshalJSON returns the JSON encoding of namespaces.
func (namespaces TemplateFuncsNamespaces) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteString("{")

	for i, ns := range namespaces {
		if i != 0 {
			buf.WriteString(",")
		}
		b, err := ns.toJSON()
		if err != nil {
			return nil, err
		}
		buf.Write(b)
	}

	buf.WriteString("}")

	return buf.Bytes(), nil
}

func (t *TemplateFuncsNamespace) toJSON() ([]byte, error) {

	var buf bytes.Buffer

	godoc := getGetTplPackagesGoDoc()[t.Name]

	var funcs []goDocFunc

	buf.WriteString(fmt.Sprintf(`%q: {`, t.Name))

	ctx := t.Context()
	ctxType := reflect.TypeOf(ctx)
	for i := 0; i < ctxType.NumMethod(); i++ {
		method := ctxType.Method(i)
		f := goDocFunc{
			Name: method.Name,
		}

		methodGoDoc := godoc[method.Name]

		if mapping, ok := t.MethodMappings[method.Name]; ok {
			f.Aliases = mapping.Aliases
			f.Examples = mapping.Examples
			f.Description = methodGoDoc.Description
			f.Args = methodGoDoc.Args
		}

		funcs = append(funcs, f)
	}

	for i, f := range funcs {
		if i != 0 {
			buf.WriteString(",")
		}
		funcStr, err := f.toJSON()
		if err != nil {
			return nil, err
		}
		buf.Write(funcStr)
	}

	buf.WriteString("}")

	return buf.Bytes(), nil
}

type methodGoDocInfo struct {
	Description string
	Args        []string
}

var (
	tplPackagesGoDoc     map[string]map[string]methodGoDocInfo
	tplPackagesGoDocInit sync.Once
)

func getGetTplPackagesGoDoc() map[string]map[string]methodGoDocInfo {
	tplPackagesGoDocInit.Do(func() {
		tplPackagesGoDoc = make(map[string]map[string]methodGoDocInfo)
		pwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		fset := token.NewFileSet()

		// pwd will be inside one of the namespace packages during tests
		var basePath string
		if strings.Contains(pwd, "tpl") {
			basePath = filepath.Join(pwd, "..")
		} else {
			basePath = filepath.Join(pwd, "tpl")
		}

		files, err := ioutil.ReadDir(basePath)
		if err != nil {
			log.Fatal(err)
		}

		for _, fi := range files {
			if !fi.IsDir() {
				continue
			}

			namespaceDoc := make(map[string]methodGoDocInfo)
			packagePath := filepath.Join(basePath, fi.Name())

			d, err := parser.ParseDir(fset, packagePath, nil, parser.ParseComments)
			if err != nil {
				log.Fatal(err)
			}

			for _, f := range d {
				p := doc.New(f, "./", 0)

				for _, t := range p.Types {
					if t.Name == "Namespace" {
						for _, tt := range t.Methods {
							var args []string
							for _, p := range tt.Decl.Type.Params.List {
								for _, pp := range p.Names {
									args = append(args, pp.Name)
								}
							}

							description := strings.TrimSpace(tt.Doc)
							di := methodGoDocInfo{Description: description, Args: args}
							namespaceDoc[tt.Name] = di
						}
					}
				}
			}

			tplPackagesGoDoc[fi.Name()] = namespaceDoc
		}
	})

	return tplPackagesGoDoc
}
