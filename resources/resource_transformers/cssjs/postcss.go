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

// Package cssjs provides resource transformations backed by some popular JS based frameworks.
package cssjs

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/gohugoio/hugo/common/collections"
	"github.com/gohugoio/hugo/common/hexec"
	"github.com/gohugoio/hugo/common/loggers"

	"github.com/gohugoio/hugo/common/hugo"

	"github.com/gohugoio/hugo/resources/internal"
	"github.com/spf13/cast"

	"github.com/mitchellh/mapstructure"

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/resources/resource"
)

// NewPostCSSClient creates a new PostCSSClient with the given specification.
func NewPostCSSClient(rs *resources.Spec) *PostCSSClient {
	return &PostCSSClient{rs: rs}
}

func decodePostCSSOptions(m map[string]any) (opts PostCSSOptions, err error) {
	if m == nil {
		return
	}
	err = mapstructure.WeakDecode(m, &opts)

	if !opts.NoMap {
		// There was for a long time a discrepancy between documentation and
		// implementation for the noMap property, so we need to support both
		// camel and snake case.
		opts.NoMap = cast.ToBool(m["no-map"])
	}

	return
}

// PostCSSClient is the client used to do PostCSS transformations.
type PostCSSClient struct {
	rs *resources.Spec
}

// Process transforms the given Resource with the PostCSS processor.
func (c *PostCSSClient) Process(res resources.ResourceTransformer, options map[string]any) (resource.Resource, error) {
	return res.Transform(&postcssTransformation{rs: c.rs, optionsm: options})
}

type InlineImports struct {
	// Service `mapstructure:",squash"`
	// Enable inlining of @import statements.
	// Does so recursively, but currently once only per file;
	// that is, it's not possible to import the same file in
	// different scopes (root, media query...)
	// Note that this import routine does not care about the CSS spec,
	// so you can have @import anywhere in the file.
	InlineImports bool

	// When InlineImports is enabled, we fail the build if an import cannot be resolved.
	// You can enable this to allow the build to continue and leave the import statement in place.
	// Note that the inline importer does not process url location or imports with media queries,
	// so those will be left as-is even without enabling this option.
	SkipInlineImportsNotFound bool
}

// Some of the options from https://github.com/postcss/postcss-cli
type PostCSSOptions struct {
	// Set a custom path to look for a config file.
	Config string

	NoMap bool // Disable the default inline sourcemaps

	InlineImports `mapstructure:",squash"`

	// Options for when not using a config file
	Use         string // List of postcss plugins to use
	Parser      string //  Custom postcss parser
	Stringifier string // Custom postcss stringifier
	Syntax      string // Custom postcss syntax
}

func (opts PostCSSOptions) toArgs() []string {
	var args []string
	if opts.NoMap {
		args = append(args, "--no-map")
	}
	if opts.Use != "" {
		args = append(args, "--use")
		args = append(args, strings.Fields(opts.Use)...)
	}
	if opts.Parser != "" {
		args = append(args, "--parser", opts.Parser)
	}
	if opts.Stringifier != "" {
		args = append(args, "--stringifier", opts.Stringifier)
	}
	if opts.Syntax != "" {
		args = append(args, "--syntax", opts.Syntax)
	}
	return args
}

type postcssTransformation struct {
	optionsm map[string]any
	rs       *resources.Spec
}

func (t *postcssTransformation) Key() internal.ResourceTransformationKey {
	return internal.NewResourceTransformationKey("postcss", t.optionsm)
}

// Transform shells out to postcss-cli to do the heavy lifting.
// For this to work, you need some additional tools. To install them globally:
// npm install -g postcss-cli
// npm install -g autoprefixer
func (t *postcssTransformation) Transform(ctx *resources.ResourceTransformationCtx) error {
	const binaryName = "postcss"

	infol := t.rs.Logger.InfoCommand(binaryName)
	infow := loggers.LevelLoggerToWriter(infol)

	ex := t.rs.ExecHelper

	var configFile string

	options, err := decodePostCSSOptions(t.optionsm)
	if err != nil {
		return err
	}

	if options.Config != "" {
		configFile = options.Config
	} else {
		configFile = "postcss.config.js"
	}

	configFile = filepath.Clean(configFile)

	// We need an absolute filename to the config file.
	if !filepath.IsAbs(configFile) {
		configFile = t.rs.BaseFs.ResolveJSConfigFile(configFile)
		if configFile == "" && options.Config != "" {
			// Only fail if the user specified config file is not found.
			return fmt.Errorf("postcss config %q not found", options.Config)
		}
	}

	var cmdArgs []any

	if configFile != "" {
		infol.Logf("use config file %q", configFile)
		cmdArgs = []any{"--config", configFile}
	}

	if optArgs := options.toArgs(); len(optArgs) > 0 {
		cmdArgs = append(cmdArgs, collections.StringSliceToInterfaceSlice(optArgs)...)
	}

	var errBuf bytes.Buffer

	stderr := io.MultiWriter(infow, &errBuf)
	cmdArgs = append(cmdArgs, hexec.WithStderr(stderr))
	cmdArgs = append(cmdArgs, hexec.WithStdout(ctx.To))
	cmdArgs = append(cmdArgs, hexec.WithEnviron(hugo.GetExecEnviron(t.rs.Cfg.BaseConfig().WorkingDir, t.rs.Cfg, t.rs.BaseFs.Assets.Fs)))

	cmd, err := ex.Npx(binaryName, cmdArgs...)
	if err != nil {
		if hexec.IsNotFound(err) {
			// This may be on a CI server etc. Will fall back to pre-built assets.
			return &herrors.FeatureNotAvailableError{Cause: err}
		}
		return err
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	src := ctx.From

	imp := newImportResolver(
		ctx.From,
		ctx.InPath,
		options.InlineImports,
		t.rs.Assets.Fs, t.rs.Logger, ctx.DependencyManager,
	)

	if options.InlineImports.InlineImports {
		var err error
		src, err = imp.resolve()
		if err != nil {
			return err
		}
	}

	go func() {
		defer stdin.Close()
		io.Copy(stdin, src)
	}()

	err = cmd.Run()
	if err != nil {
		if hexec.IsNotFound(err) {
			return &herrors.FeatureNotAvailableError{
				Cause: err,
			}
		}
		return imp.toFileError(errBuf.String())
	}

	return nil
}
