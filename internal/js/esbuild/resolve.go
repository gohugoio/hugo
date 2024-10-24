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
	"os"
	"path/filepath"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/resources"
	"github.com/spf13/afero"
)

const (
	NsHugoImport            = "ns-hugo-import"
	NsHugoImportResolveFunc = "ns-hugo-import-resolvefunc"
	nsHugoParams            = "ns-hugo-params"

	stdinImporter = "<stdin>"
)

const (
	PrefixHugoVirtual = "@hugo-virtual"
)

var extensionToLoaderMap = map[string]api.Loader{
	".js":   api.LoaderJS,
	".mjs":  api.LoaderJS,
	".cjs":  api.LoaderJS,
	".jsx":  api.LoaderJSX,
	".ts":   api.LoaderTS,
	".tsx":  api.LoaderTSX,
	".css":  api.LoaderCSS,
	".json": api.LoaderJSON,
	".txt":  api.LoaderText,
}

func loaderFromFilename(filename string) api.Loader {
	l, found := extensionToLoaderMap[filepath.Ext(filename)]
	if found {
		return l
	}
	return api.LoaderJS
}

func resolveComponentInAssets(fs afero.Fs, impPath string) *hugofs.FileMeta {
	findFirst := func(base string) *hugofs.FileMeta {
		// This is the most common sub-set of ESBuild's default extensions.
		// We assume that imports of JSON, CSS etc. will be using their full
		// name with extension.
		for _, ext := range []string{".js", ".ts", ".tsx", ".jsx"} {
			if strings.HasSuffix(impPath, ext) {
				// Import of foo.js.js need the full name.
				continue
			}
			if fi, err := fs.Stat(base + ext); err == nil {
				return fi.(hugofs.FileMetaInfo).Meta()
			}
		}

		// Not found.
		return nil
	}

	var m *hugofs.FileMeta

	// We need to check if this is a regular file imported without an extension.
	// There may be ambiguous situations where both foo.js and foo/index.js exists.
	// This import order is in line with both how Node and ESBuild's native
	// import resolver works.

	// It may be a regular file imported without an extension, e.g.
	// foo or foo/index.
	m = findFirst(impPath)
	if m != nil {
		return m
	}

	base := filepath.Base(impPath)
	if base == "index" {
		// try index.esm.js etc.
		m = findFirst(impPath + ".esm")
		if m != nil {
			return m
		}
	}

	// Check the path as is.
	fi, err := fs.Stat(impPath)

	if err == nil {
		if fi.IsDir() {
			m = findFirst(filepath.Join(impPath, "index"))
			if m == nil {
				m = findFirst(filepath.Join(impPath, "index.esm"))
			}
		} else {
			m = fi.(hugofs.FileMetaInfo).Meta()
		}
	} else if strings.HasSuffix(base, ".js") {
		m = findFirst(strings.TrimSuffix(impPath, ".js"))
	}

	return m
}

func createBuildPlugins(rs *resources.Spec, depsManager identity.Manager, opts Options) ([]api.Plugin, error) {
	fs := rs.Assets

	resolveImport := func(args api.OnResolveArgs) (api.OnResolveResult, error) {
		impPath := args.Path
		shimmed := false
		if opts.Shims != nil {
			override, found := opts.Shims[impPath]
			if found {
				impPath = override
				shimmed = true
			}
		}

		if opts.ImportOnResolveFunc != nil {
			if s := opts.ImportOnResolveFunc(depsManager, impPath, args); s != "" {
				return api.OnResolveResult{Path: s, Namespace: NsHugoImportResolveFunc}, nil
			}
		}

		importer := args.Importer

		isStdin := importer == stdinImporter
		var relDir string
		if !isStdin {
			if strings.HasPrefix(importer, PrefixHugoVirtual) {
				relDir = filepath.Dir(strings.TrimPrefix(importer, PrefixHugoVirtual))
			} else {
				rel, found := fs.MakePathRelative(importer, true)

				if !found {
					if shimmed {
						relDir = opts.SourceDir
					} else {
						// Not in any of the /assets folders.
						// This is an import from a node_modules, let
						// ESBuild resolve this.
						return api.OnResolveResult{}, nil
					}
				} else {
					relDir = filepath.Dir(rel)
				}
			}
		} else {
			relDir = opts.SourceDir
		}

		// Imports not starting with a "." is assumed to live relative to /assets.
		// Hugo makes no assumptions about the directory structure below /assets.
		if relDir != "" && strings.HasPrefix(impPath, ".") {
			impPath = filepath.Join(relDir, impPath)
		}

		m := resolveComponentInAssets(fs.Fs, impPath)

		if m != nil {
			depsManager.AddIdentity(m.PathInfo)

			// Store the source root so we can create a jsconfig.json
			// to help IntelliSense when the build is done.
			// This should be a small number of elements, and when
			// in server mode, we may get stale entries on renames etc.,
			// but that shouldn't matter too much.
			rs.JSConfigBuilder.AddSourceRoot(m.SourceRoot)
			return api.OnResolveResult{Path: m.Filename, Namespace: NsHugoImport}, nil
		}

		// Fall back to ESBuild's resolve.
		return api.OnResolveResult{}, nil
	}

	importResolver := api.Plugin{
		Name: "hugo-import-resolver",
		Setup: func(build api.PluginBuild) {
			build.OnResolve(api.OnResolveOptions{Filter: `.*`},
				func(args api.OnResolveArgs) (api.OnResolveResult, error) {
					return resolveImport(args)
				})
			build.OnLoad(api.OnLoadOptions{Filter: `.*`, Namespace: NsHugoImport},
				func(args api.OnLoadArgs) (api.OnLoadResult, error) {
					b, err := os.ReadFile(args.Path)
					if err != nil {
						return api.OnLoadResult{}, fmt.Errorf("failed to read %q: %w", args.Path, err)
					}
					c := string(b)

					return api.OnLoadResult{
						// See https://github.com/evanw/esbuild/issues/502
						// This allows all modules to resolve dependencies
						// in the main project's node_modules.
						ResolveDir: opts.ResolveDir,
						Contents:   &c,
						Loader:     loaderFromFilename(args.Path),
					}, nil
				})
			build.OnLoad(api.OnLoadOptions{Filter: `.*`, Namespace: NsHugoImportResolveFunc},
				func(args api.OnLoadArgs) (api.OnLoadResult, error) {
					c := opts.ImportOnLoadFunc(args)
					if c == "" {
						return api.OnLoadResult{}, fmt.Errorf("ImportOnLoadFunc failed to resolve %q", args.Path)
					}

					return api.OnLoadResult{
						ResolveDir: opts.ResolveDir,
						Contents:   &c,
						Loader:     loaderFromFilename(args.Path),
					}, nil
				})
		},
	}

	params := opts.Params
	if params == nil {
		// This way @params will always resolve to something.
		params = make(map[string]any)
	}

	b, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal params: %w", err)
	}
	paramsPlugin := api.Plugin{
		Name: "hugo-params-plugin",
		Setup: func(build api.PluginBuild) {
			build.OnResolve(api.OnResolveOptions{Filter: `^@params$`},
				func(args api.OnResolveArgs) (api.OnResolveResult, error) {
					return api.OnResolveResult{
						Path:      args.Importer,
						Namespace: nsHugoParams,
					}, nil
				})
			build.OnLoad(api.OnLoadOptions{Filter: `.*`, Namespace: nsHugoParams},
				func(args api.OnLoadArgs) (api.OnLoadResult, error) {
					bb := b

					if opts.ImportParamsOnLoadFunc != nil {
						if bbb := opts.ImportParamsOnLoadFunc(args); bbb != nil {
							bb = bbb
						}
					}

					s := string(bb)

					return api.OnLoadResult{
						Contents: &s,
						Loader:   api.LoaderJSON,
					}, nil
				})
		},
	}

	return []api.Plugin{importResolver, paramsPlugin}, nil
}
