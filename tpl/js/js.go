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

// Package js provides functions for building JavaScript resources
package js

import (
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/resources/resource"
	"github.com/gohugoio/hugo/resources/resource_transformers/js"
	"github.com/gohugoio/hugo/tpl/internal/resourcehelpers"
)

// New returns a new instance of the js-namespaced template functions.
func New(deps *deps.Deps) *Namespace {
	if deps.ResourceSpec == nil {
		return &Namespace{}
	}
	return &Namespace{
		client: js.New(deps.BaseFs.Assets, deps.ResourceSpec),
	}
}

// Namespace provides template functions for the "js" namespace.
type Namespace struct {
	client *js.Client
}

// Build processes the given Resource with ESBuild.
func (ns *Namespace) Build(args ...any) (resource.Resource, error) {
	var (
		r          resources.ResourceTransformer
		m          map[string]any
		targetPath string
		err        error
		ok         bool
	)

	r, targetPath, ok = resourcehelpers.ResolveIfFirstArgIsString(args)

	if !ok {
		r, m, err = resourcehelpers.ResolveArgs(args)
		if err != nil {
			return nil, err
		}
	}

	if targetPath != "" {
		m = map[string]any{"targetPath": targetPath}
	}

	return ns.client.Process(r, m)
}
