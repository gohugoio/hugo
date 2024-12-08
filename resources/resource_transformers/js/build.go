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
	"path"
	"regexp"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/gohugoio/hugo/hugolib/filesystems"
	"github.com/gohugoio/hugo/internal/js/esbuild"

	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/resources/resource"
)

// Client context for ESBuild.
type Client struct {
	c *esbuild.BuildClient
}

// New creates a new client context.
func New(fs *filesystems.SourceFilesystem, rs *resources.Spec) *Client {
	return &Client{
		c: esbuild.NewBuildClient(fs, rs),
	}
}

// Process processes a resource with the user provided options.
func (c *Client) Process(res resources.ResourceTransformer, opts map[string]any) (resource.Resource, error) {
	return res.Transform(
		&buildTransformation{c: c, optsm: opts},
	)
}

func (c *Client) transform(opts esbuild.Options, transformCtx *resources.ResourceTransformationCtx) (api.BuildResult, error) {
	if transformCtx.DependencyManager != nil {
		opts.DependencyManager = transformCtx.DependencyManager
	}

	opts.StdinSourcePath = transformCtx.SourcePath

	result, err := c.c.Build(opts)
	if err != nil {
		return result, err
	}

	if opts.ExternalOptions.SourceMap == "linked" || opts.ExternalOptions.SourceMap == "external" {
		content := string(result.OutputFiles[1].Contents)
		if opts.ExternalOptions.SourceMap == "linked" {
			symPath := path.Base(transformCtx.OutPath) + ".map"
			re := regexp.MustCompile(`//# sourceMappingURL=.*\n?`)
			content = re.ReplaceAllString(content, "//# sourceMappingURL="+symPath+"\n")
		}

		if err = transformCtx.PublishSourceMap(string(result.OutputFiles[0].Contents)); err != nil {
			return result, err
		}
		_, err := transformCtx.To.Write([]byte(content))
		if err != nil {
			return result, err
		}
	} else {
		_, err := transformCtx.To.Write(result.OutputFiles[0].Contents)
		if err != nil {
			return result, err
		}

	}
	return result, nil
}
