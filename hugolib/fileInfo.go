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

	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/source"
)

// fileInfo implements the File and ReadableFile interface.
var (
	_ source.File         = (*fileInfo)(nil)
	_ source.ReadableFile = (*fileInfo)(nil)
	_ pathLangFile        = (*fileInfo)(nil)
)

// A partial interface to prevent ambigous compiler error.
type basePather interface {
	Filename() string
	RealName() string
	BaseDir() string
}

type fileInfo struct {
	bundleTp bundleDirType

	source.ReadableFile
	basePather

	overriddenLang string

	// Set if the content language for this file is disabled.
	disabled bool
}

func (fi *fileInfo) Lang() string {
	if fi.overriddenLang != "" {
		return fi.overriddenLang
	}
	return fi.ReadableFile.Lang()
}

func (fi *fileInfo) Filename() string {
	return fi.basePather.Filename()
}

func (fi *fileInfo) isOwner() bool {
	return fi.bundleTp > bundleNot
}

func isContentFile(filename string) bool {
	return contentFileExtensionsSet[strings.TrimPrefix(helpers.Ext(filename), ".")]
}

func (fi *fileInfo) isContentFile() bool {
	return contentFileExtensionsSet[fi.Ext()]
}

func newFileInfo(sp *source.SourceSpec, baseDir, filename string, fi pathLangFileFi, tp bundleDirType) *fileInfo {

	baseFi := sp.NewFileInfo(baseDir, filename, tp == bundleLeaf, fi)
	f := &fileInfo{
		bundleTp:     tp,
		ReadableFile: baseFi,
		basePather:   fi,
	}

	lang := f.Lang()
	f.disabled = lang != "" && sp.DisabledLanguages[lang]

	return f

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
	if !isContentFile(name) {
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
