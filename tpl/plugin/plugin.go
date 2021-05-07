// Copyright 2018 The Hugo Authors. All rights reserved.
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

// Package path provides template functions for manipulating paths.
package plugin

import (
	"errors"
	"fmt"
	"plugin"
	"reflect"

	"github.com/gohugoio/hugo/deps"
	"github.com/spf13/cast"
)

// New returns a new instance of the path-namespaced template functions.
func New(deps *deps.Deps) *Namespace {
	return &Namespace{
		deps: deps,
	}
}

// Namespace provides template functions for the "os" namespace.
type Namespace struct {
	deps *deps.Deps
}

type ErrPluginNotFound struct {
	Name string
}

func (err *ErrPluginNotFound) Error() string {
	return fmt.Sprintf(`plugin "%s" not found`, err.Name)
}

// Open returns the loaded plugin named name.
func (ns *Namespace) Open(name interface{}) (*plugin.Plugin, error) {
	sname, err := cast.ToStringE(name)
	if err != nil {
		return nil, err
	}

	plugins := ns.deps.Site.Plugin()

	if p, ok := plugins[sname]; ok {
		return p.(*plugin.Plugin), nil
	}

	return nil, &ErrPluginNotFound{
		Name: sname,
	}
}

// Get returns the loaded plugin named name.
func (ns *Namespace) Get(pluginName, fieldName interface{}) (interface{}, error) {
	spluginName, err := cast.ToStringE(pluginName)
	if err != nil {
		return nil, err
	}

	plugins := ns.deps.Site.Plugin()

	p, ok := plugins[spluginName]
	if !ok {
		return nil, &ErrPluginNotFound{
			Name: spluginName,
		}
	}

	sfieldName, err := cast.ToStringE(fieldName)
	if err != nil {
		return nil, err
	}

	symbol, err := p.(*plugin.Plugin).Lookup(sfieldName)

	return symbol, err
}

func (ns *Namespace) Call(symbol interface{}, arguments ...interface{}) (interface{}, error) {
	fn := reflect.ValueOf(symbol)
	if fn.Kind() == reflect.Func {
		in := make([]reflect.Value, len(arguments))

		for i, arg := range arguments {
			in[i] = reflect.ValueOf(arg)
		}

		fnType := fn.Type()
		for i := 0; i<fnType.NumIn(); i++ {
			argType := fnType.In(i)
			if !in[i].Type().AssignableTo(argType) {
				return nil, fmt.Errorf(`Unsuported type %s for argument %d ("%v"): expected %s`, in[i].Type().Name(), i+1, arguments[i], argType.Name())
			}
		}

		result := fn.Call(in)

		switch len(result) {
			case 0:
				return nil, nil
			case 1:
				if !result[0].CanInterface() {
					return nil, errors.New("invalid signature")
				}

				return result[0].Interface(), nil
			case 2:
				if !result[0].CanInterface() {
					return nil, errors.New("invalid signature")
				}
				r0 := result[0].Interface()

				if result[1].IsNil() {
					return r0, nil
				}

				if !result[1].CanInterface() {
					return nil, errors.New("invalid signature")
				}

				r1, ok := result[1].Interface().(error)
				if !ok {
					return nil, errors.New("invalid signature")
				}

				return r0, r1
			default:
				return nil, errors.New("invalid signature")
		}
	}

	return nil, errors.New("invalid argument")
}

// Exist returns true if the plugin exist.
func (ns *Namespace) Exist(name interface{}) (bool, error) {
	sname, err := cast.ToStringE(name)
	if err != nil {
		return false, err
	}

	plugins := ns.deps.Site.Plugin()

	_, ok := plugins[sname]
	return ok, nil
}


// Has returns true if the symbol exists in plugin.
func (ns *Namespace) Has(pluginName, fieldName interface{}) (bool, error) {
	p, err := ns.Open(pluginName)
	if err != nil {
		return false, err
	}

	sfieldName, err := cast.ToStringE(fieldName)
	if err != nil {
		return false, err
	}

	_, err = p.Lookup(sfieldName)

	return err == nil, nil
}
