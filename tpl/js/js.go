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
	"errors"
	"fmt"

	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/resources/resource"
	"github.com/gohugoio/hugo/resources/resource_transformers/js"
	_errors "github.com/pkg/errors"
)

// New returns a new instance of the js-namespaced template functions.
func New(deps *deps.Deps) *Namespace {
	return &Namespace{
		client: js.New(deps.BaseFs.Assets, deps.ResourceSpec),
	}
}

// Namespace provides template functions for the "js" namespace.
type Namespace struct {
	deps   *deps.Deps
	client *js.Client
}

// Build processes the given Resource with ESBuild.
func (ns *Namespace) Build(args ...interface{}) (resource.Resource, error) {
	r, m, err := ns.resolveArgs(args)
	if err != nil {
		return nil, err
	}
	var options js.Options
	if m != nil {
		options, err = js.DecodeOptions(m)

		if err != nil {
			return nil, err
		}
	}

	return ns.client.Process(r, options)

}

// This roundabout way of doing it is needed to get both pipeline behaviour and options as arguments.
// This is a copy of tpl/resources/resolveArgs
func (ns *Namespace) resolveArgs(args []interface{}) (resources.ResourceTransformer, map[string]interface{}, error) {
	if len(args) == 0 {
		return nil, nil, errors.New("no Resource provided in transformation")
	}

	if len(args) == 1 {
		r, ok := args[0].(resources.ResourceTransformer)
		if !ok {
			return nil, nil, fmt.Errorf("type %T not supported in Resource transformations", args[0])
		}
		return r, nil, nil
	}

	r, ok := args[1].(resources.ResourceTransformer)
	if !ok {
		if _, ok := args[1].(map[string]interface{}); !ok {
			return nil, nil, fmt.Errorf("no Resource provided in transformation")
		}
		return nil, nil, fmt.Errorf("type %T not supported in Resource transformations", args[0])
	}

	m, err := maps.ToStringMapE(args[0])
	if err != nil {
		return nil, nil, _errors.Wrap(err, "invalid options type")
	}

	return r, m, nil
}
