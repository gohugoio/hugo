// Copyright 2025 The Hugo Authors. All rights reserved.
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

package vendor

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/cespare/xxhash/v2"
	"github.com/gobwas/glob"
	"github.com/gohugoio/hugo/common/hashing"
	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/common/types"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/afero"
)

var _ io.Closer = (*ResourceVendor)(nil)

type Vendorable interface {
	// The fully qualified name, using forward slash as path separator, of the function or method that produces this vendorable resource.
	// Lower case, e.g. "css/tailwindcss".
	VendorName() string

	// The key that identifies the resource transformation.
	// Make this shallow so a transformation may be shared between environments, if needed,
	// do not include any filenames/content hashes here.
	// Typically you can use the VendorKeyFromOpts function to extract this from the options.
	VendorScope() map[string]any
}

type VendorScope struct {
	// An optional key that identifies the resource transformation variant.
	// The default vendor path is deliberately shallow, so this allows multiple vendored variants of the same resource transformation
	// with different configurations.
	Key string `json:"key"`

	// A Glob pattern matching the build environment, e.g. “{production,development}” or “*”. The default is “*”.
	Environment string `json:"environment"`
}

// VendorScopeFromOpts extracts the vendor scoped from the given options map.
func VendorScopeFromOpts(m map[string]any) map[string]any {
	const vendorScopeKey = "vendorScope"

	for k, v := range m {
		if strings.EqualFold(k, vendorScopeKey) {
			return maps.ToStringMap(v)
		}
	}
	return nil
}

type ResourceVendor struct {
	// For the from hash.
	Digest *xxhash.Digest

	closeFunc func() error

	VendoredFile hugio.ReadSeekCloser

	// TODO1 names.
	FinalFrom io.Reader
	FinalTo   io.Writer
}

func (v *ResourceVendor) Close() error {
	if v.closeFunc != nil {
		return v.closeFunc()
	}
	return nil
}

func NewVendorer(vendorFs, sourceFs afero.Fs, environment string) (*ResourceVendorer, error) {
	rv := &ResourceVendorer{
		vendorFs:    vendorFs,
		sourceFs:    sourceFs,
		environment: environment,
	}

	if err := rv.init(); err != nil {
		return nil, err
	}

	return rv, nil
}

type ResourceVendorer struct {
	vendorFs    afero.Fs // Fs relative to the vendor root, top layer writatble.
	sourceFs    afero.Fs // Usually OS filesystem.
	environment string

	mu sync.Mutex

	output            outputResources
	vendoredResources map[string]*maps.Ordered[string, vendoredResource]
}

type outputResources struct {
	Resources []outputResource `json:"resources"`
}

func (v *ResourceVendorer) init() error {
	dir, err := v.vendorFs.Open(vendorResources)
	if err != nil {
		if herrors.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer dir.Close()
	fis, err := dir.(fs.ReadDirFile).ReadDir(-1)
	if err != nil {
		return err
	}

	v.mu.Lock()
	defer v.mu.Unlock()
	v.vendoredResources = make(map[string]*maps.Ordered[string, vendoredResource])

	for _, de := range fis {
		if de.Name() == vendorResourcesJSON {
			fim := de.(hugofs.FileMetaInfo)
			f, err := fim.Meta().Open()
			if err != nil {
				return err
			}
			defer f.Close()

			var resources outputResources
			if err := json.NewDecoder(f).Decode(&resources); err != nil {
				return err
			}
			vendorDir := filepath.Dir(fim.Meta().Filename)
			for _, r := range resources.Resources {
				vr := vendoredResource{
					resource:  r,
					vendorDir: vendorDir,
				}
				m := v.vendoredResources[r.BasePath]
				if m == nil {
					m = maps.NewOrdered[string, vendoredResource]()
					v.vendoredResources[r.BasePath] = m
				}
				if !m.Contains(r.ScopeHash) {
					if err := vr.init(); err != nil {
						return err
					}
					m.Set(r.ScopeHash, vr)
				}
			}

		}
	}
	return nil
}

type vendoredResource struct {
	resource  outputResource
	vendorDir string

	matchEnvironment glob.Glob
}

func (v *vendoredResource) init() error {
	v.matchEnvironment = glob.MustCompile(v.resource.Scope.Environment)
	return nil
}

func (v *ResourceVendorer) Finalize() error {
	// Sort the vendored resources to make the order deterministic.
	sort.Slice(v.output.Resources, func(i, j int) bool {
		return v.output.Resources[i].BasePath < v.output.Resources[j].BasePath
	})

	vendorFilename := filepath.Join(vendorResources, vendorResourcesJSON)

	f, err := v.vendorFs.Create(vendorFilename)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v.output); err != nil {
		return err
	}

	return nil
}

type ResourceVendorOptions struct {
	Target Vendorable
	Name   string
	InPath string
	From   io.Reader
	To     io.Writer
}

const (
	vendorRoot          = "_vendor"
	vendorModules       = "modules"
	vendorResources     = "resources" // TODO1 + modules.
	vendorResourcesJSON = "resources.json"
)

type outputResource struct {
	// BasePath to the vendored resource, relative to the vendor root.
	// Unix style path.
	// The full path starting from the vendor root is BasePath/ScopeHash/Path.
	BasePath string `json:"basePath"`

	// Path the last path element of the vendored resource.
	Path string `json:"path"`

	Scope     VendorScope `json:"scope"`
	ScopeHash string      `json:"scopeHash"`
}

func (v *ResourceVendorer) OpenVendoredFileForWriting(opts ResourceVendorOptions) (io.WriteCloser, hugio.OpenReadSeekCloser, error) {
	vs := opts.Target.VendorScope()
	scopeHash := hashing.HashStringHex(vs)
	vendorScope := VendorScope{
		Environment: "*",
	}
	if err := mapstructure.WeakDecode(vs, &vendorScope); err != nil {
		return nil, nil, err
	}

	vendorBasePath := v.vendorPath(opts)
	vendorDir := filepath.Join(vendorBasePath, scopeHash)
	if err := v.vendorFs.MkdirAll(vendorDir, 0o755); err != nil {
		return nil, nil, fmt.Errorf("failed to create directory %q: %w", vendorBasePath, err)
	}
	vendorFilename := filepath.Join(vendorDir, opts.InPath)
	f, err := helpers.OpenFileForWriting(v.vendorFs, vendorFilename)
	if err != nil {
		return nil, nil, err
	}
	open := func() (hugio.ReadSeekCloser, error) {
		return v.vendorFs.Open(vendorFilename)
	}

	closer := types.CloserFunc(func() error {
		vendoredResource := outputResource{
			BasePath:  vendorBasePath,
			Path:      opts.InPath,
			Scope:     vendorScope,
			ScopeHash: scopeHash,
		}
		v.mu.Lock()
		v.output.Resources = append(v.output.Resources, vendoredResource)
		v.mu.Unlock()

		return f.Close()
	})

	w := hugio.NewWriteCloser(f, closer)

	return w, open, nil
}

// OpenVendoredFile opens a vendored file for reading or nil if not found.
func (v *ResourceVendorer) OpenVendoredFile(opts ResourceVendorOptions) (hugio.ReadSeekCloser, hugio.OpenReadSeekCloser, error) {
	open := v.VendoredOpenReadSeekCloser(opts)
	f, err := open()
	if err != nil {
		if herrors.IsNotExist(err) {
			return nil, nil, nil
		}
		return nil, nil, err
	}
	return f, open, nil
}

func (v *ResourceVendorer) VendoredOpenReadSeekCloser(opts ResourceVendorOptions) hugio.OpenReadSeekCloser {
	vendorFilePath := v.vendorPath(opts)

	r, found := v.vendoredResources[vendorFilePath]

	return func() (hugio.ReadSeekCloser, error) {
		if !found {
			return nil, afero.ErrFileNotFound
		}

		var filename string
		r.Range(func(key string, vr vendoredResource) bool {
			if vr.matchEnvironment.Match(v.environment) {
				r := vr.resource
				filename = filepath.Join(vr.vendorDir, vendorFilePath, r.ScopeHash, r.Path)
				return false
			}
			return true
		})
		if filename == "" {
			return nil, afero.ErrFileNotFound
		}
		// resources/css/tailwindcss
		return v.sourceFs.Open(filename)
	}
}

func (v *ResourceVendorer) vendorPath(opts ResourceVendorOptions) string {
	n := filepath.ToSlash(path.Join(vendorResources, opts.Target.VendorName()))
	return n
}
