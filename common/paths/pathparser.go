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

package paths

import (
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/gohugoio/hugo/common/types"
	"github.com/gohugoio/hugo/hugofs/files"
	"github.com/gohugoio/hugo/identity"
)

// PathParser parses a path into a Path.
type PathParser struct {
	// Maps the language code to its index in the languages/sites slice.
	LanguageIndex map[string]int

	// Reports whether the given language is disabled.
	IsLangDisabled func(string) bool

	// Reports whether the given ext is a content file.
	IsContentExt func(string) bool
}

// NormalizePathString returns a normalized path string using the very basic Hugo rules.
func NormalizePathStringBasic(s string) string {
	// All lower case.
	s = strings.ToLower(s)

	// Replace spaces with hyphens.
	s = strings.ReplaceAll(s, " ", "-")

	return s
}

// ParseIdentity parses component c with path s into a StringIdentity.
func (pp *PathParser) ParseIdentity(c, s string) identity.StringIdentity {
	p := pp.parsePooled(c, s)
	defer putPath(p)
	return identity.StringIdentity(p.IdentifierBase())
}

// ParseBaseAndBaseNameNoIdentifier parses component c with path s into a base and a base name without any identifier.
func (pp *PathParser) ParseBaseAndBaseNameNoIdentifier(c, s string) (string, string) {
	p := pp.parsePooled(c, s)
	defer putPath(p)
	return p.Base(), p.BaseNameNoIdentifier()
}

func (pp *PathParser) parsePooled(c, s string) *Path {
	s = NormalizePathStringBasic(s)
	p := getPath()
	p.component = c
	p, err := pp.doParse(c, s, p)
	if err != nil {
		panic(err)
	}
	return p
}

// Parse parses component c with path s into Path using Hugo's content path rules.
func (pp *PathParser) Parse(c, s string) *Path {
	p, err := pp.parse(c, s)
	if err != nil {
		panic(err)
	}
	return p
}

func (pp *PathParser) newPath(component string) *Path {
	return &Path{
		component:             component,
		posContainerLow:       -1,
		posContainerHigh:      -1,
		posSectionHigh:        -1,
		posIdentifierLanguage: -1,
	}
}

func (pp *PathParser) parse(component, s string) (*Path, error) {
	ss := NormalizePathStringBasic(s)

	p, err := pp.doParse(component, ss, pp.newPath(component))
	if err != nil {
		return nil, err
	}

	if s != ss {
		var err error
		// Preserve the original case for titles etc.
		p.unnormalized, err = pp.doParse(component, s, pp.newPath(component))
		if err != nil {
			return nil, err
		}
	} else {
		p.unnormalized = p
	}

	return p, nil
}

func (pp *PathParser) doParse(component, s string, p *Path) (*Path, error) {
	hasLang := pp.LanguageIndex != nil
	hasLang = hasLang && (component == files.ComponentFolderContent || component == files.ComponentFolderLayouts)

	if runtime.GOOS == "windows" {
		s = path.Clean(filepath.ToSlash(s))
		if s == "." {
			s = ""
		}
	}

	if s == "" {
		s = "/"
	}

	// Leading slash, no trailing slash.
	if !strings.HasPrefix(s, "/") {
		s = "/" + s
	}

	if s != "/" && s[len(s)-1] == '/' {
		s = s[:len(s)-1]
	}

	p.s = s
	slashCount := 0

	for i := len(s) - 1; i >= 0; i-- {
		c := s[i]

		switch c {
		case '.':
			if p.posContainerHigh == -1 {
				var high int
				if len(p.identifiers) > 0 {
					high = p.identifiers[len(p.identifiers)-1].Low - 1
				} else {
					high = len(p.s)
				}
				id := types.LowHigh[string]{Low: i + 1, High: high}
				if len(p.identifiers) == 0 {
					p.identifiers = append(p.identifiers, id)
				} else if len(p.identifiers) == 1 {
					// Check for a valid language.
					s := p.s[id.Low:id.High]

					if hasLang {
						var disabled bool
						_, langFound := pp.LanguageIndex[s]
						if !langFound {
							disabled = pp.IsLangDisabled != nil && pp.IsLangDisabled(s)
							if disabled {
								p.disabled = true
								langFound = true
							}
						}
						if langFound {
							p.posIdentifierLanguage = 1
							p.identifiers = append(p.identifiers, id)
						}
					}
				}
			}
		case '/':
			slashCount++
			if p.posContainerHigh == -1 {
				p.posContainerHigh = i + 1
			} else if p.posContainerLow == -1 {
				p.posContainerLow = i + 1
			}
			if i > 0 {
				p.posSectionHigh = i
			}
		}
	}

	if len(p.identifiers) > 0 {
		isContentComponent := p.component == files.ComponentFolderContent || p.component == files.ComponentFolderArchetypes
		isContent := isContentComponent && pp.IsContentExt(p.Ext())
		id := p.identifiers[len(p.identifiers)-1]
		b := p.s[p.posContainerHigh : id.Low-1]
		if isContent {
			switch b {
			case "index":
				p.bundleType = PathTypeLeaf
			case "_index":
				p.bundleType = PathTypeBranch
			default:
				p.bundleType = PathTypeContentSingle
			}

			if slashCount == 2 && p.IsLeafBundle() {
				p.posSectionHigh = 0
			}
		} else if b == files.NameContentData && files.IsContentDataExt(p.Ext()) {
			p.bundleType = PathTypeContentData
		}
	}

	return p, nil
}

func ModifyPathBundleTypeResource(p *Path) {
	if p.IsContent() {
		p.bundleType = PathTypeContentResource
	} else {
		p.bundleType = PathTypeFile
	}
}

type PathType int

const (
	// A generic resource, e.g. a JSON file.
	PathTypeFile PathType = iota

	// All below are content files.
	// A resource of a content type with front matter.
	PathTypeContentResource

	// E.g. /blog/my-post.md
	PathTypeContentSingle

	// All below are bundled content files.

	// Leaf bundles, e.g. /blog/my-post/index.md
	PathTypeLeaf

	// Branch bundles, e.g. /blog/_index.md
	PathTypeBranch

	// Content data file, _content.gotmpl.
	PathTypeContentData
)

type Path struct {
	// Note: Any additions to this struct should also be added to the pathPool.
	s string

	posContainerLow  int
	posContainerHigh int
	posSectionHigh   int

	component  string
	bundleType PathType

	identifiers []types.LowHigh[string]

	posIdentifierLanguage int
	disabled              bool

	trimLeadingSlash bool

	unnormalized *Path
}

var pathPool = &sync.Pool{
	New: func() any {
		p := &Path{}
		p.reset()
		return p
	},
}

func getPath() *Path {
	return pathPool.Get().(*Path)
}

func putPath(p *Path) {
	p.reset()
	pathPool.Put(p)
}

func (p *Path) reset() {
	p.s = ""
	p.posContainerLow = -1
	p.posContainerHigh = -1
	p.posSectionHigh = -1
	p.component = ""
	p.bundleType = 0
	p.identifiers = p.identifiers[:0]
	p.posIdentifierLanguage = -1
	p.disabled = false
	p.trimLeadingSlash = false
	p.unnormalized = nil
}

// TrimLeadingSlash returns a copy of the Path with the leading slash removed.
func (p Path) TrimLeadingSlash() *Path {
	p.trimLeadingSlash = true
	return &p
}

func (p *Path) norm(s string) string {
	if p.trimLeadingSlash {
		s = strings.TrimPrefix(s, "/")
	}
	return s
}

// IdentifierBase satisfies identity.Identity.
func (p *Path) IdentifierBase() string {
	return p.Base()
}

// Component returns the component for this path (e.g. "content").
func (p *Path) Component() string {
	return p.component
}

// Container returns the base name of the container directory for this path.
func (p *Path) Container() string {
	if p.posContainerLow == -1 {
		return ""
	}
	return p.norm(p.s[p.posContainerLow : p.posContainerHigh-1])
}

// ContainerDir returns the container directory for this path.
// For content bundles this will be the parent directory.
func (p *Path) ContainerDir() string {
	if p.posContainerLow == -1 || !p.IsBundle() {
		return p.Dir()
	}
	return p.norm(p.s[:p.posContainerLow-1])
}

// Section returns the first path element (section).
func (p *Path) Section() string {
	if p.posSectionHigh <= 0 {
		return ""
	}
	return p.norm(p.s[1:p.posSectionHigh])
}

// IsContent returns true if the path is a content file (e.g. mypost.md).
// Note that this will also return true for content files in a bundle.
func (p *Path) IsContent() bool {
	return p.BundleType() >= PathTypeContentResource
}

// isContentPage returns true if the path is a content file (e.g. mypost.md),
// but nof if inside a leaf bundle.
func (p *Path) isContentPage() bool {
	return p.BundleType() >= PathTypeContentSingle
}

// Name returns the last element of path.
func (p *Path) Name() string {
	if p.posContainerHigh > 0 {
		return p.s[p.posContainerHigh:]
	}
	return p.s
}

// Name returns the last element of path without any extension.
func (p *Path) NameNoExt() string {
	if i := p.identifierIndex(0); i != -1 {
		return p.s[p.posContainerHigh : p.identifiers[i].Low-1]
	}
	return p.s[p.posContainerHigh:]
}

// Name returns the last element of path without any language identifier.
func (p *Path) NameNoLang() string {
	i := p.identifierIndex(p.posIdentifierLanguage)
	if i == -1 {
		return p.Name()
	}

	return p.s[p.posContainerHigh:p.identifiers[i].Low-1] + p.s[p.identifiers[i].High:]
}

// BaseNameNoIdentifier returns the logical base name for a resource without any identifier (e.g. no extension).
// For bundles this will be the containing directory's name, e.g. "blog".
func (p *Path) BaseNameNoIdentifier() string {
	if p.IsBundle() {
		return p.Container()
	}
	return p.NameNoIdentifier()
}

// NameNoIdentifier returns the last element of path without any identifier (e.g. no extension).
func (p *Path) NameNoIdentifier() string {
	if len(p.identifiers) > 0 {
		return p.s[p.posContainerHigh : p.identifiers[len(p.identifiers)-1].Low-1]
	}
	return p.s[p.posContainerHigh:]
}

// Dir returns all but the last element of path, typically the path's directory.
func (p *Path) Dir() (d string) {
	if p.posContainerHigh > 0 {
		d = p.s[:p.posContainerHigh-1]
	}
	if d == "" {
		d = "/"
	}
	d = p.norm(d)
	return
}

// Path returns the full path.
func (p *Path) Path() (d string) {
	return p.norm(p.s)
}

// Unnormalized returns the Path with the original case preserved.
func (p *Path) Unnormalized() *Path {
	return p.unnormalized
}

// PathNoLang returns the Path but with any language identifier removed.
func (p *Path) PathNoLang() string {
	return p.base(true, false)
}

// PathNoIdentifier returns the Path but with any identifier (ext, lang) removed.
func (p *Path) PathNoIdentifier() string {
	return p.base(false, false)
}

// PathRel returns the path relative to the given owner.
func (p *Path) PathRel(owner *Path) string {
	ob := owner.Base()
	if !strings.HasSuffix(ob, "/") {
		ob += "/"
	}
	return strings.TrimPrefix(p.Path(), ob)
}

// BaseRel returns the base path relative to the given owner.
func (p *Path) BaseRel(owner *Path) string {
	ob := owner.Base()
	if ob == "/" {
		ob = ""
	}
	return p.Base()[len(ob)+1:]
}

// For content files, Base returns the path without any identifiers (extension, language code etc.).
// Any 'index' as the last path element is ignored.
//
// For other files (Resources), any extension is kept.
func (p *Path) Base() string {
	return p.base(!p.isContentPage(), p.IsBundle())
}

// BaseNoLeadingSlash returns the base path without the leading slash.
func (p *Path) BaseNoLeadingSlash() string {
	return p.Base()[1:]
}

func (p *Path) base(preserveExt, isBundle bool) string {
	if len(p.identifiers) == 0 {
		return p.norm(p.s)
	}

	if preserveExt && len(p.identifiers) == 1 {
		// Preserve extension.
		return p.norm(p.s)
	}

	id := p.identifiers[len(p.identifiers)-1]
	high := id.Low - 1

	if isBundle {
		high = p.posContainerHigh - 1
	}

	if high == 0 {
		high++
	}

	if !preserveExt {
		return p.norm(p.s[:high])
	}

	// For txt files etc. we want to preserve the extension.
	id = p.identifiers[0]

	return p.norm(p.s[:high] + p.s[id.Low-1:id.High])
}

func (p *Path) Ext() string {
	return p.identifierAsString(0)
}

func (p *Path) Lang() string {
	return p.identifierAsString(1)
}

func (p *Path) Identifier(i int) string {
	return p.identifierAsString(i)
}

func (p *Path) Disabled() bool {
	return p.disabled
}

func (p *Path) Identifiers() []string {
	ids := make([]string, len(p.identifiers))
	for i, id := range p.identifiers {
		ids[i] = p.s[id.Low:id.High]
	}
	return ids
}

func (p *Path) BundleType() PathType {
	return p.bundleType
}

func (p *Path) IsBundle() bool {
	return p.bundleType >= PathTypeLeaf
}

func (p *Path) IsBranchBundle() bool {
	return p.bundleType == PathTypeBranch
}

func (p *Path) IsLeafBundle() bool {
	return p.bundleType == PathTypeLeaf
}

func (p *Path) IsContentData() bool {
	return p.bundleType == PathTypeContentData
}

func (p Path) ForBundleType(t PathType) *Path {
	p.bundleType = t
	return &p
}

func (p *Path) identifierAsString(i int) string {
	i = p.identifierIndex(i)
	if i == -1 {
		return ""
	}

	id := p.identifiers[i]
	return p.s[id.Low:id.High]
}

func (p *Path) identifierIndex(i int) int {
	if i < 0 || i >= len(p.identifiers) {
		return -1
	}
	return i
}

// HasExt returns true if the Unix styled path has an extension.
func HasExt(p string) bool {
	for i := len(p) - 1; i >= 0; i-- {
		if p[i] == '.' {
			return true
		}
		if p[i] == '/' {
			return false
		}
	}
	return false
}
