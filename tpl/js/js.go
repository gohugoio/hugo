// Copyright 2024 The Hugo Authors. All rights reserved.
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
	"errors"
	"path"

	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/internal/js/esbuild"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/resources/resource"
	"github.com/gohugoio/hugo/resources/resource_factories/create"
	"github.com/gohugoio/hugo/resources/resource_transformers/babel"
	"github.com/gohugoio/hugo/resources/resource_transformers/js"
	jstransform "github.com/gohugoio/hugo/resources/resource_transformers/js"
	"github.com/gohugoio/hugo/tpl/internal/resourcehelpers"
)

// New returns a new instance of the js-namespaced template functions.
func New(deps *deps.Deps) (*Namespace, error) {
	if deps.ResourceSpec == nil {
		return &Namespace{}, nil
	}

	batcherClient, err := esbuild.NewBatcherClient(deps)
	if err != nil {
		return nil, err
	}

	return &Namespace{
		d:                 deps,
		jsTransformClient: jstransform.New(deps.BaseFs.Assets, deps.ResourceSpec),
		jsBatcherClient:   batcherClient,
		createClient:      create.New(deps.ResourceSpec),
		babelClient:       babel.New(deps.ResourceSpec),
	}, nil
}

// Namespace provides template functions for the "js" namespace.
type Namespace struct {
	d *deps.Deps

	jsTransformClient *js.Client
	createClient      *create.Client
	babelClient       *babel.Client
	jsBatcherClient   *esbuild.BatcherClient
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

	return ns.jsTransformClient.Process(r, m)
}

func (ns *Namespace) Batch(id string, store *maps.Scratch) (esbuild.Batcher, error) {
	key := path.Join(esbuild.NsBatch, id)
	b, err := store.GetOrCreate(key, func() (any, error) {
		return ns.jsBatcherClient.New(id)
	})
	if err != nil {
		return nil, err
	}
	return b.(esbuild.Batcher), nil
}

// Babel processes the given Resource with Babel.
func (ns *Namespace) Babel(args ...any) (resource.Resource, error) {
	if len(args) > 2 {
		return nil, errors.New("must not provide more arguments than resource object and options")
	}

	r, m, err := resourcehelpers.ResolveArgs(args)
	if err != nil {
		return nil, err
	}
	var options babel.Options
	if m != nil {
		options, err = babel.DecodeOptions(m)
		if err != nil {
			return nil, err
		}
	}

	return ns.babelClient.Process(r, options)
}
