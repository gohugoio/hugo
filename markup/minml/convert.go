// Copyright 2022 The Hugo Authors. All rights reserved.
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

// Package minml converts MinML to HTML.
package minml

import (
	"bytes"

	"github.com/gohugoio/hugo/identity"

	"github.com/gohugoio/hugo/markup/converter"

	"github.com/dedis/matchertext/go/xml/minml"
	"github.com/dedis/matchertext/go/xml/ast"
)

// Provider is the package entry point.
var Provider converter.ProviderProvider = provide{}

type provide struct{}

func (p provide) New(cfg converter.ProviderConfig) (converter.Provider, error) {
	return converter.NewProvider("minml",
		func(ctx converter.DocumentContext) (
			converter.Converter, error) {

		return &minmlConverter{
			ctx: ctx,
			cfg: cfg,
		}, nil
	}), nil
}

type minmlConverter struct {
	ctx converter.DocumentContext
	cfg converter.ProviderConfig
}

type minmlResult struct {
	converter.Result
}

var converterIdentity = identity.KeyValueIdentity{Key: "minml", Value: "converter"}

func (c *minmlConverter) Convert(ctx converter.RenderContext) (
	result converter.Result, err error) {

	// Parse the MinML input to a slice of markup nodes
	ns, err := minml.Parse(bytes.NewReader(ctx.Src))
	if err != nil {
		return nil, err
	}

	// Write resulting AST to a result buffer
	buf := &bytes.Buffer{}
	enc := ast.NewEncoder(buf)
	if err := enc.Encode(ns); err != nil {
		return nil, err
	}

	return minmlResult{Result: buf}, nil
}

var featureSet = map[identity.Identity]bool{
}

func (c *minmlConverter) Supports(feature identity.Identity) bool {
	return featureSet[feature.GetIdentity()]
}

