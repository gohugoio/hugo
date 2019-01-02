// Copyright 2018 The Hugo Authors. All rights reserved.
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

package postcss

import (
	"io"
	"path/filepath"

	"github.com/gohugoio/hugo/hugofs"
	"github.com/pkg/errors"

	"os"
	"os/exec"

	"github.com/mitchellh/mapstructure"

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/resources/resource"
)

// Some of the options from https://github.com/postcss/postcss-cli
type Options struct {

	// Set a custom path to look for a config file.
	Config string

	NoMap bool `mapstructure:"no-map"` // Disable the default inline sourcemaps

	// Options for when not using a config file
	Use         string // List of postcss plugins to use
	Parser      string //  Custom postcss parser
	Stringifier string // Custom postcss stringifier
	Syntax      string // Custom postcss syntax
}

func DecodeOptions(m map[string]interface{}) (opts Options, err error) {
	if m == nil {
		return
	}
	err = mapstructure.WeakDecode(m, &opts)
	return
}

func (opts Options) toArgs() []string {
	var args []string
	if opts.NoMap {
		args = append(args, "--no-map")
	}
	if opts.Use != "" {
		args = append(args, "--use", opts.Use)
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

// Client is the client used to do PostCSS transformations.
type Client struct {
	rs *resources.Spec
}

// New creates a new Client with the given specification.
func New(rs *resources.Spec) *Client {
	return &Client{rs: rs}
}

type postcssTransformation struct {
	options Options
	rs      *resources.Spec
}

func (t *postcssTransformation) Key() resources.ResourceTransformationKey {
	return resources.NewResourceTransformationKey("postcss", t.options)
}

// Transform shells out to postcss-cli to do the heavy lifting.
// For this to work, you need some additional tools. To install them globally:
// npm install -g postcss-cli
// npm install -g autoprefixer
func (t *postcssTransformation) Transform(ctx *resources.ResourceTransformationCtx) error {

	const localPostCSSPath = "node_modules/postcss-cli/bin/"
	const binaryName = "postcss"

	// Try first in the project's node_modules.
	csiBinPath := filepath.Join(t.rs.WorkingDir, localPostCSSPath, binaryName)

	binary := csiBinPath

	if _, err := exec.LookPath(binary); err != nil {
		// Try PATH
		binary = binaryName
		if _, err := exec.LookPath(binary); err != nil {
			// This may be on a CI server etc. Will fall back to pre-built assets.
			return herrors.ErrFeatureNotAvailable
		}
	}

	var configFile string
	logger := t.rs.Logger

	if t.options.Config != "" {
		configFile = t.options.Config
	} else {
		configFile = "postcss.config.js"
	}

	configFile = filepath.Clean(configFile)

	// We need an abolute filename to the config file.
	if !filepath.IsAbs(configFile) {
		// We resolve this against the virtual Work filesystem, to allow
		// this config file to live in one of the themes if needed.
		fi, err := t.rs.BaseFs.Work.Fs.Stat(configFile)
		if err != nil {
			if t.options.Config != "" {
				// Only fail if the user specificed config file is not found.
				return errors.Wrapf(err, "postcss config %q not found:", configFile)
			}
			configFile = ""
		} else {
			configFile = fi.(hugofs.RealFilenameInfo).RealFilename()
		}
	}

	var cmdArgs []string

	if configFile != "" {
		logger.INFO.Println("postcss: use config file", configFile)
		cmdArgs = []string{"--config", configFile}
	}

	if optArgs := t.options.toArgs(); len(optArgs) > 0 {
		cmdArgs = append(cmdArgs, optArgs...)
	}

	cmd := exec.Command(binary, cmdArgs...)

	cmd.Stdout = ctx.To
	cmd.Stderr = os.Stderr

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
		return err
	}

	return nil
}

// Process transforms the given Resource with the PostCSS processor.
func (c *Client) Process(res resource.Resource, options Options) (resource.Resource, error) {
	return c.rs.Transform(
		res,
		&postcssTransformation{rs: c.rs, options: options},
	)
}
