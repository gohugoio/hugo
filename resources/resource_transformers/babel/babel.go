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

package babel

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"path"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/gohugoio/hugo/common/hexec"
	"github.com/gohugoio/hugo/common/loggers"

	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/resources/internal"

	"github.com/mitchellh/mapstructure"

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/resources/resource"
)

// Options from https://babeljs.io/docs/en/options
type Options struct {
	Config string // Custom path to config file

	Minified   bool
	NoComments bool
	Compact    *bool
	Verbose    bool
	NoBabelrc  bool
	SourceMap  string
}

// DecodeOptions decodes options to and generates command flags
func DecodeOptions(m map[string]any) (opts Options, err error) {
	if m == nil {
		return
	}
	err = mapstructure.WeakDecode(m, &opts)
	return
}

func (opts Options) sourceMapEnabled() bool {
	return opts.SourceMap != "" && opts.SourceMap != "none"
}

func (opts Options) toArgs() []any {
	var args []any

	// Always use inline source maps; Transform extracts to external file if needed.
	if opts.sourceMapEnabled() {
		args = append(args, "--source-maps=inline")
	}
	if opts.Minified {
		args = append(args, "--minified")
	}
	if opts.NoComments {
		args = append(args, "--no-comments")
	}
	if opts.Compact != nil {
		args = append(args, "--compact="+strconv.FormatBool(*opts.Compact))
	}
	if opts.Verbose {
		args = append(args, "--verbose")
	}
	if opts.NoBabelrc {
		args = append(args, "--no-babelrc")
	}
	return args
}

// Client is the client used to do Babel transformations.
type Client struct {
	rs *resources.Spec
}

// New creates a new Client with the given specification.
func New(rs *resources.Spec) *Client {
	return &Client{rs: rs}
}

type babelTransformation struct {
	options Options
	rs      *resources.Spec
}

func (t *babelTransformation) Key() internal.ResourceTransformationKey {
	return internal.NewResourceTransformationKey("babel", t.options)
}

// Transform shells out to babel-cli to do the heavy lifting.
// For this to work, you need some additional tools. To install them globally:
// npm install -g @babel/core @babel/cli
// If you want to use presets or plugins such as @babel/preset-env
// Then you should install those globally as well. e.g:
// npm install -g @babel/preset-env
// Instead of installing globally, you can also install everything as a dev-dependency (--save-dev instead of -g)
func (t *babelTransformation) Transform(ctx *resources.ResourceTransformationCtx) error {
	const binaryName = "babel"

	ex := t.rs.ExecHelper

	if err := ex.Sec().CheckAllowedExec(binaryName); err != nil {
		return err
	}

	var configFile string
	infol := t.rs.Logger.InfoCommand(binaryName)
	infoW := loggers.LevelLoggerToWriter(infol)

	var errBuf bytes.Buffer

	if t.options.Config != "" {
		configFile = t.options.Config
	} else {
		configFile = "babel.config.js"
	}

	configFile = filepath.Clean(configFile)

	// We need an absolute filename to the config file.
	if !filepath.IsAbs(configFile) {
		configFile = t.rs.BaseFs.ResolveJSConfigFile(configFile)
		if configFile == "" && t.options.Config != "" {
			// Only fail if the user specified config file is not found.
			return fmt.Errorf("babel config %q not found", configFile)
		}
	}

	ctx.ReplaceOutPathExtension(".js")

	var cmdArgs []any

	if configFile != "" {
		infol.Logf("use config file %q", configFile)
		cmdArgs = []any{"--config-file", configFile}
	}

	if optArgs := t.options.toArgs(); len(optArgs) > 0 {
		cmdArgs = append(cmdArgs, optArgs...)
	}
	cmdArgs = append(cmdArgs, "--filename="+ctx.SourcePath)

	var outBuf bytes.Buffer
	stderr := io.MultiWriter(infoW, &errBuf)
	cmdArgs = append(cmdArgs, hexec.WithStderr(stderr))
	cmdArgs = append(cmdArgs, hexec.WithStdout(&outBuf))
	cmdArgs = append(cmdArgs, hexec.WithDir(t.rs.Cfg.BaseConfig().WorkingDir))
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

	go func() {
		defer stdin.Close()
		io.Copy(stdin, ctx.From)
	}()

	err = cmd.Run()
	if err != nil {
		if hexec.IsNotFound(err) {
			return &herrors.FeatureNotAvailableError{Cause: err}
		}
		return fmt.Errorf(errBuf.String()+": %w", err)
	}

	content := outBuf.Bytes()

	if t.options.SourceMap == "external" {
		if sourceMap, stripped, ok := extractInlineSourceMap(content); ok {
			if err = ctx.PublishSourceMap(sourceMap); err != nil {
				return err
			}
			targetPath := path.Base(ctx.OutPath) + ".map"
			content = append(stripped, ("//# sourceMappingURL=" + targetPath + "\n")...)
		}
	}

	ctx.To.Write(content)

	return nil
}

// Process transforms the given Resource with the Babel processor.
func (c *Client) Process(res resources.ResourceTransformer, options Options) (resource.Resource, error) {
	return res.Transform(
		&babelTransformation{rs: c.rs, options: options},
	)
}

var inlineSourceMapRe = regexp.MustCompile(`(?m)//# sourceMappingURL=data:application/json[^,]*,([A-Za-z0-9+/=]+)\s*$`)

func extractInlineSourceMap(content []byte) (sourceMap, stripped []byte, ok bool) {
	loc := inlineSourceMapRe.FindSubmatchIndex(content)
	if loc == nil {
		return nil, content, false
	}
	decoded, err := base64.StdEncoding.DecodeString(string(content[loc[2]:loc[3]]))
	if err != nil {
		return nil, content, false
	}
	return decoded, content[:loc[0]], true
}
