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

	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/resources/internal"
)

type buildTransformation struct {
	optsm map[string]any
	opts  Options
	c     *Client
}

func (t *buildTransformation) Key() internal.ResourceTransformationKey {
	// Pick the most stable key source.
	var v any = t.optsm
	if v == nil {
		v = t.opts
	}
	return internal.NewResourceTransformationKey("jsbuild", v)
}

func (t *buildTransformation) Transform(ctx *resources.ResourceTransformationCtx) error {
	ctx.OutMediaType = media.Builtin.JavascriptType

	if t.optsm != nil {
		optsExt, err := decodeOptions(t.optsm)
		if err != nil {
			return err
		}
		t.opts.ExternalOptions = optsExt
	}

	if t.opts.TargetPath != "" {
		ctx.OutPath = t.opts.TargetPath
	} else {
		ctx.ReplaceOutPathExtension(".js")
	}

	src, err := io.ReadAll(ctx.From)
	if err != nil {
		return err
	}

	t.opts.SourceDir = filepath.FromSlash(path.Dir(ctx.SourcePath))
	t.opts.Contents = string(src)
	t.opts.MediaType = ctx.InMediaType
	t.opts.Stdin = true

	_, err = t.c.build(t.opts, ctx)

	return err
}
