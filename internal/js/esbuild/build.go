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

// Package esbuild provides functions for building JavaScript resources.
package esbuild

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/common/text"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/hugolib/filesystems"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/resources"
)

// NewBuildClient creates a new BuildClient.
func NewBuildClient(fs *filesystems.SourceFilesystem, rs *resources.Spec) *BuildClient {
	return &BuildClient{
		rs:  rs,
		sfs: fs,
	}
}

// BuildClient is a client for building JavaScript resources using esbuild.
type BuildClient struct {
	rs  *resources.Spec
	sfs *filesystems.SourceFilesystem
}

// Build builds the given JavaScript resources using esbuild with the given options.
func (c *BuildClient) Build(opts Options) (api.BuildResult, error) {
	dependencyManager := opts.DependencyManager
	if dependencyManager == nil {
		dependencyManager = identity.NopManager
	}

	opts.OutDir = c.rs.AbsPublishDir
	opts.ResolveDir = c.rs.Cfg.BaseConfig().WorkingDir // where node_modules gets resolved
	opts.AbsWorkingDir = opts.ResolveDir
	opts.TsConfig = c.rs.ResolveJSConfigFile("tsconfig.json")
	assetsResolver := newFSResolver(c.rs.Assets.Fs)

	if err := opts.validate(); err != nil {
		return api.BuildResult{}, err
	}

	if err := opts.compile(); err != nil {
		return api.BuildResult{}, err
	}

	var err error
	opts.compiled.Plugins, err = createBuildPlugins(c.rs, assetsResolver, dependencyManager, opts)
	if err != nil {
		return api.BuildResult{}, err
	}

	if opts.Inject != nil {
		// Resolve the absolute filenames.
		for i, ext := range opts.Inject {
			impPath := filepath.FromSlash(ext)
			if filepath.IsAbs(impPath) {
				return api.BuildResult{}, fmt.Errorf("inject: absolute paths not supported, must be relative to /assets")
			}

			m := assetsResolver.resolveComponent(impPath)

			if m == nil {
				return api.BuildResult{}, fmt.Errorf("inject: file %q not found", ext)
			}

			opts.Inject[i] = m.Filename

		}

		opts.compiled.Inject = opts.Inject

	}

	result := api.Build(opts.compiled)

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
					errorPath = opts.StdinSourcePath
				}

				errorMessage = msg.Text

				var namespace string
				for _, ns := range hugoNamespaces {
					if strings.HasPrefix(errorPath, ns) {
						namespace = ns
						break
					}
				}

				if namespace != "" {
					namespace += ":"
					errorMessage = strings.ReplaceAll(errorMessage, namespace, "")
					errorPath = strings.TrimPrefix(errorPath, namespace)
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

	inOutputPathToAbsFilename := opts.ResolveSourceMapSource
	opts.ResolveSourceMapSource = func(s string) string {
		if inOutputPathToAbsFilename != nil {
			if filename := inOutputPathToAbsFilename(s); filename != "" {
				return filename
			}
		}

		if m := assetsResolver.resolveComponent(s); m != nil {
			return m.Filename
		}

		return ""
	}

	for i, o := range result.OutputFiles {
		if err := fixOutputFile(&o, func(s string) string {
			if s == "<stdin>" {
				return opts.ResolveSourceMapSource(opts.StdinSourcePath)
			}
			var isNsHugo bool
			if strings.HasPrefix(s, "ns-hugo") {
				isNsHugo = true
				idxColon := strings.Index(s, ":")
				s = s[idxColon+1:]
			}

			if !strings.HasPrefix(s, PrefixHugoVirtual) {
				if !filepath.IsAbs(s) {
					s = filepath.Join(opts.OutDir, s)
				}
			}

			if isNsHugo {
				if ss := opts.ResolveSourceMapSource(s); ss != "" {
					if strings.HasPrefix(ss, PrefixHugoMemory) {
						// File not on disk, mark it for removal from the sources slice.
						return ""
					}
					return ss
				}
				return ""
			}
			return s
		}); err != nil {
			return result, err
		}
		result.OutputFiles[i] = o
	}

	return result, nil
}
