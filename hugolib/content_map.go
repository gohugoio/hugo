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
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gohugoio/hugo/resources/page"
	"github.com/pkg/errors"

	"github.com/gohugoio/hugo/hugofs/files"

	"github.com/gohugoio/hugo/hugofs"

	radix "github.com/armon/go-radix"
)

// We store the branch nodes in either the `sections` or `taxonomies` tree
// with their path as a key; Unix style slashes, a leading slash but no
// trailing slash.
//
// E.g. "/blog" or "/categories/funny"
//
// Pages that belongs to a section are stored in the `pages` tree below
// the section name and a branch separator, e.g. "/blog__hb_". A page is
// given a key using the path below the section and the base filename with no extension
// with a leaf separator added.
//
// For bundled pages (/mybundle/index.md), we use the folder name.
//
// An exmple of a full page key would be "/blog__hb_/page1__hl_"
//
// Bundled resources are stored in the `resources` having their path prefixed
// with the bundle they belong to, e.g.
// "/blog__hb_/bundle__hl_data.json".
//
// The weighted taxonomy entries extracted from page front matter are stored in
// the `taxonomyEntries` tree below /plural/term/page-key, e.g.
// "/categories/funny/blog__hb_/bundle__hl_".
const (
	cmBranchSeparator = "__hb_"
	cmLeafSeparator   = "__hl_"
)

// Used to mark ambigous keys in reverse index lookups.
var ambigousContentNode = &contentNode{}

func newContentMap(cfg contentMapConfig) *contentMap {
	m := &contentMap{
		cfg:             &cfg,
		pages:           &contentTree{Name: "pages", Tree: radix.New()},
		sections:        &contentTree{Name: "sections", Tree: radix.New()},
		taxonomies:      &contentTree{Name: "taxonomies", Tree: radix.New()},
		taxonomyEntries: &contentTree{Name: "taxonomyEntries", Tree: radix.New()},
		resources:       &contentTree{Name: "resources", Tree: radix.New()},
	}

	m.pageTrees = []*contentTree{
		m.pages, m.sections, m.taxonomies,
	}

	m.bundleTrees = []*contentTree{
		m.pages, m.sections, m.taxonomies, m.resources,
	}

	m.branchTrees = []*contentTree{
		m.sections, m.taxonomies,
	}

	addToReverseMap := func(k string, n *contentNode, m map[interface{}]*contentNode) {
		k = strings.ToLower(k)
		existing, found := m[k]
		if found && existing != ambigousContentNode {
			m[k] = ambigousContentNode
		} else if !found {
			m[k] = n
		}
	}

	m.pageReverseIndex = &contentTreeReverseIndex{
		t: []*contentTree{m.pages, m.sections, m.taxonomies},
		initFn: func(t *contentTree, m map[interface{}]*contentNode) {
			t.Walk(func(s string, v interface{}) bool {
				n := v.(*contentNode)
				if n.p != nil && !n.p.File().IsZero() {
					meta := n.p.File().FileInfo().Meta()
					if meta.Path() != meta.PathFile() {
						// Keep track of the original mount source.
						mountKey := filepath.ToSlash(filepath.Join(meta.Module(), meta.PathFile()))
						addToReverseMap(mountKey, n, m)
					}
				}
				k := strings.TrimSuffix(path.Base(s), cmLeafSeparator)
				addToReverseMap(k, n, m)
				return false
			})
		},
	}

	return m
}

type cmInsertKeyBuilder struct {
	m *contentMap

	err error

	// Builder state
	tree    *contentTree
	baseKey string // Section or page key
	key     string
}

func (b cmInsertKeyBuilder) ForPage(s string) *cmInsertKeyBuilder {
	// TODO2 fmt.Println("ForPage:", s, "baseKey:", b.baseKey, "key:", b.key)
	baseKey := b.baseKey
	b.baseKey = s

	if !strings.HasPrefix(s, "/") {
		s = "/" + s
	}

	if baseKey != "/" {
		// Don't repeat the section path in the key.
		s = strings.TrimPrefix(s, baseKey)
	}

	switch b.tree {
	case b.m.sections:
		b.tree = b.m.pages
		b.key = baseKey + cmBranchSeparator + s + cmLeafSeparator
	case b.m.taxonomies:
		b.key = path.Join(baseKey, s)
	default:
		panic("invalid state")
	}

	return &b
}

func (b cmInsertKeyBuilder) ForResource(s string) *cmInsertKeyBuilder {
	// TODO2 fmt.Println("ForResource:", s, "baseKey:", b.baseKey, "key:", b.key)

	s = strings.TrimPrefix(s, "/")
	s = strings.TrimPrefix(s, strings.TrimPrefix(b.baseKey, "/")+"/")

	switch b.tree {
	case b.m.pages:
		b.key = b.key + s
	case b.m.sections, b.m.taxonomies:
		b.key = b.key + cmLeafSeparator + s
	default:
		panic(fmt.Sprintf("invalid state: %#v", b.tree))
	}
	b.tree = b.m.resources
	return &b
}

func (b *cmInsertKeyBuilder) Insert(n *contentNode) *cmInsertKeyBuilder {
	if b.err == nil {
		b.tree.Insert(cleanTreeKey(b.key), n)
	}
	return b
}

func (b *cmInsertKeyBuilder) DeleteAll() *cmInsertKeyBuilder {
	if b.err == nil {
		b.tree.DeletePrefix(cleanTreeKey(b.key))
	}
	return b
}

func (b *cmInsertKeyBuilder) WithFile(fi hugofs.FileMetaInfo) *cmInsertKeyBuilder {
	b.newTopLevel()
	m := b.m
	meta := fi.Meta()
	p := cleanTreeKey(meta.Path())
	bundlePath := m.getBundleDir(meta)
	isBundle := meta.Classifier().IsBundle()
	if isBundle {
		panic("not implemented")
	}

	p, k := b.getBundle(p)
	if k == "" {
		b.err = errors.Errorf("no bundle header found for %q", bundlePath)
		return b
	}

	id := k + m.reduceKeyPart(p, fi.Meta().Path())
	b.tree = b.m.resources
	b.key = id
	b.baseKey = p

	return b
}

func (b *cmInsertKeyBuilder) WithSection(s string) *cmInsertKeyBuilder {
	b.newTopLevel()
	b.tree = b.m.sections
	b.baseKey = s
	b.key = s
	// TODO2 fmt.Println("WithSection:", s, "baseKey:", b.baseKey, "key:", b.key)
	return b
}

func (b *cmInsertKeyBuilder) WithTaxonomy(s string) *cmInsertKeyBuilder {
	b.newTopLevel()
	b.tree = b.m.taxonomies
	b.baseKey = s
	b.key = s
	return b
}

// getBundle gets both the key to the section and the prefix to where to store
// this page bundle and its resources.
func (b *cmInsertKeyBuilder) getBundle(s string) (string, string) {
	m := b.m
	section, _ := m.getSection(s)

	p := s
	if section != "/" {
		p = strings.TrimPrefix(s, section)
	}

	bundlePathParts := strings.Split(p, "/")[1:]
	basePath := section + cmBranchSeparator

	// Put it into an existing bundle if found.
	for i := len(bundlePathParts) - 2; i >= 0; i-- {
		bundlePath := path.Join(bundlePathParts[:i]...)
		searchKey := basePath + "/" + bundlePath + cmLeafSeparator
		if _, found := m.pages.Get(searchKey); found {
			return section + "/" + bundlePath, searchKey
		}
	}

	// Put it into the section bundle.
	return section, section + cmLeafSeparator
}

func (b *cmInsertKeyBuilder) newTopLevel() {
	b.key = ""
}

type contentBundleViewInfo struct {
	name       viewName
	termKey    string
	termOrigin string
	weight     int
	ref        *contentNode
}

func (c *contentBundleViewInfo) kind() string {
	if c.termKey != "" {
		return page.KindTaxonomy
	}
	return page.KindTaxonomyTerm
}

func (c *contentBundleViewInfo) sections() []string {
	if c.kind() == page.KindTaxonomyTerm {
		return []string{c.name.plural}
	}

	return []string{c.name.plural, c.termKey}

}

func (c *contentBundleViewInfo) term() string {
	if c.termOrigin != "" {
		return c.termOrigin
	}

	return c.termKey
}

type contentMap struct {
	cfg *contentMapConfig

	// View of regular pages, sections, and taxonomies.
	pageTrees contentTrees

	// View of pages, sections, taxonomies, and resources.
	bundleTrees contentTrees

	// View of sections and taxonomies.
	branchTrees contentTrees

	// Stores page bundles keyed by its path's directory or the base filename,
	// e.g. "blog/post.md" => "/blog/post", "blog/post/index.md" => "/blog/post"
	// These are the "regular pages" and all of them are bundles.
	pages *contentTree

	// A reverse index used as a fallback in GetPage.
	// There are currently two cases where this is used:
	// 1. Short name lookups in ref/relRef, e.g. using only "mypage.md" without a path.
	// 2. Links resolved from a remounted content directory. These are restricted to the same module.
	// Both of the above cases can  result in ambigous lookup errors.
	pageReverseIndex *contentTreeReverseIndex

	// Section nodes.
	sections *contentTree

	// Taxonomy nodes.
	taxonomies *contentTree

	// Pages in a taxonomy.
	taxonomyEntries *contentTree

	// Resources stored per bundle below a common prefix, e.g. "/blog/post__hb_".
	resources *contentTree
}

func (m *contentMap) AddFiles(fis ...hugofs.FileMetaInfo) error {
	for _, fi := range fis {
		if err := m.addFile(fi); err != nil {
			return err
		}
	}

	return nil
}

func (m *contentMap) AddFilesBundle(header hugofs.FileMetaInfo, resources ...hugofs.FileMetaInfo) error {
	var (
		meta       = header.Meta()
		classifier = meta.Classifier()
		isBranch   = classifier == files.ContentClassBranch
		bundlePath = m.getBundleDir(meta)

		n = m.newContentNodeFromFi(header)
		b = m.newKeyBuilder()

		section string
	)

	if isBranch {
		// Either a section or a taxonomy node.
		section = bundlePath
		if tc := m.cfg.getTaxonomyConfig(section); !tc.IsZero() {
			term := strings.TrimPrefix(strings.TrimPrefix(section, "/"+tc.plural), "/")

			n.viewInfo = &contentBundleViewInfo{
				name:       tc,
				termKey:    term,
				termOrigin: term,
			}

			n.viewInfo.ref = n
			b.WithTaxonomy(section).Insert(n)
		} else {
			b.WithSection(section).Insert(n)
		}
	} else {
		// A regular page. Attach it to its section.
		section, _ = m.getOrCreateSection(n, bundlePath)
		b = b.WithSection(section).ForPage(bundlePath).Insert(n)
	}

	if m.cfg.isRebuild {
		// The resource owner will be either deleted or overwritten on rebuilds,
		// but make sure we handle deletion of resources (images etc.) as well.
		b.ForResource("").DeleteAll()
	}

	for _, r := range resources {
		rb := b.ForResource(cleanTreeKey(r.Meta().Path()))
		rb.Insert(&contentNode{fi: r})
	}

	return nil

}

func (m *contentMap) CreateMissingNodes() error {
	// Create missing home and root sections
	rootSections := make(map[string]interface{})
	trackRootSection := func(s string, b *contentNode) {
		parts := strings.Split(s, "/")
		if len(parts) > 2 {
			root := strings.TrimSuffix(parts[1], cmBranchSeparator)
			if root != "" {
				if _, found := rootSections[root]; !found {
					rootSections[root] = b
				}
			}
		}
	}

	m.sections.Walk(func(s string, v interface{}) bool {
		n := v.(*contentNode)

		if s == "/" {
			return false
		}

		trackRootSection(s, n)
		return false
	})

	m.pages.Walk(func(s string, v interface{}) bool {
		trackRootSection(s, v.(*contentNode))
		return false
	})

	if _, found := rootSections["/"]; !found {
		rootSections["/"] = true
	}

	for sect, v := range rootSections {
		var sectionPath string
		if n, ok := v.(*contentNode); ok && n.path != "" {
			sectionPath = n.path
			firstSlash := strings.Index(sectionPath, "/")
			if firstSlash != -1 {
				sectionPath = sectionPath[:firstSlash]
			}
		}
		sect = cleanTreeKey(sect)
		_, found := m.sections.Get(sect)
		if !found {
			m.sections.Insert(sect, &contentNode{path: sectionPath})
		}
	}

	for _, view := range m.cfg.taxonomyConfig {
		s := cleanTreeKey(view.plural)
		_, found := m.taxonomies.Get(s)
		if !found {
			b := &contentNode{
				viewInfo: &contentBundleViewInfo{
					name: view,
				},
			}
			b.viewInfo.ref = b
			m.taxonomies.Insert(s, b)
		}
	}

	return nil

}

func (m *contentMap) getBundleDir(meta hugofs.FileMeta) string {
	dir := cleanTreeKey(filepath.Dir(meta.Path()))

	switch meta.Classifier() {
	case files.ContentClassContent:
		return path.Join(dir, meta.TranslationBaseName())
	default:
		return dir
	}
}

func (m *contentMap) newContentNodeFromFi(fi hugofs.FileMetaInfo) *contentNode {
	return &contentNode{
		fi:   fi,
		path: strings.TrimPrefix(filepath.ToSlash(fi.Meta().Path()), "/"),
	}
}

func (m *contentMap) getFirstSection(s string) (string, *contentNode) {
	for {
		k, v, found := m.sections.LongestPrefix(s)
		if !found {
			return "", nil
		}
		if strings.Count(k, "/") == 1 {
			return k, v.(*contentNode)
		}
		s = path.Dir(s)
	}
}

func (m *contentMap) newKeyBuilder() *cmInsertKeyBuilder {
	return &cmInsertKeyBuilder{m: m}
}

func (m *contentMap) getOrCreateSection(n *contentNode, s string) (string, *contentNode) {
	level := strings.Count(s, "/")
	k, b := m.getSection(s)

	mustCreate := false

	if k == "" {
		mustCreate = true
	} else if level > 1 && k == "/" {
		// We found the home section, but this page needs to be placed in
		// the root, e.g. "/blog", section.
		mustCreate = true
	}

	if mustCreate {
		k = s[:strings.Index(s[1:], "/")+1]
		if k == "" {
			k = "/"
		}

		b = &contentNode{
			path: n.rootSection(),
		}

		m.sections.Insert(k, b)
	}

	return k, b
}

func (m *contentMap) getPage(section, name string) *contentNode {
	key := section + cmBranchSeparator + "/" + name + cmLeafSeparator
	v, found := m.pages.Get(key)
	if found {
		return v.(*contentNode)
	}
	return nil
}

func (m *contentMap) getSection(s string) (string, *contentNode) {
	k, v, found := m.sections.LongestPrefix(path.Dir(s))

	if found {
		return k, v.(*contentNode)
	}
	return "", nil
}

func (m *contentMap) getTaxonomyParent(s string) (string, *contentNode) {
	s = path.Dir(s)
	if s == "/" {
		v, found := m.sections.Get(s)
		if found {
			return s, v.(*contentNode)
		}
		return "", nil
	}

	for _, tree := range []*contentTree{m.taxonomies, m.sections} {
		k, v, found := tree.LongestPrefix(s)
		if found {
			return k, v.(*contentNode)
		}
	}
	return "", nil
}

func (m *contentMap) addFile(fi hugofs.FileMetaInfo) error {
	b := m.newKeyBuilder()
	return b.WithFile(fi).Insert(m.newContentNodeFromFi(fi)).err
}

func cleanTreeKey(k string) string {
	k = "/" + strings.ToLower(strings.Trim(path.Clean(filepath.ToSlash(k)), "./"))
	return k
}

func (m *contentMap) onSameLevel(s1, s2 string) bool {
	return strings.Count(s1, "/") == strings.Count(s2, "/")
}

func (m *contentMap) deleteBundleMatching(matches func(b *contentNode) bool) {
	// Check sections first
	s := m.sections.getMatch(matches)
	if s != "" {
		m.deleteSectionByPath(s)
		return
	}

	s = m.pages.getMatch(matches)
	if s != "" {
		m.deletePage(s)
		return
	}

	s = m.resources.getMatch(matches)
	if s != "" {
		m.resources.Delete(s)
	}

}

// Deletes any empty root section that's not backed by a content file.
func (m *contentMap) deleteOrphanSections() {
	var sectionsToDelete []string

	m.sections.Walk(func(s string, v interface{}) bool {
		n := v.(*contentNode)

		if n.fi != nil {
			// Section may be empty, but is backed by a content file.
			return false
		}

		if s == "/" || strings.Count(s, "/") > 1 {
			return false
		}

		prefixBundle := s + cmBranchSeparator

		if !(m.sections.hasPrefix(s+"/") || m.pages.hasPrefix(prefixBundle) || m.resources.hasPrefix(prefixBundle)) {
			sectionsToDelete = append(sectionsToDelete, s)
		}

		return false
	})

	for _, s := range sectionsToDelete {
		m.sections.Delete(s)
	}
}

func (m *contentMap) deletePage(s string) {
	m.pages.DeletePrefix(s)
	m.resources.DeletePrefix(s)
}

func (m *contentMap) deleteSectionByPath(s string) {
	m.sections.Delete(s)
	m.sections.DeletePrefix(s + "/")
	m.pages.DeletePrefix(s + cmBranchSeparator)
	m.pages.DeletePrefix(s + "/")
	m.resources.DeletePrefix(s + cmBranchSeparator)
	m.resources.DeletePrefix(s + cmLeafSeparator)
	m.resources.DeletePrefix(s + "/")
}

func (m *contentMap) deletePageByPath(s string) {
	m.pages.Walk(func(s string, v interface{}) bool {
		fmt.Println("S", s)

		return false
	})
}

func (m *contentMap) deleteTaxonomy(s string) {
	m.taxonomies.Delete(s)
	m.taxonomies.DeletePrefix(s + "/")
}

func (m *contentMap) reduceKeyPart(dir, filename string) string {
	dir, filename = filepath.ToSlash(dir), filepath.ToSlash(filename)
	dir, filename = strings.TrimPrefix(dir, "/"), strings.TrimPrefix(filename, "/")

	return strings.TrimPrefix(strings.TrimPrefix(filename, dir), "/")
}

func (m *contentMap) splitKey(k string) []string {
	if k == "" || k == "/" {
		return nil
	}

	return strings.Split(k, "/")[1:]

}

func (m *contentMap) testDump() string {
	var sb strings.Builder

	for i, r := range []*contentTree{m.pages, m.sections, m.resources} {
		sb.WriteString(fmt.Sprintf("Tree %d:\n", i))
		r.Walk(func(s string, v interface{}) bool {
			sb.WriteString("\t" + s + "\n")
			return false
		})
	}

	for i, r := range []*contentTree{m.pages, m.sections} {

		r.Walk(func(s string, v interface{}) bool {
			c := v.(*contentNode)
			cpToString := func(c *contentNode) string {
				var sb strings.Builder
				if c.p != nil {
					sb.WriteString("|p:" + c.p.Title())
				}
				if c.fi != nil {
					sb.WriteString("|f:" + filepath.ToSlash(c.fi.Meta().Path()))
				}
				return sb.String()
			}
			sb.WriteString(path.Join(m.cfg.lang, r.Name) + s + cpToString(c) + "\n")

			resourcesPrefix := s

			if i == 1 {
				resourcesPrefix += cmLeafSeparator

				m.pages.WalkPrefix(s+cmBranchSeparator, func(s string, v interface{}) bool {
					sb.WriteString("\t - P: " + filepath.ToSlash((v.(*contentNode).fi.(hugofs.FileMetaInfo)).Meta().Filename()) + "\n")
					return false
				})
			}

			m.resources.WalkPrefix(resourcesPrefix, func(s string, v interface{}) bool {
				sb.WriteString("\t - R: " + filepath.ToSlash((v.(*contentNode).fi.(hugofs.FileMetaInfo)).Meta().Filename()) + "\n")
				return false

			})

			return false
		})
	}

	return sb.String()

}

type contentMapConfig struct {
	lang                 string
	taxonomyConfig       []viewName
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
	for _, n := range cfg.taxonomyConfig {
		if strings.HasPrefix(s, n.plural) {
			return n
		}
	}

	return
}

type contentNode struct {
	p *pageState

	// Set for taxonomy nodes.
	viewInfo *contentBundleViewInfo

	// Set if source is a file.
	// We will soon get other sources.
	fi hugofs.FileMetaInfo

	// The source path. Unix slashes. No leading slash.
	path string
}

func (b *contentNode) rootSection() string {
	if b.path == "" {
		return ""
	}
	firstSlash := strings.Index(b.path, "/")
	if firstSlash == -1 {
		return b.path
	}
	return b.path[:firstSlash]

}

type contentTree struct {
	Name string
	*radix.Tree
}

type contentTrees []*contentTree

func (t contentTrees) DeletePrefix(prefix string) int {
	var count int
	for _, tree := range t {
		tree.Walk(func(s string, v interface{}) bool {
			return false
		})
		count += tree.DeletePrefix(prefix)
	}
	return count
}

type contentTreeNodeCallback func(s string, n *contentNode) bool

func newContentTreeFilter(fn func(n *contentNode) bool) contentTreeNodeCallback {
	return func(s string, n *contentNode) bool {
		return fn(n)
	}
}

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
)

func (c *contentTree) WalkQuery(query pageMapQuery, walkFn contentTreeNodeCallback) {
	filter := query.Filter
	if filter == nil {
		filter = contentTreeNoListAlwaysFilter
	}
	if query.Prefix != "" {
		c.WalkPrefix(query.Prefix, func(s string, v interface{}) bool {
			n := v.(*contentNode)
			if filter != nil && filter(s, n) {
				return false
			}
			return walkFn(s, n)
		})

		return
	}

	c.Walk(func(s string, v interface{}) bool {
		n := v.(*contentNode)
		if filter != nil && filter(s, n) {
			return false
		}
		return walkFn(s, n)
	})
}

func (c contentTrees) WalkRenderable(fn contentTreeNodeCallback) {
	query := pageMapQuery{Filter: contentTreeNoRenderFilter}
	for _, tree := range c {
		tree.WalkQuery(query, fn)
	}
}

func (c contentTrees) Walk(fn contentTreeNodeCallback) {
	for _, tree := range c {
		tree.Walk(func(s string, v interface{}) bool {
			n := v.(*contentNode)
			return fn(s, n)
		})
	}
}

func (c contentTrees) WalkPrefix(prefix string, fn contentTreeNodeCallback) {
	for _, tree := range c {
		tree.WalkPrefix(prefix, func(s string, v interface{}) bool {
			n := v.(*contentNode)
			return fn(s, n)
		})
	}
}

func (c *contentTree) getMatch(matches func(b *contentNode) bool) string {
	var match string
	c.Walk(func(s string, v interface{}) bool {
		n, ok := v.(*contentNode)
		if !ok {
			return false
		}

		if matches(n) {
			match = s
			return true
		}

		return false
	})

	return match
}

func (c *contentTree) hasPrefix(s string) bool {
	var t bool
	c.Tree.WalkPrefix(s, func(s string, v interface{}) bool {
		t = true
		return true
	})
	return t
}

func (c *contentTree) printKeys() {
	c.Walk(func(s string, v interface{}) bool {
		fmt.Println(s)
		return false
	})
}

func (c *contentTree) printKeysPrefix(prefix string) {
	c.WalkPrefix(prefix, func(s string, v interface{}) bool {
		fmt.Println(s)
		return false
	})
}

// contentTreeRef points to a node in the given tree.
type contentTreeRef struct {
	m   *pageMap
	t   *contentTree
	n   *contentNode
	key string
}

func (c *contentTreeRef) getCurrentSection() (string, *contentNode) {
	if c.isSection() {
		return c.key, c.n
	}
	return c.getSection()
}

func (c *contentTreeRef) isSection() bool {
	return c.t == c.m.sections
}

func (c *contentTreeRef) getSection() (string, *contentNode) {
	if c.t == c.m.taxonomies {
		return c.m.getTaxonomyParent(c.key)
	}
	return c.m.getSection(c.key)
}

func (c *contentTreeRef) getPages() page.Pages {
	var pas page.Pages
	c.m.collectPages(
		pageMapQuery{
			Prefix: c.key + cmBranchSeparator,
			Filter: c.n.p.m.getListFilter(true),
		},
		func(c *contentNode) {
			pas = append(pas, c.p)
		},
	)
	page.SortByDefault(pas)

	return pas
}

func (c *contentTreeRef) getPagesRecursive() page.Pages {
	var pas page.Pages

	query := pageMapQuery{
		Filter: c.n.p.m.getListFilter(true),
	}

	query.Prefix = c.key + cmBranchSeparator
	c.m.collectPages(query, func(c *contentNode) {
		pas = append(pas, c.p)
	})

	query.Prefix = c.key + "/"
	c.m.collectPages(query, func(c *contentNode) {
		pas = append(pas, c.p)
	})

	page.SortByDefault(pas)

	return pas
}

func (c *contentTreeRef) getPagesAndSections() page.Pages {
	var pas page.Pages

	query := pageMapQuery{
		Filter: c.n.p.m.getListFilter(true),
		Prefix: c.key,
	}

	c.m.collectPagesAndSections(query, func(c *contentNode) {
		pas = append(pas, c.p)
	})

	page.SortByDefault(pas)

	return pas
}

func (c *contentTreeRef) getSections() page.Pages {
	var pas page.Pages

	query := pageMapQuery{
		Filter: c.n.p.m.getListFilter(true),
		Prefix: c.key,
	}

	c.m.collectSections(query, func(c *contentNode) {
		pas = append(pas, c.p)
	})

	page.SortByDefault(pas)

	return pas
}

type contentTreeReverseIndex struct {
	t []*contentTree
	m map[interface{}]*contentNode

	init   sync.Once
	initFn func(*contentTree, map[interface{}]*contentNode)
}

func (c *contentTreeReverseIndex) Get(key interface{}) *contentNode {
	c.init.Do(func() {
		c.m = make(map[interface{}]*contentNode)
		for _, tree := range c.t {
			c.initFn(tree, c.m)
		}
	})
	return c.m[key]
}
