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
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/common/text"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/hugolib/filesystems"
	"github.com/gohugoio/hugo/identity"

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

// Process processes a resource with the user provided options.
func (c *Client) Process(res resources.ResourceTransformer, opts map[string]any) (resource.Resource, error) {
	return res.Transform(
		&buildTransformation{c: c, optsm: opts},
	)
}

func (c *Client) BuildBundle(opts Options) (api.BuildResult, error) {
	return c.build(opts, nil)
}

// Note that transformCtx may be nil.
func (c *Client) build(opts Options, transformCtx *resources.ResourceTransformationCtx) (api.BuildResult, error) {
	dependencyManager := opts.DependencyManager
	if transformCtx != nil {
		dependencyManager = transformCtx.DependencyManager // TODO1
	}
	if dependencyManager == nil {
		dependencyManager = identity.NopManager
	}

	opts.ResolveDir = c.rs.Cfg.BaseConfig().WorkingDir // where node_modules gets resolved
	opts.TsConfig = c.rs.ResolveJSConfigFile("tsconfig.json")

	if err := opts.validate(); err != nil {
		return api.BuildResult{}, err
	}

	buildOptions, err := toBuildOptions(opts)
	if err != nil {
		return api.BuildResult{}, err
	}

	buildOptions.Plugins, err = createBuildPlugins(c, dependencyManager, opts)
	if err != nil {
		return api.BuildResult{}, err
	}

	if buildOptions.Sourcemap == api.SourceMapExternal && buildOptions.Outdir == "" {
		buildOptions.Outdir, err = os.MkdirTemp(os.TempDir(), "compileOutput")
		if err != nil {
			return api.BuildResult{}, err
		}
		defer os.Remove(buildOptions.Outdir)
	}

	if opts.Inject != nil {
		// Resolve the absolute filenames.
		for i, ext := range opts.Inject {
			impPath := filepath.FromSlash(ext)
			if filepath.IsAbs(impPath) {
				return api.BuildResult{}, fmt.Errorf("inject: absolute paths not supported, must be relative to /assets")
			}

			m := resolveComponentInAssets(c.rs.Assets.Fs, impPath)

			if m == nil {
				return api.BuildResult{}, fmt.Errorf("inject: file %q not found", ext)
			}

			opts.Inject[i] = m.Filename

		}

		buildOptions.Inject = opts.Inject

	}

	result := api.Build(buildOptions)

	if len(result.Errors) > 0 {
		createErr := func(msg api.Message) error {
			if msg.Location == nil {
				return errors.New(msg.Text)
			}
			var (
				contentr     hugio.ReadSeekCloser
				errorMessage string
				loc          = msg.Location
				errorPath    = loc.File
				err          error
			)

			var resolvedError *ErrorMessageResolved

			if opts.ErrorMessageResolveFunc != nil {
				resolvedError = opts.ErrorMessageResolveFunc(msg)
			}

			if resolvedError == nil {
				if errorPath == stdinImporter {
					errorPath = transformCtx.SourcePath
				}

				errorMessage = msg.Text
				// TODO1 handle all namespaces, make more general
				errorMessage = strings.ReplaceAll(errorMessage, NsHugoImport+":", "")

				if strings.HasPrefix(errorPath, NsHugoImport) {
					errorPath = strings.TrimPrefix(errorPath, NsHugoImport+":")
					contentr, err = hugofs.Os.Open(errorPath)
				} else {
					var fi os.FileInfo
					fi, err = c.sfs.Fs.Stat(errorPath)
					if err == nil {
						m := fi.(hugofs.FileMetaInfo).Meta()
						errorPath = m.Filename
						contentr, err = m.Open()
					}
				}
			} else {
				contentr = resolvedError.Content
				errorPath = resolvedError.Path
				errorMessage = resolvedError.Message
			}

			if contentr != nil {
				defer contentr.Close()
			}

			if err == nil {
				fe := herrors.
					NewFileErrorFromName(errors.New(errorMessage), errorPath).
					UpdatePosition(text.Position{Offset: -1, LineNumber: loc.Line, ColumnNumber: loc.Column}).
					UpdateContent(contentr, nil)

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
				c.rs.Logger.Errorf("js.Build failed: %s", err)
			}
		}

		return result, errors[0]
	}

	// TODO1 option etc.	fmt.Printf("%s", api.AnalyzeMetafile(result.Metafile, api.AnalyzeMetafileOptions{}))

	if transformCtx != nil {
		if buildOptions.Sourcemap == api.SourceMapExternal {
			content := string(result.OutputFiles[1].Contents)
			symPath := path.Base(transformCtx.OutPath) + ".map"
			re := regexp.MustCompile(`//# sourceMappingURL=.*\n?`)
			content = re.ReplaceAllString(content, "//# sourceMappingURL="+symPath+"\n")

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
	}

	return result, nil
}
