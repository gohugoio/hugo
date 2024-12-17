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

package esbuild

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/identity"

	"github.com/evanw/esbuild/pkg/api"

	"github.com/gohugoio/hugo/media"
	"github.com/mitchellh/mapstructure"
)

var (
	nameTarget = map[string]api.Target{
		"":       api.ESNext,
		"esnext": api.ESNext,
		"es5":    api.ES5,
		"es6":    api.ES2015,
		"es2015": api.ES2015,
		"es2016": api.ES2016,
		"es2017": api.ES2017,
		"es2018": api.ES2018,
		"es2019": api.ES2019,
		"es2020": api.ES2020,
		"es2021": api.ES2021,
		"es2022": api.ES2022,
		"es2023": api.ES2023,
		"es2024": api.ES2024,
	}

	// source names: https://github.com/evanw/esbuild/blob/9eca46464ed5615cb36a3beb3f7a7b9a8ffbe7cf/internal/config/config.go#L208
	nameLoader = map[string]api.Loader{
		"none":       api.LoaderNone,
		"base64":     api.LoaderBase64,
		"binary":     api.LoaderBinary,
		"copy":       api.LoaderFile,
		"css":        api.LoaderCSS,
		"dataurl":    api.LoaderDataURL,
		"default":    api.LoaderDefault,
		"empty":      api.LoaderEmpty,
		"file":       api.LoaderFile,
		"global-css": api.LoaderGlobalCSS,
		"js":         api.LoaderJS,
		"json":       api.LoaderJSON,
		"jsx":        api.LoaderJSX,
		"local-css":  api.LoaderLocalCSS,
		"text":       api.LoaderText,
		"ts":         api.LoaderTS,
		"tsx":        api.LoaderTSX,
	}
)

// DecodeExternalOptions decodes the given map into ExternalOptions.
func DecodeExternalOptions(m map[string]any) (ExternalOptions, error) {
	opts := ExternalOptions{
		SourcesContent: true,
	}

	if err := mapstructure.WeakDecode(m, &opts); err != nil {
		return opts, err
	}

	if opts.TargetPath != "" {
		opts.TargetPath = paths.ToSlashTrimLeading(opts.TargetPath)
	}

	opts.Target = strings.ToLower(opts.Target)
	opts.Format = strings.ToLower(opts.Format)

	return opts, nil
}

// ErrorMessageResolved holds a resolved error message.
type ErrorMessageResolved struct {
	Path    string
	Message string
	Content hugio.ReadSeekCloser
}

// ExternalOptions holds user facing options for the js.Build template function.
type ExternalOptions struct {
	// If not set, the source path will be used as the base target path.
	// Note that the target path's extension may change if the target MIME type
	// is different, e.g. when the source is TypeScript.
	TargetPath string

	// Whether to minify to output.
	Minify bool

	// One of "inline", "external", "linked" or "none".
	SourceMap string

	SourcesContent bool

	// The language target.
	// One of: es2015, es2016, es2017, es2018, es2019, es2020 or esnext.
	// Default is esnext.
	Target string

	// The output format.
	// One of: iife, cjs, esm
	// Default is to esm.
	Format string

	// One of browser, node, neutral.
	// Default is browser.
	// See https://esbuild.github.io/api/#platform
	Platform string

	// External dependencies, e.g. "react".
	Externals []string

	// This option allows you to automatically replace a global variable with an import from another file.
	// The filenames must be relative to /assets.
	// See https://esbuild.github.io/api/#inject
	Inject []string

	// User defined symbols.
	Defines map[string]any

	// Maps a component import to another.
	Shims map[string]string

	// Configuring a loader for a given file type lets you load that file type with an
	// import statement or a require call. For example, configuring the .png file extension
	// to use the data URL loader means importing a .png file gives you a data URL
	// containing the contents of that image
	//
	// See https://esbuild.github.io/api/#loader
	Loaders map[string]string

	// User defined params. Will be marshaled to JSON and available as "@params", e.g.
	//     import * as params from '@params';
	Params any

	// What to use instead of React.createElement.
	JSXFactory string

	// What to use instead of React.Fragment.
	JSXFragment string

	// What to do about JSX syntax.
	// See https://esbuild.github.io/api/#jsx
	JSX string

	// Which library to use to automatically import JSX helper functions from. Only works if JSX is set to automatic.
	// See https://esbuild.github.io/api/#jsx-import-source
	JSXImportSource string

	// There is/was a bug in WebKit with severe performance issue with the tracking
	// of TDZ checks in JavaScriptCore.
	//
	// Enabling this flag removes the TDZ and `const` assignment checks and
	// may improve performance of larger JS codebases until the WebKit fix
	// is in widespread use.
	//
	// See https://bugs.webkit.org/show_bug.cgi?id=199866
	// Deprecated: This no longer have any effect and will be removed.
	// TODO(bep) remove. See https://github.com/evanw/esbuild/commit/869e8117b499ca1dbfc5b3021938a53ffe934dba
	AvoidTDZ bool
}

// InternalOptions holds internal options for the js.Build template function.
type InternalOptions struct {
	MediaType     media.Type
	OutDir        string
	Contents      string
	SourceDir     string
	ResolveDir    string
	AbsWorkingDir string
	Metafile      bool

	StdinSourcePath string

	DependencyManager identity.Manager

	Stdin                   bool // Set to true to pass in the entry point as a byte slice.
	Splitting               bool
	TsConfig                string
	EntryPoints             []string
	ImportOnResolveFunc     func(string, api.OnResolveArgs) string
	ImportOnLoadFunc        func(api.OnLoadArgs) string
	ImportParamsOnLoadFunc  func(args api.OnLoadArgs) json.RawMessage
	ErrorMessageResolveFunc func(api.Message) *ErrorMessageResolved
	ResolveSourceMapSource  func(string) string // Used to resolve paths in error source maps.
}

// Options holds the options passed to Build.
type Options struct {
	ExternalOptions
	InternalOptions

	compiled api.BuildOptions
}

func (opts *Options) compile() (err error) {
	target, found := nameTarget[opts.Target]
	if !found {
		err = fmt.Errorf("invalid target: %q", opts.Target)
		return
	}

	var loaders map[string]api.Loader
	if opts.Loaders != nil {
		loaders = make(map[string]api.Loader)
		for k, v := range opts.Loaders {
			loader, found := nameLoader[v]
			if !found {
				err = fmt.Errorf("invalid loader: %q", v)
				return
			}
			loaders[k] = loader
		}
	}

	mediaType := opts.MediaType
	if mediaType.IsZero() {
		mediaType = media.Builtin.JavascriptType
	}

	var loader api.Loader
	switch mediaType.SubType {
	case media.Builtin.JavascriptType.SubType:
		loader = api.LoaderJS
	case media.Builtin.TypeScriptType.SubType:
		loader = api.LoaderTS
	case media.Builtin.TSXType.SubType:
		loader = api.LoaderTSX
	case media.Builtin.JSXType.SubType:
		loader = api.LoaderJSX
	default:
		err = fmt.Errorf("unsupported Media Type: %q", opts.MediaType)
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

	var jsx api.JSX
	switch opts.JSX {
	case "", "transform":
		jsx = api.JSXTransform
	case "preserve":
		jsx = api.JSXPreserve
	case "automatic":
		jsx = api.JSXAutomatic
	default:
		err = fmt.Errorf("unsupported jsx type: %q", opts.JSX)
		return
	}

	var platform api.Platform
	switch opts.Platform {
	case "", "browser":
		platform = api.PlatformBrowser
	case "node":
		platform = api.PlatformNode
	case "neutral":
		platform = api.PlatformNeutral
	default:
		err = fmt.Errorf("unsupported platform type: %q", opts.Platform)
		return
	}

	var defines map[string]string
	if opts.Defines != nil {
		defines = maps.ToStringMapString(opts.Defines)
	}

	// By default we only need to specify outDir and no outFile
	outDir := opts.OutDir
	outFile := ""
	var sourceMap api.SourceMap
	switch opts.SourceMap {
	case "inline":
		sourceMap = api.SourceMapInline
	case "external":
		sourceMap = api.SourceMapExternal
	case "linked":
		sourceMap = api.SourceMapLinked
	case "", "none":
		sourceMap = api.SourceMapNone
	default:
		err = fmt.Errorf("unsupported sourcemap type: %q", opts.SourceMap)
		return
	}

	sourcesContent := api.SourcesContentInclude
	if !opts.SourcesContent {
		sourcesContent = api.SourcesContentExclude
	}

	opts.compiled = api.BuildOptions{
		Outfile:       outFile,
		Bundle:        true,
		Metafile:      opts.Metafile,
		AbsWorkingDir: opts.AbsWorkingDir,

		Target:         target,
		Format:         format,
		Platform:       platform,
		Sourcemap:      sourceMap,
		SourcesContent: sourcesContent,

		Loader: loaders,

		MinifyWhitespace:  opts.Minify,
		MinifyIdentifiers: opts.Minify,
		MinifySyntax:      opts.Minify,

		Outdir:    outDir,
		Splitting: opts.Splitting,

		Define:   defines,
		External: opts.Externals,

		JSXFactory:  opts.JSXFactory,
		JSXFragment: opts.JSXFragment,

		JSX:             jsx,
		JSXImportSource: opts.JSXImportSource,

		Tsconfig: opts.TsConfig,

		EntryPoints: opts.EntryPoints,
	}

	if opts.Stdin {
		// This makes ESBuild pass `stdin` as the Importer to the import.
		opts.compiled.Stdin = &api.StdinOptions{
			Contents:   opts.Contents,
			ResolveDir: opts.ResolveDir,
			Loader:     loader,
		}
	}
	return
}

func (o Options) loaderFromFilename(filename string) api.Loader {
	ext := filepath.Ext(filename)
	if optsLoaders := o.compiled.Loader; optsLoaders != nil {
		if l, found := optsLoaders[ext]; found {
			return l
		}
	}
	l, found := extensionToLoaderMap[ext]
	if found {
		return l
	}
	return api.LoaderJS
}

func (opts *Options) validate() error {
	if opts.ImportOnResolveFunc != nil && opts.ImportOnLoadFunc == nil {
		return fmt.Errorf("ImportOnLoadFunc must be set if ImportOnResolveFunc is set")
	}
	if opts.ImportOnResolveFunc == nil && opts.ImportOnLoadFunc != nil {
		return fmt.Errorf("ImportOnResolveFunc must be set if ImportOnLoadFunc is set")
	}
	if opts.AbsWorkingDir == "" {
		return fmt.Errorf("AbsWorkingDir must be set")
	}
	return nil
}
