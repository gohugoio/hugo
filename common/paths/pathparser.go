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
	"github.com/gohugoio/hugo/resources/kinds"
)

const (
	identifierBaseof = "baseof"
)

// PathParser parses a path into a Path.
type PathParser struct {
	// Maps the language code to its index in the languages/sites slice.
	LanguageIndex map[string]int

	// Reports whether the given language is disabled.
	IsLangDisabled func(string) bool

	// IsOutputFormat reports whether the given name is a valid output format.
	// The second argument is optional.
	IsOutputFormat func(name, ext string) bool

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
	p := &Path{}
	p.reset()
	p.component = component
	return p
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

func (pp *PathParser) parseIdentifier(component, s string, p *Path, i, lastDot, numDots int, isLast bool) {
	if p.posContainerHigh != -1 {
		return
	}
	mayHaveLang := numDots > 1 && p.posIdentifierLanguage == -1 && pp.LanguageIndex != nil
	mayHaveLang = mayHaveLang && (component == files.ComponentFolderContent || component == files.ComponentFolderLayouts)
	mayHaveOutputFormat := component == files.ComponentFolderLayouts
	mayHaveKind := p.posIdentifierKind == -1 && mayHaveOutputFormat
	var mayHaveLayout bool
	if p.pathType == TypeShortcode {
		mayHaveLayout = !isLast && component == files.ComponentFolderLayouts
	} else {
		mayHaveLayout = component == files.ComponentFolderLayouts
	}

	var found bool
	var high int
	if len(p.identifiersKnown) > 0 {
		high = lastDot
	} else {
		high = len(p.s)
	}
	id := types.LowHigh[string]{Low: i + 1, High: high}
	sid := p.s[id.Low:id.High]

	if len(p.identifiersKnown) == 0 {
		// The first is always the extension.
		p.identifiersKnown = append(p.identifiersKnown, id)
		found = true

		// May also be the output format.
		if mayHaveOutputFormat && pp.IsOutputFormat(sid, "") {
			p.posIdentifierOutputFormat = 0
		}
	} else {

		var langFound bool

		if mayHaveLang {
			var disabled bool
			_, langFound = pp.LanguageIndex[sid]
			if !langFound {
				disabled = pp.IsLangDisabled != nil && pp.IsLangDisabled(sid)
				if disabled {
					p.disabled = true
					langFound = true
				}
			}
			found = langFound
			if langFound {
				p.identifiersKnown = append(p.identifiersKnown, id)
				p.posIdentifierLanguage = len(p.identifiersKnown) - 1
			}
		}

		if !found && mayHaveOutputFormat {
			// At this point we may already have resolved an output format,
			// but we need to keep looking for a more specific one, e.g. amp before html.
			// Use both name and extension to prevent
			// false positives on the form css.html.
			if pp.IsOutputFormat(sid, p.Ext()) {
				found = true
				p.identifiersKnown = append(p.identifiersKnown, id)
				p.posIdentifierOutputFormat = len(p.identifiersKnown) - 1
			}
		}

		if !found && mayHaveKind {
			if kinds.GetKindMain(sid) != "" {
				found = true
				p.identifiersKnown = append(p.identifiersKnown, id)
				p.posIdentifierKind = len(p.identifiersKnown) - 1
			}
		}

		if !found && sid == identifierBaseof {
			found = true
			p.identifiersKnown = append(p.identifiersKnown, id)
			p.posIdentifierBaseof = len(p.identifiersKnown) - 1
		}

		if !found && mayHaveLayout {
			p.identifiersKnown = append(p.identifiersKnown, id)
			p.posIdentifierLayout = len(p.identifiersKnown) - 1
			found = true
		}

		if !found {
			p.identifiersUnknown = append(p.identifiersUnknown, id)
		}

	}
}

func (pp *PathParser) doParse(component, s string, p *Path) (*Path, error) {
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
	lastDot := 0
	lastSlashIdx := strings.LastIndex(s, "/")
	numDots := strings.Count(s[lastSlashIdx+1:], ".")
	if strings.Contains(s, "/_shortcodes/") {
		p.pathType = TypeShortcode
	}

	for i := len(s) - 1; i >= 0; i-- {
		c := s[i]

		switch c {
		case '.':
			pp.parseIdentifier(component, s, p, i, lastDot, numDots, false)
			lastDot = i
		case '/':
			slashCount++
			if p.posContainerHigh == -1 {
				if lastDot > 0 {
					pp.parseIdentifier(component, s, p, i, lastDot, numDots, true)
				}
				p.posContainerHigh = i + 1
			} else if p.posContainerLow == -1 {
				p.posContainerLow = i + 1
			}
			if i > 0 {
				p.posSectionHigh = i
			}
		}
	}

	if len(p.identifiersKnown) > 0 {
		isContentComponent := p.component == files.ComponentFolderContent || p.component == files.ComponentFolderArchetypes
		isContent := isContentComponent && pp.IsContentExt(p.Ext())
		id := p.identifiersKnown[len(p.identifiersKnown)-1]

		if id.Low > p.posContainerHigh {
			b := p.s[p.posContainerHigh : id.Low-1]
			if isContent {
				switch b {
				case "index":
					p.pathType = TypeLeaf
				case "_index":
					p.pathType = TypeBranch
				default:
					p.pathType = TypeContentSingle
				}

				if slashCount == 2 && p.IsLeafBundle() {
					p.posSectionHigh = 0
				}
			} else if b == files.NameContentData && files.IsContentDataExt(p.Ext()) {
				p.pathType = TypeContentData
			}
		}
	}

	if p.pathType < TypeMarkup && component == files.ComponentFolderLayouts {
		if p.posIdentifierBaseof != -1 {
			p.pathType = TypeBaseof
		} else {
			pth := p.Path()
			if strings.Contains(pth, "/_shortcodes/") {
				p.pathType = TypeShortcode
			} else if strings.Contains(pth, "/_markup/") {
				p.pathType = TypeMarkup
			} else if strings.HasPrefix(pth, "/_partials/") {
				p.pathType = TypePartial
			}
		}
	}

	if p.pathType == TypeShortcode && p.posIdentifierLayout != -1 {
		id := p.identifiersKnown[p.posIdentifierLayout]
		if id.Low == p.posContainerHigh {
			// First identifier is shortcode name.
			p.posIdentifierLayout = -1
		}
	}

	return p, nil
}

func ModifyPathBundleTypeResource(p *Path) {
	if p.IsContent() {
		p.pathType = TypeContentResource
	} else {
		p.pathType = TypeFile
	}
}

//go:generate stringer -type Type

type Type int

const (

	// A generic resource, e.g. a JSON file.
	TypeFile Type = iota

	// All below are content files.
	// A resource of a content type with front matter.
	TypeContentResource

	// E.g. /blog/my-post.md
	TypeContentSingle

	// All below are bundled content files.

	// Leaf bundles, e.g. /blog/my-post/index.md
	TypeLeaf

	// Branch bundles, e.g. /blog/_index.md
	TypeBranch

	// Content data file, _content.gotmpl.
	TypeContentData

	// Layout types.
	TypeMarkup
	TypeShortcode
	TypePartial
	TypeBaseof
)

type Path struct {
	// Note: Any additions to this struct should also be added to the pathPool.
	s string

	posContainerLow  int
	posContainerHigh int
	posSectionHigh   int

	component string
	pathType  Type

	identifiersKnown   []types.LowHigh[string]
	identifiersUnknown []types.LowHigh[string]

	posIdentifierLanguage     int
	posIdentifierOutputFormat int
	posIdentifierKind         int
	posIdentifierLayout       int
	posIdentifierBaseof       int
	disabled                  bool

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
	p.pathType = 0
	p.identifiersKnown = p.identifiersKnown[:0]
	p.posIdentifierLanguage = -1
	p.posIdentifierOutputFormat = -1
	p.posIdentifierKind = -1
	p.posIdentifierLayout = -1
	p.posIdentifierBaseof = -1
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
	if p.Component() == files.ComponentFolderLayouts {
		return p.Path()
	}
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

func (p *Path) String() string {
	if p == nil {
		return "<nil>"
	}
	return p.Path()
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
	return p.Type() >= TypeContentResource && p.Type() <= TypeContentData
}

// isContentPage returns true if the path is a content file (e.g. mypost.md),
// but nof if inside a leaf bundle.
func (p *Path) isContentPage() bool {
	return p.Type() >= TypeContentSingle && p.Type() <= TypeContentData
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
		return p.s[p.posContainerHigh : p.identifiersKnown[i].Low-1]
	}
	return p.s[p.posContainerHigh:]
}

// Name returns the last element of path without any language identifier.
func (p *Path) NameNoLang() string {
	i := p.identifierIndex(p.posIdentifierLanguage)
	if i == -1 {
		return p.Name()
	}

	return p.s[p.posContainerHigh:p.identifiersKnown[i].Low-1] + p.s[p.identifiersKnown[i].High:]
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
	lowHigh := p.nameLowHigh()
	return p.s[lowHigh.Low:lowHigh.High]
}

func (p *Path) nameLowHigh() types.LowHigh[string] {
	if len(p.identifiersKnown) > 0 {
		lastID := p.identifiersKnown[len(p.identifiersKnown)-1]
		if p.posContainerHigh == lastID.Low {
			// The last identifier is the name.
			return lastID
		}
		return types.LowHigh[string]{
			Low:  p.posContainerHigh,
			High: p.identifiersKnown[len(p.identifiersKnown)-1].Low - 1,
		}
	}
	return types.LowHigh[string]{
		Low:  p.posContainerHigh,
		High: len(p.s),
	}
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

// PathNoLeadingSlash returns the full path without the leading slash.
func (p *Path) PathNoLeadingSlash() string {
	return p.Path()[1:]
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

// PathBeforeLangAndOutputFormatAndExt returns the path up to the first identifier that is not a language or output format.
func (p *Path) PathBeforeLangAndOutputFormatAndExt() string {
	if len(p.identifiersKnown) == 0 {
		return p.norm(p.s)
	}
	i := p.identifierIndex(0)

	if j := p.posIdentifierOutputFormat; i == -1 || (j != -1 && j < i) {
		i = j
	}
	if j := p.posIdentifierLanguage; i == -1 || (j != -1 && j < i) {
		i = j
	}

	if i == -1 {
		return p.norm(p.s)
	}

	id := p.identifiersKnown[i]
	return p.norm(p.s[:id.Low-1])
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

// Used in template lookups.
// For pages with Type set, we treat that as the section.
func (p *Path) BaseReTyped(typ string) (d string) {
	base := p.Base()
	if typ == "" || p.Section() == typ {
		return base
	}
	d = "/" + typ
	if p.posSectionHigh != -1 {
		d += base[p.posSectionHigh:]
	}
	d = p.norm(d)
	return
}

// BaseNoLeadingSlash returns the base path without the leading slash.
func (p *Path) BaseNoLeadingSlash() string {
	return p.Base()[1:]
}

func (p *Path) base(preserveExt, isBundle bool) string {
	if len(p.identifiersKnown) == 0 {
		return p.norm(p.s)
	}

	if preserveExt && len(p.identifiersKnown) == 1 {
		// Preserve extension.
		return p.norm(p.s)
	}

	var high int

	if isBundle {
		high = p.posContainerHigh - 1
	} else {
		high = p.nameLowHigh().High
	}

	if high == 0 {
		high++
	}

	if !preserveExt {
		return p.norm(p.s[:high])
	}

	// For txt files etc. we want to preserve the extension.
	id := p.identifiersKnown[0]

	return p.norm(p.s[:high] + p.s[id.Low-1:id.High])
}

func (p *Path) Ext() string {
	return p.identifierAsString(0)
}

func (p *Path) OutputFormat() string {
	return p.identifierAsString(p.posIdentifierOutputFormat)
}

func (p *Path) Kind() string {
	return p.identifierAsString(p.posIdentifierKind)
}

func (p *Path) Layout() string {
	return p.identifierAsString(p.posIdentifierLayout)
}

func (p *Path) Lang() string {
	return p.identifierAsString(p.posIdentifierLanguage)
}

func (p *Path) Identifier(i int) string {
	return p.identifierAsString(i)
}

func (p *Path) Disabled() bool {
	return p.disabled
}

func (p *Path) Identifiers() []string {
	ids := make([]string, len(p.identifiersKnown))
	for i, id := range p.identifiersKnown {
		ids[i] = p.s[id.Low:id.High]
	}
	return ids
}

func (p *Path) IdentifiersUnknown() []string {
	ids := make([]string, len(p.identifiersUnknown))
	for i, id := range p.identifiersUnknown {
		ids[i] = p.s[id.Low:id.High]
	}
	return ids
}

func (p *Path) Type() Type {
	return p.pathType
}

func (p *Path) IsBundle() bool {
	return p.pathType >= TypeLeaf && p.pathType <= TypeContentData
}

func (p *Path) IsBranchBundle() bool {
	return p.pathType == TypeBranch
}

func (p *Path) IsLeafBundle() bool {
	return p.pathType == TypeLeaf
}

func (p *Path) IsContentData() bool {
	return p.pathType == TypeContentData
}

func (p Path) ForType(t Type) *Path {
	p.pathType = t
	return &p
}

func (p *Path) identifierAsString(i int) string {
	i = p.identifierIndex(i)
	if i == -1 {
		return ""
	}

	id := p.identifiersKnown[i]
	return p.s[id.Low:id.High]
}

func (p *Path) identifierIndex(i int) int {
	if i < 0 || i >= len(p.identifiersKnown) {
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
