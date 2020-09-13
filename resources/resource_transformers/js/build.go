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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
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

type Options struct {
	// If not set, the source path will be used as the base target path.
	// Note that the target path's extension may change if the target MIME type
	// is different, e.g. when the source is TypeScript.
	TargetPath string

	// Whether to minify to output.
	Minify bool

	// Whether to write mapfiles
	SourceMap string

	// The language target.
	// One of: es2015, es2016, es2017, es2018, es2019, es2020 or esnext.
	// Default is esnext.
	Target string

	// The output format.
	// One of: iife, cjs, esm
	// Default is to esm.
	Format string

	// External dependencies, e.g. "react".
	Externals []string `hash:"set"`

	// User defined symbols.
	Defines map[string]interface{}

	// What to use instead of React.createElement.
	JSXFactory string

	// What to use instead of React.Fragment.
	JSXFragment string

	mediaType  media.Type
	outDir     string
	contents   string
	sourcefile string
	resolveDir string
	workDir    string
	tsConfig   string
}

func decodeOptions(m map[string]interface{}) (Options, error) {
	var opts Options

	if err := mapstructure.WeakDecode(m, &opts); err != nil {
		return opts, err
	}

	if opts.TargetPath != "" {
		opts.TargetPath = helpers.ToSlashTrimLeading(opts.TargetPath)
	}

	opts.Target = strings.ToLower(opts.Target)
	opts.Format = strings.ToLower(opts.Format)

	return opts, nil
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

	src, err := ioutil.ReadAll(ctx.From)
	if err != nil {
		return err
	}

	sdir, sfile := path.Split(ctx.SourcePath)
	opts.workDir = t.rs.WorkingDir
	opts.sourcefile = sfile
	opts.resolveDir = t.sfs.RealFilename(sdir)
	opts.contents = string(src)
	opts.mediaType = ctx.InMediaType

	rootPaths := make([]string, 0)
	// This is not ideal, because it will ignore possible things inside of package.json in each root.
	// But ESBuild prefers jsconfig/tsconfig anyways.
	for _, mod := range t.rs.PathSpec.Paths.AllModules {
		dir := mod.Dir()
		nodeModules := path.Join(dir, "node_modules")
		if _, err := os.Stat(nodeModules); err == nil {
			rootPaths = append(rootPaths, nodeModules+"/*")
		}
	}

	// Search for original ts/jsconfig file
	tsConfig := path.Join(sdir, "tsconfig.json")
	_, err = t.sfs.Fs.Stat(tsConfig)
	if err != nil {
		tsConfig = path.Join(sdir, "jsconfig.json")
		_, err = t.sfs.Fs.Stat(tsConfig)
		if err != nil {
			tsConfig = path.Join(opts.workDir, "tsconfig.json")
			_, err = t.sfs.Fs.Stat(tsConfig)
			if err != nil {
				tsConfig = path.Join(opts.workDir, "jsconfig.json")
				_, err = t.sfs.Fs.Stat(tsConfig)
				if err != nil {
					// Use this one by default
					tsConfig = path.Join(opts.workDir, "tsconfig.json")
				}
			}
		}
	}

	// Get real source path
	configDir, _ := path.Split(t.sfs.RealFilename(ctx.SourcePath))

	// Resolve paths for @assets and @js (@js is just an alias for assets/js)
	dirs := make([]interface{}, 0)
	jsDirs := make([]interface{}, 0)
	dirIndexs := make([]interface{}, 0)
	jsDirIndexs := make([]interface{}, 0)
	for _, dir := range t.sfs.RealDirs(".") {
		rel, _ := filepath.Rel(configDir, dir)
		dirs = append(dirs, "./"+rel+"/*")
		jsDirs = append(jsDirs, "./"+rel+"/js/*")
		dirIndexs = append(dirIndexs, "./"+rel+"/index.js")
		jsDirIndexs = append(jsDirIndexs, "./"+rel+"/js/index.js")
	}

	// Create new temporary tsconfig file
	newTSConfig, err := ioutil.TempFile(configDir, "tsconfig.*.json")
	if err != nil {
		return err
	}

	// Construct new temporary tsconfig file content
	config := make(map[string]interface{})
	if tsConfig != "" {
		oldConfig, err := ioutil.ReadFile(t.sfs.RealFilename(tsConfig))
		if err != nil {
			return err
		}
		err = json.Unmarshal(oldConfig, &config)
		if err != nil {
			return err
		}
	} else {
		config["compilerOptions"] = map[string]interface{}{
			"baseUrl": ".",
		}
	}

	// Assign new global paths to the config file while reading existing ones.
	oldCompilerOptions := config["compilerOptions"].(map[string]interface{})
	oldPaths := oldCompilerOptions["paths"].(map[string]interface{})
	if oldPaths == nil {
		oldPaths = make(map[string]interface{})
		oldCompilerOptions["paths"] = oldPaths
	}
	oldPaths["@assets/*"] = dirs
	oldPaths["@js/*"] = jsDirs
	// Make @js and @assets absolue matches search for index files
	// to get around the problem in ESBuild resolving folders as index files.
	oldPaths["@assets"] = dirIndexs
	oldPaths["@js"] = jsDirIndexs

	if len(rootPaths) > 0 {
		// This will allow import "react" to resolve a react module that's
		// either in the root node_modules or in one of the hugo mods.
		oldPaths["*"] = rootPaths
	}

	// Output the new config file
	bytes, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	// Write tsconfig file
	_, err = newTSConfig.Write(bytes)
	if err != nil {
		return err
	}
	err = newTSConfig.Close()
	if err != nil {
		return err
	}

	// Tell ESBuild about this new config file to use
	opts.tsConfig = newTSConfig.Name()

	buildOptions, err := toBuildOptions(opts)
	if err != nil {
		os.Remove(opts.tsConfig)
		return err
	}

	result := api.Build(buildOptions)

	os.Remove(opts.tsConfig)

	if len(result.Warnings) > 0 {
		for _, value := range result.Errors {
			t.rs.Logger.WARN.Println(fmt.Sprintf("%s:%d: WARN: %s",
				t.sfs.RealFilename(filepath.Join(sdir, value.Location.File)),
				value.Location.Line, value.Text))
			t.rs.Logger.WARN.Println("  ", value.Location.LineText)
		}
	}
	if len(result.Errors) > 0 {
		output := result.Errors[0].Text
		for _, value := range result.Errors {
			line := fmt.Sprintf("%s:%d ERROR: %s",
				t.sfs.RealFilename(filepath.Join(sdir, value.Location.File)),
				value.Location.Line, value.Text)
			t.rs.Logger.ERROR.Println(line)
			output = fmt.Sprintf("%s\n%s", output, line)
			t.rs.Logger.ERROR.Println("  ", value.Location.LineText)
		}
		return fmt.Errorf("%s", output)
	}
	if buildOptions.Outfile != "" {
		_, tfile := path.Split(opts.TargetPath)
		output := fmt.Sprintf("%s//# sourceMappingURL=%s\n",
			string(result.OutputFiles[1].Contents), tfile+".map")
		_, err := ctx.To.Write([]byte(output))
		if err != nil {
			return err
		}
		ctx.PublishSourceMap(string(result.OutputFiles[0].Contents))
	} else {
		ctx.To.Write(result.OutputFiles[0].Contents)
	}
	return nil
}

func (c *Client) Process(res resources.ResourceTransformer, opts map[string]interface{}) (resource.Resource, error) {
	return res.Transform(
		&buildTransformation{rs: c.rs, sfs: c.sfs, optsm: opts},
	)
}

func toBuildOptions(opts Options) (buildOptions api.BuildOptions, err error) {
	var target api.Target
	switch opts.Target {
	case "", "esnext":
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
		err = fmt.Errorf("invalid target: %q", opts.Target)
		return
	}

	mediaType := opts.mediaType
	if mediaType.IsZero() {
		mediaType = media.JavascriptType
	}

	var loader api.Loader
	switch mediaType.SubType {
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
		err = fmt.Errorf("unsupported Media Type: %q", opts.mediaType)
		return
	}

	var format api.Format
	// One of: iife, cjs, esm
	switch opts.Format {
	case "", "iife":
		format = api.FormatIIFE
	case "esm":
		format = api.FormatESModule
	case "cjs":
		format = api.FormatCommonJS
	default:
		err = fmt.Errorf("unsupported script output format: %q", opts.Format)
		return
	}

	var defines map[string]string
	if opts.Defines != nil {
		defines = cast.ToStringMapString(opts.Defines)
	}

	var outDir = opts.outDir
	var outFile = ""
	var sourceMap api.SourceMap
	switch opts.SourceMap {
	case "inline":
		sourceMap = api.SourceMapInline
	case "external":
		sourceMap = api.SourceMapExternal
		outFile = filepath.Join(opts.workDir, opts.TargetPath)
		outDir = ""
	case "":
		sourceMap = api.SourceMapNone
	default:
		err = fmt.Errorf("unsupported sourcemap type: %q", opts.SourceMap)
		return
	}

	buildOptions = api.BuildOptions{
		Outfile: outFile,
		Bundle:  true,

		Target:    target,
		Format:    format,
		Sourcemap: sourceMap,

		MinifyWhitespace:  opts.Minify,
		MinifyIdentifiers: opts.Minify,
		MinifySyntax:      opts.Minify,

		Outdir:  outDir,
		Defines: defines,

		Externals: opts.Externals,

		JSXFactory:  opts.JSXFactory,
		JSXFragment: opts.JSXFragment,

		Tsconfig: opts.tsConfig,

		Stdin: &api.StdinOptions{
			Contents:   opts.contents,
			Sourcefile: opts.sourcefile,
			ResolveDir: opts.resolveDir,
			Loader:     loader,
		},
	}
	return

}
