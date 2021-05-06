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
