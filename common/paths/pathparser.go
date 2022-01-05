// Copyright 2021 The Hugo Authors. All rights reserved.
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

package paths

import (
	"errors"
	"os"
	"runtime"
	"strings"

	"github.com/gohugoio/hugo/common/types"
	"github.com/gohugoio/hugo/hugofs/files"
	"github.com/gohugoio/hugo/identity"
)

var ForComponent = func(component string) func(b *pathBase) {
	return func(b *pathBase) {
		b.component = component
	}
}

// Parse parses s into Path using Hugo's content path rules.
func Parse(s string, parseOpts ...func(b *pathBase)) Path {
	p, err := parse(s, parseOpts...)
	if err != nil {
		panic(err)
	}
	return p
}

func parse(s string, parseOpts ...func(b *pathBase)) (*pathBase, error) {
	p := &pathBase{
		component: files.ComponentFolderContent,
		posBase:   -1,
	}

	for _, opt := range parseOpts {
		opt(p)
	}

	// All lower case.
	s = strings.ToLower(s)

	// Leading slash, no trailing slash.
	if p.component != files.ComponentFolderLayouts && !strings.HasPrefix(s, "/") {
		s = "/" + s
	}

	if s != "/" && s[len(s)-1] == '/' {
		s = s[:len(s)-1]
	}

	p.s = s

	isWindows := runtime.GOOS == "windows"

	for i := len(s) - 1; i >= 0; i-- {
		c := s[i]

		if isWindows && c == os.PathSeparator {
			return nil, errors.New("only forward slashes allowed")
		}

		switch c {
		case '.':
			if p.posBase == -1 {
				var high int
				if len(p.identifiers) > 0 {
					high = p.identifiers[len(p.identifiers)-1].Low - 1
				} else {
					high = len(p.s)
				}
				p.identifiers = append(p.identifiers, types.LowHigh{Low: i + 1, High: high})
			}
		case '/':
			if p.posBase == -1 {
				p.posBase = i + 1
			}
		}
	}

	p.isContent = p.component == files.ComponentFolderContent && files.IsContentExt(p.Ext())

	if p.isContent {
		id := p.identifiers[len(p.identifiers)-1]
		b := p.s[p.posBase : id.Low-1]
		switch b {
		case "index":
			p.bundleType = BundleTypeLeaf
		case "_index":
			p.bundleType = BundleTypeBranch
		default:
			p.bundleType = BundleTypeNone
		}
	}

	return p, nil
}

type Path interface {
	identity.Identity
	Component() string
	Name() string
	Base() string
	Dir() string
	Ext() string
	Slice(bottom, top int) string
	Identifiers() []string
	Identifier(i int) string
	IsContent() bool
	IsBundle() bool
	IsLeafBundle() bool
	IsBranchBundle() bool
	BundleType() BundleType
}

func ModifyPathBundleNone(p Path) {
	p.(*pathBase).bundleType = BundleTypeNone
}

type PathInfos []PathInfo

type PathInfo interface {
	Path
	Filename() string
}

type BundleType int

const (
	BundleTypeNone BundleType = iota
	BundleTypeLeaf
	BundleTypeBranch
)

type pathBase struct {
	s string

	posBase int

	component  string
	isContent  bool
	bundleType BundleType

	identifiers []types.LowHigh
}

type pathInfo struct {
	Path
	component string
	filename  string
}

func (p *pathInfo) Filename() string {
	return p.filename
}

func WithInfo(p Path, filename string) PathInfo {
	return &pathInfo{
		Path:     p,
		filename: filename,
	}
}

// IdentifierBase satifies identity.Identity.
// TODO1 componnt?
func (p *pathBase) IdentifierBase() interface{} {
	return p.Base()
}

func (p *pathBase) Component() string {
	return p.component
}

func (p *pathBase) IsContent() bool {
	return p.isContent
}

// Name returns the last element of path.
func (p *pathBase) Name() string {
	if p.posBase > 0 {
		return p.s[p.posBase:]
	}
	return p.s
}

func (p *pathBase) Dir() string {
	if p.posBase > 0 {
		return p.s[:p.posBase-1]
	}
	return "/"
}

func (p *pathBase) Slice(bottom, top int) string {
	if bottom == 0 && top == 0 {
		return p.s
	}

	if bottom < 0 {
		bottom = 0
	}

	if top < 0 {
		top = 0
	}

	if bottom > len(p.identifiers)+1 {
		bottom = len(p.identifiers) + 1
	}

	if top > len(p.identifiers)+1 {
		top = len(p.identifiers) + 1
	}

	// 0 : posBase
	// posBase : identifier[0].Low
	// identifier[n].Low : identifier[n].High
	var low, high int
	if bottom == 1 {
		low = p.posBase
	} else if bottom > 1 {
		low = p.identifiers[len(p.identifiers)-bottom+1].Low
	}

	if top == 0 {
		high = len(p.s)
	} else if top > 0 {
		i := top
		distance := len(p.identifiers) - i
		if distance <= 0 {
			if distance == 0 {
				high = p.identifiers[len(p.identifiers)-1].Low - 1
			} else {
				high = p.posBase - 1
			}
		} else {
			high = p.identifiers[i].High
		}
	}

	if low > high {
		return ""
	}

	return p.s[low:high]
}

// For content files, Base returns the path without any identifiers (extension, language code etc.).
// Any 'index' as the last path element is ignored.
//
// For other files (Resources), any extension is kept.
func (p *pathBase) Base() string {
	if len(p.identifiers) > 0 {
		if !p.isContent && len(p.identifiers) == 1 {
			// Preserve extension.
			return p.s
		}

		id := p.identifiers[len(p.identifiers)-1]
		high := id.Low - 1
		if p.isContent {
			if p.IsBundle() {
				high = p.posBase - 1
			}
		}

		if p.isContent {
			return p.s[:high]
		}

		// For txt files etc. we want to preserve the extension.
		id = p.identifiers[0]

		return p.s[:high] + p.s[id.Low-1:id.High]
	}
	return p.s
}

func (p *pathBase) Ext() string {
	return p.identifierAsString(0)
}

func (p *pathBase) Identifier(i int) string {
	return p.identifierAsString(i)
}

func (p *pathBase) Identifiers() []string {
	ids := make([]string, len(p.identifiers))
	for i, id := range p.identifiers {
		ids[i] = p.s[id.Low:id.High]
	}
	return ids
}

func (p *pathBase) BundleType() BundleType {
	return p.bundleType
}

func (p *pathBase) IsBundle() bool {
	return p.bundleType != BundleTypeNone
}

func (p *pathBase) IsBranchBundle() bool {
	return p.bundleType == BundleTypeBranch
}

func (p *pathBase) IsLeafBundle() bool {
	return p.bundleType == BundleTypeLeaf
}

func (p *pathBase) identifierAsString(i int) string {
	if i < 0 || i >= len(p.identifiers) {
		return ""
	}
	id := p.identifiers[i]
	return p.s[id.Low:id.High]
}
