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

package js

import (
	"io"
	"path"
	"path/filepath"

	"github.com/gohugoio/hugo/internal/js/esbuild"
	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/resources/internal"
	"github.com/gohugoio/hugo/resources/internal/vendor"
)

type buildTransformation struct {
	optsm map[string]any
	c     *Client
}

var _ vendor.Vendorable = (*buildTransformation)(nil)

func (t *buildTransformation) VendorName() string {
	return "js/build"
}

func (t *buildTransformation) VendorScope() map[string]any {
	return vendor.VendorScopeFromOpts(t.optsm)
}

func (t *buildTransformation) Key() internal.ResourceTransformationKey {
	return internal.NewResourceTransformationKey("jsbuild", t.optsm)
}

func (t *buildTransformation) Transform(ctx *resources.ResourceTransformationCtx) error {
	ctx.OutMediaType = media.Builtin.JavascriptType

	var opts esbuild.Options

	if t.optsm != nil {
		optsExt, err := esbuild.DecodeExternalOptions(t.optsm)
		if err != nil {
			return err
		}
		opts.ExternalOptions = optsExt
	}

	if opts.TargetPath != "" {
		ctx.OutPath = opts.TargetPath
	} else {
		ctx.ReplaceOutPathExtension(".js")
	}

	src, err := io.ReadAll(ctx.From)
	if err != nil {
		return err
	}

	opts.SourceDir = filepath.FromSlash(path.Dir(ctx.SourcePath))
	opts.Contents = string(src)
	opts.MediaType = ctx.InMediaType
	opts.Stdin = true

	_, err = t.c.transform(opts, ctx)

	return err
}
