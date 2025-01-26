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

package cssjs

import (
	"bytes"
	"io"
	"regexp"
	"strings"

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/hexec"
	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/resources/internal"
	"github.com/gohugoio/hugo/resources/internal/vendor"
	"github.com/gohugoio/hugo/resources/resource"
	"github.com/mitchellh/mapstructure"
)

var (
	tailwindcssImportRe   = regexp.MustCompile(`^tailwindcss/?`)
	tailwindImportExclude = func(s string) bool {
		return tailwindcssImportRe.MatchString(s) && !strings.Contains(s, ".")
	}
)

// NewTailwindCSSClient creates a new TailwindCSSClient with the given specification.
func NewTailwindCSSClient(rs *resources.Spec) *TailwindCSSClient {
	return &TailwindCSSClient{rs: rs}
}

// Client is the client used to do TailwindCSS transformations.
type TailwindCSSClient struct {
	rs *resources.Spec
}

// Process transforms the given Resource with the TailwindCSS processor.
func (c *TailwindCSSClient) Process(res resources.ResourceTransformer, options map[string]any) (resource.Resource, error) {
	return res.Transform(&tailwindcssTransformation{rs: c.rs, optionsm: options})
}

var _ vendor.Vendorable = (*tailwindcssTransformation)(nil)

type tailwindcssTransformation struct {
	optionsm map[string]any
	rs       *resources.Spec
}

func (t *tailwindcssTransformation) VendorName() string {
	return "css/tailwindcss"
}

func (t *tailwindcssTransformation) VendorKey() string {
	return vendor.VendorKeyFromOpts(t.optionsm)
}

func (t *tailwindcssTransformation) Key() internal.ResourceTransformationKey {
	return internal.NewResourceTransformationKey("tailwindcss", t.optionsm)
}

type TailwindCSSOptions struct {
	Minify        bool // Optimize and minify the output
	Optimize      bool //  Optimize the output without minifying
	InlineImports `mapstructure:",squash"`
}

func (opts TailwindCSSOptions) toArgs() []any {
	var args []any
	if opts.Minify {
		args = append(args, "--minify")
	}
	if opts.Optimize {
		args = append(args, "--optimize")
	}
	return args
}

func (t *tailwindcssTransformation) Transform(ctx *resources.ResourceTransformationCtx) error {
	const binaryName = "tailwindcss"

	options, err := decodeTailwindCSSOptions(t.optionsm)
	if err != nil {
		return err
	}

	infol := t.rs.Logger.InfoCommand(binaryName)
	infow := loggers.LevelLoggerToWriter(infol)

	ex := t.rs.ExecHelper

	workingDir := t.rs.Cfg.BaseConfig().WorkingDir

	var cmdArgs []any = []any{
		"--input=-", // Read from stdin.
		"--cwd", workingDir,
	}

	cmdArgs = append(cmdArgs, options.toArgs()...)

	var errBuf bytes.Buffer

	stderr := io.MultiWriter(infow, &errBuf)
	cmdArgs = append(cmdArgs, hexec.WithStderr(stderr))
	cmdArgs = append(cmdArgs, hexec.WithStdout(ctx.To))
	cmdArgs = append(cmdArgs, hexec.WithEnviron(hugo.GetExecEnviron(workingDir, t.rs.Cfg, t.rs.BaseFs.Assets.Fs)))

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

	imp := newImportResolver(
		ctx.From,
		ctx.InPath,
		options.InlineImports,
		t.rs.Assets.Fs, t.rs.Logger, ctx.DependencyManager,
	)

	src, err := imp.resolve()
	if err != nil {
		return err
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

func decodeTailwindCSSOptions(m map[string]any) (opts TailwindCSSOptions, err error) {
	if m == nil {
		return
	}
	err = mapstructure.WeakDecode(m, &opts)
	return
}
