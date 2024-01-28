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
	"fmt"
	"io"
	"os"
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

func (opts Options) toArgs() []any {
	var args []any

	// external is not a known constant on the babel command line
	// .sourceMaps must be a boolean, "inline", "both", or undefined
	switch opts.SourceMap {
	case "external":
		args = append(args, "--source-maps")
	case "inline":
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

	// Create compile into a real temp file:
	// 1. separate stdout/stderr messages from babel (https://github.com/gohugoio/hugo/issues/8136)
	// 2. allow generation and retrieval of external source map.
	compileOutput, err := os.CreateTemp("", "compileOut-*.js")
	if err != nil {
		return err
	}

	cmdArgs = append(cmdArgs, "--out-file="+compileOutput.Name())
	stderr := io.MultiWriter(infoW, &errBuf)
	cmdArgs = append(cmdArgs, hexec.WithStderr(stderr))
	cmdArgs = append(cmdArgs, hexec.WithStdout(stderr))
	cmdArgs = append(cmdArgs, hexec.WithEnviron(hugo.GetExecEnviron(t.rs.Cfg.BaseConfig().WorkingDir, t.rs.Cfg, t.rs.BaseFs.Assets.Fs)))

	defer os.Remove(compileOutput.Name())

	// ARGA [--no-install babel --config-file /private/var/folders/_g/j3j21hts4fn7__h04w2x8gb40000gn/T/hugo-test-babel812882892/babel.config.js --source-maps --filename=js/main2.js --out-file=/var/folders/_g/j3j21hts4fn7__h04w2x8gb40000gn/T/compileOut-2237820197.js]
	//      [--no-install babel --config-file /private/var/folders/_g/j3j21hts4fn7__h04w2x8gb40000gn/T/hugo-test-babel332846848/babel.config.js --filename=js/main.js --out-file=/var/folders/_g/j3j21hts4fn7__h04w2x8gb40000gn/T/compileOut-1451390834.js 0x10304ee60 0x10304ed60 0x10304f060]
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

	content, err := io.ReadAll(compileOutput)
	if err != nil {
		return err
	}

	mapFile := compileOutput.Name() + ".map"
	if _, err := os.Stat(mapFile); err == nil {
		defer os.Remove(mapFile)
		sourceMap, err := os.ReadFile(mapFile)
		if err != nil {
			return err
		}
		if err = ctx.PublishSourceMap(string(sourceMap)); err != nil {
			return err
		}
		targetPath := path.Base(ctx.OutPath) + ".map"
		re := regexp.MustCompile(`//# sourceMappingURL=.*\n?`)
		content = []byte(re.ReplaceAllString(string(content), "//# sourceMappingURL="+targetPath+"\n"))
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
