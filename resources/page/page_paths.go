// Copyright 2019 The Hugo Authors. All rights reserved.
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

package page

import (
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gohugoio/hugo/common/paths"

	"github.com/gohugoio/hugo/resources/page/pagekinds"

	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/output"
)

func CreateTargetPaths(d TargetPathDescriptor) (tp TargetPaths) {

	// Normalize all file Windows paths to simplify what's next.
	if helpers.FilePathSeparator != "/" {
		d.Dir = filepath.ToSlash(d.Dir)
		d.PrefixFilePath = filepath.ToSlash(d.PrefixFilePath)
	}

	if !d.Type.Root && d.URL != "" && !strings.HasPrefix(d.URL, "/") {
		// Treat this as a context relative URL
		d.ForcePrefix = true
	}

	if d.URL != "" {
		d.URL = filepath.ToSlash(d.URL)
		if strings.Contains(d.URL, "..") {
			d.URL = path.Join("/", d.URL)
		}
	}

	if d.Type.Root && !d.ForcePrefix {
		d.PrefixFilePath = ""
		d.PrefixLink = ""
	}

	pb := getPagePathBuilder(d)
	defer putPagePathBuilder(pb)

	pb.fullSuffix = d.Type.MediaType.FirstSuffix.FullSuffix

	// The top level index files, i.e. the home page etc., needs
	// the index base even when uglyURLs is enabled.
	needsBase := true

	pb.isUgly = (d.UglyURLs || d.Type.Ugly) && !d.Type.NoUgly
	pb.baseNameSameAsType = d.BaseName != "" && d.BaseName == d.Type.BaseName

	if d.ExpandedPermalink == "" && pb.baseNameSameAsType {
		pb.isUgly = true
	}

	if d.Type == output.RobotsTxtFormat {
		pb.Add(d.Type.BaseName)
		pb.noSubResources = true
	} else if d.Type == output.HTTPStatusHTMLFormat || d.Type == output.SitemapFormat {
		pb.Add(d.Kind)
		pb.noSubResources = true
	} else if d.Kind != pagekinds.Page && d.URL == "" && len(d.Sections) > 0 {
		if d.ExpandedPermalink != "" {
			pb.Add(d.ExpandedPermalink)
		} else {
			pb.Add(d.Sections...)
		}
		needsBase = false
	}

	if d.Type.Path != "" {
		pb.Add(d.Type.Path)
	}

	if d.Kind != pagekinds.Home && d.URL != "" {
		pb.Add(paths.FieldsSlash(d.URL)...)

		if d.Addends != "" {
			pb.Add(d.Addends)
		}

		hasDot := strings.Contains(d.URL, ".")
		hasSlash := strings.HasSuffix(d.URL, "/")

		if hasSlash || !hasDot {
			pb.Add(d.Type.BaseName + pb.fullSuffix)
		} else if hasDot {
			pb.fullSuffix = paths.Ext(d.URL)
		}

		if pb.IsHtmlIndex() {
			pb.linkUpperOffset = 1
		}

		if d.ForcePrefix {

			// Prepend language prefix if not already set in URL
			if d.PrefixFilePath != "" && !strings.HasPrefix(d.URL, "/"+d.PrefixFilePath) {
				pb.prefixPath = d.PrefixFilePath
			}

			if d.PrefixLink != "" && !strings.HasPrefix(d.URL, "/"+d.PrefixLink) {
				pb.prefixLink = d.PrefixLink
			}
		}
	} else if !pagekinds.IsBranch(d.Kind) {
		if d.ExpandedPermalink != "" {
			pb.Add(d.ExpandedPermalink)
		} else {
			if d.Dir != "" {
				pb.Add(d.Dir)
			}
			if d.BaseName != "" {
				pb.Add(d.BaseName)
			}
		}

		if d.Addends != "" {
			pb.Add(d.Addends)
		}

		if pb.isUgly {
			pb.ConcatLast(pb.fullSuffix)
		} else {
			pb.Add(d.Type.BaseName + pb.fullSuffix)
		}

		if pb.IsHtmlIndex() {
			pb.linkUpperOffset = 1
		}

		if d.PrefixFilePath != "" {
			pb.prefixPath = d.PrefixFilePath
		}

		if d.PrefixLink != "" {
			pb.prefixLink = d.PrefixLink
		}
	} else {
		if d.Addends != "" {
			pb.Add(d.Addends)
		}

		needsBase = needsBase && d.Addends == ""

		if needsBase || !pb.isUgly {
			pb.Add(d.Type.BaseName + pb.fullSuffix)
		} else {
			pb.ConcatLast(pb.fullSuffix)
		}

		if pb.IsHtmlIndex() {
			pb.linkUpperOffset = 1
		}

		if d.PrefixFilePath != "" {
			pb.prefixPath = d.PrefixFilePath
		}

		if d.PrefixLink != "" {
			pb.prefixLink = d.PrefixLink
		}
	}

	// if page URL is explicitly set in frontmatter,
	// preserve its value without sanitization
	if d.Kind != pagekinds.Page || d.URL == "" {
		// Note: MakePathSanitized will lower case the path if
		// disablePathToLower isn't set.
		pb.Sanitize()
	}

	link := pb.Link()
	pagePath := pb.PathFile()

	tp.TargetFilename = filepath.FromSlash(pagePath)
	if !pb.noSubResources {
		tp.SubResourceBaseTarget = filepath.FromSlash(pb.PathDir())
		tp.SubResourceBaseLink = pb.LinkDir()
	}
	if d.URL != "" {
		tp.Link = paths.URLEscape(link)
	} else {
		// This is slightly faster for when we know we don't have any
		// query or scheme etc.
		tp.Link = paths.PathEscape(link)
	}
	if tp.Link == "" {
		tp.Link = "/"
	}

	return
}

// TargetPathDescriptor describes how a file path for a given resource
// should look like on the file system. The same descriptor is then later used to
// create both the permalinks and the relative links, paginator URLs etc.
//
// The big motivating behind this is to have only one source of truth for URLs,
// and by that also get rid of most of the fragile string parsing/encoding etc.
//
//
type TargetPathDescriptor struct {
	PathSpec *helpers.PathSpec

	Type output.Format
	Kind string

	Sections []string

	// For regular content pages this is either
	// 1) the Slug, if set,
	// 2) the file base name (TranslationBaseName).
	BaseName string

	// Source directory.
	Dir string

	// Typically a language prefix added to file paths.
	PrefixFilePath string

	// Typically a language prefix added to links.
	PrefixLink string

	// If in multihost mode etc., every link/path needs to be prefixed, even
	// if set in URL.
	ForcePrefix bool

	// URL from front matter if set. Will override any Slug etc.
	URL string

	// Used to create paginator links.
	Addends string

	// The expanded permalink if defined for the section, ready to use.
	ExpandedPermalink string

	// Some types cannot have uglyURLs, even if globally enabled, RSS being one example.
	UglyURLs bool
}

// TODO(bep) move this type.
type TargetPaths struct {

	// Where to store the file on disk relative to the publish dir. OS "/"es.
	TargetFilename string

	// The directory to write sub-resources of the above.
	SubResourceBaseTarget string

	// The base for creating links to sub-resources of the above.
	SubResourceBaseLink string

	// The relative permalink to this resources. Unix "/"es.
	Link string
}

func (p TargetPaths) PermalinkForOutputFormat(s *helpers.PathSpec, f output.Format) string {
	var baseURL string
	var err error
	if f.Protocol != "" {
		baseURL, err = s.BaseURL.WithProtocol(f.Protocol)
		if err != nil {
			return ""
		}
	} else {
		baseURL = s.BaseURL.String()
	}

	return s.PermalinkForBaseURL(p.Link, baseURL)
}

func (p TargetPaths) RelPermalink(s *helpers.PathSpec) string {
	return s.PrependBasePath(p.Link, false)
}

var pagePathBuilderPool = &sync.Pool{
	New: func() interface{} {
		return &pagePathBuilder{}
	},
}

// When adding state here, remember to update putPagePathBuilder.
type pagePathBuilder struct {
	els []string

	d TargetPathDescriptor

	// Builder state.
	isUgly             bool
	baseNameSameAsType bool
	noSubResources     bool
	fullSuffix         string // File suffix including any ".".
	prefixLink         string
	prefixPath         string
	linkUpperOffset    int
}

func (p *pagePathBuilder) Add(el ...string) {
	p.els = append(p.els, el...)
}

func (p *pagePathBuilder) ConcatLast(s string) {
	if p.els == nil {
		p.Add(s)
		return
	}
	old := p.els[len(p.els)-1]
	if old[len(old)-1] == '/' {
		old = old[:len(old)-1]
	}
	p.els[len(p.els)-1] = old + s
}

func (p *pagePathBuilder) IsHtmlIndex() bool {
	return p.Last() == "index.html"
}

func (p *pagePathBuilder) Last() string {
	if p.els == nil {
		return ""
	}
	return p.els[len(p.els)-1]
}

func (p *pagePathBuilder) Link() string {
	link := p.Path(p.linkUpperOffset)

	if p.baseNameSameAsType {
		link = strings.TrimSuffix(link, p.d.BaseName)
	}

	if p.prefixLink != "" {
		link = "/" + p.prefixLink + link
	}

	if p.linkUpperOffset > 0 && !strings.HasSuffix(link, "/") {
		link += "/"
	}

	return link
}

func (p *pagePathBuilder) LinkDir() string {
	if p.noSubResources {
		return ""
	}

	pathDir := p.PathDirBase()

	if p.prefixLink != "" {
		pathDir = "/" + p.prefixLink + pathDir
	}

	return pathDir
}

func (p *pagePathBuilder) Path(upperOffset int) string {
	upper := len(p.els)
	if upperOffset > 0 {
		upper -= upperOffset
	}
	pth := path.Join(p.els[:upper]...)
	return helpers.AddLeadingSlash(pth)
}

func (p *pagePathBuilder) PathDir() string {
	dir := p.PathDirBase()
	if p.prefixPath != "" {
		dir = "/" + p.prefixPath + dir
	}
	return dir
}

func (p *pagePathBuilder) PathDirBase() string {
	if p.noSubResources {
		return ""
	}

	dir := p.Path(0)
	isIndex := strings.HasPrefix(p.Last(), p.d.Type.BaseName+".")

	if isIndex {
		dir = path.Dir(dir)
	} else {
		dir = strings.TrimSuffix(dir, p.fullSuffix)
	}

	if dir == "/" {
		dir = ""
	}

	return dir
}

func (p *pagePathBuilder) PathFile() string {
	dir := p.Path(0)
	if p.prefixPath != "" {
		dir = "/" + p.prefixPath + dir
	}
	return dir
}

func (p *pagePathBuilder) Prepend(el ...string) {
	p.els = append(p.els[:0], append(el, p.els[0:]...)...)
}

func (p *pagePathBuilder) Sanitize() {
	for i, el := range p.els {
		p.els[i] = p.d.PathSpec.MakePathSanitized(el)
	}
}

func getPagePathBuilder(d TargetPathDescriptor) *pagePathBuilder {
	b := pagePathBuilderPool.Get().(*pagePathBuilder)
	b.d = d
	return b
}

func putPagePathBuilder(b *pagePathBuilder) {
	b.els = b.els[:0]
	b.fullSuffix = ""
	b.baseNameSameAsType = false
	b.isUgly = false
	b.noSubResources = false
	b.prefixLink = ""
	b.prefixPath = ""
	b.linkUpperOffset = 0
	pagePathBuilderPool.Put(b)
}
