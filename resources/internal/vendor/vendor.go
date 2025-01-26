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
	"encoding/hex"
	"io"
	"path/filepath"
	"sort"
	"strings"

	"github.com/cespare/xxhash/v2"
	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/common/types"
	"github.com/gohugoio/hugo/helpers"
	"github.com/spf13/afero"
	"github.com/spf13/cast"
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
	VendorKey() string
}

// VendorKeyFromOpts extracts the vendor key from the given options map.
func VendorKeyFromOpts(m map[string]any) string {
	const vendorKeyKey = "vendorKey"

	for k, v := range m {
		if strings.EqualFold(k, vendorKeyKey) {
			return cast.ToString(v)
		}
	}
	return ""
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

func NewVendorer(fs afero.Fs, workingDir string) *Vendorer {
	return &Vendorer{
		fs:         fs,
		workingDir: workingDir,
	}
}

type Vendorer struct {
	fs         afero.Fs
	workingDir string

	VendoredResources []VendoredResource
}

func (v *Vendorer) Finalize() error {
	// Sort the vendored resources to make the order deterministic.
	sort.Slice(v.VendoredResources, func(i, j int) bool {
		return v.VendoredResources[i].Path < v.VendoredResources[j].Path
	})
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
	VendorRoot      = "_vendor"
	vendorModules   = "mod"
	vendorResources = "res"
	vendorScopeAll  = "_all"
)

type VendoredResource struct {
	Path string `json:"path"`
	Hash string `json:"hash"`
}

// TODO1 add --vendorScope=* or {development,production} e.g. Encode it into the root folder. Or: At least create a "all" root folder for future use.
// add resources.txt to root with file/hash listing + version.
// TODO1 remove me.
func (v *Vendorer) NewResourceVendor(opts ResourceVendorOptions) (*ResourceVendor, error) {
	digest := xxhash.New()
	finalFrom := io.TeeReader(opts.From, digest)
	vendorFileName := filepath.Join(v.workingDir, VendorRoot, vendorResources, vendorScopeAll, opts.Target.VendorName(), opts.Target.VendorKey(), opts.InPath)
	vendoredFile, err := helpers.OpenFileForWriting(v.fs, vendorFileName)
	if err != nil {
		return nil, err
	}
	finalTo := io.MultiWriter(opts.To, vendoredFile)
	return &ResourceVendor{
		Digest:    digest,
		FinalFrom: finalFrom,
		FinalTo:   finalTo,
		closeFunc: vendoredFile.Close,
	}, nil
}

func (v *Vendorer) OpenVendoredFileForWriting(opts ResourceVendorOptions) (io.WriteCloser, hugio.OpenReadSeekCloser, error) {
	// TODO1 make the fs relative to the working dir.
	vendorFileName := v.vendorFilename(opts)
	vendorFilePath := paths.TrimLeading(filepath.ToSlash(strings.TrimPrefix(vendorFileName, v.workingDir)))
	f, err := helpers.OpenFileForWriting(v.fs, vendorFileName)
	if err != nil {
		return nil, nil, err
	}
	open := v.VendoredOpenReadSeekCloser(opts)
	h := xxhash.New()
	closer := types.CloserFunc(func() error {
		hash := h.Sum(nil)
		vendoredResource := VendoredResource{
			Path: vendorFilePath,
			Hash: hex.EncodeToString(hash),
		}
		v.VendoredResources = append(v.VendoredResources, vendoredResource)
		return nil
	})

	w := hugio.NewMultiWriteCloser(
		f,
		hugio.NewWriteCloser(h, closer),
	)

	return w, open, nil
}

// OpenVendoredFile opens a vendored file for reading or nil if not found.
func (v *Vendorer) OpenVendoredFile(opts ResourceVendorOptions) (hugio.ReadSeekCloser, hugio.OpenReadSeekCloser, error) {
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

func (v *Vendorer) VendoredOpenReadSeekCloser(opts ResourceVendorOptions) hugio.OpenReadSeekCloser {
	return func() (hugio.ReadSeekCloser, error) {
		vendorFileName := v.vendorFilename(opts)
		return v.fs.Open(vendorFileName)
	}
}

func (v *Vendorer) vendorFilename(opts ResourceVendorOptions) string {
	vendorKey := paths.NormalizePathStringBasic(opts.Target.VendorKey())
	n := filepath.Join(v.workingDir, VendorRoot, vendorResources, vendorScopeAll, opts.Target.VendorName(), vendorKey, opts.InPath)
	return n
}
