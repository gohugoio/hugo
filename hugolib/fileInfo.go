// Copyright 2017-present The Hugo Authors. All rights reserved.
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

package hugolib

import (
	"strings"

	"github.com/gohugoio/hugo/hugofs/files"

	"github.com/pkg/errors"

	"github.com/gohugoio/hugo/hugofs"

	"github.com/spf13/afero"

	"github.com/gohugoio/hugo/source"
)

// fileInfo implements the File and ReadableFile interface.
var (
	_ source.File = (*fileInfo)(nil)
)

type fileInfo struct {
	source.File

	overriddenLang string
}

func (fi *fileInfo) Open() (afero.File, error) {
	f, err := fi.FileInfo().Meta().Open()
	if err != nil {
		err = errors.Wrap(err, "fileInfo")
	}

	return f, err
}

func (fi *fileInfo) Lang() string {
	if fi.overriddenLang != "" {
		return fi.overriddenLang
	}
	return fi.File.Lang()
}

func (fi *fileInfo) String() string {
	if fi == nil || fi.File == nil {
		return ""
	}
	return fi.Path()
}

// TODO(bep) rename
func newFileInfo(sp *source.SourceSpec, fi hugofs.FileMetaInfo) (*fileInfo, error) {

	baseFi, err := sp.NewFileInfo(fi)
	if err != nil {
		return nil, err
	}

	f := &fileInfo{
		File: baseFi,
	}

	return f, nil

}

type bundleDirType int

const (
	bundleNot bundleDirType = iota

	// All from here are bundles in one form or another.
	bundleLeaf
	bundleBranch
)

// Returns the given file's name's bundle type and whether it is a content
// file or not.
func classifyBundledFile(name string) (bundleDirType, bool) {
	if !files.IsContentFile(name) {
		return bundleNot, false
	}
	if strings.HasPrefix(name, "_index.") {
		return bundleBranch, true
	}

	if strings.HasPrefix(name, "index.") {
		return bundleLeaf, true
	}

	return bundleNot, true
}

func (b bundleDirType) String() string {
	switch b {
	case bundleNot:
		return "Not a bundle"
	case bundleLeaf:
		return "Regular bundle"
	case bundleBranch:
		return "Branch bundle"
	}

	return ""
}
