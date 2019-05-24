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

package scss

import (
	"github.com/bep/go-tocss/scss"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugolib/filesystems"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/resources/resource"

	"github.com/mitchellh/mapstructure"
)

type Client struct {
	rs     *resources.Spec
	sfs    *filesystems.SourceFilesystem
	workFs *filesystems.SourceFilesystem
}

func New(fs *filesystems.SourceFilesystem, rs *resources.Spec) (*Client, error) {
	return &Client{sfs: fs, workFs: rs.BaseFs.Work, rs: rs}, nil
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

	// Precision of floating point math.
	Precision int

	// When enabled, Hugo will generate a source map.
	EnableSourceMap bool
}

type options struct {
	// The options we receive from the end user.
	from Options

	// The options we send to the SCSS library.
	to scss.Options
}

func (c *Client) ToCSS(res resource.Resource, opts Options) (resource.Resource, error) {
	internalOptions := options{
		from: opts,
	}

	// Transfer values from client.
	internalOptions.to.Precision = opts.Precision
	internalOptions.to.OutputStyle = scss.OutputStyleFromString(opts.OutputStyle)

	if internalOptions.to.Precision == 0 {
		// bootstrap-sass requires 8 digits precision. The libsass default is 5.
		// https://github.com/twbs/bootstrap-sass/blob/master/README.md#sass-number-precision
		internalOptions.to.Precision = 8
	}

	return c.rs.Transform(
		res,
		&toCSSTransformation{c: c, options: internalOptions},
	)
}

type toCSSTransformation struct {
	c       *Client
	options options
}

func (t *toCSSTransformation) Key() resources.ResourceTransformationKey {
	return resources.NewResourceTransformationKey("tocss", t.options.from)
}

func DecodeOptions(m map[string]interface{}) (opts Options, err error) {
	if m == nil {
		return
	}
	err = mapstructure.WeakDecode(m, &opts)

	if opts.TargetPath != "" {
		opts.TargetPath = helpers.ToSlashTrimLeading(opts.TargetPath)
	}

	return
}
