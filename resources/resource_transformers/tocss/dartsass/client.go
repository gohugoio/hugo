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

// Package godartsass integrates with the Dass Sass Embedded protocol to transpile
// SCSS/SASS.
package dartsass

import (
	"fmt"
	"io"
	"strings"

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/hugolib/filesystems"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/resources/resource"
	"github.com/spf13/afero"

	"github.com/bep/godartsass"
	"github.com/mitchellh/mapstructure"
)

// used as part of the cache key.
const transformationName = "tocss-dart"

// See https://github.com/sass/dart-sass-embedded/issues/24
// Note: This prefix must be all lower case.
const dartSassStdinPrefix = "hugostdin:"

func New(fs *filesystems.SourceFilesystem, rs *resources.Spec) (*Client, error) {
	if !Supports() {
		return &Client{dartSassNotAvailable: true}, nil
	}

	if err := rs.ExecHelper.Sec().CheckAllowedExec(dartSassEmbeddedBinaryName); err != nil {
		return nil, err
	}

	transpiler, err := godartsass.Start(godartsass.Options{
		LogEventHandler: func(event godartsass.LogEvent) {
			message := strings.ReplaceAll(event.Message, dartSassStdinPrefix, "")
			switch event.Type {
			case godartsass.LogEventTypeDebug:
				// Log as Info for now, we may adjust this if it gets too chatty.
				rs.Logger.Infof("Dart Sass: %s", message)
			default:
				// The rest are either deprecations or @warn statements.
				rs.Logger.Warnf("Dart Sass: %s", message)
			}
		},
	})
	if err != nil {
		return nil, err
	}
	return &Client{sfs: fs, workFs: rs.BaseFs.Work, rs: rs, transpiler: transpiler}, nil
}

type Client struct {
	dartSassNotAvailable bool
	rs                   *resources.Spec
	sfs                  *filesystems.SourceFilesystem
	workFs               afero.Fs
	transpiler           *godartsass.Transpiler
}

func (c *Client) ToCSS(res resources.ResourceTransformer, args map[string]any) (resource.Resource, error) {
	if c.dartSassNotAvailable {
		return res.Transform(resources.NewFeatureNotAvailableTransformer(transformationName, args))
	}
	return res.Transform(&transform{c: c, optsm: args})
}

func (c *Client) Close() error {
	if c.transpiler == nil {
		return nil
	}
	return c.transpiler.Close()
}

func (c *Client) toCSS(args godartsass.Args, src io.Reader) (godartsass.Result, error) {
	var res godartsass.Result

	in := helpers.ReaderToString(src)
	args.Source = in

	res, err := c.transpiler.Execute(args)
	if err != nil {
		if err.Error() == "unexpected EOF" {
			return res, fmt.Errorf("got unexpected EOF when executing %q. The user running hugo must have read and execute permissions on this program. With execute permissions only, this error is thrown.", dartSassEmbeddedBinaryName)
		}
		return res, herrors.NewFileErrorFromFileInErr(err, hugofs.Os, herrors.OffsetMatcher)
	}

	return res, err
}

type Options struct {

	// Hugo, will by default, just replace the extension of the source
	// to .css, e.g. "scss/main.scss" becomes "scss/main.css". You can
	// control this by setting this, e.g. "styles/main.css" will create
	// a Resource with that as a base for RelPermalink etc.
	TargetPath string

	// Hugo automatically adds the entry directories (where the main.scss lives)
	// for project and themes to the list of include paths sent to LibSASS.
	// Any paths set in this setting will be appended. Note that these will be
	// treated as relative to the working dir, i.e. no include paths outside the
	// project/themes.
	IncludePaths []string

	// Default is nested.
	// One of nested, expanded, compact, compressed.
	OutputStyle string

	// When enabled, Hugo will generate a source map.
	EnableSourceMap bool

	// If enabled, sources will be embedded in the generated source map.
	SourceMapIncludeSources bool
}

func decodeOptions(m map[string]any) (opts Options, err error) {
	if m == nil {
		return
	}
	err = mapstructure.WeakDecode(m, &opts)

	if opts.TargetPath != "" {
		opts.TargetPath = helpers.ToSlashTrimLeading(opts.TargetPath)
	}

	return
}
