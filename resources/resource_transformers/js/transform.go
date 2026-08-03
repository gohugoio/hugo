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

	"github.com/evanw/esbuild/pkg/api"
	"github.com/gohugoio/hugo/common/hmaps"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/internal/js/esbuild"
	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/resources/internal"
	"github.com/gohugoio/hugo/resources/resource"
)

type buildTransformation struct {
	optsm map[string]any
	c     *Client
}

func (t *buildTransformation) Key() internal.ResourceTransformationKey {
	key := "jsbuild"
	if t.c.c.CssMode {
		key = "cssbuild"
	}
	return internal.NewResourceTransformationKey(key, t.optsm)
}

func (t *buildTransformation) Transform(ctx *resources.ResourceTransformationCtx) error {
	ctx.OutMediaType = media.Builtin.JavascriptType
	if ctx.InMediaType == media.Builtin.CSSType {
		ctx.OutMediaType = media.Builtin.CSSType
	}

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
		ctx.ReplaceOutPathExtension(ctx.OutMediaType.FirstSuffix.FullSuffix)
	}

	src, err := io.ReadAll(ctx.From)
	if err != nil {
		return err
	}

	opts.SourceDir = filepath.FromSlash(path.Dir(ctx.SourcePath))
	opts.Contents = string(src)
	opts.MediaType = ctx.InMediaType
	opts.Stdin = true
	opts.IsCSS = t.c.c.CssMode
	var ic resource.ResourceGetter
	if opts.ImportContext != nil {
		ic = resource.NewCachedResourceGetter(opts.ImportContext)
	}

	if ic != nil {
		resolved := hmaps.NewCache[string, resource.Resource]()
		opts.ImportOnResolveFunc = func(imp string, args api.OnResolveArgs) string {
			if r := esbuild.ResolveResource(imp, ic); r != nil {
				p := esbuild.PrefixHugoVirtual + resources.InternalResourceTargetPath(r)
				resolved.Set(p, r)
				ctx.DependencyManager.AddIdentity(identity.FirstIdentity(r))
				return p
			}
			return ""
		}
		opts.ImportOnLoadFunc = func(args api.OnLoadArgs) (string, error) {
			if r, found := resolved.Get(args.Path); found {
				return resources.InternalResourceSourceContent(ctx.Ctx, r)
			}
			return "", nil
		}
	}

	_, err = t.c.transform(opts, ctx)

	return err
}
