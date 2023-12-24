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

package js

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/afero"

	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/media"

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/text"

	"github.com/gohugoio/hugo/hugolib/filesystems"
	"github.com/gohugoio/hugo/resources/internal"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/resources/resource"
)

// Client context for ESBuild.
type Client struct {
	rs  *resources.Spec
	sfs *filesystems.SourceFilesystem
}

// New creates a new client context.
func New(fs *filesystems.SourceFilesystem, rs *resources.Spec) *Client {
	return &Client{
		rs:  rs,
		sfs: fs,
	}
}

type buildTransformation struct {
	optsm map[string]any
	c     *Client
}

func (t *buildTransformation) Key() internal.ResourceTransformationKey {
	return internal.NewResourceTransformationKey("jsbuild", t.optsm)
}

func (t *buildTransformation) Transform(ctx *resources.ResourceTransformationCtx) error {
	ctx.OutMediaType = media.Builtin.JavascriptType

	opts, err := decodeOptions(t.optsm)
	if err != nil {
		return err
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

	opts.sourceDir = filepath.FromSlash(path.Dir(ctx.SourcePath))
	opts.resolveDir = t.c.rs.Cfg.BaseConfig().WorkingDir // where node_modules gets resolved
	opts.contents = string(src)
	opts.mediaType = ctx.InMediaType
	opts.tsConfig = t.c.rs.ResolveJSConfigFile("tsconfig.json")

	buildOptions, err := toBuildOptions(opts)
	if err != nil {
		return err
	}

	buildOptions.Plugins, err = createBuildPlugins(ctx.DependencyManager, t.c, opts)
	if err != nil {
		return err
	}

	if buildOptions.Sourcemap == api.SourceMapExternal && buildOptions.Outdir == "" {
		buildOptions.Outdir, err = os.MkdirTemp(os.TempDir(), "compileOutput")
		if err != nil {
			return err
		}
		defer os.Remove(buildOptions.Outdir)
	}

	if opts.Inject != nil {
		// Resolve the absolute filenames.
		for i, ext := range opts.Inject {
			impPath := filepath.FromSlash(ext)
			if filepath.IsAbs(impPath) {
				return fmt.Errorf("inject: absolute paths not supported, must be relative to /assets")
			}

			m := resolveComponentInAssets(t.c.rs.Assets.Fs, impPath)

			if m == nil {
				return fmt.Errorf("inject: file %q not found", ext)
			}

			opts.Inject[i] = m.Filename

		}

		buildOptions.Inject = opts.Inject

	}

	result := api.Build(buildOptions)

	if len(result.Errors) > 0 {

		createErr := func(msg api.Message) error {
			loc := msg.Location
			if loc == nil {
				return errors.New(msg.Text)
			}
			path := loc.File
			if path == stdinImporter {
				path = ctx.SourcePath
			}

			errorMessage := msg.Text
			errorMessage = strings.ReplaceAll(errorMessage, nsImportHugo+":", "")

			var (
				f   afero.File
				err error
			)

			if strings.HasPrefix(path, nsImportHugo) {
				path = strings.TrimPrefix(path, nsImportHugo+":")
				f, err = hugofs.Os.Open(path)
			} else {
				var fi os.FileInfo
				fi, err = t.c.sfs.Fs.Stat(path)
				if err == nil {
					m := fi.(hugofs.FileMetaInfo).Meta()
					path = m.Filename
					f, err = m.Open()
				}

			}

			if err == nil {
				fe := herrors.
					NewFileErrorFromName(errors.New(errorMessage), path).
					UpdatePosition(text.Position{Offset: -1, LineNumber: loc.Line, ColumnNumber: loc.Column}).
					UpdateContent(f, nil)

				f.Close()
				return fe
			}

			return fmt.Errorf("%s", errorMessage)
		}

		var errors []error

		for _, msg := range result.Errors {
			errors = append(errors, createErr(msg))
		}

		// Return 1, log the rest.
		for i, err := range errors {
			if i > 0 {
				t.c.rs.Logger.Errorf("js.Build failed: %s", err)
			}
		}

		return errors[0]
	}

	if buildOptions.Sourcemap == api.SourceMapExternal {
		content := string(result.OutputFiles[1].Contents)
		symPath := path.Base(ctx.OutPath) + ".map"
		re := regexp.MustCompile(`//# sourceMappingURL=.*\n?`)
		content = re.ReplaceAllString(content, "//# sourceMappingURL="+symPath+"\n")

		if err = ctx.PublishSourceMap(string(result.OutputFiles[0].Contents)); err != nil {
			return err
		}
		_, err := ctx.To.Write([]byte(content))
		if err != nil {
			return err
		}
	} else {
		_, err := ctx.To.Write(result.OutputFiles[0].Contents)
		if err != nil {
			return err
		}
	}
	return nil
}

// Process process esbuild transform
func (c *Client) Process(res resources.ResourceTransformer, opts map[string]any) (resource.Resource, error) {
	return res.Transform(
		&buildTransformation{c: c, optsm: opts},
	)
}
