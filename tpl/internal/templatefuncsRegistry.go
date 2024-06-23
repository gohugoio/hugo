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
	"context"
	"encoding/json"
	"fmt"
	"go/doc"
	"go/parser"
	"go/token"
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
	Context func(ctx context.Context, v ...any) (any, error)

	// OnCreated is called when all the namespaces are ready.
	OnCreated func(namespaces map[string]any)

	// Additional info, aliases and examples, per method name.
	MethodMappings map[string]TemplateFuncMethodMapping
}

// TemplateFuncsNamespaces is a slice of TemplateFuncsNamespace.
type TemplateFuncsNamespaces []*TemplateFuncsNamespace

// AddMethodMapping adds a method to a template function namespace.
func (t *TemplateFuncsNamespace) AddMethodMapping(m any, aliases []string, examples [][2]string) {
	if t.MethodMappings == nil {
		t.MethodMappings = make(map[string]TemplateFuncMethodMapping)
	}

	name := methodToName(m)

	// Rewrite §§ to ` in example commands.
	for i, e := range examples {
		examples[i][0] = strings.ReplaceAll(e[0], "§§", "`")
	}

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
	Method any

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

func methodToName(m any) string {
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

// ToMap returns a limited map representation of the namespaces.
func (namespaces TemplateFuncsNamespaces) ToMap() map[string]any {
	m := make(map[string]any)
	for _, ns := range namespaces {
		mm := make(map[string]any)
		for name, mapping := range ns.MethodMappings {
			mm[name] = map[string]any{
				"Examples": mapping.Examples,
				"Aliases":  mapping.Aliases,
			}
		}
		m[ns.Name] = mm
	}
	return m
}

// MarshalJSON returns the JSON encoding of namespaces.
func (namespaces TemplateFuncsNamespaces) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteString("{")

	for i, ns := range namespaces {

		b, err := ns.toJSON(context.Background())
		if err != nil {
			return nil, err
		}
		if b != nil {
			if i != 0 {
				buf.WriteString(",")
			}
			buf.Write(b)
		}
	}

	buf.WriteString("}")

	return buf.Bytes(), nil
}

var ignoreFuncs = map[string]bool{
	"Reset": true,
}

func (t *TemplateFuncsNamespace) toJSON(ctx context.Context) ([]byte, error) {
	var buf bytes.Buffer

	godoc := getGetTplPackagesGoDoc()[t.Name]

	var funcs []goDocFunc

	buf.WriteString(fmt.Sprintf(`%q: {`, t.Name))

	tctx, err := t.Context(ctx)
	if err != nil {
		return nil, err
	}
	if tctx == nil {
		// E.g. page.
		// We should fix this, but we're going to abandon this construct in a little while.
		return nil, nil
	}
	ctxType := reflect.TypeOf(tctx)
	for i := 0; i < ctxType.NumMethod(); i++ {
		method := ctxType.Method(i)
		if ignoreFuncs[method.Name] {
			continue
		}
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

		files, err := os.ReadDir(basePath)
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
