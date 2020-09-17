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
	"reflect"
	"strings"

	"github.com/achiku/varfmt"
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

// Options esbuild configuration
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

	// User defined data (must be JSON marshall'able)
	Data interface{}

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

// Client context for esbuild
type Client struct {
	rs  *resources.Spec
	sfs *filesystems.SourceFilesystem
}

// New create new client context
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

func appendExts(list []string, rel string) []string {
	for _, ext := range []string{".tsx", ".ts", ".jsx", ".mjs", ".cjs", ".js", ".json"} {
		list = append(list, fmt.Sprintf("%s/index%s", rel, ext))
	}
	return list
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

	sdir, sfile := filepath.Split(t.sfs.RealFilename(ctx.SourcePath))
	opts.workDir, err = filepath.Abs(t.rs.WorkingDir)
	if err != nil {
		return err
	}

	opts.sourcefile = sfile
	opts.resolveDir = sdir
	opts.contents = string(src)
	opts.mediaType = ctx.InMediaType

	// Create new temporary tsconfig file
	newTSConfig, err := ioutil.TempFile("", "tsconfig.*.json")
	if err != nil {
		return err
	}

	filesToDelete := make([]*os.File, 0)

	defer func() {
		for _, file := range filesToDelete {
			os.Remove(file.Name())
		}
	}()

	filesToDelete = append(filesToDelete, newTSConfig)
	configDir, _ := filepath.Split(newTSConfig.Name())

	// Search for the innerMost tsconfig or jsconfig
	innerTsConfig := ""
	tsDir := opts.resolveDir
	baseURLAbs := configDir
	baseURL := "."
	for tsDir != "." {
		tryTsConfig := path.Join(tsDir, "tsconfig.json")
		_, err := os.Stat(tryTsConfig)
		if err != nil {
			tryTsConfig := path.Join(tsDir, "jsconfig.json")
			_, err = os.Stat(tryTsConfig)
			if err == nil {
				innerTsConfig = tryTsConfig
				baseURLAbs = tsDir
				break
			}
		} else {
			innerTsConfig = tryTsConfig
			baseURLAbs = tsDir
			break
		}
		if tsDir == opts.workDir {
			break
		}
		tsDir = path.Dir(tsDir)
	}

	// Resolve paths for @assets and @js (@js is just an alias for assets/js)
	dirs := make([]string, 0)
	rootPaths := make([]string, 0)
	for _, dir := range t.sfs.RealDirs(".") {
		rootDir := dir
		if !strings.HasSuffix(dir, "package.json") {
			dirs = append(dirs, dir)
		} else {
			rootDir, _ = path.Split(dir)
		}
		nodeModules := path.Join(rootDir, "node_modules")
		if _, err := os.Stat(nodeModules); err == nil {
			rootPaths = append(rootPaths, nodeModules)
		}
	}

	// Construct new temporary tsconfig file content
	config := make(map[string]interface{})
	if innerTsConfig != "" {
		oldConfig, err := ioutil.ReadFile(innerTsConfig)
		if err == nil {
			// If there is an error, it just means there is no config file here.
			// Since we're also using the tsConfig file path to detect where
			// to put the temp file, this is ok.
			err = json.Unmarshal(oldConfig, &config)
			if err != nil {
				return err
			}
		}
	}

	if config["compilerOptions"] == nil {
		config["compilerOptions"] = map[string]interface{}{}
	}

	// Assign new global paths to the config file while reading existing ones.
	compilerOptions := config["compilerOptions"].(map[string]interface{})

	// Handle original baseUrl if it's there
	if compilerOptions["baseUrl"] != nil {
		baseURL = compilerOptions["baseUrl"].(string)
		oldBaseURLAbs := path.Join(tsDir, baseURL)
		rel, _ := filepath.Rel(configDir, oldBaseURLAbs)
		configDir = oldBaseURLAbs
		baseURLAbs = configDir
		if "/" != helpers.FilePathSeparator {
			// On windows we need to use slashes instead of backslash
			rel = strings.ReplaceAll(rel, helpers.FilePathSeparator, "/")
		}
		if rel != "" {
			if strings.HasPrefix(rel, ".") {
				baseURL = rel
			} else {
				baseURL = fmt.Sprintf("./%s", rel)
			}
		}
		compilerOptions["baseUrl"] = baseURL
	} else {
		compilerOptions["baseUrl"] = baseURL
	}

	jsRel := func(refPath string) string {
		rel, _ := filepath.Rel(configDir, refPath)
		if "/" != helpers.FilePathSeparator {
			// On windows we need to use slashes instead of backslash
			rel = strings.ReplaceAll(rel, helpers.FilePathSeparator, "/")
		}
		if rel != "" {
			if !strings.HasPrefix(rel, ".") {
				rel = fmt.Sprintf("./%s", rel)
			}
		} else {
			rel = "."
		}
		return rel
	}

	// Handle possible extends
	if config["extends"] != nil {
		extends := config["extends"].(string)
		extendsAbs := path.Join(tsDir, extends)
		rel := jsRel(extendsAbs)
		config["extends"] = rel
	}

	var optionsPaths map[string]interface{}
	// Get original paths if they exist
	if compilerOptions["paths"] != nil {
		optionsPaths = compilerOptions["paths"].(map[string]interface{})
	} else {
		optionsPaths = make(map[string]interface{})
	}
	compilerOptions["paths"] = optionsPaths

	assets := make([]string, 0)
	assetsExact := make([]string, 0)
	js := make([]string, 0)
	jsExact := make([]string, 0)
	for _, dir := range dirs {
		rel := jsRel(dir)
		assets = append(assets, fmt.Sprintf("%s/*", rel))
		assetsExact = appendExts(assetsExact, rel)

		rel = jsRel(filepath.Join(dir, "js"))
		js = append(js, fmt.Sprintf("%s/*", rel))
		jsExact = appendExts(jsExact, rel)
	}

	optionsPaths["@assets/*"] = assets
	optionsPaths["@js/*"] = js

	// Make @js and @assets absolue matches search for index files
	// to get around the problem in ESBuild resolving folders as index files.
	optionsPaths["@assets"] = assetsExact
	optionsPaths["@js"] = jsExact

	var newDataFile *os.File
	if opts.Data != nil {
		// Create a data file
		lines := make([]string, 0)
		lines = append(lines, "// auto generated data import")
		exports := make([]string, 0)
		keys := make(map[string]bool)

		var bytes []byte

		conv := reflect.ValueOf(opts.Data)
		convType := conv.Kind()
		if convType == reflect.Interface {
			if conv.IsNil() {
				conv = reflect.Value{}
			}
		}

		if conv.Kind() != reflect.Map {
			// Write out as single JSON file
			newDataFile, err = ioutil.TempFile("", "data.*.json")
			// Output the data
			bytes, err = json.MarshalIndent(conv.InterfaceData(), "", "  ")
			if err != nil {
				return err
			}
		} else {
			// Try to allow tree shaking at the root
			newDataFile, err = ioutil.TempFile(configDir, "data.*.js")
			for _, key := range conv.MapKeys() {
				strKey := key.Interface().(string)
				if keys[strKey] {
					continue
				}
				keys[strKey] = true

				value := conv.MapIndex(key)

				keyVar := varfmt.PublicVarName(strKey)

				// Output the data
				bytes, err := json.MarshalIndent(value.Interface(), "", "  ")
				if err != nil {
					return err
				}
				jsonValue := string(bytes)

				lines = append(lines, fmt.Sprintf("export const %s = %s;", keyVar, jsonValue))
				exports = append(exports, fmt.Sprintf("  %s,", keyVar))
				if strKey != keyVar {
					exports = append(exports, fmt.Sprintf("  [\"%s\"]: %s,", strKey, keyVar))
				}
			}

			lines = append(lines, "const all = {")
			for _, line := range exports {
				lines = append(lines, line)
			}
			lines = append(lines, "};")
			lines = append(lines, "export default all;")

			bytes = []byte(strings.Join(lines, "\n"))
		}

		// Write tsconfig file
		_, err = newDataFile.Write(bytes)
		if err != nil {
			return err
		}
		err = newDataFile.Close()
		if err != nil {
			return err
		}

		// Link this file into `import data from "@data"`
		dataFiles := make([]string, 1)
		rel, _ := filepath.Rel(baseURLAbs, newDataFile.Name())
		dataFiles[0] = rel
		optionsPaths["@data"] = dataFiles

		filesToDelete = append(filesToDelete, newDataFile)
	}

	if len(rootPaths) > 0 {
		// This will allow import "react" to resolve a react module that's
		// either in the root node_modules or in one of the hugo mods.
		optionsPaths["*"] = rootPaths
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

	if len(result.Warnings) > 0 {
		for _, value := range result.Warnings {
			if value.Location != nil {
				t.rs.Logger.WARN.Println(fmt.Sprintf("%s:%d: WARN: %s",
					filepath.Join(sdir, value.Location.File),
					value.Location.Line, value.Text))
				t.rs.Logger.WARN.Println("  ", value.Location.LineText)
			} else {
				t.rs.Logger.WARN.Println(fmt.Sprintf("%s: WARN: %s",
					sdir,
					value.Text))
			}
		}
	}
	if len(result.Errors) > 0 {
		output := result.Errors[0].Text
		for _, value := range result.Errors {
			var line string
			if value.Location != nil {
				line = fmt.Sprintf("%s:%d ERROR: %s",
					filepath.Join(sdir, value.Location.File),
					value.Location.Line, value.Text)
			} else {
				line = fmt.Sprintf("%s ERROR: %s",
					sdir,
					value.Text)
			}
			t.rs.Logger.ERROR.Println(line)
			output = fmt.Sprintf("%s\n%s", output, line)
			if value.Location != nil {
				t.rs.Logger.ERROR.Println("  ", value.Location.LineText)
			}
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

// Process process esbuild transform
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

	// By default we only need to specify outDir and no outFile
	var outDir = opts.outDir
	var outFile = ""
	var sourceMap api.SourceMap
	switch opts.SourceMap {
	case "inline":
		sourceMap = api.SourceMapInline
	case "external":
		// When doing external sourcemaps we should specify
		// out file and no out dir
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
