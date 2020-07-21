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

package js

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/spf13/cast"

	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugolib/filesystems"
	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/resources/internal"

	"github.com/mitchellh/mapstructure"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/resources/resource"
)

const defaultTarget = "esnext"

type Options struct {
	// If not set, the source path will be used as the base target path.
	// Note that the target path's extension may change if the target MIME type
	// is different, e.g. when the source is TypeScript.
	TargetPath string

	// Whether to minify to output.
	Minify bool

	// The language target.
	// One of: es2015, es2016, es2017, es2018, es2019, es2020 or esnext.
	// Default is esnext.
	Target string

	// External dependencies, e.g. "react".
	Externals []string `hash:"set"`

	// User defined symbols.
	Defines map[string]interface{}

	// What to use instead of React.createElement.
	JSXFactory string

	// What to use instead of React.Fragment.
	JSXFragment string
}

func decodeOptions(m map[string]interface{}) (opts Options, err error) {
	err = mapstructure.WeakDecode(m, &opts)
	if err != nil {
		return
	}

	if opts.TargetPath != "" {
		opts.TargetPath = helpers.ToSlashTrimLeading(opts.TargetPath)
	}

	if opts.Target == "" {
		opts.Target = defaultTarget
	}

	opts.Target = strings.ToLower(opts.Target)

	return
}

type Client struct {
	rs  *resources.Spec
	sfs *filesystems.SourceFilesystem
}

func New(fs *filesystems.SourceFilesystem, rs *resources.Spec) *Client {
	return &Client{rs: rs, sfs: fs}
}

type buildTransformation struct {
	optsm map[string]interface{}
	rs    *resources.Spec
	sfs   *filesystems.SourceFilesystem
}

func (t *buildTransformation) Key() internal.ResourceTransformationKey {
	return internal.NewResourceTransformationKey("jsbuild", t.optsm)
}

func (t *buildTransformation) Transform(ctx *resources.ResourceTransformationCtx) error {
	ctx.OutMediaType = media.JavascriptType

	opts, err := decodeOptions(t.optsm)
	if err != nil {
		return err
	}

	if opts.TargetPath != "" {
		ctx.OutPath = opts.TargetPath
	} else {
		ctx.ReplaceOutPathExtension(".js")
	}

	var target api.Target
	switch opts.Target {
	case defaultTarget:
		target = api.ESNext
	case "es5":
		target = api.ES5
	case "es6", "es2015":
		target = api.ES2015
	case "es2016":
		target = api.ES2016
	case "es2017":
		target = api.ES2017
	case "es2018":
		target = api.ES2018
	case "es2019":
		target = api.ES2019
	case "es2020":
		target = api.ES2020
	default:
		return fmt.Errorf("invalid target: %q", opts.Target)
	}

	var loader api.Loader
	switch ctx.InMediaType.SubType {
	// TODO(bep) ESBuild support a set of other loaders, but I currently fail
	// to see the relevance. That may change as we start using this.
	case media.JavascriptType.SubType:
		loader = api.LoaderJS
	case media.TypeScriptType.SubType:
		loader = api.LoaderTS
	case media.TSXType.SubType:
		loader = api.LoaderTSX
	case media.JSXType.SubType:
		loader = api.LoaderJSX
	default:
		return fmt.Errorf("unsupported Media Type: %q", ctx.InMediaType)

	}

	src, err := ioutil.ReadAll(ctx.From)
	if err != nil {
		return err
	}

	sdir, sfile := path.Split(ctx.SourcePath)
	sdir = t.sfs.RealFilename(sdir)

	buildOptions := api.BuildOptions{
		Outfile: "",
		Bundle:  true,

		Target: target,

		MinifyWhitespace:  opts.Minify,
		MinifyIdentifiers: opts.Minify,
		MinifySyntax:      opts.Minify,

		Defines: cast.ToStringMapString(opts.Defines),

		Externals: opts.Externals,

		JSXFactory:  opts.JSXFactory,
		JSXFragment: opts.JSXFragment,

		//Tsconfig: opts.TSConfig,

		Stdin: &api.StdinOptions{
			Contents:   string(src),
			Sourcefile: sfile,
			ResolveDir: sdir,
			Loader:     loader,
		},
	}
	result := api.Build(buildOptions)
	if len(result.Errors) > 0 {
		return fmt.Errorf("%s", result.Errors[0].Text)
	}
	if len(result.OutputFiles) != 1 {
		return fmt.Errorf("unexpected output count: %d", len(result.OutputFiles))
	}

	ctx.To.Write(result.OutputFiles[0].Contents)
	return nil
}

func (c *Client) Process(res resources.ResourceTransformer, opts map[string]interface{}) (resource.Resource, error) {
	return res.Transform(
		&buildTransformation{rs: c.rs, sfs: c.sfs, optsm: opts},
	)
}
