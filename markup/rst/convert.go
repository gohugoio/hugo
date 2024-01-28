// Copyright 2019 The Hugo Authors. All rights reserved.
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

// Package rst converts content to HTML using the RST external helper.
package rst

import (
	"bytes"
	"runtime"

	"github.com/gohugoio/hugo/common/hexec"
	"github.com/gohugoio/hugo/htesting"

	"github.com/gohugoio/hugo/identity"

	"github.com/gohugoio/hugo/markup/converter"
	"github.com/gohugoio/hugo/markup/internal"
)

// Provider is the package entry point.
var Provider converter.ProviderProvider = provider{}

type provider struct{}

func (p provider) New(cfg converter.ProviderConfig) (converter.Provider, error) {
	return converter.NewProvider("rst", func(ctx converter.DocumentContext) (converter.Converter, error) {
		return &rstConverter{
			ctx: ctx,
			cfg: cfg,
		}, nil
	}), nil
}

type rstConverter struct {
	ctx converter.DocumentContext
	cfg converter.ProviderConfig
}

func (c *rstConverter) Convert(ctx converter.RenderContext) (converter.ResultRender, error) {
	b, err := c.getRstContent(ctx.Src, c.ctx)
	if err != nil {
		return nil, err
	}
	return converter.Bytes(b), nil
}

func (c *rstConverter) Supports(feature identity.Identity) bool {
	return false
}

// getRstContent calls the Python script rst2html as an external helper
// to convert reStructuredText content to HTML.
func (c *rstConverter) getRstContent(src []byte, ctx converter.DocumentContext) ([]byte, error) {
	logger := c.cfg.Logger
	binaryName, binaryPath := getRstBinaryNameAndPath()

	if binaryName == "" {
		logger.Println("rst2html / rst2html.py not found in $PATH: Please install.\n",
			"                 Leaving reStructuredText content unrendered.")
		return src, nil
	}

	logger.Infoln("Rendering", ctx.DocumentName, "with", binaryName, "...")

	var result []byte
	var err error

	// certain *nix based OSs wrap executables in scripted launchers
	// invoking binaries on these OSs via python interpreter causes SyntaxError
	// invoke directly so that shebangs work as expected
	// handle Windows manually because it doesn't do shebangs
	if runtime.GOOS == "windows" {
		pythonBinary, _ := internal.GetPythonBinaryAndExecPath()
		args := []string{binaryPath, "--leave-comments", "--initial-header-level=2"}
		result, err = internal.ExternallyRenderContent(c.cfg, ctx, src, pythonBinary, args)
	} else {
		args := []string{"--leave-comments", "--initial-header-level=2"}
		result, err = internal.ExternallyRenderContent(c.cfg, ctx, src, binaryName, args)
	}

	if err != nil {
		return nil, err
	}

	// TODO(bep) check if rst2html has a body only option.
	bodyStart := bytes.Index(result, []byte("<body>\n"))
	if bodyStart < 0 {
		bodyStart = -7 // compensate for length
	}

	bodyEnd := bytes.Index(result, []byte("\n</body>"))
	if bodyEnd < 0 || bodyEnd >= len(result) {
		bodyEnd = len(result) - 1
		if bodyEnd < 0 {
			bodyEnd = 0
		}
	}

	return result[bodyStart+7 : bodyEnd], err
}

var rst2Binaries = []string{"rst2html", "rst2html.py"}

func getRstBinaryNameAndPath() (string, string) {
	for _, candidate := range rst2Binaries {
		if pth := hexec.LookPath(candidate); pth != "" {
			return candidate, pth
		}
	}
	return "", ""
}

// Supports returns whether rst is (or should be) installed on this computer.
func Supports() bool {
	name, _ := getRstBinaryNameAndPath()
	hasBin := name != ""
	if htesting.SupportsAll() {
		if !hasBin {
			panic("rst not installed")
		}
		return true
	}
	return hasBin
}
