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

// Package dartsass integrates with the Dart Sass Embedded protocol to transpile
// SCSS/SASS.
package dartsass

import (
	"fmt"
	"io"
	"strings"

	"github.com/bep/godartsass/v2"
	"github.com/bep/logg"
	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/hugolib/filesystems"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/resources/resource"
	"github.com/spf13/afero"

	"github.com/mitchellh/mapstructure"
)

// used as part of the cache key.
const transformationName = "tocss-dart"

// See https://github.com/sass/dart-sass-embedded/issues/24
// Note: This prefix must be all lower case.
const dartSassStdinPrefix = "hugostdin:"

func New(fs *filesystems.SourceFilesystem, rs *resources.Spec) (*Client, error) {
	if !Supports() {
		return &Client{}, nil
	}

	if hugo.DartSassBinaryName == "" {
		return nil, fmt.Errorf("no Dart Sass binary found in $PATH")
	}

	if !hugo.IsDartSassGeV2() {
		return nil, fmt.Errorf("unsupported Dart Sass version detected, please upgrade to Dart Sass 1.63.0 or later, see https://gohugo.io/functions/css/sass/#dart-sass")
	}

	if err := rs.ExecHelper.Sec().CheckAllowedExec(hugo.DartSassBinaryName); err != nil {
		return nil, err
	}

	var (
		transpiler *godartsass.Transpiler
		err        error
		infol      = rs.Logger.InfoCommand("Dart Sass")
		warnl      = rs.Logger.WarnCommand("Dart Sass")
	)

	transpiler, err = godartsass.Start(godartsass.Options{
		DartSassEmbeddedFilename: hugo.DartSassBinaryName,
		LogEventHandler: func(event godartsass.LogEvent) {
			message := strings.ReplaceAll(event.Message, dartSassStdinPrefix, "")
			switch event.Type {
			case godartsass.LogEventTypeDebug:
				// Log as Info for now, we may adjust this if it gets too chatty.
				infol.Log(logg.String(message))
			case godartsass.LogEventTypeDeprecated:
				warnl.Logf("DEPRECATED [%s]: %s", event.DeprecationType, message)
			default:
				// The rest are @warn statements.
				warnl.Log(logg.String(message))
			}
		},
	})
	if err != nil {
		return nil, err
	}
	return &Client{sfs: fs, workFs: rs.BaseFs.Work, rs: rs, transpiler: transpiler}, nil
}

type Client struct {
	rs     *resources.Spec
	sfs    *filesystems.SourceFilesystem
	workFs afero.Fs

	// This may be nil if Dart Sass is not available.
	transpiler *godartsass.Transpiler
}

func (c *Client) ToCSS(res resources.ResourceTransformer, args map[string]any) (resource.Resource, error) {
	if c.transpiler == nil {
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
	in := helpers.ReaderToString(src)

	args.Source = in

	res, err := c.transpiler.Execute(args)
	if err != nil {
		if err.Error() == "unexpected EOF" {
			//lint:ignore ST1005 end user message.
			return res, fmt.Errorf("got unexpected EOF when executing %q. The user running hugo must have read and execute permissions on this program. With execute permissions only, this error is thrown.", hugo.DartSassBinaryName)
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

	// Vars will be available in 'hugo:vars', e.g:
	//     @use "hugo:vars";
	//     $color: vars.$color;
	Vars map[string]any

	// Deprecations IDs in this slice will be silenced.
	// The IDs can be found in the Dart Sass log output, e.g. "import" in
	//    WARN  Dart Sass: DEPRECATED [import].
	SilenceDeprecations []string

	// Whether to silence deprecation warnings from dependencies, where a
	// dependency is considered any file transitively imported through a load
	// path. This does not apply to @warn or @debug rules.
	SilenceDependencyDeprecations bool
}

func decodeOptions(m map[string]any) (opts Options, err error) {
	if m == nil {
		return
	}
	err = mapstructure.WeakDecode(m, &opts)

	if opts.TargetPath != "" {
		opts.TargetPath = paths.ToSlashTrimLeading(opts.TargetPath)
	}

	return
}
