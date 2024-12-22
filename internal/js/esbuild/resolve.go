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
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/resources/resource"
	"github.com/spf13/afero"
)

const (
	NsHugoImport            = "ns-hugo-imp"
	NsHugoImportResolveFunc = "ns-hugo-imp-func"
	nsHugoParams            = "ns-hugo-params"
	pathHugoConfigParams    = "@params/config"

	stdinImporter = "<stdin>"
)

var hugoNamespaces = []string{NsHugoImport, NsHugoImportResolveFunc, nsHugoParams}

const (
	PrefixHugoVirtual = "__hu_v"
	PrefixHugoMemory  = "__hu_m"
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

// This is a common sub-set of ESBuild's default extensions.
// We assume that imports of JSON, CSS etc. will be using their full
// name with extension.
var commonExtensions = []string{".js", ".ts", ".tsx", ".jsx"}

// ResolveComponent resolves a component using the given resolver.
func ResolveComponent[T any](impPath string, resolve func(string) (v T, found, isDir bool)) (v T, found bool) {
	findFirst := func(base string) (v T, found, isDir bool) {
		for _, ext := range commonExtensions {
			if strings.HasSuffix(impPath, ext) {
				// Import of foo.js.js need the full name.
				continue
			}
			if v, found, isDir = resolve(base + ext); found {
				return
			}
		}

		// Not found.
		return
	}

	// We need to check if this is a regular file imported without an extension.
	// There may be ambiguous situations where both foo.js and foo/index.js exists.
	// This import order is in line with both how Node and ESBuild's native
	// import resolver works.

	// It may be a regular file imported without an extension, e.g.
	// foo or foo/index.
	v, found, _ = findFirst(impPath)
	if found {
		return v, found
	}

	base := filepath.Base(impPath)
	if base == "index" {
		// try index.esm.js etc.
		v, found, _ = findFirst(impPath + ".esm")
		if found {
			return v, found
		}
	}

	// Check the path as is.
	var isDir bool
	v, found, isDir = resolve(impPath)
	if found && isDir {
		v, found, _ = findFirst(filepath.Join(impPath, "index"))
		if !found {
			v, found, _ = findFirst(filepath.Join(impPath, "index.esm"))
		}
	}

	if !found && strings.HasSuffix(base, ".js") {
		v, found, _ = findFirst(strings.TrimSuffix(impPath, ".js"))
	}

	return
}

// ResolveResource resolves a resource using the given resourceGetter.
func ResolveResource(impPath string, resourceGetter resource.ResourceGetter) (r resource.Resource) {
	resolve := func(name string) (v resource.Resource, found, isDir bool) {
		r := resourceGetter.Get(name)
		return r, r != nil, false
	}
	r, found := ResolveComponent(impPath, resolve)
	if !found {
		return nil
	}
	return r
}

func newFSResolver(fs afero.Fs) *fsResolver {
	return &fsResolver{fs: fs, resolved: maps.NewCache[string, *hugofs.FileMeta]()}
}

type fsResolver struct {
	fs       afero.Fs
	resolved *maps.Cache[string, *hugofs.FileMeta]
}

func (r *fsResolver) resolveComponent(impPath string) *hugofs.FileMeta {
	v, _ := r.resolved.GetOrCreate(impPath, func() (*hugofs.FileMeta, error) {
		resolve := func(name string) (*hugofs.FileMeta, bool, bool) {
			if fi, err := r.fs.Stat(name); err == nil {
				return fi.(hugofs.FileMetaInfo).Meta(), true, fi.IsDir()
			}
			return nil, false, false
		}
		v, _ := ResolveComponent(impPath, resolve)
		return v, nil
	})
	return v
}

func createBuildPlugins(rs *resources.Spec, assetsResolver *fsResolver, depsManager identity.Manager, opts Options) ([]api.Plugin, error) {
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

		for _, ext := range opts.Externals {
			// ESBuild will do a more thorough check for packages resolved in node_modules,
			// but we need to make sure that we don't try to resolve these in the /assets folder.
			if ext == impPath {
				return api.OnResolveResult{
					Path:     impPath,
					External: true,
				}, nil
			}
		}

		if opts.ImportOnResolveFunc != nil {
			if s := opts.ImportOnResolveFunc(impPath, args); s != "" {
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

		m := assetsResolver.resolveComponent(impPath)

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
						Loader:     opts.loaderFromFilename(args.Path),
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
						Loader:     opts.loaderFromFilename(args.Path),
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
			build.OnResolve(api.OnResolveOptions{Filter: `^@params(/config)?$`},
				func(args api.OnResolveArgs) (api.OnResolveResult, error) {
					resolvedPath := args.Importer

					if args.Path == pathHugoConfigParams {
						resolvedPath = pathHugoConfigParams
					}

					return api.OnResolveResult{
						Path:      resolvedPath,
						Namespace: nsHugoParams,
					}, nil
				})
			build.OnLoad(api.OnLoadOptions{Filter: `.*`, Namespace: nsHugoParams},
				func(args api.OnLoadArgs) (api.OnLoadResult, error) {
					bb := b
					if args.Path != pathHugoConfigParams && opts.ImportParamsOnLoadFunc != nil {
						bb = opts.ImportParamsOnLoadFunc(args)
					}
					s := string(bb)

					if s == "" {
						s = "{}"
					}

					return api.OnLoadResult{
						Contents: &s,
						Loader:   api.LoaderJSON,
					}, nil
				})
		},
	}

	return []api.Plugin{importResolver, paramsPlugin}, nil
}
