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
	"regexp"

	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugolib/filesystems"
	"github.com/gohugoio/hugo/resources"
	"github.com/spf13/afero"

	"github.com/mitchellh/mapstructure"
)

const transformationName = "tocss"

type Client struct {
	rs     *resources.Spec
	sfs    *filesystems.SourceFilesystem
	workFs afero.Fs
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

var (
	regularCSSImportTo   = regexp.MustCompile(`.*(@import "(.*\.css)";).*`)
	regularCSSImportFrom = regexp.MustCompile(`.*(\/\* HUGO_IMPORT_START (.*) HUGO_IMPORT_END \*\/).*`)
)

func replaceRegularImportsIn(s string) (string, bool) {
	replaced := regularCSSImportTo.ReplaceAllString(s, "/* HUGO_IMPORT_START $2 HUGO_IMPORT_END */")
	return replaced, s != replaced
}

func replaceRegularImportsOut(s string) string {
	return regularCSSImportFrom.ReplaceAllString(s, "@import \"$2\";")
}
