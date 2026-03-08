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

package cssjs

import (
	"io"

	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/resources/internal"
	"github.com/gohugoio/hugo/resources/resource"
	"github.com/mitchellh/mapstructure"
)

// NewInlineImportsClient creates a new InlineImportsClient with the given specification.
func NewInlineImportsClient(rs *resources.Spec) *InlineImportsClient {
	return &InlineImportsClient{rs: rs}
}

// InlineImportsClient is the client used to do CSS inline import transformations.
type InlineImportsClient struct {
	rs *resources.Spec
}

// Process transforms the given Resource by inlining CSS @import statements.
func (c *InlineImportsClient) Process(res resources.ResourceTransformer, options map[string]any) (resource.Resource, error) {
	return res.Transform(&inlineImportsTransformation{rs: c.rs, optionsm: options})
}

type inlineImportsTransformation struct {
	optionsm map[string]any
	rs       *resources.Spec
}

func (t *inlineImportsTransformation) Key() internal.ResourceTransformationKey {
	return internal.NewResourceTransformationKey("cssimports", t.optionsm)
}

// InlineImportsOptions holds options for the css.InlineImports template function.
type InlineImportsOptions struct {
	SkipInlineImportsNotFound bool
}

func (t *inlineImportsTransformation) Transform(ctx *resources.ResourceTransformationCtx) error {
	opts, err := decodeInlineImportsOptions(t.optionsm)
	if err != nil {
		return err
	}

	imp := newImportResolver(
		ctx.From,
		ctx.InPath,
		InlineImports{
			InlineImports:             true,
			SkipInlineImportsNotFound: opts.SkipInlineImportsNotFound,
		},
		t.rs.Assets.Fs, t.rs.Logger, ctx.DependencyManager,
	)

	resolved, err := imp.resolve()
	if err != nil {
		return err
	}

	_, err = io.Copy(ctx.To, resolved)
	return err
}

func decodeInlineImportsOptions(m map[string]any) (opts InlineImportsOptions, err error) {
	if m == nil {
		return
	}
	err = mapstructure.WeakDecode(m, &opts)
	return
}
