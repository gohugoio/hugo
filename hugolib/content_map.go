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

package hugolib

import (
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"unicode"

	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/parser/pageparser"

	"github.com/gohugoio/hugo/resources/page/pagekinds"

	"github.com/gobuffalo/flect"
	"github.com/gohugoio/hugo/output"

	"github.com/gohugoio/hugo/common/types"

	"github.com/gohugoio/hugo/helpers"

	"github.com/gohugoio/hugo/hugofs"
)

// Used to mark ambiguous keys in reverse index lookups.
var ambiguousContentNode = &contentNode{}

var (
	contentTreeNoListAlwaysFilter = func(s string, n *contentNode) bool {
		if n.p == nil {
			return true
		}
		return n.p.m.noListAlways()
	}

	contentTreeNoRenderFilter = func(s string, n *contentNode) bool {
		if n.p == nil {
			return true
		}
		return n.p.m.noRender()
	}

	contentTreeNoLinkFilter = func(s string, n *contentNode) bool {
		if n.p == nil {
			return true
		}
		return n.p.m.noLink()
	}
)

var (
	_ contentKindProvider = (*contentBundleViewInfo)(nil)
	_ viewInfoTrait       = (*contentBundleViewInfo)(nil)
)

var trimCutsetDotSlashSpace = func(r rune) bool {
	return r == '.' || r == '/' || unicode.IsSpace(r)
}

func newcontentTreeNodeCallbackChain(callbacks ...contentTreeNodeCallback) contentTreeNodeCallback {
	return func(s string, n *contentNode) bool {
		for i, cb := range callbacks {
			// Allow the last callback to stop the walking.
			if i == len(callbacks)-1 {
				return cb(s, n)
			}

			if cb(s, n) {
				// Skip the rest of the callbacks, but continue walking.
				return false
			}
		}
		return false
	}
}

type contentBundleViewInfo struct {
	name viewName
	term string
}

func (c *contentBundleViewInfo) Kind() string {
	if c.term != "" {
		return pagekinds.Term
	}
	return pagekinds.Taxonomy
}

func (c *contentBundleViewInfo) Term() string {
	return c.term
}

func (c *contentBundleViewInfo) ViewInfo() *contentBundleViewInfo {
	if c == nil {
		panic("ViewInfo() called on nil")
	}
	return c
}

type contentGetBranchProvider interface {
	// GetBranch returns the the current branch, which will be itself
	// for branch nodes (e.g. sections).
	// To always navigate upwards, use GetContainerBranch().
	GetBranch() *contentBranchNode
}

type contentGetContainerBranchProvider interface {
	// GetContainerBranch returns the container for pages and sections.
	GetContainerBranch() *contentBranchNode
}

type contentGetContainerNodeProvider interface {
	// GetContainerNode returns the container for resources.
	GetContainerNode() *contentNode
}

type contentGetNodeProvider interface {
	GetNode() *contentNode
}

type contentKindProvider interface {
	Kind() string
}

type contentMapConfig struct {
	lang                 string
	taxonomyConfig       taxonomiesConfigValues
	taxonomyDisabled     bool
	taxonomyTermDisabled bool
	pageDisabled         bool
	isRebuild            bool
}

func (cfg contentMapConfig) getTaxonomyConfig(s string) (v viewName) {
	s = strings.TrimPrefix(s, "/")
	if s == "" {
		return
	}
	for _, n := range cfg.taxonomyConfig.views {
		if strings.HasPrefix(s, n.plural) {
			return n
		}
	}

	return
}

var (
	_ identity.IdentityProvider          = (*contentNode)(nil)
	_ identity.DependencyManagerProvider = (*contentNode)(nil)
)

type contentNode struct {
	key string

	keyPartsInit sync.Once
	keyParts     []string

	p *pageState

	running bool

	// Additional traits for this node.
	traits interface{}

	// Tracks dependencies in server mode.
	idmInit sync.Once
	idm     identity.Manager
}

type contentNodeIdentity struct {
	n *contentNode
}

func (n *contentNode) IdentifierBase() interface{} {
	return n.key
}

func (b *contentNode) GetIdentity() identity.Identity {
	return b
}

func (b *contentNode) GetDependencyManager() identity.Manager {
	b.idmInit.Do(func() {
		// TODO1
		if true || b.running {
			b.idm = identity.NewManager(b)
		} else {
			b.idm = identity.NopManager
		}
	})
	return b.idm
}

func (b *contentNode) GetContainerNode() *contentNode {
	return b
}

func (b *contentNode) HasFi() bool {
	_, ok := b.traits.(hugofs.FileInfoProvider)
	return ok
}

func (b *contentNode) FileInfo() hugofs.FileMetaInfo {
	fip, ok := b.traits.(hugofs.FileInfoProvider)
	if !ok {
		return nil
	}
	return fip.FileInfo()
}

func (b *contentNode) Key() string {
	return b.key
}

func (b *contentNode) KeyParts() []string {
	b.keyPartsInit.Do(func() {
		if b.key != "" {
			b.keyParts = paths.FieldsSlash(b.key)
		}
	})
	return b.keyParts
}

func (b *contentNode) GetNode() *contentNode {
	return b
}

func (b *contentNode) IsStandalone() bool {
	_, ok := b.traits.(kindOutputFormat)
	return ok
}

// IsView returns whether this is a view node (a taxonomy or a term).
func (b *contentNode) IsView() bool {
	_, ok := b.traits.(viewInfoTrait)
	return ok
}

// isCascadingEdit parses any front matter and returns whether it has a cascade section and
// if that has changed.
func (n *contentNode) isCascadingEdit() bool {
	if n.p == nil {
		return false
	}
	fi := n.FileInfo()
	if fi == nil {
		return false
	}
	f, err := fi.Meta().Open()
	if err != nil {
		// File may have been removed, assume a cascading edit.
		// Some false positives are OK.
		return true
	}

	pf, err := pageparser.ParseFrontMatterAndContent(f)
	f.Close()
	if err != nil {
		return true
	}

	if n.p == nil || n.p.bucket == nil {
		return false
	}

	maps.PrepareParams(pf.FrontMatter)
	cascade1, ok := pf.FrontMatter["cascade"]
	hasCascade := n.p.bucket.cascade != nil && len(n.p.bucket.cascade) > 0
	if !ok {
		return hasCascade
	}

	if !hasCascade {
		return true
	}

	for _, v := range n.p.bucket.cascade {
		if !reflect.DeepEqual(cascade1, v) {
			return true
		}
	}

	return false
}

type contentNodeInfo struct {
	branch     *contentBranchNode
	isBranch   bool
	isResource bool
}

func (info *contentNodeInfo) SectionsEntries() []string {
	return info.branch.n.KeyParts()
}

// TDOO1 somehow document that this will now return a leading slash, "" for home page.
func (info *contentNodeInfo) SectionsPath() string {
	k := info.branch.n.Key()
	if k == "" {
		// TODO1 consider this.
		return "/"
	}
	return k
}

type contentNodeInfoProvider interface {
	SectionsEntries() []string
	SectionsPath() string
}

type contentNodeProvider interface {
	contentGetNodeProvider
	types.Identifier
}

type contentTreeNodeCallback func(s string, n *contentNode) bool

type contentTreeNodeCallbackNew func(node contentNodeProvider) bool

type contentTreeRefProvider interface {
	contentGetBranchProvider
	contentGetContainerNodeProvider
	contentNodeInfoProvider
	contentNodeProvider
}

type fileInfoHolder struct {
	fi hugofs.FileMetaInfo
}

func (f fileInfoHolder) FileInfo() hugofs.FileMetaInfo {
	return f.fi
}

type kindOutputFormat struct {
	kind   string
	output output.Format
}

func (k kindOutputFormat) Kind() string {
	return k.kind
}

func (k kindOutputFormat) OutputFormat() output.Format {
	return k.output
}

type kindOutputFormatTrait interface {
	Kind() string
	OutputFormat() output.Format
}

// bookmark3
func (m *pageMap) AddFilesBundle(header hugofs.FileMetaInfo, resources ...hugofs.FileMetaInfo) error {
	var (
		n        *contentNode
		pageTree *contentBranchNode
		pathInfo = header.Meta().PathInfo
	)

	if !pathInfo.IsBranchBundle() && m.cfg.pageDisabled {
		return nil
	}

	if pathInfo.IsBranchBundle() {
		// Apply some metadata if it's a taxonomy node.
		if tc := m.cfg.getTaxonomyConfig(pathInfo.Base()); !tc.IsZero() {
			term := strings.TrimPrefix(strings.TrimPrefix(pathInfo.Base(), "/"+tc.plural), "/")

			n = m.NewContentNode(
				viewInfoFileInfoHolder{
					&contentBundleViewInfo{
						name: tc,
						term: term,
					},
					fileInfoHolder{fi: header},
				},
				pathInfo.Base(),
			)

		} else {
			n = m.NewContentNode(
				fileInfoHolder{fi: header},
				pathInfo.Base(),
			)
		}

		pageTree = m.InsertBranch(n)

		for _, r := range resources {
			n := m.NewContentNode(
				fileInfoHolder{fi: r},
				r.Meta().Path,
			)
			pageTree.resources.nodes.Insert(n.key, n)
		}

		return nil
	}

	n = m.NewContentNode(
		fileInfoHolder{fi: header},
		pathInfo.Base(),
	)

	// A regular page. Attach it to its section.
	var created bool
	_, pageTree, created = m.getOrCreateSection(n)

	if created {
		// This means there are most likely no content file for this
		// section.
		// Apply some default metadata to the node.
		sectionName := helpers.FirstUpper(m.rootSection(pathInfo))
		var title string
		if m.s.Cfg.GetBool("pluralizeListTitles") {
			title = flect.Pluralize(sectionName)
		} else {
			title = sectionName
		}
		pageTree.defaultTitle = title
	}

	pageTree.InsertPage(n.key, n)

	for _, r := range resources {
		n := m.NewContentNode(
			fileInfoHolder{fi: r},
			r.Meta().Path,
		)
		pageTree.pageResources.nodes.Insert(n.key, n)
	}

	return nil
}

func (m *pageMap) getOrCreateSection(n *contentNode) (string, *contentBranchNode, bool) {
	level := strings.Count(n.key, "/")

	k, pageTree := m.LongestPrefix(path.Dir(n.key))

	mustCreate := false

	if pageTree == nil {
		mustCreate = true
	} else if level > 1 && k == "" {
		// We found the home section, but this page needs to be placed in
		// the root, e.g. "/blog", section.
		mustCreate = true
	} else {
		return k, pageTree, false
	}

	if !mustCreate {
		return k, pageTree, false
	}

	var keyParts []string
	if level > 1 {
		keyParts = n.KeyParts()[:1]
	}
	n = m.NewContentNode(nil, keyParts...)

	if k != "" {
		// Make sure we always have the root/home node.
		if m.Get("") == nil {
			m.InsertBranch(&contentNode{})
		}
	}

	pageTree = m.InsertBranch(n)
	return k, pageTree, true
}

func (m *pageMap) rootSection(p paths.Path) string {
	firstSlash := strings.Index(p.Dir(), "/")
	if firstSlash == -1 {
		return p.Dir()
	}
	return p.Dir()[:firstSlash]
}

type stringKindProvider string

func (k stringKindProvider) Kind() string {
	return string(k)
}

type viewInfoFileInfoHolder struct {
	viewInfoTrait
	hugofs.FileInfoProvider
}

type viewInfoTrait interface {
	Kind() string
	ViewInfo() *contentBundleViewInfo
}

// The home page is represented with the zero string.
// All other keys starts with a leading slash. No trailing slash.
// Slashes are Unix-style.
func cleanTreeKey(elem ...string) string {
	var s string
	if len(elem) > 0 {
		s = elem[0]
		if len(elem) > 1 {
			s = path.Join(elem...)
		}
	}
	s = strings.TrimFunc(s, trimCutsetDotSlashSpace)
	s = filepath.ToSlash(strings.ToLower(paths.Sanitize(s)))
	if s == "" || s == "/" {
		return ""
	}
	if s[0] != '/' {
		s = "/" + s
	}
	return s
}
