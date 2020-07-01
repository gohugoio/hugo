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

package rollup

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/resources/internal"

	"github.com/mitchellh/mapstructure"

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/resources/resource"
	"github.com/pkg/errors"
)

// Options from https://rollupjs.io/docs/en/options TODO
type Options struct {
	Config string // Custom path to config file

	Verbose bool
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

	if !opts.Verbose {
		args = append(args, "--silent")
	}
	return args
}

// Client is the client used to do Rollup transformations.
type Client struct {
	rs *resources.Spec
}

// New creates a new Client with the given specification.
func New(rs *resources.Spec) *Client {
	return &Client{rs: rs}
}

type rollupTransformation struct {
	options Options
	rs      *resources.Spec
}

func (t *rollupTransformation) Key() internal.ResourceTransformationKey {
	return internal.NewResourceTransformationKey("rollup", t.options)
}

// Transform shells out to rollup to do the heavy lifting.
func (t *rollupTransformation) Transform(ctx *resources.ResourceTransformationCtx) error {
	const localRollupPath = "node_modules/.bin/"
	const binaryName = "rollup"

	// Try first in the project's node_modules.
	csiBinPath := filepath.Join(t.rs.WorkingDir, localRollupPath, binaryName)

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
		configFile = "rollup.config.js"
	}

	configFile = filepath.Clean(configFile)

	// We need an abolute filename to the config file.
	if !filepath.IsAbs(configFile) {
		// We resolve this against the virtual Work filesystem, to allow
		// this config file to live in one of the themes if needed.
		fi, err := t.rs.BaseFs.Work.Stat(configFile)
		if err != nil {
			if t.options.Config != "" {
				// Only fail if the user specificed config file is not found.
				return errors.Wrapf(err, "rollup config %q not found:", configFile)
			}
		} else {
			configFile = fi.(hugofs.FileMetaInfo).Meta().Filename()
		}
	}

	var cmdArgs []string

	if configFile != "" {
		logger.INFO.Println("rollup: use config file", configFile)
		cmdArgs = []string{"-c", configFile}
	}

	if optArgs := t.options.toArgs(); len(optArgs) > 0 {
		cmdArgs = append(cmdArgs, optArgs...)
	}
	cmdArgs = append(cmdArgs, "--input=-")

	cmd := exec.Command(binary, cmdArgs...)

	cmd.Stdout = ctx.To
	cmd.Stderr = os.Stderr
	cmd.Env = hugo.GetExecEnviron(t.rs.Cfg)

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

// Process transforms the given Resource with the Rollup processor.
func (c *Client) Process(res resources.ResourceTransformer, options Options) (resource.Resource, error) {
	return res.Transform(
		&rollupTransformation{rs: c.rs, options: options},
	)
}
